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
  check-config                     校验运行 YAML 和三套 API YAML
  compose-config                   校验唯一 docker-compose.yaml
  up [infra|db-init|backend|all]   启动同构容器
  down                             停止 qixi_mergers 容器
  ps                               查看 qixi_mergers 容器

说明:
  - local/test 使用相同的 qixi_mergers 项目、容器、网络、IP、数据库名和 YAML；仅宿主机不同。
  - 所有密钥只填写 release/config.yaml 与 release/config/*/app.yaml，均不提交 Git。
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

load_runtime_config() {
	[[ -f "${RUNTIME_CONFIG}" ]] || {
		echo "错误: 缺少 ${RUNTIME_CONFIG}；先执行 make init-env-local" >&2
		exit 1
	}
	QIXI_MYSQL_ROOT_PASSWORD="$(yaml_scalar docker mysql_root_password "${RUNTIME_CONFIG}")"
	[[ -n "${QIXI_MYSQL_ROOT_PASSWORD}" ]] || {
		echo "错误: ${RUNTIME_CONFIG} 的 docker.mysql_root_password 不能为空" >&2
		exit 1
	}
	export QIXI_MYSQL_ROOT_PASSWORD
}

service_config_source() {
	case "$1" in
	api-platform|api-business|api-merchant|job) echo "${ROOT_DIR}/$1/conf/app.yaml" ;;
	*) echo "" ;;
	esac
}

service_release_dir() {
	case "$1" in
	api-platform) echo qixi-mergers-api-platform ;;
	api-business) echo qixi-mergers-api-business ;;
	api-merchant) echo qixi-mergers-api-merchant ;;
	job) echo qixi-mergers-job ;;
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
	*) echo "错误: 未定义 ${service} 的数据库范围" >&2; return 1 ;;
	esac
	dsn="$(database_dsn "${scope}" "${file}")"
	jwt_secret="$(yaml_scalar jwt secret "${file}")"
	[[ -n "${dsn}" ]] || { echo "错误: ${file} 的 databases.${scope}.dsn 不能为空" >&2; return 1; }
	[[ -n "${jwt_secret}" ]] || { echo "错误: ${file} 的 jwt.secret 不能为空" >&2; return 1; }
}

check_config() {
	load_runtime_config
	check_app_yaml api-platform
	check_app_yaml api-business
	check_app_yaml api-merchant
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
		GOCACHE="${GOCACHE:-${TMPDIR:-/tmp}/qixi-mergers-go-cache}" CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
			go build -o "${output}" ./cmd
	)
	chmod +x "${output}"
	echo "已构建 ${service} → ${output}"
}

pack() {
	local target="${1:-backend-all}"
	case "${target}" in
	api-platform|api-business|api-merchant) pack_one "${target}" ;;
	backend-all)
		pack_one api-platform
		pack_one api-business
		pack_one api-merchant
		;;
	job)
		echo "错误: job 仍依赖旧单库数据访问层，完成迁移前禁止构建进入新运行形态。" >&2
		exit 1
		;;
	*) echo "错误: 未知构建目标 ${target}" >&2; exit 1;;
	esac
}

compose() {
	load_runtime_config
	docker compose --project-name qixi_mergers --file "${COMPOSE_FILE}" "$@"
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
		infra) compose up -d mysql redis etcd nats ;;
		db-init) compose up db-init ;;
		backend) compose up -d mysql redis etcd nats api-platform api-business api-merchant ;;
		all) compose up -d mysql redis etcd nats api-platform api-business api-merchant ;;
		*) echo "错误: up 仅支持 infra、db-init、backend、all" >&2; exit 1;;
		esac
		;;
	down) compose down ;;
	ps) compose ps ;;
	*) usage; exit 1;;
	esac
}

main "$@"
