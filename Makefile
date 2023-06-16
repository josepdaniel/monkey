run-test:
	@mkdir -p target
	@go run main.go ./things/test.thing target/test.s
	@nasm -fmacho64 target/test.s -o target/test.o
	@ld -static -e _start -o target/test target/test.o
	./target/test

run-invalid:
	@mkdir -p target
	@go run main.go ./things/invalid.thing target/invalid.s
	@nasm -fmacho64 target/invalid.s -o target/invalid.o
	@ld -static -e _start -o target/invalid target/invalid.o
	./target/invalid

test:
	@go test ./...