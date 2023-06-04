run-test:
	@mkdir -p target
	@go run main.go test.thing target/test.s
	@nasm -fmacho64 target/test.s -o target/test.o
	@ld -static -e _start -o target/test target/test.o
	./target/test