package handlers

import (
	"time"
)

const (
	EnvironmentVariablesGroupsRoot = "/v3/environment_variable_groups"
	EnvironmentVariablesGroupPath  = "/v3/environment_variable_groups/{guid}"
)

//counterfeiter:generate -o fake -fake-name CFEnvironmentVariableGroupRepository . CFEnvironmentVariableGroupRepository
type EnvironmentVariable struct {
	value string
}

type EnvironmentVariablesGroup struct {
	name       string
	updated_at time.Time
	variables  []EnvironmentVariable
}

func NewEnvironmentVariablesGroup(name string) *EnvironmentVariablesGroup {
	variables := make(map[string]EnvironmentVariable)
	return &EnvironmentVariablesGroup{
		name:       name,
		updated_at: time.Now().UTC(),
		variables:  variables,
	}
}

/*
TODO:
- create initital running and staging secret during kind deploy (installation)
- copy the secrets into the namespaces of the app (but how do we update this? how is it done for the registry key?)
- finish update and get function



func (h *EnvironmentVariablesGroup) update(r *http.Request) (*routing.Response, error) {
	authInfo, _ := authorization.InfoFromContext(r.Context())
	logger := logr.FromContextOrDiscard(r.Context()).WithName("handlers.evg.update")

	// ensure only Admin can set environment variable groups
}

func (h *EnvironmentVariablesGroup) get(r *http.Request) (*routing.Response, error) {
	authInfo, _ := authorization.InfoFromContext(r.Context())
	logger := logr.FromContextOrDiscard(r.Context()).WithName("handlers.evg.get")

	secretName := routing.URLParam(r, "guid")
}

func (h *EnvironmentVariablesGroup) UnauthenticatedRoutes() []routing.Route {
	return nil
}

func (h *EnvironmentVariablesGroup) AuthenticatedRoutes() []routing.Route {
	return []routing.Route{
		{Method: "GET", Pattern: EnvironmentVariablesGroupPath, Handler: h.get},
		{Method: "GET", Pattern: EnvironmentVariablesGroupsRoot, Handler: h.list},
		{Method: "PATCH", Pattern: EnvironmentVariablesGroupPath, Handler: h.update},
	}
}
*/
