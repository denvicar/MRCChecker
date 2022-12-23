package main

import (
	"fmt"
	"log"
	"os"

	"github.com/denvicar/mrcchecker/check"
	"github.com/denvicar/mrcchecker/parse"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("can't load environment variables")
	}

	url := os.Args[1]

	f := parse.Parse(url)

	if check.CheckGenres(f.Challenges[0].Manga.Id, f.Challenges[0].Manga.Name, []string{"Comedy", "Romance"}) {
		fmt.Println("Check passato")
	} else {
		fmt.Println("Check non passato")
	}

}
