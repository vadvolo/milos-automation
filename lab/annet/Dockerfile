FROM ubuntu:24.04

# Install dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    python3-venv python3 git golang-go vim net-tools iproute2 inetutils-ping telnet openssh-client

# Setup Python environment
RUN python3 -m venv /venv
RUN /venv/bin/pip install annet[netbox]==v0.16.15
RUN /venv/bin/pip install gnetcli_adapter==1.0.9
RUN /venv/bin/pip install gnetclisdk==1.0.43
RUN /venv/bin/pip install icecream==2.1.3

# Set environment variables
ENV PATH=/venv/bin:$PATH
ENV ANN_CONTEXT_CONFIG_PATH=/config.yaml

# Setup Go environment
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH
RUN go install github.com/annetutil/gnetcli/cmd/gnetcli_server@latest

# Create directories
RUN mkdir -p /lab_generators
RUN mkdir -p /build/netsshsetup

# Build netsshsetup
COPY ./netsshsetup/* /build/netsshsetup/
WORKDIR /build/netsshsetup
RUN go build -o /usr/local/bin/netsshsetup

# Copy annet config
COPY config.yaml /config.yaml

# Set final workdir and default command
WORKDIR "/"
CMD ["tail", "-f", "/dev/null"]
