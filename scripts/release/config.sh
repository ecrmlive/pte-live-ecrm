#!/usr/bin/env bash
# shellcheck shell=bash
set -euo pipefail

RELEASE_CONFIG_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck disable=SC1091
source "${RELEASE_CONFIG_DIR}/lib.sh"

export QX_RELEASE_ENV="${QX_RELEASE_ENV:-local}"

release_config_dir() {
	local svc_dir="$1"
	echo "${svc_dir}/config/${QX_RELEASE_ENV}"
}

init_config_from_examples() {
	local env_name="${1:-local}" key svc_dir source_dir example target created=0
	# 仅 Docker 服务有 compose.env
	for key in db mq api-admin api-app job; do
		svc_dir="$(release_path "${key}")"
		source_dir="${svc_dir}/config/${env_name}"
		[[ -d "${source_dir}" ]] || continue
		while IFS= read -r -d '' example; do
			target="${example%.example}"
			if [[ ! -f "${target}" ]]; then
				cp "${example}" "${target}"
				echo "  + ${target}"
				created=1
			fi
		done < <(find "${source_dir}" -name '*.example' -print0 2>/dev/null)
	done
	if [[ "${created}" -eq 0 ]]; then
		echo ">> config/${env_name} 已存在，未覆盖"
	else
		echo ">> 已从 *.example 生成 config/${env_name}"
	fi
}

sync_project_app_yaml() {
	local key="$1" src dest_local dest_prod ex_prod
	src="$(service_go_conf_src "${key}")"
	[[ -f "${src}" ]] || { echo "错误: 缺少 ${src}" >&2; return 1; }
	dest_local="$(release_path "${key}")/config/local/app.yaml"
	dest_prod="$(release_path "${key}")/config/prod/app.yaml"
	ex_prod="$(release_path "${key}")/config/prod/app.yaml.example"
	mkdir -p "$(dirname "${dest_local}")" "$(dirname "${dest_prod}")"
	# local：每次 pack 从源码主配置覆盖
	cp "${src}" "${dest_local}"
	# prod 运行时：仅首次种子（从 prod example），避免 pack 冲掉本机已改密
	if [[ ! -f "${dest_prod}" ]]; then
		if [[ -f "${ex_prod}" ]]; then
			cp "${ex_prod}" "${dest_prod}"
		else
			cp "${src}" "${dest_prod}"
		fi
		echo ">> 已种子 ${dest_prod}（请按生产改密；后续 pack 不覆盖）"
	fi
	echo ">> 已同步 ${dest_local} ← ${src}"
}

bootstrap_service_config() {
	local key="$1" cfg_dir f
	case "${key}" in
	opts | admin | merchant-admin | h5 | pc | service-web) return 0 ;;
	esac
	cfg_dir="$(release_config_dir "$(release_path "${key}")")"
	[[ -d "${cfg_dir}" ]] || return 1
	local -a required=(compose.env)
	case "${key}" in
	api-admin | api-app | job) required+=(app.yaml) ;;
	esac
	for f in "${required[@]}"; do
		if [[ ! -f "${cfg_dir}/${f}" && -f "${cfg_dir}/${f}.example" ]]; then
			cp "${cfg_dir}/${f}.example" "${cfg_dir}/${f}"
			echo ">> 已从 example 生成 ${cfg_dir}/${f}"
		fi
	done
}

require_service_config() {
	local key="$1" cfg_dir f missing=0
	case "${key}" in
	opts | admin | merchant-admin | h5 | pc | service-web) return 0 ;;
	esac
	cfg_dir="$(release_config_dir "$(release_path "${key}")")"
	if [[ ! -d "${cfg_dir}" ]]; then
		echo "错误: 缺少 ${cfg_dir}，请 make init-env" >&2
		exit 1
	fi
	bootstrap_service_config "${key}"
	local -a required=(compose.env)
	case "${key}" in
	api-admin | api-app | job) required+=(app.yaml) ;;
	esac
	for f in "${required[@]}"; do
		if [[ ! -f "${cfg_dir}/${f}" ]]; then
			echo "错误: 缺少 ${cfg_dir}/${f}" >&2
			missing=1
		fi
	done
	(( missing == 0 )) || exit 1
}
