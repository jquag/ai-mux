package data

type WorkItem struct {
	Id          string
	ShortName   string
	Description string
	Status      string
	IsClosing   bool
}

type NewWorkItemMsg struct {
	WorkItem *WorkItem
}

type UpdateWorkItemMsg struct {
	WorkItem *WorkItem
}

type WorkItemRemovedMsg struct {
	WorkItem *WorkItem
}
