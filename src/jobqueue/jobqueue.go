package jobqueue

import (
	"fmt"
	"time"
)

type Job interface {
	Execute() error
}

type ExecJob struct {
	f func() (interface{}, error)
}

func NewExecJob(f func() (interface{}, error)) *ExecJob {
	return &ExecJob{
		f: f,
	}
}

func (ej ExecJob) Execute() error {
	fmt.Println("Execute job")
	_, err := ej.f()
	return err
}

type JobScheduler struct {
	JobQueue chan Job
	Interval time.Duration
}

func NewJobScheduler(interval time.Duration) *JobScheduler {
	return &JobScheduler{
		JobQueue: make(chan Job),
		Interval: interval,
	}
}

func (s *JobScheduler) Start() {
	go func() {

		for {
			select {
			case job := <-s.JobQueue:
				err := job.Execute()
				if err != nil {
					go func() {
						time.Sleep(s.Interval)
						s.JobQueue <- job
					}()
				}

			}
		}
	}()
}
