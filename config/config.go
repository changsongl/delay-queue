package config

import (
	encodeJson "encoding/json"
	"errors"
	"github.com/changsongl/delay-queue/config/decode"
	"github.com/changsongl/delay-queue/config/decode/json"
	"github.com/changsongl/delay-queue/config/decode/yaml"
	"io/ioutil"
	"os"
)

// FileType configuration file type
type FileType string

// config file type enum
const (
	// FileTypeYAML yaml type
	FileTypeYAML FileType = "yaml"
	// FileTypeJSON json type
	FileTypeJSON FileType = "json"
)

// RedisMode redis set up
type RedisMode string

const (
	// RedisModeNormal redis normal mode
	RedisModeNormal RedisMode = ""
	// RedisModeCluster redis cluster
	RedisModeCluster RedisMode = "cluster"
)

// default configurations
const (
	// delay queue configuration
	DefaultDQBindAddress       = ":8000"
	DefaultDQBucketName        = "dq_bucket"
	DefaultDQQueueName         = "dq_queue"
	DefaultDQBucketSize        = 8
	DefaultDQBucketMaxFetchNum = 200
	DefaultTimerFetchInterval  = 10000
	DefaultFetchJobBlockTime   = 5

	// redis configuration
	DefaultRedisNetwork      = "tcp"
	DefaultRedisAddress      = "127.0.0.1:6379"
	DefaultRedisDialTimeout  = 5000
	DefaultRedisReadTimeout  = 3000
	DefaultRedisWriteTimeout = 3000
)

// Conf configuration
type Conf struct {
	DelayQueue DelayQueue `yaml:"delay_queue,omitempty" json:"delay_queue,omitempty"`
	Redis      Redis      `yaml:"redis,omitempty" json:"redis,omitempty"`
}

// TODO: add duration to configs

// DelayQueue delay queue configuration
type DelayQueue struct {
	// listen address
	BindAddress string `yaml:"bind_address,omitempty" json:"bind_address,omitempty"`
	// bucket redis key name
	BucketName string `yaml:"bucket_name,omitempty" json:"bucket_name,omitempty"`
	// the number of delay queue bucket, increase number could get better concurrency.
	BucketSize uint64 `yaml:"bucket_size,omitempty" json:"bucket_size,omitempty"`
	// max fetch number of jobs in the bucket
	BucketMaxFetchNum uint64 `yaml:"bucket_max_fetch_num,omitempty" json:"bucket_max_fetch_num,omitempty"`
	// queue redis key name
	QueueName string `yaml:"queue_name,omitempty" json:"queue_name,omitempty"`
	// fetching job interval(ms), decrease interval may get better throughout.
	TimerFetchInterval int `yaml:"timer_fetch_interval,omitempty" json:"timer_fetch_interval,omitempty"`
	// fetch delay(ms), if there are still job in the bucket after the fetch
	// it will delay timer_fetch_delay ms for next fetch. Default is not wait.
	TimerFetchDelay int `yaml:"timer_fetch_delay,omitempty" json:"timer_fetch_delay,omitempty"`
}

// Redis redis configuration
type Redis struct {
	// Redis set up, "" for regular redis, "cluster" for redis cluster
	Mode RedisMode `yaml:"mode,omitempty" json:"mode,omitempty"`
	// The network type, either tcp or unix.
	// Default is tcp.
	Network string `yaml:"network,omitempty" json:"network,omitempty"`

	// host:port address.
	Address string `yaml:"address,omitempty" json:"address,omitempty"`

	// Use the specified username to authenticate the current connection
	// with one of the connections defined in the ACL list when connecting
	// to a Redis 6.0 instance, or greater, that is using the Redis ACL system.
	Username string `yaml:"username,omitempty" json:"username,omitempty"`

	// Optional password. Must match the password specified in the
	// require pass server configuration option (if connecting to a Redis 5.0 instance, or lower),
	// or the User password when connecting to a Redis 6.0 instance, or greater,
	// that is using the Redis ACL system.
	Password string `yaml:"password,omitempty" json:"password,omitempty"`

	// Database to be selected after connecting to the server.
	DB int `yaml:"db,omitempty" json:"db,omitempty"`

	// Dial timeout for establishing new connections.
	// Default is 5 seconds.
	DialTimeout int `yaml:"dial_timeout,omitempty" json:"dial_timeout,omitempty"`
	// Timeout for socket reads. If reached, commands will fail
	// with a timeout instead of blocking. Use value -1 for no timeout and 0 for default.
	// Default is 3 seconds.
	ReadTimeout int `yaml:"read_timeout,omitempty" json:"read_timeout,omitempty"`
	// Timeout for socket writes. If reached, commands will fail
	// with a timeout instead of blocking.
	// Default is ReadTimeout.
	WriteTimeout int `yaml:"write_timeout,omitempty" json:"write_timeout,omitempty"`

	// Maximum number of socket connections.
	// Default is 10 connections per every CPU as reported by runtime.NumCPU.
	PoolSize int `yaml:"pool_size,omitempty" json:"pool_size,omitempty"`
	// Minimum number of idle connections which is useful when establishing
	// new connection is slow.
	MinIdleConns int `yaml:"min_idle_conns,omitempty" json:"min_idle_conns,omitempty"`
}

// New Conf instance
func New() *Conf {
	return &Conf{
		DelayQueue: DelayQueue{
			BindAddress:        DefaultDQBindAddress,
			BucketName:         DefaultDQBucketName,
			BucketSize:         DefaultDQBucketSize,
			QueueName:          DefaultDQQueueName,
			BucketMaxFetchNum:  DefaultDQBucketMaxFetchNum,
			TimerFetchInterval: DefaultTimerFetchInterval,
		},
		Redis: Redis{
			Network:      DefaultRedisNetwork,
			Address:      DefaultRedisAddress,
			DialTimeout:  DefaultRedisDialTimeout,
			ReadTimeout:  DefaultRedisReadTimeout,
			WriteTimeout: DefaultRedisWriteTimeout,
		},
	}
}

// Load configuration
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

// load the real method
func (c *Conf) load(bts []byte, decodeFunc func([]byte, interface{}) error) error {
	err := decodeFunc(bts, c)
	if err != nil {
		return nil
	}
	return nil
}

// getDecoderByFileType get file type for decoding
func (c *Conf) getDecoderByFileType(fileType FileType) (decode.Decoder, error) {
	if fileType == FileTypeJSON {
		return json.NewDecoder(), nil
	} else if fileType == FileTypeYAML {
		return yaml.NewDecoder(), nil
	}

	return nil, errors.New("invalid file type")
}

// String config string
func (c *Conf) String() string {
	bytes, _ := encodeJson.Marshal(c)
	return string(bytes)
}

// IsCluster return current redis mode is cluster
func (r RedisMode) IsCluster() bool {
	return r == RedisModeCluster
}
