package services

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"lol-match-exporter/internal/models"
)

// GroupService gère les opérations relatives aux groupes
type GroupService struct {
	db *sql.DB
}

// NewGroupService crée une nouvelle instance du service de groupes
func NewGroupService(db *sql.DB) *GroupService {
	return &GroupService{
		db: db,
	}
}

// CreateGroup crée un nouveau groupe
func (gs *GroupService) CreateGroup(ownerID int, name, description string, privacy string) (*models.Group, error) {
	// Générer un code d'invitation unique
	inviteCode, err := gs.generateInviteCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate invite code: %w", err)
	}

	// Valider les paramètres
	if name == "" {
		return nil, fmt.Errorf("group name cannot be empty")
	}
	
	if privacy == "" {
		privacy = "private"
	}
	
	validPrivacyLevels := []string{"public", "private", "invite_only"}
	if !contains(validPrivacyLevels, privacy) {
		return nil, fmt.Errorf("invalid privacy level: %s", privacy)
	}

	// Créer le groupe
	defaultSettings := models.GetDefaultGroupSettings()
	now := time.Now()
	
	query := `
		INSERT INTO groups (name, description, owner_id, privacy, invite_code, settings, member_count, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, 1, ?, ?)
	`
	
	result, err := gs.db.Exec(query, name, description, ownerID, privacy, inviteCode, defaultSettings, now, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create group: %w", err)
	}
	
	groupID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get group ID: %w", err)
	}

	// Ajouter le propriétaire comme membre avec le rôle "owner"
	err = gs.AddMemberToGroup(int(groupID), ownerID, "owner")
	if err != nil {
		return nil, fmt.Errorf("failed to add owner to group: %w", err)
	}

	// Récupérer et retourner le groupe complet
	return gs.GetGroupByID(int(groupID))
}

// GetGroupByID récupère un groupe par son ID
func (gs *GroupService) GetGroupByID(groupID int) (*models.Group, error) {
	query := `
		SELECT g.id, g.name, g.description, g.owner_id, g.privacy, g.invite_code, 
		       g.settings, g.member_count, g.created_at, g.updated_at,
		       u.id, u.riot_id, u.riot_tag, u.region
		FROM groups g
		LEFT JOIN users u ON g.owner_id = u.id
		WHERE g.id = ?
	`
	
	row := gs.db.QueryRow(query, groupID)
	
	var group models.Group
	var owner models.User
	
	err := row.Scan(
		&group.ID, &group.Name, &group.Description, &group.OwnerID, &group.Privacy,
		&group.InviteCode, &group.Settings, &group.MemberCount, &group.CreatedAt, &group.UpdatedAt,
		&owner.ID, &owner.RiotID, &owner.RiotTag, &owner.Region,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("group not found")
		}
		return nil, fmt.Errorf("failed to get group: %w", err)
	}
	
	group.Owner = &owner
	
	// Récupérer les membres du groupe
	members, err := gs.GetGroupMembers(groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group members: %w", err)
	}
	group.Members = members
	
	return &group, nil
}

// GetGroupsByUserID récupère tous les groupes d'un utilisateur
func (gs *GroupService) GetGroupsByUserID(userID int) ([]models.Group, error) {
	query := `
		SELECT DISTINCT g.id, g.name, g.description, g.owner_id, g.privacy, g.invite_code,
		       g.settings, g.member_count, g.created_at, g.updated_at,
		       u.riot_id, u.riot_tag, u.region, u.profile_icon_id,
		       gm.role, gm.status
		FROM groups g
		JOIN group_members gm ON g.id = gm.group_id
		LEFT JOIN users u ON g.owner_id = u.id
		WHERE gm.user_id = ? AND gm.status = 'active'
		ORDER BY g.updated_at DESC
	`
	
	rows, err := gs.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user groups: %w", err)
	}
	defer rows.Close()
	
	var groups []models.Group
	
	for rows.Next() {
		var group models.Group
		var owner models.User
		var memberRole, memberStatus string
		
		err := rows.Scan(
			&group.ID, &group.Name, &group.Description, &group.OwnerID, &group.Privacy,
			&group.InviteCode, &group.Settings, &group.MemberCount, &group.CreatedAt, &group.UpdatedAt,
			&owner.RiotID, &owner.RiotTag, &owner.Region, &owner.ProfileIconID,
			&memberRole, &memberStatus,
		)
		
		if err != nil {
			return nil, fmt.Errorf("failed to scan group: %w", err)
		}
		
		group.Owner = &owner
		groups = append(groups, group)
	}
	
	return groups, nil
}

// AddMemberToGroup ajoute un membre à un groupe
func (gs *GroupService) AddMemberToGroup(groupID, userID int, role string) error {
	if role == "" {
		role = "member"
	}
	
	validRoles := []string{"owner", "admin", "member"}
	if !contains(validRoles, role) {
		return fmt.Errorf("invalid role: %s", role)
	}
	
	permissions := models.GetMemberPermissions(role)
	now := time.Now()
	
	query := `
		INSERT INTO group_members (group_id, user_id, role, status, joined_at, permissions)
		VALUES (?, ?, ?, 'active', ?, ?)
	`
	
	_, err := gs.db.Exec(query, groupID, userID, role, now, permissions)
	if err != nil {
		return fmt.Errorf("failed to add member to group: %w", err)
	}
	
	// Mettre à jour le compteur de membres du groupe
	err = gs.updateGroupMemberCount(groupID)
	if err != nil {
		return fmt.Errorf("failed to update group member count: %w", err)
	}
	
	return nil
}

// GetGroupMembers récupère tous les membres d'un groupe
func (gs *GroupService) GetGroupMembers(groupID int) ([]models.GroupMember, error) {
	query := `
		SELECT gm.id, gm.group_id, gm.user_id, gm.role, gm.status, gm.joined_at, 
		       gm.nickname, gm.permissions,
		       u.riot_id, u.riot_tag, u.region, u.profile_icon_id, u.summoner_level
		FROM group_members gm
		JOIN users u ON gm.user_id = u.id
		WHERE gm.group_id = ? AND gm.status = 'active'
		ORDER BY gm.role DESC, gm.joined_at ASC
	`
	
	rows, err := gs.db.Query(query, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group members: %w", err)
	}
	defer rows.Close()
	
	var members []models.GroupMember
	
	for rows.Next() {
		var member models.GroupMember
		var user models.User
		
		err := rows.Scan(
			&member.ID, &member.GroupID, &member.UserID, &member.Role, &member.Status,
			&member.JoinedAt, &member.Nickname, &member.Permissions,
			&user.RiotID, &user.RiotTag, &user.Region, &user.ProfileIconID,
			&user.SummonerLevel,
		)
		
		if err != nil {
			return nil, fmt.Errorf("failed to scan group member: %w", err)
		}
		
		member.User = &user
		members = append(members, member)
	}
	
	return members, nil
}

// RemoveMemberFromGroup retire un membre d'un groupe
func (gs *GroupService) RemoveMemberFromGroup(groupID, userID int) error {
	// Vérifier que l'utilisateur n'est pas le propriétaire
	group, err := gs.GetGroupByID(groupID)
	if err != nil {
		return fmt.Errorf("failed to get group: %w", err)
	}
	
	if group.OwnerID == userID {
		return fmt.Errorf("cannot remove group owner")
	}
	
	query := `
		UPDATE group_members 
		SET status = 'removed', updated_at = ?
		WHERE group_id = ? AND user_id = ?
	`
	
	_, err = gs.db.Exec(query, time.Now(), groupID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove member from group: %w", err)
	}
	
	// Mettre à jour le compteur de membres
	return gs.updateGroupMemberCount(groupID)
}

// CreateGroupInvite crée une invitation pour un groupe
func (gs *GroupService) CreateGroupInvite(groupID, inviterID int, email string, message string) (*models.GroupInvite, error) {
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // Expire dans 7 jours
	now := time.Now()
	
	query := `
		INSERT INTO group_invites (group_id, inviter_id, email, status, message, expires_at, created_at, updated_at)
		VALUES (?, ?, ?, 'pending', ?, ?, ?, ?)
	`
	
	result, err := gs.db.Exec(query, groupID, inviterID, email, message, expiresAt, now, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create group invite: %w", err)
	}
	
	inviteID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get invite ID: %w", err)
	}
	
	return gs.GetGroupInviteByID(int(inviteID))
}

// GetGroupInviteByID récupère une invitation par son ID
func (gs *GroupService) GetGroupInviteByID(inviteID int) (*models.GroupInvite, error) {
	query := `
		SELECT gi.id, gi.group_id, gi.inviter_id, gi.invitee_id, gi.email, 
		       gi.status, gi.message, gi.expires_at, gi.created_at, gi.updated_at,
		       g.name, g.description, g.privacy,
		       u.riot_id, u.riot_tag
		FROM group_invites gi
		JOIN groups g ON gi.group_id = g.id
		JOIN users u ON gi.inviter_id = u.id
		WHERE gi.id = ?
	`
	
	row := gs.db.QueryRow(query, inviteID)
	
	var invite models.GroupInvite
	var group models.Group
	var inviter models.User
	
	err := row.Scan(
		&invite.ID, &invite.GroupID, &invite.InviterID, &invite.InviteeID, &invite.Email,
		&invite.Status, &invite.Message, &invite.ExpiresAt, &invite.CreatedAt, &invite.UpdatedAt,
		&group.Name, &group.Description, &group.Privacy,
		&inviter.RiotID, &inviter.RiotTag,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invite not found")
		}
		return nil, fmt.Errorf("failed to get invite: %w", err)
	}
	
	invite.Group = &group
	invite.Inviter = &inviter
	
	return &invite, nil
}

// AcceptGroupInvite accepte une invitation de groupe
func (gs *GroupService) AcceptGroupInvite(inviteID, userID int) error {
	invite, err := gs.GetGroupInviteByID(inviteID)
	if err != nil {
		return fmt.Errorf("failed to get invite: %w", err)
	}
	
	// Vérifier que l'invitation est valide
	if invite.Status != "pending" {
		return fmt.Errorf("invite is not pending")
	}
	
	if time.Now().After(invite.ExpiresAt) {
		return fmt.Errorf("invite has expired")
	}
	
	// Ajouter l'utilisateur au groupe
	err = gs.AddMemberToGroup(invite.GroupID, userID, "member")
	if err != nil {
		return fmt.Errorf("failed to add user to group: %w", err)
	}
	
	// Marquer l'invitation comme acceptée
	query := `
		UPDATE group_invites 
		SET status = 'accepted', invitee_id = ?, updated_at = ?
		WHERE id = ?
	`
	
	_, err = gs.db.Exec(query, userID, time.Now(), inviteID)
	if err != nil {
		return fmt.Errorf("failed to update invite status: %w", err)
	}
	
	return nil
}

// UpdateGroupSettings met à jour les paramètres d'un groupe
func (gs *GroupService) UpdateGroupSettings(groupID int, settings models.GroupSettings) error {
	query := `
		UPDATE groups 
		SET settings = ?, updated_at = ?
		WHERE id = ?
	`
	
	_, err := gs.db.Exec(query, settings, time.Now(), groupID)
	if err != nil {
		return fmt.Errorf("failed to update group settings: %w", err)
	}
	
	return nil
}

// SearchPublicGroups recherche des groupes publics
func (gs *GroupService) SearchPublicGroups(searchTerm string, limit int) ([]models.Group, error) {
	if limit <= 0 {
		limit = 20
	}
	
	query := `
		SELECT g.id, g.name, g.description, g.owner_id, g.privacy, g.invite_code,
		       g.settings, g.member_count, g.created_at, g.updated_at,
		       u.riot_id, u.riot_tag, u.region, u.profile_icon_id
		FROM groups g
		LEFT JOIN users u ON g.owner_id = u.id
		WHERE g.privacy = 'public' 
		  AND (g.name LIKE ? OR g.description LIKE ?)
		ORDER BY g.member_count DESC, g.updated_at DESC
		LIMIT ?
	`
	
	searchPattern := "%" + searchTerm + "%"
	rows, err := gs.db.Query(query, searchPattern, searchPattern, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search public groups: %w", err)
	}
	defer rows.Close()
	
	var groups []models.Group
	
	for rows.Next() {
		var group models.Group
		var owner models.User
		
		err := rows.Scan(
			&group.ID, &group.Name, &group.Description, &group.OwnerID, &group.Privacy,
			&group.InviteCode, &group.Settings, &group.MemberCount, &group.CreatedAt, &group.UpdatedAt,
			&owner.RiotID, &owner.RiotTag, &owner.Region, &owner.ProfileIconID,
		)
		
		if err != nil {
			return nil, fmt.Errorf("failed to scan group: %w", err)
		}
		
		group.Owner = &owner
		groups = append(groups, group)
	}
	
	return groups, nil
}

// JoinGroupByInviteCode permet de rejoindre un groupe via son code d'invitation
func (gs *GroupService) JoinGroupByInviteCode(inviteCode string, userID int) error {
	// Récupérer le groupe par son code d'invitation
	query := `SELECT id FROM groups WHERE invite_code = ?`
	
	var groupID int
	err := gs.db.QueryRow(query, inviteCode).Scan(&groupID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("invalid invite code")
		}
		return fmt.Errorf("failed to find group: %w", err)
	}
	
	// Vérifier que l'utilisateur n'est pas déjà membre
	existingMemberQuery := `
		SELECT COUNT(*) FROM group_members 
		WHERE group_id = ? AND user_id = ? AND status = 'active'
	`
	
	var memberCount int
	err = gs.db.QueryRow(existingMemberQuery, groupID, userID).Scan(&memberCount)
	if err != nil {
		return fmt.Errorf("failed to check membership: %w", err)
	}
	
	if memberCount > 0 {
		return fmt.Errorf("user is already a member of this group")
	}
	
	// Ajouter l'utilisateur au groupe
	return gs.AddMemberToGroup(groupID, userID, "member")
}

// Helper functions

// updateGroupMemberCount met à jour le compteur de membres d'un groupe
func (gs *GroupService) updateGroupMemberCount(groupID int) error {
	query := `
		UPDATE groups 
		SET member_count = (
			SELECT COUNT(*) FROM group_members 
			WHERE group_id = ? AND status = 'active'
		),
		updated_at = ?
		WHERE id = ?
	`
	
	_, err := gs.db.Exec(query, groupID, time.Now(), groupID)
	return err
}

// generateInviteCode génère un code d'invitation unique
func (gs *GroupService) generateInviteCode() (string, error) {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	
	return strings.ToUpper(hex.EncodeToString(bytes)), nil
}

// contains vérifie si un slice contient une valeur
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}