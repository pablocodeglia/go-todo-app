package clistore

import (
	"reflect"
	"testing"
)

func TestAddParallel(t *testing.T) {
	tests := []struct {
		name            string
		todoDescription string
	}{
		{"test1", "1st new TODO"},
		{"test2", "2nd new TODO"},
		{"test3", "3rd new TODO"},
		{"test4", "4th new TODO"},
		{"test5", "5th new TODO"},
		{"test6", "6th new TODO"},
	}

	store := NewStore()

	t.Parallel()
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store.Add(tt.todoDescription)
			got := len(store.Data)
			want := i+1
			if !reflect.DeepEqual(got, want) {
				t.Errorf("expected %d, got %d", want, got)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	store := NewStore()

	store.Add("TEST NEW TASK")
	store.Add("TEST ANOTHER NEW TASK")

	got := len(store.Data)
	want := 2

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}

}

func BenchmarkClearCache(b *testing.B) {
	store := NewStore()

	for i := 0; i < b.N; i++ {
		store.ClearCache()
	}
}

func BenchmarkListTodos(b *testing.B) {
	store := NewStore()

	for i := 0; i < b.N; i++ {
		store.ListTodos()
	}
}
