# Web Crawler Challenge

This project crawls a given domain, it processes each page concurrently
but only one request is allowed at a time as a precaution to not
produce a DoS on the server.

## Configuration

The domain to be visited is set through the environment variable `CRAWLER_DOMAIN`
in Bash it can be done as

```bash
CRAWLER_DOMAIN=http://domain.com go run ./cmd/main.go
```
if the variable is omitted the program uses a default domain to showcase
how it works.

To turn off the single request at a time feature set the following variable to false
`REQUEST_THROTTLING=false`.

## Commands

The available commands are standard Go project commands, 
here's just a couple of examples:

### Execute
```bash
go run cmd/main.go
```
### Test
```bash
go test -race ./...
```

## Requirements

The project was developed in Go version 1.21.0 but might work on earlier versions.
It was only tested on Linux but there's nothing  Linux specific
that should prevent it from working in other OS's
