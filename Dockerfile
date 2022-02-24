FROM alpine
# 维护者信息
MAINTAINER hzx790254812@gmail.com
# 添加二进制文件进工作目录
ADD pushdeer-drone-plugin /bin/
# 运行证书
#RUN apk --update add ca-certificates
# 运行
ENTRYPOINT ["/bin/pushdeer-drone-plugin"]