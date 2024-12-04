#!/bin/bash

# API基础URL
BASE_URL="http://localhost:8080/api/v1"
TOKEN=""
USER_ID=""

# 颜色
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# 打印带颜色的消息
print_success() {
    echo -e "${GREEN}$1${NC}"
}

print_error() {
    echo -e "${RED}$1${NC}"
}

# 测试用户注册
test_register() {
    echo "Testing user registration..."
    response=$(curl -s -w "\n%{http_code}" -X POST "${BASE_URL}/auth/register" \
        -H "Content-Type: application/json" \
        -d '{
            "username": "testuser",
            "email": "test@gmail.com",
            "password": "Test123!@#"
        }')
    
    status_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$status_code" -eq 200 ]; then
        print_success "Registration successful: $body"
        # 提取用户ID
        USER_ID=$(echo "$body" | grep -o '"data":{"user":{"id":[0-9]*' | grep -o '[0-9]*$')
        echo "Extracted user ID: ${USER_ID}"
    else
        print_error "Registration failed: $body"
    fi
}

# 测试用户登录
test_login() {
    echo "Testing user login..."
    response=$(curl -s -w "\n%{http_code}" -X POST "${BASE_URL}/auth/login" \
        -H "Content-Type: application/json" \
        -d '{
            "email": "test@gmail.com",
            "password": "Test123!@#"
        }')
    
    status_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$status_code" -eq 200 ]; then
        # 使用jq提取token
        TOKEN=$(echo "$body" | jq -r '.data.token')
        if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
            TOKEN=$(echo "$body" | grep -o '"data":{"token":"[^"]*"' | sed 's/"data":{"token":"\([^"]*\)"/\1/')
        fi
        print_success "Login successful, token: $TOKEN"
    else
        print_error "Login failed: $body"
    fi
}

# 测试获取用户列表
test_list_users() {
    echo "Testing get users list..."
    response=$(curl -s -w "\n%{http_code}" -X GET "${BASE_URL}/users" \
        -H "Authorization: Bearer $TOKEN")
    
    status_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$status_code" -eq 200 ]; then
        print_success "Get users list successful: $body"
    else
        print_error "Get users list failed: $body"
    fi
}

# ���试获取用户详情
test_get_user() {
    echo "Testing get user details..."
    response=$(curl -s -w "\n%{http_code}" -X GET "${BASE_URL}/users/${USER_ID}" \
        -H "Authorization: Bearer $TOKEN")
    
    status_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$status_code" -eq 200 ]; then
        print_success "Get user details successful: $body"
    else
        print_error "Get user details failed: $body"
    fi
}

# 测试更新用户
test_update_user() {
    echo "Testing update user..."
    response=$(curl -s -w "\n%{http_code}" -X PUT "${BASE_URL}/users/${USER_ID}" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d '{
            "username": "updateduser",
            "email": "updated@gmail.com",
            "password": "Test123!@#"
        }')
    
    status_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$status_code" -eq 200 ]; then
        print_success "Update user successful: $body"
    else
        print_error "Update user failed: $body"
    fi
}

# 测试删除用户
test_delete_user() {
    echo "Testing delete user..."
    response=$(curl -s -w "\n%{http_code}" -X DELETE "${BASE_URL}/users/${USER_ID}" \
        -H "Authorization: Bearer $TOKEN")
    
    status_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$status_code" -eq 204 ]; then
        print_success "Delete user successful"
    else
        print_error "Delete user failed: $body"
    fi
}

# 运行所有测试
echo "Starting API tests..."
test_register
sleep 1
test_login
sleep 1
test_list_users
sleep 1
test_get_user
sleep 1
test_update_user
sleep 1
test_delete_user
echo "API tests completed." 