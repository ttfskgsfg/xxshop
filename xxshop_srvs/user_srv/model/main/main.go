package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

func main() {
	//dsn := "root:123456@tcp(localhost:8086)/xxshop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"
	//newLogger := logger.New(
	//	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	//	logger.Config{
	//		SlowThreshold: time.Second, // 慢 SQL 阈值
	//		LogLevel:      logger.Info, // Log level
	//		Colorful:      true,        // 禁用彩色打印
	//	},
	//)
	//
	//// 全局模式
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
	//	NamingStrategy: schema.NamingStrategy{
	//		SingularTable: true,
	//	},
	//	Logger: newLogger,
	//})
	//if err != nil {
	//	panic(err)
	//}
	//_ = db.AutoMigrate(&model.User{})
	//salt, encodedPwd := password.Encode("generic password", nil)
	//fmt.Println(salt)
	//fmt.Println(encodedPwd)
	//check := password.Verify("generic password", salt, encodedPwd, nil)
	//fmt.Println(check)

	// Using custom options
	//options := &password.Options{16, 100, 32, sha512.New}
	//salt, encodedPwd := password.Encode("generic password", options)
	//newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	//fmt.Println(len(newPassword)) //长度不要超过100，否则数据保存失败
	//fmt.Println(newPassword)
	//解析存储的密码
	//passwordInfo := strings.Split(newPassword, "$")
	//fmt.Println(passwordInfo)
	//check := password.Verify("generic password", passwordInfo[2], passwordInfo[3], options)
	//fmt.Println(check) // true
}

// 生成md5模版 加盐值 随机字符串+用户密码
func genMd5(code string) string {
	Md5 := md5.New()
	_, _ = io.WriteString(Md5, code)
	return hex.EncodeToString(Md5.Sum(nil))
}
