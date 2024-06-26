-- This script was generated by the ERD tool in pgAdmin 4.
-- Please log an issue at https://redmine.postgresql.org/projects/pgadmin4/issues/new if you find any bugs, including reproduction steps.
BEGIN;

CREATE TABLE IF NOT EXISTS public.customers
(
    username character varying(255) UNIQUE NOT NULL,
    password_hash character varying(60) NOT NULL,  
    first_name character varying(255) NOT NULL,
    last_name character varying(255) NOT NULL,
    email character varying(255) UNIQUE NOT NULL,
    phone_primary text NOT NULL,
    phone_other1 text,
    phone_other2 text,
    phone_other3 text,
    PRIMARY KEY (username)
);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'employee_type') THEN
        EXECUTE 'CREATE TYPE Employee_type AS ENUM (''Part-time'', ''Full-time'', ''Manager'', ''Admin'')';
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'residence_type') THEN
       EXECUTE 'CREATE TYPE Residence_type AS ENUM (''Business'', ''House'',''Apartment'',''Condo'',''Storage Unit'',''Other'')';
    END IF;
END$$;

CREATE TABLE IF NOT EXISTS public.employees
(
    username character varying(255) UNIQUE NOT NULL,
    password_hash character varying(60) NOT NULL,  
    first_name character varying(255) NOT NULL,
    last_name character varying(255) NOT NULL,
    email character varying(255) UNIQUE NOT NULL,
    phone_primary text NOT NULL,
    phone_other1 text,
    phone_other2 text,
    employee_type Employee_type NOT NULL,
    employee_priority integer NOT NULL,
    PRIMARY KEY (username)
);


CREATE TABLE IF NOT EXISTS public.jobs
(
    job_id integer NOT NULL GENERATED ALWAYS AS IDENTITY,
    estimate_id integer,

    man_hours numeric(10,2) NOT NULL DEFAULT 0,
    rate numeric(10,2),
    cost numeric(10,2),

    finalized boolean NOT NULL DEFAULT False,
    actual_man_hours numeric(10,2) NOT NULL DEFAULT 0,
    final_cost numeric(10,2), 
    amount_payed numeric(10,2),

    notes TEXT,
    PRIMARY KEY (job_id)
);

CREATE TABLE IF NOT EXISTS public.estimates
(
    estimate_id integer NOT NULL GENERATED ALWAYS AS IDENTITY,
    customer_username character varying,

    load_addr_id integer,
    unload_addr_id integer,
    start_time timestamp with time zone,
    end_time timestamp with time zone,

    rooms jsonb,    
    special jsonb, 
    small_items integer,
    medium_items integer,
    large_items integer,
    boxes integer, 
    item_load integer,
    flight_mult numeric(3,1) DEFAULT 1,

    pack boolean NOT NULL DEFAULT False, 
    unpack boolean NOT NULL DEFAULT False, 
    load boolean NOT NULL DEFAULT False, 
    unload boolean NOT NULL DEFAULT False, 

    clean boolean NOT NULL DEFAULT False, 

    need_truck boolean,
    number_workers integer DEFAULT 2, 
    dist_to_job integer NOT NULL DEFAULT 0, 
    dist_move integer NOT NULL DEFAULT 0, 

    estimated_man_hours numeric(10,2) NOT NULL DEFAULT 0, 
    estimated_rate numeric(10,2),
    estimated_cost numeric(10,2),

    customer_notes TEXT NOT NULL DEFAULT 'None',

    PRIMARY KEY (estimate_id)
);

CREATE TABLE IF NOT EXISTS public.addresses
(
    address_id integer NOT NULL GENERATED ALWAYS AS IDENTITY,
    street character varying NOT NULL,
    city character varying NOT NULL,
    state character varying NOT NULL,
    zip character varying NOT NULL,
    res_type Residence_type NOT NULL,
    square_feet integer,
    flights integer NOT NULL DEFAULT 0,
    apt_num character varying NOT NULL,
    PRIMARY KEY (address_id)
);

CREATE TABLE IF NOT EXISTS public.employee_jobs
(
    employee_username character varying,
    job_id integer,
    manager_override boolean,
    PRIMARY KEY (employee_username, job_id)
);

ALTER TABLE IF EXISTS public.estimates
    ADD FOREIGN KEY (load_addr_id)
    REFERENCES public.addresses (address_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;


ALTER TABLE IF EXISTS public.estimates
    ADD FOREIGN KEY (unload_addr_id)
    REFERENCES public.addresses (address_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;


ALTER TABLE IF EXISTS public.estimates
    ADD FOREIGN KEY (customer_username)
    REFERENCES public.customers (username) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;

ALTER TABLE IF EXISTS public.jobs
    ADD FOREIGN KEY (estimate_id)
    REFERENCES public.estimates (estimate_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;

ALTER TABLE IF EXISTS public.employee_jobs
    ADD FOREIGN KEY (employee_username)
    REFERENCES public.employees (username) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;


ALTER TABLE IF EXISTS public.employee_jobs
    ADD FOREIGN KEY (job_id)
    REFERENCES public.jobs (job_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;

END;
