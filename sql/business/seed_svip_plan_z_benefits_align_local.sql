-- 在 seed_svip_plan_local 之后执行：套餐权益快照对齐会员权益名称（演示数据）
USE `qixi_crm_business`;
SET NAMES utf8mb4;

UPDATE `qixi_crm_b_svip_plan`
SET `benefits` = JSON_ARRAY('会员专属价','专属客服'), `updated_at`=NOW()
WHERE `id`=980001;
UPDATE `qixi_crm_b_svip_plan`
SET `benefits` = JSON_ARRAY('会员专属价','专属客服','签到返利'), `updated_at`=NOW()
WHERE `id`=980002;
UPDATE `qixi_crm_b_svip_plan`
SET `benefits` = JSON_ARRAY('会员专属价','专属客服','会员优惠券'), `updated_at`=NOW()
WHERE `id`=980003;
UPDATE `qixi_crm_b_svip_plan`
SET `benefits` = JSON_ARRAY('会员专属价','专属客服','经验翻倍','会员优惠券'), `updated_at`=NOW()
WHERE `id`=980004;
