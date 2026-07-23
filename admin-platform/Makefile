# pte_live_vben_admin — Vben 迁移版平台超管（nginx + /api 反代）

SHELL := /bin/bash
.DEFAULT_GOAL := help

export DOCKER_BUILDKIT := 1
IMAGE_TAG ?= v1
export IMAGE_TAG

ADMIN_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
ROOT_DIR := $(abspath $(ADMIN_DIR)/..)
include $(ROOT_DIR)/scripts/make-deploy-env.mk

ADMIN_PORT ?= 11524
PNPM ?= pnpm
ROOT_VBEN_DIR := $(ADMIN_DIR)
COMPOSE := docker compose -p pte_live -f docker-compose.yaml

# 所有环境统一使用七禧直播品牌。
PTE_RELEASE_ENV ?= local
PTE_BRAND ?= qixi
PLATFORM_API_BASE_URL ?= $(if $(filter test,$(PTE_RELEASE_ENV)),https://api-platform.qxkejiwl.top,http://127.0.0.1:11503)
VITE_BUILD_MODE ?= $(if $(filter test,$(PTE_RELEASE_ENV)),test,development)

export ADMIN_PORT
export PTE_BRAND
export VITE_PTE_BRAND := $(PTE_BRAND)
export VITE_GLOB_API_URL := $(PLATFORM_API_BASE_URL)
export VITE_API_BASE_URL := $(PLATFORM_API_BASE_URL)

.PHONY: help install dev pack build preview require-dist up down restart deploy logs ps docker-build ensure-network ensure-api

help: ## 显示可用命令
	@echo ""
	@echo "pte_live_vben_admin — Vben 迁移版平台超管"
	@echo ""
	@grep -E '^[a-zA-Z0-9_-]+:.*##' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*## "}; {printf "  %-14s %s\n", $$1, $$2}'
	@echo ""
	@echo "本地开发: http://localhost:$(ADMIN_PORT)"
	@echo "测试发布: 请从根目录执行 make test-frontend-admin-platform"
	@echo "Token Key: qixiLiveAdminToken（与 qixi-live-admin 的 jjjSingleAdminToken 隔离）"
	@echo ""

install: ## 安装 pnpm 依赖（在 monorepo 根目录）
	cd "$(ROOT_VBEN_DIR)" && $(PNPM) install

dev: ## 启动 Vben 开发服务（端口 $(ADMIN_PORT)）
	cd "$(ROOT_VBEN_DIR)" && $(PNPM) dev:platform

pack: ## 本地 Vite 打包 → dist/（开发机执行）
	cd "$(ROOT_VBEN_DIR)" && $(PNPM) vite build --mode $(VITE_BUILD_MODE)

build: pack ## 同 pack

require-dist: ## 检查 dist/ 是否存在
	@test -f "$(ADMIN_DIR)/dist/index.html" || { \
		echo "缺少 $(ADMIN_DIR)/dist/index.html" >&2; \
		echo "请先在本机 make pack，再将 dist/ 上传到服务器同路径" >&2; \
		exit 1; \
	}

preview: pack ## 本地预览 dist/ 构建产物
	cd "$(ADMIN_DIR)" && $(PNPM) exec vite preview --host 0.0.0.0 --port $(ADMIN_PORT)

docker-build: require-dist ## 校验本机静态产物（容器只读挂载，不构建镜像）

ensure-network: ## 确保 Docker 网络 pte_live_net 存在（172.30.0.0/24）
	@bash "$(ROOT_DIR)/scripts/ensure-qixi-live-net.sh"

ensure-api: ensure-network ## 严格检查 API（失败则 exit 1）
	@bash "$(ROOT_DIR)/scripts/stack-manage.sh" check-api admin

up: ensure-network docker-build ## 启动官方 Nginx（需 dist/ 已存在）
	@$(COMPOSE) up -d
	@SOFT_CHECK=1 bash "$(ROOT_DIR)/scripts/stack-manage.sh" check-api admin || true
	@echo ">> Vben 平台超管: http://localhost:$(ADMIN_PORT)"

down: ## 停止并移除容器
	@$(COMPOSE) down

restart: down up ## 重新构建并启动

deploy: up ## 本地 Docker 启动

logs: ## 查看日志
	$(COMPOSE) logs -f pte_live_platform

ps: ## 查看容器状态
	$(COMPOSE) ps -a
