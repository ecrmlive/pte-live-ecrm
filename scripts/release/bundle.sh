#!/usr/bin/env bash
# shellcheck shell=bash
set -euo pipefail

RELEASE_BUNDLE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck disable=SC1091
source "${RELEASE_BUNDLE_DIR}/lib.sh"
# shellcheck disable=SC1091
source "${RELEASE_BUNDLE_DIR}/pack.sh"

# deploy 前校验 prod 配置齐全
require_prod_bundle() {
	local key="$1" prev
	prev="${QX_RELEASE_ENV:-local}"
	export QX_RELEASE_ENV=prod
	bootstrap_service_config "${key}"
	verify_release_bundle_ready "${key}"
	export QX_RELEASE_ENV="${prev}"
}
