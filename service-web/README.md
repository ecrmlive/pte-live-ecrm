# service-web — 客服工作台

阶段 6b 轻量：静态查单页（`public/index.html`），对接：

- `/api/service/v1/orders/:id?mer_id=`（Header `X-Service-Token: service_demo_token`）
- 可选演示 OpenAPI：`/api/open/v1/auth` + `/order/detail/:id`

发布：宿主机 Nginx `:18084` 挂载本目录 `public/`（见 release）。
