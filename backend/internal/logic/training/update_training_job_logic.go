package training

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTrainingJobLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新训练作业
func NewUpdateTrainingJobLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTrainingJobLogic {
	return &UpdateTrainingJobLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateTrainingJobLogic) UpdateTrainingJob(req *types.UpdateTrainingJobReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
