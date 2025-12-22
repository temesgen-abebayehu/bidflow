package handler

import (
	"context"
	"time"

	pb "github.com/temesgen-abebayehu/bidflow/backend/proto/pb"
	"github.com/temesgen-abebayehu/bidflow/backend/services/auction/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcHandler struct {
	pb.UnimplementedAuctionServiceServer
	service *service.AuctionService
}

func NewGrpcHandler(service *service.AuctionService) *GrpcHandler {
	return &GrpcHandler{service: service}
}

func (h *GrpcHandler) CreateAuction(ctx context.Context, req *pb.CreateAuctionRequest) (*pb.CreateAuctionResponse, error) {
	startTime := time.Unix(req.StartTime, 0)
	endTime := time.Unix(req.EndTime, 0)

	auction, err := h.service.CreateAuction(
		ctx,
		req.SellerId,
		req.Title,
		req.Description,
		req.StartPrice,
		startTime,
		endTime,
		req.Category,
		req.ImageUrl,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create auction: %v", err)
	}

	return &pb.CreateAuctionResponse{
		Auction: &pb.Auction{
			Id:           auction.ID,
			SellerId:     auction.SellerID,
			Title:        auction.Title,
			Description:  auction.Description,
			StartPrice:   auction.StartPrice,
			CurrentPrice: auction.CurrentPrice,
			Status:       string(auction.Status),
			StartTime:    auction.StartTime.Unix(),
			EndTime:      auction.EndTime.Unix(),
			Category:     auction.Category,
			ImageUrl:     auction.ImageURL,
		},
	}, nil
}

func (h *GrpcHandler) GetAuction(ctx context.Context, req *pb.GetAuctionRequest) (*pb.GetAuctionResponse, error) {
	auction, err := h.service.GetAuction(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "auction not found: %v", err)
	}

	return &pb.GetAuctionResponse{
		Auction: &pb.Auction{
			Id:           auction.ID,
			SellerId:     auction.SellerID,
			Title:        auction.Title,
			Description:  auction.Description,
			StartPrice:   auction.StartPrice,
			CurrentPrice: auction.CurrentPrice,
			Status:       string(auction.Status),
			StartTime:    auction.StartTime.Unix(),
			EndTime:      auction.EndTime.Unix(),
			Category:     auction.Category,
			ImageUrl:     auction.ImageURL,
		},
	}, nil
}

func (h *GrpcHandler) ListAuctions(ctx context.Context, req *pb.ListAuctionsRequest) (*pb.ListAuctionsResponse, error) {
	auctions, total, err := h.service.ListAuctions(ctx, int(req.Page), int(req.Limit), req.Status, req.Category)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list auctions: %v", err)
	}

	var pbAuctions []*pb.Auction
	for _, a := range auctions {
		pbAuctions = append(pbAuctions, &pb.Auction{
			Id:           a.ID,
			SellerId:     a.SellerID,
			Title:        a.Title,
			Description:  a.Description,
			StartPrice:   a.StartPrice,
			CurrentPrice: a.CurrentPrice,
			Status:       string(a.Status),
			StartTime:    a.StartTime.Unix(),
			EndTime:      a.EndTime.Unix(),
			Category:     a.Category,
			ImageUrl:     a.ImageURL,
		})
	}

	return &pb.ListAuctionsResponse{
		Auctions:   pbAuctions,
		TotalCount: total,
	}, nil
}

func (h *GrpcHandler) UpdateAuction(ctx context.Context, req *pb.UpdateAuctionRequest) (*pb.UpdateAuctionResponse, error) {
	auction, err := h.service.UpdateAuction(ctx, req.Id, req.Title, req.Description, req.ImageUrl)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update auction: %v", err)
	}

	return &pb.UpdateAuctionResponse{
		Auction: &pb.Auction{
			Id:           auction.ID,
			SellerId:     auction.SellerID,
			Title:        auction.Title,
			Description:  auction.Description,
			StartPrice:   auction.StartPrice,
			CurrentPrice: auction.CurrentPrice,
			Status:       string(auction.Status),
			StartTime:    auction.StartTime.Unix(),
			EndTime:      auction.EndTime.Unix(),
			Category:     auction.Category,
			ImageUrl:     auction.ImageURL,
		},
	}, nil
}

func (h *GrpcHandler) CloseAuction(ctx context.Context, req *pb.CloseAuctionRequest) (*pb.CloseAuctionResponse, error) {
	err := h.service.CloseAuction(ctx, req.Id)
	if err != nil {
		return &pb.CloseAuctionResponse{Success: false, Message: err.Error()}, nil
	}
	return &pb.CloseAuctionResponse{Success: true, Message: "Auction closed successfully"}, nil
}

func (h *GrpcHandler) ValidateBid(ctx context.Context, req *pb.BidRequest) (*pb.BidResponse, error) {
	isValid, msg, err := h.service.ValidateBid(ctx, req.AuctionId, req.Amount)
	if err != nil {
		// If error is "not found", return valid=false with message
		return &pb.BidResponse{IsValid: false, Message: msg}, nil
	}

	auction, _ := h.service.GetAuction(ctx, req.AuctionId)
	currentPrice := 0.0
	if auction != nil {
		currentPrice = auction.CurrentPrice
	}

	return &pb.BidResponse{
		IsValid:      isValid,
		CurrentPrice: currentPrice,
		Message:      msg,
	}, nil
}

func (h *GrpcHandler) GetAuctionStatus(ctx context.Context, req *pb.StatusRequest) (*pb.StatusResponse, error) {
	auction, err := h.service.GetAuction(ctx, req.AuctionId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "auction not found")
	}

	return &pb.StatusResponse{
		AuctionId:    auction.ID,
		Title:        auction.Title,
		CurrentPrice: auction.CurrentPrice,
		Status:       string(auction.Status),
		EndTimeUnix:  auction.EndTime.Unix(),
	}, nil
}
