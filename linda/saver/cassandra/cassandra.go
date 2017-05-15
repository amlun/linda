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
	if err := s.session.Query(`INSERT INTO tasks (id, args, create_time, frequency, func) VALUES (?, ?, ?, ?, ?)`,
		t.ID, t.Args, time.Now(), t.Frequency, t.Func).Exec(); err != nil {
		return err
	}
	return nil
}

func (s *Saver) PublishJob(t *core.Job) error {
	if err := s.session.Query(`INSERT INTO jobs (id, args, func, run_time, status, task_id) VALUES (?, ?, ?, ?, ?, ?)`,
		t.ID, t.Args, t.Func, t.RunTime, t.Status, t.TaskId).Exec(); err != nil {
		return err
	}
	return nil
}

func (s *Saver) Frequency(frequency time.Duration) error {
	if err := s.session.Query(`UPDATE frequencies SET count = count + 1 WHERE frequency = ?`,
		frequency).Exec(); err != nil {
		return err
	}
	return nil
}

func (s *Saver) Frequencies() []time.Duration {
	var frequencyList []time.Duration
	var frequency time.Duration
	iter := s.session.Query(`SELECT frequency FROM frequencies`, ).Iter()
	for iter.Scan(&frequency) {
		frequencyList = append(frequencyList, frequency)
	}
	iter.Close()
	return frequencyList
}

func (s *Saver) GetTimingTask(frequency time.Duration, tasks chan core.Task) {
	var task core.Task
	iter := s.session.Query(`SELECT id, args, frequency, func FROM tasks WHERE frequency = ?`, frequency).Iter()
	for iter.Scan(&task.ID, &task.Args, &task.Frequency, &task.Func) {
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
	iter := s.session.Query(`SELECT id, args, frequency, func FROM tasks`).PageSize(PAGE_SIZE).PageState(stateByte).Iter()
	for iter.Scan(&task.ID, &task.Args, &task.Frequency, &task.Func) {
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
