package services

import (
	"archive-api/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func insertNimbusExecutionSql(ne utils.NimbusExecution, pl *utils.Placeholder) string {
	insert_into_table_nimbus := "INSERT INTO table_nimbus_execution" +
		" (exp_id,config_name,extension,lossless,nan_value_encoding,threshold,chunks,rx,ry) VALUES "

	insert_into_table_nimbus += fmt.Sprintf("(%s,%s,%s,%s,%s,%s,%s,%s,%s)",
		pl.Get(ne.Exp_id),
		pl.Get(ne.Config_name),
		pl.Get(ne.Extension),
		pl.Get(ne.Lossless),
		pl.Get(ne.Nan_value_encoding),
		pl.Get(ne.Threshold),
		pl.Get(ne.Chunks),
		pl.Get(ne.Rx),
		pl.Get(ne.Ry),
	)
	return insert_into_table_nimbus
}

func insertVariablesSql(variables []utils.Variable, pl *utils.Placeholder) (string, error) {
	insert_into_table_variable := "INSERT INTO table_variable " +
		"(name, paths_ts, paths_mean, levels, timesteps, xsize, xfirst, xinc, ysize, yfirst, yinc, metadata) VALUES "
	for i, v := range variables {
		metadata, err := json.Marshal(v.Metadata)
		if err != nil {
			return "", err
		}
		insert_into_table_variable += fmt.Sprintf("(%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s)",
			pl.Get(v.Name),
			pl.Get(v.Paths_ts),
			pl.Get(v.Paths_mean),
			pl.Get(v.Levels),
			pl.Get(v.Timesteps),
			pl.Get(v.Xsize),
			pl.Get(v.Xfirst),
			pl.Get(v.Xinc),
			pl.Get(v.Ysize),
			pl.Get(v.Yfirst),
			pl.Get(v.Yinc),
			pl.Get(metadata),
		)
		if i < len(variables)-1 {
			insert_into_table_variable += ","
		}
	}
	return insert_into_table_variable, nil
}

func AddVariablesWithExp(exp_id string, request *utils.Request, pool *pgxpool.Pool) error {

	if err := pgx.BeginFunc(context.Background(), pool,
		func(tx pgx.Tx) error {
			pl := new(utils.Placeholder)
			pl.Build(0, 144)

			insert_into_table_nimbus := insertNimbusExecutionSql(request.Request.Table_nimbus_execution, pl)
			insert_into_table_variable, err_sql := insertVariablesSql(request.Request.Table_variable, pl)
			if err_sql != nil {
				return err_sql
			}
			sql := fmt.Sprintf("WITH nimbus_id AS (%s"+
				" ON CONFLICT (config_name, extension, lossless, nan_value_encoding, chunks, rx, ry)"+
				" DO UPDATE SET created_at = now() RETURNING id),"+
				" var_ids_name AS (%s RETURNING name,id)"+
				" INSERT INTO join_nimbus_execution_variables"+
				" SELECT nimbus_id.id AS id_nimbus_execution,"+
				" var_ids_name.name AS variable_name,"+
				" var_ids_name.id AS variable_id"+
				" FROM var_ids_name CROSS JOIN nimbus_id"+
				" ON CONFLICT (id_nimbus_execution,variable_name)"+
				" DO UPDATE SET variable_id = excluded.variable_id"+
				";",
				insert_into_table_nimbus, insert_into_table_variable)

			_, err := tx.Exec(context.Background(), sql, pl.Args...)
			return err
		},
	); err != nil {
		log.Default().Println("transactions error :", err)
		return err
	}
	log.Default().Println("insert success", exp_id)
	return nil
}

func Clean(pool *pgxpool.Pool) error {
	if err := pgx.BeginFunc(context.Background(), pool,
		func(tx pgx.Tx) error {
			sql := "DELETE FROM table_variable WHERE id NOT IN (SELECT variable_id FROM join_nimbus_execution_variables);"
			_, err := tx.Exec(context.Background(), sql)
			return err
		},
	); err != nil {
		log.Default().Println("transactions error :", err)
		return err
	}
	return nil
}
