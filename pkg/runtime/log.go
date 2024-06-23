package runtime

import "go.uber.org/zap"

var (
	z, _ = zap.NewProduction()
	log  = z.Sugar()
)
