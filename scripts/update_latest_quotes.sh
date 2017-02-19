#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
CONF_FILE=$DIR/../conf/db.conf

if [ ! -e $CONF_FILE ] ; then
    echo "Please run install.sh first!"
    exit
fi

source $CONF_FILE

if [ -z $DB_HOSTNAME ] || [ -z $DB_USERNAME ] || [ -z $DB_PASSWORD ] || [ -z $DB_NAME ] || [ -z $QUOTES_PATH ] ; then
    echo -e "Config file seems to be not filled properly!\nSome variables have not been set."
    exit
fi

TICKERS_FILE=$DIR/../conf/tickers.conf

if [ ! -e $TICKERS_FILE ] ; then
    echo "Please create tickers config file"
    exit
fi

if [ ! -e $1 ] ; then
    echo "Invalid file path provided"
    exit
fi

mapfile -t TICKERS < "$TICKERS_FILE"

for TICKER in ${TICKERS[@]} ; do
    QUOTES="$( sed -n "s/,D,/,/g;/^${TICKER^^},/p" $1 | tail -2 )"
    echo $QUOTES
    RESULT="$( psql -h $DB_HOSTNAME -U $DB_USERNAME -d $DB_NAME -c "COPY public.latest_quotes FROM STDIN (FORMAT 'csv', DELIMITER ',', HEADER)" <<< "$QUOTES" )"
    echo "${TICKER^^} ${RESULT/COPY/""}"
done

