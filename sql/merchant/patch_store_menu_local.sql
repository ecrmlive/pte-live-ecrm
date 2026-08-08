-- DEPRECATED：店铺菜单全量请用 init_menu_crmeb_full.sql（693 条 is_mer=1）。
-- 本文件仅保留历史增量补丁，勿再对新环境执行。
-- 本地 Docker MySQL 增量补丁：店铺菜单 CRMEB 对齐 + created_at 列
-- 用法（在项目根目录）：
--   docker exec -i pte_live_mysql sh -ec 'MYSQL_PWD="$MYSQL_ROOT_PASSWORD" mysql -uroot --default-character-set=utf8mb4 qixi_crm_merchant' < sql/merchant/patch_store_menu_local.sql
-- 若 created_at 列已存在，首条 ALTER 会报错，可忽略后继续执行后续语句。

SET NAMES utf8mb4;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_m_menu' AND COLUMN_NAME='created_at')=0,
    'ALTER TABLE `qixi_crm_m_menu` ADD COLUMN `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP AFTER `status`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

UPDATE `qixi_crm_m_menu` SET `name`='首页' WHERE `id`=1;

INSERT INTO `qixi_crm_m_menu` (`id`,`parent_id`,`code`,`name`,`path`,`component`,`icon`,`is_menu`,`is_route`,`sort`,`status`) VALUES
  (110,0,'staff','员工','/server','','lucide:headset',1,0,55,1),
  (120,0,'permission','权限','self','','lucide:shield-check',1,0,65,1),
  (130,0,'notice','公告','/station/notice','views/ecrm/station/notice.vue','lucide:megaphone',1,1,75,1),
  (102,100,'diy.list','店铺装修','/devise/diy/list','views/ecrm/diy/list.vue','',1,1,1,1),
  (18,10,'product.label','商品标签','/product/label','views/ecrm/product/label.vue','',1,1,8,1),
  (19,10,'product.guarantee','保障服务','/config/guarantee','views/ecrm/product/guarantee.vue','',1,1,9,1),
  (37,30,'marketing.presell','预售活动','/marketing/presell/list','views/ecrm/marketing/presell.vue','',1,1,7,1),
  (38,30,'marketing.broadcast','直播管理','/marketing/broadcast','views/ecrm/marketing/broadcast.vue','',1,1,8,1),
  (39,30,'marketing.community','社区话题','/marketing/community','views/ecrm/marketing/community.vue','',1,1,9,1),
  (311,30,'marketing.coupon.records','优惠券记录','/marketing/coupon/send','views/ecrm/marketing/coupon-records.vue','',1,1,10,1),
  (66,50,'setting.attachment','素材管理','/setting/attachment','views/ecrm/setting/attachment.vue','',1,1,13,1),
  (67,50,'setting.sms','一号通','/setting/sms/sms_config/index','views/ecrm/setting/sms-yihaotong.vue','',1,1,14,1),
  (68,50,'setting.svip','付费会员','/systemForm/Basics/svip','views/ecrm/setting/svip.vue','',1,1,15,1)
ON DUPLICATE KEY UPDATE `parent_id`=VALUES(`parent_id`),`name`=VALUES(`name`),`path`=VALUES(`path`),`component`=VALUES(`component`),`icon`=VALUES(`icon`),`is_menu`=VALUES(`is_menu`),`is_route`=VALUES(`is_route`),`sort`=VALUES(`sort`),`status`=VALUES(`status`);

UPDATE `qixi_crm_m_menu` SET `parent_id`=110, `sort`=1 WHERE `id`=52;
UPDATE `qixi_crm_m_menu` SET `parent_id`=120, `sort`=1 WHERE `id`=53;
UPDATE `qixi_crm_m_menu` SET `parent_id`=120, `sort`=2 WHERE `id`=54;
UPDATE `qixi_crm_m_menu` SET `sort`=50 WHERE `id`=80;
UPDATE `qixi_crm_m_menu` SET `sort`=60 WHERE `id`=100;
UPDATE `qixi_crm_m_menu` SET `sort`=70 WHERE `id`=50;

INSERT IGNORE INTO `qixi_crm_m_role_menu` (`role_code`,`menu_id`) VALUES
  ('owner',110),('manager',110),('owner',120),('manager',120),('owner',130),('manager',130),
  ('owner',102),('manager',102),('owner',18),('manager',18),('owner',19),('manager',19),
  ('owner',37),('manager',37),('owner',38),('manager',38),('owner',39),('manager',39),
  ('owner',311),('manager',311),('owner',66),('manager',66),('owner',67),('manager',67),('owner',68),('manager',68);
