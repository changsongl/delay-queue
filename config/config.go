package config

import (
	"errors"
	"github.com/changsongl/delay-queue/config/decode"
	"github.com/changsongl/delay-queue/config/decode/json"
	"github.com/changsongl/delay-queue/config/decode/yaml"
	"io/ioutil"
	"os"
	"time"
)

type FileType string

const (
	FileTypeYaml FileType = "yaml"
	FileTypeJson FileType = "json"
)

// configuration
type Conf struct {
	DelayQueue DelayQueue `yaml:"delay_queue,omitempty" json:"delay_queue,omitempty"`
	Redis      Redis      `yaml:"redis,omitempty" json:"redis,omitempty"`
}

// delay queue configuration
type DelayQueue struct {
	BucketSize int `yaml:"bucket_size,omitempty" json:"bucket_size,omitempty"`
}

// redis configuration
type Redis struct {
	// The network type, either tcp or unix.
	// Default is tcp.
	Network string `yaml:"network,omitempty" json:"network,omitempty"`

	// host:port address.
	Addr string `yaml:"addr,omitempty" json:"addr,omitempty"`

	// Use the specified username to authenticate the current connection
	// with one of the connections defined in the ACL list when connecting
	// to a Redis 6.0 instance, or greater, that is using the Redis ACL system.
	Username string `yaml:"username,omitempty" json:"username,omitempty"`

	// Optional password. Must match the password specified in the
	// requirepass server configuration option (if connecting to a Redis 5.0 instance, or lower),
	// or the User password when connecting to a Redis 6.0 instance, or greater,
	// that is using the Redis ACL system.
	Password string `yaml:"password,omitempty" json:"password,omitempty"`

	// Database to be selected after connecting to the server.
	DB int `yaml:"db,omitempty" json:"db,omitempty"`

	// Maximum number of retries before giving up.
	// Default is 3 retries.
	MaxRetries int `yaml:"max_retries,omitempty" json:"max_retries,omitempty"`
	// Minimum backoff between each retry.
	// Default is 8 milliseconds; -1 disables backoff.
	MinRetryBackoff time.Duration `yaml:"min_retry_backoff,omitempty" json:"min_retry_backoff,omitempty"`
	// Maximum backoff between each retry.
	// Default is 512 milliseconds; -1 disables backoff.
	MaxRetryBackoff time.Duration `yaml:"max_retry_backoff,omitempty" json:"max_retry_backoff,omitempty"`

	// Dial timeout for establishing new connections.
	// Default is 5 seconds.
	DialTimeout time.Duration `yaml:"dial_timeout,omitempty" json:"dial_timeout,omitempty"`
	// Timeout for socket reads. If reached, commands will fail
	// with a timeout instead of blocking. Use value -1 for no timeout and 0 for default.
	// Default is 3 seconds.
	ReadTimeout time.Duration `yaml:"read_timeout,omitempty" json:"read_timeout,omitempty"`
	// Timeout for socket writes. If reached, commands will fail
	// with a timeout instead of blocking.
	// Default is ReadTimeout.
	WriteTimeout time.Duration `yaml:"write_timeout,omitempty" json:"write_timeout,omitempty"`

	// Maximum number of socket connections.
	// Default is 10 connections per every CPU as reported by runtime.NumCPU.
	PoolSize int `yaml:"pool_size,omitempty" json:"pool_size,omitempty"`
	// Minimum number of idle connections which is useful when establishing
	// new connection is slow.
	MinIdleConns int `yaml:"min_idle_conns,omitempty" json:"min_idle_conns,omitempty"`
	// Connection age at which client retires (closes) the connection.
	// Default is to not close aged connections.
	MaxConnAge time.Duration `yaml:"max_conn_age,omitempty" json:"max_conn_age,omitempty"`
	// Amount of time client waits for connection if all connections
	// are busy before returning an error.
	// Default is ReadTimeout + 1 second.
	PoolTimeout time.Duration `yaml:"pool_timeout,omitempty" json:"pool_timeout,omitempty"`
	// Amount of time after which client closes idle connections.
	// Should be less than server's timeout.
	// Default is 5 minutes. -1 disables idle timeout check.
	IdleTimeout time.Duration `yaml:"idle_timeout,omitempty" json:"idle_timeout,omitempty"`
	// Frequency of idle checks made by idle connections reaper.
	// Default is 1 minute. -1 disables idle connections reaper,
	// but idle connections are still discarded by the client
	// if IdleTimeout is set.
	IdleCheckFrequency time.Duration `yaml:"idle_check_frequency,omitempty" json:"idle_check_frequency,omitempty"`
}

func New() *Conf {
	return &Conf{}
}

func (c *Conf) Load(file string, fileType FileType) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}

	bts, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	decoder, err := c.getDecoderByFileType(fileType)
	if err != nil {
		return err
	}

	err = c.load(bts, decoder.DecodeFunc())
	if err != nil {
		return err
	}

	return nil
}

func (c *Conf) load(bts []byte, decodeFunc func([]byte, interface{}) error) error {
	err := decodeFunc(bts, c)
	if err != nil {
		return nil
	}
	return nil
}

func (c *Conf) getDecoderByFileType(fileType FileType) (decode.Decoder, error) {
	if fileType == FileTypeJson {
		return json.NewDecoder(), nil
	} else if fileType == FileTypeYaml {
		return yaml.NewDecoder(), nil
	}

	return nil, errors.New("invalid file type")
}
