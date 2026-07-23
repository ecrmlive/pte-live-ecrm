-- 全量导入 CRMEB MER v4.0 eb_system_menu → qixi_system_menu
-- 来源: install/crmeb_merchant.sql (+ backup/update_4_0.sql)
-- 变换: is_mer 0→1 / 1→2; is_menu 0→2 / 1→1; 与本仓自定义 menu_id 1–172 冲突的 CRMEB 行偏移 +20000
-- 不覆盖已有自定义菜单；角色 rules 不自动灌入全量按钮（见文件末说明）
USE `qixi_mergers`;

INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20033, 9218, '/9218/', '', '控制台', '/dashboard', '[]', 100, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20033);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20034, 20118, '/110/118/', '', '系统设置', '/systemForm/Basics/system_tabs', '[]', 99, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20034);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20035, 520, '/520/', '', '配置分类', '/config/classify', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20035);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20036, 520, '/520/', '', '配置管理', '/config/setting', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20036);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20038, 20110, '/110/', '', '权限管理', '/setting', '[]', 98, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20038);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20039, 20038, '/110/38/', '', '角色权限', '/setting/systemRole', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20039);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20040, 20038, '/110/38/', '', '管理员管理', '/setting/systemAdmin', '', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20040);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20041, 1665, '/1665/', '', '素材管理', '/config/picture', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20041);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20042, 9492, '/9492/', '', '商户管理', '/merchant', '[]', 20, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20042);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20043, 1284, '/9492/1284/', '', '店铺菜单', '/merchant/system', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20043);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20044, 6370, '/9492/6370/', '', '店铺列表', '/merchant/list', '[]', 100, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20044);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20047, 520, '/520/', '', '操作日志', '/setting/systemLog', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20047);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20048, 20038, '/110/38/', '', '菜单管理', '/setting/menu', '', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20048);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20049, 526, '/526/', '', '权限管理', '/setting', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20049);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20050, 20049, '/526/49/', '', '身份管理', '/setting/systemRole', '', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20050);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20051, 20049, '/526/49/', '', '管理员管理', '/setting/systemAdmin', '', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20051);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20052, 20049, '/526/49/', '', '操作日志', '/setting/systemLog', '', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20052);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20054, 1671, '1671/', '', '素材管理', '/config/picture', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20054);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20055, 0, '/', 'house', '首页', '/dashboard', '[]', 100, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20055);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20057, 521, '/520/521/', '', '组合数据', '/group/list', '', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20057);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20058, 519, '/519/', '', '公众号', '/app/wechat', '', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20058);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20059, 20058, '/519/58/', '', '微信菜单', '/app/wechat/menus', '', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20059);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20060, 9361, '/9361/', 's-management', '文章', '/cms', '[]', 96, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20060);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20061, 20060, '/9361/60/', '', '文章管理', '/cms/article', '', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20061);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20062, 20060, '/9361/60/', '', '文章分类', '/cms/articleCategory', '', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20062);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20072, 20111, '/110/118/1368/111/', '', '短信配置', '/systemForm/Basics/message', '[]', -1, 0, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20072);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20074, 526, '/526/', '', '店铺配置', '/systemForm/Basics/mer_base', '[]', 99, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20074);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20077, 20058, '/519/58/', '', '自动回复', '/admin/app/wechat/reply', '', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20077);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20079, 20077, '/519/58/77/', '', '微信关注回复', '/app/wechat/reply/follow/subscribe', '', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20079);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20080, 20077, '/519/58/77/', '', '关键字回复', '/app/wechat/reply/keyword', '', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20080);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20081, 20077, '/519/58/77/', '', '无效关键词回复', '/app/wechat/reply/index/default', '', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20081);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20082, 20058, '/519/58/', '', '图文管理', '/app/wechat/newsCategory', '', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20082);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20087, 0, '/', 's-goods', '商品', '/product', '', 99, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20087);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20088, 20087, '/87/', '', '商品分类', '/product/classify', '', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20088);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20092, 20087, '/87/', '', '品牌管理', '/product/brand', '', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20092);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20093, 20092, '/87/92/', '', '品牌分类', '/product/band/brandClassify', '', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20093);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20094, 20092, '/87/92/', '', '品牌列表', '/product/band/brandList', '', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20094);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20095, 0, '/', 'goods', '商品', '/product', '[]', 99, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20095);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20096, 20095, '/95/', '', '商品分类', '/product/classify', '[]', 8, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20096);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20099, 20095, '/95/', '', '商品规格', '/product/attr', '[]', 7, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20099);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20100, 6370, '/9492/6370/', '', '店铺分类', '/merchant/classify', '[]', 22, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20100);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20101, 0, '/', 'user-solid', '用户', '/user', '', 96, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20101);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20102, 20101, '/101/', '', '用户分组', '/user/group', '', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20102);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20103, 20101, '/101/', '', '用户列表', '/user/list', '[]', 1, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20103);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20104, 20101, '/101/', '', '用户标签', '/user/label', '', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20104);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20105, 20095, '/95/', '', '商品列表', '/product/list', '[]', 10, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20105);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20106, 0, '/', 'bell', '营销', '/marketing', '[]', 97, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20106);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20107, 20106, '/106/', '', '优惠券', '/marketing/coupon', '', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20107);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20110, 0, '/', 's-tools', '设置', '/settings', '[]', 91, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20110);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20111, 1368, '/110/118/1368/', '', '短信设置', '/sms', '[]', 0, 0, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20111);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20112, 20111, '/110/118/1368/111/', '', '短信账户', '/sms/user', '[]', 0, 0, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20112);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20113, 20111, '/110/118/1368/111/', '', '短信模板', '/sms/template', '', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20113);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20114, 20111, '/110/118/1368/111/', '', '申请记录', '/sms/applyList', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20114);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20115, 20107, '/106/107/', '', '优惠券列表', '/marketing/coupon/list', '[]', 1, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20115);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20116, 520, '/520/', '', '安全维护', '/maintain', '', 9, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20116);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20117, 20116, '/520/116/', '', '数据备份', '/maintain/dataBackup', '', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20117);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20118, 20110, '/110/', '', '系统设置', '/sys', '[]', 100, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20118);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20119, 5124, '/110/9360/5124/', '', '物流公司', '/freight/express', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20119);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20120, 9410, '9410/', '', '店员列表', '/config/service', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20120);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20121, 20087, '/87/', '', '评论管理', '/product/comment', '', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20121);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 20122, 20107, '/106/107/', '', '领取记录', '/marketing/coupon/user', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20122);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 466, 20101, '/101/', '', '用户反馈', '/feedback', '', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 466);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 467, 466, '/101/466/', '', '反馈分类', '/feedback/classify', '', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 467);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 468, 466, '/101/466/', '', '反馈列表', '/feedback/list', '', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 468);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 512, 0, '/', 'tickets', '订单', '/order', '[]', 99, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 512);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 513, 512, '/512/', '', '订单管理', '/order/list', '[]', 127, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 513);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 514, 0, '/', 's-promotion', '分销', '/promoter', '[]', 97, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 514);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 515, 0, '/', 's-data', '财务', '/accounts', '[]', 93, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 515);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 516, 537, '/515/537/', '', '提现管理', '/accounts/extract', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 516);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 517, 537, '/515/537/', '', '充值记录', '/accounts/bill', '', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 517);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 519, 0, '/', 's-grid', '应用', '/apploction', '[]', 93, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 519);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 520, 0, '/', 's-help', '维护', '/safe', '[]', 90, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 520);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 521, 520, '/520/', '', '开发配置', '/safe/exploit', '[]', 10, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 521);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 522, 514, '/514/', '', '分销员列表', '/promoter/user', '[]', 10, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 522);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 524, 519, '/519/', '', '小程序', '/app/routine', '', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 524);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 525, 0, '/', 'pie-chart', '财务', '/accounts', '[]', 97, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 525);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 526, 0, '/', 'setting', '设置', '/config', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 526);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 528, 512, '/512/', '', '退款订单', '/order/refund', '[]', 126, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 528);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 532, 20058, '/519/58/', '', '微信模板消息', '/app/wechat/template', '[]', 0, 0, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 532);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 537, 515, '/515/', '', '用户结算', '/accounts/record', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 537);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 538, 537, '/515/537/', '', '资金记录', '/accounts/capital', '', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 538);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 539, 20087, '/87/', '', '商品管理', '/product/examine', '[]', 100, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 539);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 540, 0, '/', 's-order', '订单', '/order', '[]', 98, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 540);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 541, 540, '/540/', '', '订单列表', '/order/list', '', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 541);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 542, 540, '/540/', '', '退款订单', '/order/refund', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 542);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 544, 512, '512/', '', '商品评价', '/product/reviews', '[]', 100, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 544);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 545, 524, '/519/524/', '', '小程序订阅消息', '/app/routine/template', '[]', 0, 0, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 545);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 546, 526, '/526/', '', '店铺信息', '/systemForm/modifyStoreInfo', '[]', 100, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 546);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 666, 20116, '/520/116/', '', '商业授权', '/setting/system/maintain/auth', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 666);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 667, 5126, '/719/5126/', '', '余额设置', '/systemForm/Basics/balance', '[]', 8, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 667);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 668, 8613, '9358/1656/8613/', '', '保存配置内容', 'merchantStorePrinterSetContent', '', 0, 0, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 668);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 670, 20110, '/110/', '', '应用配置', '/app_config', '', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 670);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 673, 670, '/110/670/', '', '公众号配置', '/systemForm/Basics/wechat', '', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 673);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 674, 670, '/110/670/', '', '小程序配置', '/systemForm/Basics/smallapp', '', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 674);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 677, 514, '/514/', '', '提现银行', '/group/config/76', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 677);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 684, 9360, '/110/9360/', '', '热门搜索', '/group/config/67', '[]', 95, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 684);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 685, 514, '/514/', '', '分销特权', '/group/config/75', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 685);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 686, 514, '/514/', '', '分销海报', '/group/config/68', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 686);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 687, 5126, '/719/5126/', '', '余额充值配置', '/group/config/69', '[]', -1, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 687);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 700, 1649, '1649/', '', '运费模板', '/config/freight/shippingTemplates', '[]', 7, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 700);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 715, 520, '/520/', '', '页面链接', '/safe/pageLinks', '[]', -1, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 715);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 719, 0, '/', 's-flag', '营销', '/marketing', '[]', 97, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 719);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 720, 719, '/719/', '', '商户优惠券', '/marketing/coupon', '[]', 95, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 720);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 721, 720, '/719/720/', '', '优惠券列表', '/marketing/coupon/list', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 721);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 731, 514, '/514/', '', '分销礼包', '/promoter/gift', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 731);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 734, 720, '/719/720/', '', '领取记录', '/marketing/coupon/user', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 734);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 778, 6370, '/9492/6370/', '', '店铺入驻申请', '/merchant/application', '[]', 8, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 778);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 779, 780, '/719/780/', '', '秒杀配置', '/marketing/seckill/seckillConfig', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 779);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 780, 719, '/719/', '', '秒杀', '/marketing/seckill', '[]', 85, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 780);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 781, 782, '/719/782/', '', '直播间管理', '/marketing/studio/list', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 781);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 782, 719, '/719/', '', '直播', '/marketing2', '[]', 80, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 782);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 783, 782, '/719/782/', '', '直播商品管理', '/marketing/broadcast/list', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 783);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 784, 540, '/540/', '', '核销记录', '/order/cancellation', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 784);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 785, 20106, '/106/', '', '直播', '/', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 785);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 786, 785, '/106/785/', '', '直播间管理', '/marketing/studio/list', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 786);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 787, 785, '/106/785/', '', '直播商品管理', '/marketing/broadcast/list', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 787);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 788, 20106, '/106/', '', '秒杀', '/marketing/seckill/list', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 788);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 789, 512, '/512/', '', '核销记录', '/order/cancellation', '[]', 125, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 789);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 790, 525, '/525/', '', '资金流水', '/accounts/capitalFlow', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 790);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 791, 515, '/515/', '', '资金流水', '/accounts/capitalFlow', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 791);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 794, 780, '/719/780/', '', '秒杀管理', '/marketing/seckill/list', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 794);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1022, 719, '/719/', '', '预售', '/marketing/presell', '[]', 70, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1022);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1023, 1022, '/719/1022/', '', '预售商品', '/marketing/presell/list', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1023);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1024, 1022, '/719/1022/', '', '预售协议', '/marketing/presell/agreement', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1024);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1025, 20106, '/106/', '', '预售', '/marketing/presell/list', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1025);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1027, 0, '/', 'user', '用户', '/user', '[]', 1, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1027);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1028, 1027, '/1027/', '', '标签管理', '/user/_label', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1028);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1029, 1028, '/1028/', '', '手动标签', '/user/label', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1029);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1030, 1028, '/1028/', '', '自动标签', '/user/maticlabel', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1030);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1051, 719, '/719/', '', '助力', '/assist', '[]', 60, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1051);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1095, 1051, '/719/1051/', '', '活动商品', '/marketing/assist/goods_list', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1095);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1096, 1051, '/719/1051/', '', '助力活动', '/marketing/assist/list', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1096);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1099, 20106, '/106/', '', '助力', '/assist', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1099);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1100, 1099, '/1099/', '', '助力商品', '/marketing/assist/list', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1100);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1101, 1099, '/1099/', '', '助力活动', '/marketing/assist/assist_set', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1101);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1102, 525, '/525/', '', '发票管理', '/order/invoice', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1102);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1103, 1027, '/1027/', '', '用户管理', '/user/list', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1103);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1119, 0, '/', 'notebook-2', '公告', '/station/notice', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1119);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1120, 5125, '/110/5125/', '', '公告管理', '/station/notice', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1120);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1132, 20107, '/107/', '', '发送记录', '/marketing/coupon/send', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1132);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1135, 719, '/719/', '', '拼团', '/marketing/combination', '[]', 50, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1135);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1136, 1135, '/719/1135/', '', '拼团设置', '/marketing/combination/combination_set', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1136);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1137, 1135, '/719/1135/', '', '拼团商品列表', '/marketing/combination/combination_goods', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1137);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1138, 1135, '/719/1135/', '', '拼团活动列表', '/marketing/combination/combination_list', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1138);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1139, 20106, '/106/', '', '拼团', '/marketing/combination', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1139);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1140, 1139, '/1139/', '', '拼团商品列表', '/marketing/combination/combination_goods', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1140);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1141, 1139, '/1139/', '', '拼团活动列表', '/marketing/combination/combination_list', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1141);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1177, 9220, '/515/9220/', '', '转账记录', '/accounts/transferRecord', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1177);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1178, 9220, '/515/9220/', '', '平台账单', '/accounts/statement', '[]', 10, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1178);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1179, 9220, '/515/9220/', '', '转账设置', '/accounts/settings', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1179);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1181, 525, '/525/', '', '转账记录', '/accounts/transManagement', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1181);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1182, 525, '/525/', '', '收款方式', '/accounts/payType', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1182);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1183, 525, '/525/', '', '账单管理', '/accounts/statement', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1183);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1244, 520, '/520/', '', '导出记录', '/group/exportList', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1244);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1245, 20087, '/87/', '', '保障服务', '/product/guarantee', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1245);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1246, 20095, '/95/', '', '服务模板', '/config/guarantee', '[]', 4, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1246);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1247, 20101, '/101/', '', '用户协议', '/user/agreement', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1247);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1284, 9492, '/9492/', '', '店铺设置', '/mer/store', '[]', 21, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1284);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1285, 20101, '/101/', '', '搜索记录', '/user/searchRecord', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1285);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1286, 1027, '/1027/', '', '搜索记录', '/user/searchRecord', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1286);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1289, 719, '/719/', '', '积分', '/marketing/integral', '[]', 40, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1289);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1290, 1289, '/719/1289/', '', '积分配置', '/marketing/integral/config', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1290);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1291, 1289, '/719/1289/', '', '积分日志', '/marketing/integral/log', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1291);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1293, 20106, '/106/', '', '积分', '/marketing/integral', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1293);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1294, 1293, '/1293/', '', '积分配置', '/marketing/integral/config', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1294);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1295, 1293, '/1293/', '', '积分日志', '/marketing/integral/log', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1295);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1296, 514, '/514/', '', '佣金说明', '/promoter/commission', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1296);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1298, 515, '/515/', '', '发票管理', '/accounts/accounts', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1298);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1299, 1298, '/515/1298/', '', '发票列表', '/accounts/receipt', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1299);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1300, 1298, '/515/1298/', '', '发票说明', '/accounts/invoiceDesc', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1300);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1302, 670, '/110/670/', '', 'APP配置', '/systemForm/Basics/wechat_open_app', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1302);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1303, 6370, '/9492/6370/', '', '店铺分账申请', '/merchant/applyments', '[]', 7, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1303);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1304, 525, '/525/', '', '申请分账商户', '/systemForm/applyments', '[]', -1, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1304);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1305, 9220, '/515/9220/', '', '分账管理', '/merchant/applyList', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1305);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1366, 525, '/525/', '', '分账管理', '/systemForm/applyList', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1366);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1368, 20118, '/110/118/', '', '一号通', '/serve', '[]', 89, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1368);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1369, 1368, '/110/118/1368/', '', '登陆入口', '/setting/sms/sms_config/index', '[]', 100, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1369);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1371, 1368, '/110/118/1368/', '', '服务配置', '/service/settings', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1371);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1372, 1368, '/110/118/1368/', '', '购买记录', '/service/purchase', '[]', 0, 0, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1372);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1373, 514, '/514/', '', '分销等级', 'brokerage', '[]', 1, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1373);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1374, 1373, '/514/1373/', '', '分销员等级', '/promoter/membership_level', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1374);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1375, 1373, '/514/1373/', '', '等级规则', '/promoter/distribution', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1375);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1376, 526, '526/', '', '一号通', '/one_setting', '[]', 1, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1376);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1377, 1376, '526/1376/', '', '配置管理', '/setting/sms/dumpConfig', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1377);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1378, 1376, '526/1376/', '', '平台一号通', '/setting/sms/sms_config/index', '[]', 1, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1378);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1379, 1368, '/110/118/1368/', '', '商户结余', '/service/balance_record', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1379);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1380, 1649, '1649/', '', '物流公司', '/config/freight/express', '[]', 6, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1380);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1381, 6370, '/9492/6370/', '', '店铺类型', '/merchant/type', '[]', 20, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1381);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1382, 1284, '/9492/1284/', '', '说明提示', '/merchant/type/description', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1382);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1468, 20095, '/95/', '', '商品标签', '/product/label', '[]', 5, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1468);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1469, 20087, '/87/', '', '商品标签', '/product/label', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1469);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1470, 719, '/719/', '', '专场列表', '/group/topic/94', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1470);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1471, 20106, '/106/', '', '专场列表', '/group/topic/95', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1471);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1508, 9047, '/101/9047/', '', '等级配置', '/systemForm/Basics/members', '[]', 10, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1508);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1509, 9047, '/101/9047/', '', '等级权益', '/user/member/interests', '[]', 90, 0, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1509);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1510, 9047, '/101/9047/', '', '等级管理', '/user/member/list', '[]', 100, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1510);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1511, 9047, '/101/9047/', '', '等级说明', '/user/member/description', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1511);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1538, 9361, '/9361/', 's-cooperation', '社区', '/community', '[]', 96, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1538);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1539, 1538, '/9361/1538/', '', '社区内容', '/community/list', '[]', 8, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1539);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1540, 1538, '/9361/1538/', '', '社区分类', '/community/category', '[]', 10, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1540);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1541, 1538, '/9361/1538/', '', '社区话题', '/community/topic', '[]', 9, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1541);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1542, 1538, '/9361/1538/', '', '社区配置', '/systemForm/Basics/community', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1542);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1543, 1538, '/9361/1538/', '', '社区评论', '/community/reply', '[]', 7, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1543);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1594, 785, '/785/', '', '直播助手', '/marketing/studio/assistant', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1594);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1595, 670, '/110/670/', '', '上传校验文件', '/app/wechat/file', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1595);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1596, 524, '/519/524/', '', '小程序下载', '/app/routine/download', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1596);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1597, 5125, '/110/5125/', '', '消息管理', '/setting/notification/index', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1597);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1617, 20116, '/520/116/', '', '缓存清除', '/maintain/cache', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1617);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1628, 670, '/110/670/', '', 'APP升级配置', '/systemForm/Basics/app_version', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1628);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1629, 719, '/719/', '', '优惠套餐', '/marketing/discounts/list', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1629);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1630, 20106, '/106/', '', '优惠套餐', '/marketing/discounts/list', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1630);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1647, 1284, '/9492/1284/', '', '店铺保证金', '/merchant/deposit_list', '[]', 70, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1647);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1648, 1653, '/110/9360/5124/1653/', '', '配送配置', '/systemForm/delivery', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1648);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1649, 526, '526/', '', '快递配送', '/city', '[]', 1, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1649);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1650, 9332, '9332/', '', '发货点管理', '/delivery/store_manage', '[]', 0, 0, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1650);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1651, 9332, '9332/', '', '配送记录', '/delivery/usage_record', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1651);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1652, 9332, '9332/', '', '充值记录', '/delivery/recharge_record', '[]', 0, 0, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1652);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1653, 5124, '/110/9360/5124/', '', '第三方送', '/delivery', '[]', 92, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1653);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1654, 1653, '/110/9360/5124/1653/', '', '充值记录', '/delivery/recharge_record', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1654);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1655, 1653, '/110/9360/5124/1653/', '', '配送记录', '/delivery/usage_record', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1655);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1656, 9358, '9358/', '', '小票打印', '/setting/printer/list', '[]', 6, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1656);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1657, 719, '/719/', '', '平台优惠券', '/marketing/platform_coupon', '[]', 99, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1657);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1658, 1657, '/719/1657/', '', '优惠券列表', '/marketing/platform_coupon/list', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1658);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1659, 1657, '/719/1657/', '', '领取记录', '/marketing/platform_coupon/couponRecord', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1659);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1661, 0, '/', 'headset', '员工', '/server', '[]', 1, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1661);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1662, 1657, '/719/1657/', '', '发送记录', '/marketing/platform_coupon/couponSend', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1662);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1663, 1657, '/719/1657/', '', '使用说明', '/marketing/platform_coupon/instructions', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1663);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1664, 1665, '/1665/', '', '主题风格', '/setting/theme_style', '[]', 90, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1664);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1665, 0, '/', 's-open', '装修', '/theme', '[]', 92, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1665);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1666, 1665, '/1665/', '', '页面装修', '/setting/diy/list', '[]', 80, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1666);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1667, 9360, '/110/9360/', '', '协议规则', '/setting/agreements', '[]', 95, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1667);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1668, 5054, '/5054/', '', '客服自动回复', '/systemForm/customer_keyword', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1668);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1669, 6372, '/1665/6372/', '', '平台页面链接', '/setting/diy/links/list', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1669);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1670, 6372, '/1665/6372/', '', '平台页面分类', '/setting/diy/plantform/category/list', '[]', 1, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1670);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1671, 0, '/', 'brush', '装修', '/devise/', '[]', 1, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1671);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1672, 1671, '/1671/', '', '装修', '/devise/diy/list', '[]', 100, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1672);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1673, 6372, '/1665/6372/', '', '商户页面分类', '/setting/diy/merchant/category/list', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1673);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 1674, 6372, '/1665/6372/', '', '商户页面链接', '/setting/diy/merLink/list', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1674);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 5054, 0, '/', 's-custom', '客服', '/service', '[]', 92, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 5054);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 5088, 5054, '/5054/', '', '客服列表', '/service/customer/list', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 5088);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 5119, 9360, '/110/9360/', '', '商城设置', '/systemForm/Basics/shop_tabs', '[]', 97, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 5119);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 5120, 9360, '/110/9360/', '', '支付设置', '/systemForm/Basics/pay_tabs', '[]', 96, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 5120);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 5121, 20118, '/110/118/', '', '接口配置', '/systemForm/Basics/extend_tabs', '[]', 93, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 5121);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 5122, 514, '/514/', '', '分销配置', '/systemForm/Basics/distribution_tabs', '[]', 9, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 5122);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 5123, 9410, '9410/', '', '自动回复', '/systemForm/customer_keyword', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 5123);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 5124, 9360, '/110/9360/', '', '配送配置', '/delivery_config', '[]', 94, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 5124);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 5125, 20110, '/110/', '', '消息管理', '/notice', '[]', 92, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 5125);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 5126, 719, '/719/', '', '余额充值', '/banlace', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 5126);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 5127, 5054, '/5054/', '', '客服设置', '/systemForm/Basics/service', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 5127);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 6370, 9492, '/9492/', '', '店铺管理', '/mer/mer', '[]', 30, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 6370);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 6372, 1665, '/1665/', '', '页面链接', '/setting/page', '[]', 30, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 6372);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7649, 1303, '/9492/6370/1303/', '', '权限', '/merchant/applyments', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7649);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7650, 1305, '/515/9220/1305/', '', '权限', '/merchant/applyList', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7650);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7651, 1305, '/515/9220/1305/', '', '附加权限', 'append_/merchant/applyList', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7651);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7652, 1179, '/515/9220/1179/', '', '权限', '/accounts/settings', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7652);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7653, 516, '/515/537/516/', '', '权限', '/accounts/extract', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7653);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7654, 516, '/515/537/516/', '', '附加权限', 'append_/accounts/extract', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7654);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7655, 1299, '/515/1298/1299/', '', '权限', '/accounts/receipt', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7655);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7656, 517, '/515/537/517/', '', '权限', '/accounts/bill', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7656);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7657, 538, '/515/537/538/', '', '权限', '/accounts/capital', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7657);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7658, 538, '/515/537/538/', '', '附加权限', 'append_/accounts/capital', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7658);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7659, 1178, '/515/9220/1178/', '', '权限', '/accounts/sysBill', '', 1, 0, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7659);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7660, 1178, '/515/9220/1178/', '', '附加权限', 'append_/accounts/sysBill', '', 1, 0, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7660);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7661, 791, '/515/791/', '', '权限', '/accounts/capitalFlow', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7661);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7662, 791, '/515/791/', '', '附加权限', 'append_/accounts/capitalFlow', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7662);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7663, 1177, '/515/9220/1177/', '', '权限', '/accounts/transferRecord', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7663);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7664, 1177, '/515/9220/1177/', '', '附加权限', 'append_/accounts/transferRecord', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7664);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7665, 20062, '/9361/60/62/', '', '权限', '/cms/articleCategory', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7665);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7666, 20061, '/9361/60/61/', '', '权限', '/cms/article', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7666);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7667, 20061, '/9361/60/61/', '', '附加权限', 'append_/cms/article', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7667);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7668, 20041, '/1665/41/', '', '权限', '/config/picture', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7668);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7669, 1539, '/9361/1538/1539/', '', '权限', '/community/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7669);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7670, 1540, '/9361/1538/1540/', '', '权限', '/community/category', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7670);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7671, 1541, '/9361/1538/1541/', '', '权限', '/community/topic', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7671);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7672, 1541, '/9361/1538/1541/', '', '附加权限', 'append_/community/topic', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7672);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7673, 1543, '/9361/1538/1543/', '', '权限', '/community/reply', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7673);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7674, 673, '/110/670/673/', '', '权限', '/systemForm/Basics/wechat', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7674);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7675, 20072, '/110/118/1368/111/72/', '', '权限', '/systemForm/Basics/message', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7675);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7676, 674, '/110/670/674/', '', '权限', '/systemForm/Basics/smallapp', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7676);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7677, 667, '/719/5126/667/', '', '权限', '/systemForm/Basics/balance', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7677);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7679, 1302, '/110/670/1302/', '', '权限', '/systemForm/Basics/wechat_open_app', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7679);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7680, 20035, '/520/35/', '', '权限', '/config/classify', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7680);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7681, 20036, '/520/36/', '', '权限', '/config/setting', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7681);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7682, 721, '/719/720/721/', '', '权限', '/marketing/coupon/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7682);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7683, 734, '/719/720/734/', '', '权限', '/marketing/coupon/user', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7683);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7684, 1658, '/719/1657/1658/', '', '权限', '/marketing/platform_coupon/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7684);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7685, 1658, '/719/1657/1658/', '', '附加权限', 'append_/marketing/platform_coupon/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7685);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7686, 1659, '/719/1657/1659/', '', '权限', '/marketing/platform_coupon/couponRecord', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7686);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7687, 20103, '/101/103/', '', '权限', '/user/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7687);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7688, 20103, '/101/103/', '', '附加权限', 'append_/user/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7688);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7689, 1662, '/719/1657/1662/', '', '权限', '/marketing/platform_coupon/couponSend', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7689);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7690, 1648, '/110/9360/5124/1653/1648/', '', '权限', '/systemForm/delivery', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7690);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7691, 1655, '/110/9360/5124/1653/1655/', '', '权限', '/delivery/usage_record', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7691);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7692, 1654, '/110/9360/5124/1653/1654/', '', '权限', '/delivery/recharge_record', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7692);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7693, 20119, '/110/9360/5124/119/', '', '权限', '/freight/express', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7693);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7694, 1664, '/1665/1664/', '', '权限', '/setting/theme_style', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7694);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7695, 1670, '/1665/6372/1670/', '', '权限', '/setting/diy/plantform/category/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7695);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7696, 1673, '/1665/6372/1673/', '', '权限', '/setting/diy/merchant/category/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7696);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7697, 1669, '/1665/6372/1669/', '', '权限', '/setting/diy/links/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7697);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7698, 1674, '/1665/6372/1674/', '', '权限', '/setting/diy/merLink/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7698);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7699, 1666, '/1665/1666/', '', '权限', '/setting/diy/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7699);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7700, 1666, '/1665/1666/', '', '附加权限', 'append_/setting/diy/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7700);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7701, 20057, '/520/521/57/', '', '权限', '/group/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7701);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7702, 20057, '/520/521/57/', '', '附加权限', 'append_/group/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7702);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7703, 684, '/110/9360/684/', '', '权限', '/group/config/67', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7703);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7704, 684, '/110/9360/684/', '', '附加权限', 'append_/group/config/67', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7704);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7705, 686, '/514/686/', '', '权限', '/group/config/68', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7705);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7706, 686, '/514/686/', '', '附加权限', 'append_/group/config/68', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7706);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7707, 687, '/719/5126/687/', '', '权限', '/group/config/69', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7707);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7708, 687, '/719/5126/687/', '', '附加权限', 'append_/group/config/69', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7708);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7709, 685, '/514/685/', '', '权限', '/group/config/75', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7709);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7710, 685, '/514/685/', '', '附加权限', 'append_/group/config/75', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7710);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7711, 677, '/514/677/', '', '权限', '/group/config/76', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7711);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7712, 677, '/514/677/', '', '附加权限', 'append_/group/config/76', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7712);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7739, 1470, '/719/1470/', '', '权限', '/group/topic/94', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7739);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7740, 1470, '/719/1470/', '', '附加权限', 'append_/group/topic/94', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7740);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7741, 1290, '/719/1289/1290/', '', '权限', '/marketing/integral/config', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7741);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7742, 1291, '/719/1289/1291/', '', '权限', '/marketing/integral/log', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7742);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7743, 1291, '/719/1289/1291/', '', '附加权限', 'append_/marketing/integral/log', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7743);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7744, 1023, '/719/1022/1023/', '', '权限', '/marketing/presell/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7744);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7745, 1095, '/719/1051/1095/', '', '权限', '/marketing/assist/goods_list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7745);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7746, 1096, '/719/1051/1096/', '', '权限', '/marketing/assist/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7746);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7747, 1137, '/719/1135/1137/', '', '权限', '/marketing/combination/combination_goods', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7747);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7748, 1138, '/719/1135/1138/', '', '权限', '/marketing/combination/combination_list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7748);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7749, 1136, '/719/1135/1136/', '', '权限', '/marketing/combination/combination_set', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7749);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7750, 781, '/719/782/781/', '', '权限', '/marketing/studio/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7750);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7751, 783, '/719/782/783/', '', '权限', '/marketing/broadcast/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7751);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7752, 779, '/719/780/779/', '', '权限', '/marketing/seckill/seckillConfig', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7752);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7753, 779, '/719/780/779/', '', '附加权限', 'append_/marketing/seckill/seckillConfig', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7753);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7754, 794, '/719/780/794/', '', '权限', '/marketing/seckill/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7754);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7755, 1510, '/101/9047/1510/', '', '权限', '/user/member/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7755);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7756, 1510, '/101/9047/1510/', '', '附加权限', 'append_/user/member/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7756);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7757, 1509, '/101/9047/1509/', '', '权限', '/user/member/interests', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7757);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7758, 1509, '/101/9047/1509/', '', '附加权限', 'append_/user/member/interests', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7758);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7759, 20048, '/110/38/48/', '', '权限', '/setting/menu', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7759);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7760, 20043, '/9492/1284/43/', '', '权限', '/merchant/system', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7760);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7761, 20100, '/9492/6370/100/', '', '权限', '/merchant/classify', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7761);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7762, 778, '/9492/6370/778/', '', '权限', '/merchant/application', '', 1, 0, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7762);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7763, 20044, '/9492/6370/44/', '', '权限', '/merchant/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7763);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7764, 20044, '/9492/6370/44/', '', '附加权限', 'append_/merchant/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7764);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7765, 1381, '/9492/6370/1381/', '', '权限', '/merchant/type', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7765);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7766, 1647, '/9492/1284/1647/', '', '权限', '/merchant/deposit_list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7766);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7767, 1120, '/110/5125/1120/', '', '权限', '/station/notice', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7767);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7768, 1597, '/110/5125/1597/', '', '权限', '/setting/notification/index', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7768);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7769, 541, '/540/541/', '', '权限', '/order/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7769);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7770, 541, '/540/541/', '', '附加权限', 'append_/order/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7770);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7771, 784, '/540/784/', '', '权限', '/order/cancellation', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7771);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7772, 542, '/540/542/', '', '权限', '/order/refund', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7772);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7773, 542, '/540/542/', '', '附加权限', 'append_/order/refund', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7773);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7774, 20088, '/87/88/', '', '权限', '/product/classify', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7774);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7775, 20088, '/87/88/', '', '附加权限', 'append_/product/classify', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7775);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7776, 20093, '/87/92/93/', '', '权限', '/product/band/brandClassify', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7776);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7777, 20094, '/87/92/94/', '', '权限', '/product/band/brandList', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7777);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7778, 539, '/87/539/', '', '权限', '/product/examine', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7778);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7779, 20121, '/87/121/', '', '权限', '/product/comment', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7779);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7780, 20121, '/87/121/', '', '附加权限', 'append_/product/comment', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7780);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7781, 1245, '/87/1245/', '', '权限', '/product/guarantee', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7781);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7782, 1245, '/87/1245/', '', '附加权限', 'append_/product/guarantee', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7782);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7783, 1469, '/87/1469/', '', '权限', '/product/label', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7783);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7784, 1629, '/719/1629/', '', '权限', '/marketing/discounts/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7784);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7785, 1629, '/719/1629/', '', '附加权限', 'append_/marketing/discounts/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7785);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7786, 20039, '/110/38/39/', '', '权限', '/setting/systemRole', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7786);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7787, 20040, '/110/38/40/', '', '权限', '/setting/systemAdmin', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7787);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7788, 20040, '/110/38/40/', '', '附加权限', 'append_/setting/systemAdmin', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7788);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7789, 20047, '/520/47/', '', '权限', '/setting/systemLog', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7789);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7790, 0, '/', '', '权限', 'self', '', 1, 0, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7790);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7791, 5088, '/5054/5088/', '', '权限', '/service/customer/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7791);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7792, 1668, '/5054/1668/', '', '权限', '/systemForm/customer_keyword', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7792);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7793, 522, '/514/522/', '', '权限', '/promoter/user', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7793);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7794, 1374, '/514/1373/1374/', '', '权限', '/promoter/membership_level', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7794);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7795, 1374, '/514/1373/1374/', '', '附加权限', 'append_/promoter/membership_level', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7795);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7796, 731, '/514/731/', '', '权限', '/promoter/gift', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7796);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7797, 5122, '/514/5122/', '', '权限', '/systemForm/Basics/distribution_tabs', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7797);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7798, 1244, '/520/1244/', '', '权限', '/group/exportList', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7798);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7799, 20033, '/9218/33/', '', '权限', '/dashboard', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7799);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7800, 20117, '/520/116/117/', '', '权限', '/maintain/dataBackup', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7800);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7801, 1617, '/520/116/1617/', '', '权限', '/maintain/cache', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7801);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7802, 8613, '9358/1656/8613/', '', '获取配置内容', 'merchantStorePrinterGetContent', '', 1, 0, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7802);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7803, 1667, '/110/9360/1667/', '', '权限', '/setting/agreements', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7803);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7804, 1024, '/719/1022/1024/', '', '权限', '/marketing/presell/agreement', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7804);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7805, 1296, '/514/1296/', '', '权限', '/promoter/commission', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7805);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7806, 1375, '/514/1373/1375/', '', '权限', '/promoter/distribution', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7806);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7807, 1663, '/719/1657/1663/', '', '权限', '/marketing/platform_coupon/instructions', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7807);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7808, 1247, '/101/1247/', '', '权限', '/user/agreement', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7808);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7809, 1511, '/101/9047/1511/', '', '权限', '/user/member/description', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7809);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7810, 1382, '/9492/1284/1382/', '', '权限', '/merchant/type/description', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7810);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7811, 1300, '/515/1298/1300/', '', '权限', '/accounts/invoiceDesc', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7811);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7812, 20104, '/101/104/', '', '权限', '/user/label', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7812);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7813, 1285, '/101/1285/', '', '权限', '/user/searchRecord', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7813);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7814, 1285, '/101/1285/', '', '附加权限', 'append_/user/searchRecord', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7814);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7815, 20102, '/101/102/', '', '权限', '/user/group', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7815);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7816, 467, '/101/466/467/', '', '权限', '/feedback/classify', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7816);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7817, 468, '/101/466/468/', '', '权限', '/feedback/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7817);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7818, 1595, '/110/670/1595/', '', '权限', '/app/wechat/file', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7818);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7819, 1596, '/519/524/1596/', '', '权限', '/app/routine/download', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7819);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7820, 20059, '/519/58/59/', '', '权限', '/app/wechat/menus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7820);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7821, 20077, '/519/58/77/', '', '权限', '/admin/app/wechat/reply', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7821);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7822, 20082, '/519/58/82/', '', '权限', '/app/wechat/newsCategory', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7822);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7823, 20082, '/519/58/82/', '', '附加权限', 'append_/app/wechat/newsCategory', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7823);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7824, 532, '/519/58/532/', '', '权限', '/app/wechat/template', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7824);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7825, 545, '/519/524/545/', '', '权限', '/app/routine/template', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7825);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7826, 1369, '/110/118/1368/1369/', '', '权限', '/setting/sms/sms_config/index', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7826);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7827, 1372, '/110/118/1368/1372/', '', '权限', '/service/purchase', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7827);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7828, 1379, '/110/118/1368/1379/', '', '权限', '/service/balance_record', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7828);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7829, 1371, '/110/118/1368/1371/', '', '权限', '/service/settings', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7829);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7830, 20113, '/110/118/1368/111/113/', '', '权限', '/sms/template', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7830);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7831, 20114, '/110/118/1368/111/114/', '', '权限', '/sms/applyList', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7831);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7832, 7649, '/9492/6370/1303/7649/', '', '分账商户申请列表', 'systemMerchantApplymentsLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7832);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7833, 7649, '/9492/6370/1303/7649/', '', '分账商户申请详情', 'systemMerchantApplymentsDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7833);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7834, 7649, '/9492/6370/1303/7649/', '', '分账商户申请审核', 'systemMerchantApplymentsStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7834);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7835, 7649, '/9492/6370/1303/7649/', '', '分账商户申请备注', 'systemMerchantApplymentsMarrkSave', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7835);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7836, 7650, '/515/9220/1305/7650/', '', '列表', 'systemOrderProfitsharingLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7836);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7837, 7650, '/515/9220/1305/7650/', '', '重新分账', 'systemOrderProfitsharingAgain', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7837);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7838, 7650, '/515/9220/1305/7650/', '', '导出', 'systemOrderProfitsharingExport', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7838);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7839, 7651, '/515/9220/1305/7651/', '', '导出列表', 'systemStoreExcelLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7839);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7840, 7651, '/515/9220/1305/7651/', '', '导出下载', 'systemStoreExcelDownload', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7840);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7841, 7652, '/515/9220/1179/7652/', '', '配置信息', 'systemOrderProfitsharingGetConfig', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7841);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7842, 7652, '/515/9220/1179/7652/', '', '配置保存', 'systemOrderProfitsharingSetConfig', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7842);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7843, 7653, '/515/537/516/7653/', '', '申请列表', 'systemUserExtractLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7843);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7844, 7653, '/515/537/516/7653/', '', '审核', 'systemUserExtractSwitchStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7844);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7845, 7653, '/515/537/516/7653/', '', '导出', 'systemUserExtractExport', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7845);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7846, 7654, '/515/537/516/7654/', '', '导出列表', 'systemStoreExcelLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7846);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7847, 7654, '/515/537/516/7654/', '', '导出下载', 'systemStoreExcelDownload', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7847);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7848, 7655, '/515/1298/1299/7655/', '', '列表', 'systemOrderReceiptList', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7848);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7849, 7655, '/515/1298/1299/7655/', '', '详情', 'systemOrderReceiptDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7849);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7850, 7656, '/515/537/517/7656/', '', '列表', 'systemUserRechargeList', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7850);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7851, 7656, '/515/537/517/7656/', '', '统计', 'systemUserRechargeTotal', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7851);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7852, 7657, '/515/537/538/7657/', '', '列表', 'systemUserBillList', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7852);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7853, 7657, '/515/537/538/7657/', '', '导出', 'systemUserBillExport', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7853);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7854, 7658, '/515/537/538/7658/', '', '导出列表', 'systemStoreExcelLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7854);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7855, 7658, '/515/537/538/7658/', '', '导出下载', 'systemStoreExcelDownload', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7855);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7856, 7659, '/515/9220/1178/7659/', '', '列表', 'systemFinancialRecordLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7856);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7857, 7659, '/515/9220/1178/7659/', '', '统计', 'systemFinancialRecordTitle', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7857);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7858, 7659, '/515/9220/1178/7659/', '', '详情', 'systemFinancialRecordDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7858);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7859, 7659, '/515/9220/1178/7659/', '', '导出', 'systemFinancialRecordDetailExport', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7859);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7860, 7660, '/515/9220/1178/7660/', '', '导出列表', 'systemStoreExcelLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7860);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7861, 7660, '/515/9220/1178/7660/', '', '导出下载', 'systemStoreExcelDownload', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7861);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7862, 7661, '/515/791/7661/', '', '列表', 'systemFinancialRecordList', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7862);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7863, 7661, '/515/791/7661/', '', '导出', 'systemFinancialRecordExport', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7863);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7864, 7661, '/515/791/7661/', '', '统计', 'systemFinancialCount', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7864);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7865, 7662, '/515/791/7662/', '', '导出列表', 'systemStoreExcelLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7865);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7866, 7662, '/515/791/7662/', '', '导出下载', 'systemStoreExcelDownload', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7866);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7867, 7663, '/515/9220/1177/7663/', '', '列表', 'systemFinancialList', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7867);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7868, 7663, '/515/9220/1177/7663/', '', '详情', 'systemFinancialDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7868);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7869, 7663, '/515/9220/1177/7663/', '', '编辑', 'systemFinancialUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7869);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7870, 7663, '/515/9220/1177/7663/', '', '修改状态', 'systemFinancialSwitchStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7870);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7871, 7663, '/515/9220/1177/7663/', '', '备注', 'systemFinancialMark', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7871);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7872, 7663, '/515/9220/1177/7663/', '', '统计', 'systemFinancialTitle', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7872);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7873, 7663, '/515/9220/1177/7663/', '', '导出', 'systemFinancialExport', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7873);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7874, 7664, '/515/9220/1177/7664/', '', '导出列表', 'systemStoreExcelLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7874);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7875, 7664, '/515/9220/1177/7664/', '', '导出下载', 'systemStoreExcelDownload', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7875);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7876, 7665, '/9361/60/62/7665/', '', '文章分类列表', 'systemArticleCategoryLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7876);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7877, 7665, '/9361/60/62/7665/', '', '文章分类添加', 'systemArticleCategoryCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7877);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7878, 7665, '/9361/60/62/7665/', '', '文章分类编辑', 'systemArticleCategoryUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7878);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7879, 7665, '/9361/60/62/7665/', '', '文章分类修改状态', 'systemArticleCategoryStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7879);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7880, 7665, '/9361/60/62/7665/', '', '文章分类删除', 'systemArticleCategoryDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7880);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7881, 7665, '/9361/60/62/7665/', '', '文章分类详情', 'systemArticleCategoryDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7881);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7882, 7666, '/9361/60/61/7666/', '', '文章列表', 'systemArticlArticleLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7882);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7883, 7666, '/9361/60/61/7666/', '', '文章添加', 'systemArticleArticleCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7883);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7884, 7666, '/9361/60/61/7666/', '', '文章编辑', 'systemArticArticleleUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7884);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7885, 7666, '/9361/60/61/7666/', '', '文章删除', 'systemArticArticleleDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7885);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7886, 7666, '/9361/60/61/7666/', '', '文章详情', 'systemArticArticleleDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7886);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7887, 7666, '/9361/60/61/7666/', '', '文章修改状态', 'systemArticlArticlStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7887);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7888, 7667, '/9361/60/61/7667/', '', '上传图片', 'uploadImage', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7888);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7889, 7667, '/9361/60/61/7667/', '', '素材列表', 'systemAttachmentLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7889);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7890, 7668, '/1665/41/7668/', '', '素材分类列表', 'systemAttachmentCategoryGetFormatList', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7890);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7891, 7668, '/1665/41/7668/', '', '素材分类添加', 'systemAttachmentCategoryCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7891);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7892, 7668, '/1665/41/7668/', '', '素材编辑', 'systemAttachmentCategoryUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7892);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7893, 7668, '/1665/41/7668/', '', '素材删除', 'systemAttachmentCategoryDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7893);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7894, 7668, '/1665/41/7668/', '', '素材列表', 'systemAttachmentLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7894);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7895, 7668, '/1665/41/7668/', '', '素材删除', 'systemAttachmentDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7895);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7896, 7668, '/1665/41/7668/', '', '批量移动', 'systemAttachmentBatchChangeCategory', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7896);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7897, 7668, '/1665/41/7668/', '', '素材编辑', 'systemAttachmentUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7897);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7898, 7669, '/9361/1538/1539/7669/', '', '文章列表', 'systemCommunityLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7898);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7899, 7669, '/9361/1538/1539/7669/', '', '文章详情', 'systemCommunityDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7899);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7900, 7669, '/9361/1538/1539/7669/', '', '文章编辑', 'systemCommunityUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7900);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7901, 7669, '/9361/1538/1539/7669/', '', '文章删除', 'systemCommunityDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7901);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7902, 7669, '/9361/1538/1539/7669/', '', '修改状态', 'systemCommunityStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7902);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7903, 7669, '/9361/1538/1539/7669/', '', '文章详情', 'systemCommunityShow', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7903);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7904, 7670, '/9361/1538/1540/7670/', '', '社区分类状态', 'systemCommunityCategoryLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7904);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7905, 7670, '/9361/1538/1540/7670/', '', '社区分类添加', 'systemCommunityCategoryCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7905);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7906, 7670, '/9361/1538/1540/7670/', '', '社区分类编辑', 'systemCommunityCategoryUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7906);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7907, 7670, '/9361/1538/1540/7670/', '', '社区分类详情', 'systemCommunityCategoryDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7907);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7908, 7670, '/9361/1538/1540/7670/', '', '社区分类删除', 'systemCommunityCategoryDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7908);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7909, 7670, '/9361/1538/1540/7670/', '', '社区分类修改状态', 'systemCommunityCategoryStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7909);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7910, 7671, '/9361/1538/1541/7671/', '', '社区话题', 'systemCommunityTopicLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7910);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7911, 7671, '/9361/1538/1541/7671/', '', '社区话题添加', 'systemCommunityTopicCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7911);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7912, 7671, '/9361/1538/1541/7671/', '', '社区话题编辑', 'systemCommunityTopicUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7912);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7913, 7671, '/9361/1538/1541/7671/', '', '社区话题详情 ', 'systemCommunityTopicDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7913);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7914, 7671, '/9361/1538/1541/7671/', '', '社区话题删除', 'systemCommunityTopicDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7914);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7915, 7671, '/9361/1538/1541/7671/', '', '社区话题修改状态', 'systemCommunityTopicStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7915);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7916, 7671, '/9361/1538/1541/7671/', '', '社区话题推荐', 'systemCommunityTopicHot', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7916);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7917, 7672, '/9361/1538/1541/7672/', '', '上传图片', 'uploadImage', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7917);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7918, 7672, '/9361/1538/1541/7672/', '', '图片列表', 'systemAttachmentLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7918);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7919, 7673, '/9361/1538/1543/7673/', '', '社区评论列表', 'systemCommunityReplyLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7919);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7920, 7673, '/9361/1538/1543/7673/', '', '社区评论删除', 'systemCommunityReplyDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7920);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7921, 7673, '/9361/1538/1543/7673/', '', '社区评论审核', 'systemCommunityReplyStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7921);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7922, 7674, '/110/670/673/7674/', '', '编辑配置信息', 'configSave', 'wechat', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7922);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7923, 7675, '/110/118/1368/111/72/7675/', '', '编辑配置信息', 'configSave', 'message', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7923);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7924, 7676, '/110/670/674/7676/', '', '编辑配置信息', 'configSave', 'smallapp', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7924);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7925, 7677, '/719/5126/667/7677/', '', '编辑配置信息', 'configSave', 'balance', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7925);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7927, 7679, '/110/670/1302/7679/', '', '编辑配置信息', 'configSave', 'wechat_open_app', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7927);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7928, 7680, '/520/35/7680/', '', '配置分类添加', 'configClassifyCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7928);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7929, 7680, '/520/35/7680/', '', '配置分类删除', 'configClassifyDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7929);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7930, 7680, '/520/35/7680/', '', '配置分类编辑', 'configClassifyUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7930);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7931, 7680, '/520/35/7680/', '', '配置分类修改状态', 'configClassifySwitchStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7931);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7932, 7680, '/520/35/7680/', '', '配置分类列表', 'configClassifyLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7932);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7933, 7681, '/520/36/7681/', '', '配置添加', 'configSettingCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7933);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7934, 7681, '/520/36/7681/', '', '配置编辑', 'configSettingUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7934);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7935, 7681, '/520/36/7681/', '', '配置修改状态', 'configSettingSwitchStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7935);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7936, 7681, '/520/36/7681/', '', '配置列表', 'configSettingLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7936);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7937, 7681, '/520/36/7681/', '', '配置删除', 'configSettingDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7937);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7938, 7682, '/719/720/721/7682/', '', '列表', 'systemStoreCouponLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7938);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7939, 7682, '/719/720/721/7682/', '', '详情', 'systemCouponDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7939);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7940, 7682, '/719/720/721/7682/', '', '商品列表', 'systemCouponProduct', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7940);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7941, 7683, '/719/720/734/7683/', '', '使用记录', 'systemCouponIssue', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7941);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7942, 7684, '/719/1657/1658/7684/', '', '添加', 'systemCouponCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7942);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7943, 7684, '/719/1657/1658/7684/', '', '编辑', 'systemCouponUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7943);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7944, 7684, '/719/1657/1658/7684/', '', '删除', 'systemCouponDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7944);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7945, 7684, '/719/1657/1658/7684/', '', '修改状态', 'systemCouponStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7945);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7946, 7684, '/719/1657/1658/7684/', '', '列表', 'systemCouponList', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7946);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7947, 7684, '/719/1657/1658/7684/', '', '详情', 'systemCouponShow', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7947);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7948, 7684, '/719/1657/1658/7684/', '', '详情关联列表', 'systemCouponShowLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7948);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7949, 7685, '/719/1657/1658/7685/', '', '上传图片', 'uploadImage', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7949);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7950, 7685, '/719/1657/1658/7685/', '', '图片列表', 'systemAttachmentLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7950);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7951, 7686, '/719/1657/1659/7686/', '', '使用记录', 'systemCouponIssue', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7951);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7952, 7687, '/101/103/7687/', '', '发送优惠券', 'systemCouponSend', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7952);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7953, 7687, '/101/103/7687/', '', '用户列表', 'systemUserLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7953);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7954, 7687, '/101/103/7687/', '', '用户编辑', 'systemUserUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7954);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7955, 7687, '/101/103/7687/', '', '用户修改余额', 'systemUserChangeNowMoney', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7955);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7956, 7687, '/101/103/7687/', '', '用户修改积分', 'systemUserChangeIntegral', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7956);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7957, 7687, '/101/103/7687/', '', '用户发送图文', 'systemWechatUserSendNews', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7957);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7958, 7687, '/101/103/7687/', '', '用户详情', 'systemUserDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7958);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7959, 7687, '/101/103/7687/', '', '用户消费记录', 'systemUserOrder', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7959);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7960, 7687, '/101/103/7687/', '', '用户持有优惠券', 'systemUserCoupon', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7960);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7961, 7687, '/101/103/7687/', '', '用户余额变动列表', 'systemUserBill', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7961);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7962, 7687, '/101/103/7687/', '', '推荐人修改记录', 'systemUserSpreadLog', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7962);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7963, 7687, '/101/103/7687/', '', '修改推荐人', 'systemUserSpreadChange', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7963);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7964, 7687, '/101/103/7687/', '', '用户修改会员等级', 'systemUserMemberSave', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7964);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7965, 7687, '/101/103/7687/', '', '用户添加', 'systemUserCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7965);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7966, 7687, '/101/103/7687/', '', '用户修改密码', 'systemUserChangePassword', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7966);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7967, 7687, '/101/103/7687/', '', '用户分组编辑', 'systemUserChangeGroup', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7967);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7968, 7687, '/101/103/7687/', '', '用户分组批量编辑', 'systemUserBatchChangeGroup', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7968);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7969, 7687, '/101/103/7687/', '', '用户标签编辑', 'systemUserChangeLabel', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7969);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7970, 7687, '/101/103/7687/', '', '用户标签批量编辑', 'systemUserBatchChangeLabel', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7970);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7971, 7688, '/101/103/7688/', '', '优惠券列表', 'systemCouponList', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7971);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7972, 7688, '/101/103/7688/', '', '上传图片', 'uploadImage', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7972);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7973, 7688, '/101/103/7688/', '', '图片列表', 'systemAttachmentLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7973);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7974, 7689, '/719/1657/1662/7689/', '', '发送记录', 'systemCouponSendLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7974);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7975, 7690, '/110/9360/5124/1653/1648/7690/', '', '编辑配置', 'systemDeliveryConfigSave', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7975);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7976, 7691, '/110/9360/5124/1653/1655/7691/', '', '配送记录', 'systemDeliveryOrderLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7976);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7977, 7691, '/110/9360/5124/1653/1655/7691/', '', '配送详情', 'systemDeliveryOrderDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7977);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7978, 7692, '/110/9360/5124/1653/1654/7692/', '', '充值记录', 'systemDeliveryStationPaayyLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7978);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7979, 7692, '/110/9360/5124/1653/1654/7692/', '', '统计', 'systemDeliveryOrderTitle', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7979);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7980, 7692, '/110/9360/5124/1653/1654/7692/', '', '余额', 'systemDeliveryStationGetBalance', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7980);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7981, 7693, '/110/9360/5124/119/7693/', '', '列表', 'systemExpressLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7981);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7982, 7693, '/110/9360/5124/119/7693/', '', '修改状态', 'systemExpressSwitchStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7982);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7983, 7693, '/110/9360/5124/119/7693/', '', '编辑', 'systemExpressUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7983);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7984, 7693, '/110/9360/5124/119/7693/', '', '删除', 'systemExpressDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7984);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7985, 7693, '/110/9360/5124/119/7693/', '', '同步', 'systemExpressSync', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7985);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7986, 7693, '/110/9360/5124/119/7693/', '', '列表', 'systemServeExportLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7986);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7987, 7694, '/1665/1664/7694/', '', '一键换色保存', 'systemSetChangeColor', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7987);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7988, 7695, '/1665/6372/1670/7695/', '', '列表 ', 'systemDiyPageCategroyLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7988);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7989, 7695, '/1665/6372/1670/7695/', '', '添加', 'systemDiyPageCategroyCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7989);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7990, 7695, '/1665/6372/1670/7695/', '', '编辑', 'systemDiyPageCategroyUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7990);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7991, 7695, '/1665/6372/1670/7695/', '', '编辑状态', 'systemDiyPageCategroyStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7991);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7992, 7695, '/1665/6372/1670/7695/', '', '删除', 'systemDiyPageCategroyDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7992);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7993, 7696, '/1665/6372/1673/7696/', '', '列表 ', 'systemDiyPageMerCategroyLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7993);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7994, 7696, '/1665/6372/1673/7696/', '', '添加', 'systemDiyPageMerCategroyCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7994);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7995, 7696, '/1665/6372/1673/7696/', '', '编辑', 'systemDiyPageMerCategroyUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7995);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7996, 7696, '/1665/6372/1673/7696/', '', '编辑状态', 'systemDiyPageMerCategroyStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7996);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7997, 7696, '/1665/6372/1673/7696/', '', '删除', 'systemDiyPageMerCategroyDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7997);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7998, 7697, '/1665/6372/1669/7697/', '', '列表', 'systemDiyPageLinkLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7998);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 7999, 7697, '/1665/6372/1669/7697/', '', '添加', 'systemDiyPageLinkCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 7999);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8000, 7697, '/1665/6372/1669/7697/', '', '编辑', 'systemDiyPageLinkUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8000);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8001, 7697, '/1665/6372/1669/7697/', '', '删除', 'systemDiyPageLinkDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8001);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8002, 7698, '/1665/6372/1674/7698/', '', '列表', 'systemDiyPageLinkMerLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8002);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8003, 7698, '/1665/6372/1674/7698/', '', '添加', 'systemDiyPageLinkMerCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8003);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8004, 7698, '/1665/6372/1674/7698/', '', '编辑', 'systemDiyPageLinkMerUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8004);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8005, 7698, '/1665/6372/1674/7698/', '', '删除', 'systemDiyPageLinkMerDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8005);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8006, 7699, '/1665/1666/7699/', '', '列表 ', 'systemDiyLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8006);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8007, 7699, '/1665/1666/7699/', '', '添加/编辑', 'systemDiyCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8007);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8008, 7699, '/1665/1666/7699/', '', '使用模板', 'systemDiyStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8008);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8009, 7699, '/1665/1666/7699/', '', '设置默认', 'systemDiySetDefault', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8009);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8010, 7699, '/1665/1666/7699/', '', '重置', 'systemDiyRecovery', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8010);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8011, 7699, '/1665/1666/7699/', '', '删除', 'systemDiyDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8011);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8012, 7699, '/1665/1666/7699/', '', '商品列表', 'systemDiyProductLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8012);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8013, 7699, '/1665/1666/7699/', '', '复制', 'systemDiyCopy', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8013);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8014, 7699, '/1665/1666/7699/', '', '个人中心装修', 'systemVisualUserInfo', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8014);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8015, 7699, '/1665/1666/7699/', '', '店铺街装修', 'systemVisualStoreStreet', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8015);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8016, 7700, '/1665/1666/7700/', '', '上传图片', 'uploadImage', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8016);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8017, 7700, '/1665/1666/7700/', '', '图片列表', 'systemAttachmentLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8017);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8018, 7701, '/520/521/57/7701/', '', '组合数据配置列表', 'groupLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8018);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8019, 7701, '/520/521/57/7701/', '', '组合数据配置添加', 'groupCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8019);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8020, 7701, '/520/521/57/7701/', '', '组合数据配置编辑', 'groupUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8020);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8021, 7701, '/520/521/57/7701/', '', '详情', 'groupDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8021);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8022, 7701, '/520/521/57/7701/', '', '列表', 'groupDataLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8022);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8023, 7701, '/520/521/57/7701/', '', '添加', 'groupDataCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8023);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8024, 7701, '/520/521/57/7701/', '', '编辑', 'groupDataUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8024);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8025, 7701, '/520/521/57/7701/', '', '删除', 'groupDataDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8025);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8026, 7701, '/520/521/57/7701/', '', '修改状态', 'groupDataChangeStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8026);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8027, 7702, '/520/521/57/7702/', '', '上传图片', 'uploadImage', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8027);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8028, 7702, '/520/521/57/7702/', '', '图片列表', 'systemAttachmentLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8028);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8029, 7703, '/110/9360/684/7703/', '', '详情', 'groupDetail', '67', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8029);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8030, 7703, '/110/9360/684/7703/', '', '列表', 'groupDataLst', '67', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8030);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8031, 7703, '/110/9360/684/7703/', '', '添加', 'groupDataCreate', '67', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8031);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8032, 7703, '/110/9360/684/7703/', '', '编辑', 'groupDataUpdate', '67', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8032);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8033, 7703, '/110/9360/684/7703/', '', '删除', 'groupDataDelete', '67', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8033);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8034, 7703, '/110/9360/684/7703/', '', '修改状态', 'groupDataChangeStatus', '67', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8034);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8035, 7704, '/110/9360/684/7704/', '', '上传图片', 'uploadImage', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8035);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8036, 7704, '/110/9360/684/7704/', '', '图片列表', 'systemAttachmentLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8036);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8037, 7705, '/514/686/7705/', '', '详情', 'groupDetail', '68', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8037);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8038, 7705, '/514/686/7705/', '', '列表', 'groupDataLst', '68', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8038);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8039, 7705, '/514/686/7705/', '', '添加', 'groupDataCreate', '68', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8039);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8040, 7705, '/514/686/7705/', '', '编辑', 'groupDataUpdate', '68', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8040);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8041, 7705, '/514/686/7705/', '', '删除', 'groupDataDelete', '68', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8041);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8042, 7705, '/514/686/7705/', '', '修改状态', 'groupDataChangeStatus', '68', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8042);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8043, 7706, '/514/686/7706/', '', '上传图片', 'uploadImage', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8043);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8044, 7706, '/514/686/7706/', '', '图片列表', 'systemAttachmentLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8044);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8045, 7707, '/719/5126/687/7707/', '', '详情', 'groupDetail', '69', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8045);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8046, 7707, '/719/5126/687/7707/', '', '列表', 'groupDataLst', '69', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8046);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8047, 7707, '/719/5126/687/7707/', '', '添加', 'groupDataCreate', '69', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8047);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8048, 7707, '/719/5126/687/7707/', '', '编辑', 'groupDataUpdate', '69', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8048);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8049, 7707, '/719/5126/687/7707/', '', '删除', 'groupDataDelete', '69', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8049);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8050, 7707, '/719/5126/687/7707/', '', '修改状态', 'groupDataChangeStatus', '69', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8050);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8051, 7708, '/719/5126/687/7708/', '', '上传图片', 'uploadImage', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8051);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8052, 7708, '/719/5126/687/7708/', '', '图片列表', 'systemAttachmentLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8052);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8053, 7709, '/514/685/7709/', '', '详情', 'groupDetail', '75', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8053);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8054, 7709, '/514/685/7709/', '', '列表', 'groupDataLst', '75', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8054);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8055, 7709, '/514/685/7709/', '', '添加', 'groupDataCreate', '75', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8055);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8056, 7709, '/514/685/7709/', '', '编辑', 'groupDataUpdate', '75', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8056);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8057, 7709, '/514/685/7709/', '', '删除', 'groupDataDelete', '75', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8057);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8058, 7709, '/514/685/7709/', '', '修改状态', 'groupDataChangeStatus', '75', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8058);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8059, 7710, '/514/685/7710/', '', '上传图片', 'uploadImage', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8059);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8060, 7710, '/514/685/7710/', '', '图片列表', 'systemAttachmentLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8060);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8061, 7711, '/514/677/7711/', '', '详情', 'groupDetail', '76', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8061);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8062, 7711, '/514/677/7711/', '', '列表', 'groupDataLst', '76', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8062);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8063, 7711, '/514/677/7711/', '', '添加', 'groupDataCreate', '76', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8063);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8064, 7711, '/514/677/7711/', '', '编辑', 'groupDataUpdate', '76', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8064);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8065, 7711, '/514/677/7711/', '', '删除', 'groupDataDelete', '76', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8065);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8066, 7711, '/514/677/7711/', '', '修改状态', 'groupDataChangeStatus', '76', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8066);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8067, 7712, '/514/677/7712/', '', '上传图片', 'uploadImage', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8067);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8068, 7712, '/514/677/7712/', '', '图片列表', 'systemAttachmentLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8068);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8173, 7739, '/719/1470/7739/', '', '详情', 'groupDetail', '94', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8173);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8174, 7739, '/719/1470/7739/', '', '列表', 'groupDataLst', '94', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8174);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8175, 7739, '/719/1470/7739/', '', '添加', 'groupDataCreate', '94', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8175);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8176, 7739, '/719/1470/7739/', '', '编辑', 'groupDataUpdate', '94', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8176);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8177, 7739, '/719/1470/7739/', '', '删除', 'groupDataDelete', '94', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8177);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8178, 7739, '/719/1470/7739/', '', '修改状态', 'groupDataChangeStatus', '94', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8178);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8179, 7740, '/719/1470/7740/', '', '上传图片', 'uploadImage', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8179);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8180, 7740, '/719/1470/7740/', '', '图片列表', 'systemAttachmentLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8180);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8181, 7741, '/719/1289/1290/7741/', '', '积分配置保存', 'systemUserIntegralConfigSave', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8181);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8182, 7742, '/719/1289/1291/7742/', '', '积分统计', 'systemUserIntegralTitle', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8182);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8183, 7742, '/719/1289/1291/7742/', '', '积分日志', 'systemUserIntegralLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8183);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8184, 7742, '/719/1289/1291/7742/', '', '积分导出', 'systemUserIntegralExcel', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8184);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8185, 7743, '/719/1289/1291/7743/', '', '导出列表', 'systemStoreExcelLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8185);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8186, 7743, '/719/1289/1291/7743/', '', '导出下载', 'systemStoreExcelDownload', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8186);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8187, 7744, '/719/1022/1023/7744/', '', '列表', 'systemStoreProductPresellLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8187);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8188, 7744, '/719/1022/1023/7744/', '', '显示/隐藏', 'systemStoreProductPresellShow', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8188);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8189, 7744, '/719/1022/1023/7744/', '', '详情', 'systemStoreProductPresellDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8189);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8190, 7744, '/719/1022/1023/7744/', '', '编辑数据', 'systemStoreProductPresellGet', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8190);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8191, 7744, '/719/1022/1023/7744/', '', '编辑', 'systemStoreProductPresellUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8191);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8192, 7744, '/719/1022/1023/7744/', '', '审核', 'systemStoreProductPresellSwitchStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8192);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8193, 7744, '/719/1022/1023/7744/', '', '设置标签', 'systemStoreProductPresellLabels', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8193);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8194, 7745, '/719/1051/1095/7745/', '', '列表', 'systemStoreProductAssistLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8194);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8195, 7745, '/719/1051/1095/7745/', '', '显示/隐藏', 'systemStoreProductAssistShow', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8195);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8196, 7745, '/719/1051/1095/7745/', '', '详情', 'systemStoreProductAssistDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8196);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8197, 7745, '/719/1051/1095/7745/', '', '编辑', 'systemStoreProductAssistProductUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8197);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8198, 7745, '/719/1051/1095/7745/', '', '审核', 'systemStoreProductAssistStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8198);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8199, 7745, '/719/1051/1095/7745/', '', '编辑数据', 'systemStoreProductAssistGet', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8199);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8200, 7745, '/719/1051/1095/7745/', '', '设置标签', 'systemStoreProductAssistLabels', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8200);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8201, 7746, '/719/1051/1096/7746/', '', '列表', 'systemStoreProductAssistSetLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8201);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8202, 7746, '/719/1051/1096/7746/', '', '详情', 'systemStoreProductAssistSetDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8202);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8203, 7747, '/719/1135/1137/7747/', '', '列表', 'systemStoreProductGroupLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8203);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8204, 7747, '/719/1135/1137/7747/', '', '显示/隐藏', 'systemStoreProductGroupShow', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8204);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8205, 7747, '/719/1135/1137/7747/', '', '详情', 'systemStoreProductGroupDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8205);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8206, 7747, '/719/1135/1137/7747/', '', '编辑', 'systemStoreProductGroupProductUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8206);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8207, 7747, '/719/1135/1137/7747/', '', '审核', 'systemStoreProductGroupStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8207);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8208, 7747, '/719/1135/1137/7747/', '', '编辑数据', 'systemStoreProductGroupGet', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8208);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8209, 7747, '/719/1135/1137/7747/', '', '排序', 'systemStoreProductGroupSort', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8209);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8210, 7747, '/719/1135/1137/7747/', '', '设置标签', 'systemStoreProductGroupLabels', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8210);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8211, 7748, '/719/1135/1138/7748/', '', '列表', 'systemStoreProductGroupBuyingLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8211);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8212, 7748, '/719/1135/1138/7748/', '', '详情', 'systemStoreProductGroupBuyingDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8212);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8213, 7749, '/719/1135/1136/7749/', '', '配置保存', 'configOthersGroupBuyingUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8213);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8214, 7750, '/719/782/781/7750/', '', '列表', 'systemBroadcastRoomLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8214);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8215, 7750, '/719/782/781/7750/', '', '详情', 'systemBroadcastRoomDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8215);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8216, 7750, '/719/782/781/7750/', '', '申请', 'systemBroadcastRoomApply', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8216);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8217, 7750, '/719/782/781/7750/', '', '修改状态', 'systemBroadcastRoomChangeStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8217);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8218, 7750, '/719/782/781/7750/', '', '排序', 'systemBroadcastRoomSort', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8218);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8219, 7750, '/719/782/781/7750/', '', '修改状态', 'systemBroadcastRoomChangeLiveStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8219);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8220, 7750, '/719/782/781/7750/', '', '删除', 'systemBroadcastRoomDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8220);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8221, 7750, '/719/782/781/7750/', '', '商品列表', 'systemBroadcastRoomGoods', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8221);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8222, 7750, '/719/782/781/7750/', '', '客服开关', 'systemBroadcastRoomCloseKf', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8222);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8223, 7750, '/719/782/781/7750/', '', '禁言开关', 'systemBroadcastRoomCloseComment', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8223);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8224, 7750, '/719/782/781/7750/', '', '收录开关', 'systemBroadcastRoomClosesFeeds', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8224);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8225, 7751, '/719/782/783/7751/', '', '列表', 'systemBroadcastGoodsLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8225);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8226, 7751, '/719/782/783/7751/', '', '详情', 'systemBroadcastGoodsDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8226);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8227, 7751, '/719/782/783/7751/', '', '审核', 'systemBroadcastGoodsApply', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8227);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8228, 7751, '/719/782/783/7751/', '', '修改状态', 'systemBroadcastGoodsChangeStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8228);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8229, 7751, '/719/782/783/7751/', '', '排序', 'systemBroadcastGoodsSort', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8229);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8230, 7751, '/719/782/783/7751/', '', '删除', 'systemBroadcastGoodsDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8230);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8231, 7752, '/719/780/779/7752/', '', '列表', 'systemSeckillConfigLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8231);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8232, 7752, '/719/780/779/7752/', '', '添加', 'systemSeckillConfigCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8232);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8233, 7752, '/719/780/779/7752/', '', '编辑', 'systemSeckillConfigUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8233);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8234, 7752, '/719/780/779/7752/', '', '排序', 'systemSeckillConfigStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8234);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8235, 7752, '/719/780/779/7752/', '', '删除', 'systemSeckillConfigDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8235);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8236, 7753, '/719/780/779/7753/', '', '上传图片', 'uploadImage', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8236);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8237, 7753, '/719/780/779/7753/', '', '图片列表', 'systemAttachmentLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8237);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8238, 7754, '/719/780/794/7754/', '', '统计', 'systemStoreSeckillProductLstFilter', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8238);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8239, 7754, '/719/780/794/7754/', '', '列表', 'systemStoreSeckillProductLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8239);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8240, 7754, '/719/780/794/7754/', '', '权限', 'systemStoreSeckillProductDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8240);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8241, 7754, '/719/780/794/7754/', '', '编辑', 'systemStoreSeckillProductUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8241);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8242, 7754, '/719/780/794/7754/', '', '审核', 'systemStoreSeckillProductSwitchStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8242);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8243, 7754, '/719/780/794/7754/', '', '显示/隐藏', 'systemStoreSeckillProductChangeUsed', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8243);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8244, 7754, '/719/780/794/7754/', '', '设置标签', 'systemStoreSeckillProductLabels', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8244);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8245, 7755, '/101/9047/1510/7755/', '', '普通会员等级列表', 'systemUserMemberLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8245);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8246, 7755, '/101/9047/1510/7755/', '', '普通会员等级添加', 'systemUserMemberCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8246);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8247, 7755, '/101/9047/1510/7755/', '', '普通会员等级编辑', 'systemUserMemberUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8247);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8248, 7755, '/101/9047/1510/7755/', '', '普通会员等级删除', 'systemUserMemberDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8248);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8249, 7756, '/101/9047/1510/7756/', '', '上传图片', 'uploadImage', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8249);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8250, 7756, '/101/9047/1510/7756/', '', '图片列表', 'systemAttachmentLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8250);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8251, 7757, '/101/9047/1509/7757/', '', '会员权益', 'systemUserMemberInterestsLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8251);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8252, 7757, '/101/9047/1509/7757/', '', '会员权益添加', 'systemUserMemberInterestsCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8252);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8253, 7757, '/101/9047/1509/7757/', '', '会员权益编辑', 'systemUserMemberInterestsUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8253);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8254, 7757, '/101/9047/1509/7757/', '', '会员权益删除', 'systemUserMemberInterestsDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8254);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8255, 7758, '/101/9047/1509/7758/', '', '上传图片', 'uploadImage', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8255);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8256, 7758, '/101/9047/1509/7758/', '', '图片列表', 'systemAttachmentLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8256);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8257, 7759, '/110/38/48/7759/', '', '平台菜单/权限列表', 'systemMenuGetLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8257);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8258, 7759, '/110/38/48/7759/', '', '平台菜单/权限添加', 'systemMenuCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8258);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8259, 7759, '/110/38/48/7759/', '', '平台菜单/权限编辑', 'systemMenuUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8259);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8260, 7759, '/110/38/48/7759/', '', '平台菜单/权限删除', 'systemMenuDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8260);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8261, 7760, '/9492/1284/43/7760/', '', '商户菜单/权限列表', 'systemMerchantMenuGetLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8261);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8262, 7760, '/9492/1284/43/7760/', '', '商户菜单/权限添加', 'systemMerchantMenuCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8262);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8263, 7760, '/9492/1284/43/7760/', '', '商户菜单/权限编辑', 'systemMerchantMenuUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8263);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8264, 7760, '/9492/1284/43/7760/', '', '商户菜单/权限删除', 'systemMerchantMenuDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8264);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8265, 7761, '/9492/6370/100/7761/', '', '商户分类列表', 'systemMerchantCategoryLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8265);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8266, 7761, '/9492/6370/100/7761/', '', '商户分类添加', 'systemMerchantCategoryCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8266);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8267, 7761, '/9492/6370/100/7761/', '', '商户分类删除', 'systemMerchantCategoryDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8267);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8268, 7761, '/9492/6370/100/7761/', '', '商户分类编辑', 'systemMerchantCategoryUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8268);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8269, 7762, '/9492/6370/778/7762/', '', '列表', 'systemMerchantIntentionLst', '', 1, 0, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8269);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8270, 7762, '/9492/6370/778/7762/', '', '审核', 'systemMerchantIntentionStatus', '', 1, 0, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8270);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8271, 7762, '/9492/6370/778/7762/', '', '删除', 'systemMerchantIntentionDelete', '', 1, 0, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8271);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8274, 7762, '/9492/6370/778/7762/', '', '备注', 'systemMerchantIntentionMark', '', 1, 0, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8274);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8276, 7763, '/9492/6370/44/7763/', '', '商户列表', 'systemMerchantCreateForm', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8276);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8277, 7763, '/9492/6370/44/7763/', '', '商户列表统计', 'systemMerchantCount', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8277);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8278, 7763, '/9492/6370/44/7763/', '', '商户列表', 'systemMerchantLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8278);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8279, 7763, '/9492/6370/44/7763/', '', '商户添加', 'systemMerchantCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8279);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8280, 7763, '/9492/6370/44/7763/', '', '商户编辑', 'systemMerchantUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8280);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8281, 7763, '/9492/6370/44/7763/', '', '商户修改推荐', 'systemMerchantStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8281);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8282, 7763, '/9492/6370/44/7763/', '', '商户开启/关闭', 'systemMerchantClose', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8282);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8283, 7763, '/9492/6370/44/7763/', '', '商户删除', 'systemMerchantDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8283);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8284, 7763, '/9492/6370/44/7763/', '', '商户修改密码', 'systemMerchantAdminPassword', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8284);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8285, 7763, '/9492/6370/44/7763/', '', '商户登录', 'systemMerchantLogin', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8285);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8286, 7763, '/9492/6370/44/7763/', '', '修改采集商品次数', 'systemMerchantChangeCopy', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8286);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8287, 7764, '/9492/6370/44/7764/', '', '上传图片', 'uploadImage', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8287);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8288, 7764, '/9492/6370/44/7764/', '', '图片列表', 'systemAttachmentLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8288);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8289, 7765, '/9492/6370/1381/7765/', '', '列表', 'systemMerchantTypeLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8289);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8290, 7765, '/9492/6370/1381/7765/', '', '添加', 'systemMerchantTypeCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8290);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8291, 7765, '/9492/6370/1381/7765/', '', '编辑', 'systemMerchantTypeUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8291);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8292, 7765, '/9492/6370/1381/7765/', '', '删除', 'systemMerchantTypeDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8292);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8293, 7765, '/9492/6370/1381/7765/', '', '备注', 'systemMerchantTypeMark', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8293);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8294, 7765, '/9492/6370/1381/7765/', '', '备注', 'systemMerchantTypeDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8294);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8295, 7766, '/9492/1284/1647/7766/', '', '缴纳记录', 'systemMerchantMarginLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8295);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8296, 7766, '/9492/1284/1647/7766/', '', '扣费记录', 'systemMarginList', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8296);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8297, 7766, '/9492/1284/1647/7766/', '', '扣除保证金', 'systemMarginSet', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8297);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8298, 7766, '/9492/1284/1647/7766/', '', '退款申请列表', 'systemMarginRefundList', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8298);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8299, 7766, '/9492/1284/1647/7766/', '', '退款申请详情', 'systemMarginRefundShow', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8299);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8300, 7766, '/9492/1284/1647/7766/', '', '审核', 'systemMarginRefundSwitchStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8300);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8301, 7766, '/9492/1284/1647/7766/', '', '备注', 'systemMarginRefundMark', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8301);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8302, 7767, '/110/5125/1120/7767/', '', '系统公告列表', 'systemNoticeList', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8302);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8303, 7767, '/110/5125/1120/7767/', '', '系统公告发布', 'systemNoticeCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8303);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8304, 7768, '/110/5125/1597/7768/', '', '消息配置列表', 'systemNoticeConfigLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8304);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8305, 7768, '/110/5125/1597/7768/', '', '消息配置添加', 'systemNoticeConfigCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8305);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8306, 7768, '/110/5125/1597/7768/', '', '消息配置编辑', 'systemNoticeConfigUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8306);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8307, 7768, '/110/5125/1597/7768/', '', '消息配置详情', 'systemNoticeConfigDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8307);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8308, 7768, '/110/5125/1597/7768/', '', '消息配置删除', 'systemNoticeConfigDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8308);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8309, 7768, '/110/5125/1597/7768/', '', '消息配置修改状态', 'systemNoticeConfigStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8309);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8310, 7769, '/540/541/7769/', '', '列表', 'systemOrderLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8310);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8311, 7769, '/540/541/7769/', '', '金额统计', 'systemOrderStat', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8311);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8312, 7769, '/540/541/7769/', '', '快递查询', 'systemOrderExpress', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8312);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8313, 7769, '/540/541/7769/', '', '头部统计', 'systemOrderTitle', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8313);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8314, 7769, '/540/541/7769/', '', '详情', 'systemOrderDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8314);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8315, 7769, '/540/541/7769/', '', '导出', 'systemOrderExcel', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8315);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8316, 7770, '/540/541/7770/', '', '导出列表', 'systemStoreExcelLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8316);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8317, 7770, '/540/541/7770/', '', '导出列表', 'systemStoreExcelDownload', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8317);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8318, 7771, '/540/784/7771/', '', '核销', 'systemOrderTakeStat', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8318);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8319, 7771, '/540/784/7771/', '', '核销订单', 'systemTakeOrderLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8319);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8320, 7771, '/540/784/7771/', '', '头部统计', 'systemTakeOrderTitle', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8320);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8321, 7772, '/540/542/7772/', '', '列表', 'systemRefundOrderLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8321);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8322, 7772, '/540/542/7772/', '', '导出', 'systemRefundOrderExcel', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8322);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8323, 7773, '/540/542/7773/', '', '导出列表', 'systemStoreExcelLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8323);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8324, 7773, '/540/542/7773/', '', '导出下载', 'systemStoreExcelDownload', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8324);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8325, 7774, '/87/88/7774/', '', '编辑', 'systemStoreCategoryUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8325);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8326, 7774, '/87/88/7774/', '', '列表', 'systemStoreCategoryLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8326);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8327, 7774, '/87/88/7774/', '', '详情', 'systemStoreCategoryDtailt', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8327);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8328, 7774, '/87/88/7774/', '', '添加', 'systemStoreCategoryCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8328);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8329, 7774, '/87/88/7774/', '', '删除', 'systemStoreCategoryDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8329);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8330, 7774, '/87/88/7774/', '', '修改状态', 'systemStoreCategorySwitchStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8330);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8331, 7774, '/87/88/7774/', '', '修改推荐', 'systemStoreCategorySwitchIsHot', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8331);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8332, 7775, '/87/88/7775/', '', '上传图片', 'uploadImage', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8332);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8333, 7775, '/87/88/7775/', '', '图片列表', 'systemAttachmentLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8333);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8334, 7776, '/87/92/93/7776/', '', '编辑', 'systemStoreBrandCategoryUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8334);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8335, 7776, '/87/92/93/7776/', '', '列表', 'systemStoreBrandCategoryLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8335);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8336, 7776, '/87/92/93/7776/', '', '详情', 'systemStoreBrandCategoryDtailt', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8336);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8337, 7776, '/87/92/93/7776/', '', '添加', 'systemStoreBrandCategoryCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8337);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8338, 7776, '/87/92/93/7776/', '', '删除', 'systemStoreBrandCategoryDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8338);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8339, 7776, '/87/92/93/7776/', '', '修改状态', 'systemStoreBrandCategorySwitchStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8339);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8340, 7777, '/87/92/94/7777/', '', '列表', 'systemStoreBrandLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8340);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8341, 7777, '/87/92/94/7777/', '', '修改状态', 'systemStoreBrandSwithStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8341);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8342, 7777, '/87/92/94/7777/', '', '添加', 'systemStoreBrandCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8342);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8343, 7777, '/87/92/94/7777/', '', '编辑', 'systemStoreBrandUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8343);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8344, 7777, '/87/92/94/7777/', '', '删除', 'systemStoreBrandDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8344);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8345, 7778, '/87/539/7778/', '', '统计', 'systemStoreProductLstFilter', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8345);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8346, 7778, '/87/539/7778/', '', '列表', 'systemStoreProductLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8346);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8347, 7778, '/87/539/7778/', '', '详情', 'systemStoreProductDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8347);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8348, 7778, '/87/539/7778/', '', '编辑', 'systemStoreProductUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8348);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8349, 7778, '/87/539/7778/', '', '上下架', 'systemStoreProductSwitchStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8349);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8350, 7778, '/87/539/7778/', '', '分销状态变更商品检测', 'systemStoreProductCheck', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8350);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8351, 7778, '/87/539/7778/', '', '显示/隐藏', 'systemStoreProductChangeUsed', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8351);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8352, 7778, '/87/539/7778/', '', '虚拟销量', 'systemStoreProductAddFicti', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8352);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8353, 7778, '/87/539/7778/', '', '设置标签', 'systemStoreProductLabels', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8353);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8354, 7779, '/87/121/7779/', '', '列表', 'systemProductReplyLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8354);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8355, 7779, '/87/121/7779/', '', '添加虚拟评论', 'systemProductReplyCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8355);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8356, 7779, '/87/121/7779/', '', '排序', 'systemProductReplySort', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8356);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8357, 7779, '/87/121/7779/', '', '删除', 'systemProductReplyDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8357);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8358, 7780, '/87/121/7780/', '', '上传图片', 'uploadImage', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8358);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8359, 7780, '/87/121/7780/', '', '图片列表', 'systemAttachmentLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8359);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8360, 7781, '/87/1245/7781/', '', '列表', 'systemGuaranteeLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8360);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8361, 7781, '/87/1245/7781/', '', '添加', 'systemGuaranteeCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8361);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8362, 7781, '/87/1245/7781/', '', '编辑', 'systemGuaranteeUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8362);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8363, 7781, '/87/1245/7781/', '', '详情', 'systemGuaranteeDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8363);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8364, 7781, '/87/1245/7781/', '', '删除', 'systemGuaranteeDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8364);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8365, 7781, '/87/1245/7781/', '', '排序', 'systemGuaranteeSort', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8365);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8366, 7781, '/87/1245/7781/', '', '修改状态', 'systemGuaranteeStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8366);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8367, 7782, '/87/1245/7782/', '', '上传图片', 'uploadImage', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8367);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8368, 7782, '/87/1245/7782/', '', '图片列表', 'systemAttachmentLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8368);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8369, 7783, '/87/1469/7783/', '', '列表', 'systemStoreProductLabelLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8369);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8370, 7783, '/87/1469/7783/', '', '添加', 'systemStoreProductLabelCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8370);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8371, 7783, '/87/1469/7783/', '', '编辑', 'systemStoreProductLabelUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8371);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8372, 7783, '/87/1469/7783/', '', '详情', 'systemStoreProductLabelDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8372);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8373, 7783, '/87/1469/7783/', '', '删除', 'systemStoreProductLabelDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8373);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8374, 7783, '/87/1469/7783/', '', '修改状态', 'systemStoreProductLabelStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8374);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8375, 7784, '/719/1629/7784/', '', '优惠套餐列表', 'systemStoreDiscountsLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8375);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8376, 7784, '/719/1629/7784/', '', '优惠套餐详情', 'systemStoreDiscountsDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8376);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8377, 7784, '/719/1629/7784/', '', '优惠套餐修改状态', 'systemStoreDiscountsStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8377);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8378, 7785, '/719/1629/7785/', '', '上传图片', 'uploadImage', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8378);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8379, 7785, '/719/1629/7785/', '', '图片列表', 'systemAttachmentLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8379);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8380, 7786, '/110/38/39/7786/', '', '身份列表', 'systemRoleGetList', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8380);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8381, 7786, '/110/38/39/7786/', '', '身份添加', 'systemRoleCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8381);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8382, 7786, '/110/38/39/7786/', '', '身份编辑', 'systemRoleUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8382);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8383, 7786, '/110/38/39/7786/', '', '身份修改状态', 'systemRoleStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8383);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8384, 7786, '/110/38/39/7786/', '', '身份删除', 'systemRoleDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8384);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8385, 7787, '/110/38/40/7787/', '', '管理员列表', 'systemAdminLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8385);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8386, 7787, '/110/38/40/7787/', '', '管理员修改状态', 'systemAdminStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8386);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8387, 7787, '/110/38/40/7787/', '', '管理员添加', 'systemAdminCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8387);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8388, 7787, '/110/38/40/7787/', '', '管理员编辑', 'systemAdminUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8388);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8389, 7787, '/110/38/40/7787/', '', '管理员修改密码', 'systemAdminPassword', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8389);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8390, 7787, '/110/38/40/7787/', '', '管理员删除', 'systemAdminDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8390);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8391, 7788, '/110/38/40/7788/', '', '上传图片', 'uploadImage', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8391);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8392, 7788, '/110/38/40/7788/', '', '图片列表', 'systemAttachmentLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8392);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8393, 7789, '/520/47/7789/', '', '操作日志', 'systemAdminLog', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8393);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8394, 7790, '/7790/', '', '修改信息', 'systemAdminEdit', '', 0, 0, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8394);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8395, 7790, '/7790/', '', '修改密码', 'systemAdminEditPassword', '', 0, 0, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8395);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8396, 7791, '/5054/5088/7791/', '', '列表', 'adminServiceLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8396);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8397, 7791, '/5054/5088/7791/', '', '登录', 'adminServiceLogin', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8397);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8398, 7791, '/5054/5088/7791/', '', '添加', 'adminServiceCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8398);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8399, 7791, '/5054/5088/7791/', '', '编辑', 'adminServiceUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8399);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8400, 7791, '/5054/5088/7791/', '', '修改状态', 'adminServiceSwitchStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8400);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8401, 7791, '/5054/5088/7791/', '', '删除', 'adminServiceDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8401);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8402, 7791, '/5054/5088/7791/', '', '客服的全部用户 ', 'adminServiceServiceUserList', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8402);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8403, 7791, '/5054/5088/7791/', '', '用户与客服聊天记录', 'adminServiceServiceUserLogLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8403);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8404, 7791, '/5054/5088/7791/', '', '客服的聊天用户列表', 'adminServiceServiceMerchantUserList', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8404);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8405, 7791, '/5054/5088/7791/', '', '用户与商户聊天记录', 'adminServiceMerchantUserLogLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8405);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8406, 7792, '/5054/1668/7792/', '', '列表', 'adminServiceReplyLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8406);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8407, 7792, '/5054/1668/7792/', '', '添加', 'adminServiceReplyCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8407);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8408, 7792, '/5054/1668/7792/', '', '编辑', 'adminServiceReplyUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8408);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8409, 7792, '/5054/1668/7792/', '', '切换状态', 'adminServiceReplyStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8409);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8410, 7792, '/5054/1668/7792/', '', '删除', 'adminServiceReplyDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8410);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8411, 7793, '/514/522/7793/', '', '分销员列表', 'systemPromoterUserLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8411);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8412, 7793, '/514/522/7793/', '', '分销员统计', 'systemPromoterUserCount', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8412);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8413, 7793, '/514/522/7793/', '', '修改分销员等级', 'systemUserSpreadSave', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8413);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8414, 7793, '/514/522/7793/', '', '推广人列表', 'systemUserSpreadLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8414);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8415, 7793, '/514/522/7793/', '', '推广人订单', 'systemUserSpreadOrder', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8415);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8416, 7793, '/514/522/7793/', '', '清除推广人', 'systemUserSpreadClear', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8416);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8417, 7794, '/514/1373/1374/7794/', '', '分销员等级列表', 'systemUserBrokerageLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8417);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8418, 7794, '/514/1373/1374/7794/', '', '分销员等级添加', 'systemUserBrokerageCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8418);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8419, 7794, '/514/1373/1374/7794/', '', '分销员等级编辑', 'systemUserBrokerageUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8419);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8420, 7794, '/514/1373/1374/7794/', '', '分销员等级删除', 'systemUserBrokerageDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8420);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8421, 7795, '/514/1373/1374/7795/', '', '上传图片', 'uploadImage', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8421);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8422, 7795, '/514/1373/1374/7795/', '', '图片列表', 'systemAttachmentLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8422);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8423, 7796, '/514/731/7796/', '', '统计', 'systemStoreBagLstFilter', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8423);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8424, 7796, '/514/731/7796/', '', '列表', 'systemStoreBagLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8424);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8425, 7796, '/514/731/7796/', '', '详情', 'systemStoreBagDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8425);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8426, 7796, '/514/731/7796/', '', '编辑', 'systemStoreBagUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8426);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8427, 7796, '/514/731/7796/', '', '修改状态', 'systemStoreBagSwitchStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8427);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8428, 7796, '/514/731/7796/', '', '显示/隐藏', 'systemStoreBagChangeUsed', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8428);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8429, 7797, '/514/5122/7797/', '', '配置保存', 'configOthersSettingUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8429);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8430, 7798, '/520/1244/7798/', '', '列表', 'systemStoreExcelLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8430);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8431, 7798, '/520/1244/7798/', '', '下载', 'systemStoreExcelDownload', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8431);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8432, 7799, '/9218/33/7799/', '', '主要数据', 'systemStatisticsMain', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8432);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8433, 7799, '/9218/33/7799/', '', '当日订单', 'systemStatisticsOrder', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8433);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8434, 7799, '/9218/33/7799/', '', '当日订单数', 'systemStatisticsOrderNum', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8434);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8435, 7799, '/9218/33/7799/', '', '当日支付人数', 'systemStatisticsOrderUser', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8435);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8436, 7799, '/9218/33/7799/', '', '商户销量', 'systemStatisticsMerchantStock', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8436);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8437, 7799, '/9218/33/7799/', '', '商户访问量', 'systemStatisticsMerchantRate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8437);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8438, 7799, '/9218/33/7799/', '', '商户销售额', 'systemStatisticsMerchantVisit', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8438);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8439, 7799, '/9218/33/7799/', '', '用户数据', 'systemStatisticsUserData', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8439);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8440, 7799, '/9218/33/7799/', '', '成交用户', 'systemStatisticsUser', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8440);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8441, 7799, '/9218/33/7799/', '', '成交用户占比', 'systemStatisticsUserRate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8441);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8442, 7800, '/520/116/117/7800/', '', '数据库列表', 'systemSafetyDatabaseLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8442);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8443, 7800, '/520/116/117/7800/', '', '数据库备份列表', 'systemSafetyDatabaseFileList', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8443);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8444, 7800, '/520/116/117/7800/', '', '数据库备份详情', 'systemSafetyDatabaseDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8444);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8445, 7800, '/520/116/117/7800/', '', '备份', 'systemSafetyDatabaseBackups', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8445);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8446, 7800, '/520/116/117/7800/', '', '数据库优化', 'systemSafetyDatabaseOptimize', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8446);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8447, 7800, '/520/116/117/7800/', '', '数据库维护', 'systemSafetyDatabaseRepair', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8447);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8448, 7800, '/520/116/117/7800/', '', '数据库备份下载', 'systemSafetyDatabaseDownloadFile', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8448);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8449, 7800, '/520/116/117/7800/', '', '数据库备份删除', 'systemSafetyDatabaseDeleteFile', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8449);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8450, 7801, '/520/116/1617/7801/', '', '清除缓存', 'systemClearCache', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8450);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8451, 9330, '9330/', '', '权限', '/delivery/personnel_manage/permission', '', 1, 0, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8451);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8452, 7803, '/110/9360/1667/7803/', '', '协议列表', 'systemAgreeKeyLsy', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8452);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8453, 7803, '/110/9360/1667/7803/', '', '商户入住申请协议', 'systemAgreeSave', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8453);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8454, 7804, '/719/1022/1024/7804/', '', '预售协议', 'systemAgreeSave', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8454);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8455, 7805, '/514/1296/7805/', '', '佣金说明', 'systemAgreeSave', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8455);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8456, 7806, '/514/1373/1375/7806/', '', '等级规则', 'systemAgreeSave', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8456);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8457, 7807, '/719/1657/1663/7807/', '', '使用说明', 'systemAgreeSave', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8457);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8458, 7808, '/101/1247/7808/', '', '用户协议', 'systemAgreeSave', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8458);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8459, 7809, '/101/9047/1511/7809/', '', '会员等级规则', 'systemAgreeSave', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8459);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8460, 7810, '/9492/1284/1382/7810/', '', '店铺类型说明 ', 'systemAgreeSave', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8460);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8461, 7811, '/515/1298/1300/7811/', '', '发票说明 ', 'systemAgreeSave', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8461);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8462, 7812, '/101/104/7812/', '', '用户标签列表', 'systemUserLabelLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8462);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8463, 7812, '/101/104/7812/', '', '用户标签添加', 'systemUserLabelCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8463);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8464, 7812, '/101/104/7812/', '', '用户标签删除', 'systemUserLabelDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8464);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8465, 7812, '/101/104/7812/', '', '用户标签编辑', 'systemUserLabelUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8465);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8466, 7813, '/101/1285/7813/', '', '用户搜索记录', 'systemUserSearchLog', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8466);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8467, 7813, '/101/1285/7813/', '', '用户搜索记录导出', 'systemUserExportSearchLog', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8467);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8468, 7814, '/101/1285/7814/', '', '导出列表', 'systemStoreExcelLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8468);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8469, 7814, '/101/1285/7814/', '', '导出下载', 'systemStoreExcelDownload', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8469);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8470, 7815, '/101/102/7815/', '', '用户分组列表', 'systemUserGroupLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8470);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8471, 7815, '/101/102/7815/', '', '用户分组添加', 'systemUserGroupCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8471);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8472, 7815, '/101/102/7815/', '', '用户分组删除', 'systemUserGroupDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8472);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8473, 7815, '/101/102/7815/', '', '用户分组编辑', 'systemUserGroupUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8473);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8474, 7816, '/101/466/467/7816/', '', '列表', 'systemUserFeedBackCategoryLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8474);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8475, 7816, '/101/466/467/7816/', '', '添加', 'systemUserFeedBackCategoryCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8475);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8476, 7816, '/101/466/467/7816/', '', '编辑', 'systemUserFeedBackCategoryUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8476);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8477, 7816, '/101/466/467/7816/', '', '修改状态', 'systemUserFeedBackCategorySwitchStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8477);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8478, 7816, '/101/466/467/7816/', '', '删除', 'systemUserFeedBackCategoryDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8478);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8479, 7817, '/101/466/468/7817/', '', '列表', 'systemUserFeedBackLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8479);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8480, 7817, '/101/466/468/7817/', '', '详情', 'systemUserFeedBackDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8480);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8481, 7817, '/101/466/468/7817/', '', '回复', 'systemUserFeedBackReply', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8481);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8482, 7817, '/101/466/468/7817/', '', '删除', 'systemUserFeedBackDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8482);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8483, 7818, '/110/670/1595/7818/', '', '上传文件', 'configUpload', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8483);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8484, 7818, '/110/670/1595/7818/', '', '上传原名文件', 'configUploadName', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8484);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8485, 7818, '/110/670/1595/7818/', '', '微信校验文件上传', 'configWechatUploadSet', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8485);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8486, 7818, '/110/670/1595/7818/', '', '小程序配置', 'configRoutineConfig', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8486);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8487, 7819, '/519/524/1596/7819/', '', '小程序下载', 'configRoutineDownload', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8487);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8488, 7820, '/519/58/59/7820/', '', '微信菜单配置', 'wechatMenu', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8488);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8489, 7820, '/519/58/59/7820/', '', '保存微信菜单配置', 'saveWechatMenu', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8489);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8490, 7821, '/519/58/77/7821/', '', '详情', 'wechatReplyInfo', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8490);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8491, 7821, '/519/58/77/7821/', '', '编辑', 'saveWechatReply', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8491);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8492, 7821, '/519/58/77/7821/', '', '添加', 'createWechatReply', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8492);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8493, 7821, '/519/58/77/7821/', '', '修改', 'updateWechatReply', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8493);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8494, 7821, '/519/58/77/7821/', '', '列表', 'wechatReplyLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8494);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8495, 7821, '/519/58/77/7821/', '', '删除', 'wechatReplyDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8495);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8496, 7821, '/519/58/77/7821/', '', '修改状态', 'wechatReplyStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8496);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8497, 7821, '/519/58/77/7821/', '', '上传图片', 'wechatUploadImage', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8497);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8498, 7821, '/519/58/77/7821/', '', '上传语音', 'wechatUploadVoice', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8498);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8499, 7822, '/519/58/82/7822/', '', '添加', 'systemWechatNewsCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8499);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8500, 7822, '/519/58/82/7822/', '', '编辑', 'systemWechatNewsUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8500);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8501, 7822, '/519/58/82/7822/', '', '删除', 'systemWechatNewsDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8501);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8502, 7822, '/519/58/82/7822/', '', '列表', 'systemWechatNewsLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8502);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8503, 7822, '/519/58/82/7822/', '', '详情', 'systemWechatNewsDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8503);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8504, 7823, '/519/58/82/7823/', '', '上传图片', 'uploadImage', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8504);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8505, 7823, '/519/58/82/7823/', '', '图片列表', 'systemAttachmentLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8505);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8506, 7824, '/519/58/532/7824/', '', '同步', 'systemTemplateMessageSync', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8506);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8507, 7824, '/519/58/532/7824/', '', '列表', 'systemTemplateMessageLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8507);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8508, 7824, '/519/58/532/7824/', '', '添加', 'systemTemplateMessageCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8508);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8509, 7824, '/519/58/532/7824/', '', '编辑', 'systemTemplateMessageUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8509);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8510, 7824, '/519/58/532/7824/', '', '删除', 'systemTemplateMessageDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8510);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8511, 7824, '/519/58/532/7824/', '', '修改状态', 'systemTemplateMessageSwitchStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8511);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8512, 7825, '/519/524/545/7825/', '', '同步', 'systemTemplateMessageMinSync', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8512);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8513, 7825, '/519/524/545/7825/', '', '列表 ', 'systemTemplateMessageMinList', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8513);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8514, 7825, '/519/524/545/7825/', '', '添加', 'systemTemplateMessageMinCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8514);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8515, 7825, '/519/524/545/7825/', '', '编辑', 'systemTemplateMessageMinUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8515);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8516, 7825, '/519/524/545/7825/', '', '删除', 'systemTemplateMessageMinDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8516);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8517, 7825, '/519/524/545/7825/', '', '修改状态', 'systemTemplateMessageMinSwitchStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8517);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8518, 7826, '/110/118/1368/1369/7826/', '', '使用记录', 'systemStoreProductCopyLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8518);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8519, 7826, '/110/118/1368/1369/7826/', '', '短信发送记录', 'smsRecord', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8519);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8520, 7826, '/110/118/1368/1369/7826/', '', '退出登录', 'smsLogout', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8520);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8521, 7826, '/110/118/1368/1369/7826/', '', '获取验证码', 'systemServeCaptcha', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8521);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8522, 7826, '/110/118/1368/1369/7826/', '', '验证码校验', 'systemServeCaptchaCheck', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8522);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8523, 7826, '/110/118/1368/1369/7826/', '', '注册', 'systemServeRegister', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8523);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8524, 7826, '/110/118/1368/1369/7826/', '', '登录', 'systemServeLogin', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8524);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8525, 7826, '/110/118/1368/1369/7826/', '', '修改密码', 'systemServeChangePassword', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8525);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8526, 7826, '/110/118/1368/1369/7826/', '', '修改手机号', 'systemServeChangePhone', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8526);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8527, 7826, '/110/118/1368/1369/7826/', '', '检测登录状态', 'systemServeIsLogin', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8527);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8528, 7826, '/110/118/1368/1369/7826/', '', '使用记录', 'systemServeRecordLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8528);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8529, 7826, '/110/118/1368/1369/7826/', '', '套餐列表', 'systemServeMealLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8529);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8530, 7826, '/110/118/1368/1369/7826/', '', '购买套餐', 'systemServePayMeal', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8530);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8531, 7826, '/110/118/1368/1369/7826/', '', '开通服务', 'systemServeOpenServe', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8531);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8532, 7826, '/110/118/1368/1369/7826/', '', '修改签名', 'systemServeChangeSign', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8532);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8533, 7826, '/110/118/1368/1369/7826/', '', '模板', 'systemServeExportTemps', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8533);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8534, 7826, '/110/118/1368/1369/7826/', '', '使用记录', 'systemServeExportDumpLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8534);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8535, 7827, '/110/118/1368/1372/7827/', '', '购买记录', 'systemServePayLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8535);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8536, 7827, '/110/118/1368/1372/7827/', '', '商户购买记录', 'systemServeMerPayLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8536);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8537, 7828, '/110/118/1368/1379/7828/', '', '商户结余', 'systemServeMerLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8537);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8538, 7829, '/110/118/1368/1371/7829/', '', '列表', 'systemServeMerMealLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8538);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8539, 7829, '/110/118/1368/1371/7829/', '', '详情', 'systemServeMealDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8539);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8540, 7829, '/110/118/1368/1371/7829/', '', '添加', 'systemServeMealCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8540);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8541, 7829, '/110/118/1368/1371/7829/', '', '编辑', 'systemServeMealUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8541);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8542, 7829, '/110/118/1368/1371/7829/', '', '删除', 'systemServeMealDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8542);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8543, 7829, '/110/118/1368/1371/7829/', '', '修改状态', 'systemServeMealStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8543);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8544, 7830, '/110/118/1368/111/113/7830/', '', '短信模板', 'systemServeSmsTemps', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8544);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8545, 7830, '/110/118/1368/111/113/7830/', '', '申请模板', 'systemServeSmsApply', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8545);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8546, 7831, '/110/118/1368/111/114/7831/', '', '申请记录', 'systemServeSmsApplyRecord', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8546);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8547, 1181, '/525/1181/', '', '权限', '/accounts/transManagement', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8547);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8548, 1181, '/525/1181/', '', '附加权限', 'append_/accounts/transManagement', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8548);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8549, 1182, '/525/1182/', '', '权限', '/accounts/payType', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8549);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8550, 790, '/525/790/', '', '权限', '/accounts/capitalFlow', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8550);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8551, 790, '/525/790/', '', '附加权限', 'append_/accounts/capitalFlow', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8551);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8552, 1183, '/525/1183/', '', '权限', '/accounts/statement', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8552);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8553, 1183, '/525/1183/', '', '附加权限', 'append_/accounts/statement', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8553);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8554, 1102, '/525/1102/', '', '权限', '/order/invoice', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8554);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8555, 1366, '/525/1366/', '', '权限', '/systemForm/applyList', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8555);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8556, 1366, '/525/1366/', '', '附加权限', 'append_/systemForm/applyList', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8556);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8557, 1304, '/525/1304/', '', '权限', '/systemForm/applyments', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8557);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8558, 20054, '1671/54/', '', '权限', '/config/picture', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8558);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8559, 1103, '/1027/1103/', '', '权限', '/user/list', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8559);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8560, 1132, '/107/1132/', '', '权限', '/marketing/coupon/send', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8560);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8561, 20115, '/106/107/115/', '', '权限', '/marketing/coupon/list', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8561);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8562, 20115, '/106/107/115/', '', '附加权限', 'append_/marketing/coupon/list', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8562);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8563, 20122, '/106/107/122/', '', '权限', '/marketing/coupon/user', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8563);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8564, 1380, '1649/1380/', '', '权限', '/config/freight/express', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8564);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8565, 9452, '9332/1650/', '', '权限', '/delivery/store_manage', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8565);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8566, 1652, '9332/1652/', '', '权限', '/delivery/recharge_record', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8566);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8567, 1651, '9332/1651/', '', '权限', '/delivery/usage_record', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8567);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8568, 700, '1649/700/', '', '权限', '/config/freight/shippingTemplates', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8568);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8569, 1672, '/1671/1672/', '', '权限', '/devise/diy/list', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8569);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8570, 1672, '/1671/1672/', '', '附加权限', 'append_/devise/diy/list', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8570);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8573, 1471, '/106/1471/', '', '权限', '/group/topic/95', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8573);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8574, 1471, '/106/1471/', '', '附加权限', 'append_/group/topic/95', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8574);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8575, 9290, '9290/', '', '权限', '/marketing/seckill/product/list', '', 1, 0, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8575);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8576, 9289, '9289/', '', '附加权限', 'append_/marketing/seckill/store_seckill/list', '', 1, 0, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8576);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8577, 1025, '/106/1025/', '', '权限', '/marketing/presell/list', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8577);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8578, 1025, '/106/1025/', '', '附加权限', 'append_/marketing/presell/list', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8578);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8579, 1100, '/1099/1100/', '', '权限', '/marketing/assist/list', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8579);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8580, 1100, '/1099/1100/', '', '附加权限', 'append_/marketing/assist/list', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8580);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8581, 1101, '/1099/1101/', '', '权限', '/marketing/assist/assist_set', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8581);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8582, 1140, '/1139/1140/', '', '权限', '/marketing/combination/combination_goods', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8582);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8583, 1140, '/1139/1140/', '', '附加权限', 'append_/marketing/combination/combination_goods', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8583);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8584, 1141, '/1139/1141/', '', '权限', '/marketing/combination/combination_list', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8584);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8585, 1141, '/1139/1141/', '', '附加权限', 'append_/marketing/combination/combination_list', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8585);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8586, 786, '/106/785/786/', '', '权限', '/marketing/studio/list', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8586);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8587, 786, '/106/785/786/', '', '附加权限', 'append_/marketing/studio/list', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8587);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8588, 1594, '/785/1594/', '', '权限', '/marketing/studio/assistant', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8588);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8589, 1594, '/785/1594/', '', '附加权限', 'append_/marketing/studio/assistant', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8589);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8590, 1295, '/1293/1295/', '', '权限', '/marketing/integral/log', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8590);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8591, 1295, '/1293/1295/', '', '附加权限', 'append_/marketing/integral/log', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8591);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8592, 513, '/512/513/', '', '权限', '/order/list', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8592);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8593, 513, '/512/513/', '', '附加权限', 'append_/order/list', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8593);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8594, 789, '/512/789/', '', '权限', '/order/cancellation', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8594);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8595, 528, '/512/528/', '', '权限', '/order/refund', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8595);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8596, 528, '/512/528/', '', '附加权限', 'append_/order/refund', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8596);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8597, 20099, '/95/99/', '', '权限', '/product/attr', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8597);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8598, 20096, '/95/96/', '', '权限', '/product/classify', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8598);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8599, 20096, '/95/96/', '', '附加权限', 'append_/product/classify', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8599);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8600, 20105, '/95/105/', '', '权限', '/product/list', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8600);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8601, 20105, '/95/105/', '', '附加权限', 'append_/product/list', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8601);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8602, 544, '512/544/', '', '权限', '/product/reviews', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8602);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8603, 544, '512/544/', '', '附加权限', 'append_/product/reviews', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8603);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8604, 1468, '/95/1468/', '', '权限', '/product/label', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8604);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8605, 1630, '/106/1630/', '', '权限', '/marketing/discounts/list', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8605);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8606, 1630, '/106/1630/', '', '附加权限', 'append_/marketing/discounts/list', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8606);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8607, 1246, '/95/1246/', '', '权限', '/config/guarantee', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8607);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8608, 20051, '/526/49/51/', '', '权限', '/setting/systemAdmin', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8608);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8609, 0, '/', '', '权限', 'self', '', 1, 0, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8609);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8610, 20120, '9410/120/', '', '权限', '/config/service', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8610);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8611, 5123, '9410/5123/', '', '权限', '/systemForm/customer_keyword', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8611);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8613, 1656, '9358/1656/', '', '权限', '/setting/printer/list', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8613);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8614, 20055, '/55/', '', '权限', '/dashboard', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8614);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8615, 1119, '/1119/', '', '权限', '/station/notice', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8615);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8616, 20074, '/526/74/', '', '权限', '/systemForm/Basics/mer_base', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8616);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8617, 546, '/526/546/', '', '权限', '/systemForm/modifyStoreInfo', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8617);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8618, 1286, '/1027/1286/', '', '权限', '/user/searchRecord', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8618);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8619, 1030, '/1028/1030/', '', '权限', '/user/maticlabel', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8619);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8620, 1029, '/1028/1029/', '', '权限', '/user/label', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8620);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8621, 1378, '526/1376/1378/', '', '权限', '/setting/sms/sms_config/index', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8621);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8622, 1377, '526/1376/1377/', '', '权限', '/setting/sms/dumpConfig', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8622);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8623, 8547, '/525/1181/8547/', '', '列表', 'merchantFinancialRefundMargin', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8623);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8624, 8547, '/525/1181/8547/', '', '转账记录', 'merchantFinancialLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8624);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8625, 8547, '/525/1181/8547/', '', '详情', 'merchantFinancialDetail', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8625);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8626, 8547, '/525/1181/8547/', '', '申请', 'merchantFinancialCreateSave', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8626);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8627, 8547, '/525/1181/8547/', '', '删除', 'merchantFinancialDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8627);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8628, 8547, '/525/1181/8547/', '', '备注', 'merchantFinancialMark', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8628);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8629, 8547, '/525/1181/8547/', '', '导出', 'merchantFinancialExport', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8629);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8630, 8548, '/525/1181/8548/', '', '导出列表', 'merchantStoreExcelLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8630);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8631, 8548, '/525/1181/8548/', '', '导出下载', 'merchantStoreExcelDownload', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8631);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8632, 8549, '/525/1182/8549/', '', '收款方式', 'merchantFinancialAccountSave', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8632);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8633, 8550, '/525/790/8550/', '', '列表', 'merchantFinancialRecordList', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8633);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8634, 8550, '/525/790/8550/', '', '导出', 'merchantFinancialRecordExport', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8634);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8635, 8550, '/525/790/8550/', '', '统计', 'merchantFinancialCount', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8635);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8636, 8551, '/525/790/8551/', '', '导出列表', 'merchantStoreExcelLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8636);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8637, 8551, '/525/790/8551/', '', '导出下载', 'merchantStoreExcelDownload', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8637);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8638, 8552, '/525/1183/8552/', '', '列表', 'merchantFinanciaRecordlLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8638);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8639, 8552, '/525/1183/8552/', '', '统计', 'merchantFinancialTitle', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8639);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8640, 8552, '/525/1183/8552/', '', '详情', 'merchantFinancialRecordDetail', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8640);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8641, 8552, '/525/1183/8552/', '', '导出', 'merchantFinancialRecordDetailExport', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8641);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8642, 8553, '/525/1183/8553/', '', '导出列表', 'merchantStoreExcelLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8642);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8643, 8553, '/525/1183/8553/', '', '导出下载', 'merchantStoreExcelDownload', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8643);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8644, 8554, '/525/1102/8554/', '', '列表', 'merchantOrderReceiptLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8644);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8645, 8554, '/525/1102/8554/', '', '详情', 'merchantOrderReceiptDetail', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8645);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8646, 8554, '/525/1102/8554/', '', '开发票', 'merchantOrderReceiptSetRecipt', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8646);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8647, 8554, '/525/1102/8554/', '', '保存发票', 'merchantOrderReceiptSave', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8647);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8648, 8554, '/525/1102/8554/', '', '备注', 'merchantOrderReceiptMark', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8648);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8649, 8554, '/525/1102/8554/', '', '编辑', 'merchantOrderReceiptUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8649);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8650, 8555, '/525/1366/8555/', '', '列表', 'merchantOrderProfitsharingLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8650);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8651, 8555, '/525/1366/8555/', '', '导出', 'merchantOrderProfitsharingExport', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8651);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8652, 8556, '/525/1366/8556/', '', '导出列表', 'merchantStoreExcelLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8652);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8653, 8556, '/525/1366/8556/', '', '导出下载', 'merchantStoreExcelDownload', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8653);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8654, 8557, '/525/1304/8557/', '', '申请', 'merchantApplymentsCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8654);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8655, 8557, '/525/1304/8557/', '', '详情', 'merchantApplymentsDetail', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8655);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8656, 8557, '/525/1304/8557/', '', '编辑', 'merchantApplymentsUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8656);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8657, 8557, '/525/1304/8557/', '', '上传图片', 'merchantApplymentsUpload', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8657);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8658, 8558, '1671/54/8558/', '', '列表', 'merchantAttachmentLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8658);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8659, 8558, '1671/54/8558/', '', '删除', 'merchantAttachmentDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8659);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8660, 8558, '1671/54/8558/', '', '批量修改', 'merchantAttachmentBatchChangeCategory', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8660);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8661, 8558, '1671/54/8558/', '', '编辑', 'merchantAttachmentUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8661);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8662, 8558, '1671/54/8558/', '', '分类列表', 'merchantAttachmentCategoryGetFormatList', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8662);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8663, 8558, '1671/54/8558/', '', '添加', 'merchantAttachmentCategoryCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8663);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8664, 8558, '1671/54/8558/', '', '编辑', 'merchantAttachmentCategoryUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8664);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8665, 8558, '1671/54/8558/', '', '删除', 'merchantAttachmentCategoryDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8665);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8666, 8559, '/1027/1103/8559/', '', '优惠券可用商品', 'merchantCouponProduct', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8666);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8667, 8559, '/1027/1103/8559/', '', '发送优惠券', 'merchantCouponSendCoupon', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8667);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8668, 8559, '/1027/1103/8559/', '', '列表', 'merchantUserLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8668);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8669, 8559, '/1027/1103/8559/', '', '修改标签', 'merchantUserChangeLabel', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8669);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8670, 8559, '/1027/1103/8559/', '', '订单列表', 'merchantUserOrder', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8670);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8671, 8559, '/1027/1103/8559/', '', '优惠券', 'merchantUserCoupon', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8671);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8672, 8560, '/107/1132/8560/', '', '发送优惠券记录', 'merchantCouponSendLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8672);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8673, 8561, '/106/107/115/8561/', '', '添加', 'merchantCouponCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8673);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8674, 8561, '/106/107/115/8561/', '', '修改状态', 'merchantCouponIssueChangeStatus', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8674);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8675, 8561, '/106/107/115/8561/', '', '列表', 'merchantCouponLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8675);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8676, 8561, '/106/107/115/8561/', '', '删除', 'merchantCouponDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8676);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8677, 8561, '/106/107/115/8561/', '', '详情', 'merchantCouponDetail', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8677);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8678, 8561, '/106/107/115/8561/', '', '编辑', 'systemCouponUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8678);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8679, 8562, '/106/107/115/8562/', '', '上传图片', 'merchantUploadImage', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8679);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8680, 8562, '/106/107/115/8562/', '', '图片列表', 'merchantAttachmentLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8680);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8681, 8563, '/106/107/122/8563/', '', '使用记录', 'merchantCouponIssue', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8681);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8682, 8564, '1649/1380/8564/', '', '列表', 'merchantServeExportLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8682);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8683, 8564, '1649/1380/8564/', '', '月结账号编辑', 'merchantExpressPratnerUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8683);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8684, 8565, '9332/1650/8565/', '', '获取分类', 'merchantStoreDeliveryBusiness', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8684);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8685, 8565, '9332/1650/8565/', '', '添加', 'merchantStoreDeliveryCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8685);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8686, 8565, '9332/1650/8565/', '', '编辑', 'merchantStoreDeliveryUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8686);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8687, 8565, '9332/1650/8565/', '', '编辑状态', 'merchantStoreDeliveryStatus', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8687);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8688, 8565, '9332/1650/8565/', '', '列表', 'merchantStoreDeliveryLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8688);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8689, 8565, '9332/1650/8565/', '', '详情', 'merchantStoreDeliveryDetail', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8689);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8690, 8565, '9332/1650/8565/', '', '删除', 'merchantStoreDeliveryDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8690);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8691, 8565, '9332/1650/8565/', '', '备注', 'merchantStoreDeliveryMark', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8691);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8692, 8565, '9332/1650/8565/', '', '城市列表', 'merchantStoreDeliveryCityList', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8692);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8693, 8566, '9332/1652/8566/', '', '充值记录', 'merchantStoreDeliveryPayLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8693);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8694, 8567, '9332/1651/8567/', '', '列表', 'merchantStoreDeliveryOrderLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8694);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8695, 8567, '9332/1651/8567/', '', '取消', 'merchantStoreDeliveryOrderCancel', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8695);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8696, 8567, '9332/1651/8567/', '', '详情', 'merchantStoreDeliveryOrderDetail', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8696);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8697, 8568, '1649/700/8568/', '', '添加 ', 'merchantStoreShippingTemplateCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8697);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8698, 8568, '1649/700/8568/', '', '编辑', 'merchantStoreShippingTemplateUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8698);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8699, 8568, '1649/700/8568/', '', '详情', 'merchantStoreShippingTemplateDetail', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8699);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8700, 8568, '1649/700/8568/', '', '删除', 'merchantStoreShippingTemplateDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8700);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8701, 8569, '/1671/1672/8569/', '', '列表 ', 'merchantDiyLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8701);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8702, 8569, '/1671/1672/8569/', '', '添加/编辑', 'merchantDiyCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8702);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8703, 8569, '/1671/1672/8569/', '', '使用模板', 'merchantDiyStatus', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8703);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8704, 8569, '/1671/1672/8569/', '', '设置模版默认数据', 'merchantDiySetDefault', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8704);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8705, 8569, '/1671/1672/8569/', '', '重置模板', 'merchantDiyRecovery', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8705);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8706, 8569, '/1671/1672/8569/', '', '当前使用模板', 'merchantDiyInfo', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8706);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8707, 8569, '/1671/1672/8569/', '', '删除', 'merchantDiyDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8707);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8708, 8569, '/1671/1672/8569/', '', '店铺街装修', 'merchantDiyProductLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8708);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8709, 8569, '/1671/1672/8569/', '', '复制', 'merchantDiyCopy', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8709);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8710, 8570, '/1671/1672/8570/', '', '上传图片', 'uploadImage', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8710);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8711, 8570, '/1671/1672/8570/', '', '图片列表', 'systemAttachmentLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8711);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8720, 8573, '/106/1471/8573/', '', '数据详情', 'merchantGroupDetail', '95', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8720);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8721, 8573, '/106/1471/8573/', '', '数据列表', 'merchantGroupDataLst', '95', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8721);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8722, 8573, '/106/1471/8573/', '', '数据添加', 'merchantGroupDataCreate', '95', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8722);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8723, 8573, '/106/1471/8573/', '', '数据编辑', 'merchantGroupDataUpdate', '95', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8723);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8724, 8573, '/106/1471/8573/', '', '数据删除', 'merchantGroupDataDelete', '95', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8724);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8725, 8573, '/106/1471/8573/', '', '数据修改状态', 'merchantGroupDataChangeStatus', '95', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8725);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8726, 8574, '/106/1471/8574/', '', '上传图片', 'uploadImage', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8726);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8727, 8574, '/106/1471/8574/', '', '图片列表', 'systemAttachmentLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8727);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8728, 8575, '9290/8575/', '', '统计', 'merchantStoreSeckillProductLstFilter', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8728);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8729, 8575, '9290/8575/', '', '列表', 'merchantStoreSeckillProductLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8729);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8730, 8575, '9290/8575/', '', '添加 ', 'merchantStoreSeckillProductCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8730);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8731, 8575, '9290/8575/', '', '详情', 'merchantStoreSeckillProductDetail', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8731);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8732, 8575, '9290/8575/', '', '编辑', 'merchantStoreSeckillProductUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8732);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8733, 8575, '9290/8575/', '', '删除', 'merchantStoreSeckillProductDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8733);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8734, 8575, '9290/8575/', '', '彻底删除', 'merchantStoreSeckillProductDestory', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8734);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8735, 8575, '9290/8575/', '', '恢复', 'merchantStoreSeckillProductRestore', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8735);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8736, 8575, '9290/8575/', '', '修改状态', 'merchantStoreSeckillProductSwitchStatus', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8736);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8737, 8575, '9290/8575/', '', '排序', 'merchantStoreSeckillProductUpdateSort', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8737);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8738, 8575, '9290/8575/', '', '预览', 'merchantStoreSeckillProductPreview', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8738);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8739, 8575, '9290/8575/', '', '设置标签', 'merchantStoreSeckillProductLabels', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8739);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8740, 8576, '9289/8576/', '', '上传图片', 'merchantUploadImage', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8740);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8741, 8576, '9289/8576/', '', '图片列表', 'merchantAttachmentLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8741);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8742, 8577, '/106/1025/8577/', '', '列表', 'merchantStoreProductPresellLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8742);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8743, 8577, '/106/1025/8577/', '', '添加', 'merchantStoreProductPresellCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8743);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8744, 8577, '/106/1025/8577/', '', '详情', 'merchantStoreProductPresellDetail', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8744);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8745, 8577, '/106/1025/8577/', '', '编辑', 'merchantStoreProductPresellUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8745);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8746, 8577, '/106/1025/8577/', '', '删除', 'merchantStoreProductPresellDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8746);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8747, 8577, '/106/1025/8577/', '', '修改状态', 'merchantStoreProductPresellStatus', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8747);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8748, 8577, '/106/1025/8577/', '', '排序', 'merchantStoreProductPresellUpdateSort', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8748);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8749, 8577, '/106/1025/8577/', '', '预览', 'merchantStoreProductPresellPreview', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8749);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8750, 8577, '/106/1025/8577/', '', '设置标签', 'merchantStoreProductPreselltLabels', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8750);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8751, 8578, '/106/1025/8578/', '', '上传图片', 'merchantUploadImage', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8751);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8752, 8578, '/106/1025/8578/', '', '图片列表', 'merchantAttachmentLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8752);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8753, 8579, '/1099/1100/8579/', '', '列表 ', 'merchantStoreProductAssistLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8753);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8754, 8579, '/1099/1100/8579/', '', '添加', 'merchantStoreProductAssistCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8754);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8755, 8579, '/1099/1100/8579/', '', '详情', 'merchantStoreProductAssistDetail', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8755);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8756, 8579, '/1099/1100/8579/', '', '编辑', 'merchantStoreProductAssistUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8756);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8757, 8579, '/1099/1100/8579/', '', '删除', 'merchantStoreProductAssistDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8757);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8758, 8579, '/1099/1100/8579/', '', '修改状态', 'merchantStoreProductAssistStatus', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8758);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8759, 8579, '/1099/1100/8579/', '', '排序', 'merchantStoreProductAssistUpdateSort', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8759);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8760, 8579, '/1099/1100/8579/', '', '预览', 'merchantStoreProductAssistPreview', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8760);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8761, 8579, '/1099/1100/8579/', '', '设置标签', 'merchantStoreProductAssistLabels', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8761);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8762, 8580, '/1099/1100/8580/', '', '上传图片', 'merchantUploadImage', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8762);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8763, 8580, '/1099/1100/8580/', '', '图片列表', 'merchantAttachmentLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8763);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8764, 8581, '/1099/1101/8581/', '', '活动列表', 'merchantStoreProductAssistSetLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8764);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8765, 8581, '/1099/1101/8581/', '', '活动详情', 'merchantStoreProductAssistSetDetail', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8765);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8766, 8582, '/1139/1140/8582/', '', '列表', 'merchantStoreProductGroupLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8766);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8767, 8582, '/1139/1140/8582/', '', '添加', 'merchantStoreProductGroupCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8767);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8768, 8582, '/1139/1140/8582/', '', '详情', 'merchantStoreProductGroupDetail', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8768);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8769, 8582, '/1139/1140/8582/', '', '编辑', 'merchantStoreProductGroupUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8769);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8770, 8582, '/1139/1140/8582/', '', '删除', 'merchantStoreProductGroupDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8770);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8771, 8582, '/1139/1140/8582/', '', '修改状态', 'merchantStoreProductGroupStatus', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8771);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8772, 8582, '/1139/1140/8582/', '', '排序', 'merchantStoreProductGroupSort', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8772);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8773, 8582, '/1139/1140/8582/', '', '预览', 'merchantStoreProductGroupPreview', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8773);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8774, 8582, '/1139/1140/8582/', '', '设置标签', 'merchantStoreProductGroupLabels', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8774);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8775, 8582, '/1139/1140/8582/', '', '拼团配置', 'merchantConfigGroupBuying', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8775);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8776, 8583, '/1139/1140/8583/', '', '上传图片', 'merchantUploadImage', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8776);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8777, 8583, '/1139/1140/8583/', '', '图片列表', 'merchantAttachmentLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8777);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8778, 8584, '/1139/1141/8584/', '', '活动列表 ', 'merchantStoreProductGroupBuyingLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8778);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8779, 8584, '/1139/1141/8584/', '', '活动详情', 'merchantStoreProductGroupBuyingDetail', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8779);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8780, 8585, '/1139/1141/8585/', '', '上传图片', 'merchantUploadImage', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8780);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8781, 8585, '/1139/1141/8585/', '', '图片列表', 'merchantAttachmentLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8781);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8782, 8586, '/106/785/786/8586/', '', '列表 ', 'merchantBroadcastRoomLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8782);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8783, 8586, '/106/785/786/8586/', '', '详情', 'merchantBroadcastRoomDetail', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8783);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8784, 8586, '/106/785/786/8586/', '', '添加', 'merchantBroadcastRoomCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8784);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8785, 8586, '/106/785/786/8586/', '', '编辑', 'merchantBroadcastRoomUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8785);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8786, 8586, '/106/785/786/8586/', '', '修改状态', 'merchantBroadcastRoomChangeStatus', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8786);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8787, 8586, '/106/785/786/8586/', '', '导入商品', 'merchantBroadcastRoomExportGoods', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8787);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8788, 8586, '/106/785/786/8586/', '', '删除商品', 'merchantBroadcastRoomRmExportGoods', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8788);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8789, 8586, '/106/785/786/8586/', '', '备注', 'merchantBroadcastRoomMark', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8789);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8790, 8586, '/106/785/786/8586/', '', '商品详情', 'merchantBroadcastRoomGoods', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8790);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8791, 8586, '/106/785/786/8586/', '', '关闭客服', 'merchantBroadcastRoomCloseKf', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8791);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8792, 8586, '/106/785/786/8586/', '', '禁言', 'merchantBroadcastRoomCloseComment', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8792);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8793, 8586, '/106/785/786/8586/', '', '收录', 'merchantBroadcastRoomCloseFeeds', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8793);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8794, 8586, '/106/785/786/8586/', '', '商品上下架', 'merchantBroadcastOnSale', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8794);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8795, 8586, '/106/785/786/8586/', '', '删除', 'merchantBroadcastRoomDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8795);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8796, 8586, '/106/785/786/8586/', '', '添加 客服', 'merchantBroadcastAddAssistant', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8796);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8797, 8586, '/106/785/786/8586/', '', '消息推送', 'merchantBroadcastPushMessage', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8797);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8798, 8586, '/106/785/786/8586/', '', '列表', 'merchantBroadcastGoodsLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8798);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8799, 8586, '/106/785/786/8586/', '', '详情', 'merchantBroadcastGoodsDetail', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8799);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8800, 8586, '/106/785/786/8586/', '', '添加', 'merchantBroadcastGoodsCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8800);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8801, 8586, '/106/785/786/8586/', '', '编辑', 'merchantBroadcastGoodsUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8801);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8802, 8586, '/106/785/786/8586/', '', '修改状态', 'merchantBroadcastGoodsChangeStatus', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8802);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8803, 8586, '/106/785/786/8586/', '', '备注', 'merchantBroadcastGoodsMark', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8803);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8804, 8586, '/106/785/786/8586/', '', '删除', 'merchantBroadcastGoodsDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8804);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8805, 8586, '/106/785/786/8586/', '', '批量添加', 'merchantBroadcastGoodsbatchCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8805);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8806, 8587, '/106/785/786/8587/', '', '上传图片', 'merchantUploadImage', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8806);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8807, 8587, '/106/785/786/8587/', '', '图片列表', 'merchantAttachmentLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8807);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8808, 8588, '/785/1594/8588/', '', '列表', 'merchantBroadcastAssistantLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8808);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8809, 8588, '/785/1594/8588/', '', '添加', 'merchantBroadcastAssistantCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8809);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8810, 8588, '/785/1594/8588/', '', '编辑', 'merchantBroadcastAssistantUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8810);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8811, 8588, '/785/1594/8588/', '', '备注', 'merchantBroadcastAssistantMark', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8811);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8812, 8588, '/785/1594/8588/', '', '删除', 'merchantBroadcastAssistantDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8812);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8813, 8589, '/785/1594/8589/', '', '上传图片', 'merchantUploadImage', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8813);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8814, 8589, '/785/1594/8589/', '', '图片列表', 'merchantAttachmentLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8814);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8815, 8590, '/1293/1295/8590/', '', '列表', 'merchantIntegralList', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8815);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8816, 8590, '/1293/1295/8590/', '', '统计', 'merchantIntegralTitle', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8816);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8817, 8591, '/1293/1295/8591/', '', '配置获取', 'merchantConfigForm', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8817);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8818, 8591, '/1293/1295/8591/', '', '配置保存', 'merchantConfigSave', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8818);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8819, 8592, '/512/513/8592/', '', '预览', 'merchantServeExportTemps', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8819);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8820, 8592, '/512/513/8592/', '', '默认模板', 'merchantServeExportDumpLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8820);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8821, 8592, '/512/513/8592/', '', '导出', 'merchantStoreOrderExcel', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8821);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8822, 8592, '/512/513/8592/', '', '打印小票', 'merchantStoreOrderPrinter', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8822);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8823, 8592, '/512/513/8592/', '', '统计', 'merchantStoreOrderTitle', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8823);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8824, 8592, '/512/513/8592/', '', '列表', 'merchantStoreOrderLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8824);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8825, 8592, '/512/513/8592/', '', '快递查询', 'merchantStoreOrderExpress', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8825);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8826, 8592, '/512/513/8592/', '', '发货', 'merchantStoreOrderDelivery', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8826);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8827, 8592, '/512/513/8592/', '', '批量发货', 'merchantStoreOrderBatchDelivery', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8827);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8828, 8592, '/512/513/8592/', '', '导出发货单', 'merchantStoreOrderDeliveryExport', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8828);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8829, 8592, '/512/513/8592/', '', '头部统计', 'merchantStoreOrderStat', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8829);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8830, 8592, '/512/513/8592/', '', '编辑', 'merchantStoreOrderUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8830);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8831, 8592, '/512/513/8592/', '', '详情', 'merchantStoreOrderDetail', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8831);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8832, 8592, '/512/513/8592/', '', '操作记录', 'merchantStoreOrderLog', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8832);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8833, 8592, '/512/513/8592/', '', '备注', 'merchantStoreOrderRemark', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8833);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8834, 8592, '/512/513/8592/', '', '核销', 'merchantStoreOrderVerify', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8834);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8835, 8592, '/512/513/8592/', '', '删除', 'merchantStoreOrderDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8835);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8836, 8592, '/512/513/8592/', '', '导入', 'merchantStoreOrderDeliveryImport', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8836);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8837, 8592, '/512/513/8592/', '', '导入记录', 'merchantStoreOrderDeliveryImportLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8837);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8838, 8592, '/512/513/8592/', '', '详情', 'merchantStoreOrderDeliveryImportDetail', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8838);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8839, 8592, '/512/513/8592/', '', '导出发货记录', 'merchantStoreOrderDeliveryImportExcel', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8839);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8840, 8593, '/512/513/8593/', '', '导出列表', 'merchantStoreExcelLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8840);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8841, 8593, '/512/513/8593/', '', '导出下载', 'merchantStoreExcelDownload', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8841);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8842, 8594, '/512/789/8594/', '', '统计', 'merchantStoreOrderTakeTitle', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8842);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8843, 8594, '/512/789/8594/', '', '列表', 'merchantStoreTakeOrderLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8843);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8844, 8595, '/512/528/8595/', '', '列表', 'merchantStoreRefundOrderLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8844);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8845, 8595, '/512/528/8595/', '', '详情', 'merchantStoreRefundOrderDetail', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8845);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8846, 8595, '/512/528/8595/', '', '审核', 'merchantStoreRefundOrderSwitchStatus', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8846);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8847, 8595, '/512/528/8595/', '', '收到退回商品后确认退款', 'merchantStoreRefundOrderRefund', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8847);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8848, 8595, '/512/528/8595/', '', '删除', 'merchantStoreRefundDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8848);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8849, 8595, '/512/528/8595/', '', '备注', 'merchantStoreRefundMark', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8849);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8850, 8595, '/512/528/8595/', '', '操作记录', 'merchantStoreRefundLog', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8850);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8851, 8595, '/512/528/8595/', '', '快递查询', 'merchantStoreRefundExpress', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8851);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8852, 8595, '/512/528/8595/', '', '导出', 'merchantStoreRefundCreateExcel', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8852);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8853, 8596, '/512/528/8596/', '', '导出列表', 'merchantStoreExcelLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8853);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8854, 8596, '/512/528/8596/', '', '导出下载', 'merchantStoreExcelDownload', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8854);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8855, 8597, '/95/99/8597/', '', '列表', 'merchantStoreAttrTemplateLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8855);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8856, 8597, '/95/99/8597/', '', '添加 ', 'merchantStoreAttrTemplateCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8856);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8857, 8597, '/95/99/8597/', '', '删除', 'merchantStoreAttrTemplateDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8857);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8858, 8597, '/95/99/8597/', '', '文件类型', 'merchantStoreAttrTemplateUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8858);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8859, 8598, '/95/96/8598/', '', '编辑', 'merchantStoreCategoryUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8859);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8860, 8598, '/95/96/8598/', '', '列表', 'merchantStoreCategoryLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8860);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8861, 8598, '/95/96/8598/', '', '详情', 'merchantStoreCategoryDtailt', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8861);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8862, 8598, '/95/96/8598/', '', '添加', 'merchantStoreCategoryCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8862);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8863, 8598, '/95/96/8598/', '', '删除', 'merchantStoreCategoryDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8863);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8864, 8598, '/95/96/8598/', '', '修改状态', 'merchantStoreCategorySwitchStatus', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8864);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8865, 8599, '/95/96/8599/', '', '上传图片', 'merchantUploadImage', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8865);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8866, 8599, '/95/96/8599/', '', '图片列表', 'merchantAttachmentLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8866);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8867, 8600, '/95/105/8600/', '', '头部统计', 'merchantStoreProductLstFilter', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8867);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8868, 8600, '/95/105/8600/', '', '列表', 'merchantStoreProductLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8868);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8869, 8600, '/95/105/8600/', '', '添加', 'merchantStoreProductCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8869);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8870, 8600, '/95/105/8600/', '', '详情', 'merchantStoreProductDetail', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8870);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8871, 8600, '/95/105/8600/', '', '上传视频配置', 'merchantStoreProductTempKey', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8871);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8872, 8600, '/95/105/8600/', '', '编辑', 'merchantStoreProductUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8872);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8873, 8600, '/95/105/8600/', '', '删除', 'merchantStoreProductDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8873);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8874, 8600, '/95/105/8600/', '', '加入回收站', 'merchantStoreProductDestory', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8874);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8875, 8600, '/95/105/8600/', '', '恢复', 'merchantStoreProductRestore', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8875);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8876, 8600, '/95/105/8600/', '', '上下架', 'merchantStoreProductSwitchStatus', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8876);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8877, 8600, '/95/105/8600/', '', '排序', 'merchantStoreProductUpdateSort', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8877);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8878, 8600, '/95/105/8600/', '', '预览', 'merchantStoreProductPreview', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8878);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8879, 8600, '/95/105/8600/', '', '标签', 'merchantStoreProductLabels', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8879);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8880, 8600, '/95/105/8600/', '', '获取规格', 'merchantStoreProductAttrValue', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8880);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8881, 8600, '/95/105/8600/', '', '列表', 'merchantStoreProductCopyLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8881);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8882, 8600, '/95/105/8600/', '', '获取信息', 'merchantStoreProductCopyGet', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8882);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8883, 8600, '/95/105/8600/', '', '统计', 'merchantStoreProductCopyCount', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8883);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8884, 8600, '/95/105/8600/', '', '保存', 'merchantStoreProductCopySave', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8884);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8885, 8601, '/95/105/8601/', '', '上传图片', 'merchantUploadImage', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8885);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8886, 8601, '/95/105/8601/', '', '图片列表', 'merchantAttachmentLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8886);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8887, 8602, '512/544/8602/', '', '列表', 'merchantProductReplyLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8887);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8888, 8602, '512/544/8602/', '', '回复表单', 'merchantProductReplyForm', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8888);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8889, 8602, '512/544/8602/', '', '回复', 'merchantProductReplyReply', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8889);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8890, 8602, '512/544/8602/', '', '排序', 'merchantProductReplySort', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8890);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8891, 8603, '512/544/8603/', '', '上传图片', 'merchantUploadImage', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8891);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8892, 8603, '512/544/8603/', '', '图片列表', 'merchantAttachmentLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8892);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8893, 8604, '/95/1468/8604/', '', '列表', 'merchantStoreProductLabelLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8893);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8894, 8604, '/95/1468/8604/', '', '添加', 'merchantStoreProductLabelCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8894);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8895, 8604, '/95/1468/8604/', '', '编辑', 'merchantStoreProductLabelUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8895);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8896, 8604, '/95/1468/8604/', '', '详情', 'merchantStoreProductLabelDetail', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8896);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8897, 8604, '/95/1468/8604/', '', '删除', 'merchantStoreProductLabelDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8897);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8898, 8604, '/95/1468/8604/', '', '修改状态', 'merchantStoreProductLabelStatus', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8898);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8899, 8605, '/106/1630/8605/', '', '添加', 'merchantStoreDiscountsCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8899);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8900, 8605, '/106/1630/8605/', '', '编辑', 'merchantStoreDiscountsUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8900);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8901, 8605, '/106/1630/8605/', '', '列表', 'merchantStoreDiscountsLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8901);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8902, 8605, '/106/1630/8605/', '', '详情', 'merchantStoreDiscountsDetail', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8902);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8903, 8605, '/106/1630/8605/', '', '删除', 'merchantStoreDiscountsDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8903);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8904, 8605, '/106/1630/8605/', '', '修改状态', 'merchantStoreDiscountsStatus', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8904);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8905, 8606, '/106/1630/8606/', '', '上传图片', 'merchantUploadImage', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8905);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8906, 8606, '/106/1630/8606/', '', '图片列表', 'merchantAttachmentLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8906);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8907, 8607, '/95/1246/8607/', '', '列表', 'merchantGuaranteeLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8907);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8908, 8607, '/95/1246/8607/', '', '添加', 'smerchantGuaranteeCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8908);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8909, 8607, '/95/1246/8607/', '', '编辑', 'merchantGuaranteeUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8909);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8910, 8607, '/95/1246/8607/', '', '详情', 'merchantGuaranteeDetail', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8910);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8911, 8607, '/95/1246/8607/', '', '删除', 'merchantGuaranteeDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8911);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8912, 8607, '/95/1246/8607/', '', '排序', 'merchantGuaranteeSort', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8912);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8913, 8607, '/95/1246/8607/', '', '修改状态', 'merchantGuaranteeStatus', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8913);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8914, 8608, '/526/49/51/8608/', '', '列表', 'merchantAdminLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8914);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8915, 8608, '/526/49/51/8608/', '', '修改状态', 'merchantAdminStatus', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8915);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8916, 8608, '/526/49/51/8608/', '', '添加', 'merchantAdminCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8916);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8917, 8608, '/526/49/51/8608/', '', '编辑', 'merchantAdminUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8917);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8918, 8608, '/526/49/51/8608/', '', '修改密码', 'merchantAdminPassword', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8918);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8919, 8608, '/526/49/51/8608/', '', '删除', 'merchantAdminDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8919);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8920, 8609, '/8609/', '', '修改信息', 'merchantAdminEdit', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8920);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8921, 8609, '/8609/', '', '修改密码', 'merchantAdminEditPassword', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8921);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8922, 8610, '9410/120/8610/', '', '列表', 'merchantServiceLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8922);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8923, 8610, '9410/120/8610/', '', '登录', 'merchantServiceLogin', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8923);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8924, 8610, '9410/120/8610/', '', '添加', 'merchantServiceCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8924);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8925, 8610, '9410/120/8610/', '', '编辑', 'merchantServiceUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8925);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8926, 8610, '9410/120/8610/', '', '修改状态', 'merchantServiceSwitchStatus', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8926);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8927, 8610, '9410/120/8610/', '', '删除', 'merchantServiceDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8927);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8928, 8610, '9410/120/8610/', '', '客服的全部用户 ', 'merchantServiceServiceUserList', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8928);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8929, 8610, '9410/120/8610/', '', '用户与客服聊天记录', 'merchantServiceServiceUserLogLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8929);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8930, 8610, '9410/120/8610/', '', '客服的聊天用户列表', 'merchantServiceServiceMerchantUserList', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8930);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8931, 8610, '9410/120/8610/', '', '用户与商户聊天记录', 'merchantServiceMerchantUserLogLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8931);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8932, 8611, '9410/5123/8611/', '', '列表', 'merchantServiceReplyLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8932);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8933, 8611, '9410/5123/8611/', '', '添加', 'merchantServiceReplyCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8933);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8934, 8611, '9410/5123/8611/', '', '编辑', 'merchantServiceReplyUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8934);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8935, 8611, '9410/5123/8611/', '', '切换状态', 'merchantServiceReplyStatus', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8935);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8936, 8611, '9410/5123/8611/', '', '删除', 'merchantServiceReplyDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8936);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8940, 8613, '9358/1656/8613/', '', '列表', 'merchantStorePrinterLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8940);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8941, 8613, '9358/1656/8613/', '', '添加', 'merchantStorePrinterCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8941);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8942, 8613, '9358/1656/8613/', '', '编辑', 'merchantStorePrinterUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8942);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8943, 8613, '9358/1656/8613/', '', '取消', 'merchantStorePrinterStatus', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8943);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8944, 8613, '9358/1656/8613/', '', '删除', 'merchantStorePrinterDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8944);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8945, 8614, '/55/8614/', '', '所有数据', 'merchantStatisticsMain', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8945);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8946, 8614, '/55/8614/', '', '支付订单', 'merchantStatisticsOrder', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8946);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8947, 8614, '/55/8614/', '', '成交客户', 'merchantStatisticsUser', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8947);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8948, 8614, '/55/8614/', '', '成交客户比', 'merchantStatisticsUserRate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8948);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8949, 8614, '/55/8614/', '', '商品支付排行', 'merchantStatisticsProduct', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8949);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8950, 8614, '/55/8614/', '', '商品访问排行', 'merchantStatisticsProductVisit', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8950);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8951, 8614, '/55/8614/', '', '商品加购排行', 'merchantStatisticsProductCart', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8951);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8952, 8615, '/1119/8615/', '', '列表', 'systemNoticeLogList', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8952);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8953, 8615, '/1119/8615/', '', '已读', 'systemNoticeLogRead', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8953);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8954, 8615, '/1119/8615/', '', '删除', 'systemNoticeLogDel', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8954);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8955, 8615, '/1119/8615/', '', '未读统计', 'systemNoticeLogUnreadCount', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8955);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8956, 8616, '/526/74/8616/', '', '资料更新', 'merchantUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8956);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8957, 8617, '/526/546/8617/', '', '保存到店自提信息', 'merchantTakeUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8957);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8958, 8618, '/1027/1286/8618/', '', '搜索记录', 'merchantUserSearchLog', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8958);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8959, 8619, '/1028/1030/8619/', '', '列表', 'merchantLabelRuleLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8959);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8960, 8619, '/1028/1030/8619/', '', '添加', 'merchantLabelRuleCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8960);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8961, 8619, '/1028/1030/8619/', '', '编辑', 'merchantLabelRuleUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8961);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8962, 8619, '/1028/1030/8619/', '', '删除', 'merchantLabelRuleDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8962);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8963, 8619, '/1028/1030/8619/', '', '自动同步', 'merchantLabelRuleSync', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8963);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8964, 8620, '/1028/1029/8620/', '', '列表', 'merchantUserLabelLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8964);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8965, 8620, '/1028/1029/8620/', '', '添加', 'merchantUserLabelCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8965);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8966, 8620, '/1028/1029/8620/', '', '删除', 'merchantUserLabelDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8966);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8967, 8620, '/1028/1029/8620/', '', '编辑', 'merchantUserLabelUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8967);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8968, 8621, '526/1376/1378/8621/', '', '套餐列表', 'merchantServeMeal', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8968);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8969, 8621, '526/1376/1378/8621/', '', '支付二维码', 'merchantServeCode', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8969);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8970, 8621, '526/1376/1378/8621/', '', '购买记录', 'merchantServeLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8970);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8971, 8621, '526/1376/1378/8621/', '', '详情', 'merchantServeDetail', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8971);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8972, 8621, '526/1376/1378/8621/', '', '账号信息', 'merchantServeInfo', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8972);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8973, 8622, '526/1376/1377/8622/', '', '保存配置', 'merchantServeSetConfig', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8973);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8974, 1508, '/101/9047/1508/', '', '权限', '/systemForm/Basics/members', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8974);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8975, 1542, '/9361/1538/1542/', '', '权限', '/systemForm/Basics/community', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8975);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8976, 1628, '/110/670/1628/', '', '权限', '/systemForm/Basics/app_version', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8976);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8977, 20034, '/110/118/34/', '', '权限', '/systemForm/Basics/system_tabs', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8977);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8978, 5119, '/110/9360/5119/', '', '权限', '/systemForm/Basics/shop_tabs', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8978);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8979, 5120, '/110/9360/5120/', '', '权限', '/systemForm/Basics/pay_tabs', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8979);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8980, 5121, '/110/118/5121/', '', '权限', '/systemForm/Basics/extend_tabs', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8980);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8981, 5127, '/5054/5127/', '', '权限', '/systemForm/Basics/service', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8981);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8982, 1631, '/110/1233/1235/1631/', '', '权限', '/group/config/98', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8982);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8983, 1631, '/110/1233/1235/1631/', '', '附加权限', 'append_/group/config/98', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8983);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8984, 666, '/520/116/666/', '', '权限', '/maintain/auth', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8984);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8985, 8974, '/101/9047/1508/8974/', '', '编辑配置信息', 'configSave', 'members', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8985);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8986, 8975, '/9361/1538/1542/8975/', '', '编辑配置信息', 'configSave', 'community', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8986);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8987, 8976, '/110/670/1628/8976/', '', '编辑配置信息', 'configSave', 'app_version', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8987);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8988, 8977, '/110/118/34/8977/', '', '编辑配置信息', 'configSave', 'system_tabs', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8988);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8989, 8978, '/110/9360/5119/8978/', '', '编辑配置信息', 'configSave', 'shop_tabs', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8989);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8990, 8979, '/110/9360/5120/8979/', '', '编辑配置信息', 'configSave', 'pay_tabs', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8990);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8991, 8980, '/110/118/5121/8980/', '', '编辑配置信息', 'configSave', 'extend_tabs', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8991);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8992, 8981, '/5054/5127/8981/', '', '编辑配置信息', 'configSave', 'service', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8992);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8993, 8982, '/110/1233/1235/1631/8982/', '', '详情', 'groupDetail', '98', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8993);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8994, 8982, '/110/1233/1235/1631/8982/', '', '列表', 'groupDataLst', '98', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8994);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8995, 8982, '/110/1233/1235/1631/8982/', '', '添加', 'groupDataCreate', '98', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8995);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8996, 8982, '/110/1233/1235/1631/8982/', '', '编辑', 'groupDataUpdate', '98', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8996);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8997, 8982, '/110/1233/1235/1631/8982/', '', '删除', 'groupDataDelete', '98', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8997);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8998, 8982, '/110/1233/1235/1631/8982/', '', '修改状态', 'groupDataChangeStatus', '98', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8998);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 8999, 8983, '/110/1233/1235/1631/8983/', '', '上传图片', 'uploadImage', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 8999);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9000, 8983, '/110/1233/1235/1631/8983/', '', '图片列表', 'systemAttachmentLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9000);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9001, 7768, '/110/5125/1597/7768/', '', '消息配置修改模板ID', 'systemNoticeConfigGetChangeTempId', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9001);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9002, 7768, '/110/5125/1597/7768/', '', '消息配置修改模板ID', 'systemNoticeConfigSetChangeTempId', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9002);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9003, 8984, '/520/116/666/8984/', '', '获取去版权信息', 'systemCopyright', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9003);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9004, 8984, '/520/116/666/8984/', '', '获取授权信息', 'systemAuthCopyright', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9004);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9005, 8984, '/520/116/666/8984/', '', '保存去版权信息', 'systemSaveCopyright', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9005);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9006, 1665, '/1665/', '', '微页面', '/setting/micro/list', '[]', 60, 0, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9006);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9007, 719, '/719/', '', '活动氛围图', '/marketing/atmosphere/list', '[]', 10, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9007);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9008, 719, '/719/', '', '活动边框图', '/marketing/border/list', '[]', 10, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9008);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9009, 9006, '/1665/9006/', '', '权限', '/setting/micro/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9009);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9010, 9006, '/1665/9006/', '', '附加权限', 'append_/setting/micro/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9010);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9011, 9007, '/719/9007/', '', '权限', '/marketing/atmosphere/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9011);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9012, 9008, '/719/9008/', '', '权限', '/marketing/border/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9012);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9013, 9009, '/1665/9006/9009/', '', '列表 ', 'systemDiyMicroLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9013);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9014, 9009, '/1665/9006/9009/', '', '详情 ', 'systemDiyMicroDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9014);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9015, 9009, '/1665/9006/9009/', '', '添加/编辑', 'systemDiyMicroCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9015);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9016, 9009, '/1665/9006/9009/', '', '重置', 'systemDiyMicroRecovery', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9016);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9017, 9009, '/1665/9006/9009/', '', '删除', 'systemDiyMicroDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9017);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9018, 9010, '/1665/9006/9010/', '', '上传图片', 'uploadImage', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9018);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9019, 9010, '/1665/9006/9010/', '', '图片列表', 'systemAttachmentLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9019);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9020, 9011, '/719/9007/9011/', '', '添加', 'systemActivityAtmosphereCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9020);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9021, 9011, '/719/9007/9011/', '', '列表', 'systemActivityAtmosphereLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9021);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9022, 9011, '/719/9007/9011/', '', '编辑', 'systemActivityAtmosphereUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9022);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9023, 9011, '/719/9007/9011/', '', '详情', 'systemActivityAtmosphereDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9023);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9024, 9011, '/719/9007/9011/', '', '删除', 'systemActivityAtmosphereDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9024);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9025, 9011, '/719/9007/9011/', '', '修改状态', 'systemActivityAtmosphereStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9025);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9026, 9012, '/719/9008/9012/', '', '添加', 'systemActivityBorderCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9026);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9027, 9012, '/719/9008/9012/', '', '列表', 'systemActivityBorderLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9027);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9028, 9012, '/719/9008/9012/', '', '编辑', 'systemActivityBorderUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9028);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9029, 9012, '/719/9008/9012/', '', '详情', 'systemActivityBorderDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9029);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9030, 9012, '/719/9008/9012/', '', '删除', 'systemActivityBorderDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9030);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9031, 9012, '/719/9008/9012/', '', '修改状态', 'systemActivityBorderStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9031);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9032, 7778, '/87/539/7778/', '', '批量上下架', 'systemStoreProductSwitchBatchStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9032);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9033, 7778, '/87/539/7778/', '', '批量设置标签', 'systemStoreProductSwitchBatchLabels', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9033);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9034, 7778, '/87/539/7778/', '', '批量设置推荐', 'systemStoreProductSwitchBatchHot', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9034);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9035, 8600, '/95/105/8600/', '', '免审编辑', 'merchantStoreProductFreeTrial', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9035);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9036, 8600, '/95/105/8600/', '', '批量上下架', 'merchantStoreProductSwitchBatchStatus', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9036);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9037, 8600, '/95/105/8600/', '', '批量设置运费模板', 'merchantStoreProductSwitchBatchTemplate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9037);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9038, 8600, '/95/105/8600/', '', '批量设置标签', 'merchantStoreProductSwitchBatchLabels', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9038);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9039, 8600, '/95/105/8600/', '', '批量设置推荐', 'merchantStoreProductSwitchBatchHot', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9039);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9040, 9007, '/719/9007/', '', '附加权限', 'append_/marketing/atmosphere/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9040);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9041, 9008, '/719/9008/', '', '附加权限', 'append_/marketing/border/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9041);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9042, 9040, '/719/9007/9040/', '', '上传图片', 'uploadImage', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9042);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9043, 9040, '/719/9007/9040/', '', '图片列表', 'systemAttachmentLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9043);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9044, 9041, '/719/9008/9041/', '', '上传图片', 'uploadImage', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9044);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9045, 9041, '/719/9008/9041/', '', '图片列表', 'systemAttachmentLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9045);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9046, 20101, '/101/', '', '付费会员', '/user/svip', '[]', -1, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9046);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9047, 20101, '/101/', '', '用户等级', '/user/member', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9047);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9048, 9046, '/101/9046/', '', '会员协议', '/user/member/vipAgreement', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9048);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9049, 9046, '/101/9046/', '', '会员权益', '/user/member/equity', '[]', 80, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9049);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9050, 9046, '/101/9046/', '', '付费会员配置', '/systemForm/Basics/svip', '[]', 100, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9050);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9051, 9046, '/101/9046/', '', '会员类型', '/user/member/type', '[]', 90, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9051);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9052, 20087, '/87/', '', '商品参数', '/product/specsMain', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9052);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9053, 526, '526/', '', '付费会员', '/systemForm/Basics/svip', '[]', 97, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9053);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9054, 9046, '/101/9046/', '', '会员记录', '/user/member/record', '[]', 70, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9054);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9055, 5124, '/110/9360/5124/', '', '城市数据', '/freight/city/list', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9055);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9056, 20087, '/87/', '', '价格说明', '/product/priceDescription', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9056);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9057, 9052, '/87/9052/', '', '店铺商品参数', '/product/merSpecs', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9057);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9058, 9052, '/87/9052/', '', '平台商品参数', '/product/specs', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9058);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9059, 20095, '95/', '', '商品参数', '/product/specs', '[]', 7, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9059);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9060, 9050, '/101/9046/9050/', '', '权限', '/systemForm/Basics/svip', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9060);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9061, 7669, '/9361/1538/1539/7669/', '', '统计', 'systemCommunityTitle', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9061);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9062, 9060, '/101/9046/9050/9060/', '', '编辑配置信息', 'configSave', 'svip', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9062);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9063, 7687, '/101/103/7687/', '', '用户标签编辑', 'systemUserSvipUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9063);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9064, 7693, '/110/9360/5124/119/7693/', '', '列表', 'systemCityAreaLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9064);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9065, 7693, '/110/9360/5124/119/7693/', '', '编辑', 'systemCityAreaCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9065);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9066, 7693, '/110/9360/5124/119/7693/', '', '编辑', 'systemCityAreaUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9066);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9067, 7693, '/110/9360/5124/119/7693/', '', '删除', 'systemCityAreaDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9067);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9068, 7769, '/540/541/7769/', '', '记录', 'systemOrderStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9068);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9069, 7784, '/719/1629/7784/', '', '价格说明列表', 'systemPriceRuleLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9069);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9070, 7784, '/719/1629/7784/', '', '添加价格说明', 'systemPriceRuleCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9070);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9071, 7784, '/719/1629/7784/', '', '修改价格说明', 'systemPriceRuleUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9071);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9072, 7784, '/719/1629/7784/', '', '价格说明修改状态', 'systemPriceRuleStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9072);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9073, 7784, '/719/1629/7784/', '', '删除价格说明', 'systemPriceRuleDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9073);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9074, 8592, '/512/513/8592/', '', '核销详情', 'merchantStoreOrderVerifyDetail', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9074);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9075, 9051, '/101/9046/9051/', '', '权限', '/user/member/type', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9075);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9076, 9051, '/101/9046/9051/', '', '附加权限', 'append_/user/member/type', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9076);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9077, 9049, '/101/9046/9049/', '', '权限', '/user/member/equity', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9077);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9078, 9049, '/101/9046/9049/', '', '附加权限', 'append_/user/member/equity', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9078);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9079, 9058, '/87/9052/9058/', '', '权限', '/product/specs', '', 1, 0, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9079);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9080, 9057, '/87/9052/9057/', '', '权限', '/product/merSpecs', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9080);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9081, 9075, '/101/9046/9051/9075/', '', '列表', 'systemUserSvipLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9081);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9082, 9075, '/101/9046/9051/9075/', '', '添加', 'systemUserSvipCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9082);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9083, 9075, '/101/9046/9051/9075/', '', '编辑表单', 'systemUserSvipUpdateForm', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9083);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9084, 9075, '/101/9046/9051/9075/', '', '编辑', 'systemUserSvipTypeUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9084);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9085, 9075, '/101/9046/9051/9075/', '', '删除', 'systemUserSvipDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9085);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9086, 9075, '/101/9046/9051/9075/', '', '修改状态', 'systemUserSvipStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9086);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9087, 9076, '/101/9046/9051/9076/', '', '上传图片', 'uploadImage', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9087);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9088, 9076, '/101/9046/9051/9076/', '', '图片列表', 'systemAttachmentLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9088);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9089, 9077, '/101/9046/9049/9077/', '', '列表', 'systemUserSvipInterestsLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9089);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9090, 9077, '/101/9046/9049/9077/', '', '编辑', 'systemUserSvipInterestsUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9090);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9091, 9077, '/101/9046/9049/9077/', '', '编辑状态', 'systemUserSvipInterestsStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9091);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9092, 9078, '/101/9046/9049/9078/', '', '上传图片', 'uploadImage', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9092);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9093, 9078, '/101/9046/9049/9078/', '', '图片列表', 'systemAttachmentLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9093);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9094, 9079, '/87/9052/9058/9079/', '', '平台参数列表', 'systemStoreParameterTemplateLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9094);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9095, 9079, '/87/9052/9058/9079/', '', '详情', 'systemStoreParameterTemplateDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9095);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9096, 9079, '/87/9052/9058/9079/', '', '删除', 'systemStoreParameterTemplateDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9096);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9097, 9079, '/87/9052/9058/9079/', '', '添加', 'systemStoreParameterTemplateCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9097);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9098, 9079, '/87/9052/9058/9079/', '', '编辑', 'systemStoreParameterTemplateUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9098);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9099, 9080, '/87/9052/9057/9080/', '', '商户参数模板', 'systemStoreParameterTemplateMerLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9099);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9100, 9059, '95/9059/', '', '权限', '/product/specs', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9100);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9101, 9100, '95/9059/9100/', '', '列表', 'merchantStoreParameterTemplateLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9101);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9102, 9100, '95/9059/9100/', '', '详情', 'merchantStoreParameterTemplateDetail', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9102);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9103, 9100, '95/9059/9100/', '', '删除', 'merchantStoreParameterTemplateDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9103);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9104, 9100, '95/9059/9100/', '', '添加', 'merchantStoreParameterTemplateCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9104);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9105, 9100, '95/9059/9100/', '', '编辑', 'merchantStoreParameterTemplateUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9105);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9106, 9056, '/87/9056/', '', '权限', '/product/priceDescription', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9106);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9107, 9056, '/87/9056/', '', '附加权限', 'append_/product/priceDescription', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9107);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9108, 9106, '/87/9056/9106/', '', '价格说明列表', 'systemPriceRuleLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9108);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9109, 9106, '/87/9056/9106/', '', '添加价格说明', 'systemPriceRuleCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9109);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9110, 9106, '/87/9056/9106/', '', '修改价格说明', 'systemPriceRuleUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9110);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9111, 9106, '/87/9056/9106/', '', '价格说明修改状态', 'systemPriceRuleStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9111);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9112, 9106, '/87/9056/9106/', '', '删除价格说明', 'systemPriceRuleDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9112);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9113, 9107, '/87/9056/9107/', '', '上传图片', 'uploadImage', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9113);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9114, 9107, '/87/9056/9107/', '', '图片列表', 'systemAttachmentLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9114);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9115, 8600, '/95/105/8600/', '', '批量设置推荐', 'merchantStoreProductSwitchBatchExtension', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9115);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9116, 8600, '/95/105/8600/', '', '批量设置会员价', 'merchantStoreProductSwitchBatchSvipType', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9116);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9117, 1665, '/1665/', '', '页面配置', '/setting/system_visualization_data', '[]', 50, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9117);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9118, 1289, '/719/1289/', '', '商品分类', '/marketing/integral/classify', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9118);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9119, 1289, '/719/1289/', '', '商品列表', '/marketing/integral/proList', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9119);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9120, 1289, '/719/1289/', '', '积分订单', '/marketing/integral/orderList', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9120);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9121, 1284, '/9492/1284/', '', '保证金配置', '/systemForm/Basics/margin', '[]', 60, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9121);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9122, 9121, '/9492/1284/9121/', '', '权限', '/systemForm/Basics/margin', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9122);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9123, 9118, '/719/1289/9118/', '', '权限', '/marketing/integral/classify', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9123);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9124, 9120, '/719/1289/9120/', '', '权限', '/marketing/integral/orderList', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9124);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9125, 9122, '/9492/1284/9121/9122/', '', '编辑配置信息', 'configSave', 'margin', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9125);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9126, 7687, '/101/103/7687/', '', '积分记录', 'systemUserIntegralList', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9126);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9127, 7687, '/101/103/7687/', '', '签到记录', 'systemUserSginLog', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9127);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9128, 7687, '/101/103/7687/', '', '浏览记录', 'systemUserHistory', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9128);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9132, 7766, '/9492/1284/1647/7766/', '', '待缴列表', 'systemMarginMakeUpMarginLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9132);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9133, 7769, '/540/541/7769/', '', '关联订单', 'systemOrderChildrenList', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9133);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9134, 9123, '/719/1289/9118/9123/', '', '列表', 'pointsCateLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9134);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9135, 9123, '/719/1289/9118/9123/', '', '详情', 'pointsCateDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9135);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9136, 9123, '/719/1289/9118/9123/', '', '添加', 'pointsCateCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9136);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9137, 9123, '/719/1289/9118/9123/', '', '编辑', 'pointsCateUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9137);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9138, 9123, '/719/1289/9118/9123/', '', '修改状态', 'pointsCateStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9138);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9139, 9123, '/719/1289/9118/9123/', '', '列表', 'pointsProductLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9139);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9140, 9123, '/719/1289/9118/9123/', '', '获取规格', 'pointsCateFormatAttr', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9140);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9141, 9123, '/719/1289/9118/9123/', '', '编辑', 'pointsProductDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9141);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9142, 9123, '/719/1289/9118/9123/', '', '添加', 'pointsProductCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9142);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9143, 9123, '/719/1289/9118/9123/', '', '编辑', 'pointsProductUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9143);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9144, 9123, '/719/1289/9118/9123/', '', '修改状态', 'pointsProductStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9144);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9145, 9123, '/719/1289/9118/9123/', '', '预览', 'pointsProductPreview', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9145);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9146, 9124, '/719/1289/9120/9124/', '', '列表', 'pointsOrderLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9146);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9147, 9124, '/719/1289/9120/9124/', '', '编辑', 'pointsOrderDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9147);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9148, 9124, '/719/1289/9120/9124/', '', '发货', 'pointsOrderDelivery', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9148);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9149, 9124, '/719/1289/9120/9124/', '', '批量发货', 'pointsOrderBatchDelivery', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9149);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9150, 9124, '/719/1289/9120/9124/', '', '快递查询', 'pointsOrderExpress', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9150);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9151, 9124, '/719/1289/9120/9124/', '', '导出', 'pointsOrderExcel', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9151);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9152, 9124, '/719/1289/9120/9124/', '', '备注', 'pointsOrderMark', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9152);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9153, 9124, '/719/1289/9120/9124/', '', '记录', 'pointsOrderStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9153);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9154, 9124, '/719/1289/9120/9124/', '', '删除', 'pointsOrderDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9154);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9155, 8592, '/512/513/8592/', '', '关联订单', 'merchantStoreOrderChildrenList', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9155);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9156, 8600, '/95/105/8600/', '', '获取规格', 'merchantStoreProductFormatAttr', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9156);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9157, 9117, '/1665/9117/', '', '权限', '/setting/system_visualization_data', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9157);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9158, 9119, '/719/1289/9119/', '', '权限', '/marketing/integral/proList', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9158);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9159, 9157, '/1665/9117/9157/', '', '可视化列表', 'systemVisualStoreGetThemeKey', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9159);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9160, 9157, '/1665/9117/9157/', '', '可视化详情', 'systemVisualStoreGetTheme', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9160);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9161, 9157, '/1665/9117/9157/', '', '可视化保存', 'systemVisualSetTheme', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9161);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9162, 9158, '/719/1289/9119/9158/', '', '列表', 'pointsProductLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9162);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9163, 9158, '/719/1289/9119/9158/', '', '获取规格', 'pointsCateFormatAttr', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9163);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9164, 9158, '/719/1289/9119/9158/', '', '编辑', 'pointsProductDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9164);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9165, 9158, '/719/1289/9119/9158/', '', '添加', 'pointsProductCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9165);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9166, 9158, '/719/1289/9119/9158/', '', '编辑', 'pointsProductUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9166);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9167, 9158, '/719/1289/9119/9158/', '', '修改状态', 'pointsProductStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9167);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9168, 9158, '/719/1289/9119/9158/', '', '预览', 'pointsProductPreview', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9168);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9169, 514, '/514/', '', '分销订单', '/promoter/orderList', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9169);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9170, 1665, '/1665/', '', '店铺模板', '/setting/merchant/diyList', '[]', 70, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9170);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9172, 526, '526/', '', '开放账户', '/systemForm/openAuth/list', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9172);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9173, 7668, '/1665/41/7668/', '', '上传二维码', 'systemAttachmentScanQrcode', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9173);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9174, 7668, '/1665/41/7668/', '', '扫码上传图片', 'systemAttachmentScanImage', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9174);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9175, 7668, '/1665/41/7668/', '', '扫码上传保存', 'systemAttachmentScanImageSave', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9175);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9176, 7668, '/1665/41/7668/', '', '在线图片', 'systemAttachmentOnline', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9176);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9177, 7697, '/1665/6372/1669/7697/', '', '修改状态', 'systemDiyPageLinkStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9177);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9178, 7698, '/1665/6372/1674/7698/', '', '修改状态', 'systemDiyPageLinkMerStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9178);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9179, 7699, '/1665/1666/7699/', '', '详情 ', 'systemDiyDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9179);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9180, 7763, '/9492/6370/44/7763/', '', '详情', 'systemMerchantDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9180);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9181, 9214, '/514/9169/9214/', '', '金额统计', 'systemSpreadOrderStat', '', 1, 0, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9181);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9182, 9214, '/514/9169/9214/', '', '快递查询', 'systemSpreadOrderExpress', '', 1, 0, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9182);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9183, 9214, '/514/9169/9214/', '', '头部统计', 'systemSpreadOrderTitle', '', 1, 0, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9183);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9184, 9214, '/514/9169/9214/', '', '详情', 'systemSpreadOrderDetail', '', 1, 0, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9184);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9185, 9214, '/514/9169/9214/', '', '导出', 'systemSpreadOrderExcel', '', 1, 0, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9185);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9186, 9214, '/514/9169/9214/', '', '记录', 'systemSpreadOrderStatus', '', 1, 0, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9186);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9187, 9214, '/514/9169/9214/', '', '关联订单', 'systemSpreadOrderChildrenList', '', 1, 0, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9187);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9188, 8558, '1671/54/8558/', '', '上传二维码', 'merchantAttachmentScanQrcode', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9188);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9189, 8558, '1671/54/8558/', '', '扫码上传图片', 'merchantAttachmentScanImage', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9189);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9190, 8558, '1671/54/8558/', '', '扫码上传保存', 'merchantAttachmentScanImageSave', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9190);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9191, 8558, '1671/54/8558/', '', '在线图片', 'merchantAttachmentOnline', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9191);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9192, 8569, '/1671/1672/8569/', '', '默认模板列表 ', 'merchantDefaultDiyLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9192);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9193, 8569, '/1671/1672/8569/', '', '详情 ', 'merchantDiyDetail', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9193);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9194, 9213, '9213/', '', '列表', 'merchantOpenapiLst', '', 1, 0, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9194);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9195, 9213, '9213/', '', '添加', 'merchantOpenapiCreate', '', 1, 0, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9195);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9196, 9213, '9213/', '', '编辑', 'merchantOpenapiUpdate', '', 1, 0, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9196);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9197, 9213, '9213/', '', '修改状态', 'merchantOpenapiStatus', '', 1, 0, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9197);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9198, 9213, '9213/', '', '删除', 'merchantOpenapiDeleta', '', 1, 0, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9198);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9199, 9213, '9213/', '', '查看', 'merchantOpenapiGetSecretKey', '', 1, 0, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9199);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9200, 9213, '9213/', '', '重置', 'merchantOpenapiSetSecretKey', '', 1, 0, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9200);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9201, 8595, '/512/528/8595/', '', '创建', 'merchantStoreRefundOrderCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9201);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9202, 9170, '/1665/9170/', '', '权限', '/setting/merchant/diyList', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9202);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9203, 9202, '/1665/9170/9202/', '', '列表 ', 'systemMerDiyLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9203);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9204, 9202, '/1665/9170/9202/', '', '详情 ', 'systemMerDiyDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9204);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9205, 9202, '/1665/9170/9202/', '', '添加/编辑', 'systemMerDiyCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9205);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9207, 9202, '/1665/9170/9202/', '', '设置默认', 'systemMerDiySetDefault', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9207);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9208, 9202, '/1665/9170/9202/', '', '重置', 'systemMerDiyRecovery', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9208);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9209, 9202, '/1665/9170/9202/', '', '删除', 'systemMerDiyDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9209);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9210, 9202, '/1665/9170/9202/', '', '复制', 'systemMerDiyCopy', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9210);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9211, 9202, '/1665/9170/9202/', '', '保存适用范围', 'systemMerDiyGetScope', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9211);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9212, 9202, '/1665/9170/9202/', '', '保存适用范围', 'systemMerDiySetScope', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9212);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9213, 9172, '9172/', '', '权限', '/systemForm/openAuth/list', '', 0, 0, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9213);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9214, 9169, '/514/9169/', '', '权限', '/promoter/orderList', '', 0, 0, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9214);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9215, 20101, '/101/', '', '用户设置', '/user/setup_user', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9215);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9216, 1665, '/1665/', '', '系统表单', '/systemForm/form_list', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9216);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9217, 719, '/719/', '', '报名活动', '/marketing/application/list', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9217);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9218, 0, '/', 's-home', '首页', '/', '[]', 127, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9218);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9219, 9218, '/9218/', '', '数据大屏', '/data-screen/index', '[]', 99, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9219);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9220, 515, '/515/', '', '店铺结算', '/mer/accounts', '[]', 10, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9220);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9221, 9220, '/515/9220/', '', '店铺账单', '/accounts/merchantBill', '[]', 9, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9221);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9222, 1671, '1671/', '', '系统表单', '/systemForm/form_list', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9222);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9223, 20095, '95/', '', '商品单位', '/product/unit', '[]', 6, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9223);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9224, 1178, '/515/9220/1178/', '', '权限', '/accounts/statement', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9224);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9225, 9221, '/515/9220/9221/', '', '权限', '/accounts/merchantBill', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9225);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9226, 9219, '/9218/9219/', '', '权限', '/data-screen/index', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9226);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9227, 7653, '/515/537/516/7653/', '', '详情', 'systemUserExtractDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9227);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9228, 9224, '/515/9220/1178/9224/', '', '列表', 'systemFinancialRecordLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9228);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9229, 9224, '/515/9220/1178/9224/', '', '统计', 'systemFinancialRecordTitle', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9229);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9230, 9224, '/515/9220/1178/9224/', '', '详情', 'systemFinancialRecordDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9230);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9231, 9224, '/515/9220/1178/9224/', '', '导出', 'systemFinancialRecordDetailExport', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9231);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9232, 9225, '/515/9220/9221/9225/', '', '商户列表', 'systemFinancialRecordMerLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9232);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9233, 9225, '/515/9220/9221/9225/', '', '商户统计', 'systemFinancialRecordMerAcountsLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9233);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9234, 9225, '/515/9220/9221/9225/', '', '商户财务头部统计', 'systemFinancialRecordMerTitle', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9234);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9235, 9225, '/515/9220/9221/9225/', '', '商户财务详情', 'systemFinancialRecordMerDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9235);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9236, 9225, '/515/9220/9221/9225/', '', '商户财务导出', 'systemFinancialRecordMerExcel', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9236);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9237, 7687, '/101/103/7687/', '', '用户信息导出', 'systemUserExcel', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9237);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9238, 7687, '/101/103/7687/', '', '批量设置分销员', 'getMemberLevelBatchSpread', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9238);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9239, 7763, '/9492/6370/44/7763/', '', '操作日志', 'systemMerchantOperateList', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9239);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9240, 9214, '/514/9169/9214/', '', '列表', 'systemSpreadOrderLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9240);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9241, 7778, '/87/539/7778/', '', '获取商品操作记录', 'systemStoreProductGetOperateList', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9241);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9242, 7799, '/9218/33/7799/', '', '未处理业务统计', 'systemStatisticsAdminCount', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9242);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9243, 7799, '/9218/33/7799/', '', '待办事项', 'systemStatisticsAdminTodo', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9243);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9244, 7799, '/9218/33/7799/', '', '商户销量排行', 'systemStatisticsMerchantTop', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9244);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9245, 9226, '/9218/9219/9226/', '', '数据大屏', 'systemDataScreen', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9245);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9246, 9223, '95/9223/', '', '权限', '/product/unit', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9246);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9247, 8568, '1649/700/8568/', '', '设置默认模板', 'merchantStoreShippingTemplateSetDefault', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9247);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9248, 8600, '/95/105/8600/', '', '获取批量修改列表', 'merchantStoreProductGetBatchList', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9248);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9249, 8600, '/95/105/8600/', '', '批量修改商品属性', 'merchantStoreProductSwitchBatchProcess', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9249);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9250, 8600, '/95/105/8600/', '', '操作记录', 'merchantStoreProductGetOperateList', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9250);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9251, 9246, '95/9223/9246/', '', '商品单位列表', 'merchantStoreProductUnitLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9251);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9252, 9246, '95/9223/9246/', '', '商品单位添加', 'merchantStoreProductUnitCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9252);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9253, 9246, '95/9223/9246/', '', '商品单位编辑', 'merchantStoreProductUnitUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9253);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9254, 9246, '95/9223/9246/', '', '商品单位删除', 'merchantStoreProductUnitDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9254);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9255, 8614, '/55/8614/', '', '首页未处理业务统计', 'merchantStatisticsMerchantCount', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9255);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9256, 8614, '/55/8614/', '', '待办事项', 'merchantStatisticsMerchantTodo', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9256);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9257, 8614, '/55/8614/', '', '获取商户代办统计', 'merchantStatisticsProductSalesPriceTop', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9257);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9258, 9216, '/1665/9216/', '', '权限', '/systemForm/form_list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9258);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9259, 9217, '/719/9217/', '', '权限', '/marketing/application/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9259);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9260, 9217, '/719/9217/', '', '附加权限', 'append_/marketing/application/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9260);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9261, 9215, '/101/9215/', '', '权限', '/user/setup_user', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9261);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9262, 9258, '/1665/9216/9258/', '', '添加', 'systemFormCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9262);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9263, 9258, '/1665/9216/9258/', '', '编辑', 'systemFormUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9263);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9264, 9258, '/1665/9216/9258/', '', '编辑状态', 'systemFormStatusSwitch', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9264);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9265, 9258, '/1665/9216/9258/', '', '删除', 'systemFormDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9265);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9266, 9258, '/1665/9216/9258/', '', '详情', 'systemFormDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9266);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9267, 9258, '/1665/9216/9258/', '', '列表', 'systemFormLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9267);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9268, 9258, '/1665/9216/9258/', '', '导出', 'systemFormExcel', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9268);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9269, 9258, '/1665/9216/9258/', '', '表单提交记录', 'systemFormUserLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9269);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9270, 9259, '/719/9217/9259/', '', '添加', 'systemActivityFormCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9270);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9271, 9259, '/719/9217/9259/', '', '列表', 'systemActivityFormLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9271);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9272, 9259, '/719/9217/9259/', '', '编辑', 'systemActivityFormUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9272);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9273, 9259, '/719/9217/9259/', '', '详情', 'systemActivityFormDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9273);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9274, 9259, '/719/9217/9259/', '', '删除', 'systemActivityFormDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9274);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9275, 9259, '/719/9217/9259/', '', '修改状态', 'systemActivityFormStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9275);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9276, 9259, '/719/9217/9259/', '', '活动记录', 'systemFormActivUserLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9276);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9277, 9259, '/719/9217/9259/', '', '活动记录导出', 'systemFormActivUserExcel', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9277);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9278, 9260, '/719/9217/9260/', '', '上传图片', 'uploadImage', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9278);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9279, 9260, '/719/9217/9260/', '', '图片列表', 'systemAttachmentLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9279);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9280, 9261, '/101/9215/9261/', '', '列表', 'systemUserInfolst', '', 1, 0, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9280);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9281, 9261, '/101/9215/9261/', '', '添加', 'systemUserInfoCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9281);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9282, 9261, '/101/9215/9261/', '', '保存信息', 'systemUserInfoSaveAll', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9282);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9283, 9261, '/101/9215/9261/', '', '删除', 'systemUserInfoDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9283);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9284, 8617, '/526/546/8617/', '', '退保证金申请', 'merchantFinancialRefundMarginApply', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9284);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9285, 20118, '/110/118/', '', '存储管理', '/setting/storage', '[]', 90, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9285);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9286, 1376, '1376/', '', '自有一号通', '/setting/sms/sms_account/index', '[]', 10, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9286);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9287, 780, '/719/780/', '', '秒杀活动', '/marketing/seckill/store_seckill/list', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9287);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9289, 788, '788/', '', '秒杀活动', '/marketing/seckill/store_seckill/list', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9289);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9290, 788, '788/', '', '秒杀商品', '/marketing/seckill/product/list', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9290);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9291, 9287, '/719/780/9287/', '', '权限', '/marketing/seckill/store_seckill/list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9291);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9292, 9285, '/110/118/9285/', '', '权限', '/setting/storage', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9292);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9293, 9291, '/719/780/9287/9291/', '', '列表', 'systemSeckillActiveGetActiveList', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9293);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9294, 9291, '/719/780/9287/9291/', '', '详情', 'systemSeckillActiveGetActiveInfo', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9294);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9295, 9291, '/719/780/9287/9291/', '', '创建', 'systemSeckillActiveCreateActive', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9295);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9296, 9291, '/719/780/9287/9291/', '', '编辑', 'systemSeckillActiveUpdateActive', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9296);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9297, 9291, '/719/780/9287/9291/', '', '编辑状态', 'systemSeckillActiveUpdateActiveStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9297);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9298, 9291, '/719/780/9287/9291/', '', '删除', 'systemSeckillActiveDeleteActive', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9298);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9299, 9291, '/719/780/9287/9291/', '', '活动统计数据面板', 'systemSeckillActiveChartPanel', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9299);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9300, 9291, '/719/780/9287/9291/', '', '活动参与人统计列表', 'systemSeckillActiveChartPeople', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9300);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9301, 9291, '/719/780/9287/9291/', '', '活动订单统计列表', 'systemSeckillActiveChartOrder', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9301);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9302, 9291, '/719/780/9287/9291/', '', '活动商品统计列表', 'systemSeckillActiveChartProduct', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9302);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9303, 7754, '/719/780/794/7754/', '', '列表', 'systemStoreSeckillProductPageLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9303);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9304, 7754, '/719/780/794/7754/', '', '加入回收站', 'systemStoreSeckillProductDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9304);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9305, 7754, '/719/780/794/7754/', '', '删除', 'systemStoreSeckillProductDestory', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9305);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9306, 7759, '/110/38/48/7759/', '', '搜索获取菜单', 'getMenusList', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9306);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9307, 7778, '/87/539/7778/', '', '获取自营商品列表', 'systemStoreProductGetSelfProductList', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9307);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9308, 9292, '/110/118/9285/9292/', '', '配置信息', 'systemStorageGetConfig', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9308);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9309, 9292, '/110/118/9285/9292/', '', '提交配置', 'systemStorageSetConfig', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9309);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9310, 9292, '/110/118/9285/9292/', '', '保存云存储配置', 'systemStorageUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9310);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9311, 9292, '/110/118/9285/9292/', '', '同步存储空间', 'systemStorageSync', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9311);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9312, 9292, '/110/118/9285/9292/', '', '存储空间列表', 'systemStorageLstRegion', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9312);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9313, 9292, '/110/118/9285/9292/', '', '添加存储空间', 'systemStorageCreateRegion', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9313);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9314, 9292, '/110/118/9285/9292/', '', '删除存储空间', 'systemStorageDeleteRegion', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9314);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9315, 9292, '/110/118/9285/9292/', '', '使用存储空间', 'systemStorageRegionSwtichStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9315);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9316, 9292, '/110/118/9285/9292/', '', '修改存储空间名称', 'systemStorageUpdateDomain', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9316);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9317, 7801, '/520/116/1617/7801/', '', '替换素材域名', 'systemAttachmentReplaceHost', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9317);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9318, 9289, '788/9289/', '', '权限', '/marketing/seckill/store_seckill/list', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9318);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9319, 9318, '788/9289/9318/', '', '列表', 'merchantStoreSeckillActiveGetActiveList', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9319);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9320, 9318, '788/9289/9318/', '', '详情', 'merchantStoreSeckillActiveGetActiveInfo', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9320);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9321, 9318, '788/9289/9318/', '', '列表', 'merchantStoreSeckillActiveGetActiveAll', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9321);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9322, 9318, '788/9289/9318/', '', '活动统计数据面板', 'merchantStoreSeckillActiveChartPanel', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9322);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9323, 9318, '788/9289/9318/', '', '活动参与人统计列表', 'merchantStoreSeckillActiveChartPeople', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9323);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9324, 9318, '788/9289/9318/', '', '活动订单统计列表', 'merchantStoreSeckillActiveChartOrder', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9324);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9325, 9318, '788/9289/9318/', '', '活动商品统计列表', 'merchantStoreSeckillActiveChartProduct', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9325);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9326, 8575, '9290/8575/', '', '分页列表', 'merchantStoreSeckillProductPageLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9326);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9327, 8575, '9290/8575/', '', '商品列表', 'merchantStoreSeckillProductGetProductList', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9327);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9328, 8575, '9290/8575/', '', '设置标签', 'merchantStoreSeckillProductSetLabels', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9328);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9329, 20095, '95/', '', '卡密列表', '/product/cdkey', '[]', 9, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9329);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9330, 1661, '1661/', '', '配送人员', '/delivery/personnel_manage/index', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9330);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9331, 9410, '9410/', '', '店员配置', '/systemForm/Basics/mer_service_switch', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9331);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9332, 526, '526/', '', '同城配送', '/delivery', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9332);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9334, 7778, '/87/539/7778/', '', '批量设置分类推荐', 'systemStoreProductSwitchBatchCateHot', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9334);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9335, 9261, '/101/9215/9261/', '', '保存注册配置', 'systemUserRegisterConfig', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9335);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9336, 9261, '/101/9215/9261/', '', '新人礼优惠券列表', 'systemUserRegisterCoupon', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9336);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9337, 9329, '95/9329/', '', '权限', '/product/cdkey', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9337);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9338, 8451, '8451/', '', '列表', 'merchantDeliveryServiceLst', '', 1, 0, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9338);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9339, 8451, '8451/', '', '修改状态', 'merchantDeliveryServiceStatus', '', 1, 0, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9339);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9340, 8451, '8451/', '', '添加', 'merchantDeliveryServiceCreate', '', 1, 0, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9340);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9341, 8451, '8451/', '', '编辑', 'merchantDeliveryServiceUpdate', '', 1, 0, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9341);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9342, 8451, '8451/', '', '删除', 'merchantDeliveryServiceDelete', '', 1, 0, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9342);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9343, 8592, '/512/513/8592/', '', '线下支付', 'merchantStoreOrderOffline', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9343);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9344, 8600, '/95/105/8600/', '', '操作记录', 'merchantStoreProductUnbind', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9344);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9345, 9337, '95/9329/9337/', '', '列表', 'merchantStoreProductCdkeyLibraryLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9345);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9346, 9337, '95/9329/9337/', '', '列表', 'merchantStoreProductCdkeyLibraryDetail', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9346);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9347, 9337, '95/9329/9337/', '', '添加', 'merchantStoreProductCdkeyLibraryCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9347);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9348, 9337, '95/9329/9337/', '', '编辑表单', 'merchantStoreProductCdkeyLibraryUpdateForm', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9348);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9349, 9337, '95/9329/9337/', '', '编辑', 'merchantStoreProductCdkeyLibraryUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9349);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9350, 9337, '95/9329/9337/', '', '删除', 'merchantStoreProductCdkeyLibraryDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9350);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9351, 9337, '95/9329/9337/', '', '导出', 'merchantStoreProductCdkeyLibraryExcel', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9351);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9352, 9337, '95/9329/9337/', '', '批量导入', 'merchantStoreProductCdkeyLibraryImport', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9352);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9353, 9337, '95/9329/9337/', '', '卡密列表', 'merchantStoreProductCdkeyLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9353);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9354, 9337, '95/9329/9337/', '', '添加卡密', 'merchantStoreProductCdkeyCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9354);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9355, 9337, '95/9329/9337/', '', '编辑卡密', 'merchantStoreProductCdkeyUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9355);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9356, 9337, '95/9329/9337/', '', '删除卡密', 'merchantStoreProductCdkeyDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9356);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9357, 9337, '95/9329/9337/', '', '批量删除', 'merchantStoreProductCdkeyLibraryBatchDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9357);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9358, 526, '526/', '', '打印配置', '/setting/printer', '[]', 95, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9358);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9359, 9358, '9358/', '', '打印配置', '/systemForm/Basics/printer_tabs', '[]', 10, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9359);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9360, 20110, '/110/', '', '商城设置', '/shop', '[]', 99, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9360);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9361, 0, '/', 's-management', '内容', '/content', '[]', 96, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9361);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9362, 8592, '/512/513/8592/', '', '电子面单复打', 'merchantStoreOrderRepeatDump', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9362);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9363, 1665, '/1665/', '', '商品详情', '/setting/diy/product_detail', '[]', 70, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9363);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9365, 9367, '/42/6370/9367/', '', '权限', '/merchant/grouping', '[]', 0, 0, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9365);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9367, 6370, '/9492/6370/', '', '店铺分组', '/merchant/grouping', '[]', 21, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9367);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9368, 514, '/514/', '', '分销说明', '/promoter/retail', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9368);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9369, 512, '512/', '', '代客下单', '/order/customer', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9369);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9370, 7699, '/1665/1666/7699/', '', '商品详情 ', 'systemDiyGetProductDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9370);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9371, 7699, '/1665/1666/7699/', '', '商品详情保存 ', 'systemDiySaveProductDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9371);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9372, 9365, '/42/6370/9367/9365/', '', '添加', 'systemMerchantRegionCreate', '', 1, 0, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9372);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9373, 9365, '/42/6370/9367/9365/', '', '编辑', 'systemMerchantRegionUpdate', '', 1, 0, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9373);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9374, 9365, '/42/6370/9367/9365/', '', '列表', 'systemMerchantRegionLst', '', 1, 0, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9374);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9375, 9365, '/42/6370/9367/9365/', '', '删除', 'systemMerchantRegionDelete', '', 1, 0, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9375);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9376, 7763, '/9492/6370/44/7763/', '', '虚拟关注量', 'systemMerchantCareFicti', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9376);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9377, 9079, '/87/9052/9058/9079/', '', '删除属性', 'systemStoreParameterTemplateDeleteValue', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9377);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9378, 8592, '/512/513/8592/', '', '配货单', 'merchantStoreOrderNote', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9378);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9379, 9100, '95/9059/9100/', '', '删除属性', 'merchantStoreParameterTemplateDeleteValue', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9379);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9380, 8600, '/95/105/8600/', '', '编辑商品获取信息', 'merchantStoreProductGetEdit', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9380);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9381, 9369, '512/9369/', '', '权限', '/order/customer', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9381);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9382, 9381, '512/9369/9381/', '', '商品分类', 'behalfProductCategory', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9382);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9383, 9381, '512/9369/9381/', '', '商品列表', 'behalfProductList', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9383);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9384, 9381, '512/9369/9381/', '', '商品规格详情', 'behalfProductDetail', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9384);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9385, 9381, '512/9369/9381/', '', '会员查询', 'behalfUserQuery', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9385);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9386, 9381, '512/9369/9381/', '', '会员详情', 'behalfUserInfo', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9386);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9387, 9381, '512/9369/9381/', '', '会员添加', 'behalfUserCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9387);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9388, 9381, '512/9369/9381/', '', '地址列表', 'behalfUserAddressList', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9388);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9389, 9381, '512/9369/9381/', '', '地址添加', 'behalfUserAddressCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9389);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9390, 9381, '512/9369/9381/', '', '购物车列表', 'behalfCartList', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9390);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9391, 9381, '512/9369/9381/', '', '添加购物车', 'behalfCartCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9391);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9392, 9381, '512/9369/9381/', '', '修改购物车数据', 'behalfCartChange', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9392);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9393, 9381, '512/9369/9381/', '', '删除购物数据', 'behalfCartDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9393);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9394, 9381, '512/9369/9381/', '', '清空购物车', 'behalfCartClear', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9394);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9395, 9381, '512/9369/9381/', '', '购物车总数量', 'behalfCartCount', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9395);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9396, 9381, '512/9369/9381/', '', '修改价格', 'behalfCartUpdatePrice', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9396);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9397, 9381, '512/9369/9381/', '', '批量修改价格', 'behalfCartBatchUpdatePrice', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9397);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9398, 9381, '512/9369/9381/', '', '校验订单', 'behalfCheck', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9398);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9399, 9381, '512/9369/9381/', '', '支付配置', 'behalfPayConfig', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9399);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9400, 9381, '512/9369/9381/', '', '创建订单', 'behalfCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9400);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9401, 9381, '512/9369/9381/', '', '支付', 'behalfPay', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9401);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9402, 9381, '512/9369/9381/', '', '获取结果', 'behalfPayStatus', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9402);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9403, 9381, '512/9369/9381/', '', '余额支付获取验证码', 'behalfVerify', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9403);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9404, 1661, '1661/', '', '服务人员', '/config/service_staff', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9404);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9405, 512, '512/', '', '预约设置', '/product/reservation', '[]', 116, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9405);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9406, 512, '512/', '', '预约服务', '/order/reservation', '[]', 120, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9406);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9408, 20106, '106/', '', '逛逛社区', '/community/list', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9408);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9410, 1661, '1661/', '', '店员管理', '/server_manage', '[]', 5, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9410);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9411, 7673, '/9361/1538/1543/7673/', '', '内容评论列表', 'systemCommunityReply', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9411);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9412, 7772, '/540/542/7772/', '', '详情', 'systemRefundOrderDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9412);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9413, 7772, '/540/542/7772/', '', '日志', 'systemRefundOrderLog', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9413);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9414, 7772, '/540/542/7772/', '', '审核', 'systemRefundOrderApprove', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9414);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9415, 7813, '/101/1285/7813/', '', '清除用户搜索记录', 'systemUserClearSearchLog', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9415);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9416, 9408, '106/9408/', '', '权限', '/community/list', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9416);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9417, 9222, '1671/9222/', '', '权限', '/systemForm/form_list', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9417);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9418, 9406, '512/9406/', '', '权限', '/order/reservation', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9418);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9419, 9404, '1661/9404/', '', '权限', '/config/service_staff', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9419);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9420, 9416, '106/9408/9416/', '', '分类列表', 'merchantCommunityCateLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9420);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9421, 9416, '106/9408/9416/', '', '列表', 'merchantCommunityLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9421);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9422, 9416, '106/9408/9416/', '', '添加', 'merchantCommunityCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9422);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9423, 9416, '106/9408/9416/', '', '详情', 'merchantCommunityDetail', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9423);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9424, 9416, '106/9408/9416/', '', '编辑', 'merchantCommunityUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9424);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9425, 9416, '106/9408/9416/', '', '删除', 'merchantCommunityDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9425);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9426, 9416, '106/9408/9416/', '', '评论', 'merchantCommunityReply', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9426);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9427, 8564, '1649/1380/8564/', '', '修改状态', 'merchantExpressChangeMerStatus', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9427);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9428, 9417, '1671/9222/9417/', '', '添加', 'merFormCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9428);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9429, 9417, '1671/9222/9417/', '', '编辑', 'merFormUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9429);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9430, 9417, '1671/9222/9417/', '', '删除', 'merFormDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9430);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9431, 9417, '1671/9222/9417/', '', '详情', 'merFormDetail', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9431);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9432, 9417, '1671/9222/9417/', '', '列表', 'merFormLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9432);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9433, 8592, '/512/513/8592/', '', '修改收货信息', 'merchantStoreOrderCollectCargo', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9433);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9434, 8592, '/512/513/8592/', '', '预约订单派单', 'merchantStoreOrderReservationDispatch', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9434);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9435, 8592, '/512/513/8592/', '', '预约订单改派', 'merchantStoreOrderReservationUpdateDispatch', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9435);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9436, 8592, '/512/513/8592/', '', '预约订单改约', 'merchantStoreOrderReservationReschedule', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9436);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9437, 8592, '/512/513/8592/', '', '单独修改预约时间', 'merchantStoreOrderReservationTime', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9437);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9438, 8592, '/512/513/8592/', '', '预约订单核销', 'merchantStoreOrderReservationVerify', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9438);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9439, 9418, '512/9406/9418/', '', '预约日历', 'merchantReservationServiceList', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9439);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9440, 8600, '/95/105/8600/', '', '添加预约商品', 'merchantStoreReservationProductCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9440);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9441, 8600, '/95/105/8600/', '', '获取预约商品', 'merchantStoreReservationProductEditInfo', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9441);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9442, 8600, '/95/105/8600/', '', '编辑预约商品', 'merchantStoreReservationProductEdit', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9442);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9443, 8600, '/95/105/8600/', '', '预约商品详情', 'merchantStoreReservationProductDetail', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9443);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9444, 8600, '/95/105/8600/', '', '批量修改预约商品库存', 'merchantStoreReservationProductEditStock', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9444);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9445, 9419, '1661/9404/9419/', '', '列表', 'merchantStaffsLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9445);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9446, 9419, '1661/9404/9419/', '', '添加', 'merchantStaffsCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9446);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9447, 9419, '1661/9404/9419/', '', '编辑', 'merchantStaffsUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9447);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9448, 9419, '1661/9404/9419/', '', '修改状态', 'merchantStaffsSwitchStatus', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9448);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9449, 9419, '1661/9404/9419/', '', '删除', 'merchantStaffsDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9449);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9450, 9261, '/101/9215/9261/', '', '扩展信息表单', 'systemUserFieldSaveForm', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9450);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9451, 9261, '/101/9215/9261/', '', '添加或编辑', 'systemUserInfoFieldSave', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9451);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9452, 9332, '9332/', '', '配送门店', '/delivery/delivery_point', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9452);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9453, 1665, '/1665/', '', '悬浮菜单', '/setting/fab', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9453);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9454, 1665, '/1665/', '', '商品分类', '/setting/product_category', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9454);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9455, 1671, '1671/', '', '商品分类', '/devise/diy/product_category', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9455);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9456, 9332, '9332/', '', '配送设置', '/setting/delivery', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9456);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9459, 1665, '/1665/', '', '店铺街', '/setting/diy/store', '[]', 78, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9459);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9460, 1665, '/1665/', '', '个人中心', '/setting/diy/personal', '[]', 77, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9460);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9463, 9453, '/1665/9453/', '', '权限', '/setting/fab', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9463);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9464, 9453, '/1665/9453/', '', '附加权限', 'append_/setting/fab', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9464);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9465, 9454, '/1665/9454/', '', '权限', '/setting/product_category', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9465);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9466, 9463, '/1665/9453/9463/', '', '悬浮按钮信息', 'systemDiyFabInfo', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9466);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9467, 9463, '/1665/9453/9463/', '', '保存悬浮按钮', 'systemDiyFabSave', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9467);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9468, 9464, '/1665/9453/9464/', '', '上传图片', 'uploadImage', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9468);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9469, 9464, '/1665/9453/9464/', '', '图片列表', 'systemAttachmentLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9469);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9470, 9465, '/1665/9454/9465/', '', '商品分类信息', 'systemDiyProductCategoryInfo', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9470);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9471, 9465, '/1665/9454/9465/', '', '保存商品分类', 'systemDiyProductCategorySave', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9471);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9472, 9455, '1671/9455/', '', '权限', '/devise/diy/product_category', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9472);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9474, 9472, '1671/9455/9472/', '', '商品分类信息', 'merchantDiyProductCategoryInfo', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9474);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9475, 9472, '1671/9455/9472/', '', '保存商品分类', 'merchantDiyProductCategorySave', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9475);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9476, 8592, '/512/513/8592/', '', '获取商家寄件价格', 'merchantStoreOrderGetPrice', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9476);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9477, 8592, '/512/513/8592/', '', '获取商家寄件价格', 'merchantStoreOrderShipmentList', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9477);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9478, 8592, '/512/513/8592/', '', '取消商家寄件', 'merchantStoreOrderCancelShipment', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9478);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9479, 8592, '/512/513/8592/', '', '同城配送派单', 'merchantStoreOrderDeliveryDispatch', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9479);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9480, 8592, '/512/513/8592/', '', '同城配送改派', 'merchantStoreOrderDeliveryUpdateDispatch', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9480);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9481, 8592, '/512/513/8592/', '', '同城配送核销', 'merchantStoreOrderDeliveryConfirm', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9481);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9482, 8592, '/512/513/8592/', '', '同城配送再次同步', 'merchantStoreOrderDeliverySync', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9482);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9485, 9456, '9456/', '', '权限', '/setting/delivery', '[]', 0, 0, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9485);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9486, 9485, '9485/', '', '配送设置信息', 'merchantDeliveryConfigSettings', '', 0, 0, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9486);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9487, 9485, '9485/', '', '更新配送设置信息', 'merchantDeliveryConfigUpdate', '', 1, 0, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9487);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9488, 20087, '/87/', '', '活动标签', '/product/activityLabel', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9488);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9489, 9218, '/9218/', '', '订单统计', '/statistic/order', '[]', 97, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9489);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9490, 9218, '/9218/', '', '商品统计', '/statistic/product', '[]', 98, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9490);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9491, 9218, '/9218/', '', '用户统计', '/statistic/member', '[]', 96, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9491);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9492, 0, '/', 'place', '店铺', '/mer', '[]', 110, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9492);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9493, 9492, '/9492/', '', '区域代理', '/business-zones/manage', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9493);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9495, 9493, '/9492/9493/', '', '代理人员', '/business-zones/agents', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9495);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9496, 9493, '/9492/9493/', '', '代理审核', '/business-zones/agent-review', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9496);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9497, 9493, '/9492/9493/', '', '代理设置', '/business-zones/settings', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9497);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9499, 20055, '55/', '', '控制台', '/dashboard', '[]', 9, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9499);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9500, 20055, '55/', '', '商品统计', '/statistic/product', '[]', 8, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9500);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9501, 20055, '55/', '', '订单统计', '/statistic/order', '[]', 7, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9501);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9502, 515, '/515/', '', '商圈代理', '/accounts/zoneAgent', '[]', 20, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9502);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9503, 9502, '/515/9502/', '', '结算审核', '/accounts/zoneAgent/settlementReview', '[]', 40, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9503);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9506, 9502, '/515/9502/', '', '提成流水', '/accounts/zoneAgent/commissionRecord', '[]', 10, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9506);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9515, 0, '/', 's-tools', '设置', '/settings', '[]', 0, 1, 1, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9515);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9551, 0, '/', 's-data', '财务', '/accounts', '[]', 70, 1, 1, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9551);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9556, 9404, '9404/', '', '服务人员', '/config/service_staff', '[]', 1, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9556);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9557, 9404, '9404/', '', '服务统计', '/config/service_statistic', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9557);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9558, 9330, '9330/', '', '配送员管理', '/delivery/personnel_manage', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9558);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9559, 9330, '9330/', '', '配送统计', '/delivery/delivice_statistic', '[]', 0, 1, 2, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9559);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9576, 9506, '/515/9502/9506/', '', '权限', '/accounts/zoneAgent/commissionRecord', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9576);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9595, 9576, '/515/9502/9506/9576/', '', '商圈提成流水列表', 'systemCircleFinancialRecordList', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9595);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9596, 7767, '/110/5125/1120/7767/', '', '系统公告编辑', 'systemNoticeUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9596);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9597, 7767, '/110/5125/1120/7767/', '', '系统公告修改状态', 'systemNoticeSwitchStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9597);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9598, 7767, '/110/5125/1120/7767/', '', '系统公告详情', 'systemNoticeDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9598);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9599, 7767, '/110/5125/1120/7767/', '', '系统公告删除', 'systemNoticeDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9599);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9600, 8451, '9330/8451/', '', '统计列表', 'merchantDeliveryServiceStatisticsList', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9600);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9601, 8451, '9330/8451/', '', '统计详情', 'merchantDeliveryServiceStatisticsDetail', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9601);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9602, 9419, '1661/9404/9419/', '', '列表', 'merchantStaffsStatisticsList', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9602);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9603, 9419, '1661/9404/9419/', '', '详情', 'merchantStaffsStatisticsDetail', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9603);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9604, 8615, '/1119/8615/', '', '详情', 'systemNoticeLogDetail', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9604);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9605, 9503, '/515/9502/9503/', '', '权限', '/accounts/zoneAgent/settlementReview', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9605);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9606, 9605, '/515/9502/9503/9605/', '', '平台结算审核', 'systemCircleCheckoutAudit', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9606);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9607, 9605, '/515/9502/9503/9605/', '', '平台转账', 'systemCircleCheckoutTransfer', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9607);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9608, 9605, '/515/9502/9503/9605/', '', '平台备注', 'systemCircleCheckoutPlatformRemark', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9608);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9641, 9551, '/9551/', '', '商圈代理', '/accounts/zoneAgent', '[]', 0, 1, 1, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9641);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9643, 9641, '/9551/9641/', '', '提成流水', '/accounts/zoneAgent/commissionRecord', '[]', 0, 1, 1, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9643);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9644, 9641, '/9551/9641/', '', '结算账号', '/accounts/zoneAgent/settlementAccount', '[]', 0, 1, 1, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9644);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9645, 9515, '/9515/', '', '权限管理', '/setting', '[]', 0, 1, 1, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9645);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9646, 9645, '/9515/9645/', '', '角色权限', '/setting/systemRole', '[]', 0, 1, 1, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9646);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9647, 9645, '/9515/9645/', '', '管理员管理', '/setting/systemAdmin', '[]', 0, 1, 1, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9647);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9649, 9644, '/9551/9641/9644/', '', '权限', '/accounts/zoneAgent/settlementAccount', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9649);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9651, 9643, '/9551/9641/9643/', '', '权限', '/accounts/zoneAgent/commissionRecord', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9651);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9652, 9646, '/9515/9645/9646/', '', '权限', '/setting/systemRole', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9652);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9653, 9647, '/9515/9645/9647/', '', '权限', '/setting/systemAdmin', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9653);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9654, 9647, '/9515/9645/9647/', '', '附加权限', 'append_/setting/systemAdmin', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9654);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9655, 9649, '/9551/9641/9644/9649/', '', '结算方式get', 'systemCircleAgentGetSettlementMethod', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9655);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9656, 9649, '/9551/9641/9644/9649/', '', '结算方式post', 'systemCircleAgentSetSettlementMethod', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9656);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9661, 9651, '/9551/9641/9643/9651/', '', '商圈提成流水列表', 'systemCircleFinancialRecordList', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9661);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9662, 9652, '/9515/9645/9646/9652/', '', '身份列表', 'systemRoleGetList', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9662);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9663, 9652, '/9515/9645/9646/9652/', '', '身份添加', 'systemRoleCreate', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9663);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9664, 9652, '/9515/9645/9646/9652/', '', '身份编辑', 'systemRoleUpdate', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9664);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9665, 9652, '/9515/9645/9646/9652/', '', '身份修改状态', 'systemRoleStatus', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9665);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9666, 9652, '/9515/9645/9646/9652/', '', '身份删除', 'systemRoleDelete', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9666);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9667, 9653, '/9515/9645/9647/9653/', '', '管理员列表', 'systemAdminLst', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9667);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9668, 9653, '/9515/9645/9647/9653/', '', '管理员修改状态', 'systemAdminStatus', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9668);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9669, 9653, '/9515/9645/9647/9653/', '', '管理员添加', 'systemAdminCreate', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9669);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9670, 9653, '/9515/9645/9647/9653/', '', '管理员编辑', 'systemAdminUpdate', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9670);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9671, 9653, '/9515/9645/9647/9653/', '', '管理员修改密码', 'systemAdminPassword', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9671);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9672, 9653, '/9515/9645/9647/9653/', '', '管理员删除', 'systemAdminDelete', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9672);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9673, 9654, '/9515/9645/9647/9654/', '', '上传图片', 'uploadImage', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9673);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9674, 9654, '/9515/9645/9647/9654/', '', '图片列表', 'systemAttachmentLst', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9674);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9675, 0, '/', 's-shop', '店铺', '/mer', '[]', 95, 1, 1, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9675);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9676, 9878, '/9675/9878/', '', '商户列表', '/merchant/list', '[]', 0, 1, 1, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9676);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9677, 9676, '/9675/9878/9676/', '', '权限', '/merchant/list', '[]', 0, 0, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9677);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9679, 9676, '/9675/9878/9676/', '', '附加权限', 'append_/merchant/list', '[]', 0, 0, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9679);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9682, 0, '/', 's-home', '首页', '/', '[]', 100, 1, 1, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9682);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9685, 0, '/', 's-goods', '商品', '/product', '[]', 90, 1, 1, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9685);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9686, 9685, '/9685/', '', '商品管理', '/product/examine', '[]', 0, 1, 1, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9686);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9687, 9686, '/9685/9686/', '', '权限', '/product/examine', '', 0, 0, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9687);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9688, 0, '/', 's-order', '订单', '/order', '[]', 85, 1, 1, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9688);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9689, 9688, '/9688/', '', '订单列表', '/order/list', '[]', 0, 1, 1, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9689);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9690, 9688, '/9688/', '', '退款订单', '/order/refund', '[]', 0, 1, 1, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9690);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9691, 9688, '/9688/', '', '核销记录', '/order/cancellation', '[]', 0, 1, 1, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9691);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9692, 9689, '/9688/9689/', '', '权限', '/order/list', '', 0, 0, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9692);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9693, 9689, '/9688/9689/', '', '附加权限', 'append_/order/list', '[]', 0, 0, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9693);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9694, 9690, '/9688/9690/', '', '权限', '/order/refund', '[]', 0, 0, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9694);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9696, 9691, '/9688/9691/', '', '权限', '/order/cancellation', '', 0, 0, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9696);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9697, 0, '/', 's-flag', '营销', '/marketing', '[]', 80, 1, 1, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9697);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9698, 9697, '/9697/', '', '秒杀', '/marketing/seckill', '[]', 0, 1, 1, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9698);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9699, 9697, '/9697/', '', '预售', '/marketing/presell', '[]', 0, 1, 1, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9699);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9700, 9697, '/9697/', '', '助力', '/assist', '[]', 0, 1, 1, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9700);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9701, 9697, '/9697/', '', '拼团', '/marketing/combination', '[]', 0, 1, 1, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9701);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9703, 9701, '/9697/9701/', '', '拼团商品列表', '/marketing/combination/combination_goods', '[]', 0, 1, 1, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9703);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9704, 9701, '/9697/9701/', '', '拼团活动列表', '/marketing/combination/combination_list', '[]', 0, 1, 1, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9704);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9705, 9700, '/9697/9700/', '', '活动商品', '/marketing/assist/goods_list', '[]', 0, 1, 1, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9705);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9706, 9700, '/9697/9700/', '', '助力活动', '/marketing/assist/list', '[]', 0, 1, 1, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9706);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9707, 9699, '/9697/9699/', '', '预售商品', '/marketing/presell/list', '[]', 0, 1, 1, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9707);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9710, 9698, '/9697/9698/', '', '秒杀管理', '/marketing/seckill/list', '[]', 0, 1, 1, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9710);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9711, 9698, '/9697/9698/', '', '秒杀活动', '/marketing/seckill/store_seckill/list', '[]', 0, 1, 1, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9711);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9714, 9710, '/9697/9698/9710/', '', '权限', '/marketing/seckill/list', '', 0, 0, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9714);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9715, 9711, '/9697/9698/9711/', '', '权限', '/marketing/seckill/store_seckill/list', '', 0, 0, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9715);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9716, 9707, '/9697/9699/9707/', '', '权限', '/marketing/presell/list', '', 0, 0, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9716);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9718, 9705, '/9697/9700/9705/', '', '权限', '/marketing/assist/goods_list', '', 0, 0, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9718);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9719, 9706, '/9697/9700/9706/', '', '权限', '/marketing/assist/list', '', 0, 0, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9719);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9720, 9703, '/9697/9701/9703/', '', '权限', '/marketing/combination/combination_goods', '', 0, 0, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9720);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9721, 9704, '/9697/9701/9704/', '', '权限', '/marketing/combination/combination_list', '', 0, 0, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9721);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9722, 9690, '/9688/9690/', '', '附加权限', 'append_/order/refund', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9722);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9723, 9716, '/9697/9699/9707/9716/', '', '列表', 'systemStoreProductPresellLst', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9723);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9724, 9716, '/9697/9699/9707/9716/', '', '显示/隐藏', 'systemStoreProductPresellShow', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9724);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9725, 9716, '/9697/9699/9707/9716/', '', '详情', 'systemStoreProductPresellDetail', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9725);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9726, 9716, '/9697/9699/9707/9716/', '', '编辑数据', 'systemStoreProductPresellGet', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9726);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9727, 9716, '/9697/9699/9707/9716/', '', '编辑', 'systemStoreProductPresellUpdate', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9727);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9728, 9716, '/9697/9699/9707/9716/', '', '审核', 'systemStoreProductPresellSwitchStatus', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9728);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9729, 9716, '/9697/9699/9707/9716/', '', '设置标签', 'systemStoreProductPresellLabels', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9729);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9730, 9718, '/9697/9700/9705/9718/', '', '列表', 'systemStoreProductAssistLst', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9730);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9731, 9718, '/9697/9700/9705/9718/', '', '显示/隐藏', 'systemStoreProductAssistShow', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9731);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9732, 9718, '/9697/9700/9705/9718/', '', '详情', 'systemStoreProductAssistDetail', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9732);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9733, 9718, '/9697/9700/9705/9718/', '', '编辑', 'systemStoreProductAssistProductUpdate', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9733);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9734, 9718, '/9697/9700/9705/9718/', '', '审核', 'systemStoreProductAssistStatus', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9734);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9735, 9718, '/9697/9700/9705/9718/', '', '编辑数据', 'systemStoreProductAssistGet', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9735);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9736, 9718, '/9697/9700/9705/9718/', '', '设置标签', 'systemStoreProductAssistLabels', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9736);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9737, 9719, '/9697/9700/9706/9719/', '', '列表', 'systemStoreProductAssistSetLst', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9737);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9738, 9719, '/9697/9700/9706/9719/', '', '详情', 'systemStoreProductAssistSetDetail', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9738);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9739, 9720, '/9697/9701/9703/9720/', '', '列表', 'systemStoreProductGroupLst', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9739);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9740, 9720, '/9697/9701/9703/9720/', '', '显示/隐藏', 'systemStoreProductGroupShow', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9740);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9741, 9720, '/9697/9701/9703/9720/', '', '详情', 'systemStoreProductGroupDetail', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9741);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9742, 9720, '/9697/9701/9703/9720/', '', '编辑', 'systemStoreProductGroupProductUpdate', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9742);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9743, 9720, '/9697/9701/9703/9720/', '', '审核', 'systemStoreProductGroupStatus', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9743);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9744, 9720, '/9697/9701/9703/9720/', '', '编辑数据', 'systemStoreProductGroupGet', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9744);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9745, 9720, '/9697/9701/9703/9720/', '', '排序', 'systemStoreProductGroupSort', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9745);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9746, 9720, '/9697/9701/9703/9720/', '', '设置标签', 'systemStoreProductGroupLabels', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9746);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9747, 9721, '/9697/9701/9704/9721/', '', '列表', 'systemStoreProductGroupBuyingLst', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9747);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9748, 9721, '/9697/9701/9704/9721/', '', '详情', 'systemStoreProductGroupBuyingDetail', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9748);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9749, 9715, '/9697/9698/9711/9715/', '', '列表', 'systemSeckillActiveGetActiveList', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9749);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9750, 9715, '/9697/9698/9711/9715/', '', '详情', 'systemSeckillActiveGetActiveInfo', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9750);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9751, 9715, '/9697/9698/9711/9715/', '', '创建', 'systemSeckillActiveCreateActive', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9751);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9752, 9715, '/9697/9698/9711/9715/', '', '编辑', 'systemSeckillActiveUpdateActive', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9752);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9753, 9715, '/9697/9698/9711/9715/', '', '编辑状态', 'systemSeckillActiveUpdateActiveStatus', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9753);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9754, 9715, '/9697/9698/9711/9715/', '', '删除', 'systemSeckillActiveDeleteActive', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9754);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9755, 9715, '/9697/9698/9711/9715/', '', '活动统计数据面板', 'systemSeckillActiveChartPanel', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9755);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9756, 9715, '/9697/9698/9711/9715/', '', '活动参与人统计列表', 'systemSeckillActiveChartPeople', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9756);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9757, 9715, '/9697/9698/9711/9715/', '', '活动订单统计列表', 'systemSeckillActiveChartOrder', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9757);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9758, 9715, '/9697/9698/9711/9715/', '', '活动商品统计列表', 'systemSeckillActiveChartProduct', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9758);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9759, 9714, '/9697/9698/9710/9714/', '', '统计', 'systemStoreSeckillProductLstFilter', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9759);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9760, 9714, '/9697/9698/9710/9714/', '', '列表', 'systemStoreSeckillProductPageLst', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9760);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9761, 9714, '/9697/9698/9710/9714/', '', '列表', 'systemStoreSeckillProductLst', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9761);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9762, 9714, '/9697/9698/9710/9714/', '', '详情', 'systemStoreSeckillProductDetail', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9762);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9763, 9714, '/9697/9698/9710/9714/', '', '编辑', 'systemStoreSeckillProductUpdate', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9763);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9764, 9714, '/9697/9698/9710/9714/', '', '审核', 'systemStoreSeckillProductSwitchStatus', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9764);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9765, 9714, '/9697/9698/9710/9714/', '', '显示/隐藏', 'systemStoreSeckillProductChangeUsed', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9765);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9766, 9714, '/9697/9698/9710/9714/', '', '设置标签', 'systemStoreSeckillProductLabels', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9766);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9767, 9714, '/9697/9698/9710/9714/', '', '加入回收站', 'systemStoreSeckillProductDelete', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9767);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9768, 9714, '/9697/9698/9710/9714/', '', '删除', 'systemStoreSeckillProductDestory', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9768);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9769, 9677, '/9675/9878/9676/9677/', '', '虚拟关注量', 'systemMerchantCareFicti', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9769);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9770, 9677, '/9675/9878/9676/9677/', '', '商户列表', 'systemMerchantCreateForm', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9770);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9771, 9677, '/9675/9878/9676/9677/', '', '商户列表统计', 'systemMerchantCount', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9771);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9772, 9677, '/9675/9878/9676/9677/', '', '商户列表', 'systemMerchantLst', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9772);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9773, 9677, '/9675/9878/9676/9677/', '', '商户添加', 'systemMerchantCreate', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9773);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9774, 9677, '/9675/9878/9676/9677/', '', '商户编辑', 'systemMerchantUpdate', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9774);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9775, 9677, '/9675/9878/9676/9677/', '', '商户修改推荐', 'systemMerchantStatus', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9775);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9776, 9677, '/9675/9878/9676/9677/', '', '商户开启/关闭', 'systemMerchantClose', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9776);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9777, 9677, '/9675/9878/9676/9677/', '', '商户删除', 'systemMerchantDelete', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9777);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9778, 9677, '/9675/9878/9676/9677/', '', '商户修改密码', 'systemMerchantAdminPassword', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9778);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9779, 9677, '/9675/9878/9676/9677/', '', '商户登录', 'systemMerchantLogin', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9779);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9780, 9677, '/9675/9878/9676/9677/', '', '修改采集商品次数', 'systemMerchantChangeCopy', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9780);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9781, 9677, '/9675/9878/9676/9677/', '', '详情', 'systemMerchantDetail', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9781);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9782, 9677, '/9675/9878/9676/9677/', '', '操作日志', 'systemMerchantOperateList', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9782);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9783, 9679, '/9675/9878/9676/9679/', '', '上传图片', 'uploadImage', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9783);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9784, 9679, '/9675/9878/9676/9679/', '', '图片列表', 'systemAttachmentLst', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9784);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9785, 9692, '/9688/9689/9692/', '', '列表', 'systemOrderLst', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9785);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9786, 9692, '/9688/9689/9692/', '', '金额统计', 'systemOrderStat', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9786);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9787, 9692, '/9688/9689/9692/', '', '快递查询', 'systemOrderExpress', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9787);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9788, 9692, '/9688/9689/9692/', '', '头部统计', 'systemOrderTitle', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9788);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9789, 9692, '/9688/9689/9692/', '', '详情', 'systemOrderDetail', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9789);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9790, 9692, '/9688/9689/9692/', '', '导出', 'systemOrderExcel', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9790);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9791, 9692, '/9688/9689/9692/', '', '记录', 'systemOrderStatus', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9791);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9792, 9692, '/9688/9689/9692/', '', '关联订单', 'systemOrderChildrenList', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9792);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9793, 9693, '/9688/9689/9693/', '', '导出列表', 'systemStoreExcelLst', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9793);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9794, 9693, '/9688/9689/9693/', '', '导出列表', 'systemStoreExcelDownload', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9794);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9795, 9696, '/9688/9691/9696/', '', '核销', 'systemOrderTakeStat', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9795);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9796, 9696, '/9688/9691/9696/', '', '核销订单', 'systemTakeOrderLst', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9796);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9797, 9696, '/9688/9691/9696/', '', '头部统计', 'systemTakeOrderTitle', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9797);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9798, 9694, '/9688/9690/9694/', '', '列表', 'systemRefundOrderLst', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9798);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9799, 9694, '/9688/9690/9694/', '', '详情', 'systemRefundOrderDetail', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9799);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9800, 9694, '/9688/9690/9694/', '', '日志', 'systemRefundOrderLog', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9800);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9801, 9694, '/9688/9690/9694/', '', '审核', 'systemRefundOrderApprove', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9801);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9802, 9694, '/9688/9690/9694/', '', '导出', 'systemRefundOrderExcel', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9802);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9803, 9722, '/9688/9690/9722/', '', '导出列表', 'systemStoreExcelLst', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9803);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9804, 9722, '/9688/9690/9722/', '', '导出下载', 'systemStoreExcelDownload', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9804);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9805, 9687, '/9685/9686/9687/', '', '统计', 'systemStoreProductLstFilter', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9805);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9806, 9687, '/9685/9686/9687/', '', '列表', 'systemStoreProductLst', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9806);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9807, 9687, '/9685/9686/9687/', '', '详情', 'systemStoreProductDetail', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9807);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9808, 9687, '/9685/9686/9687/', '', '编辑', 'systemStoreProductUpdate', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9808);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9809, 9687, '/9685/9686/9687/', '', '上下架', 'systemStoreProductSwitchStatus', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9809);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9810, 9687, '/9685/9686/9687/', '', '批量上下架', 'systemStoreProductSwitchBatchStatus', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9810);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9811, 9687, '/9685/9686/9687/', '', '批量设置标签', 'systemStoreProductSwitchBatchLabels', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9811);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9812, 9687, '/9685/9686/9687/', '', '批量设置推荐', 'systemStoreProductSwitchBatchHot', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9812);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9813, 9687, '/9685/9686/9687/', '', '批量设置分类推荐', 'systemStoreProductSwitchBatchCateHot', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9813);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9814, 9687, '/9685/9686/9687/', '', '分销状态变更商品检测', 'systemStoreProductCheck', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9814);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9815, 9687, '/9685/9686/9687/', '', '显示/隐藏', 'systemStoreProductChangeUsed', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9815);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9816, 9687, '/9685/9686/9687/', '', '虚拟销量', 'systemStoreProductAddFicti', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9816);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9817, 9687, '/9685/9686/9687/', '', '设置标签', 'systemStoreProductLabels', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9817);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9818, 9687, '/9685/9686/9687/', '', '获取商品操作记录', 'systemStoreProductGetOperateList', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9818);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9819, 9687, '/9685/9686/9687/', '', '获取自营商品列表', 'systemStoreProductGetSelfProductList', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9819);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9833, 9682, '/9682/', '', '商品统计', '/statistic/product', '[]', 0, 1, 1, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9833);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9834, 9682, '/9682/', '', '订单统计', '/statistic/order', '[]', 0, 1, 1, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9834);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9835, 9489, '/9218/9489/', '', '权限', '/statistic/order', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9835);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9836, 9490, '/9218/9490/', '', '权限', '/statistic/product', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9836);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9837, 9491, '/9218/9491/', '', '权限', '/statistic/member', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9837);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9838, 9835, '/9218/9489/9835/', '', '顶部统计', 'systemAnalyticsOrderTop', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9838);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9839, 9835, '/9218/9489/9835/', '', '折线图统计', 'systemAnalyticsOrderLineChart', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9839);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9840, 9835, '/9218/9489/9835/', '', '折线图统计', 'systemAnalyticsOrderTypePieChart', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9840);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9841, 9836, '/9218/9490/9836/', '', '顶部统计', 'systemAnalyticsProductTop', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9841);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9842, 9836, '/9218/9490/9836/', '', '折线图统计', 'systemAnalyticsProductLineChart', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9842);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9843, 9836, '/9218/9490/9836/', '', '折线图统计', 'systemAnalyticsProductTypePieChart', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9843);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9844, 9837, '/9218/9491/9837/', '', '顶部统计', 'systemAnalyticsUserTop', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9844);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9845, 9837, '/9218/9491/9837/', '', '折线图统计', 'systemAnalyticsUserLineChart', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9845);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9846, 9837, '/9218/9491/9837/', '', '折线图统计', 'systemAnalyticsUserTypePieChart', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9846);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9847, 9834, '/9682/9834/', '', '权限', '/statistic/order', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9847);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9848, 9833, '/9682/9833/', '', '权限', '/statistic/product', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9848);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9849, 9847, '/9682/9834/9847/', '', '顶部统计', 'systemAnalyticsOrderTop', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9849);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9850, 9847, '/9682/9834/9847/', '', '折线图统计', 'systemAnalyticsOrderLineChart', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9850);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9851, 9847, '/9682/9834/9847/', '', '折线图统计', 'systemAnalyticsOrderTypePieChart', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9851);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9852, 9848, '/9682/9833/9848/', '', '顶部统计', 'systemAnalyticsProductTop', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9852);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9853, 9848, '/9682/9833/9848/', '', '折线图统计', 'systemAnalyticsProductLineChart', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9853);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9854, 9848, '/9682/9833/9848/', '', '折线图统计', 'systemAnalyticsProductTypePieChart', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9854);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9862, 9501, '55/9501/', '', '权限', '/statistic/order', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9862);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9863, 9500, '55/9500/', '', '权限', '/statistic/product', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9863);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9864, 9862, '55/9501/9862/', '', '顶部统计', 'merchantAnalyticsOrderTop', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9864);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9865, 9862, '55/9501/9862/', '', '折线图统计', 'merchantAnalyticsOrderLineChart', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9865);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9866, 9862, '55/9501/9862/', '', '折线图统计', 'merchantAnalyticsOrderTypePieChart', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9866);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9867, 9863, '55/9500/9863/', '', '顶部统计', 'merchantAnalyticsProductTop', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9867);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9868, 9863, '55/9500/9863/', '', '折线图统计', 'merchantAnalyticsProductLineChart', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9868);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9869, 9863, '55/9500/9863/', '', '折线图统计', 'merchantAnalyticsProductTypePieChart', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9869);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9873, 20042, '/9492/42/', '', '商户管理员', '/merchant/admin-list', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9873);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9876, 20042, '/9492/42/', '', '商户入驻审核', '/merchant/review', '[]', 20, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9876);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9877, 20042, '/9492/42/', '', '商户设置', '/merchant/apply-setting', '[]', 0, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9877);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9878, 9675, '/9675/', '', '店铺管理', '/mer/mer', '[]', 0, 1, 1, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9878);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9880, 9878, '/9675/9878/', '', '店铺入驻申请', '/merchant/application', '[]', 0, 1, 1, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9880);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9881, 9880, '/9675/9878/9880/', '', '权限', '/merchant/application', '', 0, 0, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9881);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9885, 7786, '/110/38/39/7786/', '', '身份列表', 'systemOrganizationRoleGetList', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9885);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9886, 7786, '/110/38/39/7786/', '', '身份添加', 'systemOrganizationRoleCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9886);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9887, 7786, '/110/38/39/7786/', '', '身份编辑', 'systemOrganizationRoleUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9887);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9888, 7786, '/110/38/39/7786/', '', '身份修改状态', 'systemOrganizationRoleStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9888);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9889, 7786, '/110/38/39/7786/', '', '身份删除', 'systemOrganizationRoleDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9889);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9890, 7786, '/110/38/39/7786/', '', '身份选项', 'systemOrganizationRoleOptions', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9890);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9891, 7778, '/87/539/7778/', '', '批量复制商品到店铺', 'systemStoreProductBatchCopy', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9891);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9892, 9558, '9330/9558/', '', '权限', '/delivery/personnel_manage', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9892);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9893, 8557, '/525/1304/8557/', '', '上传视频', 'merchantApplymentsUploadVideo', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9893);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9894, 9892, '9330/9558/9892/', '', '列表', 'merchantDeliveryServiceLst', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9894);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9895, 9892, '9330/9558/9892/', '', '修改状态', 'merchantDeliveryServiceStatus', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9895);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9896, 9892, '9330/9558/9892/', '', '添加', 'merchantDeliveryServiceCreate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9896);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9897, 9892, '9330/9558/9892/', '', '编辑', 'merchantDeliveryServiceUpdate', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9897);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9898, 9892, '9330/9558/9892/', '', '删除', 'merchantDeliveryServiceDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9898);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9899, 9892, '9330/9558/9892/', '', '统计列表', 'merchantDeliveryServiceStatisticsList', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9899);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9900, 9892, '9330/9558/9892/', '', '统计详情', 'merchantDeliveryServiceStatisticsDetail', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9900);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9906, 9652, '/9515/9645/9646/9652/', '', '身份列表', 'systemOrganizationRoleGetList', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9906);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9907, 9652, '/9515/9645/9646/9652/', '', '身份添加', 'systemOrganizationRoleCreate', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9907);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9908, 9652, '/9515/9645/9646/9652/', '', '身份编辑', 'systemOrganizationRoleUpdate', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9908);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9909, 9652, '/9515/9645/9646/9652/', '', '身份修改状态', 'systemOrganizationRoleStatus', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9909);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9910, 9652, '/9515/9645/9646/9652/', '', '身份删除', 'systemOrganizationRoleDelete', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9910);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9911, 9652, '/9515/9645/9646/9652/', '', '身份选项', 'systemOrganizationRoleOptions', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9911);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9912, 9881, '/9675/9878/9880/9881/', '', '列表', 'systemMerchantIntentionLst', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9912);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9913, 9881, '/9675/9878/9880/9881/', '', '审核', 'systemMerchantIntentionStatus', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9913);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9914, 9881, '/9675/9878/9880/9881/', '', '删除', 'systemMerchantIntentionDelete', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9914);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9915, 9881, '/9675/9878/9880/9881/', '', '备注', 'systemMerchantIntentionMark', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9915);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9916, 9687, '/9685/9686/9687/', '', '批量复制商品到店铺', 'systemStoreProductBatchCopy', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9916);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9917, 20042, '/9492/42/', '', '商户列表', '/merchant/index', '[]', 30, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9917);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9920, 9493, '/9492/9493/', '', '区域列表', '/business-zones/index', '[]', 50, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9920);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9921, 0, '/', 's-shop', '店铺', '/mer', '[]', 95, 1, 1, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9921);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9922, 9921, '/9921/', '', '店铺管理', '/mer/mer', '[]', 0, 1, 1, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9922);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9923, 9922, '/9921/9922/', '', '店铺列表', '/merchant/list', '[]', 0, 1, 1, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9923);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9924, 9920, '/9492/9493/9920/', '', '权限', '/business-zones/index', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9924);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9925, 9917, '/9492/42/9917/', '', '权限', '/merchant/index', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9925);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9926, 9495, '/9492/9493/9495/', '', '权限', '/business-zones/agents', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9926);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9927, 9873, '/9492/42/9873/', '', '权限', '/merchant/admin-list', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9927);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9928, 9496, '/9492/9493/9496/', '', '权限', '/business-zones/agent-review', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9928);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9929, 9876, '/9492/42/9876/', '', '权限', '/merchant/review', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9929);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9930, 9367, '/9492/6370/9367/', '', '权限', '/merchant/grouping', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9930);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9931, 9924, '/9492/9493/9920/9924/', '', '商圈列表', 'systemCircleList', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9931);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9932, 9924, '/9492/9493/9920/9924/', '', '商圈详情', 'systemCircleDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9932);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9933, 9924, '/9492/9493/9920/9924/', '', '商圈添加', 'systemCircleCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9933);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9934, 9924, '/9492/9493/9920/9924/', '', '商圈编辑', 'systemCircleUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9934);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9935, 9924, '/9492/9493/9920/9924/', '', '商圈删除', 'systemCircleDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9935);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9936, 9924, '/9492/9493/9920/9924/', '', '商圈状态切换', 'systemCircleSwitch', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9936);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9937, 9924, '/9492/9493/9920/9924/', '', '关联商户列表', 'systemCircleMerchantList', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9937);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9938, 9925, '/9492/42/9917/9925/', '', '商圈列表', 'systemCircleList', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9938);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9939, 9925, '/9492/42/9917/9925/', '', '商圈详情', 'systemCircleDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9939);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9940, 9925, '/9492/42/9917/9925/', '', '商圈添加', 'systemCircleCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9940);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9941, 9925, '/9492/42/9917/9925/', '', '商圈编辑', 'systemCircleUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9941);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9942, 9925, '/9492/42/9917/9925/', '', '商圈删除', 'systemCircleDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9942);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9943, 9925, '/9492/42/9917/9925/', '', '商圈状态切换', 'systemCircleSwitch', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9943);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9944, 9925, '/9492/42/9917/9925/', '', '关联商户列表', 'systemCircleMerchantList', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9944);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9945, 9926, '/9492/9493/9495/9926/', '', '商圈代理添加', 'systemCircleAgentCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9945);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9946, 9926, '/9492/9493/9495/9926/', '', '商圈代理编辑', 'systemCircleAgentUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9946);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9947, 9926, '/9492/9493/9495/9926/', '', '商圈代理删除', 'systemCircleAgentDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9947);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9948, 9926, '/9492/9493/9495/9926/', '', '关联商户列表', 'systemCircleAgentMerchant', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9948);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9949, 9926, '/9492/9493/9495/9926/', '', '代理选项', 'systemCircleAgentOptions', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9949);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9950, 9926, '/9492/9493/9495/9926/', '', '重置密码', 'systemCircleAgentResetPassword', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9950);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9951, 9927, '/9492/42/9873/9927/', '', '商圈代理添加', 'systemCircleAgentCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9951);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9952, 9927, '/9492/42/9873/9927/', '', '商圈代理编辑', 'systemCircleAgentUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9952);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9953, 9927, '/9492/42/9873/9927/', '', '商圈代理删除', 'systemCircleAgentDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9953);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9954, 9927, '/9492/42/9873/9927/', '', '关联商户列表', 'systemCircleAgentMerchant', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9954);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9955, 9927, '/9492/42/9873/9927/', '', '代理选项', 'systemCircleAgentOptions', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9955);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9956, 9927, '/9492/42/9873/9927/', '', '重置密码', 'systemCircleAgentResetPassword', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9956);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9957, 9928, '/9492/9493/9496/9928/', '', '商圈代理审核', 'systemCircleAgentAudit', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9957);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9958, 9929, '/9492/42/9876/9929/', '', '商圈代理审核', 'systemCircleAgentAudit', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9958);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9959, 9930, '/9492/6370/9367/9930/', '', '列表', 'systemStoreGroupLst', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9959);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9960, 9930, '/9492/6370/9367/9930/', '', '详情', 'systemStoreGroupDetail', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9960);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9961, 9930, '/9492/6370/9367/9930/', '', '添加', 'systemStoreGroupCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9961);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9962, 9930, '/9492/6370/9367/9930/', '', '编辑', 'systemStoreGroupUpdate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9962);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9963, 9930, '/9492/6370/9367/9930/', '', '删除', 'systemStoreGroupDelete', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9963);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9964, 9930, '/9492/6370/9367/9930/', '', '状态切换', 'systemStoreGroupSwitchStatus', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9964);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9965, 9930, '/9492/6370/9367/9930/', '', '设置店铺分组模板', 'systemStoreGroupSetTemplate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9965);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9966, 9930, '/9492/6370/9367/9930/', '', '关联店铺列表', 'systemStoreGroupStores', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9966);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9967, 7763, '/9492/6370/44/7763/', '', '商户添加店铺', 'systemMerchantBusinessCreate', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9967);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9994, 9923, '/9921/9922/9923/', '', '权限', '/merchant/list', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9994);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9995, 9923, '/9921/9922/9923/', '', '附加权限', 'append_/merchant/list', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9995);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9996, 0, '/', '', '权限', 'self', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9996);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9997, 9994, '/9921/9922/9923/9994/', '', '虚拟关注量', 'systemMerchantCareFicti', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9997);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9998, 9994, '/9921/9922/9923/9994/', '', '店铺列表', 'systemMerchantCreateForm', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9998);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 9999, 9994, '/9921/9922/9923/9994/', '', '店铺列表统计', 'systemMerchantCount', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 9999);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10000, 9994, '/9921/9922/9923/9994/', '', '店铺列表', 'systemMerchantLst', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10000);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10001, 9994, '/9921/9922/9923/9994/', '', '店铺添加', 'systemMerchantCreate', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10001);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10002, 9994, '/9921/9922/9923/9994/', '', '店铺编辑', 'systemMerchantUpdate', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10002);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10003, 9994, '/9921/9922/9923/9994/', '', '店铺修改推荐', 'systemMerchantStatus', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10003);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10004, 9994, '/9921/9922/9923/9994/', '', '店铺开启/关闭', 'systemMerchantClose', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10004);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10005, 9994, '/9921/9922/9923/9994/', '', '店铺删除', 'systemMerchantDelete', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10005);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10006, 9994, '/9921/9922/9923/9994/', '', '店铺修改密码', 'systemMerchantAdminPassword', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10006);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10007, 9994, '/9921/9922/9923/9994/', '', '店铺登录', 'systemMerchantLogin', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10007);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10008, 9994, '/9921/9922/9923/9994/', '', '修改采集商品次数', 'systemMerchantChangeCopy', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10008);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10009, 9994, '/9921/9922/9923/9994/', '', '详情', 'systemMerchantDetail', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10009);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10010, 9994, '/9921/9922/9923/9994/', '', '操作日志', 'systemMerchantOperateList', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10010);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10011, 9994, '/9921/9922/9923/9994/', '', '商户添加店铺', 'systemMerchantBusinessCreate', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10011);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10012, 9995, '/9921/9922/9923/9995/', '', '上传图片', 'uploadImage', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10012);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10013, 9995, '/9921/9922/9923/9995/', '', '图片列表', 'systemAttachmentLst', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10013);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10014, 9996, '/9996/', '', '修改信息', 'systemAdminEdit', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10014);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10015, 9996, '/9996/', '', '修改密码', 'systemAdminEditPassword', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10015);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10016, 0, '/', 's-home', '首页', '/', '[]', 100, 1, 1, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10016);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10017, 10016, '/10016/', '', '商品统计', '/statistic/product', '[]', 0, 1, 1, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10017);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10018, 10016, '/10016/', '', '订单统计', '/statistic/order', '[]', 0, 1, 1, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10018);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10019, 0, '/', 's-goods', '商品', '/product', '[]', 90, 1, 1, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10019);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10020, 10019, '/10019/', '', '商品管理', '/product/examine', '[]', 0, 1, 1, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10020);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10021, 0, '/', 's-order', '订单', '/order', '[]', 85, 1, 1, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10021);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10022, 10021, '/10021/', '', '订单列表', '/order/list', '[]', 0, 1, 1, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10022);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10023, 10021, '/10021/', '', '退款订单', '/order/refund', '[]', 0, 1, 1, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10023);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10024, 10021, '/10021/', '', '核销记录', '/order/cancellation', '[]', 0, 1, 1, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10024);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10025, 0, '/', 's-flag', '营销', '/marketing', '[]', 80, 1, 1, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10025);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10026, 10025, '/10025/', '', '秒杀', '/marketing/seckill', '[]', 0, 1, 1, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10026);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10027, 10025, '/10025/', '', '预售', '/marketing/presell', '[]', 0, 1, 1, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10027);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10028, 10025, '/10025/', '', '助力', '/assist', '[]', 0, 1, 1, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10028);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10029, 10025, '/10025/', '', '拼团', '/marketing/combination', '[]', 0, 1, 1, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10029);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10030, 9641, '/9551/9641/', '', '申请结算', '/accounts/zoneAgent/settlementApply', '[]', 2, 1, 1, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10030);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10031, 0, '/', 's-data', '财务', '/accounts', '[]', 70, 1, 1, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10031);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10032, 0, '/', 's-tools', '设置', '/settings', '[]', 0, 1, 1, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10032);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10033, 10032, '/10032/', '', '权限管理', '/setting', '[]', 0, 1, 1, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10033);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10034, 10033, '/10032/10033/', '', '角色权限', '/setting/systemRole', '[]', 0, 1, 1, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10034);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10035, 10033, '/10032/10033/', '', '管理员管理', '/setting/systemAdmin', '[]', 0, 1, 1, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10035);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10036, 10026, '/10025/10026/', '', '秒杀管理', '/marketing/seckill/list', '[]', 0, 1, 1, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10036);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10037, 10026, '/10025/10026/', '', '秒杀活动', '/marketing/seckill/store_seckill/list', '[]', 0, 1, 1, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10037);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10038, 10027, '/10025/10027/', '', '预售商品', '/marketing/presell/list', '[]', 0, 1, 1, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10038);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10039, 10028, '/10025/10028/', '', '活动商品', '/marketing/assist/goods_list', '[]', 0, 1, 1, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10039);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10040, 10028, '/10025/10028/', '', '助力活动', '/marketing/assist/list', '[]', 0, 1, 1, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10040);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10041, 10029, '/10025/10029/', '', '拼团商品列表', '/marketing/combination/combination_goods', '[]', 0, 1, 1, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10041);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10042, 10029, '/10025/10029/', '', '拼团活动列表', '/marketing/combination/combination_list', '[]', 0, 1, 1, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10042);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10043, 10031, '/10031/', '', '店铺结算', '/mer/accounts', '[]', 0, 1, 1, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10043);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10044, 10043, '/10031/10043/', '', '店铺账单', '/accounts/merchantBill', '[]', 0, 1, 1, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10044);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10045, 10030, '/9551/9641/10030/', '', '权限', '/accounts/zoneAgent/settlementApply', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10045);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10046, 10045, '/9551/9641/10030/10045/', '', '商圈申请结算获取余额', 'systemCircleCheckoutCreate', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10046);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10047, 10045, '/9551/9641/10030/10045/', '', '商圈申请结算提交', 'systemCircleCheckoutSave', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10047);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10048, 10045, '/9551/9641/10030/10045/', '', '商圈撤销结算', 'systemCircleCheckoutRevoke', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10048);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10049, 10045, '/9551/9641/10030/10045/', '', '商圈备注', 'systemCircleCheckoutRemark', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10049);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10050, 9677, '/9675/9878/9676/9677/', '', '商户添加店铺', 'systemMerchantBusinessCreate', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10050);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10051, 10044, '/10031/10043/10044/', '', '权限', '/accounts/merchantBill', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10051);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10052, 10018, '/10016/10018/', '', '权限', '/statistic/order', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10052);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10053, 10017, '/10016/10017/', '', '权限', '/statistic/product', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10053);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10054, 10038, '/10025/10027/10038/', '', '权限', '/marketing/presell/list', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10054);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10055, 10039, '/10025/10028/10039/', '', '权限', '/marketing/assist/goods_list', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10055);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10056, 10040, '/10025/10028/10040/', '', '权限', '/marketing/assist/list', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10056);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10057, 10041, '/10025/10029/10041/', '', '权限', '/marketing/combination/combination_goods', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10057);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10058, 10042, '/10025/10029/10042/', '', '权限', '/marketing/combination/combination_list', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10058);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10059, 10037, '/10025/10026/10037/', '', '权限', '/marketing/seckill/store_seckill/list', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10059);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10060, 10036, '/10025/10026/10036/', '', '权限', '/marketing/seckill/list', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10060);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10061, 10022, '/10021/10022/', '', '权限', '/order/list', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10061);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10062, 10022, '/10021/10022/', '', '附加权限', 'append_/order/list', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10062);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10063, 10024, '/10021/10024/', '', '权限', '/order/cancellation', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10063);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10064, 10023, '/10021/10023/', '', '权限', '/order/refund', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10064);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10065, 10023, '/10021/10023/', '', '附加权限', 'append_/order/refund', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10065);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10066, 10020, '/10019/10020/', '', '权限', '/product/examine', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10066);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10067, 10034, '/10032/10033/10034/', '', '权限', '/setting/systemRole', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10067);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10068, 10035, '/10032/10033/10035/', '', '权限', '/setting/systemAdmin', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10068);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10069, 10035, '/10032/10033/10035/', '', '附加权限', 'append_/setting/systemAdmin', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10069);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10070, 10051, '/10031/10043/10044/10051/', '', '商户列表', 'systemFinancialRecordMerLst', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10070);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10071, 10051, '/10031/10043/10044/10051/', '', '商户统计', 'systemFinancialRecordMerAcountsLst', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10071);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10072, 10051, '/10031/10043/10044/10051/', '', '商户财务头部统计', 'systemFinancialRecordMerTitle', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10072);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10073, 10051, '/10031/10043/10044/10051/', '', '商户财务详情', 'systemFinancialRecordMerDetail', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10073);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10074, 10051, '/10031/10043/10044/10051/', '', '商户财务导出', 'systemFinancialRecordMerExcel', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10074);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10075, 10052, '/10016/10018/10052/', '', '顶部统计', 'systemAnalyticsOrderTop', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10075);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10076, 10052, '/10016/10018/10052/', '', '折线图统计', 'systemAnalyticsOrderLineChart', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10076);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10077, 10052, '/10016/10018/10052/', '', '折线图统计', 'systemAnalyticsOrderTypePieChart', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10077);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10078, 10053, '/10016/10017/10053/', '', '顶部统计', 'systemAnalyticsProductTop', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10078);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10079, 10053, '/10016/10017/10053/', '', '折线图统计', 'systemAnalyticsProductLineChart', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10079);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10080, 10053, '/10016/10017/10053/', '', '折线图统计', 'systemAnalyticsProductTypePieChart', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10080);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10081, 10054, '/10025/10027/10038/10054/', '', '列表', 'systemStoreProductPresellLst', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10081);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10082, 10054, '/10025/10027/10038/10054/', '', '显示/隐藏', 'systemStoreProductPresellShow', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10082);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10083, 10054, '/10025/10027/10038/10054/', '', '详情', 'systemStoreProductPresellDetail', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10083);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10084, 10054, '/10025/10027/10038/10054/', '', '编辑数据', 'systemStoreProductPresellGet', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10084);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10085, 10054, '/10025/10027/10038/10054/', '', '编辑', 'systemStoreProductPresellUpdate', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10085);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10086, 10054, '/10025/10027/10038/10054/', '', '审核', 'systemStoreProductPresellSwitchStatus', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10086);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10087, 10054, '/10025/10027/10038/10054/', '', '设置标签', 'systemStoreProductPresellLabels', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10087);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10088, 10055, '/10025/10028/10039/10055/', '', '列表', 'systemStoreProductAssistLst', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10088);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10089, 10055, '/10025/10028/10039/10055/', '', '显示/隐藏', 'systemStoreProductAssistShow', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10089);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10090, 10055, '/10025/10028/10039/10055/', '', '详情', 'systemStoreProductAssistDetail', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10090);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10091, 10055, '/10025/10028/10039/10055/', '', '编辑', 'systemStoreProductAssistProductUpdate', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10091);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10092, 10055, '/10025/10028/10039/10055/', '', '审核', 'systemStoreProductAssistStatus', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10092);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10093, 10055, '/10025/10028/10039/10055/', '', '编辑数据', 'systemStoreProductAssistGet', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10093);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10094, 10055, '/10025/10028/10039/10055/', '', '设置标签', 'systemStoreProductAssistLabels', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10094);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10095, 10056, '/10025/10028/10040/10056/', '', '列表', 'systemStoreProductAssistSetLst', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10095);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10096, 10056, '/10025/10028/10040/10056/', '', '详情', 'systemStoreProductAssistSetDetail', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10096);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10097, 10057, '/10025/10029/10041/10057/', '', '列表', 'systemStoreProductGroupLst', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10097);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10098, 10057, '/10025/10029/10041/10057/', '', '显示/隐藏', 'systemStoreProductGroupShow', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10098);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10099, 10057, '/10025/10029/10041/10057/', '', '详情', 'systemStoreProductGroupDetail', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10099);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10100, 10057, '/10025/10029/10041/10057/', '', '编辑', 'systemStoreProductGroupProductUpdate', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10100);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10101, 10057, '/10025/10029/10041/10057/', '', '审核', 'systemStoreProductGroupStatus', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10101);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10102, 10057, '/10025/10029/10041/10057/', '', '编辑数据', 'systemStoreProductGroupGet', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10102);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10103, 10057, '/10025/10029/10041/10057/', '', '排序', 'systemStoreProductGroupSort', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10103);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10104, 10057, '/10025/10029/10041/10057/', '', '设置标签', 'systemStoreProductGroupLabels', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10104);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10105, 10058, '/10025/10029/10042/10058/', '', '列表', 'systemStoreProductGroupBuyingLst', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10105);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10106, 10058, '/10025/10029/10042/10058/', '', '详情', 'systemStoreProductGroupBuyingDetail', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10106);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10107, 10059, '/10025/10026/10037/10059/', '', '列表', 'systemSeckillActiveGetActiveList', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10107);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10108, 10059, '/10025/10026/10037/10059/', '', '详情', 'systemSeckillActiveGetActiveInfo', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10108);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10109, 10059, '/10025/10026/10037/10059/', '', '创建', 'systemSeckillActiveCreateActive', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10109);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10110, 10059, '/10025/10026/10037/10059/', '', '编辑', 'systemSeckillActiveUpdateActive', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10110);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10111, 10059, '/10025/10026/10037/10059/', '', '编辑状态', 'systemSeckillActiveUpdateActiveStatus', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10111);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10112, 10059, '/10025/10026/10037/10059/', '', '删除', 'systemSeckillActiveDeleteActive', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10112);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10113, 10059, '/10025/10026/10037/10059/', '', '活动统计数据面板', 'systemSeckillActiveChartPanel', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10113);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10114, 10059, '/10025/10026/10037/10059/', '', '活动参与人统计列表', 'systemSeckillActiveChartPeople', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10114);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10115, 10059, '/10025/10026/10037/10059/', '', '活动订单统计列表', 'systemSeckillActiveChartOrder', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10115);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10116, 10059, '/10025/10026/10037/10059/', '', '活动商品统计列表', 'systemSeckillActiveChartProduct', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10116);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10117, 10060, '/10025/10026/10036/10060/', '', '统计', 'systemStoreSeckillProductLstFilter', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10117);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10118, 10060, '/10025/10026/10036/10060/', '', '列表', 'systemStoreSeckillProductPageLst', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10118);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10119, 10060, '/10025/10026/10036/10060/', '', '列表', 'systemStoreSeckillProductLst', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10119);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10120, 10060, '/10025/10026/10036/10060/', '', '详情', 'systemStoreSeckillProductDetail', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10120);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10121, 10060, '/10025/10026/10036/10060/', '', '编辑', 'systemStoreSeckillProductUpdate', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10121);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10122, 10060, '/10025/10026/10036/10060/', '', '审核', 'systemStoreSeckillProductSwitchStatus', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10122);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10123, 10060, '/10025/10026/10036/10060/', '', '显示/隐藏', 'systemStoreSeckillProductChangeUsed', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10123);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10124, 10060, '/10025/10026/10036/10060/', '', '设置标签', 'systemStoreSeckillProductLabels', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10124);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10125, 10060, '/10025/10026/10036/10060/', '', '加入回收站', 'systemStoreSeckillProductDelete', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10125);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10126, 10060, '/10025/10026/10036/10060/', '', '删除', 'systemStoreSeckillProductDestory', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10126);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10127, 10061, '/10021/10022/10061/', '', '列表', 'systemOrderLst', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10127);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10128, 10061, '/10021/10022/10061/', '', '金额统计', 'systemOrderStat', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10128);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10129, 10061, '/10021/10022/10061/', '', '快递查询', 'systemOrderExpress', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10129);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10130, 10061, '/10021/10022/10061/', '', '头部统计', 'systemOrderTitle', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10130);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10131, 10061, '/10021/10022/10061/', '', '详情', 'systemOrderDetail', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10131);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10132, 10061, '/10021/10022/10061/', '', '导出', 'systemOrderExcel', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10132);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10133, 10061, '/10021/10022/10061/', '', '记录', 'systemOrderStatus', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10133);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10134, 10061, '/10021/10022/10061/', '', '关联订单', 'systemOrderChildrenList', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10134);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10135, 10062, '/10021/10022/10062/', '', '导出列表', 'systemStoreExcelLst', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10135);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10136, 10062, '/10021/10022/10062/', '', '导出列表', 'systemStoreExcelDownload', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10136);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10137, 10063, '/10021/10024/10063/', '', '核销', 'systemOrderTakeStat', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10137);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10138, 10063, '/10021/10024/10063/', '', '核销订单', 'systemTakeOrderLst', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10138);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10139, 10063, '/10021/10024/10063/', '', '头部统计', 'systemTakeOrderTitle', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10139);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10140, 10064, '/10021/10023/10064/', '', '列表', 'systemRefundOrderLst', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10140);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10141, 10064, '/10021/10023/10064/', '', '详情', 'systemRefundOrderDetail', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10141);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10142, 10064, '/10021/10023/10064/', '', '日志', 'systemRefundOrderLog', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10142);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10143, 10064, '/10021/10023/10064/', '', '审核', 'systemRefundOrderApprove', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10143);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10144, 10064, '/10021/10023/10064/', '', '导出', 'systemRefundOrderExcel', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10144);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10145, 10065, '/10021/10023/10065/', '', '导出列表', 'systemStoreExcelLst', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10145);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10146, 10065, '/10021/10023/10065/', '', '导出下载', 'systemStoreExcelDownload', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10146);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10147, 10066, '/10019/10020/10066/', '', '统计', 'systemStoreProductLstFilter', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10147);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10148, 10066, '/10019/10020/10066/', '', '列表', 'systemStoreProductLst', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10148);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10149, 10066, '/10019/10020/10066/', '', '详情', 'systemStoreProductDetail', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10149);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10150, 10066, '/10019/10020/10066/', '', '编辑', 'systemStoreProductUpdate', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10150);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10151, 10066, '/10019/10020/10066/', '', '上下架', 'systemStoreProductSwitchStatus', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10151);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10152, 10066, '/10019/10020/10066/', '', '批量上下架', 'systemStoreProductSwitchBatchStatus', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10152);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10153, 10066, '/10019/10020/10066/', '', '批量设置标签', 'systemStoreProductSwitchBatchLabels', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10153);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10154, 10066, '/10019/10020/10066/', '', '批量设置推荐', 'systemStoreProductSwitchBatchHot', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10154);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10155, 10066, '/10019/10020/10066/', '', '批量设置分类推荐', 'systemStoreProductSwitchBatchCateHot', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10155);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10156, 10066, '/10019/10020/10066/', '', '分销状态变更商品检测', 'systemStoreProductCheck', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10156);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10157, 10066, '/10019/10020/10066/', '', '显示/隐藏', 'systemStoreProductChangeUsed', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10157);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10158, 10066, '/10019/10020/10066/', '', '虚拟销量', 'systemStoreProductAddFicti', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10158);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10159, 10066, '/10019/10020/10066/', '', '设置标签', 'systemStoreProductLabels', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10159);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10160, 10066, '/10019/10020/10066/', '', '获取商品操作记录', 'systemStoreProductGetOperateList', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10160);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10161, 10066, '/10019/10020/10066/', '', '获取自营商品列表', 'systemStoreProductGetSelfProductList', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10161);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10162, 10066, '/10019/10020/10066/', '', '批量复制商品到店铺', 'systemStoreProductBatchCopy', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10162);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10163, 10067, '/10032/10033/10034/10067/', '', '身份列表', 'systemRoleGetList', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10163);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10164, 10067, '/10032/10033/10034/10067/', '', '身份添加', 'systemRoleCreate', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10164);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10165, 10067, '/10032/10033/10034/10067/', '', '身份编辑', 'systemRoleUpdate', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10165);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10166, 10067, '/10032/10033/10034/10067/', '', '身份修改状态', 'systemRoleStatus', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10166);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10167, 10067, '/10032/10033/10034/10067/', '', '身份删除', 'systemRoleDelete', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10167);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10168, 10068, '/10032/10033/10035/10068/', '', '管理员列表', 'systemAdminLst', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10168);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10169, 10068, '/10032/10033/10035/10068/', '', '管理员修改状态', 'systemAdminStatus', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10169);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10170, 10068, '/10032/10033/10035/10068/', '', '管理员添加', 'systemAdminCreate', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10170);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10171, 10068, '/10032/10033/10035/10068/', '', '管理员编辑', 'systemAdminUpdate', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10171);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10172, 10068, '/10032/10033/10035/10068/', '', '管理员修改密码', 'systemAdminPassword', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10172);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10173, 10068, '/10032/10033/10035/10068/', '', '管理员删除', 'systemAdminDelete', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10173);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10174, 10069, '/10032/10033/10035/10069/', '', '上传图片', 'uploadImage', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10174);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10175, 10069, '/10032/10033/10035/10069/', '', '图片列表', 'systemAttachmentLst', '', 1, 1, 1, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10175);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10176, 0, '/', '', '权限', 'self', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10176);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10177, 10176, '/10176/', '', '修改信息', 'systemAdminEdit', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10177);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10178, 10176, '/10176/', '', '修改密码', 'systemAdminEditPassword', '', 1, 1, 1, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10178);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10183, 8600, '/95/105/8600/', '', '批量加入回收站', 'merchantStoreProductBatchDelete', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10183);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10184, 8600, '/95/105/8600/', '', '批量恢复', 'merchantStoreProductBatchRestore', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10184);
INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)
SELECT 10185, 8600, '/95/105/8600/', '', '批量设置服务保障', 'merchantStoreProductBatchGuarantee', '', 1, 1, 2, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 10185);

-- 角色 rules：按导入后的全表重建（平台超管 / 商户模板角色）
SET SESSION group_concat_max_len = 1024*1024;

UPDATE `qixi_system_role` r
SET r.`rules` = (
  SELECT GROUP_CONCAT(m.`menu_id` ORDER BY m.`menu_id` SEPARATOR ',')
  FROM `qixi_system_menu` m WHERE m.`is_mer` = 1
)
WHERE r.`role_id` = 1 AND r.`mer_id` = 0;

UPDATE `qixi_system_role` r
SET r.`rules` = (
  SELECT GROUP_CONCAT(m.`menu_id` ORDER BY m.`menu_id` SEPARATOR ',')
  FROM `qixi_system_menu` m WHERE m.`is_mer` = 2
)
WHERE r.`role_id` = 2;

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'crmeb-menu-043', '全量 CRMEB 菜单/按钮导入（冲突 id +20000；rules 按 is_mer 重建）'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'crmeb-menu-043');

