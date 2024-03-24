Begin;

CREATE TABLE IF NOT EXISTS public.password_resets(
    code varchar PRIMARY KEY,
    username varchar,
    email varchar,
    role varchar,
    expires_at timestamptz
);

END;