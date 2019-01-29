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

SET search_path = public, pg_catalog;

--
-- Name: currency; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE currency AS ENUM (
    'PLN',
    'USD',
    'EUR'
);


--
-- Name: financing_operation; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE financing_operation AS ENUM (
    'deposit',
    'withdraw',
    'dividend',
    'bond interest',
    'compounding'
);


--
-- Name: quotes_source; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE quotes_source AS ENUM (
    'stooq',
    'google',
    'ingturbo',
    'alphavantage',
    'bankier',
    'sw'
);


--
-- Name: agg_first(anyarray, anyelement, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION agg_first(p_state anyarray, p_new_element anyelement, p_limit integer) RETURNS anyarray
    LANGUAGE sql
    AS $$
select case
    when coalesce( array_length( p_state, 1 ), 0 ) < p_limit
         then p_state || p_new_element
    else p_state
     end;
$$;


--
-- Name: latest_quotes_before_insert(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION latest_quotes_before_insert() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
                    BEGIN
                        DELETE FROM latest_quotes where "ticker" = NEW."ticker";
                        RETURN NEW;
                    END;
                    $$;


--
-- Name: quotes_before_insert(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION quotes_before_insert() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
                    BEGIN
                        DELETE FROM quotes where "ticker" = NEW."ticker" AND "date" = NEW."date";
                        RETURN NEW;
                    END;
                    $$;


--
-- Name: transactions_after_insert(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION transactions_after_insert() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
	v_shares transactions.shares%TYPE;
	v_remaining_shares transactions.shares%TYPE;
	v_transaction_id disposals.in_transaction_id%TYPE;
	
	BEGIN
		SELECT SUM(t.shares) - SUM(COALESCE(d.disposed_shares, 0)) INTO v_remaining_shares
		FROM transactions t
		LEFT JOIN (SELECT in_transaction_id, SUM(disposed_shares) disposed_shares FROM disposals GROUP BY in_transaction_id) d ON d.in_transaction_id = t.transaction_id
		WHERE t.portfolio_id = NEW."portfolio_id" AND t.ticker = NEW."ticker" AND t.transaction_id != NEW.transaction_id
		GROUP BY t.ticker
		HAVING (SUM(t.shares) - SUM(COALESCE(d.disposed_shares, 0))) <> 0;

		IF NOT FOUND THEN
			RETURN NEW;
		END IF;

		IF (v_remaining_shares * NEW.shares < 0) THEN -- have different sign
			v_shares := NEW."shares";

			WHILE v_shares <> 0 LOOP
				SELECT t.shares - SUM(COALESCE(d.disposed_shares, 0)), t.transaction_id INTO v_remaining_shares, v_transaction_id
				FROM transactions t
				LEFT JOIN disposals d ON t.transaction_id = d.in_transaction_id
				WHERE t.portfolio_id = NEW."portfolio_id" AND t.ticker = NEW."ticker" AND t.transaction_id != NEW.transaction_id
				GROUP BY t.transaction_id, t.shares, t."date"
				HAVING (t.shares - SUM(COALESCE(d.disposed_shares, 0))) <> 0
				ORDER BY t."date" ASC, t.transaction_id ASC
				LIMIT 1;

				IF NOT FOUND THEN
					RETURN NEW;
				END IF;

				IF ABS(v_remaining_shares) >= ABS(v_shares) THEN
					INSERT INTO disposals (in_transaction_id, out_transaction_id, disposed_shares, disposed) VALUES (v_transaction_id, NEW."transaction_id", v_shares * -1, true);
					INSERT INTO disposals (in_transaction_id, out_transaction_id, disposed_shares, disposed) VALUES (NEW."transaction_id", v_transaction_id, v_shares, false);
					v_shares := 0;
				ELSE
					INSERT INTO disposals (in_transaction_id, out_transaction_id, disposed_shares, disposed) VALUES (v_transaction_id, NEW."transaction_id", v_remaining_shares, true);
					INSERT INTO disposals (in_transaction_id, out_transaction_id, disposed_shares, disposed) VALUES (NEW."transaction_id", v_transaction_id, v_remaining_shares * -1, false);
					v_shares := v_shares + v_remaining_shares;
				END IF;
			END LOOP;
		END IF;
		RETURN NEW;
	END;
$$;


--
-- Name: first(anyelement, integer); Type: AGGREGATE; Schema: public; Owner: -
--

CREATE AGGREGATE first(anyelement, integer) (
    SFUNC = agg_first,
    STYPE = anyarray,
    INITCOND = '{}'
);


SET default_with_oids = false;

--
-- Name: annual_range; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE annual_range (
    year smallint NOT NULL,
    range daterange NOT NULL
);


--
-- Name: currencies; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW currencies AS
 SELECT e.enumlabel AS currency
   FROM (pg_type t
     JOIN pg_enum e ON ((t.oid = e.enumtypid)))
  WHERE (t.typname = 'currency'::name);


SET default_with_oids = false;

--
-- Name: disposals; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE disposals (
    in_transaction_id integer NOT NULL,
    out_transaction_id integer NOT NULL,
    disposed_shares numeric NOT NULL,
    disposed boolean NOT NULL
);


--
-- Name: financial_operation_types; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW financial_operation_types AS
 SELECT e.enumlabel AS type
   FROM (pg_type t
     JOIN pg_enum e ON ((t.oid = e.enumtypid)))
  WHERE (t.typname = 'financing_operation'::name);


--
-- Name: latest_quotes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE latest_quotes (
    ticker character(12) NOT NULL,
    date date NOT NULL,
    open numeric,
    high numeric,
    low numeric,
    close numeric,
    volume bigint,
    openint bigint
);


--
-- Name: securities; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE securities (
    ticker character(12) NOT NULL,
    short_name character varying(128) NOT NULL,
    full_name character varying(128) NOT NULL,
    market character varying(8) NOT NULL,
    leverage numeric DEFAULT 1 NOT NULL,
    quotes_source quotes_source NOT NULL,
    ticker_bankier character(12)
);


--
-- Name: transactions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE transactions (
    transaction_id integer NOT NULL,
    portfolio_id integer NOT NULL,
    date date NOT NULL,
    ticker character(12) NOT NULL,
    price numeric NOT NULL,
    currency currency NOT NULL,
    shares numeric NOT NULL,
    commision numeric NOT NULL,
    exchange_rate numeric NOT NULL,
    tax numeric DEFAULT 0 NOT NULL,
    price_alt numeric,
    CONSTRAINT transactions_commision_check CHECK ((commision >= (0)::numeric)),
    CONSTRAINT transactions_exchange_rate_check CHECK ((exchange_rate > (0)::numeric)),
    CONSTRAINT transactions_price_check CHECK ((price > (0)::numeric)),
    CONSTRAINT transactions_shares_check CHECK ((shares <> (0)::numeric))
);


--
-- Name: most_profitable; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW most_profitable AS
 SELECT
        CASE
            WHEN ("position"((tin.ticker)::text, 'F'::text) = 1) THEN ((regexp_matches((tin.ticker)::text, '^(F.{3})[UZM]{1}[0-9]{2}$'::text))::text)::bpchar
            ELSE
            CASE
                WHEN ("position"((tin.ticker)::text, 'INT'::text) = 1) THEN ((regexp_matches((tin.ticker)::text, '^INT(.{4})[0-9]{4,5}$'::text))::text)::bpchar
                ELSE tin.ticker
            END
        END AS ticker,
    sum(((((d.disposed_shares * tout.price) * tout.exchange_rate) * COALESCE(s.leverage, (1)::numeric)) - (((d.disposed_shares * tin.price) * tin.exchange_rate) * COALESCE(s.leverage, (1)::numeric)))) AS value
   FROM (((disposals d
     JOIN transactions tin ON ((tin.transaction_id = d.in_transaction_id)))
     JOIN transactions tout ON ((tout.transaction_id = d.out_transaction_id)))
     LEFT JOIN securities s ON ((s.ticker = tin.ticker)))
  WHERE d.disposed
  GROUP BY
        CASE
            WHEN ("position"((tin.ticker)::text, 'F'::text) = 1) THEN ((regexp_matches((tin.ticker)::text, '^(F.{3})[UZM]{1}[0-9]{2}$'::text))::text)::bpchar
            ELSE
            CASE
                WHEN ("position"((tin.ticker)::text, 'INT'::text) = 1) THEN ((regexp_matches((tin.ticker)::text, '^INT(.{4})[0-9]{4,5}$'::text))::text)::bpchar
                ELSE tin.ticker
            END
        END
  ORDER BY sum(((((d.disposed_shares * tout.price) * tout.exchange_rate) * COALESCE(s.leverage, (1)::numeric)) - (((d.disposed_shares * tin.price) * tin.exchange_rate) * COALESCE(s.leverage, (1)::numeric)))) DESC;


--
-- Name: operations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE operations (
    operation_id integer NOT NULL,
    portfolio_id integer NOT NULL,
    date date NOT NULL,
    type financing_operation NOT NULL,
    value numeric,
    description character varying(128),
    commision numeric DEFAULT 0 NOT NULL,
    tax numeric DEFAULT 0 NOT NULL
);


--
-- Name: operations_operation_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE operations_operation_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: operations_operation_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE operations_operation_id_seq OWNED BY operations.operation_id;


--
-- Name: owned_stocks; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE owned_stocks (
    portfolio_id integer,
    portfolio_name character varying(128),
    ticker character(12),
    short_name character varying(128),
    market character varying(8),
    shares numeric,
    last_price numeric,
    currency currency,
    exchange_rate numeric,
    last_price_base_currency numeric,
    average_price numeric,
    average_price_base_currency numeric,
    leverage numeric,
    gain numeric,
    percentage_gain numeric,
    gain_base_currency numeric,
    percentage_gain_base_currency numeric,
    market_value numeric,
    market_value_base_currency numeric,
    investment_base_currency numeric,
    average_price_adjusted numeric,
    gain_adjusted numeric,
    percentage_gain_adjusted numeric
);

ALTER TABLE ONLY owned_stocks REPLICA IDENTITY NOTHING;


--
-- Name: portfolios; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE portfolios (
    portfolio_id integer NOT NULL,
    name character varying(128) NOT NULL,
    currency currency NOT NULL
);


--
-- Name: portfolios_ext; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW portfolios_ext AS
 WITH cache AS (
         SELECT o.portfolio_id,
            round(sum((o.value * (
                CASE
                    WHEN (o.type = 'withdraw'::financing_operation) THEN (-1)
                    ELSE 1
                END)::numeric)), 2) AS value,
            sum(o.commision) AS commision
           FROM operations o
          GROUP BY o.portfolio_id
        ), dividend AS (
         SELECT o.portfolio_id,
            sum(o.value) AS value,
            sum(o.commision) AS commision,
            sum(o.tax) AS tax
           FROM operations o
          WHERE (o.type = 'dividend'::financing_operation)
          GROUP BY o.portfolio_id
        ), bondinterest AS (
         SELECT o.portfolio_id,
            sum(o.value) AS value,
            sum(o.commision) AS commision,
            sum(o.tax) AS tax
           FROM operations o
          WHERE (o.type = 'bond interest'::financing_operation)
          GROUP BY o.portfolio_id
        ), expenditure AS (
         SELECT t.portfolio_id,
            sum((((t.shares * t.price) * t.exchange_rate) * COALESCE(s.leverage, (1)::numeric))) AS value,
            sum(t.commision) AS commision,
            sum(t.tax) AS tax
           FROM (transactions t
             LEFT JOIN securities s ON ((s.ticker = t.ticker)))
          GROUP BY t.portfolio_id
        ), gain_of_sold_shares AS (
         SELECT tin.portfolio_id,
            sum(((((d_1.disposed_shares * tout.price) * tout.exchange_rate) * COALESCE(s.leverage, (1)::numeric)) - (((d_1.disposed_shares * tin.price) * tin.exchange_rate) * COALESCE(s.leverage, (1)::numeric)))) AS value
           FROM (((disposals d_1
             JOIN transactions tin ON ((tin.transaction_id = d_1.in_transaction_id)))
             JOIN transactions tout ON ((tout.transaction_id = d_1.out_transaction_id)))
             LEFT JOIN securities s ON ((s.ticker = tin.ticker)))
          WHERE d_1.disposed
          GROUP BY tin.portfolio_id
        ), owned_shares_summary AS (
         SELECT s.portfolio_id,
            sum(s.gain_base_currency) AS gain_of_owned_shares,
            sum(s.market_value_base_currency) AS market_value_base_currency
           FROM owned_stocks s
          GROUP BY s.portfolio_id
        ), annual_balance AS (
         WITH annual_disposals AS (
                 SELECT d_1.in_transaction_id,
                    d_1.out_transaction_id,
                    d_1.disposed_shares,
                    d_1.disposed
                   FROM ((disposals d_1
                     JOIN transactions tout ON ((d_1.out_transaction_id = tout.transaction_id)))
                     JOIN annual_range ar_1 ON (((ar_1.range @> tout.date) AND (ar_1.year = (date_part('year'::text, ('now'::text)::date))::smallint))))
                ), annual_transactions AS (
                 SELECT annual_disposals.in_transaction_id AS transaction_id
                   FROM annual_disposals
                UNION
                 SELECT annual_disposals.out_transaction_id AS transaction_id
                   FROM annual_disposals
                ), annual_costs AS (
                 SELECT t.portfolio_id,
                    sum(t.commision) AS commision,
                    sum(t.tax) AS tax
                   FROM (transactions t
                     JOIN annual_transactions at ON ((at.transaction_id = t.transaction_id)))
                  GROUP BY t.portfolio_id
                ), annual_revenue AS (
                 SELECT tin.portfolio_id,
                    sum((((d_1.disposed_shares * tin.price) * tin.exchange_rate) * COALESCE(s.leverage, (1)::numeric))) AS outcome,
                    sum((((d_1.disposed_shares * tout.price) * tout.exchange_rate) * COALESCE(s.leverage, (1)::numeric))) AS income
                   FROM (((annual_disposals d_1
                     JOIN transactions tout ON ((d_1.out_transaction_id = tout.transaction_id)))
                     JOIN transactions tin ON ((d_1.in_transaction_id = tin.transaction_id)))
                     LEFT JOIN securities s ON ((s.ticker = tin.ticker)))
                  WHERE d_1.disposed
                  GROUP BY tin.portfolio_id
                )
         SELECT ar.portfolio_id,
            ((ar.outcome + ac.tax) + ac.commision) AS outcome,
            ar.income
           FROM (annual_revenue ar
             JOIN annual_costs ac ON ((ac.portfolio_id = ar.portfolio_id)))
        ), month_balance AS (
         WITH month_disposals AS (
                 SELECT d_1.in_transaction_id,
                    d_1.out_transaction_id,
                    d_1.disposed_shares,
                    d_1.disposed
                   FROM (disposals d_1
                     JOIN transactions tout ON ((((d_1.out_transaction_id = tout.transaction_id) AND (date_part('year'::text, tout.date) = date_part('year'::text, ('now'::text)::date))) AND (date_part('month'::text, tout.date) = date_part('month'::text, ('now'::text)::date)))))
                ), month_transactions AS (
                 SELECT month_disposals.in_transaction_id AS transaction_id
                   FROM month_disposals
                UNION
                 SELECT month_disposals.out_transaction_id AS transaction_id
                   FROM month_disposals
                ), month_costs AS (
                 SELECT t.portfolio_id,
                    sum(t.commision) AS commision,
                    sum(t.tax) AS tax
                   FROM (transactions t
                     JOIN month_transactions mt ON ((mt.transaction_id = t.transaction_id)))
                  GROUP BY t.portfolio_id
                ), month_revenue AS (
                 SELECT tin.portfolio_id,
                    sum((((d_1.disposed_shares * tin.price) * tin.exchange_rate) * COALESCE(s.leverage, (1)::numeric))) AS outcome,
                    sum((((d_1.disposed_shares * tout.price) * tout.exchange_rate) * COALESCE(s.leverage, (1)::numeric))) AS income
                   FROM (((month_disposals d_1
                     JOIN transactions tout ON ((d_1.out_transaction_id = tout.transaction_id)))
                     JOIN transactions tin ON ((d_1.in_transaction_id = tin.transaction_id)))
                     LEFT JOIN securities s ON ((s.ticker = tin.ticker)))
                  WHERE d_1.disposed
                  GROUP BY tin.portfolio_id
                )
         SELECT mr.portfolio_id,
            ((mr.outcome + mc.tax) + mc.commision) AS outcome,
            mr.income
           FROM (month_revenue mr
             JOIN month_costs mc ON ((mc.portfolio_id = mr.portfolio_id)))
        )
 SELECT p.portfolio_id,
    p.name,
    p.currency,
    round(((((c.value - c.commision) - COALESCE(e.value, (0)::numeric)) - COALESCE(e.commision, (0)::numeric)) - COALESCE(e.tax, (0)::numeric)), 2) AS cache_value,
    round(gss.value, 2) AS gain_of_sold_shares,
    e.commision,
    e.tax,
    round(oss.gain_of_owned_shares, 2) AS gain_of_owned_shares,
    round((((COALESCE(gss.value, (0)::numeric) + COALESCE(oss.gain_of_owned_shares, (0)::numeric)) + COALESCE(d.value, (0)::numeric)) + COALESCE(bi.value, (0)::numeric)), 2) AS estimated_gain,
    round((((((((((COALESCE(gss.value, (0)::numeric) + COALESCE(oss.gain_of_owned_shares, (0)::numeric)) + COALESCE(d.value, (0)::numeric)) + COALESCE(bi.value, (0)::numeric)) - e.commision) - e.tax) - COALESCE(d.commision, (0)::numeric)) - COALESCE(d.tax, (0)::numeric)) - COALESCE(bi.commision, (0)::numeric)) - COALESCE(bi.tax, (0)::numeric)), 2) AS estimated_gain_costs_inc,
    round((COALESCE(((((c.value - c.commision) - e.value) - e.commision) - e.tax), (0)::numeric) + COALESCE(oss.market_value_base_currency, (0)::numeric)), 2) AS estimated_value,
    COALESCE(round((ab.income - ab.outcome), 2), (0)::numeric) AS annual_balance,
    COALESCE(round((mb.income - mb.outcome), 2), (0)::numeric) AS month_balance
   FROM ((((((((portfolios p
     LEFT JOIN cache c ON ((c.portfolio_id = p.portfolio_id)))
     LEFT JOIN dividend d ON ((d.portfolio_id = p.portfolio_id)))
     LEFT JOIN bondinterest bi ON ((bi.portfolio_id = p.portfolio_id)))
     LEFT JOIN expenditure e ON ((e.portfolio_id = p.portfolio_id)))
     LEFT JOIN gain_of_sold_shares gss ON ((gss.portfolio_id = p.portfolio_id)))
     LEFT JOIN owned_shares_summary oss ON ((oss.portfolio_id = p.portfolio_id)))
     LEFT JOIN annual_balance ab ON ((ab.portfolio_id = p.portfolio_id)))
     LEFT JOIN month_balance mb ON ((mb.portfolio_id = p.portfolio_id)));


--
-- Name: portfolios_portfolio_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE portfolios_portfolio_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: portfolios_portfolio_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE portfolios_portfolio_id_seq OWNED BY portfolios.portfolio_id;


--
-- Name: price_adjustments; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE price_adjustments (
    price_adjustment_id integer NOT NULL,
    transaction_id integer NOT NULL,
    date date NOT NULL,
    adjustment numeric NOT NULL
);


--
-- Name: price_adjustments_price_adjustment_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE price_adjustments_price_adjustment_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: price_adjustments_price_adjustment_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE price_adjustments_price_adjustment_id_seq OWNED BY price_adjustments.price_adjustment_id;


--
-- Name: quotes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE quotes (
    ticker character(12) NOT NULL,
    date date NOT NULL,
    open numeric,
    high numeric,
    low numeric,
    close numeric,
    volume bigint,
    openint bigint
);


--
-- Name: transactions_transaction_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE transactions_transaction_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: transactions_transaction_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE transactions_transaction_id_seq OWNED BY transactions.transaction_id;


--
-- Name: operation_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY operations ALTER COLUMN operation_id SET DEFAULT nextval('operations_operation_id_seq'::regclass);


--
-- Name: portfolio_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY portfolios ALTER COLUMN portfolio_id SET DEFAULT nextval('portfolios_portfolio_id_seq'::regclass);


--
-- Name: price_adjustment_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY price_adjustments ALTER COLUMN price_adjustment_id SET DEFAULT nextval('price_adjustments_price_adjustment_id_seq'::regclass);


--
-- Name: transaction_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY transactions ALTER COLUMN transaction_id SET DEFAULT nextval('transactions_transaction_id_seq'::regclass);


--
-- Name: annual_range_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY annual_range
    ADD CONSTRAINT annual_range_pkey PRIMARY KEY (year);


--
-- Name: disposals_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY disposals
    ADD CONSTRAINT disposals_pkey PRIMARY KEY (in_transaction_id, out_transaction_id);


--
-- Name: latest_quotes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY latest_quotes
    ADD CONSTRAINT latest_quotes_pkey PRIMARY KEY (ticker);


--
-- Name: operations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY operations
    ADD CONSTRAINT operations_pkey PRIMARY KEY (operation_id);


--
-- Name: portfolios_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY portfolios
    ADD CONSTRAINT portfolios_pkey PRIMARY KEY (portfolio_id);


--
-- Name: price_adjustments_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY price_adjustments
    ADD CONSTRAINT price_adjustments_pkey PRIMARY KEY (price_adjustment_id);


--
-- Name: quotes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY quotes
    ADD CONSTRAINT quotes_pkey PRIMARY KEY (ticker, date);


--
-- Name: securities_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY securities
    ADD CONSTRAINT securities_pkey PRIMARY KEY (ticker);


--
-- Name: transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY transactions
    ADD CONSTRAINT transactions_pkey PRIMARY KEY (transaction_id);


--
-- Name: _RETURN; Type: RULE; Schema: public; Owner: -
--

CREATE RULE "_RETURN" AS
    ON SELECT TO owned_stocks DO INSTEAD  WITH remaining_shares AS (
         SELECT t.transaction_id,
            (t.shares - sum(COALESCE(d.disposed_shares, (0)::numeric))) AS shares,
            (t.price - sum(COALESCE(pa.adjustment, (0)::numeric))) AS price_adjusted
           FROM ((transactions t
             LEFT JOIN disposals d ON ((t.transaction_id = d.in_transaction_id)))
             LEFT JOIN ( SELECT price_adjustments.transaction_id,
                    sum(price_adjustments.adjustment) AS adjustment
                   FROM price_adjustments
                  GROUP BY price_adjustments.transaction_id) pa ON ((pa.transaction_id = t.transaction_id)))
          GROUP BY t.transaction_id, t.shares
         HAVING ((t.shares - sum(COALESCE(d.disposed_shares, (0)::numeric))) <> (0)::numeric)
        ), owned_shares AS (
         SELECT t.portfolio_id,
            t.ticker,
            t.currency,
            sum(rs.shares) AS shares,
            sum(((rs.shares * t.price) * COALESCE(s_1.leverage, (1)::numeric))) AS expenditure,
            sum((((rs.shares * t.price) * t.exchange_rate) * COALESCE(s_1.leverage, (1)::numeric))) AS expenditure_base_currency,
            (sum((rs.shares * t.price)) / sum(rs.shares)) AS average_price,
            (sum(((rs.shares * t.price) * t.exchange_rate)) / sum(rs.shares)) AS average_price_base_currency,
            first(t.price, 1 ORDER BY t.date DESC, t.transaction_id DESC) AS last_purchase_price,
            (sum((rs.shares * rs.price_adjusted)) / sum(rs.shares)) AS average_price_adjusted,
            sum(((rs.shares * rs.price_adjusted) * COALESCE(s_1.leverage, (1)::numeric))) AS expenditure_adjusted
           FROM ((transactions t
             JOIN remaining_shares rs ON ((rs.transaction_id = t.transaction_id)))
             LEFT JOIN securities s_1 ON ((s_1.ticker = t.ticker)))
          GROUP BY t.portfolio_id, t.ticker, t.currency
        )
 SELECT p.portfolio_id,
    p.name AS portfolio_name,
    os.ticker,
    s.short_name,
    s.market,
    os.shares,
    q.close AS last_price,
    os.currency,
        CASE
            WHEN (os.currency = p.currency) THEN (1)::numeric
            ELSE e.close
        END AS exchange_rate,
    (q.close *
        CASE
            WHEN (os.currency = p.currency) THEN (1)::numeric
            ELSE e.close
        END) AS last_price_base_currency,
    round(os.average_price, 2) AS average_price,
    round(os.average_price_base_currency, 2) AS average_price_base_currency,
    COALESCE(s.leverage, (1)::numeric) AS leverage,
    round((((os.shares * COALESCE(q.close, os.last_purchase_price[1])) * COALESCE(s.leverage, (1)::numeric)) - os.expenditure), 2) AS gain,
    round((((((os.shares * COALESCE(q.close, os.last_purchase_price[1])) * COALESCE(s.leverage, (1)::numeric)) - os.expenditure) / abs(os.expenditure)) * (100)::numeric), 2) AS percentage_gain,
    round(((((os.shares * COALESCE(q.close, os.last_purchase_price[1])) * COALESCE(s.leverage, (1)::numeric)) *
        CASE
            WHEN (os.currency = p.currency) THEN (1)::numeric
            ELSE e.close
        END) - os.expenditure_base_currency), 2) AS gain_base_currency,
    round(((((((os.shares * COALESCE(q.close, os.last_purchase_price[1])) * COALESCE(s.leverage, (1)::numeric)) *
        CASE
            WHEN (os.currency = p.currency) THEN (1)::numeric
            ELSE e.close
        END) - os.expenditure_base_currency) / abs(os.expenditure_base_currency)) * (100)::numeric), 2) AS percentage_gain_base_currency,
    round(((os.shares * COALESCE(q.close, os.last_purchase_price[1])) * COALESCE(s.leverage, (1)::numeric)), 2) AS market_value,
    round((((os.shares * COALESCE(q.close, os.last_purchase_price[1])) * COALESCE(s.leverage, (1)::numeric)) *
        CASE
            WHEN (os.currency = p.currency) THEN (1)::numeric
            ELSE e.close
        END), 2) AS market_value_base_currency,
    round(os.expenditure_base_currency, 2) AS investment_base_currency,
    round(os.average_price_adjusted, 2) AS average_price_adjusted,
    round((((os.shares * COALESCE(q.close, os.last_purchase_price[1])) * COALESCE(s.leverage, (1)::numeric)) - os.expenditure_adjusted), 2) AS gain_adjusted,
    round((((((os.shares * COALESCE(q.close, os.last_purchase_price[1])) * COALESCE(s.leverage, (1)::numeric)) - os.expenditure_adjusted) / abs(os.expenditure_adjusted)) * (100)::numeric), 2) AS percentage_gain_adjusted
   FROM ((((owned_shares os
     JOIN portfolios p ON ((os.portfolio_id = p.portfolio_id)))
     LEFT JOIN securities s ON ((os.ticker = s.ticker)))
     LEFT JOIN latest_quotes q ON ((os.ticker = q.ticker)))
     LEFT JOIN latest_quotes e ON ((((e.ticker)::text = ((os.currency)::text || (p.currency)::text)) AND (os.currency <> p.currency))))
  ORDER BY p.portfolio_id;


--
-- Name: latest_quotes_before_insert_trigger; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER latest_quotes_before_insert_trigger BEFORE INSERT ON latest_quotes FOR EACH ROW EXECUTE PROCEDURE latest_quotes_before_insert();


--
-- Name: quotes_before_insert_trigger; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER quotes_before_insert_trigger BEFORE INSERT ON quotes FOR EACH ROW EXECUTE PROCEDURE quotes_before_insert();


--
-- Name: transactions_after_insert; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER transactions_after_insert AFTER INSERT ON transactions FOR EACH ROW EXECUTE PROCEDURE transactions_after_insert();


--
-- Name: disposals_in_transaction_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY disposals
    ADD CONSTRAINT disposals_in_transaction_id_fkey FOREIGN KEY (in_transaction_id) REFERENCES transactions(transaction_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: disposals_out_transaction_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY disposals
    ADD CONSTRAINT disposals_out_transaction_id_fkey FOREIGN KEY (out_transaction_id) REFERENCES transactions(transaction_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: operations_portfolio_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY operations
    ADD CONSTRAINT operations_portfolio_id_fkey FOREIGN KEY (portfolio_id) REFERENCES portfolios(portfolio_id);


--
-- Name: price_adjustments_transaction_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY price_adjustments
    ADD CONSTRAINT price_adjustments_transaction_id_fkey FOREIGN KEY (transaction_id) REFERENCES transactions(transaction_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: transactions_portfolio_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY transactions
    ADD CONSTRAINT transactions_portfolio_id_fkey FOREIGN KEY (portfolio_id) REFERENCES portfolios(portfolio_id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- PostgreSQL database dump complete
--

