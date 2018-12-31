'use strict';
const http = require('http');
var assert = require('assert');
const express= require('express');
const app = express();
const mustache = require('mustache');
const filesystem = require('fs');
const url = require('url');
const hbase = require('hbase-rpc-client');
const hostname = '127.0.0.1';
const port = 3192;
var bodyParser = require('body-parser');
/* Commented out lines are for running on our cluster */
var client = hbase({
    // zookeeperHosts: ["localhost:2181"],
    zookeeperHosts: ["10.0.0.2:2181"],
    zookeeperRoot: "/hbase-unsecure"
});

client.on('error', function(err) {
  console.log(err)
})

console.log("foo");
app.use(express.static('public'));
app.use(bodyParser.urlencoded({extended: false }));
app.use(bodyParser.json());
app.get('/users.html',function (req, res) {
    const ranking=req.query.ranking;
    const get = new hbase.Get(ranking);
    console.log(client.get("honglibu_hbase_excellent_users", get, function(err, row) {
        assert.ok(!err, "get returned an error: #{err}");
        if(!row){
            res.send("<html><body>No such data</body></html>");
            return;
        }
        var template = filesystem.readFileSync("excellent_users.mustache").toString();
        var html = mustache.render(template,  {
            Ranking : req.query.ranking,
            Name : row.cols["users:display_name"].value,
            Id : row.cols["users:user_id"].value,
            StartFrom : row.cols["users:creation_date"].value,
            Reputation : row.cols["users:reputation"].value,
            Upvotes : row.cols["users:up_votes"].value,
            Downvotes : row.cols["users:down_votes"].value,
            Views : row.cols["users:view"].value
        });
        res.send(html);
    }));
});

app.get('/questions.html',function (req, res) {
    const ranking= req.query.ranking;
    const get = new hbase.Get(ranking);
    console.log(client.get("honglibu_hbase_best_questions_"+req.query.year, get, function(err, row) {
        assert.ok(!err, "get returned an error: #{err}");
        if(!row){
            res.send("<html><body>No such data</body></html>");
            return;
        }
        var template = filesystem.readFileSync("best_questions.mustache").toString();
        var html = mustache.render(template,  {
            Ranking : req.query.ranking,
            Id : row.cols["questions:id"].value,
            Title : row.cols["questions:title"].value,
            AnswerCount : row.cols["questions:answer_count"].value,
            CommentCount : row.cols["questions:comment_count"].value,
            CreationDate : row.cols["questions:creation_date"].value,
            FavoriteCount : row.cols["questions:favorite_count"].value,
            OwnerId: row.cols["questions:owner_user_id"].value,
            OwnerName: row.cols["questions:owner_display_name"].value,
            Score: row.cols["questions:score"].value,
            Tags : row.cols["questions:tags"].value,
            ViewCount: row.cols["questions:view_count"].value,
            AcceptedAnswerId : typeof row.cols["questions:accepted_answer_id"] !== "undefined" && row.cols["questions:accepted_answer_id"] !== null?row.cols["questions:accepted_answer_id"].value : null,
            AcceptedAnswerDate: typeof row.cols["questions:accepted_answer_creation_date"] !== "undefined" && row.cols["questions:accepted_answer_creation_date"] !==null ?row.cols["questions:accepted_answer_creation_date"].value : null,
            AcceptedAnswerOwnerId : typeof row.cols["questions:accepted_answer_owner_user_id"] !== "undefined" && row.cols["questions:accepted_answer_owner_user_id"] !== null?row.cols["questions:accepted_answer_owner_user_id"].value :null,
            AcceptedAnswerOwnerName : typeof row.cols["questions:accepted_answer_owner_display_name"] !== "undefined" && row.cols["questions:accepted_answer_owner_display_name"]!==null?row.cols["questions:accepted_answer_owner_display_name"].value:null
        });
        res.send(html);
    }));
});

/* Send simulated weather to kafka */
var kafka = require('kafka-node');
var Producer = kafka.Producer;
var KeyedMessage = kafka.KeyedMessage;
var kafkaClient = new kafka.KafkaClient({kafkaHost: '10.0.0.2:6667'});
var kafkaProducer = new Producer(kafkaClient);

app.get('/active_users.html',function (req, res) {
    var idAndName =req.query.idAndName;
    var userid = idAndName.split("+")[0]
    var name = idAndName.split("+")[1] 
    var report = {
    id : userid,
    display_name : name,
    };

    kafkaProducer.send([{ topic: 'honglibu_active_users', messages: JSON.stringify(report)}],
               function (err, data) {
                   console.log(data);
               });
    res.redirect('submit-vote.html');
});

app.listen(port);