# 1. Set up the database


# 2. Start the database docker
docker pull postgres:17-alpine
docker run --name postgres17 -p 55432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17-alpine
docker exec -it postgres17 psql -U root
docker logs postgres17

# 3. Connect to the database to monitor
connect to db using table plus, dbeaver or anything

# 4. install go lang migrate
migrate create -ext sql -dir db/migration -seq init_schema


# 5. set up the makefile
Remove everything and set up makefile to create empty container and db inside the postgres


# 6. Migrate to the simple-bank
migrate -path db/migration -database "postgresql://root:secret@localhost:55432/simple_bank?sslmode=disable" -verbose up


# 7. Set up for the unit test
go get github.com/lib/pq
go get github.com/stretchr/testify/require

# 8 implement DB transaction -> automatic

# 9 prevent deadlock for account update balance
edit the sql script to add
GetAccount `FOR UPDATE` to make sure to release the db after one query done commit
2 tables get connected by one key -> the insert into table 1 can block the select from the table 2

# 10 Create the CI workflows with github actions 
Set up the go
install go migrate

# 11 Implement HTTP API
create API folder and set the API using Gin
use postman
emit empty slice in sql to handle empty request

# load environment variable for more secure
