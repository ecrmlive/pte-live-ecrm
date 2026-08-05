# Make 命令说明

本项目的本地、测试与生产操作都通过根目录 `Makefile` 执行。后端必须先在本机构建，Docker 只挂载构建产物，不在容器或服务器编译源码。

## 本地验收

```bash
make local-compose-check  # 校验 Compose 结构
make local-infra          # 确认 pte-live-im 共享基础设施已运行
make local-db-init        # 幂等导入三库结构与中文测试数据
make local-backend        # 启动七禧三个 API
make local-ps             # 查看服务状态
```

`make local-db-reset` 会删除并重建本地三库，仅在明确需要重置测试数据时使用。

## 打包与部署

```bash
make pack-backend
make deploy-backend-all
make deploy-frontend-all
make deploy-all
```

部署按 `release/` 目录中的本机构建产物和 YAML 配置执行；不得上传源代码到服务器构建。真实凭据仅存放在被 Git 忽略的本地 YAML 或密钥 SQL 中。
