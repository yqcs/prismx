# Prism X · 单兵渗透系统 / Prism X · Single Soldier Combat Platform

## 特性

- 集渗透前置、后置于一体的轻量型跨平台系统

### 风险扫描
![pc_view](/public/static/view.jpg)

### 主机管理
![pc_home](/public/static/pc_home.jpg)



## 启动

<Alert type="warning">
本工具仅面向合法授权的企业资产风险检测，请严格遵守法律规定，不得危害国家安全、公共利益，不得损害个人、组织的合法权益，否则应自行承担所引起的一切法律责任。
</Alert>

下载对应 OS ARCH 的软件包 [Prism X releases](https://github.com/yqcs/heartsk_community/releases/)
，解压之后赋予可执行权限之后直接运行即可。

Linux amd64 运行示例：

```bash
$ wget https://github.com/yqcs/prismx/releases/download/1.0.10/prismx_linux_amd64.zip
$ unzip prismx_linux_amd64.zip
$ cd prismx_linux_amd64
$ chmod +x prismx
$ ./prismx
```

### WEB 模式

为了方便使用，系统提供了 CLI 命令行以及更具交互性的 WEB 模式两种运行方式。WEB 模式需提供 License 文件，运行`./prismx`
命令即可启动。系统已经签发了 WEB 模式需要的 License 及公钥文件。

运行之后访问`http://yourIP:80`即可进入登录页，使用-port 参数可指定端口。 系统默认账号`prismx/prismx@passw0rd`
，首次使用请修改账户名与密码！

<img src="/public/static/guide/login.png" alt="login Page"/>

### CLI 命令行

命令行模式无需授权及公钥文件，但是只具有基础的扫描模块，无法使用 WEB 模式的扫描配置以及信息收集等高级功能。执行-h
命令可获取相关帮助。

```bash
$ ./prismx -h
$ ./prismx -t 127.0.0.1 -p 1-500,3000-6000
```

<img src="/public/static/cli.png" alt="cli Page"  width="70%"/>

### Linux For ARM（Android）

#### 具有 Root 权限可以避免百分之九十的问题！

安卓设备为例，直接使用 adb push 推送到 `/data/local/tmp/`目录，然后使用`chmod +x `赋予可执行权限即可直接运行。该方案不便随时运行，可使用终端软件
Termux 支撑。

下载终端工具[Termux](https://termux.com/) ，打开软件之后更新软件包然后安装 wget，再下载二进制程序。

```bash
$ pkg update
$ pkg upgrade
$ pkg install wget
$ wget https://github.com/yqcs/prismx/releases/download/1.0.10/prismx_linux_arm64.zip
$ unzip prismx_linux_arm64.zip
$ cd prismx_linux_arm64
$ chmod +x prismx
$ ./prismx
```

未授予 Root 权限会出现错误：` listen tcp 0.0.0.0:80: bind: permission denied`，使用-port 参数切换绑定端口即可。

执行扫描任务时出现错误：`xx on [::1]:53: read udp [::1]:37606->[::1]:53: read: connection refused`

> 有 ROOT 权限：在手机根目录的 /etc/ 文件夹下新建一个名为 resolv.conf 的文件，内容为`nameserver 8.8.8.8`（DNS 服务器），然后重启
> Termux 之后再次运行即可。
>
> 无 ROOT
> 权限：执行`pkg install proot resolv-conf && proot -b $PREFIX/etc/resolv.conf:/etc/resolv.conf ./prismx -port 8000`
> （运行参数）
> 至此，便可成功启动，在手机浏览器访问首页：http://127.0.0.1:8000 但是并不代表可以完整使用了，以非 ROOT 权限执行任务时切记将存活检测切换为
> Ping 模式！！

<img src="/public/static/guide/phone.jpg" alt="phone Page" width="30%"/>
