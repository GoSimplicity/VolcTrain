package gpu_usage

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGpuUsageRecordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGpuUsageRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGpuUsageRecordLogic {
	return &GetGpuUsageRecordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGpuUsageRecordLogic) GetGpuUsageRecord(req *types.GetGpuUsageRecordReq) (resp *types.GetGpuUsageRecordResp, err error) {
	// todo: add your logic here and delete this line

	return
}
