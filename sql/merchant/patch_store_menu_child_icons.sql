-- 店铺菜单子级 Iconify（ant-design:*）：页面/目录补图标
-- 纯按钮（is_menu=2、无 path）默认空；「权限」顶层 8609 / 首页子级 8614 与带 path 的权限页例外，须有图标
-- 用法：docker exec -i pte_live_mysql sh -ec 'MYSQL_PWD="$MYSQL_ROOT_PASSWORD" mysql -uroot --default-character-set=utf8mb4 qixi_crm_merchant' < sql/merchant/patch_store_menu_child_icons.sql

SET NAMES utf8mb4;
USE `qixi_crm_merchant`;

-- 顶层（与 patch_store_menu_top_sort 对齐；可重复执行）
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:home-outlined' WHERE `id` = 55;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:shopping-outlined' WHERE `id` = 95;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:file-text-outlined' WHERE `id` = 512;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:bell-outlined' WHERE `id` = 106;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:field-time-outlined' WHERE `id` = 525;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:user-outlined' WHERE `id` = 1027;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:customer-service-outlined' WHERE `id` = 1661;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:format-painter-outlined' WHERE `id` = 1671;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:safety-outlined' WHERE `id` = 8609;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:setting-outlined' WHERE `id` = 526;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:notification-outlined' WHERE `id` = 1119;

UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:safety-outlined' WHERE `id` = 49;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:idcard-outlined' WHERE `id` = 50;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:idcard-outlined' WHERE `id` = 51;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:profile-outlined' WHERE `id` = 52;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:picture-outlined' WHERE `id` = 54;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:setting-outlined' WHERE `id` = 74;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:appstore-outlined' WHERE `id` = 96;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:cluster-outlined' WHERE `id` = 99;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:shopping-outlined' WHERE `id` = 105;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:gift-outlined' WHERE `id` = 107;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:gift-outlined' WHERE `id` = 115;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:customer-service-outlined' WHERE `id` = 120;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:gift-outlined' WHERE `id` = 122;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:file-text-outlined' WHERE `id` = 513;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:rollback-outlined' WHERE `id` = 528;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:star-outlined' WHERE `id` = 544;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:shop-outlined' WHERE `id` = 546;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:car-outlined' WHERE `id` = 700;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:video-camera-outlined' WHERE `id` = 785;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:video-camera-outlined' WHERE `id` = 786;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:shopping-outlined' WHERE `id` = 787;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:thunderbolt-outlined' WHERE `id` = 788;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:audit-outlined' WHERE `id` = 789;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:transaction-outlined' WHERE `id` = 790;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:history-outlined' WHERE `id` = 1025;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:tags-outlined' WHERE `id` = 1028;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:tags-outlined' WHERE `id` = 1029;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:tags-outlined' WHERE `id` = 1030;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:like-outlined' WHERE `id` = 1099;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:shopping-outlined' WHERE `id` = 1100;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:like-outlined' WHERE `id` = 1101;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:file-done-outlined' WHERE `id` = 1102;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:user-outlined' WHERE `id` = 1103;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:gift-outlined' WHERE `id` = 1132;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:team-outlined' WHERE `id` = 1139;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:shopping-outlined' WHERE `id` = 1140;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:team-outlined' WHERE `id` = 1141;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:picture-outlined' WHERE `id` = 1176;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:swap-outlined' WHERE `id` = 1181;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:credit-card-outlined' WHERE `id` = 1182;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:account-book-outlined' WHERE `id` = 1183;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:file-protect-outlined' WHERE `id` = 1246;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:search-outlined' WHERE `id` = 1286;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:trophy-outlined' WHERE `id` = 1293;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:trophy-outlined' WHERE `id` = 1294;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:trophy-outlined' WHERE `id` = 1295;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:account-book-outlined' WHERE `id` = 1304;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:cluster-outlined' WHERE `id` = 1366;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:cloud-outlined' WHERE `id` = 1376;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:setting-outlined' WHERE `id` = 1377;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:cloud-outlined' WHERE `id` = 1378;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:car-outlined' WHERE `id` = 1380;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:tags-outlined' WHERE `id` = 1468;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:flag-outlined' WHERE `id` = 1471;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:video-camera-outlined' WHERE `id` = 1594;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:percentage-outlined' WHERE `id` = 1630;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:car-outlined' WHERE `id` = 1649;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:car-outlined' WHERE `id` = 1650;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:car-outlined' WHERE `id` = 1651;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:wallet-outlined' WHERE `id` = 1652;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:picture-outlined' WHERE `id` = 1656;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:format-painter-outlined' WHERE `id` = 1672;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:message-outlined' WHERE `id` = 5123;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:crown-outlined' WHERE `id` = 9053;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:form-outlined' WHERE `id` = 9059;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:key-outlined' WHERE `id` = 9172;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:form-outlined' WHERE `id` = 9222;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:gold-outlined' WHERE `id` = 9223;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:cloud-outlined' WHERE `id` = 9286;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:thunderbolt-outlined' WHERE `id` = 9289;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:shopping-outlined' WHERE `id` = 9290;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:barcode-outlined' WHERE `id` = 9329;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:car-outlined' WHERE `id` = 9330;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:customer-service-outlined' WHERE `id` = 9331;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:car-outlined' WHERE `id` = 9332;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:picture-outlined' WHERE `id` = 9358;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:picture-outlined' WHERE `id` = 9359;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:file-text-outlined' WHERE `id` = 9369;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:customer-service-outlined' WHERE `id` = 9404;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:schedule-outlined' WHERE `id` = 9405;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:calendar-outlined' WHERE `id` = 9406;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:comment-outlined' WHERE `id` = 9408;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:customer-service-outlined' WHERE `id` = 9410;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:car-outlined' WHERE `id` = 9452;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:appstore-outlined' WHERE `id` = 9455;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:car-outlined' WHERE `id` = 9456;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:picture-outlined' WHERE `id` = 9457;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:safety-outlined' WHERE `id` = 9458;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:dashboard-outlined' WHERE `id` = 9499;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:bar-chart-outlined' WHERE `id` = 9500;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:line-chart-outlined' WHERE `id` = 9501;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:customer-service-outlined' WHERE `id` = 9556;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:bar-chart-outlined' WHERE `id` = 9557;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:car-outlined' WHERE `id` = 9558;
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:bar-chart-outlined' WHERE `id` = 9559;

-- 纯按钮权限点清空（无 path）；再回填需展示的「权限」节点
UPDATE `qixi_crm_m_menu`
SET `icon` = ''
WHERE `is_menu` = 2
  AND TRIM(IFNULL(`path`, '')) = '';

UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:safety-outlined' WHERE `id` IN (8609, 8614);

UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:printer-outlined' WHERE `name` IN ('打印配置','小票打印');

-- 名称兜底：页面/目录（is_menu=1）且 icon 仍空
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:home-outlined' WHERE `is_menu` = 1 AND TRIM(IFNULL(`icon`,'')) = '' AND `name` = '首页';
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:shopping-outlined' WHERE `is_menu` = 1 AND TRIM(IFNULL(`icon`,'')) = '' AND `name` IN ('商品','商品列表','助力商品','拼团商品列表','秒杀商品','直播商品管理');
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:file-text-outlined' WHERE `is_menu` = 1 AND TRIM(IFNULL(`icon`,'')) = '' AND `name` IN ('订单','订单管理','代客下单');
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:bell-outlined' WHERE `is_menu` = 1 AND TRIM(IFNULL(`icon`,'')) = '' AND `name` = '营销';
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:field-time-outlined' WHERE `is_menu` = 1 AND TRIM(IFNULL(`icon`,'')) = '' AND `name` = '财务';
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:user-outlined' WHERE `is_menu` = 1 AND TRIM(IFNULL(`icon`,'')) = '' AND `name` IN ('用户','用户管理');
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:customer-service-outlined' WHERE `is_menu` = 1 AND TRIM(IFNULL(`icon`,'')) = '' AND `name` IN ('员工','店员列表','店员管理','店员配置','服务人员');
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:format-painter-outlined' WHERE `is_menu` = 1 AND TRIM(IFNULL(`icon`,'')) = '' AND `name` IN ('装修','配置管理');
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:safety-outlined' WHERE `is_menu` = 1 AND TRIM(IFNULL(`icon`,'')) = '' AND `name` IN ('权限','权限管理');
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:setting-outlined' WHERE `is_menu` = 1 AND TRIM(IFNULL(`icon`,'')) = '' AND `name` IN ('设置','店铺配置');
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:notification-outlined' WHERE `is_menu` = 1 AND TRIM(IFNULL(`icon`,'')) = '' AND `name` = '公告';
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:dashboard-outlined' WHERE `is_menu` = 1 AND TRIM(IFNULL(`icon`,'')) = '' AND `name` = '控制台';
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:bar-chart-outlined' WHERE `is_menu` = 1 AND TRIM(IFNULL(`icon`,'')) = '' AND `name` IN ('商品统计','服务统计','配送统计');
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:line-chart-outlined' WHERE `is_menu` = 1 AND TRIM(IFNULL(`icon`,'')) = '' AND `name` = '订单统计';
UPDATE `qixi_crm_m_menu` SET `icon` = 'ant-design:account-book-outlined' WHERE `is_menu` = 1 AND TRIM(IFNULL(`icon`,'')) = '' AND `name` IN ('账单管理','申请分账商户');
