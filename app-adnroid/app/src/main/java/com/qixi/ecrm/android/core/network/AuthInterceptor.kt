package com.qixi.ecrm.android.core.network

import com.qixi.ecrm.android.core.session.SessionStore
import okhttp3.Interceptor
import okhttp3.Response

/** C 端 API 约定只使用 Authori-zation 请求头。 */
class AuthInterceptor(
    private val sessionStore: SessionStore,
) : Interceptor {
    override fun intercept(chain: Interceptor.Chain): Response {
        val token = sessionStore.accessToken()
        val request = chain.request().newBuilder()
            .apply {
                if (token != null) {
                    header("Authori-zation", "Bearer $token")
                }
            }
            .build()
        return chain.proceed(request)
    }
}

