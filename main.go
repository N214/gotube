package main

import (
	"log"
	 _ "github.com/N214/gotube/youtubenotify"
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)


func main() {
	port := "8089"
	//if err := gotube.Run(); err != nil {
	//	log.Fatal(err)
	//}
	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}