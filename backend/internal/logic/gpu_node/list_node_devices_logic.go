package gpu_node

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListNodeDevicesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListNodeDevicesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListNodeDevicesLogic {
	return &ListNodeDevicesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListNodeDevicesLogic) ListNodeDevices(req *types.ListNodeDevicesReq) (resp *types.ListNodeDevicesResp, err error) {
	// todo: add your logic here and delete this line

	return
}
