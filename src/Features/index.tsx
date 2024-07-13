import {Col, Row} from "antd";
import React from "react";
import styles from "./index.less";

const Features = () => {
    return (
        <div className={styles.features}>
            <div className={styles.title}>集成一体化系统</div>
            <div className={styles.zdyDivC}>
                <Row>
                    <Col xs={12} lg={6} className={styles.card}>
                        <div className={styles.normal}>
                            <svg
                                width="60"
                                viewBox="0 0 1024 1024"
                                version="1.1"
                                xmlns="http://www.w3.org/2000/svg">
                                <path
                                    d="M964.024495 639.989665A52.987774 52.987774 0 0 0 921.531884 694.001357V895.96925a25.597959 25.597959 0 0 1-25.597958 25.597959h-201.967893A52.987774 52.987774 0 0 0 639.95434 964.05982 51.195917 51.195917 0 0 0 691.150257 1023.959043h230.381627a102.391834 102.391834 0 0 0 102.391835-102.391834v-230.381627a51.195917 51.195917 0 0 0-59.899224-51.195917zM332.778837 921.567209H127.995169a25.597959 25.597959 0 0 1-25.597959-25.597959v-201.967893A52.987774 52.987774 0 0 0 59.904599 639.989665 51.195917 51.195917 0 0 0 0.005376 691.185582v230.381627a102.391834 102.391834 0 0 0 102.391834 102.391834h227.565852A52.987774 52.987774 0 0 0 383.974754 981.466432 51.195917 51.195917 0 0 0 332.778837 921.567209zM921.531884 0.040701h-227.565851A52.987774 52.987774 0 0 0 639.95434 42.533312 51.195917 51.195917 0 0 0 691.150257 102.432535h204.783669a25.597959 25.597959 0 0 1 25.597958 25.597959v201.967893A52.987774 52.987774 0 0 0 964.024495 384.010079 51.195917 51.195917 0 0 0 1023.923719 332.814162V102.432535a102.391834 102.391834 0 0 0-102.391835-102.391834zM59.904599 384.010079A52.987774 52.987774 0 0 0 102.39721 329.998387V128.030494a25.597959 25.597959 0 0 1 25.597959-25.597959h201.967893A52.987774 52.987774 0 0 0 383.974754 59.939924 51.195917 51.195917 0 0 0 332.778837 0.040701H102.39721a102.391834 102.391834 0 0 0-102.391834 102.391834v230.381627a51.195917 51.195917 0 0 0 59.899223 51.195917z"
                                    fill="#13227a"
                                />
                                <path
                                    d="M724.171624 200.728696H299.757471a99.320079 99.320079 0 0 0-99.0641 99.0641v270.826401a99.320079 99.320079 0 0 0 99.0641 99.0641H460.76863V716.783541h-102.391834a51.195917 51.195917 0 0 0 0 102.391834h307.175503a51.195917 51.195917 0 0 0 0-102.391834h-76.793876v-47.100244h135.413201a99.320079 99.320079 0 0 0 99.064099-99.0641V299.792796a99.320079 99.320079 0 0 0-99.064099-99.0641z m-11.519082 358.37142H311.276552v-247.788239h401.37599z"
                                    fill="#13227a"
                                />
                            </svg>
                            <div className={styles.cardTitle}>系统风险检测</div>
                            <div className={styles.summary}>快速扫描安全隐患</div>
                        </div>
                    </Col>
                    <Col xs={12} lg={6} className={styles.card}>
                        <div className={styles.normal}>
                            <svg
                                className="icon"
                                viewBox="0 0 1024 1024"
                                version="1.1"
                                xmlns="http://www.w3.org/2000/svg"
                                width="60"
                                height="60"
                            >
                                <path
                                    d="M157 258L472.61 75.72a83.85 83.85 0 0 1 83.88 0L872.13 258a83.9 83.9 0 0 1 41.94 72.64v364.43a83.87 83.87 0 0 1-41.94 72.64L556.49 950a83.9 83.9 0 0 1-83.88 0L157 767.71a83.87 83.87 0 0 1-42-72.64V330.59A83.9 83.9 0 0 1 157 258z m359.58 622.8l317.59-185.72-2-367.92-319.61-182.25-317.63 185.68 2 367.92 315.63 182.24a4 4 0 0 0 2 0.55z"
                                    fill="#13227a"
                                />
                                <path
                                    d="M494.57 461.35l-125.21-72.29c-26.63-15.38-26.63-53.82 0-69.2l125.21-72.29a39.93 39.93 0 0 1 39.95 0l125.21 72.29c26.64 15.38 26.64 53.82 0 69.2l-125.21 72.29a40 40 0 0 1-39.95 0z"
                                    fill="#13227a"
                                />
                                <path
                                    d="M514.55 380.61l45.3-26.14-45.3-26.16-45.3 26.16 45.3 26.14zM480 555.88v144.57c0 30.76-33.3 50-59.93 34.6l-125.26-72.29a39.94 39.94 0 0 1-20-34.59V483.59c0-30.76 33.3-50 59.93-34.6L460 521.28a39.94 39.94 0 0 1 20 34.6z"
                                    fill="#13227a"
                                />
                                <path
                                    d="M400.04 578.95l-45.29-26.16-0.01 52.31 45.31 26.14-0.01-52.29zM569.12 521.28L694.33 449c26.64-15.38 59.93 3.84 59.93 34.6v144.57a39.92 39.92 0 0 1-20 34.59l-125.2 72.29c-26.64 15.38-59.93-3.84-59.93-34.6V555.88a39.93 39.93 0 0 1 19.99-34.6z"
                                    fill="#13227a"
                                />
                                <path
                                    d="M629.06 578.95l-0.01 52.29 45.3-26.14v-52.31l-45.29 26.16z"
                                    fill="#13227a"
                                />
                            </svg>
                            <div className={styles.cardTitle}>资产应用识别</div>
                            <div className={styles.summary}>
                                一键收集资产暴露面
                            </div>
                        </div>
                    </Col>
                    <Col xs={12} lg={6} className={styles.card}>
                        <div className={styles.normal}>
                            <svg
                                className="icon"
                                viewBox="0 0 1024 1024"
                                version="1.1"
                                xmlns="http://www.w3.org/2000/svg"
                                width="60">
                                <path
                                    d="M928.2 640h-64v-96.2c0-17.7-14.3-32-32-32-1.2 0-2.4 0.1-3.6 0.2H543.8v-63.8h192.3c17.7 0 32-14.3 32-32v-286c0-17.7-14.3-32-32-32H288c-17.7 0-32 14.3-32 32v286c0 17.7 14.3 32 32 32h191.8V512H192c-17.7 0-32 14.3-32 32v96H96.2c-17.7 0-32 14.3-32 32v224.3c0 17.7 14.3 32 32 32h191.1c17.7 0 32-14.3 32-32V672c0-17.7-14.3-32-32-32H224v-64h255.8v64H416c-17.7 0-32 14.3-32 32v224.3c0 17.7 14.3 32 32 32h192c17.7 0 32-14.3 32-32V672c0-17.7-14.3-32-32-32h-64.3v-64H800v64h-64c-17.7 0-32 14.3-32 32v224.3c0 17.7 14.3 32 32 32h192c17.7 0 32-14.3 32-32V672c0.2-17.7-14.1-32-31.8-32zM320 162.2h384.1v222H517c-1.7-0.3-3.4-0.4-5.2-0.4-1.8 0-3.5 0.1-5.2 0.4H320v-222z m-64.7 702.1H128.2V704h127.1v160.3z m320.8 0h-128V704h59.8c1.3 0.2 2.6 0.3 4 0.3s2.7-0.1 4-0.3h60.3v160.3z m320.1 0h-128V704h128v160.3z"
                                    fill="#13227a"
                                />
                            </svg>
                            <div className={styles.cardTitle}>远程协助</div>
                            <div className={styles.summary}>仅需安装Agent即可协助操作</div>
                        </div>
                    </Col>
                    <Col xs={12} lg={6} className={styles.card}>
                        <div className={styles.normal}>
                            <svg
                                className="icon"
                                viewBox="0 0 1024 1024"
                                version="1.1"
                                xmlns="http://www.w3.org/2000/svg"
                                width="60"
                            >
                                <path
                                    d="M518.9632 946.4832c-10.5984 0-20.992-1.7408-30.0544-5.2736-106.7008-41.1648-199.3216-107.3664-275.3536-196.8128-45.824-53.9136-70.6048-121.0368-73.6768-199.5776-1.4336-36.8128-1.0752-73.984-0.768-109.8752 0.1536-14.4384 0.256-28.8768 0.3072-43.3152-0.2048-1.4848-0.3584-3.0208-0.3584-4.608 0-15.5648 0.0512-31.0784 0.1024-46.6432 0.1024-37.4272 0.256-76.1856-0.4096-114.0736s18.3808-62.6688 55.0912-71.7312c48.4864-11.9296 96.9728-24.064 145.408-36.2496 52.4288-13.1584 106.5984-26.7264 160-39.8336 12.9536-3.1744 27.8528-2.7648 44.288 1.28 97.536 23.808 196.5568 48.4352 292.3008 72.2944l2.0992 0.512c48.0256 11.9808 61.9008 29.8496 61.9008 79.6672l0.0512 77.568c0.0512 69.5808 0.1024 141.568-0.1024 212.3776-0.3072 115.6608-46.7456 212.4288-137.8816 287.5904-58.2144 47.9744-127.3856 101.632-212.224 131.584-9.6768 3.4304-20.2752 5.12-30.72 5.12zM200.4992 382.4128c0.256 1.4848 0.3584 3.072 0.3584 4.608 0 16.1792-0.1536 32.3584-0.3072 48.4864-0.3584 35.1744-0.7168 71.5264 0.7168 106.8544 2.56 64.4608 22.4256 119.04 59.136 162.2016 69.2736 81.5104 153.6 141.8752 250.624 179.3024 3.7888 1.4848 12.2368 1.6384 18.176-0.4608 75.4176-26.624 139.4688-76.4416 193.5872-121.088 77.4656-63.8464 115.2512-142.4384 115.5584-240.3328 0.2048-70.7072 0.1536-142.592 0.1024-212.1216l-0.0512-77.6192c0-8.5504-0.3072-13.2096-0.6656-15.7184-2.2528-0.9216-6.5536-2.304-14.6432-4.3008l-2.0992-0.512c-95.6928-23.8592-194.6112-48.4864-291.9936-72.2432-7.5264-1.8432-12.8512-1.8432-15.104-1.28-53.248 13.056-107.3664 26.624-159.6928 39.7312-47.7184 11.9808-97.0752 24.3712-145.7152 36.3008-5.5808 1.3824-7.3728 2.816-7.5264 2.9184 0.0512 0-0.9216 2.1504-0.8192 8.0384 0.6656 38.5536 0.512 77.6192 0.4096 115.3536 0.0512 14.0288 0 27.9552-0.0512 41.8816z"
                                    fill="#13227a"
                                />
                                <path
                                    d="M525.4144 779.5712c-2.7136 0-5.4272-0.3584-8.1408-1.1264a30.75584 30.75584 0 0 1-22.5792-29.6448v-165.2224h-15.0016c-7.3728 0-14.9504 0.0512-22.4768 0.0512-17.0496 0.1024-34.7136 0.1536-52.2752-0.2048-26.1632-0.6144-47.4624-11.7248-58.4192-30.5152-10.5984-18.1248-10.0352-40.9088 1.4848-62.5152 34.5088-64.6144 70.2464-129.8944 104.7552-193.024 13.3632-24.4224 26.7264-48.8448 40.0384-73.2672a30.72 30.72 0 0 1 32.8704-15.4624l5.632 1.1264c14.3872 2.816 24.832 15.4624 24.832 30.1568v163.4304h4.096c25.5488 0 50.3808-0.0512 75.2128 0.0512 40.6528 0.1024 58.88 16.4864 67.0208 30.1568 8.0896 13.6704 13.7216 37.3248-5.4784 72.6528-49.7664 91.5456-95.0784 174.2336-145.152 258.3552a30.96064 30.96064 0 0 1-26.4192 15.0016z m0-257.4336c16.9472 0 30.72 13.7728 30.72 30.72v81.5616c28.6208-50.7392 56.9344-102.5536 86.8352-157.5936 2.6624-4.9152 4.0448-8.4992 4.8128-10.9056-2.56-0.5632-6.5536-1.1264-12.4928-1.1264-24.7808-0.0512-49.5616-0.0512-75.008-0.0512h-34.816c-16.9472 0-30.72-13.7728-30.72-30.72V348.7232c-30.72 56.1152-62.1056 113.664-92.5184 170.5984-0.4096 0.8192-0.768 1.536-1.024 2.0992 1.28 0.256 2.9696 0.512 5.1712 0.5632 16.6912 0.3584 33.9456 0.3072 50.5856 0.2048 7.5776-0.0512 15.104-0.1024 22.784-0.0512h45.6704z"
                                    fill="#13227a"
                                />
                            </svg>
                            <div className={styles.cardTitle}>风险验证</div>
                            <div className={styles.summary}>
                                安全风险一键验证
                            </div>
                        </div>
                    </Col>
                </Row>
            </div>
        </div>
    );
};
export default Features;