package particles

type Burst struct {
	time      float64
	nextTime  float64
	completed int

	Count    int
	Cycles   int
	Interval float64

	//Not implemented
	//StartTime   float64
	//Probability float64
}

func (b *Burst) update(dt float64) {
	b.time += dt
}

func (b *Burst) due() int {
	var n int
	for b.nextTime <= b.time && b.completed < b.Cycles{
		n += b.Count

		b.nextTime += b.Interval
		b.completed++
	}
	return n
}
