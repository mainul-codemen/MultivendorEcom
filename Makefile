postgres:
	 docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root mvdb

dropdb:
	docker exec -it postgres12 dropdb mvdb

.PHONY: postgres createdb dropdb
