FROM alpine-glibc:latest
ENV ROCK=/rock

COPY rocktest-go /rocktest-go

RUN apk add unixodbc

VOLUME /scen
WORKDIR /scen

ENTRYPOINT ["/rocktest-go"]