package handler

import (
	"batch_tx/internal/logic"
	"batch_tx/internal/svc"
	"context"

	"github.com/zeromicro/go-zero/core/service"
)

func RegisterJob(serverCtx *svc.ServiceContext, group *service.ServiceGroup) {

	group.Add(logic.NewProducerLogic(context.Background(), serverCtx))
	// group.Add(logic.NewConsumerLogic(context.Background(), serverCtx))

	group.Start()

}
