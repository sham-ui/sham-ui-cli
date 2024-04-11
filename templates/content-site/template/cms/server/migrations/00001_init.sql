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
-- TOC entry 209 (class 1259 OID 49142)
-- Name: article; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.article (
    id integer NOT NULL,
    title character varying(100),
    slug character varying(100),
    category_id integer NOT NULL,
    short_body text,
    body text,
    published_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- TOC entry 208 (class 1259 OID 49140)
-- Name: article_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.article_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 3081 (class 0 OID 0)
-- Dependencies: 208
-- Name: article_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.article_id_seq OWNED BY public.article.id;


--
-- TOC entry 211 (class 1259 OID 49165)
-- Name: article_tag; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.article_tag (
    id integer NOT NULL,
    tag_id integer,
    article_id integer
);


--
-- TOC entry 210 (class 1259 OID 49163)
-- Name: article_tag_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.article_tag_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 3082 (class 0 OID 0)
-- Dependencies: 210
-- Name: article_tag_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.article_tag_id_seq OWNED BY public.article_tag.id;


--
-- TOC entry 213 (class 1259 OID 169890)
-- Name: category; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.category (
    id integer NOT NULL,
    name character varying(100),
    slug character varying(100)
);


--
-- TOC entry 212 (class 1259 OID 169888)
-- Name: category_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.category_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 3083 (class 0 OID 0)
-- Dependencies: 212
-- Name: category_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.category_id_seq OWNED BY public.category.id;


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
-- TOC entry 207 (class 1259 OID 49129)
-- Name: tag; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tag (
    id integer NOT NULL,
    name character varying(100),
    slug character varying(100)
);


--
-- TOC entry 206 (class 1259 OID 49127)
-- Name: tag_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tag_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 3087 (class 0 OID 0)
-- Dependencies: 206
-- Name: tag_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tag_id_seq OWNED BY public.tag.id;


--
-- TOC entry 2898 (class 2604 OID 49145)
-- Name: article id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.article ALTER COLUMN id SET DEFAULT nextval('public.article_id_seq'::regclass);


--
-- TOC entry 2900 (class 2604 OID 49168)
-- Name: article_tag id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.article_tag ALTER COLUMN id SET DEFAULT nextval('public.article_tag_id_seq'::regclass);


--
-- TOC entry 2901 (class 2604 OID 169893)
-- Name: category id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.category ALTER COLUMN id SET DEFAULT nextval('public.category_id_seq'::regclass);


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
-- TOC entry 2897 (class 2604 OID 49132)
-- Name: tag id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tag ALTER COLUMN id SET DEFAULT nextval('public.tag_id_seq'::regclass);


--
-- TOC entry 2923 (class 2606 OID 49158)
-- Name: article article_name_unique; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.article
    ADD CONSTRAINT article_name_unique UNIQUE (title);


--
-- TOC entry 2925 (class 2606 OID 49151)
-- Name: article article_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.article
    ADD CONSTRAINT article_pkey PRIMARY KEY (id);


--
-- TOC entry 2928 (class 2606 OID 49160)
-- Name: article article_slug_unique; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.article
    ADD CONSTRAINT article_slug_unique UNIQUE (slug);


--
-- TOC entry 2932 (class 2606 OID 49170)
-- Name: article_tag article_tag_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.article_tag
    ADD CONSTRAINT article_tag_pkey PRIMARY KEY (id);


--
-- TOC entry 2935 (class 2606 OID 49182)
-- Name: article_tag article_tag_tag_id_article_id_unique; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.article_tag
    ADD CONSTRAINT article_tag_tag_id_article_id_unique UNIQUE (tag_id, article_id);


--
-- TOC entry 2937 (class 2606 OID 169897)
-- Name: category category_name_unique; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.category
    ADD CONSTRAINT category_name_unique UNIQUE (name);


--
-- TOC entry 2939 (class 2606 OID 169895)
-- Name: category category_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.category
    ADD CONSTRAINT category_pkey PRIMARY KEY (id);


--
-- TOC entry 2942 (class 2606 OID 169899)
-- Name: category category_slug_unique; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.category
    ADD CONSTRAINT category_slug_unique UNIQUE (slug);


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
-- TOC entry 2915 (class 2606 OID 49136)
-- Name: tag tag_name_unique; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tag
    ADD CONSTRAINT tag_name_unique UNIQUE (name);


--
-- TOC entry 2917 (class 2606 OID 49134)
-- Name: tag tag_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tag
    ADD CONSTRAINT tag_pkey PRIMARY KEY (id);


--
-- TOC entry 2920 (class 2606 OID 49138)
-- Name: tag tag_slug_unique; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tag
    ADD CONSTRAINT tag_slug_unique UNIQUE (slug);


--
-- TOC entry 2921 (class 1259 OID 49162)
-- Name: article_category_id_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX article_category_id_index ON public.article USING btree (category_id);


--
-- TOC entry 2926 (class 1259 OID 49161)
-- Name: article_slug_index; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX article_slug_index ON public.article USING btree (slug);


--
-- TOC entry 2930 (class 1259 OID 49184)
-- Name: article_tag_article_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX article_tag_article_id ON public.article_tag USING btree (article_id);


--
-- TOC entry 2933 (class 1259 OID 49183)
-- Name: article_tag_tag_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX article_tag_tag_id ON public.article_tag USING btree (tag_id);


--
-- TOC entry 2940 (class 1259 OID 169900)
-- Name: category_slug_index; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX category_slug_index ON public.category USING btree (slug);


--
-- TOC entry 2929 (class 1259 OID 169945)
-- Name: fki_article_category_id_fkey; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX fki_article_category_id_fkey ON public.article USING btree (category_id);


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


--
-- TOC entry 2918 (class 1259 OID 49139)
-- Name: tag_slug_index; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX tag_slug_index ON public.tag USING btree (slug);


--
-- TOC entry 2943 (class 2606 OID 169940)
-- Name: article article_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.article
    ADD CONSTRAINT article_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.category(id) ON DELETE CASCADE;


--
-- TOC entry 2945 (class 2606 OID 49176)
-- Name: article_tag article_tag_article_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.article_tag
    ADD CONSTRAINT article_tag_article_id_fkey FOREIGN KEY (article_id) REFERENCES public.article(id) ON DELETE CASCADE;


--
-- TOC entry 2944 (class 2606 OID 49171)
-- Name: article_tag article_tag_tag_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.article_tag
    ADD CONSTRAINT article_tag_tag_id_fkey FOREIGN KEY (tag_id) REFERENCES public.tag(id) ON DELETE CASCADE;


-- Completed on 2024-02-08 21:44:16 +07

--
-- PostgreSQL database dump complete
--
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 1;
-- +goose StatementEnd
