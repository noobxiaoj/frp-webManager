package models

import (
	"time"
)

// 用户表
type User struct {
	ID              uint   `gorm:"primarykey"` //主键
	Username        string `gorm:"type:varchar(50);not null;unique"` //用户名
	PasswordHash    string `gorm:"type:varchar(255);not null"` //Hash密码
	Email           string `gorm:"type:varchar(100);unique"` //邮箱
	Role            string `gorm:"type:varchar(20);default:'user'"` //权限
	APIKey          string `gorm:"type:varchar(64);unique"` //唯一API
	InviteCode      string `gorm:"type:varchar(20)"` //邀请码
	InvitationCode  string `gorm:"type:varchar(20)"` //被邀请码
	APIKeyCreatedAt time.Time //API创造时间
	CreatedAt       time.Time //账号创造时间
	UpdatedAt       time.Time //账号更新时间
}

// 隧道表
type Tunnel struct {
	ID         uint   `gorm:"primarykey"` //主键
	UserID     uint   `gorm:"not null"` //建立这个隧道的ID
	Name       string `gorm:"type:varchar(50);not null"` //隧道名字
	LocalIP    string `gorm:"type:varchar(50);not null"` //本地IP
	LocalPort  int    `gorm:"not null"` //本地端口
	ServerPort int    `gorm:"not null;unique"` //服务器端口
	Protocol   string `gorm:"type:varchar(10);not null"` //协议类型
	Status     string `gorm:"type:varchar(20);default:'inactive'"` //是否在线
	Remark     string `gorm:"not null"` //隧道备注
	User       User   `gorm:"foreignKey:UserID"` //用户结构体
	CreatedAt  time.Time //隧道创造时间
	UpdatedAt  time.Time //隧道更新时间
}

// 端口使用记录表
type PortUsage struct {
	ID        uint   `gorm:"primarykey"` //主键
	Port      int    `gorm:"not null"` //服务器端口
	UserID    uint   `gorm:"not null"` //使用的用户ID
	TunnelID  uint   `gorm:"not null"` //隧道表ID
	Protocol  string `gorm:"type:varchar(10);not null"` //协议类型
	User      User   `gorm:"foreignKey:UserID"` //用户结构体
	Tunnel    Tunnel `gorm:"foreignKey:TunnelID"` //隧道结构体
	CreatedAt time.Time //端口使用创造时间
	ExpiredAt time.Time //端口超时时间
}

// 客户端连接记录表
type ClientConnection struct {
	ID             uint   `gorm:"primarykey"` //主键
	UserID         uint   `gorm:"not null"` //用户ID
	ClientIP       string `gorm:"type:varchar(50);not null"` //客户端IP
	ClientVersion  string `gorm:"type:varchar(20)"` //客户端版本
	Status         string `gorm:"type:varchar(20);default:'offline'"` //是否在线
	User           User   `gorm:"foreignKey:UserID"` //用户表
	ConnectedAt    time.Time //客户端连接时间
	DisconnectedAt time.Time //客户端断开时间
}
