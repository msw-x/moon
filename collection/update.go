package collection

type updateMode int

const (
	updatePure   updateMode = 1 // only update
	updateSoft   updateMode = 2 // only update in mem, but not update db
	updateRemove updateMode = 3 // only update, but print to log remove
	updateDelete updateMode = 4 // update db, but remove from collection
)
