FROM golang:1.14-buster as agentbuild
WORKDIR /go/src/github.com/VerizonMedia/kubectl-flame
ADD . /go/src/github.com/VerizonMedia/kubectl-flame
RUN go get -d -v ./...
RUN cd agent && go build -o /go/bin/agent

FROM python:3.8-alpine AS pyspybuild
# only works with pip version 20.2.4
# https://github.com/benfred/py-spy/issues/353
RUN pip install pip==20.2.4
RUN echo 'manylinux1_compatible = True' > /usr/local/lib/python3.8/site-packages/_manylinux.py
RUN pip3 install py-spy==0.4.0.dev1

FROM python:3.8-alpine
RUN mkdir /app
COPY --from=agentbuild /go/bin/agent /app/agent
COPY --from=pyspybuild /usr/local/bin/py-spy /app/py-spy

CMD [ "/app/agent" ]