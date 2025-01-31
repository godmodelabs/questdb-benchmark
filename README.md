# questdb-benchmark

A simple questdb benchmark - single writer, single destination table.

Used table definition
```
CREATE TABLE benchmark (timestamp TIMESTAMP, tick INT, instrumentID SYMBOL, price DOUBLE, volume INT) TIMESTAMP(timestamp) PARTITION BY DAY WAL;
```

Install docker + docker-compose

# run
```
docker-compose up
```

Benchmark is executed after questb container starts. Benchmark params are currently hard-coded.
40 mil ticks in total, 2.5 mil symbols, 100k updates per second

The benchmark outputs the writerTxn and sequencerTxn delta of wal_tables() query in questdb.
Furthermore it outputs the delta of "number of ticks sent to questdb" - "actual ticks in questdb" (via count on table).
