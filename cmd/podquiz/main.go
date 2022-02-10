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

func fromEnv(name, dflt string) string {
	value := os.Getenv(name)
	if value == "" {
		value = dflt
	}
	return value
}

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

	endpoint := fromEnv("PQ_ENDPOINT", ":8080")

	templates, err := assets.MakeTemplates()
	if err != nil {
		log.Fatal(err)
	}
	done, err := done.New(fromEnv("PQ_DONEFILE", "done.txt"))
	if err != nil {
		log.Fatal(err)
	}
	server := &server.Server{
		Database:   quiz.NewDatabase(dir),
		Template:   templates,
		PictureDir: pictures,
		MerchUrl:   fromEnv("PQ_MERCHURL", "http://google.com"),
		DiscordUrl: fromEnv("PQ_DISCORDURL", "http://google.com"),
		Done:       done,
	}

	log.Printf("Database directory: %s\n", dir)
	log.Printf("Pictures directory: %s\n", pictures)
	log.Printf("Merch URL: %s\n", server.MerchUrl)
	log.Printf("Discord URL: %s\n", server.DiscordUrl)
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
