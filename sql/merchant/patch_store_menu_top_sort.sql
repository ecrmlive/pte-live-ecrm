-- 店铺菜单顶层 sort + Iconify 图标：对齐 CRMEB（https://mer.crmeb.net/admin/merchant/system）
-- 列表/权限树统一 ORDER BY sort DESC, id ASC，故数值越大越靠前。
-- 目标顺序：首页 → 商品 → 订单 → 营销 → 财务 → 用户 → 员工 → 装修 → 权限 → 设置 → 公告
-- 按名称匹配顶层（不依赖本库 ID 是否等于 CRMEB 55/95…）；sort 取互斥降序，避免并列依赖 id。
-- 用法（项目根目录）：
--   docker exec -i pte_live_mysql sh -ec 'MYSQL_PWD="$MYSQL_ROOT_PASSWORD" mysql -uroot --default-character-set=utf8mb4 qixi_crm_merchant' < sql/merchant/patch_store_menu_top_sort.sql

SET NAMES utf8mb4;
USE `qixi_crm_merchant`;

UPDATE `qixi_crm_m_menu`
SET `sort` = 100, `icon` = 'ant-design:home-outlined'
WHERE `parent_id` = 0 AND `name` = '首页';

UPDATE `qixi_crm_m_menu`
SET `sort` = 90, `icon` = 'ant-design:shopping-outlined'
WHERE `parent_id` = 0 AND `name` = '商品';

UPDATE `qixi_crm_m_menu`
SET `sort` = 80, `icon` = 'ant-design:file-text-outlined'
WHERE `parent_id` = 0 AND `name` = '订单';

UPDATE `qixi_crm_m_menu`
SET `sort` = 70, `icon` = 'ant-design:bell-outlined'
WHERE `parent_id` = 0 AND `name` = '营销';

UPDATE `qixi_crm_m_menu`
SET `sort` = 60, `icon` = 'ant-design:field-time-outlined'
WHERE `parent_id` = 0 AND `name` = '财务';

UPDATE `qixi_crm_m_menu`
SET `sort` = 50, `icon` = 'ant-design:user-outlined'
WHERE `parent_id` = 0 AND `name` = '用户';

UPDATE `qixi_crm_m_menu`
SET `sort` = 40, `icon` = 'ant-design:customer-service-outlined'
WHERE `parent_id` = 0 AND `name` = '员工';

UPDATE `qixi_crm_m_menu`
SET `sort` = 30, `icon` = 'ant-design:format-painter-outlined'
WHERE `parent_id` = 0 AND `name` = '装修';

UPDATE `qixi_crm_m_menu`
SET `sort` = 20, `icon` = 'ant-design:safety-outlined'
WHERE `parent_id` = 0 AND `name` = '权限';

UPDATE `qixi_crm_m_menu`
SET `sort` = 10, `icon` = 'ant-design:setting-outlined'
WHERE `parent_id` = 0 AND `name` = '设置';

UPDATE `qixi_crm_m_menu`
SET `sort` = 0, `icon` = 'ant-design:notification-outlined'
WHERE `parent_id` = 0 AND `name` = '公告';
