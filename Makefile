#still figuring this out --dakota
#can these be gotten from the .env? idk yet
PGUSER=postgres
PGPASS=12345
PGDATABASE=choicemovers
################################################
PGURL=postgresql://$(PGUSER):$(PGPASS)@localhost:5432/$(PGDATABASE)?sslmode=disable


backend:
	docker compose up --build -d

api:
	docker compose up --build backend -d

db:
	docker compose up --build db -d

createdb:
	docker exec -it db createdb --username=$(PGUSER) --owner=$(PGUSER) $(PGDATABASE)

dropdb:
	docker exec -it db dropdb -U $(PGUSER) $(PGDATABASE)

cleandb:
	docker exec -it db dropdb -U $(PGUSER) $(PGDATABASE)
	docker exec -it db createdb --username=$(PGUSER) --owner=$(PGUSER) $(PGDATABASE)
	migrate -path database/migration -database "$(PGURL)" -verbose up
	./insert_test_data.sh

migrateup:
	migrate -path database/migration -database "$(PGURL)" -verbose up

migrateup1:
	migrate -path database/migration -database "$(PGURL)" -verbose up 1

migratedown:
	migrate -path database/migration -database "$(PGURL)" -verbose down

migratedown1:
	migrate -path database/migration -database "$(PGURL)" -verbose down 1

newmigration:
	migrate create -ext sql -dir database/migration -seq $(name)

