create external table honglibu_hbase_active_users (
    id string,
    display_name string,
    reputation bigint,
    answer_num bigint,
    votes bigint)
STORED BY 'org.apache.hadoop.hive.hbase.HBaseStorageHandler'
WITH SERDEPROPERTIES ('hbase.columns.mapping' = ':key, users:display_name, users:reputation, users:answer_num, users:votes#b')
TBLPROPERTIES ('hbase.table.name' = 'honglibu_hbase_active_users');

insert overwrite table honglibu_hbase_active_users
select
    id,
    display_name,
    reputation,
    answer_num,
    0 AS votes
 from honglibu_active_users
 limit 10;