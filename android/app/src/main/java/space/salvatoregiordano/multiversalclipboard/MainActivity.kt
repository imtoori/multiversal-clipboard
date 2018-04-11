package space.salvatoregiordano.multiversalclipboard

import android.content.Context
import android.content.Intent
import android.content.SharedPreferences
import android.os.Bundle
import android.support.v7.app.AppCompatActivity
import android.view.View
import androidx.core.content.edit
import com.google.firebase.database.FirebaseDatabase
import kotlinx.android.synthetic.main.activity_main.*

class MainActivity : AppCompatActivity() {

    private var userId: String? = null
    private lateinit var sharedPreferences: SharedPreferences

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        sharedPreferences = getSharedPreferences(getString(R.string.key_shared_preferences), Context.MODE_PRIVATE)

        userId = sharedPreferences.getString(getString(R.string.key_user_id), "")

        if (userId != "") {
            registerButton.visibility = View.GONE

            textView.text = userId

            val intent = Intent(this, ClipboardService::class.java)
            startService(intent)
        } else {
            registerButton.setOnClickListener {
                val newRef = FirebaseDatabase.getInstance().reference.push()
                newRef.setValue("")

                sharedPreferences.edit {
                    putString(this@MainActivity.getString(R.string.key_user_id), newRef.key)
                }

                textView.text = newRef.key

                registerButton.visibility = View.GONE

                val intent = Intent(this, ClipboardService::class.java)
                startService(intent)
            }
        }
    }


}
