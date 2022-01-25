FROM golang:1.14-buster as agentbuild
WORKDIR /go/src/github.com/VerizonMedia/kubectl-flame
ADD . /go/src/github.com/VerizonMedia/kubectl-flame
RUN go get -d -v ./...
RUN cd agent && go build -o /go/bin/agent

FROM node:16.6-alpine AS perfbuild
RUN apk add --update perf
RUN npm i -g stackvis@0.5.0

FROM node:16.6-alpine
RUN mkdir /app
COPY --from=agentbuild /go/bin/agent /app/agent
COPY --from=perfbuild /usr/bin/perf /app/perf
COPY --from=perfbuild /usr/local/bin/stackvis /app/stackvis

RUN perf
RUN stackvis

CMD [ "/app/agent" ]
