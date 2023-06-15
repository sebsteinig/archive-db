package services

import (
	"archive-api/utils"
	"archive-api/utils/sql"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func insertTableExp(table_exp utils.TableExperiment, tx pgx.Tx) error {
	query, err := sql.Insert[utils.TableExperiment]("table_exp", table_exp)
	query.Suffixe(` ON CONFLICT DO NOTHING`)
	if err != nil {
		log.Default().Println("ERROR <insertTableExp>")
		return err
	}
	err = sql.Exec(context.Background(), &query, tx)
	if err != nil {
		log.Default().Println("ERROR <insertTableExp>")
		return err
	}
	return err
}

func insertTableLabels(labels []utils.Label, publication_labels []utils.Label, exp_id string, tx pgx.Tx) error {
	if len(labels) == 0 && len(publication_labels) == 0 {
		return nil
	}

	var exp_label_joins []utils.JoinExpLabel
	for _, label := range labels {
		exp_label_joins = append(exp_label_joins, utils.JoinExpLabel{
			Exp_id:   exp_id,
			Label:    label.Label,
			Metadata: label.Metadata,
		})
	}
	for _, label := range publication_labels {
		exp_label_joins = append(exp_label_joins, utils.JoinExpLabel{
			Exp_id:   exp_id,
			Label:    label.Label,
			Metadata: label.Metadata,
		})
	}

	query, err := sql.Insert[utils.JoinExpLabel]("table_labels", exp_label_joins...)
	query.Suffixe(` ON CONFLICT (exp_id,labels)  DO NOTHING`)
	if err != nil {
		log.Default().Println("ERROR <insertTableExp>")
		return err
	}
	err = sql.Exec(context.Background(), &query, tx)
	if err != nil {
		log.Default().Println("ERROR <insertTableExp>")
		return err
	}
	return err
}

func insertVariables(nimbus_execution utils.NimbusExecution, variables []utils.Variable, tx pgx.Tx) error {
	insert_into_table_nimbus, err := sql.Insert[utils.NimbusExecution]("table_nimbus_execution", nimbus_execution)
	if err != nil {
		return err
	}
	insert_into_table_variable, err2 := sql.Insert[utils.Variable]("table_variable", variables...)
	if err2 != nil {
		return err2
	}
	query, _ := sql.SQLf("WITH nimbus_id AS (")
	query.Append(
		insert_into_table_nimbus.Suffixe(` 
		ON CONFLICT (exp_id, config_name, extension, lossless, nan_value_encoding, chunks, rx, ry)
			DO UPDATE SET created_at = excluded.created_at 
				WHERE table_nimbus_execution.created_at < excluded.created_at 
		RETURNING id),
	
		var_ids_name AS 
		(`), insert_into_table_variable.Suffixe(` RETURNING name,id)

		INSERT INTO join_nimbus_execution_variables

			SELECT 
				nimbus_id.id AS id_nimbus_execution,
				var_ids_name.name AS variable_name,
				var_ids_name.id AS variable_id

			FROM var_ids_name CROSS JOIN nimbus_id

			ON CONFLICT (id_nimbus_execution,variable_name)
			DO UPDATE SET variable_id = excluded.variable_id;
		`))
	err = sql.Exec(context.Background(), &query, tx)
	if err != nil {
		log.Default().Println("ERROR <insertVariables>")
		return err
	}
	return err
}

func InsertAll(exp_id string, request *utils.Request, pool *pgxpool.Pool) error {
	if err := pgx.BeginFunc(context.Background(), pool,
		func(tx pgx.Tx) error {
			err := insertVariables(request.Request.Table_nimbus_execution, request.Request.Table_variable, tx)
			if err != nil {
				log.Default().Println("ERROR <InsertAll>")
				return err
			}

			err = insertTableExp(request.Request.Table_experiment, tx)
			if err != nil {
				log.Default().Println("ERROR <InsertAll>")
				return err
			}

			publication_labels, err_u := updatePublicationExp(exp_id, tx)
			if err_u != nil {
				log.Default().Println("ERROR <InsertAll>")
				return err_u
			}
			err = insertTableLabels(request.Request.Table_experiment.Labels, publication_labels, request.Request.Table_experiment.Exp_id, tx)
			if err != nil {
				log.Default().Println("ERROR <InsertAll>")
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

func updatePublicationExp(exp string, tx pgx.Tx) ([]utils.Label, error) {
	query, err := sql.SQLf(`
		UPDATE join_publication_exp SET requested_exp_id = NULL, exp_id = %s
		WHERE requested_exp_id = %s
		RETURNING metadata
	`, exp, exp)
	if err != nil {
		log.Default().Println("ERROR <updatePublicationExp>")
		return nil, err
	}
	type Response struct {
		Metadata []map[string]any `sql:"metadata"`
	}
	responses, err := sql.Receive[Response](context.Background(), &query, tx)
	if err != nil {
		log.Default().Println("ERROR <updatePublicationExp>")
		return nil, err
	}
	labels := make([]utils.Label, 0, len(responses))
	for _, res := range responses {
		for _, lm := range res.Metadata {
			if l, ok := lm["label"]; ok {
				l_str := fmt.Sprintf("%s", l)
				label := utils.Label{
					Label: l_str,
				}
				if m, ok := lm["metadata"]; ok {
					switch m.(type) {
					case map[string]any:
						label.Metadata = m.(map[string]any)
					}
				}
				labels = append(labels, label)
			}

		}
	}
	return labels, err
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

func AddLabelsForId(id string, labels []utils.Label, pool *pgxpool.Pool) error {
	if err := pgx.BeginFunc(context.Background(), pool,
		func(tx pgx.Tx) error {
			err := insertTableLabels(labels, []utils.Label{}, id, tx)
			if err != nil {
				log.Default().Println("ERROR <AddLabelsForId>")
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
