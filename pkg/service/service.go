package service

import "github.com/0quz/gitlab-jira-integration/pkg/model"

// dependency part
type MarcoService struct {
	MarcoModel *model.MarcoModel
}

func NewMarcoService(m *model.MarcoModel) MarcoService {
	return MarcoService{MarcoModel: m}
}

// Add an event to DB
func (m *MarcoService) CreateEvent(e *model.Event) error {
	return m.MarcoModel.CreateEvent(e)
}

// Find an event if it exists in DB
func (m *MarcoService) FindEventByMrUrl(mrUrl string) error {
	return m.MarcoModel.FindEventByMrUrl(mrUrl)
}

// Add an hotfix error to DB
func (m *MarcoService) CreateHotFixError(h *model.HotFixError) error {
	return m.MarcoModel.CreateHotFixError(h)
}
