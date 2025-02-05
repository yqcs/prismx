<h1 align="left">棱镜 X · Open Source</h1> 

<a href="README.md">`English`</a> • <a href="README_CN.md">`中文`</a>

---
**棱镜X 集资产发现、指纹识别、弱密码检测、漏洞验证于一体，采用模块化YAML插件策略配置，实现与真实攻击链高度相仿的PoC验证机制。**

- 跨平台、轻量型设计，支持多种操作系统，便于部署和使用。
- 提供主机存活扫描和资产指纹识别功能，全面掌握网络资产状况。
- 具备弱口令识别和漏洞扫描能力，及时发现安全隐患，保障系统安全。
- 内置JNDI外链服务，支持扫描 jndi、rmi 等需外联漏洞。
- 端口指纹识别框架：[**`yqcs/fingerscan`**](https://github.com/yqcs/fingerscan) 

 <h1 align="center">
    <img src="images/scan.png" width="95%" height="350">
</h1>

### 运行命令

```
Usage of prismx_cli.exe:

  -t  string
        扫描主机，格式支持192.168.1.1/24、16、8，192.168.3.1-80，prismx.io，使用英文逗号分割
  -p  string
        扫描端口，支持格式: 80,22,8000-8080
  -bip  string
        过滤主机，支持过滤网段
  -bp  string
        过滤端口，支持端口范围
  -m  string
        扫描速度，参数: s（慢速）, d（中速）, f（快速） 默认 "d"
  -ping  boolean
        低权限下可能无法发送icmp包，默认不启用 -ping=false
  -pn  boolean
        不进行主机存活检测，默认不启用 -pn=false
  -s  boolean   
        联网扫描子域名，默认不启用 -s=false
  -vul  boolean
        检测漏洞，默认启用 -vul=true
  -weak  boolean
        扫描弱口令，默认启用 -weak=true
```

### 源码结构
<Tree>
  <ul>
    <li>
      core: 系统核心
      <ul>
        <li>
          aliveCheck: 主机、端口存活检测
        </li>
        <li>
          hydra: 弱口令检测相关
        </li>
        <li>
          jsFind: 检测js文件是否存在敏感内容
        </li>
        <li>
          owaspTop10: xss、sql注入等检测工具（暂未完成，需后续优化）
        </li>
        <li>
          plugins: 插件注册中心以及插件文件
        </li>
        <li>
          subdomain: 子域名扫描
        </li>
        <li>
          vulnerability: 漏洞检测模块
        </li>
        <li>
          models: 公共模块相关依赖
        </li>
      </ul>
    </li>
    <li>
      scan: 任务调度中心
    </li>
    <li>
      utils: 工具包
          <ul>
            <li>任务列表</li>
            <li>新建任务</li>
          </ul>
    </li>
    <li>
      main.go: 程序入口
    </li>
  </ul>
</Tree>

 <h1 align="center">
    <img src="images/img.png" width="95%" height="350">
</h1>

### 编译

Tips: 推荐使用golang1.20版本进行编译（新版go不再支持windows 7及以下版本） 

```bash
  go build -ldflags "-s -w   -buildid=" -buildmode="pie"  -trimpath  
```
---

## [**`深度定制: Prismx.io`**](https://prismx.io/)

 <h1 align="center"> 
<a href="https://prismx.io/"><img src="https://prismx.io/static/pc_home.jpg"  width="90%" height="350"></a>

<a href="https://prismx.io/"><img src="https://prismx.io/static/view.jpg"  width="90%" height="350"></a>
</h1>

#### 联系 / 定制（左） | 公众号：遮天实验室（右）
<img src="images/wx.jpg" width="150"> <img src="images/wx_qrcode.jpg" width="150">