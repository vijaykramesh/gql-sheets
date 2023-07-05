# gql-sheets
![Coverage](https://img.shields.io/badge/Coverage-51.4%25-yellow)


## About

<img width="1609" alt="Screen Shot 2023-07-04 at 11 34 55 PM" src="https://github.com/vijaykramesh/gql-sheets/assets/556288/953877a4-d8f1-4d8e-b9d0-9bb1209902de">

This is a small project to power a web-based spreadsheet application.

The backend is in Golang using gqlgen and gorm, and the frontend is in React using Apollo Client.

## Features
- GraphQL API
- Websocket support for live updates - TODO
- Formula support
- Markdown support

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

## Tests
Run `go test ./... -cover` to run all tests.

## Contributing
Contributions are welcome, please fork and submit a PR.  Please follow the existing code style and conventions.

## License
This project is licensed under the MIT License - see the LICENSE file for details
