-- 平台「内容」侧栏：文章/社区二级目录 + 三级叶子；财务保持一级兄弟
-- 执行后需重新登录平台后台以刷新菜单缓存。

SET NAMES utf8mb4;
USE `qixi_crm_admin`;

UPDATE `qixi_crm_a_menu`
SET `title`='内容',
    `icon`='ant-design:read-outlined',
    `route_path`='/content',
    `kind`='directory',
    `sort`=8,
    `parent_id`=0,
    `status`=1
WHERE `id`=80 OR `code`='content';

-- 旧扁平/错挂：侧栏隐藏（素材迁到装修；公告可走消息管理）
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=990 WHERE `id`=81 OR `code`='content.notice';
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=991 WHERE `id`=83 OR `code`='content.attachment';

-- 二级目录
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (534,80,'content.article.dir','文章','lucide:newspaper','/cms','directory',1,1),
  (82,80,'content.community','社区','lucide:images','/community','directory',2,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=80,
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`='directory',
  `sort`=VALUES(`sort`),
  `status`=1;

-- 三级：文章
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (84,534,'content.article','文章管理','lucide:pen-line','/cms/article','page',1,1),
  (535,534,'content.article.category','文章分类','lucide:folder-tree','/cms/articleCategory','page',2,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=534,
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`='page',
  `sort`=VALUES(`sort`),
  `status`=1;

-- 三级：社区
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (85,82,'content.community.category','社区分类','lucide:tags','/community/category','page',1,1),
  (86,82,'content.community.topic','社区话题','lucide:hash','/community/topic','page',2,1),
  (87,82,'content.community.list','社区内容','lucide:images','/community/list','page',3,1),
  (88,82,'content.community.reply','社区评论','lucide:message-square','/community/reply','page',4,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=82,
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`='page',
  `sort`=VALUES(`sort`),
  `status`=1;

UPDATE `qixi_crm_a_menu` SET `parent_id`=84, `route_path`='cms/article' WHERE `id`=20910 OR `code`='content.article.manage';
UPDATE `qixi_crm_a_menu` SET `parent_id`=535, `route_path`='cms/articleCategory' WHERE `id`=20911 OR `code`='content.article_category.manage';
UPDATE `qixi_crm_a_menu` SET `parent_id`=87, `route_path`='community/list' WHERE `id` IN (20920,20921) OR `code` IN ('content.community.audit','content.community.delete');

INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` AS r
CROSS JOIN `qixi_crm_a_menu` AS m
WHERE r.code = 'platform'
  AND m.id IN (80,82,84,85,86,87,88,534,535);
