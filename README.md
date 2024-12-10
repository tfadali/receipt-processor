# receipt-processor

A Golang application that processes receipts submitted in JSON format to assign reward points according to a set of rules.

### Running the application

To run the application, assuming you have Golang installed, clone the repository and then type

```console
go install .
go run .
```

Alternatively, if you have Docker installed, you can type

```console
docker build -t receipts .
docker run -p 8080:8080 receipts
```

### Running the tests.

A suite of end-to-end tests is included to be run against the running application.  To run the tests, assuming Golang is installed, first make sure the application is running and then type:

```console
go test .
```
