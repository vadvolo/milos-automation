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
	docker rm -f netbox-docker-netbox-worker-1 netbox-docker-netbox-housekeeping-1 \
	netbox-docker-netbox-1 netbox-docker-annet-1 netbox-docker-redis-cache-1 \
	netbox-docker-dynamips-lab-1 netbox-docker-redis-1 netbox-docker-postgres-1

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

build_discovery:
	@cd ./discovery && go mod init discovery && go mod tidy &&  env GOOS=linux GOARCH=arm go build && cd ..

run: start_compose annet_restart

netbox_export:
	@docker exec --user root netbox-docker-netbox-1 /opt/netbox/netbox/manage.py dumpdata -a -o netbox-dump.json --exclude extras.cachedvalue --exclude extras.objectchange --exclude extras.report --exclude extras.script --exclude django_rq.queue
	
netbox_import:
	@docker exec --user root netbox-docker-netbox-1 /opt/netbox/netbox/manage.py loaddata /home/ubuntu/src/netbox-dump.json

