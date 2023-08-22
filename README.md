主要后端接口已实现，通过apiPost测试。暂无前端(学习中ing)


所用技术

go 1.19

gorm

gin

grpc

jwt

nacos

consul

jaeger

redis 5.0

mysql 8.0

elasticsearch

阿里云oss


项目整体框架: 图片优化中


快速启动
先配置nacos配置文件。nacos命令空间为user、goods、inventory、orders


docker安装配置命令 地址端口号供参考

安装mysql

docker run --name camps_mysql -e MYSQL_ROOT_PASSWORD=123456 -d -e MYSQL_DATABASE=camps_user -p 8086:3306 mysql:8.0

安装redis

docker run --name camps_redis -d -p 8089:6379 redis:6.2-rc2

安装consul

docker run -d -p 8500:8500 -p 8300:8300 -p 8301:8301 -p 8302:8302 -p 8600:8600/udp consul consul agent -dev -client=0.0.0.0

安装nacos

docker run --name nacos-standalone -e MODE=standalone -e JVM_XMS=512m -e JVM_XMS=512m -e JVM_XMN=256m -p 8848:8848 -d nacos/nacos-server:latest

在web和srv各目录下的main文件即为各个微服务的启动入口

测试
下载附录json文件后，导入apiPost即可
