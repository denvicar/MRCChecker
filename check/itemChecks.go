package check

type MangaInput struct {
	Id   int
	Name string
}

var CheckForItem = map[int]func(MangaInput, []string) bool{}
