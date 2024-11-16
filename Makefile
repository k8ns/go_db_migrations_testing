
.PHONY: test tesdb gen

gen:
	go generate -tags integration ./...

test:
	go test -v -count=1 ./...

testdb:
	go test -tags integration -v -count=1 ./...
