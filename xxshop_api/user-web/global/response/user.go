package reponse

import (
	"fmt"
	"time"
)

// 不能直接对 time.Time加方法，要加个新类型
// 调用时会自动转换这个方法
// 后端展示数据转成json交给前端时候 用struct做映射比较规范
type JsonTime time.Time

// 生成JsonTime类型的时候，调用转换json时，内部会直接调用这个方法
func (j JsonTime) MarshalJSON() ([]byte, error) {
	var stmp = fmt.Sprintf("\"%s\"", time.Time(j).Format("2006-01-02"))
	return []byte(stmp), nil
}

type UserResponse struct {
	Id       int32  `json:"id"`
	NickName string `json:"name"`
	//Birthday string `json:"birthday"`  //第一种方法 转成string
	//第二种方法 转成Jsontime型
	Birthday JsonTime `json:"birthday"` //自动生成json 格式化
	Gender   string   `json:"gender"`
	Mobile   string   `json:"mobile"`
}
