package training

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTrainingJobLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取训练作业详情
func NewGetTrainingJobLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTrainingJobLogic {
	return &GetTrainingJobLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTrainingJobLogic) GetTrainingJob(req *types.GetTrainingJobReq) (resp *types.GetTrainingJobResp, err error) {
	// todo: add your logic here and delete this line

	return
}
