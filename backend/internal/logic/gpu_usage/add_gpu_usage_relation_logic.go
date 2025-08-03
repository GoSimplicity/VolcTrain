package gpu_usage

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddGpuUsageRelationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddGpuUsageRelationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddGpuUsageRelationLogic {
	return &AddGpuUsageRelationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddGpuUsageRelationLogic) AddGpuUsageRelation(req *types.AddGpuUsageRelationReq) (resp *types.AddGpuUsageRelationResp, err error) {
	// todo: add your logic here and delete this line

	return
}
