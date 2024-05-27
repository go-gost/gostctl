package i18n

var en_US = map[Key]string{
	Server:     "Server",
	Service:    "Service",
	Chain:      "Chain",
	Hop:        "Hop",
	Auther:     "Auther",
	Admission:  "Admission",
	Bypass:     "Bypass",
	Resolver:   "Resolver",
	Hosts:      "Hosts",
	Limiter:    "Limiter",
	Ingress:    "Ingress",
	Observer:   "Observer",
	Logger:     "Logger",
	Recorder:   "Recorder",
	Plugin:     "Plugin",
	Selector:   "Selector",
	Node:       "Node",
	Nameserver: "Nameserver",

	Type:         "Type",
	Name:         "Name",
	Address:      "Address",
	Host:         "Host",
	ServerName:   "Server name",
	Path:         "Path",
	URL:          "URL",
	URLHint:      "e.g. http://localhost:8000",
	Interval:     "Interval",
	IntervalHint: "the period for obtaining configuration",
	Timeout:      "Timeout",
	TimeoutHint:  "request timeout when obtaining configuration",
	Seconds:      "Seconds",
	Filter:       "Filter",
	TLS:          "TLS",
	HTTP:         "HTTP",
	CertFile:     "Cert File",
	KeyFile:      "Key File",
	CAFile:       "CA File",

	MetadataKey:   "Key",
	MetadataValue: "Value",

	Auth:      "Auth",
	BasicAuth: "Auth",
	Username:  "Username",
	Password:  "Password",
	AutoSave:  "Auto save",

	Basic:         "Basic",
	Advanced:      "Advanced",
	AuthSimple:    "Single",
	AuthAuther:    "Auther",
	Interface:     "Interface",
	InterfaceHint: "IP or interface name",

	FileServer:           "File Server",
	SerialPortRedirector: "Serial Port Redirector",
	UnixDomainSocket:     "Unix Domain Socket",
	ReverseProxyTunnel:   "Tunnel",

	Handler:           "Handler",
	Listener:          "Listener",
	Forwarder:         "Forwarder",
	Connector:         "Connector",
	Dialer:            "Dialer",
	Protocol:          "Protocol",
	VerifyServerCert:  "Verify server's certificate",
	Nodes:             "Nodes",
	Metadata:          "Metadata",
	RewriteHostHeader: "Rewrite host header",

	DeleteServer:       "Delete server?",
	DeleteService:      "Delete service?",
	DeleteChain:        "Delete chain?",
	DeleteHop:          "Delete hop?",
	DeleteNode:         "Delete node?",
	DeleteMetadata:     "Delete metadata?",
	DeleteAuther:       "Delete auther?",
	DeleteAuth:         "Delete auth?",
	DeleteAdmission:    "Delete admission?",
	DeleteRules:        "Delete rules?",
	DeleteBypass:       "Delete bypass?",
	DeleteResolver:     "Delete resolver?",
	DeleteNameserver:   "Delete nameserver?",
	DeleteHosts:        "Delete hosts?",
	DeleteHostMappings: "Delete host mappings?",
	DeleteLimiter:      "Delete limiter?",
	DeleteLimits:       "Delete limits?",
	DeleteObserver:     "Delete observer?",
	DeleteRecorder:     "Delete recorder?",

	SelectorStrategy: "Strategy",
	SelectorRound:    "Round-Robin",
	SelectorRandom:   "Random",
	SelectorFIFO:     "HA",

	DataSource:       "Data Source",
	FileDataSource:   "File",
	FilePath:         "File path",
	RedisDataSource:  "Redis",
	RedisAddr:        "Address",
	RedisDB:          "Database",
	RedisPassword:    "Password",
	RedisKey:         "Key",
	RedisType:        "Type",
	HTTPDataSource:   "HTTP",
	HTTPURL:          "URL",
	HTTPTimeout:      "Timeout",
	TCPDataSource:    "TCP",
	TCPAddr:          "Address",
	TCPTimeout:       "Timeout",
	DataSourceReload: "Reload period",
	FileSep:          "Sep",

	PluginGRPC: "gRPC",
	PluginHTTP: "HTTP",

	TimeSecond: "s",

	DirectoryPath:  "Directory path",
	CustomHostname: "Custom hostname (rewrite HTTP Host header)",
	Hostname:       "Hostname",
	EnableTLS:      "Enalbe TLS",
	Keepalive:      "Keepalive",
	Whitelist:      "Whitelist",
	Matcher:        "Matcher",
	Rule:           "Rule",
	Rules:          "Rules",
	Auths:          "Auths",
	HostMappings:   "Host mappings",
	Mapping:        "Mapping",
	HostAlias:      "Alias",
	Async:          "Async",
	Prefer:         "Prefer",
	Only:           "Only",
	ClientIP:       "Client IP",
	Limits:         "Limits",

	ErrNameRequired: "name is required",
	ErrNameExists:   "name exists",
	ErrURLRequired:  "URL is required",
	ErrInvalidAddr:  "invalid address format, should be [IP]:PORT or [HOST]:PORT",
	ErrDigitOnly:    "Must contain only digits",
	ErrDirectory:    "is not a directory",

	OK:     "OK",
	Cancel: "Cancel",

	Running: "Running",
	Ready:   "Ready",
	Failed:  "Failed",
	Closed:  "Closed",
	Unknown: "Unknown",

	Settings: "Settings",
	Config:   "Config",
	Language: "Language",
	English:  "English",
	Chinese:  "Chinese",
	Theme:    "Theme",
	Light:    "Light",
	Dark:     "Dark",
}
