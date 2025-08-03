package training

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteJobRelationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除作业关联关系
func NewDeleteJobRelationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteJobRelationLogic {
	return &DeleteJobRelationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteJobRelationLogic) DeleteJobRelation(req *types.DeleteJobRelationReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
