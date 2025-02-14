# CDNDe 是一款判断域名是否存在 CDN 的工具

# 项目简介

CDNDe是一款使用Go语言开发的工具，它能批量判断域名是否存在CDN

项目地址：https://github.com/Msup5/CDNDe

# 选项

```
参数:
  -h, --help  查看帮助
  -f          打开文件
  -o          输出文件
  -t          保存类型, 默认ip
  -g          线程, 默认20
```



# 使用说明

简单用法

```
cdnde.exe -f target.txt  批量查看域名解析IP
```

其他用法

```
cdnde.exe -f target.txt -o ./Desktop            保存无CDN域名的IP（默认）
cdnde.exe -f target.txt -o ./Desktop -t domain	保存无CDN的域名
```

# 免责声明

本工具仅用于学习，严禁用于任何非法活动。使用本文所述技术前，请确保已获得目标系统所有者的明确授权。任何滥用信息造成的法律责任及后果均由使用者自行承担，作者不承担任何责任。
