#!/usr/bin/env bash
# shellcheck shell=bash
set -euo pipefail

RELEASE_COMPOSE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck disable=SC1091
source "${RELEASE_COMPOSE_DIR}/lib.sh"
# shellcheck disable=SC1091
source "${RELEASE_COMPOSE_DIR}/config.sh"

export QX_RELEASE_ENV="${QX_RELEASE_ENV:-local}"

compose_env_file() {
	local svc_dir="$1"
	echo "${svc_dir}/config/${QX_RELEASE_ENV}/compose.env"
}

compose_in_dir() {
	local action="$1" svc_dir="$2" project="$3" key="${4:-}" env_file
	env_file="$(compose_env_file "${svc_dir}")"
	[[ -f "${svc_dir}/docker-compose.yaml" ]] || { echo "错误: 缺少 ${svc_dir}/docker-compose.yaml" >&2; exit 1; }
	[[ -f "${env_file}" ]] || { echo "错误: 缺少 ${env_file}（make init-env）" >&2; exit 1; }
	if [[ -n "${key}" && "${action}" != ps ]]; then
		require_service_config "${key}"
	fi
	echo ">> compose ${action} (${QX_RELEASE_ENV}) $(basename "${svc_dir}")"
	(
		cd "${svc_dir}"
		export QX_RELEASE_ENV
		case "${action}" in
		up) docker compose -p "${project}" -f docker-compose.yaml --env-file "${env_file}" up -d ;;
		down) docker compose -p "${project}" -f docker-compose.yaml --env-file "${env_file}" down ;;
		restart) docker compose -p "${project}" -f docker-compose.yaml --env-file "${env_file}" restart ;;
		ps) docker compose -p "${project}" -f docker-compose.yaml --env-file "${env_file}" ps ;;
		config) docker compose -p "${project}" -f docker-compose.yaml --env-file "${env_file}" config >/dev/null ;;
		*) echo "未知 compose 动作: ${action}" >&2; exit 1 ;;
		esac
	)
}

compose_service() {
	local action="$1" key="$2" svc_dir project
	key="$(normalize_service "${key}")"
	if ! service_has_compose "${key}"; then
		echo "错误: ${key} 无 Docker（前端/Nginx 用宿主机，见 release/opts/nginx）" >&2
		exit 1
	fi
	svc_dir="$(release_path "${key}")"
	project="$(service_compose_project "${key}")"
	if [[ "${key}" != db ]]; then
		ensure_docker_network
	fi
	compose_in_dir "${action}" "${svc_dir}" "${project}" "${key}"
}

compose_remote_restart_or_up() {
	require_remote_env
	local key="$1" mode="${2:-restart}" remote project env_file containers go_bin go_preflight="" release_env
	key="$(normalize_service "${key}")"
	if ! service_has_compose "${key}"; then
		echo ">> ${key} 无远程 compose（仅 rsync dist/opts）"
		return 0
	fi
	remote="$(remote_path "${key}")"
	project="$(service_compose_project "${key}")"
	release_env="${QX_RELEASE_ENV:-prod}"
	env_file="config/${release_env}/compose.env"
	containers="$(service_container_names "${key}")"
	if go_bin="$(service_go_bin_name "${key}" 2>/dev/null)"; then
		go_preflight="
bin_path='bin/${go_bin}'
if [[ -d \"\${bin_path}\" ]]; then rm -rf \"\${bin_path}\"; fi
if [[ ! -f \"\${bin_path}\" ]]; then echo \"错误: 缺少 ${remote}/\${bin_path}\" >&2; exit 1; fi
chmod +x \"\${bin_path}\" 2>/dev/null || true
"
	fi

	echo ">> SSH 启动/重启 ${key} @ ${RELEASE_HOST} (${release_env}, ${mode})"
	ssh_release "bash -s" <<EOF
set -euo pipefail
export QX_RELEASE_ENV='${release_env}'
export QIXI_MERGERS_DOCKER_NET='${QIXI_MERGERS_DOCKER_NET:-qixi_mergers_net}'
$(ensure_docker_network_remote_snippet)
cd '${remote}'
if [[ ! -f '${env_file}' ]]; then
	echo "错误: 缺少 ${remote}/${env_file}" >&2
	exit 1
fi
${go_preflight}
running=0
exists=0
for c in ${containers}; do
	if docker ps --format '{{.Names}}' | grep -qx "\${c}"; then running=1; fi
	if docker ps -a --format '{{.Names}}' | grep -qx "\${c}"; then exists=1; fi
done
if [[ '${key}' == db || '${key}' == mq ]]; then
	docker compose -p '${project}' -f docker-compose.yaml --env-file '${env_file}' up -d
elif [[ '${mode}' == recreate ]]; then
	docker compose -p '${project}' -f docker-compose.yaml --env-file '${env_file}' up -d --force-recreate
elif [[ \${running} -eq 1 ]]; then
	docker compose -p '${project}' -f docker-compose.yaml --env-file '${env_file}' restart
elif [[ \${exists} -eq 1 ]]; then
	docker compose -p '${project}' -f docker-compose.yaml --env-file '${env_file}' up -d --force-recreate
else
	docker compose -p '${project}' -f docker-compose.yaml --env-file '${env_file}' up -d
fi
EOF
}

compose_local_all_up() {
	local key
	for key in "${START_ORDER[@]}"; do
		if [[ ! -f "$(release_path "${key}")/docker-compose.yaml" ]]; then
			echo "跳过 ${key}"
			continue
		fi
		echo ""
		echo "======== up ${key} ========"
		compose_service up "${key}"
	done
}

compose_check_all() {
	local key
	for key in db mq api-admin api-app job; do
		echo ">> compose config check ${key}"
		compose_service config "${key}"
	done
	echo ">> compose config OK（前端无 compose；Nginx 见 release/opts/nginx）"
}
