-- 直播间管理：主播微信号字段 + 平台列表中文演示数据（对齐 CRMEB 直播间管理）
USE `qixi_crm_business`;
SET NAMES utf8mb4;

SET @col_wechat := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'qixi_crm_b_broadcast_room'
    AND COLUMN_NAME = 'anchor_wechat'
);
SET @sql_wechat := IF(
  @col_wechat = 0,
  'ALTER TABLE `qixi_crm_b_broadcast_room`
     ADD COLUMN `anchor_wechat` varchar(64) NOT NULL DEFAULT '''' COMMENT ''主播微信号（演示/展示用，非真实账号）'' AFTER `anchor_name`',
  'SELECT 1'
);
PREPARE stmt_wechat FROM @sql_wechat;
EXECUTE stmt_wechat;
DEALLOCATE PREPARE stmt_wechat;

-- 刷新既有夹具微信号与可读中文字段（推流地址、主播手机号保持空）
UPDATE `qixi_crm_b_broadcast_room`
SET `anchor_wechat`='demo_anchor_xq',
    `name`='CRM Live服饰秋日穿搭直播间',
    `anchor_name`='虚构主播小七',
    `mark`='中文模拟直播审核夹具',
    `sort`=10,
    `star`=5
WHERE `broadcast_room_id`=7101;

UPDATE `qixi_crm_b_broadcast_room`
SET `anchor_wechat`='demo_anchor_xj',
    `name`='CRM Live居家香氛新品直播间',
    `anchor_name`='虚构主播小居',
    `mark`='中文模拟已审核直播夹具',
    `sort`=8,
    `star`=4
WHERE `broadcast_room_id`=7102;

-- 补充平台「直播间管理」列表演示行（utf8mb4 中文可读）
INSERT INTO `qixi_crm_b_broadcast_room` (
  `broadcast_room_id`,`mer_id`,`name`,`cover_img`,`feeds_img`,`play_url`,`push_url`,
  `start_time`,`end_time`,`anchor_name`,`anchor_wechat`,`phone`,
  `status`,`live_status`,`is_show`,`is_del`,`create_time`,`sort`,`star`,`mark`,`refusal`
) VALUES
  (
    7103, 1, 'CRM Live服饰周末上新直播间',
    '/demo/live-fashion-cover.png', '/demo/live-fashion-feed.png', '', '',
    DATE_ADD(NOW(), INTERVAL 2 DAY), DATE_ADD(DATE_ADD(NOW(), INTERVAL 2 DAY), INTERVAL 3 HOUR),
    '虚构主播小茉', 'demo_anchor_xm', '',
    2, 102, 1, 0, DATE_SUB(NOW(), INTERVAL 3 DAY), 12, 3, '中文模拟待开播直播间', ''
  ),
  (
    7104, 2, 'CRM Live居家收纳技巧直播间',
    '/demo/live-home-cover.png', '/demo/live-home-feed.png', '', '',
    DATE_SUB(NOW(), INTERVAL 2 DAY), DATE_ADD(DATE_SUB(NOW(), INTERVAL 2 DAY), INTERVAL 2 HOUR),
    '虚构主播小宁', 'demo_anchor_xn', '',
    2, 103, 0, 0, DATE_SUB(NOW(), INTERVAL 5 DAY), 6, 2, '中文模拟已结束直播间', ''
  ),
  (
    7105, 3, 'CRM Live数码开学季直播间',
    '/demo/live-digital-cover.png', '/demo/live-digital-feed.png', '', '',
    DATE_ADD(NOW(), INTERVAL 12 HOUR), DATE_ADD(NOW(), INTERVAL 15 HOUR),
    '虚构主播小数', 'demo_anchor_xs', '',
    -1, 102, 0, 0, DATE_SUB(NOW(), INTERVAL 1 DAY), 4, 0, '中文模拟驳回直播间', '封面不符合规范，请更换后重新提交'
  ),
  (
    7106, 1, 'CRM Live服饰晚间连麦直播间',
    '/demo/live-fashion-cover.png', '/demo/live-fashion-feed.png', '', '',
    DATE_SUB(NOW(), INTERVAL 30 MINUTE), DATE_ADD(NOW(), INTERVAL 90 MINUTE),
    '虚构主播小晚', 'demo_anchor_xw', '',
    2, 101, 1, 0, DATE_SUB(NOW(), INTERVAL 2 DAY), 15, 5, '中文模拟直播中直播间', ''
  )
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `cover_img`=VALUES(`cover_img`),
  `feeds_img`=VALUES(`feeds_img`),
  `start_time`=VALUES(`start_time`),
  `end_time`=VALUES(`end_time`),
  `anchor_name`=VALUES(`anchor_name`),
  `anchor_wechat`=VALUES(`anchor_wechat`),
  `status`=VALUES(`status`),
  `live_status`=VALUES(`live_status`),
  `is_show`=VALUES(`is_show`),
  `sort`=VALUES(`sort`),
  `star`=VALUES(`star`),
  `mark`=VALUES(`mark`),
  `refusal`=VALUES(`refusal`);

INSERT INTO `qixi_crm_b_broadcast_room_goods` (`broadcast_room_id`,`product_id`,`on_sale`,`sort`) VALUES
  (7103,1001,1,1),(7103,1003,1,2),
  (7104,1101,1,1),
  (7105,1201,1,1),
  (7106,1002,1,1),(7106,1001,1,2)
ON DUPLICATE KEY UPDATE `on_sale`=VALUES(`on_sale`),`sort`=VALUES(`sort`);
