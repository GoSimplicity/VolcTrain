package training

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetJobLogsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取作业日志
func NewGetJobLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetJobLogsLogic {
	return &GetJobLogsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetJobLogsLogic) GetJobLogs(req *types.GetJobLogsReq) (resp *types.GetJobLogsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
