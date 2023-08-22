package api

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
	"strings"
	"time"
	"xxshop-api/user-web/forms"
	"xxshop-api/user-web/global"
	"xxshop-api/user-web/global/response"
	"xxshop-api/user-web/middlewares"
	"xxshop-api/user-web/models"
	"xxshop-api/user-web/proto"
)

func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() { //code 拿到服务器访问状态码
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{ //返回前端
					"msg": e.Message(), //message拿到服务器访问状态
				})
			case codes.Internal: //不要把grpc错误暴露给用户，不友好
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
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
					"msg": "其他错误" + e.Message(), //内部错误
				})
			}
			return
		}
	}
}

func GetUserList(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")               //从gin中拿到claims类型
	currentUser := claims.(*models.CustomClaims) //断言 进行类型转换
	zap.S().Infof("访问用户: %d", currentUser.ID)
	zap.S().Infof("访问用户权限: %d", currentUser.AuthorityId)
	//调用接口
	pn := ctx.DefaultQuery("pn", "0") //get类型参数都是放在gin中
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("psize", "10") //默认每页10个
	pSizeInt, _ := strconv.Atoi(pSize)
	rsp, err := global.UserSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询 【用户列表失败】")
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	//从这里发起连接 grpc里面取数据出来
	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		user := reponse.UserResponse{
			Id:       value.Id,
			NickName: value.NickName,
			Birthday: reponse.JsonTime(time.Unix(int64(value.BirthDay), 0)),
			Gender:   value.Gender,
			Mobile:   value.Mobile,
		}
		result = append(result, user)
	}
	ctx.JSON(http.StatusOK, result)
}

func HandleValidatorError(ctx *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	ctx.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)),
	})
	return
}

func PassWordLogin(ctx *gin.Context) {
	//表单验证
	passWordLoginForm := forms.PasswordLoginForm{}
	if err := ctx.ShouldBind(&passWordLoginForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}
	//验证图片验证码  参数是从前端传来的
	if !store.Verify(passWordLoginForm.CaptchaId, passWordLoginForm.Captcha, true) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"captcha": "验证码错误",
		})
		return
	}
	//登录的逻辑
	//forms.PasswordLoginForm结构体里的成员是前端传来的信息
	//proto.MobileRequest是通过grpc通信从srv服务那边拿到的信息
	//设置rsp接收,通过srv服务查数据库得到的用户信息
	rsp, err := global.UserSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: passWordLoginForm.Mobile,
	})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusBadRequest, map[string]string{
					"mobile": "用户不存在",
				})
			default:
				ctx.JSON(http.StatusInternalServerError, map[string]string{
					"mobile": "登录失败",
				})
			}
			return
		}
	} else {
		//只是查询到用户而已，没有检查密码
		//调用server层的检查密码接口
		if passRsp, passErr := global.UserSrvClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
			//passWordLoginForm.PassWord: forms中没有加密的密码
			Password:          passWordLoginForm.PassWord, //前端传来的密码
			EncryptedPassword: rsp.PassWord,               //后端查询数据库的加密密码
		}); passErr != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]string{
				"password": "登录失败",
			})
		} else {
			if passRsp.Success {
				//生成token  指定jwt
				//middlewares.JWT{
				//	SigningKey:  ,
				//}
				//生成token config中配置
				j := middlewares.NewJWT()
				//对指定model签名
				calims := models.CustomClaims{
					ID:          uint(rsp.Id), //从user信息中取id
					NickName:    rsp.NickName, //将nickname返回给前端，让其在用户中心显示
					AuthorityId: uint(rsp.Role),
					StandardClaims: jwt.StandardClaims{
						NotBefore: time.Now().Unix(),               //签名的生效时间
						ExpiresAt: time.Now().Unix() + 60*60*24*30, //30天过期
						Issuer:    "xxshop",                        //签名机构
					},
				}
				token, err := j.CreateToken(calims)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"msg": "生成token失败",
					})
					return
				}

				ctx.JSON(http.StatusOK, gin.H{
					"id":         rsp.Id,
					"nick_name":  rsp.NickName,
					"token":      token,
					"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
				})
			} else {
				ctx.JSON(http.StatusBadRequest, map[string]string{
					"password": "登录失败",
				})
			}
		}
	}
}

// 用户注册
func Register(ctx *gin.Context) {
	//接收前端传来的结构体值
	registerForm := forms.RegisterFrom{}
	if err := ctx.ShouldBind(&registerForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}
	//验证码校验
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.Serverconfig.RedisInfo.Host,
			global.Serverconfig.RedisInfo.Port),
	})
	value, err := rdb.Get(context.Background(), registerForm.Mobile).Result()
	if err == redis.Nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": "验证码错误",
		})
		return
	} else {
		if value != registerForm.Code {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code": "验证码错误",
			})
		}
	}

	//生成grpc的client并调用接口
	//新建用户
	user, err := global.UserSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		NickName: registerForm.Mobile,
		PassWord: registerForm.PassWord,
		Mobile:   registerForm.Mobile,
	})
	if err != nil {
		zap.S().Errorf("[Register] 查询 【新建用户失败】:%s", err.Error())
		HandleGrpcErrorToHttp(err, ctx)
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
			Issuer:    "xxshop",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成token失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":         user.Id,
		"nick_name":  user.NickName,
		"token":      token,
		"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
	})
}
