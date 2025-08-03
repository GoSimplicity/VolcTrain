package gpu_node

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGpuNodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateGpuNodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGpuNodeLogic {
	return &CreateGpuNodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateGpuNodeLogic) CreateGpuNode(req *types.CreateGpuNodeReq) (resp *types.CreateGpuNodeResp, err error) {
	// todo: add your logic here and delete this line

	return
}
