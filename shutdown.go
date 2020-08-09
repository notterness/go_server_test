package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const Ok200Response = "{\"response\": 200}"

/*
** Send the shutdown request to the server. This should return the OK_200 response.
 */
func shutdownServer() {
	/*
	** Wait five second prior to shutting the server down
	 */
	time.Sleep(5000 * time.Millisecond)

	log.Println("Starting shutdownServer")
	resp, err := http.PostForm("http://localhost:8080/shutdown", nil)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))

	if strings.Contains(string(body), Ok200Response) {
		log.Println("Shutdown passed")
	}
}
