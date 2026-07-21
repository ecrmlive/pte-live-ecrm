# 路由补全：页面路径 → 接口级 CRUD

> 用于补齐「菜单无按钮权限节点」的页面：以 `route` 的 `_path` + `_alias` 为准。

## 平台后台

### `/app/routine/download`  （— ✓R — —，接口 1）

| 操作说明 | CRUD | HTTP | name |
| --- | --- | --- | --- |
| 小程序下载 | `R` | GET | `configRoutineDownload` |

### `/config/picture`  （✓C — — —，接口 1）

| 操作说明 | CRUD | HTTP | name |
| --- | --- | --- | --- |
| 上传图片 | `C` | POST | `uploadImage` |

### `/dashboard`  （— — — —，接口 2）

| 操作说明 | CRUD | HTTP | name |
| --- | --- | --- | --- |
| 成交用户 | `O` | GET | `systemStatisticsUser` |
| 成交用户占比 | `O` | GET | `systemStatisticsUserRate` |

### `/data-screen/index`  （— — — —，接口 1）

| 操作说明 | CRUD | HTTP | name |
| --- | --- | --- | --- |
| 数据大屏 | `O` | GET | `systemDataScreen` |

### `/delivery/recharge_record`  （✓C ✓R — —，接口 3）

| 操作说明 | CRUD | HTTP | name |
| --- | --- | --- | --- |
| 统计 | `R` | GET | `systemDeliveryOrderTitle` |
| 余额 | `R` | GET | `systemDeliveryStationGetBalance` |
| 充值记录 | `CR` | GET | `systemDeliveryStationPaayyLst` |

### `/delivery/station`  （— ✓R — —，接口 2）

| 操作说明 | CRUD | HTTP | name |
| --- | --- | --- | --- |
| 门店详情 | `R` | GET | `systemDeliveryStationDetail` |
| 门店列表 | `R` | GET | `systemDeliveryStationlst` |

### `/delivery/usage_record`  （— ✓R — —，接口 2）

| 操作说明 | CRUD | HTTP | name |
| --- | --- | --- | --- |
| 配送详情 | `R` | GET | `systemDeliveryOrderDetail` |
| 配送记录 | `R` | GET | `systemDeliveryOrderLst` |

### `/freight/express`  （— ✓R — —，接口 1）

| 操作说明 | CRUD | HTTP | name |
| --- | --- | --- | --- |
| 列表 | `R` | GET | `systemServeExportLst` |

### `/maintain/cache`  （— — — —，接口 2）

| 操作说明 | CRUD | HTTP | name |
| --- | --- | --- | --- |
| 替换素材域名 | `O` | POST | `systemAttachmentReplaceHost` |
| 清除缓存 | `O` | POST | `systemClearCache` |

### `/marketing/coupon/user`  （— ✓R — —，接口 1）

| 操作说明 | CRUD | HTTP | name |
| --- | --- | --- | --- |
| 使用记录 | `R` | GET | `systemCouponIssue` |

### `/marketing/platform_coupon/couponRecord`  （— ✓R — —，接口 1）

| 操作说明 | CRUD | HTTP | name |
| --- | --- | --- | --- |
| 使用记录 | `R` | GET | `systemCouponIssue` |

### `/marketing/platform_coupon/couponSend`  （✓C ✓R — —，接口 1）

| 操作说明 | CRUD | HTTP | name |
| --- | --- | --- | --- |
| 发送记录 | `CR` | GET | `systemCouponSendLst` |

### `/product/merSpecs`  （— ✓R — —，接口 1）

| 操作说明 | CRUD | HTTP | name |
| --- | --- | --- | --- |
| 商户参数模板 | `R` | GET | `systemStoreParameterTemplateMerLst` |

### `/service/balance_record`  （— ✓R — —，接口 1）

| 操作说明 | CRUD | HTTP | name |
| --- | --- | --- | --- |
| 商户结余 | `R` | GET | `systemServeMerLst` |

### `/service/purchase`  （✓C ✓R — —，接口 2）

| 操作说明 | CRUD | HTTP | name |
| --- | --- | --- | --- |
| 商户购买记录 | `CR` | GET | `systemServeMerPayLst` |
| 购买记录 | `CR` | GET | `systemServePayLst` |

### `/setting/sms/sms_config/index`  （— ✓R — —，接口 3）

| 操作说明 | CRUD | HTTP | name |
| --- | --- | --- | --- |
| 使用记录 | `R` | GET | `systemServeExportDumpLst` |
| 模板 | `O` | GET | `systemServeExportTemps` |
| 使用记录 | `R` | GET | `systemStoreProductCopyLst` |

### `/setting/systemLog`  （— ✓R — —，接口 1）

| 操作说明 | CRUD | HTTP | name |
| --- | --- | --- | --- |
| 操作日志 | `R` | GET | `systemAdminLog` |

### `/setting/theme_style`  （✓C ✓R — —，接口 2）

| 操作说明 | CRUD | HTTP | name |
| --- | --- | --- | --- |
| 一键换色 | `R` | GET | `systemGetChangeColor` |
| 一键换色保存 | `C` | POST | `systemSetChangeColor` |

### `/sms/applyList`  （— ✓R — —，接口 1）

| 操作说明 | CRUD | HTTP | name |
| --- | --- | --- | --- |
| 申请记录 | `R` | GET | `systemServeSmsApplyRecord` |

### `/sms/template`  （— — — —，接口 2）

| 操作说明 | CRUD | HTTP | name |
| --- | --- | --- | --- |
| 申请模板 | `O` | POST | `systemServeSmsApply` |
| 短信模板 | `O` | GET | `systemServeSmsTemps` |

### `/systemForm/Basics/upload`  （✓C — ✓U —，接口 1）

| 操作说明 | CRUD | HTTP | name |
| --- | --- | --- | --- |
| 上传配置保存 | `CU` | POST | `systemSaveUploadConfig` |

### `/systemForm/delivery`  （✓C — ✓U —，接口 1）

| 操作说明 | CRUD | HTTP | name |
| --- | --- | --- | --- |
| 编辑配置 | `CU` | POST | `systemDeliveryConfigSave` |

### `/user/list`  （✓C — — —，接口 1）

| 操作说明 | CRUD | HTTP | name |
| --- | --- | --- | --- |
| 发送优惠券 | `C` | POST | `systemCouponSend` |

### `/user/member/record`  （✓C ✓R — —，接口 2）

| 操作说明 | CRUD | HTTP | name |
| --- | --- | --- | --- |
| 统计 | `R` | GET | `systemUserSvipCountInfo` |
| 列表 | `CR` | GET | `systemUserSvipPayLst` |

## 商户后台

### `/delivery/recharge_record`  （✓C ✓R — —，接口 2）

| 操作说明 | CRUD | HTTP | name |
| --- | --- | --- | --- |
| 充值二维码 | `CR` | GET | `merchantStoreDeliveryGetQrcode` |
| 充值记录 | `CR` | GET | `merchantStoreDeliveryPayLst` |

### `/devise/diy/list`  （— ✓R — —，接口 1）

| 操作说明 | CRUD | HTTP | name |
| --- | --- | --- | --- |
| 列表 | `R` | GET | `merchantDiyPageLinkLst` |

### `/marketing/coupon/send`  （✓C ✓R — —，接口 1）

| 操作说明 | CRUD | HTTP | name |
| --- | --- | --- | --- |
| 发送优惠券记录 | `CR` | GET | `merchantCouponSendLst` |

### `/marketing/coupon/user`  （— ✓R — —，接口 1）

| 操作说明 | CRUD | HTTP | name |
| --- | --- | --- | --- |
| 使用记录 | `R` | GET | `merchantCouponIssue` |

### `/setting/sms/sms_config/index`  （— ✓R — —，接口 1）

| 操作说明 | CRUD | HTTP | name |
| --- | --- | --- | --- |
| 使用记录 | `R` | GET | `merchantServeExportDumpLst` |

### `/setting/systemLog`  （— ✓R — —，接口 1）

| 操作说明 | CRUD | HTTP | name |
| --- | --- | --- | --- |
| 操作日志 | `R` | GET | `merchantAdminLog` |

### `/setting/theme_style`  （— ✓R — —，接口 1）

| 操作说明 | CRUD | HTTP | name |
| --- | --- | --- | --- |
| 一键换色 | `R` | GET | `merchantGetChangeColor` |

### `/user/list`  （— — — —，接口 1）

| 操作说明 | CRUD | HTTP | name |
| --- | --- | --- | --- |
| 优惠券可用商品 | `O` | GET | `merchantCouponProduct` |

### `/user/searchRecord`  （— ✓R — —，接口 1）

| 操作说明 | CRUD | HTTP | name |
| --- | --- | --- | --- |
| 搜索记录 | `R` | GET | `merchantUserSearchLog` |

### `config`  （✓C ✓R ✓U —，接口 2）

| 操作说明 | CRUD | HTTP | name |
| --- | --- | --- | --- |
| 配置获取 | `RU` | GET | `merchantConfigForm` |
| 配置保存 | `CU` | POST | `merchantConfigSave` |
