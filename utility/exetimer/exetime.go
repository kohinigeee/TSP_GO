package exetimer

import "time"

type ExeTimer struct {
	startTime    time.Time
	endTime      time.Time
	isMeasureEnd bool
}

func newExeTime(startTime time.Time, endtime time.Time) *ExeTimer {
	return &ExeTimer{
		startTime:    startTime,
		endTime:      endtime,
		isMeasureEnd: true,
	}
}

func MeasureStart() *ExeTimer {
	startTime := time.Now()
	return &ExeTimer{
		startTime:    startTime,
		isMeasureEnd: false,
	}
}

func (t *ExeTimer) MeasureEnd() {
	if t.isMeasureEnd {
		panic("ExeTimer: MeasureEnd() called twice")
	}
	t.endTime = time.Now()
	t.isMeasureEnd = true
}

func (t *ExeTimer) Elaplsed() time.Duration {
	if !t.isMeasureEnd {
		panic("ExeTimer: MeasureEnd() is not called")
	}
	return t.endTime.Sub(t.startTime)
}

func (t *ExeTimer) ElapsedMilliSeconds() int64 {
	return t.Elaplsed().Milliseconds()
}

func (t *ExeTimer) ElapsedSeconds() float64 {
	return t.Elaplsed().Seconds()
}
