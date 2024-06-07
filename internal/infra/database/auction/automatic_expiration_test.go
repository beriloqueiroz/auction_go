package auction

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/beriloqueiroz/auction_go/internal/entity/auction_entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/testcontainers/testcontainers-go"
)

func getDatabase(ctx context.Context) (database *mongo.Database, close func(ctx context.Context), err error) {
	var env = map[string]string{
		"MONGO_INITDB_ROOT_USERNAME":   "root",
		"MONGO_INITDB_ROOT_PASSWORD":   "pass",
		"MONGO_INITDB_DATABASE":        "testdb",
		"TESTCONTAINERS_HOST_OVERRIDE": "host.docker.internal",
	}
	var port = "27017/tcp"

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mongo:latest",
			ExposedPorts: []string{port},
			Env:          env,
		},
		Started: true,
	}

	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	p, err := container.MappedPort(ctx, "27017")
	if err != nil {
		return nil, nil, err
	}

	log.Println("mongo container ready and running at port: ", p.Port())
	close = func(ctxc context.Context) {
		container.Terminate(ctxc)
	}

	uri := fmt.Sprintf("mongodb://root:pass@localhost:%s", p.Port())
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, close, err
	}
	database = client.Database("auctions")
	return database, close, nil
}

func TestAutomaticExpireAuction(t *testing.T) {
	t.Run("testing time when not expire auction", func(t *testing.T) {
		ctx := context.Background()

		database, close, err := getDatabase(ctx)
		if err != nil {
			assert.Fail(t, err.Error())
		}
		defer close(ctx)

		os.Setenv("AUCTION_DURATION_HOUR", "60s")
		id := uuid.NewString()
		repo := NewAuctionRepository(database)
		repo.CreateAuction(ctx, &auction_entity.Auction{
			Id:          id,
			ProductName: "teste product name",
			Category:    "category",
			Description: "description name",
			Timestamp:   time.Now(),
			Status:      0,
		})

		time.Sleep(time.Second * 1)

		auction, err := repo.FindAuctionById(ctx, id)

		assert.EqualValues(t, auction.Status, 0)
	})

	t.Run("testing time when expire auction", func(t *testing.T) {
		ctx := context.Background()

		database, close, err := getDatabase(ctx)
		if err != nil {
			assert.Fail(t, err.Error())
		}
		defer close(ctx)

		os.Setenv("AUCTION_DURATION_HOUR", "5s")
		id := uuid.NewString()
		repo := NewAuctionRepository(database)
		repo.CreateAuction(ctx, &auction_entity.Auction{
			Id:          id,
			ProductName: "teste product name",
			Category:    "category",
			Description: "description name",
			Timestamp:   time.Now(),
			Status:      0,
		})

		time.Sleep(time.Second * 6)

		auction, err := repo.FindAuctionById(ctx, id)

		assert.EqualValues(t, auction.Status, 1)
	})
}
