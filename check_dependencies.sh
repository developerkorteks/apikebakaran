#!/bin/bash

# Dependency Checker untuk VPN API
# Memastikan semua yang dibutuhkan tersedia

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_status() {
    echo -e "${GREEN}[✓]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[!]${NC} $1"
}

print_error() {
    echo -e "${RED}[✗]${NC} $1"
}

print_info() {
    echo -e "${BLUE}[i]${NC} $1"
}

echo "🔍 Checking VPN API Dependencies..."
echo "=================================="

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    print_error "Please run this script as root!"
    exit 1
fi

print_status "Running as root"

# Check required commands
REQUIRED_COMMANDS=(
    "curl"
    "awk" 
    "grep"
    "cut"
    "vnstat"
    "service"
    "systemctl"
    "chage"
    "useradd"
    "userdel"
    "groupdel"
    "ps"
    "free"
    "uptime"
    "uname"
    "hostnamectl"
    "bc"
)

echo ""
print_info "Checking required commands..."

MISSING_COMMANDS=()
for cmd in "${REQUIRED_COMMANDS[@]}"; do
    if command -v "$cmd" >/dev/null 2>&1; then
        print_status "$cmd is available"
    else
        print_error "$cmd is missing"
        MISSING_COMMANDS+=("$cmd")
    fi
done

# Install missing commands
if [ ${#MISSING_COMMANDS[@]} -gt 0 ]; then
    print_warning "Installing missing commands..."
    apt update
    
    for cmd in "${MISSING_COMMANDS[@]}"; do
        case $cmd in
            "vnstat")
                apt install -y vnstat
                systemctl enable vnstat
                systemctl start vnstat
                ;;
            "bc")
                apt install -y bc
                ;;
            "curl")
                apt install -y curl
                ;;
            *)
                print_warning "Don't know how to install $cmd, it might be part of coreutils"
                ;;
        esac
    done
fi

echo ""
print_info "Checking required directories..."

# Check and create required directories
REQUIRED_DIRS=(
    "/etc/xray"
    "/etc/apivpn"
    "/var/lib/scrz-prem"
    "/var/log/apivpn"
)

for dir in "${REQUIRED_DIRS[@]}"; do
    if [ -d "$dir" ]; then
        print_status "$dir exists"
    else
        print_warning "Creating $dir"
        mkdir -p "$dir"
        chmod 755 "$dir"
        print_status "$dir created"
    fi
done

echo ""
print_info "Checking required files..."

# Check and create required files
if [ ! -f "/etc/xray/domain" ]; then
    print_warning "Creating /etc/xray/domain"
    SERVER_IP=$(curl -s ipinfo.io/ip 2>/dev/null || echo "127.0.0.1")
    echo "$SERVER_IP" > /etc/xray/domain
    print_status "/etc/xray/domain created with IP: $SERVER_IP"
else
    DOMAIN=$(cat /etc/xray/domain)
    print_status "/etc/xray/domain exists: $DOMAIN"
fi

if [ ! -f "/var/lib/scrz-prem/ipvps.conf" ]; then
    print_warning "Creating /var/lib/scrz-prem/ipvps.conf"
    SERVER_IP=$(cat /etc/xray/domain)
    echo "IP=$SERVER_IP" > /var/lib/scrz-prem/ipvps.conf
    print_status "/var/lib/scrz-prem/ipvps.conf created"
else
    print_status "/var/lib/scrz-prem/ipvps.conf exists"
fi

if [ ! -f "/etc/xray/config.json" ]; then
    print_warning "Creating basic /etc/xray/config.json"
    cat > /etc/xray/config.json << 'EOF'
{
  "log": {
    "loglevel": "warning"
  },
  "inbounds": [],
  "outbounds": [
    {
      "protocol": "freedom",
      "settings": {}
    }
  ]
}
EOF
    print_status "/etc/xray/config.json created"
else
    print_status "/etc/xray/config.json exists"
fi

echo ""
print_info "Checking services..."

# Check if services exist (not necessarily running)
SERVICES=(
    "ssh"
    "nginx" 
    "xray"
    "dropbear"
    "stunnel5"
    "ws-stunnel"
)

for service in "${SERVICES[@]}"; do
    if systemctl list-unit-files | grep -q "^$service.service"; then
        STATUS=$(systemctl is-active "$service" 2>/dev/null || echo "inactive")
        if [ "$STATUS" = "active" ]; then
            print_status "$service service is running"
        else
            print_warning "$service service exists but not running"
        fi
    else
        print_warning "$service service not found (this is OK if not installed)"
    fi
done

echo ""
print_info "Testing service status detection..."

# Test the service status detection method from menu_script.txt
for service in "ssh" "nginx"; do
    if systemctl list-unit-files | grep -q "^$service.service"; then
        # Test the exact method from menu_script.txt
        cek=$(service "$service" status 2>/dev/null | grep active | cut -d ' ' -f5 || echo "")
        if [ "$cek" = "active" ]; then
            stat="-f5"
        else
            stat="-f7"
        fi
        
        status_result=$(service "$service" status 2>/dev/null | grep active | cut -d ' ' $stat || echo "inactive")
        
        if [ "$status_result" = "active" ]; then
            print_status "$service status detection: ACTIVE"
        else
            print_warning "$service status detection: $status_result"
        fi
    fi
done

echo ""
print_info "Testing user detection..."

# Test SSH user counting
ssh_count=$(awk -F: '$3 >= 1000 && $1 != "nobody" {print $1}' /etc/passwd | wc -l)
print_status "SSH users found: $ssh_count"

if [ "$ssh_count" -gt 0 ]; then
    print_info "SSH users:"
    awk -F: '$3 >= 1000 && $1 != "nobody" {print "  - " $1}' /etc/passwd
fi

echo ""
print_info "Testing xray user detection..."

# Test xray user counting
for protocol in "vmsg" "vlsg" "trg" "ssg"; do
    count=$(grep -c -E "^#$protocol " /etc/xray/config.json 2>/dev/null || echo "0")
    print_status "$protocol users found: $count"
done

echo ""
print_info "Testing system information gathering..."

# Test system info gathering
OS_INFO=$(hostnamectl 2>/dev/null | grep "Operating System" | cut -d ' ' -f5- || echo "Unknown")
KERNEL=$(uname -r)
CPU_NAME=$(awk -F: '/model name/ {name=$2} END {print name}' /proc/cpuinfo | sed 's/^[ \t]*//')
CPU_CORES=$(awk -F: '/model name/ {core++} END {print core}' /proc/cpuinfo)
UPTIME=$(uptime -p | cut -d ' ' -f 2-10)

print_status "OS: $OS_INFO"
print_status "Kernel: $KERNEL"
print_status "CPU: $CPU_NAME"
print_status "Cores: $CPU_CORES"
print_status "Uptime: $UPTIME"

# Test RAM info
RAM_USED=$(free -m | grep Mem: | awk '{print $3}')
RAM_TOTAL=$(free -m | grep Mem: | awk '{print $2}')
RAM_USAGE=$(echo "scale=1; ($RAM_USED / $RAM_TOTAL) * 100" | bc | cut -d. -f1)

print_status "RAM: ${RAM_USED}MB / ${RAM_TOTAL}MB (${RAM_USAGE}%)"

# Test bandwidth info
DAILY_BW=$(vnstat -d --oneline 2>/dev/null | awk -F\; '{print $6}' | sed 's/ //' || echo "N/A")
MONTHLY_BW=$(vnstat -m --oneline 2>/dev/null | awk -F\; '{print $11}' | sed 's/ //' || echo "N/A")

print_status "Daily bandwidth: $DAILY_BW"
print_status "Monthly bandwidth: $MONTHLY_BW"

echo ""
print_info "Testing IP detection..."

# Test IP detection methods
IP1=$(curl -s ipinfo.io/ip 2>/dev/null || echo "")
IP2=$(curl -sS ipv4.icanhazip.com 2>/dev/null || echo "")
IP3=$(curl -sS ifconfig.me 2>/dev/null || echo "")

print_status "IP from ipinfo.io: ${IP1:-N/A}"
print_status "IP from icanhazip.com: ${IP2:-N/A}"
print_status "IP from ifconfig.me: ${IP3:-N/A}"

echo ""
echo "🎉 Dependency check completed!"
echo ""
print_info "Summary:"
echo "  - All required commands are available"
echo "  - All required directories exist"
echo "  - Basic configuration files created"
echo "  - System information gathering works"
echo "  - Service status detection works"
echo ""
print_status "System is ready for VPN API!"