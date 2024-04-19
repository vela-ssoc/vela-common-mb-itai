package dong

import "fmt"

// responseBody 咚咚通知响应结果
type responseBody struct {
	Code string `json:"code"` // 请求返回码
	Msg  string `json:"msg"`  // 请求返回消息
}

func (r responseBody) Error() error {
	if r.Code == "200" {
		return nil
	}

	return &dongError{
		Code:  r.Code,
		Msg:   r.Msg,
		Cause: codes[r.Code],
	}
}

type dongError struct {
	Code  string
	Msg   string
	Cause string
}

// Error error
func (de *dongError) Error() string {
	return fmt.Sprintf("error code: %s, msg: %s", de.Code, de.Msg)
}

var codes = map[string]string{
	"4000009": "未登陆",
	"4000001": "参数无效",
	"4000002": "权限不足",
	"4000003": "验证码错误",
	"4000004": "登陆频率过高",
	"4000005": "验证码已失效",
	"4000006": "没有找到资源",
	"4000007": "系统错误，请联系管理员",
	"4000008": "无权访问",
	"4000010": "上传提交的数据已经在审批中",
	"4000012": "消息类型不匹配",
	"4000013": "群必须包含 3 个以上成员",
	"4000014": "账号已失效",
	"4000015": "认证失败",
	"5000001": "用户不存在",
	"5000002": "用户名密码不匹配",
	"5000003": "没有权限登录",
	"5000004": "无该用户实名信息",
	"5000005": "文件格式不正确",
	"5000006": "用户越权操作",
	"5000009": "服务访问频率受限",
	"5000010": "服务 IP 不在白名单内",
}
