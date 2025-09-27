build:
	docker compose up --build

reset-volumes:
	docker compose down
	docker volume rm multiplayer-game-server_postgres_data multiplayer-game-server_valkey_data
