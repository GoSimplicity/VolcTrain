package gpu_node

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateGpuNodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateGpuNodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateGpuNodeLogic {
	return &UpdateGpuNodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateGpuNodeLogic) UpdateGpuNode(req *types.UpdateGpuNodeReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
