#!/bin/bash

# For Execute
# sed -i 's/\r$//' deploy/setup.sh && bash deploy/setup.sh

echo "--> Setting default permissions..."
sudo chown -R www:www . 2>/dev/null || true
sudo find . -type d -exec chmod 755 {} \; 2>/dev/null || true
sudo find . -type f -exec chmod 644 {} \; 2>/dev/null || true

echo "--> Preparing directories..."
mkdir -p logs
touch logs/golang.log logs/golang-error.log 2>/dev/null || true
sudo chown -R www:www logs deploy/resources
sudo chmod -R 775 logs

echo "--> Binary permission..."
sudo chown www:www deploy/bin/api
sudo chmod 755 deploy/bin/api

echo "--> Supervisor setup..."
sudo cp ./deploy/supervisor.conf /etc/supervisor/conf.d/golang-api_novadev_myid.conf

sudo supervisorctl reread
sudo supervisorctl update
sudo supervisorctl restart golang-api_novadev_myid

echo "--> Securing env files..."
sudo chmod 600 .env .env.local .env.production artisan .well-known .git Makefile deploy/setup.sh 2>/dev/null

echo "[SUCCESS] This script has been executed successfully."
