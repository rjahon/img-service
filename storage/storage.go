package storage

import (
	"context"

	pb "github.com/rjahon/img-service/genproto/img_service"
	"github.com/rjahon/img-service/storage/model"
)

type StorageI interface {
	CloseDB()
	Img() ImgRepoI
}

type ImgRepoI interface {
	Create(ctx context.Context, req *model.CreateImgRequest) (id *string, err error)
	GetByPK(ctx context.Context, id *string) (res *pb.CreateResponse, err error)
	GetList(ctx context.Context, req *pb.GetListRequest) (res *pb.GetListResponse, err error)
}
