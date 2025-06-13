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
