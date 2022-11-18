package postgres

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

	"github.com/google/uuid"
	pb "github.com/rjahon/img-service/genproto/img_service"
	"github.com/rjahon/img-service/storage"
	"github.com/rjahon/img-service/storage/model"

	"github.com/jackc/pgx/v4/pgxpool"
)

type imgRepo struct {
	db *pgxpool.Pool
}

func NewImgRepo(db *pgxpool.Pool) storage.ImgRepoI {
	return &imgRepo{
		db: db,
	}
}

func (r *imgRepo) Create(ctx context.Context, req *model.CreateImgRequest) (id *string, err error) {
	query := `
		INSERT INTO img
		(
			id,
			title
		) VALUES
		(
			$1,
			$2
		)
	`
	uuid, err := uuid.NewUUID()
	if err != nil {
		return
	}

	_, err = r.db.Exec(
		ctx, query,
		uuid,
		req.Title,
	)
	if err != nil {
		return
	}

	idStr := uuid.String()

	return &idStr, err
}

func (r *imgRepo) GetList(ctx context.Context, req *pb.GetListRequest) (res *pb.GetListResponse, err error) {
	res = &pb.GetListResponse{}

	var (
		filter string = ` WHERE TRUE `
		params map[string]interface{}
		limit  string
		offset string
	)

	query := `
		SELECT 
			id,
			title,
			created_at,
			updated_at,
			count(*) OVER()
		FROM
			img
	`

	params = make(map[string]interface{})

	if req.GetLimit() > 0 {
		limit = ` LIMIT :limit`
		params["limit"] = req.Limit
	}

	if req.GetOffset() > 0 {
		offset = ` OFFSET :offset`
		params["offset"] = req.Offset
	}

	query += filter + limit + offset
	query, args := ReplaceQueryParams(query, params)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id         sql.NullString
			title      sql.NullString
			created_at sql.NullString
			updated_at sql.NullString
		)

		err = rows.Scan(
			&id,
			&title,
			&created_at,
			&updated_at,
			&res.Count,
		)
		if err != nil {
			return
		}

		res.Imgs = append(res.Imgs, &pb.Img{
			Id:        id.String,
			Title:     title.String,
			CreatedAt: created_at.String,
			UpdatedAt: updated_at.String,
		})
	}

	return res, err
}

func (r *imgRepo) GetByPK(ctx context.Context, id *string) (res *pb.CreateResponse, err error) {
	res = &pb.CreateResponse{}
	var (
		pk         sql.NullString
		title      sql.NullString
		created_at sql.NullString
		updated_at sql.NullString
	)

	query := `
		SELECT
			id,
			title,
			created_at,
			updated_at
		FROM
			img
		WHERE
			id = $1
	`

	err = r.db.QueryRow(ctx, query, id).Scan(
		&pk,
		&title,
		&created_at,
		&updated_at,
	)
	if err != nil {
		return
	}

	res = &pb.CreateResponse{
		Id:        pk.String,
		Title:     title.String,
		CreatedAt: created_at.String,
		UpdatedAt: updated_at.String,
	}

	return res, err
}

func ReplaceQueryParams(namedQuery string, params map[string]interface{}) (string, []interface{}) {
	var (
		i    int = 1
		args []interface{}
	)

	for k, v := range params {
		if k != "" {
			oldsize := len(namedQuery)
			namedQuery = strings.ReplaceAll(namedQuery, ":"+k, "$"+strconv.Itoa(i))

			if oldsize != len(namedQuery) {
				args = append(args, v)
				i++
			}
		}
	}

	return namedQuery, args
}
