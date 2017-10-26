# Orchestration service
Angela is a orchestration service, with http restful api, it run command with ssh.

## Config

```
[default]
http_addr: 0.0.0.0:10001

log: ../log/angela.log
level: debug

[angela]
ssh_key: /path/id_rsa
```

