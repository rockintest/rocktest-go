FROM ubuntu:latest
ENV ROCK=/rock

COPY rocktest-go /rocktest-go

RUN apt update && apt install -y unixodbc

VOLUME /scen
WORKDIR /scen

ENTRYPOINT ["/rocktest-go"]