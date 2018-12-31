DROP TABLE honglibu_posts_questions;
CREATE EXTERNAL TABLE IF NOT EXISTS honglibu_posts_questions(
    id int,
    title string,
    accepted_answer_id int,
    answer_count int,
    comment_count int,
    creation_date string,
    favorite_count int,
    owner_user_id int,
    owner_display_name string,
    score int,
    tags string,
    view_count int
)
row format delimited
fields terminated by ','
lines terminated by '\n'

STORED AS TEXTFILE
  location '/tmp/honglibu/posts_questions'

tblproperties("skip.header.line.count"="1");

-- orc format table
CREATE EXTERNAL TABLE IF NOT EXISTS honglibu_orc_posts_questions(
    id int,
    title string,
    accepted_answer_id int,
    answer_count int,
    comment_count int,
    creation_date timestamp,
    favorite_count int,
    owner_user_id int,
    owner_display_name string,
    score int,
    tags string,
    view_count int
)
stored as orc;

insert overwrite table honglibu_orc_posts_questions
select 
    id,
    title,
    accepted_answer_id,
    answer_count,
    comment_count,
    from_utc_timestamp(date_format(split(creation_date,'[\.]')[0],'yyyy-MM-dd HH:mm:ss'),'UTC') as creation_date,
    favorite_count,
    owner_user_id,
    owner_display_name,
    score,
    tags,
    view_count
from honglibu_posts_questions;