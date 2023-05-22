CREATE TABLE IF NOT EXISTS public.table_experiments
(
    exp_id text COLLATE pg_catalog."default" NOT NULL,
    description text COLLATE pg_catalog."default",
    CONSTRAINT table_experiments_pkey PRIMARY KEY (exp_id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.table_experiments
    OWNER to root;

CREATE TABLE IF NOT EXISTS public.table_collections
(
    name text COLLATE pg_catalog."default" NOT NULL,
    description text COLLATE pg_catalog."default",
    CONSTRAINT table_collections_pkey PRIMARY KEY (name)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.table_collections
    OWNER to root;


CREATE TABLE IF NOT EXISTS public.join_collections_experiments
(
    exp_id text COLLATE pg_catalog."default" NOT NULL,
    collection_name text COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT collection_name FOREIGN KEY (collection_name)
        REFERENCES public.table_collections (name) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID,
    CONSTRAINT exp_id FOREIGN KEY (exp_id)
        REFERENCES public.table_experiments (exp_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.join_collections_experiments
    OWNER to root;


CREATE TABLE IF NOT EXISTS public.table_variables
(
    id serial NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    exp_id text COLLATE pg_catalog."default" NOT NULL,
    paths text[] COLLATE pg_catalog."default" NOT NULL,
    created_at date NOT NULL DEFAULT now(),
    config_name text COLLATE pg_catalog."default" NOT NULL,
    levels integer NOT NULL,
    timesteps integer NOT NULL,
    xsize integer NOT NULL,
    xfirst real NOT NULL,
    xinc real NOT NULL,
    ysize integer NOT NULL,
    yfirst real NOT NULL,
    yinc real NOT NULL,
    extension text COLLATE pg_catalog."default" NOT NULL,
    lossless boolean NOT NULL,
    nan_value_encoding integer NOT NULL,
    threshold real NOT NULL,
    chunks integer NOT NULL,
    metadata json NOT NULL,
    CONSTRAINT table_variables_pkey PRIMARY KEY (id),
    CONSTRAINT exp_id FOREIGN KEY (exp_id)
        REFERENCES public.table_experiments (exp_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.table_variables
    OWNER to root;
