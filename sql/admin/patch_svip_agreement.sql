-- 平台「付费会员 / 会员协议」：协议键 sys_svip（CRMEB CacheRepository::SYS_SVIP）+ 按钮权限
-- 用法（项目根目录）：
--   docker exec -i pte_live_mysql sh -ec 'MYSQL_PWD="$MYSQL_ROOT_PASSWORD" mysql -uroot --default-character-set=utf8mb4' < sql/admin/patch_svip_agreement.sql
-- 或：make local-sync-sql

SET NAMES utf8mb4;

USE `qixi_crm_admin`;

-- 确保导航页存在且启用
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (533,73,'user.svip.agreement','会员协议','lucide:file-text','/user/member/vipAgreement','page',4,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

-- RequireAdminMenu 仅认 kind=button；page 码不能用于接口鉴权。
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (21070,533,'user.svip.agreement.read','查看会员协议','','user/member/vipAgreement','button',1,1),
  (21071,533,'user.svip.agreement.manage','维护会员协议','','user/member/vipAgreement','button',2,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` AS r
CROSS JOIN `qixi_crm_a_menu` AS m
WHERE r.code IN ('platform','operations')
  AND m.code IN (
    'user.svip',
    'user.svip.agreement',
    'user.svip.agreement.read',
    'user.svip.agreement.manage'
  );

INSERT INTO `qixi_crm_a_setting_cache` (`key`,`expire_time`,`result`) VALUES
  ('sys_svip',0,'<h2 style="text-align:center;">付费会员服务协议</h2><p>欢迎开通本商城付费会员服务。请您在开通或续费前仔细阅读本协议。一旦您完成支付并开通服务，即视为您已阅读、理解并同意本协议全部条款。</p><h3>一、服务说明</h3><p>1. 付费会员是本商城向注册用户提供的增值会员服务，会员类型、售价、有效期及权益内容以开通页面实时公示为准。</p><p>2. 会员权益可能包括专属折扣、专属活动、积分加成、专属标识等，具体以平台公示内容为准。平台可根据运营需要调整权益方案，并将通过页面公告等方式告知。</p><p>3. 会员服务期限自开通成功之时起算，至对应套餐到期日止；到期后未续费的，相关付费权益将自动失效。</p><h3>二、开通与续费</h3><p>1. 您应使用本人账户完成开通与续费，并确保支付账户合法合规。</p><p>2. 支付成功后，系统将按订单结果即时或在合理时间内开通/延长会员权益；如因网络或系统原因导致状态延迟，请以账户内会员状态展示为准。</p><p>3. 除法律法规另有规定或平台规则明确承诺外，付费会员服务一经开通，费用不予退还。</p><h3>三、使用规范</h3><p>1. 会员权益仅限您本人账户使用，不得转让、出租、出借或以其他方式提供给第三方使用。</p><p>2. 您不得利用会员权益从事违法违规、恶意套取优惠、扰乱交易秩序等行为；一经发现，平台有权暂停或终止相关权益，并保留追究责任的权利。</p><p>3. 若您的账户因违规被限制、注销或冻结，已开通的付费会员权益可能同步受限，平台不因此承担额外赔偿责任。</p><h3>四、变更与中止</h3><p>1. 为持续优化服务，平台可能调整会员规则、权益内容或收费标准，调整前将通过站内公告等方式提示。</p><p>2. 如遇系统维护、不可抗力或第三方服务异常，相关会员功能可能暂时无法使用，平台将尽力恢复。</p><h3>五、其他</h3><p>1. 本协议未尽事宜，适用本商城用户协议、隐私政策及相关平台规则。</p><p>2. 本协议的解释权归本商城运营方所有。如有疑问，请通过商城内客服入口联系咨询（本页不展示真实电话号码）。</p>')
ON DUPLICATE KEY UPDATE
  `expire_time`=VALUES(`expire_time`),
  `result`=IF(`result` IS NULL OR `result`='',VALUES(`result`),`result`);

USE `qixi_crm_business`;
INSERT INTO `qixi_crm_b_content_view` (`content_id`,`content_type`,`title`,`cover_url`,`body`,`status`,`version`,`published_at`,`updated_at`) VALUES
  (2111,'agreement','sys_svip','','<h2 style="text-align:center;">付费会员服务协议</h2><p>欢迎开通本商城付费会员服务。请您在开通或续费前仔细阅读本协议。一旦您完成支付并开通服务，即视为您已阅读、理解并同意本协议全部条款。</p><h3>一、服务说明</h3><p>1. 付费会员是本商城向注册用户提供的增值会员服务，会员类型、售价、有效期及权益内容以开通页面实时公示为准。</p><p>2. 会员权益可能包括专属折扣、专属活动、积分加成、专属标识等，具体以平台公示内容为准。平台可根据运营需要调整权益方案，并将通过页面公告等方式告知。</p><p>3. 会员服务期限自开通成功之时起算，至对应套餐到期日止；到期后未续费的，相关付费权益将自动失效。</p><h3>二、开通与续费</h3><p>1. 您应使用本人账户完成开通与续费，并确保支付账户合法合规。</p><p>2. 支付成功后，系统将按订单结果即时或在合理时间内开通/延长会员权益；如因网络或系统原因导致状态延迟，请以账户内会员状态展示为准。</p><p>3. 除法律法规另有规定或平台规则明确承诺外，付费会员服务一经开通，费用不予退还。</p><h3>三、使用规范</h3><p>1. 会员权益仅限您本人账户使用，不得转让、出租、出借或以其他方式提供给第三方使用。</p><p>2. 您不得利用会员权益从事违法违规、恶意套取优惠、扰乱交易秩序等行为；一经发现，平台有权暂停或终止相关权益，并保留追究责任的权利。</p><p>3. 若您的账户因违规被限制、注销或冻结，已开通的付费会员权益可能同步受限，平台不因此承担额外赔偿责任。</p><h3>四、变更与中止</h3><p>1. 为持续优化服务，平台可能调整会员规则、权益内容或收费标准，调整前将通过站内公告等方式提示。</p><p>2. 如遇系统维护、不可抗力或第三方服务异常，相关会员功能可能暂时无法使用，平台将尽力恢复。</p><h3>五、其他</h3><p>1. 本协议未尽事宜，适用本商城用户协议、隐私政策及相关平台规则。</p><p>2. 本协议的解释权归本商城运营方所有。如有疑问，请通过商城内客服入口联系咨询（本页不展示真实电话号码）。</p>',1,1,NOW(),NOW())
ON DUPLICATE KEY UPDATE
  `content_type`=VALUES(`content_type`),
  `title`=VALUES(`title`),
  `body`=IF(`body` IS NULL OR `body`='',VALUES(`body`),`body`),
  `status`=1,
  `updated_at`=NOW();
