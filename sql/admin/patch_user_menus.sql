-- 平台「用户」侧栏对齐 CRMEB 嵌套结构（二级叶子/目录同列；三级挂文件夹下）
-- 执行后需重新登录平台后台以刷新菜单缓存。

SET NAMES utf8mb4;
USE `qixi_crm_admin`;

UPDATE `qixi_crm_a_menu`
SET `title`='用户',
    `icon`='ant-design:user-outlined',
    `route_path`='/user',
    `kind`='directory',
    `sort`=7,
    `parent_id`=0,
    `status`=1
WHERE `id`=70 OR `code`='user';

-- 运营辅助页：侧栏隐藏（路由/按钮保留）
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=990 WHERE `id`=78 OR `code`='user.assets_adjustment';
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=991 WHERE `id`=79 OR `code`='user.member_adjustment';
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=992 WHERE `id`=170 OR `code`='user.coupon_operation';
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=993 WHERE `id`=171 OR `code`='user.referrer_adjustment';
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=994 WHERE `id`=172 OR `code`='user.group_assignment';
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=995 WHERE `id`=173 OR `code`='user.label_assignment';
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=996 WHERE `id`=174 OR `code`='user.status_adjustment';
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=997 WHERE `id`=175 OR `code`='user.create';
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=998 WHERE `id`=176 OR `code`='user.profile_maintenance';
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=999 WHERE `id`=177 OR `code`='user.password_reset';
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=1000 WHERE `id`=178 OR `code`='user.promoter_assignment';
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=1001 WHERE `id`=179 OR `code`='user.notification';

-- 二级：叶子 + 目录
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (77,70,'user.list','用户列表','lucide:contact-round','/user/list','page',1,1),
  (72,70,'user.group','用户分组','lucide:users','/user/group','page',2,1),
  (71,70,'user.label','用户标签','lucide:tags','/user/label','page',3,1),
  (74,70,'user.feedback','用户反馈','lucide:message-circle-warning','/user/feedback','directory',4,1),
  (182,70,'user.search_record','搜索记录','lucide:search','/user/search-record','page',5,1),
  (530,70,'user.level.dir','用户等级','lucide:medal','/user/member','directory',6,1),
  (532,70,'user.setup','用户设置','lucide:settings-2','/user/setup_user','page',7,1),
  (73,70,'user.svip','付费会员','lucide:crown','/user/svip','directory',8,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

-- 三级：用户反馈
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (75,74,'user.feedback.list','反馈列表','lucide:message-circle','/user/feedback/list','page',1,1),
  (76,74,'user.feedback.category','反馈分类','lucide:folder-tree','/user/feedback/categories','page',2,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=74,
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`='page',
  `sort`=VALUES(`sort`),
  `status`=1;

-- 三级：用户等级
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (183,530,'user.member.level','等级管理','lucide:medal','/user/member/list','page',1,1),
  (531,530,'user.level.description','等级说明','lucide:file-text','/user/member/description','page',2,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=530,
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`='page',
  `sort`=VALUES(`sort`),
  `status`=1;

-- 三级：付费会员
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (180,73,'user.svip.plan','会员类型','lucide:badge-check','/user/member/type','page',1,1),
  (184,73,'user.svip.interest','会员权益','lucide:heart-handshake','/user/member/equity','page',2,1),
  (181,73,'user.svip.record','会员记录','lucide:scroll-text','/user/member/record','page',3,1),
  (533,73,'user.svip.agreement','会员协议','lucide:file-text','/user/member/vipAgreement','page',4,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=73,
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`='page',
  `sort`=VALUES(`sort`),
  `status`=1;

UPDATE `qixi_crm_a_menu` SET `parent_id`=532, `route_path`='user/setup_user' WHERE `id`=20998 OR `code`='user.setup.manage';
UPDATE `qixi_crm_a_menu` SET `parent_id`=183, `route_path`='user/member/list' WHERE `id`=20978 OR `code`='user.member.level.manage';
UPDATE `qixi_crm_a_menu` SET `parent_id`=73, `route_path`='user/svip' WHERE `id`=20972 OR `code`='user.svip.manage';

INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` AS r
CROSS JOIN `qixi_crm_a_menu` AS m
WHERE r.code = 'platform'
  AND m.id IN (70,71,72,73,74,75,76,77,180,181,182,183,184,530,531,532,533);
