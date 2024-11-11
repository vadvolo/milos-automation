# Check if sudo is installed and if the user is not root
# Use "sudo" if both are true
SUDO := $(shell if [ -x "$$(command -v sudo)" ] && [ "$$(id -u)" -ne 0 ]; then echo "sudo"; fi)

### Services (annet and netbox) control

build:
	@git clone -b release https://github.com/netbox-community/netbox-docker.git
	@cp ./netbox-docker/docker-compose.yml ./netbox-docker/docker-compose.yml.bak
	@sed -i.bak 's/VERSION-v4.1-3.0.2/VERSION-v3.7/g' ./netbox-docker/docker-compose.yml
	@cp ./docker-compose.override.yml ./netbox-docker/docker-compose.override.yml
	@cd ./netbox-docker && docker-compose build --no-cache && cd ..
	@cd scripts/netsshsetup && env GOOS=linux GOARCH=arm go build && cd ../..

rebuild:
	rm -rf netbox-docker || true && \
	$(MAKE) build

services_start:
	cd ./netbox-docker && \
	docker-compose up -d --build

services_stop:
	cd ./netbox-docker && docker-compose down -v && cd ..

### Services utilities

services_check:
	@if ! docker ps --format '{{.Names}}' | grep -q '^netbox$$' ; then \
		echo "Error: netbox container is not running" >&2; \
		exit 1; \
	fi
	@if ! docker ps --format '{{.Names}}' | grep -q '^annet$$' ; then \
		echo "Error: annet container is not running" >&2; \
		exit 1; \
	fi

create_netbox_superuser:
	@docker exec --user root netbox /opt/netbox/netbox/manage.py createsuperuser --noinput

netbox_flush_db:
	@docker exec --user root netbox /opt/netbox/netbox/manage.py flush --noinput

annet_flush_generators:
	@docker exec annet bash -c "find lab_generators -mindepth 1 -delete || true"

lab_pre_start: services_check annet_flush_generators netbox_flush_db create_netbox_superuser

# Build and destroy architecture

lab00: services_stop
	@echo "Lab00 Cisco Basic Scenario is launching..." && \
	cp ./topologies/lab00_cisco_basic_scenario/docker-compose.yml ./netbox-docker/docker-compose.override.yml && \
	cd netbox-docker/ && docker-compose up -d && cd .. && \
	docker exec --user root netbox /opt/netbox/netbox/manage.py loaddata /home/ubuntu/src/netbox-dump.json && \
	echo "!!! Don't forget add Cisco IOS image into ../vm_images folder !!!"
	eval "docker exec -u root -t -i annet /bin/bash"

lab01: services_stop
	@cp ./topologies/lab01_frr-only-test/docker-compose.yml ./netbox-docker/docker-compose.override.yml && \
	$(SUDO) find operational_configs -mindepth 1 -not -name '.gitkeep' -delete || true && \
	cd topologies/lab01_frr-only-test && \
	cp -r default_configs/* operational_configs/ && \
	cd ../../netbox-docker/ && docker-compose up -d --build && cd .. && \
	docker exec --user root netbox /opt/netbox/netbox/manage.py loaddata /home/ubuntu/src/netbox-dump.json

lab02: services_stop
	cp ./topologies/lab02_dc_cisco/docker-compose.yml ./netbox-docker/docker-compose.override.yml && \
	cd netbox-docker/ && docker-compose up -d && cd .. && \
	docker exec --user root netbox /opt/netbox/netbox/manage.py loaddata /home/ubuntu/src/netbox-dump.json && \
	source ./netbox_env.sh && \
	cd ./netbox-docker && \
	docker-compose up -d && \
	cd .. && \
	echo "!!! Don't forget add iol image into ./lab00_cisco_basic_scenario/src/ios-7200 folder !!!"

lab03: services_stop
	cp ./topologies/lab03_multivendor/docker-compose.yml ./netbox-docker/docker-compose.override.yml && \
	cd netbox-docker/ && docker-compose up -d && cd .. && \
	docker exec --user root netbox /opt/netbox/netbox/manage.py loaddata /home/ubuntu/src/netbox-dump.json

netbox_export:
	@docker exec --user root netbox /opt/netbox/netbox/manage.py dumpdata -a -o /home/ubuntu/src/netbox-dump.json --exclude extras.cachedvalue --exclude extras.objectchange --exclude extras.report --exclude extras.script --exclude django_rq.queue --exclude auth.permission --exclude contenttypes.contenttype

netbox_import:
	@docker exec --user root netbox /opt/netbox/netbox/manage.py loaddata /home/ubuntu/src/netbox-dump.json
