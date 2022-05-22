# 构建

windows: CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o easy_up_cloud.exe
linux: CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o easy_up_cloud_linux
mac: CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o easy_up_cloud_mac


TODO
- 同账户登录后，已登陆账户挤下线
- 增加配置文件
- 订阅资源组，邮件通知

后续大改动：
1. 使用 Gorm 重构数据库部分
2. 精简Controller逻辑，提取公共部分
3. 基础函数库完善
