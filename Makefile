# 七禧 CRM：local/test 仅为执行宿主机，不产生不同的容器、网络、数据库或服务名。
SHELL := /bin/bash
.DEFAULT_GOAL := help

ROOT_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
QX := /bin/bash $(ROOT_DIR)/scripts/qixi-crm.sh

.PHONY: help init-env-local init-env-test check-config pack pack-backend local-admin-init \
	local local-infra local-db-init local-db-reset local-backend local-job local-down local-ps local-compose-check \
	local-sync local-sync-sql local-sync-api \
	test test-infra test-db-init test-backend test-job test-down test-ps test-compose-check

help:
	@$(QX) help

init-env-local:
	@$(QX) init-env local

init-env-test:
	@$(QX) init-env test

check-config:
	@$(QX) check-config

pack:
	@$(QX) pack $(SVC)

pack-backend:
	@$(QX) pack backend-all

# 仅从被 Git 忽略的 YAML 读取本机管理员初始化数据；不在 SQL 或仓库中保存密码。
local-admin-init:
	@cd api-platform && GOCACHE="$${GOCACHE:-$${TMPDIR:-/tmp}/pte-live-ecrm-go-cache}" \
		go run ./cmd/admin-init -app-config ../release/config/api-platform/app.yaml -seed-config ../release/admin-users.yaml

local: local-backend
local-infra:
	@$(QX) up infra
local-db-init:
	@$(QX) up db-init
local-db-reset:
	@QX_RELEASE_ENV=local /bin/bash $(ROOT_DIR)/scripts/release/db-reset.sh
local-backend:
	@$(QX) up backend
local-job:
	@$(QX) up job
local-down:
	@$(QX) down
local-ps:
	@$(QX) ps
local-compose-check:
	@$(QX) compose-config

# 开发同步：sql patch/seed → local MySQL；改过的 API → pack + 重启容器
#   make local-sync              # 按 git 变更自动
#   make local-sync-sql          # 仅 SQL
#   make local-sync-api SVC=api-platform
#   make local-sync MODE=all     # 全量 SQL + 三个 API
local-sync:
	@/bin/bash $(ROOT_DIR)/scripts/local-dev-sync.sh $${MODE:-auto}
local-sync-sql:
	@/bin/bash $(ROOT_DIR)/scripts/local-dev-sync.sh sql
local-sync-api:
	@/bin/bash $(ROOT_DIR)/scripts/local-dev-sync.sh api $(SVC)

# test 与 local 调用完全相同的单 Compose；差别只在命令执行的宿主机。
test: test-backend
test-infra:
	@$(QX) up infra
test-db-init:
	@$(QX) up db-init
test-backend:
	@$(QX) up backend
test-job:
	@$(QX) up job
test-down:
	@$(QX) down
test-ps:
	@$(QX) ps
test-compose-check:
	@$(QX) compose-config
