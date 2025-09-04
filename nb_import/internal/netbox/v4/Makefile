build:
	@git clone -b release https://github.com/netbox-community/netbox-docker.git
	@cp ./netbox-docker/docker-compose.yml ./netbox-docker/docker-compose.yml.bak
	@cd ./netbox-docker && docker compose build --no-cache && cd ..

run:
	@cd netbox-docker/ && docker compose up -d && cd ..