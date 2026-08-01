---
name: pte-live-ecrm
description: >-
  Maintain the pte-live-ecrm multi-merchant mall system. Use when working on
  this repository, multi-merchant marketplace features, platform/merchant admin,
  trade/catalog/finance domains, CRMEB MER v4.0 feature alignment, Vben admin,
  uni-app x / PC storefront, or Go Gin/GORM APIs for pte-live-ecrm.
disable-model-invocation: false
---

# pte-live-ecrm

多商户商城管理系统。功能基线对齐 CRMEB Merchant v4.0；实现栈为 Go（Gin/GORM）+ Vben 5+ + uni-app x。

## Read first

按顺序阅读（存在则必读）：

1. 仓库根目录 `AGENTS.md`
2. `docs/product-understanding.md`（我们要做什么：主链路/对象）
3. `docs/analysis-completeness.md`（完整度与是否允许写代码）
4. `docs/features/README.md`（各端→按钮→CRUD）与 `docs/schema/README.md`（表前缀 `qixi_m_admin_`（平台/商户）与 `qixi_m_app_`（C 端））
5. `docs/README.md` 与其它文档：
   - `docs/overview.md`
   - `docs/feature-matrix.md`
   - `docs/generated/EXTRACT-STATS.md`
   - `docs/architecture-target.md`
   - `docs/domain-modules.md`
   - `docs/roles-and-portals.md`
6. 涉及发布时：`docs/release/COMMANDS.md`、`docs/release/PACK-AND-CONFIG.md`
7. 本 Skill 参考：
   - [work-conventions.md](references/work-conventions.md)
   - [tech-stack.md](references/tech-stack.md)
   - [feature-map.md](references/feature-map.md)
   - [crmeb-pointer.md](references/crmeb-pointer.md)
   - [schema-prefix.md](references/schema-prefix.md)

## Analysis gate

功能基线**已锁定**（见 `docs/analysis-completeness.md`）：允许技术方案与业务编码。  
验收以 `docs/features/` 为准；高风险域先读 `docs/api/FUNCTIONAL-TRUTH.md`。

## Task scope

- 只做：制定方案、开发功能、查找原因、提供解决方案、修复问题。
- 验证与部署：仅在用户明确要求时执行。
- 未经用户明确授权：禁止 commit、push、创建 PR、重置分支、覆盖用户已有修改。
- 部署/服务器操作必须用户明确授权；部署完成后反馈服务器 IP。

## Product shape

三角色：

```text
平台后台 ── 商户审核 / 全局监管 / 结算 / RBAC
商户后台 ── 本店商品订单营销财务
用户端   ── 购买 / 营销 / 分销积分 / 售后客服
```

高风险域（改前先对齐状态机与幂等）：订单创建、支付回调、退款、库存、优惠券、积分、佣金、商户结算。

## Implementation rules

- 配置文件使用 `.yaml`。
- API：Gin、GORM、Swagger、MySQL、etcd、NATS。
- 管理后台：Vben 5+；C 端 H5/小程序用 uni-app x（`app-mp/`）；PC 商城用 Vue 3（`app-web/`）。
- **数据库表前缀固定为 `qixi_m_admin_`（平台/商户）与 `qixi_m_app_`（C 端）**（CRMEB `eb_` 映射）。见 `docs/schema/`。禁止新代码使用裸 `qixi_` 或 `eb_`。
- 商户接口强制 `mer_id` 隔离。
- 测试数据注意中文乱码（utf8mb4）。
- 本机构建产物 + Docker Compose 挂载运行；禁止服务器源码构建；Dockerfile 只复制本机产物。
- 不提交密码、Token、密钥、证书等敏感信息。

## CRMEB reference

需要对照原功能实现时，只读：

```text
~/Downloads/CRMEB多商户系统/CRMEB_MER_v4.0
```

重点：`app/common/repositories/`、`route/`、`install/crmeb_merchant.sql`。  
不要把该 PHP 工程提交进本仓库，不要以 ThinkPHP/Swoole 作为本仓库运行时。

## Default answer shape

用户问「怎么做 / 怎么改」时：

1. 目标与所属领域模块
2. 涉及文档 / 接口 / 表或状态机
3. 推荐实现步骤
4. 风险点（资金/库存/幂等）
5. 用户未要求则不主动部署
