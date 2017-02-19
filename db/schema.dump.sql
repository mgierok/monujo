--
-- PostgreSQL database dump
--

-- Dumped from database version 9.4.1
-- Dumped by pg_dump version 9.5.1

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

SET search_path = "public", pg_catalog;

--
-- Name: latest_quotes_before_insert(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION "latest_quotes_before_insert"() RETURNS "trigger"
    LANGUAGE "plpgsql"
    AS $$
                    BEGIN
                        DELETE FROM latest_quotes where "ticker" = NEW."ticker";
                        RETURN NEW;
                    END;
                    $$;


--
-- Name: quotes_before_insert(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION "quotes_before_insert"() RETURNS "trigger"
    LANGUAGE "plpgsql"
    AS $$
                    BEGIN
                        DELETE FROM quotes where "ticker" = NEW."ticker" AND "date" = NEW."date";
                        RETURN NEW;
                    END;
                    $$;


SET default_with_oids = false;

--
-- Name: latest_quotes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "latest_quotes" (
    "ticker" character(8) NOT NULL,
    "date" "date" NOT NULL,
    "open" real,
    "high" real,
    "low" real,
    "close" real,
    "volume" bigint,
    "openint" bigint
);

--
-- Name: quotes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "quotes" (
    "ticker" character(8) NOT NULL,
    "date" "date" NOT NULL,
    "open" real,
    "high" real,
    "low" real,
    "close" real,
    "volume" bigint,
    "openint" bigint
);

--
-- Name: latest_quotes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "latest_quotes"
    ADD CONSTRAINT "latest_quotes_pkey" PRIMARY KEY ("ticker");


--
-- Name: quotes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "quotes"
    ADD CONSTRAINT "quotes_pkey" PRIMARY KEY ("ticker", "date");


--
-- Name: latest_quotes_before_insert_trigger; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER "latest_quotes_before_insert_trigger" BEFORE INSERT ON "latest_quotes" FOR EACH ROW EXECUTE PROCEDURE "latest_quotes_before_insert"();


--
-- Name: quotes_before_insert_trigger; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER "quotes_before_insert_trigger" BEFORE INSERT ON "quotes" FOR EACH ROW EXECUTE PROCEDURE "quotes_before_insert"();


--
-- PostgreSQL database dump complete
--

