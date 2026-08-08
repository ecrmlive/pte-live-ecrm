-- 平台「客服」侧栏：三个二级叶子同列（结构已扁平正确，补齐排序/图标）
-- 执行后需重新登录平台后台以刷新菜单缓存。

SET NAMES utf8mb4;
USE `qixi_crm_admin`;

UPDATE `qixi_crm_a_menu`
SET `title`='客服',
    `icon`='ant-design:customer-service-outlined',
    `route_path`='/service',
    `kind`='directory',
    `sort`=12,
    `parent_id`=0,
    `status`=1
WHERE `id`=30 OR `code`='service';

INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (301,30,'service.auto_reply','客服自动回复','lucide:message-square-reply','/systemForm/customer_keyword','page',1,1),
  (302,30,'service.customer.list','客服列表','ant-design:customer-service-outlined','/service/customer/list','page',2,1),
  (303,30,'service.settings','客服设置','lucide:settings-2','/systemForm/Basics/service','page',3,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=30,
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`='page',
  `sort`=VALUES(`sort`),
  `status`=1;

INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` AS r
CROSS JOIN `qixi_crm_a_menu` AS m
WHERE r.code = 'platform'
  AND m.id IN (30,301,302,303);
