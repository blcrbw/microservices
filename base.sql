DROP TABLE IF EXISTS test_records CASCADE;

CREATE TABLE public.test_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    text text,
    created timestamp without time zone,
    stored timestamp without time zone
);
