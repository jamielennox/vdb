package common

type CollectionTarget struct {
	Name     CollectionName
	Id       CollectionId
	Revision RevisionID
	Type     string
	Labels   Labels
}

type Event struct {
	Operation OperationType
	Target    CollectionTarget
	Subject   UserInfo
}
