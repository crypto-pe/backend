package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/0xsequence/go-sequence/lib/prototyp"
	"github.com/crypto-pe/backend/data"
	"github.com/crypto-pe/backend/proto"
	"github.com/go-chi/httplog"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/jwt"
)

var (
	ErrUnauthorized   = proto.WrapError(proto.ErrUnauthenticated, nil, "unauthorized")
	ErrSessionExpired = proto.WrapError(proto.ErrPermissionDenied, nil, "session expired")
)

var (
	LoggerCtxKey = &contextKey{"Logger"}

	SessionTypeCtxKey = &contextKey{"SessionType"}
	WalletCtxKey      = &contextKey{"Wallet"} // Ethereum account address (string/hash)
	UserCtxKey        = &contextKey{"User"}   // user account object (*data.User)
	ServiceCtxKey     = &contextKey{"Service"}
)

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "context value " + k.name
}

// Session middleware to attach `account` or `service` sessions to the request context
func Session(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token, claims, tokenErr := jwtauth.FromContext(ctx)

		logEntry := httplog.LogEntry(r.Context())
		if token == nil {

			switch tokenErr {

			case jwtauth.ErrNoTokenFound:
				// When JWT token is not passed, we set to anonymous/public session and continue
				// to AccessControl middleware (which is expected to be next)
				httplog.LogEntrySetField(r.Context(), "auth", "anonymous")
				ctx = context.WithValue(ctx, SessionTypeCtxKey, SessionTypePublic)
				next.ServeHTTP(w, r.WithContext(ctx))

			case jwtauth.ErrExpired:
				logEntry.Info().Msgf("jwt expired")
				proto.RespondWithError(w, ErrSessionExpired)

			default:
				logEntry.Warn().Msgf("jwt unauthorized")
				proto.RespondWithError(w, ErrUnauthorized)
			}

			return
		}
		// When JWT token is found, ensure it verifies, or error
		if token == nil || jwt.Validate(token) != nil || tokenErr != nil {
			logEntry.Warn().Msgf("jwt unauthorized")
			proto.RespondWithError(w, ErrUnauthorized)
			return
		}

		// // Origin check
		// if originClaim, ok := claims["ogn"].(string); ok {
		// 	originHeader := r.Header.Get("Origin")
		// 	if originHeader != "" && originHeader != originClaim {
		// 		logEntry.Warn().Msgf("jwt unauthorized -- invalid origin. Request origin is %s but expecting %s", originClaim, originHeader)
		// 		proto.RespondWithError(w, fmt.Errorf("invalid origin claim: %w", ErrUnauthorized))
		// 	}
		// }

		accountClaim, _ := claims["account"].(string)
		if accountClaim != "" {
			httplog.LogEntrySetField(r.Context(), "jwtAccount", accountClaim)
			// user, _ := data.DB.Users.FindByAddress(context.Background(), prototyp.HashFromString(accountClaim))
			user, err := data.DB.GetUser(ctx, prototyp.HashFromString(accountClaim).String())
			// Wallet account address, from jwt claims
			fmt.Println("wtfWOiefnrwoaenwaoiwnefoainweofiweanfpawoefinaweonfwaerfoin?", err)
			ctx = context.WithValue(ctx, WalletCtxKey, accountClaim)

			if err == nil {
				// user account
				if user.Admin.Bool {
					httplog.LogEntrySetField(r.Context(), "auth", "admin")
					ctx = context.WithValue(ctx, SessionTypeCtxKey, SessionTypeAdmin)
				} else {
					httplog.LogEntrySetField(r.Context(), "auth", "user")
					ctx = context.WithValue(ctx, SessionTypeCtxKey, SessionTypeUser)
				}
				ctx = context.WithValue(ctx, UserCtxKey, user)
			} else {
				httplog.LogEntrySetField(r.Context(), "auth", "anonymous")
				ctx = context.WithValue(ctx, SessionTypeCtxKey, SessionTypePublic)
			}
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
