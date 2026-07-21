# 商户开放接口 `/openapi/`

> 对照文档。置信度：high=6 stale=0 unresolved=0。先读 [ACCURACY.md](./ACCURACY.md)。

合计 **6** 条。

## `openapi/auth`

### `POST /openapi/auth` — 获取开放接口 Token

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`openapi.Auth/auth`
- 源码：`app/controller/openapi/Auth.php` :: `auth()`
- 请求参数：
- `unique` (query/body, 必填) 唯一标识
- `expiration` (query/body, 必填) 时间戳，相差≤300秒
- `access_key` (query/body, 必填) access_key
- `signature` (query/body, 必填) 签名
- 返回：data.token, data.exp


## `openapi/order`

### `GET /openapi/order/detail/:id` — 订单详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`openapi.store.StoreOrder/detail`
- 源码：`app/controller/openapi/store/StoreOrder.php` :: `detail()`
- 请求参数：
- `unique` (query/body, 必填) 唯一标识
- `expiration` (query/body, 必填) 时间戳，相差≤300秒
- `access_key` (query/body, 必填) access_key
- `signature` (query/body, 必填) 签名
- `id` (path, 必填) 订单ID
- 返回：data 订单详情

### `GET /openapi/order/list` — 订单列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`openapi.store.StoreOrder/lst`
- 源码：`app/controller/openapi/store/StoreOrder.php` :: `lst()`
- 请求参数：
- `unique` (query/body, 必填) 唯一标识
- `expiration` (query/body, 必填) 时间戳，相差≤300秒
- `access_key` (query/body, 必填) access_key
- `signature` (query/body, 必填) 签名
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `status` (query, 可选) 0全部 1未支付 2待发货 3待收货 4待评价 5完成 6已退款 7待核销
- `date` (query, 可选) 时间筛选
- `order_sn` (query, 可选) 订单编号
- `username` (query, 可选) 用户昵称
- `order_type` (query, 可选) 0普通 1自提 2虚拟 3卡密
- `keywords` (query, 可选) 关键词
- `order_id` (query, 可选) 订单ID
- `activity_type` (query, 可选) 活动类型
- `group_order_sn` (query, 可选) 主订单编号
- `store_name` (query, 可选) 商品名模糊
- `filter_delivery` (query, 可选) 发货方式筛选
- `filter_product` (query, 可选) 1实物 2虚拟 3卡密
- `pay_type` (query, 可选) 支付方式
- 返回：data 分页订单列表


## `openapi/product`

### `POST /openapi/product/create` — 创建商品

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`openapi.store.StoreProduct/create`
- 源码：`app/controller/openapi/store/StoreProduct.php` :: `create()`
- 请求参数：
- `unique` (query/body, 必填) 唯一标识
- `expiration` (query/body, 必填) 时间戳，相差≤300秒
- `access_key` (query/body, 必填) access_key
- `signature` (query/body, 必填) 签名
- `is_gift_bag` (body, 可选) ProductRepository::CREATE_PARAMS
- `integral_rate` (body, 可选) ProductRepository::CREATE_PARAMS
- `mer_labels` (body, 可选) ProductRepository::CREATE_PARAMS
- `delivery_way` (body, 可选) ProductRepository::CREATE_PARAMS
- `delivery_free` (body, 可选) ProductRepository::CREATE_PARAMS
- `param_temp_id` (body, 可选) ProductRepository::CREATE_PARAMS
- `custom_temp_id` (body, 可选) ProductRepository::CREATE_PARAMS
- `extend` (body, 可选) ProductRepository::CREATE_PARAMS
- `mer_form_id` (body, 可选) ProductRepository::CREATE_PARAMS
- `auto_on_time` (body, 可选) ProductRepository::CREATE_PARAMS
- `auto_off_time` (body, 可选) ProductRepository::CREATE_PARAMS
- `refund_switch` (body, 可选) ProductRepository::CREATE_PARAMS
- `once_max_count` (body, 可选) ProductRepository::CREATE_PARAMS
- `once_min_count` (body, 可选) ProductRepository::CREATE_PARAMS
- `pay_limit` (body, 可选) ProductRepository::CREATE_PARAMS
- `give_coupon_ids` (body, 可选) ProductRepository::CREATE_PARAMS
- `type` (body, 可选) ProductRepository::CREATE_PARAMS
- `svip_price` (body, 可选) ProductRepository::CREATE_PARAMS
- `svip_price_type` (body, 可选) ProductRepository::CREATE_PARAMS
- `params` (body, 可选) ProductRepository::CREATE_PARAMS
- `product_type` (body, 可选) ProductRepository::CREATE_PARAMS
- `good_ids` (body, 可选) ProductRepository::CREATE_PARAMS
- `reservation_time_type` (body, 可选) ProductRepository::CREATE_PARAMS
- `reservation_start_time` (body, 可选) ProductRepository::CREATE_PARAMS
- `reservation_end_time` (body, 可选) ProductRepository::CREATE_PARAMS
- `reservation_time_interval` (body, 可选) ProductRepository::CREATE_PARAMS
- `time_period` (body, 可选) ProductRepository::CREATE_PARAMS
- `reservation_type` (body, 可选) ProductRepository::CREATE_PARAMS
- `show_num_type` (body, 可选) ProductRepository::CREATE_PARAMS
- `sale_time_type` (body, 可选) ProductRepository::CREATE_PARAMS
- `sale_time_start_day` (body, 可选) ProductRepository::CREATE_PARAMS
- `sale_time_end_day` (body, 可选) ProductRepository::CREATE_PARAMS
- `sale_time_week` (body, 可选) ProductRepository::CREATE_PARAMS
- `show_reservation_days` (body, 可选) ProductRepository::CREATE_PARAMS
- `is_advance` (body, 可选) ProductRepository::CREATE_PARAMS
- `advance_time` (body, 可选) ProductRepository::CREATE_PARAMS
- `is_cancel_reservation` (body, 可选) ProductRepository::CREATE_PARAMS
- `cancel_reservation_time` (body, 可选) ProductRepository::CREATE_PARAMS
- `reservation_form_type` (body, 可选) ProductRepository::CREATE_PARAMS
- `labels` (body, 可选) ProductRepository::CREATE_PARAMS
- `activity_label_ids` (body, 可选) ProductRepository::CREATE_PARAMS
- 返回：message 添加成功

### `GET /openapi/product/detail/:id` — 商品详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`openapi.store.StoreProduct/detail`
- 源码：`app/controller/openapi/store/StoreProduct.php` :: `detail()`
- 请求参数：
- `unique` (query/body, 必填) 唯一标识
- `expiration` (query/body, 必填) 时间戳，相差≤300秒
- `access_key` (query/body, 必填) access_key
- `signature` (query/body, 必填) 签名
- `id` (path, 必填) 商品ID
- 返回：data 商品详情

### `GET /openapi/product/list` — 商品列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`openapi.store.StoreProduct/lst`
- 源码：`app/controller/openapi/store/StoreProduct.php` :: `lst()`
- 请求参数：
- `unique` (query/body, 必填) 唯一标识
- `expiration` (query/body, 必填) 时间戳，相差≤300秒
- `access_key` (query/body, 必填) access_key
- `signature` (query/body, 必填) 签名
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `temp_id` (query/body, 可选)
- `cate_id` (query/body, 可选)
- `keyword` (query/body, 可选)
- `mer_cate_id` (query/body, 可选)
- `is_gift_bag` (query/body, 可选)
- `status` (query/body, 可选)
- `us_status` (query/body, 可选)
- `product_id` (query/body, 可选)
- `mer_labels` (query/body, 可选)
- `order` (query/body, 可选) 默认 'sort'
- `is_ficti` (query/body, 可选)
- `svip_price_type` (query/body, 可选)
- `filters_type` (query/body, 可选)
- `is_action` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：data 商品分页列表

