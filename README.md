# VPN Management API

A comprehensive Go API for managing VPN services including SSH, VMESS, VLESS, Trojan, and Shadowsocks protocols. This API is designed to work with existing VPN server setups and provides endpoints for monitoring, user management, and system administration.

## Features

- **Multi-Protocol Support**: SSH/WebSocket, VMESS, VLESS, Trojan, Shadowsocks
- **System Monitoring**: Real-time server stats, bandwidth usage, service status
- **User Management**: Create, delete, extend VPN users
- **Domain & SSL Management**: Add domains and renew SSL certificates
- **RESTful API**: Clean JSON API with proper error handling
- **JWT Authentication**: Secure admin authentication
- **Bot Integration Ready**: Perfect for WhatsApp, Telegram bots, and web interfaces

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Linux server with VPN services installed
- Root access for system operations

### Installation

1. Clone the repository:
```bash
git clone https://github.com/nabilulilalbab/apivpn.git
cd apivpn
```

2. Install dependencies:
```bash
go mod tidy
```

3. Create default admin user:
```bash
sudo go run main.go
```

4. Check default credentials:
```bash
sudo cat /etc/apivpn/default_credentials.txt
```

5. Start the API server:
```bash
sudo go run main.go
```

The API will be available at `http://localhost:8080`

## API Documentation

### Authentication

All protected endpoints require a JWT token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

### Get JWT Token

**POST** `/api/v1/auth/login`

```json
{
  "username": "admin",
  "password": "your-password"
}
```

Response:
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

### System Monitoring

#### Get System Information
**GET** `/api/v1/system/info`

Returns server information including CPU, RAM, uptime, and bandwidth usage.

#### Get Service Status
**GET** `/api/v1/system/status`

Returns status of all VPN services (SSH, Nginx, Xray, etc.).

#### Get Bandwidth Usage
**GET** `/api/v1/system/bandwidth`

Returns daily and monthly bandwidth usage.

### VPN User Management

#### Create SSH User
**POST** `/api/v1/vpn/ssh/create`

```json
{
  "username": "testuser",
  "password": "securepassword",
  "days": 30
}
```

#### Create VMESS User
**POST** `/api/v1/vpn/vmess/create`

```json
{
  "username": "testuser",
  "days": 30
}
```

#### Get All Users
**GET** `/api/v1/vpn/users/all`

Returns all VPN users across all protocols.

#### Delete User
**DELETE** `/api/v1/vpn/{protocol}/users/{username}`

#### Extend User
**PUT** `/api/v1/vpn/{protocol}/users/{username}/extend`

```json
{
  "days": 30
}
```

### Domain Management

#### Add Domain
**POST** `/api/v1/domain/add`

```json
{
  "domain": "yourdomain.com"
}
```

#### Renew SSL Certificate
**POST** `/api/v1/domain/ssl/renew`

#### Get Current Domain
**GET** `/api/v1/domain/current`

## Bot Integration Examples

### WhatsApp Bot Integration

```javascript
// Example using whatsapp-web.js
const { Client } = require('whatsapp-web.js');
const axios = require('axios');

const client = new Client();
const API_BASE = 'http://your-server:8080/api/v1';
let authToken = '';

// Login to API
async function loginToAPI() {
  const response = await axios.post(`${API_BASE}/auth/login`, {
    username: 'admin',
    password: 'your-password'
  });
  authToken = response.data.data.token;
}

// Create VPN user
async function createVPNUser(protocol, username, days) {
  const response = await axios.post(`${API_BASE}/vpn/${protocol}/create`, {
    username: username,
    password: 'generated-password',
    days: days
  }, {
    headers: { Authorization: `Bearer ${authToken}` }
  });
  return response.data;
}

client.on('message', async msg => {
  if (msg.body.startsWith('/create')) {
    const [, protocol, username, days] = msg.body.split(' ');
    try {
      const result = await createVPNUser(protocol, username, parseInt(days));
      msg.reply(`VPN user created successfully!\n\nConfig: ${JSON.stringify(result.data.config, null, 2)}`);
    } catch (error) {
      msg.reply('Failed to create VPN user: ' + error.message);
    }
  }
});
```

### Telegram Bot Integration

```python
# Example using python-telegram-bot
import requests
import json
from telegram.ext import Application, CommandHandler

API_BASE = 'http://your-server:8080/api/v1'
auth_token = ''

def login_to_api():
    global auth_token
    response = requests.post(f'{API_BASE}/auth/login', json={
        'username': 'admin',
        'password': 'your-password'
    })
    auth_token = response.json()['data']['token']

async def create_vpn_user(update, context):
    args = context.args
    if len(args) < 3:
        await update.message.reply_text('Usage: /create <protocol> <username> <days>')
        return
    
    protocol, username, days = args[0], args[1], int(args[2])
    
    headers = {'Authorization': f'Bearer {auth_token}'}
    response = requests.post(f'{API_BASE}/vpn/{protocol}/create', 
                           json={'username': username, 'days': days}, 
                           headers=headers)
    
    if response.status_code == 201:
        config = response.json()['data']
        await update.message.reply_text(f'VPN user created!\n\n```json\n{json.dumps(config, indent=2)}\n```', 
                                      parse_mode='Markdown')
    else:
        await update.message.reply_text('Failed to create VPN user')

app = Application.builder().token('YOUR_BOT_TOKEN').build()
app.add_handler(CommandHandler('create', create_vpn_user))
```

## Environment Variables

- `PORT`: API server port (default: 8080)
- `JWT_SECRET`: JWT signing secret (change this!)
- `DOMAIN`: Server domain
- `XRAY_PATH`: Path to Xray config (default: /etc/xray)
- `SSH_PATH`: Path to SSH config (default: /etc/ssh)

## Security Notes

1. **Change Default Password**: Always change the default admin password after first login
2. **Use HTTPS**: Deploy behind a reverse proxy with SSL/TLS
3. **Firewall**: Restrict API access to trusted IPs
4. **JWT Secret**: Use a strong, unique JWT secret
5. **Regular Updates**: Keep the system and dependencies updated

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support and questions:
- Create an issue on GitHub
- Contact: [your-email@example.com]

## Roadmap

- [ ] Database integration for better user management
- [ ] Real-time traffic monitoring
- [ ] User bandwidth limits
- [ ] Multi-server support
- [ ] Web dashboard
- [ ] Docker containerization
- [ ] Automated backups