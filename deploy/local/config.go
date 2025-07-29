package local

type Config struct {
	// Shell 执行环境。
	// 前置命令。
	PreCommand string `json:"pre_command,omitempty" yaml:"PreCommand" xml:"PreCommand" env:"LOCAL_PRE_COMMAND"`
	// 后置命令。
	PostCommand string `json:"post_command,omitempty" yaml:"PostCommand" xml:"PostCommand" env:"LOCAL_POST_COMMAND"`
	// 输出证书文件路径。
	CertPath string `json:"cert_path,omitempty" yaml:"CertPath" xml:"CertPath" env:"LOCAL_CERT_PATH"`
	// 输出私钥文件路径。
	KeyPath string `json:"key_path,omitempty" yaml:"KeyPath" xml:"KeyPath" env:"LOCAL_KEY_PATH"`
}
