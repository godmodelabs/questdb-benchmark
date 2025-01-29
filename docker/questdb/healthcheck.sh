#!/bin/bash

perl -e 'use IO::Socket::INET; $socket = IO::Socket::INET->new(PeerAddr=>"localhost:9003") or exit 1; print $socket "GET / HTTP/1.0\r\nHost: localhost:9003 \r\n\r\n"; $response = <$socket>; if ($response =~ /HTTP\/1\.[01] 2\d\d/) { exit 0; } else { exit 1; }'
