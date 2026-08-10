-- 预售管理：平台列表扩展字段 + 本地演示数据（全款/定金 Tab）
USE `qixi_crm_business`;
SET NAMES utf8mb4;

-- star / sys_labels / stock_count / attend 统计字段（幂等）
SET @col_star := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'qixi_crm_b_presell' AND COLUMN_NAME = 'star'
);
SET @sql_star := IF(
  @col_star = 0,
  'ALTER TABLE `qixi_crm_b_presell`
     ADD COLUMN `star` tinyint NOT NULL DEFAULT 0 COMMENT ''推荐星级'' AFTER `seles`,
     ADD COLUMN `sys_labels` varchar(255) NOT NULL DEFAULT '''' COMMENT ''平台标签ID逗号分隔'' AFTER `star`,
     ADD COLUMN `stock_count` int NOT NULL DEFAULT 0 COMMENT ''限量总数'' AFTER `sys_labels`,
     ADD COLUMN `attend_num` int NOT NULL DEFAULT 0 COMMENT ''参与人数'' AFTER `stock_count`,
     ADD COLUMN `success_num` int NOT NULL DEFAULT 0 COMMENT ''成功人数'' AFTER `attend_num`',
  'SELECT 1'
);
PREPARE stmt_star FROM @sql_star;
EXECUTE stmt_star;
DEALLOCATE PREPARE stmt_star;

-- 存量行补限量：stock_count = stock + seles（仅未写过时）
UPDATE `qixi_crm_b_presell`
SET `stock_count` = GREATEST(`stock` + `seles`, `stock`, 0)
WHERE `stock_count` = 0 AND (`stock` > 0 OR `seles` > 0);

-- 演示数据：全款约 10 条、定金约 4 条；覆盖自营/非自营、未开始/进行中/已结束、待审/通过/拒绝/下架
INSERT INTO `qixi_crm_b_presell` (
  `product_presell_id`,`start_time`,`end_time`,`final_start_time`,`final_end_time`,
  `status`,`presell_type`,`pay_count`,`delivery_type`,`delivery_day`,
  `product_id`,`price`,`down_price`,`final_price`,`stock`,`is_show`,
  `store_name`,`mer_id`,`store_info`,`is_del`,`product_status`,`refusal`,
  `action_status`,`seles`,`star`,`sys_labels`,`stock_count`,`attend_num`,`success_num`
) VALUES
  -- ===== 全款预售（10）=====
  -- 全款·进行中·出售中·非自营
  (6202, DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_ADD(NOW(), INTERVAL 60 DAY), '', '',
   1, 1, 0, 1, 10, 1104, 238.00, 0.00, 0.00, 39, 1,
   '晨间居家香氛套装', 2, '全款预售，预计十日内发货。', 0, 1, '',
   1, 37, 3, '', 76, 52, 37),
  -- 全款·未开始·待审核·自营
  (6211, DATE_ADD(NOW(), INTERVAL 3 DAY), DATE_ADD(NOW(), INTERVAL 30 DAY), '', '',
   1, 1, 2, 1, 7, 1004, 188.00, 0.00, 0.00, 50, 0,
   '精纺圆领羊毛开衫·全款预售', 1, '全款预售待审演示。', 0, 0, '',
   1, 0, 0, '', 50, 0, 0),
  -- 全款·已结束·非自营
  (6212, DATE_SUB(NOW(), INTERVAL 40 DAY), DATE_SUB(NOW(), INTERVAL 5 DAY), '', '',
   1, 1, 0, 1, 5, 1104, 168.00, 0.00, 0.00, 0, 1,
   '晨间居家香氛套装·往期全款', 2, '已结束全款预售演示。', 0, 1, '',
   -1, 80, 1, '', 80, 96, 80),
  -- 全款·进行中·出售中·自营
  (6215, DATE_SUB(NOW(), INTERVAL 2 DAY), DATE_ADD(NOW(), INTERVAL 45 DAY), '', '',
   1, 1, 0, 1, 8, 1001, 259.00, 0.00, 0.00, 62, 1,
   '轻奢羊绒针织衫·全款预售', 1, '自营全款预售进行中演示。', 0, 1, '',
   1, 28, 5, '', 90, 41, 28),
  -- 全款·进行中·仓库中·非自营
  (6216, DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_ADD(NOW(), INTERVAL 50 DAY), '', '',
   1, 1, 0, 1, 12, 1101, 198.00, 0.00, 0.00, 44, 0,
   '无火藤条香氛礼盒·全款仓库', 2, '审核通过但关闭显示演示。', 0, 1, '',
   1, 16, 2, '', 60, 22, 16),
  -- 全款·未开始·审核通过·自营
  (6217, DATE_ADD(NOW(), INTERVAL 5 DAY), DATE_ADD(NOW(), INTERVAL 40 DAY), '', '',
   1, 1, 1, 1, 10, 1002, 399.00, 0.00, 0.00, 80, 1,
   '头层牛皮通勤托特包·全款未开', 1, '未开始全款预售演示。', 0, 1, '',
   1, 0, 4, '', 80, 3, 0),
  -- 全款·进行中·出售中·非自营
  (6218, DATE_SUB(NOW(), INTERVAL 3 DAY), DATE_ADD(NOW(), INTERVAL 55 DAY), '', '',
   1, 1, 0, 1, 9, 1102, 119.00, 0.00, 0.00, 71, 1,
   '晚安助眠香薰蜡烛·全款预售', 2, '居家香氛全款预售演示。', 0, 1, '',
   1, 19, 1, '', 90, 27, 19),
  -- 全款·已结束·自营
  (6219, DATE_SUB(NOW(), INTERVAL 50 DAY), DATE_SUB(NOW(), INTERVAL 10 DAY), '', '',
   1, 1, 0, 1, 6, 1007, 149.00, 0.00, 0.00, 0, 1,
   '柔软亲肤针织披肩·往期全款', 1, '已结束自营全款预售演示。', 0, 1, '',
   -1, 65, 0, '', 65, 78, 65),
  -- 全款·进行中·审核未通过·非自营
  (6220, DATE_SUB(NOW(), INTERVAL 12 HOUR), DATE_ADD(NOW(), INTERVAL 25 DAY), '', '',
   1, 1, 0, 1, 7, 1103, 139.00, 0.00, 0.00, 30, 0,
   '恒温随行保温杯·全款驳回', 2, '审核未通过演示。', 0, -1, '价格与主图不符',
   1, 0, 0, '', 30, 0, 0),
  -- 全款·进行中·出售中·自营
  (6221, DATE_SUB(NOW(), INTERVAL 4 DAY), DATE_ADD(NOW(), INTERVAL 70 DAY), '', '',
   1, 1, 0, 1, 15, 1008, 369.00, 0.00, 0.00, 53, 1,
   '城市通勤训练跑鞋·全款预售', 1, '自营跑鞋全款预售演示。', 0, 1, '',
   1, 47, 5, '', 100, 68, 47),

  -- ===== 定金预售（4）=====
  -- 定金·进行中·出售中·自营
  (6201, DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_ADD(NOW(), INTERVAL 60 DAY),
   DATE_FORMAT(DATE_ADD(NOW(), INTERVAL 61 DAY), '%Y-%m-%d %H:%i:%s'),
   DATE_FORMAT(DATE_ADD(NOW(), INTERVAL 75 DAY), '%Y-%m-%d %H:%i:%s'),
   1, 2, 0, 1, 15, 1004, 299.00, 59.00, 240.00, 36, 1,
   '精纺圆领羊毛开衫', 1, '秋冬预售，尾款支付后按订单地址发货。', 0, 1, '',
   1, 48, 4, '', 84, 61, 48),
  -- 定金·进行中·仓库中·非自营
  (6203, DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_ADD(NOW(), INTERVAL 60 DAY),
   DATE_FORMAT(DATE_ADD(NOW(), INTERVAL 61 DAY), '%Y-%m-%d %H:%i:%s'),
   DATE_FORMAT(DATE_ADD(NOW(), INTERVAL 75 DAY), '%Y-%m-%d %H:%i:%s'),
   1, 2, 0, 1, 20, 1107, 279.00, 59.00, 220.00, 35, 0,
   '客厅氛围香薰礼盒·定金仓库', 2, '新品预售，支持定金与尾款支付。', 0, 1, '',
   1, 29, 2, '', 64, 40, 29),
  -- 定金·待审核·非自营
  (6213, DATE_SUB(NOW(), INTERVAL 2 HOUR), DATE_ADD(NOW(), INTERVAL 20 DAY),
   DATE_FORMAT(DATE_ADD(NOW(), INTERVAL 21 DAY), '%Y-%m-%d %H:%i:%s'),
   DATE_FORMAT(DATE_ADD(NOW(), INTERVAL 35 DAY), '%Y-%m-%d %H:%i:%s'),
   1, 2, 1, 2, 10, 1108, 79.00, 19.00, 60.00, 20, 1,
   '织物护理香氛喷雾·定金待审', 2, '定金预售待审演示。', 0, 0, '',
   1, 0, 0, '', 20, 0, 0),
  -- 定金·强制下架·自营
  (6214, DATE_SUB(NOW(), INTERVAL 5 DAY), DATE_ADD(NOW(), INTERVAL 15 DAY),
   DATE_FORMAT(DATE_ADD(NOW(), INTERVAL 16 DAY), '%Y-%m-%d %H:%i:%s'),
   DATE_FORMAT(DATE_ADD(NOW(), INTERVAL 28 DAY), '%Y-%m-%d %H:%i:%s'),
   0, 2, 0, 1, 12, 1004, 259.00, 49.00, 210.00, 10, 0,
   '精纺圆领羊毛开衫·平台下架', 1, '强制下架演示。', 0, -2, '违规宣传，平台强制下架',
   1, 5, 0, '', 15, 8, 5)
ON DUPLICATE KEY UPDATE
  `start_time`=VALUES(`start_time`),
  `end_time`=VALUES(`end_time`),
  `final_start_time`=VALUES(`final_start_time`),
  `final_end_time`=VALUES(`final_end_time`),
  `status`=VALUES(`status`),
  `presell_type`=VALUES(`presell_type`),
  `product_id`=VALUES(`product_id`),
  `mer_id`=VALUES(`mer_id`),
  `price`=VALUES(`price`),
  `down_price`=VALUES(`down_price`),
  `final_price`=VALUES(`final_price`),
  `stock`=VALUES(`stock`),
  `is_show`=VALUES(`is_show`),
  `store_name`=VALUES(`store_name`),
  `store_info`=VALUES(`store_info`),
  `is_del`=VALUES(`is_del`),
  `product_status`=VALUES(`product_status`),
  `refusal`=VALUES(`refusal`),
  `action_status`=VALUES(`action_status`),
  `seles`=VALUES(`seles`),
  `star`=VALUES(`star`),
  `sys_labels`=VALUES(`sys_labels`),
  `stock_count`=VALUES(`stock_count`),
  `attend_num`=VALUES(`attend_num`),
  `success_num`=VALUES(`success_num`);
