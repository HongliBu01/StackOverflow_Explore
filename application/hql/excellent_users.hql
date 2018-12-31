create external table if not exists honglibu_excellent_users (
    id int,
    display_name string,
    creation_date timestamp,
    reputation int,
    up_votes int,
    down_votes int,
    views int
) 
stored as orc;

insert overwrite table honglibu_excellent_users
  select *
  from honglibu_orc_posts_users
  order by reputation
  desc;