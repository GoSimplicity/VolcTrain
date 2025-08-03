package gpu_device

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListGpuDeviceAllocationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListGpuDeviceAllocationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListGpuDeviceAllocationsLogic {
	return &ListGpuDeviceAllocationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListGpuDeviceAllocationsLogic) ListGpuDeviceAllocations(req *types.ListGpuDeviceAllocationsReq) (resp *types.ListGpuDeviceAllocationsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
