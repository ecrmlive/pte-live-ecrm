#!/usr/bin/env bash
# shellcheck shell=bash
set -euo pipefail

RELEASE_LIB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${RELEASE_LIB_DIR}/../.." && pwd)"
RELEASE_DIR="${ROOT_DIR}/release"
ENV_FILE="${ROOT_DIR}/scripts/release.env"

# 基建：本仓无 Docker 基建容器。MySQL/Redis/NATS/etcd → pte-live-im；对象存储 → 腾讯云 COS。
# pack db 仅同步 sql/。INFRA_ALL_ORDER 保留 db 以便 pack 灌库脚本。
INFRA_ALL_ORDER=(db)
BACKEND_ALL_ORDER=(api-admin api-app job)
FRONTEND_ALL_ORDER=(admin merchant-admin h5 pc service-web manager)
ALL_ORDER=(api-admin api-app job admin merchant-admin h5 pc service-web manager)
START_ORDER=(api-admin api-app job admin merchant-admin h5 pc service-web manager)

load_release_env() {
	if [[ -f "${ENV_FILE}" ]]; then
		# shellcheck disable=SC1090
		source "${ENV_FILE}"
	fi
	: "${RELEASE_USER:=ubuntu}"
	: "${RELEASE_SSH_KEY:=}"
	: "${RELEASE_BASE_DIR:=/home/ubuntu/qixi-mergers}"
	: "${RELEASE_OPS_DIR:=${RELEASE_BASE_DIR}/opts}"
	: "${RELEASE_GOOS:=linux}"
	: "${RELEASE_GOARCH:=amd64}"
	: "${QIXI_MERGERS_DOCKER_NET:=qixi_mergers_net}"
	: "${PTE_LIVE_DOCKER_NET:=pte_live_net}"
}

require_remote_env() {
	load_release_env
	local missing=0
	for k in RELEASE_HOST RELEASE_USER RELEASE_BASE_DIR; do
		if [[ -z "${!k:-}" ]]; then
			echo "错误: 未设置 ${k}，请编辑 ${ENV_FILE}" >&2
			missing=1
		fi
	done
	if [[ ! -f "${ENV_FILE}" ]]; then
		echo "提示: cp scripts/release.env.example scripts/release.env" >&2
	fi
	(( missing == 0 )) || exit 1
}

normalize_service() {
	case "$1" in
	db) echo db ;;
	mq) echo mq ;;
	opts | ops) echo opts ;;
	api-admin | admin-api) echo api-admin ;;
	api-app | app-api | client-api) echo api-app ;;
	api)
		echo "错误: 已拆分为 api-admin 与 api-app，请指定其一（或 backend-all）" >&2
		return 1
		;;
	job) echo job ;;
	admin | platform-admin) echo admin ;;
	merchant-admin | mer-admin) echo merchant-admin ;;
	h5 | app-h5) echo h5 ;;
	pc | app-pc) echo pc ;;
	service-web | cs-web) echo service-web ;;
	manager | app-manager | staff) echo manager ;;
	infra-all) echo infra-all ;;
	backend-all) echo backend-all ;;
	frontend-all) echo frontend-all ;;
	all) echo all ;;
	gateway | nginx)
		echo "错误: 已取消 Docker Nginx / gateway；请用宿主机 Nginx（release/opts/nginx）" >&2
		return 1
		;;
	*) echo "未知服务: $1" >&2; return 1 ;;
	esac
}

service_release_dir() {
	case "$1" in
	db) echo qixi-mergers-db ;;
	mq) echo qixi-mergers-db ;;
	opts) echo opts ;;
	api-admin) echo qixi-mergers-api-admin ;;
	api-app) echo qixi-mergers-api-app ;;
	job) echo qixi-mergers-job ;;
	admin) echo qixi-mergers-admin ;;
	merchant-admin) echo qixi-mergers-merchant-admin ;;
	h5) echo qixi-mergers-h5 ;;
	pc) echo qixi-mergers-pc ;;
	service-web) echo qixi-mergers-service-web ;;
	manager) echo qixi-mergers-manager ;;
	*) echo "未知服务: $1" >&2; return 1 ;;
	esac
}

service_compose_project() {
	case "$1" in
	db | mq)
		echo "错误: 本仓无 db/mq Docker 服务（IM 基建 + 腾讯云 COS）" >&2
		return 1
		;;
	api-admin | api-app | job | admin | merchant-admin | h5 | pc | service-web | manager) echo qixi_mergers ;;
	*) echo "错误: ${1} 无 Docker compose" >&2; return 1 ;;
	esac
}

service_container_names() {
	case "$1" in
	db | mq)
		echo "错误: 本仓无 db/mq 容器" >&2
		return 1
		;;
	api-admin) echo qixi_mergers_api_admin ;;
	api-app) echo qixi_mergers_api_app ;;
	job) echo qixi_mergers_job ;;
	admin) echo qixi_mergers_admin ;;
	merchant-admin) echo qixi_mergers_merchant_admin ;;
	h5) echo qixi_mergers_h5 ;;
	pc) echo qixi_mergers_pc ;;
	service-web) echo qixi_mergers_service_web ;;
	manager) echo qixi_mergers_manager ;;
	*) echo "错误: ${1} 无容器" >&2; return 1 ;;
	esac
}

service_go_bin_name() {
	case "$1" in
	api-admin) echo api-admin ;;
	api-app) echo api-app ;;
	job) echo job ;;
	*) return 1 ;;
	esac
}

service_go_main_pkg() {
	case "$1" in
	api-admin) echo ./cmd/api-admin ;;
	api-app) echo ./cmd/api-app ;;
	job) echo ./cmd/job ;;
	*) return 1 ;;
	esac
}

service_go_conf_src() {
	case "$1" in
	api-admin) echo "${ROOT_DIR}/api/conf/admin.yaml" ;;
	api-app) echo "${ROOT_DIR}/api/conf/app.yaml" ;;
	job) echo "${ROOT_DIR}/api/conf/job.yaml" ;;
	*) return 1 ;;
	esac
}

service_has_compose() {
	case "$1" in
	api-admin | api-app | job | admin | merchant-admin | h5 | pc | service-web | manager) return 0 ;;
	db | mq) return 1 ;;
	*) return 1 ;;
	esac
}

print_im_infra_hint() {
	cat <<'EOF'
MySQL / Redis / NATS / etcd 由 pte-live-im 启动（db/ + mq/，网络 pte_live_net）。
对象存储使用腾讯云 COS（api/conf 中 cos:，密钥用 QIXI_COS_* 环境变量）。
请先在 IM 仓库启动基建，再对本仓执行：
  docker exec -i pte_live_mysql mysql -uroot -p"$MYSQL_ROOT_PASSWORD" < sql/000_shared_im_mysql_bootstrap.sql
  # 然后按 sql/README.md 灌入业务迁移
本仓 make local-db 仅确保 Docker 网络 qixi_mergers_net 存在（无容器）。
EOF
}

release_path() {
	echo "${RELEASE_DIR}/$(service_release_dir "$1")"
}

remote_path() {
	echo "${RELEASE_BASE_DIR}/$(service_release_dir "$1")"
}

ssh_release() {
	local -a opts=(-o StrictHostKeyChecking=accept-new)
	if [[ -n "${RELEASE_SSH_KEY:-}" ]]; then
		opts+=(-i "${RELEASE_SSH_KEY/#\~/$HOME}")
	fi
	ssh "${opts[@]}" "${RELEASE_USER}@${RELEASE_HOST}" "$@"
}

rsync_ssh_wrap() {
	if [[ -n "${RELEASE_SSH_KEY:-}" ]]; then
		printf 'ssh -o StrictHostKeyChecking=accept-new -i %s' "${RELEASE_SSH_KEY/#\~/$HOME}"
	else
		printf 'ssh -o StrictHostKeyChecking=accept-new'
	fi
}

rsync_release_dir() {
	local src="$1" dest="$2"
	rsync -az --delete \
		--exclude 'config/local/**' \
		-e "$(rsync_ssh_wrap)" \
		"${src}/" "${RELEASE_USER}@${RELEASE_HOST}:${dest}/"
}

go_build_linux() {
	local repo_dir="$1" out_path="$2" main_pkg="$3"
	load_release_env
	echo ">> go build $(basename "${out_path}") (${RELEASE_GOOS}/${RELEASE_GOARCH})"
	(
		cd "${repo_dir}"
		GOCACHE="${GOCACHE:-${TMPDIR:-/tmp}/qixi-mergers-go-build-cache}" \
			CGO_ENABLED=0 GOOS="${RELEASE_GOOS}" GOARCH="${RELEASE_GOARCH}" go build -o "${out_path}" "${main_pkg}"
	)
	chmod +x "${out_path}"
}

expand_service_list() {
	local target="$1"
	EXPANDED_SERVICES=()
	case "${target}" in
	infra-all) EXPANDED_SERVICES=("${INFRA_ALL_ORDER[@]}") ;;
	backend-all) EXPANDED_SERVICES=("${BACKEND_ALL_ORDER[@]}") ;;
	frontend-all) EXPANDED_SERVICES=("${FRONTEND_ALL_ORDER[@]}") ;;
	all) EXPANDED_SERVICES=("${ALL_ORDER[@]}") ;;
	*) EXPANDED_SERVICES=("${target}") ;;
	esac
}

ensure_docker_network() {
	local net="${QIXI_MERGERS_DOCKER_NET:-qixi_mergers_net}"
	local im_net="${PTE_LIVE_DOCKER_NET:-pte_live_net}"
	if ! docker network inspect "${net}" >/dev/null 2>&1; then
		echo ">> 创建 Docker 网络 ${net} (172.30.80.0/24)"
		docker network create \
			--driver bridge \
			--subnet 172.30.80.0/24 \
			"${net}" >/dev/null
	fi
	if ! docker network inspect "${im_net}" >/dev/null 2>&1; then
		echo "提示: 缺少 ${im_net}；请先在 pte-live-im 启动 db/mq（IM 会创建该网络）" >&2
	fi
}

ensure_docker_network_remote_snippet() {
	cat <<'EOS'
NET_NAME="${QIXI_MERGERS_DOCKER_NET:-qixi_mergers_net}"
IM_NET="${PTE_LIVE_DOCKER_NET:-pte_live_net}"
if ! docker network inspect "${NET_NAME}" >/dev/null 2>&1; then
  docker network create --driver bridge --subnet 172.30.80.0/24 "${NET_NAME}" >/dev/null
fi
if ! docker network inspect "${IM_NET}" >/dev/null 2>&1; then
  echo "错误: 缺少网络 ${IM_NET}；请先在 pte-live-im 部署 db/mq" >&2
  exit 1
fi
EOS
}
