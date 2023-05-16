CREATE TABLE IF NOT EXISTS public.table_conversion_info
(
    pk_id SERIAL PRIMARY KEY,
    creation_at timestamp with time zone NOT NULL DEFAULT now(),
    levels bigint NOT NULL DEFAULT 1,
    time_steps bigint NOT NULL DEFAULT 1,
    xsize bigint NOT NULL,
    ysize bigint NOT NULL,
    xfirst real NOT NULL,
    xinc real NOT NULL,
    yfirst real NOT NULL,
    yinc real NOT NULL,
    nan_value_encoding smallint NOT NULL DEFAULT 255,
    threshold real NOT NULL DEFAULT 3
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.table_conversion_info
    OWNER to root;

CREATE TABLE IF NOT EXISTS public.table_execution_info
(
    pk_id SERIAL PRIMARY KEY,
    exp_id text COLLATE pg_catalog."default" NOT NULL,
    config text COLLATE pg_catalog."default" NOT NULL,
    nimbus_version text COLLATE pg_catalog."default" NOT NULL,
    execution_date timestamp without time zone NOT NULL DEFAULT now(),
    var_clt boolean NOT NULL DEFAULT false,
    var_currents boolean NOT NULL DEFAULT false,
    var_height boolean NOT NULL DEFAULT false,
    var_liconc boolean NOT NULL DEFAULT false,
    var_mlosts boolean NOT NULL DEFAULT false,
    var_pfts boolean NOT NULL DEFAULT false,
    var_pr boolean NOT NULL DEFAULT false,
    var_sic boolean NOT NULL DEFAULT false,
    var_snc boolean NOT NULL DEFAULT false,
    var_tas boolean NOT NULL DEFAULT false,
    var_tos boolean NOT NULL DEFAULT false,
    var_winds boolean NOT NULL DEFAULT false,
    metadata json,
    fk_id_conversion_info integer NOT NULL,
    CONSTRAINT fk_id_conversion_info FOREIGN KEY (fk_id_conversion_info)
        REFERENCES public.table_conversion_info (pk_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.table_execution_info
    OWNER to root;

CREATE TABLE IF NOT EXISTS public.table_image_paths
(
    pk_id BIGSERIAL PRIMARY KEY,
    fk_id_conversion_info integer NOT NULL,
    path path NOT NULL,
    extension text COLLATE pg_catalog."default",
    metadata json,
    CONSTRAINT fk_id_conversion_info FOREIGN KEY (pk_id)
        REFERENCES public.table_conversion_info (pk_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.table_image_paths
    OWNER to root;
