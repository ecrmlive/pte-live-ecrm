-- 本地验收：文章分类 + 中文文章演示（含可访问封面图）
USE `qixi_crm_admin`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_a_article_category` (`cid`,`title`,`info`,`image`,`status`,`sort`,`is_del`,`create_time`) VALUES
  (501,'商城公告','平台公告与服务通知','https://picsum.photos/seed/qixi-article-cate-notice/120/120',1,20,0,NOW()),
  (502,'选购指南','商品选购与搭配建议','https://picsum.photos/seed/qixi-article-cate-guide/120/120',1,10,0,NOW()),
  (503,'售后须知','售后与退换货说明','https://picsum.photos/seed/qixi-article-cate-aftersale/120/120',1,5,0,NOW()),
  (504,'生活','生活灵感与日常好物','https://picsum.photos/seed/qixi-article-cate-life/120/120',1,30,0,NOW())
ON DUPLICATE KEY UPDATE
  `title`=VALUES(`title`),
  `info`=VALUES(`info`),
  `image`=VALUES(`image`),
  `status`=VALUES(`status`),
  `sort`=VALUES(`sort`),
  `is_del`=VALUES(`is_del`);

INSERT INTO `qixi_crm_a_article`
  (`article_id`,`cid`,`title`,`author`,`image`,`synopsis`,`content`,`visit`,`sort`,`status`,`is_del`,`create_time`)
VALUES
  (
    5101,501,'七禧商城秋季服务公告','虚构运营编辑',
    'https://picsum.photos/seed/qixi-article-5101/400/300',
    '本地验收用的中文公告摘要，说明客服响应与配送时效。',
    '<p>本内容仅用于统一后台中文验收，不含真实用户或商户信息。</p><p>客服工作时间：工作日 9:00–18:00；节假日值班安排以站内公告为准。</p>',
    18,20,1,0,NOW()
  ),
  (
    5102,502,'居家香氛选购小贴士','虚构运营编辑',
    'https://picsum.photos/seed/qixi-article-5102/400/300',
    '本地验收用的中文选购指南，帮助用户按空间选择香氛。',
    '<p>请根据空间大小和使用场景选择香氛产品。</p><ul><li>客厅宜选清新木质调</li><li>卧室宜选柔和花香调</li></ul>',
    7,10,1,0,NOW()
  ),
  (
    5103,502,'换季收纳清单','虚构内容编辑',
    'https://picsum.photos/seed/qixi-article-5103/400/300',
    '演示文章：按衣柜分区整理换季衣物与配饰。',
    '<p>换季前建议先清点衣物，再按「常穿 / 留存 / 捐赠」三类整理。</p><p>收纳盒请标注品类与季节，便于下次取用。</p>',
    12,8,1,0,DATE_SUB(NOW(), INTERVAL 1 DAY)
  ),
  (
    5104,503,'退换货流程说明','虚构客服编辑',
    'https://picsum.photos/seed/qixi-article-5104/400/300',
    '演示文章：说明申请退换货的步骤与注意事项。',
    '<p>提交售后申请后，请保持商品完好并保留包装。</p><p>平台审核通过后，请按页面提示寄回并填写物流单号。</p>',
    25,6,1,0,DATE_SUB(NOW(), INTERVAL 2 DAY)
  ),
  (
    5105,501,'平台活动规则摘要','虚构运营编辑',
    'https://picsum.photos/seed/qixi-article-5105/400/300',
    '演示文章：营销活动参与条件与结算说明（虚构）。',
    '<p>活动库存与优惠力度以提交订单时系统校验结果为准。</p><p>本摘要仅为本地演示，不构成真实活动承诺。</p>',
    9,4,0,0,DATE_SUB(NOW(), INTERVAL 3 DAY)
  )
ON DUPLICATE KEY UPDATE
  `cid`=VALUES(`cid`),
  `title`=VALUES(`title`),
  `author`=VALUES(`author`),
  `image`=VALUES(`image`),
  `synopsis`=VALUES(`synopsis`),
  `content`=VALUES(`content`),
  `visit`=VALUES(`visit`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `is_del`=VALUES(`is_del`);
