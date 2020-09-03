## This service provides you an API for a message scheduling

## Run tests
`make test`

## Validate dependencies
`make deps`

## Build and run
`make dockerise service-up`

## Manual testing
```
send HTTP GET request to http://127.0.0.1:8080/printMeAt
request parameters: msg(text message), ts(date when we want message to be printed, in rfc3339 format)
example: http://127.0.0.1:8080/printMeAt?msg="helloworld"&ts=2020-09-02T17:57:15+03:00
expected service log: Here is the message, scheduled for 2020-09-02 17:57:15 +0300 EEST:"helloworld"
```
