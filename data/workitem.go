package workitem

type WorkItem struct {
	Id          string
	BranchName  string
	Description string
	PlanMode    bool
}

type NewWorkItemMsg struct {
	WorkItem *WorkItem
}
