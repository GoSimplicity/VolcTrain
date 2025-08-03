package gpu_node

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveDeviceFromNodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveDeviceFromNodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveDeviceFromNodeLogic {
	return &RemoveDeviceFromNodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveDeviceFromNodeLogic) RemoveDeviceFromNode(req *types.RemoveDeviceFromNodeReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
