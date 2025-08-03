package training

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelTrainingJobLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 取消训练作业
func NewCancelTrainingJobLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelTrainingJobLogic {
	return &CancelTrainingJobLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelTrainingJobLogic) CancelTrainingJob(req *types.CancelTrainingJobReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
