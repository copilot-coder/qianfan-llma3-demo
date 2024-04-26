# 概要

- 本项目用于演示如何通过Golang接入百度千帆大模型平台的Meta Llama 3大语言模型。支持流式响应和非流式响应
- 百度云Meta Llama 3接口文档：https://cloud.baidu.com/doc/WENXINWORKSHOP/s/ilv62om62
- 鉴权方式采用基于安全认证Access Key/Secret Key进行签名计算鉴权
  - 获取Access Key和Secret Key的步骤可参考：https://mp.weixin.qq.com/s/_HLe_OK7NuFSqwK_9LwF4Q
  - 签名鉴权算法：https://cloud.baidu.com/doc/Reference/s/Njwvz1wot


# 如何运行
- 获取Access Key和Secret Key，并设置环境变量QIANFAN_ACCESS_KEY和QIANFAN_SECRET_KEY
- 执行go run main.go
- 输出结果示例：

```
北京有很多好玩的景点，以下是一些：

1. ** Forbidden City**（故宫）：中国古代皇宫，拥有丰富的历史文化价值。
2. **Great Wall of China**（长城）：中国最著名的古迹，长城的北京段是最受欢迎的旅游景点。
3. **Tiananmen Square**（天安门广场）：中国最大的城市广场，位于国家人大、国务院和中央军委三家机构的所在地。
4. **Temple of Heaven**（天坛）：中国古代皇帝举行祈年和祈雨仪式的场所，具有丰富的历史文化价值。
5. **Summer Palace**（颐和园）：中国古代皇宫之一，拥有美丽的湖泊、桥梁和亭台。
6. **Hutongs**（胡同）：北京的传统街区，拥有浓郁的中国古代文化气息。
7. **Lama Temple**（喇嘛寺）：中国佛教最重要的寺庙之一，拥有丰富的艺术和文化价值。
8. **Beijing Zoo**（北京动物园）：中国最古老的动物园，拥有各种珍稀动物。
9. **Olympic Park**（奥林匹克公园）：2008年北京奥运会的主场所，拥有多个体育设施和景点。
10. **Yanqing District**（延庆区）：北京市郊外的旅游区，拥有美丽的山水、寺庙和古迹。
11. **Ming Tombs**（明陵）：中国明朝皇帝的陵墓，位于北京市郊外。
12. **Badaling Ski Resort**（八达岭滑雪场）：北京市郊外的滑雪场，拥有多个滑雪道和其他娱乐项目。

这些景点只是北京旅游的tip，北京还有很多其他的好玩景点和活动，欢迎您来探索！
```