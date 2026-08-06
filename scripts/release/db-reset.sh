#!/usr/bin/env bash
# shellcheck shell=bash
# 仅重置当前宿主机共享 MySQL 中的七禧三库；绝不触碰 pte-live-im 的库和表。
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/../.." && pwd)"

[[ "${QX_RELEASE_ENV:-local}" == "local" ]] || {
	echo "错误: 数据库重置仅允许 QX_RELEASE_ENV=local" >&2
	exit 1
}

docker inspect --format '{{.State.Running}}' pte_live_mysql 2>/dev/null | grep -qx true || {
	echo "错误: 共享 MySQL pte_live_mysql 未运行；请先启动 pte-live-im 基础设施。" >&2
	exit 1
}

mysql_exec() {
	docker exec -i pte_live_mysql sh -ec 'MYSQL_PWD="$MYSQL_ROOT_PASSWORD" exec mysql --protocol=socket -uroot --default-character-set=utf8mb4 "$@"' sh "$@"
}

for database_name in qixi_crm_admin qixi_crm_business qixi_crm_merchant; do
	echo ">> 重建 ${database_name}（共享 MySQL 中仅七禧数据库）"
	mysql_exec -e "DROP DATABASE IF EXISTS \`${database_name}\`; CREATE DATABASE \`${database_name}\` CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;"
done

for domain in admin business merchant; do
	for phase in init_table init_config init_data init_key init_test_data; do
		sql_file="${ROOT_DIR}/sql/${domain}/${phase}.sql"
		[[ -f "${sql_file}" ]] || { echo "错误: 缺少 ${sql_file}；init_key.sql 不纳入 Git，请从 init_key.sql.example 复制后填写" >&2; exit 1; }
		echo ">> 导入 sql/${domain}/${phase}.sql"
		mysql_exec <"${sql_file}"
	done
done

# CRMEB 商户菜单全量（is_mer=1）；须在 merchant/init_data 之后覆盖写入。
merchant_menu_sql="${ROOT_DIR}/sql/merchant/init_menu_crmeb_full.sql"
[[ -f "${merchant_menu_sql}" ]] || { echo "错误: 缺少 ${merchant_menu_sql}" >&2; exit 1; }
echo ">> 导入 sql/merchant/init_menu_crmeb_full.sql"
mysql_exec <"${merchant_menu_sql}"

for pair in "qixi_crm_admin:qixi_crm_a_%" "qixi_crm_business:qixi_crm_b_%" "qixi_crm_merchant:qixi_crm_m_%"; do
	database_name="${pair%%:*}"
	prefix="${pair#*:}"
	count="$(mysql_exec -N -e "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='${database_name}' AND table_name LIKE '${prefix}';")"
	[[ "${count}" -gt 0 ]] || { echo "错误: ${database_name} 未生成 ${prefix} 表" >&2; exit 1; }
done

echo ">> 七禧 CRM 三库重置完成；pte-live-im 的数据库与表未被修改。"
