package handler

import (
	"context"

	pb "github.com/temesgen-abebayehu/bidflow/backend/proto/pb"
	"github.com/temesgen-abebayehu/bidflow/backend/services/bidding/internal/service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GrpcHandler struct {
	pb.UnimplementedBiddingServiceServer
	service *service.BiddingService
}

func NewGrpcHandler(service *service.BiddingService) *GrpcHandler {
	return &GrpcHandler{service: service}
}

func (h *GrpcHandler) PlaceBid(ctx context.Context, req *pb.PlaceBidRequest) (*pb.PlaceBidResponse, error) {
	bid, err := h.service.PlaceBid(ctx, req.AuctionId, req.BidderId, req.Amount)
	if err != nil {
		return nil, err
	}

	return &pb.PlaceBidResponse{
		Bid: &pb.Bid{
			Id:        bid.ID,
			AuctionId: bid.AuctionID,
			BidderId:  bid.BidderID,
			Amount:    bid.Amount,
			Timestamp: timestamppb.New(bid.Timestamp),
		},
	}, nil
}

func (h *GrpcHandler) GetBidsByAuction(ctx context.Context, req *pb.GetBidsByAuctionRequest) (*pb.GetBidsByAuctionResponse, error) {
	bids, err := h.service.GetBidsByAuction(ctx, req.AuctionId)
	if err != nil {
		return nil, err
	}

	var pbBids []*pb.Bid
	for _, b := range bids {
		pbBids = append(pbBids, &pb.Bid{
			Id:        b.ID,
			AuctionId: b.AuctionID,
			BidderId:  b.BidderID,
			Amount:    b.Amount,
			Timestamp: timestamppb.New(b.Timestamp),
		})
	}

	return &pb.GetBidsByAuctionResponse{
		Bids: pbBids,
	}, nil
}
