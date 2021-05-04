package provider

import (
	"blogfa/auth/database/mysql"
	"blogfa/auth/model"
	"blogfa/auth/pkg/jtrace"
	"blogfa/auth/pkg/logger"
	pb "blogfa/auth/proto"
	"context"
	"fmt"

	"go.uber.org/zap"
)

// Register method, register a provider
func (p *Provider) Register(ctx context.Context, prov model.Provider) error {
	span, _ := jtrace.Tracer.SpanFromContext(ctx, "register_provider")
	defer span.Finish()
	span.SetTag("register", "register provider model")

	tx := mysql.Storage.GetDatabase().Begin()

	if err := tx.Create(&prov).Error; err != nil {
		tx.Rollback()
		return err
	}
	defer tx.Commit()

	return nil
}

// Get method, get a Provider with table name, query and args for search
func (p *Provider) Get(ctx context.Context, table string, query interface{}, args ...interface{}) (model.Provider, error) {
	span, _ := jtrace.Tracer.SpanFromContext(ctx, "get Provider model")
	defer span.Finish()
	span.SetTag("model", "get Provider model")

	tx := mysql.Storage.GetDatabase().Begin()

	var provider = model.Provider{}
	if err := tx.Preload("User").Table(table).Where(query, args...).First(&provider).Error; err != nil {
		log := logger.GetZapLogger(false)
		logger.Prepare(log).
			Append(zap.Any("error", fmt.Sprintf("get Provider: %s", err))).
			Level(zap.ErrorLevel).
			Development().
			Commit("env")
		tx.Rollback()
		return provider, err
	}
	defer tx.Commit()

	return provider, nil
}

// Update method, update provider
// first get provider with ID
// then update provider fields with new fields and update database
func (p *Provider) Update(ctx context.Context, prov model.Provider) error {
	span, _ := jtrace.Tracer.SpanFromContext(ctx, "get provider model")
	defer span.Finish()
	span.SetTag("model", fmt.Sprintf("update provider with id: %d", prov.ID))

	// get provder with ID
	provider, err := p.Get(ctx, "providers", "id = ?", prov.ID)
	if err != nil {
		return err
	}

	// update provider attrs
	provider.FixedNumber = prov.FixedNumber
	provider.Company = prov.Company
	provider.CardNumber = prov.CardNumber
	provider.ShebaNumber = prov.ShebaNumber
	provider.Card = prov.Card
	provider.Address = prov.Address

	tx := mysql.Storage.GetDatabase().Begin()

	// start to update provider
	if err := tx.Table("providers").Where("id = ?", prov.ID).Select("*").Updates(&provider).Error; err != nil {
		log := logger.GetZapLogger(false)
		logger.Prepare(log).
			Append(zap.Any("error", fmt.Sprintf("update user error: %s", err))).
			Level(zap.ErrorLevel).
			Development().
			Commit("env")
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

// Search method for search providers
func (p Provider) Search(ctx context.Context, from, to int, search string) ([]model.Provider, error) {
	span, _ := jtrace.Tracer.SpanFromContext(ctx, "search provider model")
	defer span.Finish()
	span.SetTag("model", fmt.Sprintf("search providers"))

	tx := mysql.Storage.GetDatabase().Begin()

	// provider list
	var providers []model.Provider
	err := tx.
		Preload("User").
		Table("providers").
		Where("fixed_number LIKE ?", "%"+search+"%").
		Or("company LIKE ?", "%"+search+"%").
		Or("card LIKE ?", "%"+search+"%").
		Or("card_number LIKE ?", "%"+search+"%").
		Or("sheba_number LIKE ?", "%"+search+"%").
		Or("address LIKE ?", "%"+search+"%").
		Or("user_id LIKE ?", "%"+search+"%").
		Limit(to - from).
		Offset(from).
		Select("*").
		Find(&providers).Error
	if err != nil {
		log := logger.GetZapLogger(false)
		logger.Prepare(log).
			Append(zap.Any("error", fmt.Sprintf("update providers error: %s", err))).
			Level(zap.ErrorLevel).
			Development().
			Commit("env")
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	return providers, nil
}

func (p Provider) ToProto(prvs []model.Provider) []*pb.Providers {
	resp := make([]*pb.Providers, len(prvs))

	for i, p := range prvs {
		resp[i] = &pb.Providers{
			FixedNumber: p.FixedNumber,
			Company:     p.Company,
			Card:        p.Card,
			CardNumber:  p.CardNumber,
			ShebaNumber: p.ShebaNumber,
			Address:     p.Address,
		}
	}

	return resp
}
