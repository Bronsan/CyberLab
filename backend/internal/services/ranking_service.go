package services

import (
	"time"

	"github.com/cyberlab/backend/internal/models"
	"github.com/cyberlab/backend/internal/repository"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

type RankingService struct {
	userRepo     *repository.UserRepository
	redisClient  *redis.Client
}

func NewRankingService(userRepo *repository.UserRepository, redisClient *redis.Client) *RankingService {
	return &RankingService{
		userRepo:    userRepo,
		redisClient: redisClient,
	}
}

type RankingEntry struct {
	Rank        int    `json:"rank"`
	Username    string `json:"username"`
	Score       int    `json:"score"`
	SolvedCount int    `json:"solvedCount"`
}

type RankingResponse struct {
	List []RankingEntry `json:"list"`
}

func (s *RankingService) GetGlobalRanking() (*RankingResponse, error) {
	// Try Redis cache first
	ctx := context.Background()
	results, err := s.redisClient.ZRevRangeWithScores(ctx, "rank:global", 0, 99).Result()
	if err == nil && len(results) > 0 {
		var entries []RankingEntry
		for i, z := range results {
			entries = append(entries, RankingEntry{
				Rank:  i + 1,
				Username: z.Member.(string),
				Score: int(z.Score),
			})
		}
		return &RankingResponse{List: entries}, nil
	}

	// Fallback to MySQL
	users, err := s.userRepo.GetRanking(100)
	if err != nil {
		return nil, err
	}

	var entries []RankingEntry
	for i, user := range users {
		entries = append(entries, RankingEntry{
			Rank:        i + 1,
			Username:    user.Username,
			Score:       user.Score,
			SolvedCount: user.SolvedCount,
		})
	}

	// Update Redis cache
	s.syncRankingToRedis(users)

	return &RankingResponse{List: entries}, nil
}

func (s *RankingService) GetWeeklyRanking() (*RankingResponse, error) {
	// Try Redis cache first
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	results, err := s.redisClient.ZRevRangeWithScores(ctx, "rank:weekly", 0, 99).Result()
	if err == nil && len(results) > 0 {
		entries := make([]RankingEntry, 0, len(results))
		for i, z := range results {
			entries = append(entries, RankingEntry{
				Rank:     i + 1,
				Username:  z.Member.(string),
				Score:     int(z.Score),
			})
		}
		return &RankingResponse{List: entries}, nil
	}

	// Fallback to MySQL
	weekAgo := time.Now().AddDate(0, 0, -7).Format("2006-01-02 15:04:05")
	users, err := s.userRepo.GetRankingByDate(weekAgo, 100)
	if err != nil {
		return nil, err
	}

	entries := make([]RankingEntry, 0, len(users))
	for i, user := range users {
		entries = append(entries, RankingEntry{
			Rank:        i + 1,
			Username:    user.Username,
			Score:       user.Score,
			SolvedCount: user.SolvedCount,
		})
	}
	s.cacheRanking("rank:weekly", users, 5*time.Minute)
	return &RankingResponse{List: entries}, nil
}

func (s *RankingService) GetMonthlyRanking() (*RankingResponse, error) {
	// Try Redis cache first
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	results, err := s.redisClient.ZRevRangeWithScores(ctx, "rank:monthly", 0, 99).Result()
	if err == nil && len(results) > 0 {
		entries := make([]RankingEntry, 0, len(results))
		for i, z := range results {
			entries = append(entries, RankingEntry{
				Rank:     i + 1,
				Username:  z.Member.(string),
				Score:     int(z.Score),
			})
		}
		return &RankingResponse{List: entries}, nil
	}

	// Fallback to MySQL
	monthAgo := time.Now().AddDate(0, -1, 0).Format("2006-01-02 15:04:05")
	users, err := s.userRepo.GetRankingByDate(monthAgo, 100)
	if err != nil {
		return nil, err
	}

	entries := make([]RankingEntry, 0, len(users))
	for i, user := range users {
		entries = append(entries, RankingEntry{
			Rank:        i + 1,
			Username:    user.Username,
			Score:       user.Score,
			SolvedCount: user.SolvedCount,
		})
	}
	s.cacheRanking("rank:monthly", users, 5*time.Minute)
	return &RankingResponse{List: entries}, nil
}

// cacheRanking writes ranking entries to Redis ZSET.
func (s *RankingService) cacheRanking(key string, users []models.User, ttl time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	pipe := s.redisClient.Pipeline()
	for _, user := range users {
		pipe.ZAdd(ctx, key, &redis.Z{
			Score:  float64(user.Score),
			Member: user.Username,
		})
	}
	pipe.Expire(ctx, key, ttl)
	pipe.Exec(ctx)
}

func (s *RankingService) syncRankingToRedis(users []models.User) {
	ctx := context.Background()
	pipe := s.redisClient.Pipeline()
	for _, user := range users {
		pipe.ZAdd(ctx, "rank:global", &redis.Z{
			Score:  float64(user.Score),
			Member: user.Username,
		})
	}
	pipe.Expire(ctx, "rank:global", 5*time.Minute)
	pipe.Exec(ctx)
}

// UpdateUserScore updates a user's score and refreshes Redis cache
func (s *RankingService) UpdateUserScore(userID uint, score int) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	newScore := user.Score + score
	if err := s.userRepo.UpdateFields(userID, map[string]interface{}{
		"score": newScore,
		"solved_count": user.SolvedCount + 1,
	}); err != nil {
		return err
	}

	// Invalidate Redis cache
	ctx := context.Background()
	s.redisClient.Del(ctx, "rank:global")

	return nil
}
