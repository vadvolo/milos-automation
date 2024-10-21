# Network Automation Lab

## Installation

### Preparation

You need to install Docker and Docker Compose softaware on yours device:
- [Docker](https://docs.docker.com/engine/install/)
  - [Linux](https://docs.docker.com/desktop/install/linux/)
  - [Mac](https://docs.docker.com/desktop/install/mac-install/)
  - [Windows](https://docs.docker.com/desktop/install/windows-install/)

Pull lab environment from repository:
```
git clone https://github.com/vadvolo/milos-automation.git
cd milos-automation
```

And switch to the `temp` branch:
```
git checkout lab
```

Navigate to the lab folder:
```
cd lab
```

Build and run lab
```
make build
make run
```

### How to connect to containers

Netbox: `docker exec -u root -t -i netbox-docker-netbox-1 /bin/bash`

Annet: `docker exec -u root -t -i netbox-docker-annet-1 /bin/bash`

Dynamips: `docker exec -u root -t -i netbox-docker-dynamips-1 /bin/bash`

### NETBOX TOKEN

1. Go to [LOCAL NETBOX INSTALLATION](http://localhost:8000)
2. 

## How to use

Go to annet:
```
docker exec -u root -t -i netbox-docker-annet-1 /bin/bash
```

Run:
- diff
```
python3 -m annet.annet diff lab-r1.nh.com
```
- patch
```
python3 -m annet.annet patch lab-r1.nh.com
```
- deploy
```
python3 -m annet.annet deploy lab-r1.nh.com
```


