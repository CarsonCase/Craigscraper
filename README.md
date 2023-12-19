# Craigscraper
A Craigslist webscrapper written in GoLang with listings stored in a SQLite database
**Don't use this for any commercial purposes. Craigslist will probably sue you. This is just for funzies**

That said
## How to use it
1. Run `make db``
2. open main.go
3. adjust values in main() func
    pagesToScan = How many pages to scan for each craigslist city
    batchSize = how many cities to scrape at 1 time. Too large of a number and you will get rate limited. Adjust based off your rotating proxy
4. Set up (proxy)[https://www.webshare.io/proxy-servers]. Webshare will give you 10 proxies free and authenticate based on your IP, so no need to change code. If you would like to customize your proxy do so in the scraping.go file.
4. `go build main.go`
5. `./Craigscraper` 
6. Also running `make dev` will preform steps 4 and 5 together
7. The SQLlite database will be filled with scraping info, look at the notebook to see examples of using the data