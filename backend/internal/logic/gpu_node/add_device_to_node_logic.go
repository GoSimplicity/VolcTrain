package gpu_node

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddDeviceToNodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddDeviceToNodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddDeviceToNodeLogic {
	return &AddDeviceToNodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddDeviceToNodeLogic) AddDeviceToNode(req *types.AddDeviceToNodeReq) (resp *types.AddDeviceToNodeResp, err error) {
	// todo: add your logic here and delete this line

	return
}
