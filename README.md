# gql-sheets
![Coverage](https://img.shields.io/badge/Coverage-51.7%25-yellow)


## About

![gql-sheets-demo](https://github.com/vijaykramesh/gql-sheets/assets/556288/4b634eb1-ae53-4eb7-b9d2-06b91251a1b8)

This is a small project to power a web-based spreadsheet application.

The backend is in Golang using gqlgen and gorm, and the frontend is in React using Apollo Client.

## Features
- GraphQL API
- Websocket support for collaborative live updates
- Automatic versioning & revert
- Formula support
- Markdown support
- Prometheus metrics

## Supported Formulas
- `=A1` reference lookup
- `=SUM(A1:A5)` sum of range
- `=AVERAGE(A1:A5)` average of range

## Local Setup

### Backend
Use go 1.20 and `go mod download` to fetch dependencies.

You can use docker-compose to run a local psql instance, first copy `.env.example` to `.env` and fill it out, and then run `docker-compose up`

You must also manually setup the postgres schema for the time being, see `local/postgres/schema.sql` for the DDLs to run

Finally, run `go run ./server.go` to bring up the server.

### Frontend
See fe/README.md for more instructions on how to run the frontend.

tl;dr: `cd fe && yarn && yarn start`

### Prometheus
Run local Prometheus + Grafana via [docker-compose](https://github.com/ninadingole/docker-compose-stacks/tree/master/prometheus-grafana) and configure it to scrape the backend service.

```yaml
- job_name: app
    scrape_interval: 5s
    static_configs:
      - targets: ['host.docker.internal:8081']
```
(Note you should also change the docker Grafana to forward to port 3333 instead of 3000)

Gorm gives us a number of neat metrics out of the box:

<img width="2160" alt="Screen Shot 2023-07-05 at 5 29 12 PM" src="https://github.com/vijaykramesh/gql-sheets/assets/556288/76bc68f0-cebe-4eae-8542-68dd03a561ba">

## Tests
Run the following to run tests & coverage locally (for the golang backend)

```shell
$ go test -v ./... -covermode=count -coverprofile=coverage.out &&  go tool cover -func=coverage.out -o=coverage.out && cat coverage.out
```

## Contributing
Contributions are welcome, please fork and submit a PR.  Please follow the existing code style and conventions.

## License
This project is licensed under the MIT License - see the LICENSE file for details
