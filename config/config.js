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
