# FRR

Here is located piece of the lab responces for FRR. At it core is Dockerfile for FRR container.

To add new FRR router you need to:
- add new conf file, for instance `frr-r4.conf`
- add into Dockerfile `COPY frr-r4.conf /frr-r4.conf`
- add into docker-compose file part with `frr-r4`
```
frr-r4:
    image: frr-debian:latest
    build: ../frr
    privileged: true
    environment:
      HOSTNAME: frr-r4
      FRR_DAEMONS: "zebra bgpd ospfd ldpd"
    command: /install.sh
    tty: true
    networks:
      - default
      - ...
```

If conf file is changed it is necessary make rebuild procedure for implementation new configuration.
