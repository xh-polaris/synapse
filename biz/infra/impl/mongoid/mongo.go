package mongoid

import (
	"context"

	"github.com/xh-polaris/synapse/biz/infra/contract/id"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func New(ctx context.Context) id.IDGenerator {
	return &objectIDGenerator{}
}

type objectIDGenerator struct{}

func (i *objectIDGenerator) GenID(_ context.Context) id.ID {
	return id.ID(primitive.NewObjectID())
}

func (i *objectIDGenerator) GenMultiIDs(_ context.Context, counts int) (ids []id.ID) {
	for _ = range counts {
		ids = append(ids, id.ID(primitive.NewObjectID()))
	}
	return
}
