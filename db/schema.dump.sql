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
-- Name: currency; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE "currency" AS ENUM (
    'PLN',
    'USD',
    'EUR'
);


--
-- Name: transaction_operation_type; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE "transaction_operation_type" AS ENUM (
    'buy',
    'sell'
);


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
-- Name: portfolios; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "portfolios" (
    "portfolio_id" integer NOT NULL,
    "name" character varying(128) NOT NULL,
    "currency" "currency" NOT NULL
);


--
-- Name: portfolios_portfolio_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE "portfolios_portfolio_id_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: portfolios_portfolio_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE "portfolios_portfolio_id_seq" OWNED BY "portfolios"."portfolio_id";


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
-- Name: transactions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "transactions" (
    "transaction_id" integer NOT NULL,
    "portfolio_id" integer NOT NULL,
    "date" "date" NOT NULL,
    "ticker" character(8) NOT NULL,
    "price" real NOT NULL,
    "type" "transaction_operation_type" NOT NULL,
    "currency" "currency" NOT NULL,
    "shares" real NOT NULL,
    "commision" real NOT NULL,
    "exchange_rate" real NOT NULL,
    CONSTRAINT "transactions_commision_check" CHECK (("commision" >= (0)::double precision)),
    CONSTRAINT "transactions_exchange_rate_check" CHECK (("exchange_rate" > (0)::double precision)),
    CONSTRAINT "transactions_price_check" CHECK (("price" > (0)::double precision)),
    CONSTRAINT "transactions_shares_check" CHECK (("shares" > (0)::double precision))
);


--
-- Name: shares; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW "shares" AS
 SELECT "t"."portfolio_id",
    "t"."ticker",
    "sum"(
        CASE
            WHEN ("t"."type" = 'sell'::"transaction_operation_type") THEN ("t"."shares" * ((-1))::double precision)
            ELSE ("t"."shares")::double precision
        END) AS "shares",
    "q"."close" AS "last_price",
    "round"((("sum"(
        CASE
            WHEN ("t"."type" = 'sell'::"transaction_operation_type") THEN ("t"."shares" * ((-1))::double precision)
            ELSE ("t"."shares")::double precision
        END) * "q"."close"))::numeric, 2) AS "market_value",
    "t"."currency",
        CASE
            WHEN ("t"."currency" = "p"."currency") THEN (1)::real
            ELSE "e"."close"
        END AS "exchange_rate",
    "round"(((("sum"(
        CASE
            WHEN ("t"."type" = 'sell'::"transaction_operation_type") THEN ("t"."shares" * ((-1))::double precision)
            ELSE ("t"."shares")::double precision
        END) * "q"."close") *
        CASE
            WHEN ("t"."currency" = "p"."currency") THEN (1)::real
            ELSE "e"."close"
        END))::numeric) AS "market_value_base_currency",
    "round"((("sum"((("t"."shares" * "t"."price") * (
        CASE
            WHEN ("t"."type" = 'buy'::"transaction_operation_type") THEN 1
            ELSE (-1)
        END)::double precision)) / "sum"(
        CASE
            WHEN ("t"."type" = 'sell'::"transaction_operation_type") THEN ("t"."shares" * ((-1))::double precision)
            ELSE ("t"."shares")::double precision
        END)))::numeric, 2) AS "average_price",
    "round"(((("sum"(
        CASE
            WHEN ("t"."type" = 'sell'::"transaction_operation_type") THEN ("t"."shares" * ((-1))::double precision)
            ELSE ("t"."shares")::double precision
        END) * "q"."close") - "sum"((("t"."shares" * "t"."price") * (
        CASE
            WHEN ("t"."type" = 'buy'::"transaction_operation_type") THEN 1
            ELSE (-1)
        END)::double precision))))::numeric, 2) AS "gain",
    "round"(((((("sum"(
        CASE
            WHEN ("t"."type" = 'sell'::"transaction_operation_type") THEN ("t"."shares" * ((-1))::double precision)
            ELSE ("t"."shares")::double precision
        END) * "q"."close") - "sum"((("t"."shares" * "t"."price") * (
        CASE
            WHEN ("t"."type" = 'buy'::"transaction_operation_type") THEN 1
            ELSE (-1)
        END)::double precision))) / "sum"((("t"."shares" * "t"."price") * (
        CASE
            WHEN ("t"."type" = 'buy'::"transaction_operation_type") THEN 1
            ELSE (-1)
        END)::double precision))) * (100)::double precision))::numeric, 2) AS "percentage_gain",
    "round"((((("sum"(
        CASE
            WHEN ("t"."type" = 'sell'::"transaction_operation_type") THEN ("t"."shares" * ((-1))::double precision)
            ELSE ("t"."shares")::double precision
        END) * "q"."close") *
        CASE
            WHEN ("t"."currency" = "p"."currency") THEN (1)::real
            ELSE "e"."close"
        END) - "sum"(((("t"."shares" * "t"."price") * "t"."exchange_rate") * (
        CASE
            WHEN ("t"."type" = 'buy'::"transaction_operation_type") THEN 1
            ELSE (-1)
        END)::double precision))))::numeric, 2) AS "gain_base_currency",
    "round"((((((("sum"(
        CASE
            WHEN ("t"."type" = 'sell'::"transaction_operation_type") THEN ("t"."shares" * ((-1))::double precision)
            ELSE ("t"."shares")::double precision
        END) * "q"."close") *
        CASE
            WHEN ("t"."currency" = "p"."currency") THEN (1)::real
            ELSE "e"."close"
        END) - "sum"(((("t"."shares" * "t"."price") * "t"."exchange_rate") * (
        CASE
            WHEN ("t"."type" = 'buy'::"transaction_operation_type") THEN 1
            ELSE (-1)
        END)::double precision))) / "sum"(((("t"."shares" * "t"."price") * "t"."exchange_rate") * (
        CASE
            WHEN ("t"."type" = 'buy'::"transaction_operation_type") THEN 1
            ELSE (-1)
        END)::double precision))) * (100)::double precision))::numeric, 2) AS "percentage_gain_base_currency"
   FROM ((("transactions" "t"
     JOIN "portfolios" "p" ON (("t"."portfolio_id" = "p"."portfolio_id")))
     LEFT JOIN "latest_quotes" "q" ON (("t"."ticker" = "q"."ticker")))
     LEFT JOIN "latest_quotes" "e" ON (((("e"."ticker")::"text" = (("t"."currency")::"text" || ("p"."currency")::"text")) AND ("t"."currency" <> "p"."currency"))))
  GROUP BY "t"."portfolio_id", "t"."ticker", "q"."close", "t"."currency", "e"."close", "p"."currency"
 HAVING ("sum"(
        CASE
            WHEN ("t"."type" = 'sell'::"transaction_operation_type") THEN ("t"."shares" * ((-1))::double precision)
            ELSE ("t"."shares")::double precision
        END) > (0)::double precision);


--
-- Name: transactions_transaction_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE "transactions_transaction_id_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: transactions_transaction_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE "transactions_transaction_id_seq" OWNED BY "transactions"."transaction_id";


--
-- Name: portfolio_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY "portfolios" ALTER COLUMN "portfolio_id" SET DEFAULT "nextval"('"portfolios_portfolio_id_seq"'::"regclass");


--
-- Name: transaction_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY "transactions" ALTER COLUMN "transaction_id" SET DEFAULT "nextval"('"transactions_transaction_id_seq"'::"regclass");


--
-- Name: latest_quotes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "latest_quotes"
    ADD CONSTRAINT "latest_quotes_pkey" PRIMARY KEY ("ticker");


--
-- Name: portfolios_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "portfolios"
    ADD CONSTRAINT "portfolios_pkey" PRIMARY KEY ("portfolio_id");


--
-- Name: quotes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "quotes"
    ADD CONSTRAINT "quotes_pkey" PRIMARY KEY ("ticker", "date");


--
-- Name: transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "transactions"
    ADD CONSTRAINT "transactions_pkey" PRIMARY KEY ("transaction_id");


--
-- Name: latest_quotes_before_insert_trigger; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER "latest_quotes_before_insert_trigger" BEFORE INSERT ON "latest_quotes" FOR EACH ROW EXECUTE PROCEDURE "latest_quotes_before_insert"();


--
-- Name: quotes_before_insert_trigger; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER "quotes_before_insert_trigger" BEFORE INSERT ON "quotes" FOR EACH ROW EXECUTE PROCEDURE "quotes_before_insert"();


--
-- Name: transactions_portfolio_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "transactions"
    ADD CONSTRAINT "transactions_portfolio_id_fkey" FOREIGN KEY ("portfolio_id") REFERENCES "portfolios"("portfolio_id");


--
-- PostgreSQL database dump complete
--

