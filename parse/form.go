package parse

import (
	"bytes"
	"fmt"
)

type Form struct {
	User       string
	Challenges []Challenge
	Level      string
}

func (f Form) String() string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("User %s\nLevel %s\n", f.User, f.Level))
	for _, c := range f.Challenges {
		b.WriteString(c.String())
	}
	return b.String()
}

type Challenge struct {
	Item        int
	Description string
	Manga
}

func (c Challenge) String() string {
	return fmt.Sprintf("%d %s - %s\n", c.Item, c.Description, c.Manga)
}

type Manga struct {
	Name string
	Id   int
}

func (m Manga) String() string {
	return fmt.Sprintf("%s (%d)", m.Name, m.Id)
}
