FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    python3-venv python3 git golang-go vim net-tools iproute2 ssh
RUN python3 -m venv /venv
RUN /venv/bin/pip install -e git+https://github.com/annetutil/annet/#egg=annet

RUN /venv/bin/pip install -e git+https://github.com/annetutil/gnetcli_adapter/#egg=gnetcli_adapter
RUN /venv/bin/pip install -e git+https://github.com/annetutil/gnetcli#egg=gnetclisdk\&subdirectory=grpc_sdk/python

RUN /venv/bin/pip install annetbox requests
RUN /venv/bin/pip install --upgrade protobuf grpcio
# RUN /venv/bin/pip install protobuf==5.27.3 grpcio>=1.66.1
ENV PATH=/venv/bin:$PATH
ENV ANN_CONTEXT_CONFIG_PATH=/config.yaml
# ENV ANN_GENERATORS_PATH=/my_generators

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH

RUN go install github.com/annetutil/gnetcli/cmd/gnetcli_server@latest

WORKDIR "/"
