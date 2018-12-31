CREATE EXTERNAL TABLE IF NOT EXISTS honglibu_posts_users(
    id int,
    display_name string,
    creation_date string,
    reputation int,
    up_votes int,
    down_votes int,
    views int
)
row format delimited
fields terminated by ','
lines terminated by '\n'

STORED AS TEXTFILE
  location '/tmp/honglibu/posts_users'

tblproperties("skip.header.line.count"="1");

-- orc format table
CREATE EXTERNAL TABLE IF NOT EXISTS honglibu_orc_posts_users(
    id int,
    display_name string,
    creation_date timestamp,
    reputation int,
    up_votes int,
    down_votes int,
    views int
)

stored as orc;

insert overwrite table honglibu_orc_posts_users
select 
    id,
    display_name,
    from_utc_timestamp(date_format(split(creation_date,'[\.]')[0],'yyyy-MM-dd HH:mm:ss'),'UTC') as creation_date,
    reputation,
    up_votes,
    down_votes,
    views
from honglibu_posts_users;