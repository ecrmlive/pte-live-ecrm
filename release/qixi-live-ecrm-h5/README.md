# qixi-live-ecrm-h5

用户端 H5 静态产物（源码 `app-uni/`）。宿主机 Nginx 托管，无 Docker 容器。

```bash
make local-h5   # npm run build:h5 → 本目录 dist/
```

微信小程序产物不进本目录：`cd app-uni && npm run build:mp-weixin` → `dist/build/mp-weixin`。
