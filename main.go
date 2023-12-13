package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/html"
)

func getHTMLPage(url string) *html.Node {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal("Get Error")
	}

	defer response.Body.Close()

	doc, err := html.Parse(response.Body)
	if err != nil {
		log.Fatal("Read Response Error")
	}
	return doc
}

func findHTML(n *html.Node, foo func(n *html.Node)) {
	foo(n)
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		findHTML(c, foo)
	}
}

func printPage(url string, count *int) {

	doc := getHTMLPage(url)

	findHTML(doc, func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "li" {
			for _, a := range n.Attr {
				if a.Key == "title" {
					*count++
					fmt.Println(a.Val)
					break
				}
			}

			findHTML(n, func(n *html.Node) {
				if n.Type == html.ElementNode && n.Data == "div" {
					for _, a := range n.Attr {
						if a.Key == "class" && a.Val == "price" {
							fmt.Println(n.FirstChild.Data)
							break
						}
					}
				}
			})
		}
	})

}

func searchCity(cityUrl string, count *int) {
	fmt.Println(cityUrl)
	for i := 0; i <= 10; i++ {
		printPage(cityUrl+"/search/cta#search=1~gallery~"+string(i)+"~0", count)
	}
}

func getCities(url string) (cities []string) {
	doc := getHTMLPage(url)

	findHTML(doc, func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {

			for _, a := range n.Attr {
				if a.Key == "href" {
					cities = append(cities, a.Val)
					break
				}
			}
		}
	})
	return cities
}

func main() {
	cities := getCities("https://geo.craigslist.org/iso/us")
	count := 0
	startTime := time.Now()
	for _, city := range cities {
		searchCity(city, &count)
	}
	fmt.Println("Completed in: ", time.Since(startTime), " seconds!")
	fmt.Println("Print Complete\t#", count, "elements printed")
}
