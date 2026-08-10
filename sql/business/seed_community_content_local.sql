-- 社区内容/评论本地演示数据（中文 utf8mb4，幂等；复用分类/话题夹具）
USE `qixi_crm_business`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_b_user` (`id`,`nickname`,`mobile`,`status`) VALUES
  (9001,'晴空漫游者',NULL,1),(9002,'居家研究员',NULL,1),(9003,'通勤玩家',NULL,1),(9004,'短视频达人',NULL,1)
ON DUPLICATE KEY UPDATE `nickname`=VALUES(`nickname`),`status`=VALUES(`status`);

INSERT INTO `qixi_crm_b_social_post` (
  `community_id`,`title`,`image`,`category_id`,`topic_id`,`uid`,`mer_id`,`product_id`,
  `count_start`,`count_reply`,`status`,`is_show`,`is_hot`,`start`,`is_type`,`video_link`,
  `content`,`refusal`,`pv`,`is_del`,`status_time`
) VALUES
  (6601,'通勤针织衫的三种叠穿思路','/demo/product-knit-v1.png,/demo/product-knit-v2.png',6401,6501,9001,1,1001,
   26,3,1,1,1,5,1,'','柔软的针织衫适合和衬衫、半裙搭配，通勤和周末都能穿。','',128,0,NOW()),
  (6602,'周末短途出行的轻便穿搭','/demo/product-fragrance-v1.png',6402,6502,9002,2,1101,
   18,2,1,1,0,4,1,'','短途旅行建议选透气鞋与易叠穿外套，行李更轻便。','',96,0,NOW()),
  (6603,'桌面恒温杯垫值得入手吗','/demo/product-tumbler-v1.png',6403,6503,9003,3,1207,
   11,1,1,1,0,3,1,'','适合常坐在电脑前的人，配保温杯使用更方便。','',65,0,NOW()),
  (6604,'周末市集穿搭短视频','/demo/product-knit-v1.png',6401,6501,9004,1,1001,
   42,0,1,1,1,5,2,'https://demo.local/video/weekend-outfit.mp4','周末市集轻便穿搭，一镜到底展示叠穿层次。','',210,0,NOW()),
  (6605,'待审：阳台绿植摆放灵感','/demo/product-fragrance-v1.png',6402,6502,9002,2,0,
   0,0,0,1,0,1,1,'','分享阳台角落的绿植与香氛搭配，等待平台审核。','',0,0,NULL),
  (6606,'已下架：违规引流示例','/demo/product-tumbler-v1.png',6403,6503,9003,3,0,
   2,0,-2,0,0,1,1,'','本条仅为强制下架演示，不代表真实业务。','含站外引流信息，已强制下架',12,0,NOW())
ON DUPLICATE KEY UPDATE
  `title`=VALUES(`title`),`image`=VALUES(`image`),`category_id`=VALUES(`category_id`),`topic_id`=VALUES(`topic_id`),
  `uid`=VALUES(`uid`),`mer_id`=VALUES(`mer_id`),`product_id`=VALUES(`product_id`),
  `count_start`=VALUES(`count_start`),`count_reply`=VALUES(`count_reply`),`status`=VALUES(`status`),
  `is_show`=VALUES(`is_show`),`is_hot`=VALUES(`is_hot`),`start`=VALUES(`start`),`is_type`=VALUES(`is_type`),
  `video_link`=VALUES(`video_link`),`content`=VALUES(`content`),`refusal`=VALUES(`refusal`),
  `pv`=VALUES(`pv`),`is_del`=VALUES(`is_del`),`status_time`=VALUES(`status_time`);

INSERT INTO `qixi_crm_b_social_reply` (
  `reply_id`,`content`,`pid`,`uid`,`count_start`,`count_reply`,`community_id`,`status`,`refusal`,`is_del`,`create_time`
) VALUES
  (6701,'叠穿思路很实用，打算周末试试。',0,9002,5,1,6601,1,'',0,DATE_SUB(NOW(),INTERVAL 2 DAY)),
  (6702,'第二套更适合通勤。',6701,9001,2,0,6601,1,'',0,DATE_SUB(NOW(),INTERVAL 1 DAY)),
  (6703,'轻便穿搭很实用。',0,9001,3,0,6602,1,'',0,DATE_SUB(NOW(),INTERVAL 3 DAY)),
  (6704,'杯垫续航怎么样？',0,9002,1,0,6603,1,'',0,DATE_SUB(NOW(),INTERVAL 5 HOUR)),
  (6705,'待审评论：想看更多细节图。',0,9003,0,0,6601,0,'',0,DATE_SUB(NOW(),INTERVAL 1 HOUR)),
  (6706,'已拒绝：广告引流。',0,9004,0,0,6602,-1,'评论含站外联系方式',0,DATE_SUB(NOW(),INTERVAL 6 HOUR))
ON DUPLICATE KEY UPDATE
  `content`=VALUES(`content`),`pid`=VALUES(`pid`),`uid`=VALUES(`uid`),
  `count_start`=VALUES(`count_start`),`count_reply`=VALUES(`count_reply`),
  `community_id`=VALUES(`community_id`),`status`=VALUES(`status`),`refusal`=VALUES(`refusal`),
  `is_del`=VALUES(`is_del`),`create_time`=VALUES(`create_time`);

UPDATE `qixi_crm_b_social_post` p
SET `count_reply`=(SELECT COUNT(*) FROM `qixi_crm_b_social_reply` r WHERE r.community_id=p.community_id AND r.is_del=0)
WHERE p.community_id IN (6601,6602,6603,6604,6605,6606);
