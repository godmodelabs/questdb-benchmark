FROM questdb/questdb

COPY --chmod=755 healthcheck.sh ./healthcheck.sh

ENV TZ=Europe/Berlin
ENV APP_ENV=production

## Additional variables can be found here: https://questdb.io/docs/reference/configuration
ENV QDB_DEBUG=false
ENV QDB_LOG_W_STDOUT_LEVEL=ERROR
ENV QDB_LOG_W_FILE_LEVEL=ERROR
ENV QDB_LOG_W_HTTP_MIN_LEVEL=ERROR
ENV QDB_SHARED_WORKER_COUNT=8
ENV QDB_TELEMETRY_ENABLED=false
ENV QDB_HTTP_BIND_TO=0.0.0.0:9000
ENV QDB_HTTP_MIN_ENABLED=true
ENV QDB_HTTP_MIN_BIND_TO=0.0.0.0:9003

# Disable postgresql
ENV QDB_PG_ENABLED=false

# Disable influxdb tcp server
ENV QBD_LINE_TCP_ENABLED=false

# Disable influxdb udp server
ENV QBD_LINE_UDP_ENABLED=false
