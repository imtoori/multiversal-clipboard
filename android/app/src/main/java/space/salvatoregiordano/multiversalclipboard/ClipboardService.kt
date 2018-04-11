package space.salvatoregiordano.multiversalclipboard

import android.app.Notification
import android.app.Service
import android.content.*
import android.os.IBinder
import android.support.v4.app.NotificationCompat
import android.util.Log
import androidx.core.content.systemService
import com.google.firebase.database.*

private val TAG = "CLIPBOARD_SERVICE"

class ClipboardService : Service() {
    override fun onBind(p0: Intent?): IBinder? = null

    lateinit var clipboardManager: ClipboardManager
    lateinit var userId: String
    lateinit var sharedPreferences: SharedPreferences
    lateinit var valueEventListener: CValueEventListener
    lateinit var userRef: DatabaseReference

    private val CHANNEL_ID = "CLIPBOARD_CHANNEL"

    private val SERVICE_ID = 9999

    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {

        startForeground(SERVICE_ID, createNotification())
        Log.d(TAG, "service started")
        return START_STICKY
    }

    override fun onCreate() {
        super.onCreate()

        Log.d(TAG, "service created")

        sharedPreferences = getSharedPreferences(getString(R.string.key_shared_preferences), Context.MODE_PRIVATE)

        userId = sharedPreferences.getString(getString(R.string.key_user_id), "")

        clipboardManager = systemService<ClipboardManager>()

        clipboardManager.addPrimaryClipChangedListener {
            userRef.setValue(clipboardManager.primaryClip.getItemAt(0).text.toString())
        }

        initDatabase()
    }

    private fun createNotification(): Notification =
            NotificationCompat.Builder(this, CHANNEL_ID)
                    .setContentTitle(TAG)
                    .setContentText(clipboardManager.primaryClip.toString())
                    .build()


    private fun initDatabase() {
        val database = FirebaseDatabase.getInstance()
        userRef = database.getReference(userId)

        valueEventListener = CValueEventListener(clipboardManager, this)
        userRef.addValueEventListener(valueEventListener)
        Log.d(TAG, "database init complete. Ref = ${userRef.key}")
    }

    fun changeUserId(userId: String) {
        this.userId = userId

        userRef.removeEventListener(valueEventListener)
        userRef = FirebaseDatabase.getInstance().getReference(userId)
        userRef.addValueEventListener(valueEventListener)
    }
}

class CValueEventListener(private val clipboardManager: ClipboardManager, private val context: Context) : ValueEventListener {
    override fun onCancelled(p0: DatabaseError?) {
        Log.e(TAG, "error reading from database")
    }

    override fun onDataChange(p0: DataSnapshot?) {
        val clip = p0?.value as? String

        Log.d(TAG, "snapshot changed $clip")

        clip?.let {
            clipboardManager.primaryClip = ClipData.newPlainText(context.getString(R.string.key_clipboard), it)
        }
    }
}