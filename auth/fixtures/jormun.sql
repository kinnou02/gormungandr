--
-- PostgreSQL database dump
--

-- Dumped from database version 9.5.12
-- Dumped by pg_dump version 9.5.12

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


--
-- Name: fallback_mode; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.fallback_mode AS ENUM (
    'walking',
    'car',
    'bss',
    'bike'
);


--
-- Name: job_state; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.job_state AS ENUM (
    'pending',
    'running',
    'done',
    'failed'
);


--
-- Name: journey_order; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.journey_order AS ENUM (
    'arrival_time',
    'departure_time'
);


--
-- Name: max_duration_criteria; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.max_duration_criteria AS ENUM (
    'time',
    'duration'
);


--
-- Name: max_duration_fallback_mode; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.max_duration_fallback_mode AS ENUM (
    'walking',
    'bss',
    'bike',
    'car'
);


--
-- Name: metric_type; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.metric_type AS ENUM (
    'ed2nav',
    'fusio2ed',
    'gtfs2ed',
    'osm2ed',
    'geopal2ed',
    'synonym2ed',
    'poi2ed'
);


--
-- Name: source_address; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.source_address AS ENUM (
    'BANO',
    'OSM'
);


--
-- Name: source_admin; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.source_admin AS ENUM (
    'OSM'
);


--
-- Name: source_poi; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.source_poi AS ENUM (
    'FUSIO',
    'OSM'
);


--
-- Name: source_street; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.source_street AS ENUM (
    'OSM'
);


--
-- Name: traveler_type; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.traveler_type AS ENUM (
    'standard',
    'slow_walker',
    'fast_walker',
    'luggage',
    'wheelchair',
    'cyclist',
    'motorist'
);


--
-- Name: user_type; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.user_type AS ENUM (
    'with_free_instances',
    'without_free_instances',
    'super_user'
);


SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: alembic_version; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.alembic_version (
    version_num character varying(32) NOT NULL
);


--
-- Name: api; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.api (
    id integer NOT NULL,
    name text NOT NULL
);


--
-- Name: api_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.api_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: api_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.api_id_seq OWNED BY public.api.id;


--
-- Name: authorization; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."authorization" (
    user_id integer NOT NULL,
    instance_id integer NOT NULL,
    api_id integer NOT NULL
);


--
-- Name: autocomplete_parameter; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.autocomplete_parameter (
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone,
    id integer NOT NULL,
    name text NOT NULL,
    street public.source_street,
    address public.source_address,
    poi public.source_poi,
    admin public.source_admin,
    admin_level integer[] DEFAULT '{8}'::integer[] NOT NULL
);


--
-- Name: autocomplete_parameter_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.autocomplete_parameter_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: autocomplete_parameter_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.autocomplete_parameter_id_seq OWNED BY public.autocomplete_parameter.id;


--
-- Name: billing_plan; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.billing_plan (
    id integer NOT NULL,
    name text NOT NULL,
    max_request_count integer DEFAULT 0,
    max_object_count integer DEFAULT 0,
    "default" boolean NOT NULL,
    end_point_id integer DEFAULT 1 NOT NULL
);


--
-- Name: billing_plan_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.billing_plan_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: billing_plan_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.billing_plan_id_seq OWNED BY public.billing_plan.id;


--
-- Name: data_set; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.data_set (
    id integer NOT NULL,
    type text NOT NULL,
    name text NOT NULL,
    job_id integer,
    family_type text NOT NULL,
    uid uuid
);


--
-- Name: data_set_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.data_set_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: data_set_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.data_set_id_seq OWNED BY public.data_set.id;


--
-- Name: end_point; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.end_point (
    id integer NOT NULL,
    name text NOT NULL,
    "default" boolean NOT NULL
);


--
-- Name: end_point_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.end_point_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: end_point_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.end_point_id_seq OWNED BY public.end_point.id;


--
-- Name: host; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.host (
    id integer NOT NULL,
    value text NOT NULL,
    end_point_id integer NOT NULL
);


--
-- Name: host_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.host_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: host_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.host_id_seq OWNED BY public.host.id;


--
-- Name: instance; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.instance (
    id integer NOT NULL,
    name text NOT NULL,
    is_free boolean NOT NULL,
    scenario text DEFAULT 'default'::text NOT NULL,
    journey_order public.journey_order DEFAULT 'arrival_time'::public.journey_order NOT NULL,
    bike_speed double precision DEFAULT '4.09999999999999964'::double precision NOT NULL,
    bss_speed double precision DEFAULT '4.09999999999999964'::double precision NOT NULL,
    car_speed double precision DEFAULT '11.1099999999999994'::double precision NOT NULL,
    max_bike_duration_to_pt integer DEFAULT 900 NOT NULL,
    max_bss_duration_to_pt integer DEFAULT 900 NOT NULL,
    max_car_duration_to_pt integer DEFAULT 1800 NOT NULL,
    max_nb_transfers integer DEFAULT 10 NOT NULL,
    max_walking_duration_to_pt integer DEFAULT 900 NOT NULL,
    walking_speed double precision DEFAULT '1.12000000000000011'::double precision NOT NULL,
    min_bike integer DEFAULT 240 NOT NULL,
    min_car integer DEFAULT 300 NOT NULL,
    min_tc_with_bike integer DEFAULT 300 NOT NULL,
    min_tc_with_car integer DEFAULT 300 NOT NULL,
    min_bss integer DEFAULT 420 NOT NULL,
    min_tc_with_bss integer DEFAULT 300 NOT NULL,
    factor_too_long_journey double precision DEFAULT '4'::double precision NOT NULL,
    min_duration_too_long_journey integer DEFAULT 900 NOT NULL,
    max_duration_criteria public.max_duration_criteria DEFAULT 'time'::public.max_duration_criteria NOT NULL,
    max_duration_fallback_mode public.max_duration_fallback_mode DEFAULT 'walking'::public.max_duration_fallback_mode NOT NULL,
    max_duration integer DEFAULT 86400 NOT NULL,
    night_bus_filter_base_factor integer DEFAULT 3600 NOT NULL,
    night_bus_filter_max_factor double precision DEFAULT 3 NOT NULL,
    walking_transfer_penalty integer DEFAULT 120 NOT NULL,
    priority integer DEFAULT 0 NOT NULL,
    bss_provider boolean DEFAULT true,
    max_additional_connections integer DEFAULT 2 NOT NULL,
    successive_physical_mode_to_limit_id text DEFAULT 'physical_mode:Bus'::text NOT NULL,
    discarded boolean DEFAULT false NOT NULL,
    full_sn_geometries boolean DEFAULT false NOT NULL,
    is_open_data boolean NOT NULL,
    import_stops_in_mimir boolean DEFAULT false NOT NULL,
    car_park_provider boolean DEFAULT true NOT NULL,
    max_car_no_park_duration_to_pt integer DEFAULT 1800 NOT NULL,
    car_no_park_speed double precision DEFAULT '6.94000000000000039'::double precision NOT NULL,
    realtime_pool_size integer,
    import_ntfs_in_mimir boolean DEFAULT false NOT NULL
);


--
-- Name: instance_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.instance_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: instance_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.instance_id_seq OWNED BY public.instance.id;


--
-- Name: job; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.job (
    id integer NOT NULL,
    task_uuid text,
    instance_id integer,
    state public.job_state,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone,
    autocomplete_params_id integer
);


--
-- Name: job_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.job_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: job_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.job_id_seq OWNED BY public.job.id;


--
-- Name: key; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.key (
    id integer NOT NULL,
    user_id integer NOT NULL,
    token text NOT NULL,
    valid_until date,
    app_name text
);


--
-- Name: key_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.key_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: key_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.key_id_seq OWNED BY public.key.id;


--
-- Name: metric; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.metric (
    id integer NOT NULL,
    job_id integer NOT NULL,
    type public.metric_type NOT NULL,
    dataset_id integer,
    duration interval
);


--
-- Name: metric_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.metric_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: metric_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.metric_id_seq OWNED BY public.metric.id;


--
-- Name: poi_type_json; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.poi_type_json (
    poi_types_json text NOT NULL,
    instance_id integer NOT NULL
);


--
-- Name: traveler_profile; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.traveler_profile (
    coverage_id integer NOT NULL,
    traveler_type public.traveler_type NOT NULL,
    walking_speed double precision NOT NULL,
    bike_speed double precision NOT NULL,
    bss_speed double precision NOT NULL,
    car_speed double precision NOT NULL,
    wheelchair boolean NOT NULL,
    max_walking_duration_to_pt integer NOT NULL,
    max_bike_duration_to_pt integer NOT NULL,
    max_bss_duration_to_pt integer NOT NULL,
    max_car_duration_to_pt integer NOT NULL,
    first_section_mode public.fallback_mode[] NOT NULL,
    last_section_mode public.fallback_mode[] NOT NULL
);


--
-- Name: user; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."user" (
    id integer NOT NULL,
    login text NOT NULL,
    email text NOT NULL,
    type public.user_type DEFAULT 'with_free_instances'::public.user_type NOT NULL,
    end_point_id integer,
    billing_plan_id integer,
    block_until timestamp without time zone,
    shape text,
    default_coord text
);


--
-- Name: user_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.user_id_seq OWNED BY public."user".id;


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.api ALTER COLUMN id SET DEFAULT nextval('public.api_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.autocomplete_parameter ALTER COLUMN id SET DEFAULT nextval('public.autocomplete_parameter_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.billing_plan ALTER COLUMN id SET DEFAULT nextval('public.billing_plan_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.data_set ALTER COLUMN id SET DEFAULT nextval('public.data_set_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.end_point ALTER COLUMN id SET DEFAULT nextval('public.end_point_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.host ALTER COLUMN id SET DEFAULT nextval('public.host_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.instance ALTER COLUMN id SET DEFAULT nextval('public.instance_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.job ALTER COLUMN id SET DEFAULT nextval('public.job_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.key ALTER COLUMN id SET DEFAULT nextval('public.key_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.metric ALTER COLUMN id SET DEFAULT nextval('public.metric_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."user" ALTER COLUMN id SET DEFAULT nextval('public.user_id_seq'::regclass);


--
-- Data for Name: alembic_version; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.alembic_version (version_num) FROM stdin;
465a7431358a
465a7431358a
465a7431358a
465a7431358a
465a7431358a
465a7431358a
\.


--
-- Data for Name: api; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.api (id, name) FROM stdin;
1	ALL
\.


--
-- Name: api_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.api_id_seq', 1, true);


--
-- Data for Name: authorization; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public."authorization" (user_id, instance_id, api_id) FROM stdin;
1	4	1
\.


--
-- Data for Name: autocomplete_parameter; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.autocomplete_parameter (created_at, updated_at, id, name, street, address, poi, admin, admin_level) FROM stdin;
\.


--
-- Name: autocomplete_parameter_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.autocomplete_parameter_id_seq', 1, false);


--
-- Data for Name: billing_plan; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.billing_plan (id, name, max_request_count, max_object_count, "default", end_point_id) FROM stdin;
4	sncf_dev	3000	60000	t	2
1	nav_dev	3000	0	f	1
2	nav_ent	0	0	f	1
3	nav_ctp	0	0	t	1
5	sncf_ent	0	0	f	2
\.


--
-- Name: billing_plan_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.billing_plan_id_seq', 5, true);


--
-- Data for Name: data_set; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.data_set (id, type, name, job_id, family_type, uid) FROM stdin;
\.


--
-- Name: data_set_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.data_set_id_seq', 1, false);


--
-- Data for Name: end_point; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.end_point (id, name, "default") FROM stdin;
1	navitia.io	t
2	sncf	f
\.


--
-- Name: end_point_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.end_point_id_seq', 2, true);


--
-- Data for Name: host; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.host (id, value, end_point_id) FROM stdin;
\.


--
-- Name: host_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.host_id_seq', 1, false);


--
-- Data for Name: instance; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.instance (id, name, is_free, scenario, journey_order, bike_speed, bss_speed, car_speed, max_bike_duration_to_pt, max_bss_duration_to_pt, max_car_duration_to_pt, max_nb_transfers, max_walking_duration_to_pt, walking_speed, min_bike, min_car, min_tc_with_bike, min_tc_with_car, min_bss, min_tc_with_bss, factor_too_long_journey, min_duration_too_long_journey, max_duration_criteria, max_duration_fallback_mode, max_duration, night_bus_filter_base_factor, night_bus_filter_max_factor, walking_transfer_penalty, priority, bss_provider, max_additional_connections, successive_physical_mode_to_limit_id, discarded, full_sn_geometries, is_open_data, import_stops_in_mimir, car_park_provider, max_car_no_park_duration_to_pt, car_no_park_speed, realtime_pool_size, import_ntfs_in_mimir) FROM stdin;
3	fr-idf	t	default	arrival_time	4.09999999999999964	4.09999999999999964	11.1099999999999994	900	900	1800	10	900	1.12000000000000011	240	300	300	300	420	300	4	900	time	walking	86400	3600	3	120	0	t	2	physical_mode:Bus	f	f	t	f	t	1800	6.94000000000000039	\N	f
4	transilien	f	default	arrival_time	4.09999999999999964	4.09999999999999964	11.1099999999999994	900	900	1800	10	900	1.12000000000000011	240	300	300	300	420	300	4	900	time	walking	86400	3600	3	120	0	t	2	physical_mode:Bus	f	f	f	f	t	1800	6.94000000000000039	\N	f
5	sncf	f	default	arrival_time	4.09999999999999964	4.09999999999999964	11.1099999999999994	900	900	1800	10	900	1.12000000000000011	240	300	300	300	420	300	4	900	time	walking	86400	3600	3	120	0	t	2	physical_mode:Bus	f	f	f	f	t	1800	6.94000000000000039	\N	f
\.


--
-- Name: instance_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.instance_id_seq', 5, true);


--
-- Data for Name: job; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.job (id, task_uuid, instance_id, state, created_at, updated_at, autocomplete_params_id) FROM stdin;
\.


--
-- Name: job_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.job_id_seq', 1, false);


--
-- Data for Name: key; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.key (id, user_id, token, valid_until, app_name) FROM stdin;
1	1	115aa17b-63d3-4a31-acd6-edebebd4d415	\N	test
\.


--
-- Name: key_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.key_id_seq', 1, true);


--
-- Data for Name: metric; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.metric (id, job_id, type, dataset_id, duration) FROM stdin;
\.


--
-- Name: metric_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.metric_id_seq', 1, false);


--
-- Data for Name: poi_type_json; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.poi_type_json (poi_types_json, instance_id) FROM stdin;
\.


--
-- Data for Name: traveler_profile; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.traveler_profile (coverage_id, traveler_type, walking_speed, bike_speed, bss_speed, car_speed, wheelchair, max_walking_duration_to_pt, max_bike_duration_to_pt, max_bss_duration_to_pt, max_car_duration_to_pt, first_section_mode, last_section_mode) FROM stdin;
\.


--
-- Data for Name: user; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public."user" (id, login, email, type, end_point_id, billing_plan_id, block_until, shape, default_coord) FROM stdin;
1	testuser	testuser@example.com	with_free_instances	1	3	\N	null	\N
\.


--
-- Name: user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.user_id_seq', 1, true);


--
-- Name: api_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.api
    ADD CONSTRAINT api_name_key UNIQUE (name);


--
-- Name: api_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.api
    ADD CONSTRAINT api_pkey PRIMARY KEY (id);


--
-- Name: authorization_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."authorization"
    ADD CONSTRAINT authorization_pkey PRIMARY KEY (user_id, instance_id, api_id);


--
-- Name: autocomplete_parameter_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.autocomplete_parameter
    ADD CONSTRAINT autocomplete_parameter_name_key UNIQUE (name);


--
-- Name: autocomplete_parameter_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.autocomplete_parameter
    ADD CONSTRAINT autocomplete_parameter_pkey PRIMARY KEY (id);


--
-- Name: billing_plan_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.billing_plan
    ADD CONSTRAINT billing_plan_pkey PRIMARY KEY (id);


--
-- Name: data_set_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.data_set
    ADD CONSTRAINT data_set_pkey PRIMARY KEY (id);


--
-- Name: data_set_uid_idx; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.data_set
    ADD CONSTRAINT data_set_uid_idx UNIQUE (uid);


--
-- Name: end_point_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.end_point
    ADD CONSTRAINT end_point_name_key UNIQUE (name);


--
-- Name: end_point_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.end_point
    ADD CONSTRAINT end_point_pkey PRIMARY KEY (id);


--
-- Name: host_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.host
    ADD CONSTRAINT host_pkey PRIMARY KEY (id);


--
-- Name: instance_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.instance
    ADD CONSTRAINT instance_name_key UNIQUE (name);


--
-- Name: instance_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.instance
    ADD CONSTRAINT instance_pkey PRIMARY KEY (id);


--
-- Name: job_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.job
    ADD CONSTRAINT job_pkey PRIMARY KEY (id);


--
-- Name: key_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.key
    ADD CONSTRAINT key_pkey PRIMARY KEY (id);


--
-- Name: key_token_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.key
    ADD CONSTRAINT key_token_key UNIQUE (token);


--
-- Name: metric_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.metric
    ADD CONSTRAINT metric_pkey PRIMARY KEY (id);


--
-- Name: poi_type_json_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.poi_type_json
    ADD CONSTRAINT poi_type_json_pkey PRIMARY KEY (instance_id);


--
-- Name: traveler_profile_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.traveler_profile
    ADD CONSTRAINT traveler_profile_pkey PRIMARY KEY (coverage_id, traveler_type);


--
-- Name: user_email_end_point_idx; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."user"
    ADD CONSTRAINT user_email_end_point_idx UNIQUE (email, end_point_id);


--
-- Name: user_login_end_point_idx; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."user"
    ADD CONSTRAINT user_login_end_point_idx UNIQUE (login, end_point_id);


--
-- Name: user_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."user"
    ADD CONSTRAINT user_pkey PRIMARY KEY (id);


--
-- Name: authorization_api_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."authorization"
    ADD CONSTRAINT authorization_api_id_fkey FOREIGN KEY (api_id) REFERENCES public.api(id);


--
-- Name: authorization_instance_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."authorization"
    ADD CONSTRAINT authorization_instance_id_fkey FOREIGN KEY (instance_id) REFERENCES public.instance(id);


--
-- Name: authorization_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."authorization"
    ADD CONSTRAINT authorization_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id);


--
-- Name: data_set_job_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.data_set
    ADD CONSTRAINT data_set_job_id_fkey FOREIGN KEY (job_id) REFERENCES public.job(id);


--
-- Name: fk_billing_plan_end_point; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.billing_plan
    ADD CONSTRAINT fk_billing_plan_end_point FOREIGN KEY (end_point_id) REFERENCES public.end_point(id);


--
-- Name: fk_job_autocomplete_parameter; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.job
    ADD CONSTRAINT fk_job_autocomplete_parameter FOREIGN KEY (autocomplete_params_id) REFERENCES public.autocomplete_parameter(id);


--
-- Name: fk_user_billing_plan; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."user"
    ADD CONSTRAINT fk_user_billing_plan FOREIGN KEY (billing_plan_id) REFERENCES public.billing_plan(id);


--
-- Name: host_end_point_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.host
    ADD CONSTRAINT host_end_point_id_fkey FOREIGN KEY (end_point_id) REFERENCES public.end_point(id);


--
-- Name: job_instance_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.job
    ADD CONSTRAINT job_instance_id_fkey FOREIGN KEY (instance_id) REFERENCES public.instance(id);


--
-- Name: key_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.key
    ADD CONSTRAINT key_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id);


--
-- Name: metric_dataset_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.metric
    ADD CONSTRAINT metric_dataset_id_fkey FOREIGN KEY (dataset_id) REFERENCES public.data_set(id);


--
-- Name: metric_job_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.metric
    ADD CONSTRAINT metric_job_id_fkey FOREIGN KEY (job_id) REFERENCES public.job(id);


--
-- Name: poi_type_json_instance_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.poi_type_json
    ADD CONSTRAINT poi_type_json_instance_id_fkey FOREIGN KEY (instance_id) REFERENCES public.instance(id);


--
-- Name: traveler_profile_coverage_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.traveler_profile
    ADD CONSTRAINT traveler_profile_coverage_id_fkey FOREIGN KEY (coverage_id) REFERENCES public.instance(id);


--
-- PostgreSQL database dump complete
--

