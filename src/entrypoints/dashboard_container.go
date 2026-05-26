package entrypoints

import dashuc "github.com/educabot/alizia-inclusion-be/src/core/usecases/dashboard"

type DashboardContainer struct {
	GetMetrics dashuc.GetMetrics
	GetAIUsage dashuc.GetAIUsage
}
