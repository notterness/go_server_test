package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const Error404Response = "{\"error\": 404}"
const Error405Response = "{\"error\": 405}"
const Error412Response = "{\"error\": 412}"
const Error422Response = "{\"error\": 422}"

const AllowedHttpVerbs = "{\"Allow\": GET POST}"


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

	if !testEmptyGet() {
		return false
	}

	if !testEmptyPost() {
		return false
	}

	if !testUnexpectedPut() {
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

/*
** This checks the servers handling of an unqualified GET request
**
** This expects a METHOD_NOT_ALLOWED_405 response from the server
 */
func testEmptyGet() bool {
	log.Println("Starting empty GET test")
	resp, err := http.Get("http://localhost:8080/")
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))

	if strings.Contains(string(body), Error405Response) {
		log.Printf("Empty GET test passed\n\n")
		return true
	} else {
		log.Printf("Empty GET test failed, expected: %s\n", Error405Response)
	}

	return false
}

/*
** This checks the servers handling of an unqualified GET request.
**
** This expects a METHOD_NOT_ALLOWED_405 response from the server
 */
func testEmptyPost() bool {
	log.Println("Starting empty POST test")
	resp, err := http.PostForm("http://localhost:8080/", nil)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))

	if strings.Contains(string(body), Error405Response) {
		log.Printf("Empty POST test passed\n\n")
		return true
	} else {
		log.Printf("Empty POST test failed, expected: %s\n", Error405Response)
	}

	return false
}

/*
** This checks the servers handling of an unexpected PUT request.
**
** This expects a METHOD_NOT_ALLOWED_405 response from the server and the list of acceptable HTTP verbs
 */
func testUnexpectedPut() bool {
	passwordData := url.Values{
		"password": {"angryMonkey"},
	}

	log.Println("Starting unexpected PUT test")

	// initialize http client
	client := &http.Client{}

	// marshal User to json
	json, err := json.Marshal(passwordData)
	if err != nil {
		log.Fatalln(err)
	}

	// set the HTTP method, url, and request body
	req, err := http.NewRequest(http.MethodPut, "http://localhost:8080/", bytes.NewBuffer(json))
	if err != nil {
		log.Fatalln(err)
	}

	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))

	if strings.Contains(string(body), Error405Response) && strings.Contains(string(body), AllowedHttpVerbs) {
		log.Printf("Unexpected PUT test passed\n\n")
		return true
	} else {
		log.Printf("Unexpected PUT test failed, expected: %s and %s\n", Error405Response, AllowedHttpVerbs)
	}

	return false
}