package utils

import (
	"reflect"
	"testing"
)

func TestArgumentsCheck(t *testing.T) {
	type args struct {
		argCount int
		min      int
		max      int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args:    args{argCount: 1, min: 1, max: 2},
			wantErr: false,
		},
		{
			args:    args{argCount: 2, min: 1, max: 2},
			wantErr: false,
		},
		{
			args:    args{argCount: 3, min: 1, max: 2},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ArgumentsCheck(tt.args.argCount, tt.args.min, tt.args.max); (err != nil) != tt.wantErr {
				t.Errorf("ArgumentsCheck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQuery(t *testing.T) {
	type args struct {
		values     []string
		keys       []string
		ignoreCase bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "1",
			args: args{
				values: []string{"test1", "test2", "test3"},
				keys:   []string{"xxx", "test"},
			},
			want: true,
		},
		{
			name: "2",
			args: args{
				values: []string{"Test1", "Test2", "Test3"},
				keys:   []string{"test"},
			},
			want: false,
		},
		{
			name: "3",
			args: args{
				values:     []string{"Test1", "Test2", "Test3"},
				keys:       []string{"xxx", "test"},
				ignoreCase: true,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Query(tt.args.values, tt.args.keys, tt.args.ignoreCase); got != tt.want {
				t.Errorf("Query() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortKeys(t *testing.T) {
	type args struct {
		m map[string]string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			args: args{
				m: map[string]string{
					"ac": "", "ab": "",
					"cc": "", "cd": "",
					"bb": "", "ba": "",
				},
			},
			want: []string{"ab", "ac", "ba", "bb", "cc", "cd"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SortKeys(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseConnect(t *testing.T) {
	type args struct {
		connect string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
		want2 string
	}{
		{
			args: args{
				connect: "root@1.2.3.4:22022",
			},
			want:  "root",
			want1: "1.2.3.4",
			want2: "22022",
		},
		{
			args: args{
				connect: "root@1.2.3.4",
			},
			want:  "root",
			want1: "1.2.3.4",
			want2: "",
		},
		{
			args: args{
				connect: "1.2.3.4",
			},
			want:  "",
			want1: "1.2.3.4",
			want2: "",
		},
		{
			args: args{
				connect: "root@1.2.3.4",
			},
			want:  "root",
			want1: "1.2.3.4",
			want2: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := ParseConnect(tt.args.connect)
			if got != tt.want {
				t.Errorf("ParseConnect() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ParseConnect() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("ParseConnect() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
