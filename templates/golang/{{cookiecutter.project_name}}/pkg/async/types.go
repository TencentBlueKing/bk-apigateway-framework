package async

import (
	"sync"

	"github.com/robfig/cron/v3"
)

type entry struct {
	id   cron.EntryID
	name string
}

// 任务 ID 与 cron.EntryID + 任务名称的映射表
type taskEntryMap struct {
	mapping map[int64]entry
	sync.RWMutex
}

func (m *taskEntryMap) get(taskID int64) (entry, bool) {
	m.RLock()
	defer m.RUnlock()
	e, ok := m.mapping[taskID]
	return e, ok
}

func (m *taskEntryMap) set(taskID int64, entry entry) {
	m.Lock()
	defer m.Unlock()
	m.mapping[taskID] = entry
}

func (m *taskEntryMap) delete(taskID int64) {
	m.Lock()
	defer m.Unlock()
	delete(m.mapping, taskID)
}
