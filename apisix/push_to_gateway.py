import os
import sys
import json
import urllib.request
from io import BytesIO

domain=input("apisix_address:")
proto_id=input("proto_id:")
proto_package=input("proto_package:")
proto_dir=input("proto_dir:")
api_key='edd1c9f034335f136f87ad84b625c8f1'

host_proto= 'http://'+domain+'/apisix/admin/proto/'+proto_id
host_router= 'http://'+domain+'/apisix/admin/routes/'+proto_id
header= {
        'Content-Type':'application/json;charset=UTF-8'
        }


"""
Step 1: 合成单文件proto
"""

buffer = BytesIO()
buffer.write(b'syntax = "proto3";')
buffer.write('package {};\r\n'.format(proto_package).encode(encoding='utf8'))
try:
    for root, dirs, files in os.walk(proto_dir):
        for file in files:
            if not file.endswith('.proto'):
                continue
            file_path = os.path.join(root, file)
            print(file_path)
            with open(file_path,'rb') as f:
                for line in f:
                    if line.startswith(b'syntax'):
                        continue
                    if line.startswith(b'option'):
                        continue
                    if line.startswith(b'//'):
                        continue
                    if line.startswith(b'import'):
                        continue
                    buffer.write(line)
    print('build proto finish')
except Exception as e:
    print(e)
                
buffer.seek(0)
data = {
        'content':buffer.read().decode(encoding='utf8')
        }

"""
Step 2: 注册proto
"""
request = urllib.request.Request(host_proto, json.dumps(data, ensure_ascii=False).encode(encoding='utf8'))
request.add_header('Content-Type', 'application/json;charset=UTF-8')
request.add_header('X-API-KEY', api_key)
request.get_method = lambda:'PUT'

try:
    req = urllib.request.urlopen(request)
    reply = req.read().decode('utf-8')
    req.close()
    del req
    print(reply)
except Exception as e:
    print(e)

"""
Step 3: 注册router
"""
"""
data = {
        "methods": ["GET"],
        "uri": "/"+proto_package+"/Dummy",
        "plugins": {
                "grpc-transcode": {
                    "proto_id": proto_id,
                    "service": proto_package+"." + "Dummy",
                    "method":"Dummy"
                    }
            },
        "upstream": {
            "scheme": "grpc",
            "type": "roundrobin",
            "nodes":{
                "127.0.0.1:50051": 1
                }
            }
        }
request = urllib.request.Request(host_router, json.dumps(data, ensure_ascii=False).encode(encoding='utf8'))
request.add_header('Content-Type', 'application/json;charset=UTF-8')
request.add_header('X-API-KEY', api_key)
request.get_method = lambda:'PUT'

try:
    req = urllib.request.urlopen(request)
    reply = req.read().decode('utf-8')
    req.close()
    del req
    print(reply)
except Exception as e:
    print(e)
"""
