build:
	@git clone -b release https://github.com/netbox-community/netbox-docker.git
	@git clone https://github.com/EvilFreelancer/docker-routeros.git
	@cp ./docker-compose.yml ./netbox-docker/docker-compose.override.yml
	@docker-compose build --no-cache

rebuild:
	cd ./netbox-docker && \
	docker-compose build --no-cache

annet_restart:
	@source ./netbox_env.sh
	@docker restart netbox-docker-annet-1

start_compose:
	cd ./netbox-docker && \
	docker-compose up -d && \
	cd ..

run: start_compose annet_restart

test:
	@source ./netbox_env.sh
	@echo ${NETBOX_IP}
	@echo ${NETBOX_URL}
	
