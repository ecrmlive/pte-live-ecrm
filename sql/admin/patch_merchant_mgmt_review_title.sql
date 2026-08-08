-- 商户管理菜单：商户入驻审核 → 店铺入驻申请（可重复执行）
SET NAMES utf8mb4;
USE `qixi_crm_admin`;

UPDATE `qixi_crm_a_menu`
SET `title`='店铺入驻申请',
    `icon`='lucide:badge-check',
    `route_path`='/merchant/review',
    `kind`='page',
    `sort`=2,
    `parent_id`=25,
    `status`=1
WHERE `id`=27 OR `code`='merchant.mgmt.review';

INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`) VALUES
  (21003,12,'merchant.intention.create','新增店铺入驻申请','','merchant/audit','button',0)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `title`=VALUES(`title`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

UPDATE `qixi_crm_a_menu` SET `title`='审核店铺入驻申请' WHERE `code`='merchant.intention.audit';
UPDATE `qixi_crm_a_menu` SET `title`='分配店铺入驻申请区域' WHERE `code`='merchant.intention.assign_region';
UPDATE `qixi_crm_a_menu` SET `title`='删除店铺入驻申请' WHERE `code`='merchant.intention.delete';

INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, 21003
FROM `qixi_crm_a_role` AS r
WHERE r.code IN ('platform');
