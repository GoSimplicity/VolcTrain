package training

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTrainingQueueLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取训练队列详情
func NewGetTrainingQueueLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTrainingQueueLogic {
	return &GetTrainingQueueLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTrainingQueueLogic) GetTrainingQueue(req *types.GetTrainingQueueReq) (resp *types.GetTrainingQueueResp, err error) {
	// todo: add your logic here and delete this line

	return
}
