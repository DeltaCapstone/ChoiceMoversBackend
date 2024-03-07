--add number of workers to job
ALTER TABLE public.jobs DROP num_workers;
--add priority level to employees
ALTER TABLE public.employee DROP priority;
--rename "rooms" to "details"
ALTER TABLE public.jobs RENAME COLUMN details TO rooms;
--add notes column to jobs
ALTER TABLE public.jobs DROP notes;

--need to double check this one
ALTER TABLE public.employee_jobs ADD requested boolean;
ALTER TABLE public.employee_jobs ADD assigned boolean; 

DROP TABLE IF EXISTS public.sessions;