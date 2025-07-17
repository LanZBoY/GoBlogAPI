package user

import "wentee/blog/app/repo/icollection"

type IUserCollection interface {
	icollection.ICountDocuments
	icollection.IFind
	icollection.IFindOne
	icollection.IInsertOne
	icollection.IUpdateOne
	icollection.IDeleteOne
}
