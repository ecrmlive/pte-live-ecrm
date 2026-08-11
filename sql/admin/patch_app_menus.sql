-- 平台「应用」侧栏：公众号目录 + 三级叶子
-- 执行后需重新登录平台后台以刷新菜单缓存。

SET NAMES utf8mb4;
USE `qixi_crm_admin`;

UPDATE `qixi_crm_a_menu`
SET `title`='应用',
    `icon`='ant-design:appstore-outlined',
    `route_path`='/app',
    `kind`='directory',
    `sort`=10,
    `parent_id`=0,
    `status`=1
WHERE `id`=130 OR `code`='app';

-- 二级：公众号目录；小程序保留为叶子（截图仅展开公众号）
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (131,130,'app.wechat','公众号','lucide:message-circle','/app/wechat','directory',1,1),
  (202,130,'app.routine','小程序','lucide:smartphone','/app/routine','page',2,1),
  (217,130,'app.mobile','App','lucide:smartphone','/app/mobile','page',3,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=130,
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

-- 推送配置已合并至 App 三端页。先迁移旧角色入口，再删除原菜单及其角色关联。
INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT `role_id`, 217
FROM `qixi_crm_a_role_menu`
WHERE `menu_id`=218;

DELETE FROM `qixi_crm_a_role_menu` WHERE `menu_id`=218;
DELETE FROM `qixi_crm_a_menu` WHERE `id`=218 OR `code`='app.push';

-- 三级：公众号（微信菜单 / 自动回复 / 图文管理；模板消息已下线 status=0）
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (204,131,'app.wechat_menus','微信菜单','lucide:menu','/app/wechat/menus','page',1,1),
  (203,131,'app.wechat_reply','自动回复','lucide:message-square-reply','/admin/app/wechat/reply','page',2,1),
  (206,131,'app.wechat_news','图文管理','lucide:newspaper','/app/wechat/newsCategory','page',3,1),
  (205,131,'app.wechat_template','微信模板消息','lucide:mail','/app/wechat/template','page',99,0)
ON DUPLICATE KEY UPDATE
  `parent_id`=131,
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`='page',
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`);

UPDATE `qixi_crm_a_menu` SET `parent_id`=131, `route_path`='app/wechat' WHERE `id`=20987 OR `code`='app.wechat.manage';
UPDATE `qixi_crm_a_menu` SET `parent_id`=204, `route_path`='app/wechat/menus' WHERE `id`=20995 OR `code`='app.wechat_menus.manage';
UPDATE `qixi_crm_a_menu` SET `parent_id`=203, `route_path`='admin/app/wechat/reply' WHERE `id`=20994 OR `code`='app.wechat_reply.manage';
UPDATE `qixi_crm_a_menu` SET `parent_id`=206, `route_path`='app/wechat/newsCategory' WHERE `id`=20997 OR `code`='app.wechat_news.manage';
UPDATE `qixi_crm_a_menu` SET `parent_id`=205, `route_path`='app/wechat/template', `status`=0 WHERE `id`=20996 OR `code`='app.wechat_template.manage';

INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (21122,217,'app.mobile.manage','维护 App 配置','','app/mobile','button',1,1),
  (21123,217,'app.push.manage','维护 App 推送配置','','app/mobile','button',2,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `title`=VALUES(`title`),
  `route_path`=VALUES(`route_path`),
  `kind`='button',
  `sort`=VALUES(`sort`),
  `status`=1;

INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` AS r
CROSS JOIN `qixi_crm_a_menu` AS m
WHERE r.code = 'platform'
  AND m.id IN (130,131,202,203,204,206,217,21122,21123);
