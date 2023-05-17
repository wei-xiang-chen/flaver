package tools

type Config struct {
	Postgresql Postgresql `mapstructure:"postgresql" json:"postgresql" yaml:"postgresql"`
	Server     Server     `mapstructure:"server" json:"server" yaml:"server"`
	Zap        Zap        `mapstructure:"zap" json:"zap" yaml:"zap"`
	Gcs        Gcs        `mapstructure:"gcs" json:"gcs" yaml:"gcs"`
	Firebase   Firebase   `mapstructure:"firebase" json:"firebase" yaml:"firebase"`
	Jwt        Jwt        `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
}

func (this Config) GetPostgresql() Postgresql {
	return this.Postgresql
}

func (this Config) GetServer() Server {
	return this.Server
}

func (this Config) GetZap() Zap {
	return this.Zap
}

func (this Config) GetGcs() Gcs {
	return this.Gcs
}

func (this Config) GetFirebase() Firebase {
	return this.Firebase
}

type Server struct {
	Env  string `mapstructure:"env" json:"env" yaml:"env"`
	Addr int    `mapstructure:"addr" json:"addr" yaml:"addr"`
}

func (this Server) GetEnv() string {
	return this.Env
}

func (this Server) GetAddr() int {
	return this.Addr
}

type Postgresql struct {
	Host                 string `mapstructure:"host" json:"host" yaml:"host"`
	Port                 string `mapstructure:"port" json:"port" yaml:"port"`
	Config               string `mapstructure:"config" json:"config" yaml:"config"`
	Dbname               string `mapstructure:"db-name" json:"dbname" yaml:"db-name"`
	Username             string `mapstructure:"username" json:"username" yaml:"username"`
	Password             string `mapstructure:"password" json:"password" yaml:"password"`
	MaxIdleConns         int    `mapstructure:"max-idle-conns" json:"maxIdleConns" yaml:"max-idle-conns"`
	MaxOpenConns         int    `mapstructure:"max-open-conns" json:"maxOpenConns" yaml:"max-open-conns"`
	MaxIdleTime          int    `mapstructure:"max-idle-time" json:"maxIdleTime" yaml:"max-idle-time"`
	MaxLifeTime          int    `mapstructure:"max-life-time" json:"maxLifeTime" yaml:"max-life-time"`
	PreferSimpleProtocol bool   `mapstructure:"prefer-simple-protocol" json:"preferSimpleProtocol" yaml:"prefer-simple-protocol"`
	Logger               bool   `mapstructure:"logger" json:"logger" yaml:"logger" default:"false"`
}

func (this Postgresql) GetPreferSimpleProtocol() bool {
	return this.PreferSimpleProtocol
}

func (this Postgresql) GetLogger() bool {
	return this.Logger
}

func (this Postgresql) GetMaxIdleConns() int {
	return this.MaxIdleConns
}

func (this Postgresql) GetMaxOpenConns() int {
	return this.MaxOpenConns
}

func (this Postgresql) GetMaxIdleTime() int {
	return this.MaxIdleTime
}

func (this Postgresql) GetMaxLifeTime() int {
	return this.MaxLifeTime
}

func (this Postgresql) GetHost() string {
	return this.Host
}

func (this Postgresql) GetPort() string {
	return this.Port
}

func (this Postgresql) GetConfig() string {
	return this.Config
}

func (this Postgresql) GetDBname() string {
	return this.Dbname
}

func (this Postgresql) GetUsername() string {
	return this.Username
}

func (this Postgresql) GetPassword() string {
	return this.Password
}

type Zap struct {
	Level         string `mapstructure:"level" json:"level" yaml:"level"`
	Format        string `mapstructure:"format" json:"format" yaml:"format"`
	Prefix        string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
	Director      string `mapstructure:"director" json:"director"  yaml:"director"`
	LinkName      string `mapstructure:"link-name" json:"linkName" yaml:"link-name"`
	ShowLine      bool   `mapstructure:"show-line" json:"showLine" yaml:"showLine"`
	EncodeLevel   string `mapstructure:"encode-level" json:"encodeLevel" yaml:"encode-level"`
	StacktraceKey string `mapstructure:"stacktrace-key" json:"stacktraceKey" yaml:"stacktrace-key"`
	LogInConsole  bool   `mapstructure:"log-in-console" json:"logInConsole" yaml:"log-in-console"`
}

func (this Zap) GetLevel() string {
	return this.Level
}

func (this Zap) GetFormat() string {
	return this.Format
}

func (this Zap) GetPrefix() string {
	return this.Prefix
}

func (this Zap) GetDirector() string {
	return this.Director
}

func (this Zap) GetLinkName() string {
	return this.LinkName
}

func (this Zap) GetShowLine() bool {
	return this.ShowLine
}

func (this Zap) GetEncodeLevel() string {
	return this.EncodeLevel
}

func (this Zap) GetStacktraceKey() string {
	return this.StacktraceKey
}

func (this Zap) GetLogInConsole() bool {
	return this.LogInConsole
}

type Gcs struct {
	BucketName         string `mapstructure:"bucket_name" json:"bucket_name" yaml:"bucket_name"`
	CredentialFilePath string `mapstructure:"credential_file_path" json:"credential_file_path" yaml:"credential_file_path"`
	BaseUrl            string `mapstructure:"base_url" json:"base_url" yaml:"base_url"`
}

func (this Gcs) GetBucketName() string {
	return this.BucketName
}

func (this Gcs) GetCredentialFilePath() string {
	return this.CredentialFilePath
}

func (this Gcs) GetBaseUrl() string {
	return this.BaseUrl
}

type Firebase struct {
	CredentialFilePath string `mapstructure:"credential_file_path" json:"credential_file_path" yaml:"credential_file_path"`
}

func (this Firebase) GetCredentialFilePath() string {
	return this.CredentialFilePath
}

type Jwt struct {
	Secret                        string `mapstructure:"secret" json:"secret" yaml:"secret"`
	AccessTokenExpiredTimeInHour  string `mapstructure:"access_token_expired_time_in_hour" json:"access_token_expired_time_in_hour" yaml:"access_token_expired_time_in_hour"`
	RefreshTokenExpiredTimeInHour string `mapstructure:"refresh_token_expired_time_in_hour" json:"refresh_token_expired_time_in_hour" yaml:"refresh_token_expired_time_in_hour"`
}
