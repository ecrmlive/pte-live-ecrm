USE `qixi_crm_admin`;
SET NAMES utf8mb4;

-- 统一后台文章夹具：仅为本地验收构造的中文内容，不含个人信息或真实素材凭据。
INSERT INTO `qixi_crm_a_article_category` (`cid`,`title`,`status`,`sort`,`is_del`) VALUES
  (501,'商城公告',1,20,0),(502,'选购指南',1,10,0)
ON DUPLICATE KEY UPDATE `title`=VALUES(`title`),`status`=VALUES(`status`),`sort`=VALUES(`sort`),`is_del`=VALUES(`is_del`);
INSERT INTO `qixi_crm_a_article` (`article_id`,`cid`,`title`,`author`,`image`,`synopsis`,`content`,`visit`,`sort`,`status`,`is_del`,`create_time`) VALUES
  (5101,501,'七禧商城秋季服务公告','虚构运营编辑','','本地验收用的中文公告摘要。','<p>本内容仅用于统一后台中文验收，不含真实用户或商户信息。</p>',18,20,1,0,NOW()),
  (5102,502,'居家香氛选购小贴士','虚构运营编辑','','本地验收用的中文选购指南。','<p>请根据空间大小和使用场景选择香氛产品。</p>',7,10,1,0,NOW())
ON DUPLICATE KEY UPDATE `cid`=VALUES(`cid`),`title`=VALUES(`title`),`author`=VALUES(`author`),`image`=VALUES(`image`),`synopsis`=VALUES(`synopsis`),`content`=VALUES(`content`),`visit`=VALUES(`visit`),`sort`=VALUES(`sort`),`status`=VALUES(`status`),`is_del`=VALUES(`is_del`);

INSERT INTO `qixi_crm_a_notice` (`notice_id`,`title`,`content`,`is_show`,`sort`,`is_del`,`create_time`) VALUES
  (5201,'七禧商城本地验收公告','本公告为统一后台中文模拟数据，不包含真实活动、用户或商户信息。',1,20,0,NOW()),
  (5202,'虚构售后服务提示','本地验收环境展示售后流程说明，实际规则以已发布协议为准。',1,10,0,NOW())
ON DUPLICATE KEY UPDATE `title`=VALUES(`title`),`content`=VALUES(`content`),`is_show`=VALUES(`is_show`),`sort`=VALUES(`sort`),`is_del`=VALUES(`is_del`);
INSERT INTO `qixi_crm_a_setting_cache` (`key`,`expire_time`,`result`) VALUES
  ('sys_user_agree',0,'七禧商城本地验收用户协议：本内容为虚构中文示例。'),
  ('sys_userr_privacy',0,'七禧商城本地验收隐私政策：不包含真实个人信息。'),
  ('sys_merchant_type',0,'<p>七禧店铺类型说明（本地验收虚构文案）：用于向入驻商户解释不同类型店铺的规则与权益。</p>'),
  ('sys_merchant_category',0,'<p>七禧店铺分类说明（本地验收虚构文案）：用于向入驻商户解释店铺分类的选择指引。</p>'),
  ('sms_config',0,'{"enabled":false,"provider":"stub","sign":"七禧商城","remark":"本地验收未配置通道"}'),
  ('merchant_apply_setting',0,'{"background_image":"","form_fields":[{"id":"demo_area","type":"text","title":"营业面积","content_type":"number","default_value":"","placeholder":"请输入营业面积","required":false}]}')
ON DUPLICATE KEY UPDATE `expire_time`=VALUES(`expire_time`),`result`=VALUES(`result`);

-- 客服策略夹具只包含虚构中文业务文本，不包含 IM Token、UserSig、云密钥或真实值班信息。
INSERT INTO `qixi_crm_a_config` (`config_key`,`config_value`,`updated_by`) VALUES
  ('customer_service.settings',JSON_OBJECT('auto_reply_enabled',TRUE,'auto_reply_text','您好，虚构演示客服将在工作时间内回复您。','queue_mode','round_robin','max_sessions_per_agent',12),NULL)
ON DUPLICATE KEY UPDATE `config_value`=VALUES(`config_value`),`updated_by`=NULL,`updated_at`=CURRENT_TIMESTAMP;

-- 商户入驻分类夹具：佣金比例是平台规则，不复用商品分类，均为中文模拟数据。
INSERT INTO `qixi_crm_a_merchant_category` (`id`,`name`,`commission_rate`,`status`) VALUES
  (701,'服饰鞋包商户',8.50,1),
  (702,'家居生活商户',6.00,1)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`),`commission_rate`=VALUES(`commission_rate`),`status`=VALUES(`status`);

-- 店铺类型夹具只描述虚构运营规则，不代表真实收费或准入承诺。
INSERT INTO `qixi_crm_a_merchant_type` (`id`,`name`,`type_info`,`is_margin`,`margin`,`description`,`remark`,`status`) VALUES
  (711,'七禧基础演示店','适用于本地验收的基础虚构店铺类型',0,0.00,'本说明仅用于本地中文验收，不构成真实招商或收费规则。','虚构类型备注',1),
  (712,'七禧保证金演示店','展示保证金校验的虚构类型',1,500.00,'本类型用于验证保证金规则和店铺权限配置。','保证金为模拟金额',1)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`),`type_info`=VALUES(`type_info`),`is_margin`=VALUES(`is_margin`),`margin`=VALUES(`margin`),`description`=VALUES(`description`),`remark`=VALUES(`remark`),`status`=VALUES(`status`);

-- 默认演示区域：入驻审核通过时落到申请 region_id（本地无真实商圈数据）。
INSERT INTO `qixi_crm_a_region` (`id`,`parent_id`,`name`,`code`,`status`) VALUES
  (1,0,'七禧默认演示区域','demo-default',1)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`),`status`=VALUES(`status`);
INSERT INTO `qixi_crm_a_merchant_type_menu` (`merchant_type_id`,`menu_code`) VALUES
  (711,'merchant.dashboard'),(712,'merchant.dashboard'),(712,'merchant.catalog')
ON DUPLICATE KEY UPDATE `menu_code`=VALUES(`menu_code`);

-- 保证金夹具为本地验收的虚构资金流水，不含收款账户、银行卡、支付渠道或真实凭据。
INSERT INTO `qixi_crm_a_merchant_deposit_account` (`merchant_id`,`required_amount`,`available_amount`,`state`) VALUES
  (1,500.00,380.00,'shortfall'),(2,300.00,300.00,'funded')
ON DUPLICATE KEY UPDATE `required_amount`=VALUES(`required_amount`),`available_amount`=VALUES(`available_amount`),`state`=VALUES(`state`);
INSERT INTO `qixi_crm_a_merchant_deposit_ledger` (`id`,`merchant_id`,`entry_type`,`amount`,`balance_after`,`reason`,`idempotency_key`,`operator_admin_id`) VALUES
  (7201,1,'fund',500.00,500.00,'虚构中文保证金缴纳','fixture-deposit-fund-1',0),
  (7202,1,'deduct',120.00,380.00,'虚构中文违规扣减','fixture-deposit-deduct-1',0)
ON DUPLICATE KEY UPDATE `reason`=VALUES(`reason`),`balance_after`=VALUES(`balance_after`);
INSERT INTO `qixi_crm_a_merchant_deposit_refund` (`id`,`merchant_id`,`amount`,`status`,`reason`,`review_note`) VALUES
  (7301,2,300.00,'applied','虚构中文退还保证金申请','等待平台审核')
ON DUPLICATE KEY UPDATE `amount`=VALUES(`amount`),`status`=VALUES(`status`),`reason`=VALUES(`reason`),`review_note`=VALUES(`review_note`);

-- 分账申请夹具仅提供平台审核状态，不包含任何收款账户或渠道配置。
INSERT INTO `qixi_crm_a_merchant_profitsharing_application` (`id`,`merchant_id`,`merchant_name`,`application_no`,`applyment_id`,`status`,`description`,`message`,`review_note`) VALUES
  (7401,1,'七禧演示店铺','PS-DEMO-ZH-001','','applied','七禧演示店铺提交的虚构分账申请。','','等待平台审核')
ON DUPLICATE KEY UPDATE `merchant_name`=VALUES(`merchant_name`),`status`=VALUES(`status`),`description`=VALUES(`description`),`message`=VALUES(`message`),`review_note`=VALUES(`review_note`);

-- 商品元数据夹具全部为虚构中文展示语，不包含供应商、收款或个人信息。
INSERT INTO `qixi_crm_a_product_label` (`id`,`name`,`description`,`color`,`sort`,`status`) VALUES
  (7501,'七天无理由','适用于演示商品的售后保障标签','#16a34a',100,1),
  (7502,'新品尝鲜','本地验收用新品标识','#2563eb',90,1)
ON DUPLICATE KEY UPDATE `description`=VALUES(`description`),`color`=VALUES(`color`),`sort`=VALUES(`sort`),`status`=VALUES(`status`);
INSERT INTO `qixi_crm_a_product_guarantee` (`id`,`name`,`content`,`icon_url`,`sort`,`status`,`mer_count`,`product_count`) VALUES
  (7511,'正品保障','平台演示：商品来源与售后承诺以订单规则为准。','',100,1,0,0),
  (7512,'极速退款','平台演示：符合规则的退款进入售后状态机处理。','',90,1,0,0)
ON DUPLICATE KEY UPDATE `content`=VALUES(`content`),`icon_url`=VALUES(`icon_url`),`sort`=VALUES(`sort`),`status`=VALUES(`status`),`mer_count`=VALUES(`mer_count`),`product_count`=VALUES(`product_count`);
INSERT INTO `qixi_crm_a_product_parameter_template`
  (`id`,`name`,`cate_ids_json`,`params_json`,`values_json`,`sort`,`status`) VALUES
  (
    7521,'测试',JSON_ARRAY(7635,7628),
    JSON_ARRAY(
      JSON_OBJECT('name','颜色','values',JSON_ARRAY('中国红','竹青色','云白色'),'required',0,'sort',100),
      JSON_OBJECT('name','尺码','values',JSON_ARRAY('S','M','L'),'required',0,'sort',90)
    ),
    JSON_ARRAY('中国红','竹青色','云白色'),0,1
  ),
  (
    7522,'通用参数',JSON_ARRAY(7636,7632,7629,7621,7635,7628,7631),
    JSON_ARRAY(
      JSON_OBJECT('name','品牌','values',JSON_ARRAY('七禧','演示品牌'),'required',0,'sort',100),
      JSON_OBJECT('name','产地','values',JSON_ARRAY('中国','进口'),'required',0,'sort',90)
    ),
    JSON_ARRAY('七禧','演示品牌'),0,1
  )
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),`cate_ids_json`=VALUES(`cate_ids_json`),`params_json`=VALUES(`params_json`),
  `values_json`=VALUES(`values_json`),`sort`=VALUES(`sort`),`status`=VALUES(`status`);
INSERT INTO `qixi_crm_a_product_price_rule`
  (`id`,`name`,`cate_ids_json`,`is_default`,`content`,`sort`,`status`) VALUES
  (7531,'集成灶',JSON_ARRAY(7635,7628,7631,7632,7621),0,'<p><strong>演示：集成灶价格与规格差异说明</strong></p>',1,1),
  (7532,'MOCO针织开衫假两件',JSON_ARRAY(7635,7628,7631,7632,7621),0,'<p>演示价格说明：含优惠、规格与运费差异提示。</p>',0,1),
  (7533,'通用价格说明',JSON_ARRAY(),1,'<p>未指定分类时默认适用于全部商品。</p>',0,1),
  (7534,'服饰专场说明',JSON_ARRAY(7631),0,'<p>服饰类演示：尺码色差与运费说明。</p>',2,1)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),`cate_ids_json`=VALUES(`cate_ids_json`),`is_default`=VALUES(`is_default`),
  `content`=VALUES(`content`),`sort`=VALUES(`sort`),`status`=VALUES(`status`);

-- 店铺分组夹具用于三级分组、店铺关联和装修模板绑定验收；仅使用中文虚构店铺。
INSERT INTO `qixi_crm_a_store_group` (`id`,`parent_id`,`path`,`level`,`name`,`sort`,`status`,`diy_page_id`,`positioning_status`,`longitude`,`latitude`,`address`) VALUES
  (801,0,'/801/',0,'七禧中文演示商圈',20,1,0,1,121.4737000,31.2304000,'虚构的上海本地验收地址'),
  (802,801,'/801/802/',1,'华东精选店铺',10,1,0,0,NULL,NULL,'')
ON DUPLICATE KEY UPDATE `parent_id`=VALUES(`parent_id`),`path`=VALUES(`path`),`level`=VALUES(`level`),`name`=VALUES(`name`),`sort`=VALUES(`sort`),`status`=VALUES(`status`),`diy_page_id`=VALUES(`diy_page_id`),`positioning_status`=VALUES(`positioning_status`),`longitude`=VALUES(`longitude`),`latitude`=VALUES(`latitude`),`address`=VALUES(`address`);
INSERT INTO `qixi_crm_a_store_group_merchant` (`store_group_id`,`merchant_id`) VALUES
  (801,1),(802,2)
ON DUPLICATE KEY UPDATE `created_at`=CURRENT_TIMESTAMP;

-- 演示素材：勿写入不存在于 COS 的 /demo/* 相对路径（浏览器会显示「加载失败」）。
-- 分类本身由 init_data 种子维护；真实封面/视频请在素材库上传后选用。
DELETE FROM `qixi_crm_a_attachment_category`
 WHERE `mer_id`=0 AND `attachment_category_enname` IN ('demo-image','demo-video');
DELETE FROM `qixi_crm_a_attachment_asset`
 WHERE `attachment_id` IN (5311,5312)
    OR `attachment_src` LIKE '/demo/%';
-- 不初始化后台账号或密码。管理员必须通过受控初始化命令创建并写入密码哈希。

-- 登录日志夹具不含密码、令牌或真实 IP，仅用于平台安全审计列表验收。
INSERT INTO `qixi_crm_a_login_log` (`id`,`admin_user_id`,`username`,`role_code`,`success`,`ip`,`user_agent`,`created_at`) VALUES
  (5501,NULL,'platform-demo','platform',0,'127.0.0.1','本地验收浏览器',DATE_SUB(NOW(),INTERVAL 3 HOUR)),
  (5502,9901,'platform-demo','platform',1,'127.0.0.1','本地验收浏览器',DATE_SUB(NOW(),INTERVAL 2 HOUR))
ON DUPLICATE KEY UPDATE `admin_user_id`=VALUES(`admin_user_id`),`role_code`=VALUES(`role_code`),`success`=VALUES(`success`),`ip`=VALUES(`ip`),`user_agent`=VALUES(`user_agent`),`created_at`=VALUES(`created_at`);

-- 统一后台操作日志夹具：仅保存虚构管理员 ID、资源标识与请求号，绝不保存请求体、密码、令牌或个人资料。
INSERT INTO `qixi_crm_a_operation_log` (`id`,`admin_user_id`,`role_code`,`action`,`resource_type`,`resource_id`,`request_id`,`created_at`) VALUES
  (5601,9901,'platform','POST /api/platform/v1/user-list/9101/coupons/3002/issue','user-list','3002','fixture-admin-operation-5601',DATE_SUB(NOW(),INTERVAL 2 HOUR)),
  (5602,9901,'platform','PUT /api/platform/v1/svip/interests/1','svip','1','fixture-admin-operation-5602',DATE_SUB(NOW(),INTERVAL 1 HOUR))
ON DUPLICATE KEY UPDATE `role_code`=VALUES(`role_code`),`action`=VALUES(`action`),`resource_type`=VALUES(`resource_type`),`resource_id`=VALUES(`resource_id`),`request_id`=VALUES(`request_id`),`created_at`=VALUES(`created_at`);

-- 本地验收用监管投影：不含真实个人信息；商户事实由 api-merchant 管理。
-- mer_avatar 初始为空：勿写不存在的 /demo 占位图；已存在行勿覆盖用户真实封面。
INSERT INTO `qixi_crm_a_merchant_view`
  (`merchant_id`,`merchant_name`,`owner_name`,`contact_name`,`contact_mobile`,`address`,`category_id`,`type_id`,`business_id`,`region_id`,`status`,`is_best`,`offline_pay`,`is_trader`,`is_audit`,`is_bro_room`,`is_bro_goods`,`commission_switch`,`commission_rate`,`mer_keyword`,`mer_info`,`mer_account`,`sub_mchid`,`applyment_id`,`care_count`,`care_ficti`,`sort`,`mark`,`mer_avatar`)
VALUES
  (1,'七禧演示店铺','七禧演示商户','演示联系人','13900000000','上海市黄浦区演示路 1 号',701,711,30,10,1,1,0,1,1,0,0,1,2.50,'演示旗舰 居家','七禧本地验收旗舰店简介','demo_store_1','','',128,20,10,'本地验收备注：演示旗舰店',''),
  (2,'七禧居家优选店','七禧居家商户','演示联系人','13900000001','上海市徐汇区演示路 2 号',702,711,31,20,1,0,1,0,1,1,0,0,0.00,'居家优选','七禧居家优选店简介','demo_store_2','','',56,8,20,'本地验收备注：居家优选','')
ON DUPLICATE KEY UPDATE
  `merchant_name`=VALUES(`merchant_name`),`owner_name`=VALUES(`owner_name`),`contact_name`=VALUES(`contact_name`),
  `contact_mobile`=VALUES(`contact_mobile`),`address`=VALUES(`address`),`category_id`=VALUES(`category_id`),
  `type_id`=VALUES(`type_id`),`business_id`=VALUES(`business_id`),`region_id`=VALUES(`region_id`),`status`=VALUES(`status`),
  `is_best`=VALUES(`is_best`),`offline_pay`=VALUES(`offline_pay`),`is_trader`=VALUES(`is_trader`),
  `is_audit`=VALUES(`is_audit`),`is_bro_room`=VALUES(`is_bro_room`),`is_bro_goods`=VALUES(`is_bro_goods`),
  `commission_switch`=VALUES(`commission_switch`),`commission_rate`=VALUES(`commission_rate`),
  `mer_keyword`=VALUES(`mer_keyword`),`mer_info`=VALUES(`mer_info`),`mer_account`=VALUES(`mer_account`),
  `sub_mchid`=VALUES(`sub_mchid`),`applyment_id`=VALUES(`applyment_id`),`care_count`=VALUES(`care_count`),
  `care_ficti`=VALUES(`care_ficti`),`sort`=VALUES(`sort`),`mark`=VALUES(`mark`);

-- 区域商圈与代理夹具仅用于本地验收；结算账号刻意为空，不包含真实收款资料。
INSERT INTO `qixi_crm_a_business_zone` (`circle_id`,`pid`,`path`,`name`,`circle_agent_id`,`commission_type`,`commission_rate`,`level`,`remark`,`sort`,`status`,`type`,`role_id`,`business_store_category`,`business_store_type`) VALUES
  (10,0,'/10/','华东中文演示区域',801,1,8.50,0,'用于区域数据范围与代理审核本地验收。',20,1,0,0,0,0),
  (20,0,'/20/','华南中文演示区域',802,0,0.00,0,'仅包含虚构商圈与代理资料。',10,1,0,0,0,0),
  (30,0,'/30/','七禧演示商户',0,0,0.00,0,'店铺所属商户（type=1）本地验收夹具。',30,1,1,0,0,0),
  (31,0,'/31/','七禧居家商户',0,0,0.00,0,'店铺所属商户（type=1）本地验收夹具。',40,1,1,0,0,0)
ON DUPLICATE KEY UPDATE `pid`=VALUES(`pid`),`path`=VALUES(`path`),`name`=VALUES(`name`),`circle_agent_id`=VALUES(`circle_agent_id`),`commission_type`=VALUES(`commission_type`),`commission_rate`=VALUES(`commission_rate`),`remark`=VALUES(`remark`),`sort`=VALUES(`sort`),`status`=VALUES(`status`),`type`=VALUES(`type`);
INSERT INTO `qixi_crm_a_business_zone_agent` (`circle_agent_id`,`uid`,`name`,`phone`,`qualification`,`remark`,`status`,`payment_method`,`payment_name`,`payment_account`,`payment_bank`,`payment_qr_img`,`balance`,`type`,`business_name`,`business_store_category`,`business_store_type`) VALUES
  (801,0,'区域主管王小明','13900000010','虚构中文区域代理资质材料。','本地验收通过的区域代理。',1,0,'本地演示结算主体','','','',320.50,0,'',0,0),
  (802,0,'商圈运营赵小七','13900000020','虚构中文商圈代理申请资料。','本地验收待审核代理。',0,1,'本地演示结算主体','','','',0.00,0,'',0,0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`),`phone`=VALUES(`phone`),`qualification`=VALUES(`qualification`),`remark`=VALUES(`remark`),`status`=VALUES(`status`),`payment_method`=VALUES(`payment_method`),`payment_name`=VALUES(`payment_name`),`balance`=VALUES(`balance`),`type`=VALUES(`type`),`business_name`=VALUES(`business_name`);

-- 平台物流夹具为虚构中文物流配置，不包含真实承运商账号、接口密钥或收件信息。
INSERT INTO `qixi_crm_a_express` (`express_id`,`name`,`code`,`sort`,`is_show`,`is_del`) VALUES
  (601,'七禧演示快递','qixi-demo-express',20,1,0),
  (602,'本地验收物流','local-acceptance-logistics',10,1,0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`),`code`=VALUES(`code`),`sort`=VALUES(`sort`),`is_show`=VALUES(`is_show`),`is_del`=VALUES(`is_del`);
INSERT INTO `qixi_crm_a_city` (`city_id`,`parent_id`,`name`,`level`,`is_show`) VALUES
  (1,0,'中国',1,1),
  (310000,1,'上海市',2,1),(310100,310000,'上海市区',3,1),
  (310101,310100,'黄浦区',4,1),(310104,310100,'徐汇区',4,1),(310105,310100,'长宁区',4,1),
  (310106,310100,'静安区',4,1),(310107,310100,'普陀区',4,1),(310115,310100,'浦东新区',4,1)
ON DUPLICATE KEY UPDATE `parent_id`=VALUES(`parent_id`),`name`=VALUES(`name`),`level`=VALUES(`level`),`is_show`=VALUES(`is_show`);

-- 平台结算监管投影夹具：均为无个人信息的中文模拟数据，不能作为真实打款凭证。
INSERT INTO `qixi_crm_a_merchant_settlement_view`
  (`source_settlement_id`,`merchant_id`,`store_id`,`merchant_name`,`region_id`,`period_start`,`period_end`,`amount`,`status`,`updated_at`)
VALUES
  (7001,1,1,'七禧演示店铺',10,'2026-07-01 00:00:00','2026-07-31 23:59:59',1280.50,'withdraw_applied','2026-08-01 09:00:00'),
  (7002,2,2,'七禧居家优选店',20,'2026-06-01 00:00:00','2026-06-30 23:59:59',960.00,'paid','2026-07-05 15:30:00')
ON DUPLICATE KEY UPDATE `merchant_name`=VALUES(`merchant_name`),`amount`=VALUES(`amount`),`status`=VALUES(`status`),`updated_at`=VALUES(`updated_at`);

INSERT INTO `qixi_crm_a_diy_page` (`id`,`page_type`,`name`,`document`,`status`,`updated_by`) VALUES
  (4001,'home','七禧平台首页',JSON_OBJECT(
    'page',JSON_OBJECT('type','page','name','页面设置','params',JSON_OBJECT('name','七禧商城','title','七禧商城')),
    'items',JSON_ARRAY(
      JSON_OBJECT('type','banner','name','轮播图','data',JSON_ARRAY(
        JSON_OBJECT('imgName','七禧商城精选','imgUrl','/demo/home-hero-v1.png','linkUrl','/goods?cate_id=101'),
        JSON_OBJECT('imgName','七禧香氛家居','imgUrl','/demo/home-hero-fragrance-v1.png','linkUrl','/goods?cate_id=102'),
        JSON_OBJECT('imgName','七禧箱包配饰','imgUrl','/demo/home-hero-accessories-v1.png','linkUrl','/goods?cate_id=10102')
      )),
      JSON_OBJECT('type','product','name','服饰鞋包展示类型','params',JSON_OBJECT(
        'source','auto','auto',JSON_OBJECT('category',101,'showNum',4,'productSort','sales')
      )),
      JSON_OBJECT('type','product','name','家居生活展示类型','params',JSON_OBJECT(
        'source','auto','auto',JSON_OBJECT('category',102,'showNum',4,'productSort','sales')
      ))
    ),
    '_qixi',JSON_OBJECT('title','七禧商城','template_name','home','is_diy',1,'is_show',1,'is_default',1)
  ),'published',0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`),`document`=VALUES(`document`),`status`=VALUES(`status`);

-- 营销装饰夹具：氛围/边框/专题/报名，均为中文演示配置，不含密钥。
INSERT INTO `qixi_crm_a_marketing_decor`
  (`id`,`decor_type`,`name`,`code`,`cover_url`,`remark`,`payload`,`status`,`sort`,`starts_at`,`ends_at`,`is_del`)
VALUES
  (8001,'atmosphere','夏日焕新氛围图','summer-atmosphere','/demo/home-hero-v1.png','中文演示氛围',JSON_OBJECT('theme','summer'),1,10,DATE_SUB(NOW(),INTERVAL 1 DAY),DATE_ADD(NOW(),INTERVAL 180 DAY),0),
  (8002,'border','精选商品边框','featured-border','/demo/home-service-wide-v1.png','中文演示边框',JSON_OBJECT('style','gold'),1,20,NULL,NULL,0),
  (8003,'topic','居家香氛专题','home-fragrance-topic','/demo/home-hero-fragrance-v1.png','中文演示专题',JSON_OBJECT('category_id',102),1,30,NULL,NULL,0),
  (8004,'application','新品内测报名','early-access','','中文演示报名',JSON_OBJECT('quota',100),1,40,NULL,NULL,0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`),`code`=VALUES(`code`),`cover_url`=VALUES(`cover_url`),`remark`=VALUES(`remark`),`payload`=VALUES(`payload`),`status`=VALUES(`status`),`sort`=VALUES(`sort`),`is_del`=VALUES(`is_del`);

-- 维护/运营配置条目夹具：热搜、组合数据、系统表单、备份登记（中文演示，不含密钥）。
INSERT INTO `qixi_crm_a_config_item`
  (`id`,`item_type`,`name`,`code`,`remark`,`payload`,`status`,`sort`,`is_del`)
VALUES
  (8101,'hot_search','夏日香氛','summer-fragrance','中文演示热搜',JSON_OBJECT('weight',10),1,10,0),
  (8102,'group_data','首页金刚区','home-icon-grid','中文演示组合数据',JSON_OBJECT('slot','home'),1,20,0),
  (8103,'system_form','售后申请表','aftersale-form','中文演示系统表单',JSON_OBJECT('fields',JSON_ARRAY('reason','images')),1,30,0),
  (8104,'backup','本地日备登记','daily-local','中文演示备份登记',JSON_OBJECT('scope','admin'),1,40,0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`),`code`=VALUES(`code`),`remark`=VALUES(`remark`),`payload`=VALUES(`payload`),`status`=VALUES(`status`),`sort`=VALUES(`sort`),`is_del`=VALUES(`is_del`);
