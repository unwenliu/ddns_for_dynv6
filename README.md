# ddns_for_dynv6
自动更新dynv6的DNS记录指向本地ip
## 使用说明
```
ddnsfordynv6 [-i 网卡名] [-hostname 域名] [-token token] [-4] [-6]
选项：
    -i 网卡名                  ip所绑定的网卡
    -show_ipv4                 显示指定网卡的ipv4地址
    -show_ipv6                 显示指定网卡的ipv6地址
    -hostname 域名             你的域名
    -token token               你的token
    -4                         更新ipv4地址
    -6                         更新ipv6地址
```
