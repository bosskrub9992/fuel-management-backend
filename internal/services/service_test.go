package services

import (
	"reflect"
	"testing"

	"github.com/bosskrub9992/fuel-management-backend/internal/entities/domains"
)

func Test_getMapFuelUsageIDToFuelUsers(t *testing.T) {
	type args struct {
		fuelUsageUsers []FuelUsageUser
	}
	tests := []struct {
		name string
		args args
		want map[int64]string
	}{
		{
			name: "no data",
			args: args{
				fuelUsageUsers: []FuelUsageUser{},
			},
			want: map[int64]string{},
		},
		{
			name: "data is nil",
			args: args{
				fuelUsageUsers: nil,
			},
			want: map[int64]string{},
		},
		{
			name: "multiple fuel usages and multiple fuel usage users",
			args: args{
				fuelUsageUsers: []FuelUsageUser{
					{
						FuelUsageUser: domains.FuelUsageUser{
							FuelUsageID: 1,
							UserID:      1,
							IsPaid:      false,
						},
						Nickname: "Boss",
					},
					{
						FuelUsageUser: domains.FuelUsageUser{
							FuelUsageID: 1,
							UserID:      2,
							IsPaid:      true,
						},
						Nickname: "Best",
					},
					{
						FuelUsageUser: domains.FuelUsageUser{
							FuelUsageID: 1,
							UserID:      3,
							IsPaid:      false,
						},
						Nickname: "Pat",
					},
					{
						FuelUsageUser: domains.FuelUsageUser{
							FuelUsageID: 2,
							UserID:      1,
							IsPaid:      true,
						},
						Nickname: "Boss",
					},
					{
						FuelUsageUser: domains.FuelUsageUser{
							FuelUsageID: 3,
							UserID:      1,
							IsPaid:      true,
						},
						Nickname: "Boss",
					},
					{
						FuelUsageUser: domains.FuelUsageUser{
							FuelUsageID: 3,
							UserID:      3,
							IsPaid:      false,
						},
						Nickname: "Pat",
					},
				},
			},
			want: map[int64]string{
				1: "❌Boss ✅Best ❌Pat",
				2: "✅Boss",
				3: "✅Boss ❌Pat",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getMapFuelUsageIDToFuelUsers(tt.args.fuelUsageUsers); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getMapFuelUsageIDToFuelUsers() = %v, want %v", got, tt.want)
			}
		})
	}
}
