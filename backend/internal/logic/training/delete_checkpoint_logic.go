package training

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCheckpointLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除检查点
func NewDeleteCheckpointLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCheckpointLogic {
	return &DeleteCheckpointLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteCheckpointLogic) DeleteCheckpoint(req *types.DeleteCheckpointReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
