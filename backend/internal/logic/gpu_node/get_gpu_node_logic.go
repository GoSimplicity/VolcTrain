package gpu_node

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGpuNodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGpuNodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGpuNodeLogic {
	return &GetGpuNodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGpuNodeLogic) GetGpuNode(req *types.GetGpuNodeReq) (resp *types.GetGpuNodeResp, err error) {
	// todo: add your logic here and delete this line

	return
}
