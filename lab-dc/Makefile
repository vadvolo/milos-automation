build:
	@git clone -b release https://github.com/netbox-community/netbox-docker.git
	@cp ./netbox-docker/docker-compose.yml ./netbox-docker/docker-compose.yml.bak
	@sed -i.bak 's/VERSION-v4.1-3.0.2/VERSION-v3.7/g' ./netbox-docker/docker-compose.yml
	@cp ./docker-compose.yml ./netbox-docker/docker-compose.override.yml
	@cd ./netbox-docker && docker-compose build --no-cache

rebuild:
	cd ./netbox-docker && \
	docker-compose build --no-cache

clean:
	rm -rf netbox-docker/
	docker rm -f netbox-docker-netbox-worker-1 netbox-docker-netbox-housekeeping-1 \
	netbox-docker-netbox-1 netbox-docker-annet-1 netbox-docker-redis-cache-1 \
	netbox-docker-dynamips-lab-1 netbox-docker-redis-1 netbox-docker-postgres-1
	docker volume prune -f -a

annet_restart:
	cd ./netbox-docker && \
	NETBOX_URL="http://netbox-docker-netbox-1:8080" NETBOX_PORT="8080" docker-compose up -d

netbox_db_load:
	docker exec -ti netbox-docker-postgres-1 psql -U netbox -f /netbox_data/drop_netbox.psql
	docker exec -ti netbox-docker-postgres-1 pg_restore -v --no-owner --host=localhost --username=netbox --dbname=netbox /netbox_data/netbox.dump

start_compose:
	cd ./netbox-docker && \
	docker-compose up -d && \
	cd ..

create_netbox_root:
	@docker compose exec netbox-docker-netbox-1 /opt/netbox/netbox/manage.py createsuperuser

run: start_compose annet_restart netbox_db_load

# pg_dump -Fc -v --host=localhost --username=netbox --dbname=netbox -f netbox.dump
# psql -U netbox -f /netbox_data/drop_netbox.psql
# pg_restore -v --no-owner --host=localhost --username=netbox --dbname=netbox /netbox_data/netbox.dump