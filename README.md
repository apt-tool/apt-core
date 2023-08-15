# APT Core

![](https://img.shields.io/badge/Language-Golang-blue)
![](https://img.shields.io/badge/App-Core-blue)
![GitHub release (with filter)](https://img.shields.io/github/v/release/apt-tool/apt-core)

This is ```apt``` base system. In this service we use scanner and AI component
in order to perform our penetration testing. In ```pkg``` directory we define our
base database modules and system modules. Other components import this directory
in order to use system modules.

## config

Base core config file template is like this:

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
ftp:
  host: 'http://localhost:9091'
  secret: 'secret'
  access: 'access'
```
