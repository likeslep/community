package model

import "testing"

func TestQuestionTransition(t *testing.T) {
	tests := []struct {
		name     string
		from, to QuestionStatus
		wantErr  bool
	}{
		{"开启关闭", QuestionOpen, QuestionClosed, false},
		{"开启删除", QuestionOpen, QuestionDeleted, false},
		{"关闭重开", QuestionClosed, QuestionOpen, false},
		{"关闭删除", QuestionClosed, QuestionDeleted, false},
		{"删除是终态", QuestionDeleted, QuestionOpen, true},
		{"删除不能关闭", QuestionDeleted, QuestionClosed, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Question{Status: tt.from}
			err := q.Transition(tt.to)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Transition(%s→%s) err = %v, wantErr = %v", tt.from, tt.to, err, tt.wantErr)
			}
		})
	}
}
