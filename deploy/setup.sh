#!/bin/bash

echo "--> Fix permission..."

# Set owner
sudo chown -R www:www . 2>/dev/null

# Default perms
sudo find . -type d -exec chmod 755 {} \; 2>/dev/null
sudo find . -type f -exec chmod 644 {} \; 2>/dev/null

# Sensitive files
sudo chmod 600 .env .env.local .env.production .well-known .git 2>/dev/null
sudo chmod -R 600 ./deploy/bash 2>/dev/null

echo "--> Creating log files..."
rm -rf logs/golang.log logs/golang-error.log
mkdir -p logs
touch logs/golang.log logs/golang-error.log

sudo chmod -R 775 logs 2>/dev/null

echo "--> Binary File..."
sudo chmod +x deploy/bin/api

echo "--> Supervisor setup..."
cp ./deploy/supervisor.conf /etc/supervisor/conf.d/golang-api_novadev_myid.conf

sudo supervisorctl reread
sudo supervisorctl update

echo "--> Supervisor restart..."
sudo supervisorctl restart golang-api_novadev_myid
