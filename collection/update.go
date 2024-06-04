package collection

type updateMode int

const (
	updatePure   updateMode = 1 // only update
	updateRemove updateMode = 2 // only update, but print to log remove
	updateDelete updateMode = 3 // update db, but remove from collection
)
