package forms

// 存放与验证有关的信息  //成员是前端传来的信息
type PasswordLoginForm struct {
	//用josn或者form都行  binding代表前端显示约束
	Mobile   string `form:"mobile" json:"mobile" binding:"required,mobile"`           //手机号码格式有规范可寻，自定义validator
	PassWord string `form:"password" json:"password" binding:"required,min=3,max=20"` //binding条件里不能加空格
	//验证码长度
	Captcha   string `form:"captcha" json:"captcha" binding:"required,min=5,max=5"`
	CaptchaId string `form:"captcha_id" json:"captcha_id" binding:"required"`
}

type RegisterFrom struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required,mobile"`
	PassWord string `form:"password" json:"password" binding:"required,min=3,max=20"`
	Code     string `form:"code" json:"code" binding:"required,min=6,max=6"`
}
