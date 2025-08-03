package training

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResumeTrainingJobLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 恢复训练作业
func NewResumeTrainingJobLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResumeTrainingJobLogic {
	return &ResumeTrainingJobLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResumeTrainingJobLogic) ResumeTrainingJob(req *types.ResumeTrainingJobReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
