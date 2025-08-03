package training

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetJobCheckpointsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取作业检查点列表
func NewGetJobCheckpointsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetJobCheckpointsLogic {
	return &GetJobCheckpointsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetJobCheckpointsLogic) GetJobCheckpoints(req *types.GetJobCheckpointsReq) (resp *types.GetJobCheckpointsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
