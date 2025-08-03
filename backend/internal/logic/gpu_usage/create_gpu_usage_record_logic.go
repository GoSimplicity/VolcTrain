package gpu_usage

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGpuUsageRecordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateGpuUsageRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGpuUsageRecordLogic {
	return &CreateGpuUsageRecordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateGpuUsageRecordLogic) CreateGpuUsageRecord(req *types.CreateGpuUsageRecordReq) (resp *types.CreateGpuUsageRecordResp, err error) {
	// todo: add your logic here and delete this line

	return
}
