# 七禧统一后台 — Vben 5.7+ 本机开发与本机构建

SHELL := /bin/bash
.DEFAULT_GOAL := help

export DOCKER_BUILDKIT := 1
IMAGE_TAG ?= v1
export IMAGE_TAG

ADMIN_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
ROOT_DIR := $(abspath $(ADMIN_DIR)/..)
include $(ROOT_DIR)/scripts/make-deploy-env.mk

ADMIN_PORT ?= 15124
PNPM ?= pnpm
ROOT_VBEN_DIR := $(ADMIN_DIR)
COMPOSE := docker compose -p pte_live -f docker-compose.yaml

# 本机前端通过 Vite 反代 qixi 的 api-platform；不复用 pte-live 的前端端口或部署脚本。
VITE_BUILD_MODE ?= development

export ADMIN_PORT

.PHONY: help install dev pack build preview require-dist docker-build

help: ## 显示可用命令
	@echo ""
	@echo "七禧统一后台 — Vben 5.7+"
	@echo ""
	@grep -E '^[a-zA-Z0-9_-]+:.*##' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*## "}; {printf "  %-14s %s\n", $$1, $$2}'
	@echo ""
	@echo "本地开发: http://localhost:$(ADMIN_PORT)"
	@echo "构建产物: $(ADMIN_DIR)/dist/"
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
