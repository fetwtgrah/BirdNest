package controller

import (
	"context"
	"errors"
	"fmt"
	"github/fetwtgrah/BirdNest/configs"
	"github/fetwtgrah/BirdNest/model"
	"github/fetwtgrah/BirdNest/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-gomail/gomail"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type userController struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
}
type UserRegister struct {
	*model.User
	Code string `json:"code" binding:"required"`
}

func SendCode(c *gin.Context) {
	var user userController
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}
	code := utils.CreatCode()

	//redis存储验证码
	key := fmt.Sprintf("user:email:%s", user.Email)
	if err := configs.Rc.Set(context.Background(), key, code, 5*time.Minute).Err(); err != nil {
		c.JSON(400, gin.H{"msg": "邮箱储存出错"})
		return
	}
	m := gomail.NewMessage()
	m.SetHeader("From", "2037461470@qq.com")
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "this is a test email form Fresh...")
	m.SetBody("text/html", fmt.Sprintf(utils.EmailContext(), user.Name, code))
	d := gomail.NewDialer(
		configs.Conf.Email.Smtp_host,
		configs.Conf.Email.Smtp_port,
		configs.Conf.Email.Smtp_user,
		configs.Conf.Email.Smtp_password)
	if err := d.DialAndSend(m); err != nil {
		configs.Rc.Del(context.Background(), key)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "邮件发送失败，请稍后重试"})
		return
	}

	c.JSON(200, gin.H{"msg": "验证码发送成功"})
}

func RegisterByCode(c *gin.Context) {
	var user UserRegister
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}
	key := fmt.Sprintf("user:email:%s", user.Email)
	RC_code, err := configs.Rc.Get(context.Background(), key).Result()
	if errors.Is(err, redis.Nil) {
		c.JSON(400, gin.H{"msg": "验证码已过期，请重新获取"})
		return
	} else if err != nil {
		c.JSON(500, gin.H{"msg": "服务异常，请稍后重试"})
		return
	}
	if RC_code != user.Code {
		c.JSON(400, gin.H{"msg": "验证码错误！"})
		return
	}
	configs.Rc.Del(context.Background(), key)
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(400, gin.H{"msg": "系统出错"})
		return
	}
	err = configs.Db.Create(&model.User{Name: user.Name, Email: user.Email, Password: string(password)}).Error
	if err != nil {
		c.JSON(400, gin.H{"msg": "用户创建失败，请稍后重试"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "用户创建成功"})
}

func LoginByPassword(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	var RcUser model.User
	err := configs.Db.Where("email = ?", user.Email).First(&RcUser).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	//检查密码
	if user.Password != RcUser.Password {
		c.JSON(400, gin.H{
			"msg": "密码错误，请重新输入...",
		})
		return
	}
	//签发token
	token, _ := utils.GenerateToken(user)
	key := fmt.Sprintf("user:name:%s", user.Name)
	err = configs.Rc.Set(context.Background(), key, token, 1*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "token签发失败"})
		return
	}
	c.JSON(200, gin.H{
		"token": token,
		"msg":   "成功登录！",
	})
}
