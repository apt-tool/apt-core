# PTaaS Core

![](https://img.shields.io/badge/language-golang_v1.20-lightblue)
![GitHub release (with filter)](https://img.shields.io/github/v/release/ptaas-tool/base-api)

This is ```PTaaS``` base api system. In this service we use ```scanner```, ```ftp server```, and ```ml``` components
to perform our penetration testing stages. In ```pkg/models``` directory we defined our
base database modules and system modules to be used in all other
system components.

## Image

Core api docker image address:

```shell
docker pull amirhossein21/ptaas-tool:base-v0.2.4
```

### config

Core system config file (config.yaml) template is something like this:

```yaml
core: # core api
  port: 9090
  enable: true
  workers: 1
  secret: "secret"
mysql: # database
  host: 'localhost'
  port: 3306
  user: root
  pass: ''
  database: 'apt'
  migrate: false
migrate: # migration commands
  root: 'admin'
  pass: '12345'
  enable: false
ftp: # apt instructions service
  host: 'http://localhost:9091'
  secret: 'secret'
  access: 'access'
```

## Setup

Setup core api in docker container with following command:

```shell
docker run -d \
  -v type=bind,source=$(pwd)/config.yml,dest=/app/config.yml
  -p 80:80 \
  amirhossein21/ptaas-tool:base-v0.2.4
```
