package counter

type Counter struct {
	keyCounter chan int8
	quit       chan int8
	count      int64
}

func NewCounter() *Counter {
	toRet := &Counter{
		keyCounter: make(chan int8),
		quit:       make(chan int8),
		count:      0,
	}
	toRet.begin()
	return toRet
}

func (c *Counter) Update(count int8) {
	c.keyCounter <- count
}

func (c *Counter) Get() *int64 {
	return &c.count
}

func (c *Counter) Stop() {
	close(c.quit)
}

func (c *Counter) Reset() {
	//to ensure this probably is the last modification
	c.keyCounter <- 0
	c.count = 0
}

//it's possible ordeing could result in a negative count and still
//be eventually consistent
func (d *Counter) begin() {
	for {
		select {
		case num := <-d.keyCounter:
			d.count += int64(num)
		case <-d.quit:
			return
		}
	}
}
