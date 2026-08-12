-- 公告管理本地验收数据：仅使用虚构中文店铺与运营文本，不含真实用户、商户或联系方式。
SET NAMES utf8mb4;
USE `qixi_crm_admin`;

INSERT INTO `qixi_crm_a_notice`
  (`notice_id`,`title`,`content`,`is_show`,`sort`,`scope_type`,`is_del`,`create_time`)
VALUES
  (5203,'店铺经营服务通知','<p>本公告仅发送给已关联的本地演示店铺，用于验收指定店铺范围。</p>',1,0,'store_name',0,DATE_SUB(NOW(), INTERVAL 1 DAY)),
  (5204,'店铺类别运营提醒','<p>本公告仅发送给已关联店铺类别，用于验收类别关联范围。</p>',1,0,'store_type',0,DATE_SUB(NOW(), INTERVAL 2 DAY)),
  (5205,'店铺分类活动公告','<p>本公告仅发送给已关联店铺分类，用于验收分类关联范围。</p>',0,0,'store_category',0,DATE_SUB(NOW(), INTERVAL 3 DAY))
ON DUPLICATE KEY UPDATE
  `title`=VALUES(`title`),`content`=VALUES(`content`),`is_show`=VALUES(`is_show`),
  `scope_type`=VALUES(`scope_type`),`is_del`=0;

DELETE FROM `qixi_crm_a_notice_scope` WHERE `notice_id` IN (5203,5204,5205);
INSERT INTO `qixi_crm_a_notice_scope` (`notice_id`,`scope_kind`,`scope_id`) VALUES
  (5203,'store_name',1),(5203,'store_name',2),
  (5204,'store_type',711),
  (5205,'store_category',701);
