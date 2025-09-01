#!/usr/bin/env python3
"""
Comprehensive API Endpoint Testing Script
Tests all VPN API endpoints and reports any errors
"""

import requests
import json
import sys
import time
from datetime import datetime

# Configuration
VPS_IP = input("Enter your VPS IP address: ").strip()
API_BASE_URL = f"http://{VPS_IP}:37849/api/v1"
DEFAULT_USERNAME = "admin"
DEFAULT_PASSWORD = "db4bb47cd788"

# Global variables
auth_token = None
test_results = []

class Colors:
    GREEN = '\033[92m'
    RED = '\033[91m'
    YELLOW = '\033[93m'
    BLUE = '\033[94m'
    ENDC = '\033[0m'
    BOLD = '\033[1m'

def log_result(endpoint, method, status, message, details=None):
    """Log test result"""
    result = {
        'endpoint': endpoint,
        'method': method,
        'status': status,
        'message': message,
        'details': details,
        'timestamp': datetime.now().isoformat()
    }
    test_results.append(result)
    
    color = Colors.GREEN if status == 'PASS' else Colors.RED if status == 'FAIL' else Colors.YELLOW
    print(f"{color}[{status}]{Colors.ENDC} {method} {endpoint} - {message}")
    if details:
        print(f"    Details: {details}")

def make_request(method, endpoint, data=None, headers=None, auth_required=True):
    """Make HTTP request with proper error handling"""
    url = f"{API_BASE_URL}{endpoint}"
    
    # Add auth header if required
    if auth_required and auth_token:
        if headers is None:
            headers = {}
        headers['Authorization'] = f"Bearer {auth_token}"
    
    try:
        if method == 'GET':
            response = requests.get(url, headers=headers, timeout=10)
        elif method == 'POST':
            response = requests.post(url, json=data, headers=headers, timeout=10)
        elif method == 'PUT':
            response = requests.put(url, json=data, headers=headers, timeout=10)
        elif method == 'DELETE':
            response = requests.delete(url, headers=headers, timeout=10)
        else:
            raise ValueError(f"Unsupported method: {method}")
        
        return response
    except requests.exceptions.RequestException as e:
        return None, str(e)

def test_health_check():
    """Test health check endpoint"""
    print(f"\n{Colors.BLUE}=== Testing Health Check ==={Colors.ENDC}")
    
    # Health endpoint is at root level, not under /api/v1
    url = f"http://{VPS_IP}:37849/health"
    try:
        response = requests.get(url, timeout=10)
    except requests.exceptions.RequestException as e:
        log_result('/health', 'GET', 'FAIL', 'Connection error', str(e))
        return False
    if isinstance(response, tuple):
        log_result('/health', 'GET', 'FAIL', 'Connection error', response[1])
        return False
    
    if response.status_code == 200:
        data = response.json()
        log_result('/health', 'GET', 'PASS', 'Health check successful', f"Status: {data.get('status')}")
        return True
    else:
        log_result('/health', 'GET', 'FAIL', f'HTTP {response.status_code}', response.text)
        return False

def test_authentication():
    """Test authentication endpoints"""
    global auth_token
    print(f"\n{Colors.BLUE}=== Testing Authentication ==={Colors.ENDC}")
    
    # Test login
    login_data = {
        "username": DEFAULT_USERNAME,
        "password": DEFAULT_PASSWORD
    }
    
    response = make_request('POST', '/auth/login', login_data, auth_required=False)
    if isinstance(response, tuple):
        log_result('/auth/login', 'POST', 'FAIL', 'Connection error', response[1])
        return False
    
    if response.status_code == 200:
        data = response.json()
        if data.get('success'):
            auth_token = data['data']['token']
            log_result('/auth/login', 'POST', 'PASS', 'Login successful')
            return True
        else:
            log_result('/auth/login', 'POST', 'FAIL', 'Login failed', data.get('error'))
            return False
    else:
        log_result('/auth/login', 'POST', 'FAIL', f'HTTP {response.status_code}', response.text)
        return False

def test_user_endpoints():
    """Test user management endpoints"""
    print(f"\n{Colors.BLUE}=== Testing User Management ==={Colors.ENDC}")
    
    endpoints = [
        ('GET', '/user/profile'),
        ('GET', '/user/list'),
    ]
    
    for method, endpoint in endpoints:
        response = make_request(method, endpoint)
        if isinstance(response, tuple):
            log_result(endpoint, method, 'FAIL', 'Connection error', response[1])
            continue
        
        if response.status_code == 200:
            data = response.json()
            if data.get('success'):
                log_result(endpoint, method, 'PASS', 'Request successful')
            else:
                log_result(endpoint, method, 'FAIL', 'Request failed', data.get('error'))
        else:
            log_result(endpoint, method, 'FAIL', f'HTTP {response.status_code}', response.text)

def test_system_endpoints():
    """Test system monitoring endpoints"""
    print(f"\n{Colors.BLUE}=== Testing System Monitoring ==={Colors.ENDC}")
    
    endpoints = [
        ('GET', '/system/info'),
        ('GET', '/system/status'),
        ('GET', '/system/bandwidth'),
        ('GET', '/domain/current'),
    ]
    
    for method, endpoint in endpoints:
        response = make_request(method, endpoint)
        if isinstance(response, tuple):
            log_result(endpoint, method, 'FAIL', 'Connection error', response[1])
            continue
        
        if response.status_code == 200:
            data = response.json()
            if data.get('success'):
                log_result(endpoint, method, 'PASS', 'Request successful')
            else:
                log_result(endpoint, method, 'FAIL', 'Request failed', data.get('error'))
        else:
            log_result(endpoint, method, 'FAIL', f'HTTP {response.status_code}', response.text)

def test_vpn_endpoints():
    """Test VPN management endpoints"""
    print(f"\n{Colors.BLUE}=== Testing VPN Management ==={Colors.ENDC}")
    
    # Test getting existing users for each protocol
    protocols = ['ssh', 'vmess', 'vless', 'trojan', 'shadowsocks']
    
    for protocol in protocols:
        endpoint = f'/vpn/{protocol}/users'
        response = make_request('GET', endpoint)
        if isinstance(response, tuple):
            log_result(endpoint, 'GET', 'FAIL', 'Connection error', response[1])
            continue
        
        if response.status_code == 200:
            data = response.json()
            if data.get('success'):
                users = data.get('data', [])
                log_result(endpoint, 'GET', 'PASS', f'Found {len(users)} {protocol} users')
            else:
                log_result(endpoint, 'GET', 'FAIL', 'Request failed', data.get('error'))
        else:
            log_result(endpoint, 'GET', 'FAIL', f'HTTP {response.status_code}', response.text)
    
    # Test get all users
    response = make_request('GET', '/vpn/users/all')
    if isinstance(response, tuple):
        log_result('/vpn/users/all', 'GET', 'FAIL', 'Connection error', response[1])
    else:
        if response.status_code == 200:
            data = response.json()
            if data.get('success'):
                users = data.get('data', [])
                log_result('/vpn/users/all', 'GET', 'PASS', f'Found {len(users)} total users')
            else:
                log_result('/vpn/users/all', 'GET', 'FAIL', 'Request failed', data.get('error'))
        else:
            log_result('/vpn/users/all', 'GET', 'FAIL', f'HTTP {response.status_code}', response.text)

def test_create_user_endpoints():
    """Test creating users for each protocol"""
    print(f"\n{Colors.BLUE}=== Testing User Creation ==={Colors.ENDC}")
    
    protocols = ['ssh', 'vmess', 'vless', 'trojan', 'shadowsocks']
    test_username = f"test_user_{int(time.time())}"
    
    for protocol in protocols:
        endpoint = f'/vpn/{protocol}/create'
        user_data = {
            "username": f"{test_username}_{protocol}",
            "password": "testpass123",
            "days": 1
        }
        # Note: protocol field is set automatically by the handler, don't include it
        
        response = make_request('POST', endpoint, user_data)
        if isinstance(response, tuple):
            log_result(endpoint, 'POST', 'FAIL', 'Connection error', response[1])
            continue
        
        if response.status_code in [200, 201]:
            data = response.json()
            if data.get('success'):
                log_result(endpoint, 'POST', 'PASS', f'Created {protocol} user successfully')
            else:
                log_result(endpoint, 'POST', 'FAIL', 'User creation failed', data.get('error'))
        else:
            log_result(endpoint, 'POST', 'FAIL', f'HTTP {response.status_code}', response.text)

def print_summary():
    """Print test summary"""
    print(f"\n{Colors.BOLD}=== TEST SUMMARY ==={Colors.ENDC}")
    
    total_tests = len(test_results)
    passed_tests = len([r for r in test_results if r['status'] == 'PASS'])
    failed_tests = len([r for r in test_results if r['status'] == 'FAIL'])
    
    print(f"Total Tests: {total_tests}")
    print(f"{Colors.GREEN}Passed: {passed_tests}{Colors.ENDC}")
    print(f"{Colors.RED}Failed: {failed_tests}{Colors.ENDC}")
    
    if failed_tests > 0:
        print(f"\n{Colors.RED}FAILED TESTS:{Colors.ENDC}")
        for result in test_results:
            if result['status'] == 'FAIL':
                print(f"  - {result['method']} {result['endpoint']}: {result['message']}")
                if result['details']:
                    print(f"    {result['details']}")
    
    # Save detailed results to file
    with open('tmp_rovodev_test_results.json', 'w') as f:
        json.dump(test_results, f, indent=2)
    print(f"\nDetailed results saved to: tmp_rovodev_test_results.json")

def main():
    """Main test function"""
    print(f"{Colors.BOLD}VPN API Endpoint Testing{Colors.ENDC}")
    print(f"Testing API at: {API_BASE_URL}")
    print("=" * 50)
    
    # Test health check first
    if not test_health_check():
        print(f"{Colors.RED}Health check failed. Cannot continue testing.{Colors.ENDC}")
        return
    
    # Test authentication
    if not test_authentication():
        print(f"{Colors.RED}Authentication failed. Cannot test protected endpoints.{Colors.ENDC}")
        return
    
    # Test all endpoints
    test_user_endpoints()
    test_system_endpoints()
    test_vpn_endpoints()
    test_create_user_endpoints()
    
    # Print summary
    print_summary()

if __name__ == "__main__":
    main()