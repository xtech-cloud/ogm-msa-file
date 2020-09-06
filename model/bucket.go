package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Bucket struct {
	UUID         string `gorm:"column:uuid;type:char(32);not null;unique;primaryKey"`
	Name         string `gorm:"column:name;type:varchar(256);not null;unique"`
	Token        string `gorm:"column:token;type:char(32)"`
	TotalSize    uint64 `gorm:"column:size_total;not null;default:0"`
	UsedSize     uint64 `gorm:"column:size_used;not null;default:0"`
	Engine       int    `gorm:"column:engine"`
	Address      string `gorm:"column:address;type:varchar(512)"`
	Scope        string `gorm:"column:scope;type:varchar(512)"`
	AccessKey    string `gorm:"column:access_key;type:varchar(1024)"`
	AccessSecret string `gorm:"column:access_secret;type:varchar(1024)"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

var ErrBucketExists = errors.New("bucket exists")
var ErrBucketNotFound = errors.New("bucket not found")

func (Bucket) TableName() string {
	return "msa_file_bucket"
}

type BucketQuery struct {
	UUID string
}

type BucketDAO struct {
	conn *Conn
}

func NewBucketDAO(_conn *Conn) *BucketDAO {
	conn := DefaultConn
	if nil != _conn {
		conn = _conn
	}
	return &BucketDAO{
		conn: conn,
	}
}

func (this *BucketDAO) Count() (int64, error) {
	var count int64
	err := this.conn.DB.Model(&Bucket{}).Count(&count).Error
	return count, err
}

func (this *BucketDAO) Insert(_bucket *Bucket) error {
	var count int64
	err := this.conn.DB.Model(&Bucket{}).Where("uuid = ? OR name = ?", _bucket.UUID, _bucket.Name).Count(&count).Error
	if nil != err {
		return err
	}

	if count > 0 {
		return ErrBucketExists
	}

	return this.conn.DB.Create(_bucket).Error
}

func (this *BucketDAO) Update(_bucket *Bucket) error {
	var count int64
	err := this.conn.DB.Model(&Bucket{}).Where("uuid = ?", _bucket.UUID).Count(&count).Error
	if nil != err {
		return err
	}

	if 0 == count {
		return ErrBucketNotFound
	}

	return this.conn.DB.Updates(_bucket).Error
}

func (this *BucketDAO) Delete(_uuid string) error {
	var count int64
	err := this.conn.DB.Model(&Bucket{}).Where("uuid = ?", _uuid).Count(&count).Error
	if nil != err {
		return err
	}

	if 0 == count {
		return ErrBucketNotFound
	}

	return this.conn.DB.Where("uuid = ?", _uuid).Delete(&Bucket{}).Error
}

func (this *BucketDAO) List(_offset int64, _count int64) ([]*Bucket, error) {
	var buckets []*Bucket
	res := this.conn.DB.Offset(int(_offset)).Limit(int(_count)).Order("created_at desc").Find(&buckets)
	return buckets, res.Error
}

func (this *BucketDAO) QueryOne(_query *BucketQuery) (*Bucket, error) {
	db := this.conn.DB.Model(&Bucket{})
	hasWhere := false
	if "" != _query.UUID {
		db = db.Where("uuid = ?", _query.UUID)
		hasWhere = true
	}
	if !hasWhere {
		return nil, ErrBucketNotFound
	}

	var bucket Bucket
	err := db.First(&bucket).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrBucketNotFound
	}
	return &bucket, err
}
