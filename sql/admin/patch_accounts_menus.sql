-- 平台「财务」侧栏：店铺结算 / 用户结算 / 发票管理 二级目录
-- 执行后需重新登录平台后台以刷新菜单缓存。

SET NAMES utf8mb4;
USE `qixi_crm_admin`;

UPDATE `qixi_crm_a_menu`
SET `title`='财务',
    `icon`='ant-design:bar-chart-outlined',
    `route_path`='/accounts',
    `kind`='directory',
    `sort`=9,
    `parent_id`=0,
    `status`=1
WHERE `id`=100 OR `code`='accounts';

-- 旧扁平叶子隐藏（由下方树承接）
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=990 WHERE `id`=103 OR `code`='accounts.merchant_settlement';

-- 二级目录
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (536,100,'accounts.merchant.dir','店铺结算','lucide:landmark','/mer/accounts','directory',1,1),
  (537,100,'accounts.user.dir','用户结算','lucide:wallet','/accounts/record','directory',2,1),
  (188,100,'accounts.invoice','发票管理','lucide:receipt','/accounts/accounts','directory',3,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=100,
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`='directory',
  `sort`=VALUES(`sort`),
  `status`=1;

-- 三级：店铺结算
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (538,536,'accounts.statement','平台账单','lucide:file-spreadsheet','/accounts/statement','page',1,1),
  (539,536,'accounts.transfer','转账记录','lucide:arrow-left-right','/accounts/transferRecord','page',2,1),
  (540,536,'accounts.profitsharing','分账管理','lucide:split','/merchant/applyList','page',3,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=536,
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`='page',
  `sort`=VALUES(`sort`),
  `status`=1;

-- 三级：用户结算
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (101,537,'accounts.withdraw','提现管理','lucide:wallet','/accounts/extract','page',1,1),
  (541,537,'accounts.recharge_record','充值记录','lucide:circle-dollar-sign','/accounts/bill','page',2,1),
  (102,537,'accounts.user_assets','资金记录','lucide:wallet-cards','/accounts/capital','page',3,1),
  (542,537,'accounts.capital_flow','资金流水','lucide:list','/accounts/capitalFlow','page',4,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=537,
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`='page',
  `sort`=VALUES(`sort`),
  `status`=1;

-- 三级：发票管理
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (543,188,'accounts.invoice.list','发票列表','lucide:receipt','/accounts/receipt','page',1,1),
  (544,188,'accounts.invoice.desc','发票说明','lucide:file-text','/accounts/invoiceDesc','page',2,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=188,
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`='page',
  `sort`=VALUES(`sort`),
  `status`=1;

UPDATE `qixi_crm_a_menu` SET `parent_id`=101, `route_path`='accounts/extract' WHERE `id`=20923 OR `code`='accounts.withdraw.review';
UPDATE `qixi_crm_a_menu` SET `parent_id`=539, `route_path`='accounts/settings' WHERE `id`=20999 OR `code`='accounts.transfer_settings.manage';
UPDATE `qixi_crm_a_menu` SET `parent_id`=102, `route_path`='accounts/capital' WHERE `id`=20924 OR `code`='accounts.user_assets.read';
UPDATE `qixi_crm_a_menu` SET `parent_id`=543, `route_path`='accounts/receipt' WHERE `id`=20983 OR `code`='accounts.invoice.read';

INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` AS r
CROSS JOIN `qixi_crm_a_menu` AS m
WHERE r.code = 'platform'
  AND m.id IN (100,101,102,188,536,537,538,539,540,541,542,543,544);
