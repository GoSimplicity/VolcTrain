package gpu_device

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReleaseGpuDeviceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReleaseGpuDeviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReleaseGpuDeviceLogic {
	return &ReleaseGpuDeviceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReleaseGpuDeviceLogic) ReleaseGpuDevice(req *types.ReleaseGpuDeviceReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
