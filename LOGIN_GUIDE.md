# 🔑 Panduan Login API VPN - VPS: 128.199.227.169

## **Cara Upload dan Setup di VPS**

### **1. Upload File ke VPS Anda:**

```bash
# Di komputer lokal (terminal/cmd)
scp apivpn.tar.gz root@128.199.227.169:/root/

# Login ke VPS
ssh root@128.199.227.169

# Extract dan setup
cd /root
tar -xzf apivpn.tar.gz
cd apivpn

# Jalankan auto setup
chmod +x auto_setup.sh
./auto_setup.sh
```

## **2. Cara Menggunakan POST /api/v1/auth/login**

### **Method 1: Menggunakan curl (Terminal)**

```bash
# Login ke VPS dulu
ssh root@128.199.227.169

# Cek password default yang dibuat otomatis
cat /etc/apivpn/default_credentials.txt

# Contoh output:
# Default Admin Credentials:
# Username: admin
# Password: a1b2c3d4e5f6

# Login menggunakan curl
curl -X POST http://128.199.227.169:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "a1b2c3d4e5f6"
  }'
```

**Response yang akan Anda terima:**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjUyNzI4MDAsImlhdCI6MTcyNTE4NjQwMCwidXNlcl9pZCI6ImFkbWluIiwidXNlcm5hbWUiOiJhZG1pbiJ9.xyz123...",
    "username": "admin",
    "expires_at": "2024-09-02T14:20:00Z"
  }
}
```

### **Method 2: Menggunakan Postman/Insomnia**

**URL:** `http://128.199.227.169:8080/api/v1/auth/login`

**Method:** `POST`

**Headers:**
```
Content-Type: application/json
```

**Body (raw JSON):**
```json
{
  "username": "admin",
  "password": "password-dari-file-default"
}
```

### **Method 3: Menggunakan JavaScript (untuk bot)**

```javascript
const axios = require('axios');

async function loginToAPI() {
    try {
        const response = await axios.post('http://128.199.227.169:8080/api/v1/auth/login', {
            username: 'admin',
            password: 'password-dari-file-default'
        });
        
        console.log('Login successful!');
        console.log('Token:', response.data.data.token);
        return response.data.data.token;
    } catch (error) {
        console.error('Login failed:', error.response?.data || error.message);
        return null;
    }
}

// Gunakan token untuk request lain
async function getSystemInfo(token) {
    try {
        const response = await axios.get('http://128.199.227.169:8080/api/v1/system/info', {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        console.log('System Info:', response.data);
    } catch (error) {
        console.error('Failed to get system info:', error.response?.data || error.message);
    }
}

// Contoh penggunaan
loginToAPI().then(token => {
    if (token) {
        getSystemInfo(token);
    }
});
```

### **Method 4: Menggunakan Python (untuk bot)**

```python
import requests
import json

def login_to_api():
    url = 'http://128.199.227.169:8080/api/v1/auth/login'
    data = {
        'username': 'admin',
        'password': 'password-dari-file-default'  # Ganti dengan password sebenarnya
    }
    
    try:
        response = requests.post(url, json=data)
        response.raise_for_status()
        
        result = response.json()
        print('Login successful!')
        print('Token:', result['data']['token'])
        return result['data']['token']
    except requests.exceptions.RequestException as e:
        print('Login failed:', e)
        return None

def get_system_info(token):
    url = 'http://128.199.227.169:8080/api/v1/system/info'
    headers = {
        'Authorization': f'Bearer {token}'
    }
    
    try:
        response = requests.get(url, headers=headers)
        response.raise_for_status()
        
        result = response.json()
        print('System Info:', json.dumps(result, indent=2))
    except requests.exceptions.RequestException as e:
        print('Failed to get system info:', e)

# Contoh penggunaan
token = login_to_api()
if token:
    get_system_info(token)
```

## **3. Contoh Penggunaan Token untuk Endpoint Lain**

Setelah mendapat token, Anda bisa menggunakan untuk endpoint lain:

### **Get System Information:**
```bash
curl -X GET http://128.199.227.169:8080/api/v1/system/info \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### **Create SSH User:**
```bash
curl -X POST http://128.199.227.169:8080/api/v1/vpn/ssh/create \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "customer1",
    "password": "password123",
    "days": 30
  }'
```

### **List SSH Users:**
```bash
curl -X GET http://128.199.227.169:8080/api/v1/vpn/ssh/users \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### **Create VMESS User:**
```bash
curl -X POST http://128.199.227.169:8080/api/v1/vpn/vmess/create \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "customer2",
    "days": 30
  }'
```

## **4. Troubleshooting**

### **Jika mendapat error "Connection refused":**
```bash
# Cek apakah service berjalan
systemctl status vpn-api

# Jika tidak berjalan, start service
systemctl start vpn-api

# Cek logs jika ada error
journalctl -u vpn-api -f
```

### **Jika mendapat error "Invalid credentials":**
```bash
# Cek password default
cat /etc/apivpn/default_credentials.txt

# Atau reset admin user
rm -f /etc/apivpn/users.json
systemctl restart vpn-api
sleep 3
cat /etc/apivpn/default_credentials.txt
```

### **Jika port 8080 tidak bisa diakses:**
```bash
# Cek firewall
ufw status

# Allow port 8080
ufw allow 8080/tcp

# Cek apakah port terbuka
netstat -tlnp | grep 8080
```

## **5. Test Script Otomatis**

Setelah setup, jalankan test script:

```bash
# Di VPS
bash /root/test_api.sh
```

Script ini akan:
- ✅ Test login otomatis
- ✅ Test system info
- ✅ Test service status
- ✅ Menampilkan semua URL yang bisa digunakan

## **6. URLs yang Tersedia**

- **Login:** `http://128.199.227.169:8080/api/v1/auth/login`
- **System Info:** `http://128.199.227.169:8080/api/v1/system/info`
- **Service Status:** `http://128.199.227.169:8080/api/v1/system/status`
- **Create SSH User:** `http://128.199.227.169:8080/api/v1/vpn/ssh/create`
- **List SSH Users:** `http://128.199.227.169:8080/api/v1/vpn/ssh/users`
- **Create VMESS User:** `http://128.199.227.169:8080/api/v1/vpn/vmess/create`
- **List All Users:** `http://128.199.227.169:8080/api/v1/vpn/users/all`

**Apakah Anda ingin saya buatkan contoh bot WhatsApp atau Telegram yang langsung bisa digunakan dengan VPS ini?**