// Author : rexdu
// Time : 2020-03-22 23:14
package repositories

import "seckill/demos/irisDemo/datamodels"

type MovieRepository interface {
	GetMovieName() string
}

type MovieManager struct {
}

func NewMovieManager() MovieRepository {
	return &MovieManager{}
}

func (m MovieManager) GetMovieName() string {
	movie := &datamodels.Movie{Name: "rexdu视频"}
	return movie.Name
}
