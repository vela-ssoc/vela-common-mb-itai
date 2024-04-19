package sonatype

import (
	"context"
	"encoding/base64"
)

// Configurer 配置文件接口
type Configurer interface {
	// Load 获取 sonatype 的查询 url 和认证信息。
	// 虽然无认证也可以查询 sonatype 漏洞库，但是容易被限制调用频率。
	// 生产环境中请尽量使用认证查询，以免产生不必要的麻烦。
	Load(ctx context.Context) (addr, auth string, err error)
}

// HardConfig 这是一个随意注册的 sonatype 账号，该账号只在此程序中查询漏洞。
// 请勿充值或他用！！！
// 请勿充值或他用！！！
// 请勿充值或他用！！！
// 用户名：yuhs2s8t3o12j@ihotmails.com
// 密码：wA3%iK2{lO6(hB3[oZ5_mX8*aG
// 该用户名是无需注册认证的临时邮箱注册，可以被他人任意申请（所以名字起的很随机），并不安全。
// 临时邮箱注册网址是：https://ihotmails.com/
func HardConfig() Configurer {
	uname := "yuhs2s8t3o12j@ihotmails.com"
	passwd := "wA3%iK2{lO6(hB3[oZ5_mX8*aG"
	addr := "https://ossindex.sonatype.org/api/v3/authorized/component-report"

	auth := uname + ":" + passwd
	enc := base64.StdEncoding.EncodeToString([]byte(auth))

	return &hardConfig{
		addr: addr,
		auth: "Basic " + enc,
	}
}

type hardConfig struct {
	addr string
	auth string
}

func (hc *hardConfig) Load(context.Context) (string, string, error) {
	return hc.addr, hc.auth, nil
}
