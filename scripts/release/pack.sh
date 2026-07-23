#!/usr/bin/env bash
# shellcheck shell=bash
set -euo pipefail

RELEASE_PACK_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck disable=SC1091
source "${RELEASE_PACK_DIR}/lib.sh"
# shellcheck disable=SC1091
source "${RELEASE_PACK_DIR}/config.sh"

pack_db() {
	local dest
	dest="$(release_path db)"
	mkdir -p "${dest}/sql"
	# 整目录进 release，deploy rsync 后仍可挂载（勿依赖 ../../sql）
	rsync -a --delete --exclude '.DS_Store' "${ROOT_DIR}/sql/" "${dest}/sql/"
	[[ -f "${dest}/sql/000_init_schema.sql" ]] || {
		echo "错误: pack_db 未得到 sql/000_init_schema.sql" >&2
		exit 1
	}
	echo ">> [db] 已同步 sql/ → release/$(service_release_dir db)/sql/"
	finalize_release_bundle_after_build db
}

pack_mq() {
	echo ">> [mq] 已迁至 pte-live-im；本仓 pack mq 等同 pack db（仅同步 sql/）"
	pack_db
}

pack_opts() {
	finalize_release_bundle_after_build opts
}

pack_go_service() {
	local key="$1" dest bin_name main_pkg
	dest="$(release_path "${key}")"
	bin_name="$(service_go_bin_name "${key}")"
	main_pkg="$(service_go_main_pkg "${key}")"
	mkdir -p "${dest}/bin"
	rm -f "${dest}/bin/${bin_name}"
	go_build_linux "${ROOT_DIR}/api" "${dest}/bin/${bin_name}" "${main_pkg}"
	sync_project_app_yaml "${key}"
	finalize_release_bundle_after_build "${key}"
}

pack_static_frontend() {
	local key="$1" repo dest dist_src=""
	case "${key}" in
	admin) repo="${ROOT_DIR}/admin-platform" ;;
	merchant-admin) repo="${ROOT_DIR}/admin-merchant" ;;
	h5) repo="${ROOT_DIR}/app-uni" ;;
	pc) repo="${ROOT_DIR}/app-pc" ;;
	service-web) repo="${ROOT_DIR}/service-web" ;;
	manager) repo="${ROOT_DIR}/app-manager" ;;
	*) echo "未知前端: ${key}" >&2; exit 1 ;;
	esac
	dest="$(release_path "${key}")"

	if [[ "${key}" == "h5" && -f "${repo}/package.json" ]]; then
		echo ">> [h5] npm run build:h5"
		(
			cd "${repo}"
			if [[ ! -d node_modules ]]; then
				npm install
			fi
			npm run build:h5
		)
		if [[ -f "${repo}/dist/build/h5/index.html" ]]; then
			dist_src="${repo}/dist/build/h5"
		elif [[ -f "${repo}/dist/index.html" ]]; then
			dist_src="${repo}/dist"
		fi
	elif [[ "${key}" == "admin" || "${key}" == "merchant-admin" ]]; then
		# Vben 5 monorepo：平台在 admin-platform，商户为 workspace 成员 admin-merchant
		echo ">> [${key}] Vben pnpm build"
		(
			cd "${ROOT_DIR}/admin-platform"
			if [[ ! -d node_modules ]]; then
				pnpm install --frozen-lockfile 2>/dev/null || pnpm install
			fi
			if [[ "${key}" == "admin" ]]; then
				pnpm run build:platform
			else
				pnpm run build:merchant
			fi
		)
		if [[ -f "${repo}/dist/index.html" ]]; then
			dist_src="${repo}/dist"
		fi
	elif [[ -f "${repo}/package.json" ]]; then
		echo ">> [${key}] pnpm build"
		(
			cd "${repo}"
			if [[ ! -d node_modules ]]; then
				pnpm install --frozen-lockfile 2>/dev/null || pnpm install
			fi
			pnpm run build
		)
		if [[ -f "${repo}/dist/index.html" ]]; then
			dist_src="${repo}/dist"
		fi
	elif [[ ! -f "${repo}/dist/index.html" ]]; then
		echo ">> [${key}] 生成骨架 dist"
		mkdir -p "${repo}/dist"
		if [[ -f "${repo}/public/index.html" ]]; then
			cp "${repo}/public/index.html" "${repo}/dist/index.html"
		elif [[ -f "${repo}/public/phase0-index.html" ]]; then
			cp "${repo}/public/phase0-index.html" "${repo}/dist/index.html"
		else
			cat >"${repo}/dist/index.html" <<HTML
<!doctype html><html lang="zh-CN"><head><meta charset="utf-8"/><title>qixi-mergers ${key}</title></head>
<body><h1>qixi-mergers ${key}</h1><p>Phase 0 skeleton</p></body></html>
HTML
		fi
		dist_src="${repo}/dist"
	else
		dist_src="${repo}/dist"
	fi

	[[ -n "${dist_src}" && -f "${dist_src}/index.html" ]] || {
		echo "错误: 缺少 ${key} 前端产物 index.html" >&2
		exit 1
	}
	rm -rf "${dest}/dist"
	mkdir -p "${dest}/dist"
	cp -R "${dist_src}/." "${dest}/dist/"
	finalize_release_bundle_after_build "${key}"
}

pack_service() {
	local key
	key="$(normalize_service "$1")"
	echo ">> [${key}] 构建发布产物 → release/$(service_release_dir "${key}")/"
	case "${key}" in
	db) pack_db ;;
	mq) pack_mq ;;
	opts) pack_opts ;;
	api-admin | api-app | job) pack_go_service "${key}" ;;
	admin | merchant-admin | h5 | pc | service-web | manager) pack_static_frontend "${key}" ;;
	*) echo "未知服务: $1" >&2; exit 1 ;;
	esac
}

verify_release_bundle_ready() {
	local key="$1" dest missing=0 bin_name
	dest="$(release_path "${key}")"
	case "${key}" in
	db | mq)
		[[ -f "${dest}/docker-compose.yaml" ]] || missing=1
		[[ -f "${dest}/config/${QX_RELEASE_ENV}/compose.env" ]] || missing=1
		[[ -f "${dest}/nats/nats.conf" ]] || missing=1
		[[ -f "${dest}/sql/000_init_schema.sql" ]] || missing=1
		;;
	opts)
		[[ -f "${dest}/nginx/qixi-mergers.local.conf" ]] || missing=1
		[[ -f "${dest}/nginx/qixi-mergers.prod.conf.example" ]] || missing=1
		;;
	api-admin | api-app | job)
		bin_name="$(service_go_bin_name "${key}")"
		[[ -f "${dest}/docker-compose.yaml" ]] || missing=1
		[[ -f "${dest}/bin/${bin_name}" ]] || missing=1
		[[ -f "${dest}/config/${QX_RELEASE_ENV}/compose.env" ]] || missing=1
		[[ -f "${dest}/config/${QX_RELEASE_ENV}/app.yaml" ]] || missing=1
		;;
	admin | merchant-admin | h5 | pc | service-web | manager)
		[[ -f "${dest}/docker-compose.yaml" ]] || missing=1
		[[ -f "${dest}/dist/index.html" ]] || missing=1
		[[ -f "${dest}/config/${QX_RELEASE_ENV}/compose.env" ]] || missing=1
		[[ -f "${dest}/nginx/default.conf" ]] || missing=1
		;;
	esac
	if (( missing != 0 )); then
		echo "错误: release/$(service_release_dir "${key}") 构建产物不完整 (env=${QX_RELEASE_ENV})" >&2
		return 1
	fi
}

print_release_bundle_manifest() {
	local key="$1" dest bin_name
	dest="$(release_path "${key}")"
	echo ">> [${key}] release/$(service_release_dir "${key}") 就绪"
	case "${key}" in
	api-admin | api-app | job)
		bin_name="$(service_go_bin_name "${key}")"
		ls -lh "${dest}/bin/${bin_name}" 2>/dev/null || true
		;;
	admin | merchant-admin | h5 | pc | service-web | manager)
		ls -lh "${dest}/dist/index.html" 2>/dev/null || true
		echo "    （compose project: qixi_mergers）"
		;;
	opts)
		ls "${dest}/nginx/"*.conf* 2>/dev/null || true
		;;
	db)
		ls "${dest}/sql/"*.sql 2>/dev/null | head -5 || true
		;;
	esac
}

finalize_release_bundle_after_build() {
	local key="$1"
	bootstrap_service_config "${key}"
	verify_release_bundle_ready "${key}"
	print_release_bundle_manifest "${key}"
}
