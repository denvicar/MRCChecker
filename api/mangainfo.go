package api

type MangaInfo struct {
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
