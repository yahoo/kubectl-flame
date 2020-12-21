FROM golang:1.14-buster as agentbuild
WORKDIR /go/src/github.com/VerizonMedia/kubectl-flame
ADD . /go/src/github.com/VerizonMedia/kubectl-flame
RUN go get -d -v ./...
RUN cd agent && go build -o /go/bin/agent

FROM alpine AS pyspybuild
RUN apk add python3 py3-pip
RUN echo 'manylinux1_compatible = True' > /usr/lib/python3.8/site-packages/_manylinux.py
RUN pip3 install py-spy==0.4.0.dev1

FROM alpine
RUN mkdir /app
COPY --from=agentbuild /go/bin/agent /app/agent
COPY --from=pyspybuild /usr/bin/py-spy /app/py-spy

CMD [ "/app/agent" ]