package cassandra

import (
	"encoding/base64"
	"github.com/amlun/linda/linda/core"
	"github.com/amlun/linda/linda/saver"
	"github.com/gocql/gocql"
	neturl "net/url"
	"strings"
	"time"
)

type Saver struct {
	cassandraURL *neturl.URL
	session      *gocql.Session
}

const PAGE_SIZE = 10

// connect cassandra with urlString cassandra://127.0.0.1:9042/keyspace
func (s *Saver) Connect(url *neturl.URL) error {
	s.cassandraURL = url
	cluster := gocql.NewCluster(url.Host)
	url.Path = strings.TrimPrefix(url.Path, "/")
	cluster.Keyspace = url.Path
	user := url.User
	if user != nil {
		if password, auth := user.Password(); auth {
			cluster.Authenticator = gocql.PasswordAuthenticator{
				Username: user.Username(),
				Password: password,
			}
		}
	}
	session, err := cluster.CreateSession()
	if err != nil {
		return err
	}
	s.session = session
	return nil
}

func (s *Saver) Close() error {
	s.session.Close()
	return nil
}

func (s *Saver) PublishTask(t *core.Task) error {
	if err := s.session.Query(`INSERT INTO tasks (task_id, args, create_time, period, func) VALUES (?, ?, ?, ?, ?)`,
		t.TaskId, t.Args, time.Now(), t.Period, t.Func).Exec(); err != nil {
		return err
	}
	batch := s.session.NewBatch(gocql.CounterBatch)
	if t.Period > 0 {
		batch.Query(`UPDATE periods SET count = count + 1 WHERE period = ?`, t.Period)
	}
	batch.Query(`UPDATE queues SET count = count + 1 WHERE queue = ?`, t.Func)
	if err := s.session.ExecuteBatch(batch); err != nil {
		return err
	}
	return nil
}

func (s *Saver) PublishJob(t *core.Job) error {
	if err := s.session.Query(`INSERT INTO jobs (job_id, args, func, run_time, status, task_id) VALUES (?, ?, ?, ?, ?, ?)`,
		t.JobId, t.Args, t.Func, t.RunTime, t.Status, t.TaskId).Exec(); err != nil {
		return err
	}
	return nil
}

func (s *Saver) Periods() []int {
	var periodList []int
	var period int
	iter := s.session.Query(`SELECT period FROM periods`).Iter()
	for iter.Scan(&period) {
		periodList = append(periodList, period)
	}
	iter.Close()
	return periodList
}

func (s *Saver) Queues() []string {
	var queueList []string
	var queue string
	iter := s.session.Query(`SELECT queue FROM queues`).Iter()
	for iter.Scan(&queue) {
		queueList = append(queueList, queue)
	}
	iter.Close()
	return queueList
}

func (s *Saver) GetPeriodicTask(period int, tasks chan core.Task) {
	var task core.Task
	iter := s.session.Query(`SELECT task_id, args, period, func FROM tasks WHERE period = ?`, period).Iter()
	for iter.Scan(&task.TaskId, &task.Args, &task.Period, &task.Func) {
		tasks <- task
	}
	close(tasks)
}

func (s *Saver) ScheduleTask(id string) error {
	if err := s.session.Query(`INSERT INTO schedules (task_id, schedule_time) VALUES (?, ?)`,
		id, time.Now()).Exec(); err != nil {
		return err
	}
	return nil
}

func (s *Saver) TaskList(taskList *core.TaskList) error {
	var task core.Task
	var i int
	var tasks []core.Task
	stateByte, err := base64.URLEncoding.DecodeString(taskList.State)
	if err != nil {
		return err
	}
	iter := s.session.Query(`SELECT task_id, args, period, func FROM tasks`).PageSize(PAGE_SIZE).PageState(stateByte).Iter()
	for iter.Scan(&task.TaskId, &task.Args, &task.Period, &task.Func) {
		tasks = append(tasks, task)
		i++
	}
	iter.Close()
	taskList.Total = iter.NumRows()
	taskList.Tasks = tasks
	taskList.State = base64.URLEncoding.EncodeToString(iter.PageState())
	return nil
}

func init() {
	saver.Register("cassandra", &Saver{})
}
