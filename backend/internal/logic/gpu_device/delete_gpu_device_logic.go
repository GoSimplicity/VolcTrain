package gpu_device

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteGpuDeviceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteGpuDeviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteGpuDeviceLogic {
	return &DeleteGpuDeviceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteGpuDeviceLogic) DeleteGpuDevice(req *types.DeleteGpuDeviceReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
