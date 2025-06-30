package services

import (
	"context"
	"log"
	"qiscus-test/repository"
	"time"
)

type AllocationService struct {
	redis      *repository.RedisRepository
	apiService APIService
}

func AgentAllocationService(r *repository.RedisRepository, apiService APIService) *AllocationService {
	return &AllocationService{
		redis:      r,
		apiService: apiService,
	}
}

func (s *AllocationService) EnqueueChat(roomID string) error {
	return s.redis.EnqueueChat(context.Background(), roomID)
}

// Worker untuk memproses antrian
func (s *AllocationService) ProcessQueue() {
	for {
		ctx := context.Background()
		time.Sleep(10 * time.Second)

		roomID, err := s.redis.PopChat(ctx)

		if err != nil {
			time.Sleep(60 * time.Second)
			continue
		}
		if roomID == "" {
			time.Sleep(60 * time.Second)
			continue
		}

		assigned, _ := s.redis.GetAssignedAgent(ctx, roomID)
		if assigned != "" {
			log.Printf("Room %s sudah dialokasikan ke %s\n", roomID, assigned)
			continue
		}

		agents, err := s.apiService.GetAvailableAgents(roomID)
		if err != nil {
			log.Println("Gagal ambil agent:", err)
			continue
		}

		// Cari agent yang chat aktif < 2
		allocated := false
		for _, a := range agents.Data.Agents {
			isAvailable := a.IsAvailable
			CurrentCustomerCount := a.CurrentCustomerCount

			if !isAvailable {
				continue
			}

			count, _ := s.redis.GetAgentChatCount(ctx, a.Email)
			if count != CurrentCustomerCount {
				log.Printf("Sinkronisasi count agent %s: Redis=%d, API=%d\n", a.Email, count, CurrentCustomerCount)
				_ = s.redis.SetAgentChatCount(ctx, a.Email, CurrentCustomerCount)
				count = CurrentCustomerCount
			}

			if count >= 2 {
				continue
			}

			// Jika agent kosong, ambil 2 room sekaligus
			if count == 0 && CurrentCustomerCount == 0 {
				log.Printf("Agent %s kosong, ambil 2 room sekaligus", a.Email)

				// Assign room pertama
				err := s.apiService.AssignAgent(roomID, a.Id)
				if err != nil {
					log.Println("Gagal assign room pertama:", err)
					continue
				}
				_ = s.redis.AssignAgent(ctx, roomID, a.Email)
				_ = s.redis.IncrementAgentChatCount(ctx, a.Email)
				log.Printf("Room %s dialokasikan ke agent %s", roomID, a.Email)

				// Ambil 1 room tambahan
				additionalRoomID, err := s.redis.PopChat(ctx)
				if err == nil && additionalRoomID != "" {
					err2 := s.apiService.AssignAgent(additionalRoomID, a.Id)
					if err2 != nil {
						log.Println("Gagal assign room kedua:", err2)
					} else {
						_ = s.redis.AssignAgent(ctx, additionalRoomID, a.Email)
						_ = s.redis.IncrementAgentChatCount(ctx, a.Email)
						log.Printf("Room %s dialokasikan ke agent %s", additionalRoomID, a.Email)
					}
				} else {
					log.Println("Tidak ada room tambahan")
				}

				allocated = true
				break
			}

			err := s.apiService.AssignAgent(roomID, a.Id)
			if err != nil {
				log.Println("Gagal assign agent:", err)
			}

			_ = s.redis.AssignAgent(ctx, roomID, a.Email)
			_ = s.redis.IncrementAgentChatCount(ctx, a.Email)
			log.Printf("Room %s dialokasikan ke agent %s\n", roomID, a.Email)
			allocated = true
			break
		}

		if !allocated {
			_ = s.redis.EnqueueChat(ctx, roomID)
			time.Sleep(1 * time.Second)
		}
	}
}
