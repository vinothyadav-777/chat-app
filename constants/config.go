package constants

const (
	ConfigIDKey              = "id"
	ConfigRegionKey          = "region"
	ConfigAppKey             = "app"
	ConfigEnvKey             = "env"
	ConfigTypeKey            = "configType"
	ConfigNamesKey           = "configNames"
	ConfigDirectoryKey       = "configsDirectory"
	ConfigCredentialsModeKey = "credentialsMode"
	ConfigSecretsNamesKey    = "secretNames"
	SecretDirectoryKey       = "secretsDirectory"
	SecretNamesKey           = "secretNames"
	SecretTypeKey            = "secretType"
)

// config names
const (
	APIConfig          = "api"
	LoggerConfig       = "logger"
	ApplicationConfig  = "application"
	DatabaseConfig     = "database"
)

//application config names
const (
	ConsumerBufferLength            = "consumer_buffer_length"
	MessageProcessingTimeoutInMilli = "message_processing_timeout_in_milli"
)

// database config key names
const (
	HistorySaveRecordsDBCallTimeout = "history_save_records_db_call_timeout"
)
