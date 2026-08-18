package model

import "testing"

func TestTransition(t *testing.T) {
	tests := []struct {
		name     string
		from, to Status
		wantErr  bool
	}{
		{"草稿提交审核", StatusDraft, StatusPendingReview, false},
		{"草稿删除", StatusDraft, StatusDeleted, false},
		{"待审通过", StatusPendingReview, StatusPublished, false},
		{"待审驳回", StatusPendingReview, StatusRejected, false},
		{"待审撤回", StatusPendingReview, StatusDeleted, false},
		{"已发布隐藏", StatusPublished, StatusHidden, false},
		{"已发布删除", StatusPublished, StatusDeleted, false},
		{"驳回回草稿", StatusRejected, StatusDraft, false},
		{"驳回重新提交", StatusRejected, StatusPendingReview, false},
		{"隐藏恢复", StatusHidden, StatusPublished, false},
		{"隐藏删除", StatusHidden, StatusDeleted, false},

		// 非法流转
		{"草稿不能直接发布", StatusDraft, StatusPublished, true},
		{"草稿不能驳回", StatusDraft, StatusRejected, true},
		{"已发布不能回草稿", StatusPublished, StatusDraft, true},
		{"已发布不能重新提交", StatusPublished, StatusPendingReview, true},
		{"已删除是终态", StatusDeleted, StatusPublished, true},
		{"隐藏不能提交审核", StatusHidden, StatusPendingReview, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Article{Status: tt.from}
			err := a.Transition(tt.to)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Transition(%s→%s) err = %v, wantErr = %v", tt.from, tt.to, err, tt.wantErr)
			}
			if err == nil && a.Status != tt.to {
				t.Fatalf("状态未更新: got %s want %s", a.Status, tt.to)
			}
		})
	}
}
