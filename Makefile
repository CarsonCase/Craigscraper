makefile:

dev: 
	go build && ./Craigscraper

db:
	rm -r db -f && touch db/listings.db