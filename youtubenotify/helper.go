package youtubenotify

import (
	"bytes"
	"io"
	"bufio"
	"encoding/xml"
	"net/http"
	"context"
	"strings"
	"net/url"
	"fmt"
	"os"

	"cloud.google.com/go/storage"
	"github.com/slack-go/slack"

)

func (ei *MyEnvInitializer) Initialize() (string, string, string) {
	PUBSUB_URL := os.Getenv("PUBSUB_URL")
	CALLBACK_URL := os.Getenv("CALLBACK_URL")
	YT_TOPIC := os.Getenv("YT_TOPIC")
	return PUBSUB_URL, CALLBACK_URL, YT_TOPIC
}

func (app *application) paseXML(req *http.Request) (*Feed, error) {
	var data Feed
	bytes, _ := io.ReadAll(req.Body)

	err := xml.Unmarshal(bytes, &data)

	if err != nil {
		app.errorLog.Println(err.Error())
		return &data, err
	}
	return &data, nil
}

func (app *application) checkDataHistory(id string) *string {
	// Init GCS
	bucket := "buk-youtube-data-dev"
	object := "video-history.log"
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		app.errorLog.Println(err.Error())
	}
	defer client.Close()

	// Use GCS
	o := client.Bucket(bucket).Object(object)
	rc, err := o.NewReader(ctx)
	app.checkErr(err)
	defer rc.Close()

	history, err := io.ReadAll(rc)
	app.checkErr(err)

	// Use local storage
	//f, err := os.OpenFile("log.txt", os.O_RDWR, 0644)
	//app.checkErr(err)
	//defer f.Close()
	//tee := io.TeeReader(rc, f)

	// Flag to indicate whether the string id was found in the file
	var found bool

	// scanner := bufio.NewScanner(history)
	scanner := bufio.NewScanner(strings.NewReader(string(history)))

	for scanner.Scan() {
		line := scanner.Text()
		if line == id {
			found = true
			break
		}

		// Check if there was an error scanning the file
		err := scanner.Err()
		app.checkErr(err)
	}
	if !found {
		// Write the string id to the end of the file
		data := fmt.Sprintf("%s\n", id)

		appendStr := []byte(data)
		b3 := append(history, appendStr...)
		app.infoLog.Printf("Added %s to database\n", id)

		// Upload to GCS
		wc := o.NewWriter(ctx)
		if _, err := wc.Write(b3); err != nil {
			app.errorLog.Println(err.Error())
		}

		// Close the writer
		if err := wc.Close(); err != nil {
			app.errorLog.Printf("Failed to close writer: %v", err)
		}
		app.infoLog.Println("File uploaded to GCS")

		return &id
	} else {
		app.infoLog.Printf("Found %s in database\n", id)
		app.infoLog.Println("Video already exists")
		return nil
	}
}

func (app *application) checkErr(err error) {
	if err != nil {
		app.infoLog.Fatal(err.Error())
	}
}

func (app *application) pushtoSlack(data string) {
	secret := getSecret("951210594861", "SLACK_YT_BOT_ACCESS_TOKEN", "2").Data
	api := slack.New(string(secret), slack.OptionDebug(true))
	channel, _, err := api.PostMessage("C04NM86CKMF", slack.MsgOptionText(fmt.Sprintln(data), true))
	if err != nil {
		app.errorLog.Printf("%s\n", err)
		return
	}
	app.infoLog.Printf("Message successfully sent to channel %s", channel)
}

func (app *application) renewSub() string {

    var initializer GetEnv = &MyEnvInitializer{}
	pubsub, callback, topic := initializer.Initialize()
	mode := "subscribe"
	verify := "sync"

	data := url.Values{}
	data.Set("hub.verify", verify)
	data.Set("hub.topic", topic)
	data.Set("hub.mode", mode)
	data.Set("hub.callback", callback)
	body := data.Encode()

	r, err := http.NewRequest("POST", pubsub, bytes.NewBufferString(body))
	app.checkErr(err)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		app.errorLog.Printf(err.Error())
	}
	defer resp.Body.Close()
	return fmt.Sprintln("Response Status:", resp.Status)
}