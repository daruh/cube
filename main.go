package main

import (
	"cube/task"
	"cube/worker"
	"fmt"
	"github.com/docker/docker/client"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"log"
	"strconv"
	"time"
)

func createContainer() (*task.Docker, *task.DockerResult) {
	c := task.Config{
		Name:  "test-container-1",
		Image: "postgres:13",
		Env: []string{
			"POSTGRES_USER=cube",
			"POSTGRES_PASSWORD=secret",
		},
	}

	dc, _ := client.NewClientWithOpts(client.FromEnv)
	d := task.Docker{
		Client: dc,
		Config: c,
	}

	result := d.Run()
	if result.Error != nil {
		fmt.Printf("%v\n", result.Error)
		return nil, nil
	}

	fmt.Printf("Container %s is running with config %v\n", result.ContainerId, c)
	return &d, &result
}

func stopContainer(d *task.Docker, containerId string) *task.DockerResult {
	result := d.Stop(containerId)
	if result.Error != nil {
		fmt.Printf("%v\n", result.Error)
		return nil
	}

	fmt.Printf("Container %s has been stopped and removed\n", containerId)
	return &result
}

func main() {

	//host := os.Getenv("CUBE_HOST")
	//port, _ := strconv.Atoi(os.Getenv("CUBE_PORT"))

	host := "127.0.0.1"
	port, _ := strconv.Atoi("8089")

	fmt.Println("Starting Cube worker")

	w := worker.Worker{
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}
	api := worker.Api{Address: host, Port: port, Worker: &w}
	go runTasks(&w)
	go w.CollectStats()
	api.Start()
}

func runTasks(w *worker.Worker) {
	for {
		if w.Queue.Len() != 0 {
			result := w.RunTask()
			if result.Error != nil {
				log.Printf("Error running task: %v\n", result.Error)
			}
		} else {
			log.Printf("No tasks to process currently.\n")
		}
		log.Println("Sleeping for 10 seconds.")
		time.Sleep(10 * time.Second)
	}

}

//func main() {
//	db := make(map[uuid.UUID]*task.Task)
//	w := worker.Worker{
//		Queue: *queue.New(),
//		Db:    db,
//	}
//
//	t := task.Task{
//		ID:    uuid.New(),
//		Name:  "test-container-1",
//		State: task.Scheduled,
//		Image: "strm/helloworld-http",
//	}
//
//	// first time the worker will see the task
//	fmt.Println("starting task")
//	w.AddTask(t)
//	result := w.RunTask()
//	if result.Error != nil {
//		panic(result.Error)
//	}
//
//	t.ContainerID = result.ContainerId
//
//	fmt.Printf("task %s is running in container %s\n", t.ID, t.ContainerID)
//	fmt.Println("Sleepy time")
//	time.Sleep(time.Second * 30)
//
//	fmt.Printf("stopping task %s\n", t.ID)
//	t.State = task.Completed
//	w.AddTask(t)
//	result = w.RunTask()
//	if result.Error != nil {
//		panic(result.Error)
//	}
//}

//func main() {
//	t := task.Task{
//		ID:     uuid.New(),
//		Name:   "Task-1",
//		State:  task.Pending,
//		Image:  "Image-1",
//		Memory: 1024,
//		Disk:   1,
//	}
//
//	te := task.TaskEvent{
//		ID:        uuid.New(),
//		State:     task.Pending,
//		Timestamp: time.Now(),
//		Task:      t,
//	}
//
//	fmt.Printf("task: %v\n", t)
//	fmt.Printf("task event: %v\n", te)
//
//	w := worker.Worker{
//		Queue: *queue.New(),
//		Db:    make(map[uuid.UUID]task.Task),
//	}
//	fmt.Printf("worker: %v\n", w)
//	w.CollectStats()
//	w.RunTask()
//	w.StartTask()
//	w.StopTask()
//
//	m := manager.Manager{
//		Pending: *queue.New(),
//		TaskDb:  make(map[string][]task.Task),
//		EventDb: make(map[string][]task.TaskEvent),
//		Workers: []string{w.Name},
//	}
//
//	fmt.Printf("manager: %v\n", m)
//	m.SelectWorker()
//	m.UpdateTasks()
//	m.SendWork()
//
//	n := node.Node{
//		Name:   "Node-1",
//		Ip:     "192.168.1.1",
//		Cores:  4,
//		Memory: 1024,
//		Disk:   25,
//		Role:   "worker",
//	}
//
//	fmt.Printf("node: %v\n", n)
//
//	fmt.Printf("create a test container\n")
//	dockerTask, createResult := createContainer()
//
//	time.Sleep(time.Second * 5)
//
//	fmt.Printf("stopping container %s\n", createResult.ContainerId)
//	fmt.Println(dockerTask)
//	_ = stopContainer(dockerTask, createResult.ContainerId)
//}
