export const storageTabDefs = [
  { desc: '上传大小限制与目录', label: '基础设置', name: 'basic' },
  { desc: 'S3 兼容对象存储（自建/私有化）', label: 'MinIO', name: 'minio' },
  { desc: '七牛云对象存储参数', label: '七牛云', name: 'qiniu' },
  { desc: '阿里云对象存储参数', label: '阿里云 OSS', name: 'aliyun' },
  { desc: '腾讯云对象存储参数', label: '腾讯云 COS', name: 'qcloud' },
] as const;

export const tencentTabDefs = [
  {
    desc: '推拉流域名、鉴权 Key、回调地址；live-api 从此处加载',
    label: '云直播',
    name: 'live',
  },
  { desc: '录播上传/播放；与云直播独立配置', label: '云点播', name: 'vod' },
  {
    desc: '云直播与云点播共用的播放器 License',
    label: 'License',
    name: 'lic',
  },
] as const;

export const settingTabDefs = [
  { desc: '平台超管登录页与安全策略', label: '系统配置', name: 'admin' },
  { desc: '商户后台登录展示', label: '商城配置', name: 'shop' },
  { desc: '微信服务商支付', label: '微信配置', name: 'weixin' },
  { desc: '文件存储与上传限制', label: '上传配置', name: 'storage' },
  { desc: '云直播、云点播与 License', label: '腾讯配置', name: 'tencent' },
] as const;

export function defaultTencentLive() {
  return {
    auth_expire_sec: 86400,
    callback_key: '',
    callback_url: '',
    play_auth_key: '',
    play_domain: '',
    push_auth_key: '',
    push_domain: '',
    region: 'ap-shanghai',
    secret_id: '',
    secret_key: '',
    tencent_app_id: 0,
  };
}

export function defaultTencentVod() {
  return {
    play_domain: '',
    play_key: '',
    region: 'ap-chongqing',
    secret_id: '',
    secret_key: '',
    storage_region: 'ap-chongqing',
    sub_app_id: 0,
    tencent_app_id: 0,
  };
}

export function defaultTencentLicense() {
  return {
    license_key: '',
    license_url: '',
  };
}

export function defaultWeixinService() {
  return {
    apikey: '',
    app_id: '',
    cert_pem: '',
    is_open: 0,
    key_pem: '',
    mch_id: '',
    serial_no: '',
  };
}

export function defaultMinioEngine() {
  return {
    access_key: '',
    bucket: '',
    domain: '',
    endpoint: '',
    region: 'us-east-1',
    secret_key: '',
    use_path_style: 1,
  };
}

export function createDefaultForm() {
  return {
    admin_bg_img: '',
    admin_code: 0,
    admin_lock_enable: 1,
    admin_lock_time: 15,
    admin_max_attempts: 5,
    admin_name: '',
    shop_bg_img: '',
    shop_code: 0,
    shop_lock_enable: 1,
    shop_lock_time: 15,
    shop_logo_img: '',
    shop_max_attempts: 5,
    shop_name: '',
    storage: {
      default: 'minio',
      directory: '',
      engine: {
        aliyun: {} as Record<string, string>,
        minio: defaultMinioEngine(),
        qcloud: {} as Record<string, string>,
        qiniu: {} as Record<string, string>,
      },
      max_image: '',
      max_video: '',
    },
    tencent_license: defaultTencentLicense(),
    tencent_live: defaultTencentLive(),
    tencent_vod: defaultTencentVod(),
    weixin_service: defaultWeixinService(),
  };
}

export type SettingFormModel = ReturnType<typeof createDefaultForm>;

export function mergeSettingValues(vars: Record<string, unknown>): SettingFormModel {
  const form = createDefaultForm();
  Object.assign(form, vars);

  form.weixin_service = {
    ...defaultWeixinService(),
    ...((vars.weixin_service as object) || {}),
  };
  form.weixin_service.is_open = Number.parseInt(
    String(form.weixin_service.is_open),
    10,
  ) || 0;

  form.admin_code = Number.parseInt(String(vars.admin_code), 10) || 0;
  form.shop_code = Number.parseInt(String(vars.shop_code), 10) || 0;
  form.admin_lock_enable = Number.parseInt(
    String(vars.admin_lock_enable ?? 1),
    10,
  );
  form.admin_lock_time = Number.parseInt(String(vars.admin_lock_time ?? 15), 10);
  form.admin_max_attempts = Number.parseInt(
    String(vars.admin_max_attempts ?? 5),
    10,
  );
  form.shop_lock_enable = Number.parseInt(String(vars.shop_lock_enable ?? 1), 10);
  form.shop_lock_time = Number.parseInt(String(vars.shop_lock_time ?? 15), 10);
  form.shop_max_attempts = Number.parseInt(
    String(vars.shop_max_attempts ?? 5),
    10,
  );

  const storage = vars.storage as SettingFormModel['storage'] | undefined;
  if (storage) {
    form.storage.default = storage.default || 'minio';
    form.storage.max_image = storage.max_image || '';
    form.storage.max_video = storage.max_video || '';
    form.storage.directory = storage.directory || '';
    if (storage.engine) {
      form.storage.engine.minio = {
        ...defaultMinioEngine(),
        ...(storage.engine.minio || {}),
      };
      form.storage.engine.qiniu = storage.engine.qiniu || {};
      form.storage.engine.aliyun = storage.engine.aliyun || {};
      form.storage.engine.qcloud = storage.engine.qcloud || {};
    }
  }

  if (form.storage.default === 'local') {
    form.storage.default = 'minio';
  }
  if (!['minio', 'qiniu', 'aliyun', 'qcloud'].includes(form.storage.default)) {
    form.storage.default = 'minio';
  }
  form.storage.engine.minio.use_path_style = Number.parseInt(
    String(form.storage.engine.minio.use_path_style ?? 1),
    10,
  )
    ? 1
    : 0;

  form.tencent_live = {
    ...defaultTencentLive(),
    ...((vars.tencent_live as object) || {}),
  };
  form.tencent_vod = {
    ...defaultTencentVod(),
    ...((vars.tencent_vod as object) || {}),
  };
  form.tencent_license = {
    ...defaultTencentLicense(),
    ...((vars.tencent_license as object) || {}),
  };

  const vod = vars.tencent_vod as Record<string, string> | undefined;
  if (!form.tencent_license.license_url && vod?.license_url) {
    form.tencent_license.license_url = vod.license_url;
  }
  if (!form.tencent_license.license_key && vod?.license_key) {
    form.tencent_license.license_key = vod.license_key;
  }

  form.tencent_live.auth_expire_sec =
    Number.parseInt(String(form.tencent_live.auth_expire_sec), 10) || 86400;
  form.tencent_live.tencent_app_id =
    Number.parseInt(String(form.tencent_live.tencent_app_id), 10) || 0;
  form.tencent_vod.sub_app_id =
    Number.parseInt(String(form.tencent_vod.sub_app_id), 10) || 0;
  form.tencent_vod.tencent_app_id =
    Number.parseInt(String(form.tencent_vod.tencent_app_id), 10) || 0;

  return form;
}

export function buildSavePayload(form: SettingFormModel) {
  const params = JSON.parse(JSON.stringify(form)) as SettingFormModel;

  params.admin_code = Number.parseInt(String(params.admin_code), 10) || 0;
  params.admin_lock_enable =
    Number.parseInt(String(params.admin_lock_enable), 10) || 0;
  params.admin_lock_time =
    Number.parseInt(String(params.admin_lock_time), 10) || 15;
  params.admin_max_attempts =
    Number.parseInt(String(params.admin_max_attempts), 10) || 5;

  params.shop_code = Number.parseInt(String(params.shop_code), 10) || 0;
  params.shop_lock_enable =
    Number.parseInt(String(params.shop_lock_enable), 10) || 0;
  params.shop_lock_time =
    Number.parseInt(String(params.shop_lock_time), 10) || 15;
  params.shop_max_attempts =
    Number.parseInt(String(params.shop_max_attempts), 10) || 5;

  if (params.weixin_service) {
    params.weixin_service.is_open =
      Number.parseInt(String(params.weixin_service.is_open), 10) || 0;
  }

  if (params.tencent_live) {
    params.tencent_live.auth_expire_sec =
      Number.parseInt(String(params.tencent_live.auth_expire_sec), 10) || 86400;
    params.tencent_live.tencent_app_id =
      Number.parseInt(String(params.tencent_live.tencent_app_id), 10) || 0;
  }

  if (params.tencent_vod) {
    params.tencent_vod.sub_app_id =
      Number.parseInt(String(params.tencent_vod.sub_app_id), 10) || 0;
    params.tencent_vod.tencent_app_id =
      Number.parseInt(String(params.tencent_vod.tencent_app_id), 10) || 0;
  }

  if (params.storage?.engine?.minio) {
    params.storage.engine.minio.use_path_style = Number.parseInt(
      String(params.storage.engine.minio.use_path_style),
      10,
    )
      ? 1
      : 0;
  }

  return params;
}
