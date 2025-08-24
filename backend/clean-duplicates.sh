#!/bin/bash

# Remove duplicate entries from config.go
sed -i '/"Irelia".*"MIDDLE",/d' internal/analytics/config.go
sed -i '/"Vayne".*"ADC",/d' internal/analytics/config.go  
sed -i '/"Lux".*"SUPPORT",/d' internal/analytics/config.go

echo "Duplicates cleaned"