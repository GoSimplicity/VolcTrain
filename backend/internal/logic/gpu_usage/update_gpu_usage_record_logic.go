package gpu_usage

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateGpuUsageRecordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateGpuUsageRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateGpuUsageRecordLogic {
	return &UpdateGpuUsageRecordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateGpuUsageRecordLogic) UpdateGpuUsageRecord(req *types.UpdateGpuUsageRecordReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
