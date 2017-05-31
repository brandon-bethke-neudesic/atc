// Code generated by counterfeiter. DO NOT EDIT.
package schedulerfakes

import (
	"sync"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc"
	"github.com/concourse/atc/dbng"
	"github.com/concourse/atc/scheduler"
)

type FakeBuildStarter struct {
	TryStartPendingBuildsForJobStub        func(logger lager.Logger, job dbng.Job, resources dbng.Resources, resourceTypes atc.VersionedResourceTypes, nextPendingBuilds []dbng.Build) error
	tryStartPendingBuildsForJobMutex       sync.RWMutex
	tryStartPendingBuildsForJobArgsForCall []struct {
		logger            lager.Logger
		job               dbng.Job
		resources         dbng.Resources
		resourceTypes     atc.VersionedResourceTypes
		nextPendingBuilds []dbng.Build
	}
	tryStartPendingBuildsForJobReturns struct {
		result1 error
	}
	tryStartPendingBuildsForJobReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeBuildStarter) TryStartPendingBuildsForJob(logger lager.Logger, job dbng.Job, resources dbng.Resources, resourceTypes atc.VersionedResourceTypes, nextPendingBuilds []dbng.Build) error {
	var nextPendingBuildsCopy []dbng.Build
	if nextPendingBuilds != nil {
		nextPendingBuildsCopy = make([]dbng.Build, len(nextPendingBuilds))
		copy(nextPendingBuildsCopy, nextPendingBuilds)
	}
	fake.tryStartPendingBuildsForJobMutex.Lock()
	ret, specificReturn := fake.tryStartPendingBuildsForJobReturnsOnCall[len(fake.tryStartPendingBuildsForJobArgsForCall)]
	fake.tryStartPendingBuildsForJobArgsForCall = append(fake.tryStartPendingBuildsForJobArgsForCall, struct {
		logger            lager.Logger
		job               dbng.Job
		resources         dbng.Resources
		resourceTypes     atc.VersionedResourceTypes
		nextPendingBuilds []dbng.Build
	}{logger, job, resources, resourceTypes, nextPendingBuildsCopy})
	fake.recordInvocation("TryStartPendingBuildsForJob", []interface{}{logger, job, resources, resourceTypes, nextPendingBuildsCopy})
	fake.tryStartPendingBuildsForJobMutex.Unlock()
	if fake.TryStartPendingBuildsForJobStub != nil {
		return fake.TryStartPendingBuildsForJobStub(logger, job, resources, resourceTypes, nextPendingBuilds)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.tryStartPendingBuildsForJobReturns.result1
}

func (fake *FakeBuildStarter) TryStartPendingBuildsForJobCallCount() int {
	fake.tryStartPendingBuildsForJobMutex.RLock()
	defer fake.tryStartPendingBuildsForJobMutex.RUnlock()
	return len(fake.tryStartPendingBuildsForJobArgsForCall)
}

func (fake *FakeBuildStarter) TryStartPendingBuildsForJobArgsForCall(i int) (lager.Logger, dbng.Job, dbng.Resources, atc.VersionedResourceTypes, []dbng.Build) {
	fake.tryStartPendingBuildsForJobMutex.RLock()
	defer fake.tryStartPendingBuildsForJobMutex.RUnlock()
	return fake.tryStartPendingBuildsForJobArgsForCall[i].logger, fake.tryStartPendingBuildsForJobArgsForCall[i].job, fake.tryStartPendingBuildsForJobArgsForCall[i].resources, fake.tryStartPendingBuildsForJobArgsForCall[i].resourceTypes, fake.tryStartPendingBuildsForJobArgsForCall[i].nextPendingBuilds
}

func (fake *FakeBuildStarter) TryStartPendingBuildsForJobReturns(result1 error) {
	fake.TryStartPendingBuildsForJobStub = nil
	fake.tryStartPendingBuildsForJobReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeBuildStarter) TryStartPendingBuildsForJobReturnsOnCall(i int, result1 error) {
	fake.TryStartPendingBuildsForJobStub = nil
	if fake.tryStartPendingBuildsForJobReturnsOnCall == nil {
		fake.tryStartPendingBuildsForJobReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.tryStartPendingBuildsForJobReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeBuildStarter) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.tryStartPendingBuildsForJobMutex.RLock()
	defer fake.tryStartPendingBuildsForJobMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeBuildStarter) recordInvocation(key string, args []interface{}) {
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

var _ scheduler.BuildStarter = new(FakeBuildStarter)
