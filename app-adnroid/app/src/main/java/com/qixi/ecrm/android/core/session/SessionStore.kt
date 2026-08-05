package com.qixi.ecrm.android.core.session

import android.content.Context
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.stringPreferencesKey
import androidx.datastore.preferences.preferencesDataStore
import dagger.hilt.android.qualifiers.ApplicationContext
import javax.inject.Inject
import javax.inject.Singleton
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.flow.map
import kotlinx.coroutines.runBlocking

private val Context.sessionDataStore by preferencesDataStore(name = "ecrm_session")

@Singleton
class SessionStore @Inject constructor(
    @ApplicationContext private val context: Context,
) {
    private val accessTokenKey = stringPreferencesKey("access_token")

    val accessTokenFlow: Flow<String?> = context.sessionDataStore.data.map { it[accessTokenKey] }

    fun accessToken(): String? = runBlocking { accessTokenFlow.first() }

    suspend fun saveAccessToken(token: String) {
        context.sessionDataStore.edit { preferences ->
            preferences[accessTokenKey] = token
        }
    }

    suspend fun clear() {
        context.sessionDataStore.edit { preferences ->
            preferences.remove(accessTokenKey)
        }
    }
}
