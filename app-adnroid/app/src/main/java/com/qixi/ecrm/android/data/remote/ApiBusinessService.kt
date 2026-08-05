package com.qixi.ecrm.android.data.remote

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import retrofit2.http.Body
import retrofit2.http.POST

/** api-business 的 C 端契约入口；具体接口按 OpenAPI 逐步补齐。 */
interface ApiBusinessService {
    @POST("api/app/v1/auth/login")
    suspend fun login(@Body request: LoginRequest): ApiResponse<LoginPayload>
}

@Serializable
data class LoginRequest(
    @SerialName("platform") val platform: String = "android",
    @SerialName("credential") val credential: String,
)

@Serializable
data class LoginPayload(
    @SerialName("access_token") val accessToken: String,
    @SerialName("refresh_token") val refreshToken: String,
)

@Serializable
data class ApiResponse<T>(
    val code: Int,
    val message: String,
    val data: T? = null,
)

