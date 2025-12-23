package service

import (
	"context"

	pb "github.com/temesgen-abebayehu/bidflow/backend/proto/pb"
	"github.com/temesgen-abebayehu/bidflow/backend/services/bidding/internal/domain"
	"google.golang.org/grpc"
)

type auctionClient struct {
	client pb.AuctionServiceClient
}

func NewAuctionClient(conn *grpc.ClientConn) domain.AuctionClient {
	return &auctionClient{
		client: pb.NewAuctionServiceClient(conn),
	}
}

func (c *auctionClient) ValidateBid(ctx context.Context, auctionID string, amount float64, bidderID string) (bool, string, error) {
	req := &pb.BidRequest{
		AuctionId: auctionID,
		Amount:    amount,
		BidderId:  bidderID,
	}

	res, err := c.client.ValidateBid(ctx, req)
	if err != nil {
		return false, "", err
	}

	return res.IsValid, res.Message, nil
}

func (c *auctionClient) UpdateAuctionPrice(ctx context.Context, auctionID string, amount float64) error {
	req := &pb.UpdateAuctionPriceRequest{
		AuctionId: auctionID,
		Amount:    amount,
	}

	_, err := c.client.UpdateAuctionPrice(ctx, req)
	return err
}
