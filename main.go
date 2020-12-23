package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"link_collector/link"
)

var URL = "https://news.google.com"
var articles map[int]link.Link

func main() {
	getArticles()
	listArticles()
	for {
		switch choice := getInput(); choice {
		case "r", "refresh":
			getArticles()
			listArticles()
		case "q", "quit":
			return
		case "":
			fmt.Println("invalid choice, try again")
		default:
			openArticle(choice)
		}
	}
}

func getArticles() {
	resp, err := http.Get(URL)
	if err != nil {
		log.Fatal("error opening URL:" + err.Error())
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	links, err := link.Links(strings.NewReader(string(body)))

	base, err := url.Parse(URL)
	if err != nil {
		log.Fatal(err)
	}

	ctr := 1
	articles = make(map[int]link.Link)
	for _, l := range links {
		u, err := url.Parse(l.Url)
		if err != nil {
			log.Fatal(err)
		}
		if len(l.Text) > 10 && strings.Contains(l.Url, "/articles/") {
			l.Url = base.ResolveReference(u).String()
			articles[ctr] = l
			ctr++
		}
	}
}

func listArticles() {
	for i := 0; i < len(articles); i++ {
		fmt.Printf("%d: %s\n", i+1, articles[i+1].Text)
	}
}

func openArticle(choice string) {
	if sno, err := strconv.Atoi(choice); err == nil {
		if l, ok := articles[sno]; ok {
			fmt.Printf("opening article %d: %s\n", sno, l.Text)
			cmd := exec.Command("open", "-a", "/Applications/Firefox.app", l.Url)
			if err = cmd.Run(); err != nil {
				log.Fatal("fatal:" + err.Error())
			}
			return
		}
	}
	fmt.Println("unknown option:", choice)
}

func getInput() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("enter choice [<article number> | (q)uit | (r)efresh]: ")
	text, _ := reader.ReadString('\n')
	return strings.Trim(text, "\n \t")
}

