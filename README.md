# GoDigDomain
用Go语言编写的域名爆破工具

## Usage
GoDigDomain: A Domain Name Burst Tool

```shell script
root@Ubuntu# ./gdd
Usage:
    gdd -dn Domain
    [-ds DNSServer]

Options:
    -dn 域名
    -ds DNS服务器 (default "114.114.114.114")
    -dt 域名字典 (default "./dict.txt")
```

## Example

```shell script
root@Ubuntu# ./gdd -d example.com
root@Ubuntu# ./gdd -d example.com -ds 1.1.1.1
root@Ubuntu# ./gdd -d example.com -ds 8.8.8.8,1.1.1.1,114.114.114
root@Ubuntu# ./gdd -d example.com -dt ~/myDict.txt
root@Ubuntu# ./gdd -d example.com -ds 8.8.8.8,1.1.1.1 -dt ~/myDict.txt
```

## Todo
- [x] 支持指定字典文件
- [x] 支持多DNS服务器
- [x] 支持多线程
- [ ] 支持获取标题
