package types

type AddTodoLabelRequest struct {
	TodoID  string `json:"todo_id" binding:"required"`
	LabelID string `json:"label_id" binding:"required"`
}

type CreateLabelRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateLabelRequest struct {
	Name string `json:"name" binding:"required"`
}
type AddNoteLabelRequest struct {
	NoteID  string `json:"note_id" binding:"required"`
	LabelID string `json:"label_id" binding:"required"`
}
type AddDocLabelRequest struct {
	DocID  string `json:"doc_id" binding:"required"`
	LabelID string `json:"label_id" binding:"required"`
}
