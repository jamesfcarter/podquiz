package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jamesfcarter/podquiz/internal/assets"
	"github.com/jamesfcarter/podquiz/internal/done"
	"github.com/jamesfcarter/podquiz/internal/server"
	"github.com/jamesfcarter/podquiz/quiz"
)

func main() {
	dir := os.Getenv("PQ_DIR")
	pictures := os.Getenv("PQ_PICTURES")
	if dir == "" {
		dir = "./internal/testdata"
		pictures = "./internal/testdata"
	}
	if pictures == "" {
		pictures = dir + "/../pictures"
	}

	endpoint := os.Getenv("PQ_ENDPOINT")
	if endpoint == "" {
		endpoint = ":8080"
	}

	merchUrl := os.Getenv("PQ_MERCHURL")
	if merchUrl == "" {
		merchUrl = "http://google.com"
	}

	doneFile := os.Getenv("PQ_DONEFILE")
	if doneFile == "" {
		doneFile = "done.txt"
	}

	templates, err := assets.MakeTemplates()
	if err != nil {
		log.Fatal(err)
	}
	done, err := done.New(doneFile)
	if err != nil {
		log.Fatal(err)
	}
	server := &server.Server{
		Database:   quiz.NewDatabase(dir),
		Template:   templates,
		PictureDir: pictures,
		MerchUrl:   merchUrl,
		Done:       done,
	}

	log.Printf("Database directory: %s\n", dir)
	log.Printf("Pictures directory: %s\n", pictures)
	log.Printf("Merch URL: %s\n", merchUrl)
	log.Printf("Starting service on %s\n", endpoint)

	server.Database.Update()

	app, err := server.App()
	if err != nil {
		log.Fatal(err)
	}
	err = http.ListenAndServe(endpoint, app)
	if err != nil {
		log.Fatal(err)
	}
}
