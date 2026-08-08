-- 平台「核销记录」菜单文案/图标对齐（幂等）
SET NAMES utf8mb4;

UPDATE `qixi_crm_a_menu`
SET `title` = '核销记录',
    `icon` = 'ant-design:audit-outlined',
    `route_path` = '/order/cancellation',
    `sort` = 3
WHERE `id` = 53;

UPDATE `qixi_crm_a_menu`
SET `title` = '查看核销记录',
    `route_path` = 'order/cancellation'
WHERE `id` = 20986 AND `code` = 'order.cancellation.read';
