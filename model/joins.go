package model

type JoinsDAO struct {
	conn *Conn
}

func NewJoinsDAO(_conn *Conn) *JoinsDAO {
	conn := DefaultConn
	if nil != _conn {
		conn = _conn
	}
	return &JoinsDAO{
		conn: conn,
	}
}

type JoinsQuery struct {
    Bucket string
    Filepath string
}


func (this *JoinsDAO) SearchObject(_offset int64, _count int64, _query *JoinsQuery) (_total int64, _object []*Object, _err error) {
    _err = nil
    _total = int64(0)
    _object = make([]*Object, 0)

	db := this.conn.DB
    db = db.Joins("JOIN msa_file_bucket ON msa_file_bucket.uuid = msa_file_object.bucket")
	db = db.Where("msa_file_bucket.name LIKE ?", _query.Bucket+ "%")
	if "" != _query.Filepath{
		db = db.Where("msa_file_object.filepath LIKE ?", _query.Filepath + "%")
	}
    db = db.Model(&Object{})

	_err = db.Count(&_total).Error
    if nil != _err {
        return
    }
	_err = db.Offset(int(_offset)).Limit(int(_count)).Order("msa_file_object.created_at desc").Find(&_object).Error
	return
}
