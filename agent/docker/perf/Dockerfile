FROM golang:1.14-buster as agentbuild
WORKDIR /go/src/github.com/VerizonMedia/kubectl-flame
ADD . /go/src/github.com/VerizonMedia/kubectl-flame
RUN go get -d -v ./...
RUN cd agent && go build -o /go/bin/agent

FROM node:16.6-alpine
RUN apk add --update perf
RUN npm i -g stackvis@0.5.0
RUN npm i -g bunyan@1.8.15

RUN mkdir /app
COPY --from=agentbuild /go/bin/agent /app/agent
RUN mv /usr/bin/perf /app/perf


CMD [ "/app/agent" ]
