CREATE EXTERNAL TABLE IF NOT EXISTS honglibu_posts_tags(
    id int,
    tag_name string,
    count int,
    excerpt_post_id int,
    wiki_post_id int
)
row format delimited
fields terminated by ','
lines terminated by '\n'

STORED AS TEXTFILE
  location '/tmp/honglibu/posts_tags'

tblproperties("skip.header.line.count"="1");

-- orc format table
CREATE EXTERNAL TABLE IF NOT EXISTS honglibu_orc_posts_tags(
    id int,
    tag_name string,
    count int,
    excerpt_post_id int,
    wiki_post_id int	
)
stored as orc;

insert overwrite table honglibu_orc_posts_tags select * from honglibu_posts_tags;