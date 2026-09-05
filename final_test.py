#!/usr/bin/env python3
import urllib.request
import json
import ssl

# Test through Cloudflare with proper headers
ctx = ssl.create_default_context()
ctx.check_hostname = False
ctx.verify_mode = ssl.CERT_NONE

headers = {
    'Content-Type': 'application/json',
    'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36'
}

all_passed = False
token = None

# Test 1: Register
print("=== Test 1: Register ===")
url = 'https://ads.alaikis.com/api/v1/auth/register'
data = json.dumps({'email': 'final_test@example.com', 'password': 'password123'}).encode()
req = urllib.request.Request(url, data=data, headers=headers)
try:
    response = urllib.request.urlopen(req, context=ctx)
    result = json.loads(response.read().decode())
    if result.get('code') == 0:
        print("✅ PASSED")
        token = result['data']['access_token']
    else:
        print(f"❌ FAILED: {result}")
except urllib.error.HTTPError as e:
    print(f"❌ FAILED: {e.code}: {e.read().decode()}")

# Test 2: Login
print("\n=== Test 2: Login ===")
url = 'https://ads.alaikis.com/api/v1/auth/login'
data = json.dumps({'email': 'final_test@example.com', 'password': 'password123'}).encode()
req = urllib.request.Request(url, data=data, headers=headers)
try:
    response = urllib.request.urlopen(req, context=ctx)
    result = json.loads(response.read().decode())
    if result.get('code') == 0:
        print("✅ PASSED")
        token = result['data']['access_token']
    else:
        print(f"❌ FAILED: {result}")
except urllib.error.HTTPError as e:
    print(f"❌ FAILED: {e.code}: {e.read().decode()}")

# Add auth header for subsequent requests
auth_headers = {**headers, 'Authorization': f'Bearer {token}'}

# Test 3: Get stores
print("\n=== Test 3: Get Stores ===")
url = 'https://ads.alaikis.com/api/v1/stores'
req = urllib.request.Request(url, headers=auth_headers)
try:
    response = urllib.request.urlopen(req, context=ctx)
    result = json.loads(response.read().decode())
    if result.get('code') == 0:
        print("✅ PASSED")
    else:
        print(f"❌ FAILED: {result}")
except urllib.error.HTTPError as e:
    print(f"❌ FAILED: {e.code}: {e.read().decode()}")

# Test 4: Get campaigns
print("\n=== Test 4: Get Campaigns ===")
url = 'https://ads.alaikis.com/api/v1/advertising/campaigns'
req = urllib.request.Request(url, headers=auth_headers)
try:
    response = urllib.request.urlopen(req, context=ctx)
    result = json.loads(response.read().decode())
    if result.get('code') == 0:
        print("✅ PASSED")
    else:
        print(f"❌ FAILED: {result}")
except urllib.error.HTTPError as e:
    print(f"❌ FAILED: {e.code}: {e.read().decode()}")

# Test 5: Get metrics summary
print("\n=== Test 5: Get Metrics Summary ===")
url = 'https://ads.alaikis.com/api/v1/metrics/summary'
req = urllib.request.Request(url, headers=auth_headers)
try:
    response = urllib.request.urlopen(req, context=ctx)
    result = json.loads(response.read().decode())
    if result.get('code') == 0:
        print("✅ PASSED")
    else:
        print(f"❌ FAILED: {result}")
except urllib.error.HTTPError as e:
    print(f"❌ FAILED: {e.code}: {e.read().decode()}")

# Test 6: Get products
print("\n=== Test 6: Get Products ===")
url = 'https://ads.alaikis.com/api/v1/products'
req = urllib.request.Request(url, headers=auth_headers)
try:
    response = urllib.request.urlopen(req, context=ctx)
    result = json.loads(response.read().decode())
    if result.get('code') == 0:
        print("✅ PASSED")
    else:
        print(f"❌ FAILED: {result}")
except urllib.error.HTTPError as e:
    print(f"❌ FAILED: {e.code}: {e.read().decode()}")

# Test 7: Get feeds
print("\n=== Test 7: Get Feeds ===")
url = 'https://ads.alaikis.com/api/v1/feeds'
req = urllib.request.Request(url, headers=auth_headers)
try:
    response = urllib.request.urlopen(req, context=ctx)
    result = json.loads(response.read().decode())
    if result.get('code') == 0:
        print("✅ PASSED")
    else:
        print(f"❌ FAILED: {result}")
except urllib.error.HTTPError as e:
    print(f"❌ FAILED: {e.code}: {e.read().decode()}")

# Test 8: Get reports
print("\n=== Test 8: Get Reports ===")
url = 'https://ads.alaikis.com/api/v1/reports'
req = urllib.request.Request(url, headers=auth_headers)
try:
    response = urllib.request.urlopen(req, context=ctx)
    result = json.loads(response.read().decode())
    if result.get('code') == 0:
        print("✅ PASSED")
    else:
        print(f"❌ FAILED: {result}")
except urllib.error.HTTPError as e:
    print(f"❌ FAILED: {e.code}: {e.read().decode()}")

# Test 9: Get rules
print("\n=== Test 9: Get Rules ===")
url = 'https://ads.alaikis.com/api/v1/rules'
req = urllib.request.Request(url, headers=auth_headers)
try:
    response = urllib.request.urlopen(req, context=ctx)
    result = json.loads(response.read().decode())
    if result.get('code') == 0:
        print("✅ PASSED")
    else:
        print(f"❌ FAILED: {result}")
except urllib.error.HTTPError as e:
    print(f"❌ FAILED: {e.code}: {e.read().decode()}")

# Test 10: Get workspaces
print("\n=== Test 10: Get Workspaces ===")
url = 'https://ads.alaikis.com/api/v1/workspaces'
req = urllib.request.Request(url, headers=auth_headers)
try:
    response = urllib.request.urlopen(req, context=ctx)
    result = json.loads(response.read().decode())
    if result.get('code') == 0:
        print("✅ PASSED")
    else:
        print(f"❌ FAILED: {result}")
except urllib.error.HTTPError as e:
    print(f"❌ FAILED: {e.code}: {e.read().decode()}")

print("\n" + "=" * 40)
print("🎉 ALL TESTS COMPLETED!")
