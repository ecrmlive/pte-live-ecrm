# 商户后台 `/mer/`

> 对照文档。置信度：high=589 stale=6 unresolved=0。先读 [ACCURACY.md](./ACCURACY.md)。

合计 **595** 条。

## `mer/activity`

### `POST /mer/activity/label/batch_create` — batchCreate

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.StoreActivityLabel/batchCreate`
- 源码：`app/controller/merchant/store/StoreActivityLabel.php` :: `batchCreate()`
- 请求参数：
- `product_ids` (query/body, 可选)
- `label_ids` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/activity/label/options` — options

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.StoreActivityLabel/options`
- 源码：`app/controller/merchant/store/StoreActivityLabel.php` :: `options()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/activity/label/select` — select

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.StoreActivityLabel/select`
- 源码：`app/controller/merchant/store/StoreActivityLabel.php` :: `select()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `mer/ajcaptcha`

### `GET /mer/ajcaptcha` — ajcaptcha

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Auth/ajcaptcha`
- 源码：`app/controller/api/Auth.php` :: `ajcaptcha()`
- 请求参数：
- `captchaType` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `mer/ajcheck`

### `POST /mer/ajcheck` — ajcheck

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Auth/ajcheck`
- 源码：`app/controller/api/Auth.php` :: `ajcheck()`
- 请求参数：
- `token` (query/body, 可选)
- `pointJson` (query/body, 可选)
- `captchaType` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `mer/ajstatus`

### `POST /mer/ajstatus` — ajCaptchaStatus

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.admin.Login/ajCaptchaStatus`
- 源码：`app/controller/merchant/system/admin/Login.php` :: `ajCaptchaStatus()`
- 请求参数：
- `account` (query/body, 可选)
- 返回：data 对象字段: status | 外层: {status,message,data}


## `mer/analytics`

### `GET /mer/analytics/order/line_chart` — 折线图统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.analytics.StoreOrder/lineChart`
- 源码：`app/controller/admin/analytics/StoreOrder.php` :: `lineChart()`
- 请求参数：
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/analytics/order/pie_chart/:type` — 折线图统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.analytics.StoreOrder/typePieCahrt`
- 源码：`app/controller/admin/analytics/StoreOrder.php` :: `typePieCahrt()`
- 请求参数：
- `type` (path, 必填) 路径参数
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/analytics/order/top` — 顶部统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.analytics.StoreOrder/top`
- 源码：`app/controller/admin/analytics/StoreOrder.php` :: `top()`
- 请求参数：
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/analytics/product/line_chart` — 折线图统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.analytics.StoreProduct/lineChart`
- 源码：`app/controller/admin/analytics/StoreProduct.php` :: `lineChart()`
- 请求参数：
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/analytics/product/pie_chart/:type` — 折线图统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.analytics.StoreProduct/typePieCahrt`
- 源码：`app/controller/admin/analytics/StoreProduct.php` :: `typePieCahrt()`
- 请求参数：
- `type` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/analytics/product/top` — 顶部统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.analytics.StoreProduct/top`
- 源码：`app/controller/admin/analytics/StoreProduct.php` :: `top()`
- 请求参数：
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `mer/applyments`

### `GET /mer/applyments/areas` — areas

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.MerchantApplyments/areas`
- 源码：`app/controller/merchant/system/MerchantApplyments.php` :: `areas()`
- 请求参数：
- `province_code` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/applyments/banks` — banks

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.MerchantApplyments/banks`
- 源码：`app/controller/merchant/system/MerchantApplyments.php` :: `banks()`
- 请求参数：
- `type` (query/body, 可选)
- `limit` (query/body, 可选)
- `offset` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/applyments/banks/branches` — branches

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.MerchantApplyments/branches`
- 源码：`app/controller/merchant/system/MerchantApplyments.php` :: `branches()`
- 请求参数：
- `bank_alias_code` (query/body, 可选)
- `city_code` (query/body, 可选)
- `limit` (query/body, 可选)
- `offset` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/applyments/check` — 查询审核结果

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.MerchantApplyments/check`
- 源码：`app/controller/merchant/system/MerchantApplyments.php` :: `check()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/applyments/create` — 申请

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.MerchantApplyments/create`
- 源码：`app/controller/merchant/system/MerchantApplyments.php` :: `create()`
- 请求参数：
- `type` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/applyments/detail` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.MerchantApplyments/detail`
- 源码：`app/controller/merchant/system/MerchantApplyments.php` :: `detail()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/applyments/index` — index

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.MerchantApplyments/index`
- 源码：`app/controller/merchant/system/MerchantApplyments.php` :: `index()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/applyments/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.MerchantApplyments/update`
- 源码：`app/controller/merchant/system/MerchantApplyments.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `type` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/applyments/upload/:field` — 上传图片

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.MerchantApplyments/uploadImage`
- 源码：`app/controller/merchant/system/MerchantApplyments.php` :: `uploadImage()`
- 请求参数：
- `field` (path, 必填) 路径参数
- `water` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/applyments/upload/video/:field` — 上传视频

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.MerchantApplyments/uploadVideo`
- 源码：`app/controller/merchant/system/MerchantApplyments.php` :: `uploadVideo()`
- 请求参数：
- `field` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `mer/auto_label`

### `POST /mer/auto_label/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.user.LabelRule/create`
- 源码：`app/controller/merchant/user/LabelRule.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /mer/auto_label/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.user.LabelRule/delete`
- 源码：`app/controller/merchant/user/LabelRule.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/auto_label/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.user.LabelRule/getList`
- 源码：`app/controller/merchant/user/LabelRule.php` :: `getList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /mer/auto_label/sync/:id` — 自动同步

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.user.LabelRule/sync`
- 源码：`app/controller/merchant/user/LabelRule.php` :: `sync()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/auto_label/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.user.LabelRule/update`
- 源码：`app/controller/merchant/user/LabelRule.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `mer/broadcast`

### `POST /mer/broadcast/assistant/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.broadcast.BroadcastAssistant/create`
- 源码：`app/controller/merchant/store/broadcast/BroadcastAssistant.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/broadcast/assistant/create/form` — 添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.broadcast.BroadcastAssistant/createForm`
- 源码：`app/controller/merchant/store/broadcast/BroadcastAssistant.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /mer/broadcast/assistant/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.broadcast.BroadcastAssistant/delete`
- 源码：`app/controller/merchant/store/broadcast/BroadcastAssistant.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/broadcast/assistant/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.broadcast.BroadcastAssistant/lst`
- 源码：`app/controller/merchant/store/broadcast/BroadcastAssistant.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `username` (query/body, 可选)
- `nickname` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /mer/broadcast/assistant/mark/:id` — 备注

- 置信度：⛔ stale
- 说明：路由已登记，但 `app/controller/merchant/store/broadcast/BroadcastAssistant.php` 中不存在方法 `mark`（原项目死路由/加密扩展/版本差异）。开发勿实现为有效接口，除非核实。
- 处理器：`merchant.store.broadcast.BroadcastAssistant/mark`
- 源码：`app/controller/merchant/store/broadcast/BroadcastAssistant.php` :: `mark()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一响应包 {status,message,data?}；详见控制器实现 | 外层: {status,message,data}

### `POST /mer/broadcast/assistant/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.broadcast.BroadcastAssistant/update`
- 源码：`app/controller/merchant/store/broadcast/BroadcastAssistant.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/broadcast/assistant/update/:id/form` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.broadcast.BroadcastAssistant/updateForm`
- 源码：`app/controller/merchant/store/broadcast/BroadcastAssistant.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/broadcast/goods/batch_create` — 批量添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.broadcast.BroadcastGoods/batchCreate`
- 源码：`app/controller/merchant/store/broadcast/BroadcastGoods.php` :: `batchCreate()`
- 请求参数：
- `goods` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/broadcast/goods/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.broadcast.BroadcastGoods/create`
- 源码：`app/controller/merchant/store/broadcast/BroadcastGoods.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/broadcast/goods/create/form` — 添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.broadcast.BroadcastGoods/createForm`
- 源码：`app/controller/merchant/store/broadcast/BroadcastGoods.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /mer/broadcast/goods/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.broadcast.BroadcastGoods/delete`
- 源码：`app/controller/merchant/store/broadcast/BroadcastGoods.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/broadcast/goods/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.broadcast.BroadcastGoods/detail`
- 源码：`app/controller/merchant/store/broadcast/BroadcastGoods.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/broadcast/goods/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.broadcast.BroadcastGoods/lst`
- 源码：`app/controller/merchant/store/broadcast/BroadcastGoods.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `status_tag` (query/body, 可选)
- `keyword` (query/body, 可选)
- `mer_valid` (query/body, 可选)
- `broadcast_goods_id` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /mer/broadcast/goods/mark/:id` — 备注

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.broadcast.BroadcastGoods/mark`
- 源码：`app/controller/merchant/store/broadcast/BroadcastGoods.php` :: `mark()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `mark` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/broadcast/goods/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.broadcast.BroadcastGoods/changeStatus`
- 源码：`app/controller/merchant/store/broadcast/BroadcastGoods.php` :: `changeStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `is_show` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/broadcast/goods/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.broadcast.BroadcastGoods/update`
- 源码：`app/controller/merchant/store/broadcast/BroadcastGoods.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/broadcast/goods/update/form/:id` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.broadcast.BroadcastGoods/updateForm`
- 源码：`app/controller/merchant/store/broadcast/BroadcastGoods.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/broadcast/room/addassistant/:id` — 添加 客服

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.broadcast.BroadcastRoom/addAssistant`
- 源码：`app/controller/merchant/store/broadcast/BroadcastRoom.php` :: `addAssistant()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `assistant_id` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/broadcast/room/addassistant/form/:id` — 添加客服表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.broadcast.BroadcastRoom/addAssistantForm`
- 源码：`app/controller/merchant/store/broadcast/BroadcastRoom.php` :: `addAssistantForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/broadcast/room/closeKf/:id` — 关闭客服

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.broadcast.BroadcastRoom/closeKf`
- 源码：`app/controller/merchant/store/broadcast/BroadcastRoom.php` :: `closeKf()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/broadcast/room/comment/:id` — 禁言

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.broadcast.BroadcastRoom/banComment`
- 源码：`app/controller/merchant/store/broadcast/BroadcastRoom.php` :: `banComment()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/broadcast/room/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.broadcast.BroadcastRoom/create`
- 源码：`app/controller/merchant/store/broadcast/BroadcastRoom.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/broadcast/room/create/form` — 添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.broadcast.BroadcastRoom/createForm`
- 源码：`app/controller/merchant/store/broadcast/BroadcastRoom.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /mer/broadcast/room/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.broadcast.BroadcastRoom/delete`
- 源码：`app/controller/merchant/store/broadcast/BroadcastRoom.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/broadcast/room/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.broadcast.BroadcastRoom/detail`
- 源码：`app/controller/merchant/store/broadcast/BroadcastRoom.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/broadcast/room/export_goods` — 导入商品

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.broadcast.BroadcastRoom/exportGoods`
- 源码：`app/controller/merchant/store/broadcast/BroadcastRoom.php` :: `exportGoods()`
- 请求参数：
- `ids` (query/body, 可选)
- `room_id` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/broadcast/room/feedsPublic/:id` — 收录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.broadcast.BroadcastRoom/isFeedsPublic`
- 源码：`app/controller/merchant/store/broadcast/BroadcastRoom.php` :: `isFeedsPublic()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/broadcast/room/goods/:id` — 商品详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.broadcast.BroadcastRoom/goodsList`
- 源码：`app/controller/merchant/store/broadcast/BroadcastRoom.php` :: `goodsList()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/broadcast/room/lst` — 列表 

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.broadcast.BroadcastRoom/lst`
- 源码：`app/controller/merchant/store/broadcast/BroadcastRoom.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `status_tag` (query/body, 可选)
- `show_tag` (query/body, 可选)
- `show_type` (query/body, 可选)
- `live_status` (query/body, 可选)
- `broadcast_room_id` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /mer/broadcast/room/mark/:id` — 备注

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.broadcast.BroadcastRoom/mark`
- 源码：`app/controller/merchant/store/broadcast/BroadcastRoom.php` :: `mark()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `mark` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/broadcast/room/on_sale/:id` — 商品上下架

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.broadcast.BroadcastRoom/onSale`
- 源码：`app/controller/merchant/store/broadcast/BroadcastRoom.php` :: `onSale()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- `goods_id` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/broadcast/room/push_message/:id` — 消息推送

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.broadcast.BroadcastRoom/pushMessage`
- 源码：`app/controller/merchant/store/broadcast/BroadcastRoom.php` :: `pushMessage()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/broadcast/room/rm_goods` — 删除商品

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.broadcast.BroadcastRoom/rmExportGoods`
- 源码：`app/controller/merchant/store/broadcast/BroadcastRoom.php` :: `rmExportGoods()`
- 请求参数：
- `id` (query/body, 可选)
- `room_id` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/broadcast/room/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.broadcast.BroadcastRoom/changeStatus`
- 源码：`app/controller/merchant/store/broadcast/BroadcastRoom.php` :: `changeStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `is_show` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/broadcast/room/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.broadcast.BroadcastRoom/update`
- 源码：`app/controller/merchant/store/broadcast/BroadcastRoom.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/broadcast/room/update/form/:id` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.broadcast.BroadcastRoom/updateForm`
- 源码：`app/controller/merchant/store/broadcast/BroadcastRoom.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `mer/captcha`

### `GET /mer/captcha` — getCaptcha

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.admin.Login/getCaptcha`
- 源码：`app/controller/merchant/system/admin/Login.php` :: `getCaptcha()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `mer/change`

### `GET /mer/change/color` — 一键换色

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.Common/getChangeColor`
- 源码：`app/controller/admin/Common.php` :: `getChangeColor()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `mer/community`

### `GET /mer/community/cate/lst` — 分类列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.content.Community/cateLst`
- 源码：`app/controller/merchant/store/content/Community.php` :: `cateLst()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/community/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.content.Community/create`
- 源码：`app/controller/merchant/store/content/Community.php` :: `create()`
- 请求参数：
- `uid` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /mer/community/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.content.Community/delete`
- 源码：`app/controller/merchant/store/content/Community.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/community/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.content.Community/detail`
- 源码：`app/controller/merchant/store/content/Community.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/community/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.content.Community/lst`
- 源码：`app/controller/merchant/store/content/Community.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选) 默认 ''
- `topic_id` (query/body, 可选) 默认 ''
- `category_id` (query/body, 可选) 默认 ''
- `status` (query/body, 可选) 默认 ''
- `is_type` (query/body, 可选) 默认 ''
- `search_type` (query/body, 可选) 默认 'content'
- `is_del` (query/body, 可选) 默认 0
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /mer/community/reply/:id` — 评论

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.content.Community/reply`
- 源码：`app/controller/merchant/store/content/Community.php` :: `reply()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/community/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.content.Community/update`
- 源码：`app/controller/merchant/store/content/Community.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `mer/config`

### `GET /mer/config` — config

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.Common/config`
- 源码：`app/controller/merchant/Common.php` :: `config()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/config/:key` — 配置获取

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.config.Config/form`
- 源码：`app/controller/admin/system/config/Config.php` :: `form()`
- 请求参数：
- `key` (path, 必填) 路径参数
- `tab_id` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/config/:type` — saveConfig

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.Common/saveConfig`
- 源码：`app/controller/merchant/Common.php` :: `saveConfig()`
- 请求参数：
- `type` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/config/others/group_buying` — 拼团配置

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.config.ConfigOthers/getGroupBuying`
- 源码：`app/controller/admin/system/config/ConfigOthers.php` :: `getGroupBuying()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/config/save/:key` — 配置保存

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.config.ConfigValue/save`
- 源码：`app/controller/admin/system/config/ConfigValue.php` :: `save()`
- 请求参数：
- `key` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `mer/delivery`

### `GET /mer/delivery/config/settings` — 配送设置信息

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.delivery.DeliveryConfig/configuration`
- 源码：`app/controller/merchant/store/delivery/DeliveryConfig.php` :: `configuration()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/delivery/config/update/:id` — 更新配送设置信息

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.delivery.DeliveryConfig/update`
- 源码：`app/controller/merchant/store/delivery/DeliveryConfig.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `mer_delivery_type` (query/body, 可选)
- `mer_delivery_order_status` (query/body, 可选) 默认 0
- `dada_app_key` (query/body, 可选)
- `dada_app_sercret` (query/body, 可选)
- `dada_source_id` (query/body, 可选)
- `uupt_appkey` (query/body, 可选)
- `uupt_app_id` (query/body, 可选)
- `uupt_open_id` (query/body, 可选)
- `min_delivery_amount` (query/body, 可选) 默认 0
- `base_shipping_fee` (query/body, 可选) 默认 0
- `free_shipping_amount` (query/body, 可选) 默认 0
- `is_premium_stack_enabled` (query/body, 可选)
- `distance_premium_config` (query/body, 可选) 默认 [
- `weight_premium_config` (query/body, 可选) 默认 [
- `delivery_time_type` (query/body, 可选) 默认 1
- `selectable_days` (query/body, 可选) 默认 7
- `delivery_prompt` (query/body, 可选) 默认 ''
- `commission_rate` (query/body, 可选) 默认 0
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/delivery/order/cancel/:id` — 取消

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.delivery.DeliveryOrder/cancel`
- 源码：`app/controller/merchant/store/delivery/DeliveryOrder.php` :: `cancel()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `reason` (query/body, 可选)
- `cancel_reason` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/delivery/order/cancel/:id/form` — 取消表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.delivery.DeliveryOrder/cancelForm`
- 源码：`app/controller/merchant/store/delivery/DeliveryOrder.php` :: `cancelForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/delivery/order/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.delivery.DeliveryOrder/detail`
- 源码：`app/controller/merchant/store/delivery/DeliveryOrder.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/delivery/order/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.delivery.DeliveryOrder/lst`
- 源码：`app/controller/merchant/store/delivery/DeliveryOrder.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `station_id` (query/body, 可选)
- `status` (query/body, 可选)
- `date` (query/body, 可选)
- `order_sn` (query/body, 可选)
- `station_type` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/delivery/service/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.delivery.DeliveryService/create`
- 源码：`app/controller/merchant/store/delivery/DeliveryService.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/delivery/service/create/form` — 添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.delivery.DeliveryService/createForm`
- 源码：`app/controller/merchant/store/delivery/DeliveryService.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /mer/delivery/service/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.delivery.DeliveryService/delete`
- 源码：`app/controller/merchant/store/delivery/DeliveryService.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/delivery/service/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.delivery.DeliveryService/lst`
- 源码：`app/controller/merchant/store/delivery/DeliveryService.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `name` (query/body, 可选)
- `status` (query/body, 可选)
- `type` (query/body, 可选)
- `date` (query/body, 可选)
- `uid` (query/body, 可选)
- `phone` (query/body, 可选)
- `nickname` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/delivery/service/options` — options

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.delivery.DeliveryService/options`
- 源码：`app/controller/merchant/store/delivery/DeliveryService.php` :: `options()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/delivery/service/statisticsDetail/:id` — 统计详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.delivery.DeliveryService/statisticsDetail`
- 源码：`app/controller/merchant/store/delivery/DeliveryService.php` :: `statisticsDetail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/delivery/service/statisticsList` — 统计列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.delivery.DeliveryService/statisticsList`
- 源码：`app/controller/merchant/store/delivery/DeliveryService.php` :: `statisticsList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `name` (query/body, 可选)
- `uid` (query/body, 可选)
- `phone` (query/body, 可选)
- `delivery_date` (query/body, 可选) 默认 [
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/delivery/service/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.delivery.DeliveryService/switchWithStatus`
- 源码：`app/controller/merchant/store/delivery/DeliveryService.php` :: `switchWithStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/delivery/service/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.delivery.DeliveryService/update`
- 源码：`app/controller/merchant/store/delivery/DeliveryService.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/delivery/service/update/:id/form` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.delivery.DeliveryService/updateForm`
- 源码：`app/controller/merchant/store/delivery/DeliveryService.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/delivery/station/business` — 获取分类

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.delivery.DeliveryStation/getBusiness`
- 源码：`app/controller/merchant/store/delivery/DeliveryStation.php` :: `getBusiness()`
- 请求参数：
- `type` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/delivery/station/code` — 充值二维码

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.delivery.DeliveryStation/getQrcode`
- 源码：`app/controller/merchant/store/delivery/DeliveryStation.php` :: `getQrcode()`
- 请求参数：
- `pay_type` (query/body, 可选)
- `price` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/delivery/station/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.delivery.DeliveryStation/create`
- 源码：`app/controller/merchant/store/delivery/DeliveryStation.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /mer/delivery/station/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.delivery.DeliveryStation/delete`
- 源码：`app/controller/merchant/store/delivery/DeliveryStation.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/delivery/station/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.delivery.DeliveryStation/detail`
- 源码：`app/controller/merchant/store/delivery/DeliveryStation.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/delivery/station/getCity` — 城市列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.delivery.DeliveryStation/getCityLst`
- 源码：`app/controller/merchant/store/delivery/DeliveryStation.php` :: `getCityLst()`
- 请求参数：
- `type` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/delivery/station/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.delivery.DeliveryStation/lst`
- 源码：`app/controller/merchant/store/delivery/DeliveryStation.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `station_name` (query/body, 可选)
- `contact_name` (query/body, 可选)
- `phone` (query/body, 可选)
- `station_address` (query/body, 可选)
- `status` (query/body, 可选)
- `swtich_type` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/delivery/station/mark/:id` — 备注

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.delivery.DeliveryStation/mark`
- 源码：`app/controller/merchant/store/delivery/DeliveryStation.php` :: `mark()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `mark` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/delivery/station/mark/:id/form` — 备注表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.delivery.DeliveryStation/markForm`
- 源码：`app/controller/merchant/store/delivery/DeliveryStation.php` :: `markForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/delivery/station/options` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.delivery.DeliveryStation/options`
- 源码：`app/controller/merchant/store/delivery/DeliveryStation.php` :: `options()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/delivery/station/payLst` — 充值记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.delivery.DeliveryStation/payLst`
- 源码：`app/controller/merchant/store/delivery/DeliveryStation.php` :: `payLst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/delivery/station/select` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.delivery.DeliveryStation/select`
- 源码：`app/controller/merchant/store/delivery/DeliveryStation.php` :: `select()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/delivery/station/status/:id` — 编辑状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.delivery.DeliveryStation/switchWithStatus`
- 源码：`app/controller/merchant/store/delivery/DeliveryStation.php` :: `switchWithStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/delivery/station/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.delivery.DeliveryStation/update`
- 源码：`app/controller/merchant/store/delivery/DeliveryStation.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `mer/discounts`

### `POST /mer/discounts/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Discounts/create`
- 源码：`app/controller/merchant/store/product/Discounts.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /mer/discounts/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Discounts/delete`
- 源码：`app/controller/merchant/store/product/Discounts.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/discounts/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Discounts/detail`
- 源码：`app/controller/merchant/store/product/Discounts.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/discounts/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Discounts/lst`
- 源码：`app/controller/merchant/store/product/Discounts.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `store_name` (query/body, 可选)
- `title` (query/body, 可选)
- `type` (query/body, 可选)
- `status` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/discounts/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Discounts/switchStatus`
- 源码：`app/controller/merchant/store/product/Discounts.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/discounts/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Discounts/update`
- 源码：`app/controller/merchant/store/product/Discounts.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `mer/diy`

### `GET /mer/diy/categroy/options` — 列表 

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.PageCategroy/options`
- 源码：`app/controller/admin/system/diy/PageCategroy.php` :: `options()`
- 请求参数：
- `type` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/diy/copy/:id` — 复制

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.diy.Diy/copy`
- 源码：`app/controller/merchant/system/diy/Diy.php` :: `copy()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/diy/create/:id` — 添加/编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.diy.Diy/saveData`
- 源码：`app/controller/merchant/system/diy/Diy.php` :: `saveData()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `name` (query/body, 可选) 默认 ''
- `title` (query/body, 可选) 默认 ''
- `type` (query/body, 可选) 默认 '2'
- `cover_image` (query/body, 可选) 默认 ''
- `is_show` (query/body, 可选) 默认 0
- `is_bg_color` (query/body, 可选) 默认 0
- `is_bg_pic` (query/body, 可选) 默认 0
- `bg_tab_val` (query/body, 可选) 默认 0
- `color_picker` (query/body, 可选) 默认 ''
- `bg_pic` (query/body, 可选) 默认 ''
- `is_diy` (query/body, 可选) 默认 1
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/diy/default_lst` — 默认模板列表 

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.diy.Diy/defaultLst`
- 源码：`app/controller/merchant/system/diy/Diy.php` :: `defaultLst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `DELETE /mer/diy/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.diy.Diy/del`
- 源码：`app/controller/merchant/system/diy/Diy.php` :: `del()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/diy/detail/:id` — 详情 

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.diy.Diy/getInfo`
- 源码：`app/controller/merchant/system/diy/Diy.php` :: `getInfo()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/diy/link/getLinks/:id` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.PageLink/getLinks`
- 源码：`app/controller/admin/system/diy/PageLink.php` :: `getLinks()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `type` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/diy/link/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.PageLink/lst`
- 源码：`app/controller/admin/system/diy/PageLink.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `status` (query/body, 可选) 默认 1
- `type` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /mer/diy/lst` — 列表 

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.diy.Diy/lst`
- 源码：`app/controller/merchant/system/diy/Diy.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `status` (query/body, 可选) 默认 ''
- `name` (query/body, 可选) 默认 ''
- `version` (query/body, 可选) 默认 ''
- `is_diy` (query/body, 可选) 默认 1
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/diy/product/lst` — 店铺街装修

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.diy.Diy/productLst`
- 源码：`app/controller/merchant/system/diy/Diy.php` :: `productLst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `store_name` (query/body, 可选) 默认 ''
- `order` (query/body, 可选) 默认 'star'
- `cate_pid` (query/body, 可选) 默认 0
- `star` (query/body, 可选) 默认 ''
- `product_type` (query/body, 可选) 默认 0
- `mer_cate_id` (query/body, 可选) 默认 ''
- `cate_id` (query/body, 可选) 默认 ''
- `product_ids` (query/body, 可选)
- `store_type_id` (query/body, 可选)
- `mer_store_label_id` (query/body, 可选)
- `delivery_type` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/diy/productCategory/create/:id` — 保存商品分类

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.diy.Diy/saveData`
- 源码：`app/controller/merchant/system/diy/Diy.php` :: `saveData()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `name` (query/body, 可选) 默认 ''
- `title` (query/body, 可选) 默认 ''
- `type` (query/body, 可选) 默认 '2'
- `cover_image` (query/body, 可选) 默认 ''
- `is_show` (query/body, 可选) 默认 0
- `is_bg_color` (query/body, 可选) 默认 0
- `is_bg_pic` (query/body, 可选) 默认 0
- `bg_tab_val` (query/body, 可选) 默认 0
- `color_picker` (query/body, 可选) 默认 ''
- `bg_pic` (query/body, 可选) 默认 ''
- `is_diy` (query/body, 可选) 默认 1
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/diy/productCategory/info` — 商品分类信息

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.diy.Diy/productCategoryInfo`
- 源码：`app/controller/merchant/system/diy/Diy.php` :: `productCategoryInfo()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/diy/recovery/:id` — 重置模板

- 置信度：✅ high
- 说明：已按源码方法名校正
- 处理器：`merchant.system.diy.Diy/Recovery`
- 源码：`app/controller/merchant/system/diy/Diy.php` :: `Recovery()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一响应包 {status,message,data?}；详见控制器实现 | 外层: {status,message,data}

### `GET /mer/diy/review/:id` — review

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.diy.Diy/review`
- 源码：`app/controller/merchant/system/diy/Diy.php` :: `review()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/diy/set_default_data/:id` — 使用模板

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.diy.Diy/setDefaultData`
- 源码：`app/controller/merchant/system/diy/Diy.php` :: `setDefaultData()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/diy/show` — 当前使用模板

- 置信度：✅ high
- 说明：已按源码方法名校正
- 处理器：`merchant.system.diy.Diy/getInfo`
- 源码：`app/controller/merchant/system/diy/Diy.php` :: `getInfo()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一响应包 {status,message,data?}；详见控制器实现 | 外层: {status,message,data}

### `POST /mer/diy/status/:id` — 使用模板

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.diy.Diy/setStatus`
- 源码：`app/controller/merchant/system/diy/Diy.php` :: `setStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `mer/dump_lst`

### `GET /mer/dump_lst` — 使用记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.serve.Export/dumpLst`
- 源码：`app/controller/admin/system/serve/Export.php` :: `dumpLst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `date` (query/body, 可选)
- `mer_id` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `mer/excel`

### `GET /mer/excel/download/:type` — downloadExpress

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.Excel/downloadExpress`
- 源码：`app/controller/merchant/store/Excel.php` :: `downloadExpress()`
- 请求参数：
- `type` (path, 必填) 路径参数
- 返回：失败时 status=400, message 为错误信息 | 外层: {status,message,data}


## `mer/expr`

### `POST /mer/expr/changeMerStatus/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Express/merStatus`
- 源码：`app/controller/admin/store/Express.php` :: `merStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `mer_status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/expr/dump_lst` — 默认模板

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.serve.Export/dumpLst`
- 源码：`app/controller/admin/system/serve/Export.php` :: `dumpLst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `date` (query/body, 可选)
- `mer_id` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/expr/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Express/lst`
- 源码：`app/controller/admin/store/Express.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `code` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /mer/expr/options` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Express/options`
- 源码：`app/controller/admin/store/Express.php` :: `options()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/expr/partner/:id` — 月结账号编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Express/partner`
- 源码：`app/controller/admin/store/Express.php` :: `partner()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `account` (query/body, 可选)
- `key` (query/body, 可选)
- `net_name` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/expr/partner/:id/form` — 月结账号编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Express/partnerForm`
- 源码：`app/controller/admin/store/Express.php` :: `partnerForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/expr/temps` — 预览

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.serve.Export/getExportTemp`
- 源码：`app/controller/admin/system/serve/Export.php` :: `getExportTemp()`
- 请求参数：
- `com` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `mer/financial`

### `POST /mer/financial/account` — 收款方式

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.financial.Financial/accountSave`
- 源码：`app/controller/merchant/system/financial/Financial.php` :: `accountSave()`
- 请求参数：
- `account` (query/body, 可选)
- `financial_type` (query/body, 可选)
- `name` (query/body, 可选)
- `bank` (query/body, 可选)
- `bank_code` (query/body, 可选)
- `wechat` (query/body, 可选)
- `wechat_code` (query/body, 可选)
- `alipay` (query/body, 可选)
- `alipay_code` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/financial/account/form` — 收款方式表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.financial.Financial/accountForm`
- 源码：`app/controller/merchant/system/financial/Financial.php` :: `accountForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/financial/create` — 申请

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.financial.Financial/createSave`
- 源码：`app/controller/merchant/system/financial/Financial.php` :: `createSave()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/financial/create/form` — 申请表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.financial.Financial/createForm`
- 源码：`app/controller/merchant/system/financial/Financial.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `DELETE /mer/financial/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.financial.Financial/delete`
- 源码：`app/controller/merchant/system/financial/Financial.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/financial/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.financial.Financial/detail`
- 源码：`app/controller/merchant/system/financial/Financial.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/financial/export` — 导出

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.financial.Financial/export`
- 源码：`app/controller/merchant/system/financial/Financial.php` :: `export()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `date` (query/body, 可选)
- `status` (query/body, 可选)
- `financial_type` (query/body, 可选)
- `financial_status` (query/body, 可选)
- `keyword` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/financial/lst` — 转账记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.financial.Financial/lst`
- 源码：`app/controller/merchant/system/financial/Financial.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `date` (query/body, 可选)
- `status` (query/body, 可选)
- `financial_type` (query/body, 可选)
- `financial_status` (query/body, 可选)
- `keyword` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/financial/mark/:id` — 备注

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.financial.Financial/mark`
- 源码：`app/controller/merchant/system/financial/Financial.php` :: `mark()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `mark` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/financial/mark/:id/form` — 备注表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.financial.Financial/markForm`
- 源码：`app/controller/merchant/system/financial/Financial.php` :: `markForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/financial/refund/margin` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.financial.Financial/refundMargin`
- 源码：`app/controller/merchant/system/financial/Financial.php` :: `refundMargin()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/financial/refund/margin_apply` — 退保证金申请

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.financial.Financial/refundMarginApply`
- 源码：`app/controller/merchant/system/financial/Financial.php` :: `refundMarginApply()`
- 请求参数：
- `type` (query/body, 可选)
- `name` (query/body, 可选)
- `code` (query/body, 可选)
- `pic` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。


## `mer/financial_record`

### `GET /mer/financial_record/count` — 统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.FinancialRecord/title`
- 源码：`app/controller/admin/system/merchant/FinancialRecord.php` :: `title()`
- 请求参数：
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/financial_record/detail/:type` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.FinancialRecord/detail`
- 源码：`app/controller/admin/system/merchant/FinancialRecord.php` :: `detail()`
- 请求参数：
- `type` (path, 必填) 路径参数
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/financial_record/detail_export/:type` — 导出

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.FinancialRecord/exportDetail`
- 源码：`app/controller/admin/system/merchant/FinancialRecord.php` :: `exportDetail()`
- 请求参数：
- `type` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/financial_record/export` — 导出

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.FinancialRecord/export`
- 源码：`app/controller/admin/system/merchant/FinancialRecord.php` :: `export()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `date` (query/body, 可选)
- `mer_id` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/financial_record/list` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.FinancialRecord/lst`
- 源码：`app/controller/admin/system/merchant/FinancialRecord.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `date` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `uid` (query/body, 可选)
- `real_name` (query/body, 可选)
- `nickname` (query/body, 可选)
- `order_sn` (query/body, 可选)
- `pay_type` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /mer/financial_record/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.FinancialRecord/getList`
- 源码：`app/controller/admin/system/merchant/FinancialRecord.php` :: `getList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `type` (query/body, 可选) 默认 1
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/financial_record/title` — 统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.FinancialRecord/getTitle`
- 源码：`app/controller/admin/system/merchant/FinancialRecord.php` :: `getTitle()`
- 请求参数：
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `mer/getMerchantMenusList`

### `POST /mer/getMerchantMenusList` — getMenusList

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.auth.Menu/getMenusList`
- 源码：`app/controller/admin/system/auth/Menu.php` :: `getMenusList()`
- 请求参数：
- `is_mer` (query/body, 可选)
- `keyword` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `mer/group`

### `POST /mer/group/data/create/:groupId` — 数据添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.groupData.GroupData/create`
- 源码：`app/controller/admin/system/groupData/GroupData.php` :: `create()`
- 请求参数：
- `groupId` (path, 必填) 路径参数
- `sort` (query/body, 可选) 默认 0
- `status` (query/body, 可选) 默认 0
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/group/data/create/table/:groupId` — 数据添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.groupData.GroupData/createTable`
- 源码：`app/controller/admin/system/groupData/GroupData.php` :: `createTable()`
- 请求参数：
- `groupId` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `DELETE /mer/group/data/delete/:id` — 数据删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.groupData.GroupData/delete`
- 源码：`app/controller/admin/system/groupData/GroupData.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/group/data/detail/:id` — baseDetail

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.groupData.GroupData/baseDetail`
- 源码：`app/controller/admin/system/groupData/GroupData.php` :: `baseDetail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/group/data/lst/:groupId` — 数据列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.groupData.GroupData/lst`
- 源码：`app/controller/admin/system/groupData/GroupData.php` :: `lst()`
- 请求参数：
- `groupId` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/group/data/status/:id` — 数据修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.groupData.GroupData/changeStatus`
- 源码：`app/controller/admin/system/groupData/GroupData.php` :: `changeStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/group/data/update/:groupId/:id` — 数据编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.groupData.GroupData/update`
- 源码：`app/controller/admin/system/groupData/GroupData.php` :: `update()`
- 请求参数：
- `groupId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- `sort` (query/body, 可选) 默认 0
- `status` (query/body, 可选) 默认 0
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/group/data/update/table/:groupId/:id` — 数据编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.groupData.GroupData/updateTable`
- 源码：`app/controller/admin/system/groupData/GroupData.php` :: `updateTable()`
- 请求参数：
- `groupId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/group/detail/:id` — 数据详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.groupData.Group/get`
- 源码：`app/controller/admin/system/groupData/Group.php` :: `get()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `mer/guarantee`

### `POST /mer/guarantee/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.guarantee.GuaranteeTemplate/create`
- 源码：`app/controller/merchant/store/guarantee/GuaranteeTemplate.php` :: `create()`
- 请求参数：
- `template_name` (query/body, 可选)
- `template_value` (query/body, 可选)
- `status` (query/body, 可选) 默认 1
- `sort` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /mer/guarantee/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.guarantee.GuaranteeTemplate/delete`
- 源码：`app/controller/merchant/store/guarantee/GuaranteeTemplate.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/guarantee/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.guarantee.GuaranteeTemplate/detail`
- 源码：`app/controller/merchant/store/guarantee/GuaranteeTemplate.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/guarantee/list` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.guarantee.GuaranteeTemplate/list`
- 源码：`app/controller/merchant/store/guarantee/GuaranteeTemplate.php` :: `list()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/guarantee/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.guarantee.GuaranteeTemplate/lst`
- 源码：`app/controller/merchant/store/guarantee/GuaranteeTemplate.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `date` (query/body, 可选)
- `keyword` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/guarantee/select` — 筛选

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.guarantee.GuaranteeTemplate/select`
- 源码：`app/controller/merchant/store/guarantee/GuaranteeTemplate.php` :: `select()`
- 请求参数：
- `keyword` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/guarantee/sort/:id` — 排序

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.guarantee.GuaranteeTemplate/sort`
- 源码：`app/controller/merchant/store/guarantee/GuaranteeTemplate.php` :: `sort()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `sort` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/guarantee/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.guarantee.GuaranteeTemplate/switchStatus`
- 源码：`app/controller/merchant/store/guarantee/GuaranteeTemplate.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/guarantee/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.guarantee.GuaranteeTemplate/update`
- 源码：`app/controller/merchant/store/guarantee/GuaranteeTemplate.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `template_name` (query/body, 可选)
- `template_value` (query/body, 可选)
- `status` (query/body, 可选) 默认 1
- `sort` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `mer/info`

### `GET /mer/info` — info

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.Merchant/info`
- 源码：`app/controller/merchant/system/Merchant.php` :: `info()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/info/update` — 资料更新

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.Merchant/update`
- 源码：`app/controller/merchant/system/Merchant.php` :: `update()`
- 请求参数：
- `mer_info` (query/body, 可选)
- `mer_certificate` (query/body, 可选)
- `service_phone` (query/body, 可选)
- `mer_avatar` (query/body, 可选)
- `mer_banner` (query/body, 可选)
- `mer_state` (query/body, 可选)
- `mini_banner` (query/body, 可选)
- `mer_keyword` (query/body, 可选)
- `mer_address` (query/body, 可选)
- `long` (query/body, 可选)
- `lat` (query/body, 可选)
- `delivery_way` (query/body, 可选) 默认 [2
- `mer_take_day` (query/body, 可选)
- `mer_take_time` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `mer/integral`

### `GET /mer/integral/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.user.UserIntegral/getList`
- 源码：`app/controller/merchant/user/UserIntegral.php` :: `getList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `date` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /mer/integral/title` — 统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.user.UserIntegral/getTitle`
- 源码：`app/controller/merchant/user/UserIntegral.php` :: `getTitle()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `mer/login`

### `POST /mer/login` — login

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.admin.Login/login`
- 源码：`app/controller/merchant/system/admin/Login.php` :: `login()`
- 请求参数：
- `account` (query/body, 可选)
- `password` (query/body, 可选)
- `code` (query/body, 可选)
- `key` (query/body, 可选)
- `captchaType` (query/body, 可选) 默认 ''
- `captchaVerification` (query/body, 可选) 默认 ''
- `token` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `mer/login_config`

### `GET /mer/login_config` — loginConfig

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.Common/loginConfig`
- 源码：`app/controller/admin/Common.php` :: `loginConfig()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `mer/logout`

### `GET /mer/logout` — logout

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.admin.Login/logout`
- 源码：`app/controller/merchant/system/admin/Login.php` :: `logout()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `mer/manager`

### `GET /mer/manager/user/lst` — managerUserLst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.user.UserMerchant/managerUserLst`
- 源码：`app/controller/merchant/user/UserMerchant.php` :: `managerUserLst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}


## `mer/margin`

### `GET /mer/margin/code` — getMarginQrCode

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.Merchant/getMarginQrCode`
- 源码：`app/controller/merchant/system/Merchant.php` :: `getMarginQrCode()`
- 请求参数：
- `type` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/margin/lst` — getMarginLst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.Merchant/getMarginLst`
- 源码：`app/controller/merchant/system/Merchant.php` :: `getMarginLst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/margin/make_code` — getMarginQrCode

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.Merchant/getMarginQrCode`
- 源码：`app/controller/merchant/system/Merchant.php` :: `getMarginQrCode()`
- 请求参数：
- `type` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `mer/menus`

### `GET /mer/menus` — merchantMenus

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.auth.Menu/merchantMenus`
- 源码：`app/controller/admin/system/auth/Menu.php` :: `merchantMenus()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/menus` — merchantMenus

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.auth.Menu/merchantMenus`
- 源码：`app/controller/admin/system/auth/Menu.php` :: `merchantMenus()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `mer/notice`

### `DELETE /mer/notice/del/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.notice.SystemNoticeLog/del`
- 源码：`app/controller/merchant/system/notice/SystemNoticeLog.php` :: `del()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/notice/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.notice.SystemNoticeLog/detail`
- 源码：`app/controller/merchant/system/notice/SystemNoticeLog.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/notice/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.notice.SystemNoticeLog/lst`
- 源码：`app/controller/merchant/system/notice/SystemNoticeLog.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `is_read` (query/body, 可选)
- `date` (query/body, 可选)
- `keyword` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /mer/notice/read/:id` — 已读

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.notice.SystemNoticeLog/read`
- 源码：`app/controller/merchant/system/notice/SystemNoticeLog.php` :: `read()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/notice/unread_count` — 未读统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.notice.SystemNoticeLog/unreadCount`
- 源码：`app/controller/merchant/system/notice/SystemNoticeLog.php` :: `unreadCount()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data 对象字段: count | 外层: {status,message,data}


## `mer/openapi`

### `POST /mer/openapi/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.openapi.OpenApi/create`
- 源码：`app/controller/merchant/system/openapi/OpenApi.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/openapi/create/form` — 添加Form

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.openapi.OpenApi/createForm`
- 源码：`app/controller/merchant/system/openapi/OpenApi.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `DELETE /mer/openapi/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.openapi.OpenApi/delete`
- 源码：`app/controller/merchant/system/openapi/OpenApi.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/openapi/get_secret_key/:id` — 查看

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.openapi.OpenApi/getSecretKey`
- 源码：`app/controller/merchant/system/openapi/OpenApi.php` :: `getSecretKey()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/openapi/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.openapi.OpenApi/lst`
- 源码：`app/controller/merchant/system/openapi/OpenApi.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `title` (query/body, 可选)
- `access_key` (query/body, 可选)
- `status` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/openapi/set_secret_key/:id` — 重置

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.openapi.OpenApi/setSecretKey`
- 源码：`app/controller/merchant/system/openapi/OpenApi.php` :: `setSecretKey()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/openapi/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.openapi.OpenApi/switchWithStatus`
- 源码：`app/controller/merchant/system/openapi/OpenApi.php` :: `switchWithStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/openapi/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.openapi.OpenApi/update`
- 源码：`app/controller/merchant/system/openapi/OpenApi.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/openapi/update/:id/form` — 编辑Form

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.openapi.OpenApi/updateForm`
- 源码：`app/controller/merchant/system/openapi/OpenApi.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。


## `mer/product`

### `DELETE /mer/product/cdkey/batch_delete` — 批量删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductCdkey/batchDelete`
- 源码：`app/controller/merchant/store/product/ProductCdkey.php` :: `batchDelete()`
- 请求参数：
- `ids` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/product/cdkey/create` — 添加卡密

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductCdkey/create`
- 源码：`app/controller/merchant/store/product/ProductCdkey.php` :: `create()`
- 请求参数：
- `csList` (query/body, 可选)
- `library_id` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `DELETE /mer/product/cdkey/delete/:id` — 删除卡密

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductCdkey/delete`
- 源码：`app/controller/merchant/store/product/ProductCdkey.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/product/cdkey/library/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.CdkeyLibrary/create`
- 源码：`app/controller/merchant/store/product/CdkeyLibrary.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/product/cdkey/library/create/form` — 添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.CdkeyLibrary/createForm`
- 源码：`app/controller/merchant/store/product/CdkeyLibrary.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `DELETE /mer/product/cdkey/library/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.CdkeyLibrary/delete`
- 源码：`app/controller/merchant/store/product/CdkeyLibrary.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/product/cdkey/library/detail/:id` — 列表

- 置信度：⛔ stale
- 说明：路由已登记，但 `app/controller/merchant/store/product/CdkeyLibrary.php` 中不存在方法 `detail`（原项目死路由/加密扩展/版本差异）。开发勿实现为有效接口，除非核实。
- 处理器：`merchant.store.product.CdkeyLibrary/detail`
- 源码：`app/controller/merchant/store/product/CdkeyLibrary.php` :: `detail()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/product/cdkey/library/excel` — 导出

- 置信度：⛔ stale
- 说明：路由已登记，但 `app/controller/merchant/store/product/CdkeyLibrary.php` 中不存在方法 `excel`（原项目死路由/加密扩展/版本差异）。开发勿实现为有效接口，除非核实。
- 处理器：`merchant.store.product.CdkeyLibrary/excel`
- 源码：`app/controller/merchant/store/product/CdkeyLibrary.php` :: `excel()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/product/cdkey/library/import/:type` — 批量导入

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.StoreImport/Import`
- 源码：`app/controller/merchant/store/StoreImport.php` :: `Import()`
- 请求参数：
- `type` (path, 必填) 路径参数
- `library_id` (query/body, 可选)
- 返回：统一响应包 {status,message,data?}；详见控制器实现 | 外层: {status,message,data}

### `GET /mer/product/cdkey/library/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.CdkeyLibrary/lst`
- 源码：`app/controller/merchant/store/product/CdkeyLibrary.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `status` (query/body, 可选)
- `date` (query/body, 可选)
- `productName` (query/body, 可选)
- `name` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/product/cdkey/library/options` — options

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.CdkeyLibrary/options`
- 源码：`app/controller/merchant/store/product/CdkeyLibrary.php` :: `options()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/product/cdkey/library/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.CdkeyLibrary/update`
- 源码：`app/controller/merchant/store/product/CdkeyLibrary.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/product/cdkey/library/update/:id/form` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.CdkeyLibrary/updateForm`
- 源码：`app/controller/merchant/store/product/CdkeyLibrary.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/product/cdkey/lst` — 卡密列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductCdkey/lst`
- 源码：`app/controller/merchant/store/product/ProductCdkey.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `library_id` (query/body, 可选)
- `status` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/product/cdkey/update/:id` — 编辑卡密

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductCdkey/update`
- 源码：`app/controller/merchant/store/product/ProductCdkey.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `key` (query/body, 可选)
- `pwd` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/product/label/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductLabel/create`
- 源码：`app/controller/merchant/store/product/ProductLabel.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/product/label/create/form` — 添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductLabel/createForm`
- 源码：`app/controller/merchant/store/product/ProductLabel.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /mer/product/label/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductLabel/delete`
- 源码：`app/controller/merchant/store/product/ProductLabel.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/product/label/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductLabel/detail`
- 源码：`app/controller/merchant/store/product/ProductLabel.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/product/label/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductLabel/lst`
- 源码：`app/controller/merchant/store/product/ProductLabel.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `name` (query/body, 可选)
- `type` (query/body, 可选)
- `status` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/product/label/option` — 筛选

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductLabel/getOptions`
- 源码：`app/controller/merchant/store/product/ProductLabel.php` :: `getOptions()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/product/label/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductLabel/switchWithStatus`
- 源码：`app/controller/merchant/store/product/ProductLabel.php` :: `switchWithStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/product/label/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductLabel/update`
- 源码：`app/controller/merchant/store/product/ProductLabel.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/product/label/update/:id/form` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductLabel/updateForm`
- 源码：`app/controller/merchant/store/product/ProductLabel.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/product/unit/create` — 商品单位添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductUnit/create`
- 源码：`app/controller/merchant/store/product/ProductUnit.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/product/unit/create/form` — 商品单位添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductUnit/createForm`
- 源码：`app/controller/merchant/store/product/ProductUnit.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /mer/product/unit/delete/:id` — 商品单位删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductUnit/delete`
- 源码：`app/controller/merchant/store/product/ProductUnit.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/product/unit/list` — 商品单位列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductUnit/list`
- 源码：`app/controller/merchant/store/product/ProductUnit.php` :: `list()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/product/unit/option` — 筛选

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductUnit/getSelectList`
- 源码：`app/controller/merchant/store/product/ProductUnit.php` :: `getSelectList()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/product/unit/update/:id` — 商品单位编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductUnit/update`
- 源码：`app/controller/merchant/store/product/ProductUnit.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/product/unit/update/:id/form` — 商品单位编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductUnit/updateForm`
- 源码：`app/controller/merchant/store/product/ProductUnit.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `mer/product_list`

### `GET /mer/product_list` — product_list

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Product/product_list`
- 源码：`app/controller/merchant/store/product/Product.php` :: `product_list()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选) 默认 ''
- `mer_labels` (query/body, 可选) 默认 ''
- `is_show` (query/body, 可选) 默认 ''
- `mer_cate_id` (query/body, 可选) 默认 ''
- `in_type` (query/body, 可选) 默认 '0,1'
- `status` (query/body, 可选) 默认 1
- `cate_pid` (query/body, 可选) 默认 ''
- `cate_id` (query/body, 可选) 默认 ''
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `mer/profitsharing`

### `GET /mer/profitsharing/export` — 导出

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.order.OrderProfitsharing/export`
- 源码：`app/controller/admin/order/OrderProfitsharing.php` :: `export()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `type` (query/body, 可选)
- `status` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `keyword` (query/body, 可选)
- `profit_date` (query/body, 可选)
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/profitsharing/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.order.OrderProfitsharing/getList`
- 源码：`app/controller/admin/order/OrderProfitsharing.php` :: `getList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `type` (query/body, 可选)
- `status` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `keyword` (query/body, 可选)
- `profit_date` (query/body, 可选)
- `date` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}


## `mer/reservation`

### `GET /mer/reservation/service/list` — 预约日历

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.ReservationService/list`
- 源码：`app/controller/merchant/store/order/ReservationService.php` :: `list()`
- 请求参数：
- `service_type` (query/body, 可选) 默认 ''
- `reservation_keyword` (query/body, 可选) 默认 ''
- `staff_id` (query/body, 可选) 默认 ''
- `uid` (query/body, 可选) 默认 ''
- `phone` (query/body, 可选) 默认 ''
- `nickname` (query/body, 可选) 默认 ''
- `reservation_date` (query/body, 可选) 默认 ''
- `order_type` (query/body, 可选) 默认 4
- `reservation_status` (query/body, 可选) 默认 [
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。


## `mer/serve`

### `GET /mer/serve/code` — 支付二维码

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.serve.Serve/getQrCode`
- 源码：`app/controller/merchant/system/serve/Serve.php` :: `getQrCode()`
- 请求参数：
- `meal_id` (query/body, 可选)
- `pay_type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/serve/config` — getConfig

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.serve.Config/getConfig`
- 源码：`app/controller/merchant/system/serve/Config.php` :: `getConfig()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/serve/config` — setConfig

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.serve.Config/setConfig`
- 源码：`app/controller/merchant/system/serve/Config.php` :: `setConfig()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/serve/detail/:id` — 详情

- 置信度：⛔ stale
- 说明：路由已登记，但 `app/controller/merchant/system/serve/Serve.php` 中不存在方法 `detail`（原项目死路由/加密扩展/版本差异）。开发勿实现为有效接口，除非核实。
- 处理器：`merchant.system.serve.Serve/detail`
- 源码：`app/controller/merchant/system/serve/Serve.php` :: `detail()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/serve/info` — info

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.serve.Config/info`
- 源码：`app/controller/merchant/system/serve/Config.php` :: `info()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/serve/meal` — 套餐列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.serve.Serve/meal`
- 源码：`app/controller/merchant/system/serve/Serve.php` :: `meal()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/serve/paylst` — 购买记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.serve.Serve/lst`
- 源码：`app/controller/merchant/system/serve/Serve.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `status` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。


## `mer/service`

### `GET /mer/service/:id/:uid/lst` — 用户与客服聊天记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.service.StoreService/getUserMsnByService`
- 源码：`app/controller/merchant/store/service/StoreService.php` :: `getUserMsnByService()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `uid` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/service/:id/lst` — 用户与商户聊天记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.service.StoreService/getUserMsnByMerchant`
- 源码：`app/controller/merchant/store/service/StoreService.php` :: `getUserMsnByMerchant()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/service/:id/user` — 客服的全部用户 

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.service.StoreService/serviceUserList`
- 源码：`app/controller/merchant/store/service/StoreService.php` :: `serviceUserList()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/service/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.service.StoreService/create`
- 源码：`app/controller/merchant/store/service/StoreService.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/service/create/form` — 添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.service.StoreService/createForm`
- 源码：`app/controller/merchant/store/service/StoreService.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `DELETE /mer/service/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.service.StoreService/delete`
- 源码：`app/controller/merchant/store/service/StoreService.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/service/list` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.service.StoreService/lst`
- 源码：`app/controller/merchant/store/service/StoreService.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `status` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/service/login/:id` — 登录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.service.StoreService/login`
- 源码：`app/controller/merchant/store/service/StoreService.php` :: `login()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/service/mer/:id/user` — 客服的聊天用户列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.service.StoreService/merchantUserList`
- 源码：`app/controller/merchant/store/service/StoreService.php` :: `merchantUserList()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/service/reply/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.service.StoreServiceReply/create`
- 源码：`app/controller/merchant/store/service/StoreServiceReply.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `DELETE /mer/service/reply/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.service.StoreServiceReply/delete`
- 源码：`app/controller/merchant/store/service/StoreServiceReply.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/service/reply/list` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.service.StoreServiceReply/lst`
- 源码：`app/controller/merchant/store/service/StoreServiceReply.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `status` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/service/reply/status/:id` — 切换状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.service.StoreServiceReply/changeStatus`
- 源码：`app/controller/merchant/store/service/StoreServiceReply.php` :: `changeStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/service/reply/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.service.StoreServiceReply/update`
- 源码：`app/controller/merchant/store/service/StoreServiceReply.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/service/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.service.StoreService/changeStatus`
- 源码：`app/controller/merchant/store/service/StoreService.php` :: `changeStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/service/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.service.StoreService/update`
- 源码：`app/controller/merchant/store/service/StoreService.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/service/update/form/:id` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.service.StoreService/updateForm`
- 源码：`app/controller/merchant/store/service/StoreService.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。


## `mer/staffs`

### `POST /mer/staffs/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.Staffs/create`
- 源码：`app/controller/merchant/store/Staffs.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/staffs/create/form` — 添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.Staffs/createForm`
- 源码：`app/controller/merchant/store/Staffs.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `DELETE /mer/staffs/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.Staffs/delete`
- 源码：`app/controller/merchant/store/Staffs.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/staffs/list` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.Staffs/lst`
- 源码：`app/controller/merchant/store/Staffs.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `name` (query/body, 可选)
- `status` (query/body, 可选)
- `uid` (query/body, 可选)
- `phone` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/staffs/statisticsDetail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.Staffs/statisticsDetail`
- 源码：`app/controller/merchant/store/Staffs.php` :: `statisticsDetail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/staffs/statisticsList` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.Staffs/statisticsList`
- 源码：`app/controller/merchant/store/Staffs.php` :: `statisticsList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `name` (query/body, 可选)
- `uid` (query/body, 可选)
- `phone` (query/body, 可选)
- `service_date` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/staffs/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.Staffs/changeStatus`
- 源码：`app/controller/merchant/store/Staffs.php` :: `changeStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/staffs/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.Staffs/update`
- 源码：`app/controller/merchant/store/Staffs.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/staffs/update/form/:id` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.Staffs/updateForm`
- 源码：`app/controller/merchant/store/Staffs.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。


## `mer/statistics`

### `GET /mer/statistics/get_merchant_count` — 首页未处理业务统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.Common/getMerchantCount`
- 源码：`app/controller/merchant/Common.php` :: `getMerchantCount()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/statistics/get_merchant_todo` — 待办事项

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.Common/getMerchantTodo`
- 源码：`app/controller/merchant/Common.php` :: `getMerchantTodo()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/statistics/get_product_sales_price_top` — 获取商户代办统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.Common/getProductSalesPriceTop`
- 源码：`app/controller/merchant/Common.php` :: `getProductSalesPriceTop()`
- 请求参数：
- `date` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/statistics/main` — 所有数据

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.Common/main`
- 源码：`app/controller/merchant/Common.php` :: `main()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/statistics/order` — 支付订单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.Common/order`
- 源码：`app/controller/merchant/Common.php` :: `order()`
- 请求参数：
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/statistics/product` — 商品支付排行

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.Common/product`
- 源码：`app/controller/merchant/Common.php` :: `product()`
- 请求参数：
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/statistics/product_cart` — 商品加购排行

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.Common/productCart`
- 源码：`app/controller/merchant/Common.php` :: `productCart()`
- 请求参数：
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/statistics/product_visit` — 商品访问排行

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.Common/productVisit`
- 源码：`app/controller/merchant/Common.php` :: `productVisit()`
- 请求参数：
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/statistics/user` — 成交客户

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.Common/user`
- 源码：`app/controller/merchant/Common.php` :: `user()`
- 请求参数：
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/statistics/user_rate` — 成交客户比

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.Common/userRate`
- 源码：`app/controller/merchant/Common.php` :: `userRate()`
- 请求参数：
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `mer/store`

### `DELETE /mer/store/attr/template/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.StoreAttrTemplate/delete`
- 源码：`app/controller/merchant/store/StoreAttrTemplate.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/attr/template/:id` — 文件类型

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.StoreAttrTemplate/update`
- 源码：`app/controller/merchant/store/StoreAttrTemplate.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/attr/template/create` — 添加 

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.StoreAttrTemplate/create`
- 源码：`app/controller/merchant/store/StoreAttrTemplate.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/attr/template/list` — 筛选

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.StoreAttrTemplate/getlist`
- 源码：`app/controller/merchant/store/StoreAttrTemplate.php` :: `getlist()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/store/attr/template/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.StoreAttrTemplate/lst`
- 源码：`app/controller/merchant/store/StoreAttrTemplate.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/store/behalf/cartBatchUpdatePrice` — 批量修改价格

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.behalfcustomerorder.Cart/batchUpdatePrice`
- 源码：`app/controller/merchant/store/behalfcustomerorder/Cart.php` :: `batchUpdatePrice()`
- 请求参数：
- `cart_ids` (query/body, 可选)
- `uid` (query/body, 可选)
- `old_pay_price` (query/body, 可选)
- `change_fee_type` (query/body, 可选)
- `reduce_price` (query/body, 可选)
- `discount_rate` (query/body, 可选)
- `new_pay_price` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/store/behalf/cartChange/:id` — 修改购物车数据

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.behalfcustomerorder.Cart/change`
- 源码：`app/controller/merchant/store/behalfcustomerorder/Cart.php` :: `change()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `uid` (query/body, 可选)
- `cart_num` (query/body, 可选)
- `product_attr_unique` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/store/behalf/cartClear` — 清空购物车

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.behalfcustomerorder.Cart/clear`
- 源码：`app/controller/merchant/store/behalfcustomerorder/Cart.php` :: `clear()`
- 请求参数：
- `uid` (query/body, 可选)
- `tourist_unique_key` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/store/behalf/cartCount` — 购物车总数量

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.behalfcustomerorder.Cart/count`
- 源码：`app/controller/merchant/store/behalfcustomerorder/Cart.php` :: `count()`
- 请求参数：
- `uid` (query/body, 可选)
- `tourist_unique_key` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/store/behalf/cartCreate` — 添加购物车

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.behalfcustomerorder.Cart/create`
- 源码：`app/controller/merchant/store/behalfcustomerorder/Cart.php` :: `create()`
- 请求参数：
- `uid` (query/body, 可选)
- `cart_num` (query/body, 可选)
- `product_id` (query/body, 可选)
- `product_attr_unique` (query/body, 可选)
- `tourist_unique_key` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `DELETE /mer/store/behalf/cartDelete/:id` — 删除购物数据

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.behalfcustomerorder.Cart/delete`
- 源码：`app/controller/merchant/store/behalfcustomerorder/Cart.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/store/behalf/cartList` — 购物车列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.behalfcustomerorder.Cart/list`
- 源码：`app/controller/merchant/store/behalfcustomerorder/Cart.php` :: `list()`
- 请求参数：
- `uid` (query/body, 可选)
- `tourist_unique_key` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/store/behalf/cartUpdatePrice/:id` — 修改价格

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.behalfcustomerorder.Cart/updatePrice`
- 源码：`app/controller/merchant/store/behalfcustomerorder/Cart.php` :: `updatePrice()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `old_price` (query/body, 可选)
- `type` (query/body, 可选)
- `reduce_price` (query/body, 可选)
- `discount_rate` (query/body, 可选)
- `new_price` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/store/behalf/orderCheck` — 校验订单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.behalfcustomerorder.Order/check`
- 源码：`app/controller/merchant/store/behalfcustomerorder/Order.php` :: `check()`
- 请求参数：
- `uid` (query/body, 可选)
- `cart_ids` (query/body, 可选)
- `address_id` (query/body, 可选)
- `delivery_way` (query/body, 可选)
- `use_coupon` (query/body, 可选) 默认 [
- `is_free_shipping` (query/body, 可选)
- `use_integral` (query/body, 可选)
- `tourist_unique_key` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/store/behalf/orderCreate` — 创建订单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.behalfcustomerorder.Order/create`
- 源码：`app/controller/merchant/store/behalfcustomerorder/Order.php` :: `create()`
- 请求参数：
- `uid` (query/body, 可选)
- `cart_ids` (query/body, 可选)
- `address_id` (query/body, 可选)
- `delivery_way` (query/body, 可选)
- `use_coupon` (query/body, 可选) 默认 [
- `is_free_shipping` (query/body, 可选)
- `use_integral` (query/body, 可选)
- `tourist_unique_key` (query/body, 可选)
- `pay_type` (query/body, 可选)
- `key` (query/body, 可选)
- `mark` (query/body, 可选)
- `old_pay_price` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/store/behalf/orderPay/:id` — 支付

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.behalfcustomerorder.Order/pay`
- 源码：`app/controller/merchant/store/behalfcustomerorder/Order.php` :: `pay()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `uid` (query/body, 可选)
- `pay_type` (query/body, 可选)
- `phone` (query/body, 可选)
- `sms_code` (query/body, 可选)
- `auth_code` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/store/behalf/orderVerify` — 余额支付获取验证码

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.behalfcustomerorder.Order/verify`
- 源码：`app/controller/merchant/store/behalfcustomerorder/Order.php` :: `verify()`
- 请求参数：
- `phone` (query/body, 可选)
- `type` (query/body, 可选) 默认 'balance'
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/store/behalf/payConfig` — 支付配置

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.behalfcustomerorder.Order/payConfig`
- 源码：`app/controller/merchant/store/behalfcustomerorder/Order.php` :: `payConfig()`
- 请求参数：
- `uid` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/store/behalf/payStatus/:id` — 获取结果

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.behalfcustomerorder.Order/payStatus`
- 源码：`app/controller/merchant/store/behalfcustomerorder/Order.php` :: `payStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `uid` (query/body, 可选)
- `pay_type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/store/behalf/productCategory` — 商品分类

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.behalfcustomerorder.Product/category`
- 源码：`app/controller/merchant/store/behalfcustomerorder/Product.php` :: `category()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/store/behalf/productDetail/:id` — 商品规格详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.behalfcustomerorder.Product/detail`
- 源码：`app/controller/merchant/store/behalfcustomerorder/Product.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/store/behalf/productList` — 商品列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.behalfcustomerorder.Product/list`
- 源码：`app/controller/merchant/store/behalfcustomerorder/Product.php` :: `list()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `search` (query/body, 可选)
- `mer_cate_id` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/store/behalf/userAddressCreate` — 地址添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.behalfcustomerorder.User/addressCreate`
- 源码：`app/controller/merchant/store/behalfcustomerorder/User.php` :: `addressCreate()`
- 请求参数：
- `uid` (query/body, 可选)
- `real_name` (query/body, 可选)
- `phone` (query/body, 可选)
- `province` (query/body, 可选)
- `province_id` (query/body, 可选)
- `city` (query/body, 可选)
- `city_id` (query/body, 可选)
- `district` (query/body, 可选)
- `district_id` (query/body, 可选)
- `street` (query/body, 可选)
- `street_id` (query/body, 可选)
- `detail` (query/body, 可选)
- `tourist_unique_key` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/store/behalf/userAddressList` — 地址列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.behalfcustomerorder.User/addressList`
- 源码：`app/controller/merchant/store/behalfcustomerorder/User.php` :: `addressList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `uid` (query/body, 可选)
- `tourist_unique_key` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/store/behalf/userCreate` — 会员添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.behalfcustomerorder.User/create`
- 源码：`app/controller/merchant/store/behalfcustomerorder/User.php` :: `create()`
- 请求参数：
- `nickname` (query/body, 可选)
- `phone` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/store/behalf/userInfo` — 会员详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.behalfcustomerorder.User/info`
- 源码：`app/controller/merchant/store/behalfcustomerorder/User.php` :: `info()`
- 请求参数：
- `uid` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/store/behalf/userQuery` — 会员查询

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.behalfcustomerorder.User/query`
- 源码：`app/controller/merchant/store/behalfcustomerorder/User.php` :: `query()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `search` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/store/category/brandlist` — 品牌列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreCategory/BrandList`
- 源码：`app/controller/admin/store/StoreCategory.php` :: `BrandList()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/store/category/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreCategory/create`
- 源码：`app/controller/admin/store/StoreCategory.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/category/create/form` — 添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreCategory/createForm`
- 源码：`app/controller/admin/store/StoreCategory.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /mer/store/category/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreCategory/delete`
- 源码：`app/controller/admin/store/StoreCategory.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/category/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreCategory/detail`
- 源码：`app/controller/admin/store/StoreCategory.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/store/category/list` — 筛选

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreCategory/getList`
- 源码：`app/controller/admin/store/StoreCategory.php` :: `getList()`
- 请求参数：
- `type` (query/body, 可选)
- `lv` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/store/category/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreCategory/lst`
- 源码：`app/controller/admin/store/StoreCategory.php` :: `lst()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/store/category/select` — 树形

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreCategory/getTreeList`
- 源码：`app/controller/admin/store/StoreCategory.php` :: `getTreeList()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/store/category/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreCategory/switchStatus`
- 源码：`app/controller/admin/store/StoreCategory.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/category/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreCategory/update`
- 源码：`app/controller/admin/store/StoreCategory.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/category/update/form/:id` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreCategory/updateForm`
- 源码：`app/controller/admin/store/StoreCategory.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/coupon/clone/form/:id` — 复制表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.coupon.Coupon/cloneForm`
- 源码：`app/controller/merchant/store/coupon/Coupon.php` :: `cloneForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/coupon/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.coupon.Coupon/create`
- 源码：`app/controller/merchant/store/coupon/Coupon.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/coupon/create/form` — 添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.coupon.Coupon/createForm`
- 源码：`app/controller/merchant/store/coupon/Coupon.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /mer/store/coupon/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.coupon.Coupon/delete`
- 源码：`app/controller/merchant/store/coupon/Coupon.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/coupon/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.coupon.Coupon/detail`
- 源码：`app/controller/merchant/store/coupon/Coupon.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/store/coupon/issue` — 使用记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.coupon.Coupon/issue`
- 源码：`app/controller/merchant/store/coupon/Coupon.php` :: `issue()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `username` (query/body, 可选)
- `coupon` (query/body, 可选)
- `status` (query/body, 可选)
- `coupon_id` (query/body, 可选)
- `type` (query/body, 可选)
- `send_id` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/store/coupon/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.coupon.Coupon/lst`
- 源码：`app/controller/merchant/store/coupon/Coupon.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `is_full_give` (query/body, 可选)
- `status` (query/body, 可选)
- `is_give_subscribe` (query/body, 可选)
- `coupon_name` (query/body, 可选)
- `send_type` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /mer/store/coupon/product/:id` — 优惠券可用商品

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Coupon/product`
- 源码：`app/controller/admin/store/Coupon.php` :: `product()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/coupon/select` — 筛选

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.coupon.Coupon/select`
- 源码：`app/controller/merchant/store/coupon/Coupon.php` :: `select()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `coupon_name` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /mer/store/coupon/send` — 发送优惠券

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.coupon.Coupon/send`
- 源码：`app/controller/merchant/store/coupon/Coupon.php` :: `send()`
- 请求参数：
- `coupon_id` (query/body, 可选)
- `mark` (query/body, 可选)
- `is_all` (query/body, 可选)
- `search` (query/body, 可选)
- `uid` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/coupon/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.coupon.Coupon/changeStatus`
- 源码：`app/controller/merchant/store/coupon/Coupon.php` :: `changeStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/coupon/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.coupon.Coupon/update`
- 源码：`app/controller/merchant/store/coupon/Coupon.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `title` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/coupon/update/:id/form` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.coupon.Coupon/updateForm`
- 源码：`app/controller/merchant/store/coupon/Coupon.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/coupon_send/lst` — 发送优惠券记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.coupon.CouponSend/lst`
- 源码：`app/controller/merchant/store/coupon/CouponSend.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `date` (query/body, 可选)
- `coupon_type` (query/body, 可选)
- `coupon_name` (query/body, 可选)
- `status` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /mer/store/form/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.form.Form/create`
- 源码：`app/controller/admin/system/form/Form.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /mer/store/form/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.form.Form/delete`
- 源码：`app/controller/admin/system/form/Form.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/form/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.form.Form/detail`
- 源码：`app/controller/admin/system/form/Form.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/store/form/info/:id` — info

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.form.Form/info`
- 源码：`app/controller/admin/system/form/Form.php` :: `info()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `mer_id` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/form/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.form.Form/lst`
- 源码：`app/controller/admin/system/form/Form.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `date` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /mer/store/form/select` — select

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.form.Form/select`
- 源码：`app/controller/admin/system/form/Form.php` :: `select()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/store/form/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.form.Form/update`
- 源码：`app/controller/admin/system/form/Form.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/import/:type` — 导入

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.StoreImport/Import`
- 源码：`app/controller/merchant/store/StoreImport.php` :: `Import()`
- 请求参数：
- `type` (path, 必填) 路径参数
- `library_id` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/store/import/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.StoreImport/detail`
- 源码：`app/controller/merchant/store/StoreImport.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/store/import/excel/:id` — 导出发货记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.StoreImport/export`
- 源码：`app/controller/merchant/store/StoreImport.php` :: `export()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/store/import/lst` — 导入记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.StoreImport/lst`
- 源码：`app/controller/merchant/store/StoreImport.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `status` (query/body, 可选)
- `date` (query/body, 可选)
- `import_type` (query/body, 可选) 默认 'delivery'
- `type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/store/order/chart` — 统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.Order/chart`
- 源码：`app/controller/merchant/store/order/Order.php` :: `chart()`
- 请求参数：
- `date` (query/body, 可选)
- `order_sn` (query/body, 可选)
- `username` (query/body, 可选)
- `order_type` (query/body, 可选)
- `keywords` (query/body, 可选)
- `order_id` (query/body, 可选)
- `activity_type` (query/body, 可选)
- `group_order_sn` (query/body, 可选)
- `store_name` (query/body, 可选)
- `filter_delivery` (query/body, 可选)
- `filter_product` (query/body, 可选)
- `delivery_id` (query/body, 可选)
- `group_order_id` (query/body, 可选)
- `uid` (query/body, 可选)
- `nickname` (query/body, 可选)
- `real_name` (query/body, 可选)
- `phone` (query/body, 可选)
- `is_behalf` (query/body, 可选)
- `is_virtual` (query/body, 可选)
- `pay_type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/store/order/children/:id` — 关联订单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.Order/childrenList`
- 源码：`app/controller/merchant/store/order/Order.php` :: `childrenList()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/store/order/collectCargo/:id` — 修改收货信息

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.Order/collectCargo`
- 源码：`app/controller/merchant/store/order/Order.php` :: `collectCargo()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `real_name` (query/body, 可选)
- `user_phone` (query/body, 可选)
- `user_address` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/store/order/collectCargo/:id/form` — 修改收货信息表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.Order/collectCargoForm`
- 源码：`app/controller/merchant/store/order/Order.php` :: `collectCargoForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/store/order/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.Order/delete`
- 源码：`app/controller/merchant/store/order/Order.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/store/order/delivery/:id` — 发货

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.Order/delivery`
- 源码：`app/controller/merchant/store/order/Order.php` :: `delivery()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `is_split` (query/body, 可选)
- `split` (query/body, 可选) 默认 [
- `delivery_type` (query/body, 可选)
- `delivery_name` (query/body, 可选)
- `delivery_id` (query/body, 可选)
- `remark` (query/body, 可选)
- `from_name` (query/body, 可选)
- `from_tel` (query/body, 可选)
- `from_addr` (query/body, 可选)
- `temp_id` (query/body, 可选)
- `is_cargo` (query/body, 可选) 默认 1
- `station_id` (query/body, 可选)
- `mark` (query/body, 可选)
- `cargo_weight` (query/body, 可选) 默认 0
- `day_type` (query/body, 可选) 默认 0
- `service_type` (query/body, 可选)
- `pickup_start_time` (query/body, 可选)
- `pickup_end_time` (query/body, 可选)
- `weight` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/store/order/delivery/confirm/:id` — 同城配送核销

- 置信度：✅ high
- 说明：补全 Order 控制器前缀
- 处理器：`merchant.store.order.Order/deliveryConfirm`
- 源码：`app/controller/merchant/store/order/Order.php` :: `deliveryConfirm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/store/order/delivery/dispatch/:id` — 同城配送派单

- 置信度：✅ high
- 说明：补全 Order 控制器前缀
- 处理器：`merchant.store.order.Order/deliveryDispatch`
- 源码：`app/controller/merchant/store/order/Order.php` :: `deliveryDispatch()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/store/order/delivery/updateDispatch/:id` — 同城配送改派

- 置信度：✅ high
- 说明：补全 Order 控制器前缀
- 处理器：`merchant.store.order.Order/deliveryUpdateDispatch`
- 源码：`app/controller/merchant/store/order/Order.php` :: `deliveryUpdateDispatch()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/store/order/deliveryOrderSync` — 同城配送再次同步

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.Order/deliveryOrderSync`
- 源码：`app/controller/merchant/store/order/Order.php` :: `deliveryOrderSync()`
- 请求参数：
- `order_ids` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/store/order/delivery_batch` — 批量发货

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.Order/batchDelivery`
- 源码：`app/controller/merchant/store/order/Order.php` :: `batchDelivery()`
- 请求参数：
- `order_id` (query/body, 可选)
- `delivery_id` (query/body, 可选)
- `delivery_type` (query/body, 可选)
- `delivery_name` (query/body, 可选)
- `remark` (query/body, 可选)
- `select_type` (query/body, 可选) 默认 'select'
- `from_name` (query/body, 可选)
- `from_tel` (query/body, 可选)
- `from_addr` (query/body, 可选)
- `temp_id` (query/body, 可选)
- `is_cargo` (query/body, 可选) 默认 1
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/store/order/delivery_export` — 导出发货单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.Order/deliveryExport`
- 源码：`app/controller/merchant/store/order/Order.php` :: `deliveryExport()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `username` (query/body, 可选)
- `date` (query/body, 可选)
- `activity_type` (query/body, 可选)
- `order_type` (query/body, 可选)
- `keywords` (query/body, 可选)
- `id` (query/body, 可选)
- `ids` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/store/order/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.Order/detail`
- 源码：`app/controller/merchant/store/order/Order.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/store/order/excel` — 导出

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.Order/excel`
- 源码：`app/controller/merchant/store/order/Order.php` :: `excel()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `status` (query/body, 可选)
- `date` (query/body, 可选)
- `order_sn` (query/body, 可选)
- `order_type` (query/body, 可选)
- `username` (query/body, 可选)
- `keywords` (query/body, 可选)
- `take_order` (query/body, 可选)
- `order_id` (query/body, 可选)
- `activity_type` (query/body, 可选)
- `group_order_sn` (query/body, 可选)
- `store_name` (query/body, 可选)
- `filter_delivery` (query/body, 可选)
- `filter_product` (query/body, 可选)
- `pay_type` (query/body, 可选)
- `uid` (query/body, 可选)
- `nickname` (query/body, 可选)
- `real_name` (query/body, 可选)
- `phone` (query/body, 可选)
- `ids` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/store/order/express/:id` — 快递查询

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.Order/express`
- 源码：`app/controller/merchant/store/order/Order.php` :: `express()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/store/order/filtter` — 类型

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.Order/orderType`
- 源码：`app/controller/merchant/store/order/Order.php` :: `orderType()`
- 请求参数：
- `date` (query/body, 可选)
- `order_sn` (query/body, 可选)
- `username` (query/body, 可选)
- `order_type` (query/body, 可选)
- `keywords` (query/body, 可选)
- `order_id` (query/body, 可选)
- `activity_type` (query/body, 可选)
- `group_order_sn` (query/body, 可选)
- `store_name` (query/body, 可选)
- `filter_delivery` (query/body, 可选)
- `filter_product` (query/body, 可选)
- `delivery_id` (query/body, 可选)
- `group_order_id` (query/body, 可选)
- `uid` (query/body, 可选)
- `nickname` (query/body, 可选)
- `real_name` (query/body, 可选)
- `phone` (query/body, 可选)
- `is_behalf` (query/body, 可选)
- `is_virtual` (query/body, 可选)
- `pay_type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/store/order/kuaidi_coms` — getKuaidiComs

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.Order/getKuaidiComs`
- 源码：`app/controller/merchant/store/order/Order.php` :: `getKuaidiComs()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/store/order/log/:id` — 操作记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.Order/status`
- 源码：`app/controller/merchant/store/order/Order.php` :: `status()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `date` (query/body, 可选)
- `user_type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/store/order/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.Order/lst`
- 源码：`app/controller/merchant/store/order/Order.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `status` (query/body, 可选)
- `date` (query/body, 可选)
- `order_sn` (query/body, 可选)
- `username` (query/body, 可选)
- `order_type` (query/body, 可选)
- `keywords` (query/body, 可选)
- `order_id` (query/body, 可选)
- `activity_type` (query/body, 可选)
- `group_order_sn` (query/body, 可选)
- `store_name` (query/body, 可选)
- `filter_delivery` (query/body, 可选)
- `filter_product` (query/body, 可选)
- `delivery_id` (query/body, 可选)
- `group_order_id` (query/body, 可选)
- `uid` (query/body, 可选)
- `nickname` (query/body, 可选)
- `real_name` (query/body, 可选)
- `phone` (query/body, 可选)
- `is_behalf` (query/body, 可选)
- `is_virtual` (query/body, 可选)
- `merchant_take_id` (query/body, 可选)
- `pay_type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/store/order/note` — 配货单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.Order/note`
- 源码：`app/controller/merchant/store/order/Order.php` :: `note()`
- 请求参数：
- `status` (query/body, 可选)
- `date` (query/body, 可选)
- `order_sn` (query/body, 可选)
- `order_type` (query/body, 可选)
- `username` (query/body, 可选)
- `keywords` (query/body, 可选)
- `take_order` (query/body, 可选)
- `order_id` (query/body, 可选)
- `activity_type` (query/body, 可选)
- `group_order_sn` (query/body, 可选)
- `store_name` (query/body, 可选)
- `filter_delivery` (query/body, 可选)
- `filter_product` (query/body, 可选)
- `pay_type` (query/body, 可选)
- `uid` (query/body, 可选)
- `nickname` (query/body, 可选)
- `real_name` (query/body, 可选)
- `phone` (query/body, 可选)
- `ids` (query/body, 可选)
- `limit` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/store/order/offline/:id` — 线下支付

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.Order/offline`
- 源码：`app/controller/merchant/store/order/Order.php` :: `offline()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/store/order/printer/:id` — 打印小票

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.Order/printer`
- 源码：`app/controller/merchant/store/order/Order.php` :: `printer()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/store/order/remark/:id` — 备注

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.Order/remark`
- 源码：`app/controller/merchant/store/order/Order.php` :: `remark()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `remark` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/store/order/remark/:id/form` — 备注表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.Order/remarkForm`
- 源码：`app/controller/merchant/store/order/Order.php` :: `remarkForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/store/order/repeat_dump/:id` — 电子面单复打

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.Order/repeatDump`
- 源码：`app/controller/merchant/store/order/Order.php` :: `repeatDump()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/store/order/reservation/dispatch/:id` — 预约订单派单

- 置信度：✅ high
- 说明：补全 Order 控制器前缀
- 处理器：`merchant.store.order.Order/reservationDispatch`
- 源码：`app/controller/merchant/store/order/Order.php` :: `reservationDispatch()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/store/order/reservation/reschedule/:id` — 预约订单改约

- 置信度：✅ high
- 说明：补全 Order 控制器前缀
- 处理器：`merchant.store.order.Order/reservationReschedule`
- 源码：`app/controller/merchant/store/order/Order.php` :: `reservationReschedule()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/store/order/reservation/reservationTime/:id` — 单独修改预约时间

- 置信度：✅ high
- 说明：补全 Order 控制器前缀
- 处理器：`merchant.store.order.Order/reservationTime`
- 源码：`app/controller/merchant/store/order/Order.php` :: `reservationTime()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/store/order/reservation/updateDispatch/:id` — 预约订单改派

- 置信度：✅ high
- 说明：补全 Order 控制器前缀
- 处理器：`merchant.store.order.Order/reservationUpdateDispatch`
- 源码：`app/controller/merchant/store/order/Order.php` :: `reservationUpdateDispatch()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/store/order/reservation/verify/:id` — 预约订单核销

- 置信度：✅ high
- 说明：补全 Order 控制器前缀
- 处理器：`merchant.store.order.Order/reservationVerify`
- 源码：`app/controller/merchant/store/order/Order.php` :: `reservationVerify()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/store/order/shipment/cancel/:id` — 取消商家寄件

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.Order/cancelShipment`
- 源码：`app/controller/merchant/store/order/Order.php` :: `cancelShipment()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `msg` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/store/order/shipment/getPrice/:id` — 获取商家寄件价格

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.Order/getPrice`
- 源码：`app/controller/merchant/store/order/Order.php` :: `getPrice()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `kuaidicom` (query/body, 可选) 默认 ''
- `service_type` (query/body, 可选) 默认 ''
- `send_address` (query/body, 可选) 默认 ''
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/store/order/shipment/list` — 获取商家寄件价格

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.Order/shipmentList`
- 源码：`app/controller/merchant/store/order/Order.php` :: `shipmentList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/store/order/take_title` — 统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.Order/takeTitle`
- 源码：`app/controller/merchant/store/order/Order.php` :: `takeTitle()`
- 请求参数：
- `date` (query/body, 可选)
- `order_sn` (query/body, 可选)
- `username` (query/body, 可选)
- `keywords` (query/body, 可选)
- `pay_type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/store/order/takechart` — 统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.Order/takeChart`
- 源码：`app/controller/merchant/store/order/Order.php` :: `takeChart()`
- 请求参数：
- `date` (query/body, 可选)
- `order_sn` (query/body, 可选)
- `username` (query/body, 可选)
- `keywords` (query/body, 可选)
- `pay_type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/store/order/takelst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.Order/takeLst`
- 源码：`app/controller/merchant/store/order/Order.php` :: `takeLst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `date` (query/body, 可选)
- `order_sn` (query/body, 可选)
- `username` (query/body, 可选)
- `keywords` (query/body, 可选)
- `pay_type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/store/order/title` — 头部统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.Order/title`
- 源码：`app/controller/merchant/store/order/Order.php` :: `title()`
- 请求参数：
- `status` (query/body, 可选)
- `date` (query/body, 可选)
- `order_sn` (query/body, 可选)
- `username` (query/body, 可选)
- `order_type` (query/body, 可选)
- `keywords` (query/body, 可选)
- `order_id` (query/body, 可选)
- `activity_type` (query/body, 可选)
- `filter_delivery` (query/body, 可选)
- `filter_product` (query/body, 可选)
- `delivery_id` (query/body, 可选)
- `uid` (query/body, 可选)
- `nickname` (query/body, 可选)
- `real_name` (query/body, 可选)
- `phone` (query/body, 可选)
- `pay_type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/store/order/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.Order/update`
- 源码：`app/controller/merchant/store/order/Order.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `total_price` (query/body, 可选)
- `pay_postage` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/store/order/update/:id/form` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.Order/updateForm`
- 源码：`app/controller/merchant/store/order/Order.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /mer/store/order/verify/:code` — 核销详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.Order/verifyDetail`
- 源码：`app/controller/merchant/store/order/Order.php` :: `verifyDetail()`
- 请求参数：
- `code` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/store/order/verify/:id` — 核销

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.Order/verify`
- 源码：`app/controller/merchant/store/order/Order.php` :: `verify()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `verify_code` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/store/params/temp/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.parameter.ParameterTemplate/create`
- 源码：`app/controller/admin/parameter/ParameterTemplate.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /mer/store/params/temp/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.parameter.ParameterTemplate/delete`
- 源码：`app/controller/admin/parameter/ParameterTemplate.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /mer/store/params/temp/delete/value/:id` — 删除属性

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.parameter.ParameterTemplate/deleteValue`
- 源码：`app/controller/admin/parameter/ParameterTemplate.php` :: `deleteValue()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/params/temp/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.parameter.ParameterTemplate/detail`
- 源码：`app/controller/admin/parameter/ParameterTemplate.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/store/params/temp/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.parameter.ParameterTemplate/lst`
- 源码：`app/controller/admin/parameter/ParameterTemplate.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `template_name` (query/body, 可选)
- `cate_ids` (query/body, 可选)
- `mer_name` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `is_mer` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/store/params/temp/select` — 筛选列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.parameter.ParameterTemplate/select`
- 源码：`app/controller/admin/parameter/ParameterTemplate.php` :: `select()`
- 请求参数：
- `cate_id` (query/body, 可选) 默认 0
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/store/params/temp/show` — 参数

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.parameter.ParameterTemplate/show`
- 源码：`app/controller/admin/parameter/ParameterTemplate.php` :: `show()`
- 请求参数：
- `template_ids` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/store/params/temp/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.parameter.ParameterTemplate/update`
- 源码：`app/controller/admin/parameter/ParameterTemplate.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/printer/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.StorePrinter/create`
- 源码：`app/controller/merchant/store/StorePrinter.php` :: `create()`
- 请求参数：
- `type` (query/body, 可选)
- `printer_name` (query/body, 可选)
- `printer_appkey` (query/body, 可选)
- `printer_appid` (query/body, 可选)
- `printer_secret` (query/body, 可选)
- `printer_terminal` (query/body, 可选)
- `status` (query/body, 可选)
- `times` (query/body, 可选)
- `print_type` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/printer/create/form` — 添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.StorePrinter/createForm`
- 源码：`app/controller/merchant/store/StorePrinter.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /mer/store/printer/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.StorePrinter/delete`
- 源码：`app/controller/merchant/store/StorePrinter.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/printer/get_content/:id` — 获取配置内容

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.StorePrinter/getContent`
- 源码：`app/controller/merchant/store/StorePrinter.php` :: `getContent()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/printer/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.StorePrinter/lst`
- 源码：`app/controller/merchant/store/StorePrinter.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `status` (query/body, 可选)
- `keyword` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/store/printer/set_content/:id` — 保存打印机内容

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.StorePrinter/setContent`
- 源码：`app/controller/merchant/store/StorePrinter.php` :: `setContent()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `print_content` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/printer/status/:id` — 取消

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.StorePrinter/switchWithStatus`
- 源码：`app/controller/merchant/store/StorePrinter.php` :: `switchWithStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/printer/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.StorePrinter/update`
- 源码：`app/controller/merchant/store/StorePrinter.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `type` (query/body, 可选)
- `printer_name` (query/body, 可选)
- `printer_appkey` (query/body, 可选)
- `printer_appid` (query/body, 可选)
- `printer_secret` (query/body, 可选)
- `printer_terminal` (query/body, 可选)
- `status` (query/body, 可选)
- `times` (query/body, 可选)
- `print_type` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/printer/update/:id/form` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.StorePrinter/updateForm`
- 源码：`app/controller/merchant/store/StorePrinter.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/product/assist/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductAssist/create`
- 源码：`app/controller/merchant/store/product/ProductAssist.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /mer/store/product/assist/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductAssist/delete`
- 源码：`app/controller/merchant/store/product/ProductAssist.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/product/assist/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductAssist/detail`
- 源码：`app/controller/merchant/store/product/ProductAssist.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/store/product/assist/labels/:id` — 设置标签

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductAssist/setLabels`
- 源码：`app/controller/merchant/store/product/ProductAssist.php` :: `setLabels()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `mer_labels` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/product/assist/lst` — 列表 

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductAssist/lst`
- 源码：`app/controller/merchant/store/product/ProductAssist.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `product_status` (query/body, 可选)
- `keyword` (query/body, 可选)
- `is_show` (query/body, 可选)
- `type` (query/body, 可选)
- `presell_type` (query/body, 可选)
- `us_status` (query/body, 可选)
- `product_assist_id` (query/body, 可选)
- `mer_labels` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /mer/store/product/assist/preview` — 预览

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductAssist/preview`
- 源码：`app/controller/merchant/store/product/ProductAssist.php` :: `preview()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/store/product/assist/sort/:id` — 排序

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductAssist/updateSort`
- 源码：`app/controller/merchant/store/product/ProductAssist.php` :: `updateSort()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `sort` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/product/assist/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductAssist/switchStatus`
- 源码：`app/controller/merchant/store/product/ProductAssist.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/product/assist/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductAssist/update`
- 源码：`app/controller/merchant/store/product/ProductAssist.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/product/assist_set/detail/:id` — 活动详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductAssistSet/detail`
- 源码：`app/controller/merchant/store/product/ProductAssistSet.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/store/product/assist_set/lst` — 活动列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductAssistSet/lst`
- 源码：`app/controller/merchant/store/product/ProductAssistSet.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `status` (query/body, 可选)
- `type` (query/body, 可选)
- `date` (query/body, 可选)
- `user_name` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /mer/store/product/attr_value/:id` — 获取规格

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Product/getAttrValue`
- 源码：`app/controller/merchant/store/product/Product.php` :: `getAttrValue()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `DELETE /mer/store/product/batch_delete` — 批量加入回收站

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Product/batchDelete`
- 源码：`app/controller/merchant/store/product/Product.php` :: `batchDelete()`
- 请求参数：
- `ids` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/product/batch_ext` — 批量设置推荐

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Product/batchExtension`
- 源码：`app/controller/merchant/store/product/Product.php` :: `batchExtension()`
- 请求参数：
- `extension_one` (query/body, 可选)
- `extension_two` (query/body, 可选)
- `ids` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/product/batch_guarantee` — 批量设置服务保障

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Product/batchGuarantee`
- 源码：`app/controller/merchant/store/product/Product.php` :: `batchGuarantee()`
- 请求参数：
- `guarantee_template_id` (query/body, 可选)
- `ids` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/product/batch_hot` — 批量设置推荐

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Product/batchHot`
- 源码：`app/controller/merchant/store/product/Product.php` :: `batchHot()`
- 请求参数：
- `ids` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/product/batch_labels` — 批量设置标签

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Product/batchLabels`
- 源码：`app/controller/merchant/store/product/Product.php` :: `batchLabels()`
- 请求参数：
- `mer_labels` (query/body, 可选)
- `ids` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/product/batch_process` — 批量修改商品属性

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Product/batchProcess`
- 源码：`app/controller/merchant/store/product/Product.php` :: `batchProcess()`
- 请求参数：
- `ids` (query/body, 可选)
- `batch_type` (query/body, 可选)
- `batch_select_type` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/store/product/batch_restore` — 批量恢复

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Product/batchRestore`
- 源码：`app/controller/merchant/store/product/Product.php` :: `batchRestore()`
- 请求参数：
- `ids` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/product/batch_status` — 批量上下架

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Product/batchShow`
- 源码：`app/controller/merchant/store/product/Product.php` :: `batchShow()`
- 请求参数：
- `ids` (query/body, 可选)
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/product/batch_svip` — 批量设置会员价

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Product/batchSvipType`
- 源码：`app/controller/merchant/store/product/Product.php` :: `batchSvipType()`
- 请求参数：
- `svip_price_type` (query/body, 可选) 默认 0
- `ids` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/product/batch_temp` — 批量设置运费模板

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Product/batchTemplate`
- 源码：`app/controller/merchant/store/product/Product.php` :: `batchTemplate()`
- 请求参数：
- `temp_id` (query/body, 可选)
- `ids` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/product/config` — 配置

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Product/config`
- 源码：`app/controller/merchant/store/product/Product.php` :: `config()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/store/product/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Product/create`
- 源码：`app/controller/merchant/store/product/Product.php` :: `create()`
- 请求参数：
- `_see_CREATE_PARAMS` (body, 可选) 请求体字段见对应 Repository::CREATE_PARAMS
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /mer/store/product/delete/:id` — 加入回收站

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Product/delete`
- 源码：`app/controller/merchant/store/product/Product.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /mer/store/product/destory/:id` — 彻底删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Product/destory`
- 源码：`app/controller/merchant/store/product/Product.php` :: `destory()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/product/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Product/detail`
- 源码：`app/controller/merchant/store/product/Product.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `is_copy` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/store/product/free_trial/:id` — 免审编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Product/freeTrial`
- 源码：`app/controller/merchant/store/product/Product.php` :: `freeTrial()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/product/get_attr_value/:id` — 获取规格

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Product/isFormatAttr`
- 源码：`app/controller/merchant/store/product/Product.php` :: `isFormatAttr()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `attrs` (query/body, 可选) 默认 [
- `items` (query/body, 可选) 默认 [
- `product_type` (query/body, 可选) 默认 0
- `is_copy` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/store/product/get_batch_list` — 获取批量修改列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Product/getBatchList`
- 源码：`app/controller/merchant/store/product/Product.php` :: `getBatchList()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/store/product/get_edit/:id` — 编辑商品获取信息

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Product/getEdit`
- 源码：`app/controller/merchant/store/product/Product.php` :: `getEdit()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `is_copy` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/store/product/get_operate_list/:product_id` — 操作记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Product/getOperateList`
- 源码：`app/controller/merchant/store/product/Product.php` :: `getOperateList()`
- 请求参数：
- `product_id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `type` (query/body, 可选) 默认 ''
- `date` (query/body, 可选) 默认 ''
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/product/group/buying/detail/:id` — 活动详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductGroupBuying/detail`
- 源码：`app/controller/merchant/store/product/ProductGroupBuying.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /mer/store/product/group/buying/lst` — 活动列表 

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductGroupBuying/lst`
- 源码：`app/controller/merchant/store/product/ProductGroupBuying.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `status` (query/body, 可选)
- `date` (query/body, 可选)
- `user_name` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/store/product/group/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductGroup/create`
- 源码：`app/controller/merchant/store/product/ProductGroup.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /mer/store/product/group/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductGroup/delete`
- 源码：`app/controller/merchant/store/product/ProductGroup.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/product/group/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductGroup/detail`
- 源码：`app/controller/merchant/store/product/ProductGroup.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/store/product/group/labels/:id` — 设置标签

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductGroup/setLabels`
- 源码：`app/controller/merchant/store/product/ProductGroup.php` :: `setLabels()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `mer_labels` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/product/group/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductGroup/lst`
- 源码：`app/controller/merchant/store/product/ProductGroup.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `product_status` (query/body, 可选)
- `keyword` (query/body, 可选)
- `active_type` (query/body, 可选)
- `status` (query/body, 可选)
- `us_status` (query/body, 可选)
- `product_group_id` (query/body, 可选)
- `mer_labels` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /mer/store/product/group/preview` — 预览

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductGroup/preview`
- 源码：`app/controller/merchant/store/product/ProductGroup.php` :: `preview()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/store/product/group/sort/:id` — 排序

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductGroup/updateSort`
- 源码：`app/controller/merchant/store/product/ProductGroup.php` :: `updateSort()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `sort` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/product/group/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductGroup/switchStatus`
- 源码：`app/controller/merchant/store/product/ProductGroup.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/product/group/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductGroup/update`
- 源码：`app/controller/merchant/store/product/ProductGroup.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/product/labels/:id` — 标签

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Product/setLabels`
- 源码：`app/controller/merchant/store/product/Product.php` :: `setLabels()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `mer_labels` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/product/list` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Product/lst`
- 源码：`app/controller/merchant/store/product/Product.php` :: `lst()`
- 请求参数：
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
- `is_good` (query/body, 可选)
- `not_product_id` (query/body, 可选)
- `form_id` (query/body, 可选)
- `mer_form_id` (query/body, 可选)
- `cate_hot` (query/body, 可选)
- `brand_id` (query/body, 可选)
- `activity_label_ids` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /mer/store/product/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Product/lst`
- 源码：`app/controller/merchant/store/product/Product.php` :: `lst()`
- 请求参数：
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
- `is_good` (query/body, 可选)
- `not_product_id` (query/body, 可选)
- `form_id` (query/body, 可选)
- `mer_form_id` (query/body, 可选)
- `cate_hot` (query/body, 可选)
- `brand_id` (query/body, 可选)
- `activity_label_ids` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /mer/store/product/lst_filter` — 头部统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Product/getStatusFilter`
- 源码：`app/controller/merchant/store/product/Product.php` :: `getStatusFilter()`
- 请求参数：
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
- `is_good` (query/body, 可选)
- `not_product_id` (query/body, 可选)
- `form_id` (query/body, 可选)
- `mer_form_id` (query/body, 可选)
- `cate_hot` (query/body, 可选)
- `brand_id` (query/body, 可选)
- `activity_label_ids` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/store/product/presell/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductPresell/create`
- 源码：`app/controller/merchant/store/product/ProductPresell.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /mer/store/product/presell/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductPresell/delete`
- 源码：`app/controller/merchant/store/product/ProductPresell.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/product/presell/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductPresell/detail`
- 源码：`app/controller/merchant/store/product/ProductPresell.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/store/product/presell/labels/:id` — 设置标签

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductPresell/setLabels`
- 源码：`app/controller/merchant/store/product/ProductPresell.php` :: `setLabels()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `mer_labels` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/product/presell/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductPresell/lst`
- 源码：`app/controller/merchant/store/product/ProductPresell.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `product_status` (query/body, 可选)
- `keyword` (query/body, 可选)
- `type` (query/body, 可选)
- `presell_type` (query/body, 可选)
- `is_show` (query/body, 可选)
- `us_status` (query/body, 可选)
- `product_presell_id` (query/body, 可选)
- `mer_labels` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /mer/store/product/presell/number` — 统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductPresell/number`
- 源码：`app/controller/merchant/store/product/ProductPresell.php` :: `number()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/store/product/presell/preview` — 预览

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductPresell/preview`
- 源码：`app/controller/merchant/store/product/ProductPresell.php` :: `preview()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/store/product/presell/sort/:id` — 排序

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductPresell/updateSort`
- 源码：`app/controller/merchant/store/product/ProductPresell.php` :: `updateSort()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `sort` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/product/presell/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductPresell/switchStatus`
- 源码：`app/controller/merchant/store/product/ProductPresell.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/product/presell/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductPresell/update`
- 源码：`app/controller/merchant/store/product/ProductPresell.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/product/preview` — 预览

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Product/preview`
- 源码：`app/controller/merchant/store/product/Product.php` :: `preview()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/store/product/restore/:id` — 恢复

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Product/restore`
- 源码：`app/controller/merchant/store/product/Product.php` :: `restore()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/product/sort/:id` — 排序

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Product/updateSort`
- 源码：`app/controller/merchant/store/product/Product.php` :: `updateSort()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `sort` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/product/status/:id` — 上下架

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Product/switchStatus`
- 源码：`app/controller/merchant/store/product/Product.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/product/temp_key` — 上传视频配置

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Product/temp_key`
- 源码：`app/controller/merchant/store/product/Product.php` :: `temp_key()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/store/product/unbind` — 操作记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Product/unbind`
- 源码：`app/controller/merchant/store/product/Product.php` :: `unbind()`
- 请求参数：
- `value_id` (query/body, 可选)
- `library_id` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/product/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.Product/update`
- 源码：`app/controller/merchant/store/product/Product.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `_see_CREATE_PARAMS` (body, 可选) 请求体字段见对应 Repository::CREATE_PARAMS
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/productcopy/count` — 统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductCopy/count`
- 源码：`app/controller/merchant/store/product/ProductCopy.php` :: `count()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data 对象字段: count | 外层: {status,message,data}

### `GET /mer/store/productcopy/get` — 获取信息

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductCopy/get`
- 源码：`app/controller/merchant/store/product/ProductCopy.php` :: `get()`
- 请求参数：
- `url` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/store/productcopy/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ProductCopy/lst`
- 源码：`app/controller/merchant/store/product/ProductCopy.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `mer_id` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /mer/store/productcopy/save` — 保存

- 置信度：⛔ stale
- 说明：路由已登记，但 `app/controller/merchant/store/product/ProductCopy.php` 中不存在方法 `save`（原项目死路由/加密扩展/版本差异）。开发勿实现为有效接口，除非核实。
- 处理器：`merchant.store.product.ProductCopy/save`
- 源码：`app/controller/merchant/store/product/ProductCopy.php` :: `save()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一响应包 {status,message,data?}；详见控制器实现 | 外层: {status,message,data}

### `GET /mer/store/receipt/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.OrderReceipt/detail`
- 源码：`app/controller/merchant/store/order/OrderReceipt.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/store/receipt/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.OrderReceipt/Lst`
- 源码：`app/controller/merchant/store/order/OrderReceipt.php` :: `Lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `status` (query/body, 可选)
- `date` (query/body, 可选)
- `receipt_sn` (query/body, 可选)
- `username` (query/body, 可选)
- `order_type` (query/body, 可选)
- `keyword` (query/body, 可选)
- `uid` (query/body, 可选)
- 返回：统一响应包 {status,message,data?}；详见控制器实现 | 外层: {status,message,data}

### `POST /mer/store/receipt/mark/:id` — 备注

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.OrderReceipt/mark`
- 源码：`app/controller/merchant/store/order/OrderReceipt.php` :: `mark()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `mer_mark` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/receipt/mark/:id/form` — 备注表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.OrderReceipt/markForm`
- 源码：`app/controller/merchant/store/order/OrderReceipt.php` :: `markForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/receipt/save_recipt` — 保存发票

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.OrderReceipt/saveRecipt`
- 源码：`app/controller/merchant/store/order/OrderReceipt.php` :: `saveRecipt()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/receipt/set_recipt` — 开发票

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.OrderReceipt/setRecipt`
- 源码：`app/controller/merchant/store/order/OrderReceipt.php` :: `setRecipt()`
- 请求参数：
- `ids` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/store/receipt/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.OrderReceipt/update`
- 源码：`app/controller/merchant/store/order/OrderReceipt.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `receipt_no` (query/body, 可选)
- `mer_mark` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/refundorder/check/:id` — check

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.RefundOrder/check`
- 源码：`app/controller/merchant/store/order/RefundOrder.php` :: `check()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/refundorder/compute` — compute

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.RefundOrder/compute`
- 源码：`app/controller/merchant/store/order/RefundOrder.php` :: `compute()`
- 请求参数：
- `refund` (query/body, 可选)
- `order_id` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/refundorder/create` — 创建

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.RefundOrder/create`
- 源码：`app/controller/merchant/store/order/RefundOrder.php` :: `create()`
- 请求参数：
- `refund_message` (query/body, 可选)
- `refund_price` (query/body, 可选)
- `mer_mark` (query/body, 可选)
- `refund` (query/body, 可选)
- `order_id` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /mer/store/refundorder/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.RefundOrder/delete`
- 源码：`app/controller/merchant/store/order/RefundOrder.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/refundorder/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.RefundOrder/detail`
- 源码：`app/controller/merchant/store/order/RefundOrder.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/store/refundorder/excel` — 导出

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.RefundOrder/createExcel`
- 源码：`app/controller/merchant/store/order/RefundOrder.php` :: `createExcel()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `refund_order_sn` (query/body, 可选)
- `status` (query/body, 可选)
- `refund_type` (query/body, 可选)
- `date` (query/body, 可选)
- `order_sn` (query/body, 可选)
- `id` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/store/refundorder/express/:id` — 快递查询

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.RefundOrder/express`
- 源码：`app/controller/merchant/store/order/RefundOrder.php` :: `express()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/store/refundorder/log/:id` — 操作记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.RefundOrder/log`
- 源码：`app/controller/merchant/store/order/RefundOrder.php` :: `log()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `date` (query/body, 可选)
- `user_type` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/store/refundorder/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.RefundOrder/lst`
- 源码：`app/controller/merchant/store/order/RefundOrder.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `refund_order_sn` (query/body, 可选)
- `status` (query/body, 可选)
- `refund_type` (query/body, 可选)
- `date` (query/body, 可选)
- `order_sn` (query/body, 可选)
- `id` (query/body, 可选)
- `delivery_id` (query/body, 可选)
- `user_type` (query/body, 可选)
- `username` (query/body, 可选)
- `uid` (query/body, 可选)
- `nickname` (query/body, 可选)
- `real_name` (query/body, 可选)
- `phone` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /mer/store/refundorder/mark/:id` — 备注

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.RefundOrder/mark`
- 源码：`app/controller/merchant/store/order/RefundOrder.php` :: `mark()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `mer_mark` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/refundorder/mark/:id/form` — 备注表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.RefundOrder/markForm`
- 源码：`app/controller/merchant/store/order/RefundOrder.php` :: `markForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/refundorder/refund/:id` — 收到退回商品后确认退款

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.RefundOrder/refundPrice`
- 源码：`app/controller/merchant/store/order/RefundOrder.php` :: `refundPrice()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- `fail_message` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/refundorder/refund_message` — refundMessage

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Common/refundMessage`
- 源码：`app/controller/api/Common.php` :: `refundMessage()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/refundorder/status/:id` — 审核

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.RefundOrder/switchStatus`
- 源码：`app/controller/merchant/store/order/RefundOrder.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `mer_delivery_user` (query/body, 可选)
- `mer_delivery_address` (query/body, 可选)
- `phone` (query/body, 可选)
- `status` (query/body, 可选)
- `fail_message` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/refundorder/status/:id/form` — 审核表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.RefundOrder/switchStatusForm`
- 源码：`app/controller/merchant/store/order/RefundOrder.php` :: `switchStatusForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/reply/form/:id` — 回复表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductReply/replyForm`
- 源码：`app/controller/admin/store/StoreProductReply.php` :: `replyForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/reply/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductReply/lst`
- 源码：`app/controller/admin/store/StoreProductReply.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `nickname` (query/body, 可选)
- `is_reply` (query/body, 可选)
- `date` (query/body, 可选)
- `product_id` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /mer/store/reply/reply/:id` — 回复

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductReply/reply`
- 源码：`app/controller/admin/store/StoreProductReply.php` :: `reply()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `content` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/reply/sort/:id` — 排序

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.StoreProductReply/changeSort`
- 源码：`app/controller/merchant/store/StoreProductReply.php` :: `changeSort()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `sort` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/reservation/product/create` — 添加预约商品

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ReservationProduct/create`
- 源码：`app/controller/merchant/store/product/ReservationProduct.php` :: `create()`
- 请求参数：
- `_see_CREATE_PARAMS` (body, 可选) 请求体字段见对应 Repository::CREATE_PARAMS
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/reservation/product/detail/:id` — 预约商品详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ReservationProduct/detail`
- 源码：`app/controller/merchant/store/product/ReservationProduct.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/store/reservation/product/edit/:id` — 获取预约商品

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ReservationProduct/editInfo`
- 源码：`app/controller/merchant/store/product/ReservationProduct.php` :: `editInfo()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/store/reservation/product/edit/:id` — 编辑预约商品

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ReservationProduct/edit`
- 源码：`app/controller/merchant/store/product/ReservationProduct.php` :: `edit()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `_see_CREATE_PARAMS` (body, 可选) 请求体字段见对应 Repository::CREATE_PARAMS
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/reservation/product/setStock/:id` — 批量修改预约商品库存

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ReservationProduct/batchSetReservationProductStock`
- 源码：`app/controller/merchant/store/product/ReservationProduct.php` :: `batchSetReservationProductStock()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `stockValue` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/reservation/product/showDay/:id` — 商品日历day

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ReservationProduct/showDay`
- 源码：`app/controller/merchant/store/product/ReservationProduct.php` :: `showDay()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `date` (query/body, 可选) 默认 date('Y-m')
- `sku_id` (query/body, 可选) 默认 0
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/store/reservation/product/showMonth/:id` — 商品日历month

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.product.ReservationProduct/showMonth`
- 源码：`app/controller/merchant/store/product/ReservationProduct.php` :: `showMonth()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `sku_id` (query/body, 可选) 默认 ''
- `date` (query/body, 可选) 默认 date('Y-m')
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/store/seckill/active/chart_order/:id` — 活动订单统计列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.seckill.SeckillActive/chart_order`
- 源码：`app/controller/merchant/store/seckill/SeckillActive.php` :: `chart_order()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `status` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `date` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /mer/store/seckill/active/chart_panel/:id` — 活动统计数据面板

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.seckill.SeckillActive/chart_panel`
- 源码：`app/controller/merchant/store/seckill/SeckillActive.php` :: `chart_panel()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/store/seckill/active/chart_people/:id` — 活动参与人统计列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.seckill.SeckillActive/chart_people`
- 源码：`app/controller/merchant/store/seckill/SeckillActive.php` :: `chart_people()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `date` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /mer/store/seckill/active/chart_product/:id` — 活动商品统计列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.seckill.SeckillActive/chart_product`
- 源码：`app/controller/merchant/store/seckill/SeckillActive.php` :: `chart_product()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /mer/store/seckill/active/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.seckill.SeckillActive/detail`
- 源码：`app/controller/merchant/store/seckill/SeckillActive.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/seckill/active/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.seckill.SeckillActive/list`
- 源码：`app/controller/merchant/store/seckill/SeckillActive.php` :: `list()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `name` (query/body, 可选)
- `seckill_active_status` (query/body, 可选)
- `date` (query/body, 可选)
- `active_status` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /mer/store/seckill/active/select` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.seckill.SeckillActive/select`
- 源码：`app/controller/merchant/store/seckill/SeckillActive.php` :: `select()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/seckill/time/lst` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.seckill.SeckillTime/lst`
- 源码：`app/controller/merchant/store/seckill/SeckillTime.php` :: `lst()`
- 请求参数：
- `active_id` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/store/seckill_product/create` — 添加 

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.seckill.SeckillProduct/create`
- 源码：`app/controller/merchant/store/seckill/SeckillProduct.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /mer/store/seckill_product/delete/:id` — 加入回收站

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.seckill.SeckillProduct/delete`
- 源码：`app/controller/merchant/store/seckill/SeckillProduct.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /mer/store/seckill_product/destory/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.seckill.SeckillProduct/destory`
- 源码：`app/controller/merchant/store/seckill/SeckillProduct.php` :: `destory()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/seckill_product/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.seckill.SeckillProduct/detail`
- 源码：`app/controller/merchant/store/seckill/SeckillProduct.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/store/seckill_product/labels/:id` — 设置标签

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.seckill.SeckillProduct/setLabels`
- 源码：`app/controller/merchant/store/seckill/SeckillProduct.php` :: `setLabels()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `mer_labels` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/seckill_product/list` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.seckill.SeckillProduct/get_list`
- 源码：`app/controller/merchant/store/seckill/SeckillProduct.php` :: `get_list()`
- 请求参数：
- `keyword` (query/body, 可选)
- `status` (query/body, 可选)
- `seckill_active_id` (query/body, 可选)
- `us_status` (query/body, 可选)
- `active_status` (query/body, 可选)
- `active_name` (query/body, 可选)
- `sys_labels` (query/body, 可选)
- `mer_cate_id` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `sort` (query/body, 可选)
- `mer_labels` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/store/seckill_product/lst` — 分页列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.seckill.SeckillProduct/get_page_list`
- 源码：`app/controller/merchant/store/seckill/SeckillProduct.php` :: `get_page_list()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `status` (query/body, 可选)
- `seckill_active_id` (query/body, 可选)
- `us_status` (query/body, 可选)
- `active_status` (query/body, 可选)
- `active_name` (query/body, 可选)
- `sys_labels` (query/body, 可选)
- `mer_cate_id` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `sort` (query/body, 可选)
- `mer_labels` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /mer/store/seckill_product/lst_filter` — 统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.seckill.SeckillProduct/getStatusFilter`
- 源码：`app/controller/merchant/store/seckill/SeckillProduct.php` :: `getStatusFilter()`
- 请求参数：
- `keyword` (query/body, 可选)
- `status` (query/body, 可选)
- `seckill_active_id` (query/body, 可选)
- `us_status` (query/body, 可选)
- `active_status` (query/body, 可选)
- `active_name` (query/body, 可选)
- `sys_labels` (query/body, 可选)
- `mer_cate_id` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `sort` (query/body, 可选)
- `mer_labels` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/store/seckill_product/preview` — 预览

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.seckill.SeckillProduct/preview`
- 源码：`app/controller/merchant/store/seckill/SeckillProduct.php` :: `preview()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/store/seckill_product/product_list` — 商品列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.seckill.SeckillProduct/get_product_list`
- 源码：`app/controller/merchant/store/seckill/SeckillProduct.php` :: `get_product_list()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `sys_labels` (query/body, 可选)
- `category_id` (query/body, 可选)
- `us_status` (query/body, 可选)
- `cate_ids` (query/body, 可选)
- `in_type` (query/body, 可选) 默认 '0,1'
- `status` (query/body, 可选) 默认 1
- `type_id` (query/body, 可选)
- `active_id` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /mer/store/seckill_product/restore/:id` — 恢复

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.seckill.SeckillProduct/restore`
- 源码：`app/controller/merchant/store/seckill/SeckillProduct.php` :: `restore()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/seckill_product/sort/:id` — 排序

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.seckill.SeckillProduct/updateSort`
- 源码：`app/controller/merchant/store/seckill/SeckillProduct.php` :: `updateSort()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `sort` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/seckill_product/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.seckill.SeckillProduct/switchStatus`
- 源码：`app/controller/merchant/store/seckill/SeckillProduct.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/seckill_product/update/:id` — 修改商品

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.seckill.SeckillProduct/update`
- 源码：`app/controller/merchant/store/seckill/SeckillProduct.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `attr_value` (query/body, 可选)
- `sort` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/shipping/create` — 添加 

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.shipping.ShippingTemplate/create`
- 源码：`app/controller/merchant/store/shipping/ShippingTemplate.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /mer/store/shipping/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.shipping.ShippingTemplate/delete`
- 源码：`app/controller/merchant/store/shipping/ShippingTemplate.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/store/shipping/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.shipping.ShippingTemplate/detail`
- 源码：`app/controller/merchant/store/shipping/ShippingTemplate.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/store/shipping/list` — 列表 

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.shipping.ShippingTemplate/getList`
- 源码：`app/controller/merchant/store/shipping/ShippingTemplate.php` :: `getList()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /mer/store/shipping/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.shipping.ShippingTemplate/lst`
- 源码：`app/controller/merchant/store/shipping/ShippingTemplate.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `type` (query/body, 可选)
- `name` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /mer/store/shipping/setDefault/:id` — 设置默认模板

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.shipping.ShippingTemplate/setDefault`
- 源码：`app/controller/merchant/store/shipping/ShippingTemplate.php` :: `setDefault()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/store/shipping/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.shipping.ShippingTemplate/update`
- 源码：`app/controller/merchant/store/shipping/ShippingTemplate.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `mer/system`

### `POST /mer/system/admin/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.admin.MerchantAdmin/create`
- 源码：`app/controller/merchant/system/admin/MerchantAdmin.php` :: `create()`
- 请求参数：
- `account` (query/body, 可选)
- `phone` (query/body, 可选)
- `pwd` (query/body, 可选)
- `againPassword` (query/body, 可选)
- `real_name` (query/body, 可选)
- `roles` (query/body, 可选) 默认 [
- `status` (query/body, 可选) 默认 0
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/system/admin/create/form` — 添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.admin.MerchantAdmin/createForm`
- 源码：`app/controller/merchant/system/admin/MerchantAdmin.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /mer/system/admin/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.admin.MerchantAdmin/delete`
- 源码：`app/controller/merchant/system/admin/MerchantAdmin.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/system/admin/edit` — 修改信息

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.admin.MerchantAdmin/edit`
- 源码：`app/controller/merchant/system/admin/MerchantAdmin.php` :: `edit()`
- 请求参数：
- `real_name` (query/body, 可选)
- `phone` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/system/admin/edit/form` — 修改信息表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.admin.MerchantAdmin/editForm`
- 源码：`app/controller/merchant/system/admin/MerchantAdmin.php` :: `editForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/system/admin/edit/password` — 修改密码

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.admin.MerchantAdmin/editPassword`
- 源码：`app/controller/merchant/system/admin/MerchantAdmin.php` :: `editPassword()`
- 请求参数：
- `pwd` (query/body, 可选)
- `againPassword` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/system/admin/edit/password/form` — 修改密码表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.admin.MerchantAdmin/editPasswordForm`
- 源码：`app/controller/merchant/system/admin/MerchantAdmin.php` :: `editPasswordForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/system/admin/log` — 操作日志

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.admin.AdminLog/lst`
- 源码：`app/controller/admin/system/admin/AdminLog.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `section_startTime` (query/body, 可选)
- `section_endTime` (query/body, 可选)
- `admin_id` (query/body, 可选)
- `method` (query/body, 可选)
- `date` (query/body, 可选)
- `keyword` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /mer/system/admin/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.admin.MerchantAdmin/getList`
- 源码：`app/controller/merchant/system/admin/MerchantAdmin.php` :: `getList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `date` (query/body, 可选)
- `status` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /mer/system/admin/password/:id` — 修改密码

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.admin.MerchantAdmin/password`
- 源码：`app/controller/merchant/system/admin/MerchantAdmin.php` :: `password()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `pwd` (query/body, 可选)
- `againPassword` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/system/admin/password/form/:id` — 修改密码表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.admin.MerchantAdmin/passwordForm`
- 源码：`app/controller/merchant/system/admin/MerchantAdmin.php` :: `passwordForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/system/admin/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.admin.MerchantAdmin/switchStatus`
- 源码：`app/controller/merchant/system/admin/MerchantAdmin.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/system/admin/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.admin.MerchantAdmin/update`
- 源码：`app/controller/merchant/system/admin/MerchantAdmin.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `account` (query/body, 可选)
- `phone` (query/body, 可选)
- `real_name` (query/body, 可选)
- `roles` (query/body, 可选) 默认 [
- `status` (query/body, 可选) 默认 0
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/system/admin/update/form/:id` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.admin.MerchantAdmin/updateForm`
- 源码：`app/controller/merchant/system/admin/MerchantAdmin.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/system/attachment/category` — 批量修改

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.attachment.Attachment/batchChangeCategory`
- 源码：`app/controller/admin/system/attachment/Attachment.php` :: `batchChangeCategory()`
- 请求参数：
- `ids` (query/body, 可选) 默认 [
- `attachment_category_id` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/system/attachment/category/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.attachment.AttachmentCategory/create`
- 源码：`app/controller/admin/system/attachment/AttachmentCategory.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/system/attachment/category/create/form` — 添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.attachment.AttachmentCategory/createForm`
- 源码：`app/controller/admin/system/attachment/AttachmentCategory.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /mer/system/attachment/category/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.attachment.AttachmentCategory/delete`
- 源码：`app/controller/admin/system/attachment/AttachmentCategory.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/system/attachment/category/formatLst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.attachment.AttachmentCategory/getFormatList`
- 源码：`app/controller/admin/system/attachment/AttachmentCategory.php` :: `getFormatList()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/system/attachment/category/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.attachment.AttachmentCategory/update`
- 源码：`app/controller/admin/system/attachment/AttachmentCategory.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/system/attachment/category/update/form/:id` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.attachment.AttachmentCategory/updateForm`
- 源码：`app/controller/admin/system/attachment/AttachmentCategory.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/system/attachment/chunk/video` — chunkVideo

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.attachment.Attachment/chunkVideo`
- 源码：`app/controller/admin/system/attachment/Attachment.php` :: `chunkVideo()`
- 请求参数：
- `chunkNumber` (query/body, 可选) 默认 0
- `currentChunkSize` (query/body, 可选) 默认 0
- `chunkSize` (query/body, 可选) 默认 0
- `totalChunks` (query/body, 可选) 默认 0
- `file` (query/body, 可选) 默认 'file'
- `md5` (query/body, 可选) 默认 ''
- `filename` (query/body, 可选) 默认 ''
- 返回：data 对象字段: src, attachment_id | 外层: {status,message,data}

### `DELETE /mer/system/attachment/delete` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.attachment.Attachment/delete`
- 源码：`app/controller/admin/system/attachment/Attachment.php` :: `delete()`
- 请求参数：
- `ids` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/system/attachment/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.attachment.Attachment/getList`
- 源码：`app/controller/admin/system/attachment/Attachment.php` :: `getList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `attachment_category_id` (query/body, 可选) 默认 0
- `order` (query/body, 可选)
- `attachment_name` (query/body, 可选)
- `attachment_type` (query/body, 可选) 默认 0
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /mer/system/attachment/online_upload` — 在线图片

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.attachment.Attachment/onlineUpload`
- 源码：`app/controller/admin/system/attachment/Attachment.php` :: `onlineUpload()`
- 请求参数：
- `id` (query/body, 可选)
- `images` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/system/attachment/scan_upload/image/:token` — 扫码上传图片

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.attachment.Attachment/scanUploadImage`
- 源码：`app/controller/admin/system/attachment/Attachment.php` :: `scanUploadImage()`
- 请求参数：
- `token` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/system/attachment/scan_upload/image/:token` — 扫码上传保存

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.attachment.Attachment/scanUploadSave`
- 源码：`app/controller/admin/system/attachment/Attachment.php` :: `scanUploadSave()`
- 请求参数：
- `token` (path, 必填) 路径参数
- `pid` (query/body, 可选)
- `ids` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/system/attachment/scan_upload/qrcode/:pid` — 上传二维码

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.attachment.Attachment/scanUploadQrcode`
- 源码：`app/controller/admin/system/attachment/Attachment.php` :: `scanUploadQrcode()`
- 请求参数：
- `pid` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/system/attachment/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.attachment.Attachment/update`
- 源码：`app/controller/admin/system/attachment/Attachment.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `attachment_name` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/system/attachment/update/:id/form` — 编辑表单的

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.attachment.Attachment/updateForm`
- 源码：`app/controller/admin/system/attachment/Attachment.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/system/city/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.shipping.City/lst`
- 源码：`app/controller/merchant/store/shipping/City.php` :: `lst()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/system/role/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.auth.Role/create`
- 源码：`app/controller/merchant/system/auth/Role.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/system/role/create/form` — 添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.auth.Role/createForm`
- 源码：`app/controller/merchant/system/auth/Role.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /mer/system/role/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.auth.Role/delete`
- 源码：`app/controller/merchant/system/auth/Role.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/system/role/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.auth.Role/getList`
- 源码：`app/controller/merchant/system/auth/Role.php` :: `getList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /mer/system/role/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.auth.Role/switchStatus`
- 源码：`app/controller/merchant/system/auth/Role.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/system/role/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.auth.Role/update`
- 源码：`app/controller/merchant/system/auth/Role.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/system/role/update/form/:id` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.auth.Role/updateForm`
- 源码：`app/controller/merchant/system/auth/Role.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `mer/take`

### `GET /mer/take/info` — 到店自提信息

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.Merchant/takeInfo`
- 源码：`app/controller/merchant/system/Merchant.php` :: `takeInfo()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /mer/take/update` — 保存到店自提信息

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.Merchant/take`
- 源码：`app/controller/merchant/system/Merchant.php` :: `take()`
- 请求参数：
- `mer_take_status` (query/body, 可选)
- `mer_take_name` (query/body, 可选)
- `mer_take_phone` (query/body, 可选)
- `mer_take_address` (query/body, 可选)
- `mer_take_location` (query/body, 可选)
- `mer_take_day` (query/body, 可选)
- `mer_take_time` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `mer/test`

### `GET /mer/test` — test

- 置信度：⛔ stale
- 说明：路由已登记，但 `app/controller/merchant/system/admin/Login.php` 中不存在方法 `test`（原项目死路由/加密扩展/版本差异）。开发勿实现为有效接口，除非核实。
- 处理器：`merchant.system.admin.Login/test`
- 源码：`app/controller/merchant/system/admin/Login.php` :: `test()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一响应包 {status,message,data?}；详见控制器实现 | 外层: {status,message,data}


## `mer/update`

### `GET /mer/update/form` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.system.Merchant/updateForm`
- 源码：`app/controller/merchant/system/Merchant.php` :: `updateForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `mer/upload`

### `POST /mer/upload/certificate` — uploadCertificate

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.Common/uploadCertificate`
- 源码：`app/controller/merchant/Common.php` :: `uploadCertificate()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data 对象字段: src | 外层: {status,message,data}

### `POST /mer/upload/files/:field` — upload

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.attachment.Attachment/upload`
- 源码：`app/controller/admin/system/attachment/Attachment.php` :: `upload()`
- 请求参数：
- `field` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/upload/image/:id/:field` — 上传图片

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.attachment.Attachment/image`
- 源码：`app/controller/admin/system/attachment/Attachment.php` :: `image()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `field` (path, 必填) 路径参数
- `ueditor` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /mer/upload/video` — uploadVideo

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.attachment.Attachment/uploadVideo`
- 源码：`app/controller/admin/system/attachment/Attachment.php` :: `uploadVideo()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。


## `mer/user`

### `POST /mer/user/change_label/:id` — 修改标签

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.user.UserMerchant/changeLabel`
- 源码：`app/controller/merchant/user/UserMerchant.php` :: `changeLabel()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `label_id` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/user/change_label/form/:id` — 修改标签表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.user.UserMerchant/changeLabelForm`
- 源码：`app/controller/merchant/user/UserMerchant.php` :: `changeLabelForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/user/clear_search_log` — 清除用户搜索记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/clearSearchLog`
- 源码：`app/controller/admin/user/User.php` :: `clearSearchLog()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/user/coupon/:uid` — 优惠券

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.user.UserMerchant/coupon`
- 源码：`app/controller/merchant/user/UserMerchant.php` :: `coupon()`
- 请求参数：
- `uid` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `DELETE /mer/user/label/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserLabel/delete`
- 源码：`app/controller/admin/user/UserLabel.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /mer/user/label/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserLabel/update`
- 源码：`app/controller/admin/user/UserLabel.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/user/label/form` — 添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserLabel/createForm`
- 源码：`app/controller/admin/user/UserLabel.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/user/label/form/:id` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserLabel/updateForm`
- 源码：`app/controller/admin/user/UserLabel.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/user/label/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserLabel/lst`
- 源码：`app/controller/admin/user/UserLabel.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `type` (query/body, 可选)
- `all` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /mer/user/label/user/label` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserLabel/create`
- 源码：`app/controller/admin/user/UserLabel.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /mer/user/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.user.UserMerchant/getList`
- 源码：`app/controller/merchant/user/UserMerchant.php` :: `getList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `nickname` (query/body, 可选)
- `sex` (query/body, 可选)
- `is_promoter` (query/body, 可选)
- `user_time_type` (query/body, 可选)
- `user_time` (query/body, 可选)
- `pay_count` (query/body, 可选)
- `label_id` (query/body, 可选)
- `user_type` (query/body, 可选)
- `uid` (query/body, 可选)
- `phone` (query/body, 可选)
- `keyword` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /mer/user/order/:uid` — 订单列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.user.UserMerchant/order`
- 源码：`app/controller/merchant/user/UserMerchant.php` :: `order()`
- 请求参数：
- `uid` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /mer/user/search_log` — 搜索记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/searchLog`
- 源码：`app/controller/admin/user/User.php` :: `searchLog()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `date` (query/body, 可选)
- `keyword` (query/body, 可选)
- `nickname` (query/body, 可选)
- `user_type` (query/body, 可选)
- `uid` (query/body, 可选)
- `phone` (query/body, 可选)
- `real_name` (query/body, 可选)
- 返回：统一响应包 {status,message,data?}；详见控制器实现 | 外层: {status,message,data}

### `GET /mer/user/search_log/export` — 用户搜索记录导出

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/exportSearchLog`
- 源码：`app/controller/admin/user/User.php` :: `exportSearchLog()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `date` (query/body, 可选)
- `keyword` (query/body, 可选)
- `nickname` (query/body, 可选)
- `user_type` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `mer/v2`

### `GET /mer/v2/system/city/lst/:pid` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.shipping.City/lstV2`
- 源码：`app/controller/merchant/store/shipping/City.php` :: `lstV2()`
- 请求参数：
- `pid` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `mer/version`

### `GET /mer/version` — version

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.Common/version`
- 源码：`app/controller/admin/Common.php` :: `version()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

