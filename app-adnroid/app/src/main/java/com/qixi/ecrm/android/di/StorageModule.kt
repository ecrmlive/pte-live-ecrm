package com.qixi.ecrm.android.di

import android.content.Context
import androidx.room.Room
import com.qixi.ecrm.android.data.local.AppDatabase
import com.qixi.ecrm.android.data.local.CachedStoreDao
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.android.qualifiers.ApplicationContext
import dagger.hilt.components.SingletonComponent
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
object StorageModule {
    @Provides
    @Singleton
    fun provideDatabase(@ApplicationContext context: Context): AppDatabase =
        Room.databaseBuilder(context, AppDatabase::class.java, "ecrm.db").build()

    @Provides
    fun provideCachedStoreDao(database: AppDatabase): CachedStoreDao = database.cachedStoreDao()
}

