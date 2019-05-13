package template

import (
	"testing"
)

func TestRenderValue(t *testing.T) {
	type args struct {
		v Value
	}
	tests := []struct {
		name       string
		args       args
		wantResult string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := RenderValue(tt.args.v); gotResult != tt.wantResult {
				t.Errorf("RenderValue() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestTemplate(t *testing.T) {
	type args struct {
		value interface{}
		temp  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Template(tt.args.value, tt.args.temp); got != tt.want {
				t.Errorf("Template() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRender(t *testing.T) {
	type args struct {
		value []byte
	}
	tests := []struct {
		name       string
		args       args
		wantResult string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := Render(tt.args.value); gotResult != tt.wantResult {
				t.Errorf("Render() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
