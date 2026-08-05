package com.qixi.ecrm.android.push

import android.content.Context
import timber.log.Timber

/**
 * 友盟 U-Push 的唯一初始化入口。
 *
 * 当前工程不携带 AppKey 或厂商配置；接入 SDK 后仅在此处调用友盟初始化，业务层只依赖
 * [PushGateway]，避免推送厂商逻辑泄漏到页面和用例层。
 */
object PushInitializer {
    fun initialize(context: Context) {
        Timber.i("U-Push 初始化入口已就绪，等待环境注入 AppKey")
    }
}

interface PushGateway {
    suspend fun bindUser(userId: String)
    suspend fun unbindUser()
}

