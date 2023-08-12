package api

import (
	"errors"
	"rhymald/mag-zeta/play"
	// "rhymald/mag-zeta/base"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	// "rhymald/mag-zeta/base"
	// "go.opentelemetry.io/otel/trace"
	// "fmt"
	// "math"
	// "rhymald/mag-zeta/connect"
	"fmt"
)

// func npcRegen(hps *base.Life, ids *map[string]int, span *trace.Span) {
// 	hp := 32
// 	hps.HealDamage(hp)
// 	(*ids)["Life"] = base.Epoch()
// 	(*span).AddEvent(fmt.Sprintf("0|+0[none|-1]+HP|%+d", hp))
// }
