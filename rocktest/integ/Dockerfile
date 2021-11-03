FROM alpine
ENV ROCK=/rock

COPY rocktest-go /rocktest-go

RUN apk add unixodbc libc6-compat

VOLUME /scen
WORKDIR /scen

ENTRYPOINT ["/rocktest-go"]