#!/usr/bin/env bash
# 本机开发同步：把配置初始化、sql patch/seed 导入 local MySQL，并把改过的 API pack 后重启容器。
# 用法：
#   scripts/local-dev-sync.sh              # 自动：git 变更涉及的 sql + api
#   scripts/local-dev-sync.sh sql          # 仅导入配置初始化、patch/seed
#   scripts/local-dev-sync.sh api          # 仅按 git 变更 pack+重启
#   scripts/local-dev-sync.sh api api-platform api-business
#   scripts/local-dev-sync.sh all          # 全量 patch/seed + 三个 API
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
QX="${ROOT_DIR}/scripts/qixi-crm.sh"

die() { echo "错误: $*" >&2; exit 1; }

require_mysql() {
	docker inspect --format '{{.State.Running}}' pte_live_mysql 2>/dev/null | grep -qx true \
		|| die "pte_live_mysql 未运行"
}

mysql_db() {
	local db="$1"
	shift
	docker exec -i pte_live_mysql sh -ec \
		'MYSQL_PWD="$MYSQL_ROOT_PASSWORD" exec mysql --protocol=socket -uroot --default-character-set=utf8mb4 "$@"' \
		sh "$@" "$db"
}

apply_sql_file() {
	local db="$1" file="$2"
	echo ">> SQL [${db}] $(basename "${file}")"
	if ! mysql_db "${db}" <"${file}"; then
		echo "!! SQL 失败（已跳过继续）：${file}" >&2
		return 0
	fi
}

# 已废弃或不应在 auto-sync 全量重跑的文件（仍可用 make local-db-init / 手工执行）
sql_skip() {
	case "$(basename "$1")" in
	patch_store_menu_local.sql) return 0 ;; # DEPRECATED，全量菜单见 init_menu_crmeb_full.sql
	esac
	return 1
}

apply_domain_sql() {
	local domain="$1" db="$2"
	local dir="${ROOT_DIR}/sql/${domain}"
	[[ -d "${dir}" ]] || return 0
	local f
	# patch 优先（结构/菜单），再 seed（模拟数据）
	while IFS= read -r f; do
		[[ -n "${f}" ]] || continue
		if sql_skip "${f}"; then
			echo "-- 跳过 $(basename "${f}")"
			continue
		fi
		apply_sql_file "${db}" "${f}"
	done < <(find "${dir}" -maxdepth 1 -type f \( -name 'patch_*.sql' -o -name 'seed_*_local.sql' -o -name 'seed_*.sql' \) | LC_ALL=C sort)
}

apply_domain_init_defaults() {
	local domain="$1" db="$2"
	local dir="${ROOT_DIR}/sql/${domain}"
	[[ -d "${dir}" ]] || return 0
	local phase file
	# init_key.sql 只在完整 db-init 中导入，避免本地同步读取或覆盖本机密钥。
	for phase in init_config init_data init_file; do
		file="${dir}/${phase}.sql"
		[[ -f "${file}" ]] || continue
		apply_sql_file "${db}" "${file}"
	done
}

apply_sql() {
	require_mysql
	echo "==> 导入 local patch / 配置初始化 / seed"
	# 结构补丁优先：新默认配置可能依赖本次新增的表。
	apply_domain_sql admin qixi_crm_admin
	apply_domain_sql business qixi_crm_business
	apply_domain_sql merchant qixi_crm_merchant
	apply_domain_init_defaults admin qixi_crm_admin
	apply_domain_init_defaults business qixi_crm_business
	apply_domain_init_defaults merchant qixi_crm_merchant
	echo "==> SQL 同步完成"
}

container_for_svc() {
	case "$1" in
	api-platform) echo pte_live_ecrm_api_platform ;;
	api-business) echo pte_live_ecrm_api_business ;;
	api-merchant) echo pte_live_ecrm_api_merchant ;;
	*) die "未知 API 服务: $1" ;;
	esac
}

pack_and_restart() {
	local svc
	[[ "$#" -gt 0 ]] || die "未指定要部署的 API"
	for svc in "$@"; do
		case "${svc}" in
		api-platform|api-business|api-merchant) ;;
		*) die "仅支持 api-platform|api-business|api-merchant，收到: ${svc}" ;;
		esac
	done
	for svc in "$@"; do
		echo "==> pack ${svc}"
		/bin/bash "${QX}" pack "${svc}"
	done
	local containers=()
	for svc in "$@"; do
		containers+=("$(container_for_svc "${svc}")")
	done
	echo "==> restart ${containers[*]}"
	docker restart "${containers[@]}" >/dev/null
	sleep 2
	for c in "${containers[@]}"; do
		docker inspect --format '{{.State.Status}}' "${c}" | grep -qx running \
			|| die "${c} 未处于 running"
		echo "   ${c}: running"
	done
	echo "==> API local 部署完成"
}

detect_changed_apis() {
	local paths
	paths="$(
		cd "${ROOT_DIR}"
		{
			git diff --name-only HEAD 2>/dev/null || true
			git diff --name-only --cached 2>/dev/null || true
			git ls-files --others --exclude-standard 2>/dev/null || true
		} | sort -u
	)"
	local out=()
	echo "${paths}" | grep -qE '^api-platform/' && out+=(api-platform)
	echo "${paths}" | grep -qE '^api-business/' && out+=(api-business)
	echo "${paths}" | grep -qE '^api-merchant/' && out+=(api-merchant)
	printf '%s\n' "${out[@]+"${out[@]}"}"
}

detect_sql_changed() {
	cd "${ROOT_DIR}"
	{
		git diff --name-only HEAD 2>/dev/null || true
		git diff --name-only --cached 2>/dev/null || true
		git ls-files --others --exclude-standard 2>/dev/null || true
	} | grep -qE '^sql/(admin|business|merchant)/' && return 0
	return 1
}

auto_sync() {
	local apis=()
	local line
	while IFS= read -r line; do
		[[ -n "${line}" ]] && apis+=("${line}")
	done < <(detect_changed_apis)

	if detect_sql_changed || [[ "${FORCE_SQL:-0}" == "1" ]]; then
		apply_sql
	else
		echo "==> 未检测到 sql/ 变更，跳过 SQL（FORCE_SQL=1 可强制）"
	fi

	if [[ "${#apis[@]}" -gt 0 ]]; then
		pack_and_restart "${apis[@]}"
	else
		echo "==> 未检测到 api-* 变更，跳过 API 部署（可显式传服务名）"
	fi
}

cmd="${1:-auto}"
shift || true

case "${cmd}" in
sql|db)
	apply_sql
	;;
api)
	if [[ "$#" -gt 0 ]]; then
		pack_and_restart "$@"
	else
		mapfile -t apis < <(detect_changed_apis)
		[[ "${#apis[@]}" -gt 0 ]] || die "未检测到 api 变更；请显式指定：api api-platform ..."
		pack_and_restart "${apis[@]}"
	fi
	;;
all)
	apply_sql
	pack_and_restart api-platform api-business api-merchant
	;;
auto)
	auto_sync
	;;
*)
	die "用法: $0 [auto|sql|api [svc...]|all]"
	;;
esac
