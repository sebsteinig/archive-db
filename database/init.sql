CREATE TABLE IF NOT EXISTS public.table_variable
(
    id serial NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    paths_ts text[] NOT NULL,
    paths_mean text[] NOT NULL,
    levels integer NOT NULL,
    timesteps integer NOT NULL,
    xsize integer NOT NULL,
    xfirst real NOT NULL,
    xinc real NOT NULL,
    ysize integer NOT NULL,
    yfirst real NOT NULL,
    yinc real NOT NULL,
    metadata json NOT NULL,
    CONSTRAINT table_variable_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.table_variable
    OWNER to root;

CREATE TABLE IF NOT EXISTS public.table_nimbus_execution
(
    id serial NOT NULL,
    created_at date NOT NULL DEFAULT now(),
    config_name text COLLATE pg_catalog."default" NOT NULL,
    extension text COLLATE pg_catalog."default" NOT NULL,
    lossless boolean NOT NULL,
    nan_value_encoding integer NOT NULL,
    chunks integer NOT NULL,
    rx real NOT NULL,
    ry real NOT NULL,
    exp_id text COLLATE pg_catalog."default" NOT NULL,
    threshold real NOT NULL,
    CONSTRAINT table_nimbus_execution_pkey PRIMARY KEY (id),
    CONSTRAINT unique_config UNIQUE (config_name, extension, lossless, nan_value_encoding, chunks, rx, ry)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.table_nimbus_execution
    OWNER to root;


CREATE TABLE IF NOT EXISTS public.table_exp
(
    exp_id text COLLATE pg_catalog."default" NOT NULL,
    age bigint,
    metadata json,
    CONSTRAINT exp_id PRIMARY KEY (exp_id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.table_exp
    OWNER to root;

CREATE TABLE IF NOT EXISTS public.join_nimbus_execution_variables
(
    id_nimbus_execution serial NOT NULL,
    variable_name text COLLATE pg_catalog."default" NOT NULL,
    variable_id serial NOT NULL,
    CONSTRAINT unique_set UNIQUE (id_nimbus_execution, variable_name, variable_id),
    CONSTRAINT fk_nimbus_execution FOREIGN KEY (id_nimbus_execution)
        REFERENCES public.table_nimbus_execution (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID,
    CONSTRAINT fk_variable FOREIGN KEY (variable_id)
        REFERENCES public.table_variable (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.join_nimbus_execution_variables
    OWNER to root;
