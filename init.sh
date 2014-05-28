#!/bin/bash
$HOME/neo4j/bin/neo4j start
/usr/sbin/nginx -c /etc/nginx/nginx.conf
/usr/local/go/bin/go run /usr/share/nginx/www/api.go 