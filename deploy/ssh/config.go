package ssh

type Config struct {
	// SSH 主机。
	// 零值时默认为 "localhost"。
	Host string `json:"host,omitempty" yaml:"Host" xml:"Host" env:"SSH_HOST"`
	// SSH 端口。
	// 零值时默认为 22。
	Port int32 `json:"port,omitempty" yaml:"Port" xml:"Port" env:"SSH_PORT"`
	// SSH 登录用户名。
	Username string `json:"username,omitempty" yaml:"Username" xml:"Username" env:"SSH_USERNAME"`
	// SSH 登录密码。
	Password string `json:"password,omitempty" yaml:"Password" xml:"Password" env:"SSH_PASSWORD"`
	// SSH 登录私钥。
	Key string `json:"key,omitempty" yaml:"Key" xml:"Key" env:"SSH_KEY"`
	// SSH 登录私钥口令。
	KeyPassphrase string `json:"key_passphrase,omitempty" yaml:"KeyPassphrase" xml:"KeyPassphrase" env:"SSH_KEY_PASSPHRASE"`
	// 是否回退使用 SCP。
	UseSCP bool `json:"use_scp,omitempty" yaml:"UseSCP" xml:"UseSCP" env:"SSH_USE_SCP"`
	// 前置命令。
	PreCommand string `json:"pre_command,omitempty" yaml:"PreCommand" xml:"PreCommand" env:"SSH_PRE_COMMAND"`
	// 后置命令。
	PostCommand string `json:"post_command,omitempty" yaml:"PostCommand" xml:"PostCommand" env:"SSH_POST_COMMAND"`

	// 输出证书文件路径。
	CertPath string `json:"cert_path,omitempty" yaml:"CertPath" xml:"CertPath" env:"SSH_CERT_PATH"`
	// 输出私钥文件路径。
	KeyPath string `json:"key_path,omitempty" yaml:"KeyPath" xml:"KeyPath" env:"SSH_KEY_PATH"`
}
