#!/usr/bin/env bash
# shellcheck shell=bash
# 仅重置当前宿主机的七禧 CRM 三库；local/test 不得在同一宿主机并行运行。
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/../.." && pwd)"
CONFIG_FILE="${ROOT_DIR}/release/config.yaml"
COMPOSE_FILE="${ROOT_DIR}/docker-compose.yaml"

[[ "${QX_RELEASE_ENV:-local}" == "local" ]] || {
	echo "错误: 数据库重置仅允许 QX_RELEASE_ENV=local" >&2
	exit 1
}
[[ -f "${CONFIG_FILE}" ]] || { echo "错误: 缺少 ${CONFIG_FILE}" >&2; exit 1; }

MYSQL_ROOT_PASSWORD="$(awk '
	/^docker:[[:space:]]*$/ { in_docker=1; next }
	in_docker && /^[^[:space:]]/ { exit }
	in_docker && /^[[:space:]]+mysql_root_password:[[:space:]]*/ {
		sub(/^[[:space:]]+mysql_root_password:[[:space:]]*/, "")
		gsub(/^[[:space:]]*["'"'"']|["'"'"'][[:space:]]*$/, "")
		print; exit
	}
' "${CONFIG_FILE}")"
[[ -n "${MYSQL_ROOT_PASSWORD}" ]] || { echo "错误: docker.mysql_root_password 不能为空" >&2; exit 1; }

compose() {
	QIXI_MYSQL_ROOT_PASSWORD="${MYSQL_ROOT_PASSWORD}" docker compose --project-name qixi_mergers --file "${COMPOSE_FILE}" "$@"
}

compose ps --status running --services | grep -qx mysql || {
	echo "错误: qixi_mergers_mysql 未运行；先执行 make local-infra" >&2
	exit 1
}

mysql_exec() {
	compose exec -T -e "MYSQL_PWD=${MYSQL_ROOT_PASSWORD}" mysql \
		mysql --protocol=TCP -h127.0.0.1 -uroot --default-character-set=utf8mb4 "$@"
}

for database in qixi_crm_admin qixi_crm_business qixi_crm_merchant; do
	echo ">> 重建 ${database}"
	mysql_exec -e "DROP DATABASE IF EXISTS \`${database}\`; CREATE DATABASE \`${database}\` CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;"
done

for domain in admin business merchant; do
	for phase in 01_table 02_data 03_config 04_key 05_test_data; do
		path="${ROOT_DIR}/sql/${domain}/${phase}.sql"
		[[ -f "${path}" ]] || { echo "错误: 缺少 ${path}" >&2; exit 1; }
		echo ">> 导入 sql/${domain}/${phase}.sql"
		mysql_exec <"${path}"
	done
done

for pair in "qixi_crm_admin:qixi_crm_a_%" "qixi_crm_business:qixi_crm_b_%" "qixi_crm_merchant:qixi_crm_m_%"; do
	database="${pair%%:*}"
	prefix="${pair#*:}"
	count="$(mysql_exec -N -e "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='${database}' AND table_name LIKE '${prefix}';")"
	[[ "${count}" -gt 0 ]] || { echo "错误: ${database} 未生成 ${prefix} 表" >&2; exit 1; }
done

echo ">> 七禧 CRM 三库重置完成（admin/business/merchant）"
