package training

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTrainingQueuesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取训练队列列表
func NewListTrainingQueuesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTrainingQueuesLogic {
	return &ListTrainingQueuesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListTrainingQueuesLogic) ListTrainingQueues(req *types.ListTrainingQueuesReq) (resp *types.ListTrainingQueuesResp, err error) {
	// todo: add your logic here and delete this line

	return
}
