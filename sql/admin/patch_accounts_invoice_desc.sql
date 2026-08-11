-- 平台「发票说明」：协议键 sys_receipt_agree（CRMEB CacheRepository::RECEIPT_AGREE）+ 按钮权限
SET NAMES utf8mb4;

USE `qixi_crm_admin`;

INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (544,188,'accounts.invoice.desc','发票说明','lucide:file-text','/accounts/invoiceDesc','page',2,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

-- RequireAdminMenu 仅认 kind=button
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (21105,544,'accounts.invoice.desc.read','查看发票说明','','accounts/invoiceDesc','button',1,1),
  (21106,544,'accounts.invoice.desc.manage','维护发票说明','','accounts/invoiceDesc','button',2,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` AS r
CROSS JOIN `qixi_crm_a_menu` AS m
WHERE r.code IN ('platform','operations')
  AND m.code IN (
    'accounts.invoice',
    'accounts.invoice.desc',
    'accounts.invoice.desc.read',
    'accounts.invoice.desc.manage'
  );

INSERT INTO `qixi_crm_a_setting_cache` (`key`,`expire_time`,`result`) VALUES
  ('sys_receipt_agree',0,'<p>发票说明测试-发票说明测试-发票说明测试-发票说明测试-发票说明测试-</p>')
ON DUPLICATE KEY UPDATE
  `expire_time`=VALUES(`expire_time`),
  `result`=IF(`result` IS NULL OR `result`='',VALUES(`result`),`result`);

USE `qixi_crm_business`;
INSERT INTO `qixi_crm_b_content_view` (`content_id`,`content_type`,`title`,`cover_url`,`body`,`status`,`version`,`published_at`,`updated_at`) VALUES
  (2112,'agreement','sys_receipt_agree','','<p>发票说明测试-发票说明测试-发票说明测试-发票说明测试-发票说明测试-</p>',1,1,NOW(),NOW())
ON DUPLICATE KEY UPDATE
  `content_type`=VALUES(`content_type`),
  `title`=VALUES(`title`),
  `body`=IF(`body` IS NULL OR `body`='',VALUES(`body`),`body`),
  `status`=1,
  `updated_at`=NOW();
