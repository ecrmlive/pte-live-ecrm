USE `qixi_crm_admin`;
-- 不初始化后台账号或密码。管理员必须通过受控初始化命令创建并写入密码哈希。

-- 本地验收用监管投影：不含真实个人信息；商户事实由 api-merchant 管理。
INSERT INTO `qixi_crm_a_merchant_view`
  (`merchant_id`,`merchant_name`,`contact_name`,`contact_mobile`,`region_id`,`status`)
VALUES (1,'七禧演示店铺','演示联系人','13900000000',NULL,1)
ON DUPLICATE KEY UPDATE `merchant_name`=VALUES(`merchant_name`),`status`=VALUES(`status`);
