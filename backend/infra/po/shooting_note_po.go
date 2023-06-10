package po

import (
	"time"

	"code-shooting/infra/database/pg/sql"
	"code-shooting/infra/logger"
)

type ShootingNotePo struct {
	ID       string `gorm:"primary_key;not null" json:"id"`
	UserID   string `json:"userid"`
	UserName string `json:"username"`
	TargetID string `json:"targetid"`
	RangeID  string `json:"rangeid"`
	Records  []byte `json:"records"`

	UpdateTime time.Time `json:"updatedtime"`
}

type ShootingNoteDB struct {
	*sql.GormDB
}

func (m *ShootingNoteDB) Save(a *ShootingNotePo) error {
	result := m.Model(a).Save(a)
	logger.Infof("ShootingNoteDB Save, %s %s %s recordLen:%d, result:[v%]", a.UserID, a.TargetID, a.RangeID, len(a.Records), result.Error)
	return result.Error
}

func (m *ShootingNoteDB) Delete(userid, targetid, rangeid string) error {
	result := m.Model(&ShootingNotePo{}).Where("target_id=? and user_id=? and range_id=?", targetid, userid, rangeid).Delete(&ShootingNotePo{})
	return result.Error
}

func (m *ShootingNoteDB) FindByMultiId(userid, targetid, rangeid string) (*ShootingNotePo, error) {
	pos := []*ShootingNotePo{}
	result := m.Model(&ShootingNotePo{}).Where("target_id=? and user_id=? and range_id=?", targetid, userid, rangeid).Find(&pos)

	if result.Error != nil {
		return nil, result.Error
	}
	return m.findLatest(pos), nil
}

func (m *ShootingNoteDB) findLatest(pos []*ShootingNotePo) *ShootingNotePo {
	if len(pos) == 0 {
		return nil
	}
	u := pos[0]
	for i := range pos {
		if pos[i].UpdateTime.Unix() > u.UpdateTime.Unix() {
			u = pos[i]
		}
	}
	return u
}

func (m *ShootingNoteDB) FindMulti(targetid string) ([]*ShootingNotePo, error) {
	notes := make([]*ShootingNotePo, 0)
	result := m.Find(&notes, "target_id=?", targetid)
	return m.findLatests(notes), result.Error
}

func (m *ShootingNoteDB) findLatests(notes []*ShootingNotePo) []*ShootingNotePo {
	userRangeMap := map[string]*ShootingNotePo{}

	for i := range notes {
		key := notes[i].UserID + notes[i].RangeID
		if _, exist := userRangeMap[key]; !exist {
			userRangeMap[key] = notes[i]
			continue
		}
		if notes[i].UpdateTime.Unix() > userRangeMap[key].UpdateTime.Unix() {
			userRangeMap[key] = notes[i]
		}
	}

	results := make([]*ShootingNotePo, 0, len(userRangeMap))
	for _, v := range userRangeMap {
		results = append(results, v)
	}
	return results
}

func (m *ShootingNoteDB) TableName() string {
	return "shootingNote"
}
