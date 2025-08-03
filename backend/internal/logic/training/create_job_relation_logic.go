package training

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateJobRelationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建作业关联关系
func NewCreateJobRelationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateJobRelationLogic {
	return &CreateJobRelationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateJobRelationLogic) CreateJobRelation(req *types.CreateJobRelationReq) (resp *types.CreateJobRelationResp, err error) {
	// todo: add your logic here and delete this line

	return
}
