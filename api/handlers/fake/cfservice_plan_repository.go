// Code generated by counterfeiter. DO NOT EDIT.
package fake

import (
	"context"
	"sync"

	"code.cloudfoundry.org/korifi/api/authorization"
	"code.cloudfoundry.org/korifi/api/handlers"
	"code.cloudfoundry.org/korifi/api/repositories"
)

type CFServicePlanRepository struct {
	ApplyPlanVisibilityStub        func(context.Context, authorization.Info, repositories.ApplyServicePlanVisibilityMessage) (repositories.ServicePlanRecord, error)
	applyPlanVisibilityMutex       sync.RWMutex
	applyPlanVisibilityArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.ApplyServicePlanVisibilityMessage
	}
	applyPlanVisibilityReturns struct {
		result1 repositories.ServicePlanRecord
		result2 error
	}
	applyPlanVisibilityReturnsOnCall map[int]struct {
		result1 repositories.ServicePlanRecord
		result2 error
	}
	GetPlanStub        func(context.Context, authorization.Info, string) (repositories.ServicePlanRecord, error)
	getPlanMutex       sync.RWMutex
	getPlanArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 string
	}
	getPlanReturns struct {
		result1 repositories.ServicePlanRecord
		result2 error
	}
	getPlanReturnsOnCall map[int]struct {
		result1 repositories.ServicePlanRecord
		result2 error
	}
	ListPlansStub        func(context.Context, authorization.Info, repositories.ListServicePlanMessage) ([]repositories.ServicePlanRecord, error)
	listPlansMutex       sync.RWMutex
	listPlansArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.ListServicePlanMessage
	}
	listPlansReturns struct {
		result1 []repositories.ServicePlanRecord
		result2 error
	}
	listPlansReturnsOnCall map[int]struct {
		result1 []repositories.ServicePlanRecord
		result2 error
	}
	UpdatePlanVisibilityStub        func(context.Context, authorization.Info, repositories.UpdateServicePlanVisibilityMessage) (repositories.ServicePlanRecord, error)
	updatePlanVisibilityMutex       sync.RWMutex
	updatePlanVisibilityArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.UpdateServicePlanVisibilityMessage
	}
	updatePlanVisibilityReturns struct {
		result1 repositories.ServicePlanRecord
		result2 error
	}
	updatePlanVisibilityReturnsOnCall map[int]struct {
		result1 repositories.ServicePlanRecord
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *CFServicePlanRepository) ApplyPlanVisibility(arg1 context.Context, arg2 authorization.Info, arg3 repositories.ApplyServicePlanVisibilityMessage) (repositories.ServicePlanRecord, error) {
	fake.applyPlanVisibilityMutex.Lock()
	ret, specificReturn := fake.applyPlanVisibilityReturnsOnCall[len(fake.applyPlanVisibilityArgsForCall)]
	fake.applyPlanVisibilityArgsForCall = append(fake.applyPlanVisibilityArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.ApplyServicePlanVisibilityMessage
	}{arg1, arg2, arg3})
	stub := fake.ApplyPlanVisibilityStub
	fakeReturns := fake.applyPlanVisibilityReturns
	fake.recordInvocation("ApplyPlanVisibility", []interface{}{arg1, arg2, arg3})
	fake.applyPlanVisibilityMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFServicePlanRepository) ApplyPlanVisibilityCallCount() int {
	fake.applyPlanVisibilityMutex.RLock()
	defer fake.applyPlanVisibilityMutex.RUnlock()
	return len(fake.applyPlanVisibilityArgsForCall)
}

func (fake *CFServicePlanRepository) ApplyPlanVisibilityCalls(stub func(context.Context, authorization.Info, repositories.ApplyServicePlanVisibilityMessage) (repositories.ServicePlanRecord, error)) {
	fake.applyPlanVisibilityMutex.Lock()
	defer fake.applyPlanVisibilityMutex.Unlock()
	fake.ApplyPlanVisibilityStub = stub
}

func (fake *CFServicePlanRepository) ApplyPlanVisibilityArgsForCall(i int) (context.Context, authorization.Info, repositories.ApplyServicePlanVisibilityMessage) {
	fake.applyPlanVisibilityMutex.RLock()
	defer fake.applyPlanVisibilityMutex.RUnlock()
	argsForCall := fake.applyPlanVisibilityArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFServicePlanRepository) ApplyPlanVisibilityReturns(result1 repositories.ServicePlanRecord, result2 error) {
	fake.applyPlanVisibilityMutex.Lock()
	defer fake.applyPlanVisibilityMutex.Unlock()
	fake.ApplyPlanVisibilityStub = nil
	fake.applyPlanVisibilityReturns = struct {
		result1 repositories.ServicePlanRecord
		result2 error
	}{result1, result2}
}

func (fake *CFServicePlanRepository) ApplyPlanVisibilityReturnsOnCall(i int, result1 repositories.ServicePlanRecord, result2 error) {
	fake.applyPlanVisibilityMutex.Lock()
	defer fake.applyPlanVisibilityMutex.Unlock()
	fake.ApplyPlanVisibilityStub = nil
	if fake.applyPlanVisibilityReturnsOnCall == nil {
		fake.applyPlanVisibilityReturnsOnCall = make(map[int]struct {
			result1 repositories.ServicePlanRecord
			result2 error
		})
	}
	fake.applyPlanVisibilityReturnsOnCall[i] = struct {
		result1 repositories.ServicePlanRecord
		result2 error
	}{result1, result2}
}

func (fake *CFServicePlanRepository) GetPlan(arg1 context.Context, arg2 authorization.Info, arg3 string) (repositories.ServicePlanRecord, error) {
	fake.getPlanMutex.Lock()
	ret, specificReturn := fake.getPlanReturnsOnCall[len(fake.getPlanArgsForCall)]
	fake.getPlanArgsForCall = append(fake.getPlanArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 string
	}{arg1, arg2, arg3})
	stub := fake.GetPlanStub
	fakeReturns := fake.getPlanReturns
	fake.recordInvocation("GetPlan", []interface{}{arg1, arg2, arg3})
	fake.getPlanMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFServicePlanRepository) GetPlanCallCount() int {
	fake.getPlanMutex.RLock()
	defer fake.getPlanMutex.RUnlock()
	return len(fake.getPlanArgsForCall)
}

func (fake *CFServicePlanRepository) GetPlanCalls(stub func(context.Context, authorization.Info, string) (repositories.ServicePlanRecord, error)) {
	fake.getPlanMutex.Lock()
	defer fake.getPlanMutex.Unlock()
	fake.GetPlanStub = stub
}

func (fake *CFServicePlanRepository) GetPlanArgsForCall(i int) (context.Context, authorization.Info, string) {
	fake.getPlanMutex.RLock()
	defer fake.getPlanMutex.RUnlock()
	argsForCall := fake.getPlanArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFServicePlanRepository) GetPlanReturns(result1 repositories.ServicePlanRecord, result2 error) {
	fake.getPlanMutex.Lock()
	defer fake.getPlanMutex.Unlock()
	fake.GetPlanStub = nil
	fake.getPlanReturns = struct {
		result1 repositories.ServicePlanRecord
		result2 error
	}{result1, result2}
}

func (fake *CFServicePlanRepository) GetPlanReturnsOnCall(i int, result1 repositories.ServicePlanRecord, result2 error) {
	fake.getPlanMutex.Lock()
	defer fake.getPlanMutex.Unlock()
	fake.GetPlanStub = nil
	if fake.getPlanReturnsOnCall == nil {
		fake.getPlanReturnsOnCall = make(map[int]struct {
			result1 repositories.ServicePlanRecord
			result2 error
		})
	}
	fake.getPlanReturnsOnCall[i] = struct {
		result1 repositories.ServicePlanRecord
		result2 error
	}{result1, result2}
}

func (fake *CFServicePlanRepository) ListPlans(arg1 context.Context, arg2 authorization.Info, arg3 repositories.ListServicePlanMessage) ([]repositories.ServicePlanRecord, error) {
	fake.listPlansMutex.Lock()
	ret, specificReturn := fake.listPlansReturnsOnCall[len(fake.listPlansArgsForCall)]
	fake.listPlansArgsForCall = append(fake.listPlansArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.ListServicePlanMessage
	}{arg1, arg2, arg3})
	stub := fake.ListPlansStub
	fakeReturns := fake.listPlansReturns
	fake.recordInvocation("ListPlans", []interface{}{arg1, arg2, arg3})
	fake.listPlansMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFServicePlanRepository) ListPlansCallCount() int {
	fake.listPlansMutex.RLock()
	defer fake.listPlansMutex.RUnlock()
	return len(fake.listPlansArgsForCall)
}

func (fake *CFServicePlanRepository) ListPlansCalls(stub func(context.Context, authorization.Info, repositories.ListServicePlanMessage) ([]repositories.ServicePlanRecord, error)) {
	fake.listPlansMutex.Lock()
	defer fake.listPlansMutex.Unlock()
	fake.ListPlansStub = stub
}

func (fake *CFServicePlanRepository) ListPlansArgsForCall(i int) (context.Context, authorization.Info, repositories.ListServicePlanMessage) {
	fake.listPlansMutex.RLock()
	defer fake.listPlansMutex.RUnlock()
	argsForCall := fake.listPlansArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFServicePlanRepository) ListPlansReturns(result1 []repositories.ServicePlanRecord, result2 error) {
	fake.listPlansMutex.Lock()
	defer fake.listPlansMutex.Unlock()
	fake.ListPlansStub = nil
	fake.listPlansReturns = struct {
		result1 []repositories.ServicePlanRecord
		result2 error
	}{result1, result2}
}

func (fake *CFServicePlanRepository) ListPlansReturnsOnCall(i int, result1 []repositories.ServicePlanRecord, result2 error) {
	fake.listPlansMutex.Lock()
	defer fake.listPlansMutex.Unlock()
	fake.ListPlansStub = nil
	if fake.listPlansReturnsOnCall == nil {
		fake.listPlansReturnsOnCall = make(map[int]struct {
			result1 []repositories.ServicePlanRecord
			result2 error
		})
	}
	fake.listPlansReturnsOnCall[i] = struct {
		result1 []repositories.ServicePlanRecord
		result2 error
	}{result1, result2}
}

func (fake *CFServicePlanRepository) UpdatePlanVisibility(arg1 context.Context, arg2 authorization.Info, arg3 repositories.UpdateServicePlanVisibilityMessage) (repositories.ServicePlanRecord, error) {
	fake.updatePlanVisibilityMutex.Lock()
	ret, specificReturn := fake.updatePlanVisibilityReturnsOnCall[len(fake.updatePlanVisibilityArgsForCall)]
	fake.updatePlanVisibilityArgsForCall = append(fake.updatePlanVisibilityArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.UpdateServicePlanVisibilityMessage
	}{arg1, arg2, arg3})
	stub := fake.UpdatePlanVisibilityStub
	fakeReturns := fake.updatePlanVisibilityReturns
	fake.recordInvocation("UpdatePlanVisibility", []interface{}{arg1, arg2, arg3})
	fake.updatePlanVisibilityMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFServicePlanRepository) UpdatePlanVisibilityCallCount() int {
	fake.updatePlanVisibilityMutex.RLock()
	defer fake.updatePlanVisibilityMutex.RUnlock()
	return len(fake.updatePlanVisibilityArgsForCall)
}

func (fake *CFServicePlanRepository) UpdatePlanVisibilityCalls(stub func(context.Context, authorization.Info, repositories.UpdateServicePlanVisibilityMessage) (repositories.ServicePlanRecord, error)) {
	fake.updatePlanVisibilityMutex.Lock()
	defer fake.updatePlanVisibilityMutex.Unlock()
	fake.UpdatePlanVisibilityStub = stub
}

func (fake *CFServicePlanRepository) UpdatePlanVisibilityArgsForCall(i int) (context.Context, authorization.Info, repositories.UpdateServicePlanVisibilityMessage) {
	fake.updatePlanVisibilityMutex.RLock()
	defer fake.updatePlanVisibilityMutex.RUnlock()
	argsForCall := fake.updatePlanVisibilityArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFServicePlanRepository) UpdatePlanVisibilityReturns(result1 repositories.ServicePlanRecord, result2 error) {
	fake.updatePlanVisibilityMutex.Lock()
	defer fake.updatePlanVisibilityMutex.Unlock()
	fake.UpdatePlanVisibilityStub = nil
	fake.updatePlanVisibilityReturns = struct {
		result1 repositories.ServicePlanRecord
		result2 error
	}{result1, result2}
}

func (fake *CFServicePlanRepository) UpdatePlanVisibilityReturnsOnCall(i int, result1 repositories.ServicePlanRecord, result2 error) {
	fake.updatePlanVisibilityMutex.Lock()
	defer fake.updatePlanVisibilityMutex.Unlock()
	fake.UpdatePlanVisibilityStub = nil
	if fake.updatePlanVisibilityReturnsOnCall == nil {
		fake.updatePlanVisibilityReturnsOnCall = make(map[int]struct {
			result1 repositories.ServicePlanRecord
			result2 error
		})
	}
	fake.updatePlanVisibilityReturnsOnCall[i] = struct {
		result1 repositories.ServicePlanRecord
		result2 error
	}{result1, result2}
}

func (fake *CFServicePlanRepository) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.applyPlanVisibilityMutex.RLock()
	defer fake.applyPlanVisibilityMutex.RUnlock()
	fake.getPlanMutex.RLock()
	defer fake.getPlanMutex.RUnlock()
	fake.listPlansMutex.RLock()
	defer fake.listPlansMutex.RUnlock()
	fake.updatePlanVisibilityMutex.RLock()
	defer fake.updatePlanVisibilityMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *CFServicePlanRepository) recordInvocation(key string, args []interface{}) {
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

var _ handlers.CFServicePlanRepository = new(CFServicePlanRepository)
