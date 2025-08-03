package gpu_usage

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListGpuUsageRelationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListGpuUsageRelationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListGpuUsageRelationsLogic {
	return &ListGpuUsageRelationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListGpuUsageRelationsLogic) ListGpuUsageRelations(req *types.ListGpuUsageRelationsReq) (resp *types.ListGpuUsageRelationsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
