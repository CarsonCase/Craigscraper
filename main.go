package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
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

type Counter struct {
	InProgress int
	Complete   int
}

func (c *Counter) incrementInProgress() {
	c.InProgress++
}

func (c *Counter) incrementComplete() {
	c.Complete++
}

func scrapePage(url string, counter *Counter, listingChan chan Listing, pageChan chan bool) {
	doc := getHTMLPage(url)
	listing := Listing{}
	findHTML(doc, func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "li" {
			for _, a := range n.Attr {
				if a.Key == "title" {
					counter.incrementInProgress()
					listing = Listing{Title: a.Val}
					break
				}
			}

			findHTML(n, func(n *html.Node) {
				if n.Type == html.ElementNode && n.Data == "div" {
					for _, a := range n.Attr {
						if a.Key == "class" && a.Val == "price" {
							counter.incrementComplete()
							listing.Price = n.FirstChild.Data
							listingChan <- listing
							break
						}
					}
				}
			})
		}
	})
	pageChan <- true
}

func searchCity(cityUrl string, pagesToSearch int, counter *Counter, listings chan Listing) {
	defer close(listings)

	doneScrapingPage := make(chan bool)
	for i := 0; i <= pagesToSearch; i++ {
		go scrapePage(cityUrl+"/search/cta#search=1~gallery~"+string(i)+"~0", counter, listings, doneScrapingPage)
	}

	for i := 0; i <= pagesToSearch; i++ {
		fmt.Println(<-doneScrapingPage)
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

type Admin struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type Response struct {
	Admin Admin  `json:"admin"`
	Token string `json:"token"`
}

type Listing struct {
	Title string `json:"title"`
	Price string `json:"price"`
}

func post(url string, payload []byte, client *http.Client, headerFunc func(req *http.Request)) []byte {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		log.Fatal("Error creating Auth Token request:", err)
		return []byte{}
	}
	headerFunc(req)

	resp, err := client.Do(req)

	if resp.StatusCode != http.StatusOK {
		log.Fatal(resp)
	}

	if err != nil {
		log.Fatal("Error fetching Auth Token")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading auth response:", err)
		return []byte{}
	}
	return body

}

func postData(pbUrl string, listingData Listing) {
	pbAuthUrl := pbUrl + "/api/admins/auth-with-password"
	pbListingUrl := pbUrl + "/api/collections/listing/records"

	// Create an HTTP client
	client := &http.Client{}

	// Get Auth Token
	payload := []byte(`{"identity":"carsonpcase@gmail.com", "password": "7G}!>LDt.6sjhmf"}`)

	result := post(pbAuthUrl, payload, client, func(req *http.Request) {
		req.Header.Set("Content-Type", "application/json")
	})

	var authResponse Response

	err := json.Unmarshal(result, &authResponse)
	if err != nil {
		log.Fatal("Error unmarshaling auth token response ", err)
	}
	auth := authResponse.Token
	payload, err = json.Marshal(listingData)

	if err != nil {
		log.Fatal("Error marshalling JSON:", err)
		return
	}
	post(pbListingUrl, payload, client, func(req *http.Request) {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+auth)
	})

}

var wg sync.WaitGroup

func printValues(lc chan Listing, done chan bool) {
	for {
		val, ok := <-lc

		if !ok {
			fmt.Println("closed")
			break
		}

		fmt.Println(val)
	}
	done <- true
}

func main() {
	pbUrl := "http://127.0.0.1:8090"
	cities := getCities("https://geo.craigslist.org/iso/us")
	counter := Counter{0, 0}
	startTime := time.Now()
	listOfListings := []Listing{}
	lc := make(chan Listing)
	done := make(chan bool)
	go printValues(lc, done)

	for _, city := range cities {
		fmt.Println("Searching:\t" + city)
		go searchCity(city, 10, &counter, lc)

		if len(os.Args) > 1 && os.Args[1] == "--store" {
			for j, listing := range listOfListings {
				wg.Add(1)
				go func(l Listing) {
					defer wg.Done()
					fmt.Println("Storing:\t"+city+"\tElement: ", j, "/", len(listOfListings))
					postData(pbUrl, l)
				}(listing)
			}
		}
	}
	<-done
	// best 268.873824ms
	wg.Wait()
	fmt.Println("Completed in: ", time.Since(startTime), " seconds!")
	fmt.Println("Print Complete\t#", counter.Complete, "elements printed")
}
