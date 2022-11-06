package ctx

import (
	"fmt"
	"testing"
	"time"
)

func Test_time_after(t *testing.T) {
	fmt.Println(time.Now())
	a := time.After(3 * time.Second)
	fmt.Println(<-a)
	fmt.Println(a)
}

/**
select: select 是golang可以处理多个通道之间的机制， select 和 io select一样是多路复用，是用来监听和channel有关的io操作，当io操作发生时
触发相应的动作。select只能用于channel操作，既可用于channel的接受也可用于channel的发送。如果select的多个分支都满足，则会随机选取一个满足条件
的分支进行操作。
*/

func Test_case1(t *testing.T) {
	c1 := make(chan int)
	//c2 := make(chan int)
	select {
	case <-c1:
		fmt.Println("read c1")
	case c1 <- 1:
		fmt.Println("write c1")
	}
}

/**
select 中的deadlock, 这个时候channel都是停止状态，导致select无法监听到io操作，直接就死锁了。
select中永远不可能有满足的条件，就会死锁，select就处于永久阻塞状态。
*/

type Hello struct {
	cancel chan struct{}
}

func (h *Hello) isQuit() bool {
	select {
	case <-h.cancel:
		fmt.Println("cancel")
		return true
	default:
		return false
	}
}

func Test_case2(t *testing.T) {
	s := &Hello{cancel: make(chan struct{})}
	timer := time.After(3 * time.Second)

	for {
		select {
		case <-timer:
			t.Log("cancel")
			close(s.cancel)
		default:
			quit := s.isQuit()
			t.Log("is Quit:", quit)
			<-time.After(1 * time.Second)
		}
	}
}
