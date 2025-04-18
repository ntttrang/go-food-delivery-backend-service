package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestGenerateCode_Generate(t *testing.T) {
	userId := uuid.MustParse("019615db-9adb-7eff-ba03-45017274084c")
	tests := []struct {
		name    string
		g       *GenerateCode
		want    string
		wantErr bool
	}{
		{
			name:    "TC 1",
			g:       &GenerateCode{},
			want:    "a",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GenerateCode{}
			got, err := g.Execute(context.Background(), userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateCode.Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateCode.Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}
