package hydra

import (
	"context"
	"fmt"

	"github.com/yqcs/prismx/core/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func MongodbWeakPass(res any) any {
	var (
		t   = res.(models.HydraTask)
		msg = models.MSG{
			Name: "Mongodb WeakPassword",
			Type: "WeakPassword",
			Payload: models.Dict{
				User:     t.Dict.User,
				Password: t.Dict.Password,
			},
			Target: t.Target,
		}
	)
	ctx, cancel := context.WithTimeout(context.Background(), t.Config.Timeout)
	opt := options.Client()
	opt.SetDialer(&proxyDialer{
		timeout: t.Config.Timeout,
	})
	opt.ApplyURI(fmt.Sprintf("mongodb://%v:%v@%v/ichunt?authMechanism=SCRAM-SHA-1", t.Dict.User, t.Dict.Password, t.Target))
	client, err := mongo.Connect(ctx, opt)
	defer cancel()
	if err != nil {
		return nil
	}
	defer client.Disconnect(ctx)
	if err = client.Ping(context.Background(), readpref.Primary()); err != nil {
		return nil
	}
	return msg
}
