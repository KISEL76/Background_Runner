package store

import (
	"testing"

	"queue-svc/internal/model"
)

// Проверяем, что можно сохранить задачу, прочитать ее и обновить статус
func TestStore_Put_Get_Update(t *testing.T) {
	s := New()

	// Добавляем задачу
	tk := &model.Task{ID: "1", Status: model.StatusQueued}
	s.Put(tk)

	// Достаем обратно
	got := s.Get("1")
	if got == nil {
		t.Fatal("task not found")
	}
	if got.Status != model.StatusQueued {
		t.Errorf("expected queued, got %v", got.Status)
	}

	// Обновляем статус
	s.UpdateStatus("1", model.StatusDone, "")
	got = s.Get("1")
	if got.Status != model.StatusDone {
		t.Errorf("expected done, got %v", got.Status)
	}

	// Смотрим снапшот
	snap := s.Snapshot()
	if len(snap) != 1 {
		t.Errorf("expected snapshot size 1, got %d", len(snap))
	}
}
