package invitation

type Status string

const (
	Pending   Status = "pending"
	Accepted  Status = "accepted"
	Finalized Status = "finalized"
)
