package main

import (
	"fmt"
	"os"

	"github.com/denvicar/mrcchecker/check"
	"github.com/denvicar/mrcchecker/parse"
)

func main() {
	url := os.Args[1]
	form := parse.Parse(url)
	if check.CheckGenres(f.Challenges[0].Manga.Id, f.Challenges[0].Manga.Name, []string{"Comedy", "Romance"}) {
		fmt.Println("Check passato")
	} else {
		fmt.Println("Check non passato")
	}

}
