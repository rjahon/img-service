package model

import (
	"context"

	"github.com/rjahon/img-service/genproto/img_service"
)

type CreateImgRequest struct {
	Title string `json:"title"`
}

type Img struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type WorkerReq struct {
	ID  string
	Ctx context.Context
}

type WorkerRes struct {
	Img img_service.Img
	Err error
}
