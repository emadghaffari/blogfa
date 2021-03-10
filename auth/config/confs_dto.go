package config

import (
	zapLogger "blogfa/auth/pkg/logger"
	"encoding/json"

	"go.uber.org/zap"
)

// Set method
// you can set new key in switch for manage config with config server
func (g *global) Set(key string, query []byte) error {
	logger := zapLogger.GetZapLogger(false)
	if err := json.Unmarshal(query, &Global); err != nil {
		zapLogger.Prepare(logger).
			Append(zap.Any("key", key)).
			Append(zap.Any("value", string(query))).
			Development().
			Level(zap.ErrorLevel).
			Commit(err.Error())
		return err
	}
	
	return nil
}
