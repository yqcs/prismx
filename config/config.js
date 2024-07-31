import { defineConfig } from "dumi";

export default defineConfig({
  sitemap: {
    hostname: "https://prismx.io",
  },
  hash: true,
  title: "Prism X / 棱镜X · 单兵渗透平台",
  mode: "site",
  logo: "/static/scan.png",
  runtimePublicPath: true,
  favicon: "/static/scan.svg",
  metas: [
    {
      name: "keywords",
      content:
        "棱镜X,PrismX,prismx,棱镜,prism,漏洞扫描器,漏洞扫描工具,渗透工具,渗透平台",
    },
    {
      name: "description",
      content: ":: 棱镜 X · 自动化企业网络安全风险检测、漏洞扫描工具。",
    },
  ],
  exportStatic: {},
  extraBabelPlugins: [
    [
      "import",
      {
        libraryName: "antd",
        libraryDirectory: "es",
        style: true,
      },
    ],
  ],
  mfsu: {},
});
