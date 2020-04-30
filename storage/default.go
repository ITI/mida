package storage

const (
	// Output Parameters
	DefaultOutputPath             = "results"
	DefaultFileSubdir             = "files"
	DefaultScriptSubdir           = "scripts"
	DefaultCoverageSubdir         = "coverage"
	DefaultCrawlMetadataFile      = "metadata.json"
	DefaultResourceMetadataFile   = "resource_metadata.json"
	DefaultScriptMetadataFile     = "script_metadata.json"
	DefaultJSTracePath            = "js_trace.json"
	DefaultResourceTreePath       = "resource_tree.json"
	DefaultWebSocketTrafficFile   = "websocket_data.json"
	DefaultEventSourceDataFile    = "event_source_data.json"
	DefaultBrowserLogFileName     = "browser.log"
	DefaultNetworkStraceFileName  = "network.strace"
	DefaultScreenShotFileName     = "screenshot.png"
	MongoStorageTimeoutSeconds    = 90
	MongoStorageJSBufferLen       = 10000
	MongoStorageResourceBufferLen = 100
	TempDir                       = ".tmp"
	MaxInt64                      = 9223372036854775807

	DefaultPostgresPort = 54330
	DefaultPos
)
