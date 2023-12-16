package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"golang.org/x/net/html"
)

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

// searchCity searches the specified city for listings.
//
// The function takes a city URL, a number of pages to search, a counter, a
// listing channel, and a done channel as arguments. The counter is used to
// track the number of listings that are currently being scraped. The listing
// channel is used to send listings to the main function. The done channel is
// used to signal that a city has been scraped.
func searchCity(cityUrl string, pagesToSearch int, counter *Counter, listings chan Listing, done chan bool) {

	doneScrapingPage := make(chan bool)
	for i := 0; i <= pagesToSearch; i++ {
		go scrapePage(cityUrl+"/search/cta#search=1~gallery~"+string(i)+"~0", counter, listings, doneScrapingPage)
	}

	for i := 0; i <= pagesToSearch; i++ {
		<-doneScrapingPage
		fmt.Println("Finished page:\t", i)
	}
	done <- true

}

// getCities gets the cities from the specified URL.
//
// The function takes a URL as an argument. It returns a slice of strings,
// where each string is the URL of a city.
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

// storeValues stores the listings in the specified database.
//
// The function takes a listing channel, a done channel, and a database URL
// as arguments. The listing channel is used to receive listings from the
// scraping functions. The done channel is used to signal that all listings
// have been stored. The database URL is the URL of the database to store the
// listings in.
func storeValues(lc chan Listing, done chan bool, pbUrl string) {
	var wg sync.WaitGroup
	for {
		val, ok := <-lc

		if !ok {
			break
		}
		count := 0
		go func(l Listing) {
			wg.Add(1)
			postData(pbUrl, l)

			// Move the cursor to the beginning of the line
			fmt.Print("\r")

			// Print the current count without a newline
			fmt.Printf("Total Listings Published To Database: ", count)

			// Flush the output to ensure it gets printed immediately
			os.Stdout.Sync()

			time.Sleep(500 * time.Millisecond) // Simulate some work

			// Increment the count
			count++

			wg.Done()
		}(val)
	}
	wg.Wait()
	done <- true
}

// main is the entry point for the program.
//
// The function does the following:
//
//  1. Gets the cities from the specified URL.
//  2. Creates a counter to track the number of listings that are currently
//     being scraped.
//  3. Creates a listing channel to send listings to the main function.
//  4. Creates a done channel to signal that all listings have been scraped.
//  5. Creates a goroutine to store the listings in the database.
//  6. Creates a goroutine to search each city for listings.
//  7. Waits for all goroutines to finish.
//  8. Prints the total number of listings that were scraped.
func main() {

	startCityIndex := 0
	endCityIndex := 5
	pagesToScan := 3

	pbUrl := "http://127.0.0.1:8090"
	cities := getCities("https://geo.craigslist.org/iso/us")[0:][startCityIndex:endCityIndex]
	counter := Counter{0, 0}
	startTime := time.Now()
	lc := make(chan Listing)
	done := make(chan bool)
	go storeValues(lc, done, pbUrl)

	cityStatus := make(chan bool)
	for _, city := range cities {
		fmt.Println("Searching:\t" + city)
		go searchCity(city, pagesToScan, &counter, lc, cityStatus)
	}

	for cityI := range cities {
		<-cityStatus
		fmt.Println("Finished city:\t", cityI)
	}
	close(lc)
	fmt.Println("Scraping Complete")
	<-done
	fmt.Println("Completed in: ", time.Since(startTime), " seconds!")
	fmt.Println("Print Complete\t#", counter.Complete, "elements printed")
}
