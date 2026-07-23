<template>

	<div v-loading="loading" class="setting-layout">
		<aside class="setting-sidebar">
			<nav class="setting-nav">
				<div v-for="tab in settingTabDefs" :key="tab.name" class="setting-nav__group">
					<button
						type="button"
						:class="['setting-nav__item', settingTab === tab.name ? 'is-active' : '']"
						@click="selectMainTab(tab.name)"
					>
						<span class="setting-nav__label">{{ tab.label }}</span>
						<span v-if="tab.desc" class="setting-nav__desc">{{ tab.desc }}</span>
					</button>

					<div v-if="tab.name === 'storage' && settingTab === 'storage'" class="setting-subnav">
						<button
							v-for="sub in storageTabDefs"
							:key="sub.name"
							type="button"
							:class="['setting-subnav__item', storageTab === sub.name ? 'is-active' : '']"
							@click="selectStorageTab(sub.name)"
						>
							{{ sub.label }}
						</button>
					</div>

					<div v-if="tab.name === 'tencent' && settingTab === 'tencent'" class="setting-subnav">
						<button
							v-for="sub in tencentTabDefs"
							:key="sub.name"
							type="button"
							:class="['setting-subnav__item', tencentTab === sub.name ? 'is-active' : '']"
							@click="tencentTab = sub.name"
						>
							{{ sub.label }}
						</button>
					</div>

				</div>
			</nav>
		</aside>

		<section class="setting-main">
			<header class="setting-main__header">
				<h3 class="setting-main__title">{{ currentTabLabel }}</h3>
				<p v-if="currentTabDesc" class="setting-main__desc">{{ currentTabDesc }}</p>
			</header>

		<el-form ref="formRef" :model="form" label-width="150px" class="setting-form">
			<template v-if="settingTab === 'admin'">
			<el-form-item label="商城运营名称" :rules="[{required: true,message: ' '}]" prop="shop_name">
				<el-input v-model="form.admin_name" placeholder="商城运营名称" class="max-w460"></el-input>
				<div class="tips">
					saas端商城运营名称，显示在登录页
				</div>
			</el-form-item>
			<el-form-item label="商城运营登录背景" prop="shop_bg_img">
				<el-input v-model="form.admin_bg_img" placeholder="商城运营登录背景" class="max-w460"></el-input>
				<div class="tips">
					saas端商城运营登录背景，不填则为系统默认登录背景，填写网络地址
				</div>
			</el-form-item>
			<el-form-item label="是否开启登录验证码">
				<el-switch
					v-model="form.admin_code"
					:active-value="1"
					:inactive-value="0"
					active-text="开启"
					inactive-text="关闭"
					inline-prompt
				/>
			</el-form-item>
			<el-form-item label="是否开启账号锁定">
				<el-switch
					v-model="form.admin_lock_enable"
					:active-value="1"
					:inactive-value="0"
					active-text="开启"
					inactive-text="关闭"
					inline-prompt
				/>
				<div class="tips">
					开启后，登录失败达到最大次数将锁定账号
				</div>
			</el-form-item>
			<el-form-item label="最大失败次数" v-if="form.admin_lock_enable == 1">
				<el-input v-model.number="form.admin_max_attempts" type="number" placeholder="最大失败次数" class="max-w460">
					<template v-slot:append>次</template>
				</el-input>
				<div class="tips">
					允许连续登录失败的最大次数，建议设置3-10次
				</div>
			</el-form-item>
			<el-form-item label="账号锁定时间" v-if="form.admin_lock_enable == 1">
				<el-input v-model.number="form.admin_lock_time" type="number" placeholder="账号锁定时间" class="max-w460">
					<template v-slot:append>分钟</template>
				</el-input>
				<div class="tips">
					登录失败达到最大次数后，账号将被锁定的时间，建议不低于5分钟
				</div>
			</el-form-item>
			</template>

			<template v-if="settingTab === 'shop'">
			<el-form-item label="商城系统名称" :rules="[{required: true,message: ' '}]" prop="shop_name">
				<el-input v-model="form.shop_name" placeholder="商城名称" class="max-w460"></el-input>
				<div class="tips">
					shop端商城名称，显示在登录页
				</div>
			</el-form-item>
			<el-form-item label="商城登录背景" prop="shop_bg_img">
				<el-input v-model="form.shop_bg_img" placeholder="商城登录背景" class="max-w460"></el-input>
				<div class="tips">
					shop端商城登录背景，不填则为系统默认登录背景，填写网络地址
				</div>
			</el-form-item>
			<el-form-item label="商城登录logo" prop="shop_logo_img">
				<el-input v-model="form.shop_logo_img" placeholder="商城登录logo" class="max-w460"></el-input>
				<div class="tips">
					shop端商城登录logo，不填则为系统默认登录log，填写网络地址
				</div>
			</el-form-item>
			<el-form-item label="是否开启登录验证码">
				<el-switch
					v-model="form.shop_code"
					:active-value="1"
					:inactive-value="0"
					active-text="开启"
					inactive-text="关闭"
					inline-prompt
				/>
			</el-form-item>
			<el-form-item label="是否开启账号锁定">
				<el-switch
					v-model="form.shop_lock_enable"
					:active-value="1"
					:inactive-value="0"
					active-text="开启"
					inactive-text="关闭"
					inline-prompt
				/>
				<div class="tips">
					开启后，登录失败达到最大次数将锁定账号
				</div>
			</el-form-item>
			<el-form-item label="最大失败次数" v-if="form.shop_lock_enable == 1">
				<el-input v-model.number="form.shop_max_attempts" type="number" placeholder="最大失败次数" class="max-w460">
					<template v-slot:append>次</template>
				</el-input>
				<div class="tips">
					允许连续登录失败的最大次数，建议设置3-10次
				</div>
			</el-form-item>
			<el-form-item label="账号锁定时间" v-if="form.shop_lock_enable == 1">
				<el-input v-model.number="form.shop_lock_time" type="number" placeholder="账号锁定时间" class="max-w460">
					<template v-slot:append>分钟</template>
				</el-input>
				<div class="tips">
					登录失败达到最大次数后，账号将被锁定的时间，建议不低于5分钟
				</div>
			</el-form-item>
			</template>

			<template v-if="settingTab === 'weixin'">
			<el-form-item label="是否开启服务商支付">
				<el-switch
					v-model="form.weixin_service.is_open"
					:active-value="1"
					:inactive-value="0"
					active-text="开启"
					inactive-text="关闭"
					inline-prompt
				/>
			</el-form-item>
			<div v-if="form.weixin_service.is_open == 1">
				<el-form-item label="服务商户号">
					<el-input v-model="form.weixin_service.mch_id" placeholder="服务商户号" class="max-w460"></el-input>
					<div class="tips">
						填写服务商户号、10位数字
					</div>
				</el-form-item>
				<el-form-item label="服务商支付秘钥apikey">
					<el-input v-model="form.weixin_service.apikey" placeholder="服务商支付秘钥apikey"
						class="max-w460"></el-input>
					<div class="tips">
						填写服务商户支付秘钥apikey
					</div>
				</el-form-item>
				<el-form-item label="服务商appid">
					<el-input v-model="form.weixin_service.app_id" placeholder="服务商appid" class="max-w460"></el-input>
					<div class="tips">
						填写服务商户号绑定的公众号appid
					</div>
				</el-form-item>
				<el-form-item label="微信支付公钥">
					<el-input v-model="form.weixin_service.serial_no" placeholder="微信支付公钥" class="max-w460"></el-input>
					<div class="tips">
					  微信商户平台-账号中心-API安全-微信支付公钥-公钥ID
					</div>
				</el-form-item>
				<el-form-item label="apiclient_cert.pem">
					<el-input type="textarea" :rows="4" placeholder="使用文本编辑器打开apiclient_cert.pem文件，将文件的全部内容复制进来"
						v-model="form.weixin_service.cert_pem" class="max-w460"></el-input>
					<div class="tips">使用文本编辑器打开apiclient_key.pem文件，将文件的全部内容复制进来</div>
				</el-form-item>
				<el-form-item label="apiclient_key.pem">
					<el-input type="textarea" :rows="4" placeholder="使用文本编辑器打开apiclient_cert.pem文件，将文件的全部内容复制进来"
						v-model="form.weixin_service.key_pem" class="max-w460"></el-input>
					<div class="tips">使用文本编辑器打开apiclient_key.pem文件，将文件的全部内容复制进来</div>
				</el-form-item>
			</div>
			</template>

			<template v-if="settingTab === 'storage'">
			<template v-if="storageTab === 'basic'">
			<el-form-item label="最大图片上传">
				<el-input v-model="form.storage.max_image" class="max-w460">
					<template v-slot:append>M</template>
				</el-input>
				<div class="tips">修改后请修改php上传配置后生效</div>
			</el-form-item>
			<el-form-item label="最大视频上传">
				<el-input v-model="form.storage.max_video" class="max-w460">
					<template v-slot:append>M</template>
				</el-input>
				<div class="tips">修改后请修改php上传配置后生效</div>
			</el-form-item>
			<el-form-item label="上传目录">
				<el-input v-model.trim="form.storage.directory" class="max-w460"></el-input>
				<div class="tips">上传目录仅对 MinIO、七牛云、阿里云、腾讯云有效</div>
			</el-form-item>
			<el-form-item label="默认上传方式">
				<el-select
					v-model="form.storage.default"
					class="max-w460"
					placeholder="请选择默认上传方式"
					teleported
					popper-class="setting-select-popper"
				>
					<el-option
						v-for="item in storageEngineOptions"
						:key="item.name"
						:label="item.label"
						:value="item.name"
					/>
				</el-select>
				<div class="tips">可分别配置多个云存储参数，此处选择实际上传时使用的引擎</div>
			</el-form-item>
			</template>

			<template v-if="storageTab === 'minio'">
				<div class="tips block-tips">MinIO 为 S3 兼容对象存储，适用于自建或私有化部署。设为默认方式请在「基础设置」中选择。</div>
				<el-form-item label="Endpoint">
					<el-input v-model="form.storage.engine.minio.endpoint" class="max-w460" placeholder="http://127.0.0.1:9000"></el-input>
					<div class="tips">API 地址，可带 http(s)://；内网示例：http://minio:9000</div>
				</el-form-item>
				<el-form-item label="存储空间名称 Bucket">
					<el-input v-model="form.storage.engine.minio.bucket" class="max-w460"></el-input>
				</el-form-item>
				<el-form-item label="Access Key">
					<el-input v-model="form.storage.engine.minio.access_key" class="max-w460"></el-input>
				</el-form-item>
				<el-form-item label="Secret Key">
					<el-input v-model="form.storage.engine.minio.secret_key" type="password" show-password class="max-w460" placeholder="留空表示不修改"></el-input>
				</el-form-item>
				<el-form-item label="Region">
					<el-input v-model="form.storage.engine.minio.region" class="max-w460" placeholder="us-east-1"></el-input>
					<div class="tips">MinIO 通常使用 us-east-1</div>
				</el-form-item>
				<el-form-item label="Path Style">
					<el-switch v-model="form.storage.engine.minio.use_path_style" :active-value="1" :inactive-value="0"></el-switch>
					<div class="tips">多数 MinIO 部署需开启 Path Style</div>
				</el-form-item>
				<el-form-item label="访问域名 Domain">
					<el-input v-model="form.storage.engine.minio.domain" class="max-w460"></el-input>
					<div class="tips">请补全 https://，用于文件公网访问，例如：https://minio-cdn.qxkejiwl.top</div>
				</el-form-item>
			</template>

			<!--七牛云存储-->
			<template v-if="storageTab === 'qiniu'">
				<el-form-item label="存储空间名称 Bucket"><el-input v-model="form.storage.engine.qiniu.bucket"
						class="max-w460"></el-input></el-form-item>
				<el-form-item label="ACCESS_KEY AK"><el-input v-model="form.storage.engine.qiniu.access_key"
						class="max-w460"></el-input></el-form-item>
				<el-form-item label="SECRET_KEY SK"><el-input v-model="form.storage.engine.qiniu.secret_key"
						class="max-w460"></el-input></el-form-item>
				<el-form-item label="空间域名 Domain">
					<el-input v-model="form.storage.engine.qiniu.domain" class="max-w460"></el-input>
					<div class="tips">
						请补全https://，例如：https://static.cloud.com
					</div>
				</el-form-item>
			</template>
			<!--阿里云OSS-->
			<template v-if="storageTab === 'aliyun'">
				<el-form-item label="存储空间名称 Bucket"><el-input v-model="form.storage.engine.aliyun.bucket"
						class="max-w460"></el-input></el-form-item>
				<el-form-item label="AccessKeyId"><el-input v-model="form.storage.engine.aliyun.access_key_id"
						class="max-w460"></el-input></el-form-item>
				<el-form-item label="AccessKeySecret"><el-input v-model="form.storage.engine.aliyun.access_key_secret"
						class="max-w460"></el-input></el-form-item>
				<el-form-item label="空间域名 Domain">
					<el-input v-model="form.storage.engine.aliyun.domain" class="max-w460"></el-input>
					<div class="tips">
						请补全https://，例如：https://static.cloud.com
					</div>
				</el-form-item>
			</template>
			<!--腾讯云COS-->
			<template v-if="storageTab === 'qcloud'">
				<el-form-item label="存储空间名称 Bucket">
					<el-input v-model="form.storage.engine.qcloud.bucket" class="max-w460"></el-input>
				</el-form-item>
				<el-form-item label="所属地域 Region">
					<el-input v-model="form.storage.engine.qcloud.region" class="max-w460" placeholder="ap-shanghai"></el-input>
					<div class="tips">
						请填写地域简称，例如：ap-beijing、ap-hongkong、eu-frankfurt
					</div>
				</el-form-item>
				<el-form-item label="SecretId">
					<el-input v-model="form.storage.engine.qcloud.secret_id" class="max-w460"></el-input>
				</el-form-item>
				<el-form-item label="SecretKey">
					<el-input v-model="form.storage.engine.qcloud.secret_key" type="password" show-password class="max-w460" placeholder="留空表示不修改"></el-input>
				</el-form-item>
				<el-form-item label="空间域名 Domain">
					<el-input v-model="form.storage.engine.qcloud.domain" class="max-w460"></el-input>
					<div class="tips">
						请补全https://，例如：https://static.cloud.com
					</div>
				</el-form-item>
			</template>
			</template>

			<template v-if="settingTab === 'tencent'">
			<template v-if="tencentTab === 'live'">
				<div class="tips block-tips">推拉流域名、鉴权 Key、回调地址；live-api 从此处加载。AppName 按租户 app_id 自动生成。</div>
				<el-form-item label="SecretId">
					<el-input v-model="form.tencent_live.secret_id" class="max-w460"></el-input>
				</el-form-item>
				<el-form-item label="SecretKey">
					<el-input v-model="form.tencent_live.secret_key" type="password" show-password class="max-w460" placeholder="留空表示不修改"></el-input>
				</el-form-item>
				<el-form-item label="Region">
					<el-input v-model="form.tencent_live.region" class="max-w460" placeholder="ap-shanghai"></el-input>
				</el-form-item>
				<el-form-item label="推流域名">
					<el-input v-model="form.tencent_live.push_domain" class="max-w460" placeholder="push.qxkejiwl.top"></el-input>
				</el-form-item>
				<el-form-item label="播放域名">
					<el-input v-model="form.tencent_live.play_domain" class="max-w460" placeholder="play.qxkejiwl.top"></el-input>
				</el-form-item>
				<el-form-item label="推流鉴权 Key">
					<el-input v-model="form.tencent_live.push_auth_key" class="max-w460"></el-input>
				</el-form-item>
				<el-form-item label="播放鉴权 Key">
					<el-input v-model="form.tencent_live.play_auth_key" class="max-w460"></el-input>
				</el-form-item>
				<el-form-item label="签名有效期(秒)">
					<el-input-number
						v-model="form.tencent_live.auth_expire_sec"
						class="max-w460"
						controls-position="right"
						:min="60"
						:max="604800"
					/>
				</el-form-item>
				<el-form-item label="回调密钥">
					<el-input v-model="form.tencent_live.callback_key" class="max-w460"></el-input>
				</el-form-item>
				<el-form-item label="回调 URL">
					<el-input v-model="form.tencent_live.callback_url" class="max-w460"></el-input>
					<div class="tips">须与腾讯云直播控制台回调地址一致</div>
				</el-form-item>
				<el-form-item label="腾讯云 AppId">
					<el-input-number
						v-model="form.tencent_live.tencent_app_id"
						class="max-w460"
						controls-position="right"
						:min="0"
					/>
				</el-form-item>
			</template>
			<template v-if="tencentTab === 'vod'">
				<div class="tips block-tips">录播上传/播放；与云直播独立配置。PlayKey 留空表示不修改。</div>
				<el-form-item label="SecretId">
					<el-input v-model="form.tencent_vod.secret_id" class="max-w460"></el-input>
				</el-form-item>
				<el-form-item label="SecretKey">
					<el-input v-model="form.tencent_vod.secret_key" type="password" show-password class="max-w460" placeholder="留空表示不修改"></el-input>
				</el-form-item>
				<el-form-item label="Region">
					<el-input v-model="form.tencent_vod.region" class="max-w460" placeholder="ap-chongqing"></el-input>
					<div class="tips">VOD API 接入地域；多数场景与存储地域一致</div>
				</el-form-item>
				<el-form-item label="存储地域">
					<el-input v-model="form.tencent_vod.storage_region" class="max-w460" placeholder="ap-chongqing"></el-input>
				</el-form-item>
				<el-form-item label="SubAppId">
					<el-input-number
						v-model="form.tencent_vod.sub_app_id"
						class="max-w460"
						controls-position="right"
						:min="0"
					/>
				</el-form-item>
				<el-form-item label="腾讯云 AppId">
					<el-input-number
						v-model="form.tencent_vod.tencent_app_id"
						class="max-w460"
						controls-position="right"
						:min="0"
					/>
				</el-form-item>
				<el-form-item label="播放域名">
					<el-input v-model="form.tencent_vod.play_domain" class="max-w460" placeholder="可选"></el-input>
				</el-form-item>
				<el-form-item label="播放 Key">
					<el-input v-model="form.tencent_vod.play_key" type="password" show-password class="max-w460" placeholder="留空表示不修改"></el-input>
				</el-form-item>
			</template>
			<template v-if="tencentTab === 'lic'">
				<div class="tips block-tips">云直播与云点播共用；License Key 留空表示不修改。</div>
				<el-form-item label="License URL">
					<el-input v-model="form.tencent_license.license_url" class="max-w460" placeholder="播放器 License 地址"></el-input>
				</el-form-item>
				<el-form-item label="License Key">
					<el-input v-model="form.tencent_license.license_key" type="password" show-password class="max-w460" placeholder="留空表示不修改"></el-input>
				</el-form-item>
			</template>
			</template>

			<div class="button-wrapper">
				<el-button
					v-access:code="'platform:setting:edit'"
					type="primary"
					:loading="saving"
					@click="saveCurrentModule"
				>
					保存{{ saveActionLabel }}
				</el-button>
			</div>
		</el-form>
		</section>
	</div>

</template>

<script lang="ts" setup>
import { computed, onMounted, ref } from 'vue';

import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElOption,
  ElSelect,
  ElSwitch,
} from 'element-plus';

import LiveConfigApi from '#/api/core/live-config';
import SettingApi from '#/api/core/setting';
import { useAccessStore } from '@vben/stores';

import {
  buildSavePayload,
  createDefaultForm,
  mergeSettingValues,
  settingTabDefs,
  storageTabDefs,
  tencentTabDefs,
} from './setting-form-data';

defineOptions({ name: 'SettingForm' });

const formRef = ref<InstanceType<typeof ElForm>>();
const loading = ref(false);
const saving = ref(false);
const settingTab = ref('admin');
const storageTab = ref('basic');
const tencentTab = ref('live');
const form = ref(createDefaultForm());

const storageEngineOptions = computed(() =>
  storageTabDefs.filter((item) => item.name !== 'basic'),
);

const currentSubTabLabel = computed(() => {
  if (settingTab.value === 'storage') {
    const sub = storageTabDefs.find((item) => item.name === storageTab.value);
    return sub ? sub.label : '';
  }
  if (settingTab.value === 'tencent') {
    const sub = tencentTabDefs.find((item) => item.name === tencentTab.value);
    return sub ? sub.label : '';
  }
  return '';
});

const currentTabLabel = computed(() => {
  const tab = settingTabDefs.find((item) => item.name === settingTab.value);
  if (!tab) {
    return '';
  }
  if (currentSubTabLabel.value) {
    return `${tab.label} / ${currentSubTabLabel.value}`;
  }
  return tab.label;
});

const currentTabDesc = computed(() => {
  const tab = settingTabDefs.find((item) => item.name === settingTab.value);
  if (!tab) {
    return '';
  }
  if (settingTab.value === 'storage') {
    const sub = storageTabDefs.find((item) => item.name === storageTab.value);
    return sub ? sub.desc : tab.desc;
  }
  if (settingTab.value === 'tencent') {
    const sub = tencentTabDefs.find((item) => item.name === tencentTab.value);
    return sub ? sub.desc : tab.desc;
  }
  return tab.desc;
});

const saveActionLabel = computed(() => currentTabLabel.value || '配置');

function selectMainTab(name: string) {
  settingTab.value = name;
}

function selectStorageTab(name: string) {
  storageTab.value = name;
}

async function getParams() {
  loading.value = true;
  try {
    const res = await SettingApi.serviceDetail({}, true);
    const vars = (res.data.vars.values || {}) as Record<string, unknown>;
    form.value = mergeSettingValues(vars);
  } catch {
    // handled by request interceptor
  } finally {
    loading.value = false;
  }
}

async function saveCurrentModule() {
  try {
    await formRef.value?.validate();
  } catch {
    return;
  }

  saving.value = true;
  try {
    const params = buildSavePayload(form.value);
    await SettingApi.editService(params, true);
    ElMessage.success(`${saveActionLabel.value}保存成功`);
    const accessStore = useAccessStore();
    const codes = accessStore.accessCodes ?? [];
    if (
      settingTab.value === 'tencent' &&
      codes.includes('platform:setting:reloadTencent')
    ) {
      LiveConfigApi.reloadTencent().catch(() => {});
    }
    await getParams();
  } finally {
    saving.value = false;
  }
}

onMounted(() => {
  getParams();
});
</script>
<style lang="scss">
@use './setting-popper.scss';
@use './setting.scss';
</style>