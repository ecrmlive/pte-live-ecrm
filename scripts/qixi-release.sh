#!/usr/bin/env bash
# qixi-live-mergers release — 本机 pack → release/ → Docker(基建+后端) / 宿主机 Nginx(前端)
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

# shellcheck disable=SC1091
source "${SCRIPT_DIR}/release/lib.sh"
# shellcheck disable=SC1091
source "${SCRIPT_DIR}/release/config.sh"
# shellcheck disable=SC1091
source "${SCRIPT_DIR}/release/pack.sh"
# shellcheck disable=SC1091
source "${SCRIPT_DIR}/release/compose.sh"
# shellcheck disable=SC1091
source "${SCRIPT_DIR}/release/bundle.sh"

export QX_RELEASE_ENV="${QX_RELEASE_ENV:-local}"

usage() {
	cat <<EOF
用法: bash scripts/qixi-release.sh <命令> [服务]

命令:
  init-env [local|prod]   从 *.example 生成 Docker 服务运行时配置
  check-config <服务>
  pack <服务>
  local <服务>            pack；有 compose 则 up（Go 再 restart）；前端仅 pack
  release <服务>          pack + rsync；有 compose 则远程 up
  up|down|restart|ps <服务>   仅 db/mq/api/job
  local-up
  compose-check
  upload-opts             rsync release/opts（含宿主机 Nginx 配置）

服务:
  db | mq | api-admin | api-app | job | admin | merchant-admin | h5 | pc | service-web | opts
  infra-all | backend-all | frontend-all | all

说明:
  - 后台 API 与 C 端 API 分立：api-admin / api-app
  - C 端前端：h5（app-uni）· pc（app-pc），均反代 api-app
  - 不使用 Docker Nginx；见 release/opts/nginx/
EOF
}

rsync_service() {
	require_remote_env
	local key="$1" local_dir remote
	key="$(normalize_service "${key}")"
	local_dir="$(release_path "${key}")"
	remote="$(remote_path "${key}")"
	[[ -d "${local_dir}" ]] || { echo "错误: 本地目录不存在 ${local_dir}" >&2; exit 1; }
	if service_has_compose "${key}"; then
		require_prod_bundle "${key}"
	else
		verify_release_bundle_ready "${key}"
	fi
	echo ""
	echo ">> rsync release/$(service_release_dir "${key}")/ → ${RELEASE_USER}@${RELEASE_HOST}:${remote}/"
	ssh_release "mkdir -p '${remote}'"
	rsync_release_dir "${local_dir}" "${remote}"
}

upload_opts() {
	require_remote_env
	pack_service opts
	echo ""
	echo ">> rsync release/opts/ → ${RELEASE_USER}@${RELEASE_HOST}:$(remote_path opts)/"
	ssh_release "mkdir -p '$(remote_path opts)'"
	rsync_release_dir "$(release_path opts)" "$(remote_path opts)"
	echo ">> 已上传 opts；请在服务器将 nginx conf include 后 reload（勿在服务器手改业务 conf 源）"
}

run_for_services() {
	local cmd="$1" target="$2" s normalized
	normalized="$(normalize_service "${target}")" || exit 1
	expand_service_list "${normalized}"
	for s in "${EXPANDED_SERVICES[@]}"; do
		if [[ "${s}" == opts && "${cmd}" != pack && "${cmd}" != release ]]; then
			echo "跳过 opts（用 upload-opts）"
			continue
		fi
		echo ""
		echo "======== ${cmd} ${s} ========"
		case "${cmd}" in
		pack) pack_service "${s}" ;;
		local)
			pack_service "${s}"
			if service_has_compose "${s}"; then
				compose_service up "${s}"
				case "${s}" in
				api-admin | api-app | job) compose_service restart "${s}" ;;
				esac
			else
				echo ">> [${s}] 已 pack dist；请用宿主机 Nginx（release/opts/nginx/qixi-mergers.local.conf）"
			fi
			;;
		up | down | restart | ps)
			compose_service "${cmd}" "${s}"
			;;
		check-config) require_service_config "${s}" ;;
		release)
			if [[ "${s}" == opts ]]; then
				upload_opts
			elif service_has_compose "${s}"; then
				export QX_RELEASE_ENV=prod
				pack_service "${s}"
				rsync_service "${s}"
				compose_remote_restart_or_up "${s}" recreate
				export QX_RELEASE_ENV=local
			else
				# 前端：local pack + rsync dist（校验用 local env 的 dist 产物）
				pack_service "${s}"
				export QX_RELEASE_ENV=prod
				rsync_service "${s}"
				export QX_RELEASE_ENV=local
				echo ">> [${s}] dist 已上传；宿主机 Nginx root 指向该目录后 reload"
			fi
			;;
		*) echo "未知命令: ${cmd}" >&2; exit 1 ;;
		esac
	done
}

main() {
	local cmd="${1:-}" target="${2:-}"
	case "${cmd}" in
	"" | -h | --help | help) usage ;;
	init-env) init_config_from_examples "${target:-local}" ;;
	compose-check) compose_check_all ;;
	local-up) compose_local_all_up ;;
	upload-opts) upload_opts ;;
	pack | local | up | down | restart | ps | release | check-config)
		[[ -n "${target}" ]] || { usage; exit 1; }
		run_for_services "${cmd}" "${target}"
		;;
	*) usage; exit 1 ;;
	esac
}

main "$@"
