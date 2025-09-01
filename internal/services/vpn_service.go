package services

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nabilulilalbab/apivpn/internal/models"
)

type VPNService struct{}

func NewVPNService() *VPNService {
	return &VPNService{}
}

// SSH User Management
func (v *VPNService) CreateSSHUser(req *models.CreateUserRequest) (*models.VPNConfig, error) {
	expiryDate := time.Now().AddDate(0, 0, req.Days)
	
	// Create system user
	commands := []string{
		fmt.Sprintf("useradd -e %s -s /bin/false -M %s", expiryDate.Format("2006-01-02"), req.Username),
		fmt.Sprintf("echo '%s:%s' | chpasswd", req.Username, req.Password),
	}

	for _, cmd := range commands {
		if err := v.executeCommand(cmd); err != nil {
			return nil, fmt.Errorf("failed to create SSH user: %v", err)
		}
	}

	// Get server domain/IP
	domainCmd := exec.Command("bash", "-c", "cat /etc/xray/domain 2>/dev/null || curl -s ipinfo.io/ip")
	domainOutput, _ := domainCmd.Output()
	domain := strings.TrimSpace(string(domainOutput))

	config := &models.VPNConfig{
		Protocol: "ssh",
		Server:   domain,
		Port:     22,
		Username: req.Username,
		Password: req.Password,
		Config: map[string]string{
			"ssh_port":    "22",
			"ssl_port":    "443",
			"ws_port":     "80",
			"stunnel_port": "444",
		},
	}

	return config, nil
}

func (v *VPNService) GetSSHUsers() ([]models.User, error) {
	// Use the exact same command as menu_script.txt
	output, err := v.executeCommandWithOutput("awk -F: '$3 >= 1000 && $1 != \"nobody\" {print $1}' /etc/passwd | wc -l")
	if err != nil {
		return nil, err
	}

	count, _ := strconv.Atoi(strings.TrimSpace(output))
	
	// Get actual usernames
	usernamesOutput, err := v.executeCommandWithOutput("awk -F: '$3 >= 1000 && $1 != \"nobody\" {print $1}' /etc/passwd")
	if err != nil {
		return nil, err
	}

	var users []models.User
	if count > 0 {
		usernames := strings.Split(strings.TrimSpace(usernamesOutput), "\n")
		
		for _, username := range usernames {
			if username == "" || username == "nobody" {
				continue
			}
			
			// Get user expiry date using chage
			expiryCmd := exec.Command("bash", "-c", fmt.Sprintf("chage -l %s 2>/dev/null | grep 'Account expires' | cut -d: -f2", username))
			expiryOutput, _ := expiryCmd.Output()
			expiry := strings.TrimSpace(string(expiryOutput))
			
			var expiryDate time.Time
			var isActive bool = true
			
			if expiry != "never" && expiry != "" && expiry != "Account expires" {
				if parsedDate, err := time.Parse("Jan 02, 2006", strings.TrimSpace(expiry)); err == nil {
					expiryDate = parsedDate
					isActive = expiryDate.After(time.Now())
				} else if parsedDate, err := time.Parse("2006-01-02", strings.TrimSpace(expiry)); err == nil {
					expiryDate = parsedDate
					isActive = expiryDate.After(time.Now())
				}
			}

			users = append(users, models.User{
				Username:    username,
				Protocol:    "ssh",
				ExpiryDate:  expiryDate,
				IsActive:    isActive,
				CreatedDate: time.Now(),
			})
		}
	}

	return users, nil
}

func (v *VPNService) DeleteSSHUser(username string) error {
	commands := []string{
		fmt.Sprintf("userdel -f %s", username),
		fmt.Sprintf("groupdel %s 2>/dev/null || true", username),
	}

	for _, cmd := range commands {
		v.executeCommand(cmd) // Ignore errors for cleanup
	}

	return nil
}

// VMESS User Management
func (v *VPNService) CreateVmessUser(req *models.CreateUserRequest) (*models.VPNConfig, error) {
	userUUID := uuid.New().String()
	expiryDate := time.Now().AddDate(0, 0, req.Days)

	// Add user to xray config
	if err := v.addXrayUser("vmess", req.Username, userUUID, expiryDate); err != nil {
		return nil, err
	}

	domainCmd := exec.Command("bash", "-c", "cat /etc/xray/domain 2>/dev/null || curl -s ipinfo.io/ip")
	domainOutput, _ := domainCmd.Output()
	domain := strings.TrimSpace(string(domainOutput))

	// Get vmess port from log-install.txt like original script
	portCmd := exec.Command("bash", "-c", "cat ~/log-install.txt | grep -w 'Vmess TLS' | cut -d: -f2 | sed 's/ //g' 2>/dev/null || echo '443'")
	portOutput, _ := portCmd.Output()
	port := strings.TrimSpace(string(portOutput))
	if port == "" {
		port = "443"
	}

	// Generate links exactly like original script
	vmessWS := fmt.Sprintf("vmess://%s@%s:%s?type=ws&encryption=none&security=tls&host=%s&path=/vmess&allowInsecure=1&sni=%s#VMESS_WS_%s", 
		userUUID, domain, port, domain, domain, req.Username)
	vmessGRPC := fmt.Sprintf("vmess://%s@%s:%s?mode=gun&security=tls&encryption=none&type=grpc&serviceName=vmess-grpc&sni=%s#VMESS_GRPC_%s", 
		userUUID, domain, port, domain, req.Username)
	configURL := fmt.Sprintf("http://%s:81/vmess-%s.txt", domain, req.Username)

	config := &models.VPNConfig{
		Protocol: "vmess",
		Server:   domain,
		Port:     443,
		Username: req.Username,
		UUID:     userUUID,
		Config: map[string]string{
			"remarks":     req.Username,
			"host":        domain,
			"port":        port,
			"uuid":        userUUID,
			"alterId":     "0",
			"security":    "auto",
			"network":     "ws/grpc",
			"path":        "/vmess",
			"serviceName": "vmess-grpc",
			"link_ws":     vmessWS,
			"link_grpc":   vmessGRPC,
			"config_url":  configURL,
			"expired_on":  expiryDate.Format("2006-01-02"),
		},
	}

	return config, nil
}

func (v *VPNService) GetVmessUsers() ([]models.User, error) {
	return v.getXrayUsers("vmess", "#vmsg")
}

func (v *VPNService) DeleteVmessUser(username string) error {
	return v.deleteXrayUser("vmess", username, "#vmsg")
}

// VLESS User Management
func (v *VPNService) CreateVlessUser(req *models.CreateUserRequest) (*models.VPNConfig, error) {
	userUUID := uuid.New().String()
	expiryDate := time.Now().AddDate(0, 0, req.Days)

	if err := v.addXrayUser("vless", req.Username, userUUID, expiryDate); err != nil {
		return nil, err
	}

	domainCmd := exec.Command("bash", "-c", "cat /etc/xray/domain 2>/dev/null || curl -s ipinfo.io/ip")
	domainOutput, _ := domainCmd.Output()
	domain := strings.TrimSpace(string(domainOutput))

	// Get vless ports from log-install.txt like original script
	tlsPortCmd := exec.Command("bash", "-c", "cat ~/log-install.txt | grep -w 'Vless TLS' | cut -d: -f2 | sed 's/ //g' 2>/dev/null || echo '443'")
	tlsPortOutput, _ := tlsPortCmd.Output()
	tlsPort := strings.TrimSpace(string(tlsPortOutput))
	if tlsPort == "" {
		tlsPort = "443"
	}

	nonePortCmd := exec.Command("bash", "-c", "cat ~/log-install.txt | grep -w 'Vless None TLS' | cut -d: -f2 | sed 's/ //g' 2>/dev/null || echo '80'")
	nonePortOutput, _ := nonePortCmd.Output()
	nonePort := strings.TrimSpace(string(nonePortOutput))
	if nonePort == "" {
		nonePort = "80"
	}

	// Generate links exactly like original script
	vlessTLS := fmt.Sprintf("vless://%s@%s:%s?type=ws&encryption=none&security=tls&host=%s&path=/vless&allowInsecure=1&sni=%s#XRAY_VLESS_TLS_%s", 
		userUUID, domain, tlsPort, domain, domain, req.Username)
	vlessNTLS := fmt.Sprintf("vless://%s@%s:%s?type=ws&encryption=none&security=none&host=%s&path=/vless#XRAY_VLESS_NTLS_%s", 
		userUUID, domain, nonePort, domain, req.Username)
	vlessGRPC := fmt.Sprintf("vless://%s@%s:%s?mode=gun&security=tls&encryption=none&type=grpc&serviceName=vless-grpc&sni=%s#VLESS_GRPC_%s", 
		userUUID, domain, tlsPort, domain, req.Username)
	configURL := fmt.Sprintf("http://%s:81/vless-%s.txt", domain, req.Username)

	config := &models.VPNConfig{
		Protocol: "vless",
		Server:   domain,
		Port:     443,
		Username: req.Username,
		UUID:     userUUID,
		Config: map[string]string{
			"remarks":     req.Username,
			"host":        domain,
			"port_tls":    tlsPort,
			"port_ntls":   nonePort,
			"uuid":        userUUID,
			"encryption":  "none",
			"network":     "ws/grpc",
			"path":        "/vless",
			"serviceName": "vless-grpc",
			"link_tls":    vlessTLS,
			"link_ntls":   vlessNTLS,
			"link_grpc":   vlessGRPC,
			"config_url":  configURL,
			"expired_on":  expiryDate.Format("2006-01-02"),
		},
	}

	return config, nil
}

func (v *VPNService) GetVlessUsers() ([]models.User, error) {
	return v.getXrayUsers("vless", "#vlsg")
}

func (v *VPNService) DeleteVlessUser(username string) error {
	return v.deleteXrayUser("vless", username, "#vlsg")
}

// Trojan User Management
func (v *VPNService) CreateTrojanUser(req *models.CreateUserRequest) (*models.VPNConfig, error) {
	userUUID := uuid.New().String()
	expiryDate := time.Now().AddDate(0, 0, req.Days)

	if err := v.addXrayUser("trojan", req.Username, userUUID, expiryDate); err != nil {
		return nil, err
	}

	// Get domain and port info like original script
	domainCmd := exec.Command("bash", "-c", "cat /etc/xray/domain 2>/dev/null || curl -s ipinfo.io/ip")
	domainOutput, _ := domainCmd.Output()
	domain := strings.TrimSpace(string(domainOutput))
	
	// Get trojan port from log-install.txt like original script
	portCmd := exec.Command("bash", "-c", "cat ~/log-install.txt | grep -w 'Trojan WS ' | cut -d: -f2 | sed 's/ //g' 2>/dev/null || echo '443'")
	portOutput, _ := portCmd.Output()
	port := strings.TrimSpace(string(portOutput))
	if port == "" {
		port = "443"
	}

	// Generate links exactly like original script
	trojanWS := fmt.Sprintf("trojan://%s@%s:%s?path=%%2Ftrojan-ws&security=tls&host=%s&type=ws&sni=%s#TROJAN_WS_%s", 
		userUUID, domain, port, domain, domain, req.Username)
	trojanGO := fmt.Sprintf("trojan-go://%s@%s:%s?path=%%2Ftrojan-ws&security=tls&host=%s&type=ws&sni=%s#TROJANGO_%s", 
		userUUID, domain, port, domain, domain, req.Username)
	trojanGRPC := fmt.Sprintf("trojan://%s@%s:%s?mode=gun&security=tls&type=grpc&serviceName=trojan-grpc&sni=%s#TROJAN_GRPC_%s", 
		userUUID, domain, port, domain, req.Username)
	configURL := fmt.Sprintf("http://%s:81/trojan-%s.txt", domain, req.Username)

	config := &models.VPNConfig{
		Protocol: "trojan",
		Server:   domain,
		Port:     443,
		Username: req.Username,
		Password: userUUID,
		Config: map[string]string{
			"remarks":     req.Username,
			"host":        domain,
			"port":        port,
			"key":         userUUID,
			"network":     "ws/grpc",
			"path":        "/trojan-ws",
			"serviceName": "trojan-grpc",
			"link_ws":     trojanWS,
			"link_go":     trojanGO,
			"link_grpc":   trojanGRPC,
			"config_url":  configURL,
			"expired_on":  expiryDate.Format("2006-01-02"),
		},
	}

	return config, nil
}

func (v *VPNService) GetTrojanUsers() ([]models.User, error) {
	return v.getXrayUsers("trojan", "#trg")
}

func (v *VPNService) DeleteTrojanUser(username string) error {
	return v.deleteXrayUser("trojan", username, "#trg")
}

// Shadowsocks User Management
func (v *VPNService) CreateShadowsocksUser(req *models.CreateUserRequest) (*models.VPNConfig, error) {
	userUUID := uuid.New().String()
	expiryDate := time.Now().AddDate(0, 0, req.Days)

	if err := v.addXrayUser("shadowsocks", req.Username, userUUID, expiryDate); err != nil {
		return nil, err
	}

	domainCmd := exec.Command("bash", "-c", "cat /etc/xray/domain 2>/dev/null || curl -s ipinfo.io/ip")
	domainOutput, _ := domainCmd.Output()
	domain := strings.TrimSpace(string(domainOutput))

	// Get shadowsocks port from log-install.txt like original script
	portCmd := exec.Command("bash", "-c", "cat ~/log-install.txt | grep -w 'Shadowsocks WS' | cut -d: -f2 | sed 's/ //g' 2>/dev/null || echo '443'")
	portOutput, _ := portCmd.Output()
	port := strings.TrimSpace(string(portOutput))
	if port == "" {
		port = "443"
	}

	// Generate links exactly like original script
	ssWS := fmt.Sprintf("ss://%s@%s:%s?type=ws&security=tls&host=%s&path=/ss&sni=%s#SS_WS_%s", 
		userUUID, domain, port, domain, domain, req.Username)
	ssGRPC := fmt.Sprintf("ss://%s@%s:%s?mode=gun&security=tls&type=grpc&serviceName=ss-grpc&sni=%s#SS_GRPC_%s", 
		userUUID, domain, port, domain, req.Username)
	configURL := fmt.Sprintf("http://%s:81/ss-%s.txt", domain, req.Username)

	config := &models.VPNConfig{
		Protocol: "shadowsocks",
		Server:   domain,
		Port:     443,
		Username: req.Username,
		Password: userUUID,
		Config: map[string]string{
			"remarks":     req.Username,
			"host":        domain,
			"port":        port,
			"password":    userUUID,
			"method":      "aes-256-gcm",
			"network":     "ws/grpc",
			"path":        "/ss",
			"serviceName": "ss-grpc",
			"link_ws":     ssWS,
			"link_grpc":   ssGRPC,
			"config_url":  configURL,
			"expired_on":  expiryDate.Format("2006-01-02"),
		},
	}

	return config, nil
}

func (v *VPNService) GetShadowsocksUsers() ([]models.User, error) {
	return v.getXrayUsers("shadowsocks", "#ssg")
}

func (v *VPNService) DeleteShadowsocksUser(username string) error {
	return v.deleteXrayUser("shadowsocks", username, "#ssg")
}

// Extension methods
func (v *VPNService) ExtendUser(protocol, username string, days int) error {
	switch protocol {
	case "ssh":
		expiryDate := time.Now().AddDate(0, 0, days)
		return v.executeCommand(fmt.Sprintf("chage -E %s %s", expiryDate.Format("2006-01-02"), username))
	case "vmess", "vless", "trojan", "shadowsocks":
		// For xray users, we need to update the expiry in our tracking system
		// This would typically be stored in a database or file
		return v.updateXrayUserExpiry(username, days)
	}
	return fmt.Errorf("unsupported protocol: %s", protocol)
}

// Get user traffic
func (v *VPNService) GetUserTraffic(username string) (*models.UserTraffic, error) {
	// This would typically query vnstat or iptables for user-specific traffic
	// For now, return placeholder data
	return &models.UserTraffic{
		Username: username,
		Upload:   "0 MB",
		Download: "0 MB",
		Total:    "0 MB",
	}, nil
}

// Cleanup expired users
func (v *VPNService) CleanupExpiredUsers() error {
	// SSH users
	sshUsers, _ := v.GetSSHUsers()
	for _, user := range sshUsers {
		if !user.IsActive && !user.ExpiryDate.IsZero() {
			v.DeleteSSHUser(user.Username)
		}
	}

	// Xray users would need similar cleanup
	protocols := []string{"vmess", "vless", "trojan", "shadowsocks"}
	for _, protocol := range protocols {
		users, _ := v.getXrayUsers(protocol, v.getXrayPrefix(protocol))
		for _, user := range users {
			if !user.IsActive && !user.ExpiryDate.IsZero() {
				v.deleteXrayUser(protocol, user.Username, v.getXrayPrefix(protocol))
			}
		}
	}

	return nil
}

// Helper methods
func (v *VPNService) executeCommand(command string) error {
	cmd := exec.Command("bash", "-c", command)
	return cmd.Run()
}

func (v *VPNService) executeCommandWithOutput(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.Output()
	return string(output), err
}

func (v *VPNService) addXrayUser(protocol, username, uuid string, expiry time.Time) error {
	// Implement the exact same logic as the original scripts but non-interactively
	expiryStr := expiry.Format("2006-01-02")
	
	// Check if user already exists
	checkCmd := fmt.Sprintf("grep -w %s /etc/xray/config.json | wc -l", username)
	output, err := v.executeCommandWithOutput(checkCmd)
	if err != nil {
		return fmt.Errorf("failed to check existing user: %v", err)
	}
	
	if strings.TrimSpace(output) != "0" {
		return fmt.Errorf("user %s already exists", username)
	}
	
	// Generate UUID if not provided
	if uuid == "" {
		uuidCmd := "cat /proc/sys/kernel/random/uuid"
		uuidOutput, err := v.executeCommandWithOutput(uuidCmd)
		if err != nil {
			return fmt.Errorf("failed to generate UUID: %v", err)
		}
		uuid = strings.TrimSpace(uuidOutput)
	}
	
	var commands []string
	switch protocol {
	case "vmess":
		// Based on add-ws script
		commands = []string{
			fmt.Sprintf(`sed -i '/#vmess$/a\#vms %s %s\
},{"id": "%s","email": "%s"' /etc/xray/config.json`, username, expiryStr, uuid, username),
			fmt.Sprintf(`sed -i '/#vmessgrpc$/a\#vmsg %s %s\
},{"id": "%s","email": "%s"' /etc/xray/config.json`, username, expiryStr, uuid, username),
		}
		
	case "vless":
		// Based on add-vless script
		commands = []string{
			fmt.Sprintf(`sed -i '/#vless$/a\#vls %s %s\
},{"id": "%s","email": "%s"' /etc/xray/config.json`, username, expiryStr, uuid, username),
			fmt.Sprintf(`sed -i '/#vlessgrpc$/a\#vlsg %s %s\
},{"id": "%s","email": "%s"' /etc/xray/config.json`, username, expiryStr, uuid, username),
		}
		
	case "trojan":
		// Based on add-tr script
		commands = []string{
			fmt.Sprintf(`sed -i '/#trojanws$/a\#tr %s %s\
},{"password": "%s","email": "%s"' /etc/xray/config.json`, username, expiryStr, uuid, username),
			fmt.Sprintf(`sed -i '/#trojangrpc$/a\#trg %s %s\
},{"password": "%s","email": "%s"' /etc/xray/config.json`, username, expiryStr, uuid, username),
		}
		
	case "shadowsocks":
		// Based on add-ssws script
		commands = []string{
			fmt.Sprintf(`sed -i '/#shadowsocks$/a\#ss %s %s\
},{"password": "%s","email": "%s"' /etc/xray/config.json`, username, expiryStr, uuid, username),
			fmt.Sprintf(`sed -i '/#shadowsocksgrpc$/a\#ssg %s %s\
},{"password": "%s","email": "%s"' /etc/xray/config.json`, username, expiryStr, uuid, username),
		}
		
	default:
		return fmt.Errorf("unsupported protocol: %s", protocol)
	}
	
	// Execute sed commands
	for _, cmd := range commands {
		if err := v.executeCommand(cmd); err != nil {
			return fmt.Errorf("failed to add %s user to config: %v", protocol, err)
		}
	}
	
	// Restart xray service
	if err := v.executeCommand("systemctl restart xray"); err != nil {
		return fmt.Errorf("failed to restart xray service: %v", err)
	}
	
	return nil
}

func (v *VPNService) addXrayUserManual(protocol, username, uuid string, expiry time.Time) error {
	// Simple manual approach - just add comment line to track user
	configPath := "/etc/xray/config.json"
	expiryStr := expiry.Format("2006-01-02")
	
	var commentPrefix string
	switch protocol {
	case "vmess":
		commentPrefix = "#vmsg"
	case "vless":
		commentPrefix = "#vlsg"
	case "trojan":
		commentPrefix = "#trg"
	case "shadowsocks":
		commentPrefix = "#ssg"
	}
	
	// Add comment line to track the user
	cmd := fmt.Sprintf(`echo "%s %s %s" >> %s`, commentPrefix, username, expiryStr, configPath)
	return v.executeCommand(cmd)
}

func (v *VPNService) getXrayUsers(protocol, prefix string) ([]models.User, error) {
	// Use the exact same pattern as menu_script.txt but avoid JSON parsing errors
	var grepPattern string
	switch protocol {
	case "vmess":
		grepPattern = "grep -c -E \"^#vmsg \" /etc/xray/config.json 2>/dev/null || echo 0"
	case "vless":
		grepPattern = "grep -c -E \"^#vlsg \" /etc/xray/config.json 2>/dev/null || echo 0"
	case "trojan":
		grepPattern = "grep -c -E \"^#trg \" /etc/xray/config.json 2>/dev/null || echo 0"
	case "shadowsocks":
		grepPattern = "grep -c -E \"^#ssg \" /etc/xray/config.json 2>/dev/null || echo 0"
	default:
		return []models.User{}, nil
	}

	output, err := v.executeCommandWithOutput(grepPattern)
	if err != nil {
		return []models.User{}, nil
	}

	count, _ := strconv.Atoi(strings.TrimSpace(output))
	
	// Get actual usernames from config file
	var users []models.User
	if count > 0 {
		// Extract usernames from xray config, avoid JSON parsing
		usernamesOutput, err := v.executeCommandWithOutput(fmt.Sprintf("grep -E \"^%s \" /etc/xray/config.json 2>/dev/null | awk '{print $2}' || true", prefix))
		if err == nil && usernamesOutput != "" {
			usernames := strings.Split(strings.TrimSpace(usernamesOutput), "\n")
			for _, username := range usernames {
				if username != "" && strings.TrimSpace(username) != "" {
					users = append(users, models.User{
						Username: strings.TrimSpace(username),
						Protocol: protocol,
						IsActive: true,
						ExpiryDate: time.Now().AddDate(0, 1, 0), // Default 1 month
					})
				}
			}
		}
	}

	return users, nil
}

func (v *VPNService) deleteXrayUser(protocol, username, prefix string) error {
	// Remove user from xray config and restart service
	// This is simplified - you'd need to properly modify the JSON config
	return v.executeCommand("systemctl restart xray")
}

func (v *VPNService) updateXrayUserExpiry(username string, days int) error {
	// Update user expiry in tracking system
	// This would typically update a database or file
	return nil
}

func (v *VPNService) getXrayPrefix(protocol string) string {
	prefixes := map[string]string{
		"vmess":       "#vmsg",
		"vless":       "#vlsg",
		"trojan":      "#trg",
		"shadowsocks": "#ssg",
	}
	return prefixes[protocol]
}