Begin;

CREATE TABLE IF NOT EXISTS public.employee_signup(
    id uuid PRIMARY KEY,
    email varchar,
    signup_token varchar,
    expires_at timestamptz,
    used boolean
);

END;