#!/usr/bin/env bash
# 七禧 CRM 单 Compose 工作流：local/test 只表示运行宿主机，运行形态完全相同。
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
RELEASE_DIR="${ROOT_DIR}/release"
RUNTIME_CONFIG="${RELEASE_DIR}/config.yaml"
COMPOSE_FILE="${ROOT_DIR}/docker-compose.yaml"

usage() {
	cat <<'EOF'
用法: bash scripts/qixi-crm.sh <命令> [服务]

命令:
  init-env [local|test]            生成同构的本机运行 YAML（不写入密钥）
  pack [api-platform|api-business|api-merchant|job|backend-all]
  check-config                     校验三套 API YAML
  compose-config                   校验唯一 docker-compose.yaml
  up [infra|db-init|backend|job|all]   校验共享基础设施、初始化数据库或启动 API
  down                             停止 pte_live_ecrm 容器
  ps                               查看 pte_live_ecrm 容器

说明:
  - local/test 使用相同的 pte_live_ecrm 项目、容器、固定 IP、数据库名和 YAML；仅宿主机不同。
  - MySQL、Redis、etcd、NATS 统一复用 pte_live_net 中的 pte_live_* 容器，七禧不启动重复基础设施。
  - JWT、数据库等运行密钥只填写 release/config.yaml 与 release/config/*/app.yaml，均不提交 Git；云服务、支付、小程序密钥只写入被 Git 忽略的 sql/*/init_key.sql。
  - local/test 必须使用内容完全相同的 sql/*/init_key.sql；db-init 会按顺序自动导入。
  - PC/H5/小程序的本地开发使用各自 Vite/HBuilderX 服务，不安装 Nginx。
EOF
}

yaml_scalar() {
	local parent="$1" key="$2" file="$3"
	awk -v parent="${parent}" -v key="${key}" '
		$0 ~ "^" parent ":[[:space:]]*$" { in_parent=1; next }
		in_parent && /^[^[:space:]]/ { exit }
		in_parent && $0 ~ "^[[:space:]]+" key ":[[:space:]]*" {
			sub("^[[:space:]]+" key ":[[:space:]]*", "")
			gsub(/^[[:space:]]*["'"'"']|["'"'"'][[:space:]]*$/, "")
			print; exit
		}
	' "${file}"
}

database_dsn() {
	local scope="$1" file="$2"
	awk -v scope="${scope}" '
		/^databases:[[:space:]]*$/ { in_databases=1; next }
		in_databases && /^[^[:space:]]/ { exit }
		in_databases && $0 ~ "^[[:space:]]{2}" scope ":[[:space:]]*$" { in_scope=1; next }
		in_scope && $0 ~ "^[[:space:]]{2}[[:alnum:]_-]+:[[:space:]]*$" { exit }
		in_scope && /^[[:space:]]{4}dsn:[[:space:]]*/ {
			sub(/^[[:space:]]{4}dsn:[[:space:]]*/, "")
			gsub(/^[[:space:]]*["'"'"']|["'"'"'][[:space:]]*$/, "")
			print; exit
		}
	' "${file}"
}

service_config_source() {
	case "$1" in
	api-platform|api-business|api-merchant|job) echo "${ROOT_DIR}/$1/conf/app.yaml" ;;
	*) echo "" ;;
	esac
}

service_release_dir() {
	case "$1" in
	api-platform) echo pte-live-ecrm-api-platform ;;
	api-business) echo pte-live-ecrm-api-business ;;
	api-merchant) echo pte-live-ecrm-api-merchant ;;
	job) echo pte-live-ecrm-job ;;
	*) return 1 ;;
	esac
}

init_env() {
	local target="${1:-local}" service src dest
	case "${target}" in local|test) ;; *) echo "错误: 环境仅支持 local 或 test" >&2; exit 1;; esac
	mkdir -p "${RELEASE_DIR}/config"
	if [[ ! -f "${RUNTIME_CONFIG}" ]]; then
		cp "${RELEASE_DIR}/config.yaml.example" "${RUNTIME_CONFIG}"
		echo "+ ${RUNTIME_CONFIG}"
	fi
	for service in api-platform api-business api-merchant job; do
		src="$(service_config_source "${service}")"
		dest="${RELEASE_DIR}/config/${service}/app.yaml"
		mkdir -p "$(dirname "${dest}")"
		if [[ ! -f "${dest}" ]]; then
			cp "${src}" "${dest}"
			echo "+ ${dest}"
		fi
	done
	echo "已生成 ${target} 宿主机使用的同构 YAML；local/test 不创建不同容器或不同数据库。"
}

check_app_yaml() {
	local service file dsn jwt_secret scope
	service="$1"
	file="${RELEASE_DIR}/config/${service}/app.yaml"
	[[ -f "${file}" ]] || { echo "错误: 缺少 ${file}" >&2; return 1; }
	case "${service}" in
	api-platform) scope=admin ;;
	api-business) scope=business ;;
	api-merchant) scope=merchant ;;
	job) scope=business ;;
	*) echo "错误: 未定义 ${service} 的数据库范围" >&2; return 1 ;;
	esac
	dsn="$(database_dsn "${scope}" "${file}")"
	jwt_secret="$(yaml_scalar jwt secret "${file}")"
	[[ -n "${dsn}" ]] || { echo "错误: ${file} 的 databases.${scope}.dsn 不能为空" >&2; return 1; }
	[[ -n "${jwt_secret}" ]] || { echo "错误: ${file} 的 jwt.secret 不能为空" >&2; return 1; }
}

check_config() {
	check_app_yaml api-platform
	check_app_yaml api-business
	check_app_yaml api-merchant
	check_app_yaml job
	echo "YAML 配置完整：local/test 共用相同配置约定。"
}

pack_one() {
	local service="$1" output source_dir
	output="${RELEASE_DIR}/$(service_release_dir "${service}")/bin/${service}"
	source_dir="${ROOT_DIR}/${service}"
	[[ -f "${source_dir}/go.mod" ]] || {
		echo "错误: 独立服务项目不存在：${source_dir}" >&2
		return 1
	}
	mkdir -p "$(dirname "${output}")"
	(
		cd "${source_dir}"
		GOCACHE="${GOCACHE:-${TMPDIR:-/tmp}/pte-live-ecrm-go-cache}" CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
			go build -o "${output}" ./cmd
	)
	chmod +x "${output}"
	echo "已构建 ${service} → ${output}"
}

pack() {
	local target="${1:-backend-all}"
	case "${target}" in
	api-platform|api-business|api-merchant|job) pack_one "${target}" ;;
	backend-all)
		pack_one api-platform
		pack_one api-business
		pack_one api-merchant
		;;
	*) echo "错误: 未知构建目标 ${target}" >&2; exit 1;;
	esac
}

compose() {
	docker compose --project-name pte_live_ecrm --file "${COMPOSE_FILE}" "$@"
}

require_job_infra() {
	docker inspect --format '{{.State.Running}}' pte_live_mysql 2>/dev/null | grep -qx true || {
		echo "错误: 共享基础设施 pte_live_mysql 未运行；job 的业务库任务无法启动。" >&2
		exit 1
	}
}

require_shared_infra() {
	local container_name
	for container_name in pte_live_mysql pte_live_redis pte_live_etcd1 pte_live_etcd2 pte_live_etcd3 pte_live_nats1 pte_live_nats2 pte_live_nats3; do
		docker inspect --format '{{.State.Running}}' "${container_name}" 2>/dev/null | grep -qx true || {
			echo "错误: 共享基础设施 ${container_name} 未运行；请先在 pte-live-im 项目启动 pte_live_db 与 pte_live_mq。" >&2
			exit 1
		}
	done
	docker network inspect pte_live_net >/dev/null 2>&1 || {
		echo "错误: 共享网络 pte_live_net 不存在；请先在 pte-live-im 项目启动基础设施。" >&2
		exit 1
	}
}

shared_mysql() {
	docker exec -i pte_live_mysql sh -ec 'MYSQL_PWD="$MYSQL_ROOT_PASSWORD" exec mysql --protocol=socket -uroot --default-character-set=utf8mb4 "$@"' sh "$@"
}

shared_database_password() {
	yaml_scalar shared_infrastructure database_password "${RUNTIME_CONFIG}"
}

ensure_shared_database_config() {
	local database_password config_tmp
	[[ -f "${RUNTIME_CONFIG}" ]] || {
		echo "错误: 缺少 ${RUNTIME_CONFIG}；先执行 make init-env-local" >&2
		exit 1
	}
	database_password="$(shared_database_password)"
	if [[ -z "${database_password}" ]]; then
		database_password="$(openssl rand -hex 24)"
		config_tmp="$(mktemp)"
		awk -v value="${database_password}" '
			/^shared_infrastructure:[[:space:]]*$/ { in_shared=1 }
			in_shared && /^[^[:space:]]/ && $0 !~ /^shared_infrastructure:/ { in_shared=0 }
			in_shared && /^[[:space:]]+mysql_container:[[:space:]]*/ {
				print
				print "  database_user: qixi_crm"
				print "  database_password: \"" value "\""
				next
			}
			{ print }
		' "${RUNTIME_CONFIG}" >"${config_tmp}"
		mv "${config_tmp}" "${RUNTIME_CONFIG}"
		echo "已在被 Git 忽略的 release/config.yaml 生成七禧共享数据库账号；test 宿主机必须复制同一份 YAML。"
	fi
}

sync_service_dsn() {
	local service="$1" scope="$2" database_name="$3" file database_password dsn config_tmp
	file="${RELEASE_DIR}/config/${service}/app.yaml"
	[[ -f "${file}" ]] || { echo "错误: 缺少 ${file}" >&2; exit 1; }
	database_password="$(shared_database_password)"
	dsn="qixi_crm:${database_password}@tcp(pte_live_mysql:3306)/${database_name}?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai&time_zone=%27%2B08%3A00%27"
	config_tmp="$(mktemp)"
	awk -v target_scope="${scope}" -v target_dsn="${dsn}" '
		/^databases:[[:space:]]*$/ { in_databases=1 }
		in_databases && /^[^[:space:]]/ && $0 !~ /^databases:/ { in_databases=0; in_scope=0 }
		in_databases && $0 ~ "^[[:space:]]{2}" target_scope ":[[:space:]]*$" { in_scope=1 }
		in_scope && /^[[:space:]]{2}[[:alnum:]_-]+:[[:space:]]*$/ && $0 !~ "^[[:space:]]{2}" target_scope ":[[:space:]]*$" { in_scope=0 }
		in_scope && /^[[:space:]]{4}dsn:[[:space:]]*/ { print "    dsn: \"" target_dsn "\""; next }
		{ print }
	' "${file}" >"${config_tmp}"
	mv "${config_tmp}" "${file}"
}

provision_shared_database_user() {
	local database_password
	ensure_shared_database_config
	database_password="$(shared_database_password)"
	[[ "${database_password}" =~ ^[A-Za-z0-9]+$ ]] || {
		echo "错误: shared_infrastructure.database_password 仅允许字母和数字。" >&2
		exit 1
	}
	shared_mysql -e "SET GLOBAL time_zone = '+08:00'; SET time_zone = '+08:00';"
	shared_mysql -e "CREATE USER IF NOT EXISTS 'qixi_crm'@'%' IDENTIFIED BY '${database_password}'; ALTER USER 'qixi_crm'@'%' IDENTIFIED BY '${database_password}'; GRANT ALL PRIVILEGES ON qixi_crm_admin.* TO 'qixi_crm'@'%'; GRANT ALL PRIVILEGES ON qixi_crm_business.* TO 'qixi_crm'@'%'; GRANT ALL PRIVILEGES ON qixi_crm_merchant.* TO 'qixi_crm'@'%'; FLUSH PRIVILEGES;"
	sync_service_dsn api-platform admin qixi_crm_admin
	sync_service_dsn api-platform business qixi_crm_business
	sync_service_dsn api-business business qixi_crm_business
	sync_service_dsn job business qixi_crm_business
	sync_service_dsn api-merchant merchant qixi_crm_merchant
	sync_service_dsn api-merchant business qixi_crm_business
}

initialize_databases() {
	local domain phase sql_file
	require_shared_infra
	for domain in admin business merchant; do
		for phase in init_table init_config init_data init_key init_test_data; do
			sql_file="${ROOT_DIR}/sql/${domain}/${phase}.sql"
			[[ -f "${sql_file}" ]] || { echo "错误: 缺少 ${sql_file}；init_key.sql 不纳入 Git，请从 init_key.sql.example 复制后填写" >&2; exit 1; }
			echo ">> 导入 sql/${domain}/${phase}.sql 到 pte_live_mysql"
			shared_mysql <"${sql_file}"
		done
	done
	screen_demo_sql="${ROOT_DIR}/sql/business/init_data_screen_demo.sql"
	if [[ -f "${screen_demo_sql}" ]]; then
		echo ">> 导入 sql/business/init_data_screen_demo.sql 到 pte_live_mysql"
		shared_mysql <"${screen_demo_sql}"
	fi
	product_stat_demo_sql="${ROOT_DIR}/sql/business/init_product_stat_demo.sql"
	if [[ -f "${product_stat_demo_sql}" ]]; then
		echo ">> 导入 sql/business/init_product_stat_demo.sql 到 pte_live_mysql"
		shared_mysql <"${product_stat_demo_sql}"
	fi
	merchant_menu_sql="${ROOT_DIR}/sql/merchant/init_menu_crmeb_full.sql"
	[[ -f "${merchant_menu_sql}" ]] || { echo "错误: 缺少 ${merchant_menu_sql}" >&2; exit 1; }
	echo ">> 导入 sql/merchant/init_menu_crmeb_full.sql 到 pte_live_mysql"
	shared_mysql <"${merchant_menu_sql}"
	provision_shared_database_user
	echo "七禧三库已导入共享 MySQL：qixi_crm_admin、qixi_crm_business、qixi_crm_merchant。"
}

main() {
	local command="${1:-}" target="${2:-}"
	case "${command}" in
	""|help|-h|--help) usage ;;
	init-env) init_env "${target:-local}" ;;
	check-config) check_config ;;
	pack) pack "${target:-backend-all}" ;;
	compose-config) compose config >/dev/null; echo "docker-compose.yaml 结构校验通过" ;;
	up)
		case "${target:-all}" in
		infra) require_shared_infra; echo "已复用 pte_live_net 共享基础设施；七禧不会创建数据库或中间件容器。" ;;
		db-init) initialize_databases ;;
		backend) require_shared_infra; compose up -d api-platform api-business api-merchant ;;
		job) require_job_infra; compose --profile job up -d job ;;
		all) require_shared_infra; compose up -d api-platform api-business api-merchant ;;
		*) echo "错误: up 仅支持 infra、db-init、backend、job、all" >&2; exit 1;;
		esac
		;;
	down) compose down ;;
	ps) compose ps ;;
	*) usage; exit 1;;
	esac
}

main "$@"
