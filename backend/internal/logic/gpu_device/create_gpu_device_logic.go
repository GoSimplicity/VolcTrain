package gpu_device

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGpuDeviceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateGpuDeviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGpuDeviceLogic {
	return &CreateGpuDeviceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateGpuDeviceLogic) CreateGpuDevice(req *types.CreateGpuDeviceReq) (resp *types.CreateGpuDeviceResp, err error) {
	// todo: add your logic here and delete this line

	return
}
