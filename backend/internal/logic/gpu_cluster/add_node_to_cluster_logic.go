package gpu_cluster

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddNodeToClusterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddNodeToClusterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddNodeToClusterLogic {
	return &AddNodeToClusterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddNodeToClusterLogic) AddNodeToCluster(req *types.AddNodeToClusterReq) (resp *types.AddNodeToClusterResp, err error) {
	// todo: add your logic here and delete this line

	return
}
