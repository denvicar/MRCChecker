package parse

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

func Parse(url string) Form {
	msgid := url[strings.IndexByte(url, '#')+1:]
	res, e := http.Get(url)
	if e != nil {
		return Form{}
	}
	defer res.Body.Close()

	doc, _ := html.Parse(res.Body)

	forumPost := elementByID(doc, msgid)
	user := getUser(forumPost)
	msgText := elementByID(forumPost, fmt.Sprintf("%s%s", "message", msgid[3:]))
	level, challenges := parsePost(msgText)

	return Form{User: user, Challenges: challenges, Level: level}
}

func parsePost(post *html.Node) (string, []Challenge) {
	var b bytes.Buffer
	challenges := make([]Challenge, 0, 140)

	visit := func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" && strings.Contains(a.Val, "http://myanimelist.net/manga/") {
					b.WriteString(fmt.Sprintf("#ID#%s ", a.Val[strings.LastIndexByte(a.Val, '/')+1:]))
				}
			}
		}
		if n.Type == html.TextNode {
			b.WriteString(n.Data)
		}
		return false
	}
	forEachNode(post, visit, nil)

	form := b.String()
	form = strings.ReplaceAll(form, "\n#", " #")

	itemRe, _ := regexp.Compile(`\((\d{1,3})\) ([\w\s\-\/\.,:;]+) #ID#(\d+) (.*)`)
	res := itemRe.FindAllStringSubmatch(form, -1)

	for _, item := range res {
		id, _ := strconv.Atoi(item[3])
		itemNumber, _ := strconv.Atoi(item[1])
		m := Manga{Name: item[4], Id: id}
		c := Challenge{Item: itemNumber, Description: item[2], Manga: m}
		challenges = append(challenges, c)
	}

	return "Genius", challenges
}

func getUser(n *html.Node) string {
	for _, a := range n.Attr {
		if a.Key == "data-user" {
			return a.Val
		}
	}
	return ""
}

func forEachNode(n *html.Node, pre, post func(n *html.Node) bool) {
	if pre != nil {
		if pre(n) {
			return
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}

	if post != nil {
		if post(n) {
			return
		}
	}
}

func elementByID(doc *html.Node, id string) *html.Node {
	var node *html.Node
	verify := func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, a := range n.Attr {
				if a.Key == "id" && a.Val == id {
					node = n
					return true
				}
			}
		}
		return false
	}

	forEachNode(doc, verify, nil)

	return node
}
