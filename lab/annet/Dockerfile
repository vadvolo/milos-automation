FROM ubuntu:24.04

# Install dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    python3-venv python3 git golang-go vim net-tools iproute2 inetutils-ping telnet

# Setup Python environment
RUN python3 -m venv /venv
RUN /venv/bin/pip install -e git+https://github.com/annetutil/annet/#egg=annet
RUN /venv/bin/pip install -e git+https://github.com/annetutil/gnetcli_adapter/#egg=gnetcli_adapter
RUN /venv/bin/pip install -e git+https://github.com/annetutil/gnetcli#egg=gnetclisdk\&subdirectory=grpc_sdk/python
RUN /venv/bin/pip install annetbox requests
RUN /venv/bin/pip install --upgrade protobuf grpcio

# Set nvironment variables
ENV PATH=/venv/bin:$PATH
ENV ANN_CONTEXT_CONFIG_PATH=/config.yaml
# TODO: move token to annet config
ENV NETBOX_TOKEN=a630dcefcb191982869e7576190e79bfd569d33c
ENV NETBOX_URL=http://netbox:8080

# Setup Go environment
ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH
RUN go install github.com/annetutil/gnetcli/cmd/gnetcli_server@latest

# Create directories
RUN mkdir /lab_generators
RUN mkdir -p /home/ubuntu/scripts

# Build netsshsetup
COPY ../scripts/netsshsetup /build/
RUN cd /build && go build -o /usr/local/bin/netsshsetup

# Copy annet config
COPY config.yaml /config.yaml

WORKDIR "/"

CMD ["tail", "-f", "/dev/null"]
