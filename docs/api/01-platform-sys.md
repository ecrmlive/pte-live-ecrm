# 平台后台 `/sys/`

> 对照文档。置信度：high=959 stale=12 unresolved=1。先读 [ACCURACY.md](./ACCURACY.md)、[FUNCTIONAL-TRUTH.md](./FUNCTIONAL-TRUTH.md)。

合计 **972** 条。

## `sys/activity`

### `POST /sys/activity/atmosphere/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在；参数由 checkParams/Validate 自动补全
- 处理器：`admin.store.marketing.StoreAtmosphere/create`
- 源码：`app/controller/admin/store/marketing/StoreAtmosphere.php` :: `create()`
- 请求参数：
- `scope_type` (body, 可选) 来自 checkParams/Validate（自动补全）
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/activity/atmosphere/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.marketing.StoreAtmosphere/delete`
- 源码：`app/controller/admin/store/marketing/StoreAtmosphere.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/activity/atmosphere/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.marketing.StoreAtmosphere/detail`
- 源码：`app/controller/admin/store/marketing/StoreAtmosphere.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/activity/atmosphere/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.marketing.StoreAtmosphere/lst`
- 源码：`app/controller/admin/store/marketing/StoreAtmosphere.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `status` (query/body, 可选)
- `date` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /sys/activity/atmosphere/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.marketing.StoreAtmosphere/statusSwitch`
- 源码：`app/controller/admin/store/marketing/StoreAtmosphere.php` :: `statusSwitch()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/activity/atmosphere/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.marketing.StoreAtmosphere/update`
- 源码：`app/controller/admin/store/marketing/StoreAtmosphere.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/activity/border/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在；参数由 checkParams/Validate 自动补全
- 处理器：`admin.store.marketing.StoreBorder/create`
- 源码：`app/controller/admin/store/marketing/StoreBorder.php` :: `create()`
- 请求参数：
- `scope_type` (body, 可选) 来自 checkParams/Validate（自动补全）
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/activity/border/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.marketing.StoreBorder/delete`
- 源码：`app/controller/admin/store/marketing/StoreBorder.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/activity/border/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.marketing.StoreBorder/detail`
- 源码：`app/controller/admin/store/marketing/StoreBorder.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/activity/border/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.marketing.StoreBorder/lst`
- 源码：`app/controller/admin/store/marketing/StoreBorder.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `status` (query/body, 可选)
- `date` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /sys/activity/border/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.marketing.StoreBorder/statusSwitch`
- 源码：`app/controller/admin/store/marketing/StoreBorder.php` :: `statusSwitch()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/activity/border/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.marketing.StoreBorder/update`
- 源码：`app/controller/admin/store/marketing/StoreBorder.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/activity/cate/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在；参数由 checkParams/Validate 自动补全
- 处理器：`admin.store.StoreActivityCate/create`
- 源码：`app/controller/admin/store/StoreActivityCate.php` :: `create()`
- 请求参数：
- `name` (body, 可选) 来自 checkParams/Validate（自动补全）
- `pic` (body, 可选) 来自 checkParams/Validate（自动补全）
- `status` (body, 可选) 来自 checkParams/Validate（自动补全）
- `sort` (body, 可选) 来自 checkParams/Validate（自动补全）
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/activity/cate/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreActivityCate/delete`
- 源码：`app/controller/admin/store/StoreActivityCate.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/activity/cate/form/create` — 添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreActivityCate/createForm`
- 源码：`app/controller/admin/store/StoreActivityCate.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/activity/cate/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreActivityCate/lst`
- 源码：`app/controller/admin/store/StoreActivityCate.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /sys/activity/cate/select` — systemActivityCateSelect

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreActivityCate/select`
- 源码：`app/controller/admin/store/StoreActivityCate.php` :: `select()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/activity/cate/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreActivityCate/update`
- 源码：`app/controller/admin/store/StoreActivityCate.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/activity/cate/update/:id/form` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreActivityCate/updateForm`
- 源码：`app/controller/admin/store/StoreActivityCate.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/activity/form/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在；参数由 checkParams/Validate 自动补全
- 处理器：`admin.store.marketing.StoreForm/create`
- 源码：`app/controller/admin/store/marketing/StoreForm.php` :: `create()`
- 请求参数：
- `form_id` (body, 可选) 来自 checkParams/Validate（自动补全）
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/activity/form/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.marketing.StoreForm/delete`
- 源码：`app/controller/admin/store/marketing/StoreForm.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/activity/form/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.marketing.StoreForm/detail`
- 源码：`app/controller/admin/store/marketing/StoreForm.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/activity/form/excel/:id` — 活动记录导出

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.marketing.StoreForm/activUserExcel`
- 源码：`app/controller/admin/store/marketing/StoreForm.php` :: `activUserExcel()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/activity/form/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.marketing.StoreForm/lst`
- 源码：`app/controller/admin/store/marketing/StoreForm.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `status` (query/body, 可选)
- `date` (query/body, 可选)
- `form_id` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /sys/activity/form/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.marketing.StoreForm/statusSwitch`
- 源码：`app/controller/admin/store/marketing/StoreForm.php` :: `statusSwitch()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/activity/form/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.marketing.StoreForm/update`
- 源码：`app/controller/admin/store/marketing/StoreForm.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/activity/form/user/lst/:id` — 活动记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.marketing.StoreForm/activUserLst`
- 源码：`app/controller/admin/store/marketing/StoreForm.php` :: `activUserLst()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/activity/label/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在；参数由 checkParams/Validate 自动补全
- 处理器：`admin.store.StoreActivityLabel/create`
- 源码：`app/controller/admin/store/StoreActivityLabel.php` :: `create()`
- 请求参数：
- `type` (body, 可选) 来自 checkParams/Validate（自动补全）
- `label_cate` (body, 可选) 来自 checkParams/Validate（自动补全）
- `label_name` (body, 可选) 来自 checkParams/Validate（自动补全）
- `style_type` (body, 可选) 来自 checkParams/Validate（自动补全）
- `color` (body, 可选) 来自 checkParams/Validate（自动补全）
- `bg_color` (body, 可选) 来自 checkParams/Validate（自动补全）
- `border_color` (body, 可选) 来自 checkParams/Validate（自动补全）
- `icon` (body, 可选) 来自 checkParams/Validate（自动补全）
- `is_show` (body, 可选) 来自 checkParams/Validate（自动补全）
- `status` (body, 可选) 来自 checkParams/Validate（自动补全）
- `sort` (body, 可选) 来自 checkParams/Validate（自动补全）
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/activity/label/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreActivityLabel/delete`
- 源码：`app/controller/admin/store/StoreActivityLabel.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/activity/label/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreActivityLabel/detail`
- 源码：`app/controller/admin/store/StoreActivityLabel.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/activity/label/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreActivityLabel/lst`
- 源码：`app/controller/admin/store/StoreActivityLabel.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `label_cate` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /sys/activity/label/options` — 下拉列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreActivityLabel/options`
- 源码：`app/controller/admin/store/StoreActivityLabel.php` :: `options()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/activity/label/status/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreActivityLabel/status`
- 源码：`app/controller/admin/store/StoreActivityLabel.php` :: `status()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/activity/label/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreActivityLabel/update`
- 源码：`app/controller/admin/store/StoreActivityLabel.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `sys/agreement`

### `GET /sys/agreement/:key` — 协议

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.Cache/getAgree`
- 源码：`app/controller/admin/system/Cache.php` :: `getAgree()`
- 请求参数：
- `key` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/agreement/:key` — 协议保存

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.Cache/saveAgree`
- 源码：`app/controller/admin/system/Cache.php` :: `saveAgree()`
- 请求参数：
- `key` (path, 必填) 路径参数
- `agree` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/agreement/keylst` — 协议列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.Cache/getKeyLst`
- 源码：`app/controller/admin/system/Cache.php` :: `getKeyLst()`
- 请求参数：
- `type` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `sys/ajcaptcha`

### `GET /sys/ajcaptcha` — ajcaptcha

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Auth/ajcaptcha`
- 源码：`app/controller/api/Auth.php` :: `ajcaptcha()`
- 请求参数：
- `captchaType` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `sys/ajcheck`

### `POST /sys/ajcheck` — ajcheck

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Auth/ajcheck`
- 源码：`app/controller/api/Auth.php` :: `ajcheck()`
- 请求参数：
- `token` (query/body, 可选)
- `pointJson` (query/body, 可选)
- `captchaType` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `sys/ajstatus`

### `POST /sys/ajstatus` — ajCaptchaStatus

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.admin.Login/ajCaptchaStatus`
- 源码：`app/controller/admin/system/admin/Login.php` :: `ajCaptchaStatus()`
- 请求参数：
- `account` (query/body, 可选)
- 返回：data 对象字段: status | 外层: {status,message,data}


## `sys/analytics`

### `GET /sys/analytics/order/line_chart` — 折线图统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.analytics.StoreOrder/lineChart`
- 源码：`app/controller/admin/analytics/StoreOrder.php` :: `lineChart()`
- 请求参数：
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/analytics/order/pie_chart/:type` — 折线图统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.analytics.StoreOrder/typePieCahrt`
- 源码：`app/controller/admin/analytics/StoreOrder.php` :: `typePieCahrt()`
- 请求参数：
- `type` (path, 必填) 路径参数
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/analytics/order/top` — 顶部统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.analytics.StoreOrder/top`
- 源码：`app/controller/admin/analytics/StoreOrder.php` :: `top()`
- 请求参数：
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/analytics/product/line_chart` — 折线图统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.analytics.StoreProduct/lineChart`
- 源码：`app/controller/admin/analytics/StoreProduct.php` :: `lineChart()`
- 请求参数：
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/analytics/product/pie_chart/:type` — 折线图统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.analytics.StoreProduct/typePieCahrt`
- 源码：`app/controller/admin/analytics/StoreProduct.php` :: `typePieCahrt()`
- 请求参数：
- `type` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/analytics/product/top` — 顶部统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.analytics.StoreProduct/top`
- 源码：`app/controller/admin/analytics/StoreProduct.php` :: `top()`
- 请求参数：
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/analytics/user/line_chart` — 折线图统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.analytics.User/lineChart`
- 源码：`app/controller/admin/analytics/User.php` :: `lineChart()`
- 请求参数：
- `type` (query/body, 可选)
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/analytics/user/pie_chart` — 折线图统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.analytics.User/typePieCahrt`
- 源码：`app/controller/admin/analytics/User.php` :: `typePieCahrt()`
- 请求参数：
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/analytics/user/top` — 顶部统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.analytics.User/top`
- 源码：`app/controller/admin/analytics/User.php` :: `top()`
- 请求参数：
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `sys/auth`

### `GET /sys/auth` — auth

- 置信度：⛔ stale
- 说明：路由已登记，但 `app/controller/admin/Common.php` 中不存在方法 `auth`（原项目死路由/加密扩展/版本差异）。开发勿实现为有效接口，除非核实。
- 处理器：`admin.Common/auth`
- 源码：`app/controller/admin/Common.php` :: `auth()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一响应包 {status,message,data?}；详见控制器实现 | 外层: {status,message,data}


## `sys/auth_apply`

### `POST /sys/auth_apply` — auth_apply

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.Common/auth_apply`
- 源码：`app/controller/admin/Common.php` :: `auth_apply()`
- 请求参数：
- `company_name` (query/body, 可选) 默认 ''
- `domain_name` (query/body, 可选) 默认 ''
- `order_id` (query/body, 可选) 默认 ''
- `phone` (query/body, 可选) 默认 ''
- `label` (query/body, 可选) 默认 10
- `captcha` (query/body, 可选) 默认 ''
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `sys/bill`

### `GET /sys/bill/brokerage` — brokerage_list

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserBill/brokerage_list`
- 源码：`app/controller/admin/user/UserBill.php` :: `brokerage_list()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `date` (query/body, 可选)
- `type` (query/body, 可选)
- `uid` (query/body, 可选)
- `phone` (query/body, 可选)
- `real_name` (query/body, 可选)
- `nickname` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /sys/bill/export` — 导出

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserBill/export`
- 源码：`app/controller/admin/user/UserBill.php` :: `export()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `date` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/bill/list` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserBill/getList`
- 源码：`app/controller/admin/user/UserBill.php` :: `getList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `date` (query/body, 可选)
- `type` (query/body, 可选)
- `uid` (query/body, 可选)
- `phone` (query/body, 可选)
- `real_name` (query/body, 可选)
- `nickname` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /sys/bill/type` — 类型

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserBill/type`
- 源码：`app/controller/admin/user/UserBill.php` :: `type()`
- 请求参数：
- `category` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `sys/broadcast`

### `POST /sys/broadcast/goods/apply/:id` — 审核

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.BroadcastGoods/apply`
- 源码：`app/controller/admin/store/BroadcastGoods.php` :: `apply()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- `msg` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/broadcast/goods/apply/form/:id` — 审核表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.BroadcastGoods/applyForm`
- 源码：`app/controller/admin/store/BroadcastGoods.php` :: `applyForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/broadcast/goods/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.BroadcastGoods/delete`
- 源码：`app/controller/admin/store/BroadcastGoods.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/broadcast/goods/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.BroadcastGoods/detail`
- 源码：`app/controller/admin/store/BroadcastGoods.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/broadcast/goods/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.BroadcastGoods/lst`
- 源码：`app/controller/admin/store/BroadcastGoods.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `status_tag` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `mer_valid` (query/body, 可选)
- `broadcast_goods_id` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /sys/broadcast/goods/sort/:id` — 排序

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.BroadcastGoods/sort`
- 源码：`app/controller/admin/store/BroadcastGoods.php` :: `sort()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `sort` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/broadcast/goods/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.BroadcastGoods/changeStatus`
- 源码：`app/controller/admin/store/BroadcastGoods.php` :: `changeStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `is_show` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/broadcast/room/apply/:id` — 申请

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.BroadcastRoom/apply`
- 源码：`app/controller/admin/store/BroadcastRoom.php` :: `apply()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- `msg` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/broadcast/room/apply/form/:id` — 申请审核表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.BroadcastRoom/applyForm`
- 源码：`app/controller/admin/store/BroadcastRoom.php` :: `applyForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/broadcast/room/closeKf/:id` — 客服开关

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.BroadcastRoom/closeKf`
- 源码：`app/controller/admin/store/BroadcastRoom.php` :: `closeKf()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/broadcast/room/comment/:id` — 禁言开关

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.BroadcastRoom/banComment`
- 源码：`app/controller/admin/store/BroadcastRoom.php` :: `banComment()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/broadcast/room/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.BroadcastRoom/delete`
- 源码：`app/controller/admin/store/BroadcastRoom.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/broadcast/room/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.BroadcastRoom/detail`
- 源码：`app/controller/admin/store/BroadcastRoom.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/broadcast/room/feedsPublic/:id` — 收录开关

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.BroadcastRoom/isFeedsPublic`
- 源码：`app/controller/admin/store/BroadcastRoom.php` :: `isFeedsPublic()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/broadcast/room/goods/:id` — 商品列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.BroadcastRoom/goodsList`
- 源码：`app/controller/admin/store/BroadcastRoom.php` :: `goodsList()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/broadcast/room/live_status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.BroadcastRoom/changeLiveStatus`
- 源码：`app/controller/admin/store/BroadcastRoom.php` :: `changeLiveStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `replay_status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/broadcast/room/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.BroadcastRoom/lst`
- 源码：`app/controller/admin/store/BroadcastRoom.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `status_tag` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `show_type` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `live_status` (query/body, 可选)
- `star` (query/body, 可选)
- `broadcast_room_id` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /sys/broadcast/room/sort/:id` — 排序

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.BroadcastRoom/sort`
- 源码：`app/controller/admin/store/BroadcastRoom.php` :: `sort()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `sort` (query/body, 可选)
- `star` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/broadcast/room/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.BroadcastRoom/changeStatus`
- 源码：`app/controller/admin/store/BroadcastRoom.php` :: `changeStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `is_show` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `sys/captcha`

### `GET /sys/captcha` — getCaptcha

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.admin.Login/getCaptcha`
- 源码：`app/controller/admin/system/admin/Login.php` :: `getCaptcha()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `sys/change`

### `GET /sys/change/color` — 一键换色

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.Common/getChangeColor`
- 源码：`app/controller/admin/Common.php` :: `getChangeColor()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/change/color` — 一键换色保存

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.Common/setChangeColor`
- 源码：`app/controller/admin/Common.php` :: `setChangeColor()`
- 请求参数：
- `global_theme` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `sys/check`

### `GET /sys/check/queue` — queue

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.Common/queue`
- 源码：`app/controller/admin/Common.php` :: `queue()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一响应包 {status,message,data?}；详见控制器实现 | 外层: {status,message,data}


## `sys/check_auth`

### `GET /sys/check_auth` — check_auth

- 置信度：⛔ stale
- 说明：路由已登记，但 `app/controller/admin/Common.php` 中不存在方法 `check_auth`（原项目死路由/加密扩展/版本差异）。开发勿实现为有效接口，除非核实。
- 处理器：`admin.Common/check_auth`
- 源码：`app/controller/admin/Common.php` :: `check_auth()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一响应包 {status,message,data?}；详见控制器实现 | 外层: {status,message,data}


## `sys/circle`

### `POST /sys/circle/agent/audit/:id` — 商圈代理审核

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.circle.CircleAgent/audit`
- 源码：`app/controller/admin/circle/CircleAgent.php` :: `audit()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- `audit_reason` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/circle/agent/create` — 商圈代理添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.circle.CircleAgent/create`
- 源码：`app/controller/admin/circle/CircleAgent.php` :: `create()`
- 请求参数：
- `type` (query/body, 可选)
- `name` (query/body, 可选)
- `phone` (query/body, 可选)
- `qualification` (query/body, 可选)
- `remark` (query/body, 可选)
- `uid` (query/body, 可选)
- `business_name` (query/body, 可选)
- `account` (query/body, 可选)
- `password` (query/body, 可选)
- `extend` (query/body, 可选) 默认 [
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/circle/agent/delete/:id` — 商圈代理删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.circle.CircleAgent/delete`
- 源码：`app/controller/admin/circle/CircleAgent.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/circle/agent/detail/:id` — 商圈代理详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.circle.CircleAgent/detail`
- 源码：`app/controller/admin/circle/CircleAgent.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/circle/agent/list` — 商圈代理列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.circle.CircleAgent/list`
- 源码：`app/controller/admin/circle/CircleAgent.php` :: `list()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `type` (query/body, 可选)
- `name` (query/body, 可选)
- `phone` (query/body, 可选)
- `uid` (query/body, 可选)
- `user_phone` (query/body, 可选)
- `business_name` (query/body, 可选)
- `nickname` (query/body, 可选)
- `is_apply` (query/body, 可选)
- `status` (query/body, 可选)
- `create_time` (query/body, 可选) 默认 [
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/circle/agent/merchantList/:id` — 关联商户列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.circle.CircleAgent/associatedMerchantList`
- 源码：`app/controller/admin/circle/CircleAgent.php` :: `associatedMerchantList()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `circle_id` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/circle/agent/options` — 代理选项

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.circle.CircleAgent/options`
- 源码：`app/controller/admin/circle/CircleAgent.php` :: `options()`
- 请求参数：
- `type` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/circle/agent/resetPwd/:id` — 重置密码

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.circle.CircleAgent/resetPassword`
- 源码：`app/controller/admin/circle/CircleAgent.php` :: `resetPassword()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/circle/agent/settlementMethod/:id` — 结算方式get

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.circle.CircleAgent/getSettlementMethod`
- 源码：`app/controller/admin/circle/CircleAgent.php` :: `getSettlementMethod()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/circle/agent/settlementMethod/:id` — 结算方式post

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.circle.CircleAgent/setSettlementMethod`
- 源码：`app/controller/admin/circle/CircleAgent.php` :: `setSettlementMethod()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `payment_method` (query/body, 可选)
- `payment_name` (query/body, 可选)
- `payment_account` (query/body, 可选)
- `payment_bank` (query/body, 可选)
- `payment_qr_img` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/circle/agent/update/:id` — 商圈代理编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.circle.CircleAgent/update`
- 源码：`app/controller/admin/circle/CircleAgent.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `type` (query/body, 可选)
- `name` (query/body, 可选)
- `phone` (query/body, 可选)
- `qualification` (query/body, 可选)
- `remark` (query/body, 可选)
- `uid` (query/body, 可选)
- `business_name` (query/body, 可选)
- `extend` (query/body, 可选) 默认 [
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/circle/checkout/audit/:id` — 平台结算审核

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.circle.CircleBrokerageCheckout/audit`
- 源码：`app/controller/admin/circle/CircleBrokerageCheckout.php` :: `audit()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `audit_status` (query/body, 可选)
- `audit_reason` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/circle/checkout/create` — 商圈申请结算获取余额

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.circle.CircleBrokerageCheckout/create`
- 源码：`app/controller/admin/circle/CircleBrokerageCheckout.php` :: `create()`
- 请求参数：
- `agent_id` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/circle/checkout/create` — 商圈申请结算提交

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.circle.CircleBrokerageCheckout/save`
- 源码：`app/controller/admin/circle/CircleBrokerageCheckout.php` :: `save()`
- 请求参数：
- `agent_id` (query/body, 可选)
- `withdrawal_amount` (query/body, 可选)
- `withdrawal_type` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/circle/checkout/detail/:id` — 结算记录详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.circle.CircleBrokerageCheckout/detail`
- 源码：`app/controller/admin/circle/CircleBrokerageCheckout.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/circle/checkout/list` — 结算记录列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.circle.CircleBrokerageCheckout/list`
- 源码：`app/controller/admin/circle/CircleBrokerageCheckout.php` :: `list()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `agent_id` (query/body, 可选)
- `agent_phone` (query/body, 可选)
- `create_time` (query/body, 可选)
- `audit_status` (query/body, 可选)
- `status` (query/body, 可选)
- `withdrawal_type` (query/body, 可选)
- `withdrawal_sn` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/circle/checkout/platformRemark/:id` — 平台备注

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.circle.CircleBrokerageCheckout/platformRemark`
- 源码：`app/controller/admin/circle/CircleBrokerageCheckout.php` :: `platformRemark()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `platform_remark` (query/body, 可选) 默认 ''
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/circle/checkout/remark/:id` — 商圈备注

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.circle.CircleBrokerageCheckout/remark`
- 源码：`app/controller/admin/circle/CircleBrokerageCheckout.php` :: `remark()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `remark` (query/body, 可选) 默认 ''
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/circle/checkout/revoke/:id` — 商圈撤销结算

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.circle.CircleBrokerageCheckout/revoke`
- 源码：`app/controller/admin/circle/CircleBrokerageCheckout.php` :: `revoke()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/circle/checkout/transfer/:id` — 平台转账

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.circle.CircleBrokerageCheckout/transfer`
- 源码：`app/controller/admin/circle/CircleBrokerageCheckout.php` :: `transfer()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `transfer_voucher` (query/body, 可选)
- `transfer_remark` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/circle/create` — 商圈添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.circle.Circle/create`
- 源码：`app/controller/admin/circle/Circle.php` :: `create()`
- 请求参数：
- `pid` (query/body, 可选) 默认 0
- `name` (query/body, 可选)
- `circle_agent_id` (query/body, 可选)
- `commission_type` (query/body, 可选)
- `commission_rate` (query/body, 可选)
- `sort` (query/body, 可选) 默认 0
- `status` (query/body, 可选) 默认 1
- `merchant_ids` (query/body, 可选)
- `type` (query/body, 可选)
- `role_id` (query/body, 可选)
- `business_store_category` (query/body, 可选)
- `business_store_type` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/circle/delete/:id` — 商圈删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.circle.Circle/delete`
- 源码：`app/controller/admin/circle/Circle.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/circle/detail/:id` — 商圈详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.circle.Circle/detail`
- 源码：`app/controller/admin/circle/Circle.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/circle/financialRecord/list` — 商圈提成流水列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.circle.CircleFinancialRecord/list`
- 源码：`app/controller/admin/circle/CircleFinancialRecord.php` :: `list()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `circle_id` (query/body, 可选)
- `mer_name` (query/body, 可选)
- `agent_id` (query/body, 可选)
- `order_time` (query/body, 可选)
- `order_sn` (query/body, 可选)
- `order_id` (query/body, 可选)
- `order_status` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/circle/list` — 商圈列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.circle.Circle/list`
- 源码：`app/controller/admin/circle/Circle.php` :: `list()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `circle_agent_id` (query/body, 可选)
- `name` (query/body, 可选)
- `status` (query/body, 可选)
- `type` (query/body, 可选) 默认 1
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/circle/merchantList/:id` — 关联商户列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.circle.Circle/associatedMerchantList`
- 源码：`app/controller/admin/circle/Circle.php` :: `associatedMerchantList()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/circle/options` — options

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.circle.Circle/options`
- 源码：`app/controller/admin/circle/Circle.php` :: `options()`
- 请求参数：
- `type` (query/body, 可选) 默认 0
- `status` (query/body, 可选) 默认 1
- `circle_agent_id` (query/body, 可选) 默认 ''
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/circle/switch/:id` — 商圈状态切换

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.circle.Circle/switch`
- 源码：`app/controller/admin/circle/Circle.php` :: `switch()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/circle/update/:id` — 商圈编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.circle.Circle/update`
- 源码：`app/controller/admin/circle/Circle.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `pid` (query/body, 可选) 默认 0
- `name` (query/body, 可选)
- `circle_agent_id` (query/body, 可选)
- `commission_type` (query/body, 可选)
- `commission_rate` (query/body, 可选)
- `sort` (query/body, 可选) 默认 0
- `status` (query/body, 可选) 默认 1
- `merchant_ids` (query/body, 可选)
- `type` (query/body, 可选)
- `role_id` (query/body, 可选)
- `business_store_category` (query/body, 可选)
- `business_store_type` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `sys/clear`

### `POST /sys/clear/cache` — 清除缓存

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.Cache/clearCache`
- 源码：`app/controller/admin/system/Cache.php` :: `clearCache()`
- 请求参数：
- `type` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `sys/community`

### `POST /sys/community/category/create` — 社区分类添加

- 置信度：✅ high
- 说明：控制器方法已校验存在；参数由 checkParams/Validate 自动补全
- 处理器：`admin.community.CommunityCategory/create`
- 源码：`app/controller/admin/community/CommunityCategory.php` :: `create()`
- 请求参数：
- `pid` (body, 可选) 来自 checkParams/Validate（自动补全）
- `cate_name` (body, 可选) 来自 checkParams/Validate（自动补全）
- `is_show` (body, 可选) 来自 checkParams/Validate（自动补全）
- `sort` (body, 可选) 来自 checkParams/Validate（自动补全）
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/community/category/create/form` — 社区分类添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.community.CommunityCategory/createForm`
- 源码：`app/controller/admin/community/CommunityCategory.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/community/category/delete/:id` — 社区分类删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.community.CommunityCategory/delete`
- 源码：`app/controller/admin/community/CommunityCategory.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/community/category/detail/:id` — 社区分类详情

- 置信度：⛔ stale
- 说明：路由已登记，但 `app/controller/admin/community/CommunityCategory.php` 中不存在方法 `detail`（原项目死路由/加密扩展/版本差异）。开发勿实现为有效接口，除非核实。
- 处理器：`admin.community.CommunityCategory/detail`
- 源码：`app/controller/admin/community/CommunityCategory.php` :: `detail()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一响应包 {status,message,data?}；详见控制器实现 | 外层: {status,message,data}

### `GET /sys/community/category/lst` — 社区分类状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.community.CommunityCategory/lst`
- 源码：`app/controller/admin/community/CommunityCategory.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `cate_name` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /sys/community/category/option` — 社区分类

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.community.CommunityCategory/getOptions`
- 源码：`app/controller/admin/community/CommunityCategory.php` :: `getOptions()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/community/category/status/:id` — 社区分类修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.community.CommunityCategory/switchStatus`
- 源码：`app/controller/admin/community/CommunityCategory.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/community/category/update/:id` — 社区分类编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.community.CommunityCategory/update`
- 源码：`app/controller/admin/community/CommunityCategory.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/community/category/update/:id/form` — 社区分类编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.community.CommunityCategory/updateForm`
- 源码：`app/controller/admin/community/CommunityCategory.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/community/delete/:id` — 文章删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.community.Community/delete`
- 源码：`app/controller/admin/community/Community.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/community/detail/:id` — 文章详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.community.Community/detail`
- 源码：`app/controller/admin/community/Community.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/community/lst` — 文章列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.community.Community/lst`
- 源码：`app/controller/admin/community/Community.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `status` (query/body, 可选)
- `username` (query/body, 可选)
- `category_id` (query/body, 可选)
- `topic_id` (query/body, 可选)
- `is_show` (query/body, 可选)
- `is_type` (query/body, 可选)
- `uid` (query/body, 可选)
- `phone` (query/body, 可选)
- `real_name` (query/body, 可选)
- `nickname` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /sys/community/reply/:id` — 内容评论列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.community.CommunityReply/reply`
- 源码：`app/controller/admin/community/CommunityReply.php` :: `reply()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `DELETE /sys/community/reply/delete/:id` — 社区评论删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.community.CommunityReply/delete`
- 源码：`app/controller/admin/community/CommunityReply.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/community/reply/lst` — 社区评论列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.community.CommunityReply/lst`
- 源码：`app/controller/admin/community/CommunityReply.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `date` (query/body, 可选)
- `username` (query/body, 可选)
- `community_id` (query/body, 可选)
- `pid` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /sys/community/reply/status/:id` — 社区评论审核

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.community.CommunityReply/switchStatus`
- 源码：`app/controller/admin/community/CommunityReply.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- `refusal` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/community/reply/status/:id/form` — 社区评论审核表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.community.CommunityReply/statusForm`
- 源码：`app/controller/admin/community/CommunityReply.php` :: `statusForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/community/show/:id` — 文章详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.community.Community/switchShow`
- 源码：`app/controller/admin/community/Community.php` :: `switchShow()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/community/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.community.Community/switchStatus`
- 源码：`app/controller/admin/community/Community.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- `refusal` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/community/status/:id/form` — 修改状态表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.community.Community/showForm`
- 源码：`app/controller/admin/community/Community.php` :: `showForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/community/title` — 统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.community.Community/title`
- 源码：`app/controller/admin/community/Community.php` :: `title()`
- 请求参数：
- `keyword` (query/body, 可选)
- `status` (query/body, 可选)
- `username` (query/body, 可选)
- `category_id` (query/body, 可选)
- `topic_id` (query/body, 可选)
- `is_show` (query/body, 可选)
- `is_type` (query/body, 可选)
- `uid` (query/body, 可选)
- `phone` (query/body, 可选)
- `real_name` (query/body, 可选)
- `nickname` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/community/topic/create` — 社区话题添加

- 置信度：✅ high
- 说明：控制器方法已校验存在；参数由 checkParams/Validate 自动补全
- 处理器：`admin.community.CommunityTopic/create`
- 源码：`app/controller/admin/community/CommunityTopic.php` :: `create()`
- 请求参数：
- `category_id` (body, 可选) 来自 checkParams/Validate（自动补全）
- `topic_name` (body, 可选) 来自 checkParams/Validate（自动补全）
- `is_hot` (body, 可选) 来自 checkParams/Validate（自动补全）
- `status` (body, 可选) 来自 checkParams/Validate（自动补全）
- `sort` (body, 可选) 来自 checkParams/Validate（自动补全）
- `pic` (body, 可选) 来自 checkParams/Validate（自动补全）
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/community/topic/create/form` — 社区话题添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.community.CommunityTopic/createForm`
- 源码：`app/controller/admin/community/CommunityTopic.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/community/topic/delete/:id` — 社区话题删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.community.CommunityTopic/delete`
- 源码：`app/controller/admin/community/CommunityTopic.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/community/topic/detail/:id` — 社区话题详情 

- 置信度：⛔ stale
- 说明：路由已登记，但 `app/controller/admin/community/CommunityTopic.php` 中不存在方法 `detail`（原项目死路由/加密扩展/版本差异）。开发勿实现为有效接口，除非核实。
- 处理器：`admin.community.CommunityTopic/detail`
- 源码：`app/controller/admin/community/CommunityTopic.php` :: `detail()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一响应包 {status,message,data?}；详见控制器实现 | 外层: {status,message,data}

### `POST /sys/community/topic/hot/:id` — 社区话题推荐

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.community.CommunityTopic/switchHot`
- 源码：`app/controller/admin/community/CommunityTopic.php` :: `switchHot()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/community/topic/lst` — 社区话题

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.community.CommunityTopic/lst`
- 源码：`app/controller/admin/community/CommunityTopic.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `topic_name` (query/body, 可选)
- `category_id` (query/body, 可选)
- `status` (query/body, 可选)
- `is_hot` (query/body, 可选)
- `is_del` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /sys/community/topic/option` — 社区话题

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.community.CommunityTopic/getOptions`
- 源码：`app/controller/admin/community/CommunityTopic.php` :: `getOptions()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/community/topic/status/:id` — 社区话题修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.community.CommunityTopic/switchStatus`
- 源码：`app/controller/admin/community/CommunityTopic.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/community/topic/update/:id` — 社区话题编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.community.CommunityTopic/update`
- 源码：`app/controller/admin/community/CommunityTopic.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/community/topic/update/:id/form` — 社区话题编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.community.CommunityTopic/updateForm`
- 源码：`app/controller/admin/community/CommunityTopic.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/community/update/:id` — 文章编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.community.Community/update`
- 源码：`app/controller/admin/community/Community.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `start` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/community/update/:id/form` — 文章编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.community.Community/updateForm`
- 源码：`app/controller/admin/community/Community.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `sys/config`

### `GET /sys/config` — config

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.Common/config`
- 源码：`app/controller/admin/Common.php` :: `config()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/config/:key` — 获取配置信息

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.config.Config/form`
- 源码：`app/controller/admin/system/config/Config.php` :: `form()`
- 请求参数：
- `key` (path, 必填) 路径参数
- `tab_id` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/config/classify/create` — 配置分类添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.config.ConfigClassify/create`
- 源码：`app/controller/admin/system/config/ConfigClassify.php` :: `create()`
- 请求参数：
- `pid` (query/body, 可选)
- `classify_name` (query/body, 可选)
- `classify_key` (query/body, 可选)
- `info` (query/body, 可选)
- `status` (query/body, 可选)
- `icon` (query/body, 可选)
- `sort` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/config/classify/create/table` — 配置分类添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.config.ConfigClassify/createTable`
- 源码：`app/controller/admin/system/config/ConfigClassify.php` :: `createTable()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/config/classify/delete/:id` — 配置分类删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.config.ConfigClassify/delete`
- 源码：`app/controller/admin/system/config/ConfigClassify.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/config/classify/lst` — 配置分类列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.config.ConfigClassify/lst`
- 源码：`app/controller/admin/system/config/ConfigClassify.php` :: `lst()`
- 请求参数：
- `status` (query/body, 可选)
- `classify_name` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/config/classify/options` — 配置分类筛选

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.config.ConfigClassify/getOptions`
- 源码：`app/controller/admin/system/config/ConfigClassify.php` :: `getOptions()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/config/classify/status/:id` — 配置分类修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.config.ConfigClassify/switchStatus`
- 源码：`app/controller/admin/system/config/ConfigClassify.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/config/classify/update/:id` — 配置分类编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.config.ConfigClassify/update`
- 源码：`app/controller/admin/system/config/ConfigClassify.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `pid` (query/body, 可选)
- `classify_name` (query/body, 可选)
- `classify_key` (query/body, 可选)
- `info` (query/body, 可选)
- `status` (query/body, 可选)
- `icon` (query/body, 可选)
- `sort` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/config/classify/update/table/:id` — 配置分类编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.config.ConfigClassify/updateTable`
- 源码：`app/controller/admin/system/config/ConfigClassify.php` :: `updateTable()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/config/others/group_buying` — 配置信息

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.config.ConfigOthers/getGroupBuying`
- 源码：`app/controller/admin/system/config/ConfigOthers.php` :: `getGroupBuying()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/config/others/group_buying` — 配置保存

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.config.ConfigOthers/setGroupBuying`
- 源码：`app/controller/admin/system/config/ConfigOthers.php` :: `setGroupBuying()`
- 请求参数：
- `ficti_status` (query/body, 可选)
- `group_buying_rate` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/config/others/update` — 配置保存

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.config.ConfigOthers/update`
- 源码：`app/controller/admin/system/config/ConfigOthers.php` :: `update()`
- 请求参数：
- `extension_status` (query/body, 可选)
- `extension_two_rate` (query/body, 可选)
- `extension_one_rate` (query/body, 可选)
- `extension_self` (query/body, 可选)
- `extension_limit` (query/body, 可选)
- `extension_limit_day` (query/body, 可选)
- `sys_extension_type` (query/body, 可选)
- `lock_brokerage_timer` (query/body, 可选)
- `max_bag_number` (query/body, 可选)
- `promoter_explain` (query/body, 可选)
- `user_extract_min` (query/body, 可选)
- `withdraw_type` (query/body, 可选)
- `promoter_type` (query/body, 可选)
- `promoter_low_money` (query/body, 可选)
- `extract_switch` (query/body, 可选)
- `extension_pop` (query/body, 可选)
- `transfer_scene_id` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/config/save/:key` — 编辑配置信息

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.config.ConfigValue/save`
- 源码：`app/controller/admin/system/config/ConfigValue.php` :: `save()`
- 请求参数：
- `key` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/config/setting/create` — 配置添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.config.Config/create`
- 源码：`app/controller/admin/system/config/Config.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/config/setting/create/table` — 配置添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.config.Config/createTable`
- 源码：`app/controller/admin/system/config/Config.php` :: `createTable()`
- 请求参数：
- `config_classify_id` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/config/setting/delete/:id` — 配置删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.config.Config/delete`
- 源码：`app/controller/admin/system/config/Config.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/config/setting/lst` — 配置列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.config.Config/lst`
- 源码：`app/controller/admin/system/config/Config.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `config_classify_id` (query/body, 可选)
- `user_type` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/config/setting/routine/config` — 小程序配置

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.config.Config/getRoutineConfig`
- 源码：`app/controller/admin/system/config/Config.php` :: `getRoutineConfig()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/config/setting/routine/downloadTemp` — 小程序下载

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.config.Config/downloadTemp`
- 源码：`app/controller/admin/system/config/Config.php` :: `downloadTemp()`
- 请求参数：
- `is_live` (query/body, 可选)
- `is_menu` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/config/setting/status/:id` — 配置修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.config.Config/switchStatus`
- 源码：`app/controller/admin/system/config/Config.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/config/setting/update/:id` — 配置编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.config.Config/update`
- 源码：`app/controller/admin/system/config/Config.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/config/setting/update/table/:id` — 配置编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.config.Config/updateTable`
- 源码：`app/controller/admin/system/config/Config.php` :: `updateTable()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/config/setting/update_name/:field` — 上传原名文件

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.config.Config/uploadAsName`
- 源码：`app/controller/admin/system/config/Config.php` :: `uploadAsName()`
- 请求参数：
- `field` (path, 必填) 路径参数
- 返回：data 对象字段: src | 外层: {status,message,data}

### `POST /sys/config/setting/upload_file/:field` — 上传文件

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.config.Config/upload`
- 源码：`app/controller/admin/system/config/Config.php` :: `upload()`
- 请求参数：
- `field` (path, 必填) 路径参数
- 返回：data 对象字段: src | 外层: {status,message,data}

### `GET /sys/config/setting/wechat/file/form` — 微信校验文件上传表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.config.Config/uploadWechatForm`
- 源码：`app/controller/admin/system/config/Config.php` :: `uploadWechatForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/config/setting/wechat_set` — 微信校验文件上传

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.config.Config/uploadWechatSet`
- 源码：`app/controller/admin/system/config/Config.php` :: `uploadWechatSet()`
- 请求参数：
- `wechat_chekc_file` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `sys/config_classify`

### `GET /sys/config_classify/:key` — 获取配置

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.config.Config/getConfig`
- 源码：`app/controller/admin/system/config/Config.php` :: `getConfig()`
- 请求参数：
- `key` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。


## `sys/copyright`

### `GET /sys/copyright/auth` — 获取授权信息

- 置信度：⛔ stale
- 说明：路由已登记，但 `app/controller/admin/Common.php` 中不存在方法 `authCopyright`（原项目死路由/加密扩展/版本差异）。开发勿实现为有效接口，除非核实。
- 处理器：`admin.Common/authCopyright`
- 源码：`app/controller/admin/Common.php` :: `authCopyright()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一响应包 {status,message,data?}；详见控制器实现 | 外层: {status,message,data}

### `GET /sys/copyright/get` — 获取去版权信息

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.Common/svaeCopyright`
- 源码：`app/controller/admin/Common.php` :: `svaeCopyright()`
- 请求参数：
- `copyright_context` (query/body, 可选)
- `copyright_image` (query/body, 可选)
- 返回：统一响应包 {status,message,data?}；详见控制器实现 | 外层: {status,message,data}

### `POST /sys/copyright/save` — 保存去版权信息

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.Common/svaeCopyright`
- 源码：`app/controller/admin/Common.php` :: `svaeCopyright()`
- 请求参数：
- `copyright_context` (query/body, 可选)
- `copyright_image` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `sys/data_screen`

### `GET /sys/data_screen/:key` — 数据大屏

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.Common/data_screen`
- 源码：`app/controller/admin/Common.php` :: `data_screen()`
- 请求参数：
- `key` (path, 必填) 路径参数
- `pid` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `sys/delivery`

### `GET /sys/delivery/belence` — 余额

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.delivery.DeliveryStation/getBalance`
- 源码：`app/controller/admin/delivery/DeliveryStation.php` :: `getBalance()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/delivery/config/form` — 配置表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.delivery.DeliveryStation/deliveryForm`
- 源码：`app/controller/admin/delivery/DeliveryStation.php` :: `deliveryForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/delivery/config/save` — 编辑配置

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.delivery.DeliveryStation/saveDeliveryConfig`
- 源码：`app/controller/admin/delivery/DeliveryStation.php` :: `saveDeliveryConfig()`
- 请求参数：
- `delivery_type` (query/body, 可选)
- `dada_app_key` (query/body, 可选)
- `dada_app_sercret` (query/body, 可选)
- `dada_source_id` (query/body, 可选)
- `uupt_appkey` (query/body, 可选)
- `uupt_app_id` (query/body, 可选)
- `uupt_open_id` (query/body, 可选)
- `delivery_status` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/delivery/order/detail/:id` — 配送详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.delivery.DeliveryOrder/detail`
- 源码：`app/controller/admin/delivery/DeliveryOrder.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/delivery/order/lst` — 配送记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.delivery.DeliveryOrder/lst`
- 源码：`app/controller/admin/delivery/DeliveryOrder.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `station_name` (query/body, 可选)
- `status` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `date` (query/body, 可选)
- `order_sn` (query/body, 可选)
- `station_type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/delivery/recharge` — 充值

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.delivery.DeliveryStation/getRecharge`
- 源码：`app/controller/admin/delivery/DeliveryStation.php` :: `getRecharge()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/delivery/station/detail/:id` — 门店详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.delivery.DeliveryStation/detail`
- 源码：`app/controller/admin/delivery/DeliveryStation.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/delivery/station/lst` — 门店列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.delivery.DeliveryStation/lst`
- 源码：`app/controller/admin/delivery/DeliveryStation.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `station_name` (query/body, 可选)
- `status` (query/body, 可选)
- `mer_id` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/delivery/station/options` — 门店筛选

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.delivery.DeliveryStation/options`
- 源码：`app/controller/admin/delivery/DeliveryStation.php` :: `options()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/delivery/station/payLst` — 充值记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.delivery.DeliveryStation/payLst`
- 源码：`app/controller/admin/delivery/DeliveryStation.php` :: `payLst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `mer_id` (query/body, 可选)
- `date` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/delivery/title` — 统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.delivery.DeliveryOrder/title`
- 源码：`app/controller/admin/delivery/DeliveryOrder.php` :: `title()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。


## `sys/discounts`

### `GET /sys/discounts/detail/:id` — 优惠套餐详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Discounts/detail`
- 源码：`app/controller/admin/store/Discounts.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/discounts/lst` — 优惠套餐列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Discounts/lst`
- 源码：`app/controller/admin/store/Discounts.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `store_name` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `title` (query/body, 可选)
- `status` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/discounts/status/:id` — 优惠套餐修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Discounts/switchStatus`
- 源码：`app/controller/admin/store/Discounts.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `sys/diy`

### `GET /sys/diy/categroy/:id/form` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.PageCategroy/updateForm`
- 源码：`app/controller/admin/system/diy/PageCategroy.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/diy/categroy/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.PageCategroy/create`
- 源码：`app/controller/admin/system/diy/PageCategroy.php` :: `create()`
- 请求参数：
- `pid` (query/body, 可选)
- `type` (query/body, 可选)
- `name` (query/body, 可选)
- `status` (query/body, 可选)
- `sort` (query/body, 可选)
- `is_mer` (query/body, 可选)
- `level` (query/body, 可选) 默认 3
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/diy/categroy/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.PageCategroy/delete`
- 源码：`app/controller/admin/system/diy/PageCategroy.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/diy/categroy/form` — 添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.PageCategroy/createForm`
- 源码：`app/controller/admin/system/diy/PageCategroy.php` :: `createForm()`
- 请求参数：
- `type` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/diy/categroy/lst` — 列表 

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.PageCategroy/lst`
- 源码：`app/controller/admin/system/diy/PageCategroy.php` :: `lst()`
- 请求参数：
- `type` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/diy/categroy/options` — 列表 

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.PageCategroy/options`
- 源码：`app/controller/admin/system/diy/PageCategroy.php` :: `options()`
- 请求参数：
- `type` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/diy/categroy/status/:id` — 编辑状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.PageCategroy/switchStatus`
- 源码：`app/controller/admin/system/diy/PageCategroy.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/diy/categroy/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.PageCategroy/update`
- 源码：`app/controller/admin/system/diy/PageCategroy.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `pid` (query/body, 可选)
- `type` (query/body, 可选)
- `name` (query/body, 可选)
- `status` (query/body, 可选)
- `sort` (query/body, 可选)
- `is_mer` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/diy/copy/:id` — 复制

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.Diy/copy`
- 源码：`app/controller/admin/system/diy/Diy.php` :: `copy()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/diy/create/:id` — 添加/编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.Diy/saveData`
- 源码：`app/controller/admin/system/diy/Diy.php` :: `saveData()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `name` (query/body, 可选) 默认 ''
- `title` (query/body, 可选) 默认 ''
- `type` (query/body, 可选) 默认 '1'
- `cover_image` (query/body, 可选) 默认 ''
- `is_show` (query/body, 可选) 默认 0
- `is_bg_color` (query/body, 可选) 默认 0
- `is_bg_pic` (query/body, 可选) 默认 0
- `bg_tab_val` (query/body, 可选) 默认 0
- `color_picker` (query/body, 可选) 默认 ''
- `bg_pic` (query/body, 可选) 默认 ''
- `is_diy` (query/body, 可选) 默认 1
- `is_default` (query/body, 可选) 默认 0
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `DELETE /sys/diy/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.Diy/del`
- 源码：`app/controller/admin/system/diy/Diy.php` :: `del()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/diy/detail/:id` — 详情 

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.Diy/getInfo`
- 源码：`app/controller/admin/system/diy/Diy.php` :: `getInfo()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/diy/fab/create/:id` — 保存悬浮按钮

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.Diy/saveData`
- 源码：`app/controller/admin/system/diy/Diy.php` :: `saveData()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `name` (query/body, 可选) 默认 ''
- `title` (query/body, 可选) 默认 ''
- `type` (query/body, 可选) 默认 '1'
- `cover_image` (query/body, 可选) 默认 ''
- `is_show` (query/body, 可选) 默认 0
- `is_bg_color` (query/body, 可选) 默认 0
- `is_bg_pic` (query/body, 可选) 默认 0
- `bg_tab_val` (query/body, 可选) 默认 0
- `color_picker` (query/body, 可选) 默认 ''
- `bg_pic` (query/body, 可选) 默认 ''
- `is_diy` (query/body, 可选) 默认 1
- `is_default` (query/body, 可选) 默认 0
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/diy/fab/info` — 悬浮按钮信息

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.Diy/fabInfo`
- 源码：`app/controller/admin/system/diy/Diy.php` :: `fabInfo()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/diy/get_product_detail` — 商品详情 

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.Diy/get_product_detail`
- 源码：`app/controller/admin/system/diy/Diy.php` :: `get_product_detail()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/diy/get_theme/:key` — 可视化详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.VisualConfig/getTheme`
- 源码：`app/controller/admin/system/diy/VisualConfig.php` :: `getTheme()`
- 请求参数：
- `key` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/diy/get_theme_key` — 可视化列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.VisualConfig/getThemeKey`
- 源码：`app/controller/admin/system/diy/VisualConfig.php` :: `getThemeKey()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/diy/link/:id/form` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.PageLink/updateForm`
- 源码：`app/controller/admin/system/diy/PageLink.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `type` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/diy/link/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.PageLink/create`
- 源码：`app/controller/admin/system/diy/PageLink.php` :: `create()`
- 请求参数：
- `cate_id` (query/body, 可选)
- `name` (query/body, 可选)
- `url` (query/body, 可选)
- `example` (query/body, 可选)
- `status` (query/body, 可选)
- `sort` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/diy/link/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.PageLink/delete`
- 源码：`app/controller/admin/system/diy/PageLink.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/diy/link/form` — 添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.PageLink/createForm`
- 源码：`app/controller/admin/system/diy/PageLink.php` :: `createForm()`
- 请求参数：
- `type` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/diy/link/getLinks/:id` — 列表

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

### `GET /sys/diy/link/list` — lst

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

### `GET /sys/diy/link/lst` — 列表

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

### `POST /sys/diy/link/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.PageLink/switchStatus`
- 源码：`app/controller/admin/system/diy/PageLink.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/diy/link/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.PageLink/update`
- 源码：`app/controller/admin/system/diy/PageLink.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `cate_id` (query/body, 可选)
- `name` (query/body, 可选)
- `url` (query/body, 可选)
- `example` (query/body, 可选)
- `status` (query/body, 可选)
- `sort` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/diy/lst` — 列表 

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.Diy/lst`
- 源码：`app/controller/admin/system/diy/Diy.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `type` (query/body, 可选)
- `is_diy` (query/body, 可选)
- `name` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/diy/mer_categroy/:id/form` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.PageCategroy/updateForm`
- 源码：`app/controller/admin/system/diy/PageCategroy.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/diy/mer_categroy/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.PageCategroy/create`
- 源码：`app/controller/admin/system/diy/PageCategroy.php` :: `create()`
- 请求参数：
- `pid` (query/body, 可选)
- `type` (query/body, 可选)
- `name` (query/body, 可选)
- `status` (query/body, 可选)
- `sort` (query/body, 可选)
- `is_mer` (query/body, 可选)
- `level` (query/body, 可选) 默认 3
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/diy/mer_categroy/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.PageCategroy/delete`
- 源码：`app/controller/admin/system/diy/PageCategroy.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/diy/mer_categroy/form` — 添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.PageCategroy/createForm`
- 源码：`app/controller/admin/system/diy/PageCategroy.php` :: `createForm()`
- 请求参数：
- `type` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/diy/mer_categroy/lst` — 列表 

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.PageCategroy/lst`
- 源码：`app/controller/admin/system/diy/PageCategroy.php` :: `lst()`
- 请求参数：
- `type` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/diy/mer_categroy/status/:id` — 编辑状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.PageCategroy/switchStatus`
- 源码：`app/controller/admin/system/diy/PageCategroy.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/diy/mer_categroy/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.PageCategroy/update`
- 源码：`app/controller/admin/system/diy/PageCategroy.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `pid` (query/body, 可选)
- `type` (query/body, 可选)
- `name` (query/body, 可选)
- `status` (query/body, 可选)
- `sort` (query/body, 可选)
- `is_mer` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/diy/mer_link/:id/form` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.PageLink/updateForm`
- 源码：`app/controller/admin/system/diy/PageLink.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `type` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/diy/mer_link/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.PageLink/create`
- 源码：`app/controller/admin/system/diy/PageLink.php` :: `create()`
- 请求参数：
- `cate_id` (query/body, 可选)
- `name` (query/body, 可选)
- `url` (query/body, 可选)
- `example` (query/body, 可选)
- `status` (query/body, 可选)
- `sort` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/diy/mer_link/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.PageLink/delete`
- 源码：`app/controller/admin/system/diy/PageLink.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/diy/mer_link/form` — 添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.PageLink/createForm`
- 源码：`app/controller/admin/system/diy/PageLink.php` :: `createForm()`
- 请求参数：
- `type` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/diy/mer_link/lst` — 列表

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

### `POST /sys/diy/mer_link/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.PageLink/switchStatus`
- 源码：`app/controller/admin/system/diy/PageLink.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/diy/mer_link/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.PageLink/update`
- 源码：`app/controller/admin/system/diy/PageLink.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `cate_id` (query/body, 可选)
- `name` (query/body, 可选)
- `url` (query/body, 可选)
- `example` (query/body, 可选)
- `status` (query/body, 可选)
- `sort` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/diy/product/lst` — 商品列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.Diy/productLst`
- 源码：`app/controller/admin/system/diy/Diy.php` :: `productLst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `store_name` (query/body, 可选) 默认 ''
- `order` (query/body, 可选) 默认 'star'
- `cate_pid` (query/body, 可选) 默认 0
- `star` (query/body, 可选) 默认 ''
- `cate_id` (query/body, 可选) 默认 ''
- `product_type` (query/body, 可选)
- `mer_cate_id` (query/body, 可选)
- `product_ids` (query/body, 可选)
- `store_type_id` (query/body, 可选)
- `store_label_id` (query/body, 可选)
- `delivery_type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/diy/productCategory/create/:id` — 保存商品分类

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.Diy/saveData`
- 源码：`app/controller/admin/system/diy/Diy.php` :: `saveData()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `name` (query/body, 可选) 默认 ''
- `title` (query/body, 可选) 默认 ''
- `type` (query/body, 可选) 默认 '1'
- `cover_image` (query/body, 可选) 默认 ''
- `is_show` (query/body, 可选) 默认 0
- `is_bg_color` (query/body, 可选) 默认 0
- `is_bg_pic` (query/body, 可选) 默认 0
- `bg_tab_val` (query/body, 可选) 默认 0
- `color_picker` (query/body, 可选) 默认 ''
- `bg_pic` (query/body, 可选) 默认 ''
- `is_diy` (query/body, 可选) 默认 1
- `is_default` (query/body, 可选) 默认 0
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/diy/productCategory/info` — 商品分类信息

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.Diy/productCategoryInfo`
- 源码：`app/controller/admin/system/diy/Diy.php` :: `productCategoryInfo()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/diy/recovery/:id` — 重置

- 置信度：✅ high
- 说明：已按源码方法名校正
- 处理器：`admin.system.diy.Diy/Recovery`
- 源码：`app/controller/admin/system/diy/Diy.php` :: `Recovery()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/diy/save_product_detail` — 商品详情保存 

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.Diy/save_product_detail`
- 源码：`app/controller/admin/system/diy/Diy.php` :: `save_product_detail()`
- 请求参数：
- `name` (query/body, 可选) 默认 ''
- `title` (query/body, 可选) 默认 ''
- `product_detail_diy` (query/body, 可选) 默认 [
- `type` (query/body, 可选) 默认 '1'
- `cover_image` (query/body, 可选) 默认 ''
- `is_show` (query/body, 可选) 默认 0
- `is_bg_color` (query/body, 可选) 默认 0
- `is_bg_pic` (query/body, 可选) 默认 0
- `bg_tab_val` (query/body, 可选) 默认 0
- `color_picker` (query/body, 可选) 默认 ''
- `bg_pic` (query/body, 可选) 默认 ''
- `is_diy` (query/body, 可选) 默认 1
- `is_default` (query/body, 可选) 默认 0
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/diy/select` — select

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.Diy/select`
- 源码：`app/controller/admin/system/diy/Diy.php` :: `select()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/diy/set_default_data/:id` — 设置默认

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.Diy/setDefaultData`
- 源码：`app/controller/admin/system/diy/Diy.php` :: `setDefaultData()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/diy/set_theme/:key` — 可视化保存

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.VisualConfig/setTheme`
- 源码：`app/controller/admin/system/diy/VisualConfig.php` :: `setTheme()`
- 请求参数：
- `key` (path, 必填) 路径参数
- `config` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/diy/status/:id` — 使用模板

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.Diy/setStatus`
- 源码：`app/controller/admin/system/diy/Diy.php` :: `setStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/diy/store_street` — 店铺街装修

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.VisualConfig/storeStreet`
- 源码：`app/controller/admin/system/diy/VisualConfig.php` :: `storeStreet()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/diy/store_street` — 店铺街装修

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.VisualConfig/setStoreStreet`
- 源码：`app/controller/admin/system/diy/VisualConfig.php` :: `setStoreStreet()`
- 请求参数：
- `mer_location` (query/body, 可选)
- `store_street_theme` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/diy/user_index` — 个人中心装修

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.VisualConfig/userIndex`
- 源码：`app/controller/admin/system/diy/VisualConfig.php` :: `userIndex()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/diy/user_index` — 个人中心装修

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.VisualConfig/setUserIndex`
- 源码：`app/controller/admin/system/diy/VisualConfig.php` :: `setUserIndex()`
- 请求参数：
- `my_banner` (query/body, 可选)
- `my_menus` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。


## `sys/excel`

### `GET /sys/excel/download/:id` — 下载

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.Excel/downloadExpress`
- 源码：`app/controller/merchant/store/Excel.php` :: `downloadExpress()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一响应包 {status,message,data?}；详见控制器实现 | 外层: {status,message,data}

### `GET /sys/excel/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.Excel/lst`
- 源码：`app/controller/merchant/store/Excel.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `type` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/excel/type` — 类型

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.Excel/type`
- 源码：`app/controller/merchant/store/Excel.php` :: `type()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `sys/financial`

### `GET /sys/financial/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.financial.Financial/detail`
- 源码：`app/controller/admin/system/financial/Financial.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/financial/export` — 导出

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.financial.Financial/export`
- 源码：`app/controller/admin/system/financial/Financial.php` :: `export()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `date` (query/body, 可选)
- `status` (query/body, 可选)
- `financial_type` (query/body, 可选)
- `financial_status` (query/body, 可选)
- `keyword` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `mer_id` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/financial/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.financial.Financial/lst`
- 源码：`app/controller/admin/system/financial/Financial.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `date` (query/body, 可选)
- `status` (query/body, 可选)
- `financial_type` (query/body, 可选)
- `financial_status` (query/body, 可选)
- `keyword` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `mer_id` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/financial/mark/:id` — 备注

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.financial.Financial/mark`
- 源码：`app/controller/admin/system/financial/Financial.php` :: `mark()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `admin_mark` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/financial/mark/:id/form` — 备注表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.financial.Financial/markForm`
- 源码：`app/controller/admin/system/financial/Financial.php` :: `markForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/financial/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.financial.Financial/switchStatus`
- 源码：`app/controller/admin/system/financial/Financial.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选) 默认 0
- `refusal` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/financial/title` — 统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.financial.Financial/title`
- 源码：`app/controller/admin/system/financial/Financial.php` :: `title()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/financial/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.financial.Financial/update`
- 源码：`app/controller/admin/system/financial/Financial.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `image` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。


## `sys/financial_record`

### `GET /sys/financial_record/count` — 统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.FinancialRecord/title`
- 源码：`app/controller/admin/system/merchant/FinancialRecord.php` :: `title()`
- 请求参数：
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/financial_record/detail/:type` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.FinancialRecord/detail`
- 源码：`app/controller/admin/system/merchant/FinancialRecord.php` :: `detail()`
- 请求参数：
- `type` (path, 必填) 路径参数
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/financial_record/detail_export/:type` — 导出

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

### `GET /sys/financial_record/export` — 导出

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

### `GET /sys/financial_record/list` — 列表

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

### `GET /sys/financial_record/lst` — 列表

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

### `GET /sys/financial_record/mer_detail/:type` — 商户财务详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.FinancialRecord/merDetail`
- 源码：`app/controller/admin/system/merchant/FinancialRecord.php` :: `merDetail()`
- 请求参数：
- `type` (path, 必填) 路径参数
- `date` (query/body, 可选)
- `mer_id` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/financial_record/mer_excel/:type` — 商户财务导出

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.FinancialRecord/merExportDetail`
- 源码：`app/controller/admin/system/merchant/FinancialRecord.php` :: `merExportDetail()`
- 请求参数：
- `type` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `date` (query/body, 可选)
- `mer_id` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/financial_record/mer_list/:id` — 商户统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.FinancialRecord/merAcountsList`
- 源码：`app/controller/admin/system/merchant/FinancialRecord.php` :: `merAcountsList()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `type` (query/body, 可选) 默认 1
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/financial_record/mer_lst` — 商户列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.FinancialRecord/merchantFinancial`
- 源码：`app/controller/admin/system/merchant/FinancialRecord.php` :: `merchantFinancial()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `mer_id` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/financial_record/mer_title/:id` — 商户财务头部统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.FinancialRecord/merAcountsTitle`
- 源码：`app/controller/admin/system/merchant/FinancialRecord.php` :: `merAcountsTitle()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/financial_record/title` — 统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.FinancialRecord/getTitle`
- 源码：`app/controller/admin/system/merchant/FinancialRecord.php` :: `getTitle()`
- 请求参数：
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `sys/group`

### `POST /sys/group/create` — 组合数据配置添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.groupData.Group/create`
- 源码：`app/controller/admin/system/groupData/Group.php` :: `create()`
- 请求参数：
- `group_name` (query/body, 可选)
- `group_info` (query/body, 可选)
- `user_type` (query/body, 可选)
- `group_key` (query/body, 可选)
- `fields` (query/body, 可选) 默认 [
- `sort` (query/body, 可选) 默认 0
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/group/create/table` — 组合数据配置添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.groupData.Group/createTable`
- 源码：`app/controller/admin/system/groupData/Group.php` :: `createTable()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/group/data/create/:groupId` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.groupData.GroupData/create`
- 源码：`app/controller/admin/system/groupData/GroupData.php` :: `create()`
- 请求参数：
- `groupId` (path, 必填) 路径参数
- `sort` (query/body, 可选) 默认 0
- `status` (query/body, 可选) 默认 0
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/group/data/create/table/:groupId` — 添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.groupData.GroupData/createTable`
- 源码：`app/controller/admin/system/groupData/GroupData.php` :: `createTable()`
- 请求参数：
- `groupId` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `DELETE /sys/group/data/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.groupData.GroupData/delete`
- 源码：`app/controller/admin/system/groupData/GroupData.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/group/data/detail/:id` — baseDetail

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.groupData.GroupData/baseDetail`
- 源码：`app/controller/admin/system/groupData/GroupData.php` :: `baseDetail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/group/data/lst/:groupId` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.groupData.GroupData/lst`
- 源码：`app/controller/admin/system/groupData/GroupData.php` :: `lst()`
- 请求参数：
- `groupId` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/group/data/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.groupData.GroupData/changeStatus`
- 源码：`app/controller/admin/system/groupData/GroupData.php` :: `changeStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/group/data/update/:groupId/:id` — 编辑

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

### `GET /sys/group/data/update/table/:groupId/:id` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.groupData.GroupData/updateTable`
- 源码：`app/controller/admin/system/groupData/GroupData.php` :: `updateTable()`
- 请求参数：
- `groupId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/group/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.groupData.Group/get`
- 源码：`app/controller/admin/system/groupData/Group.php` :: `get()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/group/lst` — 组合数据配置列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.groupData.Group/lst`
- 源码：`app/controller/admin/system/groupData/Group.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/group/update/:id` — 组合数据配置编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.groupData.Group/update`
- 源码：`app/controller/admin/system/groupData/Group.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `group_name` (query/body, 可选)
- `group_info` (query/body, 可选)
- `user_type` (query/body, 可选)
- `group_key` (query/body, 可选)
- `fields` (query/body, 可选) 默认 [
- `sort` (query/body, 可选) 默认 0
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/group/update/table/:id` — 组合数据配置编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.groupData.Group/updateTable`
- 源码：`app/controller/admin/system/groupData/Group.php` :: `updateTable()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `sys/guarantee`

### `POST /sys/guarantee/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Guarantee/create`
- 源码：`app/controller/admin/store/Guarantee.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/guarantee/create/form` — 添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Guarantee/createForm`
- 源码：`app/controller/admin/store/Guarantee.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/guarantee/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Guarantee/delete`
- 源码：`app/controller/admin/store/Guarantee.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/guarantee/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Guarantee/detail`
- 源码：`app/controller/admin/store/Guarantee.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/guarantee/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Guarantee/lst`
- 源码：`app/controller/admin/store/Guarantee.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `date` (query/body, 可选)
- `keyword` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/guarantee/sort/:id` — 排序

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Guarantee/sort`
- 源码：`app/controller/admin/store/Guarantee.php` :: `sort()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `sort` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/guarantee/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Guarantee/switchStatus`
- 源码：`app/controller/admin/store/Guarantee.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/guarantee/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Guarantee/update`
- 源码：`app/controller/admin/store/Guarantee.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/guarantee/update/:id/form` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Guarantee/updateForm`
- 源码：`app/controller/admin/store/Guarantee.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `sys/info`

### `GET /sys/info` — getSystemInfo

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.Common/getSystemInfo`
- 源码：`app/controller/admin/Common.php` :: `getSystemInfo()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `sys/login`

### `POST /sys/login` — login

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.admin.Login/login`
- 源码：`app/controller/admin/system/admin/Login.php` :: `login()`
- 请求参数：
- `account` (query/body, 可选)
- `password` (query/body, 可选)
- `code` (query/body, 可选)
- `key` (query/body, 可选)
- `captchaType` (query/body, 可选) 默认 ''
- `captchaVerification` (query/body, 可选) 默认 ''
- `token` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `sys/login_config`

### `GET /sys/login_config` — loginConfig

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.Common/loginConfig`
- 源码：`app/controller/admin/Common.php` :: `loginConfig()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `sys/logout`

### `GET /sys/logout` — logout

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.admin.Login/logout`
- 源码：`app/controller/admin/system/admin/Login.php` :: `logout()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `sys/margin`

### `GET /sys/margin/list/:id` — 扣费记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.MerchantMargin/getMarginLst`
- 源码：`app/controller/admin/system/merchant/MerchantMargin.php` :: `getMarginLst()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/margin/local/:id` — 扣除保证金表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.MerchantMargin/localMarginSet`
- 源码：`app/controller/admin/system/merchant/MerchantMargin.php` :: `localMarginSet()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `number` (query/body, 可选)
- `mark` (query/body, 可选)
- `status` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/margin/local/:id/form` — 扣除保证金表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.MerchantMargin/localMarginForm`
- 源码：`app/controller/admin/system/merchant/MerchantMargin.php` :: `localMarginForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/margin/lst` — 缴纳记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.MerchantMargin/lst`
- 源码：`app/controller/admin/system/merchant/MerchantMargin.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `date` (query/body, 可选)
- `keyword` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `category_id` (query/body, 可选)
- `type_id` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/margin/make_up` — 待缴列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.Merchant/makeUpMarginLst`
- 源码：`app/controller/admin/system/merchant/Merchant.php` :: `makeUpMarginLst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `date` (query/body, 可选)
- `status` (query/body, 可选)
- `statusTag` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `category_id` (query/body, 可选)
- `type_id` (query/body, 可选)
- `order` (query/body, 可选) 默认 'create_time'
- `is_best` (query/body, 可选)
- `offline_switch` (query/body, 可选)
- `region_id` (query/body, 可选)
- `business_id` (query/body, 可选)
- `mer_state` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/margin/refund/lst` — 退款申请列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.financial.Financial/getMarginLst`
- 源码：`app/controller/admin/system/financial/Financial.php` :: `getMarginLst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `date` (query/body, 可选)
- `status` (query/body, 可选)
- `keyword` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `type_id` (query/body, 可选)
- `category_id` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/margin/refund/mark/:id` — 备注

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.financial.Financial/mark`
- 源码：`app/controller/admin/system/financial/Financial.php` :: `mark()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `admin_mark` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/margin/refund/mark/:id/form` — 备注表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.financial.Financial/markMarginForm`
- 源码：`app/controller/admin/system/financial/Financial.php` :: `markMarginForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/margin/refund/show/:id` — 退款申请详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.financial.Financial/refundShow`
- 源码：`app/controller/admin/system/financial/Financial.php` :: `refundShow()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/margin/refund/status/:id` — 审核

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.financial.Financial/switchStatus`
- 源码：`app/controller/admin/system/financial/Financial.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选) 默认 0
- `refusal` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/margin/refund/status/:id/form` — 审核表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.financial.Financial/statusForm`
- 源码：`app/controller/admin/system/financial/Financial.php` :: `statusForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/margin/set` — 扣除保证金

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.MerchantMargin/setMargin`
- 源码：`app/controller/admin/system/merchant/MerchantMargin.php` :: `setMargin()`
- 请求参数：
- `mer_id` (query/body, 可选)
- `number` (query/body, 可选)
- `type` (query/body, 可选) 默认 'mer_margin'
- `mark` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/margin/set/:id/form` — 扣除保证金表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.MerchantMargin/setMarginForm`
- 源码：`app/controller/admin/system/merchant/MerchantMargin.php` :: `setMarginForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。


## `sys/marketing`

### `GET /sys/marketing/spu/lst` — markLst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.marketing.StoreAtmosphere/markLst`
- 源码：`app/controller/admin/store/marketing/StoreAtmosphere.php` :: `markLst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `cate_id` (query/body, 可选)
- `cate_pid` (query/body, 可选)
- `brand_id` (query/body, 可选)
- `product_type` (query/body, 可选)
- `spu_ids` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `not_type` (query/body, 可选) 默认 1
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `sys/member`

### `POST /sys/member/interests/create` — 会员权益添加

- 置信度：✅ high
- 说明：控制器方法已校验存在；参数由 checkParams/Validate 自动补全
- 处理器：`admin.user..MemberInterests/create`
- 源码：`app/controller/admin/user/MemberInterests.php` :: `create()`
- 请求参数：
- `brokerage_level` (body, 可选) 来自 checkParams/Validate（自动补全）
- `name` (body, 可选) 来自 checkParams/Validate（自动补全）
- `info` (body, 可选) 来自 checkParams/Validate（自动补全）
- `pic` (body, 可选) 来自 checkParams/Validate（自动补全）
- `type` (body, 可选) 来自 checkParams/Validate（自动补全）
- `has_type` (body, 可选) 来自 checkParams/Validate（自动补全）
- `link` (body, 可选) 来自 checkParams/Validate（自动补全）
- `value` (body, 可选) 来自 checkParams/Validate（自动补全）
- `on_pic` (body, 可选) 来自 checkParams/Validate（自动补全）
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/member/interests/create/form` — 会员权益添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user..MemberInterests/createForm`
- 源码：`app/controller/admin/user/MemberInterests.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `DELETE /sys/member/interests/delete/:id` — 会员权益删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user..MemberInterests/delete`
- 源码：`app/controller/admin/user/MemberInterests.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/member/interests/detail/:id` — 会员权益详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user..MemberInterests/detail`
- 源码：`app/controller/admin/user/MemberInterests.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/member/interests/lst` — 会员权益

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user..MemberInterests/getLst`
- 源码：`app/controller/admin/user/MemberInterests.php` :: `getLst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `name` (query/body, 可选)
- `type` (query/body, 可选) 默认 $this->repository::TYPE_FREE
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/member/interests/options` — 会员权益筛选

- 置信度：⛔ stale
- 说明：路由已登记，但 `app/controller/admin/user/MemberInterests.php` 中不存在方法 `options`（原项目死路由/加密扩展/版本差异）。开发勿实现为有效接口，除非核实。
- 处理器：`admin.user..MemberInterests/options`
- 源码：`app/controller/admin/user/MemberInterests.php` :: `options()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/member/interests/update/:id` — 会员权益编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user..MemberInterests/update`
- 源码：`app/controller/admin/user/MemberInterests.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/member/interests/update/:id/form` — 会员权益编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user..MemberInterests/updateForm`
- 源码：`app/controller/admin/user/MemberInterests.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。


## `sys/menus`

### `GET /sys/menus` — menus

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.auth.Menu/menus`
- 源码：`app/controller/admin/system/auth/Menu.php` :: `menus()`
- 请求参数：
- `circle_id` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `sys/mer_diy`

### `GET /sys/mer_diy/copy/:id` — 复制

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.Diy/copy`
- 源码：`app/controller/admin/system/diy/Diy.php` :: `copy()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/mer_diy/create/:id` — 添加/编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.Diy/saveData`
- 源码：`app/controller/admin/system/diy/Diy.php` :: `saveData()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `name` (query/body, 可选) 默认 ''
- `title` (query/body, 可选) 默认 ''
- `type` (query/body, 可选) 默认 '1'
- `cover_image` (query/body, 可选) 默认 ''
- `is_show` (query/body, 可选) 默认 0
- `is_bg_color` (query/body, 可选) 默认 0
- `is_bg_pic` (query/body, 可选) 默认 0
- `bg_tab_val` (query/body, 可选) 默认 0
- `color_picker` (query/body, 可选) 默认 ''
- `bg_pic` (query/body, 可选) 默认 ''
- `is_diy` (query/body, 可选) 默认 1
- `is_default` (query/body, 可选) 默认 0
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `DELETE /sys/mer_diy/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.Diy/del`
- 源码：`app/controller/admin/system/diy/Diy.php` :: `del()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/mer_diy/detail/:id` — 详情 

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.Diy/getInfo`
- 源码：`app/controller/admin/system/diy/Diy.php` :: `getInfo()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/mer_diy/lst` — 列表 

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.Diy/lst`
- 源码：`app/controller/admin/system/diy/Diy.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `type` (query/body, 可选)
- `is_diy` (query/body, 可选)
- `name` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/mer_diy/recovery/:id` — 重置

- 置信度：✅ high
- 说明：已按源码方法名校正
- 处理器：`admin.system.diy.Diy/Recovery`
- 源码：`app/controller/admin/system/diy/Diy.php` :: `Recovery()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/mer_diy/scope/:id` — 保存适用范围

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.Diy/getScope`
- 源码：`app/controller/admin/system/diy/Diy.php` :: `getScope()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/mer_diy/scope/:id` — 保存适用范围

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.Diy/setScope`
- 源码：`app/controller/admin/system/diy/Diy.php` :: `setScope()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `scope_type` (query/body, 可选)
- `scope_value` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/mer_diy/set_default_data/:id` — 设置默认

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.Diy/setDefaultData`
- 源码：`app/controller/admin/system/diy/Diy.php` :: `setDefaultData()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。


## `sys/merchant`

### `DELETE /sys/merchant/intention/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.MerchantIntention/delete`
- 源码：`app/controller/admin/system/merchant/MerchantIntention.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/merchant/intention/excel` — excel

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.MerchantIntention/excel`
- 源码：`app/controller/admin/system/merchant/MerchantIntention.php` :: `excel()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `mer_name` (query/body, 可选)
- `status` (query/body, 可选)
- `date` (query/body, 可选)
- `keyword` (query/body, 可选)
- `mer_intention_id` (query/body, 可选)
- `category_id` (query/body, 可选)
- `type_id` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/merchant/intention/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.MerchantIntention/lst`
- 源码：`app/controller/admin/system/merchant/MerchantIntention.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `mer_name` (query/body, 可选)
- `status` (query/body, 可选)
- `date` (query/body, 可选)
- `keyword` (query/body, 可选)
- `mer_intention_id` (query/body, 可选)
- `category_id` (query/body, 可选)
- `type_id` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /sys/merchant/intention/mark/:id` — 备注

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.MerchantIntention/mark`
- 源码：`app/controller/admin/system/merchant/MerchantIntention.php` :: `mark()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `mark` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/merchant/intention/mark/:id/form` — 备注

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.MerchantIntention/form`
- 源码：`app/controller/admin/system/merchant/MerchantIntention.php` :: `form()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/merchant/intention/status/:id` — 审核

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.MerchantIntention/switchStatus`
- 源码：`app/controller/admin/system/merchant/MerchantIntention.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- `fail_msg` (query/body, 可选)
- `create_mer` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/merchant/intention/status/:id/form` — 申请店铺

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.MerchantIntention/statusForm`
- 源码：`app/controller/admin/system/merchant/MerchantIntention.php` :: `statusForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/merchant/menu/create` — 商户菜单/权限添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.auth.Menu/create`
- 源码：`app/controller/admin/system/auth/Menu.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/merchant/menu/create/form` — 商户菜单/权限添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.auth.Menu/createForm`
- 源码：`app/controller/admin/system/auth/Menu.php` :: `createForm()`
- 请求参数：
- `is_agent` (query/body, 可选) 默认 0
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/merchant/menu/delete/:id` — 商户菜单/权限删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.auth.Menu/delete`
- 源码：`app/controller/admin/system/auth/Menu.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/merchant/menu/lst` — 商户菜单/权限列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.auth.Menu/getList`
- 源码：`app/controller/admin/system/auth/Menu.php` :: `getList()`
- 请求参数：
- `is_agent` (query/body, 可选) 默认 0
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/merchant/menu/update/:id` — 商户菜单/权限编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.auth.Menu/update`
- 源码：`app/controller/admin/system/auth/Menu.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/merchant/menu/update/form/:id` — 商户菜单/权限编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.auth.Menu/updateForm`
- 源码：`app/controller/admin/system/auth/Menu.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `is_agent` (query/body, 可选) 默认 0
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/merchant/mer_auth` — mer_auth

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.MerchantType/mer_auth`
- 源码：`app/controller/admin/system/merchant/MerchantType.php` :: `mer_auth()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/merchant/type/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.MerchantType/create`
- 源码：`app/controller/admin/system/merchant/MerchantType.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/merchant/type/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.MerchantType/delete`
- 源码：`app/controller/admin/system/merchant/MerchantType.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/merchant/type/detail/:id` — 备注

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.MerchantType/detail`
- 源码：`app/controller/admin/system/merchant/MerchantType.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/merchant/type/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.MerchantType/lst`
- 源码：`app/controller/admin/system/merchant/MerchantType.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /sys/merchant/type/mark/:id` — 备注

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.MerchantType/markForm`
- 源码：`app/controller/admin/system/merchant/MerchantType.php` :: `markForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/merchant/type/mark/:id` — 备注

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.MerchantType/mark`
- 源码：`app/controller/admin/system/merchant/MerchantType.php` :: `mark()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `mark` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/merchant/type/mer_auth` — 权限

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.MerchantType/mer_auth`
- 源码：`app/controller/admin/system/merchant/MerchantType.php` :: `mer_auth()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/merchant/type/options` — 筛选

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.MerchantType/options`
- 源码：`app/controller/admin/system/merchant/MerchantType.php` :: `options()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/merchant/type/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.MerchantType/update`
- 源码：`app/controller/admin/system/merchant/MerchantType.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `sys/micro`

### `POST /sys/micro/create/:id` — 添加/编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.Diy/saveData`
- 源码：`app/controller/admin/system/diy/Diy.php` :: `saveData()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `name` (query/body, 可选) 默认 ''
- `title` (query/body, 可选) 默认 ''
- `type` (query/body, 可选) 默认 '1'
- `cover_image` (query/body, 可选) 默认 ''
- `is_show` (query/body, 可选) 默认 0
- `is_bg_color` (query/body, 可选) 默认 0
- `is_bg_pic` (query/body, 可选) 默认 0
- `bg_tab_val` (query/body, 可选) 默认 0
- `color_picker` (query/body, 可选) 默认 ''
- `bg_pic` (query/body, 可选) 默认 ''
- `is_diy` (query/body, 可选) 默认 1
- `is_default` (query/body, 可选) 默认 0
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `DELETE /sys/micro/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.Diy/del`
- 源码：`app/controller/admin/system/diy/Diy.php` :: `del()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/micro/detail/:id` — 详情 

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.Diy/getInfo`
- 源码：`app/controller/admin/system/diy/Diy.php` :: `getInfo()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/micro/lst` — 列表 

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.diy.Diy/lst`
- 源码：`app/controller/admin/system/diy/Diy.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `type` (query/body, 可选)
- `is_diy` (query/body, 可选)
- `name` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/micro/recovery/:id` — 重置

- 置信度：❓ unresolved
- 说明：未能可靠映射到控制器，开发时勿直接照抄，需对照 route 源码
- 处理器：`Diy/recovery/`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。


## `sys/notice`

### `GET /sys/notice/config/change/:id/form` — 消息配置修改模板ID

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.notice.SystemNoticeConfig/getTemplateId`
- 源码：`app/controller/admin/system/notice/SystemNoticeConfig.php` :: `getTemplateId()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/notice/config/change/:id/save` — 消息配置修改模板ID

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.notice.SystemNoticeConfig/setTemplateId`
- 源码：`app/controller/admin/system/notice/SystemNoticeConfig.php` :: `setTemplateId()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `sms_tempid` (query/body, 可选)
- `sms_ali_tempid` (query/body, 可选)
- `routine_tempid` (query/body, 可选)
- `wechat_tempid` (query/body, 可选)
- `notice_routine` (query/body, 可选) 默认 -1
- `notice_wechat` (query/body, 可选) 默认 -1
- `notice_sms` (query/body, 可选) 默认 -1
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/notice/config/create` — 消息配置添加

- 置信度：✅ high
- 说明：控制器方法已校验存在；参数由 checkParams/Validate 自动补全
- 处理器：`admin.system.notice.SystemNoticeConfig/create`
- 源码：`app/controller/admin/system/notice/SystemNoticeConfig.php` :: `create()`
- 请求参数：
- `notice_title` (body, 可选) 来自 checkParams/Validate（自动补全）
- `notice_key` (body, 可选) 来自 checkParams/Validate（自动补全）
- `notice_info` (body, 可选) 来自 checkParams/Validate（自动补全）
- `notice_sys` (body, 可选) 来自 checkParams/Validate（自动补全）
- `notice_wechat` (body, 可选) 来自 checkParams/Validate（自动补全）
- `notice_routine` (body, 可选) 来自 checkParams/Validate（自动补全）
- `notice_sms` (body, 可选) 来自 checkParams/Validate（自动补全）
- `type` (body, 可选) 来自 checkParams/Validate（自动补全）
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/notice/config/create/form` — 消息配置添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.notice.SystemNoticeConfig/createForm`
- 源码：`app/controller/admin/system/notice/SystemNoticeConfig.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/notice/config/delete/:id` — 消息配置删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.notice.SystemNoticeConfig/delete`
- 源码：`app/controller/admin/system/notice/SystemNoticeConfig.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/notice/config/detail/:id` — 消息配置详情

- 置信度：⛔ stale
- 说明：路由已登记，但 `app/controller/admin/system/notice/SystemNoticeConfig.php` 中不存在方法 `detail`（原项目死路由/加密扩展/版本差异）。开发勿实现为有效接口，除非核实。
- 处理器：`admin.system.notice.SystemNoticeConfig/detail`
- 源码：`app/controller/admin/system/notice/SystemNoticeConfig.php` :: `detail()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一响应包 {status,message,data?}；详见控制器实现 | 外层: {status,message,data}

### `GET /sys/notice/config/lst` — 消息配置列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.notice.SystemNoticeConfig/lst`
- 源码：`app/controller/admin/system/notice/SystemNoticeConfig.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /sys/notice/config/option` — 消息配置筛选

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.notice.SystemNoticeConfig/getOptions`
- 源码：`app/controller/admin/system/notice/SystemNoticeConfig.php` :: `getOptions()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/notice/config/status/:id` — 消息配置修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.notice.SystemNoticeConfig/switchStatus`
- 源码：`app/controller/admin/system/notice/SystemNoticeConfig.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- `key` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/notice/config/update/:id` — 消息配置编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.notice.SystemNoticeConfig/update`
- 源码：`app/controller/admin/system/notice/SystemNoticeConfig.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/notice/config/update/:id/form` — 消息配置编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.notice.SystemNoticeConfig/updateForm`
- 源码：`app/controller/admin/system/notice/SystemNoticeConfig.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/notice/create` — 系统公告发布

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.notice.SystemNotice/create`
- 源码：`app/controller/admin/system/notice/SystemNotice.php` :: `create()`
- 请求参数：
- `type` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `category_id` (query/body, 可选)
- `notice_title` (query/body, 可选)
- `notice_content` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/notice/delete/:id` — 系统公告删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.notice.SystemNotice/delete`
- 源码：`app/controller/admin/system/notice/SystemNotice.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/notice/detail/:id` — 系统公告详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.notice.SystemNotice/detail`
- 源码：`app/controller/admin/system/notice/SystemNotice.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/notice/lst` — 系统公告列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.notice.SystemNotice/lst`
- 源码：`app/controller/admin/system/notice/SystemNotice.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `date` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /sys/notice/switchStatus/:id` — 系统公告修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.notice.SystemNotice/switchStatus`
- 源码：`app/controller/admin/system/notice/SystemNotice.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/notice/update/:id` — 系统公告编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.notice.SystemNotice/update`
- 源码：`app/controller/admin/system/notice/SystemNotice.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `type` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `category_id` (query/body, 可选)
- `notice_title` (query/body, 可选)
- `notice_content` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `sys/order`

### `GET /sys/order/chart` — 头部统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.order.Order/chart`
- 源码：`app/controller/admin/order/Order.php` :: `chart()`
- 请求参数：
- `type` (query/body, 可选)
- `date` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `keywords` (query/body, 可选)
- `username` (query/body, 可选)
- `order_sn` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `activity_type` (query/body, 可选)
- `group_order_sn` (query/body, 可选)
- `store_name` (query/body, 可选)
- `spread_name` (query/body, 可选)
- `top_spread_name` (query/body, 可选)
- `filter_delivery` (query/body, 可选)
- `filter_product` (query/body, 可选)
- `nickname` (query/body, 可选)
- `uid` (query/body, 可选)
- `phone` (query/body, 可选)
- `real_name` (query/body, 可选)
- `delivery_name` (query/body, 可选)
- `delivery_phone` (query/body, 可选)
- `pay_type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/order/children/:id` — 关联订单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.order.Order/childrenList`
- 源码：`app/controller/admin/order/Order.php` :: `childrenList()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/order/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.order.Order/detail`
- 源码：`app/controller/admin/order/Order.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/order/excel` — 导出

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.order.Order/excel`
- 源码：`app/controller/admin/order/Order.php` :: `excel()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `type` (query/body, 可选)
- `date` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `keywords` (query/body, 可选)
- `status` (query/body, 可选)
- `username` (query/body, 可选)
- `order_sn` (query/body, 可选)
- `take_order` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `activity_type` (query/body, 可选)
- `group_order_sn` (query/body, 可选)
- `store_name` (query/body, 可选)
- `filter_delivery` (query/body, 可选)
- `filter_product` (query/body, 可选)
- `pay_type` (query/body, 可选)
- `uid` (query/body, 可选)
- `phone` (query/body, 可选)
- `real_name` (query/body, 可选)
- `delivery_name` (query/body, 可选)
- `delivery_phone` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/order/express/:id` — 快递查询

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.order.Order/express`
- 源码：`app/controller/admin/order/Order.php` :: `express()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/order/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.order.Order/getAllList`
- 源码：`app/controller/admin/order/Order.php` :: `getAllList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `type` (query/body, 可选)
- `date` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `keywords` (query/body, 可选)
- `status` (query/body, 可选)
- `username` (query/body, 可选)
- `order_sn` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `activity_type` (query/body, 可选)
- `group_order_sn` (query/body, 可选)
- `store_name` (query/body, 可选)
- `spread_name` (query/body, 可选)
- `top_spread_name` (query/body, 可选)
- `filter_delivery` (query/body, 可选)
- `filter_product` (query/body, 可选)
- `nickname` (query/body, 可选)
- `uid` (query/body, 可选)
- `phone` (query/body, 可选)
- `real_name` (query/body, 可选)
- `delivery_name` (query/body, 可选)
- `delivery_phone` (query/body, 可选)
- `pay_type` (query/body, 可选)
- `is_spread` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/order/refund/approve/:id` — 审核

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.order.RefundOrder/approve`
- 源码：`app/controller/admin/order/RefundOrder.php` :: `approve()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- `platform_mark` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/order/refund/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.order.RefundOrder/detail`
- 源码：`app/controller/admin/order/RefundOrder.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/order/refund/excel` — 导出

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.order.RefundOrder/excel`
- 源码：`app/controller/admin/order/RefundOrder.php` :: `excel()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `refund_order_sn` (query/body, 可选)
- `status` (query/body, 可选)
- `refund_type` (query/body, 可选)
- `date` (query/body, 可选)
- `order_sn` (query/body, 可选)
- `id` (query/body, 可选)
- `mer_id` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/order/refund/log/:id` — 日志

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.order.RefundOrder/log`
- 源码：`app/controller/admin/order/RefundOrder.php` :: `log()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `date` (query/body, 可选)
- `user_type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/order/refund/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.order.RefundOrder/getAllList`
- 源码：`app/controller/admin/order/RefundOrder.php` :: `getAllList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `refund_order_sn` (query/body, 可选)
- `status` (query/body, 可选)
- `refund_type` (query/body, 可选)
- `date` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `order_sn` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `uid` (query/body, 可选)
- `phone` (query/body, 可选)
- `real_name` (query/body, 可选)
- `nickname` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/order/status/:id` — 记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.order.Order/status`
- 源码：`app/controller/admin/order/Order.php` :: `status()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `date` (query/body, 可选)
- `user_type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/order/take_title` — 核销

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.order.Order/takeTitle`
- 源码：`app/controller/admin/order/Order.php` :: `takeTitle()`
- 请求参数：
- `date` (query/body, 可选)
- `order_sn` (query/body, 可选)
- `keywords` (query/body, 可选)
- `username` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `pay_type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/order/takechart` — 头部统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.order.Order/takeChart`
- 源码：`app/controller/admin/order/Order.php` :: `takeChart()`
- 请求参数：
- `date` (query/body, 可选)
- `order_sn` (query/body, 可选)
- `keywords` (query/body, 可选)
- `username` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `pay_type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/order/takelst` — 核销订单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.order.Order/getTakeList`
- 源码：`app/controller/admin/order/Order.php` :: `getTakeList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `date` (query/body, 可选)
- `order_sn` (query/body, 可选)
- `keywords` (query/body, 可选)
- `username` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `pay_type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/order/title` — 金额统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.order.Order/title`
- 源码：`app/controller/admin/order/Order.php` :: `title()`
- 请求参数：
- `type` (query/body, 可选)
- `date` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `keywords` (query/body, 可选)
- `status` (query/body, 可选)
- `username` (query/body, 可选)
- `order_sn` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `activity_type` (query/body, 可选)
- `filter_delivery` (query/body, 可选)
- `filter_product` (query/body, 可选)
- `nickname` (query/body, 可选)
- `uid` (query/body, 可选)
- `phone` (query/body, 可选)
- `real_name` (query/body, 可选)
- `is_spread` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。


## `sys/pay`

### `GET /sys/pay/auth` — payAuth

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.Common/payAuth`
- 源码：`app/controller/admin/Common.php` :: `payAuth()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `sys/points`

### `POST /sys/points/cate/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在；参数由 checkParams/Validate 自动补全
- 处理器：`admin.points.Category/create`
- 源码：`app/controller/admin/points/Category.php` :: `create()`
- 请求参数：
- `cate_name` (body, 可选) 来自 checkParams/Validate（自动补全）
- `is_show` (body, 可选) 来自 checkParams/Validate（自动补全）
- `pic` (body, 可选) 来自 checkParams/Validate（自动补全）
- `sort` (body, 可选) 来自 checkParams/Validate（自动补全）
- `type` (body, 可选) 来自 checkParams/Validate（自动补全）
- `pid` (body, 可选) 来自 checkParams/Validate（自动补全）
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/points/cate/create/form` — 添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.points.Category/createForm`
- 源码：`app/controller/admin/points/Category.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/points/cate/delete/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.points.Category/delete`
- 源码：`app/controller/admin/points/Category.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/points/cate/detail/:id` — 详情

- 置信度：⛔ stale
- 说明：路由已登记，但 `app/controller/admin/points/Category.php` 中不存在方法 `detail`（原项目死路由/加密扩展/版本差异）。开发勿实现为有效接口，除非核实。
- 处理器：`admin.points.Category/detail`
- 源码：`app/controller/admin/points/Category.php` :: `detail()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一响应包 {status,message,data?}；详见控制器实现 | 外层: {status,message,data}

### `GET /sys/points/cate/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.points.Category/lst`
- 源码：`app/controller/admin/points/Category.php` :: `lst()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/points/cate/select` — 筛选

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.points.Category/select`
- 源码：`app/controller/admin/points/Category.php` :: `select()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/points/cate/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.points.Category/switchStatus`
- 源码：`app/controller/admin/points/Category.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/points/cate/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.points.Category/update`
- 源码：`app/controller/admin/points/Category.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/points/cate/update/form/:id` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.points.Category/updateForm`
- 源码：`app/controller/admin/points/Category.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/points/order/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.points.Order/delete`
- 源码：`app/controller/admin/points/Order.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/points/order/delivery/:id` — 发货

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.points.Order/delivery`
- 源码：`app/controller/admin/points/Order.php` :: `delivery()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `is_split` (query/body, 可选)
- `split` (query/body, 可选) 默认 [
- `delivery_type` (query/body, 可选)
- `delivery_name` (query/body, 可选)
- `delivery_id` (query/body, 可选)
- `remark` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/points/order/delivery_batch` — 批量发货

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.points.Order/batchDelivery`
- 源码：`app/controller/admin/points/Order.php` :: `batchDelivery()`
- 请求参数：
- `order_id` (query/body, 可选)
- `delivery_id` (query/body, 可选)
- `delivery_type` (query/body, 可选)
- `delivery_name` (query/body, 可选)
- `remark` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/points/order/detail/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.points.Order/detail`
- 源码：`app/controller/admin/points/Order.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/points/order/excel` — 导出

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.points.Order/excel`
- 源码：`app/controller/admin/points/Order.php` :: `excel()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `date` (query/body, 可选)
- `status` (query/body, 可选)
- `order_sn` (query/body, 可选)
- `product_id` (query/body, 可选)
- 返回：统一响应包 {status,message,data?}；详见控制器实现 | 外层: {status,message,data}

### `GET /sys/points/order/express/:id` — 快递查询

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.points.Order/express`
- 源码：`app/controller/admin/points/Order.php` :: `express()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/points/order/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.points.Order/lst`
- 源码：`app/controller/admin/points/Order.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keywords` (query/body, 可选)
- `date` (query/body, 可选)
- `status` (query/body, 可选)
- `order_sn` (query/body, 可选)
- `product_id` (query/body, 可选)
- `nickname` (query/body, 可选)
- `phone` (query/body, 可选)
- `uid` (query/body, 可选)
- `store_name` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /sys/points/order/mark/:id` — 备注

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.points.Order/remark`
- 源码：`app/controller/admin/points/Order.php` :: `remark()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `remark` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/points/order/mark/:id/form` — 备注表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.points.Order/remarkForm`
- 源码：`app/controller/admin/points/Order.php` :: `remarkForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/points/order/status/:id` — 记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.points.Order/getStatus`
- 源码：`app/controller/admin/points/Order.php` :: `getStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `user_type` (query/body, 可选)
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/points/product/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在；参数由 checkParams/Validate 自动补全
- 处理器：`admin.points.Product/create`
- 源码：`app/controller/admin/points/Product.php` :: `create()`
- 请求参数：
- `is_used` (body, 可选) 来自 checkParams/Validate（自动补全）
- `is_hot` (body, 可选) 来自 checkParams/Validate（自动补全）
- `spec_type` (body, 可选) 来自 checkParams/Validate（自动补全）
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/points/product/delete/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.points.Product/delete`
- 源码：`app/controller/admin/points/Product.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/points/product/detail/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.points.Product/detail`
- 源码：`app/controller/admin/points/Product.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/points/product/get_attr_value/:id` — 获取规格

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.points.Product/isFormatAttr`
- 源码：`app/controller/admin/points/Product.php` :: `isFormatAttr()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `attrs` (query/body, 可选) 默认 [
- `items` (query/body, 可选) 默认 [
- `product_type` (query/body, 可选) 默认 0
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/points/product/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.points.Product/lst`
- 源码：`app/controller/admin/points/Product.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `cate_id` (query/body, 可选)
- `keyword` (query/body, 可选)
- `is_used` (query/body, 可选)
- `date` (query/body, 可选)
- `store_name` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /sys/points/product/preview` — 预览

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.points.Product/preview`
- 源码：`app/controller/admin/points/Product.php` :: `preview()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/points/product/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.points.Product/switchStatus`
- 源码：`app/controller/admin/points/Product.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/points/product/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.points.Product/update`
- 源码：`app/controller/admin/points/Product.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `sys/price_rule`

### `POST /sys/price_rule/create` — 添加价格说明

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.PriceRule/create`
- 源码：`app/controller/admin/store/PriceRule.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/price_rule/del/:id` — 删除价格说明

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.PriceRule/delete`
- 源码：`app/controller/admin/store/PriceRule.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/price_rule/lst` — 价格说明列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.PriceRule/lst`
- 源码：`app/controller/admin/store/PriceRule.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `cate_id` (query/body, 可选)
- `keyword` (query/body, 可选)
- `is_show` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /sys/price_rule/status/:id` — 价格说明修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.PriceRule/switchStatus`
- 源码：`app/controller/admin/store/PriceRule.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `is_show` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/price_rule/update/:id` — 修改价格说明

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.PriceRule/update`
- 源码：`app/controller/admin/store/PriceRule.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `sys/product`

### `POST /sys/product/label/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.ProductLabel/create`
- 源码：`app/controller/admin/store/ProductLabel.php` :: `create()`
- 请求参数：
- `label_name` (query/body, 可选)
- `status` (query/body, 可选)
- `sort` (query/body, 可选)
- `info` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/product/label/create/form` — 添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.ProductLabel/createForm`
- 源码：`app/controller/admin/store/ProductLabel.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/product/label/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.ProductLabel/delete`
- 源码：`app/controller/admin/store/ProductLabel.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/product/label/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.ProductLabel/detail`
- 源码：`app/controller/admin/store/ProductLabel.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/product/label/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.ProductLabel/lst`
- 源码：`app/controller/admin/store/ProductLabel.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `name` (query/body, 可选)
- `type` (query/body, 可选)
- `mer_id` (query/body, 可选) 默认 0
- `status` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/product/label/option` — 筛选

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.ProductLabel/getOptions`
- 源码：`app/controller/admin/store/ProductLabel.php` :: `getOptions()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/product/label/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.ProductLabel/switchWithStatus`
- 源码：`app/controller/admin/store/ProductLabel.php` :: `switchWithStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/product/label/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.ProductLabel/update`
- 源码：`app/controller/admin/store/ProductLabel.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `label_name` (query/body, 可选)
- `status` (query/body, 可选)
- `sort` (query/body, 可选)
- `info` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/product/label/update/:id/form` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.ProductLabel/updateForm`
- 源码：`app/controller/admin/store/ProductLabel.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `sys/profitsharing`

### `POST /sys/profitsharing/again/:id` — 重新分账

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.order.OrderProfitsharing/again`
- 源码：`app/controller/admin/order/OrderProfitsharing.php` :: `again()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/profitsharing/config` — 配置信息

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.config.ConfigOthers/getProfitsharing`
- 源码：`app/controller/admin/system/config/ConfigOthers.php` :: `getProfitsharing()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/profitsharing/config` — 配置保存

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.config.ConfigOthers/setProfitsharing`
- 源码：`app/controller/admin/system/config/ConfigOthers.php` :: `setProfitsharing()`
- 请求参数：
- `extract_maxmum_num` (query/body, 可选)
- `extract_minimum_line` (query/body, 可选)
- `extract_minimum_num` (query/body, 可选)
- `open_wx_combine` (query/body, 可选)
- `open_wx_sub_mch` (query/body, 可选)
- `mer_lock_time` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/profitsharing/export` — 导出

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

### `GET /sys/profitsharing/lst` — 列表

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


## `sys/receipt`

### `GET /sys/receipt/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.OrderReceipt/detail`
- 源码：`app/controller/merchant/store/order/OrderReceipt.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/receipt/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.order.OrderReceipt/getList`
- 源码：`app/controller/merchant/store/order/OrderReceipt.php` :: `getList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `status` (query/body, 可选)
- `date` (query/body, 可选)
- `receipt_sn` (query/body, 可选)
- `nickname` (query/body, 可选)
- `order_type` (query/body, 可选)
- `keyword` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `uid` (query/body, 可选)
- `phone` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}


## `sys/replace`

### `POST /sys/replace/image_host` — 替换素材域名

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.attachment.Attachment/replaceHost`
- 源码：`app/controller/admin/system/attachment/Attachment.php` :: `replaceHost()`
- 请求参数：
- `origin` (query/body, 可选)
- `replace` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `sys/safety`

### `POST /sys/safety/database/backups` — 备份

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.safety.Database/backups`
- 源码：`app/controller/admin/system/safety/Database.php` :: `backups()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/safety/database/delete` — 数据库备份删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.safety.Database/deleteFile`
- 源码：`app/controller/admin/system/safety/Database.php` :: `deleteFile()`
- 请求参数：
- `feilname` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/safety/database/detail/:name` — 数据库备份详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.safety.Database/detail`
- 源码：`app/controller/admin/system/safety/Database.php` :: `detail()`
- 请求参数：
- `name` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/safety/database/download/:feilname` — 数据库备份下载

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.safety.Database/downloadFile`
- 源码：`app/controller/admin/system/safety/Database.php` :: `downloadFile()`
- 请求参数：
- `feilname` (path, 必填) 路径参数
- 返回：失败时 status=400, message 为错误信息 | 外层: {status,message,data}

### `GET /sys/safety/database/fileList` — 数据库备份列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.safety.Database/fileList`
- 源码：`app/controller/admin/system/safety/Database.php` :: `fileList()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/safety/database/lst` — 数据库列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.safety.Database/lst`
- 源码：`app/controller/admin/system/safety/Database.php` :: `lst()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/safety/database/optimize` — 数据库优化

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.safety.Database/optimize`
- 源码：`app/controller/admin/system/safety/Database.php` :: `optimize()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/safety/database/repair` — 数据库维护

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.safety.Database/repair`
- 源码：`app/controller/admin/system/safety/Database.php` :: `repair()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `sys/seckill`

### `GET /sys/seckill/active/chart_order/:id` — 活动订单统计列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreSeckillActive/chart_order`
- 源码：`app/controller/admin/store/StoreSeckillActive.php` :: `chart_order()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `status` (query/body, 可选)
- `date` (query/body, 可选)
- `mer_id` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /sys/seckill/active/chart_panel/:id` — 活动统计数据面板

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreSeckillActive/chart_panel`
- 源码：`app/controller/admin/store/StoreSeckillActive.php` :: `chart_panel()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/seckill/active/chart_people/:id` — 活动参与人统计列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreSeckillActive/chart_people`
- 源码：`app/controller/admin/store/StoreSeckillActive.php` :: `chart_people()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `date` (query/body, 可选)
- `mer_id` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /sys/seckill/active/chart_product/:id` — 活动商品统计列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreSeckillActive/chart_product`
- 源码：`app/controller/admin/store/StoreSeckillActive.php` :: `chart_product()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `mer_id` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /sys/seckill/active/create` — 创建

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreSeckillActive/create`
- 源码：`app/controller/admin/store/StoreSeckillActive.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/seckill/active/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreSeckillActive/delete`
- 源码：`app/controller/admin/store/StoreSeckillActive.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/seckill/active/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreSeckillActive/detail`
- 源码：`app/controller/admin/store/StoreSeckillActive.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/seckill/active/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreSeckillActive/list`
- 源码：`app/controller/admin/store/StoreSeckillActive.php` :: `list()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `name` (query/body, 可选)
- `date` (query/body, 可选)
- `active_status` (query/body, 可选)
- `seckill_active_status` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /sys/seckill/active/select` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreSeckillActive/select`
- 源码：`app/controller/admin/store/StoreSeckillActive.php` :: `select()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/seckill/active/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreSeckillActive/update`
- 源码：`app/controller/admin/store/StoreSeckillActive.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/seckill/active/update_status/:id` — 编辑状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreSeckillActive/update_status`
- 源码：`app/controller/admin/store/StoreSeckillActive.php` :: `update_status()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/seckill/config/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在；参数由 checkParams/Validate 自动补全
- 处理器：`admin.store.StoreSeckill/create`
- 源码：`app/controller/admin/store/StoreSeckill.php` :: `create()`
- 请求参数：
- `start_time` (body, 可选) 来自 checkParams/Validate（自动补全）
- `end_time` (body, 可选) 来自 checkParams/Validate（自动补全）
- `status` (body, 可选) 来自 checkParams/Validate（自动补全）
- `title` (body, 可选) 来自 checkParams/Validate（自动补全）
- `pic` (body, 可选) 来自 checkParams/Validate（自动补全）
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/seckill/config/create/form` — 添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreSeckill/createForm`
- 源码：`app/controller/admin/store/StoreSeckill.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/seckill/config/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreSeckill/delete`
- 源码：`app/controller/admin/store/StoreSeckill.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/seckill/config/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreSeckill/lst`
- 源码：`app/controller/admin/store/StoreSeckill.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `title` (query/body, 可选)
- `status` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /sys/seckill/config/select` — 筛选

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreSeckill/select`
- 源码：`app/controller/admin/store/StoreSeckill.php` :: `select()`
- 请求参数：
- `active_id` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/seckill/config/status/:id` — 排序

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreSeckill/switchStatus`
- 源码：`app/controller/admin/store/StoreSeckill.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/seckill/config/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreSeckill/update`
- 源码：`app/controller/admin/store/StoreSeckill.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/seckill/config/update/:id/form` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreSeckill/updateForm`
- 源码：`app/controller/admin/store/StoreSeckill.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/seckill/product/change/:id` — 显示/隐藏

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductSeckill/changeUsed`
- 源码：`app/controller/admin/store/StoreProductSeckill.php` :: `changeUsed()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/seckill/product/delete/:id` — 加入回收站

- 置信度：⛔ stale
- 说明：路由已登记，但 `app/controller/admin/store/StoreProductSeckill.php` 中不存在方法 `delete`（原项目死路由/加密扩展/版本差异）。开发勿实现为有效接口，除非核实。
- 处理器：`admin.store.StoreProductSeckill/delete`
- 源码：`app/controller/admin/store/StoreProductSeckill.php` :: `delete()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一响应包 {status,message,data?}；详见控制器实现 | 外层: {status,message,data}

### `POST /sys/seckill/product/destory/:id` — 删除

- 置信度：⛔ stale
- 说明：路由已登记，但 `app/controller/admin/store/StoreProductSeckill.php` 中不存在方法 `destory`（原项目死路由/加密扩展/版本差异）。开发勿实现为有效接口，除非核实。
- 处理器：`admin.store.StoreProductSeckill/destory`
- 源码：`app/controller/admin/store/StoreProductSeckill.php` :: `destory()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一响应包 {status,message,data?}；详见控制器实现 | 外层: {status,message,data}

### `GET /sys/seckill/product/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductSeckill/detail`
- 源码：`app/controller/admin/store/StoreProductSeckill.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/seckill/product/labels/:id` — 设置标签

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductSeckill/setLabels`
- 源码：`app/controller/admin/store/StoreProductSeckill.php` :: `setLabels()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `sys_labels` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/seckill/product/list` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductSeckill/get_list`
- 源码：`app/controller/admin/store/StoreProductSeckill.php` :: `get_list()`
- 请求参数：
- `seckill_active_id` (query/body, 可选)
- `active_name` (query/body, 可选)
- `us_status` (query/body, 可选)
- `active_status` (query/body, 可选)
- `sys_labels` (query/body, 可选)
- `mer_cate_id` (query/body, 可选)
- `keyword` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `star` (query/body, 可选)
- `status` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/seckill/product/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductSeckill/get_page_list`
- 源码：`app/controller/admin/store/StoreProductSeckill.php` :: `get_page_list()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `seckill_active_id` (query/body, 可选)
- `active_name` (query/body, 可选)
- `us_status` (query/body, 可选)
- `active_status` (query/body, 可选)
- `sys_labels` (query/body, 可选)
- `mer_cate_id` (query/body, 可选)
- `keyword` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `star` (query/body, 可选)
- `status` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /sys/seckill/product/lst_filter` — 统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductSeckill/getStatusFilter`
- 源码：`app/controller/admin/store/StoreProductSeckill.php` :: `getStatusFilter()`
- 请求参数：
- `seckill_active_id` (query/body, 可选)
- `active_name` (query/body, 可选)
- `us_status` (query/body, 可选)
- `active_status` (query/body, 可选)
- `sys_labels` (query/body, 可选)
- `mer_cate_id` (query/body, 可选)
- `keyword` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `star` (query/body, 可选)
- `status` (query/body, 可选)
- `mer_id` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/seckill/product/mer_select` — 列表 

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductSeckill/lists`
- 源码：`app/controller/admin/store/StoreProductSeckill.php` :: `lists()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/seckill/product/status` — 审核

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductSeckill/switchStatus`
- 源码：`app/controller/admin/store/StoreProductSeckill.php` :: `switchStatus()`
- 请求参数：
- `status` (query/body, 可选)
- `refusal` (query/body, 可选)
- `id` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/seckill/product/status/:id/form` — 强制下架表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductSeckill/down_product_status_form`
- 源码：`app/controller/admin/store/StoreProductSeckill.php` :: `down_product_status_form()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/seckill/product/switchStatus/:id/form` — 审核表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductSeckill/get_switch_status_form`
- 源码：`app/controller/admin/store/StoreProductSeckill.php` :: `get_switch_status_form()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/seckill/product/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductSeckill/update`
- 源码：`app/controller/admin/store/StoreProductSeckill.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `sys/serve`

### `POST /sys/serve/captcha` — 验证码校验

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.serve.Login/checkCode`
- 源码：`app/controller/admin/system/serve/Login.php` :: `checkCode()`
- 请求参数：
- `phone` (query/body, 可选)
- `verify_code` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/serve/captcha/:phone` — 获取验证码

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.serve.Login/captcha`
- 源码：`app/controller/admin/system/serve/Login.php` :: `captcha()`
- 请求参数：
- `phone` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/serve/change_password` — 修改密码

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.serve.Serve/changePassword`
- 源码：`app/controller/admin/system/serve/Serve.php` :: `changePassword()`
- 请求参数：
- `phone` (query/body, 可选)
- `account` (query/body, 可选)
- `password` (query/body, 可选)
- `verify_code` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/serve/change_phone` — 修改手机号

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.serve.Serve/updatePhone`
- 源码：`app/controller/admin/system/serve/Serve.php` :: `updatePhone()`
- 请求参数：
- `phone` (query/body, 可选)
- `account` (query/body, 可选)
- `verify_code` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/serve/change_sign` — 修改签名

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.serve.Sms/changeSign`
- 源码：`app/controller/admin/system/serve/Sms.php` :: `changeSign()`
- 请求参数：
- `phone` (query/body, 可选)
- `sign` (query/body, 可选)
- `verify_code` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/serve/expr/dump_lst` — 使用记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.serve.Export/dumpLst`
- 源码：`app/controller/admin/system/serve/Export.php` :: `dumpLst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `date` (query/body, 可选)
- `mer_id` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/serve/expr/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.serve.Export/getExportAll`
- 源码：`app/controller/admin/system/serve/Export.php` :: `getExportAll()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/serve/expr/temps` — 模板

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.serve.Export/getExportTemp`
- 源码：`app/controller/admin/system/serve/Export.php` :: `getExportTemp()`
- 请求参数：
- `com` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/serve/login` — 登录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.serve.Login/login`
- 源码：`app/controller/admin/system/serve/Login.php` :: `login()`
- 请求参数：
- `account` (query/body, 可选)
- `password` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/serve/meal/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.serve.Config/create`
- 源码：`app/controller/admin/system/serve/Config.php` :: `create()`
- 请求参数：
- `name` (query/body, 可选)
- `price` (query/body, 可选)
- `num` (query/body, 可选)
- `type` (query/body, 可选)
- `status` (query/body, 可选)
- `sort` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/serve/meal/create/form` — 添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.serve.Config/createForm`
- 源码：`app/controller/admin/system/serve/Config.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/serve/meal/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.serve.Config/detail`
- 源码：`app/controller/admin/system/serve/Config.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `DELETE /sys/serve/meal/detele/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.serve.Config/detele`
- 源码：`app/controller/admin/system/serve/Config.php` :: `detele()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/serve/meal/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.serve.Config/lst`
- 源码：`app/controller/admin/system/serve/Config.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `status` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/serve/meal/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.serve.Config/switchStatus`
- 源码：`app/controller/admin/system/serve/Config.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/serve/meal/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.serve.Config/update`
- 源码：`app/controller/admin/system/serve/Config.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `name` (query/body, 可选)
- `price` (query/body, 可选)
- `num` (query/body, 可选)
- `type` (query/body, 可选)
- `status` (query/body, 可选)
- `sort` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/serve/meal/update/:id/form` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.serve.Config/updateForm`
- 源码：`app/controller/admin/system/serve/Config.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/serve/mealList/:type` — 套餐列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.serve.Serve/mealList`
- 源码：`app/controller/admin/system/serve/Serve.php` :: `mealList()`
- 请求参数：
- `type` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/serve/mer/lst` — 商户结余

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.serve.Serve/merLst`
- 源码：`app/controller/admin/system/serve/Serve.php` :: `merLst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `mer_id` (query/body, 可选)
- `date` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/serve/mer/paylst` — 商户购买记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.serve.Serve/merPaylst`
- 源码：`app/controller/admin/system/serve/Serve.php` :: `merPaylst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `type` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `date` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/serve/open` — 开通服务

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.serve.Serve/openServe`
- 源码：`app/controller/admin/system/serve/Serve.php` :: `openServe()`
- 请求参数：
- `type` (query/body, 可选)
- `sign` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/serve/paylst` — 购买记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.serve.Serve/paylst`
- 源码：`app/controller/admin/system/serve/Serve.php` :: `paylst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/serve/paymeal` — 购买套餐

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.serve.Serve/payMeal`
- 源码：`app/controller/admin/system/serve/Serve.php` :: `payMeal()`
- 请求参数：
- `meal_id` (query/body, 可选)
- `price` (query/body, 可选)
- `num` (query/body, 可选)
- `type` (query/body, 可选)
- `pay_type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/serve/record` — 使用记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.serve.Serve/getRecord`
- 源码：`app/controller/admin/system/serve/Serve.php` :: `getRecord()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/serve/register` — 注册

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.serve.Login/register`
- 源码：`app/controller/admin/system/serve/Login.php` :: `register()`
- 请求参数：
- `phone` (query/body, 可选)
- `account` (query/body, 可选)
- `password` (query/body, 可选)
- `verify_code` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/serve/sms/apply` — 申请模板

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.serve.Sms/apply`
- 源码：`app/controller/admin/system/serve/Sms.php` :: `apply()`
- 请求参数：
- `title` (query/body, 可选)
- `type` (query/body, 可选)
- `content` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/serve/sms/apply_record` — 申请记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.serve.Sms/applyRecord`
- 源码：`app/controller/admin/system/serve/Sms.php` :: `applyRecord()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `temp_type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/serve/sms/temps` — 短信模板

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.serve.Sms/temps`
- 源码：`app/controller/admin/system/serve/Sms.php` :: `temps()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/serve/us_lst` — 使用记录

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

### `GET /sys/serve/user/info` — 账号信息

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.serve.Serve/getUserInfo`
- 源码：`app/controller/admin/system/serve/Serve.php` :: `getUserInfo()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/serve/user/is_login` — 检测登录状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.serve.Serve/is_login`
- 源码：`app/controller/admin/system/serve/Serve.php` :: `is_login()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。


## `sys/service`

### `GET /sys/service/:id/:uid/lst` — 用户与客服聊天记录

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

### `GET /sys/service/:id/lst` — 用户与商户聊天记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.service.StoreService/getUserMsnByMerchant`
- 源码：`app/controller/merchant/store/service/StoreService.php` :: `getUserMsnByMerchant()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/service/:id/user` — 客服的全部用户 

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.service.StoreService/serviceUserList`
- 源码：`app/controller/merchant/store/service/StoreService.php` :: `serviceUserList()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/service/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在；参数由 checkParams/Validate 自动补全
- 处理器：`merchant.store.service.StoreService/create`
- 源码：`app/controller/merchant/store/service/StoreService.php` :: `create()`
- 请求参数：
- `uid` (body, 可选) 来自 checkParams/Validate（自动补全）
- `nickname` (body, 可选) 来自 checkParams/Validate（自动补全）
- `account` (body, 可选) 来自 checkParams/Validate（自动补全）
- `pwd` (body, 可选) 来自 checkParams/Validate（自动补全）
- `confirm_pwd` (body, 可选) 来自 checkParams/Validate（自动补全）
- `is_open` (body, 可选) 来自 checkParams/Validate（自动补全）
- `status` (body, 可选) 来自 checkParams/Validate（自动补全）
- `customer` (body, 可选) 来自 checkParams/Validate（自动补全）
- `is_verify` (body, 可选) 来自 checkParams/Validate（自动补全）
- `is_goods` (body, 可选) 来自 checkParams/Validate（自动补全）
- `notify` (body, 可选) 来自 checkParams/Validate（自动补全）
- `avatar` (body, 可选) 来自 checkParams/Validate（自动补全）
- `phone` (body, 可选) 来自 checkParams/Validate（自动补全）
- `sort` (body, 可选) 来自 checkParams/Validate（自动补全）
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/service/create/form` — 添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.service.StoreService/createForm`
- 源码：`app/controller/merchant/store/service/StoreService.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `DELETE /sys/service/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.service.StoreService/delete`
- 源码：`app/controller/merchant/store/service/StoreService.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/service/list` — 列表

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

### `POST /sys/service/login/:id` — 登录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.service.StoreService/login`
- 源码：`app/controller/merchant/store/service/StoreService.php` :: `login()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/service/mer/:id/user` — 客服的聊天用户列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.service.StoreService/merchantUserList`
- 源码：`app/controller/merchant/store/service/StoreService.php` :: `merchantUserList()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/service/reply/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在；参数由 checkParams/Validate 自动补全
- 处理器：`merchant.store.service.StoreServiceReply/create`
- 源码：`app/controller/merchant/store/service/StoreServiceReply.php` :: `create()`
- 请求参数：
- `keyword` (body, 可选) 来自 checkParams/Validate（自动补全）
- `status` (body, 可选) 来自 checkParams/Validate（自动补全）
- `content` (body, 可选) 来自 checkParams/Validate（自动补全）
- `type` (body, 可选) 来自 checkParams/Validate（自动补全）
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `DELETE /sys/service/reply/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.service.StoreServiceReply/delete`
- 源码：`app/controller/merchant/store/service/StoreServiceReply.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/service/reply/list` — 列表

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

### `POST /sys/service/reply/status/:id` — 切换状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.service.StoreServiceReply/changeStatus`
- 源码：`app/controller/merchant/store/service/StoreServiceReply.php` :: `changeStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/service/reply/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.service.StoreServiceReply/update`
- 源码：`app/controller/merchant/store/service/StoreServiceReply.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/service/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.service.StoreService/changeStatus`
- 源码：`app/controller/merchant/store/service/StoreService.php` :: `changeStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/service/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.service.StoreService/update`
- 源码：`app/controller/merchant/store/service/StoreService.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/service/update/form/:id` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.service.StoreService/updateForm`
- 源码：`app/controller/merchant/store/service/StoreService.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/service/user_lst` — 用户

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.service.StoreService/getUserList`
- 源码：`app/controller/merchant/store/service/StoreService.php` :: `getUserList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。


## `sys/sms`

### `POST /sys/sms/captcha` — smsCaptcha

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.sms..Sms/captcha`
- 源码：`app/controller/admin/system/sms/Sms.php` :: `captcha()`
- 请求参数：
- `phone` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/sms/change_password` — smsChangePassword

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.sms..Sms/changePassword`
- 源码：`app/controller/admin/system/sms/Sms.php` :: `changePassword()`
- 请求参数：
- `password` (query/body, 可选)
- `phone` (query/body, 可选)
- `code` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/sms/change_sign` — smsChangeSign

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.sms..Sms/changeSign`
- 源码：`app/controller/admin/system/sms/Sms.php` :: `changeSign()`
- 请求参数：
- `sign` (query/body, 可选)
- `phone` (query/body, 可选)
- `code` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/sms/config` — smsLogin

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.sms..Sms/save_basics`
- 源码：`app/controller/admin/system/sms/Sms.php` :: `save_basics()`
- 请求参数：
- `account` (query/body, 可选)
- `password` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/sms/data` — smsData

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.sms..Sms/data`
- 源码：`app/controller/admin/system/sms/Sms.php` :: `data()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/sms/is_login` — smsLogout

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.sms..Sms/is_login`
- 源码：`app/controller/admin/system/sms/Sms.php` :: `is_login()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/sms/logout` — smsLogout

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.sms..Sms/logout`
- 源码：`app/controller/admin/system/sms/Sms.php` :: `logout()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `ANY /sys/sms/notice` — notice

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.sms.SmsPay/notice`
- 源码：`app/controller/admin/system/sms/SmsPay.php` :: `notice()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一响应包 {status,message,data?}；详见控制器实现 | 外层: {status,message,data}

### `GET /sys/sms/number` — smsNumber

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.sms..SmsPay/number`
- 源码：`app/controller/admin/system/sms/SmsPay.php` :: `number()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/sms/pay_code` — smsPay

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.sms..SmsPay/pay`
- 源码：`app/controller/admin/system/sms/SmsPay.php` :: `pay()`
- 请求参数：
- `payType` (query/body, 可选) 默认 'weixin'
- `mealId` (query/body, 可选) 默认 0
- `price` (query/body, 可选) 默认 0
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/sms/price` — smsPrice

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.sms..SmsPay/price`
- 源码：`app/controller/admin/system/sms/SmsPay.php` :: `price()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/sms/public` — smsPublicTemplate

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.sms..SmsTemplate/public`
- 源码：`app/controller/admin/system/sms/SmsTemplate.php` :: `public()`
- 请求参数：
- `is_have` (query/body, 可选) 默认 ''
- `page` (query/body, 可选) 默认 1
- `limit` (query/body, 可选) 默认 20
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/sms/record` — smsRecord

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.sms..Sms/record`
- 源码：`app/controller/admin/system/sms/Sms.php` :: `record()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/sms/register` — smsSave

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.sms..Sms/save`
- 源码：`app/controller/admin/system/sms/Sms.php` :: `save()`
- 请求参数：
- `account` (query/body, 可选)
- `password` (query/body, 可选)
- `phone` (query/body, 可选)
- `code` (query/body, 可选)
- `url` (query/body, 可选)
- `sign` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/sms/temp` — smsTemplate

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.sms..SmsTemplate/template`
- 源码：`app/controller/admin/system/sms/SmsTemplate.php` :: `template()`
- 请求参数：
- `status` (query/body, 可选) 默认 ''
- `title` (query/body, 可选) 默认 ''
- `temp_type` (query/body, 可选) 默认 ''
- `page` (query/body, 可选) 默认 1
- `limit` (query/body, 可选) 默认 20
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/sms/temp` — smsCreate

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.sms..SmsTemplate/apply`
- 源码：`app/controller/admin/system/sms/SmsTemplate.php` :: `apply()`
- 请求参数：
- `title` (query/body, 可选)
- `content` (query/body, 可选)
- `type` (query/body, 可选) 默认 0
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/sms/temp/form` — smsCreateForm

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.sms..SmsTemplate/form`
- 源码：`app/controller/admin/system/sms/SmsTemplate.php` :: `form()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。


## `sys/spread`

### `GET /sys/spread/order/chart` — 头部统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.order.Order/spreadChart`
- 源码：`app/controller/admin/order/Order.php` :: `spreadChart()`
- 请求参数：
- `type` (query/body, 可选)
- `date` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `keywords` (query/body, 可选)
- `username` (query/body, 可选)
- `order_sn` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `activity_type` (query/body, 可选)
- `group_order_sn` (query/body, 可选)
- `store_name` (query/body, 可选)
- `spread_name` (query/body, 可选)
- `top_spread_name` (query/body, 可选)
- `filter_delivery` (query/body, 可选)
- `filter_product` (query/body, 可选)
- `nickname` (query/body, 可选)
- `uid` (query/body, 可选)
- `phone` (query/body, 可选)
- `real_name` (query/body, 可选)
- `delivery_name` (query/body, 可选)
- `delivery_phone` (query/body, 可选)
- `pay_type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/spread/order/children/:id` — 关联订单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.order.Order/childrenList`
- 源码：`app/controller/admin/order/Order.php` :: `childrenList()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/spread/order/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.order.Order/detail`
- 源码：`app/controller/admin/order/Order.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/spread/order/excel` — 导出

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.order.Order/excel`
- 源码：`app/controller/admin/order/Order.php` :: `excel()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `type` (query/body, 可选)
- `date` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `keywords` (query/body, 可选)
- `status` (query/body, 可选)
- `username` (query/body, 可选)
- `order_sn` (query/body, 可选)
- `take_order` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `activity_type` (query/body, 可选)
- `group_order_sn` (query/body, 可选)
- `store_name` (query/body, 可选)
- `filter_delivery` (query/body, 可选)
- `filter_product` (query/body, 可选)
- `pay_type` (query/body, 可选)
- `uid` (query/body, 可选)
- `phone` (query/body, 可选)
- `real_name` (query/body, 可选)
- `delivery_name` (query/body, 可选)
- `delivery_phone` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/spread/order/express/:id` — 快递查询

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.order.Order/express`
- 源码：`app/controller/admin/order/Order.php` :: `express()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/spread/order/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.order.Order/getAllList`
- 源码：`app/controller/admin/order/Order.php` :: `getAllList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `type` (query/body, 可选)
- `date` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `keywords` (query/body, 可选)
- `status` (query/body, 可选)
- `username` (query/body, 可选)
- `order_sn` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `activity_type` (query/body, 可选)
- `group_order_sn` (query/body, 可选)
- `store_name` (query/body, 可选)
- `spread_name` (query/body, 可选)
- `top_spread_name` (query/body, 可选)
- `filter_delivery` (query/body, 可选)
- `filter_product` (query/body, 可选)
- `nickname` (query/body, 可选)
- `uid` (query/body, 可选)
- `phone` (query/body, 可选)
- `real_name` (query/body, 可选)
- `delivery_name` (query/body, 可选)
- `delivery_phone` (query/body, 可选)
- `pay_type` (query/body, 可选)
- `is_spread` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/spread/order/status/:id` — 记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.order.Order/status`
- 源码：`app/controller/admin/order/Order.php` :: `status()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `date` (query/body, 可选)
- `user_type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/spread/order/title` — 金额统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.order.Order/title`
- 源码：`app/controller/admin/order/Order.php` :: `title()`
- 请求参数：
- `type` (query/body, 可选)
- `date` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `keywords` (query/body, 可选)
- `status` (query/body, 可选)
- `username` (query/body, 可选)
- `order_sn` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `activity_type` (query/body, 可选)
- `filter_delivery` (query/body, 可选)
- `filter_product` (query/body, 可选)
- `nickname` (query/body, 可选)
- `uid` (query/body, 可选)
- `phone` (query/body, 可选)
- `real_name` (query/body, 可选)
- `is_spread` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。


## `sys/statistics`

### `GET /sys/statistics/get_admin_count` — 未处理业务统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.Common/getAdminCount`
- 源码：`app/controller/admin/Common.php` :: `getAdminCount()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/statistics/get_admin_todo` — 待办事项

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.Common/getAdminTodo`
- 源码：`app/controller/admin/Common.php` :: `getAdminTodo()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/statistics/get_merchant_top` — 商户销量排行

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.Common/getMerchantTop`
- 源码：`app/controller/admin/Common.php` :: `getMerchantTop()`
- 请求参数：
- `date` (query/body, 可选)
- `type` (query/body, 可选)
- `sort` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/statistics/main` — 主要数据

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.Common/main`
- 源码：`app/controller/admin/Common.php` :: `main()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/statistics/merchant_rate` — 商户访问量

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.Common/merchantRate`
- 源码：`app/controller/admin/Common.php` :: `merchantRate()`
- 请求参数：
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/statistics/merchant_stock` — 商户销量

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.Common/merchantStock`
- 源码：`app/controller/admin/Common.php` :: `merchantStock()`
- 请求参数：
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/statistics/merchant_visit` — 商户销售额

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.Common/merchantVisit`
- 源码：`app/controller/admin/Common.php` :: `merchantVisit()`
- 请求参数：
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/statistics/order` — 当日订单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.Common/order`
- 源码：`app/controller/admin/Common.php` :: `order()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/statistics/order_num` — 当日订单数

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.Common/orderNum`
- 源码：`app/controller/admin/Common.php` :: `orderNum()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/statistics/order_user` — 当日支付人数

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.Common/orderUser`
- 源码：`app/controller/admin/Common.php` :: `orderUser()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/statistics/user` — 成交用户

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.Common/user`
- 源码：`app/controller/merchant/Common.php` :: `user()`
- 请求参数：
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/statistics/user_data` — 用户数据

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.Common/userData`
- 源码：`app/controller/admin/Common.php` :: `userData()`
- 请求参数：
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/statistics/user_rate` — 成交用户占比

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.Common/userRate`
- 源码：`app/controller/merchant/Common.php` :: `userRate()`
- 请求参数：
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `sys/store`

### `POST /sys/store/bag/change/:id` — 显示/隐藏

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProduct/changeUsed`
- 源码：`app/controller/admin/store/StoreProduct.php` :: `changeUsed()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/store/bag/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProduct/detail`
- 源码：`app/controller/admin/store/StoreProduct.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/store/bag/list` — 列表 

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProduct/lst`
- 源码：`app/controller/admin/store/StoreProduct.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `cate_id` (query/body, 可选)
- `keyword` (query/body, 可选)
- `mer_cate_id` (query/body, 可选)
- `pid` (query/body, 可选)
- `store_name` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `us_status` (query/body, 可选)
- `product_id` (query/body, 可选)
- `star` (query/body, 可选)
- `sys_labels` (query/body, 可选)
- `hot_type` (query/body, 可选)
- `svip_price_type` (query/body, 可选)
- `is_ficti` (query/body, 可选)
- `product_ids` (query/body, 可选)
- `form_id` (query/body, 可选)
- `cate_hot` (query/body, 可选)
- `brand_id` (query/body, 可选)
- `activity_label_ids` (query/body, 可选)
- `type` (query/body, 可选)
- `mer_id` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /sys/store/bag/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProduct/bagList`
- 源码：`app/controller/admin/store/StoreProduct.php` :: `bagList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `cate_id` (query/body, 可选)
- `keyword` (query/body, 可选)
- `mer_cate_id` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `us_status` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /sys/store/bag/lst_filter` — 统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProduct/getBagStatusFilter`
- 源码：`app/controller/admin/store/StoreProduct.php` :: `getBagStatusFilter()`
- 请求参数：
- `cate_id` (query/body, 可选)
- `keyword` (query/body, 可选)
- `mer_cate_id` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `us_status` (query/body, 可选)
- `mer_id` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/store/bag/mer_select` — 商户列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProduct/lists`
- 源码：`app/controller/admin/store/StoreProduct.php` :: `lists()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/store/bag/status` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProduct/switchStatus`
- 源码：`app/controller/admin/store/StoreProduct.php` :: `switchStatus()`
- 请求参数：
- `status` (query/body, 可选)
- `refusal` (query/body, 可选)
- `id` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/store/bag/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProduct/update`
- 源码：`app/controller/admin/store/StoreProduct.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/store/brand/category/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在；参数由 checkParams/Validate 自动补全
- 处理器：`admin.store.StoreBrandCategory/create`
- 源码：`app/controller/admin/store/StoreBrandCategory.php` :: `create()`
- 请求参数：
- `pid` (body, 可选) 来自 checkParams/Validate（自动补全）
- `cate_name` (body, 可选) 来自 checkParams/Validate（自动补全）
- `is_show` (body, 可选) 来自 checkParams/Validate（自动补全）
- `sort` (body, 可选) 来自 checkParams/Validate（自动补全）
- `data` (body, 可选) 来自 checkParams/Validate（自动补全）
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/store/brand/category/create/form` — 添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreBrandCategory/createForm`
- 源码：`app/controller/admin/store/StoreBrandCategory.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/store/brand/category/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreBrandCategory/delete`
- 源码：`app/controller/admin/store/StoreBrandCategory.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/store/brand/category/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreBrandCategory/detail`
- 源码：`app/controller/admin/store/StoreBrandCategory.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/store/brand/category/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreBrandCategory/lst`
- 源码：`app/controller/admin/store/StoreBrandCategory.php` :: `lst()`
- 请求参数：
- `status` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/store/brand/category/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreBrandCategory/switchStatus`
- 源码：`app/controller/admin/store/StoreBrandCategory.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/store/brand/category/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreBrandCategory/update`
- 源码：`app/controller/admin/store/StoreBrandCategory.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/store/brand/category/update/form/:id` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreBrandCategory/updateForm`
- 源码：`app/controller/admin/store/StoreBrandCategory.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/store/brand/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在；参数由 checkParams/Validate 自动补全
- 处理器：`admin.store.StoreBrand/create`
- 源码：`app/controller/admin/store/StoreBrand.php` :: `create()`
- 请求参数：
- `brand_category_id` (body, 可选) 来自 checkParams/Validate（自动补全）
- `brand_name` (body, 可选) 来自 checkParams/Validate（自动补全）
- `is_show` (body, 可选) 来自 checkParams/Validate（自动补全）
- `sort` (body, 可选) 来自 checkParams/Validate（自动补全）
- `pic` (body, 可选) 来自 checkParams/Validate（自动补全）
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/store/brand/create/form` — 添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreBrand/createForm`
- 源码：`app/controller/admin/store/StoreBrand.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/store/brand/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreBrand/delete`
- 源码：`app/controller/admin/store/StoreBrand.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/store/brand/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreBrand/lst`
- 源码：`app/controller/admin/store/StoreBrand.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `brand_category_id` (query/body, 可选)
- `brand_name` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /sys/store/brand/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreBrand/switchStatus`
- 源码：`app/controller/admin/store/StoreBrand.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/store/brand/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreBrand/update`
- 源码：`app/controller/admin/store/StoreBrand.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/store/brand/update/form/:id` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreBrand/updateForm`
- 源码：`app/controller/admin/store/StoreBrand.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/store/category/brandlist` — 品牌列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreCategory/BrandList`
- 源码：`app/controller/admin/store/StoreCategory.php` :: `BrandList()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/store/category/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在；参数由 checkParams/Validate 自动补全
- 处理器：`admin.store.StoreCategory/create`
- 源码：`app/controller/admin/store/StoreCategory.php` :: `create()`
- 请求参数：
- `pid` (body, 可选) 来自 checkParams/Validate（自动补全）
- `cate_name` (body, 可选) 来自 checkParams/Validate（自动补全）
- `is_show` (body, 可选) 来自 checkParams/Validate（自动补全）
- `pic` (body, 可选) 来自 checkParams/Validate（自动补全）
- `sort` (body, 可选) 来自 checkParams/Validate（自动补全）
- `data` (body, 可选) 来自 checkParams/Validate（自动补全）
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/store/category/create/form` — 添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreCategory/createForm`
- 源码：`app/controller/admin/store/StoreCategory.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/store/category/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreCategory/delete`
- 源码：`app/controller/admin/store/StoreCategory.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/store/category/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreCategory/detail`
- 源码：`app/controller/admin/store/StoreCategory.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/store/category/is_hot/:id` — 修改推荐

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreCategory/switchIsHot`
- 源码：`app/controller/admin/store/StoreCategory.php` :: `switchIsHot()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/store/category/list` — 筛选

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreCategory/getList`
- 源码：`app/controller/admin/store/StoreCategory.php` :: `getList()`
- 请求参数：
- `type` (query/body, 可选)
- `lv` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/store/category/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreCategory/lst`
- 源码：`app/controller/admin/store/StoreCategory.php` :: `lst()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/store/category/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreCategory/switchStatus`
- 源码：`app/controller/admin/store/StoreCategory.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/store/category/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreCategory/update`
- 源码：`app/controller/admin/store/StoreCategory.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/store/category/update/form/:id` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreCategory/updateForm`
- 源码：`app/controller/admin/store/StoreCategory.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/store/city/create` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在；参数由 checkParams/Validate 自动补全
- 处理器：`admin.store.CityArea/create`
- 源码：`app/controller/admin/store/CityArea.php` :: `create()`
- 请求参数：
- `parent_id` (body, 可选) 来自 checkParams/Validate（自动补全）
- `level` (body, 可选) 来自 checkParams/Validate（自动补全）
- `name` (body, 可选) 来自 checkParams/Validate（自动补全）
- `path` (body, 可选) 来自 checkParams/Validate（自动补全）
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/store/city/create/form/:id` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.CityArea/createForm`
- 源码：`app/controller/admin/store/CityArea.php` :: `createForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/store/city/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.CityArea/delete`
- 源码：`app/controller/admin/store/CityArea.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/store/city/lst/:id` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.CityArea/lst`
- 源码：`app/controller/admin/store/CityArea.php` :: `lst()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /sys/store/city/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.CityArea/update`
- 源码：`app/controller/admin/store/CityArea.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/store/city/update/:id/form` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.CityArea/updateForm`
- 源码：`app/controller/admin/store/CityArea.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/store/coupon/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在；参数由 checkParams/Validate 自动补全
- 处理器：`admin.store.Coupon/create`
- 源码：`app/controller/admin/store/Coupon.php` :: `create()`
- 请求参数：
- `use_type` (body, 可选) 来自 checkParams/Validate（自动补全）
- `title` (body, 可选) 来自 checkParams/Validate（自动补全）
- `coupon_price` (body, 可选) 来自 checkParams/Validate（自动补全）
- `use_min_price` (body, 可选) 来自 checkParams/Validate（自动补全）
- `coupon_type` (body, 可选) 来自 checkParams/Validate（自动补全）
- `coupon_time` (body, 可选) 来自 checkParams/Validate（自动补全）
- `use_start_time` (body, 可选) 来自 checkParams/Validate（自动补全）
- `sort` (body, 可选) 来自 checkParams/Validate（自动补全）
- `status` (body, 可选) 来自 checkParams/Validate（自动补全）
- `type` (body, 可选) 来自 checkParams/Validate（自动补全）
- `product_id` (body, 可选) 来自 checkParams/Validate（自动补全）
- `range_date` (body, 可选) 来自 checkParams/Validate（自动补全）
- `send_type` (body, 可选) 来自 checkParams/Validate（自动补全）
- `full_reduction` (body, 可选) 来自 checkParams/Validate（自动补全）
- `is_limited` (body, 可选) 来自 checkParams/Validate（自动补全）
- `is_timeout` (body, 可选) 来自 checkParams/Validate（自动补全）
- `total_count` (body, 可选) 来自 checkParams/Validate（自动补全）
- `cate_ids` (body, 可选) 来自 checkParams/Validate（自动补全）
- `mer_type` (body, 可选) 来自 checkParams/Validate（自动补全）
- `is_trader` (body, 可选) 来自 checkParams/Validate（自动补全）
- `category_id` (body, 可选) 来自 checkParams/Validate（自动补全）
- `type_id` (body, 可选) 来自 checkParams/Validate（自动补全）
- `mer_ids` (body, 可选) 来自 checkParams/Validate（自动补全）
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/store/coupon/create/form` — 添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Coupon/createForm`
- 源码：`app/controller/admin/store/Coupon.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/store/coupon/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Coupon/delete`
- 源码：`app/controller/admin/store/Coupon.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/store/coupon/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Coupon/detail`
- 源码：`app/controller/admin/store/Coupon.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/store/coupon/issue` — 使用记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Coupon/issue`
- 源码：`app/controller/admin/store/Coupon.php` :: `issue()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `username` (query/body, 可选)
- `coupon_id` (query/body, 可选)
- `coupon` (query/body, 可选)
- `coupon_type` (query/body, 可选)
- `status` (query/body, 可选)
- `type` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `is_mer` (query/body, 可选) 默认 1
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/store/coupon/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Coupon/lst`
- 源码：`app/controller/admin/store/Coupon.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `is_full_give` (query/body, 可选)
- `status` (query/body, 可选)
- `is_give_subscribe` (query/body, 可选)
- `coupon_name` (query/body, 可选)
- `mer_id` (query/body, 可选) 默认 null
- `is_trader` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /sys/store/coupon/platformLst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Coupon/platformLst`
- 源码：`app/controller/admin/store/Coupon.php` :: `platformLst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `is_full_give` (query/body, 可选)
- `status` (query/body, 可选)
- `is_give_subscribe` (query/body, 可选)
- `coupon_name` (query/body, 可选)
- `send_type` (query/body, 可选)
- `type` (query/body, 可选)
- `not_send_type` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/store/coupon/product/:id` — 商品列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Coupon/product`
- 源码：`app/controller/admin/store/Coupon.php` :: `product()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/store/coupon/send` — 发送优惠券

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Coupon/send`
- 源码：`app/controller/admin/store/Coupon.php` :: `send()`
- 请求参数：
- `coupon_id` (query/body, 可选)
- `mark` (query/body, 可选)
- `is_all` (query/body, 可选)
- `search` (query/body, 可选)
- `uid` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/store/coupon/send/lst` — 发送记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Coupon/sendLst`
- 源码：`app/controller/admin/store/Coupon.php` :: `sendLst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `date` (query/body, 可选)
- `coupon_type` (query/body, 可选)
- `coupon_name` (query/body, 可选)
- `status` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/store/coupon/show/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Coupon/detail`
- 源码：`app/controller/admin/store/Coupon.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/store/coupon/show_lst/:id` — 详情关联列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Coupon/showLst`
- 源码：`app/controller/admin/store/Coupon.php` :: `showLst()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/store/coupon/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Coupon/switchStatus`
- 源码：`app/controller/admin/store/Coupon.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/store/coupon/sys/clone/:id/form` — 复制表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Coupon/cloneForm`
- 源码：`app/controller/admin/store/Coupon.php` :: `cloneForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/store/coupon/sys/issue` — 使用记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Coupon/platformIssue`
- 源码：`app/controller/admin/store/Coupon.php` :: `platformIssue()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `username` (query/body, 可选)
- `coupon_id` (query/body, 可选)
- `coupon` (query/body, 可选)
- `status` (query/body, 可选)
- `coupon_type` (query/body, 可选)
- `type` (query/body, 可选)
- `send_id` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/store/coupon/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Coupon/update`
- 源码：`app/controller/admin/store/Coupon.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `title` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/store/coupon/update/:id/form` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Coupon/updateForm`
- 源码：`app/controller/admin/store/Coupon.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/store/express/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Express/delete`
- 源码：`app/controller/admin/store/Express.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/store/express/lst` — 列表

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

### `GET /sys/store/express/options` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Express/options`
- 源码：`app/controller/admin/store/Express.php` :: `options()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/store/express/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Express/switchStatus`
- 源码：`app/controller/admin/store/Express.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `is_show` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/store/express/sync` — 同步

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Express/syncAll`
- 源码：`app/controller/admin/store/Express.php` :: `syncAll()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/store/express/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Express/update`
- 源码：`app/controller/admin/store/Express.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `name` (query/body, 可选)
- `code` (query/body, 可选)
- `is_show` (query/body, 可选)
- `sort` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/store/express/update/form/:id` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.Express/updateForm`
- 源码：`app/controller/admin/store/Express.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/store/params/temp/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在；参数由 checkParams/Validate 自动补全
- 处理器：`admin.parameter.ParameterTemplate/create`
- 源码：`app/controller/admin/parameter/ParameterTemplate.php` :: `create()`
- 请求参数：
- `template_name` (body, 可选) 来自 checkParams/Validate（自动补全）
- `cate_ids` (body, 可选) 来自 checkParams/Validate（自动补全）
- `sort` (body, 可选) 来自 checkParams/Validate（自动补全）
- `params` (body, 可选) 来自 checkParams/Validate（自动补全）
- `delete_params` (body, 可选) 来自 checkParams/Validate（自动补全）
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/store/params/temp/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.parameter.ParameterTemplate/delete`
- 源码：`app/controller/admin/parameter/ParameterTemplate.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/store/params/temp/delete/value/:id` — 删除属性

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.parameter.ParameterTemplate/deleteValue`
- 源码：`app/controller/admin/parameter/ParameterTemplate.php` :: `deleteValue()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/store/params/temp/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.parameter.ParameterTemplate/detail`
- 源码：`app/controller/admin/parameter/ParameterTemplate.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/store/params/temp/lst` — 平台参数列表

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

### `GET /sys/store/params/temp/merlst` — 商户参数模板

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

### `POST /sys/store/params/temp/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.parameter.ParameterTemplate/update`
- 源码：`app/controller/admin/parameter/ParameterTemplate.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/store/product/assist/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductAssist/detail`
- 源码：`app/controller/admin/store/StoreProductAssist.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/store/product/assist/get/:id` — 编辑数据

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductAssist/get`
- 源码：`app/controller/admin/store/StoreProductAssist.php` :: `get()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/store/product/assist/is_show/:id` — 显示/隐藏

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductAssist/switchStatus`
- 源码：`app/controller/admin/store/StoreProductAssist.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/store/product/assist/labels/:id` — 设置标签

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductAssist/setLabels`
- 源码：`app/controller/admin/store/StoreProductAssist.php` :: `setLabels()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `sys_labels` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/store/product/assist/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductAssist/lst`
- 源码：`app/controller/admin/store/StoreProductAssist.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `product_status` (query/body, 可选)
- `keyword` (query/body, 可选)
- `status` (query/body, 可选)
- `type` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `us_status` (query/body, 可选)
- `star` (query/body, 可选)
- `product_assist_id` (query/body, 可选)
- `sys_labels` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/store/product/assist/set/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductAssistSet/detail`
- 源码：`app/controller/admin/store/StoreProductAssistSet.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/store/product/assist/set/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductAssistSet/lst`
- 源码：`app/controller/admin/store/StoreProductAssistSet.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `status` (query/body, 可选)
- `type` (query/body, 可选)
- `date` (query/body, 可选)
- `user_name` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `is_trader` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/store/product/assist/status` — 审核

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductAssist/switchAudit`
- 源码：`app/controller/admin/store/StoreProductAssist.php` :: `switchAudit()`
- 请求参数：
- `status` (query/body, 可选)
- `refusal` (query/body, 可选)
- `id` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/store/product/assist/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductAssist/update`
- 源码：`app/controller/admin/store/StoreProductAssist.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/store/product/batchCopyProduct` — 批量复制商品到店铺

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProduct/batchCopyProductsToStore`
- 源码：`app/controller/admin/store/StoreProduct.php` :: `batchCopyProductsToStore()`
- 请求参数：
- `product_ids` (query/body, 可选) 默认 [
- `store_id` (query/body, 可选) 默认 ''
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/store/product/batch_cate_hot` — 批量设置分类推荐

- 置信度：⛔ stale
- 说明：路由已登记，但 `app/controller/admin/store/StoreProduct.php` 中不存在方法 `batchCateHot`（原项目死路由/加密扩展/版本差异）。开发勿实现为有效接口，除非核实。
- 处理器：`admin.store.StoreProduct/batchCateHot`
- 源码：`app/controller/admin/store/StoreProduct.php` :: `batchCateHot()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一响应包 {status,message,data?}；详见控制器实现 | 外层: {status,message,data}

### `POST /sys/store/product/batch_hot` — 批量设置推荐

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProduct/batchHot`
- 源码：`app/controller/admin/store/StoreProduct.php` :: `batchHot()`
- 请求参数：
- `is_hot` (query/body, 可选) 默认 0
- `is_benefit` (query/body, 可选) 默认 0
- `is_best` (query/body, 可选) 默认 0
- `is_new` (query/body, 可选) 默认 0
- `ids` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/store/product/batch_labels` — 批量设置标签

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProduct/batchLabels`
- 源码：`app/controller/admin/store/StoreProduct.php` :: `batchLabels()`
- 请求参数：
- `sys_labels` (query/body, 可选)
- `ids` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/store/product/batch_status` — 批量上下架

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProduct/batchShow`
- 源码：`app/controller/admin/store/StoreProduct.php` :: `batchShow()`
- 请求参数：
- `ids` (query/body, 可选)
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/store/product/change/:id` — 显示/隐藏

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProduct/changeUsed`
- 源码：`app/controller/admin/store/StoreProduct.php` :: `changeUsed()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/store/product/check` — 分销状态变更商品检测

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProduct/checkProduct`
- 源码：`app/controller/admin/store/StoreProduct.php` :: `checkProduct()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/store/product/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProduct/detail`
- 源码：`app/controller/admin/store/StoreProduct.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/store/product/ficti/:id` — 虚拟销量

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProduct/addFicti`
- 源码：`app/controller/admin/store/StoreProduct.php` :: `addFicti()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `type` (query/body, 可选)
- `ficti` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/store/product/ficti/form/:id` — 虚拟销量表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProduct/addFictiForm`
- 源码：`app/controller/admin/store/StoreProduct.php` :: `addFictiForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/store/product/get_operate_list/:product_id` — 获取商品操作记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProduct/getOperateList`
- 源码：`app/controller/admin/store/StoreProduct.php` :: `getOperateList()`
- 请求参数：
- `product_id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `type` (query/body, 可选) 默认 ''
- `date` (query/body, 可选) 默认 ''
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/store/product/group/buying/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductGroupBuying/detail`
- 源码：`app/controller/admin/store/StoreProductGroupBuying.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /sys/store/product/group/buying/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductGroupBuying/lst`
- 源码：`app/controller/admin/store/StoreProductGroupBuying.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `status` (query/body, 可选)
- `date` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `user_name` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/store/product/group/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductGroup/detail`
- 源码：`app/controller/admin/store/StoreProductGroup.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/store/product/group/get/:id` — 编辑数据

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductGroup/get`
- 源码：`app/controller/admin/store/StoreProductGroup.php` :: `get()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/store/product/group/is_show/:id` — 显示/隐藏

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductGroup/switchStatus`
- 源码：`app/controller/admin/store/StoreProductGroup.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/store/product/group/labels/:id` — 设置标签

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductGroup/setLabels`
- 源码：`app/controller/admin/store/StoreProductGroup.php` :: `setLabels()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `sys_labels` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/store/product/group/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductGroup/lst`
- 源码：`app/controller/admin/store/StoreProductGroup.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `product_status` (query/body, 可选)
- `keyword` (query/body, 可选)
- `status` (query/body, 可选)
- `active_type` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `level` (query/body, 可选)
- `us_status` (query/body, 可选)
- `star` (query/body, 可选)
- `product_group_id` (query/body, 可选)
- `sys_labels` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/store/product/group/sort/:id` — 排序

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductGroup/updateSort`
- 源码：`app/controller/admin/store/StoreProductGroup.php` :: `updateSort()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `sort` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/store/product/group/status` — 审核

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductGroup/switchAudit`
- 源码：`app/controller/admin/store/StoreProductGroup.php` :: `switchAudit()`
- 请求参数：
- `status` (query/body, 可选)
- `refusal` (query/body, 可选)
- `id` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/store/product/group/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductGroup/update`
- 源码：`app/controller/admin/store/StoreProductGroup.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/store/product/labels/:id` — 设置标签

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProduct/setLabels`
- 源码：`app/controller/admin/store/StoreProduct.php` :: `setLabels()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `sys_labels` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/store/product/list` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProduct/lst`
- 源码：`app/controller/admin/store/StoreProduct.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `cate_id` (query/body, 可选)
- `keyword` (query/body, 可选)
- `mer_cate_id` (query/body, 可选)
- `pid` (query/body, 可选)
- `store_name` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `us_status` (query/body, 可选)
- `product_id` (query/body, 可选)
- `star` (query/body, 可选)
- `sys_labels` (query/body, 可选)
- `hot_type` (query/body, 可选)
- `svip_price_type` (query/body, 可选)
- `is_ficti` (query/body, 可选)
- `product_ids` (query/body, 可选)
- `form_id` (query/body, 可选)
- `cate_hot` (query/body, 可选)
- `brand_id` (query/body, 可选)
- `activity_label_ids` (query/body, 可选)
- `type` (query/body, 可选)
- `mer_id` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /sys/store/product/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProduct/lst`
- 源码：`app/controller/admin/store/StoreProduct.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `cate_id` (query/body, 可选)
- `keyword` (query/body, 可选)
- `mer_cate_id` (query/body, 可选)
- `pid` (query/body, 可选)
- `store_name` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `us_status` (query/body, 可选)
- `product_id` (query/body, 可选)
- `star` (query/body, 可选)
- `sys_labels` (query/body, 可选)
- `hot_type` (query/body, 可选)
- `svip_price_type` (query/body, 可选)
- `is_ficti` (query/body, 可选)
- `product_ids` (query/body, 可选)
- `form_id` (query/body, 可选)
- `cate_hot` (query/body, 可选)
- `brand_id` (query/body, 可选)
- `activity_label_ids` (query/body, 可选)
- `type` (query/body, 可选)
- `mer_id` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /sys/store/product/lst_filter` — 统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProduct/getStatusFilter`
- 源码：`app/controller/admin/store/StoreProduct.php` :: `getStatusFilter()`
- 请求参数：
- `cate_id` (query/body, 可选)
- `keyword` (query/body, 可选)
- `mer_cate_id` (query/body, 可选)
- `pid` (query/body, 可选)
- `store_name` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `us_status` (query/body, 可选)
- `product_id` (query/body, 可选)
- `sys_labels` (query/body, 可选)
- `hot_type` (query/body, 可选)
- `svip_price_type` (query/body, 可选)
- `is_ficti` (query/body, 可选)
- `product_ids` (query/body, 可选)
- `form_id` (query/body, 可选)
- `cate_hot` (query/body, 可选)
- `brand_id` (query/body, 可选)
- `activity_label_ids` (query/body, 可选)
- `mer_id` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/store/product/mer_select` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProduct/lists`
- 源码：`app/controller/admin/store/StoreProduct.php` :: `lists()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/store/product/presell/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductPresell/detail`
- 源码：`app/controller/admin/store/StoreProductPresell.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/store/product/presell/get/:id` — 编辑数据

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductPresell/get`
- 源码：`app/controller/admin/store/StoreProductPresell.php` :: `get()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/store/product/presell/is_show/:id` — 显示/隐藏

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductPresell/switchStatus`
- 源码：`app/controller/admin/store/StoreProductPresell.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/store/product/presell/labels/:id` — 设置标签

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductPresell/setLabels`
- 源码：`app/controller/admin/store/StoreProductPresell.php` :: `setLabels()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `sys_labels` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/store/product/presell/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductPresell/lst`
- 源码：`app/controller/admin/store/StoreProductPresell.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `product_status` (query/body, 可选)
- `keyword` (query/body, 可选)
- `status` (query/body, 可选)
- `type` (query/body, 可选)
- `presell_type` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `us_status` (query/body, 可选)
- `star` (query/body, 可选)
- `product_presell_id` (query/body, 可选)
- `sys_labels` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/store/product/presell/status` — 审核

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductPresell/switchAudit`
- 源码：`app/controller/admin/store/StoreProductPresell.php` :: `switchAudit()`
- 请求参数：
- `status` (query/body, 可选)
- `refusal` (query/body, 可选)
- `id` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/store/product/presell/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductPresell/update`
- 源码：`app/controller/admin/store/StoreProductPresell.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/store/product/self_lst` — 获取自营商品列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProduct/get_self_product_list`
- 源码：`app/controller/admin/store/StoreProduct.php` :: `get_self_product_list()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `sys_labels` (query/body, 可选)
- `us_status` (query/body, 可选)
- `cate_id` (query/body, 可选)
- `active_id` (query/body, 可选)
- `type_id` (query/body, 可选)
- `category_id` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `level_one_cate_ids` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/store/product/status` — 上下架

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProduct/switchStatus`
- 源码：`app/controller/admin/store/StoreProduct.php` :: `switchStatus()`
- 请求参数：
- `status` (query/body, 可选)
- `refusal` (query/body, 可选)
- `id` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/store/product/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProduct/update`
- 源码：`app/controller/admin/store/StoreProduct.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/store/reply/create` — 添加虚拟评论

- 置信度：✅ high
- 说明：控制器方法已校验存在；参数由 checkParams/Validate 自动补全
- 处理器：`admin.store.StoreProductReply/virtualReply`
- 源码：`app/controller/admin/store/StoreProductReply.php` :: `virtualReply()`
- 请求参数：
- `product_id` (body, 可选) 来自 checkParams/Validate（自动补全）
- `nickname` (body, 可选) 来自 checkParams/Validate（自动补全）
- `comment` (body, 可选) 来自 checkParams/Validate（自动补全）
- `sort` (body, 可选) 来自 checkParams/Validate（自动补全）
- `product_score` (body, 可选) 来自 checkParams/Validate（自动补全）
- `service_score` (body, 可选) 来自 checkParams/Validate（自动补全）
- `postage_score` (body, 可选) 来自 checkParams/Validate（自动补全）
- `avatar` (body, 可选) 来自 checkParams/Validate（自动补全）
- `pics` (body, 可选) 来自 checkParams/Validate（自动补全）
- `create_time` (body, 可选) 来自 checkParams/Validate（自动补全）
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/store/reply/create/form/:id?` — 添加虚拟评论表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductReply/virtualForm`
- 源码：`app/controller/admin/store/StoreProductReply.php` :: `virtualForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/store/reply/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductReply/delete`
- 源码：`app/controller/admin/store/StoreProductReply.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/store/reply/lst` — 列表

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

### `POST /sys/store/reply/sort/:id` — 排序

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.store.StoreProductReply/sort`
- 源码：`app/controller/admin/store/StoreProductReply.php` :: `sort()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `sort` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `sys/svip`

### `GET /sys/svip/interests/:id/form` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.MemberInterests/updateSvipForm`
- 源码：`app/controller/admin/user/MemberInterests.php` :: `updateSvipForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/svip/interests/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.MemberInterests/getSvipInterests`
- 源码：`app/controller/admin/user/MemberInterests.php` :: `getSvipInterests()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/svip/interests/status/:id` — 编辑状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.MemberInterests/switchWithStatus`
- 源码：`app/controller/admin/user/MemberInterests.php` :: `switchWithStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/svip/interests/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.MemberInterests/update`
- 源码：`app/controller/admin/user/MemberInterests.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `sys/system`

### `POST /sys/system/admin/create` — 管理员添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.admin..Admin/create`
- 源码：`app/controller/admin/system/admin/Admin.php` :: `create()`
- 请求参数：
- `account` (query/body, 可选)
- `pwd` (query/body, 可选)
- `phone` (query/body, 可选)
- `againPassword` (query/body, 可选)
- `real_name` (query/body, 可选)
- `roles` (query/body, 可选) 默认 [
- `status` (query/body, 可选) 默认 0
- `is_agent` (query/body, 可选) 默认 0
- `region_ids` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/system/admin/create/form` — 管理员添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.admin..Admin/createForm`
- 源码：`app/controller/admin/system/admin/Admin.php` :: `createForm()`
- 请求参数：
- `is_agent` (query/body, 可选) 默认 0
- `region_ids` (query/body, 可选) 默认 ''
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `DELETE /sys/system/admin/delete/:id` — 管理员删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.admin..Admin/delete`
- 源码：`app/controller/admin/system/admin/Admin.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/system/admin/edit` — 修改信息

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.admin..Admin/edit`
- 源码：`app/controller/admin/system/admin/Admin.php` :: `edit()`
- 请求参数：
- `real_name` (query/body, 可选)
- `phone` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/system/admin/edit/form` — 修改信息表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.admin..Admin/editForm`
- 源码：`app/controller/admin/system/admin/Admin.php` :: `editForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/system/admin/edit/password` — 修改密码

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.admin..Admin/editPassword`
- 源码：`app/controller/admin/system/admin/Admin.php` :: `editPassword()`
- 请求参数：
- `pwd` (query/body, 可选)
- `againPassword` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/system/admin/edit/password/form` — 修改密码表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.admin..Admin/editPasswordForm`
- 源码：`app/controller/admin/system/admin/Admin.php` :: `editPasswordForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/system/admin/log` — 操作日志

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.admin..AdminLog/lst`
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
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/system/admin/lst` — 管理员列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.admin..Admin/getList`
- 源码：`app/controller/admin/system/admin/Admin.php` :: `getList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `date` (query/body, 可选)
- `status` (query/body, 可选)
- `region_ids` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/system/admin/password/:id` — 管理员修改密码

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.admin..Admin/password`
- 源码：`app/controller/admin/system/admin/Admin.php` :: `password()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `pwd` (query/body, 可选)
- `againPassword` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/system/admin/password/form/:id` — 管理员修改密码表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.admin..Admin/passwordForm`
- 源码：`app/controller/admin/system/admin/Admin.php` :: `passwordForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/system/admin/status/:id` — 管理员修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.admin..Admin/switchStatus`
- 源码：`app/controller/admin/system/admin/Admin.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/system/admin/update/:id` — 管理员编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.admin..Admin/update`
- 源码：`app/controller/admin/system/admin/Admin.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `account` (query/body, 可选)
- `real_name` (query/body, 可选)
- `phone` (query/body, 可选)
- `roles` (query/body, 可选) 默认 [
- `status` (query/body, 可选) 默认 0
- `is_agent` (query/body, 可选) 默认 0
- `region_ids` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/system/admin/update/form/:id` — 管理员编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.admin..Admin/updateForm`
- 源码：`app/controller/admin/system/admin/Admin.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `is_agent` (query/body, 可选) 默认 0
- `region_ids` (query/body, 可选) 默认 ''
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/system/applyments/detail/:id` — 分账商户申请详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.MerchantApplyments/detail`
- 源码：`app/controller/admin/system/merchant/MerchantApplyments.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/system/applyments/lst` — 分账商户申请列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.MerchantApplyments/lst`
- 源码：`app/controller/admin/system/merchant/MerchantApplyments.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `mer_name` (query/body, 可选)
- `status` (query/body, 可选)
- `date` (query/body, 可选)
- `mer_applyments_id` (query/body, 可选)
- `out_request_no` (query/body, 可选)
- `applyment_id` (query/body, 可选)
- `mer_id` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /sys/system/applyments/mark/:id` — 分账商户申请备注

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.MerchantApplyments/mark`
- 源码：`app/controller/admin/system/merchant/MerchantApplyments.php` :: `mark()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `mark` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/system/applyments/mark/:id/form` — 分账商户申请备注表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.MerchantApplyments/markForm`
- 源码：`app/controller/admin/system/merchant/MerchantApplyments.php` :: `markForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/system/applyments/merchant/:id` — 分账商户审核查询

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.MerchantApplyments/getMerchant`
- 源码：`app/controller/admin/system/merchant/MerchantApplyments.php` :: `getMerchant()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/system/applyments/status/:id` — 分账商户申请审核

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.MerchantApplyments/switchWithStatus`
- 源码：`app/controller/admin/system/merchant/MerchantApplyments.php` :: `switchWithStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- `message` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/system/article/article/create` — 文章添加

- 置信度：✅ high
- 说明：控制器方法已校验存在；参数由 checkParams/Validate 自动补全
- 处理器：`admin.article.Article/create`
- 源码：`app/controller/admin/article/Article.php` :: `create()`
- 请求参数：
- `cid` (body, 可选) 来自 checkParams/Validate（自动补全）
- `title` (body, 可选) 来自 checkParams/Validate（自动补全）
- `content` (body, 可选) 来自 checkParams/Validate（自动补全）
- `author` (body, 可选) 来自 checkParams/Validate（自动补全）
- `image_input` (body, 可选) 来自 checkParams/Validate（自动补全）
- `status` (body, 可选) 来自 checkParams/Validate（自动补全）
- `sort` (body, 可选) 来自 checkParams/Validate（自动补全）
- `synopsis` (body, 可选) 来自 checkParams/Validate（自动补全）
- `is_hot` (body, 可选) 来自 checkParams/Validate（自动补全）
- `is_banner` (body, 可选) 来自 checkParams/Validate（自动补全）
- `url` (body, 可选) 来自 checkParams/Validate（自动补全）
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/system/article/article/delete/:id` — 文章删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.article.Article/delete`
- 源码：`app/controller/admin/article/Article.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/system/article/article/detail/:id` — 文章详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.article.Article/detail`
- 源码：`app/controller/admin/article/Article.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/system/article/article/lst` — 文章列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.article.Article/getList`
- 源码：`app/controller/admin/article/Article.php` :: `getList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `cid` (query/body, 可选)
- `title` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /sys/system/article/article/status/:id` — 文章修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.article.Article/switchStatus`
- 源码：`app/controller/admin/article/Article.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/system/article/article/update/:id` — 文章编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.article.Article/update`
- 源码：`app/controller/admin/article/Article.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/system/article/category/create` — 文章分类添加

- 置信度：✅ high
- 说明：控制器方法已校验存在；参数由 checkParams/Validate 自动补全
- 处理器：`admin.article.ArticleCategory/create`
- 源码：`app/controller/admin/article/ArticleCategory.php` :: `create()`
- 请求参数：
- `pid` (body, 可选) 来自 checkParams/Validate（自动补全）
- `title` (body, 可选) 来自 checkParams/Validate（自动补全）
- `info` (body, 可选) 来自 checkParams/Validate（自动补全）
- `status` (body, 可选) 来自 checkParams/Validate（自动补全）
- `image` (body, 可选) 来自 checkParams/Validate（自动补全）
- `sort` (body, 可选) 来自 checkParams/Validate（自动补全）
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/system/article/category/create/form` — 文章分类添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.article.ArticleCategory/createForm`
- 源码：`app/controller/admin/article/ArticleCategory.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/system/article/category/delete/:id` — 文章分类删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.article.ArticleCategory/delete`
- 源码：`app/controller/admin/article/ArticleCategory.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/system/article/category/detail/:id` — 文章分类详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.article.ArticleCategory/detail`
- 源码：`app/controller/admin/article/ArticleCategory.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/system/article/category/lst` — 文章分类列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.article.ArticleCategory/lst`
- 源码：`app/controller/admin/article/ArticleCategory.php` :: `lst()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/system/article/category/select` — 文章分类筛选

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.article.ArticleCategory/select`
- 源码：`app/controller/admin/article/ArticleCategory.php` :: `select()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/system/article/category/status/:id` — 文章分类修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.article.ArticleCategory/switchStatus`
- 源码：`app/controller/admin/article/ArticleCategory.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/system/article/category/update/:id` — 文章分类编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.article.ArticleCategory/update`
- 源码：`app/controller/admin/article/ArticleCategory.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/system/article/category/update/form/:id` — 文章分类编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.article.ArticleCategory/updateForm`
- 源码：`app/controller/admin/article/ArticleCategory.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/system/attachment/category` — 批量移动

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.attachment.Attachment/batchChangeCategory`
- 源码：`app/controller/admin/system/attachment/Attachment.php` :: `batchChangeCategory()`
- 请求参数：
- `ids` (query/body, 可选) 默认 [
- `attachment_category_id` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/system/attachment/category/create` — 素材分类添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.attachment.AttachmentCategory/create`
- 源码：`app/controller/admin/system/attachment/AttachmentCategory.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/system/attachment/category/create/form` — 素材分类添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.attachment.AttachmentCategory/createForm`
- 源码：`app/controller/admin/system/attachment/AttachmentCategory.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/system/attachment/category/delete/:id` — 素材删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.attachment.AttachmentCategory/delete`
- 源码：`app/controller/admin/system/attachment/AttachmentCategory.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/system/attachment/category/formatLst` — 素材分类列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.attachment.AttachmentCategory/getFormatList`
- 源码：`app/controller/admin/system/attachment/AttachmentCategory.php` :: `getFormatList()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/system/attachment/category/update/:id` — 素材编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.attachment.AttachmentCategory/update`
- 源码：`app/controller/admin/system/attachment/AttachmentCategory.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/system/attachment/category/update/form/:id` — 素材分类编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.attachment.AttachmentCategory/updateForm`
- 源码：`app/controller/admin/system/attachment/AttachmentCategory.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/system/attachment/delete` — 素材删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.attachment.Attachment/delete`
- 源码：`app/controller/admin/system/attachment/Attachment.php` :: `delete()`
- 请求参数：
- `ids` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/system/attachment/lst` — 素材列表

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

### `POST /sys/system/attachment/online_upload` — 在线图片

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.attachment.Attachment/onlineUpload`
- 源码：`app/controller/admin/system/attachment/Attachment.php` :: `onlineUpload()`
- 请求参数：
- `id` (query/body, 可选)
- `images` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/system/attachment/scan_upload/image/:token` — 扫码上传图片

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.attachment.Attachment/scanUploadImage`
- 源码：`app/controller/admin/system/attachment/Attachment.php` :: `scanUploadImage()`
- 请求参数：
- `token` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/system/attachment/scan_upload/image/:token` — 扫码上传保存

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.attachment.Attachment/scanUploadSave`
- 源码：`app/controller/admin/system/attachment/Attachment.php` :: `scanUploadSave()`
- 请求参数：
- `token` (path, 必填) 路径参数
- `pid` (query/body, 可选)
- `ids` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/system/attachment/scan_upload/qrcode/:pid` — 上传二维码

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.attachment.Attachment/scanUploadQrcode`
- 源码：`app/controller/admin/system/attachment/Attachment.php` :: `scanUploadQrcode()`
- 请求参数：
- `pid` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/system/attachment/update/:id` — 素材编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.attachment.Attachment/update`
- 源码：`app/controller/admin/system/attachment/Attachment.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `attachment_name` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/system/attachment/update/:id/form` — 素材编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.attachment.Attachment/updateForm`
- 源码：`app/controller/admin/system/attachment/Attachment.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/system/city/lst` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.shipping.City/lst`
- 源码：`app/controller/merchant/store/shipping/City.php` :: `lst()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/system/form/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在；参数由 checkParams/Validate 自动补全
- 处理器：`admin.system.form.Form/create`
- 源码：`app/controller/admin/system/form/Form.php` :: `create()`
- 请求参数：
- `name` (body, 可选) 来自 checkParams/Validate（自动补全）
- `value` (body, 可选) 来自 checkParams/Validate（自动补全）
- `status` (body, 可选) 来自 checkParams/Validate（自动补全）
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/system/form/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.form.Form/delete`
- 源码：`app/controller/admin/system/form/Form.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/system/form/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.form.Form/detail`
- 源码：`app/controller/admin/system/form/Form.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/system/form/excel` — 导出

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.form.Form/excel`
- 源码：`app/controller/admin/system/form/Form.php` :: `excel()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `id` (query/body, 可选)
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/system/form/info/:id` — info

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.form.Form/info`
- 源码：`app/controller/admin/system/form/Form.php` :: `info()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `mer_id` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/system/form/lst` — 列表

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

### `GET /sys/system/form/select` — select

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.form.Form/select`
- 源码：`app/controller/admin/system/form/Form.php` :: `select()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/system/form/status/:id` — 编辑状态

- 置信度：⛔ stale
- 说明：路由已登记，但 `app/controller/admin/system/form/Form.php` 中不存在方法 `statusSwitch`（原项目死路由/加密扩展/版本差异）。开发勿实现为有效接口，除非核实。
- 处理器：`admin.system.form.Form/statusSwitch`
- 源码：`app/controller/admin/system/form/Form.php` :: `statusSwitch()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一响应包 {status,message,data?}；详见控制器实现 | 外层: {status,message,data}

### `POST /sys/system/form/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.form.Form/update`
- 源码：`app/controller/admin/system/form/Form.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/system/form/user_lst/:id` — 表单提交记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.form.Form/formUserList`
- 源码：`app/controller/admin/system/form/Form.php` :: `formUserList()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/system/menu/create` — 平台菜单/权限添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.auth.Menu/create`
- 源码：`app/controller/admin/system/auth/Menu.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/system/menu/create/form` — 平台菜单/权限添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.auth.Menu/createForm`
- 源码：`app/controller/admin/system/auth/Menu.php` :: `createForm()`
- 请求参数：
- `is_agent` (query/body, 可选) 默认 0
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/system/menu/delete/:id` — 平台菜单/权限删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.auth.Menu/delete`
- 源码：`app/controller/admin/system/auth/Menu.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/system/menu/getAdminMenusList` — getMenusList

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.auth.Menu/getMenusList`
- 源码：`app/controller/admin/system/auth/Menu.php` :: `getMenusList()`
- 请求参数：
- `is_mer` (query/body, 可选)
- `keyword` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/system/menu/lst` — 平台菜单/权限列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.auth.Menu/getList`
- 源码：`app/controller/admin/system/auth/Menu.php` :: `getList()`
- 请求参数：
- `is_agent` (query/body, 可选) 默认 0
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/system/menu/update/:id` — 平台菜单/权限编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.auth.Menu/update`
- 源码：`app/controller/admin/system/auth/Menu.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/system/menu/update/form/:id` — 平台菜单/权限编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.auth.Menu/updateForm`
- 源码：`app/controller/admin/system/auth/Menu.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `is_agent` (query/body, 可选) 默认 0
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/system/merchant/businessCreate` — 商户添加店铺

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant..Merchant/businessCreate`
- 源码：`app/controller/admin/system/merchant/Merchant.php` :: `businessCreate()`
- 请求参数：
- `mer_name` (query/body, 可选)
- `mer_password` (query/body, 可选)
- `mer_phone` (query/body, 可选)
- `mer_account` (query/body, 可选)
- `category_id` (query/body, 可选)
- `type_id` (query/body, 可选)
- `real_name` (query/body, 可选)
- `mark` (query/body, 可选) 默认 ''
- `business_id` (query/body, 可选) 默认 0
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/system/merchant/care_ficti/:id` — 虚拟关注量

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant..Merchant/careFicti`
- 源码：`app/controller/admin/system/merchant/Merchant.php` :: `careFicti()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `type` (query/body, 可选)
- `num` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/system/merchant/care_ficti/form/:id` — 虚拟关注量表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant..Merchant/careFictiForm`
- 源码：`app/controller/admin/system/merchant/Merchant.php` :: `careFictiForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/system/merchant/category` — 店铺分类添加

- 置信度：✅ high
- 说明：控制器方法已校验存在；参数由 checkParams/Validate 自动补全
- 处理器：`admin.system.merchant.MerchantCategory/create`
- 源码：`app/controller/admin/system/merchant/MerchantCategory.php` :: `create()`
- 请求参数：
- `category_name` (body, 可选) 来自 checkParams/Validate（自动补全）
- `commission_rate` (body, 可选) 来自 checkParams/Validate（自动补全）
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/system/merchant/category/:id` — 店铺分类删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.MerchantCategory/delete`
- 源码：`app/controller/admin/system/merchant/MerchantCategory.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/system/merchant/category/:id` — 店铺分类编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.MerchantCategory/update`
- 源码：`app/controller/admin/system/merchant/MerchantCategory.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/system/merchant/category/form` — 店铺分类添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.MerchantCategory/createForm`
- 源码：`app/controller/admin/system/merchant/MerchantCategory.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/system/merchant/category/form/:id` — 店铺分类编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.MerchantCategory/updateForm`
- 源码：`app/controller/admin/system/merchant/MerchantCategory.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/system/merchant/category/lst` — 店铺分类列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.MerchantCategory/lst`
- 源码：`app/controller/admin/system/merchant/MerchantCategory.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /sys/system/merchant/category/options` — 店铺分类筛选

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.MerchantCategory/getOptions`
- 源码：`app/controller/admin/system/merchant/MerchantCategory.php` :: `getOptions()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/system/merchant/category_lst` — 店铺分类列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.MerchantCategory/lst`
- 源码：`app/controller/admin/system/merchant/MerchantCategory.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /sys/system/merchant/changecopy/:id` — 修改采集商品次数

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant..Merchant/changeCopyNum`
- 源码：`app/controller/admin/system/merchant/Merchant.php` :: `changeCopyNum()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `type` (query/body, 可选)
- `num` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/system/merchant/changecopy/:id/form` — 修改采集商品次数表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant..Merchant/changeCopyNumForm`
- 源码：`app/controller/admin/system/merchant/Merchant.php` :: `changeCopyNumForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/system/merchant/close/:id` — 店铺开启/关闭

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant..Merchant/switchClose`
- 源码：`app/controller/admin/system/merchant/Merchant.php` :: `switchClose()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/system/merchant/count` — 店铺列表统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant..Merchant/count`
- 源码：`app/controller/admin/system/merchant/Merchant.php` :: `count()`
- 请求参数：
- `keyword` (query/body, 可选)
- `date` (query/body, 可选)
- `status` (query/body, 可选)
- `statusTag` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `category_id` (query/body, 可选)
- `type_id` (query/body, 可选)
- `is_best` (query/body, 可选)
- `offline_switch` (query/body, 可选)
- `region_id` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/system/merchant/create` — 店铺添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant..Merchant/create`
- 源码：`app/controller/admin/system/merchant/Merchant.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/system/merchant/create/form` — 店铺列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant..Merchant/createForm`
- 源码：`app/controller/admin/system/merchant/Merchant.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/system/merchant/delete/:id` — 店铺删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant..Merchant/delete`
- 源码：`app/controller/admin/system/merchant/Merchant.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/system/merchant/delete/:id/form` — 店铺删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant..Merchant/deleteForm`
- 源码：`app/controller/admin/system/merchant/Merchant.php` :: `deleteForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/system/merchant/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant..Merchant/detail`
- 源码：`app/controller/admin/system/merchant/Merchant.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/system/merchant/get_operate_list/:merchant_id` — 操作日志

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant..Merchant/getOperateList`
- 源码：`app/controller/admin/system/merchant/Merchant.php` :: `getOperateList()`
- 请求参数：
- `merchant_id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `type` (query/body, 可选) 默认 ''
- `date` (query/body, 可选) 默认 ''
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/system/merchant/group/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.StoreGroup/create`
- 源码：`app/controller/admin/system/merchant/StoreGroup.php` :: `create()`
- 请求参数：
- `pid` (query/body, 可选)
- `name` (query/body, 可选)
- `sort` (query/body, 可选) 默认 0
- `status` (query/body, 可选) 默认 1
- `positioning_status` (query/body, 可选) 默认 0
- `longitude` (query/body, 可选)
- `latitude` (query/body, 可选)
- `merchant_ids` (query/body, 可选)
- `address` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/system/merchant/group/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.StoreGroup/delete`
- 源码：`app/controller/admin/system/merchant/StoreGroup.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/system/merchant/group/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.StoreGroup/detail`
- 源码：`app/controller/admin/system/merchant/StoreGroup.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/system/merchant/group/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.StoreGroup/list`
- 源码：`app/controller/admin/system/merchant/StoreGroup.php` :: `list()`
- 请求参数：
- `name` (query/body, 可选)
- `status` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/system/merchant/group/options` — 筛选

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.StoreGroup/options`
- 源码：`app/controller/admin/system/merchant/StoreGroup.php` :: `options()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/system/merchant/group/setTemplate/:id` — 设置店铺分组模板

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.StoreGroup/setTemplate`
- 源码：`app/controller/admin/system/merchant/StoreGroup.php` :: `setTemplate()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `diy_temp_id` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/system/merchant/group/status/:id` — 状态切换

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.StoreGroup/switchStatus`
- 源码：`app/controller/admin/system/merchant/StoreGroup.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/system/merchant/group/stores/:id` — 关联店铺列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.StoreGroup/stores`
- 源码：`app/controller/admin/system/merchant/StoreGroup.php` :: `stores()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/system/merchant/group/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant.StoreGroup/update`
- 源码：`app/controller/admin/system/merchant/StoreGroup.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `pid` (query/body, 可选)
- `name` (query/body, 可选)
- `sort` (query/body, 可选) 默认 0
- `status` (query/body, 可选) 默认 1
- `positioning_status` (query/body, 可选) 默认 0
- `longitude` (query/body, 可选)
- `latitude` (query/body, 可选)
- `merchant_ids` (query/body, 可选)
- `address` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/system/merchant/login/:id` — 店铺登录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant..Merchant/login`
- 源码：`app/controller/admin/system/merchant/Merchant.php` :: `login()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/system/merchant/lst` — 店铺列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant..Merchant/lst`
- 源码：`app/controller/admin/system/merchant/Merchant.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `date` (query/body, 可选)
- `status` (query/body, 可选)
- `statusTag` (query/body, 可选)
- `is_trader` (query/body, 可选)
- `category_id` (query/body, 可选)
- `type_id` (query/body, 可选)
- `order` (query/body, 可选) 默认 'create_time'
- `is_best` (query/body, 可选)
- `offline_switch` (query/body, 可选)
- `region_id` (query/body, 可选)
- `business_id` (query/body, 可选)
- `mer_state` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/system/merchant/mer_select` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant..Merchant/mer_select`
- 源码：`app/controller/admin/system/merchant/Merchant.php` :: `mer_select()`
- 请求参数：
- `keyword` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/system/merchant/password/:id` — 店铺修改密码

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant..MerchantAdmin/password`
- 源码：`app/controller/admin/system/merchant/MerchantAdmin.php` :: `password()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `pwd` (query/body, 可选)
- `againPassword` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/system/merchant/password/form/:id` — 店铺修改密码表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant..MerchantAdmin/passwordForm`
- 源码：`app/controller/admin/system/merchant/MerchantAdmin.php` :: `passwordForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/system/merchant/status/:id` — 店铺修改推荐

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant..Merchant/switchStatus`
- 源码：`app/controller/admin/system/merchant/Merchant.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/system/merchant/update/:id` — 店铺编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant..Merchant/update`
- 源码：`app/controller/admin/system/merchant/Merchant.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/system/merchant/update/form/:id` — 店铺编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.merchant..Merchant/updateForm`
- 源码：`app/controller/admin/system/merchant/Merchant.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/system/role/create` — 身份添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.auth.Role/create`
- 源码：`app/controller/admin/system/auth/Role.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/system/role/create/form` — 身份添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.auth.Role/createForm`
- 源码：`app/controller/admin/system/auth/Role.php` :: `createForm()`
- 请求参数：
- `is_agent` (query/body, 可选) 默认 0
- `circle_id` (query/body, 可选) 默认 0
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/system/role/delete/:id` — 身份删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.auth.Role/delete`
- 源码：`app/controller/admin/system/auth/Role.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/system/role/lst` — 身份列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.auth.Role/getList`
- 源码：`app/controller/admin/system/auth/Role.php` :: `getList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `role_name` (query/body, 可选)
- `is_agent` (query/body, 可选) 默认 0
- `circle_id` (query/body, 可选) 默认 0
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /sys/system/role/status/:id` — 身份修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.auth.Role/switchStatus`
- 源码：`app/controller/admin/system/auth/Role.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/system/role/update/:id` — 身份编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.auth.Role/update`
- 源码：`app/controller/admin/system/auth/Role.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/system/role/update/form/:id` — 身份编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.auth.Role/updateForm`
- 源码：`app/controller/admin/system/auth/Role.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `is_agent` (query/body, 可选) 默认 0
- `circle_id` (query/body, 可选) 默认 0
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/system/storage/:type/form` — 获取云存储配置表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.SystemStorage/form`
- 源码：`app/controller/admin/system/SystemStorage.php` :: `form()`
- 请求参数：
- `type` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/system/storage/config` — 配置信息

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.SystemStorage/getConfig`
- 源码：`app/controller/admin/system/SystemStorage.php` :: `getConfig()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/system/storage/config` — 提交配置

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.SystemStorage/setConfig`
- 源码：`app/controller/admin/system/SystemStorage.php` :: `setConfig()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/system/storage/domain/update/:id` — 修改存储空间名称

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.SystemStorage/editDomain`
- 源码：`app/controller/admin/system/SystemStorage.php` :: `editDomain()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/system/storage/domain/update/:id/form` — 修改存储空间名称表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.SystemStorage/editDomainForm`
- 源码：`app/controller/admin/system/SystemStorage.php` :: `editDomainForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/system/storage/region/create/:type` — 添加存储空间

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.SystemStorage/createRegion`
- 源码：`app/controller/admin/system/SystemStorage.php` :: `createRegion()`
- 请求参数：
- `type` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/system/storage/region/create/:type/form` — 添加存储空间表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.SystemStorage/createRegionForm`
- 源码：`app/controller/admin/system/SystemStorage.php` :: `createRegionForm()`
- 请求参数：
- `type` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/system/storage/region/delete/:id` — 删除存储空间

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.SystemStorage/deleteRegion`
- 源码：`app/controller/admin/system/SystemStorage.php` :: `deleteRegion()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/system/storage/region/lst/:type` — 存储空间列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.SystemStorage/lstRegion`
- 源码：`app/controller/admin/system/SystemStorage.php` :: `lstRegion()`
- 请求参数：
- `type` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/system/storage/region/status/:id` — 使用存储空间

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.SystemStorage/swtichStatus`
- 源码：`app/controller/admin/system/SystemStorage.php` :: `swtichStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/system/storage/set_key` — 保存云存储配置

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.SystemStorage/setForm`
- 源码：`app/controller/admin/system/SystemStorage.php` :: `setForm()`
- 请求参数：
- `upload_type` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/system/storage/sync/:type` — 同步存储空间

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.SystemStorage/sync`
- 源码：`app/controller/admin/system/SystemStorage.php` :: `sync()`
- 请求参数：
- `type` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/system/storage/type_list` — typeList

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.SystemStorage/typeList`
- 源码：`app/controller/admin/system/SystemStorage.php` :: `typeList()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `sys/upload`

### `GET /sys/upload/config` — 上传配置

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.Common/uploadConfig`
- 源码：`app/controller/admin/Common.php` :: `uploadConfig()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/upload/config` — 上传配置保存

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.Common/saveUploadConfig`
- 源码：`app/controller/admin/Common.php` :: `saveUploadConfig()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/upload/image/:id/:field` — 上传图片

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.attachment.Attachment/image`
- 源码：`app/controller/admin/system/attachment/Attachment.php` :: `image()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `field` (path, 必填) 路径参数
- `ueditor` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/upload/temp_key` — 上传视屏KEY

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.Common/temp_key`
- 源码：`app/controller/admin/Common.php` :: `temp_key()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/upload/video` — uploadVideo

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.attachment.Attachment/uploadVideo`
- 源码：`app/controller/admin/system/attachment/Attachment.php` :: `uploadVideo()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。


## `sys/user`

### `POST /sys/user/batch_change_group` — 用户分组批量编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/batchChangeGroup`
- 源码：`app/controller/admin/user/User.php` :: `batchChangeGroup()`
- 请求参数：
- `group_id` (query/body, 可选)
- `ids` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/user/batch_change_group/form` — 用户分组批量编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/batchChangeGroupForm`
- 源码：`app/controller/admin/user/User.php` :: `batchChangeGroupForm()`
- 请求参数：
- `ids` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/user/batch_change_label` — 用户标签批量编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/batchChangeLabel`
- 源码：`app/controller/admin/user/User.php` :: `batchChangeLabel()`
- 请求参数：
- `label_id` (query/body, 可选)
- `ids` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/user/batch_change_label/form` — 用户标签批量编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/batchChangeLabelForm`
- 源码：`app/controller/admin/user/User.php` :: `batchChangeLabelForm()`
- 请求参数：
- `ids` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/user/batch_spread` — 批量设置分销员

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/batchSpread`
- 源码：`app/controller/admin/user/User.php` :: `batchSpread()`
- 请求参数：
- `uids` (query/body, 可选)
- `is_promoter` (query/body, 可选)
- `promoter_switch` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/user/batch_spread_form` — getMemberLevelBatchSpreadForm

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/batchSpreadForm`
- 源码：`app/controller/admin/user/User.php` :: `batchSpreadForm()`
- 请求参数：
- `ids` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/user/bill/:id` — 用户余额变动列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/bill`
- 源码：`app/controller/admin/user/User.php` :: `bill()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/user/brokerage/create` — 分销员等级添加

- 置信度：✅ high
- 说明：控制器方法已校验存在；参数由 checkParams/Validate 自动补全
- 处理器：`admin.user..UserBrokerage/create`
- 源码：`app/controller/admin/user/UserBrokerage.php` :: `create()`
- 请求参数：
- `brokerage_level` (body, 可选) 来自 checkParams/Validate（自动补全）
- `brokerage_name` (body, 可选) 来自 checkParams/Validate（自动补全）
- `brokerage_icon` (body, 可选) 来自 checkParams/Validate（自动补全）
- `brokerage_rule` (body, 可选) 来自 checkParams/Validate（自动补全）
- `extension_one` (body, 可选) 来自 checkParams/Validate（自动补全）
- `extension_two` (body, 可选) 来自 checkParams/Validate（自动补全）
- `image` (body, 可选) 来自 checkParams/Validate（自动补全）
- `value` (body, 可选) 来自 checkParams/Validate（自动补全）
- `type` (body, 可选) 来自 checkParams/Validate（自动补全）
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `DELETE /sys/user/brokerage/delete/:id` — 分销员等级删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user..UserBrokerage/delete`
- 源码：`app/controller/admin/user/UserBrokerage.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/user/brokerage/detail/:id` — 分销员等级详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user..UserBrokerage/detail`
- 源码：`app/controller/admin/user/UserBrokerage.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/user/brokerage/lst` — 分销员等级列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user..UserBrokerage/getLst`
- 源码：`app/controller/admin/user/UserBrokerage.php` :: `getLst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `brokerage_name` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/user/brokerage/options` — 分销员等级筛选

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user..UserBrokerage/options`
- 源码：`app/controller/admin/user/UserBrokerage.php` :: `options()`
- 请求参数：
- `type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/user/brokerage/update/:id` — 分销员等级编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user..UserBrokerage/update`
- 源码：`app/controller/admin/user/UserBrokerage.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/user/change_group/:id` — 用户分组编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/changeGroup`
- 源码：`app/controller/admin/user/User.php` :: `changeGroup()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `group_id` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/user/change_group/form/:id` — 用户分组编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/changeGroupForm`
- 源码：`app/controller/admin/user/User.php` :: `changeGroupForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/user/change_integral/:id` — 用户修改积分

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/changeIntegral`
- 源码：`app/controller/admin/user/User.php` :: `changeIntegral()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `now_money` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/user/change_integral/form/:id` — 用户修改积分表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/changeIntegralForm`
- 源码：`app/controller/admin/user/User.php` :: `changeIntegralForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/user/change_label/:id` — 用户标签编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/changeLabel`
- 源码：`app/controller/admin/user/User.php` :: `changeLabel()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `label_id` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/user/change_label/form/:id` — 用户标签编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/changeLabelForm`
- 源码：`app/controller/admin/user/User.php` :: `changeLabelForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/user/change_now_money/:id` — 用户修改余额

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/changeNowMoney`
- 源码：`app/controller/admin/user/User.php` :: `changeNowMoney()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `now_money` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/user/change_now_money/form/:id` — 用户修改余额表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/changeNowMoneyForm`
- 源码：`app/controller/admin/user/User.php` :: `changeNowMoneyForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/user/change_password/:id` — 用户修改密码

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/changePassword`
- 源码：`app/controller/admin/user/User.php` :: `changePassword()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `pwd` (query/body, 可选)
- `repwd` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/user/change_password/form/:id` — 用户修改密码表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/changePasswordForm`
- 源码：`app/controller/admin/user/User.php` :: `changePasswordForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/user/change_spread/:id` — 修改推荐人

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/spread`
- 源码：`app/controller/admin/user/User.php` :: `spread()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `spid` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/user/change_spread_form/:id` — 修改推荐人表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/spreadForm`
- 源码：`app/controller/admin/user/User.php` :: `spreadForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/user/clear_search_log` — 清除用户搜索记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/clearSearchLog`
- 源码：`app/controller/admin/user/User.php` :: `clearSearchLog()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/user/coupon/:id` — 用户持有优惠券

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/coupon`
- 源码：`app/controller/admin/user/User.php` :: `coupon()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/user/create` — 用户添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/createForm`
- 源码：`app/controller/admin/user/User.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/user/create` — 用户添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/create`
- 源码：`app/controller/admin/user/User.php` :: `create()`
- 请求参数：
- `account` (query/body, 可选)
- `pwd` (query/body, 可选)
- `repwd` (query/body, 可选)
- `nickname` (query/body, 可选)
- `avatar` (query/body, 可选)
- `real_name` (query/body, 可选)
- `phone` (query/body, 可选)
- `sex` (query/body, 可选)
- `status` (query/body, 可选)
- `card_id` (query/body, 可选)
- `is_promoter` (query/body, 可选) 默认 0
- `extend_info` (query/body, 可选) 默认 [
- `promoter_switch` (query/body, 可选) 默认 1
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/user/detail/:id` — 用户详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/detail`
- 源码：`app/controller/admin/user/User.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/user/excel` — 用户信息导出

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/excel`
- 源码：`app/controller/admin/user/User.php` :: `excel()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `label_id` (query/body, 可选)
- `user_type` (query/body, 可选)
- `sex` (query/body, 可选)
- `is_promoter` (query/body, 可选)
- `country` (query/body, 可选)
- `pay_count` (query/body, 可选)
- `user_time_type` (query/body, 可选)
- `user_time` (query/body, 可选)
- `nickname` (query/body, 可选)
- `province` (query/body, 可选)
- `city` (query/body, 可选)
- `group_id` (query/body, 可选)
- `phone` (query/body, 可选)
- `uid` (query/body, 可选)
- `is_svip` (query/body, 可选)
- `fields_type` (query/body, 可选)
- `fields_value` (query/body, 可选)
- `ids` (query/body, 可选)
- `filter_conditions` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/user/extract/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserExtract/detail`
- 源码：`app/controller/admin/user/UserExtract.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/user/extract/export` — 导出

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserExtract/export`
- 源码：`app/controller/admin/user/UserExtract.php` :: `export()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `status` (query/body, 可选)
- `keyword` (query/body, 可选)
- `date` (query/body, 可选)
- `extract_type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/user/extract/lst` — 申请列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserExtract/lst`
- 源码：`app/controller/admin/user/UserExtract.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `status` (query/body, 可选)
- `keyword` (query/body, 可选)
- `date` (query/body, 可选)
- `extract_type` (query/body, 可选)
- `uid` (query/body, 可选)
- `phone` (query/body, 可选)
- `real_name` (query/body, 可选)
- `nickname` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/user/extract/status/:id` — 审核

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserExtract/switchStatus`
- 源码：`app/controller/admin/user/UserExtract.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `fail_msg` (query/body, 可选)
- `mark` (query/body, 可选)
- `status` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/user/extract/status_form/:id` — 审核表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserExtract/switchStatusForm`
- 源码：`app/controller/admin/user/UserExtract.php` :: `switchStatusForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/user/feedback/category/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.FeedBackCategory/create`
- 源码：`app/controller/admin/user/FeedBackCategory.php` :: `create()`
- 请求参数：
- `pid` (query/body, 可选)
- `cate_name` (query/body, 可选)
- `sort` (query/body, 可选)
- `pic` (query/body, 可选)
- `is_show` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/user/feedback/category/create/form` — 添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.FeedBackCategory/createForm`
- 源码：`app/controller/admin/user/FeedBackCategory.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/user/feedback/category/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.FeedBackCategory/delete`
- 源码：`app/controller/admin/user/FeedBackCategory.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/user/feedback/category/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.FeedBackCategory/lst`
- 源码：`app/controller/admin/user/FeedBackCategory.php` :: `lst()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/user/feedback/category/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.FeedBackCategory/switchStatus`
- 源码：`app/controller/admin/user/FeedBackCategory.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/user/feedback/category/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.FeedBackCategory/update`
- 源码：`app/controller/admin/user/FeedBackCategory.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `pid` (query/body, 可选)
- `cate_name` (query/body, 可选)
- `sort` (query/body, 可选)
- `pic` (query/body, 可选)
- `is_show` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/user/feedback/category/update/:id/form` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.FeedBackCategory/updateForm`
- 源码：`app/controller/admin/user/FeedBackCategory.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/user/feedback/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.FeedBack/delete`
- 源码：`app/controller/admin/user/FeedBack.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/user/feedback/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.FeedBack/detail`
- 源码：`app/controller/admin/user/FeedBack.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/user/feedback/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.FeedBack/lst`
- 源码：`app/controller/admin/user/FeedBack.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `type` (query/body, 可选)
- `status` (query/body, 可选)
- `realname` (query/body, 可选)
- `is_del` (query/body, 可选) 默认 0
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/user/feedback/reply/:id` — 回复

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.FeedBack/reply`
- 源码：`app/controller/admin/user/FeedBack.php` :: `reply()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `reply` (query/body, 可选)
- `remake` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/user/feedback/reply/:id/form` — 回复表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.FeedBack/replyForm`
- 源码：`app/controller/admin/user/FeedBack.php` :: `replyForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/user/fields/save/:uid` — 添加或编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserFields/save`
- 源码：`app/controller/admin/user/UserFields.php` :: `save()`
- 请求参数：
- `uid` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/user/fields/save_form/:uid` — 扩展信息表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserFields/saveForm`
- 源码：`app/controller/admin/user/UserFields.php` :: `saveForm()`
- 请求参数：
- `uid` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/user/filters` — filters

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserSearch/filters`
- 源码：`app/controller/admin/user/UserSearch.php` :: `filters()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/user/get_fields` — 用户扩展信息表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/getFields`
- 源码：`app/controller/admin/user/User.php` :: `getFields()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `DELETE /sys/user/group/:id` — 用户分组删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserGroup/delete`
- 源码：`app/controller/admin/user/UserGroup.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/user/group/:id` — 用户分组编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserGroup/update`
- 源码：`app/controller/admin/user/UserGroup.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/user/group/form` — 用户分组添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserGroup/createForm`
- 源码：`app/controller/admin/user/UserGroup.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/user/group/form/:id` — 用户分组编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserGroup/updateForm`
- 源码：`app/controller/admin/user/UserGroup.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/user/group/lst` — 用户分组列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserGroup/lst`
- 源码：`app/controller/admin/user/UserGroup.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /sys/user/group/user/group` — 用户分组添加

- 置信度：✅ high
- 说明：控制器方法已校验存在；参数由 checkParams/Validate 自动补全
- 处理器：`admin.user.UserGroup/create`
- 源码：`app/controller/admin/user/UserGroup.php` :: `create()`
- 请求参数：
- `group_name` (body, 可选) 来自 checkParams/Validate（自动补全）
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/user/history/:id` — 浏览记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/history`
- 源码：`app/controller/admin/user/User.php` :: `history()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/user/info/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserInfo/create`
- 源码：`app/controller/admin/user/UserInfo.php` :: `create()`
- 请求参数：
- `field` (query/body, 可选)
- `title` (query/body, 可选)
- `is_used` (query/body, 可选)
- `is_require` (query/body, 可选)
- `is_show` (query/body, 可选)
- `type` (query/body, 可选)
- `msg` (query/body, 可选)
- `content` (query/body, 可选) 默认 [
- `sort` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/user/info/create_from` — 添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserInfo/createFrom`
- 源码：`app/controller/admin/user/UserInfo.php` :: `createFrom()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/user/info/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserInfo/delete`
- 源码：`app/controller/admin/user/UserInfo.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/user/info/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserInfo/lst`
- 源码：`app/controller/admin/user/UserInfo.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `date` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /sys/user/info/save_all` — 保存信息

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserInfo/saveAll`
- 源码：`app/controller/admin/user/UserInfo.php` :: `saveAll()`
- 请求参数：
- `avatar` (query/body, 可选) 默认 ''
- `user_extend_info` (query/body, 可选) 默认 [
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/user/info/select_list` — 下拉列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserInfo/getSelectList`
- 源码：`app/controller/admin/user/UserInfo.php` :: `getSelectList()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/user/info/type` — 类型

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserInfo/getType`
- 源码：`app/controller/admin/user/UserInfo.php` :: `getType()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/user/integral/:id` — 积分记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/integralList`
- 源码：`app/controller/admin/user/User.php` :: `integralList()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/user/integral/config` — 积分配置获取

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user..UserIntegral/getConfig`
- 源码：`app/controller/admin/user/UserIntegral.php` :: `getConfig()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/user/integral/config` — 积分配置保存

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user..UserIntegral/saveConfig`
- 源码：`app/controller/admin/user/UserIntegral.php` :: `saveConfig()`
- 请求参数：
- `integral_status` (query/body, 可选)
- `integral_clear_time` (query/body, 可选)
- `integral_order_rate` (query/body, 可选)
- `integral_freeze` (query/body, 可选)
- `integral_user_give` (query/body, 可选)
- `integral_money` (query/body, 可选)
- `integral_community_give` (query/body, 可选)
- `integral_community_give_limit` (query/body, 可选)
- `rule` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/user/integral/excel` — 积分导出

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user..UserIntegral/excel`
- 源码：`app/controller/admin/user/UserIntegral.php` :: `excel()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `date` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/user/integral/lst` — 积分日志

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user..UserIntegral/getList`
- 源码：`app/controller/admin/user/UserIntegral.php` :: `getList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `date` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/user/integral/title` — 积分统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user..UserIntegral/getTitle`
- 源码：`app/controller/admin/user/UserIntegral.php` :: `getTitle()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `DELETE /sys/user/label/:id` — 用户标签删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserLabel/delete`
- 源码：`app/controller/admin/user/UserLabel.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/user/label/:id` — 用户标签编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserLabel/update`
- 源码：`app/controller/admin/user/UserLabel.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/user/label/form` — 用户标签添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserLabel/createForm`
- 源码：`app/controller/admin/user/UserLabel.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/user/label/form/:id` — 用户标签编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserLabel/updateForm`
- 源码：`app/controller/admin/user/UserLabel.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/user/label/lst` — 用户标签列表

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

### `POST /sys/user/label/user/label` — 用户标签添加

- 置信度：✅ high
- 说明：控制器方法已校验存在；参数由 checkParams/Validate 自动补全
- 处理器：`admin.user.UserLabel/create`
- 源码：`app/controller/admin/user/UserLabel.php` :: `create()`
- 请求参数：
- `label_name` (body, 可选) 来自 checkParams/Validate（自动补全）
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/user/lst` — 用户列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/lst`
- 源码：`app/controller/admin/user/User.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `nickname` (query/body, 可选)
- `phone` (query/body, 可选)
- `uid` (query/body, 可选)
- `user_type` (query/body, 可选)
- `label_id` (query/body, 可选)
- `sex` (query/body, 可选)
- `is_promoter` (query/body, 可选)
- `country` (query/body, 可选)
- `pay_count` (query/body, 可选)
- `user_time_type` (query/body, 可选)
- `user_time` (query/body, 可选)
- `province` (query/body, 可选)
- `city` (query/body, 可选)
- `group_id` (query/body, 可选)
- `is_svip` (query/body, 可选)
- `fields_type` (query/body, 可选)
- `fields_value` (query/body, 可选)
- `member_level` (query/body, 可选)
- `keyword` (query/body, 可选)
- `birthday` (query/body, 可选)
- `filter_conditions` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /sys/user/member/:id/form` — 用户修改会员等级表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/memberForm`
- 源码：`app/controller/admin/user/User.php` :: `memberForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/user/member/:id/save` — 用户修改会员等级

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/memberSave`
- 源码：`app/controller/admin/user/User.php` :: `memberSave()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `member_level` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/user/member/create` — 普通会员等级添加

- 置信度：✅ high
- 说明：控制器方法已校验存在；参数由 checkParams/Validate 自动补全
- 处理器：`admin.user..UserBrokerage/create`
- 源码：`app/controller/admin/user/UserBrokerage.php` :: `create()`
- 请求参数：
- `brokerage_level` (body, 可选) 来自 checkParams/Validate（自动补全）
- `brokerage_name` (body, 可选) 来自 checkParams/Validate（自动补全）
- `brokerage_icon` (body, 可选) 来自 checkParams/Validate（自动补全）
- `brokerage_rule` (body, 可选) 来自 checkParams/Validate（自动补全）
- `extension_one` (body, 可选) 来自 checkParams/Validate（自动补全）
- `extension_two` (body, 可选) 来自 checkParams/Validate（自动补全）
- `image` (body, 可选) 来自 checkParams/Validate（自动补全）
- `value` (body, 可选) 来自 checkParams/Validate（自动补全）
- `type` (body, 可选) 来自 checkParams/Validate（自动补全）
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/user/member/create/form` — 普通会员等级添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user..UserBrokerage/createForm`
- 源码：`app/controller/admin/user/UserBrokerage.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `DELETE /sys/user/member/delete/:id` — 普通会员等级删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user..UserBrokerage/delete`
- 源码：`app/controller/admin/user/UserBrokerage.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/user/member/detail/:id` — 普通会员等级详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user..UserBrokerage/detail`
- 源码：`app/controller/admin/user/UserBrokerage.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/user/member/lst` — 普通会员等级列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user..UserBrokerage/getLst`
- 源码：`app/controller/admin/user/UserBrokerage.php` :: `getLst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `brokerage_name` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/user/member/options` — 普通会员等级筛选

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user..UserBrokerage/options`
- 源码：`app/controller/admin/user/UserBrokerage.php` :: `options()`
- 请求参数：
- `type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/user/member/update/:id` — 普通会员等级编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user..UserBrokerage/update`
- 源码：`app/controller/admin/user/UserBrokerage.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/user/member/update/:id/form` — 普通会员等级编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user..UserBrokerage/updateForm`
- 源码：`app/controller/admin/user/UserBrokerage.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/user/member_log` — 用户搜索记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserBill/getMembers`
- 源码：`app/controller/admin/user/UserBill.php` :: `getMembers()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `date` (query/body, 可选)
- `type` (query/body, 可选)
- `category` (query/body, 可选) 默认 $this->repository::CATEGORY_SYS_MEMBERS
- `uid` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/user/member_select_list` — 获取用户的等级下拉列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/getMemberLevelSelectList`
- 源码：`app/controller/admin/user/User.php` :: `getMemberLevelSelectList()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/user/news/push` — 用户发送图文

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/sendNews`
- 源码：`app/controller/admin/user/User.php` :: `sendNews()`
- 请求参数：
- `ids` (query/body, 可选)
- `news_id` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/user/order/:id` — 用户消费记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/order`
- 源码：`app/controller/admin/user/User.php` :: `order()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/user/promoter/count` — 分销员统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/promoterCount`
- 源码：`app/controller/admin/user/User.php` :: `promoterCount()`
- 请求参数：
- `keyword` (query/body, 可选)
- `date` (query/body, 可选)
- `brokerage_level` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/user/promoter/lst` — 分销员列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/promoterList`
- 源码：`app/controller/admin/user/User.php` :: `promoterList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `date` (query/body, 可选)
- `brokerage_level` (query/body, 可选)
- `uid` (query/body, 可选)
- `nickname` (query/body, 可选)
- `phone` (query/body, 可选)
- `real_name` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /sys/user/recharge/list` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserRecharge/getList`
- 源码：`app/controller/admin/user/UserRecharge.php` :: `getList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `date` (query/body, 可选)
- `paid` (query/body, 可选)
- `keyword` (query/body, 可选)
- `uid` (query/body, 可选)
- `phone` (query/body, 可选)
- `real_name` (query/body, 可选)
- `nickname` (query/body, 可选)
- `recharge_type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/user/recharge/total` — 统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserRecharge/total`
- 源码：`app/controller/admin/user/UserRecharge.php` :: `total()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /sys/user/register/config` — 保存注册配置

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/saveRegisterConfig`
- 源码：`app/controller/admin/user/User.php` :: `saveRegisterConfig()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/user/register/coupon` — 新人礼优惠券列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/getRegisterCoupon`
- 源码：`app/controller/admin/user/User.php` :: `getRegisterCoupon()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/user/search/:key` — search

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.UserSearch/search`
- 源码：`app/controller/admin/user/UserSearch.php` :: `search()`
- 请求参数：
- `key` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/user/search_log` — 用户搜索记录

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
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/user/search_log/export` — 用户搜索记录导出

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
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /sys/user/sign_log/:id` — 签到记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/sign_log`
- 源码：`app/controller/admin/user/User.php` :: `sign_log()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/user/spread/:id/form` — 修改分销员等级表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/spreadLevelForm`
- 源码：`app/controller/admin/user/User.php` :: `spreadLevelForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/user/spread/:id/save` — 修改分销员等级

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/spreadLevelSave`
- 源码：`app/controller/admin/user/User.php` :: `spreadLevelSave()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `brokerage_level` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/user/spread/clear/:uid` — 清除推广人

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/clearSpread`
- 源码：`app/controller/admin/user/User.php` :: `clearSpread()`
- 请求参数：
- `uid` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/user/spread/lst/:uid` — 推广人列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/spreadList`
- 源码：`app/controller/admin/user/User.php` :: `spreadList()`
- 请求参数：
- `uid` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `level` (query/body, 可选)
- `keyword` (query/body, 可选)
- `date` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /sys/user/spread/order/:uid` — 推广人订单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/spreadOrder`
- 源码：`app/controller/admin/user/User.php` :: `spreadOrder()`
- 请求参数：
- `uid` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `level` (query/body, 可选)
- `keyword` (query/body, 可选)
- `date` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /sys/user/spread_log/:id` — 推荐人修改记录

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/spreadLog`
- 源码：`app/controller/admin/user/User.php` :: `spreadLog()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/user/svip/:id` — 用户标签编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/svipUpdate`
- 源码：`app/controller/admin/user/User.php` :: `svipUpdate()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `is_svip` (query/body, 可选)
- `add_time` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/user/svip/:id/form` — systemUserSvipForm

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/svipForm`
- 源码：`app/controller/admin/user/User.php` :: `svipForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/user/svip/count_info` — 统计

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.Svip/countInfo`
- 源码：`app/controller/admin/user/Svip.php` :: `countInfo()`
- 请求参数：
- `pay_type` (query/body, 可选)
- `title` (query/body, 可选)
- `date` (query/body, 可选)
- `nickname` (query/body, 可选)
- `keyword` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/user/svip/order_lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.Svip/payList`
- 源码：`app/controller/admin/user/Svip.php` :: `payList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `pay_type` (query/body, 可选)
- `title` (query/body, 可选)
- `date` (query/body, 可选)
- `nickname` (query/body, 可选)
- `keyword` (query/body, 可选)
- `svip_type` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/user/svip/type/:id/form` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.Svip/updateTypeCreateForm`
- 源码：`app/controller/admin/user/Svip.php` :: `updateTypeCreateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/user/svip/type/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.groupData.GroupData/delete`
- 源码：`app/controller/admin/system/groupData/GroupData.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/user/svip/type/form` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.Svip/createTypeCreateForm`
- 源码：`app/controller/admin/user/Svip.php` :: `createTypeCreateForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/user/svip/type/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.Svip/getTypeLst`
- 源码：`app/controller/admin/user/Svip.php` :: `getTypeLst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/user/svip/type/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.groupData.GroupData/changeStatus`
- 源码：`app/controller/admin/system/groupData/GroupData.php` :: `changeStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/user/svip/update/:groupId/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.groupData.GroupData/update`
- 源码：`app/controller/admin/system/groupData/GroupData.php` :: `update()`
- 请求参数：
- `groupId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- `sort` (query/body, 可选) 默认 0
- `status` (query/body, 可选) 默认 0
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/user/update/:id` — 用户编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/update`
- 源码：`app/controller/admin/user/User.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `phone` (query/body, 可选)
- `card_id` (query/body, 可选)
- `mark` (query/body, 可选)
- `group_id` (query/body, 可选)
- `label_id` (query/body, 可选) 默认 [
- `is_promoter` (query/body, 可选) 默认 0
- `status` (query/body, 可选) 默认 0
- `member_level` (query/body, 可选) 默认 ''
- `extend_info` (query/body, 可选) 默认 [
- `promoter_switch` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/user/update/form/:id` — 用户编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.user.User/updateForm`
- 源码：`app/controller/admin/user/User.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `sys/version`

### `GET /sys/version` — version

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.Common/version`
- 源码：`app/controller/admin/Common.php` :: `version()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `sys/wechat`

### `GET /sys/wechat/menu` — 微信菜单配置

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.wechat.WechatMenu/info`
- 源码：`app/controller/admin/wechat/WechatMenu.php` :: `info()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/wechat/menu` — 保存微信菜单配置

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.wechat.WechatMenu/save`
- 源码：`app/controller/admin/wechat/WechatMenu.php` :: `save()`
- 请求参数：
- `button` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/wechat/news/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在；参数由 checkParams/Validate 自动补全
- 处理器：`admin.wechat.WechatNews/create`
- 源码：`app/controller/admin/wechat/WechatNews.php` :: `create()`
- 请求参数：
- `status` (body, 可选) 来自 checkParams/Validate（自动补全）
- `data` (body, 可选) 来自 checkParams/Validate（自动补全）
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/wechat/news/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.wechat.WechatNews/delete`
- 源码：`app/controller/admin/wechat/WechatNews.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/wechat/news/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.wechat.WechatNews/detail`
- 源码：`app/controller/admin/wechat/WechatNews.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/wechat/news/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.wechat.WechatNews/lst`
- 源码：`app/controller/admin/wechat/WechatNews.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `cate_name` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/wechat/news/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.wechat.WechatNews/update`
- 源码：`app/controller/admin/wechat/WechatNews.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `DELETE /sys/wechat/reply/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.wechat.WechatReply/delete`
- 源码：`app/controller/admin/wechat/WechatReply.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/wechat/reply/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.wechat.WechatReply/create`
- 源码：`app/controller/admin/wechat/WechatReply.php` :: `create()`
- 请求参数：
- `key` (query/body, 可选)
- `type` (query/body, 可选)
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/wechat/reply/detail/:id` — 详情

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.wechat.WechatReply/info`
- 源码：`app/controller/admin/wechat/WechatReply.php` :: `info()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `type` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /sys/wechat/reply/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.wechat.WechatReply/lst`
- 源码：`app/controller/admin/wechat/WechatReply.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /sys/wechat/reply/save/:key` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.wechat.WechatReply/save`
- 源码：`app/controller/admin/wechat/WechatReply.php` :: `save()`
- 请求参数：
- `key` (path, 必填) 路径参数
- `type` (query/body, 可选)
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/wechat/reply/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.wechat.WechatReply/changeStatus`
- 源码：`app/controller/admin/wechat/WechatReply.php` :: `changeStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/wechat/reply/update/:id` — 修改

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.wechat.WechatReply/update`
- 源码：`app/controller/admin/wechat/WechatReply.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `key` (query/body, 可选)
- `type` (query/body, 可选)
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/wechat/reply/upload/image` — 上传图片

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.wechat.WechatReply/uploadImage`
- 源码：`app/controller/admin/wechat/WechatReply.php` :: `uploadImage()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data 对象字段: src | 外层: {status,message,data}

### `POST /sys/wechat/reply/upload/voice` — 上传语音

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.wechat.WechatReply/uploadVoice`
- 源码：`app/controller/admin/wechat/WechatReply.php` :: `uploadVoice()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data 对象字段: src | 外层: {status,message,data}

### `POST /sys/wechat/template/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.wechat.TemplateMessage/create`
- 源码：`app/controller/admin/wechat/TemplateMessage.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/wechat/template/create/form` — 添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.wechat.TemplateMessage/createForm`
- 源码：`app/controller/admin/wechat/TemplateMessage.php` :: `createForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一响应包 {status,message,data?}；详见控制器实现 | 外层: {status,message,data}

### `DELETE /sys/wechat/template/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.wechat.TemplateMessage/delete`
- 源码：`app/controller/admin/wechat/TemplateMessage.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/wechat/template/lst` — 列表

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.wechat.TemplateMessage/lst`
- 源码：`app/controller/admin/wechat/TemplateMessage.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `status` (query/body, 可选)
- `keyword` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /sys/wechat/template/min/create` — 添加

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.wechat.TemplateMessage/create`
- 源码：`app/controller/admin/wechat/TemplateMessage.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/wechat/template/min/create/form` — 添加表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.wechat.TemplateMessage/createMinForm`
- 源码：`app/controller/admin/wechat/TemplateMessage.php` :: `createMinForm()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一响应包 {status,message,data?}；详见控制器实现 | 外层: {status,message,data}

### `DELETE /sys/wechat/template/min/delete/:id` — 删除

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.wechat.TemplateMessage/delete`
- 源码：`app/controller/admin/wechat/TemplateMessage.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/wechat/template/min/lst` — 列表 

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.wechat.TemplateMessage/minList`
- 源码：`app/controller/admin/wechat/TemplateMessage.php` :: `minList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `status` (query/body, 可选)
- `keyword` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /sys/wechat/template/min/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.wechat.TemplateMessage/switchStatus`
- 源码：`app/controller/admin/wechat/TemplateMessage.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/wechat/template/min/sync` — 同步

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.wechat.TemplateMessage/sync`
- 源码：`app/controller/admin/wechat/TemplateMessage.php` :: `sync()`
- 请求参数：
- `type` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/wechat/template/min/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.wechat.TemplateMessage/update`
- 源码：`app/controller/admin/wechat/TemplateMessage.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `tempid` (query/body, 可选)
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/wechat/template/min/update/:id/form` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.wechat.TemplateMessage/updateForm`
- 源码：`app/controller/admin/wechat/TemplateMessage.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /sys/wechat/template/status/:id` — 修改状态

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.wechat.TemplateMessage/switchStatus`
- 源码：`app/controller/admin/wechat/TemplateMessage.php` :: `switchStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/wechat/template/sync` — 同步

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.wechat.TemplateMessage/sync`
- 源码：`app/controller/admin/wechat/TemplateMessage.php` :: `sync()`
- 请求参数：
- `type` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /sys/wechat/template/update/:id` — 编辑

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.wechat.TemplateMessage/update`
- 源码：`app/controller/admin/wechat/TemplateMessage.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `tempid` (query/body, 可选)
- `status` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /sys/wechat/template/update/:id/form` — 编辑表单

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.wechat.TemplateMessage/updateForm`
- 源码：`app/controller/admin/wechat/TemplateMessage.php` :: `updateForm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

