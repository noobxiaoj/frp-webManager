package frps

import (
    "github.com/fatedier/frp/server"
    "github.com/fatedier/frp/pkg/config"
    "github.com/fatedier/frp/pkg/auth"
)

func main() {
    // 创建基础配置
    cfg := config.ServerCommonConf{
        BindAddr: "0.0.0.0",
        BindPort: 7000,
        // 认证token
        AuthToken: "12345678",
        
        // 可选的dashboard配置
        DashboardPort: 7500,
        DashboardUser: "admin",
        DashboardPwd: "admin",
    }

    // 创建frp server实例
    svr, err := server.NewService(cfg)
    if err != nil {
        panic(err)
    }

    // 运行服务
    err = svr.Run()
    if err != nil {
        panic(err)
    }
}