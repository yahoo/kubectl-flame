FROM openjdk:8 as asyncprofiler
RUN echo deb http://deb.debian.org/debian stretch-backports main > /etc/apt/sources.list.d/stretch-backports.list && \
    apt-get update && \
    apt-get -y install build-essential
RUN git clone https://github.com/edeNFed/async-profiler.git
RUN cd async-profiler && make

FROM bitnami/minideb:stretch
RUN mkdir -p /app/async-profiler/build
ADD agent /app
COPY --from=asyncprofiler /async-profiler/build /app/async-profiler/build
COPY --from=asyncprofiler /async-profiler/profiler.sh /app/async-profiler
CMD [ "/app/agent" ]

