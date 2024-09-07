package api

import (
	"LFshop-api/user-web/forms"
	"LFshop-api/user-web/global"
	"LFshop-api/user-web/global/reponse"
	"LFshop-api/user-web/middlewares"
	"LFshop-api/user-web/models"
	"LFshop-api/user-web/proto"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 这个函数实现翻译的json把前缀去掉
func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}

// grpc内部也有类型http状态码的code，grpc一旦报错，给前端返回的得是http状态码
// 这个函数实现把grpc的err转化成http状态码
func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	if err == nil {
		return
	}
	if e, ok := status.FromError(err); ok {
		switch e.Code() {
		case codes.NotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"msg": e.Message(),
			})
		case codes.Internal:
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg:": "内部错误",
			})
		case codes.InvalidArgument:
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "参数错误",
			})
		case codes.Unavailable:
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "用户服务不可用",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": e.Code(),
			})
		}
	}
	return
}

// 表单验证的逻辑，封装为一个函数
func HandleValidatorError(c *gin.Context, err error) {
	// 获取validator.ValidationErrors类型的errors
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		//非validator.ValidationErrors类型错误直接返回
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
		return
	}
	//validator.ValidationErrors类型错误则进行翻译
	c.JSON(http.StatusOK, gin.H{
		"msg": removeTopStruct(errs.Translate(global.Trans)),
	})
	return
}

func GetUserList(ctx *gin.Context) {
	// 获取访问的用户是谁（jwt校验通过后会把claims存进ctx里）
	claims, _ := ctx.Get("claims")
	currentUser := claims.(*models.CustomClaims)
	zap.S().Infof("访问用户: %d", currentUser.ID)

	// 接收前端想要的pn和psize
	Pn, _ := strconv.Atoi(ctx.DefaultQuery("Pn", "0"))
	PSize, _ := strconv.Atoi(ctx.DefaultQuery("PSize", "0"))
	// 远程调用方法

	rsp, err := global.UserSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(Pn),
		PSize: uint32(PSize),
	})
	if err != nil {
		zap.S().Errorw("查询用户列表失败")
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	zap.S().Info("获取用户列表页")

	// 接收数据
	result := make([]any, 0)
	for _, value := range rsp.Data {
		user := reponse.UserResponse{
			Id:       value.Id,
			NickName: value.NickName,
			//Birthday: time.Time(time.Unix(int64(value.BirthDay), 0)).Format("2006-01-02"),
			Birthday: reponse.JsonTime(time.Unix(int64(value.BirthDay), 0)),
			Gender:   value.Gender,
			Phone:    value.Phone,
		}
		result = append(result, user)
	}
	ctx.JSON(http.StatusOK, result)
}

// 密码登录
func PassWordLogin(c *gin.Context) {
	// 表单验证
	PwdForm := forms.PassWordLoginForm{}
	if err := c.ShouldBind(&PwdForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	// 图片验证码验证
	// store和此文件在同一个包下，可以直接用
	if !store.Verify(PwdForm.CaptchaId, PwdForm.Captcha, false) {
		c.JSON(http.StatusBadRequest, gin.H{
			"captcha": "验证码错误",
		})
		return
	}

	addr := fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port)
	userConn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("PassWordLogin服务连接失败",
			"msg", err,
		)
	}
	userSrvClient := proto.NewUserClient(userConn)
	rsp, err := userSrvClient.GetUserByPhone(context.Background(), &proto.PhoneRequest{Phone: PwdForm.Phone})
	// grpc调用返回错误，用grpc/status库来解析
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusBadRequest, map[string]string{
					"Phone": "用户不存在",
				})
			default:
				c.JSON(http.StatusInternalServerError, map[string]string{
					"Phone": "登录失败",
				})
			}
			return
		}
	} else {
		//只是查询到用户了而已，并没有检查密码
		pwdRsp, pwdErr := userSrvClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
			Password:          PwdForm.PassWord,
			EncryptedPassword: rsp.PassWord,
		})
		if pwdErr != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"password": "登录失败",
			})
		} else {
			if pwdRsp.Success {
				//生成token
				j := middlewares.NewJWT()
				claims := models.CustomClaims{
					ID:          uint(rsp.Id),
					NickName:    rsp.NickName,
					AuthorityId: uint(rsp.Role),
					StandardClaims: jwt.StandardClaims{
						NotBefore: time.Now().Unix(),               //签名的生效时间
						ExpiresAt: time.Now().Unix() + 60*60*24*30, //30天过期
						Issuer:    "LFshop",
					},
				}
				token, err := j.CreateToken(claims)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"msg": "生成token失败",
					})
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"id":         rsp.Id,
					"nick_name":  rsp.NickName,
					"token":      token,
					"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
				})
			} else {
				c.JSON(http.StatusBadRequest, map[string]string{
					"password": "登录失败",
				})
			}
		}
	}
}

// 用户注册
func Register(c *gin.Context) {
	registerForm := forms.RegisterForm{}
	if err := c.ShouldBind(&registerForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	//验证码
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})
	value, err := rdb.Get(registerForm.Phone).Result()
	if err == redis.Nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "验证码错误",
		})
		return
	} else {
		if value != registerForm.Code {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": "验证码错误",
			})
			return
		}
	}
	addr := fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port)
	userConn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("CreatedUser服务连接失败",
			"msg", err,
		)
	}
	userSrvClient := proto.NewUserClient(userConn)
	user, err := userSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		NickName: registerForm.Name,
		Phone:    registerForm.Phone,
		PassWord: registerForm.PassWord,
	})
	if err != nil {
		zap.S().Errorf("[Register] 查询 【新建用户失败】失败: %s", err.Error())
		HandleGrpcErrorToHttp(err, c)
		return
	}

	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:          uint(user.Id),
		NickName:    user.NickName,
		AuthorityId: uint(user.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),               //签名的生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*30, //30天过期
			Issuer:    "LFshop",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成token失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.Id,
		"nick_name":  user.NickName,
		"token":      token,
		"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
	})
}
