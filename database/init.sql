CREATE TABLE IF NOT EXISTS public.table_variable
(
    id serial NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    paths_ts json NOT NULL,
    paths_mean json NOT NULL,
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
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    config_name text COLLATE pg_catalog."default" NOT NULL,
    extension text COLLATE pg_catalog."default" NOT NULL,
    lossless boolean NOT NULL,
    nan_value_encoding integer NOT NULL,
    -- chunks_time integer NOT NULL,
    -- chunks_vertical integer NOT NULL,
    rx real NOT NULL,
    ry real NOT NULL,
    exp_id text COLLATE pg_catalog."default" NOT NULL,
    threshold real NOT NULL,
    CONSTRAINT table_nimbus_execution_pkey PRIMARY KEY (id),
    CONSTRAINT unique_config UNIQUE (exp_id, config_name, extension, lossless, nan_value_encoding, rx, ry)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.table_nimbus_execution
    OWNER to root;
CREATE INDEX ON table_nimbus_execution (exp_id);


CREATE TABLE IF NOT EXISTS public.table_exp
(
    exp_id text COLLATE pg_catalog."default" NOT NULL,
    co2 real,
    realistic boolean,
    coast_line_id bigint,
    gmst real,
    date_wp_created date,
    date_wp_updated date,
    metadata json,
    CONSTRAINT exp_id PRIMARY KEY (exp_id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.table_exp
    OWNER to root;
CREATE INDEX ON table_exp (exp_id);

CREATE TABLE IF NOT EXISTS public.join_nimbus_execution_variables
(
    id_nimbus_execution serial NOT NULL,
    variable_name text COLLATE pg_catalog."default" NOT NULL,
    variable_id serial NOT NULL,
    CONSTRAINT unique_set UNIQUE (id_nimbus_execution, variable_name),
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

CREATE TABLE IF NOT EXISTS public.table_labels
(
    exp_id text COLLATE pg_catalog."default" NOT NULL,
    labels text COLLATE pg_catalog."default" NOT NULL,
    metadata json NOT NULL DEFAULT '{}'::json,
    CONSTRAINT table_labels_exp_id_label_key UNIQUE (exp_id, labels),
    CONSTRAINT table_labels_exp_id_fkey FOREIGN KEY (exp_id)
        REFERENCES public.table_exp (exp_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

TABLESPACE pg_default;
ALTER TABLE IF EXISTS public.table_labels
    OWNER to root;

CREATE INDEX ON table_labels (exp_id);
CREATE INDEX ON table_labels (labels);

CREATE TABLE IF NOT EXISTS public.table_publication
(
    id serial NOT NULL,
    title text COLLATE pg_catalog."default" NOT NULL DEFAULT ''::text,
    authors_short text COLLATE pg_catalog."default" NOT NULL DEFAULT ''::text,
    authors_full text COLLATE pg_catalog."default" NOT NULL DEFAULT ''::text,
    journal text COLLATE pg_catalog."default" NOT NULL DEFAULT ''::text,
    year bigint NOT NULL DEFAULT 0,
    owner_name text COLLATE pg_catalog."default" NOT NULL DEFAULT ''::text,
    owner_email text COLLATE pg_catalog."default" NOT NULL DEFAULT ''::text,
    brief_desc text COLLATE pg_catalog."default" NOT NULL DEFAULT ''::text,
    abstract text COLLATE pg_catalog."default" NOT NULL DEFAULT ''::text,
    expts_paper text[] COLLATE pg_catalog."default" NOT NULL DEFAULT '{}'::text[],
    CONSTRAINT table_publication_pkey PRIMARY KEY (id),
    CONSTRAINT table_publication_title_journal_year_owner_name_key UNIQUE (title, journal, year, owner_name)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.table_publication
    OWNER to root;
CREATE TABLE IF NOT EXISTS public.join_publication_exp
(
    publication_id serial NOT NULL,
    requested_exp_id text COLLATE pg_catalog."default",
    exp_id text COLLATE pg_catalog."default",
    metadata json NOT NULL DEFAULT '{}'::json,
    CONSTRAINT join_publication_exp_exp_id_publication_id_key UNIQUE (exp_id, publication_id),
    CONSTRAINT join_publication_expid_expid_fkey FOREIGN KEY (exp_id)
        REFERENCES public.table_exp (exp_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT join_publication_expid_publication_id_fkey FOREIGN KEY (publication_id)
        REFERENCES public.table_publication (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.join_publication_exp
    OWNER to root;
CREATE INDEX ON join_publication_exp (exp_id);
CREATE INDEX ON join_publication_exp (publication_id);