# 商户后台 — 功能点清单（页面 → 按钮/操作，含增删改查）

> 主来源：`eb_system_menu` 按钮权限。
> 若页面无按钮权限节点，则用同页 `route` 的 `_path` 接口补全（标注来源=路由）。
> CRUD：`C`创建 · `R`读取 · `U`更新 · `D`删除 · `O`其他。
>
> 本页的页面覆盖与操作统计为生成时的展示快照。全量逐项验收以 [`../generated/features-master.tsv`](../generated/features-master.tsv) 为唯一口径，其中商户后台为 **615** 项；总验收入口见 [`../CRMEB-FULL-FUNCTION-CHECKLIST.md`](../CRMEB-FULL-FUNCTION-CHECKLIST.md)。

## 统计

| 项 | 数量 |
| --- | ---: |
| 页面 | 102 |
| 操作（按钮/接口） | 600 |
| 其中由路由补全 | 100 |
| C | 121 |
| R | 256 |
| U | 158 |
| D | 57 |

## 页面 CRUD 覆盖一览

| 模块路径 | 页面路由 | C | R | U | D | 操作数 |
| --- | --- | --- | --- | --- | --- | ---: |
| 公告 | `/station/notice` | — | ✓ | — | ✓ | 5 |
| 员工 | `/server` | — | ✓ | — | — | 0 |
| 员工 / 店员管理 | `/server_manage` | — | ✓ | — | — | 0 |
| 员工 / 店员管理 / 店员列表 | `/config/service` | ✓ | ✓ | ✓ | ✓ | 10 |
| 员工 / 店员管理 / 店员配置 | `/systemForm/Basics/mer_service_switch` | — | ✓ | — | — | 0 |
| 员工 / 店员管理 / 自动回复 | `/systemForm/customer_keyword` | ✓ | ✓ | ✓ | ✓ | 5 |
| 员工 / 服务人员 | `/config/service_staff` | ✓ | ✓ | ✓ | ✓ | 7 |
| 员工 / 服务人员 / 服务人员 | `/config/service_staff` | — | ✓ | — | — | 2 |
| 员工 / 服务人员 / 服务统计 | `/config/service_statistic` | ✓ | ✓ | ✓ | — | 5 |
| 员工 / 配送人员 | `/delivery/personnel_manage/index` | ✓ | ✓ | ✓ | ✓ | 7 |
| 员工 / 配送人员 / 配送员管理 | `/delivery/personnel_manage` | ✓ | ✓ | ✓ | ✓ | 7 |
| 员工 / 配送人员 / 配送统计 | `/delivery/delivice_statistic` | — | ✓ | — | — | 0 |
| 商品 | `/product` | ✓ | ✓ | — | — | 2 |
| 商品 / 卡密列表 | `/product/cdkey` | ✓ | ✓ | ✓ | ✓ | 13 |
| 商品 / 商品分类 | `/product/classify` | ✓ | ✓ | ✓ | ✓ | 8 |
| 商品 / 商品列表 | `/product/list` | ✓ | ✓ | ✓ | ✓ | 41 |
| 商品 / 商品单位 | `/product/unit` | ✓ | ✓ | ✓ | ✓ | 4 |
| 商品 / 商品参数 | `/product/specs` | ✓ | ✓ | ✓ | ✓ | 6 |
| 商品 / 商品标签 | `/product/label` | ✓ | ✓ | ✓ | ✓ | 6 |
| 商品 / 商品规格 | `/product/attr` | ✓ | ✓ | ✓ | ✓ | 4 |
| 商品 / 服务模板 | `/config/guarantee` | ✓ | ✓ | ✓ | ✓ | 7 |
| 用户 | `/user` | ✓ | ✓ | ✓ | — | 6 |
| 用户 / 搜索记录 | `/user/searchRecord` | — | ✓ | — | — | 1 |
| 用户 / 标签管理 | `/user/_label` | — | ✓ | — | — | 0 |
| 用户 / 标签管理 / 手动标签 | `/user/label` | ✓ | ✓ | ✓ | ✓ | 4 |
| 用户 / 标签管理 / 自动标签 | `/user/maticlabel` | ✓ | ✓ | ✓ | ✓ | 5 |
| 用户 / 用户管理 | `/user/list` | ✓ | ✓ | ✓ | — | 6 |
| 营销 | `/marketing` | ✓ | ✓ | — | — | 4 |
| 营销 / 专场列表 | `/group/topic/95` | ✓ | ✓ | ✓ | ✓ | 8 |
| 营销 / 优惠券 | `/marketing/coupon` | ✓ | ✓ | — | — | 4 |
| 营销 / 优惠券 / 优惠券列表 | `/marketing/coupon/list` | ✓ | ✓ | ✓ | ✓ | 8 |
| 营销 / 优惠券 / 发送记录 | `/marketing/coupon/send` | ✓ | ✓ | — | — | 1 |
| 营销 / 优惠券 / 领取记录 | `/marketing/coupon/user` | — | ✓ | — | — | 1 |
| 营销 / 优惠套餐 | `/marketing/discounts/list` | ✓ | ✓ | ✓ | ✓ | 8 |
| 营销 / 助力 | `/assist` | — | ✓ | — | — | 0 |
| 营销 / 助力 / 助力商品 | `/marketing/assist/list` | ✓ | ✓ | ✓ | ✓ | 11 |
| 营销 / 助力 / 助力活动 | `/marketing/assist/assist_set` | — | ✓ | — | — | 2 |
| 营销 / 拼团 | `/marketing/combination` | ✓ | — | — | — | 2 |
| 营销 / 拼团 / 拼团商品列表 | `/marketing/combination/combination_goods` | ✓ | ✓ | ✓ | ✓ | 12 |
| 营销 / 拼团 / 拼团活动列表 | `/marketing/combination/combination_list` | ✓ | ✓ | — | — | 4 |
| 营销 / 直播 | `/` | ✓ | ✓ | — | — | 5 |
| 营销 / 直播 / 直播助手 | `/marketing/studio/assistant` | ✓ | ✓ | ✓ | ✓ | 7 |
| 营销 / 直播 / 直播商品管理 | `/marketing/broadcast/list` | — | ✓ | — | — | 0 |
| 营销 / 直播 / 直播间管理 | `/marketing/studio/list` | ✓ | ✓ | ✓ | ✓ | 26 |
| 营销 / 秒杀 | `/marketing/seckill/list` | — | ✓ | — | — | 0 |
| 营销 / 秒杀 / 秒杀商品 | `/marketing/seckill/product/list` | ✓ | ✓ | ✓ | ✓ | 15 |
| 营销 / 秒杀 / 秒杀活动 | `/marketing/seckill/store_seckill/list` | ✓ | ✓ | — | — | 9 |
| 营销 / 积分 | `/marketing/integral` | — | ✓ | — | — | 2 |
| 营销 / 积分 / 积分日志 | `/marketing/integral/log` | ✓ | ✓ | ✓ | — | 4 |
| 营销 / 积分 / 积分配置 | `/marketing/integral/config` | — | ✓ | — | — | 0 |
| 营销 / 逛逛社区 | `/community/list` | ✓ | ✓ | ✓ | ✓ | 7 |
| 营销 / 预售 | `/marketing/presell/list` | ✓ | ✓ | ✓ | ✓ | 11 |
| 装修 | `/devise/` | ✓ | ✓ | ✓ | ✓ | 12 |
| 装修 / 商品分类 | `/devise/diy/product_category` | ✓ | ✓ | — | — | 2 |
| 装修 / 系统表单 | `/systemForm/form_list` | ✓ | ✓ | ✓ | ✓ | 5 |
| 装修 / 素材管理 | `/config/picture` | ✓ | ✓ | ✓ | ✓ | 12 |
| 装修 / 装修 | `/devise/diy/list` | ✓ | ✓ | ✓ | ✓ | 13 |
| 订单 | `/order` | ✓ | — | ✓ | — | 2 |
| 订单 / 代客下单 | `/order/customer` | ✓ | ✓ | ✓ | ✓ | 22 |
| 订单 / 商品评价 | `/product/reviews` | ✓ | ✓ | ✓ | — | 6 |
| 订单 / 核销记录 | `/order/cancellation` | — | ✓ | — | — | 2 |
| 订单 / 订单管理 | `/order/list` | ✓ | ✓ | ✓ | ✓ | 41 |
| 订单 / 退款订单 | `/order/refund` | ✓ | ✓ | ✓ | ✓ | 12 |
| 订单 / 预约服务 | `/order/reservation` | — | — | — | — | 1 |
| 订单 / 预约设置 | `/product/reservation` | — | ✓ | — | — | 0 |
| 设置 | `/config` | ✓ | ✓ | ✓ | ✓ | 7 |
| 设置 / 一号通 | `/one_setting` | — | ✓ | — | — | 0 |
| 设置 / 一号通 / 平台一号通 | `/setting/sms/sms_config/index` | — | ✓ | — | — | 5 |
| 设置 / 一号通 / 自有一号通 | `/setting/sms/sms_account/index` | — | ✓ | — | — | 0 |
| 设置 / 一号通 / 配置管理 | `/setting/sms/dumpConfig` | ✓ | — | ✓ | — | 1 |
| 设置 / 付费会员 | `/systemForm/Basics/svip` | — | ✓ | — | — | 0 |
| 设置 / 同城配送 | `/delivery` | — | ✓ | — | — | 3 |
| 设置 / 同城配送 / 充值记录 | `/delivery/recharge_record` | ✓ | ✓ | — | — | 1 |
| 设置 / 同城配送 / 发货点管理 | `/delivery/store_manage` | ✓ | ✓ | ✓ | ✓ | 8 |
| 设置 / 同城配送 / 配送记录 | `/delivery/usage_record` | — | ✓ | — | ✓ | 3 |
| 设置 / 同城配送 / 配送设置 | `/setting/delivery` | — | — | ✓ | — | 2 |
| 设置 / 同城配送 / 配送门店 | `/delivery/delivery_point` | ✓ | ✓ | ✓ | ✓ | 9 |
| 设置 / 店铺信息 | `/systemForm/modifyStoreInfo` | ✓ | — | ✓ | — | 2 |
| 设置 / 店铺配置 | `/systemForm/Basics/mer_base` | — | — | ✓ | — | 1 |
| 设置 / 开放账户 | `/systemForm/openAuth/list` | ✓ | ✓ | ✓ | ✓ | 7 |
| 设置 / 快递配送 | `/city` | — | ✓ | — | — | 0 |
| 设置 / 快递配送 / 物流公司 | `/config/freight/express` | — | ✓ | ✓ | — | 3 |
| 设置 / 快递配送 / 运费模板 | `/config/freight/shippingTemplates` | ✓ | ✓ | ✓ | ✓ | 5 |
| 设置 / 打印配置 | `/setting/printer` | — | ✓ | — | — | 8 |
| 设置 / 打印配置 / 小票打印 | `/setting/printer/list` | ✓ | ✓ | ✓ | ✓ | 7 |
| 设置 / 打印配置 / 打印配置 | `/systemForm/Basics/printer_tabs` | — | ✓ | — | — | 0 |
| 设置 / 权限管理 | `/setting` | — | ✓ | — | — | 1 |
| 设置 / 权限管理 / 操作日志 | `/setting/systemLog` | — | ✓ | ✓ | — | 4 |
| 设置 / 权限管理 / 管理员管理 | `/setting/systemAdmin` | ✓ | ✓ | ✓ | ✓ | 6 |
| 设置 / 权限管理 / 身份管理 | `/setting/systemRole` | ✓ | ✓ | ✓ | ✓ | 10 |
| 财务 | `/accounts` | ✓ | ✓ | — | ✓ | 8 |
| 财务 / 分账管理 | `/systemForm/applyList` | — | ✓ | — | — | 4 |
| 财务 / 发票管理 | `/order/invoice` | ✓ | ✓ | ✓ | — | 6 |
| 财务 / 收款方式 | `/accounts/payType` | ✓ | — | — | — | 1 |
| 财务 / 申请分账商户 | `/systemForm/applyments` | ✓ | ✓ | ✓ | — | 5 |
| 财务 / 账单管理 | `/accounts/statement` | — | ✓ | — | — | 6 |
| 财务 / 资金流水 | `/accounts/capitalFlow` | — | ✓ | — | — | 5 |
| 财务 / 转账记录 | `/accounts/transManagement` | ✓ | ✓ | ✓ | ✓ | 9 |
| 首页 | `/dashboard` | — | ✓ | ✓ | — | 10 |
| 首页 / 商品统计 | `/statistic/product` | — | ✓ | — | — | 3 |
| 首页 / 控制台 | `/dashboard` | ✓ | ✓ | — | ✓ | 5 |
| 首页 / 订单统计 | `/statistic/order` | — | ✓ | — | — | 3 |

## 分页面操作明细

### 公告

#### 公告

- 页面路由：`/station/notice`
- CRUD：C=— R=✓ U=— D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemNoticeLogList` | 菜单权限 |
| 已读 | `R` | `systemNoticeLogRead` | 菜单权限 |
| 删除 | `D` | `systemNoticeLogDel` | 菜单权限 |
| 未读统计 | `R` | `systemNoticeLogUnreadCount` | 菜单权限 |
| 详情 | `R` | `systemNoticeLogDetail` | 菜单权限 |

### 员工

#### 员工

- 页面路由：`/server`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 员工 / 店员管理

- 页面路由：`/server_manage`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 员工 / 店员管理 / 店员列表

- 页面路由：`/config/service`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `merchantServiceLst` | 菜单权限 |
| 登录 | `O` | `merchantServiceLogin` | 菜单权限 |
| 添加 | `C` | `merchantServiceCreate` | 菜单权限 |
| 编辑 | `U` | `merchantServiceUpdate` | 菜单权限 |
| 修改状态 | `U` | `merchantServiceSwitchStatus` | 菜单权限 |
| 删除 | `D` | `merchantServiceDelete` | 菜单权限 |
| 客服的全部用户  | `O` | `merchantServiceServiceUserList` | 菜单权限 |
| 用户与客服聊天记录 | `R` | `merchantServiceServiceUserLogLst` | 菜单权限 |
| 客服的聊天用户列表 | `R` | `merchantServiceServiceMerchantUserList` | 菜单权限 |
| 用户与商户聊天记录 | `R` | `merchantServiceMerchantUserLogLst` | 菜单权限 |

#### 员工 / 店员管理 / 店员配置

- 页面路由：`/systemForm/Basics/mer_service_switch`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 员工 / 店员管理 / 自动回复

- 页面路由：`/systemForm/customer_keyword`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `merchantServiceReplyLst` | 菜单权限 |
| 添加 | `C` | `merchantServiceReplyCreate` | 菜单权限 |
| 编辑 | `U` | `merchantServiceReplyUpdate` | 菜单权限 |
| 切换状态 | `U` | `merchantServiceReplyStatus` | 菜单权限 |
| 删除 | `D` | `merchantServiceReplyDelete` | 菜单权限 |

#### 员工 / 服务人员

- 页面路由：`/config/service_staff`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `merchantStaffsLst` | 菜单权限 |
| 添加 | `C` | `merchantStaffsCreate` | 菜单权限 |
| 编辑 | `U` | `merchantStaffsUpdate` | 菜单权限 |
| 修改状态 | `U` | `merchantStaffsSwitchStatus` | 菜单权限 |
| 删除 | `D` | `merchantStaffsDelete` | 菜单权限 |
| 列表 | `R` | `merchantStaffsStatisticsList` | 菜单权限 |
| 详情 | `R` | `merchantStaffsStatisticsDetail` | 菜单权限 |

#### 员工 / 服务人员 / 服务人员

- 页面路由：`/config/service_staff`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `merchantStaffsStatisticsList` | 路由 |
| 详情 | `R` | `merchantStaffsStatisticsDetail` | 路由 |

#### 员工 / 服务人员 / 服务统计

- 页面路由：`/config/service_statistic`
- CRUD：C=✓ R=✓ U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 客服的聊天用户列表 | `R` | `merchantServiceServiceMerchantUserList` | 路由 |
| 用户与商户聊天记录 | `R` | `merchantServiceMerchantUserLogLst` | 路由 |
| 列表 | `R` | `merchantServiceReplyLst` | 路由 |
| 添加 | `C` | `merchantServiceReplyCreate` | 路由 |
| 编辑 | `U` | `merchantServiceReplyUpdate` | 路由 |

#### 员工 / 配送人员

- 页面路由：`/delivery/personnel_manage/index`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `merchantDeliveryServiceLst` | 菜单权限 |
| 修改状态 | `U` | `merchantDeliveryServiceStatus` | 菜单权限 |
| 添加 | `C` | `merchantDeliveryServiceCreate` | 菜单权限 |
| 编辑 | `U` | `merchantDeliveryServiceUpdate` | 菜单权限 |
| 删除 | `D` | `merchantDeliveryServiceDelete` | 菜单权限 |
| 统计列表 | `R` | `merchantDeliveryServiceStatisticsList` | 菜单权限 |
| 统计详情 | `R` | `merchantDeliveryServiceStatisticsDetail` | 菜单权限 |

#### 员工 / 配送人员 / 配送员管理

- 页面路由：`/delivery/personnel_manage`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `merchantDeliveryServiceLst` | 菜单权限 |
| 修改状态 | `U` | `merchantDeliveryServiceStatus` | 菜单权限 |
| 添加 | `C` | `merchantDeliveryServiceCreate` | 菜单权限 |
| 编辑 | `U` | `merchantDeliveryServiceUpdate` | 菜单权限 |
| 删除 | `D` | `merchantDeliveryServiceDelete` | 菜单权限 |
| 统计列表 | `R` | `merchantDeliveryServiceStatisticsList` | 菜单权限 |
| 统计详情 | `R` | `merchantDeliveryServiceStatisticsDetail` | 菜单权限 |

#### 员工 / 配送人员 / 配送统计

- 页面路由：`/delivery/delivice_statistic`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

### 商品

#### 商品

- 页面路由：`/product`
- CRUD：C=✓ R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `merchantStoreAttrTemplateLst` | 路由 |
| 添加  | `C` | `merchantStoreAttrTemplateCreate` | 路由 |

#### 商品 / 卡密列表

- 页面路由：`/product/cdkey`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `merchantStoreProductCdkeyLibraryLst` | 菜单权限 |
| 列表 | `R` | `merchantStoreProductCdkeyLibraryDetail` | 菜单权限 |
| 添加 | `C` | `merchantStoreProductCdkeyLibraryCreate` | 菜单权限 |
| 编辑表单 | `U` | `merchantStoreProductCdkeyLibraryUpdateForm` | 菜单权限 |
| 编辑 | `U` | `merchantStoreProductCdkeyLibraryUpdate` | 菜单权限 |
| 删除 | `D` | `merchantStoreProductCdkeyLibraryDelete` | 菜单权限 |
| 导出 | `R` | `merchantStoreProductCdkeyLibraryExcel` | 菜单权限 |
| 批量导入 | `C` | `merchantStoreProductCdkeyLibraryImport` | 菜单权限 |
| 卡密列表 | `R` | `merchantStoreProductCdkeyLst` | 菜单权限 |
| 添加卡密 | `C` | `merchantStoreProductCdkeyCreate` | 菜单权限 |
| 编辑卡密 | `U` | `merchantStoreProductCdkeyUpdate` | 菜单权限 |
| 删除卡密 | `D` | `merchantStoreProductCdkeyDelete` | 菜单权限 |
| 批量删除 | `D` | `merchantStoreProductCdkeyLibraryBatchDelete` | 菜单权限 |

#### 商品 / 商品分类

- 页面路由：`/product/classify`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 编辑 | `U` | `merchantStoreCategoryUpdate` | 菜单权限 |
| 列表 | `R` | `merchantStoreCategoryLst` | 菜单权限 |
| 详情 | `R` | `merchantStoreCategoryDtailt` | 菜单权限 |
| 添加 | `C` | `merchantStoreCategoryCreate` | 菜单权限 |
| 删除 | `D` | `merchantStoreCategoryDelete` | 菜单权限 |
| 修改状态 | `U` | `merchantStoreCategorySwitchStatus` | 菜单权限 |
| 上传图片 | `C` | `merchantUploadImage` | 菜单权限 |
| 图片列表 | `R` | `merchantAttachmentLst` | 菜单权限 |

#### 商品 / 商品列表

- 页面路由：`/product/list`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 头部统计 | `R` | `merchantStoreProductLstFilter` | 菜单权限 |
| 列表 | `R` | `merchantStoreProductLst` | 菜单权限 |
| 添加 | `C` | `merchantStoreProductCreate` | 菜单权限 |
| 详情 | `R` | `merchantStoreProductDetail` | 菜单权限 |
| 上传视频配置 | `CU` | `merchantStoreProductTempKey` | 菜单权限 |
| 编辑 | `U` | `merchantStoreProductUpdate` | 菜单权限 |
| 删除 | `D` | `merchantStoreProductDelete` | 菜单权限 |
| 加入回收站 | `O` | `merchantStoreProductDestory` | 菜单权限 |
| 恢复 | `O` | `merchantStoreProductRestore` | 菜单权限 |
| 上下架 | `U` | `merchantStoreProductSwitchStatus` | 菜单权限 |
| 排序 | `U` | `merchantStoreProductUpdateSort` | 菜单权限 |
| 预览 | `R` | `merchantStoreProductPreview` | 菜单权限 |
| 标签 | `O` | `merchantStoreProductLabels` | 菜单权限 |
| 获取规格 | `R` | `merchantStoreProductAttrValue` | 菜单权限 |
| 列表 | `CR` | `merchantStoreProductCopyLst` | 菜单权限 |
| 获取信息 | `CR` | `merchantStoreProductCopyGet` | 菜单权限 |
| 统计 | `CR` | `merchantStoreProductCopyCount` | 菜单权限 |
| 保存 | `C` | `merchantStoreProductCopySave` | 菜单权限 |
| 上传图片 | `C` | `merchantUploadImage` | 菜单权限 |
| 图片列表 | `R` | `merchantAttachmentLst` | 菜单权限 |
| 免审编辑 | `U` | `merchantStoreProductFreeTrial` | 菜单权限 |
| 批量上下架 | `U` | `merchantStoreProductSwitchBatchStatus` | 菜单权限 |
| 批量设置运费模板 | `U` | `merchantStoreProductSwitchBatchTemplate` | 菜单权限 |
| 批量设置标签 | `U` | `merchantStoreProductSwitchBatchLabels` | 菜单权限 |
| 批量设置推荐 | `U` | `merchantStoreProductSwitchBatchHot` | 菜单权限 |
| 批量设置推荐 | `U` | `merchantStoreProductSwitchBatchExtension` | 菜单权限 |
| 批量设置会员价 | `U` | `merchantStoreProductSwitchBatchSvipType` | 菜单权限 |
| 获取规格 | `R` | `merchantStoreProductFormatAttr` | 菜单权限 |
| 获取批量修改列表 | `RU` | `merchantStoreProductGetBatchList` | 菜单权限 |
| 批量修改商品属性 | `U` | `merchantStoreProductSwitchBatchProcess` | 菜单权限 |
| 操作记录 | `R` | `merchantStoreProductGetOperateList` | 菜单权限 |
| 操作记录 | `R` | `merchantStoreProductUnbind` | 菜单权限 |
| 编辑商品获取信息 | `RU` | `merchantStoreProductGetEdit` | 菜单权限 |
| 添加预约商品 | `C` | `merchantStoreReservationProductCreate` | 菜单权限 |
| 获取预约商品 | `RU` | `merchantStoreReservationProductEditInfo` | 菜单权限 |
| 编辑预约商品 | `U` | `merchantStoreReservationProductEdit` | 菜单权限 |
| 预约商品详情 | `R` | `merchantStoreReservationProductDetail` | 菜单权限 |
| 批量修改预约商品库存 | `U` | `merchantStoreReservationProductEditStock` | 菜单权限 |
| 批量加入回收站 | `D` | `merchantStoreProductBatchDelete` | 菜单权限 |
| 批量恢复 | `O` | `merchantStoreProductBatchRestore` | 菜单权限 |
| 批量设置服务保障 | `U` | `merchantStoreProductBatchGuarantee` | 菜单权限 |

#### 商品 / 商品单位

- 页面路由：`/product/unit`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 商品单位列表 | `R` | `merchantStoreProductUnitLst` | 菜单权限 |
| 商品单位添加 | `C` | `merchantStoreProductUnitCreate` | 菜单权限 |
| 商品单位编辑 | `U` | `merchantStoreProductUnitUpdate` | 菜单权限 |
| 商品单位删除 | `D` | `merchantStoreProductUnitDelete` | 菜单权限 |

#### 商品 / 商品参数

- 页面路由：`/product/specs`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `merchantStoreParameterTemplateLst` | 菜单权限 |
| 详情 | `R` | `merchantStoreParameterTemplateDetail` | 菜单权限 |
| 删除 | `D` | `merchantStoreParameterTemplateDelete` | 菜单权限 |
| 添加 | `C` | `merchantStoreParameterTemplateCreate` | 菜单权限 |
| 编辑 | `U` | `merchantStoreParameterTemplateUpdate` | 菜单权限 |
| 删除属性 | `D` | `merchantStoreParameterTemplateDeleteValue` | 菜单权限 |

#### 商品 / 商品标签

- 页面路由：`/product/label`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `merchantStoreProductLabelLst` | 菜单权限 |
| 添加 | `C` | `merchantStoreProductLabelCreate` | 菜单权限 |
| 编辑 | `U` | `merchantStoreProductLabelUpdate` | 菜单权限 |
| 详情 | `R` | `merchantStoreProductLabelDetail` | 菜单权限 |
| 删除 | `D` | `merchantStoreProductLabelDelete` | 菜单权限 |
| 修改状态 | `U` | `merchantStoreProductLabelStatus` | 菜单权限 |

#### 商品 / 商品规格

- 页面路由：`/product/attr`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `merchantStoreAttrTemplateLst` | 菜单权限 |
| 添加  | `C` | `merchantStoreAttrTemplateCreate` | 菜单权限 |
| 删除 | `D` | `merchantStoreAttrTemplateDelete` | 菜单权限 |
| 文件类型 | `U` | `merchantStoreAttrTemplateUpdate` | 菜单权限 |

#### 商品 / 服务模板

- 页面路由：`/config/guarantee`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `merchantGuaranteeLst` | 菜单权限 |
| 添加 | `C` | `smerchantGuaranteeCreate` | 菜单权限 |
| 编辑 | `U` | `merchantGuaranteeUpdate` | 菜单权限 |
| 详情 | `R` | `merchantGuaranteeDetail` | 菜单权限 |
| 删除 | `D` | `merchantGuaranteeDelete` | 菜单权限 |
| 排序 | `U` | `merchantGuaranteeSort` | 菜单权限 |
| 修改状态 | `U` | `merchantGuaranteeStatus` | 菜单权限 |

### 用户

#### 用户

- 页面路由：`/user`
- CRUD：C=✓ R=✓ U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 优惠券可用商品 | `R` | `merchantCouponProduct` | 路由 |
| 订单列表 | `R` | `merchantUserOrder` | 路由 |
| 优惠券 | `R` | `merchantUserCoupon` | 路由 |
| 列表 | `R` | `merchantLabelRuleLst` | 路由 |
| 添加 | `C` | `merchantLabelRuleCreate` | 路由 |
| 编辑 | `U` | `merchantLabelRuleUpdate` | 路由 |

#### 用户 / 搜索记录

- 页面路由：`/user/searchRecord`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 搜索记录 | `R` | `merchantUserSearchLog` | 菜单权限 |

#### 用户 / 标签管理

- 页面路由：`/user/_label`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 用户 / 标签管理 / 手动标签

- 页面路由：`/user/label`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `merchantUserLabelLst` | 菜单权限 |
| 添加 | `C` | `merchantUserLabelCreate` | 菜单权限 |
| 删除 | `D` | `merchantUserLabelDelete` | 菜单权限 |
| 编辑 | `U` | `merchantUserLabelUpdate` | 菜单权限 |

#### 用户 / 标签管理 / 自动标签

- 页面路由：`/user/maticlabel`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `merchantLabelRuleLst` | 菜单权限 |
| 添加 | `C` | `merchantLabelRuleCreate` | 菜单权限 |
| 编辑 | `U` | `merchantLabelRuleUpdate` | 菜单权限 |
| 删除 | `D` | `merchantLabelRuleDelete` | 菜单权限 |
| 自动同步 | `U` | `merchantLabelRuleSync` | 菜单权限 |

#### 用户 / 用户管理

- 页面路由：`/user/list`
- CRUD：C=✓ R=✓ U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 优惠券可用商品 | `O` | `merchantCouponProduct` | 菜单权限 |
| 发送优惠券 | `C` | `merchantCouponSendCoupon` | 菜单权限 |
| 列表 | `R` | `merchantUserLst` | 菜单权限 |
| 修改标签 | `U` | `merchantUserChangeLabel` | 菜单权限 |
| 订单列表 | `R` | `merchantUserOrder` | 菜单权限 |
| 优惠券 | `O` | `merchantUserCoupon` | 菜单权限 |

### 营销

#### 营销

- 页面路由：`/marketing`
- CRUD：C=✓ R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 发送优惠券记录 | `CR` | `merchantCouponSendLst` | 路由 |
| 添加表单 | `C` | `merchantCouponCreateForm` | 路由 |
| 复制表单 | `C` | `merchantCouponIssueCloneForm` | 路由 |
| 添加 | `C` | `merchantCouponCreate` | 路由 |

#### 营销 / 专场列表

- 页面路由：`/group/topic/95`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 数据详情 | `R` | `merchantGroupDetail` | 菜单权限 |
| 数据列表 | `R` | `merchantGroupDataLst` | 菜单权限 |
| 数据添加 | `C` | `merchantGroupDataCreate` | 菜单权限 |
| 数据编辑 | `U` | `merchantGroupDataUpdate` | 菜单权限 |
| 数据删除 | `D` | `merchantGroupDataDelete` | 菜单权限 |
| 数据修改状态 | `U` | `merchantGroupDataChangeStatus` | 菜单权限 |
| 上传图片 | `C` | `uploadImage` | 菜单权限 |
| 图片列表 | `R` | `systemAttachmentLst` | 菜单权限 |

#### 营销 / 优惠券

- 页面路由：`/marketing/coupon`
- CRUD：C=✓ R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 发送优惠券记录 | `CR` | `merchantCouponSendLst` | 路由 |
| 添加表单 | `C` | `merchantCouponCreateForm` | 路由 |
| 复制表单 | `C` | `merchantCouponIssueCloneForm` | 路由 |
| 添加 | `C` | `merchantCouponCreate` | 路由 |

#### 营销 / 优惠券 / 优惠券列表

- 页面路由：`/marketing/coupon/list`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 添加 | `C` | `merchantCouponCreate` | 菜单权限 |
| 修改状态 | `U` | `merchantCouponIssueChangeStatus` | 菜单权限 |
| 列表 | `R` | `merchantCouponLst` | 菜单权限 |
| 删除 | `D` | `merchantCouponDelete` | 菜单权限 |
| 详情 | `R` | `merchantCouponDetail` | 菜单权限 |
| 编辑 | `U` | `systemCouponUpdate` | 菜单权限 |
| 上传图片 | `C` | `merchantUploadImage` | 菜单权限 |
| 图片列表 | `R` | `merchantAttachmentLst` | 菜单权限 |

#### 营销 / 优惠券 / 发送记录

- 页面路由：`/marketing/coupon/send`
- CRUD：C=✓ R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 发送优惠券记录 | `CR` | `merchantCouponSendLst` | 菜单权限 |

#### 营销 / 优惠券 / 领取记录

- 页面路由：`/marketing/coupon/user`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 使用记录 | `R` | `merchantCouponIssue` | 菜单权限 |

#### 营销 / 优惠套餐

- 页面路由：`/marketing/discounts/list`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 添加 | `C` | `merchantStoreDiscountsCreate` | 菜单权限 |
| 编辑 | `U` | `merchantStoreDiscountsUpdate` | 菜单权限 |
| 列表 | `R` | `merchantStoreDiscountsLst` | 菜单权限 |
| 详情 | `R` | `merchantStoreDiscountsDetail` | 菜单权限 |
| 删除 | `D` | `merchantStoreDiscountsDelete` | 菜单权限 |
| 修改状态 | `U` | `merchantStoreDiscountsStatus` | 菜单权限 |
| 上传图片 | `C` | `merchantUploadImage` | 菜单权限 |
| 图片列表 | `R` | `merchantAttachmentLst` | 菜单权限 |

#### 营销 / 助力

- 页面路由：`/assist`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 营销 / 助力 / 助力商品

- 页面路由：`/marketing/assist/list`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表  | `R` | `merchantStoreProductAssistLst` | 菜单权限 |
| 添加 | `C` | `merchantStoreProductAssistCreate` | 菜单权限 |
| 详情 | `R` | `merchantStoreProductAssistDetail` | 菜单权限 |
| 编辑 | `U` | `merchantStoreProductAssistUpdate` | 菜单权限 |
| 删除 | `D` | `merchantStoreProductAssistDelete` | 菜单权限 |
| 修改状态 | `U` | `merchantStoreProductAssistStatus` | 菜单权限 |
| 排序 | `U` | `merchantStoreProductAssistUpdateSort` | 菜单权限 |
| 预览 | `R` | `merchantStoreProductAssistPreview` | 菜单权限 |
| 设置标签 | `U` | `merchantStoreProductAssistLabels` | 菜单权限 |
| 上传图片 | `C` | `merchantUploadImage` | 菜单权限 |
| 图片列表 | `R` | `merchantAttachmentLst` | 菜单权限 |

#### 营销 / 助力 / 助力活动

- 页面路由：`/marketing/assist/assist_set`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 活动列表 | `R` | `merchantStoreProductAssistSetLst` | 菜单权限 |
| 活动详情 | `R` | `merchantStoreProductAssistSetDetail` | 菜单权限 |

#### 营销 / 拼团

- 页面路由：`/marketing/combination`
- CRUD：C=✓ R=— U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 预览 | `C` | `merchantStoreProductGroupPreview` | 路由 |
| 设置标签 | `C` | `merchantStoreProductGroupLabels` | 路由 |

#### 营销 / 拼团 / 拼团商品列表

- 页面路由：`/marketing/combination/combination_goods`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `merchantStoreProductGroupLst` | 菜单权限 |
| 添加 | `C` | `merchantStoreProductGroupCreate` | 菜单权限 |
| 详情 | `R` | `merchantStoreProductGroupDetail` | 菜单权限 |
| 编辑 | `U` | `merchantStoreProductGroupUpdate` | 菜单权限 |
| 删除 | `D` | `merchantStoreProductGroupDelete` | 菜单权限 |
| 修改状态 | `U` | `merchantStoreProductGroupStatus` | 菜单权限 |
| 排序 | `U` | `merchantStoreProductGroupSort` | 菜单权限 |
| 预览 | `R` | `merchantStoreProductGroupPreview` | 菜单权限 |
| 设置标签 | `U` | `merchantStoreProductGroupLabels` | 菜单权限 |
| 拼团配置 | `U` | `merchantConfigGroupBuying` | 菜单权限 |
| 上传图片 | `C` | `merchantUploadImage` | 菜单权限 |
| 图片列表 | `R` | `merchantAttachmentLst` | 菜单权限 |

#### 营销 / 拼团 / 拼团活动列表

- 页面路由：`/marketing/combination/combination_list`
- CRUD：C=✓ R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 活动列表  | `R` | `merchantStoreProductGroupBuyingLst` | 菜单权限 |
| 活动详情 | `R` | `merchantStoreProductGroupBuyingDetail` | 菜单权限 |
| 上传图片 | `C` | `merchantUploadImage` | 菜单权限 |
| 图片列表 | `R` | `merchantAttachmentLst` | 菜单权限 |

#### 营销 / 直播

- 页面路由：`/`
- CRUD：C=✓ R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 退保证金检测 | `C` | `merchantFinancialRefundMargin` | 路由 |
| 退保证金申请 | `C` | `merchantFinancialRefundMarginApply` | 路由 |
| 收款方式表单 | `R` | `merchantFinancialAccountForm` | 路由 |
| 到店自提信息 | `R` | `merchantTakeInfo` | 路由 |
| 保存到店自提信息 | `C` | `merchantTakeUpdate` | 路由 |

#### 营销 / 直播 / 直播助手

- 页面路由：`/marketing/studio/assistant`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `merchantBroadcastAssistantLst` | 菜单权限 |
| 添加 | `C` | `merchantBroadcastAssistantCreate` | 菜单权限 |
| 编辑 | `U` | `merchantBroadcastAssistantUpdate` | 菜单权限 |
| 备注 | `U` | `merchantBroadcastAssistantMark` | 菜单权限 |
| 删除 | `D` | `merchantBroadcastAssistantDelete` | 菜单权限 |
| 上传图片 | `C` | `merchantUploadImage` | 菜单权限 |
| 图片列表 | `R` | `merchantAttachmentLst` | 菜单权限 |

#### 营销 / 直播 / 直播商品管理

- 页面路由：`/marketing/broadcast/list`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 营销 / 直播 / 直播间管理

- 页面路由：`/marketing/studio/list`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表  | `R` | `merchantBroadcastRoomLst` | 菜单权限 |
| 详情 | `R` | `merchantBroadcastRoomDetail` | 菜单权限 |
| 添加 | `C` | `merchantBroadcastRoomCreate` | 菜单权限 |
| 编辑 | `U` | `merchantBroadcastRoomUpdate` | 菜单权限 |
| 修改状态 | `U` | `merchantBroadcastRoomChangeStatus` | 菜单权限 |
| 导入商品 | `CR` | `merchantBroadcastRoomExportGoods` | 菜单权限 |
| 删除商品 | `RD` | `merchantBroadcastRoomRmExportGoods` | 菜单权限 |
| 备注 | `U` | `merchantBroadcastRoomMark` | 菜单权限 |
| 商品详情 | `R` | `merchantBroadcastRoomGoods` | 菜单权限 |
| 关闭客服 | `O` | `merchantBroadcastRoomCloseKf` | 菜单权限 |
| 禁言 | `O` | `merchantBroadcastRoomCloseComment` | 菜单权限 |
| 收录 | `O` | `merchantBroadcastRoomCloseFeeds` | 菜单权限 |
| 商品上下架 | `U` | `merchantBroadcastOnSale` | 菜单权限 |
| 删除 | `D` | `merchantBroadcastRoomDelete` | 菜单权限 |
| 添加 客服 | `C` | `merchantBroadcastAddAssistant` | 菜单权限 |
| 消息推送 | `O` | `merchantBroadcastPushMessage` | 菜单权限 |
| 列表 | `R` | `merchantBroadcastGoodsLst` | 菜单权限 |
| 详情 | `R` | `merchantBroadcastGoodsDetail` | 菜单权限 |
| 添加 | `C` | `merchantBroadcastGoodsCreate` | 菜单权限 |
| 编辑 | `U` | `merchantBroadcastGoodsUpdate` | 菜单权限 |
| 修改状态 | `U` | `merchantBroadcastGoodsChangeStatus` | 菜单权限 |
| 备注 | `U` | `merchantBroadcastGoodsMark` | 菜单权限 |
| 删除 | `D` | `merchantBroadcastGoodsDelete` | 菜单权限 |
| 批量添加 | `C` | `merchantBroadcastGoodsbatchCreate` | 菜单权限 |
| 上传图片 | `C` | `merchantUploadImage` | 菜单权限 |
| 图片列表 | `R` | `merchantAttachmentLst` | 菜单权限 |

#### 营销 / 秒杀

- 页面路由：`/marketing/seckill/list`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 营销 / 秒杀 / 秒杀商品

- 页面路由：`/marketing/seckill/product/list`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 统计 | `R` | `merchantStoreSeckillProductLstFilter` | 菜单权限 |
| 列表 | `R` | `merchantStoreSeckillProductLst` | 菜单权限 |
| 添加  | `C` | `merchantStoreSeckillProductCreate` | 菜单权限 |
| 详情 | `R` | `merchantStoreSeckillProductDetail` | 菜单权限 |
| 编辑 | `U` | `merchantStoreSeckillProductUpdate` | 菜单权限 |
| 删除 | `D` | `merchantStoreSeckillProductDelete` | 菜单权限 |
| 彻底删除 | `D` | `merchantStoreSeckillProductDestory` | 菜单权限 |
| 恢复 | `O` | `merchantStoreSeckillProductRestore` | 菜单权限 |
| 修改状态 | `U` | `merchantStoreSeckillProductSwitchStatus` | 菜单权限 |
| 排序 | `U` | `merchantStoreSeckillProductUpdateSort` | 菜单权限 |
| 预览 | `R` | `merchantStoreSeckillProductPreview` | 菜单权限 |
| 设置标签 | `U` | `merchantStoreSeckillProductLabels` | 菜单权限 |
| 分页列表 | `R` | `merchantStoreSeckillProductPageLst` | 菜单权限 |
| 商品列表 | `R` | `merchantStoreSeckillProductGetProductList` | 菜单权限 |
| 设置标签 | `U` | `merchantStoreSeckillProductSetLabels` | 菜单权限 |

#### 营销 / 秒杀 / 秒杀活动

- 页面路由：`/marketing/seckill/store_seckill/list`
- CRUD：C=✓ R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 上传图片 | `C` | `merchantUploadImage` | 菜单权限 |
| 图片列表 | `R` | `merchantAttachmentLst` | 菜单权限 |
| 列表 | `R` | `merchantStoreSeckillActiveGetActiveList` | 菜单权限 |
| 详情 | `R` | `merchantStoreSeckillActiveGetActiveInfo` | 菜单权限 |
| 列表 | `R` | `merchantStoreSeckillActiveGetActiveAll` | 菜单权限 |
| 活动统计数据面板 | `R` | `merchantStoreSeckillActiveChartPanel` | 菜单权限 |
| 活动参与人统计列表 | `R` | `merchantStoreSeckillActiveChartPeople` | 菜单权限 |
| 活动订单统计列表 | `R` | `merchantStoreSeckillActiveChartOrder` | 菜单权限 |
| 活动商品统计列表 | `R` | `merchantStoreSeckillActiveChartProduct` | 菜单权限 |

#### 营销 / 积分

- 页面路由：`/marketing/integral`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `merchantIntegralList` | 路由 |
| 统计 | `R` | `merchantIntegralTitle` | 路由 |

#### 营销 / 积分 / 积分日志

- 页面路由：`/marketing/integral/log`
- CRUD：C=✓ R=✓ U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `merchantIntegralList` | 菜单权限 |
| 统计 | `R` | `merchantIntegralTitle` | 菜单权限 |
| 配置获取 | `RU` | `merchantConfigForm` | 菜单权限 |
| 配置保存 | `CU` | `merchantConfigSave` | 菜单权限 |

#### 营销 / 积分 / 积分配置

- 页面路由：`/marketing/integral/config`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 营销 / 逛逛社区

- 页面路由：`/community/list`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 分类列表 | `R` | `merchantCommunityCateLst` | 菜单权限 |
| 列表 | `R` | `merchantCommunityLst` | 菜单权限 |
| 添加 | `C` | `merchantCommunityCreate` | 菜单权限 |
| 详情 | `R` | `merchantCommunityDetail` | 菜单权限 |
| 编辑 | `U` | `merchantCommunityUpdate` | 菜单权限 |
| 删除 | `D` | `merchantCommunityDelete` | 菜单权限 |
| 评论 | `O` | `merchantCommunityReply` | 菜单权限 |

#### 营销 / 预售

- 页面路由：`/marketing/presell/list`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `merchantStoreProductPresellLst` | 菜单权限 |
| 添加 | `C` | `merchantStoreProductPresellCreate` | 菜单权限 |
| 详情 | `R` | `merchantStoreProductPresellDetail` | 菜单权限 |
| 编辑 | `U` | `merchantStoreProductPresellUpdate` | 菜单权限 |
| 删除 | `D` | `merchantStoreProductPresellDelete` | 菜单权限 |
| 修改状态 | `U` | `merchantStoreProductPresellStatus` | 菜单权限 |
| 排序 | `U` | `merchantStoreProductPresellUpdateSort` | 菜单权限 |
| 预览 | `R` | `merchantStoreProductPresellPreview` | 菜单权限 |
| 设置标签 | `U` | `merchantStoreProductPreselltLabels` | 菜单权限 |
| 上传图片 | `C` | `merchantUploadImage` | 菜单权限 |
| 图片列表 | `R` | `merchantAttachmentLst` | 菜单权限 |

### 装修

#### 装修

- 页面路由：`/devise/`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `merchantDiyPageLinkLst` | 路由 |
| 列表  | `R` | `merchantDiyLst` | 路由 |
| 默认模板列表  | `R` | `merchantDefaultDiyLst` | 路由 |
| 详情  | `R` | `merchantDiyDetail` | 路由 |
| 添加/编辑 | `CU` | `merchantDiyCreate` | 路由 |
| 使用模板 | `C` | `merchantDiyStatus` | 路由 |
| 使用模板 | `C` | `merchantDiySetDefault` | 路由 |
| 重置模板 | `R` | `merchantDiyRecovery` | 路由 |
| 当前使用模板 | `R` | `merchantDiyInfo` | 路由 |
| 删除 | `D` | `merchantDiyDelete` | 路由 |
| 店铺街装修 | `R` | `merchantDiyProductLst` | 路由 |
| 复制 | `C` | `merchantDiyCopy` | 路由 |

#### 装修 / 商品分类

- 页面路由：`/devise/diy/product_category`
- CRUD：C=✓ R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 商品分类信息 | `R` | `merchantDiyProductCategoryInfo` | 菜单权限 |
| 保存商品分类 | `C` | `merchantDiyProductCategorySave` | 菜单权限 |

#### 装修 / 系统表单

- 页面路由：`/systemForm/form_list`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 添加 | `C` | `merFormCreate` | 菜单权限 |
| 编辑 | `U` | `merFormUpdate` | 菜单权限 |
| 删除 | `D` | `merFormDelete` | 菜单权限 |
| 详情 | `R` | `merFormDetail` | 菜单权限 |
| 列表 | `R` | `merFormLst` | 菜单权限 |

#### 装修 / 素材管理

- 页面路由：`/config/picture`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `merchantAttachmentLst` | 菜单权限 |
| 删除 | `D` | `merchantAttachmentDelete` | 菜单权限 |
| 批量修改 | `U` | `merchantAttachmentBatchChangeCategory` | 菜单权限 |
| 编辑 | `U` | `merchantAttachmentUpdate` | 菜单权限 |
| 分类列表 | `R` | `merchantAttachmentCategoryGetFormatList` | 菜单权限 |
| 添加 | `C` | `merchantAttachmentCategoryCreate` | 菜单权限 |
| 编辑 | `U` | `merchantAttachmentCategoryUpdate` | 菜单权限 |
| 删除 | `D` | `merchantAttachmentCategoryDelete` | 菜单权限 |
| 上传二维码 | `C` | `merchantAttachmentScanQrcode` | 菜单权限 |
| 扫码上传图片 | `C` | `merchantAttachmentScanImage` | 菜单权限 |
| 扫码上传保存 | `C` | `merchantAttachmentScanImageSave` | 菜单权限 |
| 在线图片 | `O` | `merchantAttachmentOnline` | 菜单权限 |

#### 装修 / 装修

- 页面路由：`/devise/diy/list`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表  | `R` | `merchantDiyLst` | 菜单权限 |
| 添加/编辑 | `CU` | `merchantDiyCreate` | 菜单权限 |
| 使用模板 | `U` | `merchantDiyStatus` | 菜单权限 |
| 设置模版默认数据 | `U` | `merchantDiySetDefault` | 菜单权限 |
| 重置模板 | `O` | `merchantDiyRecovery` | 菜单权限 |
| 当前使用模板 | `R` | `merchantDiyInfo` | 菜单权限 |
| 删除 | `D` | `merchantDiyDelete` | 菜单权限 |
| 店铺街装修 | `R` | `merchantDiyProductLst` | 菜单权限 |
| 复制 | `C` | `merchantDiyCopy` | 菜单权限 |
| 上传图片 | `C` | `uploadImage` | 菜单权限 |
| 图片列表 | `R` | `systemAttachmentLst` | 菜单权限 |
| 默认模板列表  | `R` | `merchantDefaultDiyLst` | 菜单权限 |
| 详情  | `R` | `merchantDiyDetail` | 菜单权限 |

### 订单

#### 订单

- 页面路由：`/order`
- CRUD：C=✓ R=— U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 备注 | `C` | `merchantOrderReceiptMark` | 路由 |
| 编辑 | `U` | `merchantOrderReceiptUpdate` | 路由 |

#### 订单 / 代客下单

- 页面路由：`/order/customer`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 商品分类 | `O` | `behalfProductCategory` | 菜单权限 |
| 商品列表 | `R` | `behalfProductList` | 菜单权限 |
| 商品规格详情 | `R` | `behalfProductDetail` | 菜单权限 |
| 会员查询 | `R` | `behalfUserQuery` | 菜单权限 |
| 会员详情 | `R` | `behalfUserInfo` | 菜单权限 |
| 会员添加 | `C` | `behalfUserCreate` | 菜单权限 |
| 地址列表 | `CR` | `behalfUserAddressList` | 菜单权限 |
| 地址添加 | `C` | `behalfUserAddressCreate` | 菜单权限 |
| 购物车列表 | `R` | `behalfCartList` | 菜单权限 |
| 添加购物车 | `C` | `behalfCartCreate` | 菜单权限 |
| 修改购物车数据 | `U` | `behalfCartChange` | 菜单权限 |
| 删除购物数据 | `D` | `behalfCartDelete` | 菜单权限 |
| 清空购物车 | `D` | `behalfCartClear` | 菜单权限 |
| 购物车总数量 | `O` | `behalfCartCount` | 菜单权限 |
| 修改价格 | `U` | `behalfCartUpdatePrice` | 菜单权限 |
| 批量修改价格 | `U` | `behalfCartBatchUpdatePrice` | 菜单权限 |
| 校验订单 | `O` | `behalfCheck` | 菜单权限 |
| 支付配置 | `U` | `behalfPayConfig` | 菜单权限 |
| 创建订单 | `C` | `behalfCreate` | 菜单权限 |
| 支付 | `O` | `behalfPay` | 菜单权限 |
| 获取结果 | `RU` | `behalfPayStatus` | 菜单权限 |
| 余额支付获取验证码 | `R` | `behalfVerify` | 菜单权限 |

#### 订单 / 商品评价

- 页面路由：`/product/reviews`
- CRUD：C=✓ R=✓ U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `merchantProductReplyLst` | 菜单权限 |
| 回复表单 | `U` | `merchantProductReplyForm` | 菜单权限 |
| 回复 | `U` | `merchantProductReplyReply` | 菜单权限 |
| 排序 | `U` | `merchantProductReplySort` | 菜单权限 |
| 上传图片 | `C` | `merchantUploadImage` | 菜单权限 |
| 图片列表 | `R` | `merchantAttachmentLst` | 菜单权限 |

#### 订单 / 核销记录

- 页面路由：`/order/cancellation`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 统计 | `R` | `merchantStoreOrderTakeTitle` | 菜单权限 |
| 列表 | `R` | `merchantStoreTakeOrderLst` | 菜单权限 |

#### 订单 / 订单管理

- 页面路由：`/order/list`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 预览 | `R` | `merchantServeExportTemps` | 菜单权限 |
| 默认模板 | `R` | `merchantServeExportDumpLst` | 菜单权限 |
| 导出 | `R` | `merchantStoreOrderExcel` | 菜单权限 |
| 打印小票 | `O` | `merchantStoreOrderPrinter` | 菜单权限 |
| 统计 | `R` | `merchantStoreOrderTitle` | 菜单权限 |
| 列表 | `R` | `merchantStoreOrderLst` | 菜单权限 |
| 快递查询 | `R` | `merchantStoreOrderExpress` | 菜单权限 |
| 发货 | `U` | `merchantStoreOrderDelivery` | 菜单权限 |
| 批量发货 | `U` | `merchantStoreOrderBatchDelivery` | 菜单权限 |
| 导出发货单 | `RU` | `merchantStoreOrderDeliveryExport` | 菜单权限 |
| 头部统计 | `R` | `merchantStoreOrderStat` | 菜单权限 |
| 编辑 | `U` | `merchantStoreOrderUpdate` | 菜单权限 |
| 详情 | `R` | `merchantStoreOrderDetail` | 菜单权限 |
| 操作记录 | `R` | `merchantStoreOrderLog` | 菜单权限 |
| 备注 | `U` | `merchantStoreOrderRemark` | 菜单权限 |
| 核销 | `U` | `merchantStoreOrderVerify` | 菜单权限 |
| 删除 | `D` | `merchantStoreOrderDelete` | 菜单权限 |
| 导入 | `C` | `merchantStoreOrderDeliveryImport` | 菜单权限 |
| 导入记录 | `CR` | `merchantStoreOrderDeliveryImportLst` | 菜单权限 |
| 详情 | `CR` | `merchantStoreOrderDeliveryImportDetail` | 菜单权限 |
| 导出发货记录 | `CRU` | `merchantStoreOrderDeliveryImportExcel` | 菜单权限 |
| 导出列表 | `R` | `merchantStoreExcelLst` | 菜单权限 |
| 导出下载 | `R` | `merchantStoreExcelDownload` | 菜单权限 |
| 核销详情 | `RU` | `merchantStoreOrderVerifyDetail` | 菜单权限 |
| 关联订单 | `O` | `merchantStoreOrderChildrenList` | 菜单权限 |
| 线下支付 | `O` | `merchantStoreOrderOffline` | 菜单权限 |
| 电子面单复打 | `O` | `merchantStoreOrderRepeatDump` | 菜单权限 |
| 配货单 | `O` | `merchantStoreOrderNote` | 菜单权限 |
| 修改收货信息 | `U` | `merchantStoreOrderCollectCargo` | 菜单权限 |
| 预约订单派单 | `O` | `merchantStoreOrderReservationDispatch` | 菜单权限 |
| 预约订单改派 | `U` | `merchantStoreOrderReservationUpdateDispatch` | 菜单权限 |
| 预约订单改约 | `O` | `merchantStoreOrderReservationReschedule` | 菜单权限 |
| 单独修改预约时间 | `U` | `merchantStoreOrderReservationTime` | 菜单权限 |
| 预约订单核销 | `U` | `merchantStoreOrderReservationVerify` | 菜单权限 |
| 获取商家寄件价格 | `R` | `merchantStoreOrderGetPrice` | 菜单权限 |
| 获取商家寄件价格 | `R` | `merchantStoreOrderShipmentList` | 菜单权限 |
| 取消商家寄件 | `D` | `merchantStoreOrderCancelShipment` | 菜单权限 |
| 同城配送派单 | `O` | `merchantStoreOrderDeliveryDispatch` | 菜单权限 |
| 同城配送改派 | `U` | `merchantStoreOrderDeliveryUpdateDispatch` | 菜单权限 |
| 同城配送核销 | `U` | `merchantStoreOrderDeliveryConfirm` | 菜单权限 |
| 同城配送再次同步 | `U` | `merchantStoreOrderDeliverySync` | 菜单权限 |

#### 订单 / 退款订单

- 页面路由：`/order/refund`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `merchantStoreRefundOrderLst` | 菜单权限 |
| 详情 | `R` | `merchantStoreRefundOrderDetail` | 菜单权限 |
| 审核 | `U` | `merchantStoreRefundOrderSwitchStatus` | 菜单权限 |
| 收到退回商品后确认退款 | `D` | `merchantStoreRefundOrderRefund` | 菜单权限 |
| 删除 | `D` | `merchantStoreRefundDelete` | 菜单权限 |
| 备注 | `U` | `merchantStoreRefundMark` | 菜单权限 |
| 操作记录 | `R` | `merchantStoreRefundLog` | 菜单权限 |
| 快递查询 | `R` | `merchantStoreRefundExpress` | 菜单权限 |
| 导出 | `CR` | `merchantStoreRefundCreateExcel` | 菜单权限 |
| 导出列表 | `R` | `merchantStoreExcelLst` | 菜单权限 |
| 导出下载 | `R` | `merchantStoreExcelDownload` | 菜单权限 |
| 创建 | `C` | `merchantStoreRefundOrderCreate` | 菜单权限 |

#### 订单 / 预约服务

- 页面路由：`/order/reservation`
- CRUD：C=— R=— U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 预约日历 | `O` | `merchantReservationServiceList` | 菜单权限 |

#### 订单 / 预约设置

- 页面路由：`/product/reservation`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

### 设置

#### 设置

- 页面路由：`/config`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 编辑 | `U` | `merchantAttachmentUpdate` | 路由 |
| 列表 | `R` | `merchantAttachmentCategoryGetFormatList` | 路由 |
| 添加表单 | `C` | `merchantAttachmentCategoryCreateForm` | 路由 |
| 编辑表单 | `U` | `merchantAttachmentCategoryUpdateForm` | 路由 |
| 添加 | `C` | `merchantAttachmentCategoryCreate` | 路由 |
| 编辑 | `U` | `merchantAttachmentCategoryUpdate` | 路由 |
| 删除 | `D` | `merchantAttachmentCategoryDelete` | 路由 |

#### 设置 / 一号通

- 页面路由：`/one_setting`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 设置 / 一号通 / 平台一号通

- 页面路由：`/setting/sms/sms_config/index`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 套餐列表 | `R` | `merchantServeMeal` | 菜单权限 |
| 支付二维码 | `O` | `merchantServeCode` | 菜单权限 |
| 购买记录 | `R` | `merchantServeLst` | 菜单权限 |
| 详情 | `R` | `merchantServeDetail` | 菜单权限 |
| 账号信息 | `R` | `merchantServeInfo` | 菜单权限 |

#### 设置 / 一号通 / 自有一号通

- 页面路由：`/setting/sms/sms_account/index`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 设置 / 一号通 / 配置管理

- 页面路由：`/setting/sms/dumpConfig`
- CRUD：C=✓ R=— U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 保存配置 | `CU` | `merchantServeSetConfig` | 菜单权限 |

#### 设置 / 付费会员

- 页面路由：`/systemForm/Basics/svip`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 设置 / 同城配送

- 页面路由：`/delivery`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 城市列表 | `R` | `merchantStoreDeliveryCityList` | 路由 |
| 充值记录 | `R` | `merchantStoreDeliveryPayLst` | 路由 |
| 充值二维码 | `R` | `merchantStoreDeliveryGetQrcode` | 路由 |

#### 设置 / 同城配送 / 充值记录

- 页面路由：`/delivery/recharge_record`
- CRUD：C=✓ R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 充值记录 | `CR` | `merchantStoreDeliveryPayLst` | 菜单权限 |

#### 设置 / 同城配送 / 发货点管理

- 页面路由：`/delivery/store_manage`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `merchantDeliveryServiceLst` | 路由 |
| 修改状态 | `U` | `merchantDeliveryServiceStatus` | 路由 |
| 添加表单 | `C` | `merchantDeliveryServiceCreateForm` | 路由 |
| 添加 | `C` | `merchantDeliveryServiceCreate` | 路由 |
| 编辑 | `U` | `merchantDeliveryServiceUpdateForm` | 路由 |
| 编辑 | `U` | `merchantDeliveryServiceUpdate` | 路由 |
| 删除 | `D` | `merchantDeliveryServiceDelete` | 路由 |
| 统计列表 | `R` | `merchantDeliveryServiceStatisticsList` | 路由 |

#### 设置 / 同城配送 / 配送记录

- 页面路由：`/delivery/usage_record`
- CRUD：C=— R=✓ U=— D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `merchantStoreDeliveryOrderLst` | 菜单权限 |
| 取消 | `D` | `merchantStoreDeliveryOrderCancel` | 菜单权限 |
| 详情 | `R` | `merchantStoreDeliveryOrderDetail` | 菜单权限 |

#### 设置 / 同城配送 / 配送设置

- 页面路由：`/setting/delivery`
- CRUD：C=— R=— U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 更新配送设置信息 | `U` | `merchantDeliveryConfigUpdate` | 菜单权限 |
| 配送设置信息 | `U` | `merchantDeliveryConfigSettings` | 菜单权限 |

#### 设置 / 同城配送 / 配送门店

- 页面路由：`/delivery/delivery_point`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 获取分类 | `R` | `merchantStoreDeliveryBusiness` | 菜单权限 |
| 添加 | `C` | `merchantStoreDeliveryCreate` | 菜单权限 |
| 编辑 | `U` | `merchantStoreDeliveryUpdate` | 菜单权限 |
| 编辑状态 | `U` | `merchantStoreDeliveryStatus` | 菜单权限 |
| 列表 | `R` | `merchantStoreDeliveryLst` | 菜单权限 |
| 详情 | `R` | `merchantStoreDeliveryDetail` | 菜单权限 |
| 删除 | `D` | `merchantStoreDeliveryDelete` | 菜单权限 |
| 备注 | `U` | `merchantStoreDeliveryMark` | 菜单权限 |
| 城市列表 | `R` | `merchantStoreDeliveryCityList` | 菜单权限 |

#### 设置 / 店铺信息

- 页面路由：`/systemForm/modifyStoreInfo`
- CRUD：C=✓ R=— U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 保存到店自提信息 | `CU` | `merchantTakeUpdate` | 菜单权限 |
| 退保证金申请 | `C` | `merchantFinancialRefundMarginApply` | 菜单权限 |

#### 设置 / 店铺配置

- 页面路由：`/systemForm/Basics/mer_base`
- CRUD：C=— R=— U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 资料更新 | `U` | `merchantUpdate` | 菜单权限 |

#### 设置 / 开放账户

- 页面路由：`/systemForm/openAuth/list`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `merchantOpenapiLst` | 菜单权限 |
| 添加 | `C` | `merchantOpenapiCreate` | 菜单权限 |
| 编辑 | `U` | `merchantOpenapiUpdate` | 菜单权限 |
| 修改状态 | `U` | `merchantOpenapiStatus` | 菜单权限 |
| 删除 | `D` | `merchantOpenapiDeleta` | 菜单权限 |
| 查看 | `R` | `merchantOpenapiGetSecretKey` | 菜单权限 |
| 重置 | `O` | `merchantOpenapiSetSecretKey` | 菜单权限 |

#### 设置 / 快递配送

- 页面路由：`/city`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 设置 / 快递配送 / 物流公司

- 页面路由：`/config/freight/express`
- CRUD：C=— R=✓ U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `merchantServeExportLst` | 菜单权限 |
| 月结账号编辑 | `U` | `merchantExpressPratnerUpdate` | 菜单权限 |
| 修改状态 | `U` | `merchantExpressChangeMerStatus` | 菜单权限 |

#### 设置 / 快递配送 / 运费模板

- 页面路由：`/config/freight/shippingTemplates`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 添加  | `C` | `merchantStoreShippingTemplateCreate` | 菜单权限 |
| 编辑 | `U` | `merchantStoreShippingTemplateUpdate` | 菜单权限 |
| 详情 | `R` | `merchantStoreShippingTemplateDetail` | 菜单权限 |
| 删除 | `D` | `merchantStoreShippingTemplateDelete` | 菜单权限 |
| 设置默认模板 | `U` | `merchantStoreShippingTemplateSetDefault` | 菜单权限 |

#### 设置 / 打印配置

- 页面路由：`/setting/printer`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 所有数据 | `R` | `merchantStatisticsMain` | 路由 |
| 支付订单 | `R` | `merchantStatisticsOrder` | 路由 |
| 成交客户 | `R` | `merchantStatisticsUser` | 路由 |
| 成交客户比 | `R` | `merchantStatisticsUserRate` | 路由 |
| 商品支付排行 | `R` | `merchantStatisticsProduct` | 路由 |
| 商品访问排行 | `R` | `merchantStatisticsProductVisit` | 路由 |
| 商品加购排行 | `R` | `merchantStatisticsProductCart` | 路由 |
| 首页未处理业务统计 | `R` | `merchantStatisticsMerchantCount` | 路由 |

#### 设置 / 打印配置 / 小票打印

- 页面路由：`/setting/printer/list`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 获取配置内容 | `RU` | `merchantStorePrinterGetContent` | 菜单权限 |
| 列表 | `R` | `merchantStorePrinterLst` | 菜单权限 |
| 添加 | `C` | `merchantStorePrinterCreate` | 菜单权限 |
| 编辑 | `U` | `merchantStorePrinterUpdate` | 菜单权限 |
| 取消 | `UD` | `merchantStorePrinterStatus` | 菜单权限 |
| 删除 | `D` | `merchantStorePrinterDelete` | 菜单权限 |
| 保存配置内容 | `CU` | `merchantStorePrinterSetContent` | 菜单权限 |

#### 设置 / 打印配置 / 打印配置

- 页面路由：`/systemForm/Basics/printer_tabs`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 设置 / 权限管理

- 页面路由：`/setting`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 一键换色 | `R` | `merchantGetChangeColor` | 路由 |

#### 设置 / 权限管理 / 操作日志

- 页面路由：`/setting/systemLog`
- CRUD：C=— R=✓ U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 操作日志 | `R` | `merchantAdminLog` | 路由 |
| 修改信息表单 | `U` | `merchantAdminEditForm` | 路由 |
| 修改信息 | `U` | `merchantAdminEdit` | 路由 |
| 修改密码表单 | `U` | `merchantAdminEditPasswordForm` | 路由 |

#### 设置 / 权限管理 / 管理员管理

- 页面路由：`/setting/systemAdmin`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `merchantAdminLst` | 菜单权限 |
| 修改状态 | `U` | `merchantAdminStatus` | 菜单权限 |
| 添加 | `C` | `merchantAdminCreate` | 菜单权限 |
| 编辑 | `U` | `merchantAdminUpdate` | 菜单权限 |
| 修改密码 | `U` | `merchantAdminPassword` | 菜单权限 |
| 删除 | `D` | `merchantAdminDelete` | 菜单权限 |

#### 设置 / 权限管理 / 身份管理

- 页面路由：`/setting/systemRole`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 修改状态 | `U` | `merchantRoleStatus` | 路由 |
| 删除 | `D` | `merchantRoleDelete` | 路由 |
| 列表 | `R` | `merchantAdminLst` | 路由 |
| 修改状态 | `U` | `merchantAdminStatus` | 路由 |
| 添加 | `C` | `merchantAdminCreate` | 路由 |
| 添加表单 | `C` | `merchantAdminCreateForm` | 路由 |
| 编辑 | `U` | `merchantAdminUpdate` | 路由 |
| 编辑表单 | `U` | `merchantAdminUpdateForm` | 路由 |
| 修改密码 | `U` | `merchantAdminPassword` | 路由 |
| 修改密码表单 | `U` | `merchantAdminPasswordForm` | 路由 |

### 财务

#### 财务

- 页面路由：`/accounts`
- CRUD：C=✓ R=✓ U=— D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 收款方式 | `C` | `merchantFinancialAccountSave` | 路由 |
| 转账记录 | `R` | `merchantFinancialLst` | 路由 |
| 详情 | `R` | `merchantFinancialDetail` | 路由 |
| 申请表单 | `C` | `merchantFinancialCreateForm` | 路由 |
| 申请 | `C` | `merchantFinancialCreateSave` | 路由 |
| 删除 | `D` | `merchantFinancialDelete` | 路由 |
| 备注表单 | `R` | `merchantFinancialMarkForm` | 路由 |
| 备注 | `C` | `merchantFinancialMark` | 路由 |

#### 财务 / 分账管理

- 页面路由：`/systemForm/applyList`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `merchantOrderProfitsharingLst` | 菜单权限 |
| 导出 | `R` | `merchantOrderProfitsharingExport` | 菜单权限 |
| 导出列表 | `R` | `merchantStoreExcelLst` | 菜单权限 |
| 导出下载 | `R` | `merchantStoreExcelDownload` | 菜单权限 |

#### 财务 / 发票管理

- 页面路由：`/order/invoice`
- CRUD：C=✓ R=✓ U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `merchantOrderReceiptLst` | 菜单权限 |
| 详情 | `R` | `merchantOrderReceiptDetail` | 菜单权限 |
| 开发票 | `O` | `merchantOrderReceiptSetRecipt` | 菜单权限 |
| 保存发票 | `C` | `merchantOrderReceiptSave` | 菜单权限 |
| 备注 | `U` | `merchantOrderReceiptMark` | 菜单权限 |
| 编辑 | `U` | `merchantOrderReceiptUpdate` | 菜单权限 |

#### 财务 / 收款方式

- 页面路由：`/accounts/payType`
- CRUD：C=✓ R=— U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 收款方式 | `C` | `merchantFinancialAccountSave` | 菜单权限 |

#### 财务 / 申请分账商户

- 页面路由：`/systemForm/applyments`
- CRUD：C=✓ R=✓ U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 申请 | `C` | `merchantApplymentsCreate` | 菜单权限 |
| 详情 | `R` | `merchantApplymentsDetail` | 菜单权限 |
| 编辑 | `U` | `merchantApplymentsUpdate` | 菜单权限 |
| 上传图片 | `C` | `merchantApplymentsUpload` | 菜单权限 |
| 上传视频 | `C` | `merchantApplymentsUploadVideo` | 菜单权限 |

#### 财务 / 账单管理

- 页面路由：`/accounts/statement`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `merchantFinanciaRecordlLst` | 菜单权限 |
| 统计 | `R` | `merchantFinancialTitle` | 菜单权限 |
| 详情 | `R` | `merchantFinancialRecordDetail` | 菜单权限 |
| 导出 | `R` | `merchantFinancialRecordDetailExport` | 菜单权限 |
| 导出列表 | `R` | `merchantStoreExcelLst` | 菜单权限 |
| 导出下载 | `R` | `merchantStoreExcelDownload` | 菜单权限 |

#### 财务 / 资金流水

- 页面路由：`/accounts/capitalFlow`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `merchantFinancialRecordList` | 菜单权限 |
| 导出 | `R` | `merchantFinancialRecordExport` | 菜单权限 |
| 统计 | `R` | `merchantFinancialCount` | 菜单权限 |
| 导出列表 | `R` | `merchantStoreExcelLst` | 菜单权限 |
| 导出下载 | `R` | `merchantStoreExcelDownload` | 菜单权限 |

#### 财务 / 转账记录

- 页面路由：`/accounts/transManagement`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `merchantFinancialRefundMargin` | 菜单权限 |
| 转账记录 | `R` | `merchantFinancialLst` | 菜单权限 |
| 详情 | `R` | `merchantFinancialDetail` | 菜单权限 |
| 申请 | `C` | `merchantFinancialCreateSave` | 菜单权限 |
| 删除 | `D` | `merchantFinancialDelete` | 菜单权限 |
| 备注 | `U` | `merchantFinancialMark` | 菜单权限 |
| 导出 | `R` | `merchantFinancialExport` | 菜单权限 |
| 导出列表 | `R` | `merchantStoreExcelLst` | 菜单权限 |
| 导出下载 | `R` | `merchantStoreExcelDownload` | 菜单权限 |

### 首页

#### 首页

- 页面路由：`/dashboard`
- CRUD：C=— R=✓ U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 所有数据 | `O` | `merchantStatisticsMain` | 菜单权限 |
| 支付订单 | `O` | `merchantStatisticsOrder` | 菜单权限 |
| 成交客户 | `O` | `merchantStatisticsUser` | 菜单权限 |
| 成交客户比 | `O` | `merchantStatisticsUserRate` | 菜单权限 |
| 商品支付排行 | `O` | `merchantStatisticsProduct` | 菜单权限 |
| 商品访问排行 | `O` | `merchantStatisticsProductVisit` | 菜单权限 |
| 商品加购排行 | `O` | `merchantStatisticsProductCart` | 菜单权限 |
| 首页未处理业务统计 | `RU` | `merchantStatisticsMerchantCount` | 菜单权限 |
| 待办事项 | `O` | `merchantStatisticsMerchantTodo` | 菜单权限 |
| 获取商户代办统计 | `R` | `merchantStatisticsProductSalesPriceTop` | 菜单权限 |

#### 首页 / 商品统计

- 页面路由：`/statistic/product`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 顶部统计 | `R` | `merchantAnalyticsProductTop` | 菜单权限 |
| 折线图统计 | `R` | `merchantAnalyticsProductLineChart` | 菜单权限 |
| 折线图统计 | `R` | `merchantAnalyticsProductTypePieChart` | 菜单权限 |

#### 首页 / 控制台

- 页面路由：`/dashboard`
- CRUD：C=✓ R=✓ U=— D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 待办事项 | `R` | `merchantStatisticsMerchantTodo` | 路由 |
| 获取商户代办统计 | `R` | `merchantStatisticsProductSalesPriceTop` | 路由 |
| 列表 | `R` | `systemNoticeLogList` | 路由 |
| 已读 | `C` | `systemNoticeLogRead` | 路由 |
| 删除 | `D` | `systemNoticeLogDel` | 路由 |

#### 首页 / 订单统计

- 页面路由：`/statistic/order`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 顶部统计 | `R` | `merchantAnalyticsOrderTop` | 菜单权限 |
| 折线图统计 | `R` | `merchantAnalyticsOrderLineChart` | 菜单权限 |
| 折线图统计 | `R` | `merchantAnalyticsOrderTypePieChart` | 菜单权限 |
