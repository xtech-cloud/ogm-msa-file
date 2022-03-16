package model

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Object struct {
	UUID      string `gorm:"column:uuid;type:char(32);not null;unique;primaryKey"`
	Bucket    string `gorm:"column:bucket;type:char(32);not null"`
	Path      string `gorm:"column:path;type:varchar(256);not null"`
	Hash      string `gorm:"column:hash;type:char(32)"`
	URL       string `gorm:"column:url;type:varchar(1024)"`
	Size      uint64 `gorm:"column:size"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

var ErrObjectExists = errors.New("object exists")
var ErrObjectNotFound = errors.New("object not found")

func (Object) TableName() string {
	return "ogm_file_object"
}

type ObjectDAO struct {
	conn *Conn
}

type ObjectQuery struct {
	Bucket string
	Path   string
	Hash   string
}

func NewObjectDAO(_conn *Conn) *ObjectDAO {
	conn := DefaultConn
	if nil != _conn {
		conn = _conn
	}
	return &ObjectDAO{
		conn: conn,
	}
}

func (this *ObjectDAO) Count() (int64, error) {
	var count int64
	err := this.conn.DB.Model(&Object{}).Count(&count).Error
	return count, err
}

func (this *ObjectDAO) CountOfBucket(_bucket string) (int64, error) {
	var count int64
	err := this.conn.DB.Model(&Object{}).Where("bucket = ?", _bucket).Count(&count).Error
	return count, err
}

func (this *ObjectDAO) SumOfBucket(_bucket string) (uint64, error) {
	return 0, nil
}

func (this *ObjectDAO) Insert(_object *Object) error {
	var count int64
	err := this.conn.DB.Model(&Object{}).Where("uuid = ?", _object.UUID).Count(&count).Error
	if nil != err {
		return err
	}

	if count > 0 {
		return ErrObjectExists
	}

	return this.conn.DB.Create(_object).Error
}

func (this *ObjectDAO) Upsert(_object *Object) error {
	// uuid冲突时，更新所有列
	this.conn.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		UpdateAll: true,
	}).Create(_object)
	return nil
}

func (this *ObjectDAO) Update(_object *Object) error {
	var count int64
	err := this.conn.DB.Model(&Object{}).Where("path = ?", _object.Path).Count(&count).Error
	if nil != err {
		return err
	}

	if 0 == count {
		return ErrObjectNotFound
	}

	// 使用select选定更新字段，零值也会被更新
	return this.conn.DB.Select("path", "url", "size", "hash").Updates(_object).Error
}

func (this *ObjectDAO) Delete(_uuid string) error {
	return this.conn.DB.Where("uuid = ?", _uuid).Delete(&Object{}).Error
}

func (this *ObjectDAO) List(_offset int64, _count int64, _bucket string) (_total int64, _object []*Object, _err error) {
	_total = int64(0)
	_err = nil
	_object = make([]*Object, 0)

	db := this.conn.DB.Model(&Object{})
	if "" != _bucket {
		db = db.Where("bucket = ?", _bucket)
	}
	_err = db.Count(&_total).Error
	if nil != _err {
		return
	}
	_err = db.Offset(int(_offset)).Limit(int(_count)).Order("created_at desc").Find(&_object).Error
	return
}

func (this *ObjectDAO) Search(_offset int64, _count int64, _bucket string, _prefix string, _name string) (_total int64, _object []*Object, _err error) {
	_total = int64(0)
	_err = nil
	_object = make([]*Object, 0)

	db := this.conn.DB.Model(&Object{})
	if "" != _bucket {
		db = db.Where("bucket = ?", _bucket)
	}
	path := ""
	if "" != _prefix {
		path = _prefix + "%"
	}
	if "" != _name {
		if !strings.HasSuffix(path, "%") {
			path = path + "%"
		}
		path = path + _name + "%"
	}
	db = db.Where("path LIKE ?", path)
	_err = db.Count(&_total).Error
	if nil != _err {
		return
	}
	_err = db.Offset(int(_offset)).Limit(int(_count)).Order("created_at desc").Find(&_object).Error
	return
}

func (this *ObjectDAO) WherePath(_bucket string, _like []string, _notlike []string, _prefix string) (_object []Object, _err error) {
	db_like := this.conn.DB.Model(&Object{})
	if nil != _like {
		for _, like := range _like {
			db_like = db_like.Or("path LIKE ?", like)
		}
	}
	db_notlike := this.conn.DB.Model(&Object{})
	if nil != _notlike {
		for _, notlike := range _notlike {
			db_notlike = db_notlike.Not("path LIKE ?", notlike)
		}
	}
	db := this.conn.DB.Model(&Object{})
	var object []Object
	db = db.Where("bucket = ?", _bucket)
	if "" != _prefix {
		db = db.Where("path LIKE ?", _prefix+"%")
	}
	db = db.Where(db_like).Where(db_notlike).Order("created_at desc").Find(&object)
	err := db.Statement.Error
	return object, err
}

func (this *ObjectDAO) QueryOne(_query *ObjectQuery) (*Object, error) {
	db := this.conn.DB.Model(&Object{})
	hasWhere := false
	if "" != _query.Bucket {
		db = db.Where("bucket = ?", _query.Bucket)
		hasWhere = true
	}
	if "" != _query.Path {
		db = db.Where("path = ?", _query.Path)
		hasWhere = true
	}
	if "" != _query.Hash {
		db = db.Where("hash = ?", _query.Hash)
		hasWhere = true
	}
	if !hasWhere {
		return nil, ErrObjectNotFound
	}

	var object Object
	err := db.Limit(1).Find(&object).Error
	if "" == object.UUID {
		return nil, ErrObjectNotFound
	}
	return &object, err
}

func (this *ObjectDAO) Get(_uuid string) (*Object, error) {
	var object Object
	db := this.conn.DB.Model(&Object{}).Where("uuid = ?", _uuid)
	err := db.First(&object).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrObjectNotFound
	}
	return &object, err
}

func (this *ObjectDAO) Clean(_bucket string, _prefix string) error {
	db := this.conn.DB.Where("bucket = ?", _bucket)
	if "" != _prefix {
		db = db.Where("path LIKE ?", _prefix+"%")
	}
	return db.Delete(&Object{}).Error
}
