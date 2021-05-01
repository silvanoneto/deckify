package user

import (
	"reflect"
	"testing"

	"github.com/zmb3/spotify"
)

func TestNewUserRepoInMemoryImpl(t *testing.T) {
	tests := []struct {
		name string
		want UserRepo
	}{
		{
			name: "createAnInstance",
			want: &userRepoInMemoryImpl{users: make(map[string]User)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserRepoInMemoryImpl(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserRepoInMemoryImpl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userRepoInMemoryImpl_InsertOrUpdate(t *testing.T) {
	type fields struct {
		users map[string]User
	}
	type args struct {
		user User
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "insertEntity",
			fields: fields{users: make(map[string]User)},
			args:   args{user: User{UserInfo: spotify.PrivateUser{User: spotify.User{ID: "ABC"}}, Active: true}},
		},
		{
			name: "updateEntity",
			fields: fields{
				users: map[string]User{
					"ABC": {UserInfo: spotify.PrivateUser{User: spotify.User{ID: "ABC"}}, Active: false},
				},
			},
			args: args{user: User{UserInfo: spotify.PrivateUser{User: spotify.User{ID: "ABC"}}, Active: true}},
		},
		{
			name:   "insertOrUpdateEmptyEntity",
			fields: fields{users: make(map[string]User)},
			args:   args{user: User{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &userRepoInMemoryImpl{
				users: tt.fields.users,
			}
			repo.InsertOrUpdate(tt.args.user)
		})
	}
}

func Test_userRepoInMemoryImpl_Remove(t *testing.T) {
	type fields struct {
		users map[string]User
	}
	type args struct {
		ID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "removeNonExistentUser",
			fields:  fields{users: make(map[string]User)},
			args:    args{ID: "ABC"},
			wantErr: true,
		},
		{
			name: "removeAnExistentUser",
			fields: fields{
				users: map[string]User{
					"ABC": {UserInfo: spotify.PrivateUser{User: spotify.User{ID: "ABC"}}, Active: true},
				},
			},
			args:    args{ID: "ABC"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &userRepoInMemoryImpl{
				users: tt.fields.users,
			}
			if err := repo.Remove(tt.args.ID); (err != nil) != tt.wantErr {
				t.Errorf("userRepoInMemoryImpl.Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userRepoInMemoryImpl_GetByID(t *testing.T) {
	type fields struct {
		users map[string]User
	}
	type args struct {
		ID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    User
		wantErr bool
	}{
		{
			name:    "getNonExistentUser",
			fields:  fields{users: make(map[string]User)},
			args:    args{ID: "ABC"},
			want:    User{},
			wantErr: true,
		},
		{
			name: "getAnExistentUser",
			fields: fields{
				users: map[string]User{
					"ABC": {UserInfo: spotify.PrivateUser{User: spotify.User{ID: "ABC"}}, Active: true},
				},
			},
			args:    args{ID: "ABC"},
			want:    User{UserInfo: spotify.PrivateUser{User: spotify.User{ID: "ABC"}}, Active: true},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &userRepoInMemoryImpl{
				users: tt.fields.users,
			}
			got, err := repo.GetByID(tt.args.ID)
			if (err != nil) != tt.wantErr {
				t.Errorf("userRepoInMemoryImpl.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userRepoInMemoryImpl.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userRepoInMemoryImpl_GetAllActive(t *testing.T) {
	type fields struct {
		users map[string]User
	}
	type args struct {
		pageNumber uint
		pageSize   uint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []User
	}{
		{
			name:   "getWhenUserListIsEmpty",
			fields: fields{users: make(map[string]User)},
			args:   args{pageNumber: 0, pageSize: 1},
			want:   []User{},
		},
		{
			name: "getUserListFirstPage",
			fields: fields{
				users: map[string]User{
					"ABC": {UserInfo: spotify.PrivateUser{User: spotify.User{ID: "ABC"}}, Active: true},
				},
			},
			args: args{pageNumber: 0, pageSize: 10},
			want: []User{{UserInfo: spotify.PrivateUser{User: spotify.User{ID: "ABC"}}, Active: true}},
		},
		{
			name: "getUserListEmptyPage",
			fields: fields{
				users: map[string]User{
					"ABC": {UserInfo: spotify.PrivateUser{User: spotify.User{ID: "ABC"}}, Active: true},
				},
			},
			args: args{pageNumber: 1, pageSize: 1},
			want: []User{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &userRepoInMemoryImpl{
				users: tt.fields.users,
			}
			if got := repo.GetAllActive(tt.args.pageNumber, tt.args.pageSize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userRepoInMemoryImpl.GetAllActive() = %v, want %v", got, tt.want)
			}
		})
	}
}
