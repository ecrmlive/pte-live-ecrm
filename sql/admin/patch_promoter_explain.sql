-- 平台「分销说明」协议键夹具（CRMEB CacheRepository::PROMOTER_EXPLAIN = promoter_explain）
-- 用法：make local-sync-sql
SET NAMES utf8mb4;

USE `qixi_crm_admin`;
INSERT INTO `qixi_crm_a_setting_cache` (`key`,`expire_time`,`result`) VALUES
  ('promoter_explain',0,'<p>这个是分销说明</p>')
ON DUPLICATE KEY UPDATE
  `expire_time`=VALUES(`expire_time`),
  `result`=IF(`result` IS NULL OR `result`='',VALUES(`result`),`result`);

USE `qixi_crm_business`;
INSERT INTO `qixi_crm_b_content_view` (`content_id`,`content_type`,`title`,`cover_url`,`body`,`status`,`version`,`published_at`,`updated_at`) VALUES
  (2107,'agreement','promoter_explain','','<p>这个是分销说明</p>',1,1,NOW(),NOW())
ON DUPLICATE KEY UPDATE
  `content_type`=VALUES(`content_type`),
  `title`=VALUES(`title`),
  `body`=IF(`body` IS NULL OR `body`='',VALUES(`body`),`body`),
  `status`=1,
  `updated_at`=NOW();
