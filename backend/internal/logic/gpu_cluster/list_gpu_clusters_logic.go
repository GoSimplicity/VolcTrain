package gpu_cluster

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListGpuClustersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListGpuClustersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListGpuClustersLogic {
	return &ListGpuClustersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListGpuClustersLogic) ListGpuClusters(req *types.ListGpuClustersReq) (resp *types.ListGpuClustersResp, err error) {
	// todo: add your logic here and delete this line

	return
}
