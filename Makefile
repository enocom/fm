test: clean
	go install
	go generate

clean:
	rm -f spy_test.go
