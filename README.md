# PTaaS base API

![](https://img.shields.io/badge/language-golang_v1.20-lightblue)
![GitHub release (with filter)](https://img.shields.io/github/v/release/ptaas-tool/base-api)

This is ```PTaaS``` base api system. In this service we use ```scanner```, ```ftp server```, and ```ml``` components
to perform our penetration testing stages. In ```pkg/models``` directory we defined our
base database modules and system modules to be used in all other system components.

## Image

Base API docker image address:

```shell
docker pull amirhossein21/ptaas-tool:base-v0.X.X
```

### config

Base API system config file (```config.yaml```) template is something like this:

```yaml
core:
  port: 9090
  enable: true
  workers: 1
  secret: "secret"
mysql:
  host: 'localhost'
  port: 3306
  user: root
  pass: ''
  database: 'apt'
  migrate: false
migrate:
  root: 'admin'
  pass: '12345'
  enable: false
ai:
  enable: true
  method: "svm"
  "limit": 10
  "factor": 7
scanner:
  enable: true
  defaults:
    - "2fa"
  command: "python3 scanner.py"
  flags:
    - "host"
    - "endpoints"
    - "type"
    - "protocol"
ftp:
  host: 'http://localhost:9091'
  secret: 'secret'
  access: 'access'
```

## Setup

Setup base API in docker container with following command:

```shell
docker run -d \
  -v type=bind,source=$(pwd)/config.yml,dest=/app/config.yml
  -p 80:80 \
  amirhossein21/ptaas-tool:base-v0.X.X
```
