package serv

type Task func()

type QueueOption func(*Queue)

func WithSize(size int) QueueOption {
	return func(d *Queue) {
		d.size = size
	}
}

func WithCapacity(capacity int) QueueOption {
	return func(d *Queue) {
		d.capacity = capacity
	}
}

type Queue struct {
	ch       chan Task
	size     int
	capacity int
}

func NewQueue(opts ...QueueOption) *Queue {
	q := &Queue{capacity: 10, size: 1}
	for _, o := range opts {
		o(q)
	}
	q.ch = make(chan Task, q.capacity)
	for i := 0; i < q.size; i++ {
		go func() {
			for {
				select {
				case fn := <-q.ch:
					fn()
				}
			}
		}()
	}
	return q
}

func (my *Queue) Add(t Task) {
	go func() {
		my.ch <- t
	}()
}
