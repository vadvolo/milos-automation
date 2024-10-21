build:
	@git clone -b release https://github.com/netbox-community/netbox-docker.git
	@git clone https://github.com/EvilFreelancer/docker-routeros.git
	@cp ./netbox-docker/docker-compose.yml ./netbox-docker/docker-compose.yml.bak
	@sed -i.bak 's/VERSION-v4.1-3.0.2/VERSION-v3.7/g' ./netbox-docker/docker-compose.yml
	@cp ./docker-compose.yml ./netbox-docker/docker-compose.override.yml
	@cd ./netbox-docker && docker-compose build --no-cache

rebuild:
	cd ./netbox-docker && \
	docker-compose build --no-cache

clean:
	rm -rf docker-routeros/ netbox-docker/

annet_restart:
	@./netbox_env.sh

start_compose:
	cd ./netbox-docker && \
	docker-compose up -d && \
	cd ..

run: start_compose annet_restart

test:
	@source ./netbox_env.sh
	@echo ${NETBOX_IP}
	@echo ${NETBOX_URL}
	
