package user

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/tkane/tkblog/internal/pkg/log"
	pb "github.com/tkane/tkblog/pkg/proto/tkblog/v1"
)

func (ctrl *UserCtrl) ListUser(ctx context.Context, r *pb.ListUserRequest) (*pb.ListUserResponse, error) {
	log.C(ctx).Infow("ListUser function is called!")

	resp, err := ctrl.b.Users().List(ctx, int(r.Offset), int(r.Limit))
	if err != nil {
		return nil, err
	}
	users := make([]*pb.UserInfo, 0, len(resp.Users))
	for _, u := range resp.Users {
		createAt, _ := time.Parse("2006-01-02 15:04:05", u.CreateAt)
		updateAt, _ := time.Parse("2006-01-02 15:04:05", u.UpdateAt)
		users = append(users, &pb.UserInfo{
			Username: u.Username,
			Nickname: u.Nickname,
			Email: u.Email,
			Phone: u.Phone,
			PostCount: u.PostCount,
			CreateAt: timestamppb.New(createAt),
			UdpateAt: timestamppb.New(updateAt),
		})
	}
	ret := &pb.ListUserResponse{
		TotalCount: resp.TotalCount,
		Users: users,
	}
	return ret, nil
}