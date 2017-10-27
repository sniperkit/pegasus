# Run the tests

* Run all tests `go test ./...` without benchmarks
* Run all tests and benchmarks of **HTTP** and **GRPC** `go test tests/stress/general_test.go -bench=. -stress`
* Run all tests and benchmarks of **HTTP**, **GRPC** and **AMQP**. **Requires a AMQP server**
`go test tests/stress/general_test.go -bench=. -stress -amqp`