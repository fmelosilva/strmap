PACKAGES=`go list ./... | grep -v example`

test:
	go test -v -cover ${PACKAGES}

format:
	go fmt github.com/fmelosilva/strmap/...

.PHONEY: test