package ssoauth

import "fmt"

var ssoCodes = map[string]string{
	"01": "密码错误",
	"03": "用户不存在",
	"04": "设备类型错误",
	"09": "用户名或密码为空",
}

// reply sso 认证服务的响应报文
type ssoReply struct {
	Code string `json:"rspCde"` // 业务响应码
	Text string `json:"rspMsg"` // 响应消息
}

func (sr ssoReply) Error() error {
	code := sr.Code
	if code == "00" {
		return nil
	}

	msg, ok := ssoCodes[code]
	if !ok {
		msg = sr.Text
	}

	return fmt.Errorf("sso 认证错误：%s", msg)
}
