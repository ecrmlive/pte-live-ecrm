#!/usr/bin/env bash
# shellcheck shell=bash
set -euo pipefail

RELEASE_LIB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${RELEASE_LIB_DIR}/../.." && pwd)"
RELEASE_DIR="${ROOT_DIR}/release"
ENV_FILE="${ROOT_DIR}/scripts/release.env"

INFRA_ALL_ORDER=(db mq)
BACKEND_ALL_ORDER=(api-admin api-app job)
FRONTEND_ALL_ORDER=(admin merchant-admin h5 pc service-web)
ALL_ORDER=(api-admin api-app job admin merchant-admin h5 pc service-web)
START_ORDER=(db mq api-admin api-app job)

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
	mq) echo qixi-mergers-mq ;;
	opts) echo opts ;;
	api-admin) echo qixi-mergers-api-admin ;;
	api-app) echo qixi-mergers-api-app ;;
	job) echo qixi-mergers-job ;;
	admin) echo qixi-mergers-admin ;;
	merchant-admin) echo qixi-mergers-merchant-admin ;;
	h5) echo qixi-mergers-h5 ;;
	pc) echo qixi-mergers-pc ;;
	service-web) echo qixi-mergers-service-web ;;
	*) echo "未知服务: $1" >&2; return 1 ;;
	esac
}

service_compose_project() {
	case "$1" in
	db) echo qixi_mergers_db ;;
	mq) echo qixi_mergers_mq ;;
	api-admin | api-app | job) echo qixi_mergers ;;
	*) echo "错误: ${1} 无 Docker compose（前端由宿主机 Nginx 托管）" >&2; return 1 ;;
	esac
}

service_container_names() {
	case "$1" in
	db) echo "qixi_mergers_mysql qixi_mergers_redis qixi_mergers_etcd qixi_mergers_minio" ;;
	mq) echo "qixi_mergers_nats" ;;
	api-admin) echo qixi_mergers_api_admin ;;
	api-app) echo qixi_mergers_api_app ;;
	job) echo qixi_mergers_job ;;
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
	db | mq | api-admin | api-app | job) return 0 ;;
	*) return 1 ;;
	esac
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
	if ! docker network inspect "${net}" >/dev/null 2>&1; then
		echo "提示: Docker 网络 ${net} 尚未创建；先 make local-db" >&2
	fi
}

ensure_docker_network_remote_snippet() {
	cat <<'EOS'
NET_NAME="${QIXI_MERGERS_DOCKER_NET:-qixi_mergers_net}"
if ! docker network inspect "${NET_NAME}" >/dev/null 2>&1; then
  echo "错误: 缺少网络 ${NET_NAME}，请先 deploy-db-reload" >&2
  exit 1
fi
EOS
}
