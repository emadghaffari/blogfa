# auth
auth grpc micro service for blogfa application
 
 users:
 - register
 - login with phone-SMS
 - loign with user-pass
 
 admins:
 - search for users
 - verify user
 - send sms to spesific user
 - send notif, sms to spesific group(blue-red-green)
 - get lezzy users server streaming
 - manage users
#### services:
 - service tracer: jaeger
 - config server: etcd
 - cmd: gobra
 - configs: from file(dev), from config server(prod)
 - logger: zap

#### rename config.example.yaml to config.yaml


