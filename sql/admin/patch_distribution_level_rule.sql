-- 平台「分销等级规则」协议键夹具（CRMEB CacheRepository::SYS_BROKERAGE = sys_brokerage）
-- 用法（项目根目录）：
--   docker exec -i pte_live_mysql sh -ec 'MYSQL_PWD="$MYSQL_ROOT_PASSWORD" mysql -uroot --default-character-set=utf8mb4' < sql/admin/patch_distribution_level_rule.sql

SET NAMES utf8mb4;

USE `qixi_crm_admin`;
INSERT INTO `qixi_crm_a_setting_cache` (`key`,`expire_time`,`result`) VALUES
  ('sys_brokerage',0,'<ol><li><p>第一级</p></li><li><p>第二级</p></li></ol>')
ON DUPLICATE KEY UPDATE
  `expire_time`=VALUES(`expire_time`),
  `result`=IF(`result` IS NULL OR `result`='',VALUES(`result`),`result`);

USE `qixi_crm_business`;
INSERT INTO `qixi_crm_b_content_view` (`content_id`,`content_type`,`title`,`cover_url`,`body`,`status`,`version`,`published_at`,`updated_at`) VALUES
  (2105,'agreement','sys_brokerage','','<ol><li><p>第一级</p></li><li><p>第二级</p></li></ol>',1,1,NOW(),NOW())
ON DUPLICATE KEY UPDATE
  `content_type`=VALUES(`content_type`),
  `title`=VALUES(`title`),
  `body`=IF(`body` IS NULL OR `body`='',VALUES(`body`),`body`),
  `status`=1,
  `updated_at`=NOW();
