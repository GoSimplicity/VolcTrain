package training

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetJobInstanceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取实例详情
func NewGetJobInstanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetJobInstanceLogic {
	return &GetJobInstanceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetJobInstanceLogic) GetJobInstance(req *types.GetJobInstanceReq) (resp *types.GetJobInstanceResp, err error) {
	// todo: add your logic here and delete this line

	return
}
