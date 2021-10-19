FROM golang:1.14-buster as agentbuild
WORKDIR /go/src/github.com/VerizonMedia/kubectl-flame
ADD . /go/src/github.com/VerizonMedia/kubectl-flame
RUN go get -d -v ./...
RUN cd agent && go build -o /go/bin/agent

FROM openjdk:8 as asyncprofiler
RUN curl -o async-profiler-2.5-linux-x64.tar.gz -L \
    https://github.com/jvm-profiling-tools/async-profiler/releases/download/v2.5/async-profiler-2.5-linux-x64.tar.gz
RUN tar -xvf async-profiler-2.5-linux-x64.tar.gz && mv async-profiler-2.5-linux-x64 async-profiler

FROM bitnami/minideb:stretch
RUN mkdir -p /app/async-profiler/build
COPY --from=agentbuild /go/bin/agent /app
COPY --from=asyncprofiler /async-profiler/build /app/async-profiler/build
COPY --from=asyncprofiler /async-profiler/profiler.sh /app/async-profiler
CMD [ "/app/agent" ]