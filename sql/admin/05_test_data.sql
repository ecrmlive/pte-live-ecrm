USE `qixi_crm_admin`;
SET NAMES utf8mb4;
-- 不初始化后台账号或密码。管理员必须通过受控初始化命令创建并写入密码哈希。

-- 本地验收用监管投影：不含真实个人信息；商户事实由 api-merchant 管理。
INSERT INTO `qixi_crm_a_merchant_view`
  (`merchant_id`,`merchant_name`,`contact_name`,`contact_mobile`,`region_id`,`status`)
VALUES (1,'七禧演示店铺','演示联系人','13900000000',NULL,1)
ON DUPLICATE KEY UPDATE `merchant_name`=VALUES(`merchant_name`),`status`=VALUES(`status`);

INSERT INTO `qixi_crm_a_diy_page` (`id`,`page_type`,`name`,`document`,`status`,`updated_by`) VALUES
  (4001,'home','七禧平台首页',JSON_OBJECT(
    'page',JSON_OBJECT('type','page','name','页面设置','params',JSON_OBJECT('name','七禧商城','title','七禧商城')),
    'items',JSON_ARRAY(JSON_OBJECT('type','banner','name','轮播图','data',JSON_ARRAY(JSON_OBJECT('imgName','七禧商城精选','imgUrl','/demo/home-hero-v1.png','linkUrl','/pages/goods/list')))),
    '_qixi',JSON_OBJECT('title','七禧商城','template_name','home','is_diy',1,'is_show',1,'is_default',1)
  ),'published',0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`),`document`=VALUES(`document`),`status`=VALUES(`status`);
