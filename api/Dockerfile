FROM golang:alpine as build

RUN mkdir /src

WORKDIR /src

RUN export GO111MODULE=on

COPY . /src

RUN cd /src

RUN go mod tidy

RUN go build

FROM alpine:latest

WORKDIR /api

COPY --from=build /src/archive-api .
COPY --from=build /src/.env .

EXPOSE 3000

ENTRYPOINT ["./archive-api"]
CMD ["./archive-api"]

