package infra

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/itomofumi/go-gin-xorm-starter/util"
)

const (
	kvsExpireSecondEnv   string = "KVS_EXPIRE_SECOND"
	defaultExpireSeconds uint   = 300
)

// KVSClientInterface is key-value store interface.
type KVSClientInterface interface {
	SetStruct(key string, structPtr interface{}) error
	GetStruct(key string, structPtr interface{}) error
}

// KVSClient is key-value store client.
type KVSClient struct {
	Conn          redis.Conn
	namespace     string
	expireSeconds uint
	done          chan struct{}
}

// NewKVSClient initializes key-value store client.
func NewKVSClient() *KVSClient {
	client := &KVSClient{}
	client.namespace = os.Getenv("KVS_NAMESPACE")

	secStr := os.Getenv(kvsExpireSecondEnv)
	if sec, err := strconv.Atoi(secStr); err == nil {
		client.expireSeconds = uint(sec)
	} else {
		client.expireSeconds = 300
		logger := util.GetLogger()
		logger.Infof("%s expects uint value, but %q was given. use default %v [sec]",
			kvsExpireSecondEnv, secStr, defaultExpireSeconds)

	}

	// first try
	c, err := connect()
	if err == nil {
		client.Conn = c
	}
	client.runConnect()
	return client
}

// Close closes connection.
func (kc *KVSClient) Close() {
	if kc.Conn != nil {
		close(kc.done)
	}
}

// SetStruct store go struct object by key.
func (kc *KVSClient) SetStruct(key string, structPtr interface{}) error {
	if !kc.isConnected() {
		return fmt.Errorf("not connected")
	}
	b, err := json.Marshal(structPtr)
	if err != nil {
		return err
	}
	kc.Conn.Send("MULTI")
	kc.Conn.Send("SET", kc.namespace+key, string(b))
	kc.Conn.Send("EXPIRE", kc.namespace+key, kc.expireSeconds)
	_, err = kc.Conn.Do("EXEC")
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// GetStruct load go struct object by key.
func (kc *KVSClient) GetStruct(key string, structPtr interface{}) error {
	if !kc.isConnected() {
		return fmt.Errorf("not connected")
	}

	str, err := redis.String(kc.Conn.Do("GET", kc.namespace+key))
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = json.Unmarshal([]byte(str), structPtr)
	if err != nil {
		return err
	}

	return nil
}

func (kc *KVSClient) isConnected() bool {
	return kc.Conn != nil
}

func (kc *KVSClient) runConnect() {
	if kc.done != nil {
		return
	}

	go func() {
		for {
			select {
			case <-kc.done:
				if kc.Conn != nil {
					kc.Conn.Close()
					kc.Conn = nil
				}
				return
			default:
				if !kc.isConnected() {
					c, err := connect()
					if err == nil {
						kc.Conn = c
					} else {
						fmt.Println(err)
					}
				}
				time.Sleep(time.Second * 5)
			}
		}
	}()
}

func connect() (redis.Conn, error) {
	host := os.Getenv("KVS_HOST")
	if host == "" {
		return nil, fmt.Errorf("KVS_HOST is not set")
	}
	c, err := redis.Dial("tcp", host)
	if err != nil {
		return nil, err
	}
	fmt.Println("[success] connect to redis")
	return c, nil
}
