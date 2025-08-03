package gpu_device

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListGpuDevicesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListGpuDevicesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListGpuDevicesLogic {
	return &ListGpuDevicesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListGpuDevicesLogic) ListGpuDevices(req *types.ListGpuDevicesReq) (resp *types.ListGpuDevicesResp, err error) {
	// todo: add your logic here and delete this line

	return
}
