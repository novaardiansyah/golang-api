#!/bin/bash
set -e

echo "--> Preparing directories..."
mkdir -p logs
touch logs/golang.log logs/golang-error.log
sudo chown -R www:www logs deploy/resources
sudo chmod -R 775 logs

sudo find deploy/resources -type f -exec chmod 644 {} +
sudo find deploy/resources -type d -exec chmod 755 {} +

echo "--> Securing env files..."
sudo chmod 600 .env .env.local .env.production 2>/dev/null || true

echo "--> Binary permission..."
sudo chown www:www deploy/bin/api
sudo chmod 755 deploy/bin/api

echo "--> Supervisor setup..."
sudo cp ./deploy/supervisor.conf /etc/supervisor/conf.d/golang-api_novadev_myid.conf

sudo supervisorctl reread
sudo supervisorctl update
sudo supervisorctl restart golang-api_novadev_myid

echo "--> Done."
