package training

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SuspendTrainingJobLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 暂停训练作业
func NewSuspendTrainingJobLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SuspendTrainingJobLogic {
	return &SuspendTrainingJobLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SuspendTrainingJobLogic) SuspendTrainingJob(req *types.SuspendTrainingJobReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
