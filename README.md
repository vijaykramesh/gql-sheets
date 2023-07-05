# gql-sheets

## About
This is a small project to power a web-based spreadsheet application.

The backend is in Golang using gqlgen and gorm, and the frontend is in React using Apollo Client.

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