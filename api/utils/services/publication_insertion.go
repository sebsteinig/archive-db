package services

import (
	"archive-api/utils"
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Publication struct {
	Title         string `json:"title" sql:"title"`
	Authors_short string `json:"authors_short" sql:"authors_short"`
	Authors_full  string `json:"authors_full" sql:"authors_full"`
	Journal       string `json:"journal" sql:"journal"`
	Year          int    `json:"year" sql:"year"`
	//Volume        string `json:"volume" sql:"volume"`
	//Pages         int `json:"pages" sql:"pages"`
	//Doi           string `json:"doi" sql:"doi"`
	Owner_name  string   `json:"owner_name" sql:"owner_name"`
	Owner_email string   `json:"owner_email" sql:"owner_email"`
	Abstract    string   `json:"abstract" sql:"abstract"`
	Brief_desc  string   `json:"brief_desc" sql:"brief_desc"`
	Expts_paper []string `json:"expts_paper" sql:"expts_paper"`
	Expts_web   []string `json:"expts_web"`
}

type JoinPublicationExp struct {
	PublicationId int    `sql:"publication_id"`
	Exp_id        string `sql:"exp_id"`
}

func PublicationInsert(c *fiber.Ctx, publications []Publication, pool *pgxpool.Pool) error {
	pl := new(utils.Placeholder)
	pl.Build(0, len(publications)*14)
	insert_sql, err := utils.BuildSQLInsertAll[Publication]("table_publication", publications, pl)
	if err != nil {
		log.Default().Println("error : ", err)
		return err
	}
	insert_sql += ` RETURNING ID`
	rows, err_exec := pool.Query(context.Background(), insert_sql, pl.Args...)
	if err_exec != nil {
		log.Default().Println("Unable to query:", insert_sql, "error :", err)
		return err
	}
	type Id struct {
		Id int `sql:"id"`
	}
	defer rows.Close()
	ids, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (Id, error) {
		var res Id
		err := utils.BuildSQLResponse(row, &res)
		return res, err
	})
	if len(ids) != len(publications) {
		return fmt.Errorf("retrieved ids and submitted publication have different size")
	}
	pl = new(utils.Placeholder)
	pl.Build(0, len(publications)*2)
	var joins []JoinPublicationExp
	for i, id := range ids {
		for _, exp_id := range publications[i].Expts_web {
			joins = append(joins,
				JoinPublicationExp{
					PublicationId: id.Id,
					Exp_id:        exp_id,
				})
		}
	}
	join_sql, err := utils.BuildSQLInsertAll[JoinPublicationExp]("join_publication_exp", joins, pl)
	join_sql += " RETURNING excluded.exp_id"
	rows, err_exec = pool.Query(context.Background(), join_sql, pl.Args...)
	if err_exec != nil {
		log.Default().Println("Unable to query:", insert_sql, "error :", err)
		return err
	}
	type Exp_id struct {
		Exp_id string `sql:"exp_id"`
	}
	defer rows.Close()
	invalid_expids, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (Exp_id, error) {
		var res Exp_id
		err := utils.BuildSQLResponse(row, &res)
		return res, err
	})
	if err != nil {
		log.Default().Println("error :", err)
		return err
	}
	if len(invalid_expids) != 0 {
		return c.Status(fiber.StatusConflict).JSON(invalid_expids)
	}
	return nil
}
