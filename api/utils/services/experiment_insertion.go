package services

import (
	"archive-api/utils"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func insertTableExp(table_exp utils.TableExperiment, tx pgx.Tx) error {
	pl := new(utils.Placeholder)
	pl.Build(0, 10)
	insert_into_table_exp, err := utils.BuildSQLInsert[utils.TableExperiment]("table_exp", table_exp, pl)
	if err != nil {
		log.Default().Println("error :", err)
		return err
	}
	insert_into_table_exp += ` 
		ON CONFLICT 
		DO NOTHING`

	_, err = tx.Exec(context.Background(), insert_into_table_exp, pl.Args...)
	return err
}

func insertTableLabels(labels []string, exp_id string, tx pgx.Tx) error {
	pl := new(utils.Placeholder)
	pl.Build(0, len(labels)*2)
	insert_into_table_labels := `
		INSERT INTO table_labels
			(exp_id,labels) 
		VALUES `
	for i, label := range labels {
		insert_into_table_labels += fmt.Sprintf("(%s,%s)",
			pl.Get(exp_id),
			pl.Get(strings.ToLower(label)),
		)
		if i < len(labels)-1 {
			insert_into_table_labels += ","
		}
	}
	insert_into_table_labels += `
		ON CONFLICT (exp_id,labels) 
		DO NOTHING`
	_, err := tx.Exec(context.Background(), insert_into_table_labels, pl.Args...)
	return err
}

func insertVariables(nimbus_execution utils.NimbusExecution, variables []utils.Variable, tx pgx.Tx) error {
	pl := new(utils.Placeholder)
	pl.Build(0, 144)

	insert_into_table_nimbus, err := utils.BuildSQLInsert[utils.NimbusExecution]("table_nimbus_execution", nimbus_execution, pl)
	insert_into_table_variable, err_sql := utils.BuildSQLInsertAll[utils.Variable]("table_variable", variables, pl)
	if err_sql != nil {
		return err_sql
	}
	sql := fmt.Sprintf(`
		WITH 
			nimbus_id AS 
				(%s ON CONFLICT (exp_id, config_name, extension, lossless, nan_value_encoding, chunks, rx, ry)
					DO UPDATE SET created_at = excluded.created_at 
						WHERE table_nimbus_execution.created_at < excluded.created_at 
				RETURNING id),

			var_ids_name AS 
				(%s RETURNING name,id)

		INSERT INTO join_nimbus_execution_variables

			SELECT 
				nimbus_id.id AS id_nimbus_execution,
				var_ids_name.name AS variable_name,
				var_ids_name.id AS variable_id

			FROM var_ids_name CROSS JOIN nimbus_id

			ON CONFLICT (id_nimbus_execution,variable_name)
			DO UPDATE SET variable_id = excluded.variable_id;`,
		insert_into_table_nimbus, insert_into_table_variable)

	_, err = tx.Exec(context.Background(), sql, pl.Args...)
	return err
}

func InsertAll(exp_id string, request *utils.Request, pool *pgxpool.Pool) error {
	if err := pgx.BeginFunc(context.Background(), pool,
		func(tx pgx.Tx) error {
			err := insertVariables(request.Request.Table_nimbus_execution, request.Request.Table_variable, tx)
			if err != nil {
				log.Default().Println("error : ", err)
				return err
			}

			err = insertTableExp(request.Request.Table_experiment, tx)
			if err != nil {
				log.Default().Println("error : ", err)
				return err
			}

			err = insertTableLabels(request.Request.Table_experiment.Labels, request.Request.Table_experiment.Exp_id, tx)
			if err != nil {
				log.Default().Println("error : ", err)
				return err
			}
			return nil
		},
	); err != nil {
		log.Default().Println("transactions error :", err)
		return err
	}
	return nil
}

func Clean(pool *pgxpool.Pool) error {
	if err := pgx.BeginFunc(context.Background(), pool,
		func(tx pgx.Tx) error {
			sql := `
			DELETE FROM table_variable 
			WHERE id NOT IN (
				SELECT variable_id 
				FROM join_nimbus_execution_variables
			);`
			_, err := tx.Exec(context.Background(), sql)
			return err
		},
	); err != nil {
		log.Default().Println("transactions error :", err)
		return err
	}
	return nil
}

func AddLabelsForId(id string, labels []string, pool *pgxpool.Pool) error {
	if err := pgx.BeginFunc(context.Background(), pool,
		func(tx pgx.Tx) error {
			err := insertTableLabels(labels, id, tx)
			if err != nil {
				log.Default().Println("error : ", err)
				return err
			}
			return nil
		},
	); err != nil {
		log.Default().Println("transactions error :", err)
		return err
	}
	log.Default().Println("insert success", id)
	return nil
}
