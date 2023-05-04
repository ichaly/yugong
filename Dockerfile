# 表示依赖 alpine 最新版
FROM alpine:latest
MAINTAINER Shi Chaly<i@eunut.com>
ENV VERSION 1.0

# 在容器根目录 创建一个 apps 目录
WORKDIR /apps

# 挂载容器目录
VOLUME ["/apps/conf"]

# 拷贝app可以执行文件
COPY main/dist/app /apps/app

# 拷贝配置文件到容器中
COPY conf/prod.yml /apps/conf/prod.yml

# 设置时区为上海
RUN apk --no-cache add tzdata
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone

# 设置编码
ENV LANG C.UTF-8

# 暴露端口
EXPOSE 3000

# 运行golang程序的命令
ENTRYPOINT ["/apps/app" ,"run","-c","/apps/conf/prod.yml"]