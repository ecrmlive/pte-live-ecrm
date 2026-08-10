-- 秒杀管理：场次活动扩展字段（对齐 CRMEB 秒杀商品列表）
USE `qixi_crm_business`;
SET NAMES utf8mb4;

SET @col_is_show := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'qixi_crm_b_seckill_active' AND COLUMN_NAME = 'is_show'
);
SET @sql_is_show := IF(
  @col_is_show = 0,
  'ALTER TABLE `qixi_crm_b_seckill_active`
     ADD COLUMN `is_show` tinyint NOT NULL DEFAULT 1 COMMENT ''是否显示'' AFTER `status`,
     ADD COLUMN `product_status` tinyint NOT NULL DEFAULT 1 COMMENT ''1通过 0待审 -1未通过 -2强制下架'' AFTER `is_show`,
     ADD COLUMN `star` tinyint NOT NULL DEFAULT 0 COMMENT ''推荐星级'' AFTER `product_status`,
     ADD COLUMN `sort` int NOT NULL DEFAULT 0 COMMENT ''排序'' AFTER `star`,
     ADD COLUMN `stock` int NOT NULL DEFAULT 0 COMMENT ''限量剩余'' AFTER `sort`,
     ADD COLUMN `sales` int NOT NULL DEFAULT 0 COMMENT ''已售数量'' AFTER `stock`,
     ADD COLUMN `sys_labels` varchar(255) NOT NULL DEFAULT '''' COMMENT ''平台标签ID逗号分隔'' AFTER `sales`,
     ADD COLUMN `refusal` varchar(500) NOT NULL DEFAULT '''' COMMENT ''拒绝/下架原因'' AFTER `sys_labels`',
  'SELECT 1'
);
PREPARE stmt_is_show FROM @sql_is_show;
EXECUTE stmt_is_show;
DEALLOCATE PREPARE stmt_is_show;

-- 演示：出售中 / 仓库中 / 待审核 / 回收站
UPDATE `qixi_crm_b_seckill_active`
SET `name`='全天', `seckill_time_ids`='1,2,3,4', `is_show`=1, `product_status`=1, `star`=3, `sort`=10,
    `stock`=88, `sales`=12, `sys_labels`='', `delete_time`=NULL, `status`=1, `active_status`=1
WHERE `seckill_active_id`=6001;

UPDATE `qixi_crm_b_seckill_active`
SET `name`='晚上', `seckill_time_ids`='4', `is_show`=0, `product_status`=1, `star`=1, `sort`=8,
    `stock`=40, `sales`=3, `sys_labels`='', `delete_time`=NULL, `status`=1, `active_status`=1
WHERE `seckill_active_id`=6002;

UPDATE `qixi_crm_b_seckill_active`
SET `name`='下午', `seckill_time_ids`='3', `is_show`=1, `product_status`=0, `star`=0, `sort`=5,
    `stock`=20, `sales`=0, `sys_labels`='', `delete_time`=NULL, `status`=1, `active_status`=1
WHERE `seckill_active_id`=6003;

UPDATE `qixi_crm_b_seckill_active`
SET `name`='早上', `seckill_time_ids`='2', `is_show`=0, `product_status`=1, `star`=0, `sort`=1,
    `stock`=0, `sales`=6, `sys_labels`='', `delete_time`=UNIX_TIMESTAMP(), `status`=0, `active_status`=0,
    `refusal`='活动结束清理'
WHERE `seckill_active_id`=6004;

-- product_id 必须同时存在于：
--   商户库 qixi_crm_m_product（详情/编辑 /products/:id/edit）
--   业务库 qixi_crm_b_product_view（列表 enrich）
-- local 演示商品以 1001/1002 为准（两端皆有）
UPDATE `qixi_crm_b_seckill_active`
SET `product_id`=1001, `mer_id`=1
WHERE `seckill_active_id` IN (6001,6003,6005,6007);
UPDATE `qixi_crm_b_seckill_active`
SET `product_id`=1002, `mer_id`=1
WHERE `seckill_active_id` IN (6002,6004,6006);

INSERT INTO `qixi_crm_b_seckill_active`
  (`seckill_active_id`,`name`,`seckill_time_ids`,`start_day`,`end_day`,`mer_id`,`product_id`,`seckill_price`,
   `once_pay_count`,`all_pay_count`,`active_status`,`status`,`is_show`,`product_status`,`star`,`sort`,`stock`,`sales`,
   `sys_labels`,`refusal`,`create_time`,`update_time`,`delete_time`)
VALUES
  (6005,'午夜','1',DATE_SUB(CURDATE(),INTERVAL 1 DAY),DATE_ADD(CURDATE(),INTERVAL 90 DAY),1,1001,99.00,
   1,0,1,1,1,-1,0,2,15,0,'','价格过低未通过',UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),NULL),
  (6006,'全天','1,2,3,4',DATE_SUB(CURDATE(),INTERVAL 1 DAY),DATE_ADD(CURDATE(),INTERVAL 90 DAY),1,1002,159.00,
   1,0,1,1,1,1,5,20,56,8,'', '',UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),NULL),
  (6007,'晚上','4',DATE_SUB(CURDATE(),INTERVAL 1 DAY),DATE_ADD(CURDATE(),INTERVAL 90 DAY),1,1001,129.00,
   1,0,1,1,1,0,0,3,30,0,'','',UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),NULL)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`), `seckill_time_ids`=VALUES(`seckill_time_ids`),
  `mer_id`=VALUES(`mer_id`), `product_id`=VALUES(`product_id`),
  `is_show`=VALUES(`is_show`),
  `product_status`=VALUES(`product_status`), `star`=VALUES(`star`), `sort`=VALUES(`sort`),
  `stock`=VALUES(`stock`), `sales`=VALUES(`sales`), `refusal`=VALUES(`refusal`),
  `delete_time`=VALUES(`delete_time`), `update_time`=VALUES(`update_time`);
