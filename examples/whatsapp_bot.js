// WhatsApp Bot Example for VPN API Integration
// Install dependencies: npm install whatsapp-web.js axios qrcode-terminal

const { Client, LocalAuth } = require('whatsapp-web.js');
const axios = require('axios');
const qrcode = require('qrcode-terminal');

// Configuration
const API_BASE_URL = 'http://localhost:8080/api/v1';
const ADMIN_USERNAME = 'admin';
const ADMIN_PASSWORD = 'your-password'; // Change this!

class VPNBot {
    constructor() {
        this.client = new Client({
            authStrategy: new LocalAuth()
        });
        this.authToken = '';
        this.authorizedUsers = ['6281234567890@c.us']; // Add authorized WhatsApp numbers
        
        this.setupEventHandlers();
    }

    async init() {
        await this.loginToAPI();
        this.client.initialize();
    }

    async loginToAPI() {
        try {
            const response = await axios.post(`${API_BASE_URL}/auth/login`, {
                username: ADMIN_USERNAME,
                password: ADMIN_PASSWORD
            });
            this.authToken = response.data.data.token;
            console.log('✅ Successfully logged in to VPN API');
        } catch (error) {
            console.error('❌ Failed to login to VPN API:', error.message);
            process.exit(1);
        }
    }

    setupEventHandlers() {
        this.client.on('qr', (qr) => {
            console.log('📱 Scan this QR code with WhatsApp:');
            qrcode.generate(qr, { small: true });
        });

        this.client.on('ready', () => {
            console.log('🚀 WhatsApp Bot is ready!');
        });

        this.client.on('message', async (message) => {
            await this.handleMessage(message);
        });
    }

    isAuthorized(userId) {
        return this.authorizedUsers.includes(userId);
    }

    async handleMessage(message) {
        if (!this.isAuthorized(message.from)) {
            return;
        }

        const text = message.body.toLowerCase().trim();
        const args = text.split(' ');
        const command = args[0];

        try {
            switch (command) {
                case '/help':
                    await this.sendHelp(message);
                    break;
                case '/status':
                    await this.getSystemStatus(message);
                    break;
                case '/create':
                    await this.createVPNUser(message, args);
                    break;
                case '/list':
                    await this.listUsers(message, args);
                    break;
                case '/delete':
                    await this.deleteUser(message, args);
                    break;
                case '/extend':
                    await this.extendUser(message, args);
                    break;
                case '/traffic':
                    await this.getUserTraffic(message, args);
                    break;
                case '/info':
                    await this.getServerInfo(message);
                    break;
                default:
                    if (text.startsWith('/')) {
                        message.reply('❌ Unknown command. Type /help for available commands.');
                    }
                    break;
            }
        } catch (error) {
            console.error('Error handling message:', error);
            message.reply('❌ An error occurred while processing your request.');
        }
    }

    async sendHelp(message) {
        const helpText = `
🤖 *VPN Bot Commands*

📊 *System Commands:*
/status - Check service status
/info - Get server information

👥 *User Management:*
/create <protocol> <username> <days> - Create VPN user
/list <protocol> - List users by protocol
/delete <protocol> <username> - Delete user
/extend <protocol> <username> <days> - Extend user
/traffic <username> - Get user traffic

📋 *Supported Protocols:*
• ssh - SSH/WebSocket
• vmess - VMESS
• vless - VLESS  
• trojan - Trojan
• shadowsocks - Shadowsocks

💡 *Examples:*
/create ssh john 30
/list vmess
/delete ssh john
/extend vless alice 15
        `;
        message.reply(helpText);
    }

    async getSystemStatus(message) {
        try {
            const response = await this.apiRequest('GET', '/system/status');
            const status = response.data;
            
            const statusText = `
🖥️ *System Status*

🔐 SSH: ${status.ssh ? '✅' : '❌'}
🌐 Nginx: ${status.nginx ? '✅' : '❌'}
⚡ Xray: ${status.xray ? '✅' : '❌'}
🔒 Dropbear: ${status.dropbear ? '✅' : '❌'}
🔐 Stunnel: ${status.stunnel ? '✅' : '❌'}
🌐 SSH-WS: ${status.ssh_websocket ? '✅' : '❌'}
            `;
            message.reply(statusText);
        } catch (error) {
            message.reply('❌ Failed to get system status');
        }
    }

    async createVPNUser(message, args) {
        if (args.length < 4) {
            message.reply('❌ Usage: /create <protocol> <username> <days>\nExample: /create ssh john 30');
            return;
        }

        const [, protocol, username, days] = args;
        
        try {
            const response = await this.apiRequest('POST', `/vpn/${protocol}/create`, {
                username: username,
                password: this.generatePassword(),
                days: parseInt(days)
            });

            const config = response.data;
            const configText = this.formatVPNConfig(protocol, config);
            
            message.reply(`✅ *VPN User Created Successfully!*\n\n${configText}`);
        } catch (error) {
            message.reply(`❌ Failed to create ${protocol} user: ${error.response?.data?.error || error.message}`);
        }
    }

    async listUsers(message, args) {
        if (args.length < 2) {
            message.reply('❌ Usage: /list <protocol>\nExample: /list ssh');
            return;
        }

        const protocol = args[1];
        
        try {
            const response = await this.apiRequest('GET', `/vpn/${protocol}/users`);
            const users = response.data;
            
            if (users.length === 0) {
                message.reply(`📋 No ${protocol} users found.`);
                return;
            }

            let userList = `📋 *${protocol.toUpperCase()} Users:*\n\n`;
            users.forEach((user, index) => {
                const status = user.is_active ? '✅' : '❌';
                const expiry = user.expiry_date ? new Date(user.expiry_date).toLocaleDateString() : 'Never';
                userList += `${index + 1}. ${status} ${user.username} (Expires: ${expiry})\n`;
            });

            message.reply(userList);
        } catch (error) {
            message.reply(`❌ Failed to list ${protocol} users`);
        }
    }

    async deleteUser(message, args) {
        if (args.length < 3) {
            message.reply('❌ Usage: /delete <protocol> <username>\nExample: /delete ssh john');
            return;
        }

        const [, protocol, username] = args;
        
        try {
            await this.apiRequest('DELETE', `/vpn/${protocol}/users/${username}`);
            message.reply(`✅ Successfully deleted ${protocol} user: ${username}`);
        } catch (error) {
            message.reply(`❌ Failed to delete ${protocol} user: ${username}`);
        }
    }

    async extendUser(message, args) {
        if (args.length < 4) {
            message.reply('❌ Usage: /extend <protocol> <username> <days>\nExample: /extend ssh john 30');
            return;
        }

        const [, protocol, username, days] = args;
        
        try {
            await this.apiRequest('PUT', `/vpn/${protocol}/users/${username}/extend`, {
                days: parseInt(days)
            });
            message.reply(`✅ Successfully extended ${protocol} user: ${username} by ${days} days`);
        } catch (error) {
            message.reply(`❌ Failed to extend ${protocol} user: ${username}`);
        }
    }

    async getUserTraffic(message, args) {
        if (args.length < 2) {
            message.reply('❌ Usage: /traffic <username>\nExample: /traffic john');
            return;
        }

        const username = args[1];
        
        try {
            const response = await this.apiRequest('GET', `/vpn/users/${username}/traffic`);
            const traffic = response.data;
            
            const trafficText = `
📊 *Traffic Usage for ${username}*

⬆️ Upload: ${traffic.upload}
⬇️ Download: ${traffic.download}
📈 Total: ${traffic.total}
            `;
            message.reply(trafficText);
        } catch (error) {
            message.reply(`❌ Failed to get traffic for user: ${username}`);
        }
    }

    async getServerInfo(message) {
        try {
            const response = await this.apiRequest('GET', '/system/info');
            const info = response.data;
            
            const infoText = `
🖥️ *Server Information*

💻 OS: ${info.os}
🔧 Kernel: ${info.kernel}
⚡ CPU: ${info.cpu_name}
🧠 Cores: ${info.cpu_cores}
📊 CPU Usage: ${info.cpu_usage}
💾 RAM: ${info.ram_used_mb}MB / ${info.ram_total_mb}MB (${info.ram_usage_percent})
⏰ Uptime: ${info.uptime}
🌐 Domain: ${info.domain}
📍 IP: ${info.ip}
📈 Daily Bandwidth: ${info.daily_bandwidth}
📊 Monthly Bandwidth: ${info.monthly_bandwidth}
            `;
            message.reply(infoText);
        } catch (error) {
            message.reply('❌ Failed to get server information');
        }
    }

    formatVPNConfig(protocol, config) {
        let configText = `*Protocol:* ${protocol.toUpperCase()}\n`;
        configText += `*Server:* ${config.server}\n`;
        configText += `*Port:* ${config.port}\n`;
        configText += `*Username:* ${config.username}\n`;
        
        if (config.password) {
            configText += `*Password:* ${config.password}\n`;
        }
        
        if (config.uuid) {
            configText += `*UUID:* ${config.uuid}\n`;
        }

        if (config.config) {
            configText += `\n*Additional Config:*\n`;
            Object.entries(config.config).forEach(([key, value]) => {
                configText += `${key}: ${value}\n`;
            });
        }

        return configText;
    }

    generatePassword() {
        return Math.random().toString(36).slice(-8);
    }

    async apiRequest(method, endpoint, data = null) {
        const config = {
            method,
            url: `${API_BASE_URL}${endpoint}`,
            headers: {
                'Authorization': `Bearer ${this.authToken}`,
                'Content-Type': 'application/json'
            }
        };

        if (data) {
            config.data = data;
        }

        return await axios(config);
    }
}

// Initialize and start the bot
const bot = new VPNBot();
bot.init().catch(console.error);

// Handle graceful shutdown
process.on('SIGINT', () => {
    console.log('\n👋 Shutting down WhatsApp Bot...');
    process.exit(0);
});