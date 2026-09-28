# opg-sirius-supervision-lay-deputy-hub

### Major dependencies

- [Go](https://golang.org/) (>= 1.26.2)
- [docker compose](https://docs.docker.com/compose/install/) (>= 2.26.0)

#### Installing dependencies locally:
(This is only necessary if running without docker)

- `npm install`
- `go mod download`
---

## Local development

The application ran using the following commands. This hosts it on `localhost:1234/supervision/deputies/lay/1`
To enable debugging and hot-reloading of Go files:

`make dev-up`

Hot-reloading is managed independently and should happen seamlessly. Debugging is available on port :2345
