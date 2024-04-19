package dong

// requestBody 咚咚通知请求报文
type requestBody struct {
	UserIDs  string `json:"userIds"`   // 消息接收用户，多个 以,隔开（groupIds 为 空时必填）
	GroupIDs string `json:"group_ids"` // 消息接收群，多个以,隔开（userIds 为空时必填）
	Title    string `json:"title"`     // 标题（必填）长度：300 个字符（150 个中文）
	Detail   string `json:"detail"`    // 卡片消息详细（必填）长度限制：2000 个字符（中文：1000 个）支持 html 标签
}
