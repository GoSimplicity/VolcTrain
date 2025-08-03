package gpu_usage

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListGpuUsageRecordsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListGpuUsageRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListGpuUsageRecordsLogic {
	return &ListGpuUsageRecordsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListGpuUsageRecordsLogic) ListGpuUsageRecords(req *types.ListGpuUsageRecordsReq) (resp *types.ListGpuUsageRecordsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
