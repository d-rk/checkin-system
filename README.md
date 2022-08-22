# run postgres in docker:

```
docker run --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres

docker exec -it postgres psql -U postgres -c "CREATE DATABASE checkin;"

docker exec -it postgres psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE checkin TO postgres;"

```
