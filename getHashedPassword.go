package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const ExpectedAngryMonkeyHash = "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="

const SingleStatsResponse = "{\"total\": 1, \"average\": "

/*
** This sends the POST /hash command and checks that the response is "1". It will immediately send the GET /hash/1
**   command and will expect the NOT_FOUND_404 response. It will then wait 5 seconds and re-issue the GET /hash/1
**   command and should receive the expected hash.
** The final check is the GET /stats verb to make sure that there is only a single statistic returned.
 */
func testHashedPassword() bool {
	if !testValidPost() {
		return false
	}

	if !testGetTooSoon() {
		return false
	}

	if !testGetHashedPassword() {
		return false
	}

	if !testGetSingleStats() {
		return false
	}

	return true
}

/*
** Send the POST /hash command with the form data that contains the "angryMonkey" password. Since this is the first
**   time the POST /hash is called, the return value must be "1".
 */
func testValidPost() bool {
	formData := url.Values{
		"password": {"angryMonkey"},
	}

	log.Println("Starting POST /hash")
	resp, err := http.PostForm("http://localhost:8080/hash", formData)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))

	if strings.Contains(string(body), "1") {
		log.Printf("POST /hash passed\n\n")
		return true
	} else {
		log.Printf("POST /hash expected \"1\" in the response\n")
	}

	return false
}

/*
** This tests the GET /hash/<identifier> that has an <identifier> that has had the corresponding POST
**   performed, but the request is prior to the 5 seconds it takes to generate the hash.
** The expectation is that the GET /hash/<identifier> will only succeed if the POST /hash
**   has returned the <identifier> and five seconds have elapsed since the POST /hash command returned the
**   <identifier>.
** This request will return NOT_FOUND_404
 */
func testGetTooSoon() bool {
	log.Println("Starting GET /hash/1 testGetTooSoon")
	resp, err := http.Get("http://localhost:8080/hash/1")
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))

	if strings.Contains(string(body), Error404Response) {
		log.Printf("GET /hash/1 getTooSoon test passed\n\n")
		return true
	} else {
		log.Printf("GET /hash/1 getTooSoon test failed, expected: %s\n", Error404Response)
	}

	return false
}

/*
** This test waits 5 seconds prior to sending the GET /hash/1 request and it expects the base64 encoded
**   SHA512 hash of "angryMonkey" to be returned.
 */
func testGetHashedPassword() bool {

	/*
	** Wait five second prior to sending the GET request
	 */
	time.Sleep(5000 * time.Millisecond)

	log.Println("Starting GET /hash/1 testGetHashedPassword")
	resp, err := http.Get("http://localhost:8080/hash/1")
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))

	if strings.Contains(string(body), ExpectedAngryMonkeyHash) {
		log.Printf("GET /hash/1 getHashedPassword test passed\n\n")
		return true
	} else {
		log.Printf("GET /hash/1 getHashedPassword test failed.\n  Expected: %s\n", ExpectedAngryMonkeyHash)
	}

	return false
}

/*
** Send the GET /stats request and make sure that there was only a single request that has had stats generated
**   at this point.
 */
func testGetSingleStats() bool {

	/*
	** Wait five second prior to sending the GET request
	 */
	time.Sleep(5000 * time.Millisecond)

	log.Println("Starting GET /stats testGetSingleStats")
	resp, err := http.Get("http://localhost:8080/stats")
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))

	if strings.Contains(string(body), SingleStatsResponse) {
		log.Printf("GET /stats testGetSingleStats test passed\n\n")
		return true
	} else {
		log.Printf("GET /stats testGetSingleStats test failed.\n  Expected: %s\n", SingleStatsResponse)
	}

	return false
}
