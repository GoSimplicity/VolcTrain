package training

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTrainingJobsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取训练作业列表
func NewListTrainingJobsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTrainingJobsLogic {
	return &ListTrainingJobsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListTrainingJobsLogic) ListTrainingJobs(req *types.ListTrainingJobsReq) (resp *types.ListTrainingJobsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
