-- 公众号图文本地演示数据
SET NAMES utf8mb4;
USE `qixi_crm_admin`;

INSERT INTO `qixi_crm_a_wechat_news` (`wechat_news_id`, `status`, `items`, `create_time`, `update_time`)
SELECT
  9001,
  1,
  CAST('[{"title":"七禧商城秋季上新","author":"七禧运营","synopsis":"本地验收用的公众号图文摘要，介绍秋季精选好物。","image":"https://picsum.photos/seed/qixi-wechat-news-9001/800/450","content":"<p>本内容仅用于统一后台中文验收。</p><p>欢迎选购秋季上新商品。</p>"}]' AS JSON),
  NOW(),
  NOW()
FROM DUAL
WHERE NOT EXISTS (
  SELECT 1 FROM `qixi_crm_a_wechat_news` WHERE `wechat_news_id` = 9001
);

INSERT INTO `qixi_crm_a_wechat_news` (`wechat_news_id`, `status`, `items`, `create_time`, `update_time`)
SELECT
  9002,
  1,
  CAST('[{"title":"客服服务时间调整","author":"七禧客服","synopsis":"演示图文：说明节假日客服值班安排。","image":"https://picsum.photos/seed/qixi-wechat-news-9002a/800/450","content":"<p>工作日 9:00–18:00 在线答疑。</p>"},{"title":"配送时效说明","author":"七禧物流","synopsis":"演示多图文第二条。","image":"https://picsum.photos/seed/qixi-wechat-news-9002b/800/450","content":"<p>偏远地区配送时效以物流公司为准。</p>"}]' AS JSON),
  DATE_SUB(NOW(), INTERVAL 1 DAY),
  DATE_SUB(NOW(), INTERVAL 1 DAY)
FROM DUAL
WHERE NOT EXISTS (
  SELECT 1 FROM `qixi_crm_a_wechat_news` WHERE `wechat_news_id` = 9002
);
