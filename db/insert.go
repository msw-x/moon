package db

func BatchInsert[T any](db *Db, l []T, batchSize int) (err error) {
	if batchSize == 0 {
		batchSize = 1000
	}
	for {
		n := len(l)
		if n == 0 {
			break
		}
		var batch []T
		if n > batchSize {
			batch = l[:batchSize]
			l = l[batchSize:]
		} else {
			batch = l
			l = l[:0]
		}
		err = db.Insert(&batch)
		if err != nil {
			break
		}
	}
	return
}
