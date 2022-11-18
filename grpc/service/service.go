package service

import (
	"bufio"
	"bytes"
	"context"
	"image"
	"image/png"
	"log"
	"os"

	"github.com/rjahon/img-service/config"
	"github.com/rjahon/img-service/genproto/img_service"
	"github.com/rjahon/img-service/grpc/client"
	"github.com/rjahon/img-service/storage"
	"github.com/rjahon/img-service/storage/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var count chan struct{}

type imgService struct {
	cfg      config.Config
	strg     storage.StorageI
	services client.ServiceManagerI
	img_service.UnimplementedServiceServer
}

func NewImgService(cfg config.Config, strg storage.StorageI, svcs client.ServiceManagerI) *imgService {
	count = make(chan struct{}, 10)
	return &imgService{
		cfg:      cfg,
		strg:     strg,
		services: svcs,
	}
}

func (b *imgService) Create(ctx context.Context, req *img_service.CreateRequest) (resp *img_service.CreateResponse, err error) {
	log.Println("#CreateImg ", req.Title)

	err = Btoi(req.Body, "./out/"+req.Title)
	if err != nil {
		log.Println("!CreateImg ", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	img := model.CreateImgRequest{
		Title: req.Title,
	}

	id, err := b.strg.Img().Create(ctx, &img)
	if err != nil {
		log.Println("!CreateImg ", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	res, err := b.strg.Img().GetByPK(ctx, id)
	if err != nil {
		log.Println("!CreateImg ", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return res, err
}

func (b *imgService) Get(ctx context.Context, req *img_service.ImgPrimaryKey) (res *img_service.Img, err error) {
	count <- struct{}{}

	log.Println("#GetImg ", req.Id)
	r, err := b.strg.Img().GetByPK(ctx, &req.Id)
	if err != nil {
		log.Println("!GetImg: ", err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	imgData, err := Itob("./out/" + r.Title)
	if err != nil {
		log.Println("!GetImg: ", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	res = &img_service.Img{
		Id:        r.Id,
		Title:     r.Title,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		Body:      imgData,
	}

	<-count
	return res, err
}

func (b *imgService) GetList(ctx context.Context, req *img_service.GetListRequest) (resp *img_service.GetListResponse, err error) {
	log.Println("#GetImgList ")

	resp, err = b.strg.Img().GetList(ctx, req)

	if err != nil {
		log.Println("!GetImgList: ", err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return resp, err
}

func Itob(location string) (imgData []byte, err error) {
	file, err := os.Open(location)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	var size int64 = fileInfo.Size()

	imgData = make([]byte, size)

	buffer := bufio.NewReader(file)
	_, err = buffer.Read(imgData)
	if err != nil {
		return nil, err
	}

	return imgData, nil
}

func Btoi(data []byte, location string) (err error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return err
	}

	file, err := os.Create(location)
	if err != nil {
		return err
	}

	if err := png.Encode(file, img); err != nil {
		return err
	}

	return nil
}
