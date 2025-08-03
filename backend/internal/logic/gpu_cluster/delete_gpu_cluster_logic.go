package gpu_cluster

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteGpuClusterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteGpuClusterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteGpuClusterLogic {
	return &DeleteGpuClusterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteGpuClusterLogic) DeleteGpuCluster(req *types.DeleteGpuClusterReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
