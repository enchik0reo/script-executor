package commander

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/enchik0reo/commandApi/internal/logs"
	"github.com/enchik0reo/commandApi/internal/models"
	"github.com/enchik0reo/commandApi/internal/services/commander/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCommander_CreateNewCommand(t *testing.T) {
	type fields struct {
		Storager  *mocks.Storager
		Executor  *mocks.Executor
		log       *logs.CustomLog
		stopChans *sync.Map
	}
	type args struct {
		ctx    context.Context
		script string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
		prepare func(args2 args, fields *fields)
	}{
		{
			name: "test_1, no error",
			args: args{
				ctx:    context.Background(),
				script: "whoami",
			},
			want: 1,
			prepare: func(args2 args, fields *fields) {
				script := "whoami"

				fields.Storager.On("CreateNew", mock.Anything, mock.Anything).Return(int64(1), nil)
				fields.Executor.On("RunScript", script, script, mock.Anything).
					Return(make(<-chan string), make(<-chan error))
			},
		},
		{
			name: "test_2, with db error",
			args: args{
				ctx:    context.Background(),
				script: "whoami",
			},
			want:    -1,
			wantErr: true,
			prepare: func(args2 args, fields *fields) {
				fields.Storager.On("CreateNew", mock.Anything, mock.Anything).
					Return(int64(0), errors.New("some db error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dlog := logs.NewDiscardLogger()

			f := fields{
				Storager:  mocks.NewStorager(t),
				Executor:  mocks.NewExecutor(t),
				log:       dlog,
				stopChans: &sync.Map{},
			}

			tt.prepare(tt.args, &f)

			c := &Commander{
				cmdStorage: f.Storager,
				exec:       f.Executor,
				log:        f.log,
				stopChans:  f.stopChans,
			}

			got, err := c.CreateNewCommand(tt.args.ctx, tt.args.script)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			if err != nil && !tt.wantErr {
				t.Errorf("Unexpected error: %s", err.Error())
			}

			require.Equal(t, tt.want, got)
		})
	}
}

func TestCommander_GetCommandList(t *testing.T) {
	type fields struct {
		Storager *mocks.MockStorager
		Executor *mocks.MockExecutor
		log      *logs.CustomLog
	}
	type args struct {
		ctx   context.Context
		limit int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.Command
		wantErr bool
		prepare func(args2 args, fields fields)
	}{
		{
			name: "test_1, no error",
			args: args{
				ctx:   context.Background(),
				limit: 1,
			},
			want: []models.Command{{}},
			prepare: func(args2 args, fields fields) {
				fields.Storager.EXPECT().GetList(
					context.Background(),
					int64(1)).
					Return([]models.Command{{}}, nil)
			},
		},
		{
			name: "test_2, no error",
			args: args{
				ctx:   context.Background(),
				limit: 5,
			},
			want: []models.Command{{}, {}, {}, {}, {}},
			prepare: func(args2 args, fields fields) {
				fields.Storager.EXPECT().GetList(
					context.Background(),
					int64(5)).
					Return([]models.Command{{}, {}, {}, {}, {}}, nil)
			},
		},
		{
			name: "test_3, with error",
			args: args{
				ctx:   context.Background(),
				limit: -1,
			},
			want:    nil,
			wantErr: true,
			prepare: func(args2 args, fields fields) {
				fields.Storager.EXPECT().GetList(
					context.Background(),
					int64(-1)).
					Return(nil, errors.New("limit might be more then 0"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dlog := logs.NewDiscardLogger()

			ctrl := gomock.NewController(t)

			f := fields{
				Storager: mocks.NewMockStorager(ctrl),
				Executor: mocks.NewMockExecutor(ctrl),
				log:      dlog,
			}

			tt.prepare(tt.args, f)

			c := NewCommander(dlog, f.Storager, f.Executor)

			got, err := c.GetCommandList(tt.args.ctx, tt.args.limit)

			if err != nil {
				if !tt.wantErr {
					require.Error(t, err)
				}
			}

			require.Equal(t, tt.want, got)
		})
	}
}

func TestCommander_GetOneCommandDescription(t *testing.T) {
	type fields struct {
		Storager *mocks.MockStorager
		Executor *mocks.MockExecutor
		log      *logs.CustomLog
	}
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Command
		wantErr bool
		prepare func(args2 args, fields fields)
	}{
		{
			name: "test_1, no error",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want: &models.Command{
				ID:        1,
				Name:      "whoami",
				Output:    []string{"root"},
				IsWorking: false,
			},
			prepare: func(args2 args, fields fields) {
				fields.Storager.EXPECT().GetOne(
					context.Background(),
					int64(1)).Return(&models.Command{
					ID:        1,
					Name:      "whoami",
					Output:    []string{"root"},
					IsWorking: false,
				}, nil)
			},
		}, {
			name: "test_2, with error",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want:    nil,
			wantErr: true,
			prepare: func(args2 args, fields fields) {
				fields.Storager.EXPECT().GetOne(
					context.Background(),
					int64(1)).Return(nil, errors.New("there's no script"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dlog := logs.NewDiscardLogger()

			ctrl := gomock.NewController(t)

			f := fields{
				Storager: mocks.NewMockStorager(ctrl),
				Executor: mocks.NewMockExecutor(ctrl),
				log:      dlog,
			}

			tt.prepare(tt.args, f)

			c := NewCommander(dlog, f.Storager, f.Executor)

			got, err := c.GetOneCommandDescription(tt.args.ctx, tt.args.id)

			if err != nil {
				if !tt.wantErr {
					require.Error(t, err)
				}
			}

			require.Equal(t, tt.want, got)
		})
	}
}

func TestCommander_StopCommand(t *testing.T) {
	type fields struct {
		Storager  *mocks.MockStorager
		Executor  *mocks.MockExecutor
		log       *logs.CustomLog
		stopChans *sync.Map
	}
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
		prepare func(args2 args, fields *fields)
	}{
		{
			name: "test_1, no error",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want: 1,
			prepare: func(args2 args, fields *fields) {
				ch := make(chan struct{})

				go func() {
					<-ch
					close(ch)
				}()

				fields.stopChans.Store(int64(1), ch)

				fields.Storager.EXPECT().StopOne(
					context.Background(),
					int64(1)).Return(int64(1), nil)
			},
		},
		{
			name: "test_2, with error",
			args: args{
				ctx: context.Background(),
				id:  10,
			},
			want:    0,
			wantErr: true,
			prepare: func(args2 args, fields *fields) {},
		},
		{
			name: "test_3, with db error",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want: 0,
			prepare: func(args2 args, fields *fields) {
				ch := make(chan struct{})

				go func() {
					<-ch
					close(ch)
				}()

				fields.stopChans.Store(int64(1), ch)

				fields.Storager.EXPECT().StopOne(
					context.Background(),
					int64(1)).Return(int64(0), errors.New("db error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dlog := logs.NewDiscardLogger()

			ctrl := gomock.NewController(t)

			f := fields{
				Storager:  mocks.NewMockStorager(ctrl),
				Executor:  mocks.NewMockExecutor(ctrl),
				log:       dlog,
				stopChans: &sync.Map{},
			}

			tt.prepare(tt.args, &f)

			c := &Commander{
				cmdStorage: f.Storager,
				exec:       f.Executor,
				log:        f.log,
				stopChans:  f.stopChans,
			}

			got, err := c.StopCommand(tt.args.ctx, tt.args.id)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			if err != nil && !tt.wantErr {
				t.Errorf("Unexpected error: %s", err.Error())
			}

			require.Equal(t, tt.want, got)
		})
	}
}

func TestCommander_StopAllRunningScripts(t *testing.T) {
	type fields struct {
		Storager  *mocks.MockStorager
		Executor  *mocks.MockExecutor
		log       *logs.CustomLog
		stopChans *sync.Map
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		prepare func(args2 args, fields *fields)
	}{
		{
			name: "test_1, no error",
			args: args{
				ctx: context.Background(),
			},
			prepare: func(args2 args, fields *fields) {
				ch := make(chan struct{})

				go func() {
					<-ch
					close(ch)
				}()

				ch2 := make(chan struct{})

				go func() {
					<-ch2
					close(ch2)
				}()

				fields.stopChans.Store(int64(1), ch)
				fields.stopChans.Store(int64(2), ch2)

				fields.Storager.EXPECT().StopOne(
					context.Background(),
					int64(1)).Return(int64(1), nil)

				fields.Storager.EXPECT().StopOne(
					context.Background(),
					int64(2)).Return(int64(2), nil)
			},
		},
		{
			name: "test_2, no error empty line",
			args: args{
				ctx: context.Background(),
			},
			prepare: func(args2 args, fields *fields) {},
		},
		{
			name: "test_3, with db error",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
			prepare: func(args2 args, fields *fields) {
				ch := make(chan struct{})

				go func() {
					<-ch
					close(ch)
				}()

				fields.stopChans.Store(int64(1), ch)

				fields.Storager.EXPECT().StopOne(
					context.Background(),
					int64(1)).Return(int64(0), errors.New("db error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dlog := logs.NewDiscardLogger()

			ctrl := gomock.NewController(t)

			f := fields{
				Storager:  mocks.NewMockStorager(ctrl),
				Executor:  mocks.NewMockExecutor(ctrl),
				log:       dlog,
				stopChans: &sync.Map{},
			}

			tt.prepare(tt.args, &f)

			c := &Commander{
				cmdStorage: f.Storager,
				exec:       f.Executor,
				log:        f.log,
				stopChans:  f.stopChans,
			}

			err := c.StopAllRunningScripts(tt.args.ctx)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			if err != nil && !tt.wantErr {
				t.Errorf("Unexpected error: %s", err.Error())
			}
		})
	}
}

func TestScriptName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "test_1, empty script",
			input: "",
			want:  "",
		},
		{
			name:  "test_2, short script",
			input: "ls",
			want:  "ls",
		},
		{
			name:  "test_2, long script",
			input: "echo \"Hello, World! I'm a very long script.\"",
			want:  "echo \"Hello, World! I'm a v...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := scriptName(tt.input)

			if got != tt.want {
				t.Errorf("Expected: %q, got %q\n", tt.want, got)
			}
		})
	}
}
