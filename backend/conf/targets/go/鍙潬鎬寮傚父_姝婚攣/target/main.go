package main

import (
	"sync"
)

// --------------------------------------- 数据库OP（go协程） Begin ---------------------------------------
type DaMgr struct {
	rwLock sync.RWMutex
	dbData *DbData
}

var daMgr DaMgr

// TransBegin ... 事务访问加锁
func TransBegin() {
	m := getDaMgr()
	m.rLock()
}

// TransEnd ... 事务访问解锁
func TransEnd() {
	m := getDaMgr()
	m.rUnlock()
}

// GetBoardChips ...
func GetBoardChips(addr BoardAddr) []Chip {
	m := getDaMgr()

	m.rLock()
	chips := m.getBoardChips(addr)
	m.rUnlock()

	return chips
}

func getDaMgr() *DaMgr {
	return &daMgr
}

func (m *DaMgr) rLock() {
	m.rwLock.RLock()
}

func (m *DaMgr) rUnlock() {
	m.rwLock.RUnlock()
}

func (m *DaMgr) lock() {
	m.rwLock.Lock()
}

func (m *DaMgr) unlock() {
	m.rwLock.Unlock()
}

func (m *DaMgr) updateMoData(moUpdateData *NfomaDataUpdate) error {
	m.lock()
	defer m.unlock()

	err := m.dbData.update(moUpdateData)
	if err != nil {
		return err
	}
	m.checkDbUpdate()

	return nil
}

// --------------------------------------- 数据库OP（go协程） End ---------------------------------------

// --------------------------------------- 业务OP（go协程） Begin ---------------------------------------
type tttProc struct {
	macIPList []MacIPRes
	boardList []BoardRes
}

type BoardRes struct {
	boardAddr    BoardAddr
	chipRes      ChipRes
	chipLinkList []ChipLink
}

func (t *tttProc) updateRes() {
	TransBegin()
	for i := range t.macIPList {
		t.macIPList[i].update()
	}
	for i := range t.boardList {
		t.boardList[i].update()
	}
	TransEnd()
}

func (b *BoardRes) update() {
	chips := GetBoardChips(b.boardAddr)
	b.chipRes.update(chips)
	for i := range b.chipLinkList {
		b.chipLinkList[i].update(b.chipRes)
	}
}

// --------------------------------------- 业务OP（go协程） End ---------------------------------------
