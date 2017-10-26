# Orchestration service
Angela is a orchestration service. 
With http restful api, it runs command with ssh protocol. 
It just executes command where you want, without records execution result, but it will callback a http url to report result 

## Config

```
[default]
http_addr: 0.0.0.0:10001

log: ../log/angela.log
level: debug

[angela]
ssh_key: /path/id_rsa
```

