SET NAMES utf8mb4;
USE `qixi_crm_merchant`;

-- 商品、店铺、用户及登录账号均不是系统默认配置，统一仅在 init_test_data.sql 提供演示夹具。

-- 店铺菜单（qixi_crm_m_menu / role_menu）由 init_menu_crmeb_full.sql 全量导入：
-- 源：CRMEB 线上 GET /sys/merchant/menu/lst（https://mer.crmeb.net/admin/merchant/system），共 711 条。
-- 由 scripts/release/db-reset.sh 与 scripts/qixi-crm.sh db-init 在 merchant/init_data 之后执行。
