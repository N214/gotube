package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"

	//"io/ioutil"
	"log"

	//"os"
	"strings"

	// "net"
	"encoding/xml"
	"net/http"
	"net/url"

	"cloud.google.com/go/storage"
	"github.com/gorilla/mux"

	"github.com/slack-go/slack"
)

// Define a struct to represent the XML data
type Author struct {
	Name string `xml:"name"`
	URI  string `xml:"uri"`
}

type Data struct {
	VideoId   string `xml:"id"`
	YtVideoID string `xml:"http://www.youtube.com/xml/schemas/2015 videoId"`
	YtChnID   string `xml:"http://www.youtube.com/xml/schemas/2015 channelId"`
	Title     string `xml:"title"`
	Link      struct {
		Href string `xml:"href,attr"`
	} `xml:"link"`
	Author    Author `xml:"author"`
	Published string `xml:"published"`
	Updated   string `xml:"updated"`
}

type Subscribe struct {
	Sub int `json:"sub"`
}

type Feed struct {
	Data Data `xml:"entry"`
}

func webhook(res http.ResponseWriter, req *http.Request) {
	challenge := req.URL.Query().Get("hub.challenge")
	if challenge == "" {
		contentType := req.Header.Get("Content-Type")
		if contentType == "application/json" {
			renew := renewSub()
			fmt.Println("Subscribtion renewed")
			fmt.Println(renew)
			return
		}

		data, err := paseXML(req)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Titles: %s\n", data.Data.Title)
		fmt.Printf("Author: %s\n", data.Data.Author.Name)
		fmt.Printf("URL: %s\n", data.Data.Link.Href)

		vidToSend := checkDataHistory(data.Data.YtVideoID)
		if vidToSend == nil {
			return
		} else {
			fmt.Printf("Pushing %s to Slack\n", data.Data.Link.Href)
			pushtoSlack(data.Data.Link.Href)
			// Add function to update data history
		}
	}
	// Renew subscription if there is a challenge
	res.Write([]byte(challenge))
}

func paseXML(req *http.Request) (*Feed, error) {
	var data Feed
	bytes, _ := io.ReadAll(req.Body)

	err := xml.Unmarshal(bytes, &data)

	if err != nil {
		fmt.Println("Error parsing XML:", err)
		return &data, err
	}
	return &data, nil
}

func checkDataHistory(id string) *string {
	// Check Data History from GCS
	// Download data from GCS

	// Init GCS
	bucket := "buk-youtube-data-dev"
	object := "video-history.log"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	checkErr(err)
	defer client.Close()

	o := client.Bucket(bucket).Object(object)
	rc, err := o.NewReader(ctx)
	//if err != nil {
	//	// Create new object
	//}
	checkErr(err)
	defer rc.Close()

	//f, err := os.OpenFile("log.txt", os.O_RDWR, 0644)
	//checkErr(err)
	//defer f.Close()

	history, err := io.ReadAll(rc)
	checkErr(err)

	//tee := io.TeeReader(rc, f)
	// Flag to indicate whether the string id was found in the file
	var found bool

	// scanner := bufio.NewScanner(history)
	scanner := bufio.NewScanner(strings.NewReader(string(history)))

	for scanner.Scan() {
		line := scanner.Text()
		if line == id {
			//if strings.Contains(line, id) {
			found = true
			break
		}
		// data := fmt.Sprintf("%s\n", id)
		// appendStr := []byte(data)
		// defer f.Close()
		// checkErr(err)
		// f.Write(appendStr)

		// Check if there was an error scanning the file
		err := scanner.Err()
		checkErr(err)
	}
	if !found {
		// Write the string id to the end of the file
		data := fmt.Sprintf("%s\n", id)

		appendStr := []byte(data)
		b3 := append(history, appendStr...)
		//fmt.Println(string(b3))

		//f, err := os.OpenFile(string(history), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		//checkErr(err)
		//defer f.Close()

		//_, err = f.WriteString(data)
		//checkErr(err)
		fmt.Printf("Added %s to database\n", id)

		// Upload to GCS
		// Create a writer to upload the data
		wc := o.NewWriter(ctx)
		if _, err := wc.Write(b3); err != nil {
			log.Fatalf("Failed to write object: %v", err)
		}

		// Close the writer
		if err := wc.Close(); err != nil {
			log.Fatalf("Failed to close writer: %v", err)
		}
		fmt.Println("File uploaded to GCS")

		return &id
	} else {
		fmt.Printf("Found %s in database\n", id)
		fmt.Println("Video already exists")
		return nil
	}

	// Upload new data to GCS
}

func checkErr(err error) {
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func pushtoSlack(data string) {
	secret := getSecret("951210594861", "SLACK_YT_BOT_ACCESS_TOKEN", "2").Data
	api := slack.New(string(secret), slack.OptionDebug(true))
	channel, _, err := api.PostMessage("C04NM86CKMF", slack.MsgOptionText(fmt.Sprintln(data), true))
	if err != nil {
		log.Printf("%s\n", err)
		return
	}
	log.Printf("Message successfully sent to channel %s", channel)
}

func renewSub() string {
	pubsub := "https://pubsubhubbub.appspot.com/"

	data := url.Values{}

	callback := "https://northamerica-northeast1-development-brainfinance.cloudfunctions.net/cf-yt-notification-bot"
	topic := "https://www.youtube.com/xml/feeds/videos.xml?channel_id=UCsBjURrPoezykLs9EqgamOA"
	mode := "subscribe"
	verify := "sync"

	data.Set("hub.verify", verify)
	data.Set("hub.topic", topic)
	data.Set("hub.mode", mode)
	data.Set("hub.callback", callback)
	body := data.Encode()

	r, err := http.NewRequest("POST", pubsub, bytes.NewBufferString(body))
	checkErr(err)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		fmt.Printf("Error : %s", err)
	}
	defer resp.Body.Close()
	return fmt.Sprintln("Response Status:", resp.Status)
}

func main() {
	// router
	r := mux.NewRouter()
	r.HandleFunc("/webhooks", webhook).Methods("POST")
	//r.HandleFunc("/sub", subscribe).Methods("POST")
	port := ":8080"

	// http.HandleFunc("/", webhook)

	if err := http.ListenAndServe(port, r); err != nil {
		log.Panicln("Error starting server: ", err)
	}
}
