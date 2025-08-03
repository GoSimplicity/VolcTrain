package training

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetJobInstanceLogsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取实例日志
func NewGetJobInstanceLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetJobInstanceLogsLogic {
	return &GetJobInstanceLogsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetJobInstanceLogsLogic) GetJobInstanceLogs(req *types.GetJobInstanceLogsReq) (resp *types.GetJobInstanceLogsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
