# monujo
The investment portfolio

##HowTo Import Quotes?

 1. Go to http://stooq.com/db/h/
 2. Download daily quotes in ASCII format
 3. Unzip downloaded file
 4. Run below for each requested ticker name:
```
ticker=[ticker]; find . -name "${ticker,,}.txt" -exec sed "s/^/${ticker^^},/g" {} \; | psql -h [hostname] -U [username] -d [dbname] -c "COPY public.quotes FROM STDIN (FORMAT 'csv', DELIMITER ',', HEADER)"
```
