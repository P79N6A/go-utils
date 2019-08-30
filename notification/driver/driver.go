package driver

type Driver interface {
	NewNotification(config interface{}) (Notification, error)
}

type Contact struct {
	Mail        []string `json:"mail,omitempty"`
	MailGroup   []string `json:"mail_group,omitempty"`
	Sms         []string `json:"sms,omitempty"`
	SmsGroup    []string `json:"sms_group,omitempty"`
	Ivr         []string `json:"ivr,omitempty"`
	IvrGroup    []string `json:"ivr_group,omitempty"`
	Wechat      []string `json:"wechat,omitempty"`
	WechatGroup []string `json:"wechat_group,omitempty"`
	Weibo       []string `json:"weibo,omitempty"`
	WeiboGroup  []string `json:"weibo_group,omitempty"`
	Push        []string `json:"push,omitempty"`
	PushGroup   []string `json:"push_group,omitempty"`
}

type AsyncSend struct {
}

type NotificationMessage struct {
	Subject         string
	Content         string
	IsHtmlFormatter bool
	Extra           map[string]interface{}
}

type Notification interface {
	SyncSend(contact Contact, notificationMessage NotificationMessage) error

	AsyncSend(contact Contact, notificationMessage NotificationMessage) (AsyncSend, error)
}
