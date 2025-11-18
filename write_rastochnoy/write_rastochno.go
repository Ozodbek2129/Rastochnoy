package writerastochnoy

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	rass "rastochnoy/genproto/rastochnoy"
)

type RastochnoyRepo struct {
	db *sql.DB
}

func NewRastochnoyRepo(db *sql.DB) *RastochnoyRepo {
	return &RastochnoyRepo{db: db}
}

func (ras *RastochnoyRepo) WriteRastochnoy(ctx context.Context, req *rass.WriteRastochnoyReq) (*rass.WriteRastochnoyRes, error) {
	query := `update rastochnoy_write set value = $1 where key = $2`

	result, err := ras.db.ExecContext(ctx, query, req.Value, req.Key)
	if err != nil {
		return nil, errors.New("error updating value")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, errors.ErrUnsupported
	}

	if rowsAffected == 0 {
		return nil, errors.New("no key found")
	}

	return &rass.WriteRastochnoyRes{
		Key:     req.Key,
		Message: "information updated",
	}, nil
}

func (ras *RastochnoyRepo) ReadWriteRastochnoy(ctx context.Context, req *rass.ReadWriteRastochnoyReq) (*rass.ReadWriteRastochnoyRes, error) {
	query := `
        SELECT 
            id, 
            key, 
            COALESCE(offsett, 0) AS offsett, 
            value
        FROM rastochnoy_write
        ORDER BY offsett
    `

	rows, err := ras.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*rass.ReadWriteRastoshnoyItem

	for rows.Next() {
		var (
			id      string
			key     string
			offsett float64
			value   bool
		)

		if err := rows.Scan(&id, &key, &offsett, &value); err != nil {
			return nil, fmt.Errorf("query scan error: %w", err)
		}

		item := &rass.ReadWriteRastoshnoyItem{
			Id:      &id,
			Key:     &key,
			Offsett: &offsett,
			Value:   &value, // endi bool tipida
		}

		items = append(items, item)
	}

	return &rass.ReadWriteRastochnoyRes{
		Data: items,
	}, nil
}

func (ras *RastochnoyRepo) ReadRastochnoy(ctx context.Context, req *rass.ReadRastochnoyReq) (*rass.ReadRastochnoyRes, error) {
	query := `SELECT 
        id, 
        key, 
        COALESCE(offsett, 0) AS offsett, 
        value
    FROM rastochnoy_read
    ORDER BY offsett`

	rows, err := ras.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*rass.RastochnoyItem

	for rows.Next() {
		var (
			id      string
			key     string
			offsett float64
			value   bool
		)

		if err := rows.Scan(&id, &key, &offsett, &value); err != nil {
			return nil, err
		}

		item := &rass.RastochnoyItem{
			Id:      &id,
			Key:     &key,
			Offsett: &offsett,
			Value:   &value,
		}
		items = append(items, item)
	}

	return &rass.ReadRastochnoyRes{
		Data: items,
	}, nil
}

func (ras *RastochnoyRepo) WrtieRastochnoyDb37(ctx context.Context, req *rass.WrtieRastochnoyDb37Req) (*rass.WrtieRastochnoyDb37Res, error) {
	query := `update rastochnoy_read set value = $1 where key = $2`
	//ROUND(offsett::numeric, 3) = ROUND($3::numeric, 3)

	result, err := ras.db.ExecContext(ctx, query, req.Value, req.Key)
	if err != nil {
		return nil, errors.New("error updating value")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, errors.ErrUnsupported
	}

	if rowsAffected == 0 {
		return nil, errors.New("no key found")
	}

	return &rass.WrtieRastochnoyDb37Res{
		Message: "ok db37",
	}, nil
}


func (ras *RastochnoyRepo) WriteRastochnoyDB33(ctx context.Context, req *rass.WriteRastochnoydb33Req) (*rass.WriteRastochnoydb33Res, error) {
	query := `update rastochnoy_writedb33 set value = $1 where key = $2`

	result, err := ras.db.ExecContext(ctx, query, req.Value, req.Key)
	if err != nil {
		return nil, errors.New("error updating value")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, errors.ErrUnsupported
	}

	if rowsAffected == 0 {
		return nil, errors.New("no key found")
	}

	return &rass.WriteRastochnoydb33Res{
		Key:     req.Key,
		Message: "information updated",
	}, nil
}

func (ras *RastochnoyRepo) ReadWriteRastochnoyDB33(ctx context.Context, req *rass.ReadWriteRastochnoyDB33Req) (*rass.ReadWriteRastochnoyDB33Res, error) {
	query := `
        SELECT 
            id, 
            key, 
            COALESCE(offsett, 0) AS offsett, 
            value
        FROM rastochnoy_writedb33
        ORDER BY offsett
    `

	rows, err := ras.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*rass.ReadWriteRastoshnoyDB33Item

	for rows.Next() {
		var (
			id      string
			key     string
			offsett float64
			value   float32
		)

		if err := rows.Scan(&id, &key, &offsett, &value); err != nil {
			return nil, fmt.Errorf("query scan error: %w", err)
		}

		item := &rass.ReadWriteRastoshnoyDB33Item{
			Id:      &id,
			Key:     &key,
			Offsett: &offsett,
			Value:   &value, // endi bool tipida
		}

		items = append(items, item)
	}

	return &rass.ReadWriteRastochnoyDB33Res{
		Data: items,
	}, nil
}