package reponse

import (
	"fmt"
	"time"
)

type JsonTime time.Time

// 想让time.Time格式在json化的时候变成想要的格式
// 重新type成JsonTime，然后给JsonTime定义固定格式的MarshalJSON方法
// 有点类似重载，然后JsonTime在json化的时候就会变成指定的格式
func (j JsonTime) MarshalJSON() ([]byte, error) {
	var stmp = fmt.Sprintf("\"%s\"", time.Time(j).Format("2006-01-02"))
	return []byte(stmp), nil
}

type UserResponse struct {
	Id       int32  `json:"id"`
	NickName string `json:"name"`
	//Birthday string `json:"birthday"`
	Birthday JsonTime `json:"birthday"`
	Gender   string   `json:"gender"`
	Phone    string   `json:"phone"`
}
