package training

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetJobRelationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取作业关联关系
func NewGetJobRelationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetJobRelationsLogic {
	return &GetJobRelationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetJobRelationsLogic) GetJobRelations(req *types.GetJobRelationsReq) (resp *types.GetJobRelationsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
