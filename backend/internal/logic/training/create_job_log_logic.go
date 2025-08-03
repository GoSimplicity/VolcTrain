package training

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateJobLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建作业日志
func NewCreateJobLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateJobLogLogic {
	return &CreateJobLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateJobLogLogic) CreateJobLog(req *types.CreateJobLogReq) (resp *types.CreateJobLogResp, err error) {
	// todo: add your logic here and delete this line

	return
}
