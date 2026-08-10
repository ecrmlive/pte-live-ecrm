-- 平台「佣金说明」协议键夹具（CRMEB CacheRepository::EXTENSION_AGREE = sys_extension_agree）
-- 用法：make local-sync-sql 或 scripts/local-dev-sync.sh sql
SET NAMES utf8mb4;

USE `qixi_crm_admin`;
INSERT INTO `qixi_crm_a_setting_cache` (`key`,`expire_time`,`result`) VALUES
  ('sys_extension_agree',0,'<p>我是佣金说明</p>')
ON DUPLICATE KEY UPDATE
  `expire_time`=VALUES(`expire_time`),
  `result`=IF(`result` IS NULL OR `result`='',VALUES(`result`),`result`);

USE `qixi_crm_business`;
INSERT INTO `qixi_crm_b_content_view` (`content_id`,`content_type`,`title`,`cover_url`,`body`,`status`,`version`,`published_at`,`updated_at`) VALUES
  (2106,'agreement','sys_extension_agree','','<p>我是佣金说明</p>',1,1,NOW(),NOW())
ON DUPLICATE KEY UPDATE
  `content_type`=VALUES(`content_type`),
  `title`=VALUES(`title`),
  `body`=IF(`body` IS NULL OR `body`='',VALUES(`body`),`body`),
  `status`=1,
  `updated_at`=NOW();
