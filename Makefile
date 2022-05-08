PACKAGES=`go list ./... | grep -v example`

test:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ${PACKAGES}

format:
	go fmt github.com/fmelosilva/strmap/...

.PHONEY: test