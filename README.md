# WCH OtD API
This repository contains the code which is shared by all WCH OtD app netlify
functions. The exposed functions are in [`fetch.go`](/fetch.go).

## Testing
Tests can be run locally by first starting a redis container in docker with
`docker run --publish 6379:6379 --detach redis:alpine` before running
`go test`.
