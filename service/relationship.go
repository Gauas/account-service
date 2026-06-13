package service

import (
// "context"
// "errors"
// "net/http"
// "strings"
//
// "github.com/gauas/account-service/dto/response"
// "github.com/gauas/account-service/model"
// "github.com/gauas/account-service/packages/httpresp"
// "gorm.io/gorm"
)

//func (s *Service) ListRelationships(ctx context.Context, user *model.User, status model.RelationshipStatus) ([]responses2.RelationshipResponse, error) {
//	//query := s.Repository.Relationship.WithContext(ctx).
//	//	Preload("Actor").
//	//	Preload("Partner").
//	//	Where("actor_id = ? OR partner_id = ?", user.ID, user.ID).
//	//	Order("updated_at DESC")
//	//
//	//if status != "" {
//	//	query = query.Where("status = ?", status)
//	//}
//	//
//	//var relationships []model.Relationship
//	//if err := query.Find(&relationships).Error; err != nil {
//	//	return nil, err
//	//}
//	//
//	//out := make([]response.Relationship, 0, len(relationships))
//	//for i := range relationships {
//	//	out = append(out, relationshipResponse(&relationships[i]))
//	//}
//	//
//	//return out, nil
//	return nil, nil
//}
//
//func (s *Service) RequestRelationship(ctx context.Context, actorKey string, partnerKey string) (*responses2.RelationshipResponse, error) {
//	actor, partner, err := s.relationshipUserPair(ctx, actorKey, partnerKey)
//	if err != nil {
//		return nil, err
//	}
//
//	existing, err := s.findRelationshipBetween(ctx, actor.ID, partner.ID)
//	if err == nil {
//		return nil, relationshipConflict(existing)
//	}
//	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
//		return nil, err
//	}
//
//	relationship := &model.Relationship{
//		ActorID:   actor.ID,
//		PartnerID: partner.ID,
//		Status:    model.RelationshipStatusPending,
//	}
//
//	if _, err := s.Repository.Relationship.Create(ctx, relationship); err != nil {
//		return nil, err
//	}
//
//	relationship.Actor = *actor
//	relationship.Partner = *partner
//
//	out := relationshipResponse(relationship)
//	return &out, nil
//}
//
//func (s *Service) AcceptRelationship(ctx context.Context, userKey string, actorKey string) (*responses2.RelationshipResponse, error) {
//	user, actor, relationship, err := s.relationshipWithTarget(ctx, userKey, actorKey)
//	if err != nil {
//		return nil, err
//	}
//
//	if relationship.PartnerID != user.ID || relationship.ActorID != actor.ID {
//		return nil, httpresp.NewError(http.StatusForbidden, "only the requested user can accept this relationship")
//	}
//	if relationship.Status != model.RelationshipStatusPending {
//		return nil, httpresp.NewError(http.StatusConflict, "relationship is not pending")
//	}
//
//	relationship.Status = model.RelationshipStatusActive
//	if err := s.Repository.Relationship.WithContext(ctx).
//		Model(&model.Relationship{}).
//		Where("id = ?", relationship.ID).
//		Update("status", relationship.Status).Error; err != nil {
//		return nil, err
//	}
//
//	updated, err := s.relationshipByID(ctx, relationship.ID)
//	if err != nil {
//		return nil, err
//	}
//
//	out := relationshipResponse(updated)
//	return &out, nil
//}
//
//func (s *Service) DeclineRelationship(ctx context.Context, userKey string, actorKey string) error {
//	user, actor, relationship, err := s.relationshipWithTarget(ctx, userKey, actorKey)
//	if err != nil {
//		return err
//	}
//
//	if relationship.PartnerID != user.ID || relationship.ActorID != actor.ID {
//		return httpresp.NewError(http.StatusForbidden, "only the requested user can decline this relationship")
//	}
//	if relationship.Status != model.RelationshipStatusPending {
//		return httpresp.NewError(http.StatusConflict, "relationship is not pending")
//	}
//
//	return s.Repository.Relationship.Delete(ctx, "id = ?", relationship.ID)
//}
//
//func (s *Service) CancelRelationship(ctx context.Context, userKey string, partnerKey string) error {
//	user, partner, relationship, err := s.relationshipWithTarget(ctx, userKey, partnerKey)
//	if err != nil {
//		return err
//	}
//
//	if relationship.ActorID != user.ID || relationship.PartnerID != partner.ID {
//		return httpresp.NewError(http.StatusForbidden, "only the requesting user can cancel this relationship")
//	}
//	if relationship.Status != model.RelationshipStatusPending {
//		return httpresp.NewError(http.StatusConflict, "only pending relationships can be cancelled")
//	}
//
//	return s.Repository.Relationship.Delete(ctx, "id = ?", relationship.ID)
//}
//
//func (s *Service) relationshipUserPair(ctx context.Context, userKey string, targetUserKey string) (*model.User, *model.User, error) {
//	if strings.TrimSpace(userKey) == "" {
//		return nil, nil, httpresp.NewError(http.StatusUnauthorized, "unauthorized")
//	}
//
//	user, err := s.GetProfileByKey(ctx, userKey)
//	if err != nil {
//		return nil, nil, err
//	}
//
//	target, err := s.GetProfileByKey(ctx, targetUserKey)
//	if err != nil {
//		return nil, nil, err
//	}
//	if user.ID == target.ID {
//		return nil, nil, httpresp.NewError(http.StatusBadRequest, "cannot create relationship with yourself")
//	}
//
//	return user, target, nil
//}
//
//func (s *Service) relationshipWithTarget(ctx context.Context, userKey string, targetUserKey string) (*model.User, *model.User, *model.Relationship, error) {
//	user, target, err := s.relationshipUserPair(ctx, userKey, targetUserKey)
//	if err != nil {
//		return nil, nil, nil, err
//	}
//
//	relationship, err := s.findRelationshipBetween(ctx, user.ID, target.ID)
//	if errors.Is(err, gorm.ErrRecordNotFound) {
//		return nil, nil, nil, httpresp.NewError(http.StatusNotFound, "relationship not found")
//	}
//	if err != nil {
//		return nil, nil, nil, err
//	}
//
//	return user, target, relationship, nil
//}
//
//func (s *Service) relationshipByID(ctx context.Context, id int64) (*model.Relationship, error) {
//	var relationship model.Relationship
//	err := s.Repository.Relationship.WithContext(ctx).
//		Preload("Actor").
//		Preload("Partner").
//		Take(&relationship, "id = ?", id).Error
//	if errors.Is(err, gorm.ErrRecordNotFound) {
//		return nil, httpresp.NewError(http.StatusNotFound, "relationship not found")
//	}
//	if err != nil {
//		return nil, err
//	}
//
//	return &relationship, nil
//}
//
//func (s *Service) findRelationshipBetween(ctx context.Context, actorID int64, partnerID int64) (*model.Relationship, error) {
//	return s.Repository.Relationship.Take(
//		ctx,
//		"(actor_id = ? AND partner_id = ?) OR (actor_id = ? AND partner_id = ?)",
//		actorID,
//		partnerID,
//		partnerID,
//		actorID,
//	)
//}
//
//func relationshipConflict(relationship *model.Relationship) error {
//	switch relationship.Status {
//	case model.RelationshipStatusActive:
//		return httpresp.NewError(http.StatusConflict, "relationship already exists")
//	case model.RelationshipStatusPending:
//		return httpresp.NewError(http.StatusConflict, "relationship request already exists")
//	default:
//		return httpresp.NewError(http.StatusConflict, "relationship already exists")
//	}
//}
//
//func relationshipResponse(relationship *model.Relationship) responses2.RelationshipResponse {
//	return response.Relationship.RelationshipResponse{
//		Key:       relationship.Key.String(),
//		Status:    string(relationship.Status),
//		Actor:     httpresp.Refine[*model.User, responses2.ProfileResponse](&relationship.Actor),
//		Partner:   httpresp.Refine[*model.User, responses2.ProfileResponse](&relationship.Partner),
//		CreatedAt: relationship.CreatedAt,
//		UpdatedAt: relationship.UpdatedAt,
//	}
//}
