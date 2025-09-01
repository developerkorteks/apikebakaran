# 🚀 Setup Guide untuk VPS: 128.199.227.169

## **Step 1: Upload File ke VPS**

```bash
# Di komputer lokal, upload file ke VPS
scp apivpn.tar.gz root@128.199.227.169:/root/

# Login ke VPS
ssh root@128.199.227.169

# Extract file
cd /root
tar -xzf apivpn.tar.gz
cd apivpn
```

## **Step 2: Install Dependencies**

```bash
# Update system
apt update && apt upgrade -y

# Install Go (jika belum ada)
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Install dependencies
apt install -y curl vnstat bc nginx

# Enable vnstat
systemctl enable vnstat
systemctl start vnstat
```

## **Step 3: Build dan Setup API**

```bash
# Build API
go mod tidy
go build -o vpn-api main.go

# Setup directories
mkdir -p /etc/apivpn
mkdir -p /etc/xray
mkdir -p /var/lib/scrz-prem

# Set permissions
chmod +x vpn-api
chmod +x scripts/init.sh

# Setup domain (ganti dengan domain Anda atau gunakan IP)
echo "128.199.227.169" > /etc/xray/domain
echo "IP=128.199.227.169" > /var/lib/scrz-prem/ipvps.conf
```

## **Step 4: Create Systemd Service**

```bash
# Create service file
cat > /etc/systemd/system/vpn-api.service << 'EOF'
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
Environment=JWT_SECRET=your-super-secret-jwt-key-change-this

[Install]
WantedBy=multi-user.target
EOF

# Enable and start service
systemctl daemon-reload
systemctl enable vpn-api
systemctl start vpn-api

# Check status
systemctl status vpn-api
```

## **Step 5: Test API Login**

### **Method 1: Menggunakan curl**

```bash
# Login untuk mendapatkan token
curl -X POST http://128.199.227.169:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "lihat-password-default"
  }'
```

**Catatan:** Password default akan dibuat otomatis. Cek di:
```bash
cat /etc/apivpn/default_credentials.txt
```

### **Method 2: Menggunakan Postman/Insomnia**

**URL:** `http://128.199.227.169:8080/api/v1/auth/login`
**Method:** POST
**Headers:**
```
Content-Type: application/json
```
**Body (JSON):**
```json
{
  "username": "admin",
  "password": "password-dari-file-default"
}
```

### **Response yang diharapkan:**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "username": "admin",
    "expires_at": "2024-01-02T15:04:05Z"
  }
}
```

## **Step 6: Test API Endpoints Lainnya**

### **Get System Info (butuh token):**
```bash
# Ganti YOUR_TOKEN dengan token dari login
curl -X GET http://128.199.227.169:8080/api/v1/system/info \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### **Create SSH User:**
```bash
curl -X POST http://128.199.227.169:8080/api/v1/vpn/ssh/create \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "testpass123",
    "days": 30
  }'
```

### **List SSH Users:**
```bash
curl -X GET http://128.199.227.169:8080/api/v1/vpn/ssh/users \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## **Step 7: Setup Firewall (Opsional)**

```bash
# Allow API port
ufw allow 8080/tcp

# Allow SSH (pastikan tidak terkunci)
ufw allow 22/tcp

# Enable firewall
ufw --force enable
```

## **Step 8: Setup Nginx Reverse Proxy (Opsional)**

```bash
# Create nginx config
cat > /etc/nginx/sites-available/vpn-api << 'EOF'
server {
    listen 80;
    server_name 128.199.227.169;

    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
EOF

# Enable site
ln -s /etc/nginx/sites-available/vpn-api /etc/nginx/sites-enabled/
nginx -t
systemctl restart nginx
```

## **Troubleshooting**

### **Jika API tidak jalan:**
```bash
# Check logs
journalctl -u vpn-api -f

# Check if port is open
netstat -tlnp | grep 8080

# Restart service
systemctl restart vpn-api
```

### **Jika login gagal:**
```bash
# Check default credentials
cat /etc/apivpn/default_credentials.txt

# Reset admin user (hapus file users.json)
rm -f /etc/apivpn/users.json
systemctl restart vpn-api
```

## **URLs untuk Testing:**

- **API Base:** `http://128.199.227.169:8080/api/v1`
- **Login:** `http://128.199.227.169:8080/api/v1/auth/login`
- **System Info:** `http://128.199.227.169:8080/api/v1/system/info`
- **Create SSH User:** `http://128.199.227.169:8080/api/v1/vpn/ssh/create`

## **Next Steps:**

1. ✅ Setup API di VPS
2. 🔑 Test login dan dapatkan token
3. 🤖 Setup WhatsApp/Telegram bot
4. 🌐 Buat web interface
5. 📊 Monitor penggunaan

**Apakah Anda ingin saya buatkan script otomatis untuk setup ini?**