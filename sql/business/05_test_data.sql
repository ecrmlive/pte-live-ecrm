USE `qixi_crm_business`;
-- 仅提供无个人信息的消费读模型夹具。生产数据必须经商户事件同步写入，不能依赖本文件。
INSERT INTO `qixi_crm_b_store_view` (`store_id`,`merchant_id`,`store_app_id`,`store_name`,`status`) VALUES
  (1,1,'qixi.store.demo.1','七禧演示店铺',1)
ON DUPLICATE KEY UPDATE `merchant_id`=VALUES(`merchant_id`),`store_app_id`=VALUES(`store_app_id`),`store_name`=VALUES(`store_name`),`status`=VALUES(`status`);

INSERT INTO `qixi_crm_b_category_view` (`category_id`,`parent_id`,`name`,`sort`,`status`) VALUES
  (101,0,'服饰鞋包',10,1),(102,0,'家居生活',20,1),(103,0,'数码家电',30,1),
  (104,0,'美妆个护',40,1),(105,0,'食品生鲜',50,1),(106,0,'运动户外',60,1)
ON DUPLICATE KEY UPDATE `parent_id`=VALUES(`parent_id`),`name`=VALUES(`name`),`sort`=VALUES(`sort`),`status`=VALUES(`status`);

INSERT INTO `qixi_crm_b_product_view` (`product_id`,`merchant_id`,`store_id`,`merchant_name`,`store_name`,`category_id`,`title`,`cover_url`,`price`,`original_price`,`product_type`,`sales`,`stock`,`sale_status`,`version`,`updated_at`) VALUES
  (1001,1,1,'七禧演示商户','七禧演示店铺',101,'七禧臻选羊绒针织衫','',299.00,399.00,0,158,60,1,1,NOW()),
  (1002,1,1,'七禧演示商户','七禧演示店铺',101,'七禧头层牛皮通勤包','',469.00,599.00,0,126,32,1,1,NOW()),
  (1003,1,1,'七禧演示商户','七禧演示店铺',106,'七禧轻量跑步鞋','',369.00,459.00,0,97,48,1,1,NOW()),
  (1004,1,1,'七禧演示商户','七禧演示店铺',103,'七禧智能保温杯','',159.00,219.00,0,86,80,1,1,NOW()),
  (1005,1,1,'七禧演示商户','七禧演示店铺',104,'七禧真丝方巾','',129.00,169.00,0,72,90,1,1,NOW()),
  (1006,1,1,'七禧演示商户','七禧演示店铺',102,'七禧香氛礼盒','',239.00,299.00,0,61,40,1,1,NOW())
ON DUPLICATE KEY UPDATE `category_id`=VALUES(`category_id`),`title`=VALUES(`title`),`cover_url`=VALUES(`cover_url`),`price`=VALUES(`price`),`original_price`=VALUES(`original_price`),`sales`=VALUES(`sales`),`stock`=VALUES(`stock`),`sale_status`=VALUES(`sale_status`),`updated_at`=NOW();
