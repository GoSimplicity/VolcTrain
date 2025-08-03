package gpu_device

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AllocateGpuDeviceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAllocateGpuDeviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllocateGpuDeviceLogic {
	return &AllocateGpuDeviceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AllocateGpuDeviceLogic) AllocateGpuDevice(req *types.AllocateGpuDeviceReq) (resp *types.AllocateGpuDeviceResp, err error) {
	// todo: add your logic here and delete this line

	return
}
