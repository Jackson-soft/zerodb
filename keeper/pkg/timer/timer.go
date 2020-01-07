package timer

import (
	"container/list"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

//Node 时间轮槽的链表节点
type Node struct {
	id      string
	tType   TimerType
	delay   time.Duration          // 延迟时间
	circle  int                    // 时间轮需要转动几圈
	handler func(args interface{}) // 任务函数
	args    interface{}            // 函数参数
}

//TimeWheel 时间轮
type TimeWheel struct {
	slots      []*list.List // 时间轮的槽
	currentPos int          // 当前指针指向哪一个槽
	slotNum    int          // 槽数量
}

//Timer 定时器
type Timer struct {
	mutex  sync.Mutex
	tick   time.Duration //最小粒度
	ticker *time.Ticker

	timeWheel *TimeWheel
	timerMap  sync.Map //存储定时器id对应的槽位置
	stop      chan bool
}

//TimerType 定时器类型
type TimerType uint8

const (
	Single     TimerType = 1 //单次
	Repetition TimerType = 2 //循环

	defaultTick = time.Second
	defaultNum  = 5
)

//NewTimer 创建一个指定槽数量和最小粒度的定时器
func NewTimer(tick time.Duration, num int) *Timer {
	if tick <= 0 {
		tick = defaultTick
	}
	if num <= 0 {
		num = defaultNum
	}

	t := new(Timer)

	t.tick = tick

	t.timeWheel = new(TimeWheel)

	t.timeWheel.slotNum = num

	t.timeWheel.slots = make([]*list.List, t.timeWheel.slotNum)

	t.timerMap = sync.Map{}

	t.timeWheel.currentPos = 0

	t.stop = make(chan bool)

	for i := 0; i < t.timeWheel.slotNum; i++ {
		t.timeWheel.slots[i] = list.New()
	}

	go t.run()

	return t
}

//Register 注册一个定时器，返回timerID
func (t *Timer) Register(tType TimerType, delay time.Duration, handler func(args interface{}), args interface{}) (string, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	uid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	pos, circle := t.getPositionAndCircle(delay)

	node := new(Node)

	node.tType = tType
	node.circle = circle
	node.handler = handler
	node.args = args
	node.delay = delay

	node.id = uid.String()

	t.timeWheel.slots[pos].PushBack(node)

	t.timerMap.Store(node.id, pos)

	return node.id, nil
}

func (t *Timer) registerV1(node *Node) {
	pos, circle := t.getPositionAndCircle(node.delay)

	node.circle = circle

	t.timeWheel.slots[pos].PushBack(node)

	t.timerMap.Store(node.id, pos)
}

//Remove 删除指定定时器
func (t *Timer) Remove(timerID string) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if len(timerID) == 0 {
		return errors.New("id is nil")
	}

	pos, ok := t.timerMap.Load(timerID)
	if !ok {
		return errors.New("id do not exist")
	}

	l := t.timeWheel.slots[pos.(int)]
	for e := l.Front(); e != nil; {
		job := e.Value.(*Node)
		if job.id == timerID {
			t.timerMap.Delete(timerID)
			l.Remove(e)
			break
		}
		e = e.Next()
	}
	return nil
}

//Reset 重新设置定时器
func (t *Timer) Reset(timerID string) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if len(timerID) == 0 {
		return errors.New("id is nil")
	}

	pos, ok := t.timerMap.Load(timerID)
	if !ok {
		return errors.New("id do not exist")
	}

	l := t.timeWheel.slots[pos.(int)]
	for e := l.Front(); e != nil; {
		job := e.Value.(*Node)
		if job.id == timerID {
			l.Remove(e)
			t.registerV1(job)
			break
		}
		e = e.Next()
	}
	return nil
}

// 获取定时器在槽中的位置, 时间轮需要转动的圈数
func (t *Timer) getPositionAndCircle(d time.Duration) (pos int, circle int) {
	delaySeconds := int(d.Seconds())
	intervalSeconds := int(t.tick.Seconds())
	circle = int(delaySeconds / intervalSeconds / t.timeWheel.slotNum)
	pos = int(t.timeWheel.currentPos+delaySeconds/intervalSeconds) % t.timeWheel.slotNum

	return
}

func (t *Timer) step() {
	l := t.timeWheel.slots[t.timeWheel.currentPos]
	for e := l.Front(); e != nil; {
		job := e.Value.(*Node)
		if job.circle > 0 {
			job.circle--
			e = e.Next()
			continue
		}
		go job.handler(job.args)

		next := e.Next()
		if job.tType == Repetition {
			//循环的重新注册
			t.registerV1(job)
		} else {
			t.timerMap.Delete(job.id)
		}

		l.Remove(e)
		e = next
	}

	if t.timeWheel.currentPos == t.timeWheel.slotNum-1 {
		t.timeWheel.currentPos = 0
	} else {
		t.timeWheel.currentPos++
	}
}

//Stop 停止
func (t *Timer) Stop() {
	t.stop <- true
}

//Run 主循环
func (t *Timer) run() {
	t.ticker = time.NewTicker(t.tick)
	for {
		select {
		case <-t.ticker.C:
			t.step()
		case stop := <-t.stop:
			if stop {
				t.ticker.Stop()
			}
		}
	}
}
