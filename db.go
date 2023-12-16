package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

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
	Link  string `json:"link"`
}

// post is a helper function to make HTTP POST requests.
//
// It takes a URL, payload, HTTP client, and header function as arguments.
// The header function is used to set any additional headers on the request.
//
// The function returns the response body as a byte slice.
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

// postData posts data to the specified pocketbase server.
//
// It takes a pocketbase server URL and a listing data object as arguments.
// The function first gets an auth token from the server using the
// `auth-with-password` endpoint. It then uses the auth token to post
// the listing data to the `listing/records` endpoint.
func postData(pbUrl string, listingData Listing, onComplete func()) {
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
	onComplete()

}
