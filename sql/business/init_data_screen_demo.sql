-- 平台「数据大屏」本地演示数据（可重复导入）。
-- 用法：
--   docker exec -i pte_live_mysql sh -ec 'MYSQL_PWD="$MYSQL_ROOT_PASSWORD" exec mysql --protocol=socket -uroot --default-character-set=utf8mb4' < sql/business/init_data_screen_demo.sql
-- 全部为明确标识的虚构中文记录；封面使用 admin-platform/public/demo/data-screen/thumbs。
USE `qixi_crm_business`;
SET NAMES utf8mb4;
SET time_zone = '+08:00';

INSERT INTO `qixi_crm_b_store_view` (`store_id`,`merchant_id`,`store_app_id`,`store_name`,`status`,`integral_enabled`,`integral_points_per_yuan`,`integral_max_deduction_bps`) VALUES
  (11,11,'qixi.screen.demo.11','云岚风物旗舰店',1,0,100,2000),
  (12,12,'qixi.screen.demo.12','星野旅行生活馆',1,0,100,2000),
  (13,13,'qixi.screen.demo.13','山海影像精选店',1,0,100,2000),
  (14,14,'qixi.screen.demo.14','林泉露营装备店',1,0,100,2000),
  (15,15,'qixi.screen.demo.15','晴川艺术画廊',1,0,100,2000),
  (16,16,'qixi.screen.demo.16','远方自然美学店',1,0,100,2000),
  (17,17,'qixi.screen.demo.17','云端摄影器材店',1,0,100,2000),
  (18,18,'qixi.screen.demo.18','晨雾家居优选店',1,0,100,2000),
  (19,19,'qixi.screen.demo.19','青禾鲜果旗舰店',1,0,100,2000),
  (20,20,'qixi.screen.demo.20','暖阳母婴生活馆',1,0,100,2000),
  (21,21,'qixi.screen.demo.21','澄波茶饮礼盒店',1,0,100,2000),
  (22,22,'qixi.screen.demo.22','逐风运动优选店',1,0,100,2000)
ON DUPLICATE KEY UPDATE `merchant_id`=VALUES(`merchant_id`),`store_app_id`=VALUES(`store_app_id`),`store_name`=VALUES(`store_name`),`status`=VALUES(`status`),`updated_at`=NOW();

INSERT INTO `qixi_crm_b_product_view` (`product_id`,`merchant_id`,`store_id`,`merchant_name`,`store_name`,`category_id`,`title`,`cover_url`,`price`,`original_price`,`product_type`,`sales`,`stock`,`sale_status`,`version`,`updated_at`) VALUES
  (1501,11,11,'云岚风物商户','云岚风物旗舰店',10601,'云海雪山摄影艺术画','/demo/data-screen/thumbs/thumb-01.jpg',698.00,898.00,0,428,80,1,1,NOW()),
  (1502,12,12,'星野旅行商户','星野旅行生活馆',10601,'湖畔晨光旅行画册','/demo/data-screen/thumbs/thumb-02.jpg',488.00,628.00,0,396,96,1,1,NOW()),
  (1503,13,13,'山海影像商户','山海影像精选店',10601,'旷野日落高清装饰画','/demo/data-screen/thumbs/thumb-03.jpg',788.00,998.00,0,371,72,1,1,NOW()),
  (1504,14,14,'林泉露营商户','林泉露营装备店',10601,'森林秘境露营礼盒','/demo/data-screen/thumbs/thumb-04.jpg',968.00,1288.00,0,318,58,1,1,NOW()),
  (1505,15,15,'晴川艺术商户','晴川艺术画廊',10601,'霞光山脉收藏级画芯','/demo/data-screen/thumbs/thumb-05.jpg',1288.00,1588.00,0,286,48,1,1,NOW()),
  (1506,16,16,'远方自然商户','远方自然美学店',10601,'碧海晴空艺术微喷','/demo/data-screen/thumbs/thumb-06.jpg',1088.00,1388.00,0,265,45,1,1,NOW()),
  (1507,17,17,'云端摄影商户','云端摄影器材店',10601,'山谷云层风景摄影集','/demo/data-screen/thumbs/thumb-07.jpg',898.00,1188.00,0,243,61,1,1,NOW()),
  (1508,18,18,'晨雾家居商户','晨雾家居优选店',10601,'自然风光沉浸式画框','/demo/data-screen/thumbs/thumb-08.jpg',568.00,728.00,0,221,88,1,1,NOW()),
  (1509,19,19,'青禾鲜果商户','青禾鲜果旗舰店',105,'演示鲜果礼盒','/demo/data-screen/thumbs/thumb-09.jpg',168.00,228.00,0,512,120,1,1,NOW()),
  (1510,20,20,'暖阳母婴商户','暖阳母婴生活馆',101,'演示母婴安心套装','/demo/data-screen/thumbs/thumb-10.jpg',298.00,398.00,0,468,90,1,1,NOW()),
  (1511,21,21,'澄波茶饮商户','澄波茶饮礼盒店',105,'演示茶饮礼盒','/demo/data-screen/thumbs/thumb-11.jpg',218.00,288.00,0,401,110,1,1,NOW()),
  (1512,22,22,'逐风运动商户','逐风运动优选店',10601,'演示运动装备组合','/demo/data-screen/thumbs/thumb-12.jpg',458.00,598.00,0,355,70,1,1,NOW())
ON DUPLICATE KEY UPDATE `merchant_id`=VALUES(`merchant_id`),`store_id`=VALUES(`store_id`),`merchant_name`=VALUES(`merchant_name`),`store_name`=VALUES(`store_name`),`title`=VALUES(`title`),`cover_url`=VALUES(`cover_url`),`price`=VALUES(`price`),`original_price`=VALUES(`original_price`),`sales`=VALUES(`sales`),`stock`=VALUES(`stock`),`sale_status`=VALUES(`sale_status`),`updated_at`=NOW();

INSERT INTO `qixi_crm_b_user` (`id`,`nickname`,`mobile`,`status`,`group_id`,`auth_version`,`created_at`) WITH RECURSIVE `screen_seq` AS (
  SELECT 1 AS `n` UNION ALL SELECT `n` + 1 FROM `screen_seq` WHERE `n` < 48
)
SELECT 97200 + `n`, CONCAT('大屏演示用户', LPAD(`n`, 2, '0')), CONCAT('SCREEN-DEMO-', LPAD(`n`, 4, '0')), 1, 0, 1,
  CASE WHEN `n` <= 16 THEN NOW() ELSE DATE_SUB(NOW(), INTERVAL 20 DAY) END
FROM `screen_seq`
ON DUPLICATE KEY UPDATE `nickname`=VALUES(`nickname`),`mobile`=VALUES(`mobile`),`status`=VALUES(`status`),`created_at`=VALUES(`created_at`),`updated_at`=NOW();

-- 360 笔演示订单：前 36 笔落在「今天」供实时订单/今日排行/小时图；其余分布在本月供地图与月销售。
INSERT INTO `qixi_crm_b_group_order` (`id`,`order_no`,`user_id`,`total_amount`,`discount_amount`,`freight_amount`,`pay_amount`,`total_quantity`,`activity_type`,`points_amount`,`recipient_snapshot`,`pay_channel`,`pay_status`,`paid_at`,`idempotency_key`,`remark`) WITH RECURSIVE `screen_seq` AS (
  SELECT 1 AS `n` UNION ALL SELECT `n` + 1 FROM `screen_seq` WHERE `n` < 360
)
SELECT
  975000 + `n`,
  CONCAT('SCREEN-DEMO-G-', DATE_FORMAT(CURDATE(), '%Y%m%d'), '-', LPAD(`n`, 3, '0')),
  97200 + MOD(`n` - 1, 48) + 1,
  168.00 + MOD(`n`, 12) * 86.50,
  0.00, 0.00,
  168.00 + MOD(`n`, 12) * 86.50,
  1 + MOD(`n`, 4),
  0, 0,
  JSON_OBJECT(
    'recipient', '大屏演示收件人',
    'province', ELT(
      CASE
        WHEN MOD(`n`, 100) < 32 THEN 1  -- 广东省
        WHEN MOD(`n`, 100) < 52 THEN 2  -- 浙江省
        WHEN MOD(`n`, 100) < 68 THEN 3  -- 江苏省
        WHEN MOD(`n`, 100) < 78 THEN 4  -- 山东省
        WHEN MOD(`n`, 100) < 86 THEN 5  -- 四川省
        WHEN MOD(`n`, 100) < 91 THEN 6  -- 河南省
        WHEN MOD(`n`, 100) < 94 THEN 7  -- 湖北省
        WHEN MOD(`n`, 100) < 96 THEN 8  -- 福建省
        WHEN MOD(`n`, 100) < 97 THEN 9  -- 湖南省
        WHEN MOD(`n`, 100) < 98 THEN 10 -- 安徽省
        WHEN MOD(`n`, 100) < 99 THEN 11 -- 北京市
        ELSE 12                          -- 上海市
      END,
      '广东省','浙江省','江苏省','山东省','四川省','河南省','湖北省','福建省','湖南省','安徽省','北京市','上海市'
    )
  ),
  ELT(MOD(`n` - 1, 3) + 1, 'wechat', 'alipay', 'balance'),
  'paid',
  CASE
    -- 最近 24 小时整点序列（供「订单支付情况」按小时统计）
    WHEN `n` <= 48 THEN DATE_FORMAT(DATE_SUB(NOW(), INTERVAL MOD(`n` - 1, 24) HOUR), '%Y-%m-%d %H:00:00')
         + INTERVAL MOD(`n` * 7, 50) MINUTE + INTERVAL MOD(`n` * 11, 50) SECOND
    ELSE DATE_FORMAT(CURDATE(), '%Y-%m-01') + INTERVAL LEAST(DAY(CURDATE()) - 1, MOD(`n`, GREATEST(DAY(CURDATE()), 1))) DAY
         + INTERVAL (8 + MOD(`n`, 12)) HOUR + INTERVAL MOD(`n` * 3, 50) MINUTE
  END,
  CONCAT('screen-demo-group-', LPAD(`n`, 3, '0')),
  '数据大屏本地演示订单'
FROM `screen_seq`
ON DUPLICATE KEY UPDATE
  `user_id`=VALUES(`user_id`),`total_amount`=VALUES(`total_amount`),`pay_amount`=VALUES(`pay_amount`),
  `total_quantity`=VALUES(`total_quantity`),`recipient_snapshot`=VALUES(`recipient_snapshot`),
  `pay_channel`=VALUES(`pay_channel`),`pay_status`='paid',`paid_at`=VALUES(`paid_at`),
  `remark`=VALUES(`remark`),`updated_at`=NOW();

INSERT INTO `qixi_crm_b_order` (`id`,`group_order_id`,`order_no`,`merchant_id`,`merchant_name_snapshot`,`store_id`,`store_name_snapshot`,`user_id`,`total_amount`,`discount_amount`,`freight_amount`,`pay_amount`,`total_quantity`,`activity_type`,`points_amount`,`recipient_snapshot`,`remark`,`status`,`paid_at`) WITH RECURSIVE `screen_seq` AS (
  SELECT 1 AS `n` UNION ALL SELECT `n` + 1 FROM `screen_seq` WHERE `n` < 360
)
SELECT
  975000 + `n`,
  975000 + `n`,
  CONCAT('SCREEN-DEMO-O-', DATE_FORMAT(CURDATE(), '%Y%m%d'), '-', LPAD(`n`, 3, '0')),
  10 + MOD(`n` - 1, 12) + 1,
  ELT(MOD(`n` - 1, 12) + 1,'云岚风物商户','星野旅行商户','山海影像商户','林泉露营商户','晴川艺术商户','远方自然商户','云端摄影商户','晨雾家居商户','青禾鲜果商户','暖阳母婴商户','澄波茶饮商户','逐风运动商户'),
  10 + MOD(`n` - 1, 12) + 1,
  ELT(MOD(`n` - 1, 12) + 1,'云岚风物旗舰店','星野旅行生活馆','山海影像精选店','林泉露营装备店','晴川艺术画廊','远方自然美学店','云端摄影器材店','晨雾家居优选店','青禾鲜果旗舰店','暖阳母婴生活馆','澄波茶饮礼盒店','逐风运动优选店'),
  97200 + MOD(`n` - 1, 48) + 1,
  168.00 + MOD(`n`, 12) * 86.50,
  0.00, 0.00,
  168.00 + MOD(`n`, 12) * 86.50,
  1 + MOD(`n`, 4),
  0, 0,
  JSON_OBJECT(
    'recipient', '大屏演示收件人',
    'province', ELT(
      CASE
        WHEN MOD(`n`, 100) < 32 THEN 1
        WHEN MOD(`n`, 100) < 52 THEN 2
        WHEN MOD(`n`, 100) < 68 THEN 3
        WHEN MOD(`n`, 100) < 78 THEN 4
        WHEN MOD(`n`, 100) < 86 THEN 5
        WHEN MOD(`n`, 100) < 91 THEN 6
        WHEN MOD(`n`, 100) < 94 THEN 7
        WHEN MOD(`n`, 100) < 96 THEN 8
        WHEN MOD(`n`, 100) < 97 THEN 9
        WHEN MOD(`n`, 100) < 98 THEN 10
        WHEN MOD(`n`, 100) < 99 THEN 11
        ELSE 12
      END,
      '广东省','浙江省','江苏省','山东省','四川省','河南省','湖北省','福建省','湖南省','安徽省','北京市','上海市'
    )
  ),
  '数据大屏本地演示订单',
  ELT(MOD(`n` - 1, 4) + 1,'paid','fulfilling','shipped','completed'),
  CASE
    -- 最近 24 小时整点序列（供「订单支付情况」按小时统计）
    WHEN `n` <= 48 THEN DATE_FORMAT(DATE_SUB(NOW(), INTERVAL MOD(`n` - 1, 24) HOUR), '%Y-%m-%d %H:00:00')
         + INTERVAL MOD(`n` * 7, 50) MINUTE + INTERVAL MOD(`n` * 11, 50) SECOND
    ELSE DATE_FORMAT(CURDATE(), '%Y-%m-01') + INTERVAL LEAST(DAY(CURDATE()) - 1, MOD(`n`, GREATEST(DAY(CURDATE()), 1))) DAY
         + INTERVAL (8 + MOD(`n`, 12)) HOUR + INTERVAL MOD(`n` * 3, 50) MINUTE
  END
FROM `screen_seq`
ON DUPLICATE KEY UPDATE
  `merchant_id`=VALUES(`merchant_id`),`merchant_name_snapshot`=VALUES(`merchant_name_snapshot`),
  `store_id`=VALUES(`store_id`),`store_name_snapshot`=VALUES(`store_name_snapshot`),
  `user_id`=VALUES(`user_id`),`total_amount`=VALUES(`total_amount`),`pay_amount`=VALUES(`pay_amount`),
  `total_quantity`=VALUES(`total_quantity`),`recipient_snapshot`=VALUES(`recipient_snapshot`),
  `status`=VALUES(`status`),`paid_at`=VALUES(`paid_at`),`updated_at`=NOW();

INSERT INTO `qixi_crm_b_order_item` (`id`,`order_id`,`product_id`,`merchant_sku_id`,`sku_key`,`title_snapshot`,`cover_url_snapshot`,`spec_snapshot`,`unit_price`,`quantity`,`refund_quantity`) WITH RECURSIVE `screen_seq` AS (
  SELECT 1 AS `n` UNION ALL SELECT `n` + 1 FROM `screen_seq` WHERE `n` < 360
)
SELECT
  976000 + `n`,
  975000 + `n`,
  1500 + MOD(`n` - 1, 12) + 1,
  0,
  CONCAT('screen-', MOD(`n` - 1, 12) + 1),
  ELT(MOD(`n` - 1, 12) + 1,'云海雪山摄影艺术画','湖畔晨光旅行画册','旷野日落高清装饰画','森林秘境露营礼盒','霞光山脉收藏级画芯','碧海晴空艺术微喷','山谷云层风景摄影集','自然风光沉浸式画框','演示鲜果礼盒','演示母婴安心套装','演示茶饮礼盒','演示运动装备组合'),
  ELT(MOD(`n` - 1, 12) + 1,
    '/demo/data-screen/thumbs/thumb-01.jpg',
    '/demo/data-screen/thumbs/thumb-02.jpg',
    '/demo/data-screen/thumbs/thumb-03.jpg',
    '/demo/data-screen/thumbs/thumb-04.jpg',
    '/demo/data-screen/thumbs/thumb-05.jpg',
    '/demo/data-screen/thumbs/thumb-06.jpg',
    '/demo/data-screen/thumbs/thumb-07.jpg',
    '/demo/data-screen/thumbs/thumb-08.jpg',
    '/demo/data-screen/thumbs/thumb-09.jpg',
    '/demo/data-screen/thumbs/thumb-10.jpg',
    '/demo/data-screen/thumbs/thumb-11.jpg',
    '/demo/data-screen/thumbs/thumb-12.jpg'
  ),
  JSON_OBJECT('规格','大屏演示款'),
  168.00 + MOD(`n`, 12) * 86.50,
  1 + MOD(`n`, 4),
  0
FROM `screen_seq`
ON DUPLICATE KEY UPDATE
  `product_id`=VALUES(`product_id`),`sku_key`=VALUES(`sku_key`),`title_snapshot`=VALUES(`title_snapshot`),
  `cover_url_snapshot`=VALUES(`cover_url_snapshot`),`spec_snapshot`=VALUES(`spec_snapshot`),
  `unit_price`=VALUES(`unit_price`),`quantity`=VALUES(`quantity`),`refund_quantity`=0;

-- 浏览量 / 访客：约 240 条今日浏览，覆盖多用户。
INSERT INTO `qixi_crm_b_user_browse_history` (`id`,`user_id`,`product_id`,`store_id`,`viewed_at`) WITH RECURSIVE `screen_seq` AS (
  SELECT 1 AS `n` UNION ALL SELECT `n` + 1 FROM `screen_seq` WHERE `n` < 240
)
SELECT
  977000 + `n`,
  97200 + MOD(`n` - 1, 48) + 1,
  1500 + MOD(`n` - 1, 12) + 1,
  10 + MOD(`n` - 1, 12) + 1,
  TIMESTAMP(CURDATE()) + INTERVAL MOD(`n`, 23) HOUR + INTERVAL MOD(`n` * 5, 55) MINUTE
FROM `screen_seq`
ON DUPLICATE KEY UPDATE `user_id`=VALUES(`user_id`),`product_id`=VALUES(`product_id`),`store_id`=VALUES(`store_id`),`viewed_at`=VALUES(`viewed_at`);
