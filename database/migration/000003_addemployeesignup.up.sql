Begin;

CREATE TABLE IF NOT EXISTS public.employee_signup(
    id uuid PRIMARY KEY,
    email varchar,
    employee_type Employee_type,
    employee_priority integer,
    signup_token varchar,
    expires_at timestamptz,
    used boolean
);

END;