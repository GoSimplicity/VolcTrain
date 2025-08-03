package training

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteTrainingQueueLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除训练队列
func NewDeleteTrainingQueueLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTrainingQueueLogic {
	return &DeleteTrainingQueueLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteTrainingQueueLogic) DeleteTrainingQueue(req *types.DeleteTrainingQueueReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
