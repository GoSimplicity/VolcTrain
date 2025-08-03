package gpu_cluster

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveNodeFromClusterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveNodeFromClusterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveNodeFromClusterLogic {
	return &RemoveNodeFromClusterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveNodeFromClusterLogic) RemoveNodeFromCluster(req *types.RemoveNodeFromClusterReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
