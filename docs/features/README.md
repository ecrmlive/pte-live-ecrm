# 功能点清单（各端 → 页面 → 按钮/操作 → CRUD）

按你的要求：功能点梳理到**各端**、**具体功能**、**落实到按钮**，增删改查尽量落全。

## 怎么读

| 符号 | 含义 |
| --- | --- |
| C | 创建 / 新增 / 导入 / 提交 / 上传等 |
| R | 列表 / 详情 / 导出 / 统计等 |
| U | 编辑 / 状态开关 / 审核 / 排序等 |
| D | 删除 / 取消 / 清空 / 退款关闭等 |
| O | 特殊操作（支付回调、扫码登录等） |

后台以「菜单按钮权限」为主；缺按钮时用同页路由 `_path` 补全（文档内标注来源）。  
用户端与角色工作台以「一条 API 路由 = 一个操作」落实（等同按钮）。角色工作台必须归入统一后台或店铺管理系统，不能另建前端系统。

> **验收数量唯一口径**：[`../generated/features-master.tsv`](../generated/features-master.tsv) 共 2409 项（平台 1333、商户 615、用户 342、店员 95、客服 18、OpenAPI 6）。下表的页面覆盖统计保留其生成时的展示快照，用于按页面浏览，不能替代总表逐项验收。

## 分端文档

| 端 | 文档 | 页面/模块 | 操作数（约） |
| --- | --- | ---: | ---: |
| 统一后台：平台 / 商户 / 区域 / 客服 / 运营 | [01-platform-admin.md](./01-platform-admin.md) | CRMEB `/admin/` 角色菜单展示快照 | 1333（总表） |
| 店铺管理系统 | [02-merchant-admin.md](./02-merchant-admin.md) | CRMEB `/merchant/` 页面展示快照 | 615（总表） |
| 用户端（C 端 API） | [03-user-app.md](./03-user-app.md) | 183 小程序页附录 | 342 |
| 店铺管理系统内：员工履约工作台 | [04-merchant-mobile.md](./04-merchant-mobile.md) | CRMEB manager API 来源 | 95 |
| 统一后台内：客服工作台 | [05-customer-service.md](./05-customer-service.md) | CRMEB service API 来源 | 18 |
| OpenAPI | [06-openapi.md](./06-openapi.md) | openapi | 6 |
| 路由按菜单路径补全 | [07-route-crud-by-menu-path.md](./07-route-crud-by-menu-path.md) | 对照用 | — |
| 缺口清单 | [08-gaps.md](./08-gaps.md) | 空页+CRUD 矩阵真缺口均为 0 | ✅ 基线已锁定 |

机器可读总表：[`../generated/features-master.tsv`](../generated/features-master.tsv)（2409 项）；验收流程见 [CRMEB 全端功能验收总清单](../CRMEB-FULL-FUNCTION-CHECKLIST.md)。

## 数据来源

1. `eb_system_menu`（2419 节点：导航 + 按钮权限）  
2. `route/admin`、`route/merchant`（带 `_alias` / `_path`）  
3. `route/api/*.php`、`route/service.php`、`route/openapi*`  
4. 小程序 `extend/mp-weixin/v4.0/app.json` 页面树  

外部源码根：`~/Downloads/CRMEB多商户系统/CRMEB_MER_v4.0`

## 完整度说明

- **已做到**：各端拆分；后台到按钮权限名；用户/客服/店员到接口操作；每条操作带 CRUD 分类；总表可检索。  
- **功能基线**：**已锁定**（2026-07-21 用户确认）。允许按文档进入技术方案与业务编码。  
- **实现时对照**：`docs/api/FUNCTIONAL-TRUTH.md`（下单/支付/退款/计价）、`docs/schema/`（三库前缀）。

可选后续贴标（不阻塞编码）：用户端 UI 文案 ↔ API；配送员/服务人员能力归类。

## 建议开发对照顺序

1. 平台后台：商品 / 订单 / 商户 / 财务 / 营销  
2. 商户后台：商品 / 订单 / 营销 / 财务  
3. 用户端：登录、地址、购物车、v2 下单支付、售后、分销  
