## Drone pushDeer plugin

The result of the plug-in is a message that the bot will send you.:
```
自动化部署成功

🏠仓库：api

⭕版本：1

🎅提交者：hongzx

🕙耗时：0分3秒

📖提交分支：master

📃提交信息：通知
```

Variables
  - *url* - You sender address Example format : *https://api2.pushdeer.com/message/push* - Required
  - *pushkey* - Your equipment token - Required
  - *type* - Support send type : text、image、markdown - Required
  - *content* - Send Content

Example pipeline
```yml
kind: pipeline
name: project-go-api

steps:
  - name: build
    image: golang:latest
    pull: if-not-exists
    environment:
      GOPROXY: "https://goproxy.cn,direct" 
    volumes:
      - name: pkgdeps
        path: /go/pkg
    commands:
      - CGO_ENABLED=0 go build -o project-go-api
      
  - name: pushdeer
    image: hongzhuangxian/pushdeer-drone-plugin
    settings:
      pushkey: your pushkey 
      url: https://api2.pushdeer.com/message/push
      type: text
```
Build packed:

    set GOOS=linux
    set GOARCH=amd64
    set CGO_ENABLED=0
    go build -o pushdeer-drone-plugin

Build image:

    docker build -t hongzhuangxian/pushdeer-drone-plugin .

Push image:

    docker push hongzhuangxian/pushdeer-drone-plugin
