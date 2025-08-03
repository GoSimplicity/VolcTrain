package training

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetQueueOptionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取队列选项
func NewGetQueueOptionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetQueueOptionsLogic {
	return &GetQueueOptionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetQueueOptionsLogic) GetQueueOptions(req *types.EmptyReq) (resp *types.GetQueueOptionsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
