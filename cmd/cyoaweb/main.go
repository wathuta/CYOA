package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	cyoa "github.com/wathuta/choose_your_own_adventure"
)

func main() {
	port := flag.Int("port", 9090, "a port to listen through")
	file := flag.String("file", "gopher.json", "the .JSON file with the story")
	flag.Parse()
	f, err := os.Open(*file)
	if err != nil {
		log.Printf("unable to open file %s", *file)
	}
	story, err := cyoa.JSONStory(f)
	if err != nil {
		log.Println("unable to unmarshal json")
	}
	sh := cyoa.NewHandler(story, cyoa.WithTemplate(nil))

	log.Fatal(http.ListenAndServe("localhost:"+fmt.Sprintf("%d", *port), sh))
}
