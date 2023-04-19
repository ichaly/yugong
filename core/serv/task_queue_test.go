package serv

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

var wg sync.WaitGroup

type DownloadFile struct {
	Url  string
	File string
}

func (my *DownloadFile) process(worker string) {
	s := rand.Intn(1000)
	time.Sleep(time.Millisecond * time.Duration(s))
	my.File = fmt.Sprintf("%s is downloaded %s use %dms", worker, my.Url, s)
	wg.Done()
}

func TestTaskQueue(t *testing.T) {
	q := NewTaskQueue()

	for i := 0; i < 50; i++ {
		wg.Add(1)
		d := DownloadFile{
			Url: fmt.Sprintf("http://example.com/%d", i),
		}
		q.Add(&d)
		println(d.File)
	}

	wg.Wait()
}
