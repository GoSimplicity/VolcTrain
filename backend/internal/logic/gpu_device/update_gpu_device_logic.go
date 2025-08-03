package gpu_device

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateGpuDeviceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateGpuDeviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateGpuDeviceLogic {
	return &UpdateGpuDeviceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateGpuDeviceLogic) UpdateGpuDevice(req *types.UpdateGpuDeviceReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
