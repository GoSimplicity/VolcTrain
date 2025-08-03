package training

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateJobMetricLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建作业指标
func NewCreateJobMetricLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateJobMetricLogic {
	return &CreateJobMetricLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateJobMetricLogic) CreateJobMetric(req *types.CreateJobMetricReq) (resp *types.CreateJobMetricResp, err error) {
	// todo: add your logic here and delete this line

	return
}
