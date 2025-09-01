package services

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/nabilulilalbab/apivpn/internal/models"
)

type SystemService struct{}

func NewSystemService() *SystemService {
	return &SystemService{}
}

// GetSystemInfo retrieves system information similar to the bash script
func (s *SystemService) GetSystemInfo() (*models.SystemInfo, error) {
	info := &models.SystemInfo{}

	// Get OS information
	if osInfo, err := s.executeCommand("hostnamectl | grep 'Operating System' | cut -d ' ' -f5-"); err == nil {
		info.OS = strings.TrimSpace(osInfo)
	}

	// Get kernel version
	if kernel, err := s.executeCommand("uname -r"); err == nil {
		info.Kernel = strings.TrimSpace(kernel)
	}

	// Get CPU information
	if cpuName, err := s.getCPUName(); err == nil {
		info.CPUName = cpuName
	}

	if cores, err := s.getCPUCores(); err == nil {
		info.CPUCores = cores
	}

	if usage, err := s.getCPUUsage(); err == nil {
		info.CPUUsage = usage
	}

	// Get RAM information
	if ramUsed, ramTotal, ramUsage, err := s.getRAMInfo(); err == nil {
		info.RAMUsed = ramUsed
		info.RAMTotal = ramTotal
		info.RAMUsage = ramUsage
	}

	// Get uptime
	if uptime, err := s.executeCommand("uptime -p | cut -d ' ' -f 2-10"); err == nil {
		info.Uptime = strings.TrimSpace(uptime)
	}

	// Get domain - exactly like menu_script.txt
	if domain, err := s.executeCommand("cat /etc/xray/domain 2>/dev/null"); err == nil {
		info.Domain = strings.TrimSpace(domain)
	} else {
		info.Domain = "Not configured"
	}

	// Get IP - use multiple methods like menu_script.txt
	var ip string
	if ipOutput, err := s.executeCommand("curl -s ipinfo.io/ip"); err == nil {
		ip = strings.TrimSpace(ipOutput)
	}
	if ip == "" {
		if ipOutput, err := s.executeCommand("curl -sS ipv4.icanhazip.com"); err == nil {
			ip = strings.TrimSpace(ipOutput)
		}
	}
	if ip == "" {
		if ipOutput, err := s.executeCommand("curl -sS ifconfig.me"); err == nil {
			ip = strings.TrimSpace(ipOutput)
		}
	}
	info.IP = ip

	// Get bandwidth usage
	if daily, err := s.executeCommand("vnstat -d --oneline | awk -F\\; '{print $6}' | sed 's/ //'"); err == nil {
		info.DailyBandwidth = strings.TrimSpace(daily)
	}

	if monthly, err := s.executeCommand("vnstat -m --oneline | awk -F\\; '{print $11}' | sed 's/ //'"); err == nil {
		info.MonthlyBandwidth = strings.TrimSpace(monthly)
	}

	return info, nil
}

// GetServiceStatus checks the status of VPN services
func (s *SystemService) GetServiceStatus() (*models.ServiceStatus, error) {
	status := &models.ServiceStatus{}

	// Use the same method as menu_script.txt
	status.SSH = s.isServiceActiveScript("ssh")
	status.Nginx = s.isServiceActiveScript("nginx")
	status.Xray = s.isServiceActiveScript("xray")
	status.Dropbear = s.isServiceActiveScript("dropbear")
	status.Stunnel = s.isServiceActiveScript("stunnel5")
	status.SSHWebSocket = s.isServiceActiveScript("ws-stunnel")

	return status, nil
}

// AddDomain adds a new domain to the system
func (s *SystemService) AddDomain(domain string) error {
	// Remove existing domain file
	os.Remove("/etc/xray/domain")

	// Write new domain to config
	if err := s.writeToFile("/var/lib/scrz-prem/ipvps.conf", fmt.Sprintf("IP=%s", domain)); err != nil {
		return err
	}

	// Write domain to xray config
	if err := s.writeToFile("/etc/xray/domain", domain); err != nil {
		return err
	}

	return nil
}

// RenewSSL renews SSL certificate
func (s *SystemService) RenewSSL() error {
	commands := []string{
		"systemctl stop nginx",
		"systemctl stop xray",
		"/root/.acme.sh/acme.sh --upgrade",
		"/root/.acme.sh/acme.sh --upgrade --auto-upgrade",
		"/root/.acme.sh/acme.sh --set-default-ca --server letsencrypt",
	}

	domain, err := s.executeCommand("cat /var/lib/scrz-prem/ipvps.conf | cut -d'=' -f2")
	if err != nil {
		return err
	}
	domain = strings.TrimSpace(domain)

	// Issue new certificate
	issueCmd := fmt.Sprintf("/root/.acme.sh/acme.sh --issue -d %s --standalone -k ec-256", domain)
	commands = append(commands, issueCmd)

	// Install certificate
	installCmd := fmt.Sprintf("~/.acme.sh/acme.sh --installcert -d %s --fullchainpath /etc/xray/xray.crt --keypath /etc/xray/xray.key --ecc", domain)
	commands = append(commands, installCmd)

	// Restart services
	commands = append(commands, []string{
		"systemctl start nginx",
		"systemctl start xray",
	}...)

	for _, cmd := range commands {
		if _, err := s.executeCommand(cmd); err != nil {
			return fmt.Errorf("failed to execute command '%s': %v", cmd, err)
		}
	}

	return nil
}

// Reboot system
func (s *SystemService) Reboot() error {
	_, err := s.executeCommand("reboot")
	return err
}

// RestartServices restarts VPN services
func (s *SystemService) RestartServices() error {
	services := []string{"ssh", "nginx", "xray", "dropbear", "stunnel5", "ws-stunnel"}
	
	for _, service := range services {
		if _, err := s.executeCommand(fmt.Sprintf("systemctl restart %s", service)); err != nil {
			return fmt.Errorf("failed to restart %s: %v", service, err)
		}
	}

	return nil
}

// Helper methods
func (s *SystemService) executeCommand(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.Output()
	return string(output), err
}

func (s *SystemService) isServiceActive(service string) bool {
	output, err := s.executeCommand(fmt.Sprintf("systemctl is-active %s", service))
	return err == nil && strings.TrimSpace(output) == "active"
}

// isServiceActiveScript uses the same method as menu_script.txt
func (s *SystemService) isServiceActiveScript(service string) bool {
	// First check to determine the field position
	cekOutput, err := s.executeCommand(fmt.Sprintf("service %s status | grep active | cut -d ' ' -f5", service))
	if err != nil {
		return false
	}
	
	var stat string
	if strings.TrimSpace(cekOutput) == "active" {
		stat = "-f5"
	} else {
		stat = "-f7"
	}
	
	// Get the actual status
	statusOutput, err := s.executeCommand(fmt.Sprintf("service %s status | grep active | cut -d ' ' %s", service, stat))
	if err != nil {
		return false
	}
	
	return strings.TrimSpace(statusOutput) == "active"
}

func (s *SystemService) getCPUName() (string, error) {
	output, err := s.executeCommand("awk -F: '/model name/ {name=$2} END {print name}' /proc/cpuinfo")
	return strings.TrimSpace(output), err
}

func (s *SystemService) getCPUCores() (int, error) {
	output, err := s.executeCommand("awk -F: '/model name/ {core++} END {print core}' /proc/cpuinfo")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(output))
}

func (s *SystemService) getCPUUsage() (string, error) {
	// Use the exact same calculation as menu_script.txt
	cpu_usage1_output, err := s.executeCommand("ps aux | awk 'BEGIN {sum=0} {sum+=$3}; END {print sum}'")
	if err != nil {
		return "", err
	}
	
	cores, err := s.getCPUCores()
	if err != nil || cores == 0 {
		cores = 1
	}
	
	cpu_usage1 := strings.TrimSpace(cpu_usage1_output)
	// Remove decimal part and divide by cores (like menu_script.txt)
	if dotIndex := strings.Index(cpu_usage1, "."); dotIndex != -1 {
		cpu_usage1 = cpu_usage1[:dotIndex]
	}
	
	if usage, err := strconv.Atoi(cpu_usage1); err == nil {
		finalUsage := usage / cores
		return fmt.Sprintf("%d%%", finalUsage), nil
	}
	
	return "0%", nil
}

func (s *SystemService) getRAMInfo() (int, int, string, error) {
	// Get used RAM
	usedOutput, err := s.executeCommand("free -m | grep Mem: | awk '{print $3}'")
	if err != nil {
		return 0, 0, "", err
	}
	used, _ := strconv.Atoi(strings.TrimSpace(usedOutput))

	// Get total RAM
	totalOutput, err := s.executeCommand("free -m | grep Mem: | awk '{print $2}'")
	if err != nil {
		return 0, 0, "", err
	}
	total, _ := strconv.Atoi(strings.TrimSpace(totalOutput))

	// Calculate usage percentage
	usage := fmt.Sprintf("%.1f%%", float64(used)/float64(total)*100)

	return used, total, usage, nil
}

func (s *SystemService) writeToFile(filename, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}