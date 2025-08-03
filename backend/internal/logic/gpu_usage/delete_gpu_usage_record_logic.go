package gpu_usage

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteGpuUsageRecordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteGpuUsageRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteGpuUsageRecordLogic {
	return &DeleteGpuUsageRecordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteGpuUsageRecordLogic) DeleteGpuUsageRecord(req *types.DeleteGpuUsageRecordReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
