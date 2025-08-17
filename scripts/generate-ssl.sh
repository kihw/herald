#!/bin/sh

# Create SSL directory
mkdir -p /etc/nginx/ssl

# Generate self-signed certificate
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    -keyout /etc/nginx/ssl/key.pem \
    -out /etc/nginx/ssl/cert.pem \
    -subj "/C=FR/ST=France/L=Paris/O=Herald/OU=IT/CN=herald.lol"

echo "SSL certificates generated for herald.lol"
