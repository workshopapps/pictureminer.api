package middleware

import (
	"fmt"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/workshopapps/pictureminer.api/internal/config"
)

func SentryLogger() {
	getConfig := config.GetConfig()
	_ = sentry.Init(sentry.ClientOptions{
		Dsn: getConfig.Sentry.Dsn,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			if hint.Context != nil {
				if req, ok := hint.Context.Value(sentry.RequestContextKey).(*http.Request); ok {
					// You have access to the original Request
					fmt.Println(req)
				}
			}
			fmt.Println(event)
			return event
		},
		Debug:            true,
		AttachStacktrace: true,
	})

}
