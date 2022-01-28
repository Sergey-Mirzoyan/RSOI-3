package uc_implementation

import (
	"lab2/src/rating/repositories"
)

const (
	minRating = 0
	maxRating = 100
)

type RatingUsecase struct {
	rr repositories.IRatingRepository
}

func NewRatingUsecase(repo repositories.IRatingRepository) *RatingUsecase {
	return &RatingUsecase{rr: repo}
}

func (rc *RatingUsecase) GetByUser(username string) (int, error) {
	return rc.rr.GetByUser(username)
}

func (rc *RatingUsecase) AlterByUser(diff int, username string) error {
	val, err := rc.rr.GetByUser(username)
	if err != nil {
		return err
	}
	val += diff
	if val < minRating {
		val = minRating
	} else if val > maxRating {
		val = maxRating
	}
	return rc.rr.SetByUser(val, username)
}
