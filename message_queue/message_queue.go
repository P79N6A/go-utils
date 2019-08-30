package messageQueue

import (
	"fmt"
	"github.com/kilingzhang/go-utils/message_queue/driver"
	"sort"
	"sync"
)

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]driver.Driver)
)

func Register(name string, driver driver.Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if driver == nil {
		panic("message-queue: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("message-queue: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

func unregisterAllDrivers() {
	driversMu.Lock()
	defer driversMu.Unlock()
	// For tests.
	drivers = make(map[string]driver.Driver)
}

func Drivers() []string {
	driversMu.RLock()
	defer driversMu.RUnlock()
	var list []string
	for name := range drivers {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

func NewSyncProducer(driverName string, addrs []string, config interface{}) (syncProducer driver.SyncProducer, err error) {
	driversMu.RLock()
	driveri, ok := drivers[driverName]
	driversMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("sql: unknown driver %q (forgotten import?)", driverName)
	}

	return driveri.NewSyncProducer(addrs, config)
}

func NewConsumerGroup(driverName string, addrs []string, groupID string, config interface{}) (driver.ConsumerGroup, error) {
	driversMu.RLock()
	driveri, ok := drivers[driverName]
	driversMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("sql: unknown driver %q (forgotten import?)", driverName)
	}

	return driveri.NewConsumerGroup(addrs, groupID, config)
}
