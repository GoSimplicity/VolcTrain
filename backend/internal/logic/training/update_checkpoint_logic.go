package training

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCheckpointLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新检查点
func NewUpdateCheckpointLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCheckpointLogic {
	return &UpdateCheckpointLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCheckpointLogic) UpdateCheckpoint(req *types.UpdateCheckpointReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
