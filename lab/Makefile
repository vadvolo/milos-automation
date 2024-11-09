# Check if sudo is installed and if the user is not root
# Use "sudo" if both are true
SUDO := $(shell if [ -x "$$(command -v sudo)" ] && [ "$$(id -u)" -ne 0 ]; then echo "sudo"; fi)

# Build `netbox` and `annet` docker images
build:
	@git clone -b release https://github.com/netbox-community/netbox-docker.git
	@cp ./netbox-docker/docker-compose.yml ./netbox-docker/docker-compose.yml.bak
	@sed -i.bak 's/VERSION-v4.1-3.0.2/VERSION-v3.7/g' ./netbox-docker/docker-compose.yml
	@cp ./docker-compose.override.yml ./netbox-docker/docker-compose.override.yml
	@cd ./netbox-docker && docker-compose build --no-cache

rebuild:
	$(SUDO) rm -r netbox-docker || true && \
	$(MAKE) build

# Create `netbox` and `annet` containers
services_start:
	cd ./netbox-docker && \
	docker-compose up -d --build

services_stop:
	cd ./netbox-docker && \
	docker-compose down -v

services_check:
	@if ! docker ps --format '{{.Names}}' | grep -q '^netbox$$' ; then \
		echo "Error: netbox container is not running" >&2; \
		exit 1; \
	fi
	@if ! docker ps --format '{{.Names}}' | grep -q '^annet$$' ; then \
		echo "Error: annet container is not running" >&2; \
		exit 1; \
	fi

lab00_start: services_check annet_flush_generators netbox_flush_db create_netbox_superuser
	@cd ./topologies/lab00_cisco_base && \
	docker cp ./src/lab_generators/. annet:/lab_generators/ && \
	docker cp ./src/netbox-dump.json netbox:/home/ubuntu/netbox-dump.json && \
	docker exec --user root netbox /opt/netbox/netbox/manage.py loaddata /home/ubuntu/netbox-dump.json && \
	docker-compose up -d --build

lab00_stop:
	@cd ./topologies/lab00_cisco_base && \
	docker-compose down

lab01_start: services_check annet_flush_generators netbox_flush_db create_netbox_superuser
	@cd ./topologies/lab01_frr-only-test && \
	$(SUDO) find operational_configs -mindepth 1 -not -name '.gitkeep' -delete || true && \
	cp -r default_configs/* operational_configs/ && \
	docker cp ./src/lab_generators/. annet:/lab_generators/ && \
	docker cp ./src/netbox-dump.json netbox:/home/ubuntu/netbox-dump.json && \
	docker exec --user root netbox /opt/netbox/netbox/manage.py loaddata /home/ubuntu/netbox-dump.json && \
	docker-compose up -d --build

lab01_stop:
	@cd ./topologies/lab01_frr-only-test && \
	docker-compose down

lab02_start: services_check annet_flush_generators netbox_flush_db create_netbox_superuser
	@cd ./topologies/lab02_ceos-test && \
	$(SUDO) find operational_configs -mindepth 1 -not -name '.gitkeep' -delete || true && \
	cp -r default_configs/* operational_configs/ && \
	docker cp ./src/lab_generators/. annet:/lab_generators/ && \
	docker cp ./src/netbox-dump.json netbox:/home/ubuntu/netbox-dump.json && \
	docker exec --user root netbox /opt/netbox/netbox/manage.py loaddata /home/ubuntu/netbox-dump.json && \
	docker-compose up -d --build

lab02_stop:
	@cd ./topologies/lab02_ceos-test && \
	docker-compose down

lab03_start: services_check annet_flush_generators netbox_flush_db create_netbox_superuser
	@cd ./topologies/lab03_dc_cisco && \
	docker cp ./src/lab_generators/. annet:/lab_generators/ && \
	docker cp ./src/netbox-dump.json netbox:/home/ubuntu/netbox-dump.json && \
	docker exec --user root netbox /opt/netbox/netbox/manage.py loaddata /home/ubuntu/netbox-dump.json && \
	docker-compose up -d --build

lab03_stop:
	@cd ./topologies/lab03_dc_cisco && \
	docker-compose down

clean:
	rm -rf netbox-docker/
	docker rm -f netbox-docker-netbox-worker-1 netbox-docker-netbox-housekeeping-1 \
	netbox annet netbox-docker-redis-cache-1 \
	dynamips netbox-docker-redis-1 netbox-docker-postgres-1
	# docker volume prune -f -a

annet_restart:
	source ./netbox_env.sh && \
	cd ./netbox-docker && \
	docker-compose up -d

create_netbox_superuser:
	@docker exec --user root netbox /opt/netbox/netbox/manage.py createsuperuser --noinput

netbox_flush_db:
	@docker exec --user root netbox /opt/netbox/netbox/manage.py flush --noinput

annet_flush_generators:
	@docker exec annet bash -c "find lab_generators -mindepth 1 -delete || true"

build_discovery:
	@cd ./discovery && go mod init discovery && go mod tidy &&  env GOOS=linux GOARCH=arm go build && cd ..

netbox_export:
	@docker exec --user root netbox /opt/netbox/netbox/manage.py dumpdata -a -o /home/ubuntu/src/netbox-dump.json --exclude extras.cachedvalue --exclude extras.objectchange --exclude extras.report --exclude extras.script --exclude django_rq.queue

netbox_import:
	@docker exec --user root netbox /opt/netbox/netbox/manage.py loaddata /home/ubuntu/src/netbox-dump.json

run: start_compose annet_restart
