﻿# 自动抓取网络节点

## 声明
```
此程序自己编码较少，应用别人的代码较多
```

## 如何使用
```
git clone https://github.com/sanzhang007/crawl_nodes.git
cd crawl_nodes
go run main.go -test all
```

### 抓取的节点已分享到github

[节点分享链接](https://github.com/sanzhang007/node_free)


### urls 
存放要抓取link 可解析出网页中节点，不过大部分节点来自sub类网站
```yaml
index: []
urls:
  ########################自己付费节点(备用)###############################
  #
  ########################################################################

  ###########################优质节点#####################################
  - link
  ########################################################################
```

## 感谢
### 代码来源
- [clash](https://github.com/Dreamacro/clash)
- [LiteSpeedTest](https://github.com/xxf098/LiteSpeedTest)


