package i18n

import "sync"

type Lang struct {
	Name    Key
	Value   string
	content map[Key]string
}

var langs = []Lang{
	{
		Name:    English,
		Value:   "en_US",
		content: en_US,
	},
	{
		Name:    Chinese,
		Value:   "zh_CN",
		content: zh_CN,
	},
}

func Langs() []Lang {
	return langs
}

var (
	currentLang Lang = langs[0]
	mux         sync.RWMutex
)

func Current() Lang {
	mux.RLock()
	defer mux.RUnlock()

	return currentLang
}

func Set(lang string) {
	mux.Lock()
	defer mux.Unlock()

	for i := range langs {
		if langs[i].Value == lang {
			currentLang = langs[i]
			return
		}
	}
	currentLang = langs[0]
}

const (
	Server     Key = "server"
	Service    Key = "service"
	Chain      Key = "chain"
	Hop        Key = "hop"
	Auther     Key = "auther"
	Admission  Key = "admission"
	Bypass     Key = "bypass"
	Resolver   Key = "resolver"
	Hosts      Key = "hosts"
	Limiter    Key = "limiter"
	Ingress    Key = "ingress"
	Observer   Key = "observer"
	Logger     Key = "logger"
	Recorder   Key = "recorder"
	Plugin     Key = "plugin"
	Selector   Key = "selector"
	Node       Key = "node"
	Nameserver Key = "nameserver"

	Type         Key = "type"
	Name         Key = "name"
	Address      Key = "address"
	Host         Key = "host"
	ServerName   Key = "serverName"
	Path         Key = "path"
	URL          Key = "url"
	URLHint      Key = "urlHint"
	Interval     Key = "interval"
	IntervalHint Key = "intervalHint"
	Timeout      Key = "timeout"
	TimeoutHint  Key = "timeoutHint"
	Seconds      Key = "seconds"
	Filter       Key = "filter"
	TLS          Key = "tls"
	HTTP         Key = "http"
	CertFile     Key = "certFile"
	KeyFile      Key = "keyFile"
	CAFile       Key = "caFile"

	MetadataKey   Key = "metadataKey"
	MetadataValue Key = "metadataValue"

	Auth      Key = "auth"
	BasicAuth Key = "basicAuth"
	Username  Key = "username"
	Password  Key = "password"

	AutoSave Key = "autoSave"
	Readonly Key = "readonly"

	Basic      Key = "basic"
	Advanced   Key = "advanced"
	AuthSimple Key = "authSimple"
	AuthAuther Key = "authAuther"

	Interface     Key = "interface"
	InterfaceHint Key = "interfaceHint"

	Handler           Key = "handler"
	Listener          Key = "listener"
	Forwarder         Key = "forwarder"
	Connector         Key = "connector"
	Dialer            Key = "dialer"
	Protocol          Key = "protocol"
	Nodes             Key = "nodes"
	Metadata          Key = "metadata"
	VerifyServerCert  Key = "verifyServerCert"
	RewriteHostHeader Key = "rewriteHostHeader"

	FileServer           Key = "fileServer"
	SerialPortRedirector Key = "serialPortRedirector"
	UnixDomainSocket     Key = "unixDomainSocket"
	ReverseProxyTunnel   Key = "reverseProxyTunnel"

	DirectoryPath  Key = "dirPath"
	CustomHostname Key = "customHostname"
	Hostname       Key = "hostname"
	EnableTLS      Key = "enableTLS"
	Keepalive      Key = "keepalive"
	Whitelist      Key = "whitelist"
	Matcher        Key = "matcher"
	Rule           Key = "rule"
	Rules          Key = "rules"
	Auths          Key = "auths"
	HostMappings   Key = "hostMappings"
	Mapping        Key = "mapping"
	HostAlias      Key = "hostAlias"
	Async          Key = "async"
	Prefer         Key = "prefer"
	Only           Key = "only"
	ClientIP       Key = "clientIP"
	Limits         Key = "limits"
	Record         Key = "record"

	DeleteServer       Key = "deleteServer"
	DeleteService      Key = "deleteService"
	DeleteChain        Key = "deleteChain"
	DeleteHop          Key = "deleteHop"
	DeleteNode         Key = "deleteNode"
	DeleteMetadata     Key = "deleteMetadata"
	DeleteAuther       Key = "deleteAuther"
	DeleteAuth         Key = "deleteAuth"
	DeleteAdmission    Key = "deleteAdmission"
	DeleteRules        Key = "deleteRules"
	DeleteBypass       Key = "deleteBypass"
	DeleteResolver     Key = "deleteResolver"
	DeleteNameserver   Key = "deleteNameserver"
	DeleteHosts        Key = "deleteHosts"
	DeleteHostMappings Key = "deleteHostMappings"
	DeleteLimiter      Key = "deleteLimiter"
	DeleteLimits       Key = "deleteLimits"
	DeleteObserver     Key = "deleteObserver"
	DeleteRecorder     Key = "deleteRecorder"
	DeleteRecord       Key = "deleteRecord"

	SelectorStrategy    Key = "selectorStrategy"
	SelectorRound       Key = "selectorRound"
	SelectorRandom      Key = "selectorRandom"
	SelectorFIFO        Key = "selectorFIFO"
	SelectorMaxFails    Key = "selectorMaxFails"
	SelectorFailTimeout Key = "selectorFailTimeout"

	DataSource       Key = "dataSource"
	FileDataSource   Key = "fileDataSource"
	FilePath         Key = "filePath"
	RedisDataSource  Key = "redisDataSource"
	RedisAddr        Key = "redisAddr"
	RedisDB          Key = "redisDB"
	RedisPassword    Key = "redisPassword"
	RedisKey         Key = "redisKey"
	RedisType        Key = "redisType"
	HTTPDataSource   Key = "httpDataSource"
	HTTPURL          Key = "httpURL"
	HTTPTimeout      Key = "httpTimeout"
	TCPDataSource    Key = "tcpDataSource"
	TCPAddr          Key = "tcpAddr"
	TCPTimeout       Key = "tcpTimeout"
	DataSourceReload Key = "dataSourceReload"
	FileSep          Key = "fileSep"

	PluginGRPC Key = "pluginGRPC"
	PluginHTTP Key = "pluginHTTP"

	TimeSecond Key = "timeSecond"

	ErrNameRequired Key = "errNameRequired"
	ErrNameExists   Key = "errNameExists"
	ErrURLRequired  Key = "errURLRequired"
	ErrInvalidAddr  Key = "errInvalidAddr"
	ErrDigitOnly    Key = "errDigitOnly"
	ErrDirectory    Key = "errDir"

	OK     Key = "ok"
	Cancel Key = "cancel"

	Running Key = "running"
	Ready   Key = "ready"
	Failed  Key = "failed"
	Closed  Key = "closed"
	Unknown Key = "unknown"

	Config   Key = "config"
	Event    Key = "event"
	Settings Key = "settings"
	Language Key = "language"
	English  Key = "english"
	Chinese  Key = "chinese"
	Theme    Key = "theme"
	Light    Key = "light"
	Dark     Key = "dark"
)

type Key string

func (k Key) Value() string {
	if k == "" {
		return ""
	}

	return get(k)
}

func get(key Key) string {
	mux.RLock()
	defer mux.RUnlock()

	if v := currentLang.content[key]; v != "" {
		return v
	}

	return langs[0].content[key]
}
