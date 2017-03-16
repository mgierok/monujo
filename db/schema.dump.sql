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
    'dividend'
);


--
-- Name: transaction_operation_type; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE transaction_operation_type AS ENUM (
    'buy',
    'sell'
);


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
	v_remaining_shares remaining_shares.shares%TYPE;
	v_transaction_id remaining_shares.transaction_id%TYPE;
	
	BEGIN
		IF NEW."type" = 'sell'::transaction_operation_type THEN

			SELECT SUM(t.shares) - SUM(COALESCE(d.disposed_shares, 0))  INTO v_remaining_shares
			FROM transactions t
			LEFT JOIN (SELECT buy_transaction_id, SUM(disposed_shares) disposed_shares FROM disposals GROUP BY buy_transaction_id) d ON d.buy_transaction_id = t.transaction_id
			WHERE t.portfolio_id = NEW."portfolio_id" AND t.ticker = NEW."ticker" AND t.type = 'buy'::transaction_operation_type 
			GROUP BY t.ticker
			HAVING (SUM(t.shares) - SUM(COALESCE(d.disposed_shares, 0))) > 0;

			IF NOT FOUND THEN
				RAISE EXCEPTION 'You have no shares of % in your portfolio!', NEW."ticker";
			END IF;

			IF v_remaining_shares < NEW."shares" THEN
				RAISE EXCEPTION 'You have only % shares of % in your portfolio!', v_remaining_shares, NEW."ticker";
			END IF;

			v_shares := NEW."shares";	
			WHILE v_shares > 0 LOOP
				SELECT t.shares - SUM(COALESCE(d.disposed_shares, 0)), t.transaction_id INTO v_remaining_shares, v_transaction_id
				FROM transactions t
				LEFT JOIN disposals d ON t.transaction_id = d.buy_transaction_id
				WHERE t.portfolio_id = NEW."portfolio_id" AND t.ticker = NEW."ticker" AND t.type = 'buy'::transaction_operation_type
				GROUP BY t.transaction_id, t.shares, t."date"
				HAVING (t.shares - SUM(COALESCE(d.disposed_shares, 0))) > 0
				ORDER BY t."date" ASC, t.transaction_id ASC
				LIMIT 1;

				IF v_remaining_shares >= v_shares THEN
					INSERT INTO disposals (buy_transaction_id, sell_transaction_id, disposed_shares) VALUES (v_transaction_id, NEW."transaction_id", v_shares);
					v_shares := 0;
				ELSE
					INSERT INTO disposals (buy_transaction_id, sell_transaction_id, disposed_shares) VALUES (v_transaction_id, NEW."transaction_id", v_remaining_shares);
					v_shares := v_shares - v_remaining_shares;
				END IF;
			END LOOP;
		END IF;
		RETURN NEW;
	END;
$$;


SET default_with_oids = false;

--
-- Name: disposals; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE disposals (
    buy_transaction_id integer NOT NULL,
    sell_transaction_id integer NOT NULL,
    disposed_shares numeric NOT NULL
);


--
-- Name: latest_quotes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE latest_quotes (
    ticker character(8) NOT NULL,
    date date NOT NULL,
    open numeric,
    high numeric,
    low numeric,
    close numeric,
    volume bigint,
    openint bigint
);


--
-- Name: operations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE operations (
    operation_id integer NOT NULL,
    portfolio_id integer NOT NULL,
    date date NOT NULL,
    type financing_operation NOT NULL,
    value numeric
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
-- Name: portfolios; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE portfolios (
    portfolio_id integer NOT NULL,
    name character varying(128) NOT NULL,
    currency currency NOT NULL
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
    type transaction_operation_type NOT NULL,
    currency currency NOT NULL,
    shares numeric NOT NULL,
    commision numeric NOT NULL,
    exchange_rate numeric NOT NULL,
    tax numeric DEFAULT 0 NOT NULL,
    CONSTRAINT transactions_commision_check CHECK (((commision)::double precision >= (0)::double precision)),
    CONSTRAINT transactions_exchange_rate_check CHECK (((exchange_rate)::double precision > (0)::double precision)),
    CONSTRAINT transactions_price_check CHECK (((price)::double precision > (0)::double precision)),
    CONSTRAINT transactions_shares_check CHECK (((shares)::double precision > (0)::double precision))
);


--
-- Name: remaining_shares; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW remaining_shares AS
 SELECT t.transaction_id,
    (t.shares - sum(COALESCE(d.disposed_shares, (0)::numeric))) AS shares
   FROM (transactions t
     LEFT JOIN disposals d ON ((t.transaction_id = d.buy_transaction_id)))
  WHERE (t.type = 'buy'::transaction_operation_type)
  GROUP BY t.transaction_id, t.shares;


--
-- Name: securities; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE securities (
    ticker character(12) NOT NULL,
    short_name character varying(128) NOT NULL,
    full_name character varying(128) NOT NULL,
    market character varying(8) NOT NULL
);


--
-- Name: shares; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW shares AS
 SELECT t.portfolio_id,
    t.ticker,
    s.short_name,
    s.full_name,
    s.market,
    sum(rs.shares) AS shares,
    q.close AS last_price,
    round((sum(rs.shares) * q.close), 2) AS market_value,
    t.currency,
        CASE
            WHEN (t.currency = p.currency) THEN ((1)::real)::numeric
            ELSE e.close
        END AS exchange_rate,
    round(((sum(rs.shares) * q.close) *
        CASE
            WHEN (t.currency = p.currency) THEN ((1)::real)::numeric
            ELSE e.close
        END)) AS market_value_base_currency,
    round((sum((rs.shares * t.price)) / sum(rs.shares)), 2) AS average_price,
    round(((sum(rs.shares) * q.close) - sum((rs.shares * t.price))), 2) AS gain,
    round(((((sum(rs.shares) * q.close) - sum((rs.shares * t.price))) / sum((rs.shares * t.price))) * (100)::numeric), 2) AS percentage_gain,
    round((((sum(rs.shares) * q.close) *
        CASE
            WHEN (t.currency = p.currency) THEN ((1)::real)::numeric
            ELSE e.close
        END) - sum(((rs.shares * t.price) * t.exchange_rate))), 2) AS gain_base_currency,
    round((((((sum(rs.shares) * q.close) *
        CASE
            WHEN (t.currency = p.currency) THEN ((1)::real)::numeric
            ELSE e.close
        END) - sum(((rs.shares * t.price) * t.exchange_rate))) / sum(((rs.shares * t.price) * t.exchange_rate))) * (100)::numeric), 2) AS percentage_gain_base_currency
   FROM (((((transactions t
     JOIN portfolios p ON ((t.portfolio_id = p.portfolio_id)))
     LEFT JOIN securities s ON ((t.ticker = s.ticker)))
     LEFT JOIN latest_quotes q ON ((t.ticker = q.ticker)))
     LEFT JOIN latest_quotes e ON ((((e.ticker)::text = ((t.currency)::text || (p.currency)::text)) AND (t.currency <> p.currency))))
     LEFT JOIN remaining_shares rs ON ((rs.transaction_id = t.transaction_id)))
  GROUP BY t.portfolio_id, t.ticker, s.short_name, s.full_name, s.market, q.close, t.currency, e.close, p.currency
 HAVING (sum(
        CASE
            WHEN (t.type = 'sell'::transaction_operation_type) THEN (t.shares * ((-1))::numeric)
            ELSE t.shares
        END) > (0)::numeric)
  ORDER BY t.portfolio_id;


--
-- Name: portfolio_summary; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW portfolio_summary AS
 SELECT p.portfolio_id,
    p.name,
    p.currency,
    (fo.cache_value - t.invested_value) AS cache_value,
    "out".outgoings,
    "in".incomings,
    ("in".incomings - "out".outgoings) AS gain_of_sold_shares,
    ("in".commision + "out".commision) AS commision,
    os.gain_of_owned_shares
   FROM (((((portfolios p
     LEFT JOIN ( SELECT o.portfolio_id,
            round(sum((o.value * (
                CASE
                    WHEN (o.type = 'withdraw'::financing_operation) THEN (-1)
                    ELSE 1
                END)::numeric)), 2) AS cache_value
           FROM operations o
          GROUP BY o.portfolio_id) fo ON ((fo.portfolio_id = p.portfolio_id)))
     LEFT JOIN ( SELECT t_1.portfolio_id,
            round(sum((t_1.commision + (((t_1.shares * t_1.price) * t_1.exchange_rate) * (
                CASE
                    WHEN (t_1.type = 'buy'::transaction_operation_type) THEN 1
                    ELSE (-1)
                END)::numeric))), 2) AS invested_value
           FROM transactions t_1
          GROUP BY t_1.portfolio_id) t ON ((t.portfolio_id = p.portfolio_id)))
     LEFT JOIN ( SELECT t_1.portfolio_id,
            round(sum((((t_1.shares - rs.shares) * t_1.price) * t_1.exchange_rate)), 2) AS outgoings,
            sum(t_1.commision) AS commision
           FROM (transactions t_1
             LEFT JOIN remaining_shares rs ON ((t_1.transaction_id = rs.transaction_id)))
          WHERE ((t_1.shares - rs.shares) > (0)::numeric)
          GROUP BY t_1.portfolio_id) "out" ON (("out".portfolio_id = p.portfolio_id)))
     LEFT JOIN ( SELECT t_1.portfolio_id,
            round(sum(((t_1.shares * t_1.price) * t_1.exchange_rate)), 2) AS incomings,
            sum(t_1.commision) AS commision
           FROM transactions t_1
          WHERE (t_1.type = 'sell'::transaction_operation_type)
          GROUP BY t_1.portfolio_id) "in" ON (("in".portfolio_id = p.portfolio_id)))
     LEFT JOIN ( SELECT s.portfolio_id,
            sum(s.gain_base_currency) AS gain_of_owned_shares
           FROM shares s
          GROUP BY s.portfolio_id) os ON ((os.portfolio_id = p.portfolio_id)));


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
-- Name: quotes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE quotes (
    ticker character(8) NOT NULL,
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
-- Name: transaction_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY transactions ALTER COLUMN transaction_id SET DEFAULT nextval('transactions_transaction_id_seq'::regclass);


--
-- Name: disposals_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY disposals
    ADD CONSTRAINT disposals_pkey PRIMARY KEY (buy_transaction_id, sell_transaction_id);


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
-- Name: disposals_sell_transaction_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX disposals_sell_transaction_id_idx ON disposals USING btree (sell_transaction_id);


--
-- Name: latest_quotes_before_insert_trigger; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER latest_quotes_before_insert_trigger BEFORE INSERT ON latest_quotes FOR EACH ROW EXECUTE PROCEDURE latest_quotes_before_insert();


--
-- Name: quotes_before_insert_trigger; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER quotes_before_insert_trigger BEFORE INSERT ON quotes FOR EACH ROW EXECUTE PROCEDURE quotes_before_insert();


--
-- Name: transactions_after_insert_trigger; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER transactions_after_insert_trigger AFTER INSERT ON transactions FOR EACH ROW EXECUTE PROCEDURE transactions_after_insert();


--
-- Name: disposals_buy_transaction_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY disposals
    ADD CONSTRAINT disposals_buy_transaction_id_fkey FOREIGN KEY (buy_transaction_id) REFERENCES transactions(transaction_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: disposals_sell_transaction_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY disposals
    ADD CONSTRAINT disposals_sell_transaction_id_fkey FOREIGN KEY (sell_transaction_id) REFERENCES transactions(transaction_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: operations_portfolio_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY operations
    ADD CONSTRAINT operations_portfolio_id_fkey FOREIGN KEY (portfolio_id) REFERENCES portfolios(portfolio_id);


--
-- Name: transactions_portfolio_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY transactions
    ADD CONSTRAINT transactions_portfolio_id_fkey FOREIGN KEY (portfolio_id) REFERENCES portfolios(portfolio_id);


--
-- PostgreSQL database dump complete
--

