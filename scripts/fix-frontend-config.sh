#!/bin/sh

# Script pour corriger la configuration API dans le frontend compilé
echo "Fixing frontend API configuration..."

# Remplacer localhost:8004 par l'URL relative dans tous les fichiers JS
find /usr/share/nginx/html -name "*.js" -exec sed -i 's|http://localhost:8004|""| g' {} \;

# Alternative: remplacer par le domaine complet
# find /usr/share/nginx/html -name "*.js" -exec sed -i 's|http://localhost:8004|https://herald.lol|g' {} \;

echo "Frontend API configuration updated to use relative URLs"
