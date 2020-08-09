# go_server_test

This is a set of prebuilt test routines that exercise different parts of the go_server executable. To run the tests, from the directory the different
*.go files are, perform the following:

go build

That will build the executable, "go_server_test", which will be in the directory along with the rest of the *.go files.

To run the executable (first make sure the go_server is up an running):

./go_server_test

If there are any test failures, there will be a line indicating which test failed. The final step of the test is to shut down the go_server.

The following tests are present:
  - Missing password form data for the POST /hash request
  - Missing missing <identifier> for the GET /hash/<identifier> request
  - Invalid <identifier> for the GET /hash/<identifier> request. This is due to the POST /hash not being sent that would have returned the <identifier>
  - Empty Get request. This is when there is no method qualifier that follows the HTTP verb.
  - Empty POST request. This is when there is no method qualifier that follows the HTTP verb.
  - Unexpected PUT request. The go_server only supports GET and POST requests.
  - A valid POST /hash request
  - A GET /hash/1 request that is sent prior to 5 seconds having elapsed from when the POST /hash request that returned 1 was sent.
  - A GET /hash/1 request that is sent after 5 seconds have elapsed. This will return the hashed password.
  - 42 POST /hash requests send in blocks of 7. This means 7 are sent via different "go" func calls at a time. Then the code waits for the 7 to complete and then
      starts on the next block. After all 42 are sent, there is a check to insure that all 42 different <identifier> values have been returned.
  - A GET /hash/42 to verify that the hashed password for <identifier> 42 is returned.
  - A GET /stats request to verify that the statistics can be returned and that there is a "total" of 43. This is the 42 POST /hash calls done in parallel and the one done earlier.
  - A POST /shutdown request to verify that the go_server cleanly shuts down.
 
