# cypherscan-server

## Spin a local postgres docker container
```bash
docker run --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:10.6-alpine
```
After that, connect via:
- host: localhost
- port: 5432
- database: postgres
- user: postgres
- password: postgres

### Create a database for dev
`create database scan`

### Create a database for local test
`create database scan_test`

## Test
`go get github.com/eaburns/Watch`
`watch -t -p ./ go test`
