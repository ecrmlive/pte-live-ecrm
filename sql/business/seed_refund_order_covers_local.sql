-- 退款订单列表缩略图：演示库里的 /demo/*.png 在 COS 上 404，改为本仓已存在的素材 URL。
-- 幂等：仅替换已知失效的 demo 相对路径。

UPDATE qixi_crm_b_order_item
SET cover_url_snapshot = 'https://cos.qxkejiwl.top/pte-live-ecrm/platform/20260807/7dc13a394086786f9aba4a9606ad1eb2.png'
WHERE cover_url_snapshot IN (
  '/demo/product-knit-v1.png',
  '/demo/product-tea-v1.png',
  '/demo/product-fragrance-v1.png',
  '/demo/product-tumbler-v1.png',
  'demo/product-knit-v1.png'
)
   OR cover_url_snapshot LIKE '/demo/data-screen/thumbs/%';
