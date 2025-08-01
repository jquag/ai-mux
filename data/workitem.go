package workitem

type WorkItem struct {
	Id          string
	BranchName  string
	Description string
	PlanMode    bool
	Status      string
}

type NewWorkItemMsg struct {
	WorkItem *WorkItem
}
