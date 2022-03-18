package monitoring

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/getsentry/sentry-go"
)

func SendToSentry(tags map[string]string, extras map[string]interface{}, errCategory interface{}) {
	isSentryActive, errEnv := strconv.ParseBool(os.Getenv("SENTRY_ACTIVE"))
	if isSentryActive && errEnv == nil {
		sentry.Init(sentry.ClientOptions{
			Dsn:         os.Getenv("SENTRY_DSN"),
			Environment: os.Getenv("APP_ENV"),
			Debug:       false,
			AttachStacktrace: true,
		})
		defer sentry.Flush(2 * time.Second)

		sentry.ConfigureScope(func(scope *sentry.Scope) {
			if len(tags) > 0 {
				for k, v := range tags {
					scope.SetTag(k, v)
				}
			}
			if len(extras) > 0 {
				for k, v := range extras {
					scope.SetExtra(k, v)
				}
			}

			title := fmt.Sprintf("%s | %s", os.Getenv("APP_NAME"), errCategory)
			sentry.CaptureMessage(title)
		})
	}
}

