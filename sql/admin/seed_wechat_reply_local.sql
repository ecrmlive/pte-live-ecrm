-- 公众号自动回复演示：欢迎引导 + 数字菜单 1/2/3/4
SET NAMES utf8mb4;
USE `qixi_crm_admin`;

INSERT INTO `qixi_crm_a_wechat_reply`
  (`wechat_reply_id`,`reply_key`,`reply_type`,`content`,`status`,`sort`,`create_time`,`update_time`)
SELECT 8001,'subscribe','text',
'欢迎关注七禧商城！\n请回复数字获取服务：\n1. 商城首页\n2. 热门活动\n3. 售后须知\n4. 联系客服',
1,0,NOW(),NOW()
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM `qixi_crm_a_wechat_reply` WHERE `reply_key`='subscribe');

INSERT INTO `qixi_crm_a_wechat_reply`
  (`wechat_reply_id`,`reply_key`,`reply_type`,`content`,`status`,`sort`,`create_time`,`update_time`)
SELECT 8002,'default','text',
'暂未识别您的指令。\n请回复数字：\n1. 商城首页\n2. 热门活动\n3. 售后须知\n4. 联系客服',
1,0,NOW(),NOW()
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM `qixi_crm_a_wechat_reply` WHERE `reply_key`='default');

INSERT INTO `qixi_crm_a_wechat_reply`
  (`wechat_reply_id`,`reply_key`,`reply_type`,`content`,`status`,`sort`,`create_time`,`update_time`)
SELECT 8011,'1','text','【商城首页】\n打开商城选购好物：https://mer.crmeb.net',1,1,NOW(),NOW()
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM `qixi_crm_a_wechat_reply` WHERE `reply_key`='1');

INSERT INTO `qixi_crm_a_wechat_reply`
  (`wechat_reply_id`,`reply_key`,`reply_type`,`content`,`status`,`sort`,`create_time`,`update_time`)
SELECT 8012,'2','text','【热门活动】\n本周上新与限时优惠，请进入商城「活动」频道查看。',1,2,NOW(),NOW()
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM `qixi_crm_a_wechat_reply` WHERE `reply_key`='2');

INSERT INTO `qixi_crm_a_wechat_reply`
  (`wechat_reply_id`,`reply_key`,`reply_type`,`content`,`status`,`sort`,`create_time`,`update_time`)
SELECT 8013,'3','text','【售后须知】\n提交售后请保留包装；审核通过后按页面提示寄回并填写物流单号。',1,3,NOW(),NOW()
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM `qixi_crm_a_wechat_reply` WHERE `reply_key`='3');

INSERT INTO `qixi_crm_a_wechat_reply`
  (`wechat_reply_id`,`reply_key`,`reply_type`,`content`,`status`,`sort`,`create_time`,`update_time`)
SELECT 8014,'4','text','【联系客服】\n客服时间：工作日 9:00–18:00。\n也可在商城「我的-联系客服」发起会话。',1,4,NOW(),NOW()
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM `qixi_crm_a_wechat_reply` WHERE `reply_key`='4');
