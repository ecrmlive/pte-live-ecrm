# 数据库设计边界

完整规则见 [SYSTEM-ARCHITECTURE.md](../SYSTEM-ARCHITECTURE.md)。旧 `qixi_m_*` 映射、旧运行 DDL 与旧种子数据已失效，不能再执行。

| 数据库 | 前缀 | 第一批对象 |
| --- | --- | --- |
| `qixi_crm_admin` | `qixi_crm_a_` | 后台账号、角色、菜单、权限、数据范围、区域、运营配置、后台审计 |
| `qixi_crm_business` | `qixi_crm_b_` | C 端用户/身份、商品消费视图、购物车、订单、支付、售后、营销、资金 |
| `qixi_crm_merchant` | `qixi_crm_m_` | 店铺、店铺账号/员工、商品经营、库存、履约、物流、装修和店铺日志 |
| pte-live-im | 以 `pte-live-im/sql/init_im_schema.sql` 为准 | 仅 IM 管理与 IM 会话消息；由 pte-live-im 仓库维护 |

SQL 将在下一阶段按三库集中重建：每库分别提供结构、基础数据、配置、密钥占位和测试数据。所有 DDL 使用 utf8mb4；不提交真实密钥、Token、密码或真实用户数据。
