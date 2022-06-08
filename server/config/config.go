package config

type Server struct {
	JWT     JWT     `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Zap     Zap     `mapstructure:"zap" json:"zap" yaml:"zap"`
	Redis   Redis   `mapstructure:"redis" json:"redis" yaml:"redis"`
	Email   Email   `mapstructure:"email" json:"email" yaml:"email"`
	Casbin  Casbin  `mapstructure:"casbin" json:"casbin" yaml:"casbin"`
	System  System  `mapstructure:"system" json:"system" yaml:"system"`
	Captcha Captcha `mapstructure:"captcha" json:"captcha" yaml:"captcha"`
	// auto
	AutoCode Autocode `mapstructure:"autoCode" json:"autoCode" yaml:"autoCode"`
	// gorm
	Mysql  Mysql `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Pgsql  Pgsql `mapstructure:"pgsql" json:"pgsql" yaml:"pgsql"`
	DBList []DB  `mapstructure:"db-list" json:"db-list" yaml:"db-list"`
	// oss
	Local      Local      `mapstructure:"local" json:"local" yaml:"local"`
	Qiniu      Qiniu      `mapstructure:"qiniu" json:"qiniu" yaml:"qiniu"`
	AliyunOSS  AliyunOSS  `mapstructure:"aliyun-oss" json:"aliyunOSS" yaml:"aliyun-oss"`
	HuaWeiObs  HuaWeiObs  `mapstructure:"hua-wei-obs" json:"huaWeiObs" yaml:"hua-wei-obs"`
	TencentCOS TencentCOS `mapstructure:"tencent-cos" json:"tencentCOS" yaml:"tencent-cos"`
	AwsS3      AwsS3      `mapstructure:"aws-s3" json:"awsS3" yaml:"aws-s3"`

	Excel Excel `mapstructure:"excel" json:"excel" yaml:"excel"`
	Timer Timer `mapstructure:"timer" json:"timer" yaml:"timer"`

	// 跨域配置
	Cors CORS `mapstructure:"cors" json:"cors" yaml:"cors"`

	DBFilter DBFilter `mapstructure:"db-filter" json:"db-filter" yaml:"db-filter"`

	Login Login `mapstructure:"login" json:"login" yaml:"login"`
}

type Login struct {
	Model string `mapstructure:"model"`

	Ldap struct {
		Host    string `mapstructure:"host"`
		Port    int    `mapstructure:"port"`
		BaseDn  string `mapstructure:"base-dn"`
		UseSll  bool   `mapstructure:"use-sll"`
		BindPwd string `mapstructure:"bind-pwd"`
		BindDn  string `mapstructure:"bind-dn"`
	} `yaml:ldap`

	LoginPath          string `mapstructure:"login-path"`
	TokenSecret        string `mapstructure:"token-secret"`
	TokenEffectiveHour int    `mapstructure:"token-effective-hour"`
}

type DBFilter struct {
	LogLevel     string   `mapstructure:"log-level"`
	LogDir       string   `mapstructure:"log-dir"`
	NumOnceLimit int      `mapstructure:"num-once-limit"`
	ExecNoBackup bool     `mapstructure:"exec-no-backup"`
	AesKey       string   `mapstructure:"aes-key"` //初始化全是空的
	AesIv        string   `mapstructure:"aes-iv"`
	Reviewers    []string `mapstructure:"reviewers"`
	ReadNeedAuth bool     `mapstructure:"read-need-auth"`
}
