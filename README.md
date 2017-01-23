# monujo
The investment portfolio

##Install

 1. Run **./install.sh**

##HowTo Import Quotes?

 1. Go to http://stooq.com/db/h/
 2. Download daily quotes in ASCII format
 3. Unzip downloaded file
 4. Run **./scripts/import_quotes.sh** and pass comma separated list of tickers

##HowTo Update Quotes?

 1. Please make sure tickers.conf file is populated
 2. Go to http://stooq.com/db/
 3. Download latest daily quotes
 4. Run **./scripts/update_quotes [path_to_downloaded_file]**
