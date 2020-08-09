package main

func main() {

	/*
	** Start by testing the error paths through the different handlers
	 */
	passed := runErrorTestCases()

	/*
	** Test that the POST /hash and GET /hash/<identifier> handlers behave as expected
	 */
	if passed {
		passed = testHashedPassword()
	}

	/*
	** Test sending multiple POST /hash requests in parallel
	 */
	if passed {
		passed = testParallelPosts()
	}

	/*
	** Shutdown the server to finish the testing
	 */
	shutdownServer()
}
