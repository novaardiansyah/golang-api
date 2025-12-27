#!/bin/bash
set -e

echo "--> Preparing logs..."
mkdir -p logs
touch logs/golang.log logs/golang-error.log
sudo chown -R www:www logs
sudo chmod -R 775 logs

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
