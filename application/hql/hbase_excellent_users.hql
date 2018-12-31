create external table honglibu_hbase_excellent_users (
   ranking bigint,
   user_id bigint,
   display_name string,
   creation_date timestamp,
   reputation bigint,
   up_votes bigint,
   down_votes bigint,
   views bigint)
STORED BY 'org.apache.hadoop.hive.hbase.HBaseStorageHandler'
WITH SERDEPROPERTIES ('hbase.columns.mapping' = ':key,users:user_id,users:display_name,users:creation_date,users:reputation,users:up_votes,users:down_votes,users:view')
TBLPROPERTIES ('hbase.table.name' = 'honglibu_hbase_excellent_users');

insert overwrite table honglibu_hbase_excellent_users
select
   row_number()over(order by reputation desc) AS ranking,
   id,
   display_name,
   creation_date,
   reputation,
   up_votes,
   down_votes,
   views
from honglibu_excellent_users
limit 10000;