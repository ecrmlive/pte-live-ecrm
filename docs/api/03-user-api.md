# 用户端 `/api/`

> 对照文档。置信度：high=427 stale=8 unresolved=1。先读 [ACCURACY.md](./ACCURACY.md)。

合计 **436** 条。

## `api/activity`

### `GET /api/activity/info/:id` — activityInfo

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Common/activityInfo`
- 源码：`app/controller/api/Common.php` :: `activityInfo()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/activity/lst/:id` — activityLst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Common/activityLst`
- 源码：`app/controller/api/Common.php` :: `activityLst()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `api/admin`

### `POST /api/admin/:merId/delivery/:id` — delivery

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreOrder/delivery`
- 源码：`app/controller/api/server/StoreOrder.php` :: `delivery()`
- 请求参数：
- `merId` (path, 必填) 路径参数
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
- `station_id` (query/body, 可选)
- `mark` (query/body, 可选)
- `cargo_weight` (query/body, 可选) 默认 0
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /api/admin/:merId/delivery/confirm/:id` — deliveryConfirm

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreOrder/deliveryConfirm`
- 源码：`app/controller/api/server/StoreOrder.php` :: `deliveryConfirm()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /api/admin/:merId/delivery/dispatch/:id` — deliveryDispatch

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreOrder/deliveryDispatch`
- 源码：`app/controller/api/server/StoreOrder.php` :: `deliveryDispatch()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- `service_id` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /api/admin/:merId/delivery/options` — options

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreOrder/options`
- 源码：`app/controller/api/server/StoreOrder.php` :: `options()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/admin/:merId/delivery/person` — deliveryPersonList

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreOrder/deliveryPersonList`
- 源码：`app/controller/api/server/StoreOrder.php` :: `deliveryPersonList()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /api/admin/:merId/delivery/updateDispatch/:id` — deliveryUpdateDispatch

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreOrder/deliveryUpdateDispatch`
- 源码：`app/controller/api/server/StoreOrder.php` :: `deliveryUpdateDispatch()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- `service_id` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /api/admin/:merId/delivery_config` — getDeliveryConfig

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreOrder/getDeliveryConfig`
- 源码：`app/controller/api/server/StoreOrder.php` :: `getDeliveryConfig()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/admin/:merId/delivery_options` — getDeliveryOptions

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreOrder/getDeliveryOptions`
- 源码：`app/controller/api/server/StoreOrder.php` :: `getDeliveryOptions()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/admin/:merId/dump_temp` — getFormData

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreOrder/getFormData`
- 源码：`app/controller/api/server/StoreOrder.php` :: `getFormData()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /api/admin/:merId/mark/:id` — mark

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreOrder/mark`
- 源码：`app/controller/api/server/StoreOrder.php` :: `mark()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- `remark` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /api/admin/:merId/mer_form` — getFormData

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreOrder/getFormData`
- 源码：`app/controller/api/server/StoreOrder.php` :: `getFormData()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /api/admin/:merId/offline/:id` — offline

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreOrder/offline`
- 源码：`app/controller/api/server/StoreOrder.php` :: `offline()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /api/admin/:merId/order/:id` — order

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreOrder/order`
- 源码：`app/controller/api/server/StoreOrder.php` :: `order()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/admin/:merId/order_list` — orderList

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreOrder/orderList`
- 源码：`app/controller/api/server/StoreOrder.php` :: `orderList()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `status` (query/body, 可选)
- `is_verify` (query/body, 可选)
- `store_name` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /api/admin/:merId/order_price` — orderDetail

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreOrder/orderDetail`
- 源码：`app/controller/api/server/StoreOrder.php` :: `orderDetail()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `start` (query/body, 可选) 默认 strtotime(date('Y-m'))
- `stop` (query/body, 可选) 默认 time()
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /api/admin/:merId/pay_number` — payNumber

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreOrder/payNumber`
- 源码：`app/controller/api/server/StoreOrder.php` :: `payNumber()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `start` (query/body, 可选) 默认 strtotime(date('Y-m'))
- `stop` (query/body, 可选) 默认 time()
- `month` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/admin/:merId/pay_price` — payPrice

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreOrder/payPrice`
- 源码：`app/controller/api/server/StoreOrder.php` :: `payPrice()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `start` (query/body, 可选) 默认 strtotime(date('Y-m'))
- `stop` (query/body, 可选) 默认 time()
- `month` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /api/admin/:merId/price/:id` — price

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreOrder/price`
- 源码：`app/controller/api/server/StoreOrder.php` :: `price()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- `total_price` (query/body, 可选)
- `pay_postage` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /api/admin/:merId/reservation/staffs` — staffList

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreOrder/staffList`
- 源码：`app/controller/api/server/StoreOrder.php` :: `staffList()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /api/admin/:merId/reservationconfig` — reservationConfig

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreOrder/reservationConfig`
- 源码：`app/controller/api/server/StoreOrder.php` :: `reservationConfig()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /api/admin/:merId/reservationdispatch/:id` — reservationDispatch

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreOrder/reservationDispatch`
- 源码：`app/controller/api/server/StoreOrder.php` :: `reservationDispatch()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- `staffs_id` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /api/admin/:merId/reservationreschedule/:id` — reservationReschedule

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreOrder/reservationReschedule`
- 源码：`app/controller/api/server/StoreOrder.php` :: `reservationReschedule()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- `order_type` (query/body, 可选)
- `reservation_date` (query/body, 可选)
- `real_name` (query/body, 可选)
- `user_phone` (query/body, 可选)
- `user_address` (query/body, 可选)
- `order_extend` (query/body, 可选)
- `part_start` (query/body, 可选)
- `part_end` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /api/admin/:merId/reservationupdateDispatch/:id` — reservationUpdateDispatch

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreOrder/reservationUpdateDispatch`
- 源码：`app/controller/api/server/StoreOrder.php` :: `reservationUpdateDispatch()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- `staffs_id` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /api/admin/:merId/reservationverify/:id` — reservationVerify

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreOrder/reservationVerify`
- 源码：`app/controller/api/server/StoreOrder.php` :: `reservationVerify()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /api/admin/:merId/statistics` — orderStatistics

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreOrder/orderStatistics`
- 源码：`app/controller/api/server/StoreOrder.php` :: `orderStatistics()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /api/admin/:merId/verify/:id` — verify

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreOrder/verify`
- 源码：`app/controller/api/server/StoreOrder.php` :: `verify()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- `verify_code` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `api/agreement`

### `GET /api/agreement/:key` — getAgree

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.Cache/getAgree`
- 源码：`app/controller/admin/system/Cache.php` :: `getAgree()`
- 请求参数：
- `key` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `api/agreement_lst`

### `GET /api/agreement_lst` — getKeyLst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.Cache/getKeyLst`
- 源码：`app/controller/admin/system/Cache.php` :: `getKeyLst()`
- 请求参数：
- `type` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `api/ajcaptcha`

### `GET /api/ajcaptcha` — ajcaptcha

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Auth/ajcaptcha`
- 源码：`app/controller/api/Auth.php` :: `ajcaptcha()`
- 请求参数：
- `captchaType` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `api/ajcheck`

### `POST /api/ajcheck` — ajcheck

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Auth/ajcheck`
- 源码：`app/controller/api/Auth.php` :: `ajcheck()`
- 请求参数：
- `token` (query/body, 可选)
- `pointJson` (query/body, 可选)
- `captchaType` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `api/appVersion`

### `GET /api/appVersion` — appVersion

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Common/appVersion`
- 源码：`app/controller/api/Common.php` :: `appVersion()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `api/article`

### `GET /api/article/category/lst` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.article.ArticleCategory/lst`
- 源码：`app/controller/api/article/ArticleCategory.php` :: `lst()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/article/detail/:id` — detail

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.article.Article/detail`
- 源码：`app/controller/api/article/Article.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/article/list` — list

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.article.Article/list`
- 源码：`app/controller/api/article/Article.php` :: `list()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/article/lst/:cid` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.article.Article/lst`
- 源码：`app/controller/api/article/Article.php` :: `lst()`
- 请求参数：
- `cid` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。


## `api/auth`

### `POST /api/auth` — authLogin

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Auth/authLogin`
- 源码：`app/controller/api/Auth.php` :: `authLogin()`
- 请求参数：
- `auth` (query/body, 可选)
- 返回：失败时 status=400, message 为错误信息 | 外层: {status,message,data}

### `POST /api/auth/app` — appAuth

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Auth/appAuth`
- 源码：`app/controller/api/Auth.php` :: `appAuth()`
- 请求参数：
- `userInfo` (query/body, 可选) 默认 [
- `authResult` (query/body, 可选) 默认 [
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /api/auth/apple` — appleAuth

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Auth/appleAuth`
- 源码：`app/controller/api/Auth.php` :: `appleAuth()`
- 请求参数：
- `openId` (query/body, 可选)
- `nickname` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /api/auth/login` — login

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Auth/login`
- 源码：`app/controller/api/Auth.php` :: `login()`
- 请求参数：
- `account` (query/body, 可选)
- `auth_token` (query/body, 可选)
- `password` (query/body, 可选)
- `spread` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /api/auth/mp` — mpAuth

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Auth/mpAuth`
- 源码：`app/controller/api/Auth.php` :: `mpAuth()`
- 请求参数：
- `code` (query/body, 可选)
- `cache_key` (query/body, 可选)
- `spread_spid` (query/body, 可选) 默认 0
- `spread_code` (query/body, 可选) 默认 ''
- `iv` (query/body, 可选) 默认 ''
- `encryptedData` (query/body, 可选) 默认 ''
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /api/auth/mp_login_type` — mpLoginType

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Auth/mpLoginType`
- 源码：`app/controller/api/Auth.php` :: `mpLoginType()`
- 请求参数：
- `code` (query/body, 可选)
- `spread` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /api/auth/mp_phone` — mpPhone

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Auth/mpPhone`
- 源码：`app/controller/api/Auth.php` :: `mpPhone()`
- 请求参数：
- `code` (query/body, 可选)
- `auth_token` (query/body, 可选)
- `iv` (query/body, 可选)
- `encryptedData` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /api/auth/register` — register

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Auth/register`
- 源码：`app/controller/api/Auth.php` :: `register()`
- 请求参数：
- `phone` (query/body, 可选)
- `sms_code` (query/body, 可选)
- `spread` (query/body, 可选)
- `pwd` (query/body, 可选)
- `auth_token` (query/body, 可选)
- `user_type` (query/body, 可选) 默认 'h5'
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /api/auth/smslogin` — smsLogin

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Auth/smsLogin`
- 源码：`app/controller/api/Auth.php` :: `smsLogin()`
- 请求参数：
- `phone` (query/body, 可选)
- `sms_code` (query/body, 可选)
- `spread` (query/body, 可选)
- `auth_token` (query/body, 可选)
- `user_type` (query/body, 可选) 默认 'h5'
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /api/auth/verify` — verify

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Auth/verify`
- 源码：`app/controller/api/Auth.php` :: `verify()`
- 请求参数：
- `phone` (query/body, 可选)
- `type` (query/body, 可选) 默认 'login'
- `captchaType` (query/body, 可选) 默认 'clickWord'
- `captchaVerification` (query/body, 可选) 默认 ''
- `token` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /api/auth/wechat` — auth

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Auth/auth`
- 源码：`app/controller/api/Auth.php` :: `auth()`
- 请求参数：
- `spread` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `api/broadcast`

### `GET /api/broadcast/hot` — hot

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.broadcast.BroadcastRoom/hot`
- 源码：`app/controller/api/store/broadcast/BroadcastRoom.php` :: `hot()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `mer_id` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/broadcast/lst` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.broadcast.BroadcastRoom/lst`
- 源码：`app/controller/api/store/broadcast/BroadcastRoom.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。


## `api/captcha`

### `GET /api/captcha` — getCaptcha

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Auth/getCaptcha`
- 源码：`app/controller/api/Auth.php` :: `getCaptcha()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `api/circleAgent`

### `POST /api/circleAgent/create` — create

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.circle.CircleAgent/create`
- 源码：`app/controller/api/circle/CircleAgent.php` :: `create()`
- 请求参数：
- `name` (query/body, 可选)
- `phone` (query/body, 可选)
- `qualification` (query/body, 可选)
- `remark` (query/body, 可选)
- `type` (query/body, 可选)
- `business_name` (query/body, 可选)
- `business_store_category` (query/body, 可选)
- `business_store_type` (query/body, 可选)
- `extend` (query/body, 可选) 默认 [
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /api/circleAgent/detail/:id` — detail

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.circle.CircleAgent/detail`
- 源码：`app/controller/api/circle/CircleAgent.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/circleAgent/list` — list

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.circle.CircleAgent/list`
- 源码：`app/controller/api/circle/CircleAgent.php` :: `list()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `type` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /api/circleAgent/revoke/:id` — revoke

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.circle.CircleAgent/revoke`
- 源码：`app/controller/api/circle/CircleAgent.php` :: `revoke()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /api/circleAgent/update/:id` — update

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.circle.CircleAgent/update`
- 源码：`app/controller/api/circle/CircleAgent.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `name` (query/body, 可选)
- `phone` (query/body, 可选)
- `qualification` (query/body, 可选)
- `remark` (query/body, 可选)
- `type` (query/body, 可选)
- `business_name` (query/body, 可选)
- `business_store_category` (query/body, 可选)
- `business_store_type` (query/body, 可选)
- `extend` (query/body, 可选) 默认 [
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `api/command`

### `GET /api/command/copy` — getCommand

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Common/getCommand`
- 源码：`app/controller/api/Common.php` :: `getCommand()`
- 请求参数：
- `key` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `api/common`

### `POST /api/common/base64` — get_image_base64

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Common/get_image_base64`
- 源码：`app/controller/api/Common.php` :: `get_image_base64()`
- 请求参数：
- `image` (query/body, 可选) 默认 ''
- `code` (query/body, 可选) 默认 ''
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /api/common/commuunity/hot_keyword` — hotKeyword

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Common/hotKeyword`
- 源码：`app/controller/api/Common.php` :: `hotKeyword()`
- 请求参数：
- `type` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/common/express` — express

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Common/express`
- 源码：`app/controller/api/Common.php` :: `express()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/common/feedback_type` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.FeedBackCategory/lst`
- 源码：`app/controller/api/user/FeedBackCategory.php` :: `lst()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/common/home` — home

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Common/home`
- 源码：`app/controller/api/Common.php` :: `home()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /api/common/hot_banner/:type` — hotBanner

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Common/hotBanner`
- 源码：`app/controller/api/Common.php` :: `hotBanner()`
- 请求参数：
- `type` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/common/hot_keyword` — hotKeyword

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Common/hotKeyword`
- 源码：`app/controller/api/Common.php` :: `hotKeyword()`
- 请求参数：
- `type` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/common/menus` — menus

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Common/menus`
- 源码：`app/controller/api/Common.php` :: `menus()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data 对象字段: global_theme, banner, menu | 外层: {status,message,data}

### `GET /api/common/pay_key/:key` — pay_key

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Common/pay_key`
- 源码：`app/controller/api/Common.php` :: `pay_key()`
- 请求参数：
- `key` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/common/recharge_quota` — userRechargeQuota

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Common/userRechargeQuota`
- 源码：`app/controller/api/Common.php` :: `userRechargeQuota()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /api/common/refund_message` — refundMessage

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Common/refundMessage`
- 源码：`app/controller/api/Common.php` :: `refundMessage()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /api/common/visit` — visit

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Common/visit`
- 源码：`app/controller/api/Common.php` :: `visit()`
- 请求参数：
- `page` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `api/community`

### `GET /api/community/:id/reply` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.community.CommunityReply/lst`
- 源码：`app/controller/api/community/CommunityReply.php` :: `lst()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/community/category/lst` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.community.CommunityCategory/lst`
- 源码：`app/controller/api/community/CommunityCategory.php` :: `lst()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/community/create` — create

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.community.Community/create`
- 源码：`app/controller/api/community/Community.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/community/delete/:id` — delete

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.community.Community/delete`
- 源码：`app/controller/api/community/Community.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/community/fans/:id` — setFocus

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.community.Community/setFocus`
- 源码：`app/controller/api/community/Community.php` :: `setFocus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/community/fans/lst` — getUserFans

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.community.Community/getUserFans`
- 源码：`app/controller/api/community/Community.php` :: `getUserFans()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/community/focus/lst` — getUserFocus

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.community.Community/getUserFocus`
- 源码：`app/controller/api/community/Community.php` :: `getUserFocus()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/community/focuslst` — focuslst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.community.Community/focuslst`
- 源码：`app/controller/api/community/Community.php` :: `focuslst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/community/hist_product/lst` — historyList

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.community.Community/historyList`
- 源码：`app/controller/api/community/Community.php` :: `historyList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/community/lst` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.community.Community/lst`
- 源码：`app/controller/api/community/Community.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `topic_id` (query/body, 可选)
- `is_hot` (query/body, 可选)
- `category_id` (query/body, 可选)
- `spu_id` (query/body, 可选)
- `search_type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/community/order/:id` — getSpuByOrder

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.community.Community/getSpuByOrder`
- 源码：`app/controller/api/community/Community.php` :: `getSpuByOrder()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/community/pay_product/lst` — payList

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.community.Community/payList`
- 源码：`app/controller/api/community/Community.php` :: `payList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/community/qrcode/:id` — qrcode

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.community.Community/qrcode`
- 源码：`app/controller/api/community/Community.php` :: `qrcode()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/community/rela_product/lst` — relationList

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.community.Community/relationList`
- 源码：`app/controller/api/community/Community.php` :: `relationList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/community/reply/create/:id` — create

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.community.CommunityReply/create`
- 源码：`app/controller/api/community/CommunityReply.php` :: `create()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `content` (query/body, 可选)
- `reply_id` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/community/reply/start/:id` — start

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.community.CommunityReply/start`
- 源码：`app/controller/api/community/CommunityReply.php` :: `start()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/community/show/:id` — show

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.community.Community/show`
- 源码：`app/controller/api/community/Community.php` :: `show()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/community/start/:id` — startCommunity

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.community.Community/startCommunity`
- 源码：`app/controller/api/community/Community.php` :: `startCommunity()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/community/start/lst` — getUserStartCommunity

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.community.Community/getUserStartCommunity`
- 源码：`app/controller/api/community/Community.php` :: `getUserStartCommunity()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/community/update/:id` — update

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.community.Community/update`
- 源码：`app/controller/api/community/Community.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/community/user/community/:id` — userCommunitylst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.community.Community/userCommunitylst`
- 源码：`app/controller/api/community/Community.php` :: `userCommunitylst()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/community/user/community_video/:id` — userCommunityVideolst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.community.Community/userCommunityVideolst`
- 源码：`app/controller/api/community/Community.php` :: `userCommunityVideolst()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `is_star` (query/body, 可选)
- `community_id` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/community/user/info/:id` — userInfo

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.community.Community/userInfo`
- 源码：`app/controller/api/community/Community.php` :: `userInfo()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/community/user_lst` — userList

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.community.Community/userList`
- 源码：`app/controller/api/community/Community.php` :: `userList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/community/video_lst` — videoShow

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.community.Community/videoShow`
- 源码：`app/controller/api/community/Community.php` :: `videoShow()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `id` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。


## `api/config`

### `GET /api/config` — config

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Common/config`
- 源码：`app/controller/api/Common.php` :: `config()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `api/copyright`

### `GET /api/copyright` — copyright

- 置信度：⛔ stale
- 说明：路由已登记，但 `app/controller/api/Common.php` 中不存在方法 `copyright`（原项目死路由/加密扩展/版本差异）。开发勿实现为有效接口，除非核实。
- 处理器：`api.Common/copyright`
- 源码：`app/controller/api/Common.php` :: `copyright()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一响应包 {status,message,data?}；详见控制器实现 | 外层: {status,message,data}


## `api/coupon`

### `GET /api/coupon/getlst` — getList

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreCoupon/getList`
- 源码：`app/controller/api/store/product/StoreCoupon.php` :: `getList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `type` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `product` (query/body, 可选)
- `region_id` (query/body, 可选)
- `is_pc` (query/body, 可选)
- `send_type` (query/body, 可选) 默认 0
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/coupon/list` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreCoupon/lst`
- 源码：`app/controller/api/store/product/StoreCoupon.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `statusTag` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/coupon/new_people` — newPeople

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreCoupon/newPeople`
- 源码：`app/controller/api/store/product/StoreCoupon.php` :: `newPeople()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/coupon/product` — coupon

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreCoupon/coupon`
- 源码：`app/controller/api/store/product/StoreCoupon.php` :: `coupon()`
- 请求参数：
- `ids` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /api/coupon/receive/:id` — receiveCoupon

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreCoupon/receiveCoupon`
- 源码：`app/controller/api/store/product/StoreCoupon.php` :: `receiveCoupon()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /api/coupon/store/:id` — merCoupon

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreCoupon/merCoupon`
- 源码：`app/controller/api/store/product/StoreCoupon.php` :: `merCoupon()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `all` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `api/delivery`

### `GET /api/delivery/order/:id` — deliveryOrderDetail

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.service.Delivery/orderDetail`
- 源码：`app/controller/api/store/service/Delivery.php` :: `orderDetail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /api/delivery/order/:id/confirm` — confirm

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.service.Delivery/confirm`
- 源码：`app/controller/api/store/service/Delivery.php` :: `confirm()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `mer_id` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /api/delivery/order/:id/mark` — mark

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.service.Delivery/mark`
- 源码：`app/controller/api/store/service/Delivery.php` :: `mark()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `remark` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /api/delivery/order/:id/receive` — receive

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.service.Delivery/receive`
- 源码：`app/controller/api/store/service/Delivery.php` :: `receive()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /api/delivery/order_lst` — deliveryOrderLst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.service.Delivery/order_lst`
- 源码：`app/controller/api/store/service/Delivery.php` :: `order_lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `type` (query/body, 可选) 默认 0
- `status` (query/body, 可选) 默认 0
- `delivery_keywords` (query/body, 可选)
- 返回：data 对象字段: count, list | 外层: {status,message,data}


## `api/discounts`

### `GET /api/discounts/lst` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.Discounts/lst`
- 源码：`app/controller/api/store/product/Discounts.php` :: `lst()`
- 请求参数：
- `product_id` (query/body, 可选) 默认 0
- `limit` (query/body, 可选) 默认 5
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `api/diy`

### `GET /api/diy` — diy

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Common/diy`
- 源码：`app/controller/api/Common.php` :: `diy()`
- 请求参数：
- `id` (query/body, 可选)
- `did` (query/body, 可选)
- 返回：统一响应包 {status,message,data?}；详见控制器实现 | 外层: {status,message,data}

### `GET /api/diy/assist` — assist

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Diy/assist`
- 源码：`app/controller/api/Diy.php` :: `assist()`
- 请求参数：
- `limit` (query/body, 可选)
- `region_id` (query/body, 可选)
- `mer_id` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/diy/broadcast` — broadcast

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Diy/broadcast`
- 源码：`app/controller/api/Diy.php` :: `broadcast()`
- 请求参数：
- `mer_id` (query/body, 可选)
- `limit` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/diy/category` — category

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Diy/category`
- 源码：`app/controller/api/Diy.php` :: `category()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/diy/community` — community

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Diy/community`
- 源码：`app/controller/api/Diy.php` :: `community()`
- 请求参数：
- `limit` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/diy/coupon` — coupon

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Diy/coupon`
- 源码：`app/controller/api/Diy.php` :: `coupon()`
- 请求参数：
- `limit` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `region_id` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/diy/fab` — fab

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Diy/fab`
- 源码：`app/controller/api/Diy.php` :: `fab()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/diy/group` — group

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Diy/group`
- 源码：`app/controller/api/Diy.php` :: `group()`
- 请求参数：
- `limit` (query/body, 可选)
- `region_id` (query/body, 可选)
- `mer_id` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/diy/hot_top` — hot_top

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Diy/hot_top`
- 源码：`app/controller/api/Diy.php` :: `hot_top()`
- 请求参数：
- `cate_pid` (query/body, 可选)
- `region_id` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/diy/presell` — presell

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Diy/presell`
- 源码：`app/controller/api/Diy.php` :: `presell()`
- 请求参数：
- `limit` (query/body, 可选)
- `region_id` (query/body, 可选)
- `mer_id` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/diy/productCategory` — productCategory

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Diy/productCategory`
- 源码：`app/controller/api/Diy.php` :: `productCategory()`
- 请求参数：
- `mer_id` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/diy/product_detail` — productDetail

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Diy/productDetail`
- 源码：`app/controller/api/Diy.php` :: `productDetail()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：失败时 status=400, message 为错误信息 | 外层: {status,message,data}

### `GET /api/diy/seckill` — seckill

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Diy/seckill`
- 源码：`app/controller/api/Diy.php` :: `seckill()`
- 请求参数：
- `mer_id` (query/body, 可选)
- `region_id` (query/body, 可选)
- `limit` (query/body, 可选)
- 返回：data 对象字段: list | 外层: {status,message,data}

### `GET /api/diy/spu` — spu

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Diy/spu`
- 源码：`app/controller/api/Diy.php` :: `spu()`
- 请求参数：
- `cate_pid` (query/body, 可选)
- `product_ids` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `mer_cate_id` (query/body, 可选)
- `order` (query/body, 可选)
- `latitude` (query/body, 可选)
- `longitude` (query/body, 可选)
- `limit` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/diy/store` — store

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Diy/store`
- 源码：`app/controller/api/Diy.php` :: `store()`
- 请求参数：
- `type_id` (query/body, 可选)
- `category_id` (query/body, 可选)
- `is_best` (query/body, 可选)
- `order` (query/body, 可选)
- `latitude` (query/body, 可选)
- `longitude` (query/body, 可选)
- `sort` (query/body, 可选)
- `region_id` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `limit` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `api/excel`

### `GET /api/excel/download/:id` — download

- 置信度：⛔ stale
- 说明：路由已登记，但 `app/controller/merchant/store/order/Order.php` 中不存在方法 `download`（原项目死路由/加密扩展/版本差异）。开发勿实现为有效接口，除非核实。
- 处理器：`merchant.store.order.Order/download`
- 源码：`app/controller/merchant/store/order/Order.php` :: `download()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一响应包 {status,message,data?}；详见控制器实现 | 外层: {status,message,data}


## `api/getVersion`

### `ANY /api/getVersion` — getVersion

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Common/getVersion`
- 源码：`app/controller/api/Common.php` :: `getVersion()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data 对象字段: version, host, system, php | 外层: {status,message,data}


## `api/has_service`

### `GET /api/has_service/:id` — hasService

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.service.Service/hasService`
- 源码：`app/controller/api/store/service/Service.php` :: `hasService()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `api/intention`

### `GET /api/intention/business` — business

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.merchant.MerchantIntention/business`
- 源码：`app/controller/api/store/merchant/MerchantIntention.php` :: `business()`
- 请求参数：
- `name` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/intention/cate` — cateLst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.merchant.MerchantIntention/cateLst`
- 源码：`app/controller/api/store/merchant/MerchantIntention.php` :: `cateLst()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/intention/circles` — circles

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.merchant.MerchantIntention/circles`
- 源码：`app/controller/api/store/merchant/MerchantIntention.php` :: `circles()`
- 请求参数：
- `type` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /api/intention/create` — create

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.merchant.MerchantIntention/create`
- 源码：`app/controller/api/store/merchant/MerchantIntention.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /api/intention/detail/:id` — detail

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.merchant.MerchantIntention/detail`
- 源码：`app/controller/api/store/merchant/MerchantIntention.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/intention/lst` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.merchant.MerchantIntention/lst`
- 源码：`app/controller/api/store/merchant/MerchantIntention.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/intention/type` — typeLst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.merchant.MerchantIntention/typeLst`
- 源码：`app/controller/api/store/merchant/MerchantIntention.php` :: `typeLst()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /api/intention/update/:id` — update

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.merchant.MerchantIntention/update`
- 源码：`app/controller/api/store/merchant/MerchantIntention.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `api/lbs`

### `GET /api/lbs/address` — lbs_address

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Common/lbs_address`
- 源码：`app/controller/api/Common.php` :: `lbs_address()`
- 请求参数：
- `region` (query/body, 可选)
- `address` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/lbs/geocoder` — lbs_geocoder

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Common/lbs_geocoder`
- 源码：`app/controller/api/Common.php` :: `lbs_geocoder()`
- 请求参数：
- `location` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `api/logout`

### `POST /api/logout` — logout

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Auth/logout`
- 源码：`app/controller/api/Auth.php` :: `logout()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `api/micro`

### `GET /api/micro` — micro

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Common/micro`
- 源码：`app/controller/api/Common.php` :: `micro()`
- 请求参数：
- `version` (query/body, 可选)
- `id` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `api/navigation`

### `GET /api/navigation` — getNavigation

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Common/getNavigation`
- 源码：`app/controller/api/Common.php` :: `getNavigation()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `api/notice`

### `ANY /api/notice/:type` — wechatNotify

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Common/notify`
- 源码：`app/controller/api/Common.php` :: `notify()`
- 请求参数：
- `type` (path, 必填) 路径参数
- 返回：统一响应包 {status,message,data?}；详见控制器实现 | 外层: {status,message,data}

### `ANY /api/notice/callback` — deliveryNotify

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Common/deliveryNotify`
- 源码：`app/controller/api/Common.php` :: `deliveryNotify()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一响应包 {status,message,data?}；详见控制器实现 | 外层: {status,message,data}

### `ANY /api/notice/mchNotify/:type` — mchNotify

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Common/mchNotify`
- 源码：`app/controller/api/Common.php` :: `mchNotify()`
- 请求参数：
- `type` (path, 必填) 路径参数
- 返回：统一响应包 {status,message,data?}；详见控制器实现 | 外层: {status,message,data}

### `ANY /api/notice/pay/alipay` — alipayNotify

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Common/alipayNotify`
- 源码：`app/controller/api/Common.php` :: `alipayNotify()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一响应包 {status,message,data?}；详见控制器实现 | 外层: {status,message,data}


## `api/open_screen`

### `GET /api/open_screen` — open_screen

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Common/open_screen`
- 源码：`app/controller/api/Common.php` :: `open_screen()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `api/order`

### `POST /api/order/cancel/:id` — cancelGroupOrder

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreOrder/cancelGroupOrder`
- 源码：`app/controller/api/store/order/StoreOrder.php` :: `cancelGroupOrder()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /api/order/cashier_order/:id` — getCashierOrder

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreOrder/getCashierOrder`
- 源码：`app/controller/api/store/order/StoreOrder.php` :: `getCashierOrder()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /api/order/check` — checkOrder

- 置信度：⛔ stale
- 说明：路由指向 `checkOrder()`，但当前源码 StoreOrder.php **不存在该方法**（仅有 v2CreateOrder/v2CheckOrder）。此 CRMEB 构建下请以 `/api/v2/order/*` 为准；勿按本条实现。
- 处理器：`api.store.order.StoreOrder/checkOrder`
- 源码：`app/controller/api/store/order/StoreOrder.php` :: `checkOrder()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：原路由可能不可用；见 doc_note

### `POST /api/order/create` — createOrder

- 置信度：⛔ stale
- 说明：路由指向 `createOrder()`，但当前源码 StoreOrder.php **不存在该方法**（仅有 v2CreateOrder/v2CheckOrder）。此 CRMEB 构建下请以 `/api/v2/order/*` 为准；勿按本条实现。
- 处理器：`api.store.order.StoreOrder/createOrder`
- 源码：`app/controller/api/store/order/StoreOrder.php` :: `createOrder()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：原路由可能不可用；见 doc_note

### `POST /api/order/del/:id` — del

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreOrder/del`
- 源码：`app/controller/api/store/order/StoreOrder.php` :: `del()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /api/order/delivery/:id` — getOrderDelivery

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreOrder/getOrderDelivery`
- 源码：`app/controller/api/store/order/StoreOrder.php` :: `getOrderDelivery()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/order/deliverySetings` — deliveryConfig

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreOrder/deliveryConfig`
- 源码：`app/controller/api/store/order/StoreOrder.php` :: `deliveryConfig()`
- 请求参数：
- `mer_id` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/order/deliveryStation/list` — deliveryStationList

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreOrder/deliveryStationList`
- 源码：`app/controller/api/store/order/StoreOrder.php` :: `deliveryStationList()`
- 请求参数：
- `switch_city` (query/body, 可选) 默认 0
- `switch_take` (query/body, 可选) 默认 0
- `status` (query/body, 可选) 默认 1
- `mer_id` (query/body, 可选)
- `address_id` (query/body, 可选)
- `name_and_address_search` (query/body, 可选)
- `longitude` (query/body, 可选)
- `latitude` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/order/deliveryTrack/:id` — deliveryTrack

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreOrder/deliveryTrack`
- 源码：`app/controller/api/store/order/StoreOrder.php` :: `deliveryTrack()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/order/detail/:id` — detail

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreOrder/detail`
- 源码：`app/controller/api/store/order/StoreOrder.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /api/order/express/:id` — express

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreOrder/express`
- 源码：`app/controller/api/store/order/StoreOrder.php` :: `express()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /api/order/group_order_detail/:id` — groupOrderDetail

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreOrder/groupOrderDetail`
- 源码：`app/controller/api/store/order/StoreOrder.php` :: `groupOrderDetail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/order/group_order_list` — groupOrderList

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreOrder/groupOrderList`
- 源码：`app/controller/api/store/order/StoreOrder.php` :: `groupOrderList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /api/order/list` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreOrder/lst`
- 源码：`app/controller/api/store/order/StoreOrder.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `status` (query/body, 可选)
- `store_name` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/order/number` — number

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreOrder/number`
- 源码：`app/controller/api/store/order/StoreOrder.php` :: `number()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /api/order/pay/:id` — groupOrderPay

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreOrder/groupOrderPay`
- 源码：`app/controller/api/store/order/StoreOrder.php` :: `groupOrderPay()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `type` (query/body, 可选)
- `is_points` (query/body, 可选)
- `return_url` (query/body, 可选)
- 返回：失败时 status=400, message 为错误信息 | 外层: {status,message,data}

### `POST /api/order/points/pay/:id` — groupOrderPay

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreOrder/groupOrderPay`
- 源码：`app/controller/api/store/order/StoreOrder.php` :: `groupOrderPay()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `type` (query/body, 可选)
- `is_points` (query/body, 可选)
- `return_url` (query/body, 可选)
- 返回：失败时 status=400, message 为错误信息 | 外层: {status,message,data}

### `POST /api/order/receipt/:id` — createReceipt

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreOrder/createReceipt`
- 源码：`app/controller/api/store/order/StoreOrder.php` :: `createReceipt()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `receipt_type` (query/body, 可选)
- `receipt_title` (query/body, 可选)
- `duty_paragraph` (query/body, 可选)
- `receipt_title_type` (query/body, 可选)
- `bank_name` (query/body, 可选)
- `bank_code` (query/body, 可选)
- `address` (query/body, 可选)
- `tel` (query/body, 可选)
- `email` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /api/order/self/cancel/:id` — cancelOrder

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreOrder/cancelOrder`
- 源码：`app/controller/api/store/order/StoreOrder.php` :: `cancelOrder()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /api/order/status/:id` — groupOrderStatus

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreOrder/groupOrderStatus`
- 源码：`app/controller/api/store/order/StoreOrder.php` :: `groupOrderStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /api/order/take/:id` — take

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreOrder/take`
- 源码：`app/controller/api/store/order/StoreOrder.php` :: `take()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /api/order/v3/check` — beforCheck

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.PointsOrder/beforCheck`
- 源码：`app/controller/api/store/order/PointsOrder.php` :: `beforCheck()`
- 请求参数：
- `cart_id` (query/body, 可选)
- `address_id` (query/body, 可选)
- `use_integral` (query/body, 可选)
- `use_coupon` (query/body, 可选)
- `takes` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/order/v3/create` — createOrder

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.PointsOrder/createOrder`
- 源码：`app/controller/api/store/order/PointsOrder.php` :: `createOrder()`
- 请求参数：
- `cart_id` (query/body, 可选)
- `address_id` (query/body, 可选)
- `use_integral` (query/body, 可选)
- `mark` (query/body, 可选)
- `pay_type` (query/body, 可选)
- `return_url` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/order/v3/pay/:id` — orderPay

- 置信度：⛔ stale
- 说明：路由已登记，但 `app/controller/api/store/order/PointsOrder.php` 中不存在方法 `orderPay`（原项目死路由/加密扩展/版本差异）。开发勿实现为有效接口，除非核实。
- 处理器：`api.store.order.PointsOrder/orderPay`
- 源码：`app/controller/api/store/order/PointsOrder.php` :: `orderPay()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/order/verify_code/:id` — verifyCode

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreOrder/verifyCode`
- 源码：`app/controller/api/store/order/StoreOrder.php` :: `verifyCode()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data 对象字段: qrcode | 外层: {status,message,data}


## `api/order_call_back`

### `ANY /api/order_call_back` — mchNotify

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Common/callBackNotify`
- 源码：`app/controller/api/Common.php` :: `callBackNotify()`
- 请求参数：
- `type` (query/body, 可选) 默认 ''
- `id` (query/body, 可选) 默认 0
- 返回：统一响应包 {status,message,data?}；详见控制器实现 | 外层: {status,message,data}


## `api/pay`

### `GET /api/pay/config` — payConfig

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreOrder/payConfig`
- 源码：`app/controller/api/store/order/StoreOrder.php` :: `payConfig()`
- 请求参数：
- `id` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `api/points`

### `GET /api/points/detail/:id` — detail

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.PointsProduct/detail`
- 源码：`app/controller/api/store/product/PointsProduct.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/points/home` — home

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.PointsProduct/home`
- 源码：`app/controller/api/store/product/PointsProduct.php` :: `home()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /api/points/lst` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.PointsProduct/lst`
- 源码：`app/controller/api/store/product/PointsProduct.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `scope` (query/body, 可选)
- `order` (query/body, 可选) 默认 'sort'
- `price` (query/body, 可选)
- `sales` (query/body, 可选)
- `keyword` (query/body, 可选)
- `cate_id` (query/body, 可选)
- `is_hot` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /api/points/order/deleate/:id` — del

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.PointsOrder/del`
- 源码：`app/controller/api/store/order/PointsOrder.php` :: `del()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/points/order/detail/:id` — detail

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.PointsOrder/detail`
- 源码：`app/controller/api/store/order/PointsOrder.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/points/order/lst` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.PointsOrder/lst`
- 源码：`app/controller/api/store/order/PointsOrder.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `pay_type` (query/body, 可选)
- `paid` (query/body, 可选)
- `status` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/points/order/take/:id` — take

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.PointsOrder/take`
- 源码：`app/controller/api/store/order/PointsOrder.php` :: `take()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/points/scope` — points_mall_scope

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.PointsProduct/points_mall_scope`
- 源码：`app/controller/api/store/product/PointsProduct.php` :: `points_mall_scope()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `api/presell`

### `POST /api/presell/pay/:id` — pay

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.PresellOrder/pay`
- 源码：`app/controller/api/store/order/PresellOrder.php` :: `pay()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `type` (query/body, 可选)
- `return_url` (query/body, 可选)
- 返回：失败时 status=400, message 为错误信息 | 外层: {status,message,data}


## `api/product`

### `GET /api/product/spu/active/category/:type` — activeCategory

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreSpu/activeCategory`
- 源码：`app/controller/api/store/product/StoreSpu.php` :: `activeCategory()`
- 请求参数：
- `type` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/product/spu/bag` — bag

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreSpu/bag`
- 源码：`app/controller/api/store/product/StoreSpu.php` :: `bag()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/product/spu/bag/recommend` — bagRecommend

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreSpu/bagRecommend`
- 源码：`app/controller/api/store/product/StoreSpu.php` :: `bagRecommend()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/product/spu/copy` — copy

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreSpu/copy`
- 源码：`app/controller/api/store/product/StoreSpu.php` :: `copy()`
- 请求参数：
- `id` (query/body, 可选)
- `product_type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/product/spu/coupon_product` — getProductByCoupon

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreSpu/getProductByCoupon`
- 源码：`app/controller/api/store/product/StoreSpu.php` :: `getProductByCoupon()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `cate_id` (query/body, 可选)
- `cate_pid` (query/body, 可选)
- `order` (query/body, 可选)
- `price_on` (query/body, 可选)
- `price_off` (query/body, 可选)
- `brand_id` (query/body, 可选)
- `pid` (query/body, 可选)
- `mer_cate_id` (query/body, 可选)
- `coupon_id` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/product/spu/get/:id` — get

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreSpu/get`
- 源码：`app/controller/api/store/product/StoreSpu.php` :: `get()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/product/spu/get_hot_ranking` — getHotRanking

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreSpu/getHotRanking`
- 源码：`app/controller/api/store/product/StoreSpu.php` :: `getHotRanking()`
- 请求参数：
- `cate_pid` (query/body, 可选)
- `limit` (query/body, 可选)
- `region_id` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/product/spu/hot/:type` — hot

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreSpu/hot`
- 源码：`app/controller/api/store/product/StoreSpu.php` :: `hot()`
- 请求参数：
- `type` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `common` (query/body, 可选)
- `mer_id` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/product/spu/hot_lst` — getHotList

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreProduct/getHotList`
- 源码：`app/controller/api/store/product/StoreProduct.php` :: `getHotList()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/product/spu/hot_top` — getHotTop

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreProduct/getHotTop`
- 源码：`app/controller/api/store/product/StoreProduct.php` :: `getHotTop()`
- 请求参数：
- `limit` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/product/spu/labels` — labelsLst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreSpu/labelsLst`
- 源码：`app/controller/api/store/product/StoreSpu.php` :: `labelsLst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `mer_id` (query/body, 可选)
- `labels` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/product/spu/local/:id` — local

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreSpu/local`
- 源码：`app/controller/api/store/product/StoreSpu.php` :: `local()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/product/spu/lst` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreSpu/lst`
- 源码：`app/controller/api/store/product/StoreSpu.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `cate_id` (query/body, 可选)
- `cate_pid` (query/body, 可选)
- `order` (query/body, 可选)
- `price_on` (query/body, 可选)
- `price_off` (query/body, 可选)
- `brand_id` (query/body, 可选)
- `pid` (query/body, 可选)
- `mer_cate_id` (query/body, 可选)
- `product_type` (query/body, 可选)
- `action` (query/body, 可选)
- `common` (query/body, 可选)
- `is_trader` (query/body, 可选) 默认 ''
- `product_ids` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `filter_params` (query/body, 可选)
- `mer_type_id` (query/body, 可选)
- `ids` (query/body, 可选)
- `latitude` (query/body, 可选)
- `longitude` (query/body, 可选)
- `store_type_id` (query/body, 可选)
- `store_label_id` (query/body, 可选)
- `delivery_type` (query/body, 可选)
- `region_id` (query/body, 可选) 默认 0
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/product/spu/merchant/:id` — merProductLst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreSpu/merProductLst`
- 源码：`app/controller/api/store/product/StoreSpu.php` :: `merProductLst()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `cate_id` (query/body, 可选)
- `order` (query/body, 可选)
- `price_on` (query/body, 可选)
- `price_off` (query/body, 可选)
- `brand_id` (query/body, 可选)
- `pid` (query/body, 可选)
- `mer_cate_id` (query/body, 可选)
- `product_type` (query/body, 可选) 默认 0
- `action` (query/body, 可选)
- `common` (query/body, 可选)
- `ids` (query/body, 可选)
- `product_ids` (query/body, 可选)
- `store_type_id` (query/body, 可选)
- `mer_store_label_id` (query/body, 可选)
- `delivery_type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/product/spu/params` — select

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreParams/select`
- 源码：`app/controller/api/store/product/StoreParams.php` :: `select()`
- 请求参数：
- `keyword` (query/body, 可选)
- `cate_id` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `mer_cate_id` (query/body, 可选)
- `cate_pid` (query/body, 可选)
- `is_pc` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/product/spu/params_value/:id` — getValue

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreParams/getValue`
- 源码：`app/controller/api/store/product/StoreParams.php` :: `getValue()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/product/spu/recommend` — recommend

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreSpu/recommend`
- 源码：`app/controller/api/store/product/StoreSpu.php` :: `recommend()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `common` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `latitude` (query/body, 可选)
- `longitude` (query/body, 可选)
- `region_id` (query/body, 可选) 默认 0
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。


## `api/refund`

### `POST /api/refund/apply/:id` — refund

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreRefundOrder/refund`
- 源码：`app/controller/api/store/order/StoreRefundOrder.php` :: `refund()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `type` (query/body, 可选)
- `refund_type` (query/body, 可选)
- `refund_price` (query/body, 可选)
- `num` (query/body, 可选)
- `ids` (query/body, 可选)
- `refund_message` (query/body, 可选)
- `mark` (query/body, 可选)
- `pics` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /api/refund/back_goods/:id` — back_goods

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreRefundOrder/back_goods`
- 源码：`app/controller/api/store/order/StoreRefundOrder.php` :: `back_goods()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `delivery_type` (query/body, 可选)
- `delivery_id` (query/body, 可选)
- `delivery_phone` (query/body, 可选)
- `delivery_mark` (query/body, 可选)
- `delivery_pics` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /api/refund/batch_product/:id` — batchProduct

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreRefundOrder/batchProduct`
- 源码：`app/controller/api/store/order/StoreRefundOrder.php` :: `batchProduct()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /api/refund/cancel/:id` — cancel

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreRefundOrder/cancel`
- 源码：`app/controller/api/store/order/StoreRefundOrder.php` :: `cancel()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /api/refund/compute` — compute

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreRefundOrder/compute`
- 源码：`app/controller/api/store/order/StoreRefundOrder.php` :: `compute()`
- 请求参数：
- `refund` (query/body, 可选)
- `order_id` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /api/refund/del/:id` — del

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreRefundOrder/del`
- 源码：`app/controller/api/store/order/StoreRefundOrder.php` :: `del()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /api/refund/detail/:id` — detail

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreRefundOrder/detail`
- 源码：`app/controller/api/store/order/StoreRefundOrder.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/refund/express/:id` — express

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreRefundOrder/express`
- 源码：`app/controller/api/store/order/StoreRefundOrder.php` :: `express()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /api/refund/list` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreRefundOrder/lst`
- 源码：`app/controller/api/store/order/StoreRefundOrder.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `type` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /api/refund/platform/:id` — platformIntervene

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreRefundOrder/platformIntervene`
- 源码：`app/controller/api/store/order/StoreRefundOrder.php` :: `platformIntervene()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /api/refund/product/:id` — product

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreRefundOrder/product`
- 源码：`app/controller/api/store/order/StoreRefundOrder.php` :: `product()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `ids` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `api/reply`

### `POST /api/reply/:id` — reply

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreReply/reply`
- 源码：`app/controller/api/store/product/StoreReply.php` :: `reply()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `comment` (query/body, 可选)
- `product_score` (query/body, 可选)
- `service_score` (query/body, 可选)
- `postage_score` (query/body, 可选)
- `pics` (query/body, 可选) 默认 [
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /api/reply/product/:id` — product

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreReply/product`
- 源码：`app/controller/api/store/product/StoreReply.php` :: `product()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `api/route`

### `GET /api/route/list` — list

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Route/list`
- 源码：`app/controller/api/Route.php` :: `list()`
- 请求参数：
- `type` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `api/scan_upload`

### `POST /api/scan_upload/image/:field/:token` — scanUploadImage

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Common/scanUploadImage`
- 源码：`app/controller/api/Common.php` :: `scanUploadImage()`
- 请求参数：
- `field` (path, 必填) 路径参数
- `token` (path, 必填) 路径参数
- `pid` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `name` (query/body, 可选)
- 返回：data 对象字段: src | 外层: {status,message,data}


## `api/script`

### `GET /api/script` — script

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Common/script`
- 源码：`app/controller/api/Common.php` :: `script()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `api/server`

### `POST /api/server/:merId/attr/create` — create

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreProductAttrTemplate/create`
- 源码：`app/controller/api/server/StoreProductAttrTemplate.php` :: `create()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/server/:merId/attr/delete` — batchDelete

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreProductAttrTemplate/batchDelete`
- 源码：`app/controller/api/server/StoreProductAttrTemplate.php` :: `batchDelete()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `ids` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/server/:merId/attr/detail/:id` — detail

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreProductAttrTemplate/detail`
- 源码：`app/controller/api/server/StoreProductAttrTemplate.php` :: `detail()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/server/:merId/attr/list` — getlist

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreProductAttrTemplate/getlist`
- 源码：`app/controller/api/server/StoreProductAttrTemplate.php` :: `getlist()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/server/:merId/attr/lst` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreProductAttrTemplate/lst`
- 源码：`app/controller/api/server/StoreProductAttrTemplate.php` :: `lst()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/server/:merId/attr/update/:id` — update

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreProductAttrTemplate/update`
- 源码：`app/controller/api/server/StoreProductAttrTemplate.php` :: `update()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/server/:merId/category/brandlist` — BrandList

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreCategory/BrandList`
- 源码：`app/controller/api/server/StoreCategory.php` :: `BrandList()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/server/:merId/category/create` — create

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreCategory/create`
- 源码：`app/controller/api/server/StoreCategory.php` :: `create()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/server/:merId/category/delete/:id` — delete

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreCategory/delete`
- 源码：`app/controller/api/server/StoreCategory.php` :: `delete()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/server/:merId/category/detail/:id` — detail

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreCategory/detail`
- 源码：`app/controller/api/server/StoreCategory.php` :: `detail()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/server/:merId/category/list` — getList

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreCategory/getList`
- 源码：`app/controller/api/server/StoreCategory.php` :: `getList()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/server/:merId/category/lst` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreCategory/lst`
- 源码：`app/controller/api/server/StoreCategory.php` :: `lst()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/server/:merId/category/select` — getTreeList

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreCategory/getTreeList`
- 源码：`app/controller/api/server/StoreCategory.php` :: `getTreeList()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/server/:merId/category/status/:id` — switchStatus

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreCategory/switchStatus`
- 源码：`app/controller/api/server/StoreCategory.php` :: `switchStatus()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/server/:merId/category/update/:id` — update

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreCategory/update`
- 源码：`app/controller/api/server/StoreCategory.php` :: `update()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/server/:merId/product/config` — config

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreProduct/config`
- 源码：`app/controller/api/server/StoreProduct.php` :: `config()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/server/:merId/product/create` — create

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreProduct/create`
- 源码：`app/controller/api/server/StoreProduct.php` :: `create()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `_see_CREATE_PARAMS` (body, 可选) 请求体字段见对应 Repository::CREATE_PARAMS
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/server/:merId/product/delete/:id` — delete

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreProduct/delete`
- 源码：`app/controller/api/server/StoreProduct.php` :: `delete()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/server/:merId/product/destory/:id` — destory

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreProduct/destory`
- 源码：`app/controller/api/server/StoreProduct.php` :: `destory()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/server/:merId/product/detail/:id` — detail

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreProduct/detail`
- 源码：`app/controller/api/server/StoreProduct.php` :: `detail()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/server/:merId/product/edit_cate/:id` — editCate

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreProduct/editCate`
- 源码：`app/controller/api/server/StoreProduct.php` :: `editCate()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- `cate_id` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/server/:merId/product/edit_mer_cate/:id` — editMerCate

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreProduct/editMerCate`
- 源码：`app/controller/api/server/StoreProduct.php` :: `editMerCate()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- `ids` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/server/:merId/product/good/:id` — updateGood

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreProduct/updateGood`
- 源码：`app/controller/api/server/StoreProduct.php` :: `updateGood()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- `is_good` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/server/:merId/product/lst` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreProduct/lst`
- 源码：`app/controller/api/server/StoreProduct.php` :: `lst()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `cate_id` (query/body, 可选)
- `keyword` (query/body, 可选)
- `mer_cate_id` (query/body, 可选)
- `is_gift_bag` (query/body, 可选)
- `status` (query/body, 可选)
- `us_status` (query/body, 可选)
- `product_id` (query/body, 可选)
- `mer_labels` (query/body, 可选)
- `order` (query/body, 可选) 默认 'sort'
- `type` (query/body, 可选) 默认 ''
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/server/:merId/product/restore/:id` — restore

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreProduct/restore`
- 源码：`app/controller/api/server/StoreProduct.php` :: `restore()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/server/:merId/product/status/:id` — switchStatus

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreProduct/switchStatus`
- 源码：`app/controller/api/server/StoreProduct.php` :: `switchStatus()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/server/:merId/product/title` — title

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreProduct/title`
- 源码：`app/controller/api/server/StoreProduct.php` :: `title()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/server/:merId/product/update/:id` — update

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreProduct/update`
- 源码：`app/controller/api/server/StoreProduct.php` :: `update()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- `_see_CREATE_PARAMS` (body, 可选) 请求体字段见对应 Repository::CREATE_PARAMS
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/server/:merId/product/value/:id` — getValue

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreProduct/getValue`
- 源码：`app/controller/api/server/StoreProduct.php` :: `getValue()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/server/:merId/product/value/:id` — setValue

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreProduct/setValue`
- 源码：`app/controller/api/server/StoreProduct.php` :: `setValue()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- `attr_value` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/server/:merId/refund/check/:id` — check

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreRefundOrder/check`
- 源码：`app/controller/api/server/StoreRefundOrder.php` :: `check()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /api/server/:merId/refund/compute` — compute

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreRefundOrder/compute`
- 源码：`app/controller/api/server/StoreRefundOrder.php` :: `compute()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `refund` (query/body, 可选)
- `order_id` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /api/server/:merId/refund/confirm/:id` — refundPrice

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreRefundOrder/refundPrice`
- 源码：`app/controller/api/server/StoreRefundOrder.php` :: `refundPrice()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- `status` (query/body, 可选)
- `fail_message` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /api/server/:merId/refund/create` — create

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreRefundOrder/create`
- 源码：`app/controller/api/server/StoreRefundOrder.php` :: `create()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `refund_message` (query/body, 可选)
- `refund_price` (query/body, 可选)
- `mer_mark` (query/body, 可选)
- `refund` (query/body, 可选)
- `order_id` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /api/server/:merId/refund/detail/:id` — detail

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreRefundOrder/detail`
- 源码：`app/controller/api/server/StoreRefundOrder.php` :: `detail()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/server/:merId/refund/express/:id` — express

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreRefundOrder/express`
- 源码：`app/controller/api/server/StoreRefundOrder.php` :: `express()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/server/:merId/refund/get/:id` — getRefundPrice

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreRefundOrder/getRefundPrice`
- 源码：`app/controller/api/server/StoreRefundOrder.php` :: `getRefundPrice()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/server/:merId/refund/lst` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreRefundOrder/lst`
- 源码：`app/controller/api/server/StoreRefundOrder.php` :: `lst()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `order_type` (query/body, 可选)
- `delivery_id` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /api/server/:merId/refund/mark/:id` — mark

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreRefundOrder/mark`
- 源码：`app/controller/api/server/StoreRefundOrder.php` :: `mark()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- `mer_mark` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /api/server/:merId/refund/status/:id` — switchStatus

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.StoreRefundOrder/switchStatus`
- 源码：`app/controller/api/server/StoreRefundOrder.php` :: `switchStatus()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- `mer_delivery_user` (query/body, 可选)
- `mer_delivery_address` (query/body, 可选)
- `phone` (query/body, 可选)
- `status` (query/body, 可选)
- `fail_message` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /api/server/:merId/template/create` — create

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.ShippingTemplate/create`
- 源码：`app/controller/api/server/ShippingTemplate.php` :: `create()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/server/:merId/template/delete` — batchDelete

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.ShippingTemplate/batchDelete`
- 源码：`app/controller/api/server/ShippingTemplate.php` :: `batchDelete()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `ids` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/server/:merId/template/detail/:id` — detail

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.ShippingTemplate/detail`
- 源码：`app/controller/api/server/ShippingTemplate.php` :: `detail()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/server/:merId/template/lst` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.ShippingTemplate/lst`
- 源码：`app/controller/api/server/ShippingTemplate.php` :: `lst()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `type` (query/body, 可选)
- `name` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/server/:merId/template/select` — getList

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.ShippingTemplate/getList`
- 源码：`app/controller/api/server/ShippingTemplate.php` :: `getList()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/server/:merId/template/update/:id` — update

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.server.ShippingTemplate/update`
- 源码：`app/controller/api/server/ShippingTemplate.php` :: `update()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。


## `api/service`

### `GET /api/service/history/:id` — chatHistory

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.service.Service/chatHistory`
- 源码：`app/controller/api/store/service/Service.php` :: `chatHistory()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/service/info/:id` — merchantInfo

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.service.Service/merchantInfo`
- 源码：`app/controller/api/store/service/Service.php` :: `merchantInfo()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/service/list` — getList

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.service.Service/getList`
- 源码：`app/controller/api/store/service/Service.php` :: `getList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /api/service/mark/:merId/:uid` — mark

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.service.Service/mark`
- 源码：`app/controller/api/store/service/Service.php` :: `mark()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `uid` (path, 必填) 路径参数
- `mark` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /api/service/mer_history/:merId/:id` — serviceHistory

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.service.Service/serviceHistory`
- 源码：`app/controller/api/store/service/Service.php` :: `serviceHistory()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /api/service/scan_login/:key` — scanLogin

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.service.Service/scanLogin`
- 源码：`app/controller/api/store/service/Service.php` :: `scanLogin()`
- 请求参数：
- `key` (path, 必填) 路径参数
- `service_id` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /api/service/user/:merId/:uid` — user

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.service.Service/user`
- 源码：`app/controller/api/store/service/Service.php` :: `user()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `uid` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/service/user_list/:merId` — serviceUserList

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.service.Service/serviceUserList`
- 源码：`app/controller/api/store/service/Service.php` :: `serviceUserList()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `api/staffs`

### `GET /api/staffs/order/:id` — orderDetail

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.service.Staffs/orderDetail`
- 源码：`app/controller/api/store/service/Staffs.php` :: `orderDetail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /api/staffs/order/:id/check` — checkIn

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.service.Staffs/checkIn`
- 源码：`app/controller/api/store/service/Staffs.php` :: `checkIn()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `clock_in_info` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /api/staffs/order/:id/dispatch` — reservationDispatch

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.service.Staffs/reservationDispatch`
- 源码：`app/controller/api/store/service/Staffs.php` :: `reservationDispatch()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /api/staffs/order/:id/mark` — mark

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.service.Staffs/mark`
- 源码：`app/controller/api/store/service/Staffs.php` :: `mark()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `remark` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /api/staffs/order/:id/trace` — addTrace

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.service.Staffs/addTrace`
- 源码：`app/controller/api/store/service/Staffs.php` :: `addTrace()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `reservation_service_voucher` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /api/staffs/order/:id/verifier` — verify

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.service.Staffs/verify`
- 源码：`app/controller/api/store/service/Staffs.php` :: `verify()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `mer_id` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /api/staffs/order_lst` — order_lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.service.Staffs/order_lst`
- 源码：`app/controller/api/store/service/Staffs.php` :: `order_lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `paid` (query/body, 可选) 默认 1
- `status` (query/body, 可选) 默认 ''
- `assigned` (query/body, 可选) 默认 ''
- `is_del` (query/body, 可选) 默认 0
- `filter_product` (query/body, 可选) 默认 4
- `order_type` (query/body, 可选) 默认 0
- `store_name` (query/body, 可选) 默认 ''
- `date` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/staffs/reservation/config` — reservationConfig

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.service.Staffs/reservationConfig`
- 源码：`app/controller/api/store/service/Staffs.php` :: `reservationConfig()`
- 请求参数：
- `mer_id` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `api/store`

### `POST /api/store/certificate/:merId` — getMerCertificate

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Auth/getMerCertificate`
- 源码：`app/controller/api/Auth.php` :: `getMerCertificate()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `key` (query/body, 可选)
- `code` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/store/expr/temps` — getExportTemp

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.serve.Export/getExportTemp`
- 源码：`app/controller/admin/system/serve/Export.php` :: `getExportTemp()`
- 请求参数：
- `com` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/store/group/options` — options

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.StoreGroup/options`
- 源码：`app/controller/api/store/StoreGroup.php` :: `options()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/store/group/recommend` — recommendList

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.StoreGroup/recommendList`
- 源码：`app/controller/api/store/StoreGroup.php` :: `recommendList()`
- 请求参数：
- `latitude` (query/body, 可选)
- `longitude` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/store/merchant/category/lst/:id` — categoryList

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.merchant.Merchant/categoryList`
- 源码：`app/controller/api/store/merchant/Merchant.php` :: `categoryList()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/merchant/detail/0` — systemDetail

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.merchant.Merchant/systemDetail`
- 源码：`app/controller/api/store/merchant/Merchant.php` :: `systemDetail()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/merchant/detail/:id` — detail

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.merchant.Merchant/detail`
- 源码：`app/controller/api/store/merchant/Merchant.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/merchant/local` — localLst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.merchant.Merchant/localLst`
- 源码：`app/controller/api/store/merchant/Merchant.php` :: `localLst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `order` (query/body, 可选)
- `is_best` (query/body, 可选)
- `location` (query/body, 可选)
- `category_id` (query/body, 可选)
- `type_id` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/merchant/localDetail/:id` — localDetail

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.merchant.Merchant/localDetail`
- 源码：`app/controller/api/store/merchant/Merchant.php` :: `localDetail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `latitude` (query/body, 可选)
- `longitude` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/merchant/lst` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.merchant.Merchant/lst`
- 源码：`app/controller/api/store/merchant/Merchant.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `order` (query/body, 可选)
- `is_best` (query/body, 可选)
- `location` (query/body, 可选)
- `category_id` (query/body, 可选)
- `type_id` (query/body, 可选)
- `is_trader` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/merchant/product/lst/:id` — productList

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.merchant.Merchant/productList`
- 源码：`app/controller/api/store/merchant/Merchant.php` :: `productList()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `order` (query/body, 可选)
- `mer_cate_id` (query/body, 可选)
- `cate_id` (query/body, 可选)
- `price_on` (query/body, 可选)
- `price_off` (query/body, 可选)
- `brand_id` (query/body, 可选)
- `pid` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/merchant/qrcode/:id` — qrcode

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.merchant.Merchant/qrcode`
- 源码：`app/controller/api/store/merchant/Merchant.php` :: `qrcode()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/product/assist/count` — userCount

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreProductAssist/userCount`
- 源码：`app/controller/api/store/product/StoreProductAssist.php` :: `userCount()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /api/store/product/assist/create/:id` — create

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreProductAssistSet/create`
- 源码：`app/controller/api/store/product/StoreProductAssistSet.php` :: `create()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/product/assist/detail/:id` — detail

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreProductAssistSet/detail`
- 源码：`app/controller/api/store/product/StoreProductAssistSet.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/product/assist/lst` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreProductAssist/lst`
- 源码：`app/controller/api/store/product/StoreProductAssist.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `type` (query/body, 可选)
- `star` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `region_id` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/store/product/assist/set/:id` — set

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreProductAssistSet/set`
- 源码：`app/controller/api/store/product/StoreProductAssistSet.php` :: `set()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/store/product/assist/set/delete/:id` — delete

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreProductAssistSet/delete`
- 源码：`app/controller/api/store/product/StoreProductAssistSet.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/product/assist/set/lst` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreProductAssistSet/lst`
- 源码：`app/controller/api/store/product/StoreProductAssistSet.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/product/assist/share/:id` — shareNum

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreProductAssistSet/shareNum`
- 源码：`app/controller/api/store/product/StoreProductAssistSet.php` :: `shareNum()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/product/assist/user/:id` — userList

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreProductAssistSet/userList`
- 源码：`app/controller/api/store/product/StoreProductAssistSet.php` :: `userList()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/product/bag/explain` — getBagExplain

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreProduct/getBagExplain`
- 源码：`app/controller/api/store/product/StoreProduct.php` :: `getBagExplain()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/product/brand/lst` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreBrand/lst`
- 源码：`app/controller/api/store/product/StoreBrand.php` :: `lst()`
- 请求参数：
- `keyword` (query/body, 可选)
- `cate_id` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `mer_cate_id` (query/body, 可选)
- `pid` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/product/cate_hot` — cateHotList

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreProduct/cateHotList`
- 源码：`app/controller/api/store/product/StoreProduct.php` :: `cateHotList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `product_id` (query/body, 可选)
- `cate_id` (query/body, 可选)
- `cate_pid` (query/body, 可选)
- `filter` (query/body, 可选) 默认 ''
- `mer_id` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/product/category` — children

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreCategory/children`
- 源码：`app/controller/api/store/product/StoreCategory.php` :: `children()`
- 请求参数：
- `pid` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/product/category/hotranking` — cateHotRanking

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreCategory/cateHotRanking`
- 源码：`app/controller/api/store/product/StoreCategory.php` :: `cateHotRanking()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/product/category/lst` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreCategory/lst`
- 源码：`app/controller/api/store/product/StoreCategory.php` :: `lst()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/product/detail/:id` — detail

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreProduct/detail`
- 源码：`app/controller/api/store/product/StoreProduct.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/product/get_attr_value/:id` — getAttrValue

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreProduct/getAttrValue`
- 源码：`app/controller/api/store/product/StoreProduct.php` :: `getAttrValue()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/product/get_spec/:id` — getSpec

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreProduct/getSpec`
- 源码：`app/controller/api/store/product/StoreProduct.php` :: `getSpec()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/product/good_list/:id` — getGoodList

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreProduct/getGoodList`
- 源码：`app/controller/api/store/product/StoreProduct.php` :: `getGoodList()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/store/product/group/cancel` — cancel

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreProductGroup/cancel`
- 源码：`app/controller/api/store/product/StoreProductGroup.php` :: `cancel()`
- 请求参数：
- `group_buying_id` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /api/store/product/group/category` — category

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreProductGroup/category`
- 源码：`app/controller/api/store/product/StoreProductGroup.php` :: `category()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/product/group/count` — userCount

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreProductGroup/userCount`
- 源码：`app/controller/api/store/product/StoreProductGroup.php` :: `userCount()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/product/group/detail/:id` — detail

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreProductGroup/detail`
- 源码：`app/controller/api/store/product/StoreProductGroup.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/product/group/get/:id` — groupBuying

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreProductGroup/groupBuying`
- 源码：`app/controller/api/store/product/StoreProductGroup.php` :: `groupBuying()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/product/group/lst` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreProductGroup/lst`
- 源码：`app/controller/api/store/product/StoreProductGroup.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `active_type` (query/body, 可选) 默认 1
- `store_category_id` (query/body, 可选)
- `star` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `region_id` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/product/guarantee/:id` — guaranteeTemplate

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreProduct/guaranteeTemplate`
- 源码：`app/controller/api/store/product/StoreProduct.php` :: `guaranteeTemplate()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/store/product/increase_take` — setIncreaseTake

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreProduct/setIncreaseTake`
- 源码：`app/controller/api/store/product/StoreProduct.php` :: `setIncreaseTake()`
- 请求参数：
- `product_id` (query/body, 可选)
- `unique` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/product/presell/agree` — getAgree

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreProductPresell/getAgree`
- 源码：`app/controller/api/store/product/StoreProductPresell.php` :: `getAgree()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/product/presell/detail/:id` — detail

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreProductPresell/detail`
- 源码：`app/controller/api/store/product/StoreProductPresell.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/product/presell/lst` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreProductPresell/lst`
- 源码：`app/controller/api/store/product/StoreProductPresell.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `type` (query/body, 可选) 默认 4
- `star` (query/body, 可选)
- `mer_id` (query/body, 可选)
- `region_id` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/product/preview` — preview

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreProduct/preview`
- 源码：`app/controller/api/store/product/StoreProduct.php` :: `preview()`
- 请求参数：
- `key` (query/body, 可选)
- `id` (query/body, 可选)
- `product_type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/product/price_rule/:id` — priceRule

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreProduct/priceRule`
- 源码：`app/controller/api/store/product/StoreProduct.php` :: `priceRule()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/product/qrcode/:id` — qrcode

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreProduct/qrcode`
- 源码：`app/controller/api/store/product/StoreProduct.php` :: `qrcode()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `type` (query/body, 可选)
- `product_type` (query/body, 可选) 默认 0
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/product/recommendProduct` — recommendProduct

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreProduct/recommendProduct`
- 源码：`app/controller/api/store/product/StoreProduct.php` :: `recommendProduct()`
- 请求参数：
- `recommend_num` (query/body, 可选) 默认 1
- `product_id` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/product/reply/lst/:id` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreReply/lst`
- 源码：`app/controller/api/store/product/StoreReply.php` :: `lst()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/store/product/reservation/checkRange` — checkRange

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StorePrdouctReservation/checkRange`
- 源码：`app/controller/api/store/product/StorePrdouctReservation.php` :: `checkRange()`
- 请求参数：
- `address_id` (query/body, 可选)
- `mer_id` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/product/reservation/getDay/:id` — showDay

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StorePrdouctReservation/showDay`
- 源码：`app/controller/api/store/product/StorePrdouctReservation.php` :: `showDay()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `date` (query/body, 可选) 默认 date('Y-m')
- `sku_id` (query/body, 可选) 默认 0
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/product/reservation/getMonth/:id` — showMonth

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StorePrdouctReservation/showMonth`
- 源码：`app/controller/api/store/product/StorePrdouctReservation.php` :: `showMonth()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `sku_id` (query/body, 可选) 默认 ''
- `date` (query/body, 可选) 默认 date('Y-m')
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/product/seckill/detail/:id` — detail

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreProductSeckill/detail`
- 源码：`app/controller/api/store/product/StoreProductSeckill.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `seckill_time_id` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/product/seckill/lst` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreProductSeckill/lst`
- 源码：`app/controller/api/store/product/StoreProductSeckill.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `mer_id` (query/body, 可选)
- `start_time` (query/body, 可选)
- `end_time` (query/body, 可选)
- `region_id` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/product/seckill/select` — select

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreProductSeckill/select`
- 源码：`app/controller/api/store/product/StoreProductSeckill.php` :: `select()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/store/product/show/:id` — show

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.product.StoreProduct/show`
- 源码：`app/controller/api/store/product/StoreProduct.php` :: `show()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `ANY /api/store/test` — test

- 置信度：❓ unresolved
- 说明：未能可靠映射到控制器，开发时勿直接照抄，需对照 route 源码
- 处理器：`api.Test/test`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。


## `api/subscribe`

### `GET /api/subscribe` — subscribe

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Common/subscribe`
- 源码：`app/controller/api/Common.php` :: `subscribe()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data 对象字段: subscribe | 外层: {status,message,data}


## `api/svip`

### `GET /api/svip/coupon_lst` — svipCoupon

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.Svip/svipCoupon`
- 源码：`app/controller/api/user/Svip.php` :: `svipCoupon()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /api/svip/coupon_receive/:id` — receiveCoupon

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.Svip/receiveCoupon`
- 源码：`app/controller/api/user/Svip.php` :: `receiveCoupon()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /api/svip/pay/:id` — createOrder

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.Svip/createOrder`
- 源码：`app/controller/api/user/Svip.php` :: `createOrder()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `pay_type` (query/body, 可选)
- `return_url` (query/body, 可选)
- 返回：失败时 status=400, message 为错误信息 | 外层: {status,message,data}

### `GET /api/svip/pay_lst` — getTypeLst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.Svip/getTypeLst`
- 源码：`app/controller/api/user/Svip.php` :: `getTypeLst()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data 对象字段: def, list | 外层: {status,message,data}

### `GET /api/svip/product_lst` — svipProductList

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.Svip/svipProductList`
- 源码：`app/controller/api/user/Svip.php` :: `svipProductList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/svip/user_info` — svipUserInfo

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.Svip/svipUserInfo`
- 源码：`app/controller/api/user/Svip.php` :: `svipUserInfo()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `api/system`

### `GET /api/system/city/lst` — getlist

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.shipping.City/getlist`
- 源码：`app/controller/merchant/store/shipping/City.php` :: `getlist()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/system/form/detail/:id` — detail

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.form.Form/detail`
- 源码：`app/controller/api/store/form/Form.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：失败时 status=400, message 为错误信息 | 外层: {status,message,data}

### `GET /api/system/form/info/:form_id` — getFormInfo

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.form.Form/getFormInfo`
- 源码：`app/controller/api/store/form/Form.php` :: `getFormInfo()`
- 请求参数：
- `form_id` (path, 必填) 路径参数
- 返回：失败时 status=400, message 为错误信息 | 外层: {status,message,data}

### `GET /api/system/form/lst` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.form.Form/lst`
- 源码：`app/controller/api/store/form/Form.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /api/system/form/share_posters/:id` — getSharePosters

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.form.Form/getSharePosters`
- 源码：`app/controller/api/store/form/Form.php` :: `getSharePosters()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `type` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `api/upload`

### `POST /api/upload/certificate/:field` — uploadCertificate

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Common/uploadCertificate`
- 源码：`app/controller/api/Common.php` :: `uploadCertificate()`
- 请求参数：
- `field` (path, 必填) 路径参数
- 返回：data 对象字段: path | 外层: {status,message,data}

### `POST /api/upload/image/:field` — uploadImage

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Common/uploadImage`
- 源码：`app/controller/api/Common.php` :: `uploadImage()`
- 请求参数：
- `field` (path, 必填) 路径参数
- `name` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `POST /api/upload/video` — uploadVideo

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.system.attachment.Attachment/uploadVideo`
- 源码：`app/controller/admin/system/attachment/Attachment.php` :: `uploadVideo()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data 对象字段: src, attachment_id | 外层: {status,message,data}


## `api/user`

### `GET /api/user` — userInfo

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Auth/userInfo`
- 源码：`app/controller/api/Auth.php` :: `userInfo()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/user/account` — account

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.User/account`
- 源码：`app/controller/api/user/User.php` :: `account()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/user/address/create` — create

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.UserAddress/create`
- 源码：`app/controller/api/user/UserAddress.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/user/address/delete/:id` — delete

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.UserAddress/delete`
- 源码：`app/controller/api/user/UserAddress.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/user/address/detail/:id` — detail

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.UserAddress/detail`
- 源码：`app/controller/api/user/UserAddress.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/user/address/lst` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.UserAddress/lst`
- 源码：`app/controller/api/user/UserAddress.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/user/address/update/:id` — editDefault

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.UserAddress/editDefault`
- 源码：`app/controller/api/user/UserAddress.php` :: `editDefault()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/user/bill` — bill

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.User/bill`
- 源码：`app/controller/api/user/User.php` :: `bill()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/user/binding` — binding

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.User/binding`
- 源码：`app/controller/api/user/User.php` :: `binding()`
- 请求参数：
- `phone` (query/body, 可选)
- `sms_code` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/user/brokerage/all` — brokerage_all

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.User/brokerage_all`
- 源码：`app/controller/api/user/User.php` :: `brokerage_all()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/user/brokerage/info` — brokerage_info

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.User/brokerage_info`
- 源码：`app/controller/api/user/User.php` :: `brokerage_info()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/user/brokerage/notice` — notice

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.User/notice`
- 源码：`app/controller/api/user/User.php` :: `notice()`
- 请求参数：
- `type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/user/brokerage_list` — brokerage_list

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.User/brokerage_list`
- 源码：`app/controller/api/user/User.php` :: `brokerage_list()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `start` (query/body, 可选)
- `stop` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/user/brokerage_top` — brokerage_top

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.User/brokerage_top`
- 源码：`app/controller/api/user/User.php` :: `brokerage_top()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/user/cancel` — cancel

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Auth/cancel`
- 源码：`app/controller/api/Auth.php` :: `cancel()`
- 请求参数：
- `key` (query/body, 可选)
- 返回：失败时 status=400, message 为错误信息 | 外层: {status,message,data}

### `POST /api/user/cart/again` — again

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreCart/again`
- 源码：`app/controller/api/store/order/StoreCart.php` :: `again()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/user/cart/batchCreate` — batchCreate

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreCart/batchCreate`
- 源码：`app/controller/api/store/order/StoreCart.php` :: `batchCreate()`
- 请求参数：
- `discount_id` (query/body, 可选)
- `is_new` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/user/cart/change/:id` — change

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreCart/change`
- 源码：`app/controller/api/store/order/StoreCart.php` :: `change()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `cart_num` (query/body, 可选)
- `product_attr_unique` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/user/cart/check/:id` — checkCerate

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreCart/check`
- 源码：`app/controller/api/store/order/StoreCart.php` :: `check()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/user/cart/clear` — clear

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreCart/clear`
- 源码：`app/controller/api/store/order/StoreCart.php` :: `clear()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/user/cart/count` — cartCount

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreCart/cartCount`
- 源码：`app/controller/api/store/order/StoreCart.php` :: `cartCount()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/user/cart/create` — create

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreCart/create`
- 源码：`app/controller/api/store/order/StoreCart.php` :: `create()`
- 请求参数：
- `source` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/user/cart/delete` — batchDelete

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreCart/batchDelete`
- 源码：`app/controller/api/store/order/StoreCart.php` :: `batchDelete()`
- 请求参数：
- `cart_id` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/user/cart/lst` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreCart/lst`
- 源码：`app/controller/api/store/order/StoreCart.php` :: `lst()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/user/change/info` — updateBaseInfo

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.User/updateBaseInfo`
- 源码：`app/controller/api/user/User.php` :: `updateBaseInfo()`
- 请求参数：
- `nickname` (query/body, 可选)
- `avatar` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/user/change/password` — changePassword

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.User/changePassword`
- 源码：`app/controller/api/user/User.php` :: `changePassword()`
- 请求参数：
- `repassword` (query/body, 可选)
- `password` (query/body, 可选)
- `sms_code` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/user/change/phone` — changePhone

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.User/changePhone`
- 源码：`app/controller/api/user/User.php` :: `changePhone()`
- 请求参数：
- `phone` (query/body, 可选)
- `sms_code` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/user/change_pwd` — changePassword

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Auth/changePassword`
- 源码：`app/controller/api/Auth.php` :: `changePassword()`
- 请求参数：
- `phone` (query/body, 可选)
- `sms_code` (query/body, 可选)
- `pwd` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /api/user/extract/banklst` — bankLst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.UserExtract/bankLst`
- 源码：`app/controller/api/user/UserExtract.php` :: `bankLst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/user/extract/create` — create

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.UserExtract/create`
- 源码：`app/controller/api/user/UserExtract.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/user/extract/detail/:id` — detail

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.UserExtract/detail`
- 源码：`app/controller/api/user/UserExtract.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/user/extract/history_bank` — historyBank

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.UserExtract/historyBank`
- 源码：`app/controller/api/user/UserExtract.php` :: `historyBank()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/user/extract/lst` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.UserExtract/lst`
- 源码：`app/controller/api/user/UserExtract.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `status` (query/body, 可选)
- `start` (query/body, 可选)
- `stop` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/user/feedback` — feedback

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.Feedback/feedback`
- 源码：`app/controller/api/user/Feedback.php` :: `feedback()`
- 请求参数：
- `type` (query/body, 可选)
- `content` (query/body, 可选)
- `images` (query/body, 可选) 默认 [
- `realname` (query/body, 可选)
- `contact` (query/body, 可选)
- `status` (query/body, 可选) 默认 0
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/user/feedback/detail/:id` — detail

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.Feedback/detail`
- 源码：`app/controller/api/user/Feedback.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/user/feedback/list` — feedbackList

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.Feedback/feedbackList`
- 源码：`app/controller/api/user/Feedback.php` :: `feedbackList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `DELETE /api/user/fields/delete` — delete

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.UserFields/delete`
- 源码：`app/controller/api/user/UserFields.php` :: `delete()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/user/fields/info` — info

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.UserFields/info`
- 源码：`app/controller/api/user/UserFields.php` :: `info()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/user/fields/save` — save

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.UserFields/save`
- 源码：`app/controller/api/user/UserFields.php` :: `save()`
- 请求参数：
- `extend_info` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/user/form/create/:id` — create

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.form.FormRelated/create`
- 源码：`app/controller/api/store/form/FormRelated.php` :: `create()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /api/user/form/detail/:id` — detail

- 置信度：⛔ stale
- 说明：路由已登记，但 `app/controller/api/store/form/FormRelated.php` 中不存在方法 `detail`（原项目死路由/加密扩展/版本差异）。开发勿实现为有效接口，除非核实。
- 处理器：`api.store.form.FormRelated/detail`
- 源码：`app/controller/api/store/form/FormRelated.php` :: `detail()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一响应包 {status,message,data?}；详见控制器实现 | 外层: {status,message,data}

### `GET /api/user/form/lst` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.form.FormRelated/lst`
- 源码：`app/controller/api/store/form/FormRelated.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `GET /api/user/form/show/:id` — show

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.form.FormRelated/show`
- 源码：`app/controller/api/store/form/FormRelated.php` :: `show()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一响应包 {status,message,data?}；详见控制器实现 | 外层: {status,message,data}

### `GET /api/user/history` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.UserHistory/lst`
- 源码：`app/controller/api/user/UserHistory.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/user/history/batch/delete` — deleteHistoryBatch

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.UserHistory/deleteHistoryBatch`
- 源码：`app/controller/api/user/UserHistory.php` :: `deleteHistoryBatch()`
- 请求参数：
- `history_id` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/user/history/delete/:id` — deleteHistory

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.UserHistory/deleteHistory`
- 源码：`app/controller/api/user/UserHistory.php` :: `deleteHistory()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/user/integral/info` — integralInfo

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.User/integralInfo`
- 源码：`app/controller/api/user/User.php` :: `integralInfo()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/user/integral/lst` — integralList

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.User/integralList`
- 源码：`app/controller/api/user/User.php` :: `integralList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/user/member/info` — memberInfo

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.User/memberInfo`
- 源码：`app/controller/api/user/User.php` :: `memberInfo()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/user/member/log` — getMemberValue

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.Member/getMemberValue`
- 源码：`app/controller/api/user/Member.php` :: `getMemberValue()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/user/mp/binding` — mpPhone

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.User/mpPhone`
- 源码：`app/controller/api/user/User.php` :: `mpPhone()`
- 请求参数：
- `code` (query/body, 可选)
- `iv` (query/body, 可选)
- `encryptedData` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/user/receipt/create` — create

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.UserReceipt/create`
- 源码：`app/controller/api/user/UserReceipt.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/user/receipt/delete/:id` — delete

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.UserReceipt/delete`
- 源码：`app/controller/api/user/UserReceipt.php` :: `delete()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/user/receipt/detail/:id` — detail

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.UserReceipt/detail`
- 源码：`app/controller/api/user/UserReceipt.php` :: `detail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/user/receipt/is_default/:id` — isDefault

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.UserReceipt/isDefault`
- 源码：`app/controller/api/user/UserReceipt.php` :: `isDefault()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/user/receipt/lst` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.UserReceipt/lst`
- 源码：`app/controller/api/user/UserReceipt.php` :: `lst()`
- 请求参数：
- `receipt_title_type` (query/body, 可选)
- `receipt_type` (query/body, 可选)
- `is_default` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/user/receipt/order` — order

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.UserReceipt/order`
- 源码：`app/controller/api/user/UserReceipt.php` :: `order()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `status` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/user/receipt/order/:id` — orderDetail

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.UserReceipt/orderDetail`
- 源码：`app/controller/api/user/UserReceipt.php` :: `orderDetail()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/user/receipt/update/:id` — update

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.UserReceipt/update`
- 源码：`app/controller/api/user/UserReceipt.php` :: `update()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/user/recharge` — recharge

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.UserRecharge/recharge`
- 源码：`app/controller/api/user/UserRecharge.php` :: `recharge()`
- 请求参数：
- `type` (query/body, 可选)
- `price` (query/body, 可选)
- `recharge_id` (query/body, 可选)
- `return_url` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/user/recharge/brokerage` — brokerage

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.UserRecharge/brokerage`
- 源码：`app/controller/api/user/UserRecharge.php` :: `brokerage()`
- 请求参数：
- `brokerage` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/user/relation/batch/create` — batchCreate

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.UserRelation/batchCreate`
- 源码：`app/controller/api/user/UserRelation.php` :: `batchCreate()`
- 请求参数：
- `type_id` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/user/relation/batch/delete` — batchDelete

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.UserRelation/batchDelete`
- 源码：`app/controller/api/user/UserRelation.php` :: `batchDelete()`
- 请求参数：
- `type_id` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/user/relation/create` — create

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.UserRelation/create`
- 源码：`app/controller/api/user/UserRelation.php` :: `create()`
- 请求参数：
- `type_id` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/user/relation/delete` — delete

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.UserRelation/delete`
- 源码：`app/controller/api/user/UserRelation.php` :: `delete()`
- 请求参数：
- `type_id` (query/body, 可选)
- `type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/user/relation/merchant/lst` — merchantList

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.UserRelation/merchantList`
- 源码：`app/controller/api/user/UserRelation.php` :: `merchantList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/user/relation/product/lst` — productList

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.UserRelation/productList`
- 源码：`app/controller/api/user/UserRelation.php` :: `productList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/user/services` — services

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.User/services`
- 源码：`app/controller/api/user/User.php` :: `services()`
- 请求参数：
- `is_verify` (query/body, 可选)
- `customer` (query/body, 可选)
- `is_goods` (query/body, 可选)
- `is_open` (query/body, 可选) 默认 1
- `is_sys` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/user/sign/create` — create

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.UserSign/create`
- 源码：`app/controller/api/user/UserSign.php` :: `create()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/user/sign/info` — info

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.UserSign/info`
- 源码：`app/controller/api/user/UserSign.php` :: `info()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/user/sign/lst` — lst

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.UserSign/lst`
- 源码：`app/controller/api/user/UserSign.php` :: `lst()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/user/sign/month` — month

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.UserSign/month`
- 源码：`app/controller/api/user/UserSign.php` :: `month()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/user/spread` — spread

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Auth/spread`
- 源码：`app/controller/api/Auth.php` :: `spread()`
- 请求参数：
- `spread_spid` (query/body, 可选) 默认 0
- `spread_code` (query/body, 可选) 默认 null
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /api/user/spread_image` — spread_image

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.User/spread_image`
- 源码：`app/controller/api/user/User.php` :: `spread_image()`
- 请求参数：
- `type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/user/spread_info` — spread_info

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.User/spread_info`
- 源码：`app/controller/api/user/User.php` :: `spread_info()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/user/spread_level` — spread_info

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.User/spread_info`
- 源码：`app/controller/api/user/User.php` :: `spread_info()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/user/spread_list` — spread_list

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.User/spread_list`
- 源码：`app/controller/api/user/User.php` :: `spread_list()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `sort` (query/body, 可选)
- `keyword` (query/body, 可选)
- `start` (query/body, 可选)
- `stop` (query/body, 可选)
- `level` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/user/spread_order` — spread_order

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.User/spread_order`
- 源码：`app/controller/api/user/User.php` :: `spread_order()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- `start` (query/body, 可选)
- `stop` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/user/spread_top` — spread_top

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.User/spread_top`
- 源码：`app/controller/api/user/User.php` :: `spread_top()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `POST /api/user/switch` — switchUser

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.User/switchUser`
- 源码：`app/controller/api/user/User.php` :: `switchUser()`
- 请求参数：
- `uid` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。

### `GET /api/user/v2/spread_image` — spread_image_v2

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.user.User/spread_image_v2`
- 源码：`app/controller/api/user/User.php` :: `spread_image_v2()`
- 请求参数：
- `type` (query/body, 可选)
- 返回：统一 JSON：`{ status: number, message: string, data?: any }`；成功 status 通常为 200，失败 400。


## `api/v2`

### `POST /api/v2/order/check` — 订单核对

- 置信度：✅ high
- 说明：控制器方法已存在；参数来自方法体 param/params
- 处理器：`api.store.order.StoreOrder/v2CheckOrder`
- 源码：`app/controller/api/store/order/StoreOrder.php` :: `v2CheckOrder()`
- 请求参数：
- `cart_id` (query/body, 可选)
- `address_id` (query/body, 可选)
- `use_coupon` (query/body, 可选)
- `takes` (query/body, 可选)
- `city_takes` (query/body, 可选)
- `use_integral` (query/body, 可选)
- 返回：data 为下单核对/创建结果；外层 {status,message,data}

### `POST /api/v2/order/create` — 创建订单

- 置信度：✅ high
- 说明：控制器方法已存在；参数来自方法体 param/params
- 处理器：`api.store.order.StoreOrder/v2CreateOrder`
- 源码：`app/controller/api/store/order/StoreOrder.php` :: `v2CreateOrder()`
- 请求参数：
- `cart_id` (query/body, 可选)
- `address_id` (query/body, 可选)
- `use_coupon` (query/body, 可选)
- `takes` (query/body, 可选)
- `use_integral` (query/body, 可选)
- `receipt_data` (query/body, 可选)
- `extend` (query/body, 可选)
- `mark` (query/body, 可选)
- `pay_type` (query/body, 可选)
- `key` (query/body, 可选)
- `post` (query/body, 可选)
- `return_url` (query/body, 可选)
- 返回：data 为下单核对/创建结果；外层 {status,message,data}

### `GET /api/v2/system/city` — cityList

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.shipping.City/cityList`
- 源码：`app/controller/merchant/store/shipping/City.php` :: `cityList()`
- 请求参数：
- `address` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /api/v2/system/city/lst/:pid` — lstV2

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`merchant.store.shipping.City/lstV2`
- 源码：`app/controller/merchant/store/shipping/City.php` :: `lstV2()`
- 请求参数：
- `pid` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `api/v3`

### `POST /api/v3/order/check` — v3CheckOrder

- 置信度：⛔ stale
- 说明：路由指向 `v3CheckOrder()`，但当前源码 StoreOrder.php **不存在该方法**（仅有 v2CreateOrder/v2CheckOrder）。此 CRMEB 构建下请以 `/api/v2/order/*` 为准；勿按本条实现。
- 处理器：`api.store.order.StoreOrder/v3CheckOrder`
- 源码：`app/controller/api/store/order/StoreOrder.php` :: `v3CheckOrder()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：原路由可能不可用；见 doc_note

### `POST /api/v3/order/create` — v3CreateOrder

- 置信度：⛔ stale
- 说明：路由指向 `v3CreateOrder()`，但当前源码 StoreOrder.php **不存在该方法**（仅有 v2CreateOrder/v2CheckOrder）。此 CRMEB 构建下请以 `/api/v2/order/*` 为准；勿按本条实现。
- 处理器：`api.store.order.StoreOrder/v3CreateOrder`
- 源码：`app/controller/api/store/order/StoreOrder.php` :: `v3CreateOrder()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：原路由可能不可用；见 doc_note


## `api/verifier`

### `POST /api/verifier/:merId/:id` — verify

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreOrderVerify/verify`
- 源码：`app/controller/api/store/order/StoreOrderVerify.php` :: `verify()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- `verify_code` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

### `GET /api/verifier/:merId/order/:id` — detail

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.store.order.StoreOrderVerify/detail`
- 源码：`app/controller/api/store/order/StoreOrderVerify.php` :: `detail()`
- 请求参数：
- `merId` (path, 必填) 路径参数
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `api/version`

### `GET /api/version` — version

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`admin.Common/version`
- 源码：`app/controller/admin/Common.php` :: `version()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `api/wechat`

### `GET /api/wechat/config` — jsConfig

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`api.Wechat/jsConfig`
- 源码：`app/controller/api/Wechat.php` :: `jsConfig()`
- 请求参数：
- `url` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

