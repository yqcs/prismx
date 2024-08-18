<h1 align="center">

<a href="https://prismx.io/"><img src="public/static/scan.png" width="200px"></a>

</h1>

<h1 align="center">:: 棱镜 X · 轻量型跨平台单兵渗透系统</h1>

---

<p align="center">
  <a href="https://prismx.io/guide" target="_blank">使用文档</a> ·
  <a href="https://prismx.io/guide">远程管理</a> ·
  <a href="https://prismx.io/guide">风险扫描</a> ·
  <a href="https://prismx.io/guide">邮件测试</a> ·
  <a href="https://prismx.io/guide">一键应急</a>
</p>

## 启动

### · WEB 系统

##### 依赖文件：

- lib.zip： web 版依赖库，CLI 模式无需下载。

存储仓库： https://oss.prismx.io Linux Amd64 运行示例：

```bash
$ wget https://oss.prismx.io/lib.zip
$ wget https://oss.prismx.io/prismx_linux_amd64
$ unzip lib.zip
$ chmod +x prismx_linux_amd64
$ ./prismx_linux_amd64
```

启动后访问`https://yourIP:443`即可进入登录页，使用 -port 参数可指定端口。系统默认账号`prismx/prismx@passw0rd`
，首次使用请修改账户名与密码！

#### 主页：

<img src="public/static/pc_home.jpg" alt="pc_home"/>

#### 数据大屏：

<img src="public/static/view.jpg" alt="pc_home"/>

### · CLI 命令行

命令行模式无需任何依赖文件，只具有基础的扫描模块。执行-h 命令可获取相关帮助。

```bash
$ ./prismx_linux_amd64_cli -h
$ ./prismx_linux_amd64_cli -t 127.0.0.1 -p 1-500,3000-6000
```

<img src="public/static/cli.png" alt="pc_home"/>

---

## QQ 安全研究群：

### [点击加入：528118163](https://jq.qq.com/?_wv=1027&k=azWZhmSy)

## 加群 / 联系（左） | 公众号：遮天实验室（右）

<img src="public/static/wx.jpg" width="200"><img src="public/static/wx_qrcode.jpg" width="200">
