package training

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTrainingQueueLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新训练队列
func NewUpdateTrainingQueueLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTrainingQueueLogic {
	return &UpdateTrainingQueueLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateTrainingQueueLogic) UpdateTrainingQueue(req *types.UpdateTrainingQueueReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
