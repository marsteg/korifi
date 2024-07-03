// Code generated by counterfeiter. DO NOT EDIT.
package fake

import (
	"context"
	"sync"

	"code.cloudfoundry.org/korifi/api/authorization"
	"code.cloudfoundry.org/korifi/api/handlers"
	"code.cloudfoundry.org/korifi/api/repositories"
)

type CFProcessRepository struct {
	CreateProcessStub        func(context.Context, authorization.Info, repositories.CreateProcessMessage) error
	createProcessMutex       sync.RWMutex
	createProcessArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.CreateProcessMessage
	}
	createProcessReturns struct {
		result1 error
	}
	createProcessReturnsOnCall map[int]struct {
		result1 error
	}
	GetAppRevisionStub        func(context.Context, authorization.Info, string) (string, error)
	getAppRevisionMutex       sync.RWMutex
	getAppRevisionArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 string
	}
	getAppRevisionReturns struct {
		result1 string
		result2 error
	}
	getAppRevisionReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	GetProcessStub        func(context.Context, authorization.Info, string) (repositories.ProcessRecord, error)
	getProcessMutex       sync.RWMutex
	getProcessArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 string
	}
	getProcessReturns struct {
		result1 repositories.ProcessRecord
		result2 error
	}
	getProcessReturnsOnCall map[int]struct {
		result1 repositories.ProcessRecord
		result2 error
	}
	GetProcessByAppTypeAndSpaceStub        func(context.Context, authorization.Info, string, string, string) (repositories.ProcessRecord, error)
	getProcessByAppTypeAndSpaceMutex       sync.RWMutex
	getProcessByAppTypeAndSpaceArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 string
		arg4 string
		arg5 string
	}
	getProcessByAppTypeAndSpaceReturns struct {
		result1 repositories.ProcessRecord
		result2 error
	}
	getProcessByAppTypeAndSpaceReturnsOnCall map[int]struct {
		result1 repositories.ProcessRecord
		result2 error
	}
	ListProcessesStub        func(context.Context, authorization.Info, repositories.ListProcessesMessage) ([]repositories.ProcessRecord, error)
	listProcessesMutex       sync.RWMutex
	listProcessesArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.ListProcessesMessage
	}
	listProcessesReturns struct {
		result1 []repositories.ProcessRecord
		result2 error
	}
	listProcessesReturnsOnCall map[int]struct {
		result1 []repositories.ProcessRecord
		result2 error
	}
	PatchProcessStub        func(context.Context, authorization.Info, repositories.PatchProcessMessage) (repositories.ProcessRecord, error)
	patchProcessMutex       sync.RWMutex
	patchProcessArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.PatchProcessMessage
	}
	patchProcessReturns struct {
		result1 repositories.ProcessRecord
		result2 error
	}
	patchProcessReturnsOnCall map[int]struct {
		result1 repositories.ProcessRecord
		result2 error
	}
	ScaleProcessStub        func(context.Context, authorization.Info, repositories.ScaleProcessMessage) (repositories.ProcessRecord, error)
	scaleProcessMutex       sync.RWMutex
	scaleProcessArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.ScaleProcessMessage
	}
	scaleProcessReturns struct {
		result1 repositories.ProcessRecord
		result2 error
	}
	scaleProcessReturnsOnCall map[int]struct {
		result1 repositories.ProcessRecord
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *CFProcessRepository) CreateProcess(arg1 context.Context, arg2 authorization.Info, arg3 repositories.CreateProcessMessage) error {
	fake.createProcessMutex.Lock()
	ret, specificReturn := fake.createProcessReturnsOnCall[len(fake.createProcessArgsForCall)]
	fake.createProcessArgsForCall = append(fake.createProcessArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.CreateProcessMessage
	}{arg1, arg2, arg3})
	stub := fake.CreateProcessStub
	fakeReturns := fake.createProcessReturns
	fake.recordInvocation("CreateProcess", []interface{}{arg1, arg2, arg3})
	fake.createProcessMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *CFProcessRepository) CreateProcessCallCount() int {
	fake.createProcessMutex.RLock()
	defer fake.createProcessMutex.RUnlock()
	return len(fake.createProcessArgsForCall)
}

func (fake *CFProcessRepository) CreateProcessCalls(stub func(context.Context, authorization.Info, repositories.CreateProcessMessage) error) {
	fake.createProcessMutex.Lock()
	defer fake.createProcessMutex.Unlock()
	fake.CreateProcessStub = stub
}

func (fake *CFProcessRepository) CreateProcessArgsForCall(i int) (context.Context, authorization.Info, repositories.CreateProcessMessage) {
	fake.createProcessMutex.RLock()
	defer fake.createProcessMutex.RUnlock()
	argsForCall := fake.createProcessArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFProcessRepository) CreateProcessReturns(result1 error) {
	fake.createProcessMutex.Lock()
	defer fake.createProcessMutex.Unlock()
	fake.CreateProcessStub = nil
	fake.createProcessReturns = struct {
		result1 error
	}{result1}
}

func (fake *CFProcessRepository) CreateProcessReturnsOnCall(i int, result1 error) {
	fake.createProcessMutex.Lock()
	defer fake.createProcessMutex.Unlock()
	fake.CreateProcessStub = nil
	if fake.createProcessReturnsOnCall == nil {
		fake.createProcessReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.createProcessReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *CFProcessRepository) GetAppRevision(arg1 context.Context, arg2 authorization.Info, arg3 string) (string, error) {
	fake.getAppRevisionMutex.Lock()
	ret, specificReturn := fake.getAppRevisionReturnsOnCall[len(fake.getAppRevisionArgsForCall)]
	fake.getAppRevisionArgsForCall = append(fake.getAppRevisionArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 string
	}{arg1, arg2, arg3})
	stub := fake.GetAppRevisionStub
	fakeReturns := fake.getAppRevisionReturns
	fake.recordInvocation("GetAppRevision", []interface{}{arg1, arg2, arg3})
	fake.getAppRevisionMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFProcessRepository) GetAppRevisionCallCount() int {
	fake.getAppRevisionMutex.RLock()
	defer fake.getAppRevisionMutex.RUnlock()
	return len(fake.getAppRevisionArgsForCall)
}

func (fake *CFProcessRepository) GetAppRevisionCalls(stub func(context.Context, authorization.Info, string) (string, error)) {
	fake.getAppRevisionMutex.Lock()
	defer fake.getAppRevisionMutex.Unlock()
	fake.GetAppRevisionStub = stub
}

func (fake *CFProcessRepository) GetAppRevisionArgsForCall(i int) (context.Context, authorization.Info, string) {
	fake.getAppRevisionMutex.RLock()
	defer fake.getAppRevisionMutex.RUnlock()
	argsForCall := fake.getAppRevisionArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFProcessRepository) GetAppRevisionReturns(result1 string, result2 error) {
	fake.getAppRevisionMutex.Lock()
	defer fake.getAppRevisionMutex.Unlock()
	fake.GetAppRevisionStub = nil
	fake.getAppRevisionReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *CFProcessRepository) GetAppRevisionReturnsOnCall(i int, result1 string, result2 error) {
	fake.getAppRevisionMutex.Lock()
	defer fake.getAppRevisionMutex.Unlock()
	fake.GetAppRevisionStub = nil
	if fake.getAppRevisionReturnsOnCall == nil {
		fake.getAppRevisionReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.getAppRevisionReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *CFProcessRepository) GetProcess(arg1 context.Context, arg2 authorization.Info, arg3 string) (repositories.ProcessRecord, error) {
	fake.getProcessMutex.Lock()
	ret, specificReturn := fake.getProcessReturnsOnCall[len(fake.getProcessArgsForCall)]
	fake.getProcessArgsForCall = append(fake.getProcessArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 string
	}{arg1, arg2, arg3})
	stub := fake.GetProcessStub
	fakeReturns := fake.getProcessReturns
	fake.recordInvocation("GetProcess", []interface{}{arg1, arg2, arg3})
	fake.getProcessMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFProcessRepository) GetProcessCallCount() int {
	fake.getProcessMutex.RLock()
	defer fake.getProcessMutex.RUnlock()
	return len(fake.getProcessArgsForCall)
}

func (fake *CFProcessRepository) GetProcessCalls(stub func(context.Context, authorization.Info, string) (repositories.ProcessRecord, error)) {
	fake.getProcessMutex.Lock()
	defer fake.getProcessMutex.Unlock()
	fake.GetProcessStub = stub
}

func (fake *CFProcessRepository) GetProcessArgsForCall(i int) (context.Context, authorization.Info, string) {
	fake.getProcessMutex.RLock()
	defer fake.getProcessMutex.RUnlock()
	argsForCall := fake.getProcessArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFProcessRepository) GetProcessReturns(result1 repositories.ProcessRecord, result2 error) {
	fake.getProcessMutex.Lock()
	defer fake.getProcessMutex.Unlock()
	fake.GetProcessStub = nil
	fake.getProcessReturns = struct {
		result1 repositories.ProcessRecord
		result2 error
	}{result1, result2}
}

func (fake *CFProcessRepository) GetProcessReturnsOnCall(i int, result1 repositories.ProcessRecord, result2 error) {
	fake.getProcessMutex.Lock()
	defer fake.getProcessMutex.Unlock()
	fake.GetProcessStub = nil
	if fake.getProcessReturnsOnCall == nil {
		fake.getProcessReturnsOnCall = make(map[int]struct {
			result1 repositories.ProcessRecord
			result2 error
		})
	}
	fake.getProcessReturnsOnCall[i] = struct {
		result1 repositories.ProcessRecord
		result2 error
	}{result1, result2}
}

func (fake *CFProcessRepository) GetProcessByAppTypeAndSpace(arg1 context.Context, arg2 authorization.Info, arg3 string, arg4 string, arg5 string) (repositories.ProcessRecord, error) {
	fake.getProcessByAppTypeAndSpaceMutex.Lock()
	ret, specificReturn := fake.getProcessByAppTypeAndSpaceReturnsOnCall[len(fake.getProcessByAppTypeAndSpaceArgsForCall)]
	fake.getProcessByAppTypeAndSpaceArgsForCall = append(fake.getProcessByAppTypeAndSpaceArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 string
		arg4 string
		arg5 string
	}{arg1, arg2, arg3, arg4, arg5})
	stub := fake.GetProcessByAppTypeAndSpaceStub
	fakeReturns := fake.getProcessByAppTypeAndSpaceReturns
	fake.recordInvocation("GetProcessByAppTypeAndSpace", []interface{}{arg1, arg2, arg3, arg4, arg5})
	fake.getProcessByAppTypeAndSpaceMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4, arg5)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFProcessRepository) GetProcessByAppTypeAndSpaceCallCount() int {
	fake.getProcessByAppTypeAndSpaceMutex.RLock()
	defer fake.getProcessByAppTypeAndSpaceMutex.RUnlock()
	return len(fake.getProcessByAppTypeAndSpaceArgsForCall)
}

func (fake *CFProcessRepository) GetProcessByAppTypeAndSpaceCalls(stub func(context.Context, authorization.Info, string, string, string) (repositories.ProcessRecord, error)) {
	fake.getProcessByAppTypeAndSpaceMutex.Lock()
	defer fake.getProcessByAppTypeAndSpaceMutex.Unlock()
	fake.GetProcessByAppTypeAndSpaceStub = stub
}

func (fake *CFProcessRepository) GetProcessByAppTypeAndSpaceArgsForCall(i int) (context.Context, authorization.Info, string, string, string) {
	fake.getProcessByAppTypeAndSpaceMutex.RLock()
	defer fake.getProcessByAppTypeAndSpaceMutex.RUnlock()
	argsForCall := fake.getProcessByAppTypeAndSpaceArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5
}

func (fake *CFProcessRepository) GetProcessByAppTypeAndSpaceReturns(result1 repositories.ProcessRecord, result2 error) {
	fake.getProcessByAppTypeAndSpaceMutex.Lock()
	defer fake.getProcessByAppTypeAndSpaceMutex.Unlock()
	fake.GetProcessByAppTypeAndSpaceStub = nil
	fake.getProcessByAppTypeAndSpaceReturns = struct {
		result1 repositories.ProcessRecord
		result2 error
	}{result1, result2}
}

func (fake *CFProcessRepository) GetProcessByAppTypeAndSpaceReturnsOnCall(i int, result1 repositories.ProcessRecord, result2 error) {
	fake.getProcessByAppTypeAndSpaceMutex.Lock()
	defer fake.getProcessByAppTypeAndSpaceMutex.Unlock()
	fake.GetProcessByAppTypeAndSpaceStub = nil
	if fake.getProcessByAppTypeAndSpaceReturnsOnCall == nil {
		fake.getProcessByAppTypeAndSpaceReturnsOnCall = make(map[int]struct {
			result1 repositories.ProcessRecord
			result2 error
		})
	}
	fake.getProcessByAppTypeAndSpaceReturnsOnCall[i] = struct {
		result1 repositories.ProcessRecord
		result2 error
	}{result1, result2}
}

func (fake *CFProcessRepository) ListProcesses(arg1 context.Context, arg2 authorization.Info, arg3 repositories.ListProcessesMessage) ([]repositories.ProcessRecord, error) {
	fake.listProcessesMutex.Lock()
	ret, specificReturn := fake.listProcessesReturnsOnCall[len(fake.listProcessesArgsForCall)]
	fake.listProcessesArgsForCall = append(fake.listProcessesArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.ListProcessesMessage
	}{arg1, arg2, arg3})
	stub := fake.ListProcessesStub
	fakeReturns := fake.listProcessesReturns
	fake.recordInvocation("ListProcesses", []interface{}{arg1, arg2, arg3})
	fake.listProcessesMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFProcessRepository) ListProcessesCallCount() int {
	fake.listProcessesMutex.RLock()
	defer fake.listProcessesMutex.RUnlock()
	return len(fake.listProcessesArgsForCall)
}

func (fake *CFProcessRepository) ListProcessesCalls(stub func(context.Context, authorization.Info, repositories.ListProcessesMessage) ([]repositories.ProcessRecord, error)) {
	fake.listProcessesMutex.Lock()
	defer fake.listProcessesMutex.Unlock()
	fake.ListProcessesStub = stub
}

func (fake *CFProcessRepository) ListProcessesArgsForCall(i int) (context.Context, authorization.Info, repositories.ListProcessesMessage) {
	fake.listProcessesMutex.RLock()
	defer fake.listProcessesMutex.RUnlock()
	argsForCall := fake.listProcessesArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFProcessRepository) ListProcessesReturns(result1 []repositories.ProcessRecord, result2 error) {
	fake.listProcessesMutex.Lock()
	defer fake.listProcessesMutex.Unlock()
	fake.ListProcessesStub = nil
	fake.listProcessesReturns = struct {
		result1 []repositories.ProcessRecord
		result2 error
	}{result1, result2}
}

func (fake *CFProcessRepository) ListProcessesReturnsOnCall(i int, result1 []repositories.ProcessRecord, result2 error) {
	fake.listProcessesMutex.Lock()
	defer fake.listProcessesMutex.Unlock()
	fake.ListProcessesStub = nil
	if fake.listProcessesReturnsOnCall == nil {
		fake.listProcessesReturnsOnCall = make(map[int]struct {
			result1 []repositories.ProcessRecord
			result2 error
		})
	}
	fake.listProcessesReturnsOnCall[i] = struct {
		result1 []repositories.ProcessRecord
		result2 error
	}{result1, result2}
}

func (fake *CFProcessRepository) PatchProcess(arg1 context.Context, arg2 authorization.Info, arg3 repositories.PatchProcessMessage) (repositories.ProcessRecord, error) {
	fake.patchProcessMutex.Lock()
	ret, specificReturn := fake.patchProcessReturnsOnCall[len(fake.patchProcessArgsForCall)]
	fake.patchProcessArgsForCall = append(fake.patchProcessArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.PatchProcessMessage
	}{arg1, arg2, arg3})
	stub := fake.PatchProcessStub
	fakeReturns := fake.patchProcessReturns
	fake.recordInvocation("PatchProcess", []interface{}{arg1, arg2, arg3})
	fake.patchProcessMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFProcessRepository) PatchProcessCallCount() int {
	fake.patchProcessMutex.RLock()
	defer fake.patchProcessMutex.RUnlock()
	return len(fake.patchProcessArgsForCall)
}

func (fake *CFProcessRepository) PatchProcessCalls(stub func(context.Context, authorization.Info, repositories.PatchProcessMessage) (repositories.ProcessRecord, error)) {
	fake.patchProcessMutex.Lock()
	defer fake.patchProcessMutex.Unlock()
	fake.PatchProcessStub = stub
}

func (fake *CFProcessRepository) PatchProcessArgsForCall(i int) (context.Context, authorization.Info, repositories.PatchProcessMessage) {
	fake.patchProcessMutex.RLock()
	defer fake.patchProcessMutex.RUnlock()
	argsForCall := fake.patchProcessArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFProcessRepository) PatchProcessReturns(result1 repositories.ProcessRecord, result2 error) {
	fake.patchProcessMutex.Lock()
	defer fake.patchProcessMutex.Unlock()
	fake.PatchProcessStub = nil
	fake.patchProcessReturns = struct {
		result1 repositories.ProcessRecord
		result2 error
	}{result1, result2}
}

func (fake *CFProcessRepository) PatchProcessReturnsOnCall(i int, result1 repositories.ProcessRecord, result2 error) {
	fake.patchProcessMutex.Lock()
	defer fake.patchProcessMutex.Unlock()
	fake.PatchProcessStub = nil
	if fake.patchProcessReturnsOnCall == nil {
		fake.patchProcessReturnsOnCall = make(map[int]struct {
			result1 repositories.ProcessRecord
			result2 error
		})
	}
	fake.patchProcessReturnsOnCall[i] = struct {
		result1 repositories.ProcessRecord
		result2 error
	}{result1, result2}
}

func (fake *CFProcessRepository) ScaleProcess(arg1 context.Context, arg2 authorization.Info, arg3 repositories.ScaleProcessMessage) (repositories.ProcessRecord, error) {
	fake.scaleProcessMutex.Lock()
	ret, specificReturn := fake.scaleProcessReturnsOnCall[len(fake.scaleProcessArgsForCall)]
	fake.scaleProcessArgsForCall = append(fake.scaleProcessArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.ScaleProcessMessage
	}{arg1, arg2, arg3})
	stub := fake.ScaleProcessStub
	fakeReturns := fake.scaleProcessReturns
	fake.recordInvocation("ScaleProcess", []interface{}{arg1, arg2, arg3})
	fake.scaleProcessMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFProcessRepository) ScaleProcessCallCount() int {
	fake.scaleProcessMutex.RLock()
	defer fake.scaleProcessMutex.RUnlock()
	return len(fake.scaleProcessArgsForCall)
}

func (fake *CFProcessRepository) ScaleProcessCalls(stub func(context.Context, authorization.Info, repositories.ScaleProcessMessage) (repositories.ProcessRecord, error)) {
	fake.scaleProcessMutex.Lock()
	defer fake.scaleProcessMutex.Unlock()
	fake.ScaleProcessStub = stub
}

func (fake *CFProcessRepository) ScaleProcessArgsForCall(i int) (context.Context, authorization.Info, repositories.ScaleProcessMessage) {
	fake.scaleProcessMutex.RLock()
	defer fake.scaleProcessMutex.RUnlock()
	argsForCall := fake.scaleProcessArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFProcessRepository) ScaleProcessReturns(result1 repositories.ProcessRecord, result2 error) {
	fake.scaleProcessMutex.Lock()
	defer fake.scaleProcessMutex.Unlock()
	fake.ScaleProcessStub = nil
	fake.scaleProcessReturns = struct {
		result1 repositories.ProcessRecord
		result2 error
	}{result1, result2}
}

func (fake *CFProcessRepository) ScaleProcessReturnsOnCall(i int, result1 repositories.ProcessRecord, result2 error) {
	fake.scaleProcessMutex.Lock()
	defer fake.scaleProcessMutex.Unlock()
	fake.ScaleProcessStub = nil
	if fake.scaleProcessReturnsOnCall == nil {
		fake.scaleProcessReturnsOnCall = make(map[int]struct {
			result1 repositories.ProcessRecord
			result2 error
		})
	}
	fake.scaleProcessReturnsOnCall[i] = struct {
		result1 repositories.ProcessRecord
		result2 error
	}{result1, result2}
}

func (fake *CFProcessRepository) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createProcessMutex.RLock()
	defer fake.createProcessMutex.RUnlock()
	fake.getAppRevisionMutex.RLock()
	defer fake.getAppRevisionMutex.RUnlock()
	fake.getProcessMutex.RLock()
	defer fake.getProcessMutex.RUnlock()
	fake.getProcessByAppTypeAndSpaceMutex.RLock()
	defer fake.getProcessByAppTypeAndSpaceMutex.RUnlock()
	fake.listProcessesMutex.RLock()
	defer fake.listProcessesMutex.RUnlock()
	fake.patchProcessMutex.RLock()
	defer fake.patchProcessMutex.RUnlock()
	fake.scaleProcessMutex.RLock()
	defer fake.scaleProcessMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *CFProcessRepository) recordInvocation(key string, args []interface{}) {
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

var _ handlers.CFProcessRepository = new(CFProcessRepository)
