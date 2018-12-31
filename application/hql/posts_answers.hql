CREATE EXTERNAL TABLE IF NOT EXISTS honglibu_posts_answers(
    id int,
    comment_count int,
    creation_date string,
    owner_display_name string,
    owner_user_id int,
    parent_id int,
    score int
)
row format delimited
fields terminated by ','
lines terminated by '\n'

STORED AS TEXTFILE
  location '/tmp/honglibu/posts_answers'

tblproperties("skip.header.line.count"="1");

-- orc format table
CREATE EXTERNAL TABLE IF NOT EXISTS honglibu_orc_posts_answers(
    id int,
    comment_count int,
    creation_date timestamp,
    owner_display_name string,
    owner_user_id int,
    parent_id int,
    score int
)

stored as orc;

insert overwrite table honglibu_orc_posts_answers
select 
    id,
    comment_count,
    from_utc_timestamp(date_format(split(creation_date,'[\.]')[0],'yyyy-MM-dd HH:mm:ss'),'UTC') as creation_date,
    owner_display_name,
    owner_user_id,
    parent_id,
    score
from honglibu_posts_answers;
