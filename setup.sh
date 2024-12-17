#!/bin/bash

echo "Deploying..."

SERVICE="monitor-api.service"
DEPLOY_PATH="/opt/monitor_api"

sudo mkdir $DEPLOY_PATH
sudo chown -R $USER:$USER $DEPLOY_PATH
cp * $DEPLOY_PATH
cd $DEPLOY_PATH

echo "Building app..."

go mod download
go build -o counter

echo "Built"
echo "Creating db..."

./counter --createdb

echo "Starting service..."

sudo cp ./$SERVICE /etc/systemd/system/$SERVICE
sudo chmod +x ./monitor-api.sh

sudo systemctl enable $SERVICE
sudo systemctl start $SERVICE
sudo systemctl status $SERVICE

echo "Deployed successfully!"
