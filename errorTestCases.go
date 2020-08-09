package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const Error404Response = "{\"error\": 404}"
const Error412Response = "{\"error\": 412}"
const Error422Response = "{\"error\": 422}"

/*
** Run all the different error checking functions
 */
func runErrorTestCases() bool {
	if !testMissingPassword() {
		return false
	}

	if !testMissingIdentifier() {
		return false
	}

	if !testInvalidIdentifier() {
		return false
	}

	return true
}

/*
** This tests that a PRECONDITION_FAILED_412 is returned if the password form data is missing for the POST /hash
 */
func testMissingPassword() bool {
	/* The following is the form data the POST request expects
	formData := url.Values{
		"password": {"angryMonkey"},
	}
	*/

	log.Println("Starting POST /hash testMissingPassword")
	resp, err := http.PostForm("http://localhost:8080/hash", nil)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))

	if strings.Contains(string(body), Error412Response) {
		log.Printf("Missing password form data test passed\n\n")
		return true
	} else {
		log.Printf("POST /hash missing password form data test failed, expected: %s\n", Error412Response)
	}

	return false
}

/*
** This tests the GET /hash that is missing the <identifier>
** This request will return UNPROCESSABLE_ENTITY_422
 */
func testMissingIdentifier() bool {
	log.Println("Starting GET /hash testMissingIdentifier")
	resp, err := http.Get("http://localhost:8080/hash")
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))

	if strings.Contains(string(body), Error422Response) {
		log.Printf("GET /hash missing <identifier> test passed\n\n")
		return true
	} else {
		log.Printf("GET /hash missing <identifier> test failed, expected: %s\n", Error422Response)
	}

	return false
}

/*
** This tests the GET /hash/<identifier> that has an <identifier> that has not had the corresponding POST
**   performed. The expectation is that the GET /hash/<identifier> will only succeed if the POST /hash
**   has returned the <identifier> and five seconds have elapsed since the POST /hash command returned the
**   <identifier>.
** This request will return NOT_FOUND_404
 */
func testInvalidIdentifier() bool {
	log.Println("Starting GET /hash testInvalidIdentifier")
	resp, err := http.Get("http://localhost:8080/hash/2")
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
		log.Printf("GET /hash invalid <identifier> test passed\n\n")
		return true
	} else {
		log.Printf("GET /hash invalid <identifier> test failed, expected: %s\n", Error404Response)
	}

	return false
}
