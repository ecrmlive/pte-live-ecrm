# 功能解析完整度说明

**当前标准（按你的要求）**：功能点必须落到各端、具体功能、按钮级，增删改查尽量齐全。

权威清单入口：[CRMEB 全端功能验收总清单](./CRMEB-FULL-FUNCTION-CHECKLIST.md)。
操作数量和逐项记录以 `generated/features-master.tsv` 为唯一口径；各端 Markdown 的页面覆盖统计是较早的展示快照，不能覆盖总表。

## 锁定状态

| 项 | 状态 |
| --- | --- |
| 功能基线锁定 | **已锁定**（用户确认，2026-07-21） |
| 允许写业务代码 | **是**（按方案开发；高风险域先对齐状态机） |
| 验收主文档 | `docs/features/` |
| 接口/状态机对照 | `docs/api/FUNCTIONAL-TRUTH.md`、`docs/api/ACCURACY.md` |
| 表结构 | `docs/schema/`（`qixi_crm_a_` / `qixi_crm_b_` / `qixi_crm_m_`，详见系统架构总则） |

## 1. 已落实

| 端 | 粒度 | 文档 | 操作数（约） |
| --- | --- | --- | ---: |
| 平台后台 | 页面 → 按钮权限（不足则路由补全） | `generated/features-master.tsv` | 1333 |
| 商户后台 | 同上 | `generated/features-master.tsv` | 615 |
| 用户端 | API 路由 = 操作；附录小程序页 | `features/03-user-app.md` | 342 |
| 店铺系统员工履约角色 | API 路由 = 操作 | `features/04-merchant-mobile.md` | 95 |
| 统一后台客服角色 | API 路由 = 操作 | `features/05-customer-service.md` | 18 |
| OpenAPI | API 路由 = 操作 | `features/06-openapi.md` | 6 |
| 总表 | TSV 可检索 | `generated/features-master.tsv` | 2409 |

每条操作带 CRUD 分类（C/R/U/D/O）。后台页面有 CRUD 覆盖一览表。

## 2. 缺口结案（锁定依据）

见 [`features/08-gaps.md`](./features/08-gaps.md)：

| 类型 | 数量 | 状态 |
| --- | ---: | --- |
| 原「无操作页面」 | 53 | ✅ 真缺口 0 |
| 叶子菜单无按钮 | 14 | ✅ 已结案 |
| CRUD 不完整的页面 | 169 | ✅ 真缺口 0 |

可选后续抽检（不阻塞编码）：

1. 用户端页面按钮文案 ↔ API 操作名。  
2. 配送员 / 服务人员能力是否都在 `manager` / 用户 API。  
3. 功能表 PNG「标准多商户版」红字项抽检。

### API / 源码核对进度（文档正确性）

| 项 | 状态 |
| --- | --- |
| 路由→控制器方法存在性 | ✅ 1999 high / 26 stale / 2 unresolved（0 false-high） |
| 下单真实入口（普通 v2 / 积分 order/v3 / 死路由） | ✅ `docs/api/FUNCTIONAL-TRUTH.md` |
| 购物车/地址关键参数（checkParams） | ✅ 已按源码补全 |
| 支付回调渠道 + PAY_TYPE | ✅ FUNCTIONAL-TRUTH §5 |
| 退款/提现状态机 | ✅ FUNCTIONAL-TRUTH §6–7 |
| schema 165 表与 install SQL 一一对应 | ✅ |
| `08-gaps` 空页 + CRUD 矩阵 | ✅ 真缺口 0 |
| 营销计价主规则（券/SVIP/积分/单独购） | ✅ FUNCTIONAL-TRUTH §8 |
| 响应 `data` 字段全量建模 | ❌（动态结构，实现时对照 Repository） |

## 3. 与「可以写代码」的关系

| 条件 | 状态 |
| --- | --- |
| 各端按钮/操作级清单已生成 | ✅ |
| 缺口清单已列出并结案 | ✅ |
| 功能基线锁定 | ✅ |
| 允许写业务代码 | **是** |

## 4. 锁定后工作方式

1. 方案 → 编码；验收以 `docs/features/` 为准。  
2. 高风险域（下单/支付/退款/库存/券/积分/佣金/结算）先读 `FUNCTIONAL-TRUTH.md` 与状态机。  
3. 表前缀按系统架构总则使用 `qixi_crm_a_`、`qixi_crm_b_` 或 `qixi_crm_m_`；店铺接口强制店铺范围隔离。
4. 验证与部署仅在你明确要求时执行。
