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
	opt := option.WithCredentialsFile("firebase-credentials.json")
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
