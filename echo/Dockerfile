FROM golang:1.20.5-alpine3.18 as builder
RUN mkdir /build
ADD *.go go.* /build/
WORKDIR /build
RUN go build .

FROM alpine:3.18
RUN mkdir /app
WORKDIR /app
COPY --from=builder /build/echo .

CMD /app/echo