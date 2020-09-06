package shortUrls

import (
	"log"
	"time"

	"github.com/grafana/grafana/pkg/api/dtos"
	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/util"
)

type ShortUrlService interface {
	GetFullUrlByUID(uid string) (string, error)
	CreateShortUrl(cmd *dtos.CreateShortUrlForm) (string, error)
}

type shortUrlServiceImpl struct {
	user *models.SignedInUser
	log  log.Logger
}

var NewShortUrlService = func(orgId int64, user *models.SignedInUser) ShortUrlService {
	return &shortUrlServiceImpl{
		user: user,
	}
}

func (dr *shortUrlServiceImpl) buildCreateShortUrlCommand(path string) (*models.CreateShortUrlCommand, error) {
	cmd := &models.CreateShortUrlCommand{
		Uid:       util.GenerateShortUID(),
		Path:      path,
		CreatedBy: dr.user.UserId,
		CreatedAt: time.Now(),
	}

	return cmd, nil
}

func (dr *shortUrlServiceImpl) GetFullUrlByUID(uid string) (string, error) {
	query := models.GetFullUrlQuery{Uid: uid}
	if err := bus.Dispatch(&query); err != nil {
		return "", err
	}

	if query.Result.Path == "" {
		return "", models.ErrShortUrlNotFound
	}

	return query.Result.Path, nil
}

func (dr *shortUrlServiceImpl) CreateShortUrl(cmd *dtos.CreateShortUrlForm) (string, error) {
	createShortUrlCmd, err := dr.buildCreateShortUrlCommand(cmd.Path)
	if err != nil {
		return "", err
	}

	err = bus.Dispatch(createShortUrlCmd)
	if err != nil {
		return "", err
	}

	return createShortUrlCmd.Result.Uid, nil
}