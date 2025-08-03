package gpu_node

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListGpuNodesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListGpuNodesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListGpuNodesLogic {
	return &ListGpuNodesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListGpuNodesLogic) ListGpuNodes(req *types.ListGpuNodesReq) (resp *types.ListGpuNodesResp, err error) {
	// todo: add your logic here and delete this line

	return
}
