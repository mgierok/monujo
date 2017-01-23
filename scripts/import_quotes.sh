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

if [ -z $1 ] ; then
    mapfile -t TICKERS < $DIR/../conf/tickers.conf
else
    IFS=',' read -r -a TICKERS <<< "$1"
    unset IFS
fi

for TICKER in ${TICKERS[@]} ; do
    QUOTES="$( find $QUOTES_PATH -name "${TICKER,,}.txt" -exec sed "s/^/${TICKER^^},/g" {} \; )"
    RESULT="$( psql -h $DB_HOSTNAME -U $DB_USERNAME -d $DB_NAME -c "COPY public.quotes FROM STDIN (FORMAT 'csv', DELIMITER ',', HEADER)" <<< "$QUOTES" )"
    echo "${TICKER^^} ${RESULT/COPY/""}"
done

