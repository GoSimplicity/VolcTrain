package gpu_cluster

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListClusterNodesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListClusterNodesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListClusterNodesLogic {
	return &ListClusterNodesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListClusterNodesLogic) ListClusterNodes(req *types.ListClusterNodesReq) (resp *types.ListClusterNodesResp, err error) {
	// todo: add your logic here and delete this line

	return
}
