// Author : rexdu
// Time : 2020-03-22 23:19
package services

import "seckill/demos/irisDemo/repositories"

type MovieService interface {
	ShowMovieName() string
}

type MovieServiceManger struct {
	repo repositories.MovieRepository
}

func NewMovieServiceManger(repo repositories.MovieRepository) MovieService {
	return &MovieServiceManger{repo: repo}
}

func (m *MovieServiceManger) ShowMovieName() string {
	return "视频名字是：" + m.repo.GetMovieName()

}
