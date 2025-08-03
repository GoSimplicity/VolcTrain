package training

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCheckpointLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取检查点详情
func NewGetCheckpointLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCheckpointLogic {
	return &GetCheckpointLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCheckpointLogic) GetCheckpoint(req *types.GetCheckpointReq) (resp *types.GetCheckpointResp, err error) {
	// todo: add your logic here and delete this line

	return
}
