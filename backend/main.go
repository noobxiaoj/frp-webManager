package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	//"strings"
	"math/rand"

	"backend/config"
	"backend/DB"
	"backend/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"

)

//生成随机邀请码
func generateInviteCode() string {
	//字符集
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const codeLength = 5

	//创建随机数生成器
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	//生成随机邀请码
	b := make([]byte, codeLength)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	
	return string(b)
}
//生成API
func generateApiCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const codeLength = 10

	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, codeLength)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}

//检查管理员邀请码
func isAdminCodeUsed(adminCode string) (bool, error) {
    var count int64
    //查询数据库中是否存在使用该邀请码注册的用户
    err := DB.DB.Model(&models.User{}).
        Where("invitation_code = ?", adminCode).
        Count(&count).Error
    
    if err != nil {
        return false, fmt.Errorf("查询管理员邀请码使用状态失败: %v", err)
    }
    
    return count > 0, nil
}

//检查是否已存在管理员
func hasAdminUser() (bool, error) {
    var count int64
    
    //查询数据库中是否存在角色为admin的用户
    err := DB.DB.Model(&models.User{}).
        Where("role = ?", "admin").
        Count(&count).Error
    
    if err != nil {
        return false, fmt.Errorf("查询管理员用户失败: %v", err)
    }
    
    return count > 0, nil
}

func main() {
	//加载config
	cfg, err := config.LoadConfig("./config/config.json")
	if err != nil {
		log.Fatalf("加载: %v", err)
	}

	port := cfg.Server.Port
	Admin := cfg.Admincode.Invite_Code
	fmt.Printf("后端端口:%d \n", port)

	//初始化数据库连接
    DB.InitDB(
        cfg.Database.Username,
        cfg.Database.Password,
        cfg.Database.Host,
        fmt.Sprintf("%d", cfg.Database.Port),
        cfg.Database.DBName,
    )

	//创建Gin引擎
	r := gin.Default()

	//配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	
	//注册api
	r.POST("/api/register",func(c *gin.Context){
		//定义注册请求结构
		var registerRequest struct {
			Username       string `json:"username" binding:"required"`
			Password       string `json:"password" binding:"required"`
			Email          string `json:"email" binding:"required,email"`
			InvitationCode string `json:"invitation_code"`
		}

		if err := c.ShouldBindJSON(&registerRequest);err !=nil {
			c.JSON(http.StatusBadRequest,gin.H{"error":"无效的请求参数"})
			return
		}
		//检查用户名是否存在
		var existingUser models.User
		if err := DB.DB.Where("username = ?", registerRequest.Username).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
			return
		}
		//检查邮箱是否已注册
		if err := DB.DB.Where("email = ?", registerRequest.Email).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "邮箱已被注册"})
			return
		}
		
		Roleturn := "user"

		//验证邀请码（如果系统需要邀请码才能注册）
		if registerRequest.InvitationCode != "" {
			//在注册API中修改验证邀请码的部分
			if registerRequest.InvitationCode != "" {
			    //检查是否是管理员邀请码
			    if registerRequest.InvitationCode == Admin {
			        //检查管理员邀请码是否已被使用
			        adminUsed, err := isAdminCodeUsed(Admin)
			        if err != nil {
			            c.JSON(http.StatusInternalServerError, gin.H{"error": "验证邀请码失败"})
			            return
			        }
        
			        //检查是否已存在管理员用户
			        hasAdmin, err := hasAdminUser()
			        if err != nil {
			            c.JSON(http.StatusInternalServerError, gin.H{"error": "验证管理员状态失败"})
			            return
			        }
        
			        if adminUsed || hasAdmin {
			            c.JSON(http.StatusBadRequest, gin.H{"error": "管理员邀请码已被使用"})
			            return
			        }
			        
			        //设置用户角色为管理员
			        Roleturn = "admin"
			    } else {
			        //普通邀请码验证逻辑
			        var inviter models.User
			        if err := DB.DB.Where("invite_code = ?", registerRequest.InvitationCode).First(&inviter).Error; err != nil {
			            c.JSON(http.StatusBadRequest, gin.H{"error": "无效的邀请码"})
			            return
			        }
			        Roleturn = "user"
			    }
			} else {
			    // 如果没有提供邀请码，设置为普通用户
			    Roleturn = "user"
			}
		}

		//生成邀请码和API
		inviteCode := generateInviteCode()
		APICode := generateApiCode()

		//使用bcrypt加密密码
		hash, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
			return
		}

		//创建新用户
		newUser := models.User{
			Username:       registerRequest.Username,
			PasswordHash:   string(hash),
			Email:          registerRequest.Email,
			Role:           Roleturn, 
			APIKey:         APICode,
			InviteCode:     inviteCode,
			InvitationCode: registerRequest.InvitationCode,
			APIKeyCreatedAt:time.Now(),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		//存储用户信息
		if err := DB.DB.Create(&newUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "用户注册失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message":     "注册成功",
			"invite_code": inviteCode,
		})
	})

	//登录API
	r.POST("/api/login", func(c *gin.Context) {
		//定义登录请求结构
    	var loginRequest struct {
        	Username string `json:"username" binding:"required"`
        	Password string `json:"password" binding:"required"`
    	}

		//绑定并验证请求参数
    	if err := c.ShouldBindJSON(&loginRequest); err != nil {
        	c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
        	return
    	}

		//查找用户
    	var user models.User
    	if err := DB.DB.Where("username = ?", loginRequest.Username).First(&user).Error; err != nil {
        	if err == gorm.ErrRecordNotFound {
            	c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
            	return
        	}
        	c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户失败"})
        	return
    	}

		// 验证密码
    	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginRequest.Password)); err != nil {
        	c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
        	return
    	}

		//创建JWT令牌
    	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        	"user_id": user.ID,
        	"username": user.Username,
        	"role": user.Role,
        	"exp": time.Now().Add(time.Hour * 24).Unix(), // 令牌有效期24小时
    	})

		//使用密钥签名令牌
    	tokenString, err := token.SignedString([]byte("your_jwt_secret_key")) // 在实际应用中应该从配置中读取密钥
    	if err != nil {
        	c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
        	return
    	}

		//更新用户最后登录时间
    	DB.DB.Model(&user).Update("last_login", time.Now())

		// 返回令牌和用户信息
    	c.JSON(http.StatusOK, gin.H{
        	"token": tokenString,
        	"user": gin.H{
            	"id": user.ID,
            	"username": user.Username,
            	"email": user.Email,
            	"role": user.Role,
            	"invite_code": user.InviteCode,
            	"api_key": user.APIKey,
        	},
    	})

	})

	// 获取用户个人信息
	r.GET("/api/user/profile", func(c *gin.Context) {
    	// 验证token
    	tokenString := c.GetHeader("Authorization")
    	if tokenString == "" {
        	c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
        	return
    	}

    	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
    		tokenString = tokenString[7:]
    	}

    	// 解析JWT令牌
    	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            	return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        	}
        	return []byte("your_jwt_secret_key"), nil
    	})

    	if err != nil {
        	c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌"})
        return
    	}

    	claims, ok := token.Claims.(jwt.MapClaims)
    	if !ok || !token.Valid {
        	c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌"})
        	return
    	}

    	userID, ok := claims["user_id"].(float64)
    	if !ok {
        	c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的用户信息"})
        	return
    	}

    	// 查询用户信息
    	var user models.User
    	if err := DB.DB.First(&user, uint(userID)).Error; err != nil {
        	c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
        	return
    	}

    	// 返回用户信息
    	c.JSON(http.StatusOK, gin.H{
        	"user": gin.H{
            	"id":             user.ID,
            	"username":       user.Username,
            	"email":         user.Email,
            	"role":          user.Role,
            	"invite_code":   user.InviteCode,
            	"api_key":       user.APIKey,
            	"created_at":    user.CreatedAt,
            	"updated_at":    user.UpdatedAt,
        	},
    	})
	})

	// 修改个人信息
	r.PUT("/api/user/profile", func(c *gin.Context) {
    	// 验证token
    	tokenString := c.GetHeader("Authorization")
    	if tokenString == "" {
        	c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
        	return
    	}

    	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
        	tokenString = tokenString[7:]
    	}

    	// 解析JWT令牌
    	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            	return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        	}
        	return []byte("your_jwt_secret_key"), nil
    	})

    	if err != nil {
        	c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌"})
        	return
    	}

    	claims, ok := token.Claims.(jwt.MapClaims)
    	if !ok || !token.Valid {
        	c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌"})
        	return
    	}

    	userID, ok := claims["user_id"].(float64)
    	if !ok {
        	c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的用户信息"})
        	return
    	}

    	// 定义更新请求结构
    	var updateRequest struct {
        	Username    string `json:"username"`
        	OldPassword string `json:"old_password"`
        	NewPassword string `json:"new_password"`
    	}

    	// 绑定请求参数
    	if err := c.ShouldBindJSON(&updateRequest); err != nil {
        	c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
        	return
    	}

    	// 查询用户
    	var user models.User
    	if err := DB.DB.First(&user, uint(userID)).Error; err != nil {
        	c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
        	return
    	}

    	// 如果要更改密码，验证旧密码
    	if updateRequest.NewPassword != "" {
        	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(updateRequest.OldPassword)); err != nil {
            	c.JSON(http.StatusUnauthorized, gin.H{"error": "旧密码错误"})
            	return
        	}

        	// 生成新密码的hash
        	hash, err := bcrypt.GenerateFromPassword([]byte(updateRequest.NewPassword), bcrypt.DefaultCost)
        	if err != nil {
            	c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
            	return
        	}
        	user.PasswordHash = string(hash)
    	}

    	// 如果要更改用户名，检查是否已存在
    	if updateRequest.Username != "" && updateRequest.Username != user.Username {
        	var existingUser models.User
        	if err := DB.DB.Where("username = ?", updateRequest.Username).First(&existingUser).Error; err == nil {
            	c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
            	return
        	}
        	user.Username = updateRequest.Username
    	}

    	user.UpdatedAt = time.Now()

    	// 更新用户信息
    	if err := DB.DB.Save(&user).Error; err != nil {
        	c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户信息失败"})
        	return
    	}

    	c.JSON(http.StatusOK, gin.H{
        	"message": "个人信息更新成功",
        	"user": gin.H{
            "id":        user.ID,
            "username":  user.Username,
            "email":    user.Email,
            "updated_at": user.UpdatedAt,
        	},
    	})
	})

	//注册隧道API
	r.POST("/api/tunnel/register", func(c *gin.Context) {
		//从请求头中获取token
    	tokenString := c.GetHeader("Authorization")
    	if tokenString == "" {
        	c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
        	return
    	}

		//去掉Bearer前缀
    	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
        	tokenString = tokenString[7:]
    	}

		//解析JWT令牌
    	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            	return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        	}
        	return []byte("your_jwt_secret_key"), nil
    	})
		if err != nil {
        	c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌"})
        	return
    	}

		//验证令牌有效性并提取用户ID
    	claims, ok := token.Claims.(jwt.MapClaims)
    	if !ok || !token.Valid {
        	c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌"})
        	return
    	}

		//获取用户ID
    	userID, ok := claims["user_id"].(float64)
    	if !ok {
        	c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的用户信息"})
        	return
    	}

		//定义隧道注册请求结构
    	var tunnelRequest struct {
        	Name       string `json:"name" binding:"required"`
        	LocalIP    string `json:"local_ip" binding:"required,ip"`
        	LocalPort  int    `json:"local_port" binding:"required,min=1,max=65535"`
        	ServerPort int    `json:"server_port" binding:"required,min=1000,max=65535"`
        	Protocol   string `json:"protocol" binding:"required"`
			Remark     string `json:"remark"`
    	}

		//绑定并验证请求参数
    	if err := c.ShouldBindJSON(&tunnelRequest); err != nil {
        	c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数", "details": err.Error()})
        	return
    	}

		tcp_true := cfg.Protocols.TCP
		udp_true := cfg.Protocols.UDP
		http_true := cfg.Protocols.HTTP
		https_true := cfg.Protocols.HTTPS


		//验证协议类型
    	validProtocols := map[string]bool{
        	"tcp": tcp_true,
        	"udp": udp_true,
        	"http": http_true,
        	"https": https_true,
    	}

		protocols_show := ""

		if tcp_true&&protocols_show == ""{protocols_show += "tcp"}
		if udp_true&&protocols_show == ""{protocols_show += "udp"}else{protocols_show += ",udp"}
		if http_true&&protocols_show == ""{protocols_show += "http"}else{protocols_show += ",http"}
		if https_true&&protocols_show == ""{protocols_show += "https"}else{protocols_show += ",https"}

		error_show := ""
		if protocols_show != ""{error_show = "不支持的协议类型，支持的协议有：" + protocols_show}else{error_show = "网站协议开放错误，请联系管理人员"}

		//检查端口
		if !validProtocols[tunnelRequest.Protocol] {
        	c.JSON(http.StatusBadRequest, gin.H{"error": error_show})
        	return
    	}

		//检查服务器端口是否已被使用
    	var existingTunnel models.Tunnel
    	if err := DB.DB.Where("server_port = ?", tunnelRequest.ServerPort).First(&existingTunnel).Error; err == nil {
        	c.JSON(http.StatusConflict, gin.H{"error": "服务器端口已被占用"})
        	return
    	}

		//检查用户是否存在
    	var user models.User
    	if err := DB.DB.First(&user, uint(userID)).Error; err != nil {
        	c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
        	return
    	}
    
		//创建新隧道
    	newTunnel := models.Tunnel{
        	UserID:     uint(userID),
        	Name:       tunnelRequest.Name,
        	LocalIP:    tunnelRequest.LocalIP,
        	LocalPort:  tunnelRequest.LocalPort,
        	ServerPort: tunnelRequest.ServerPort,
        	Protocol:   tunnelRequest.Protocol,
        	Status:     "inactive", //初始未连接
			Remark:     tunnelRequest.Remark,
        	CreatedAt:  time.Now(),
        	UpdatedAt:  time.Now(),
    	}

		//保存隧道信息到数据库
    	if err := DB.DB.Create(&newTunnel).Error; err != nil {
        	c.JSON(http.StatusInternalServerError, gin.H{"error": "创建隧道失败", "details": err.Error()})
        	return
    	}

		year_out := cfg.RotExpired.YEAR
		month_out := cfg.RotExpired.MONTH
		day_out := cfg.RotExpired.DAY

		//创建端口使用记录
    	portUsage := models.PortUsage{
        	Port:      tunnelRequest.ServerPort,
        	UserID:    uint(userID),
        	TunnelID:  newTunnel.ID,
        	CreatedAt: time.Now(),
        	ExpiredAt: time.Now().AddDate(year_out, month_out, day_out),
    	}
		


		if err := DB.DB.Create(&portUsage).Error; err != nil {
        	//如果创建端口使用记录失败，回滚隧道创建
        	DB.DB.Delete(&newTunnel)
        	c.JSON(http.StatusInternalServerError, gin.H{"error": "创建端口使用记录失败"})
        	return
    	}

		//返回成功信息和隧道详情
    	c.JSON(http.StatusOK, gin.H{
        	"message": "隧道注册成功",
        	"tunnel": gin.H{
            	"id":          newTunnel.ID,
            	"name":        newTunnel.Name,
            	"local_ip":    newTunnel.LocalIP,
            	"local_port":  newTunnel.LocalPort,
            	"server_port": newTunnel.ServerPort,
            	"protocol":    newTunnel.Protocol,
            	"status":      newTunnel.Status,
				"remark":      newTunnel.Remark,
            	"created_at":  newTunnel.CreatedAt,
        	},
    	})

	})

	//获取用户隧道列表
	r.GET("/api/tunnels", func(c *gin.Context) {
		//从请求头中获取token
    	tokenString := c.GetHeader("Authorization")
    	if tokenString == "" {
        	c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
        	return
    	}

		//去掉Bearer前缀
    	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
        	tokenString = tokenString[7:]
    	}
    
    	//解析JWT令牌
    	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            	return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        	}
        	return []byte("your_jwt_secret_key"), nil
    	})
    
    	if err != nil {
        	c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌"})
        	return
    	}

		//验证令牌有效性并提取用户ID
    	claims, ok := token.Claims.(jwt.MapClaims)
    	if !ok || !token.Valid {
        	c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌"})
        	return
    	}

		//获取用户ID
    	userID, ok := claims["user_id"].(float64)
    	if !ok {
        	c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的用户信息"})
        	return
    	}

		//获取用户角色
    	role, _ := claims["role"].(string)
    
    	var tunnels []models.Tunnel

		//查询角色隧道
    	if role == "admin" {
			//如果是管理员，可以查看所有隧道
        	if err := DB.DB.Find(&tunnels).Error; err != nil {
            	c.JSON(http.StatusInternalServerError, gin.H{"error": "获取隧道列表失败"})
            	return
        	}
    	} else {
        	//普通用户只能查看自己的隧道
        	if err := DB.DB.Where("user_id = ?", uint(userID)).Find(&tunnels).Error; err != nil {
            	c.JSON(http.StatusInternalServerError, gin.H{"error": "获取隧道列表失败"})
            	return
        	}
    	}

		//构建响应数据
    	var tunnelList []gin.H
    	for _, tunnel := range tunnels {
        	tunnelList = append(tunnelList, gin.H{
            	"id":          tunnel.ID,
            	"name":        tunnel.Name,
            	"local_ip":    tunnel.LocalIP,
            	"local_port":  tunnel.LocalPort,
            	"server_port": tunnel.ServerPort,
            	"protocol":    tunnel.Protocol,
            	"status":      tunnel.Status,
				"remark":      tunnel.Remark,
            	"created_at":  tunnel.CreatedAt,
            	"updated_at":  tunnel.UpdatedAt,
        	})
    	}

		c.JSON(http.StatusOK, gin.H{
        	"tunnels": tunnelList,
        	"count":   len(tunnelList),
    	})
	})
	
	//删除隧道
	r.DELETE("/api/tunnel/:id", func(c *gin.Context) {
    	//从请求头中获取token
    	tokenString := c.GetHeader("Authorization")
    	if tokenString == "" {
        	c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
        	return
    	}

		//去掉Bearer前缀
    	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
        	tokenString = tokenString[7:]
    	}
    
		//解析JWT令牌
    	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            	return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        	}
        	return []byte("your_jwt_secret_key"), nil
    	})

		if err != nil {
        	c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌"})
        	return
    	}

		//验证令牌有效性并提取用户ID
    	claims, ok := token.Claims.(jwt.MapClaims)
    	if !ok || !token.Valid {
        	c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌"})
        	return
    	}

		//获取用户ID和角色
    	userID, ok := claims["user_id"].(float64)
    	if !ok {
        	c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的用户信息"})
        	return
    	}

		role, _ := claims["role"].(string)

		//获取隧道ID
    	tunnelID := c.Param("id")

		//查找隧道
    	var tunnel models.Tunnel
    	if err := DB.DB.First(&tunnel, tunnelID).Error; err != nil {
        	if err == gorm.ErrRecordNotFound {
            	c.JSON(http.StatusNotFound, gin.H{"error": "隧道不存在"})
        	} else {
            	c.JSON(http.StatusInternalServerError, gin.H{"error": "查询隧道失败"})
        	}
        	return
    	}

		//检查权限：只有隧道所有者或管理员可以删除
    	if tunnel.UserID != uint(userID) && role != "admin" {
        	c.JSON(http.StatusForbidden, gin.H{"error": "没有权限删除此隧道"})
        	return
    	}

		//开始事务
    	tx := DB.DB.Begin()

		//删除相关的端口使用记录
    	if err := tx.Where("tunnel_id = ?", tunnel.ID).Delete(&models.PortUsage{}).Error; err != nil {
        	tx.Rollback()
        	c.JSON(http.StatusInternalServerError, gin.H{"error": "删除端口使用记录失败"})
        	return
    	}

		//删除隧道
    	if err := tx.Delete(&tunnel).Error; err != nil {
        	tx.Rollback()
        	c.JSON(http.StatusInternalServerError, gin.H{"error": "删除隧道失败"})
        	return
    	}

		//提交事务
    	tx.Commit()

		c.JSON(http.StatusOK, gin.H{
        	"message": "隧道删除成功",
    	})
	})

	r.PUT("/api/tunnel/:id", func(c *gin.Context) {
		// 验证token
    	tokenString := c.GetHeader("Authorization")
    	if tokenString == "" {
        	c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
        	return
    	}

		// 去掉Bearer前缀
    	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
        	tokenString = tokenString[7:]
    	}

		// 解析JWT令牌
    	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            	return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        	}
        	return []byte("your_jwt_secret_key"), nil
    	})

		if err != nil {
        	c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌"})
        	return
    	}

		// 验证令牌并获取用户信息
    	claims, ok := token.Claims.(jwt.MapClaims)
    	if !ok || !token.Valid {
        	c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌"})
        	return
    	}

		userID, ok := claims["user_id"].(float64)
    	if !ok {
        	c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的用户信息"})
        	return
    	}

		role, _ := claims["role"].(string)

		// 获取隧道ID
    	tunnelID := c.Param("id")

		// 定义更新请求结构
    	var updateRequest struct {
        	Name      string `json:"name"`
        	LocalIP   string `json:"local_ip"`
        	LocalPort int    `json:"local_port"`
        	Protocol  string `json:"protocol"`
    	}

		// 绑定请求参数
    	if err := c.ShouldBindJSON(&updateRequest); err != nil {
        	c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
        	return
    	}

		// 查找现有隧道
    	var tunnel models.Tunnel
    	if err := DB.DB.First(&tunnel, tunnelID).Error; err != nil {
        	c.JSON(http.StatusNotFound, gin.H{"error": "隧道不存在"})
        	return
    	}

		// 检查权限
    	if tunnel.UserID != uint(userID) && role != "admin" {
        	c.JSON(http.StatusForbidden, gin.H{"error": "没有权限修改此隧道"})
        	return
    	}

		// 更新隧道信息
    	updates := models.Tunnel{
        	Name:      updateRequest.Name,
        	LocalIP:   updateRequest.LocalIP,
        	LocalPort: updateRequest.LocalPort,
        	Protocol:  updateRequest.Protocol,
        	UpdatedAt: time.Now(),
    	}

		if err := DB.DB.Model(&tunnel).Updates(updates).Error; err != nil {
        	c.JSON(http.StatusInternalServerError, gin.H{"error": "更新隧道失败"})
        	return
    	}

		c.JSON(http.StatusOK, gin.H{
        	"message": "隧道更新成功",
    		"tunnel": gin.H{
            	"id":          tunnel.ID,
            	"name":        updates.Name,
            	"local_ip":    updates.LocalIP,
            	"local_port":  updates.LocalPort,
            	"server_port": tunnel.ServerPort,
            	"protocol":    updates.Protocol,
            	"status":      tunnel.Status,
            	"updated_at":  updates.UpdatedAt,
        	},
    	})
	})





	//时间api
	r.GET("/api/time", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"timestamp":    time.Now().Format(time.RFC3339),
		})
	})
	//端口api
	r.GET("/api/port", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"port": port,
		})
	})

	// 启动服务器
	r.Run(fmt.Sprintf(":%d", port))
}
