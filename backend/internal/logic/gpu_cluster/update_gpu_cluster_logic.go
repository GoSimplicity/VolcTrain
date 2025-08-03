package gpu_cluster

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateGpuClusterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateGpuClusterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateGpuClusterLogic {
	return &UpdateGpuClusterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateGpuClusterLogic) UpdateGpuCluster(req *types.UpdateGpuClusterReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
