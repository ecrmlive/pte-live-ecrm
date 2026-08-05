package com.qixi.ecrm.android

import android.app.Application
import com.qixi.ecrm.android.push.PushInitializer
import dagger.hilt.android.HiltAndroidApp
import timber.log.Timber

@HiltAndroidApp
class EcrmApplication : Application() {
    override fun onCreate() {
        super.onCreate()
        if (BuildConfig.DEBUG) {
            Timber.plant(Timber.DebugTree())
        }
        PushInitializer.initialize(this)
    }
}

