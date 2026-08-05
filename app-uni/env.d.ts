import {
  computed as computedOrigin,
  createApp as createVueAppOrigin,
  nextTick as nextTickOrigin,
  onMounted as onMountedOrigin,
  ref as refOrigin,
  watch as watchOrigin,
} from 'vue'
import {
  onLaunch as onLaunchOrigin,
  onLoad as onLoadOrigin,
  onReachBottom as onReachBottomOrigin,
  onShow as onShowOrigin,
  onUnload as onUnloadOrigin,
} from '@dcloudio/uni-app'

declare module "vue" {
  export const createSSRApp: typeof createVueAppOrigin
  export const createVueApp: typeof createVueAppOrigin
}

declare global {
  /** uni-app x 页面脚本使用的组合式 API 由编译器自动注入。 */
  const ref: typeof refOrigin
  const computed: typeof computedOrigin
  const watch: typeof watchOrigin
  const nextTick: typeof nextTickOrigin
  const onMounted: typeof onMountedOrigin

  /** uni-app x 页面生命周期由编译器自动注入。 */
  const onLaunch: typeof onLaunchOrigin
  const onLoad: typeof onLoadOrigin
  const onReachBottom: typeof onReachBottomOrigin
  const onShow: typeof onShowOrigin
  const onUnload: typeof onUnloadOrigin

  /** H5 端以普通 JSON 对象承载 UTSJSONObject。 */
  type UTSJSONObject = Record<string, any>

  /** H5 入口转换使用的 Vue 应用工厂别名。 */
  const createVueApp: typeof createVueAppOrigin
}
