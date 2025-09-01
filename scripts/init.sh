#!/bin/bash

# VPN API Initialization Script
# This script sets up the necessary directories and permissions

echo "Initializing VPN API..."

# Create necessary directories
sudo mkdir -p /etc/apivpn
sudo mkdir -p /etc/xray
sudo mkdir -p /var/lib/scrz-prem

# Set proper permissions
sudo chmod 755 /etc/apivpn
sudo chmod 600 /etc/apivpn/users.json 2>/dev/null || true

# Create default domain file if it doesn't exist
if [ ! -f /etc/xray/domain ]; then
    echo "Setting up default domain..."
    EXTERNAL_IP=$(curl -s ipinfo.io/ip)
    echo "$EXTERNAL_IP" | sudo tee /etc/xray/domain > /dev/null
    echo "IP=$EXTERNAL_IP" | sudo tee /var/lib/scrz-prem/ipvps.conf > /dev/null
fi

# Install required system packages
echo "Installing required packages..."
sudo apt update
sudo apt install -y curl vnstat bc

# Enable vnstat service
sudo systemctl enable vnstat
sudo systemctl start vnstat

# Build and install the API
echo "Building VPN API..."
go mod tidy
go build -o vpn-api main.go

# Create systemd service
echo "Creating systemd service..."
sudo tee /etc/systemd/system/vpn-api.service > /dev/null <<EOF
[Unit]
Description=VPN Management API
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=$(pwd)
ExecStart=$(pwd)/vpn-api
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

# Enable and start the service
sudo systemctl daemon-reload
sudo systemctl enable vpn-api

echo "VPN API initialized successfully!"
echo "To start the service: sudo systemctl start vpn-api"
echo "To check status: sudo systemctl status vpn-api"
echo "To view logs: sudo journalctl -u vpn-api -f"
echo ""
echo "Default admin credentials will be created on first run."
echo "Check /etc/apivpn/default_credentials.txt after starting the service."