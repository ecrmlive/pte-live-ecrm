-- 本地演示：付费会员类型（青铜/白银/黄金/钻石），幂等
USE `qixi_crm_business`;
SET NAMES utf8mb4;

-- 确保原价列存在（兼容未先跑 patch 的导入顺序）
SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_b_svip_plan' AND COLUMN_NAME='cost_price')=0,
    'ALTER TABLE `qixi_crm_b_svip_plan` ADD COLUMN `cost_price` decimal(12,2) NOT NULL DEFAULT 0.00 AFTER `name`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

INSERT INTO `qixi_crm_b_svip_plan`
  (`id`,`name`,`cost_price`,`price`,`plan_type`,`duration_days`,`benefits`,`status`,`sort`)
VALUES
  (980001,'青铜会员',99.00,29.00,'period',30,JSON_ARRAY('会员专享价'),1,10),
  (980002,'白银会员',199.00,79.00,'period',90,JSON_ARRAY('会员专享价','专属客服优先响应'),1,20),
  (980003,'黄金会员',399.00,299.00,'period',365,JSON_ARRAY('会员专享价','专属客服优先响应','年度权益礼包'),1,30),
  (980004,'钻石会员',999.00,599.00,'lifetime',NULL,JSON_ARRAY('会员专享价','专属客服优先响应','年度权益礼包'),1,40)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `cost_price`=VALUES(`cost_price`),
  `price`=VALUES(`price`),
  `plan_type`=VALUES(`plan_type`),
  `duration_days`=VALUES(`duration_days`),
  `benefits`=VALUES(`benefits`),
  `status`=VALUES(`status`),
  `sort`=VALUES(`sort`),
  `updated_at`=NOW();

-- 同步演示订单里的套餐名（不伪造支付成功）
UPDATE `qixi_crm_b_svip_order` SET `plan_name`='青铜会员', `plan_type`='period', `duration_days`=30, `amount`=29.00
WHERE `id`=980010;
UPDATE `qixi_crm_b_svip_order` SET `plan_name`='白银会员', `plan_type`='period', `duration_days`=90, `amount`=79.00
WHERE `id`=980011;
