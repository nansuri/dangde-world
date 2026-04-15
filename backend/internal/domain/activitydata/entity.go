package activitydata

type Item struct {
	ID           uint   `json:"id"`
	ActivityID   uint   `json:"activityId"`
	AssignmentID uint   `json:"assignmentId"`
	KidID        uint   `json:"kidId"`
	Key          string `json:"key"`
	Value        string `json:"value"`
}
