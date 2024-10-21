build:
	@git clone -b release https://github.com/netbox-community/netbox-docker.git
	@git clone https://github.com/EvilFreelancer/docker-routeros.git
	@cp ./netbox-docker/docker-compose.yml ./netbox-docker/docker-compose.yml.bak
	@sed -i.bak 's/VERSION-v4.1-3.0.2/VERSION-v3.7/g' ./netbox-docker/docker-compose.yml
	@cd ./discovery && env GOOS=linux GOARCH=arm go build && cd ..
	@cp ./docker-compose.yml ./netbox-docker/docker-compose.override.yml
	@cd ./netbox-docker && docker-compose build --no-cache

rebuild:
	cd ./netbox-docker && \
	docker-compose build --no-cache

clean:
	rm -rf docker-routeros/ netbox-docker/

annet_restart:
	source ./netbox_env.sh && \
	cd ./netbox-docker && \
	docker-compose up -d

start_compose:
	cd ./netbox-docker && \
	docker-compose up -d && \
	cd ..

create_netbox_superuser:
	@docker compose exec netbox-docker-netbox-1 /opt/netbox/netbox/manage.py createsuperuser

run: start_compose annet_restart

test:
	@source ./netbox_env.sh
	@echo ${NETBOX_IP}
	@echo ${NETBOX_URL}
	
