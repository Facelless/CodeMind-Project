package usecase

import (
	"context"
	"miservicegolang/core/domain/match"
	"miservicegolang/infrastructure/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QueueUsecase struct {
	Queue    chan *match.Queue
	Match    *MatchUsecase
	database repository.UserDatabaseRepo
	Hub      *Hub
}

func NewQueueUsecase(database repository.UserDatabaseRepo, match_u *MatchUsecase, hub *Hub) *QueueUsecase {
	return &QueueUsecase{
		Queue:    make(chan *match.Queue, 100),
		Match:    match_u,
		database: database,
		Hub:      hub,
	}
}

func (q *QueueUsecase) QueueMaker() {
	for {
		p1 := <-q.Queue
		p2 := <-q.Queue

		users, _ := q.database.FindAllUsersById(
			context.Background(),
			[]primitive.ObjectID{p1.Id, p2.Id},
		)
		u1 := users[0]
		u2 := users[1]

		if u1.Elo == u2.Elo {
			roomID := "match-" + primitive.NewObjectID().Hex()

			q.Hub.Join(roomID, p1.Conn)
			q.Hub.Join(roomID, p2.Conn)

			match := &match.Match{
				Id:    roomID,
				User1: u1,
				User2: u2,
			}

			q.Match.Match <- match
		} else {
			q.Queue <- p1
			q.Queue <- p2
		}
	}
}
