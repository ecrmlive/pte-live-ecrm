package com.qixi.ecrm.android

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import com.qixi.ecrm.android.presentation.navigation.EcrmApp
import com.qixi.ecrm.android.presentation.theme.EcrmTheme
import dagger.hilt.android.AndroidEntryPoint

@AndroidEntryPoint
class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContent {
            EcrmTheme {
                EcrmApp()
            }
        }
    }
}

