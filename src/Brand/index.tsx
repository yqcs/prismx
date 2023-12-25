import { Divider } from "antd";
import React from "react";
import styles from "./index.less";

const Brand = () => {
  return (
    <div className={styles.brand}>
      <div>
        <div className={styles.content}>
          <span className={styles.title}>
            轻量，跨平台
          </span>
          <div className={styles.desc}>
            支持Windows、Linux、MacOS，甚至可以在Raspberry Pi、安卓手机上构建您的风险检测系统
          </div>
          <div className={styles.homeDiv}>
            <div className={styles.phone}>
              <img src="/static/phone_home.png" alt="phone home" />
            </div>
            <div className={styles.pc}>
              <img src="/static/pc_home.jpg" alt="pc home" />
            </div>
          </div>
        </div>
      </div>
      <div>
        <div className={styles.content2}>
          <span className={styles.title}>CLI / WEB 切换</span>
          <div className={styles.desc}>
            以CLI命令行临时扫描，亦可以WEB服务常驻
          </div>
          <div className={styles.homeDiv}>
            <div className={styles.cli}>
              <img src="/static/cli.png" alt="cli" />
            </div>
            <div className={styles.web}>
              <img src="/static/pc2.png" alt="pc home" />
            </div>
          </div>
        </div>
      </div>
      <Divider style={{ marginTop: "5%" }} />
    </div>
  );
};

export default Brand;
