package config

import (
	"encoding/json"
	"fmt"
	"os"
)

//映射配置文件中的JSON数据
type Config struct {
	Server struct {//后端
		Port int    `json:"port"` // 端口号
		Host string `json:"host"` // 地址
	} `json:"server"`
	Database struct {//数据库
		Type     string `json:"type"`     // 数据库类型
		Host     string `json:"host"`     // 地址
		Port     int    `json:"port"`     // 端口号
		Username string `json:"username"` // 用户名
		Password string `json:"password"` // 密码
		DBName   string `json:"dbname"`   // 名称
	} `json:"database"`
	Admincode struct {//管理
		Invite_Code string `json:"invite_code"` //邀请码
	} `json:"admincode"`
}

//加载配置文件
func LoadConfig(filePath string) (*Config, error) {
	//打开配置文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开配置文件: %w", err)
	}
	//确保在函数结束时关闭文件
	defer file.Close()

	//创建一个Config结构体实例
	var config Config
	
	//创建JSON解码器
	decoder := json.NewDecoder(file)
	
	//解码JSON数据到config结构体
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return &config, nil
}