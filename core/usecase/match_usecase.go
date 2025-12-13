package usecase

import (
	"miservicegolang/core/domain/match"
)

type MatchUsecase struct {
	Match    chan *match.Match
	Generate AiUsecase
	Hub      *Hub
}

func NewMatchUsecase(generate AiUsecase, hub *Hub) *MatchUsecase {
	return &MatchUsecase{
		Match:    make(chan *match.Match, 100),
		Generate: generate,
		Hub:      hub,
	}
}

func (m *MatchUsecase) MatchMaker() {
	res, _ := m.Generate.MainGenerate("Gere 1 desafio simples de programacao.")
	for match := range m.Match {
		msg := match.User1.Username + "e" + match.User2.Username + res
		m.Hub.Broadcast(match.Id, []byte(msg))
	}
}
