package services

import (
	"archive-api/utils"
	"context"
	"log"

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
	//Pages         string `json:"pages" sql:"pages"`
	//Doi           string `json:"doi" sql:"doi"`
	Owner_name  string   `json:"owner_name" sql:"owner_name"`
	Owner_email string   `json:"owner_email" sql:"owner_email"`
	Abstract    string   `json:"abstract" sql:"abstract"`
	Brief_desc  string   `json:"brief_desc" sql:"brief_desc"`
	Expts_paper string   `json:"expts_paper" sql:"expts_paper"`
	Expts_web   []string `json:"expts_web"`
}

func PublicationInsertMultiple(publications []Publication, pool *pgxpool.Pool) error {
	if err := pgx.BeginFunc(context.Background(), pool,
		func(tx pgx.Tx) error {
			pl := new(utils.Placeholder)
			pl.Build(0, len(publications)*14)
			sql, err := utils.BuildSQLInsertAll[Publication]("table_publication", publications, pl)
			if err != nil {
				log.Default().Println("error : ", err)
				return err
			}
			_, err_exec := tx.Exec(context.Background(), sql, pl.Args...)
			return err_exec
		},
	); err != nil {
		log.Default().Println("transactions error :", err)
		return err
	}
	return nil
}

func PublicationInsertSingle(publication Publication, pool *pgxpool.Pool) error {
	if err := pgx.BeginFunc(context.Background(), pool,
		func(tx pgx.Tx) error {
			pl := new(utils.Placeholder)
			pl.Build(0, 14)
			sql, err := utils.BuildSQLInsert[Publication]("table_publication", publication, pl)
			if err != nil {
				log.Default().Println("error : ", err)
				return err
			}
			_, err_exec := tx.Exec(context.Background(), sql, pl.Args...)
			return err_exec
		},
	); err != nil {
		log.Default().Println("transactions error :", err)
		return err
	}
	return nil
}
