package training

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCheckpointLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建检查点
func NewCreateCheckpointLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCheckpointLogic {
	return &CreateCheckpointLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateCheckpointLogic) CreateCheckpoint(req *types.CreateCheckpointReq) (resp *types.CreateCheckpointResp, err error) {
	// todo: add your logic here and delete this line

	return
}
