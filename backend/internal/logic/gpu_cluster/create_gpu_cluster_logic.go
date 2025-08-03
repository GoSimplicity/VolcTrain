package gpu_cluster

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGpuClusterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateGpuClusterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGpuClusterLogic {
	return &CreateGpuClusterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateGpuClusterLogic) CreateGpuCluster(req *types.CreateGpuClusterReq) (resp *types.CreateGpuClusterResp, err error) {
	// todo: add your logic here and delete this line

	return
}
