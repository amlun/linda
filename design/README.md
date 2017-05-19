# Design

## System Design

![system-design](https://rawgit.com/amlun/linda/master/design/linda-system-design.jpg)

## Model Design

![model-design](https://rawgit.com/amlun/linda/master/design/linda-model-design.jpg)

## Saver Design

### Table Tasks

| field | type | description |
|:-------|:--------|:--------|
|task_id|text|the identify of the task|
|args|list<text>|runtime args of the task|
|create_time|time|the task created time|
|func|text|the task runtime func name|
|period|int|scheduled time of seconds|
|queue|text|job will enqueue in this queue name|


### Table Jobs

| field | type | description |
|:-------|:--------|:--------|
|job_id|text|the identify of the job|
|args|list<text>|runtime args of the job|
|func|text|the job runtime func name|
|queue|text|job queue name|
|run_time|time|the job created time|
|status|int|the job status|
|task_id|text|source of the task id|
