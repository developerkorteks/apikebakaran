# 🚀 VPN API - Production Ready dengan SQLite Database

## ✅ **Status: 100% SIAP PRODUKSI**

Project ini sudah dianalisis ulang berdasarkan `menu_script.txt` dan dipastikan:
- ✅ **Service detection** menggunakan method yang sama dengan script asli
- ✅ **User counting** untuk SSH dan Xray sesuai dengan pattern asli  
- ✅ **System monitoring** menggunakan command yang identik
- ✅ **Path dan file** sesuai dengan struktur VPS yang ada
- ✅ **Error handling** yang robust untuk troubleshooting minimal
- ✅ **Dependency checker** untuk memastikan semua requirement terpenuhi

---

## 📦 **File Production Ready**

**File:** `apivpn-production-ready.tar.gz`

**Berisi:**
- ✅ Go API Server dengan SQLite database
- ✅ Port 37849 (sangat jarang digunakan)
- ✅ Environment configuration (.env)
- ✅ Auto setup script dengan dependency check
- ✅ Database models dan migrations
- ✅ JWT authentication dengan bcrypt
- ✅ Rate limiting dan security features
- ✅ Bot integration examples (WhatsApp & Telegram)
- ✅ Dependency checker script
- ✅ Production-ready configuration

---

## 🚀 **Setup di VPS: 128.199.227.169**

### **1. Upload ke VPS**
```bash
# Upload file
scp apivpn-production-ready.tar.gz root@128.199.227.169:/root/

# Login ke VPS
ssh root@128.199.227.169

# Extract
cd /root
tar -xzf apivpn-production-ready.tar.gz
cd apivpn
```

### **2. Auto Setup (Satu Perintah)**
```bash
chmod +x auto_setup.sh
./auto_setup.sh
```

**Script akan otomatis:**
1. 🔍 **Check dependencies** - memastikan semua command tersedia
2. 📦 **Install Go 1.21** dan dependencies
3. 🗂️ **Create directories** dan files yang dibutuhkan
4. 🔨 **Build API** dengan SQLite
5. ⚙️ **Setup systemd service** dengan environment variables
6. 🔥 **Configure firewall** untuk port 37849
7. 🌐 **Setup Nginx** reverse proxy
8. 🔑 **Create default admin** user
9. 🧪 **Test API** endpoints

### **3. Verifikasi Setup**
```bash
# Check service
systemctl status vpn-api

# Check logs
journalctl -u vpn-api -f

# Test API
bash /root/test_api.sh

# Health check
curl http://128.199.227.169:37849/health
```

---

## 🔑 **Login ke API**

### **Cek Password Default**
```bash
cat /etc/apivpn/default_credentials.txt
```

### **Test Login**
```bash
curl -X POST http://128.199.227.169:37849/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin", 
    "password": "password-dari-file-default"
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

---

## 📡 **API Endpoints Lengkap**

### **🔐 Authentication**
- `POST /api/v1/auth/login` - Login admin
- `POST /api/v1/auth/register` - Register admin baru

### **👤 User Management**
- `GET /api/v1/user/profile` - Profile user
- `PUT /api/v1/user/password` - Change password
- `GET /api/v1/user/list` - List admin users
- `PUT /api/v1/user/:username/status` - Enable/disable user
- `DELETE /api/v1/user/:username` - Delete user

### **📊 System Monitoring**
- `GET /api/v1/system/info` - Info sistem (CPU, RAM, uptime, dll)
- `GET /api/v1/system/status` - Status services (SSH, Nginx, Xray, dll)
- `GET /api/v1/system/bandwidth` - Usage bandwidth
- `POST /api/v1/system/reboot` - Reboot sistem
- `POST /api/v1/system/restart` - Restart services
- `GET /health` - Health check

### **🔧 VPN Management**

**SSH/WebSocket:**
- `POST /api/v1/vpn/ssh/create` - Create SSH user
- `GET /api/v1/vpn/ssh/users` - List SSH users
- `DELETE /api/v1/vpn/ssh/users/:username` - Delete SSH user
- `PUT /api/v1/vpn/ssh/users/:username/extend` - Extend SSH user

**VMESS:**
- `POST /api/v1/vpn/vmess/create` - Create VMESS user
- `GET /api/v1/vpn/vmess/users` - List VMESS users
- `DELETE /api/v1/vpn/vmess/users/:username` - Delete VMESS user
- `PUT /api/v1/vpn/vmess/users/:username/extend` - Extend VMESS user

**VLESS:**
- `POST /api/v1/vpn/vless/create` - Create VLESS user
- `GET /api/v1/vpn/vless/users` - List VLESS users
- `DELETE /api/v1/vpn/vless/users/:username` - Delete VLESS user
- `PUT /api/v1/vpn/vless/users/:username/extend` - Extend VLESS user

**Trojan:**
- `POST /api/v1/vpn/trojan/create` - Create Trojan user
- `GET /api/v1/vpn/trojan/users` - List Trojan users
- `DELETE /api/v1/vpn/trojan/users/:username` - Delete Trojan user
- `PUT /api/v1/vpn/trojan/users/:username/extend` - Extend Trojan user

**Shadowsocks:**
- `POST /api/v1/vpn/shadowsocks/create` - Create Shadowsocks user
- `GET /api/v1/vpn/shadowsocks/users` - List Shadowsocks users
- `DELETE /api/v1/vpn/shadowsocks/users/:username` - Delete Shadowsocks user
- `PUT /api/v1/vpn/shadowsocks/users/:username/extend` - Extend Shadowsocks user

**General:**
- `GET /api/v1/vpn/users/all` - List semua users
- `GET /api/v1/vpn/users/:username/traffic` - Get user traffic
- `POST /api/v1/vpn/users/cleanup-expired` - Cleanup expired users

### **🌐 Domain Management**
- `POST /api/v1/domain/add` - Add domain
- `POST /api/v1/domain/ssl/renew` - Renew SSL certificate
- `GET /api/v1/domain/current` - Get current domain

---

## 🤖 **Bot Integration**

### **Environment untuk Bot**
```bash
# Untuk WhatsApp Bot (JavaScript)
API_BASE_URL="http://128.199.227.169:37849/api/v1"
ADMIN_USERNAME="admin"
ADMIN_PASSWORD="password-dari-file-default"

# Untuk Telegram Bot (Python)
API_BASE_URL = 'http://128.199.227.169:37849/api/v1'
ADMIN_USERNAME = 'admin'
ADMIN_PASSWORD = 'password-dari-file-default'
```

### **Contoh Create User via API**
```bash
# Login untuk dapat token
TOKEN=$(curl -s -X POST http://128.199.227.169:37849/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password-sebenarnya"}' | \
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

---

## 💾 **Database SQLite**

**Path:** `/etc/apivpn/vpnapi.db`

**Tables:**
- `users` - Admin users dengan authentication
- `vpn_users` - VPN customers dengan expiry tracking
- `traffic_logs` - Bandwidth usage tracking
- `system_logs` - System events dan audit logs

**Backup:**
```bash
cp /etc/apivpn/vpnapi.db /root/backup-$(date +%Y%m%d-%H%M%S).db
```

**Reset database:**
```bash
rm /etc/apivpn/vpnapi.db
systemctl restart vpn-api
cat /etc/apivpn/default_credentials.txt
```

---

## 🔒 **Security Features**

- ✅ **Port non-standard** (37849) - mengurangi scan otomatis
- ✅ **JWT Authentication** dengan expiry 24 jam
- ✅ **Password bcrypt hashing** dengan cost 12
- ✅ **Rate limiting** 100 requests per menit per IP
- ✅ **Login attempt limiting** max 5 attempts
- ✅ **Database event logging** untuk audit trail
- ✅ **Input validation** dan sanitization
- ✅ **CORS protection** dan security headers

---

## 📊 **Monitoring & Troubleshooting**

### **Check Status**
```bash
# Service status
systemctl status vpn-api

# Real-time logs
journalctl -u vpn-api -f

# Database stats
curl -H "Authorization: Bearer $TOKEN" \
  http://128.199.227.169:37849/health

# System info
curl -H "Authorization: Bearer $TOKEN" \
  http://128.199.227.169:37849/api/v1/system/info
```

### **Common Issues**

**Service tidak start:**
```bash
journalctl -u vpn-api --no-pager -l
systemctl restart vpn-api
```

**Port tidak bisa diakses:**
```bash
ufw status
ufw allow 37849/tcp
netstat -tlnp | grep 37849
```

**Database error:**
```bash
ls -la /etc/apivpn/
chmod 755 /etc/apivpn
systemctl restart vpn-api
```

**Reset admin:**
```bash
rm /etc/apivpn/vpnapi.db
systemctl restart vpn-api
cat /etc/apivpn/default_credentials.txt
```

---

## 🌐 **URLs Lengkap**

- **API Direct:** `http://128.199.227.169:37849/api/v1`
- **Via Nginx:** `http://128.199.227.169/api/v1`
- **Health Check:** `http://128.199.227.169:37849/health`
- **Login:** `http://128.199.227.169:37849/api/v1/auth/login`

---

## 🎯 **Next Steps**

1. ✅ **Upload dan setup** di VPS
2. 🔑 **Test login** dan endpoints
3. 🤖 **Setup bot** (WhatsApp/Telegram)
4. 🌐 **Buat web dashboard** (opsional)
5. 📊 **Monitor usage** dan performance

---

## 📞 **Support**

Jika ada masalah:
1. Cek logs: `journalctl -u vpn-api -f`
2. Run dependency check: `./check_dependencies.sh`
3. Test API: `./test_api.sh`
4. Reset jika perlu: hapus database dan restart service

**Project ini sudah 100% production ready dan siap digunakan!** 🚀