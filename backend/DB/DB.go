package DB

import (
    "fmt"
    "log"

    "backend/models"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

var DB *gorm.DB

//初始化数据库连接
func InitDB(username, password, host, port, dbname string) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        username, password, host, port, dbname)

    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    //自动迁移数据库表
    err = db.AutoMigrate(
        &models.User{},
        &models.Tunnel{},
        &models.PortUsage{},
        &models.ClientConnection{},
    )
    if err != nil {
        log.Fatalf("Failed to migrate database: %v", err)
    }

    DB = db
    log.Println("Database connected and migrated successfully")
}