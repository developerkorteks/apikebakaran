#!/usr/bin/env python3
"""
Telegram Bot Example for VPN API Integration
Install dependencies: pip install python-telegram-bot requests
"""

import asyncio
import json
import logging
import requests
from datetime import datetime
from telegram import Update, InlineKeyboardButton, InlineKeyboardMarkup
from telegram.ext import Application, CommandHandler, CallbackQueryHandler, ContextTypes

# Configuration
API_BASE_URL = 'http://localhost:8080/api/v1'
ADMIN_USERNAME = 'admin'
ADMIN_PASSWORD = 'your-password'  # Change this!
BOT_TOKEN = 'YOUR_BOT_TOKEN'  # Get from @BotFather
AUTHORIZED_USERS = [123456789, 987654321]  # Add authorized Telegram user IDs

# Setup logging
logging.basicConfig(
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    level=logging.INFO
)
logger = logging.getLogger(__name__)

class VPNTelegramBot:
    def __init__(self):
        self.auth_token = ''
        self.protocols = ['ssh', 'vmess', 'vless', 'trojan', 'shadowsocks']
    
    async def login_to_api(self):
        """Login to VPN API and get JWT token"""
        try:
            response = requests.post(f'{API_BASE_URL}/auth/login', json={
                'username': ADMIN_USERNAME,
                'password': ADMIN_PASSWORD
            })
            response.raise_for_status()
            self.auth_token = response.json()['data']['token']
            logger.info('✅ Successfully logged in to VPN API')
            return True
        except Exception as e:
            logger.error(f'❌ Failed to login to VPN API: {e}')
            return False
    
    def is_authorized(self, user_id: int) -> bool:
        """Check if user is authorized to use the bot"""
        return user_id in AUTHORIZED_USERS
    
    def api_request(self, method: str, endpoint: str, data=None):
        """Make API request with authentication"""
        headers = {
            'Authorization': f'Bearer {self.auth_token}',
            'Content-Type': 'application/json'
        }
        
        url = f'{API_BASE_URL}{endpoint}'
        
        if method.upper() == 'GET':
            response = requests.get(url, headers=headers)
        elif method.upper() == 'POST':
            response = requests.post(url, json=data, headers=headers)
        elif method.upper() == 'DELETE':
            response = requests.delete(url, headers=headers)
        elif method.upper() == 'PUT':
            response = requests.put(url, json=data, headers=headers)
        else:
            raise ValueError(f'Unsupported HTTP method: {method}')
        
        response.raise_for_status()
        return response.json()

    async def start_command(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        """Handle /start command"""
        if not self.is_authorized(update.effective_user.id):
            await update.message.reply_text('❌ You are not authorized to use this bot.')
            return
        
        welcome_text = """
🤖 **VPN Management Bot**

Welcome! I can help you manage your VPN server.

Available commands:
/help - Show all commands
/status - Check service status
/info - Get server information
/create - Create VPN user
/list - List users
/delete - Delete user
/extend - Extend user
/traffic - Get user traffic

Use /help for detailed command usage.
        """
        await update.message.reply_text(welcome_text, parse_mode='Markdown')

    async def help_command(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        """Handle /help command"""
        if not self.is_authorized(update.effective_user.id):
            return
        
        help_text = """
🤖 **VPN Bot Commands**

📊 **System Commands:**
/status - Check service status
/info - Get server information

👥 **User Management:**
/create - Create VPN user (interactive)
/list - List users by protocol
/delete - Delete user (interactive)
/extend - Extend user (interactive)
/traffic <username> - Get user traffic

📋 **Supported Protocols:**
• ssh - SSH/WebSocket
• vmess - VMESS
• vless - VLESS
• trojan - Trojan
• shadowsocks - Shadowsocks

💡 **Examples:**
/traffic john
/list
/create
        """
        await update.message.reply_text(help_text, parse_mode='Markdown')

    async def status_command(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        """Handle /status command"""
        if not self.is_authorized(update.effective_user.id):
            return
        
        try:
            response = self.api_request('GET', '/system/status')
            status = response['data']
            
            status_text = f"""
🖥️ **System Status**

🔐 SSH: {'✅' if status['ssh'] else '❌'}
🌐 Nginx: {'✅' if status['nginx'] else '❌'}
⚡ Xray: {'✅' if status['xray'] else '❌'}
🔒 Dropbear: {'✅' if status['dropbear'] else '❌'}
🔐 Stunnel: {'✅' if status['stunnel'] else '❌'}
🌐 SSH-WS: {'✅' if status['ssh_websocket'] else '❌'}
            """
            await update.message.reply_text(status_text, parse_mode='Markdown')
        except Exception as e:
            await update.message.reply_text(f'❌ Failed to get system status: {str(e)}')

    async def info_command(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        """Handle /info command"""
        if not self.is_authorized(update.effective_user.id):
            return
        
        try:
            response = self.api_request('GET', '/system/info')
            info = response['data']
            
            info_text = f"""
🖥️ **Server Information**

💻 OS: {info['os']}
🔧 Kernel: {info['kernel']}
⚡ CPU: {info['cpu_name']}
🧠 Cores: {info['cpu_cores']}
📊 CPU Usage: {info['cpu_usage']}
💾 RAM: {info['ram_used_mb']}MB / {info['ram_total_mb']}MB ({info['ram_usage_percent']})
⏰ Uptime: {info['uptime']}
🌐 Domain: {info['domain']}
📍 IP: {info['ip']}
📈 Daily Bandwidth: {info['daily_bandwidth']}
📊 Monthly Bandwidth: {info['monthly_bandwidth']}
            """
            await update.message.reply_text(info_text, parse_mode='Markdown')
        except Exception as e:
            await update.message.reply_text(f'❌ Failed to get server information: {str(e)}')

    async def create_command(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        """Handle /create command with inline keyboard"""
        if not self.is_authorized(update.effective_user.id):
            return
        
        keyboard = []
        for protocol in self.protocols:
            keyboard.append([InlineKeyboardButton(protocol.upper(), callback_data=f'create_{protocol}')])
        
        reply_markup = InlineKeyboardMarkup(keyboard)
        await update.message.reply_text(
            '🔧 **Create VPN User**\n\nSelect protocol:',
            reply_markup=reply_markup,
            parse_mode='Markdown'
        )

    async def list_command(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        """Handle /list command with inline keyboard"""
        if not self.is_authorized(update.effective_user.id):
            return
        
        keyboard = []
        for protocol in self.protocols:
            keyboard.append([InlineKeyboardButton(protocol.upper(), callback_data=f'list_{protocol}')])
        keyboard.append([InlineKeyboardButton('ALL USERS', callback_data='list_all')])
        
        reply_markup = InlineKeyboardMarkup(keyboard)
        await update.message.reply_text(
            '📋 **List Users**\n\nSelect protocol:',
            reply_markup=reply_markup,
            parse_mode='Markdown'
        )

    async def traffic_command(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        """Handle /traffic command"""
        if not self.is_authorized(update.effective_user.id):
            return
        
        if not context.args:
            await update.message.reply_text('❌ Usage: /traffic <username>\nExample: /traffic john')
            return
        
        username = context.args[0]
        
        try:
            response = self.api_request('GET', f'/vpn/users/{username}/traffic')
            traffic = response['data']
            
            traffic_text = f"""
📊 **Traffic Usage for {username}**

⬆️ Upload: {traffic['upload']}
⬇️ Download: {traffic['download']}
📈 Total: {traffic['total']}
            """
            await update.message.reply_text(traffic_text, parse_mode='Markdown')
        except Exception as e:
            await update.message.reply_text(f'❌ Failed to get traffic for user: {username}')

    async def button_callback(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        """Handle inline keyboard button callbacks"""
        query = update.callback_query
        await query.answer()
        
        if not self.is_authorized(query.from_user.id):
            return
        
        data = query.data
        
        if data.startswith('create_'):
            protocol = data.replace('create_', '')
            context.user_data['action'] = 'create'
            context.user_data['protocol'] = protocol
            
            await query.edit_message_text(
                f'🔧 **Creating {protocol.upper()} User**\n\n'
                'Please send username and days in format:\n'
                '`username days`\n\n'
                'Example: `john 30`',
                parse_mode='Markdown'
            )
        
        elif data.startswith('list_'):
            protocol = data.replace('list_', '')
            
            try:
                if protocol == 'all':
                    response = self.api_request('GET', '/vpn/users/all')
                    all_users = response['data']
                    
                    user_text = '📋 **All Users:**\n\n'
                    for proto, users in all_users.items():
                        user_text += f'**{proto.upper()}:**\n'
                        if users:
                            for i, user in enumerate(users, 1):
                                status = '✅' if user['is_active'] else '❌'
                                expiry = 'Never'
                                if user['expiry_date']:
                                    expiry = datetime.fromisoformat(user['expiry_date'].replace('Z', '+00:00')).strftime('%Y-%m-%d')
                                user_text += f'{i}. {status} {user["username"]} (Expires: {expiry})\n'
                        else:
                            user_text += 'No users found\n'
                        user_text += '\n'
                else:
                    response = self.api_request('GET', f'/vpn/{protocol}/users')
                    users = response['data']
                    
                    if not users:
                        user_text = f'📋 No {protocol.upper()} users found.'
                    else:
                        user_text = f'📋 **{protocol.upper()} Users:**\n\n'
                        for i, user in enumerate(users, 1):
                            status = '✅' if user['is_active'] else '❌'
                            expiry = 'Never'
                            if user['expiry_date']:
                                expiry = datetime.fromisoformat(user['expiry_date'].replace('Z', '+00:00')).strftime('%Y-%m-%d')
                            user_text += f'{i}. {status} {user["username"]} (Expires: {expiry})\n'
                
                await query.edit_message_text(user_text, parse_mode='Markdown')
            except Exception as e:
                await query.edit_message_text(f'❌ Failed to list users: {str(e)}')

    async def handle_text_message(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        """Handle text messages for user input"""
        if not self.is_authorized(update.effective_user.id):
            return
        
        user_data = context.user_data
        
        if user_data.get('action') == 'create':
            try:
                parts = update.message.text.strip().split()
                if len(parts) != 2:
                    await update.message.reply_text('❌ Invalid format. Use: `username days`', parse_mode='Markdown')
                    return
                
                username, days = parts[0], int(parts[1])
                protocol = user_data['protocol']
                
                # Generate random password for non-SSH protocols
                password = self.generate_password() if protocol == 'ssh' else None
                
                create_data = {
                    'username': username,
                    'days': days
                }
                if password:
                    create_data['password'] = password
                
                response = self.api_request('POST', f'/vpn/{protocol}/create', create_data)
                config = response['data']
                
                config_text = self.format_vpn_config(protocol, config)
                await update.message.reply_text(
                    f'✅ **{protocol.upper()} User Created Successfully!**\n\n{config_text}',
                    parse_mode='Markdown'
                )
                
                # Clear user data
                context.user_data.clear()
                
            except ValueError:
                await update.message.reply_text('❌ Invalid number of days. Please enter a valid number.')
            except Exception as e:
                await update.message.reply_text(f'❌ Failed to create user: {str(e)}')
                context.user_data.clear()

    def format_vpn_config(self, protocol: str, config: dict) -> str:
        """Format VPN configuration for display"""
        config_text = f'**Protocol:** {protocol.upper()}\n'
        config_text += f'**Server:** {config["server"]}\n'
        config_text += f'**Port:** {config["port"]}\n'
        config_text += f'**Username:** {config["username"]}\n'
        
        if config.get('password'):
            config_text += f'**Password:** `{config["password"]}`\n'
        
        if config.get('uuid'):
            config_text += f'**UUID:** `{config["uuid"]}`\n'
        
        if config.get('config'):
            config_text += '\n**Additional Config:**\n'
            for key, value in config['config'].items():
                config_text += f'**{key}:** `{value}`\n'
        
        return config_text

    def generate_password(self) -> str:
        """Generate random password"""
        import random
        import string
        return ''.join(random.choices(string.ascii_letters + string.digits, k=8))

async def main():
    """Main function to run the bot"""
    bot = VPNTelegramBot()
    
    # Login to API
    if not await bot.login_to_api():
        logger.error('Failed to login to API. Exiting...')
        return
    
    # Create application
    application = Application.builder().token(BOT_TOKEN).build()
    
    # Add handlers
    application.add_handler(CommandHandler('start', bot.start_command))
    application.add_handler(CommandHandler('help', bot.help_command))
    application.add_handler(CommandHandler('status', bot.status_command))
    application.add_handler(CommandHandler('info', bot.info_command))
    application.add_handler(CommandHandler('create', bot.create_command))
    application.add_handler(CommandHandler('list', bot.list_command))
    application.add_handler(CommandHandler('traffic', bot.traffic_command))
    application.add_handler(CallbackQueryHandler(bot.button_callback))
    application.add_handler(MessageHandler(filters.TEXT & ~filters.COMMAND, bot.handle_text_message))
    
    # Start the bot
    logger.info('🚀 Starting Telegram Bot...')
    await application.run_polling()

if __name__ == '__main__':
    asyncio.run(main())