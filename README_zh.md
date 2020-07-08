# nginxConfigFormatterGo

使用go语言开发的, 优雅的nginx配置文件格式化工具

受 <https://github.com/1connect/nginx-config-formatter.git> 项目激发和鼓励  
之前也有给该项目提交补丁, 限于python水平有限, 于是决定用go重新造一遍轮子.

## 项目目标和特性

- 可预测的格式化结果.
- 所有的注释都独立一行.
- 连续的多个空行合并为一个空行.
- 花括号使用Java的习惯.
- 所有的行使用统一的方式进行缩进, 缩进的空格数由用户指定 (默认 4 个空格).
- 多余的空白字符合并为一个空格, 但是在注释和引号中的空白字符不进行任何处理.

## 编译要求

go 1.14.4+ (or go 1.13.12+)

## 安装

### 1. go get 方式

```shell
go get github.com/rwx------/nginxConfigFormatterGo

# 可能会被安装在如下目录
$HOME/go/bin/nginxConfigFormatterGo
```

### 2. go build 方式

```shell
git clone https://github.com/rwx------/nginxConfigFormatterGo.git
cd nginxConfigFormatterGo
go build
```

### 3. 预编译好的二进制包

你可以在 [发布页面](https://github.com/rwx------/nginxConfigFormatterGo/releases) 获取预编译的二进制包.

## 使用方法

```code
NAME:
   nginxConfigFormatterGo - nginx 格式化工具

USAGE:
   ./nginxConfigFormatterGo [-s 2] [-c utf-8] [-b] [-v] [-t] <filelists>

DESCRIPTION:
   nginx 格式化工具

AUTHOR:
   github.com/rwx------

COMMANDS:
   help, h  显示命令列表或单个命令的帮助

GLOBAL OPTIONS:
   --charset value, -c value  当前支持的字符集: gbk, gb18030, windows-1252, utf-8 (默认: "utf-8")
   --space value, -s value    缩进的空格数 (默认: 4)
   --backup, -b               备份原始的配置文件
   --verbose, -v              冗长模式
   --testing, -t              只进行测试, 不真正执行
   --help, -h                 显示本页的帮助信息
```
