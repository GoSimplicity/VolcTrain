package gpu_device

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGpuDeviceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGpuDeviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGpuDeviceLogic {
	return &GetGpuDeviceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGpuDeviceLogic) GetGpuDevice(req *types.GetGpuDeviceReq) (resp *types.GetGpuDeviceResp, err error) {
	// todo: add your logic here and delete this line

	return
}
