# 🚀 VPN API - Setup Lengkap dengan SQLite Database

## **📋 Informasi Penting**

- **Port API:** `37849` (port yang sangat jarang digunakan)
- **Database:** SQLite (otomatis dibuat)
- **Environment:** Menggunakan file `.env`
- **Security:** JWT dengan bcrypt, rate limiting
- **VPS:** 128.199.227.169

## **🔧 Cara Setup di VPS**

### **1. Upload File ke VPS**

```bash
# Di komputer lokal
scp apivpn-ready.tar.gz root@128.199.227.169:/root/

# Login ke VPS
ssh root@128.199.227.169

# Extract file
cd /root
tar -xzf apivpn-ready.tar.gz
cd apivpn
```

### **2. Auto Setup (Satu Perintah)**

```bash
# Jalankan auto setup
chmod +x auto_setup.sh
./auto_setup.sh
```

**Script akan otomatis:**
- ✅ Install Go 1.21
- ✅ Install dependencies (vnstat, nginx, dll)
- ✅ Build API dengan SQLite
- ✅ Setup database otomatis
- ✅ Create systemd service
- ✅ Configure firewall (port 37849)
- ✅ Setup Nginx reverse proxy
- ✅ Create default admin user

### **3. Cek Status Setup**

```bash
# Cek service status
systemctl status vpn-api

# Cek logs
journalctl -u vpn-api -f

# Test API
bash /root/test_api.sh

# Health check
curl http://128.199.227.169:37849/health
```

## **🔑 Login ke API**

### **Cek Password Default**

```bash
cat /etc/apivpn/default_credentials.txt
```

**Output contoh:**
```
Default Admin Credentials:
Username: admin
Password: a1b2c3d4e5f6

Please change this password after first login!
```

### **Test Login**

```bash
curl -X POST http://128.199.227.169:37849/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "a1b2c3d4e5f6"
  }'
```

**Response sukses:**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "username": "admin",
    "expires_at": "2024-09-02T14:20:00Z"
  }
}
```

## **📡 Endpoint API Utama**

### **Authentication**
- `POST /api/v1/auth/login` - Login admin
- `POST /api/v1/auth/register` - Register admin baru

### **System Monitoring**
- `GET /api/v1/system/info` - Info sistem
- `GET /api/v1/system/status` - Status service
- `GET /api/v1/system/bandwidth` - Usage bandwidth
- `GET /health` - Health check

### **VPN Management**
- `POST /api/v1/vpn/ssh/create` - Create SSH user
- `POST /api/v1/vpn/vmess/create` - Create VMESS user
- `POST /api/v1/vpn/vless/create` - Create VLESS user
- `POST /api/v1/vpn/trojan/create` - Create Trojan user
- `POST /api/v1/vpn/shadowsocks/create` - Create Shadowsocks user
- `GET /api/v1/vpn/users/all` - List semua users
- `DELETE /api/v1/vpn/{protocol}/users/{username}` - Delete user

### **User Management**
- `GET /api/v1/user/profile` - Profile user
- `PUT /api/v1/user/password` - Change password
- `GET /api/v1/user/list` - List admin users

## **🤖 Contoh untuk Bot**

### **Environment Variables untuk Bot**

```bash
# Untuk WhatsApp Bot
API_BASE_URL="http://128.199.227.169:37849/api/v1"
ADMIN_USERNAME="admin"
ADMIN_PASSWORD="a1b2c3d4e5f6"  # Ganti dengan password sebenarnya

# Untuk Telegram Bot
API_BASE_URL = 'http://128.199.227.169:37849/api/v1'
ADMIN_USERNAME = 'admin'
ADMIN_PASSWORD = 'a1b2c3d4e5f6'  # Ganti dengan password sebenarnya
```

### **Contoh Create User via API**

```bash
# Login dulu untuk dapat token
TOKEN=$(curl -s -X POST http://128.199.227.169:37849/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"a1b2c3d4e5f6"}' | \
  jq -r '.data.token')

# Create SSH user
curl -X POST http://128.199.227.169:37849/api/v1/vpn/ssh/create \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "customer1",
    "password": "password123",
    "days": 30
  }'

# Create VMESS user
curl -X POST http://128.199.227.169:37849/api/v1/vpn/vmess/create \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "customer2",
    "days": 30
  }'
```

## **💾 Database SQLite**

Database otomatis dibuat di: `/etc/apivpn/vpnapi.db`

**Tables:**
- `users` - Admin users
- `vpn_users` - VPN customers
- `traffic_logs` - Bandwidth tracking
- `system_logs` - System events

**Backup database:**
```bash
cp /etc/apivpn/vpnapi.db /root/backup-$(date +%Y%m%d).db
```

## **🔒 Security Features**

- ✅ JWT Authentication dengan expiry
- ✅ Password hashing dengan bcrypt
- ✅ Rate limiting (100 requests/minute)
- ✅ Login attempt limiting
- ✅ Database logging semua events
- ✅ Port non-standard (37849)

## **📊 Monitoring**

```bash
# Cek database stats
curl -H "Authorization: Bearer $TOKEN" \
  http://128.199.227.169:37849/health

# Cek system info
curl -H "Authorization: Bearer $TOKEN" \
  http://128.199.227.169:37849/api/v1/system/info

# Cek VPN stats
curl -H "Authorization: Bearer $TOKEN" \
  http://128.199.227.169:37849/api/v1/vpn/stats
```

## **🛠️ Troubleshooting**

### **Service tidak jalan:**
```bash
systemctl status vpn-api
journalctl -u vpn-api -f
systemctl restart vpn-api
```

### **Port tidak bisa diakses:**
```bash
ufw status
ufw allow 37849/tcp
netstat -tlnp | grep 37849
```

### **Database error:**
```bash
ls -la /etc/apivpn/
chmod 755 /etc/apivpn
rm /etc/apivpn/vpnapi.db  # Reset database
systemctl restart vpn-api
```

### **Reset admin password:**
```bash
rm /etc/apivpn/vpnapi.db
systemctl restart vpn-api
cat /etc/apivpn/default_credentials.txt
```

## **🌐 URLs Lengkap**

- **API Direct:** `http://128.199.227.169:37849/api/v1`
- **Via Nginx:** `http://128.199.227.169/api/v1`
- **Health Check:** `http://128.199.227.169:37849/health`
- **Login:** `http://128.199.227.169:37849/api/v1/auth/login`

## **📱 Ready untuk Bot Integration!**

API sudah siap untuk:
- 🤖 WhatsApp Bot
- 📱 Telegram Bot  
- 🌐 Web Dashboard
- 📊 Mobile App

**Apakah Anda ingin saya buatkan bot WhatsApp atau Telegram yang langsung connect ke API ini?**