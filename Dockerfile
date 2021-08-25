FROM golang:1.16 AS build
RUN mkdir /build
COPY . /build/
WORKDIR /build
RUN go build cmd/server/main.go

FROM debian

ENV OPENWEATHERMAP_KEY=4e419a6f5ff541f749fe6c9a1efdd383
ENV WEATHERSTACK_KEY=2b58b9b3c8e1a15ff7423988f682582a
ENV CACHE_TTL=3

RUN useradd app

USER app
COPY --from=build /build/main /home/app/main

WORKDIR /home/app
ENTRYPOINT ./main

EXPOSE 8080
