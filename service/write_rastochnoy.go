package service

import (
	"context"
	"log"
	rass "rastochnoy/genproto/rastochnoy"
	writerastochnoy "rastochnoy/write_rastochnoy"
)

type RastochnoyService struct {
	rass.UnimplementedRastochnoyServer
	Ras *writerastochnoy.RastochnoyRepo
}

func NewRastochnoyService(ras *writerastochnoy.RastochnoyRepo) *RastochnoyService {
	return &RastochnoyService{
		Ras: ras,
	}
}

func (r *RastochnoyService) WriteRastochnoy(ctx context.Context, req *rass.WriteRastochnoyReq) (*rass.WriteRastochnoyRes, error) {
	resp, err := r.Ras.WriteRastochnoy(ctx,req)
	if err != nil {
		log.Println("Yozish service da xatolik?????", err)
		return nil, err
	}
	return resp,nil
}

func (r *RastochnoyService) ReadWriteRastochnoy(ctx context.Context, req *rass.ReadWriteRastochnoyReq) (*rass.ReadWriteRastochnoyRes, error) {
	resp, err := r.Ras.ReadWriteRastochnoy(ctx, req)
	if err != nil {
		log.Println("Yozish serviceni o'qishida xatolik?????", err)
		return nil, err
	}
	return resp,nil
}

func (r *RastochnoyService) ReadRastochnoy(ctx context.Context, req *rass.ReadRastochnoyReq) (*rass.ReadRastochnoyRes, error) {
	resp, err := r.Ras.ReadRastochnoy(ctx, req)
	if err != nil {
		log.Println("O'qish service da xatolik????", err)
		return nil, err
	}
	return resp, nil
}

func (r *RastochnoyService) WrtieRastochnoyDb37(ctx context.Context, req *rass.WrtieRastochnoyDb37Req) (*rass.WrtieRastochnoyDb37Res, error) {
	resp, err := r.Ras.WrtieRastochnoyDb37(ctx, req)
	if err != nil {
		log.Println("Yozish service da xatolik????", err)
		return nil, err
	}
	return resp, nil
}

func (r *RastochnoyService) WriteRastochnoyDB33(ctx context.Context, req *rass.WriteRastochnoydb33Req) (*rass.WriteRastochnoydb33Res, error) {
	resp, err := r.Ras.WriteRastochnoyDB33(ctx,req)
	if err != nil {
		log.Println("Yozish service da xatolik?????", err)
		return nil, err
	}
	return resp,nil
}

func (r *RastochnoyService) ReadWriteRastochnoyDB33(ctx context.Context, req *rass.ReadWriteRastochnoyDB33Req) (*rass.ReadWriteRastochnoyDB33Res, error) {
	resp, err := r.Ras.ReadWriteRastochnoyDB33(ctx, req)
	if err != nil {
		log.Println("Yozish serviceni o'qishida xatolik?????", err)
		return nil, err
	}
	return resp,nil
}