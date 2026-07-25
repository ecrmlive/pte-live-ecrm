# sql/

可执行迁移目录。权威字段以 `docs/schema/` 为准；本目录存放按阶段拆分的迁移脚本。

## 约定

- 表前缀固定 `qixi_`
- 文件名：`NNN_description.sql`（递增）
- Phase 0：`000_init_schema.sql` 从 `docs/schema/qixi_schema_reference.sql` 摘取/同步入口（全量参考仍在 docs）

## 应用方式（本机）

MySQL 由 **pte-live-im** 的 `pte_live_mysql` 提供（宿主口仍 `13306`）。首次：

```bash
# 创建库与账号 qixi/qixi_local（需 IM root 密码）
docker exec -i pte_live_mysql mysql -uroot -p"$MYSQL_ROOT_PASSWORD" < sql/000_shared_im_mysql_bootstrap.sql
```

然后灌迁移：

> 导入时必须显式使用 `--default-character-set=utf8mb4`，避免中文种子数据按错误连接字符集写入数据库。

```bash
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/000_init_schema.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/001_identity.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/002_platform_merchant_catalog.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/003_merchant_product_app_catalog.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/004_trade_cart_order.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/005_aftersale_finance.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/006_promotion_coupon_spread.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/007_loyalty_points.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/008_openapi.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/009_manager_service.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/010_diy_seckill.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/011_combination.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/012_svip.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/013_reservation.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/014_presell.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/015_presell_deposit.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/016_broadcast.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/017_community.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/018_assist.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/019_stage7_menu_fix.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/020_stage7_hide_settings.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/021_stage7_settings_crud.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/022_stage7_rbac_menu.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/023_stage7_merchant_admins.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/024_stage7_button_perms.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/025_stage7_refund_button_perms.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/026_stage7_attachment.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/027_stage7_platform_button_perms.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/028_stage7_product_button_perms.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/029_stage7_product_community_button_perms.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/030_stage7_marketing_broadcast_button_perms.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/031_stage7_coupon_button_perms.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/032_stage7_presell_assist_button_perms.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/033_stage7_coupon_create_button_perms.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/034_stage7_coupon_activity_crud_button_perms.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/035_stage7_reservation_svip_broadcast_button_perms.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/036_stage7_broadcast_attachment_button_perms.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/037_stage7_diy_staff_button_perms.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/038_stage7_setting_write_button_perms.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/039_stage7_service_reply.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/040_stage7_merchant_community_write.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/041_stage7_agreement_notice.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/042_gap_fill_cs_delivery_invoice.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/043_crmeb_system_menu_full.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/044_im_remote_bridge.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/045_gap_logistics_productmeta_article_usertag.sql
mysql --default-character-set=utf8mb4 -h127.0.0.1 -P13306 -uqixi -pqixi_local qixi_mergers < sql/046_diy_visual_doc.sql
```

容器内主机名：`pte_live_mysql:3306`（库 `qixi_mergers`）。

## 阶段进度

| 文件 | 阶段 | 说明 |
| --- | --- | --- |
| `000_init_schema.sql` | 0 | 库与 meta |
| `001_identity.sql` | 1 | 平台/商户管理员、C 端用户、角色菜单、演示种子 |
| `002_platform_merchant_catalog.sql` | 2 | 入驻申请、平台类目品牌、商品审核 |
| `003_merchant_product_app_catalog.sql` | 2 | 商户 SKU 表 + C 端可售种子 |
| `004_trade_cart_order.sql` | 3 | 地址/购物车/主子单/订单商品 + 双店种子 |
| `005_aftersale_finance.sql` | 4 | 退款单/提现 + 菜单；商户 mer_money 种子 |
| `006_promotion_coupon_spread.sql` | 5 | 券/用户券/分销日志/账单；平台+店铺券种子 |
| `007_loyalty_points.sql` | 6 | 积分字段/积分商品(type=1,product_type=20) + 公告 |
| `008_openapi.sql` | 6c | `qixi_open_auth` + 演示 access_key |
| `009_manager_service.sql` | 6d | 店员 `qixi_store_service` + 核销码补齐 |
| `010_diy_seckill.sql` | 6a/活动 | DIY 首页 + 秒杀场次/活动 + `activity_id` + 菜单 |
| `011_combination.sql` | 6/活动 | 拼团商品/团次/成员 + 菜单 |
| `012_svip.sql` | 6/SVIP | 会员价字段、`svip_coupon_merge`、demo 永久 SVIP |
| `013_reservation.sql` | 6e | 预约商品/时段 |
| `014_presell.sql` | 6/活动 | 全款预售活动 |
| `015_presell_deposit.sql` | 6/活动 | 定金预售字段 + `qixi_presell_order` 尾款单 |
| `016_broadcast.sql` | 6e | 直播间 + 挂货（无微信推流） |
| `017_community.sql` | 6e | 社区种草（分类/话题/帖/评论） |
| `018_assist.sql` | 6/活动 | 助力活动/助力单/助力人；商品13 演示 + 菜单 |
| `019_stage7_menu_fix.sql` | 7 | 补平台退款/财务/提现菜单；社区 pid→内容 |
| `020_stage7_hide_settings.sql` | 7 | （过渡）隐藏设置占位；商户社区 path；DIY 预约 |
| `021_stage7_settings_crud.sql` | 7 | 恢复设置菜单可见（管理员/角色/菜单/店铺/店员） |
| `022_stage7_rbac_menu.sql` | 7 | 商户「角色权限」菜单 + role rules |
| `023_stage7_merchant_admins.sql` | 7 | 子账号菜单 `/setting/admins`；演示角色3 + mersub |
| `024_stage7_button_perms.sql` | 7 | 发货/核销按钮（is_menu=2）；mersub 仅发货 |
| `025_stage7_refund_button_perms.sql` | 7 | 同意/拒绝退款按钮；mersub 仅同意 |
| `026_stage7_attachment.sql` | 7 | 素材库表 + 平台菜单33 / 商户菜单131 |
| `027_stage7_platform_button_perms.sql` | 7 | 平台退款/提现按钮 + 演示 auditor |
| `028_stage7_product_button_perms.sql` | 7 | 商品上下架/库存按钮 + 演示 merprod |
| `029_stage7_product_community_button_perms.sql` | 7 | 商品发布/删除 + 社区审帖/删帖；auditor 可审不可删 |
| `030_stage7_marketing_broadcast_button_perms.sql` | 7 | 秒杀/拼团启停 + 直播审房；演示 meract |
| `031_stage7_coupon_button_perms.sql` | 7 | 平台/商户优惠券启停；meract/auditor 可启停 |
| `032_stage7_presell_assist_button_perms.sql` | 7 | 预售/助力上下架；meract 不赋权（不对称） |
| `033_stage7_coupon_create_button_perms.sql` | 7 | 新建优惠券；meract/auditor 可启停不可新建 |
| `034_stage7_coupon_activity_crud_button_perms.sql` | 7 | 券删除 + 秒杀/拼团/预售/助力创建删除；meract 不赋权 |
| `035_stage7_reservation_svip_broadcast_button_perms.sql` | 7 | 直播建房/删除 + 预约配置 + 商户/平台 SVIP；meract/auditor 无写 |
| `036_stage7_broadcast_attachment_button_perms.sql` | 7 | 直播开播/挂货 + 平台/商户素材上传删除；meract/auditor 无写 |
| `037_stage7_diy_staff_button_perms.sql` | 7 | DIY create/update/delete/active/pick + staff1 发货(is_goods) |
| `038_stage7_setting_write_button_perms.sql` | 7 | 商户店铺/店员/子账号/角色写按钮；meract/mersub 无 |
| `039_stage7_service_reply.sql` | 7 | 客服快捷回复表 + 菜单163/按钮164；mer1 种子；service 按 mer_id 读库 |
| `040_stage7_merchant_community_write.sql` | 7 | 商户社区发帖/编辑/删除按钮；meract/mersub 无写 |
| `041_stage7_agreement_notice.sql` | 7 | 协议 qixi_cache + 菜单52/按钮53；C 端公告闭环 |
| `046_diy_visual_doc.sql` | DIY | 页面外观字段 + `{page,items[]}` 文档协议种子升级 |
