# pte-live-ecrm 文档索引

本仓库目标：基于 CRMEB 多商户商城（MER v4.0）的功能基线，重建七禧 CRM 多商户商城系统。

## 阅读顺序

0. **[系统架构总则](./SYSTEM-ARCHITECTURE.md)** — **系统边界、账号/JWT、数据库、目录和迁移顺序的唯一口径**
1. **[JWT、店铺 AppId 与商户 IM SDK AppId 契约](./auth-store-appid-im-contract.md)** — 请求头、店铺上下文、IM 多 AppId 与迁移唯一口径
2. [我们要做什么（产品全景）](./product-understanding.md) — 业务主链路与对象
2. [功能解析完整度](./analysis-completeness.md) — 完整度与是否允许写代码
3. **[功能点清单（各端→按钮→CRUD）](./features/README.md)** — 验收主文档
4. **[原项目接口文档](./api/README.md)** — CRMEB 只读对照资料
5. **[数据边界](./schema/README.md)** — 三库、表前缀与 IM 边界
6. [项目总览](./overview.md) — 系统组成与建设范围
7. [功能地图（模块级摘要）](./feature-matrix.md) — 功能表 + 脑图摘要
8. [目标技术架构](./architecture-target.md) — 目标目录与服务边界
9. [角色与端入口](./roles-and-portals.md) — 菜单角色、账号与数据范围
12. **[服务命名 · Docker 网络/IP](./release/SERVICE-MATRIX.md)** — 各端命名、固定 IP、发布形态
13. [发布命令](./release/COMMANDS.md) · [Pack/Config](./release/PACK-AND-CONFIG.md)
14. **[接入 pte-live-im（客服）](./integration-pte-live-im.md)** — 与现有 IM 项目对接分析
16. **[云服务配置中心](./cloud-config-center.md)** — 支付、腾讯云、直播与 IM 的加密数据库配置
17. **[CRMEB → Vben 后台对齐验收](./crmeb-vben-parity.md)** — 各后台角色的真实实现状态与逐页验收标准
17b. **[管理后台布局强制标准（100%）](./acceptance/LAYOUT-FIDELITY-CHECKLIST.md)** — 以店铺列表为金标准；列表/抽屉页必须遵守
18. **[H5 / 小程序 DIY 装修运行链路](./diy-h5-runtime.md)** — 后台配置、启用发布与 uni-app x 动态渲染
19. **[CRMEB 全端功能验收总清单](./CRMEB-FULL-FUNCTION-CHECKLIST.md)** — 各端模块、2409 条逐操作基线入口与完成判定
20. **[CRMEB 全端实施与验收台账](./acceptance/CRMEB-PORTAL-PARITY-LEDGER.md)** — 当前实现证据、逐端阻断项与关闭规则

> 功能基线以 `docs/features/` 为准；系统与数据架构以 [`SYSTEM-ARCHITECTURE.md`](./SYSTEM-ARCHITECTURE.md) 为准。`features/08-gaps.md` 只说明 CRMEB 基线解析，不代表功能已实现。

## 原始素材

| 素材 | 路径 |
| --- | --- |
| 前端脑图 | [assets/前端脑图.png](./assets/前端脑图.png) |
| 功能表 1（移动端/用户侧） | [assets/多商户商城系统功能表1.png](./assets/多商户商城系统功能表1.png) |
| 功能表 2（后台管理） | [assets/多商户商城系统功能表2.png](./assets/多商户商城系统功能表2.png) |
| 功能表 3（商城端标准版） | [assets/多商户商城系统功能表3.png](./assets/多商户商城系统功能表3.png) |
| 功能表 4（PC 端） | [assets/多商户商城系统功能表4.png](./assets/多商户商城系统功能表4.png) |
| CRMEB 源码（外部参考，不入库） | `~/Downloads/CRMEB多商户系统/CRMEB_MER_v4.0` |

## Agent 约定

- 仓库根目录 [`AGENTS.md`](../AGENTS.md)
- Cursor Skill：[`.cursor/skills/pte-live-ecrm/SKILL.md`](../.cursor/skills/pte-live-ecrm/SKILL.md)
- Codex Skill 镜像：[`codex-skills/pte-live-ecrm/SKILL.md`](../codex-skills/pte-live-ecrm/SKILL.md)

发布与运行约定见 [服务命名 · Docker 网络/IP](./release/SERVICE-MATRIX.md)、[发布命令](./release/COMMANDS.md) 和 [Pack/Config](./release/PACK-AND-CONFIG.md)：local/test 共用 `pte_live_ecrm` Docker 身份，只能切换运行。
