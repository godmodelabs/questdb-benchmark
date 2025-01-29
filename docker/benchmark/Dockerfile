FROM golang:1.23 AS build-env
RUN addgroup benchmark && adduser --system --group benchmark

ENV TZ=Europe/Berlin

WORKDIR /app

COPY . /app/

RUN CGO_ENABLED=0 go build -tags timetzdata -a -installsuffix cgo -o /app/benchmark

# main container
FROM scratch

COPY --from=build-env /etc/passwd /etc/passwd
USER benchmark

ENV TZ=Europe/Berlin

WORKDIR /app

COPY --from=build-env /app/benchmark /app/benchmark
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/app/benchmark"]
