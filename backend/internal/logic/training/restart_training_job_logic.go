package training

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RestartTrainingJobLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 重启训练作业
func NewRestartTrainingJobLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RestartTrainingJobLogic {
	return &RestartTrainingJobLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RestartTrainingJobLogic) RestartTrainingJob(req *types.RestartTrainingJobReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
