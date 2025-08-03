package training

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetJobOptionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取作业选项
func NewGetJobOptionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetJobOptionsLogic {
	return &GetJobOptionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetJobOptionsLogic) GetJobOptions(req *types.EmptyReq) (resp *types.GetJobOptionsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
