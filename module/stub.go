package module

import "sync/atomic"

type ModuleInternal interface {
	Module
	IncrCalledCount()
	IncrAcceptedCount()
	IncrCompletedCount()
	IncrHandlingNumber()
	DecrHandlingNumber()
	Clear()
}

type myModule struct {
	calledCount    int64
	acceptedCount  int64
	completedCount int64
	handlingNumber int64
}

/*
 * create an instance of ModuleInternal
 */
func NewModuleInternal() (mi ModuleInternal) {
	return &myModule{}
}

/*Get*/

func (m *myModule) CalledCount() int64 {
	return atomic.LoadInt64(&m.calledCount)
}

func (m *myModule) AcceptedCount() int64 {
	return atomic.LoadInt64(&m.acceptedCount)
}

func (m *myModule) CompletedCount() int64 {
	return atomic.LoadInt64(&m.completedCount)
}

func (m *myModule) HandlingNumber() int64 {
	return atomic.LoadInt64(&m.handlingNumber)
}

func (m *myModule) Counts() Counts {
	return Counts{
		CalledCount:    m.CalledCount(),
		AcceptedCount:  m.AcceptedCount(),
		CompletedCount: m.CompletedCount(),
		HandlingNumber: m.HandlingNumber(),
	}
}

func (m *myModule) Summary() SummaryStruct {
	counts := m.Counts()
	return SummaryStruct{
		Called:    counts.CalledCount,
		Accepted:  counts.AcceptedCount,
		Completed: counts.CompletedCount,
		Handling:  counts.HandlingNumber,
		Extra:     nil,
	}
}

/*Update*/

func (m *myModule) IncrCalledCount() {
	atomic.AddInt64(&m.calledCount, 1)
}

func (m *myModule) IncrAcceptedCount() {
	atomic.AddInt64(&m.acceptedCount, 1)
}

func (m *myModule) IncrCompletedCount() {
	atomic.AddInt64(&m.completedCount, 1)
}

func (m *myModule) IncrHandlingNumber() {
	atomic.AddInt64(&m.handlingNumber, 1)
}

func (m *myModule) DecrHandlingNumber() {
	atomic.AddInt64(&m.handlingNumber, -1)
}

func (m *myModule) Clear() {
	atomic.StoreInt64(&m.calledCount, 0)
	atomic.StoreInt64(&m.acceptedCount, 0)
	atomic.StoreInt64(&m.completedCount, 0)
	atomic.StoreInt64(&m.handlingNumber, 0)
}
