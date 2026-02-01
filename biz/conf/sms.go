package conf

type SMS struct {
	MaxInPeriod int               // 周期内最多发送的次数
	Period      int               // 以分钟为单位
	Provider    string            // 供应商
	Account     string            // 账号
	Token       string            // 鉴权token
	Extra       map[string]string // 不同渠道的额外信息
	AppConf     SMSAppConf        `json:",default={}"` // 以AppName为名获取到到对应App的配置, 可从其中通过原因获取对应的模板
}

type SMSAppConf = map[string]CauseToTemplate

type CauseToTemplate = map[string]string
