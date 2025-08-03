package gpu_node

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteGpuNodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteGpuNodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteGpuNodeLogic {
	return &DeleteGpuNodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteGpuNodeLogic) DeleteGpuNode(req *types.DeleteGpuNodeReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
