package auction

import (
	"context"
	"os"
	"time"

	"github.com/beriloqueiroz/auction_go/configuration/logger"
	"github.com/beriloqueiroz/auction_go/internal/entity/auction_entity"
	"github.com/beriloqueiroz/auction_go/internal/internal_error"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}
type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	repo := &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
	go func() {
		for {
			repo.expireAuction(context.Background())
		}
	}()
	return repo
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	return nil
}

func (ar *AuctionRepository) expireAuction(ctx context.Context) *internal_error.InternalError {
	filter := bson.M{"status": auction_entity.Active}

	list, err := ar.Collection.Find(ctx, filter)
	if err != nil {
		logger.Error("Error finding auctions to expire", err)
		return internal_error.NewInternalServerError("Error finding auctions")
	}
	defer list.Close(ctx)

	var auctionsMongo []AuctionEntityMongo
	if err := list.All(ctx, &auctionsMongo); err != nil {
		logger.Error("Error decoding auctions", err)
		return internal_error.NewInternalServerError("Error decoding auctions")
	}

	auctionLimitDuration := os.Getenv("AUCTION_DURATION_HOUR")
	duration, err := time.ParseDuration(auctionLimitDuration)
	if err != nil {
		duration = time.Hour * 24
	}

	for _, auction := range auctionsMongo {
		tm := time.Unix(auction.Timestamp, 0)
		if tm.Add(duration).Before(time.Now()) {
			filter := bson.D{{Key: "_id", Value: auction.Id}}
			update := bson.D{{Key: "$set",
				Value: bson.D{
					{Key: "status", Value: auction_entity.Completed},
				},
			}}
			_, err := ar.Collection.UpdateOne(ctx, filter, update)
			if err != nil {
				logger.Error("Error trying to expire auction", err)
			}
		}
	}

	return nil
}
