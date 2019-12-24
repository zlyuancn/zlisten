/*
-------------------------------------------------
   Author :       zlyuan
   date：         2019/9/23
   Description :
-------------------------------------------------
*/

package zlisten

import (
    "fmt"
    "net"
)

type Config struct {
    BindIP   string // bind的ip
    BindPort int    // bind的端口

    AdvertiseIP       string // 公告ip地址(优先于AdvertiseIPPrefix
    AdvertiseIPPrefix string // 公告ip地址前缀
    AdvertisePort     int    // 公告端口
}

type Listen struct {
    Listener      net.Listener
    AdvertiseIP   string
    AdvertisePort int
}

// 获取所有能检测到的本地ip, 127.*不会返回
func GetLocalIPs() []string {
    var ips []string
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return ips
    }

    for _, address := range addrs {
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                ips = append(ips, ipnet.IP.String())
            }
        }
    }
    return ips
}

// 获取所有能检测到的本地ip, 返回第一个匹配 prefix 的ip, 否则返回空字符串
func GetLocalIPLinkPrefix(prefix string) string {
    ips := GetLocalIPs()
    le := len(prefix)
    if le == 0 {
        return ips[0]
    }

    for _, ip := range ips {
        if len(ip) >= le && ip[:le] == prefix {
            return ip
        }
    }
    return ""
}

// 构造监听器
func MakeTcpListen(conf *Config) (*Listen, error) {
    address := fmt.Sprintf("%s:%d", conf.BindIP, conf.BindPort)
    l, err := net.Listen("tcp", address)
    if err != nil {
        return nil, err
    }

    ip := conf.AdvertiseIP
    if ip == "" {
        ip = conf.BindIP
        if conf.AdvertiseIPPrefix != "" {
            ip = GetLocalIPLinkPrefix(conf.AdvertiseIPPrefix)
        }
        if ip == "" {
            ips := GetLocalIPs()
            if len(ips) > 1 {
                ip = ips[0]
            }
        }
    }

    port := conf.AdvertisePort
    if port == 0 {
        port = l.Addr().(*net.TCPAddr).Port
    }

    return &Listen{
        Listener:      l,
        AdvertiseIP:   ip,
        AdvertisePort: port,
    }, nil
}
