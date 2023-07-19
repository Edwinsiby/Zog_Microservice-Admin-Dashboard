package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"service2/pb"
	"service2/pkg/entity"
	repo "service2/pkg/repository"

	"github.com/golang/protobuf/ptypes"
)

type AdminDashboard struct {
	pb.UnimplementedAdminDashboardServer
}

func (s *AdminDashboard) MyMethod(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	log.Println("Microservice1: MyMethod called")

	result := "Hello, " + req.Data
	return &pb.Response{Result: result}, nil
}

func (s *AdminDashboard) UserList(ctx context.Context, req *pb.UserListRequest) (*pb.UserListResponse, error) {
	offset := 1
	limit := 10
	users, err := repo.GetAllUsers(offset, limit)
	if err != nil {
		return nil, err
	}
	var userList []*pb.User
	for _, u := range users {
		pbUser := &pb.User{
			Id:         int32(u.ID),
			Firstname:  u.FirstName,
			Lastname:   u.LastName,
			Email:      u.Email,
			Phone:      u.Phone,
			Wallet:     int32(u.Wallet),
			Permission: u.Permission,
		}
		userList = append(userList, pbUser)
	}
	response := &pb.UserListResponse{
		Users: userList,
	}
	return response, nil
}

func (s *AdminDashboard) SortUserByPermission(ctx context.Context, req *pb.SortUserRequest) (*pb.SortUserResponse, error) {
	offset := 1
	limit := 10
	var sortby bool
	if req.Permission == "true" {
		sortby = true
	} else {
		sortby = false
	}
	users, err := repo.GetAllUsersByPermission(offset, limit, sortby)
	if err != nil {
		return nil, err
	}
	var userList []*pb.User
	for _, u := range users {
		pbUser := &pb.User{
			Id:         int32(u.ID),
			Firstname:  u.FirstName,
			Lastname:   u.LastName,
			Email:      u.Email,
			Phone:      u.Phone,
			Wallet:     int32(u.Wallet),
			Permission: u.Permission,
		}
		userList = append(userList, pbUser)
	}
	response := &pb.SortUserResponse{
		Users: userList,
	}
	return response, nil
}

func (s *AdminDashboard) SearchUserByid(ctx context.Context, req *pb.SearchUserByidRequest) (*pb.SearchUserByidResponse, error) {
	users, err := repo.GetAllUsersById(int(req.Userid))
	if err != nil {
		return nil, err
	}
	var userList []*pb.User
	for _, u := range users {
		pbUser := &pb.User{
			Id:         int32(u.ID),
			Firstname:  u.FirstName,
			Lastname:   u.LastName,
			Email:      u.Email,
			Phone:      u.Phone,
			Wallet:     int32(u.Wallet),
			Permission: u.Permission,
		}
		userList = append(userList, pbUser)
	}
	response := &pb.SearchUserByidResponse{
		Users: userList,
	}
	return response, nil
}

func (s *AdminDashboard) SearchUserByname(ctx context.Context, req *pb.SearchUserBynameRequest) (*pb.SearchUserBynameResponse, error) {
	users, err := repo.GetAllUsersByName(req.Name)
	if err != nil {
		return nil, err
	}
	var userList []*pb.User
	for _, u := range users {
		pbUser := &pb.User{
			Id:         int32(u.ID),
			Firstname:  u.FirstName,
			Lastname:   u.LastName,
			Email:      u.Email,
			Phone:      u.Phone,
			Wallet:     int32(u.Wallet),
			Permission: u.Permission,
		}
		userList = append(userList, pbUser)
	}
	response := &pb.SearchUserBynameResponse{
		Users: userList,
	}
	return response, nil
}

func (s *AdminDashboard) TogglePermission(ctx context.Context, req *pb.TogglePermissionRequest) (*pb.TogglePermissionResponse, error) {
	result, err := repo.GetByID(int(req.Userid))
	if err != nil {
		return nil, err
	}
	result.Permission = !result.Permission
	err = repo.Update(result)
	if err != nil {
		return nil, errors.New("user permission toggling failed")
	}
	resp := &pb.TogglePermissionResponse{
		Result: "User permission toggled succesfuly",
	}
	return resp, nil
}

func (s *AdminDashboard) CreateApparel(ctx context.Context, req *pb.CreateApparelRequest) (*pb.CreateApparelResponse, error) {
	jwtToken, ok := ctx.Value("jwtToken").(string)
	if !ok {
		fmt.Println("unable to retrieve JWT token from context")
	} else {
		fmt.Println(jwtToken)
	}
	err := repo.GetByApparelName(req.Name)
	if err == nil {
		return nil, errors.New("product already exists")
	}
	newApparel := &entity.Apparel{
		Name:     req.Name,
		Price:    int(req.Price),
		ImageURL: req.Image,
		AdminId:  int(req.Adminid),
	}
	_, err1 := repo.CreateApparel(newApparel)
	if err1 != nil {
		return nil, err1
	} else {
		resp := &pb.CreateApparelResponse{
			Result: "Apparel added succesfuly",
		}
		return resp, nil
	}
}

func (s *AdminDashboard) EditApparel(ctx context.Context, req *pb.EditApparelResquest) (*pb.EditApparelResponse, error) {
	newApparel, err := repo.GetApparelByID(int(req.Id))
	if err != nil {
		return nil, err
	}
	newApparel = &entity.Apparel{
		Name:     req.Name,
		Price:    int(req.Price),
		ImageURL: req.Image,
		AdminId:  int(req.Adminid),
	}
	err = repo.UpdateApparel(newApparel)
	if err != nil {
		return nil, err
	} else {
		resp := &pb.EditApparelResponse{
			Result: "Apparel edited succesfuly",
		}
		return resp, nil
	}
}

func (s *AdminDashboard) DeleteApparel(ctx context.Context, req *pb.DeleteApparelRequest) (*pb.DeleteApparelResponse, error) {
	result, err := repo.GetApparelByID(int(req.Id))
	if err != nil {
		return nil, err
	}
	result.Removed = !result.Removed
	err = repo.UpdateApparel(result)
	if err != nil {
		return nil, err
	} else {
		resp := &pb.DeleteApparelResponse{
			Result: "Apparel deleted succesfuly",
		}
		return resp, nil
	}

}

func (s *AdminDashboard) AddCoupon(ctx context.Context, req *pb.AddCouponRequest) (*pb.AddCouponResponse, error) {
	validTime, err := ptypes.Timestamp(req.Valid)
	if err != nil {
		return nil, errors.New("Creating Timestamp failed")
	}
	coupon := &entity.Coupon{
		Code:       req.Code,
		Type:       req.Type,
		Category:   req.Category,
		Amount:     int(req.Amount),
		UsageLimit: int(req.Limit),
		ValidUntil: validTime,
	}
	err = repo.CreateCoupon(coupon)
	if err != nil {
		return nil, errors.New("Creating Coupon failed")
	} else {
		resp := &pb.AddCouponResponse{
			Result: "Coupon Added succesfuly",
		}
		return resp, nil
	}
}

func (s *AdminDashboard) AddOffer(ctx context.Context, req *pb.AddOfferRequest) (*pb.AddOfferResponse, error) {
	offer := &entity.Offer{
		Name:       req.Code,
		Type:       req.Type,
		Category:   req.Category,
		Amount:     int(req.Amount),
		UsageLimit: int(req.Limit),
		MinPrice:   int(req.Minprice),
	}
	err := repo.CreateOffer(offer)
	if err != nil {
		return nil, errors.New("Creating Offer failed")
	} else {
		resp := &pb.AddOfferResponse{
			Result: "Coupon added succesfuly",
		}
		return resp, nil
	}
}
