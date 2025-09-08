--
-- PostgreSQL database dump
--

\restrict hFOQ19oQcS76ZnJNwtNg1a2ikZMuLsDw41AnUvfwn6FfUtll8unwiDhCQzvwrcv

-- Dumped from database version 15.14
-- Dumped by pg_dump version 15.14

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
-- Name: admins; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.admins (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    email character varying(255) NOT NULL,
    password character varying(255) NOT NULL,
    token_version bigint DEFAULT 1,
    last_login timestamp with time zone
);


ALTER TABLE public.admins OWNER TO postgres;

--
-- Name: admins_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.admins_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.admins_id_seq OWNER TO postgres;

--
-- Name: admins_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.admins_id_seq OWNED BY public.admins.id;


--
-- Name: articles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.articles (
    id bigint NOT NULL,
    title character varying(255) NOT NULL,
    content text NOT NULL,
    slug character varying(255) NOT NULL,
    summary character varying(500),
    status character varying(20) DEFAULT 'draft'::character varying,
    published_at timestamp with time zone,
    admin_id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.articles OWNER TO postgres;

--
-- Name: articles_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.articles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.articles_id_seq OWNER TO postgres;

--
-- Name: articles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.articles_id_seq OWNED BY public.articles.id;


--
-- Name: devices; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.devices (
    id bigint NOT NULL,
    device_id character varying(100) NOT NULL,
    name character varying(255) NOT NULL,
    api_key character varying(255) NOT NULL,
    token_version bigint DEFAULT 1,
    last_seen timestamp with time zone,
    status character varying(50) DEFAULT 'inactive'::character varying,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.devices OWNER TO postgres;

--
-- Name: devices_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.devices_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.devices_id_seq OWNER TO postgres;

--
-- Name: devices_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.devices_id_seq OWNED BY public.devices.id;


--
-- Name: refresh_tokens; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.refresh_tokens (
    id bigint NOT NULL,
    token character varying(255) NOT NULL,
    user_id bigint NOT NULL,
    user_type character varying(50) NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    is_revoked boolean DEFAULT false,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.refresh_tokens OWNER TO postgres;

--
-- Name: refresh_tokens_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.refresh_tokens_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.refresh_tokens_id_seq OWNER TO postgres;

--
-- Name: refresh_tokens_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.refresh_tokens_id_seq OWNED BY public.refresh_tokens.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    name character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    password character varying(255) NOT NULL,
    token_version bigint DEFAULT 1,
    last_login timestamp with time zone,
    admin_id bigint,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: admins id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.admins ALTER COLUMN id SET DEFAULT nextval('public.admins_id_seq'::regclass);


--
-- Name: articles id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.articles ALTER COLUMN id SET DEFAULT nextval('public.articles_id_seq'::regclass);


--
-- Name: devices id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.devices ALTER COLUMN id SET DEFAULT nextval('public.devices_id_seq'::regclass);


--
-- Name: refresh_tokens id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.refresh_tokens ALTER COLUMN id SET DEFAULT nextval('public.refresh_tokens_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: admins; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.admins (id, created_at, updated_at, deleted_at, email, password, token_version, last_login) FROM stdin;
1	2025-09-08 18:45:57.435412+00	2025-09-08 18:45:57.435412+00	\N	admin@example.com	$2a$12$cohD33AH6.A2XYoG2aVGguPOHmjQD3KC6jMoMTnzRAZxu3tSfa82C	1	2025-09-08 18:45:57.434976+00
\.


--
-- Data for Name: articles; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.articles (id, title, content, slug, summary, status, published_at, admin_id, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- Data for Name: devices; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.devices (id, device_id, name, api_key, token_version, last_seen, status, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- Data for Name: refresh_tokens; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.refresh_tokens (id, token, user_id, user_type, expires_at, is_revoked, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, name, email, password, token_version, last_login, admin_id, created_at, updated_at, deleted_at) FROM stdin;
1	John Doe	john.doe@example.com	$2a$11$Zm2vva.DdPrz6lL6.Fml..CzDheleMLQJdRA589zttCDk6Nk.y92C	1	2025-09-08 18:45:57.577078+00	1	2025-09-08 18:45:57.577115+00	2025-09-08 18:45:57.577115+00	\N
2	Jane Smith	jane.smith@example.com	$2a$11$VfEkOKrrE2yA30Rzoy7.v.CWXWvaIA8YtmMFQ1OgKmH8qNb33rUYG	1	2025-09-08 18:45:57.716406+00	1	2025-09-08 18:45:57.716416+00	2025-09-08 18:45:57.716416+00	\N
3	Bob Johnson	bob.johnson@example.com	$2a$11$hvUr8rOBu3jsBf0Bp6avcuSuZXD2iP8yxW.93ZkrzOrsdOXCS72ri	1	2025-09-08 18:45:57.854115+00	1	2025-09-08 18:45:57.854126+00	2025-09-08 18:45:57.854126+00	\N
4	Alice Williams	alice.williams@example.com	$2a$11$gXLisVaZzRa1mL4xtqSLW.FFPCSM/wAMJG16i3Xr3HRltdhLB25cy	1	2025-09-08 18:45:57.994878+00	1	2025-09-08 18:45:57.994909+00	2025-09-08 18:45:57.994909+00	\N
5	Charlie Brown	charlie.brown@example.com	$2a$11$ieH.fLoJi2MQUD8OzbyF8ewJ0pbDe5viCs1Gj8mVau4bjGorfWlGC	1	2025-09-08 18:45:58.134868+00	1	2025-09-08 18:45:58.134893+00	2025-09-08 18:45:58.134893+00	\N
6	Sam Wilson	sam.wilson@example.com	$2a$11$4lvYKPha1J81KflJFRvW6.HHlA0OknTAt3/Vg2GVMhRXfDqz6nOAe	1	2025-09-08 18:45:58.274502+00	0	2025-09-08 18:45:58.274529+00	2025-09-08 18:45:58.274529+00	\N
7	Maria Rodriguez	maria.rodriguez@example.com	$2a$11$83gSLDSelvRQN2LxSmo52.LsEjSoJKlI1fJoQVy6LboC25xJDxP5C	1	2025-09-08 18:45:58.41374+00	0	2025-09-08 18:45:58.413751+00	2025-09-08 18:45:58.413751+00	\N
8	David Kim	david.kim@example.com	$2a$11$TCE1SjdX.L7xDoR8BG.WEuqt9NKUpCttwOJQPWQ0GM5d/yPu8J0ee	1	2025-09-08 18:45:58.55304+00	0	2025-09-08 18:45:58.553059+00	2025-09-08 18:45:58.553059+00	\N
\.


--
-- Name: admins_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.admins_id_seq', 1, true);


--
-- Name: articles_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.articles_id_seq', 1, false);


--
-- Name: devices_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.devices_id_seq', 1, false);


--
-- Name: refresh_tokens_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.refresh_tokens_id_seq', 1, false);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 8, true);


--
-- Name: admins admins_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.admins
    ADD CONSTRAINT admins_pkey PRIMARY KEY (id);


--
-- Name: articles articles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.articles
    ADD CONSTRAINT articles_pkey PRIMARY KEY (id);


--
-- Name: devices devices_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.devices
    ADD CONSTRAINT devices_pkey PRIMARY KEY (id);


--
-- Name: refresh_tokens refresh_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_admins_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_admins_deleted_at ON public.admins USING btree (deleted_at);


--
-- Name: idx_admins_email; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_admins_email ON public.admins USING btree (email);


--
-- Name: idx_articles_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_articles_deleted_at ON public.articles USING btree (deleted_at);


--
-- Name: idx_articles_slug; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_articles_slug ON public.articles USING btree (slug);


--
-- Name: idx_devices_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_devices_deleted_at ON public.devices USING btree (deleted_at);


--
-- Name: idx_devices_device_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_devices_device_id ON public.devices USING btree (device_id);


--
-- Name: idx_refresh_tokens_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_refresh_tokens_deleted_at ON public.refresh_tokens USING btree (deleted_at);


--
-- Name: idx_refresh_tokens_token; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_refresh_tokens_token ON public.refresh_tokens USING btree (token);


--
-- Name: idx_users_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_users_deleted_at ON public.users USING btree (deleted_at);


--
-- Name: idx_users_email; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_users_email ON public.users USING btree (email);


--
-- PostgreSQL database dump complete
--

\unrestrict hFOQ19oQcS76ZnJNwtNg1a2ikZMuLsDw41AnUvfwn6FfUtll8unwiDhCQzvwrcv

