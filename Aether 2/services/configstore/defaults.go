// Services > ConfigStore
// This module handles saving and reading values from a config user file.

package configstore

import (
	"time"
)

// Defaults

const (
	defaultNetworkHeadDays                         = 14
	defaultNetworkMemoryDays                       = 180
	defaultLocalMemoryDays                         = 180
	defaultPoWBailoutTimeSeconds                   = 30
	defaultTimeBlockSizeMinutes                    = 5
	defaultPastBlocksToCheck                       = 3
	defaultCacheGenerationIntervalHours            = 6
	defaultCacheDurationHours                      = 6
	defaultPOSTResponseExpiryMinutes               = 540 // 9h
	defaultPOSTResponseIneligibilityMinutes        = 480 // 8h
	defaultConnectionTimeout                       = 60 * time.Second
	defaultTCPConnectTimeout                       = 3 * time.Second
	defaultTLSHandshakeTimeout                     = 1 * time.Second
	defaultPingerPageSize                          = 100
	defaultOnlineAddressFinderPageSize             = 99
	defaultDispatchExclusionExpiryForLiveAddress   = 5 * time.Second // this is normally minute TODO
	defaultDispatchExclusionExpiryForStaticAddress = 72 * time.Hour
	defaultPowStrength                             = 20
	defaultExternalIp                              = "0.0.0.0" // Localhost, if this is still 0.0.0.0 at any point in the future we failed at finding this out.
	defaultExternalIpType                          = 4         // IPv4
	defaultExternalPort                            = 49999
	defaultDbEngine                                = "sqlite" // 'sqlite' or 'mysql'
	defaultDBIp                                    = "127.0.0.1"
	defaultDbPort                                  = 3306
	defaultDbUsername                              = "aether-app-db-access-user"
	defaultDbPassword                              = "exventoveritas"
	defaultNeighbourCount                          = 100
	defaultMaxAddressTableSize                     = 1000
)

// Default entity page sizes

const (
	defaultBoardsPageSize      = 2000  // 0.2x
	defaultThreadsPageSize     = 400   // 1x
	defaultPostsPageSize       = 400   // 1x
	defaultVotesPageSize       = 2000  // 0.2x 2000
	defaultKeysPageSize        = 2000  // 0.2x
	defaultTruststatesPageSize = 6000  // 0.025x
	defaultAddressesPageSize   = 16000 // 0.025x

	defaultBoardIndexesPageSize      = 8000  // 0.025x
	defaultThreadIndexesPageSize     = 16000 // 0.025x
	defaultPostIndexesPageSize       = 12000 // 0.033x
	defaultVoteIndexesPageSize       = 4000  // 0.1x
	defaultKeyIndexesPageSize        = 20000 // 0.02x
	defaultTruststateIndexesPageSize = 15000 // 0.01x
	defaultAddressIndexesPageSize    = 16000 // 0.025x - Address is its own index
	// Every regular page is about 500kb that way.
	// Every index page is about 1mb.

	defaultBoardManifestsPageSize      = 30000
	defaultThreadManifestsPageSize     = 30000
	defaultPostManifestsPageSize       = 30000
	defaultVoteManifestsPageSize       = 30000
	defaultKeyManifestsPageSize        = 30000
	defaultTruststateManifestsPageSize = 30000
	defaultAddressManifestsPageSize    = 30000
	// Manifests are all the same size, so they're all the same.

)

// Hardcoded version numbers specific to this build
const (
	clientVersionMajor   = 2
	clientVersionMinor   = 0
	clientVersionPatch   = 0
	clientName           = "Aether"
	protocolVersionMajor = 1
	protocolVersionMinor = 0
)
