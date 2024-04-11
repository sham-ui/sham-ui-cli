-- +goose Up
-- +goose StatementBegin
--
-- PostgreSQL database dump
--
-- Started on 2024-02-08 21:44:16 +07

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_table_access_method = heap;

--
-- TOC entry 201 (class 1259 OID 49080)
-- Name: http_sessions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.http_sessions (
    id bigint NOT NULL,
    key bytea,
    data bytea,
    created_on timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    modified_on timestamp with time zone,
    expires_on timestamp with time zone
);


--
-- TOC entry 200 (class 1259 OID 49078)
-- Name: http_sessions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.http_sessions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 3084 (class 0 OID 0)
-- Dependencies: 200
-- Name: http_sessions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.http_sessions_id_seq OWNED BY public.http_sessions.id;


--
-- TOC entry 205 (class 1259 OID 49103)
-- Name: members; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.members (
    id integer NOT NULL,
    name character varying(100),
    email character varying(100),
    password character varying(255),
    is_superuser boolean DEFAULT false NOT NULL
);


--
-- TOC entry 204 (class 1259 OID 49101)
-- Name: members_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.members_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 3085 (class 0 OID 0)
-- Dependencies: 204
-- Name: members_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.members_id_seq OWNED BY public.members.id;

--
-- TOC entry 2891 (class 2604 OID 49083)
-- Name: http_sessions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.http_sessions ALTER COLUMN id SET DEFAULT nextval('public.http_sessions_id_seq'::regclass);


--
-- TOC entry 2895 (class 2604 OID 49106)
-- Name: members id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.members ALTER COLUMN id SET DEFAULT nextval('public.members_id_seq'::regclass);

--
-- TOC entry 2905 (class 2606 OID 49089)
-- Name: http_sessions http_sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.http_sessions
    ADD CONSTRAINT http_sessions_pkey PRIMARY KEY (id);


--
-- TOC entry 2911 (class 2606 OID 49111)
-- Name: members member_email_unique; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.members
    ADD CONSTRAINT member_email_unique UNIQUE (email);


--
-- TOC entry 2913 (class 2606 OID 49109)
-- Name: members members_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.members
    ADD CONSTRAINT members_pkey PRIMARY KEY (id);

--
-- TOC entry 2902 (class 1259 OID 169871)
-- Name: http_sessions_expiry_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX http_sessions_expiry_idx ON public.http_sessions USING btree (expires_on);


--
-- TOC entry 2903 (class 1259 OID 169872)
-- Name: http_sessions_key_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX http_sessions_key_idx ON public.http_sessions USING btree (key);



-- Completed on 2024-02-08 21:44:16 +07

--
-- PostgreSQL database dump complete
--
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 1;
-- +goose StatementEnd
