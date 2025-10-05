package conf

type SMS struct {
	Provider string
	Account  string
	Token    string
	AppConf  SMSAppConf `json:",default={}"` // 以AppName为名获取到到对应App的配置, 可从其中通过原因获取对应的模板
}

type SMSAppConf = map[string]CauseToTemplate

type CauseToTemplate = map[string]string
