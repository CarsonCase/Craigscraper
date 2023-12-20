# Craigscraper
A Craigslist webscrapper written in GoLang with listings stored in a SQLite database and an AI vector database interpreter.
**Don't use this for any commercial purposes. Craigslist will probably sue you. This is just for funzies**

That said
## How to scrape
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

## How to make natural text queries with Langchain
The ai.ipynb file contains necessary code to use AI to interpret your Craigslist data!

First start a python environment and `pip install -r requirements. txt` to install requirements

(this)[https://docs.datastax.com/en/astra-serverless/docs/vector-search/quickstart.html] tutorial shows how to set up an Astra Cassandra server with all the files needed in .env (shows in .env.example)

Next run `python ai.py` to start publishing the scraped sqlite data to Asta.

Finally use ai_query.ipynb to make some queries!
