--
-- PostgreSQL database dump
--

-- Dumped from database version 16.1 (Debian 16.1-1.pgdg120+1)
-- Dumped by pg_dump version 16.1 (Debian 16.1-1.pgdg120+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: test_records; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public.test_records (
    id integer,
    text text,
    created timestamp without time zone,
    stored timestamp without time zone
);


ALTER TABLE public.test_records OWNER TO "user";

--
-- Name: test_records_id_seq; Type: SEQUENCE; Schema: public; Owner: user
--

CREATE SEQUENCE public.test_records_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.test_records_id_seq OWNER TO "user";

--
-- Name: test_records_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: user
--

ALTER SEQUENCE public.test_records_id_seq OWNED BY public.test_records.id;


--
-- Name: test_records id; Type: DEFAULT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.test_records ALTER COLUMN id SET DEFAULT nextval('public.test_records_id_seq'::regclass);


--
-- Name: test_records test_records_pkey; Type: CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.test_records
    ADD CONSTRAINT test_records_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--
