# Craigscraper
Don't use this for any commercial purposes. Craigslist will probably sue you. This is just for funzies

That said
## How to use it
1. download [https://pocketbase.io/](Pocketbase) into this directory and run `./pocketbase serve`
2. Use pocketbase to create a "listing" collection with "title" and "price" string fields 
3. open main.go
4. adjust values in main() func
5. `go build main.go`
6. `./main`
7. Pocketbase will be filled with listing elements of cars. Feel free to adjust main.go to fit other listing scraping