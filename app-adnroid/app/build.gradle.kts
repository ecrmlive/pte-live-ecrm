plugins {
    alias(libs.plugins.android.application)
    alias(libs.plugins.kotlin.android)
    alias(libs.plugins.kotlin.compose)
    alias(libs.plugins.kotlin.serialization)
    alias(libs.plugins.ksp)
    alias(libs.plugins.hilt)
}

val apiBaseUrl = providers.gradleProperty("ECRM_API_BASE_URL")
    .orElse("https://api.example.invalid/")
    .get()
val umengAppKey = providers.gradleProperty("UMENG_APPKEY").orElse("").get()
val applicationIdValue = providers.gradleProperty("ECRM_APPLICATION_ID")
    .orElse("com.qixi.ecrm.android")
    .get()
val versionCodeValue = providers.gradleProperty("ECRM_VERSION_CODE")
    .orElse("1")
    .get()
    .toIntOrNull()
    ?: error("ECRM_VERSION_CODE 必须为正整数")
require(versionCodeValue > 0) { "ECRM_VERSION_CODE 必须为正整数" }
val versionNameValue = providers.gradleProperty("ECRM_VERSION_NAME")
    .orElse("0.1.0")
    .get()
val releaseStoreFile = providers.gradleProperty("ECRM_RELEASE_STORE_FILE").orNull
val releaseStorePassword = providers.gradleProperty("ECRM_RELEASE_STORE_PASSWORD").orNull
val releaseKeyAlias = providers.gradleProperty("ECRM_RELEASE_KEY_ALIAS").orNull
val releaseKeyPassword = providers.gradleProperty("ECRM_RELEASE_KEY_PASSWORD").orNull
val hasReleaseSigning = listOf(
    releaseStoreFile,
    releaseStorePassword,
    releaseKeyAlias,
    releaseKeyPassword,
).all { !it.isNullOrBlank() }

android {
    namespace = "com.qixi.ecrm.android"
    compileSdk = 36

    defaultConfig {
        applicationId = applicationIdValue
        minSdk = 31
        targetSdk = 36
        versionCode = versionCodeValue
        versionName = versionNameValue

        testInstrumentationRunner = "androidx.test.runner.AndroidJUnitRunner"
    }

    signingConfigs {
        if (hasReleaseSigning) {
            create("release") {
                storeFile = file(requireNotNull(releaseStoreFile))
                storePassword = requireNotNull(releaseStorePassword)
                keyAlias = requireNotNull(releaseKeyAlias)
                keyPassword = requireNotNull(releaseKeyPassword)
            }
        }
    }

    buildTypes {
        debug {
            buildConfigField("String", "API_BASE_URL", "\"$apiBaseUrl\"")
            manifestPlaceholders["UMENG_APPKEY"] = umengAppKey
        }
        release {
            isMinifyEnabled = false
            proguardFiles(
                getDefaultProguardFile("proguard-android-optimize.txt"),
                "proguard-rules.pro",
            )
            buildConfigField("String", "API_BASE_URL", "\"$apiBaseUrl\"")
            manifestPlaceholders["UMENG_APPKEY"] = umengAppKey
            if (hasReleaseSigning) {
                signingConfig = signingConfigs.getByName("release")
            }
        }
    }

    buildFeatures {
        buildConfig = true
        compose = true
    }

    compileOptions {
        sourceCompatibility = JavaVersion.VERSION_17
        targetCompatibility = JavaVersion.VERSION_17
    }
}

tasks.matching { it.name == "assembleRelease" || it.name == "bundleRelease" }.configureEach {
    doFirst {
        check(hasReleaseSigning) {
            "Release 打包必须通过 Gradle 属性提供 ECRM_RELEASE_STORE_FILE、ECRM_RELEASE_STORE_PASSWORD、ECRM_RELEASE_KEY_ALIAS、ECRM_RELEASE_KEY_PASSWORD"
        }
    }
}

kotlin {
    jvmToolchain(17)
}

dependencies {
    implementation(libs.androidx.core.ktx)
    implementation(libs.androidx.lifecycle.runtime.ktx)
    implementation(libs.androidx.lifecycle.runtime.compose)
    implementation(libs.androidx.lifecycle.viewmodel.compose)
    implementation(libs.androidx.activity.compose)
    implementation(libs.androidx.navigation.compose)
    implementation(libs.androidx.hilt.navigation.compose)

    implementation(platform(libs.androidx.compose.bom))
    implementation(libs.androidx.compose.ui)
    implementation(libs.androidx.compose.ui.graphics)
    implementation(libs.androidx.compose.ui.tooling.preview)
    implementation(libs.androidx.compose.material3)
    debugImplementation(libs.androidx.compose.ui.tooling)

    implementation(libs.hilt.android)
    ksp(libs.hilt.compiler)

    implementation(libs.androidx.room.runtime)
    implementation(libs.androidx.room.ktx)
    ksp(libs.androidx.room.compiler)
    implementation(libs.androidx.datastore.preferences)

    implementation(libs.retrofit.core)
    implementation(libs.retrofit.kotlinx.serialization)
    implementation(libs.okhttp.core)
    implementation(libs.okhttp.logging)
    implementation(libs.kotlinx.serialization.json)
    implementation(libs.timber)

    testImplementation(libs.junit)
    androidTestImplementation(libs.androidx.junit)
    androidTestImplementation(libs.androidx.espresso.core)
    androidTestImplementation(platform(libs.androidx.compose.bom))
}
