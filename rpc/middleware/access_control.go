package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/crypto-pe/backend/proto"
	"github.com/go-chi/httplog"
)

type SessionType uint

// SECURITY NOTE: the order of this list is important, as there is a security
// hierarchy, greater number gets higher privilege
const (
	SessionTypeUnknown SessionType = iota // 0
	SessionTypePublic                     // 1
	SessionTypeUser                       // 2
	SessionTypeAdmin                      // 3
)

var serviceNames = []string{"API"} // /rpc/{serviceName}/*

// accessMap for /rpc/API/_METHODS_
var (
	accessMap = map[string][]SessionType{
		//
		"Ping":    {SessionTypePublic, SessionTypeUser, SessionTypeAdmin},
		"Version": {SessionTypePublic, SessionTypeUser, SessionTypeAdmin},

		"GetSupportedTokens": {SessionTypePublic, SessionTypeUser, SessionTypeAdmin},

		"CreateAccount": {SessionTypePublic},
		"Login":         {SessionTypePublic},

		"GetAccount":                {SessionTypeUser, SessionTypeAdmin},
		"UpdateAccount":             {SessionTypeUser, SessionTypeAdmin},
		"DeleteAccount":             {SessionTypeUser, SessionTypeAdmin},
		"CreateOrganization":        {SessionTypeUser, SessionTypeAdmin},
		"GetOrganization":           {SessionTypeUser, SessionTypeAdmin},
		"UpdateOrganization":        {SessionTypeUser, SessionTypeAdmin},
		"DeleteOrganization":        {SessionTypeUser, SessionTypeAdmin},
		"GetAllOrganizations":       {SessionTypeUser, SessionTypeAdmin},
		"CreateOrganizationMember":  {SessionTypeUser, SessionTypeAdmin},
		"GetOrganizationMember":     {SessionTypeUser, SessionTypeAdmin},
		"GetAllOrganizationMembers": {SessionTypeUser, SessionTypeAdmin},
		"UpdateOrganizationMember":  {SessionTypeUser, SessionTypeAdmin},
		"DeleteOrganizationMember":  {SessionTypeUser, SessionTypeAdmin},
	}
	accessMapACL map[string]map[SessionType]bool
)

func init() {
	// convert accessMap into more searchable format
	accessMapACL = make(map[string]map[SessionType]bool)
	for method, perms := range accessMap {
		accessMapACL[method] = make(map[SessionType]bool)
		for _, perm := range perms {
			accessMapACL[method][perm] = true
		}
	}
}

func AccessControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logEntry := httplog.LogEntry(r.Context())

		webrpcReq, err := newRequest(r.URL.Path)
		if err != nil {
			logEntry.Warn().Msg("invalid rpc method called")
			proto.RespondWithError(w, ErrUnauthorized)
			return
		}
		httplog.LogEntrySetField(r.Context(), "rpc_service", webrpcReq.ServiceName)
		httplog.LogEntrySetField(r.Context(), "rpc_method", webrpcReq.MethodName)

		err = webrpcReq.AuthorizeRequest(r)
		if err != nil {
			logEntry.Warn().Msg("unauthorized request")
			proto.RespondWithError(w, ErrUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type webrpcRequest struct {
	PackageName string
	ServiceName string
	MethodName  string
}

func newRequest(path string) (*webrpcRequest, error) {
	p1 := strings.Split(path, "/")
	if len(p1) != 4 {
		return nil, errors.New("access: unexpected method")
	}
	t := &webrpcRequest{
		PackageName: p1[1],
		ServiceName: p1[2],
		MethodName:  p1[3],
	}
	if t.PackageName == "" || t.ServiceName == "" || t.MethodName == "" {
		return nil, errors.New("access: unexpected method")
	}
	return t, nil
}

func (t *webrpcRequest) AuthorizeRequest(r *http.Request) error {
	if t.PackageName != "rpc" {
		return ErrUnauthorized
	}

	serviceOk := false
	for _, s := range serviceNames {
		if t.ServiceName == s {
			serviceOk = true
			break
		}
	}
	if !serviceOk {
		return ErrUnauthorized
	}

	// get method's ACL
	perms, ok := accessMapACL[t.MethodName]
	if !ok {
		// unable to find method in rules list. deny.
		return ErrUnauthorized
	}

	ctx := r.Context()
	sessionType, ok := ctx.Value(SessionTypeCtxKey).(SessionType)
	if !ok {
		// should never happen as previous middleware will set it,
		// but lets be specific
		return ErrUnauthorized
	}

	// authorize using methods's ACL
	if !perms[sessionType] {
		return ErrUnauthorized
	}

	return nil
}
