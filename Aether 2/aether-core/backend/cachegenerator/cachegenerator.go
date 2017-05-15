// Backend > CacheGenerator
// This package provides a main method, when given a time range and an entity, will generate a full-fledged cache for that entity within the given timeframe.

package cachegenerator

import (
	"aether-core/io/api"
	"aether-core/services/globals"
)

// GenerateCache is the high level API serving the outside world. Given an entity type, a start and an end, it will generate the cache and save it to the proper location for access.
func GenerateCache(etype string, start api.Timestamp, end api.Timestamp) {

}
