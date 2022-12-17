

DB_CONTAINER = postgres12-vf
DB_USERNAME = root
DB_PASSWORD = root
DB_DBNAME = vanillefraise
DB_PORT = 5676


startdb:
	docker run $(DB_CONTAINER)

stopdb:
	docker stop $(DB_CONTAINER)

postgres:
	docker run --name $(DB_CONTAINER) -p $(DB_PORT):$(DB_PORT) -e POSTGRES_USER=$(DB_USERNAME) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -d postgres:12-alpine 


