/**
 * 直播运行时 HTTP — 与商城 API 同源（api-platform :11503）。
 * native 路径 `/api/v1/shop/live/*`、`/api/v1/room/*` 等（目标 api-platform 进程内 Go；
 * 由 api-platform 进程内 Go 路由处理。
 * 浏览器不直连 api-live :11501（api-live 仅 C 端 / 主播端）。
 */
export { requestClient as liveRequestClient } from '#/api/request';
