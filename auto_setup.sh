#!/bin/bash

# Auto Setup Script untuk VPN API
# Untuk VPS: 128.199.227.169

set -e

echo "🚀 Starting VPN API Auto Setup..."
echo "=================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo -e "${BLUE}$1${NC}"
}

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    print_error "Please run this script as root!"
    exit 1
fi

# Get server IP
SERVER_IP=$(curl -s ipinfo.io/ip || echo "128.199.227.169")
print_status "Detected server IP: $SERVER_IP"

# Step 1: Update system
print_header "📦 Step 1: Updating system packages..."
apt update && apt upgrade -y

# Step 2: Install Go if not exists
print_header "🔧 Step 2: Installing Go..."
if ! command -v go &> /dev/null; then
    print_status "Installing Go 1.21..."
    cd /tmp
    wget -q https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
    rm -rf /usr/local/go
    tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
    
    # Add Go to PATH
    if ! grep -q "/usr/local/go/bin" ~/.bashrc; then
        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    fi
    export PATH=$PATH:/usr/local/go/bin
    
    print_status "Go installed successfully"
else
    print_status "Go is already installed: $(go version)"
fi

# Step 3: Install system dependencies
print_header "📋 Step 3: Installing system dependencies..."
apt install -y curl vnstat bc nginx jq

# Enable and start vnstat
systemctl enable vnstat
systemctl start vnstat
print_status "vnstat service enabled and started"

# Step 4: Setup directories
print_header "📁 Step 4: Creating directories..."
mkdir -p /etc/apivpn
mkdir -p /etc/xray
mkdir -p /var/lib/scrz-prem

# Setup domain configuration
echo "$SERVER_IP" > /etc/xray/domain
echo "IP=$SERVER_IP" > /var/lib/scrz-prem/ipvps.conf
print_status "Domain configuration created with IP: $SERVER_IP"

# Step 5: Build API
print_header "🔨 Step 5: Building VPN API..."
cd /root/apivpn || {
    print_error "Please make sure you're in the correct directory (/root/apivpn)"
    exit 1
}

# Download dependencies and build
export PATH=$PATH:/usr/local/go/bin
go mod tidy
go build -o vpn-api main.go
chmod +x vpn-api

print_status "VPN API built successfully"

# Step 6: Generate JWT secret
JWT_SECRET=$(openssl rand -hex 32)
print_status "Generated JWT secret"

# Step 7: Create systemd service
print_header "⚙️  Step 7: Creating systemd service..."
cat > /etc/systemd/system/vpn-api.service << EOF
[Unit]
Description=VPN Management API
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/root/apivpn
ExecStart=/root/apivpn/vpn-api
Restart=always
RestartSec=5
Environment=PORT=8080
Environment=JWT_SECRET=$JWT_SECRET

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd and start service
systemctl daemon-reload
systemctl enable vpn-api
systemctl start vpn-api

# Wait a moment for service to start
sleep 3

# Check service status
if systemctl is-active --quiet vpn-api; then
    print_status "VPN API service started successfully"
else
    print_error "Failed to start VPN API service"
    print_status "Checking logs..."
    journalctl -u vpn-api --no-pager -l
    exit 1
fi

# Step 8: Setup firewall
print_header "🔥 Step 8: Configuring firewall..."
ufw --force enable
ufw allow 22/tcp
ufw allow 8080/tcp
ufw allow 80/tcp
ufw allow 443/tcp
print_status "Firewall configured"

# Step 9: Setup Nginx reverse proxy
print_header "🌐 Step 9: Setting up Nginx reverse proxy..."
cat > /etc/nginx/sites-available/vpn-api << EOF
server {
    listen 80;
    server_name $SERVER_IP _;

    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }

    location / {
        return 200 'VPN API Server is running! Use /api/ endpoints.';
        add_header Content-Type text/plain;
    }
}
EOF

# Enable nginx site
ln -sf /etc/nginx/sites-available/vpn-api /etc/nginx/sites-enabled/
rm -f /etc/nginx/sites-enabled/default

# Test nginx config
if nginx -t; then
    systemctl restart nginx
    print_status "Nginx configured successfully"
else
    print_error "Nginx configuration failed"
fi

# Step 10: Wait for API to be ready and get default credentials
print_header "🔑 Step 10: Getting default admin credentials..."
sleep 5

# Check if credentials file exists
CRED_FILE="/etc/apivpn/default_credentials.txt"
if [ -f "$CRED_FILE" ]; then
    print_status "Default credentials found:"
    cat "$CRED_FILE"
else
    print_warning "Default credentials not found. Creating admin user..."
    # Force restart to create default user
    systemctl restart vpn-api
    sleep 5
    
    if [ -f "$CRED_FILE" ]; then
        print_status "Default credentials created:"
        cat "$CRED_FILE"
    else
        print_error "Failed to create default credentials"
    fi
fi

# Step 11: Test API
print_header "🧪 Step 11: Testing API..."

# Test if API is responding
if curl -s http://localhost:8080/api/v1/auth/login > /dev/null; then
    print_status "API is responding on port 8080"
else
    print_error "API is not responding"
    print_status "Checking service status..."
    systemctl status vpn-api --no-pager
fi

# Step 12: Create test script
print_header "📝 Step 12: Creating test scripts..."
cat > /root/test_api.sh << 'EOF'
#!/bin/bash

# Test script for VPN API
SERVER_IP=$(curl -s ipinfo.io/ip)
API_URL="http://$SERVER_IP:8080/api/v1"

echo "🧪 Testing VPN API..."
echo "API URL: $API_URL"
echo ""

# Read default credentials
if [ -f "/etc/apivpn/default_credentials.txt" ]; then
    echo "📋 Default credentials:"
    cat /etc/apivpn/default_credentials.txt
    echo ""
    
    # Extract password from credentials file
    DEFAULT_PASS=$(grep "Password:" /etc/apivpn/default_credentials.txt | cut -d' ' -f2)
    
    echo "🔑 Testing login..."
    LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/auth/login" \
        -H "Content-Type: application/json" \
        -d "{\"username\":\"admin\",\"password\":\"$DEFAULT_PASS\"}")
    
    echo "Login response:"
    echo "$LOGIN_RESPONSE" | jq . 2>/dev/null || echo "$LOGIN_RESPONSE"
    echo ""
    
    # Extract token
    TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.token' 2>/dev/null)
    
    if [ "$TOKEN" != "null" ] && [ "$TOKEN" != "" ]; then
        echo "✅ Login successful! Token obtained."
        echo ""
        
        echo "📊 Testing system info..."
        curl -s -X GET "$API_URL/system/info" \
            -H "Authorization: Bearer $TOKEN" | jq . 2>/dev/null
        echo ""
        
        echo "🔧 Testing service status..."
        curl -s -X GET "$API_URL/system/status" \
            -H "Authorization: Bearer $TOKEN" | jq . 2>/dev/null
        echo ""
        
        echo "✅ API is working correctly!"
        echo ""
        echo "🌐 You can now access the API at:"
        echo "   - Direct: http://$SERVER_IP:8080/api/v1"
        echo "   - Via Nginx: http://$SERVER_IP/api/v1"
        echo ""
        echo "📱 Ready for bot integration!"
    else
        echo "❌ Login failed"
    fi
else
    echo "❌ Default credentials file not found"
fi
EOF

chmod +x /root/test_api.sh

# Final summary
print_header "🎉 Setup Complete!"
echo "=================================="
print_status "VPN API has been successfully installed and configured!"
echo ""
echo "📋 Summary:"
echo "   - API URL: http://$SERVER_IP:8080/api/v1"
echo "   - Nginx URL: http://$SERVER_IP/api/v1"
echo "   - Service: systemctl status vpn-api"
echo "   - Logs: journalctl -u vpn-api -f"
echo "   - Test: /root/test_api.sh"
echo ""
echo "🔑 Default admin credentials:"
if [ -f "$CRED_FILE" ]; then
    cat "$CRED_FILE"
else
    echo "   Check: cat /etc/apivpn/default_credentials.txt"
fi
echo ""
echo "🧪 Run test script:"
echo "   bash /root/test_api.sh"
echo ""
echo "📱 Next steps:"
echo "   1. Test the API with the test script"
echo "   2. Setup WhatsApp/Telegram bot"
echo "   3. Create web interface"
echo ""
print_status "Setup completed successfully! 🚀"