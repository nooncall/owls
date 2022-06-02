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

	DBFilter DBFilter `mapstructure:"db_filter" json:"db_filter" yaml:"db_filter"`

	Login Login `mapstructure:"login" json:"login" yaml:"login"`
}

type Login struct {
	Ldap struct {
		Host    string `json:"host"`
		Port    int    `json:"port"`
		BaseDn  string `json:"base_dn"`
		UseSll  bool   `json:"use_sll"`
		BindPwd string `json:"bind_pwd"`
		BindDn  string `json:"bind_dn"`
	}

	LoginPath          string `json:"login_path"`
	TokenSecret        string `json:"token_secret"`
	TokenEffectiveHour int    `json:"token_effective_hour"`
}

type DBFilter struct {
	LogLevel     string   `json:"log_level"`
	LogDir       string   `json:"log_dir"`
	NumOnceLimit int      `json:"num_once_limit"`
	ExecNoBackup bool     `json:"exec_no_backup"`
	AesKey       string   `json:"aes_key"`
	AesIv        string   `json:"aes_iv"`
	Reviewers    []string `json:"reviewers"`
	ReadNeedAuth bool     `json:"read_need_auth"`
}
