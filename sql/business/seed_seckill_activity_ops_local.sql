-- local：秒杀活动「用户操作」数据（销量/库存波动，供统计抽屉）
-- 依赖 patch_seckill_activity.sql 已写入活动场 1~5 与商品 6201~6231
USE `qixi_crm_business`;
SET NAMES utf8mb4;

-- 抬高进行中活动的销量，形成可见操作痕迹
UPDATE `qixi_crm_b_seckill_active`
SET `sales` = GREATEST(`sales`, 40),
    `stock` = GREATEST(`stock`, 10),
    `update_time` = UNIX_TIMESTAMP()
WHERE `activity_id` = 5 AND `delete_time` IS NULL;

UPDATE `qixi_crm_b_seckill_active`
SET `sales` = GREATEST(`sales`, 25),
    `update_time` = UNIX_TIMESTAMP()
WHERE `activity_id` = 4 AND `delete_time` IS NULL;

UPDATE `qixi_crm_b_seckill_active`
SET `sales` = GREATEST(`sales`, 20),
    `update_time` = UNIX_TIMESTAMP()
WHERE `activity_id` = 3 AND `delete_time` IS NULL;

-- 回填活动商品/店铺数（与真实挂载一致）
UPDATE `qixi_crm_b_seckill_activity` a
SET
  `product_count` = (
    SELECT COUNT(*) FROM `qixi_crm_b_seckill_active` p
    WHERE p.`activity_id` = a.`seckill_activity_id` AND p.`delete_time` IS NULL
  ),
  `merchant_count` = (
    SELECT COUNT(DISTINCT p.`mer_id`) FROM `qixi_crm_b_seckill_active` p
    WHERE p.`activity_id` = a.`seckill_activity_id` AND p.`delete_time` IS NULL
  ),
  `active_status` = CASE
    WHEN CURDATE() < a.`start_day` THEN 0
    WHEN CURDATE() > a.`end_day` THEN -1
    ELSE 1
  END
WHERE a.`seckill_activity_id` IN (1,2,3,4,5) AND a.`delete_time` IS NULL;
