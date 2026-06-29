 # My-Online-Store


## Deploy / 在线部署

### 方式：Render

[Render](https://render.com) 提供免费 Web 服务托管，适合本应用。

部署步骤：
1. 在 [render.com](https://render.com) 注册账号（推荐用 GitHub 账号登录）
2. 点击 **New +** -> **Web Service**
3. 连接 GitHub，选择仓库 **lilelin111/My-Online-Store**
4. 填写以下配置：
   - **Name**: my-online-store
   - **Runtime**: Go
   - **Build Command**: go build -o server ./web/
   - **Start Command**: ./server
5. 选择 **Free** 套餐
6. 点击 **Create Web Service**

部署完成后会生成 *.onrender.com 域名，可直接访问。

 
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
