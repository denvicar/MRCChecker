package check

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

const mal_api = "https://api.myanimelist.net/v2/manga/{id}?fields=title,start_date,end_date,mean,rank,popularity,num_list_users,genres,media_type"
const mu_api = "https://api.mangaupdates.com/v1/series/"

var api_key = os.Getenv("API_KEY")

var client = new(http.Client)

type mangaInfo struct {
	Title        string
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	Mean         float32
	Rank         int
	Popularity   int
	NumListUsers int `json:"num_list_users"`
	Genres       []mangaGenre
	MediaType    string `json:"media_type"`
}

type mangaGenre struct {
	Id   int
	Name string
}

type muResponse struct {
	Results []muResult
}

type muResult struct {
	Record muRecord
}

type muRecord struct {
	Series_id int
	Title     string
	Genres    []struct{ Genre string }
}

func (m *muRecord) GenresSlice() []string {
	s := make([]string, len(m.Genres))
	for i, genre := range m.Genres {
		s[i] = genre.Genre
	}
	return s
}

func (g *mangaInfo) GenresSlice() []string {
	s := make([]string, len(g.Genres))
	for i, genre := range g.Genres {
		s[i] = genre.Name
	}
	return s
}

func CheckGenres(id int, name string, genres []string) bool {
	mal := getMALInfo(id)
	mu := getMUInfo(name)

	for _, genre := range genres {
		if slices.Contains(mal.GenresSlice(), genre) || slices.Contains(mu.GenresSlice(), genre) {
			return true
		}
	}

	return false
}

func getMALInfo(id int) *mangaInfo {
	m := &mangaInfo{}
	url := strings.Replace(mal_api, "{id}", strconv.Itoa(id), 1)
	data := getMALResponse(url)

	json.Unmarshal(data, m)

	return m
}

func getMUInfo(name string) muRecord {
	var ans muRecord
	postData, _ := json.Marshal(map[string]string{"search": name})

	req, _ := http.NewRequest("POST", mu_api+"search", bytes.NewBuffer(postData))
	req.Header.Add("Content-Type", "application/json")

	res, _ := client.Do(req)

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var data = &muResponse{}
	json.Unmarshal(body, data)

	for i, result := range data.Results {
		if result.Record.Title == name {
			ans = data.Results[i].Record
		}
	}

	return ans

}

func getMALResponse(url string) []byte {
	req, e := http.NewRequest("GET", url, nil)

	if e != nil {
		fmt.Printf("Error while preparing request for endpoing '%s'\n", url)
		return []byte{}
	}

	req.Header.Add("X-MAL-CLIENT-ID", api_key)
	r, e := client.Do(req)

	if e != nil {
		fmt.Printf("Error while making the request for '%s'\n", url)
		return []byte{}
	}

	defer r.Body.Close()

	data, e := io.ReadAll(r.Body)
	if e != nil {
		fmt.Printf("Error while reading the server response\n")
		return []byte{}
	}

	return data
}
