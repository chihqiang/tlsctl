package common

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chihqiang/logx"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tcssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
)

func WaitForDeploy(ctx context.Context, client *tcssl.Client, recordId uint64) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		describeHostDeployRecordDetailReq := tcssl.NewDescribeHostDeployRecordDetailRequest()
		describeHostDeployRecordDetailReq.DeployRecordId = common.StringPtr(fmt.Sprintf("%d", recordId))
		resp, err := client.DescribeHostDeployRecordDetail(describeHostDeployRecordDetailReq)
		if err != nil {
			return fmt.Errorf("failed to execute sdk request 'ssl.DescribeHostDeployRecordDetail': %w", err)
		}

		if resp.Response.TotalCount == nil {
			return errors.New("unexpected deployment job status")
		}

		var runningCount, succeededCount, failedCount, totalCount int64
		if resp.Response.RunningTotalCount != nil {
			runningCount = *resp.Response.RunningTotalCount
		}
		if resp.Response.SuccessTotalCount != nil {
			succeededCount = *resp.Response.SuccessTotalCount
		}
		if resp.Response.FailedTotalCount != nil {
			failedCount = *resp.Response.FailedTotalCount
		}
		if resp.Response.TotalCount != nil {
			totalCount = *resp.Response.TotalCount
		}

		if succeededCount+failedCount == totalCount {
			return nil
		}

		logx.Info("waiting for deployment job completion (running: %d, succeeded: %d, failed: %d, total: %d) ...", runningCount, succeededCount, failedCount, totalCount)
		time.Sleep(time.Second * 5)
	}
}
