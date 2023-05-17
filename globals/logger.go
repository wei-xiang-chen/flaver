package globals

import (
	"flaver/globals/tools"

	"go.uber.org/zap"
)

func GetLogger() *zap.SugaredLogger {
	return tools.GetLogger()
}
