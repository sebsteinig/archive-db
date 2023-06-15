package services

import (
	"archive-api/utils"
	"archive-api/utils/sql"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Exp_with_Label struct {
	Exp_id string        `json:"exp_id"`
	Labels []utils.Label `json:"labels"`
}

type Publication struct {
	Title         string `json:"title" sql:"title"`
	Authors_short string `json:"authors_short" sql:"authors_short"`
	Authors_full  string `json:"authors_full" sql:"authors_full"`
	Journal       string `json:"journal" sql:"journal"`
	Year          int    `json:"year" sql:"year"`
	//Volume        string `json:"volume" sql:"volume,ignore"`
	//Pages         int `json:"pages" sql:"pages,ignore"`
	//Doi           string `json:"doi" sql:"doi,ignore"`
	Owner_name  string `json:"owner_name" sql:"owner_name"`
	Owner_email string `json:"owner_email" sql:"owner_email"`
	Abstract    string `json:"abstract" sql:"abstract"`
	Brief_desc  string `json:"brief_desc" sql:"brief_desc"`
	//Expts_paper []string         `json:"expts_paper" sql:"expts_paper"`
	Expts_web []Exp_with_Label `json:"expts_web"`
}

type JoinPublicationExp struct {
	PublicationId    int           `sql:"publication_id"`
	Requested_exp_id string        `sql:"requested_exp_id,nullable"`
	Exp_id           string        `sql:"exp_id,nullable"`
	Metadata         []utils.Label `sql:"metadata"`
}

func selectRequestedIds(exp_ids []string, pool *pgxpool.Pool) (map[string]struct{}, error) {
	pl := sql.BuildPlaceholder(len(exp_ids))
	values_exp_ids := make([]string, len(exp_ids))
	for i, exp_id := range exp_ids {
		values_exp_ids[i] = fmt.Sprintf("(%s)", pl.Push(exp_id))
	}
	select_requested_ids_sql := fmt.Sprintf(`
		select exp_id as requested_exp_id
		from (
				values %s
			) as exps_publication (exp_id)
		except
		select exp_id
		from table_exp
	`, strings.Join(values_exp_ids, ","))
	rows, err := pool.Query(context.Background(), select_requested_ids_sql, pl.Args...)
	if err != nil {
		log.Default().Println("Unable to query:", select_requested_ids_sql, "error :", err)
		return nil, err
	}
	defer rows.Close()
	var map_exp map[string]struct{} = make(map[string]struct{})
	var res string
	_, err = pgx.ForEachRow(rows, []any{&res}, func() error {
		map_exp[res] = struct{}{}
		return nil
	})
	return map_exp, err
}

type Id struct {
	Id int `sql:"id"`
}

func insertPublication(publications []Publication, ids *[]int, tx pgx.Tx) error {

	query, err := sql.Insert[Publication]("table_publication", publications...)
	if err != nil {
		log.Default().Println("ERROR <insertPublication>")
		return err
	}
	query.Suffixe(` 
		ON CONFLICT (title, journal, year, owner_name) DO NOTHING
		RETURNING id
	`)
	type Id struct {
		Id int `sql:"id"`
	}
	responses, err := sql.Receive[Id](context.Background(), &query, tx)
	if err != nil {
		log.Default().Println("ERROR <SearchExperimentLike>")
		return err
	}
	for _, id := range responses {
		*ids = append(*ids, id.Id)
	}
	return err
}

func PublicationInsert(c *fiber.Ctx, exp_ids []string, publications []Publication, pool *pgxpool.Pool) error {
	requested_ids, err := selectRequestedIds(exp_ids, pool)
	if err != nil {
		log.Default().Println("error : ", err)
		return err
	}

	if err := pgx.BeginFunc(context.Background(), pool,
		func(tx pgx.Tx) error {
			ids := make([]int, 0, len(publications))
			err = insertPublication(publications, &ids, tx)
			if err != nil {
				log.Default().Println("error : ", err)
				return err
			}
			if len(ids) == 0 {
				return nil
			}
			var joins []JoinPublicationExp
			var exp_label_joins []utils.JoinExpLabel
			for i, id := range ids {
				for _, exp_id := range publications[i].Expts_web {
					if _, ok := requested_ids[exp_id.Exp_id]; ok {
						joins = append(joins,
							JoinPublicationExp{
								PublicationId:    id,
								Requested_exp_id: exp_id.Exp_id,
								Metadata:         exp_id.Labels,
							})
					} else {
						joins = append(joins,
							JoinPublicationExp{
								PublicationId: id,
								Exp_id:        exp_id.Exp_id,
								Metadata:      exp_id.Labels,
							})
						for _, label := range exp_id.Labels {
							exp_label_joins = append(exp_label_joins, utils.JoinExpLabel{
								Exp_id:   exp_id.Exp_id,
								Label:    label.Label,
								Metadata: label.Metadata,
							})
						}
					}
				}
			}

			query, err := sql.Insert[JoinPublicationExp]("join_publication_exp", joins...)
			query.Suffixe(`
				ON CONFLICT (exp_id, publication_id)
				DO NOTHING
			`)
			if err != nil {
				log.Default().Println("ERROR <PublicationInsert>")
				return err
			}
			err = sql.Exec(context.Background(), &query, tx)
			if err != nil {
				log.Default().Println("ERROR <PublicationInsert>")
				return err
			}

			if len(exp_label_joins) > 0 {
				query, err := sql.Insert[utils.JoinExpLabel]("table_labels", exp_label_joins...)
				query.Suffixe(`
					ON CONFLICT (exp_id,labels) 
					DO NOTHING
				`)
				if err != nil {
					log.Default().Println("ERROR <PublicationInsert>")
					return err
				}
				err = sql.Exec(context.Background(), &query, tx)
				if err != nil {
					log.Default().Println("ERROR <PublicationInsert>")
					return err
				}
			}

			return nil
		},
	); err != nil {
		log.Default().Println("transactions error :", err)
		return err
	}
	if len(requested_ids) != 0 {
		invalid_expids := make([]string, len(requested_ids))
		i := 0
		for exp_id := range requested_ids {
			invalid_expids[i] = exp_id
			i++
		}
		type Requested_exp_id struct {
			Invalid_expids []string `json:"requested_id"`
		}
		return c.Status(fiber.StatusConflict).JSON(Requested_exp_id{Invalid_expids: invalid_expids})
	}
	return nil
}
