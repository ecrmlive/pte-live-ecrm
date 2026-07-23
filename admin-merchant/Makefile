# pte_live_vben_admin_shop — Vben 迁移版商户商城后台（nginx + /api 反代）

SHELL := /bin/bash
.DEFAULT_GOAL := help

export DOCKER_BUILDKIT := 1
IMAGE_TAG ?= v1
export IMAGE_TAG

ADMIN_SHOP_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
ROOT_DIR := $(abspath $(ADMIN_SHOP_DIR)/..)

ADMIN_SHOP_PORT ?= 11525
PNPM ?= pnpm
ROOT_VBEN_DIR := $(abspath $(ADMIN_SHOP_DIR)/../admin-platform)
COMPOSE := docker compose -p pte_live -f docker-compose.yaml

# 所有环境统一使用七禧直播品牌。
PTE_RELEASE_ENV ?= local
PTE_BRAND ?= qixi
PLATFORM_API_BASE_URL ?= $(if $(filter test,$(PTE_RELEASE_ENV)),https://api-platform.qxkejiwl.top,http://127.0.0.1:11503)
LIVE_IM_WS_URL ?= $(if $(filter test,$(PTE_RELEASE_ENV)),wss://wss.qxkejiwl.top/ws,ws://127.0.0.1:11510/ws)
VITE_BUILD_MODE ?= $(if $(filter test,$(PTE_RELEASE_ENV)),test,development)
# 商户后台在测试站挂载于 /shop/。Vite 由共享 workspace 根目录启动，
# 因此必须显式导出 base，不能只依赖 admin-merchant/.env.test。
VITE_BASE ?= $(if $(filter test,$(PTE_RELEASE_ENV)),/shop/,/)

export ADMIN_SHOP_PORT
export PTE_BRAND
export VITE_PTE_BRAND := $(PTE_BRAND)
export VITE_GLOB_API_URL := $(PLATFORM_API_BASE_URL)
export VITE_API_BASE_URL := $(PLATFORM_API_BASE_URL)
export VITE_LIVE_IM_WS_URL := $(LIVE_IM_WS_URL)
export VITE_BASE

.PHONY: help install dev pack build preview require-dist up down restart deploy logs ps docker-build ensure-network ensure-api

help: ## 显示可用命令
	@echo ""
	@echo "pte_live_vben_admin_shop — Vben 迁移版商户商城后台"
	@echo ""
	@grep -E '^[a-zA-Z0-9_-]+:.*##' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*## "}; {printf "  %-14s %s\n", $$1, $$2}'
	@echo ""
	@echo "本地开发: http://localhost:$(ADMIN_SHOP_PORT)"
	@echo "测试发布: 请从根目录执行 make test-frontend-admin-merchant"
	@echo "直播中控: /live/control/center（独立全屏页，保持现有 UI）"
	@echo "DIY 装修: views/page/** 保持现有实现，Vben 迁移时不改动"
	@echo ""

install: ## 安装 pnpm 依赖（在 monorepo 根目录）
	cd "$(ROOT_VBEN_DIR)" && $(PNPM) install

dev: ## 启动 Vben 开发服务（端口 $(ADMIN_SHOP_PORT)）
	cd "$(ROOT_VBEN_DIR)" && $(PNPM) dev:merchant

pack: ## 本地 Vite 打包 → dist/（开发机执行）
	cd "$(ROOT_VBEN_DIR)" && $(PNPM) -F @pte-live/admin-merchant exec vite build --mode $(VITE_BUILD_MODE)

build: pack ## 同 pack

require-dist: ## 检查 dist/ 是否存在
	@test -f "$(ADMIN_SHOP_DIR)/dist/index.html" || { \
		echo "缺少 $(ADMIN_SHOP_DIR)/dist/index.html" >&2; \
		echo "请先在本机 make pack，再将 dist/ 上传到服务器同路径" >&2; \
		exit 1; \
	}

preview: pack ## 本地预览 dist/ 构建产物
	cd "$(ADMIN_SHOP_DIR)" && $(PNPM) exec vite preview --host 0.0.0.0 --port $(ADMIN_SHOP_PORT)

docker-build: require-dist ## 校验本机静态产物（容器只读挂载，不构建镜像）

ensure-network: ## 确保 Docker 网络 pte_live_net 存在（172.30.0.0/24）
	@bash "$(ROOT_DIR)/scripts/ensure-qixi-live-net.sh"

ensure-api: ensure-network ## 严格检查 API（失败则 exit 1）
	@bash "$(ROOT_DIR)/scripts/stack-manage.sh" check-api shop

up: ensure-network docker-build ## 启动官方 Nginx（需 dist/ 已存在）
	@$(COMPOSE) up -d
	@SOFT_CHECK=1 bash "$(ROOT_DIR)/scripts/stack-manage.sh" check-api shop || true
	@echo ">> Vben 商户后台: http://localhost:$(ADMIN_SHOP_PORT)"

down: ## 停止并移除容器
	@$(COMPOSE) down

restart: down up ## 重新构建并启动

deploy: up ## 本地 Docker 启动

logs: ## 查看日志
	$(COMPOSE) logs -f pte_live_merchant

ps: ## 查看容器状态
	$(COMPOSE) ps -a
