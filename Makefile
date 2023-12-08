app.docker.stop:
	docker compose --file ./docker-compose.yml  down --remove-orphans

app.docker.start:
	docker compose --file ./docker-compose.yml  up -d

run.dev.watch: app.docker.start 
		air -c .air.toml
