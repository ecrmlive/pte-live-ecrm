-- 预售审核 Drawer「商品详情」Tab 读 qixi_crm_a_product_ops.content_html
USE `qixi_crm_admin`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_a_product_ops`
  (`product_id`,`is_used`,`star`,`rank_sort`,`refund_switch`,`once_min_count`,`content_html`,`updated_by`)
VALUES
  (1001,1,0,10,1,1,'<p>轻奢羊绒针织衫演示详情：保暖亲肤，适合秋冬通勤。</p>',0),
  (1002,1,0,20,1,1,'<p>头层牛皮通勤托特包演示详情：大容量分区，日常出行好搭档。</p>',0),
  (1004,1,0,30,1,1,'<p>精纺圆领羊毛开衫演示详情：全款/定金预售审核关联商品。</p>',0),
  (1007,1,0,40,1,1,'<p>柔软亲肤针织披肩演示详情。</p>',0),
  (1008,1,0,50,1,1,'<p>城市通勤训练跑鞋演示详情。</p>',0),
  (1101,1,0,60,1,1,'<p>无火藤条香氛礼盒演示详情。</p>',0),
  (1102,1,0,70,1,1,'<p>晚安助眠香薰蜡烛演示详情。</p>',0),
  (1103,1,0,80,1,1,'<p>恒温随行保温杯演示详情。</p>',0),
  (1104,1,0,90,1,1,'<p>晨间居家香氛套装演示详情。</p>',0),
  (1107,1,0,100,1,1,'<p>客厅氛围香薰礼盒演示详情。</p>',0),
  (1108,1,0,110,1,1,'<p>织物护理香氛喷雾演示详情：定金预售待审关联商品。</p>',0)
ON DUPLICATE KEY UPDATE
  `content_html`=VALUES(`content_html`),
  `is_used`=VALUES(`is_used`),
  `refund_switch`=VALUES(`refund_switch`),
  `once_min_count`=VALUES(`once_min_count`);
