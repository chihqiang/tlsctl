package fc

import (
	"fmt"
	aliopen "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	alifc3 "github.com/alibabacloud-go/fc-20230330/v4/client"
	alifc2 "github.com/alibabacloud-go/fc-open-20210406/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"strings"
)

type Clients struct {
	FC2 *alifc2.Client
	FC3 *alifc3.Client
}

func NewClients(accessKeyId, accessKeySecret, region string) (*Clients, error) {
	// 接入点一览 https://api.aliyun.com/product/FC-Open
	var fc2Endpoint string
	switch region {
	case "":
		fc2Endpoint = "fc.aliyuncs.com"
	case "cn-hangzhou-finance":
		fc2Endpoint = fmt.Sprintf("%s.fc.aliyuncs.com", region)
	default:
		fc2Endpoint = fmt.Sprintf("fc.%s.aliyuncs.com", region)
	}

	fc2Config := &aliopen.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String(fc2Endpoint),
	}
	fc2Client, err := alifc2.NewClient(fc2Config)
	if err != nil {
		return nil, err
	}

	// 接入点一览 https://api.aliyun.com/product/FC-Open
	fc3Endpoint := strings.ReplaceAll(fmt.Sprintf("fcv3.%s.aliyuncs.com", region), "..", ".")
	fc3Config := &aliopen.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String(fc3Endpoint),
	}
	fc3Client, err := alifc3.NewClient(fc3Config)
	if err != nil {
		return nil, err
	}
	return &Clients{FC2: fc2Client, FC3: fc3Client}, nil
}
