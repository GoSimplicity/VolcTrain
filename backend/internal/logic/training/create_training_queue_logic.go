package training

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTrainingQueueLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建训练队列
func NewCreateTrainingQueueLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTrainingQueueLogic {
	return &CreateTrainingQueueLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateTrainingQueueLogic) CreateTrainingQueue(req *types.CreateTrainingQueueReq) (resp *types.CreateTrainingQueueResp, err error) {
	// todo: add your logic here and delete this line

	return
}
