# qixi-live-mergers — 根目录 Make
# 约定见 docs/release/COMMANDS.md · SERVICE-MATRIX.md
# API 分立：api-admin（后台）/ api-app（C 端）；Nginx 仅宿主机

SHELL := /bin/bash
.DEFAULT_GOAL := help

ROOT_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
QX       := bash $(ROOT_DIR)/scripts/qixi-release.sh

ENV ?= local
SVC ?= api-admin
export QX_RELEASE_ENV := $(ENV)

.PHONY: help init-env init-env-prod check-config \
	pack pack-infra pack-backend pack-frontend pack-all \
	local local-db local-mq local-api-admin local-api-app local-job \
	local-admin local-merchant-admin local-h5 local-pc local-service-web \
	local-backend-all local-frontend-all local-up local-compose-check \
	up down restart ps \
	deploy-db-reload deploy-mq-reload deploy-api-admin deploy-api-app deploy-job \
	deploy-admin deploy-merchant-admin deploy-h5 deploy-pc deploy-service-web \
	deploy-backend-all deploy-frontend-all deploy-all update-nginx

help:
	@echo ""
	@echo "qixi-live-mergers — release 命令"
	@echo ""
	@echo "  make init-env / init-env-prod"
	@echo "  make local-db / local-mq"
	@echo "  make local-api-admin    # 后台 :18080"
	@echo "  make local-api-app      # C 端 :18085"
	@echo "  make local-job"
	@echo "  make local-admin / local-merchant-admin / local-h5 / local-pc / local-service-web"
	@echo "  make local-backend-all / local-frontend-all"
	@echo "  make update-nginx"
	@echo ""
	@echo "矩阵: docs/release/SERVICE-MATRIX.md"
	@echo ""

init-env:
	@$(QX) init-env local

init-env-prod:
	@$(QX) init-env prod

check-config:
	@$(QX) check-config $(SVC)

pack:
	@$(QX) pack $(SVC)

pack-infra:
	@$(QX) pack infra-all

pack-backend:
	@$(QX) pack backend-all

pack-frontend:
	@$(QX) pack frontend-all

pack-all:
	@$(QX) pack infra-all
	@$(QX) pack backend-all
	@$(QX) pack frontend-all
	@$(QX) pack opts

local:
	@$(QX) local $(SVC)

local-db:
	@$(QX) local db

local-mq:
	@$(QX) local mq

local-api-admin:
	@$(QX) local api-admin

local-api-app:
	@$(QX) local api-app

local-job:
	@$(QX) local job

local-admin:
	@$(QX) local admin

local-merchant-admin:
	@$(QX) local merchant-admin

local-h5:
	@$(QX) local h5

local-pc:
	@$(QX) local pc

local-service-web:
	@$(QX) local service-web

local-backend-all:
	@$(QX) local backend-all

local-frontend-all:
	@$(QX) local frontend-all

local-up:
	@$(QX) local-up

local-compose-check:
	@$(QX) init-env local
	@$(QX) compose-check

up:
	@$(QX) up $(SVC)

down:
	@$(QX) down $(SVC)

restart:
	@$(QX) restart $(SVC)

ps:
	@$(QX) ps $(SVC)

deploy-db-reload:
	@$(QX) release db

deploy-mq-reload:
	@$(QX) release mq

deploy-api-admin:
	@$(QX) release api-admin

deploy-api-app:
	@$(QX) release api-app

deploy-job:
	@$(QX) release job

deploy-admin:
	@$(QX) release admin

deploy-merchant-admin:
	@$(QX) release merchant-admin

deploy-h5:
	@$(QX) release h5

deploy-pc:
	@$(QX) release pc

deploy-service-web:
	@$(QX) release service-web

deploy-backend-all:
	@$(QX) release backend-all

deploy-frontend-all:
	@$(QX) release frontend-all

deploy-all:
	@$(QX) release backend-all
	@$(QX) release frontend-all

update-nginx:
	@$(QX) upload-opts
