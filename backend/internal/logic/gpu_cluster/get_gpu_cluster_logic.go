package gpu_cluster

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGpuClusterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGpuClusterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGpuClusterLogic {
	return &GetGpuClusterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGpuClusterLogic) GetGpuCluster(req *types.GetGpuClusterReq) (resp *types.GetGpuClusterResp, err error) {
	// todo: add your logic here and delete this line

	return
}
