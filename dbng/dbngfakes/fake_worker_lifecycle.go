// This file was generated by counterfeiter
package dbngfakes

import (
	"sync"

	"github.com/concourse/atc/dbng"
)

type FakeWorkerLifecycle struct {
	StallUnresponsiveWorkersStub        func() ([]string, error)
	stallUnresponsiveWorkersMutex       sync.RWMutex
	stallUnresponsiveWorkersArgsForCall []struct{}
	stallUnresponsiveWorkersReturns     struct {
		result1 []string
		result2 error
	}
	LandFinishedLandingWorkersStub        func() ([]string, error)
	landFinishedLandingWorkersMutex       sync.RWMutex
	landFinishedLandingWorkersArgsForCall []struct{}
	landFinishedLandingWorkersReturns     struct {
		result1 []string
		result2 error
	}
	DeleteFinishedRetiringWorkersStub        func() ([]string, error)
	deleteFinishedRetiringWorkersMutex       sync.RWMutex
	deleteFinishedRetiringWorkersArgsForCall []struct{}
	deleteFinishedRetiringWorkersReturns     struct {
		result1 []string
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeWorkerLifecycle) StallUnresponsiveWorkers() ([]string, error) {
	fake.stallUnresponsiveWorkersMutex.Lock()
	fake.stallUnresponsiveWorkersArgsForCall = append(fake.stallUnresponsiveWorkersArgsForCall, struct{}{})
	fake.recordInvocation("StallUnresponsiveWorkers", []interface{}{})
	fake.stallUnresponsiveWorkersMutex.Unlock()
	if fake.StallUnresponsiveWorkersStub != nil {
		return fake.StallUnresponsiveWorkersStub()
	} else {
		return fake.stallUnresponsiveWorkersReturns.result1, fake.stallUnresponsiveWorkersReturns.result2
	}
}

func (fake *FakeWorkerLifecycle) StallUnresponsiveWorkersCallCount() int {
	fake.stallUnresponsiveWorkersMutex.RLock()
	defer fake.stallUnresponsiveWorkersMutex.RUnlock()
	return len(fake.stallUnresponsiveWorkersArgsForCall)
}

func (fake *FakeWorkerLifecycle) StallUnresponsiveWorkersReturns(result1 []string, result2 error) {
	fake.StallUnresponsiveWorkersStub = nil
	fake.stallUnresponsiveWorkersReturns = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeWorkerLifecycle) LandFinishedLandingWorkers() ([]string, error) {
	fake.landFinishedLandingWorkersMutex.Lock()
	fake.landFinishedLandingWorkersArgsForCall = append(fake.landFinishedLandingWorkersArgsForCall, struct{}{})
	fake.recordInvocation("LandFinishedLandingWorkers", []interface{}{})
	fake.landFinishedLandingWorkersMutex.Unlock()
	if fake.LandFinishedLandingWorkersStub != nil {
		return fake.LandFinishedLandingWorkersStub()
	} else {
		return fake.landFinishedLandingWorkersReturns.result1, fake.landFinishedLandingWorkersReturns.result2
	}
}

func (fake *FakeWorkerLifecycle) LandFinishedLandingWorkersCallCount() int {
	fake.landFinishedLandingWorkersMutex.RLock()
	defer fake.landFinishedLandingWorkersMutex.RUnlock()
	return len(fake.landFinishedLandingWorkersArgsForCall)
}

func (fake *FakeWorkerLifecycle) LandFinishedLandingWorkersReturns(result1 []string, result2 error) {
	fake.LandFinishedLandingWorkersStub = nil
	fake.landFinishedLandingWorkersReturns = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeWorkerLifecycle) DeleteFinishedRetiringWorkers() ([]string, error) {
	fake.deleteFinishedRetiringWorkersMutex.Lock()
	fake.deleteFinishedRetiringWorkersArgsForCall = append(fake.deleteFinishedRetiringWorkersArgsForCall, struct{}{})
	fake.recordInvocation("DeleteFinishedRetiringWorkers", []interface{}{})
	fake.deleteFinishedRetiringWorkersMutex.Unlock()
	if fake.DeleteFinishedRetiringWorkersStub != nil {
		return fake.DeleteFinishedRetiringWorkersStub()
	} else {
		return fake.deleteFinishedRetiringWorkersReturns.result1, fake.deleteFinishedRetiringWorkersReturns.result2
	}
}

func (fake *FakeWorkerLifecycle) DeleteFinishedRetiringWorkersCallCount() int {
	fake.deleteFinishedRetiringWorkersMutex.RLock()
	defer fake.deleteFinishedRetiringWorkersMutex.RUnlock()
	return len(fake.deleteFinishedRetiringWorkersArgsForCall)
}

func (fake *FakeWorkerLifecycle) DeleteFinishedRetiringWorkersReturns(result1 []string, result2 error) {
	fake.DeleteFinishedRetiringWorkersStub = nil
	fake.deleteFinishedRetiringWorkersReturns = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeWorkerLifecycle) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.stallUnresponsiveWorkersMutex.RLock()
	defer fake.stallUnresponsiveWorkersMutex.RUnlock()
	fake.landFinishedLandingWorkersMutex.RLock()
	defer fake.landFinishedLandingWorkersMutex.RUnlock()
	fake.deleteFinishedRetiringWorkersMutex.RLock()
	defer fake.deleteFinishedRetiringWorkersMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeWorkerLifecycle) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ dbng.WorkerLifecycle = new(FakeWorkerLifecycle)
