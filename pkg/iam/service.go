package iam

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/tendant/simple-idm/pkg/iam/iamdb"
	"github.com/tendant/simple-idm/pkg/utils"
	"golang.org/x/exp/slog"
)

type IamService struct {
	queries *iamdb.Queries
}

func NewIamService(queries *iamdb.Queries) *IamService {
	return &IamService{
		queries: queries,
	}
}

func (s *IamService) CreateUser(ctx context.Context, email, username, name string, roleIds []uuid.UUID, loginID string) (iamdb.GetUserWithRolesRow, error) {
	// Validate email
	if email == "" {
		return iamdb.GetUserWithRolesRow{}, fmt.Errorf("email is required")
	}
	// Validate username
	if username == "" {
		return iamdb.GetUserWithRolesRow{}, fmt.Errorf("username is required")
	}

	// Create the user first
	nullString := sql.NullString{String: name, Valid: name != ""}
	nullLoginID := sql.NullString{String: loginID, Valid: loginID != ""}

	// Note: Username field is removed as it doesn't exist in the struct
	user, err := s.queries.CreateUser(ctx, iamdb.CreateUserParams{
		Email:   email,
		Name:    nullString,
		LoginID: utils.NullStringToNullUUID(nullLoginID),
	})
	if err != nil {
		return iamdb.GetUserWithRolesRow{}, fmt.Errorf("failed to create user: %w", err)
	}

	// If there are roles to assign, create the user-role associations
	if len(roleIds) > 0 {
		slog.Info("Assigning roles to user", "userId", user.ID, "roleIds", roleIds)
		// Insert role assignments one by one
		for _, roleId := range roleIds {
			slog.Info("Assigning role", "userId", user.ID, "roleId", roleId)
			_, err = s.queries.CreateUserRole(ctx, iamdb.CreateUserRoleParams{
				UserID: user.ID,
				RoleID: roleId,
			})
			if err != nil {
				slog.Error("Failed to assign role", "error", err, "userId", user.ID, "roleId", roleId)
				return iamdb.GetUserWithRolesRow{}, fmt.Errorf("failed to assign role: %w", err)
			}
		}
	} else {
		slog.Info("No roles to assign", "userId", user.ID)
	}

	// Get the user with roles
	userWithRoles, err := s.queries.GetUserWithRoles(ctx, user.ID)
	if err != nil {
		return iamdb.GetUserWithRolesRow{}, fmt.Errorf("failed to get user with roles: %w", err)
	}

	return userWithRoles, nil
}

func (s *IamService) FindUsers(ctx context.Context) ([]iamdb.FindUsersWithRolesRow, error) {
	return s.queries.FindUsersWithRoles(ctx)
}

func (s *IamService) GetUser(ctx context.Context, userId uuid.UUID) (iamdb.GetUserWithRolesRow, error) {
	return s.queries.GetUserWithRoles(ctx, userId)
}

func (s *IamService) UpdateUser(ctx context.Context, userId uuid.UUID, name string, roleIds []uuid.UUID, loginId *uuid.UUID) (iamdb.GetUserWithRolesRow, error) {
	// Update the user's name and login ID if provided
	nullString := sql.NullString{String: name, Valid: name != ""}

	// Create update params
	updateParams := iamdb.UpdateUserParams{
		ID:   userId,
		Name: nullString,
	}

	// Update the user
	_, err := s.queries.UpdateUser(ctx, updateParams)
	if err != nil {
		return iamdb.GetUserWithRolesRow{}, err
	}

	// If loginId is provided, update the user's login ID
	if loginId != nil {
		// Create a NullUUID for the login ID
		nullUUID := uuid.NullUUID{UUID: *loginId, Valid: true}

		// Update the user's login ID
		_, err := s.queries.UpdateUserLoginID(ctx, iamdb.UpdateUserLoginIDParams{
			ID:      userId,
			LoginID: nullUUID,
		})
		if err != nil {
			return iamdb.GetUserWithRolesRow{}, fmt.Errorf("failed to update user login ID: %w", err)
		}
	}

	// Delete existing roles
	err = s.queries.DeleteUserRoles(ctx, userId)
	if err != nil {
		return iamdb.GetUserWithRolesRow{}, err
	}

	// If there are new roles to assign, create the user-role associations
	if len(roleIds) > 0 {
		// Insert role assignments one by one
		for _, roleId := range roleIds {
			_, err = s.queries.CreateUserRole(ctx, iamdb.CreateUserRoleParams{
				UserID: userId,
				RoleID: roleId,
			})
			if err != nil {
				return iamdb.GetUserWithRolesRow{}, err
			}
		}
	}

	// Return the updated user with roles
	return s.queries.GetUserWithRoles(ctx, userId)
}

func (s *IamService) DeleteUser(ctx context.Context, userId uuid.UUID) error {
	// Check if user exists
	_, err := s.queries.GetUserWithRoles(ctx, userId)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	return s.queries.DeleteUser(ctx, userId)
}
