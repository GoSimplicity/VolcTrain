package training

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetJobMetricsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取作业指标
func NewGetJobMetricsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetJobMetricsLogic {
	return &GetJobMetricsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetJobMetricsLogic) GetJobMetrics(req *types.GetJobMetricsReq) (resp *types.GetJobMetricsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
