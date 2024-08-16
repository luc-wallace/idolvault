--
-- PostgreSQL database dump
--

-- Dumped from database version 15.7 (Debian 15.7-0+deb12u1)
-- Dumped by pg_dump version 16.1

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
-- Name: biases; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.biases (
    username character varying(25) NOT NULL,
    idol_name character varying(30) NOT NULL,
    group_name character varying(30) NOT NULL
);


ALTER TABLE public.biases OWNER TO postgres;

--
-- Name: cards; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.cards (
    id integer NOT NULL,
    variant character varying(30) NOT NULL,
    idol_name character varying(30) NOT NULL,
    collection_name character varying(30) NOT NULL,
    group_name character varying(30) NOT NULL
);


ALTER TABLE public.cards OWNER TO postgres;

--
-- Name: cards_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.cards_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.cards_id_seq OWNER TO postgres;

--
-- Name: cards_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.cards_id_seq OWNED BY public.cards.id;


--
-- Name: collections; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.collections (
    name character varying(30) NOT NULL,
    group_name character varying(30) NOT NULL,
    release_date date NOT NULL
);


ALTER TABLE public.collections OWNER TO postgres;

--
-- Name: followers; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.followers (
    follower character varying(25) NOT NULL,
    following character varying(25) NOT NULL,
    followed_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.followers OWNER TO postgres;

--
-- Name: groups; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.groups (
    name character varying(30) NOT NULL,
    country character varying(30) NOT NULL,
    fandom character varying(20) NOT NULL,
    spotify_id character varying(30) NOT NULL,
    popularity integer NOT NULL,
    followers bigint NOT NULL,
    image_url character varying(80) NOT NULL,
    genres character varying(20)[] NOT NULL,
    updated_at timestamp with time zone
);


ALTER TABLE public.groups OWNER TO postgres;

--
-- Name: idols; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.idols (
    stage_name character varying(30) NOT NULL,
    real_name character varying(30) NOT NULL,
    group_name character varying(30) NOT NULL,
    birthday date NOT NULL,
    country character varying(30) NOT NULL,
    mbti character(4) NOT NULL
);


ALTER TABLE public.idols OWNER TO postgres;

--
-- Name: sessions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sessions (
    token text NOT NULL,
    data bytea NOT NULL,
    expiry timestamp with time zone NOT NULL
);


ALTER TABLE public.sessions OWNER TO postgres;

--
-- Name: songs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.songs (
    spotify_id character varying(30) NOT NULL,
    name character varying(50) NOT NULL,
    group_name character varying(30) NOT NULL,
    popularity integer NOT NULL,
    album_image_url character varying(80) NOT NULL
);


ALTER TABLE public.songs OWNER TO postgres;

--
-- Name: user_cards; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_cards (
    username character varying(25) NOT NULL,
    card_id integer NOT NULL
);


ALTER TABLE public.user_cards OWNER TO postgres;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    email character varying(254) NOT NULL,
    username character varying(25) NOT NULL,
    provider character varying(30) NOT NULL,
    avatar_url character varying(200) NOT NULL,
    bio character varying(300) NOT NULL
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: cards id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cards ALTER COLUMN id SET DEFAULT nextval('public.cards_id_seq'::regclass);


--
-- Name: biases bias_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.biases
    ADD CONSTRAINT bias_pkey PRIMARY KEY (username, idol_name, group_name);


--
-- Name: cards cards_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cards
    ADD CONSTRAINT cards_pkey PRIMARY KEY (id);


--
-- Name: collections collections_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.collections
    ADD CONSTRAINT collections_pkey PRIMARY KEY (name, group_name);


--
-- Name: followers followers_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.followers
    ADD CONSTRAINT followers_pkey PRIMARY KEY (follower, following);


--
-- Name: groups groups_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.groups
    ADD CONSTRAINT groups_pkey PRIMARY KEY (name);


--
-- Name: idols idols_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.idols
    ADD CONSTRAINT idols_pkey PRIMARY KEY (stage_name, group_name);


--
-- Name: sessions sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (token);


--
-- Name: songs songs_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.songs
    ADD CONSTRAINT songs_pkey PRIMARY KEY (spotify_id);


--
-- Name: user_cards user_cards_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_cards
    ADD CONSTRAINT user_cards_pkey PRIMARY KEY (username, card_id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (email);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: sessions_expiry_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX sessions_expiry_idx ON public.sessions USING btree (expiry);


--
-- Name: biases bias_group_name_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.biases
    ADD CONSTRAINT bias_group_name_fkey FOREIGN KEY (group_name) REFERENCES public.groups(name);


--
-- Name: biases bias_idol_name_group_name_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.biases
    ADD CONSTRAINT bias_idol_name_group_name_fkey FOREIGN KEY (idol_name, group_name) REFERENCES public.idols(stage_name, group_name);


--
-- Name: biases bias_username_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.biases
    ADD CONSTRAINT bias_username_fkey FOREIGN KEY (username) REFERENCES public.users(username);


--
-- Name: cards cards_collection_name_group_name_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cards
    ADD CONSTRAINT cards_collection_name_group_name_fkey FOREIGN KEY (collection_name, group_name) REFERENCES public.collections(name, group_name);


--
-- Name: cards cards_group_name_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cards
    ADD CONSTRAINT cards_group_name_fkey FOREIGN KEY (group_name) REFERENCES public.groups(name);


--
-- Name: cards cards_idol_name_group_name_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cards
    ADD CONSTRAINT cards_idol_name_group_name_fkey FOREIGN KEY (idol_name, group_name) REFERENCES public.idols(stage_name, group_name);


--
-- Name: collections collections_group_name_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.collections
    ADD CONSTRAINT collections_group_name_fkey FOREIGN KEY (group_name) REFERENCES public.groups(name);


--
-- Name: followers followers_follower_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.followers
    ADD CONSTRAINT followers_follower_fkey FOREIGN KEY (follower) REFERENCES public.users(username);


--
-- Name: followers followers_following_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.followers
    ADD CONSTRAINT followers_following_fkey FOREIGN KEY (following) REFERENCES public.users(username);


--
-- Name: idols idols_group_name_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.idols
    ADD CONSTRAINT idols_group_name_fkey FOREIGN KEY (group_name) REFERENCES public.groups(name);


--
-- Name: songs songs_group_name_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.songs
    ADD CONSTRAINT songs_group_name_fkey FOREIGN KEY (group_name) REFERENCES public.groups(name);


--
-- Name: user_cards user_cards_card_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_cards
    ADD CONSTRAINT user_cards_card_id_fkey FOREIGN KEY (card_id) REFERENCES public.cards(id);


--
-- Name: user_cards user_cards_username_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_cards
    ADD CONSTRAINT user_cards_username_fkey FOREIGN KEY (username) REFERENCES public.users(username);


--
-- PostgreSQL database dump complete
--

