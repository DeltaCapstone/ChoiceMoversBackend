--add number of workers to job
ALTER TABLE public.jobs ADD num_workers integer;
--add priority level to employees
ALTER TABLE public.employee ADD priority integer;
--rename "rooms" to "details"
ALTER TABLE public.jobs RENAME COLUMN rooms TO details;
--add notes column to jobs
ALTER TABLE public.jobs ADD notes TEXT;

--need to double check this one
ALTER TABLE public.employee_jobs DROP requested;
ALTER TABLE public.employee_jobs DROP assigned; 

-- ADDING sessions
CREATE TABLE IF NOT EXISTS public.sessions(

)