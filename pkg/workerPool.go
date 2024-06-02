package pkg

type WorkerPool struct {
	tasks    chan Task
	stopChan chan struct{}
}

type Task struct {
	fn   func(args ...interface{})
	args []interface{}
}

// 创建新的协程池
func NewWorkerPool(size int) *WorkerPool {
	pool := &WorkerPool{
		tasks:    make(chan Task, 16384),
		stopChan: make(chan struct{}),
	}

	for i := 0; i < size; i++ {
		go pool.worker()
	}
	return pool
}

// 每个 worker 从任务通道中获取任务并执行
func (p *WorkerPool) worker() {
	for {
		select {
		case task := <-p.tasks:
			{
				func() {
					defer func() {
						recover()
					}()
					task.fn(task.args...)
				}()
			}
		case <-p.stopChan:
			return
		}
	}
}

// 提交任务到协程池
func (p *WorkerPool) Submit(fn func(args ...interface{}), args ...interface{}) {
	p.tasks <- Task{fn: fn, args: args}
}

// 关闭协程池
func (p *WorkerPool) Close() {
	close(p.stopChan)
	close(p.tasks)
}
