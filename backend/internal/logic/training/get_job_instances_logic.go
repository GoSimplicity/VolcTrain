package training

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetJobInstancesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取作业实例列表
func NewGetJobInstancesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetJobInstancesLogic {
	return &GetJobInstancesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetJobInstancesLogic) GetJobInstances(req *types.GetJobInstancesReq) (resp *types.GetJobInstancesResp, err error) {
	// todo: add your logic here and delete this line

	return
}
