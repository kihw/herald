#!/bin/sh

# Generate SSL certificates if they don't exist
if [ ! -f /etc/nginx/ssl/cert.pem ]; then
    echo "Generating SSL certificates..."
    mkdir -p /etc/nginx/ssl
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
        -keyout /etc/nginx/ssl/key.pem \
        -out /etc/nginx/ssl/cert.pem \
        -subj "/C=FR/ST=France/L=Paris/O=Herald/OU=IT/CN=herald.lol"
    echo "SSL certificates generated"
fi

# Fix frontend API configuration
echo "Fixing frontend API configuration..."
find /usr/share/nginx/html -name "*.js" -exec sed -i 's|http://localhost:8004||g' {} \;
echo "Frontend API configuration updated"

# Start the Go backend in the background
echo "Starting Go backend..."
/usr/local/bin/server &

# Start nginx in the foreground
echo "Starting nginx..."
nginx -g "daemon off;"
