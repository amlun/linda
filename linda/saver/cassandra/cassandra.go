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

func (s *Saver) SaveTask(t *core.Task) error {
	if err := s.session.Query(`INSERT INTO tasks (task_id, args, create_time, period, func, queue) VALUES (?, ?, ?, ?, ?, ?) IF NOT EXISTS`,
		t.TaskId, t.Args, time.Now(), t.Period, t.Func, t.Queue).Exec(); err != nil {
		return err
	}
	return nil
}

func (s *Saver) SaveJob(t *core.Job) error {
	if err := s.session.Query(`INSERT INTO jobs (job_id, args, func, run_time, status, task_id, queue) VALUES (?, ?, ?, ?, ?, ?, ?) IF NOT EXISTS`,
		t.JobId, t.Args, t.Func, t.RunTime, t.Status, t.TaskId, t.Queue).Exec(); err != nil {
		return err
	}
	return nil
}

func (s *Saver) GetTask(taskId string) (*core.Task, error) {
	var task core.Task
	if err := s.session.Query(`SELECT task_id, args, period, func, queue FROM tasks WHERE task_id = ?`,
		taskId).Scan(&task.TaskId, &task.Args, &task.Period, &task.Func, &task.Queue); err != nil {
		return nil, err
	}
	return &task, nil
}

func (s *Saver) TaskList(taskList *core.TaskList) error {
	var task core.Task
	var i int
	var tasks []core.Task
	stateByte, err := base64.URLEncoding.DecodeString(taskList.State)
	if err != nil {
		return err
	}
	iter := s.session.Query(`SELECT task_id, args, period, func, queue FROM tasks`).PageSize(PAGE_SIZE).PageState(stateByte).Iter()
	for iter.Scan(&task.TaskId, &task.Args, &task.Period, &task.Func, &task.Queue) {
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
