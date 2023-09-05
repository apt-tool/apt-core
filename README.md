# APT Core

![](https://img.shields.io/badge/Language-Golang-blue)
![](https://img.shields.io/badge/App-Core-blue)
![GitHub release (with filter)](https://img.shields.io/github/v/release/apt-tool/apt-core)

This is ```apt``` base api system. In this service we use ```apt scanner```, ```apt instructions```, and ```apt AI``` components
to perform our penetration testing stages. In ```pkg/models``` directory we defined our
base database modules and system modules to be used in all other
system components.

## config

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
