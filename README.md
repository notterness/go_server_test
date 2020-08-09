# go_server_test

This is a set of prebuilt test routines that exercise different parts of the go_server executable. To run the tests, from the directory the different
*.go files are, perform the following:

go build

That will build the executable, "go_server_test", which will be in the directory along with the rest of the *.go files.

To run the executable (first make sure the go_server is up an running):

./go_server_test

If there are any test failures, there will be a line indicating which test failed. The final step of the test is to shut down the go_server.
