package main

import (
	"context"
	"firebase.google.com/go"
	"firebase.google.com/go/db"
	"fmt"
	"github.com/getlantern/systray"
	"github.com/robfig/cron"
	"google.golang.org/api/option"
	"os"
	"os/exec"
	"golang.org/x/oauth2/google"
	"encoding/json"
)

var (
	app          *firebase.App
	database     *db.Client
	userDbRef    *db.Ref
	userId       string
	lastSentClip string
)

const (
	DATABASEURL = "https://multiversal-clipboard.firebaseio.com/"
)

var CREDENTIALS = struct {
	Type                        string `json:"type"`
	Project_id                  string `json:"project_id"`
	Private_key_id              string `json:"private_key_id"`
	Private_key                 string `json:"private_key"`
	Client_email                string `json:"client_email"`
	Client_id                   string `json:"client_id"`
	Auth_uri                    string `json:"auth_uri"`
	Token_uri                   string `json:"token_uri"`
	Auth_provider_x509_cert_url string `json:"auth_provider_x509_cert_url"`
	Client_x509_cert_url        string `json:"client_x509_cert_url"`
}{
!!!!SET HERE THE CREDENTIALS
}

func main() {
	if userId = os.Args[1]; userId != "" {
		systray.Run(onReady, onExit)
	} else {
		fmt.Println("Missing user reference. Usage ./multiversal-clipboard <user reference>")
	}
}

func onReady() {
	fmt.Println("Starting..")
	systray.SetTitle("Multiversal clipboard")

	clipboardCron := cron.New()

	initFCM()

	clipboardCron.AddFunc("@every 3s", pollingFunction)
	clipboardCron.Start()
}

func pollingFunction() {
	var clipboard string
	if clipboard = readClipboard(); clipboard == "" {
		return
	}

	if clipboard != lastSentClip {
		lastSentClip = clipboard
		setRemoteClipboard()
		return
	}

	if clipboard = getRemoteClipboard(); clipboard == "" {
		return
	}

	if clipboard != lastSentClip {
		lastSentClip = clipboard
		setClipboard(clipboard)
	}

}
func setClipboard(clipboard string) error {
	cmd := fmt.Sprint("echo -n ", clipboard, " | ", "pbcopy")
	print(cmd)
	if o, e := exec.Command("bash", "-c", cmd).Output(); e != nil {
		fmt.Println("Error reading clipboard", e.Error())
		return e
	} else {
		fmt.Println(string(o))
		fmt.Println("Clipboard set to", clipboard)
		return nil
	}
}

func readClipboard() string {
	if output, e := exec.Command("pbpaste", "").Output(); e != nil {
		fmt.Println("Error reading clipboard", e.Error())
		return ""
	} else {
		return string(output)
	}
}

func getRemoteClipboard() string {
	var clipboard string
	if e := userDbRef.Get(context.Background(), &clipboard); e != nil {
		fmt.Println("Error getting clipboard.", e)
		return ""
	} else {
		return clipboard
	}
}

func setRemoteClipboard() {
	fmt.Println("Clipboard:", lastSentClip)
	if err := userDbRef.Set(context.Background(), lastSentClip); err != nil {
		fmt.Println("Error updating clipboard", err)
		return
	}
	fmt.Println("Remote clipboard set to", lastSentClip)
}

func initFCM() {
	fmt.Println("Init fcm..")

	var err error
	var j []byte

	if j, err = json.Marshal(CREDENTIALS); err != nil {
		fmt.Println("Error reading credentials", err)
	}

	credentials, err := google.CredentialsFromJSON(context.Background(), j,"https://www.googleapis.com/auth/firebase","https://www.googleapis.com/auth/userinfo.email","https://www.googleapis.com/auth/firebase.database")

	opt := option.WithCredentials(&google.Credentials{
		TokenSource: credentials.TokenSource,
		JSON:        credentials.JSON,
		ProjectID:   credentials.ProjectID,
	})

	config := &firebase.Config{
		DatabaseURL: DATABASEURL,
	}
	app, err = firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		fmt.Println("error initializing app:", err)
		return
	}

	if database, err = app.Database(context.Background()); err != nil {
		fmt.Println("error initializing db:", err)
		return
	}

	userDbRef = database.NewRef(userId)
}

func onExit() {
	fmt.Println("Exiting..")
}
