package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

const PartialPostResponse = "POST /hash HTTP/1.1\n"

/*
** The value for "total" is the number of requestsToSend + 1
 */
const MultipleStatsResponse = "{\"total\": 43, \"average\": "

const requestLoops = 6
const maxRequestsInParallel = 7
const requestsToSend = requestLoops * maxRequestsInParallel

/*
** The "+ 2" is to handle the fact the indexing actually starts at 1 and that 42 requests are being sent
**   after the initial request was already sent.
 */
var postResponseReceived [requestsToSend + 2]bool

var outstandingPosts int32

/*
** This tests sending the POST /hash requests in parallel and validates that all of the expected responses have been
**   received.
** If all of the responses are received, then it checks if it can obtain the hashed password for request 42.
** If the hashed password request succeeds, then it requests the stats
 */
func testParallelPosts() bool {
	/*
	** Since the prior tests for getHashedPassword sent the first POST /hash request, mark the first value as
	**   received
	 */
	postResponseReceived[1] = true

	var loopCount = 0
	for {
		var i = 0
		for {
			atomic.AddInt32(&outstandingPosts, 1)

			go testPost()

			i++
			if i == maxRequestsInParallel {
				break
			}
		}

		/*
		** Wait for the responses to finish
		 */
		for {
			requests := atomic.LoadInt32(&outstandingPosts)
			if requests == 0 {
				break
			} else {
				/*
				** Wait 100 milli-second prior to checking if all the requests have been processed
				 */
				time.Sleep(100 * time.Millisecond)
			}
		}

		/*
		** Check if another pass through the loop is required
		 */
		loopCount++
		if loopCount == requestLoops {
			break
		}
	}

	/*
	** Now verify that all the expected responses were received (start at index 2)
	 */
	var passed = true
	for i := 2; i < requestsToSend+2; i++ {
		if postResponseReceived[i] != true {
			log.Printf("Response for request %d is missing\n")
			passed = false
			break
		}
	}

	if passed {
		/*
		** Perform the GET /hash/42 to obtain the hashed passowrd for <identifier> 42
		 */
		passed = testGetHashedPassword42()
	}

	if passed {
		/*
		** Perform the GET /stats request. This should return {"total": 43 since there were 42 requests sent
		**   in the above loop, plus the one sent as part of the getHashedPassword set of tests.
		 */
		passed = testGetStats()
	}

	return passed
}

/*
** This just sends a POST /hash and then extracts the response <identifier> and uses that to set a bool in
**   a response array.
** This is run using the "go testPost()" call so they can take place in parallel.
 */
func testPost() {
	formData := url.Values{
		"password": {"angryMonkey"},
	}

	resp, err := http.PostForm("http://localhost:8080/hash", formData)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	identifier := extractIdentifier(string(body))
	if identifier != -1 {
		postResponseReceived[identifier] = true
	}

	atomic.AddInt32(&outstandingPosts, -1)
}

/*
** This parses the response to the POST /hash request and pulls out the integer <identifier> that was returned.
**   The <identifier> is used to mark the array of booleans that can be checked after all the requests have been sent
**   to insure that the <identifier> returned is always unique.
 */
func extractIdentifier(response string) int64 {
	pos := strings.LastIndex(response, PartialPostResponse)
	if pos == -1 {
		return -1
	}

	adjustedPos := pos + len(PartialPostResponse)
	if adjustedPos >= len(response) {
		return -1
	}

	intString := strings.TrimRight(response[adjustedPos:len(response)], "\n")

	/*
	** Now convert the intString to an actual integer
	 */
	i, err := strconv.ParseInt(intString, 10, 32)
	if err != nil {
		return -1
	}

	/* DEBUG
	log.Printf("extracted intString: \"%s\" %d ", intString, i)
	*/

	return i
}

/*
** This test waits 5 seconds prior to sending the GET /hash/1 request and it expects the base64 encoded
**   SHA512 hash of "angryMonkey" to be returned.
 */
func testGetHashedPassword42() bool {

	/*
	** Wait five second prior to sending the GET request
	 */
	time.Sleep(5000 * time.Millisecond)

	log.Println("Starting GET /hash/42 testGetHashedPassword")
	resp, err := http.Get("http://localhost:8080/hash/42")
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
		log.Printf("GET /hash/42 getHashedPassword test passed\n\n")
		return true
	} else {
		log.Printf("GET /hash/42 getHashedPassword test failed.\n  Expected: %s\n", ExpectedAngryMonkeyHash)
	}

	return false
}

/*
** Send the GET /stats request and make sure that there was only a single request that has had stats generated
**   at this point.
 */
func testGetStats() bool {

	/*
	** Wait five second prior to sending the GET request
	 */
	time.Sleep(5000 * time.Millisecond)

	log.Println("Starting GET /stats testGetStats")
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

	if strings.Contains(string(body), MultipleStatsResponse) {
		log.Printf("GET /stats testGetStats test passed\n\n")
		return true
	} else {
		log.Printf("GET /stats testGetStats test failed.\n  Expected: %s\n", SingleStatsResponse)
	}

	return false
}
