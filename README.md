 # My-Online-Store


## Deploy / 在线部署

点击下方按钮一键部署到 Railway：

[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/template/My-Online-Store)

或者手动部署：
1. 在 [railway.app](https://railway.app) 注册账号
2. 点击 **New Project** -> **Deploy from GitHub repo**
3. 选择本仓库：**lilelin111/My-Online-Store**
4. Railway 会自动识别 Go 项目并部署
5. 部署完成后会生成公开的 .railway.app 域名

 
 一个在线记账应用，支持用户注册、登录、收入/支出记录管理（增删查）。
 
 ## Screenshots / 界面截图
 
 ### 用户注册
 ![注册页面](screenshot-register.png)
 
 ### 记账记录管理
 ![记账记录页面](screenshot-records.png)
 
 ## How to run / 如何运行
 
 ```bash
 go run web/main.go
 ```
 
 服务启动后访问 **http://localhost:8080** 即可使用。
