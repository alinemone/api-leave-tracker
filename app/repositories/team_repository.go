package repositories

import (
	"github.com/goravel/framework/facades"

	teamresp "leave/app/http/responses/team"
	"leave/app/models"
)

type TeamRepository struct{}

func NewTeamRepository() *TeamRepository {
	return &TeamRepository{}
}

// Create makes a team and registers its creator as an approved manager member.
func (r *TeamRepository) Create(name string, creatorID uint, webhookURL string) (models.Team, error) {
	team := models.Team{Name: name, CreatorID: creatorID, WebhookURL: webhookURL}
	if err := facades.Orm().Query().Create(&team); err != nil {
		return models.Team{}, err
	}

	member := models.TeamMember{TeamID: team.ID, UserID: creatorID, IsManager: true, IsApproved: true}
	if err := facades.Orm().Query().Create(&member); err != nil {
		return models.Team{}, err
	}

	return team, nil
}

func (r *TeamRepository) Find(id uint) (models.Team, error) {
	var team models.Team
	err := facades.Orm().Query().Where("id", id).First(&team)
	return team, err
}

func (r *TeamRepository) Update(id uint, name string, webhookURL string) error {
	_, err := facades.Orm().Query().Model(&models.Team{}).Where("id", id).
		Update(map[string]any{"name": name, "webhook_url": webhookURL})
	return err
}

// Delete soft-deletes the team and removes its memberships.
func (r *TeamRepository) Delete(id uint) error {
	if _, err := facades.Orm().Query().Where("team_id", id).Delete(&models.TeamMember{}); err != nil {
		return err
	}
	_, err := facades.Orm().Query().Where("id", id).Delete(&models.Team{})
	return err
}

// MemberTeams returns approved teams the user belongs to, with role and company.
func (r *TeamRepository) MemberTeams(userID uint) ([]teamresp.TeamListRow, error) {
	var rows []teamresp.TeamListRow
	err := facades.Orm().Query().Table("team_members").
		Select(`
			teams.id AS team_id,
			teams.name AS name,
			team_members.is_manager AS is_manager
		`).
		Join("JOIN teams ON teams.id = team_members.team_id AND teams.deleted_at IS NULL").
		Where("team_members.user_id = ?", userID).
		Where("team_members.is_approved = ?", true).
		OrderBy("teams.name").
		Find(&rows)
	return rows, err
}

func (r *TeamRepository) IsManager(userID uint, teamID uint) bool {
	count, _ := facades.Orm().Query().Model(&models.TeamMember{}).
		Where("team_id", teamID).Where("user_id", userID).
		Where("is_manager", true).Where("is_approved", true).Count()
	return count > 0
}

func (r *TeamRepository) IsMember(userID uint, teamID uint) bool {
	count, _ := facades.Orm().Query().Model(&models.TeamMember{}).
		Where("team_id", teamID).Where("user_id", userID).Where("is_approved", true).Count()
	return count > 0
}

// Members returns approved members of a team.
func (r *TeamRepository) Members(teamID uint) ([]teamresp.MemberRow, error) {
	var rows []teamresp.MemberRow
	err := facades.Orm().Query().Table("team_members").
		Select(`
			users.id AS user_id,
			users.name AS name,
			users.mobile AS mobile,
			team_members.is_manager AS is_manager
		`).
		Join("JOIN users ON users.id = team_members.user_id").
		Where("team_members.team_id = ?", teamID).
		Where("team_members.is_approved = ?", true).
		OrderByDesc("team_members.is_manager").
		OrderBy("users.name").
		Find(&rows)
	return rows, err
}

// PendingRequests returns members awaiting approval for a team.
func (r *TeamRepository) PendingRequests(teamID uint) ([]teamresp.MemberRow, error) {
	var rows []teamresp.MemberRow
	err := facades.Orm().Query().Table("team_members").
		Select(`
			users.id AS user_id,
			users.name AS name,
			users.mobile AS mobile,
			team_members.is_manager AS is_manager
		`).
		Join("JOIN users ON users.id = team_members.user_id").
		Where("team_members.team_id = ?", teamID).
		Where("team_members.is_approved = ?", false).
		OrderBy("users.name").
		Find(&rows)
	return rows, err
}

// SearchAddableUsers finds registered users matching q (name/mobile) who are
// not already linked to the team.
func (r *TeamRepository) SearchAddableUsers(teamID uint, q string) ([]teamresp.MemberRow, error) {
	var rows []teamresp.MemberRow
	like := "%" + q + "%"
	err := facades.Orm().Query().Table("users").
		Select("users.id AS user_id, users.name AS name, users.mobile AS mobile").
		Where("users.name ILIKE ? OR users.mobile ILIKE ?", like, like).
		Where("users.id NOT IN (SELECT user_id FROM team_members WHERE team_id = ?)", teamID).
		OrderBy("users.name").
		Limit(10).
		Find(&rows)
	return rows, err
}

// AddMember links a user to a team as an approved member (manager action).
func (r *TeamRepository) AddMember(teamID uint, userID uint, isManager bool) error {
	var existing models.TeamMember
	_ = facades.Orm().Query().Where("team_id", teamID).Where("user_id", userID).First(&existing)
	if existing.ID != 0 {
		existing.IsManager = isManager
		existing.IsApproved = true
		return facades.Orm().Query().Save(&existing)
	}

	member := models.TeamMember{TeamID: teamID, UserID: userID, IsManager: isManager, IsApproved: true}
	return facades.Orm().Query().Create(&member)
}

// RequestJoin creates a pending membership for a user (self-service join).
func (r *TeamRepository) RequestJoin(teamID uint, userID uint) error {
	var existing models.TeamMember
	_ = facades.Orm().Query().Where("team_id", teamID).Where("user_id", userID).First(&existing)
	if existing.ID != 0 {
		return nil // already a member or already requested
	}
	member := models.TeamMember{TeamID: teamID, UserID: userID, IsManager: false, IsApproved: false}
	return facades.Orm().Query().Create(&member)
}

// ApproveMember approves a pending membership.
func (r *TeamRepository) ApproveMember(teamID uint, userID uint) error {
	_, err := facades.Orm().Query().Model(&models.TeamMember{}).
		Where("team_id", teamID).Where("user_id", userID).
		Update("is_approved", true)
	return err
}

func (r *TeamRepository) RemoveMember(teamID uint, userID uint) error {
	_, err := facades.Orm().Query().
		Where("team_id", teamID).Where("user_id", userID).
		Delete(&models.TeamMember{})
	return err
}

// SetManager promotes/demotes an approved member to/from manager.
func (r *TeamRepository) SetManager(teamID uint, userID uint, isManager bool) error {
	_, err := facades.Orm().Query().Model(&models.TeamMember{}).
		Where("team_id", teamID).Where("user_id", userID).Where("is_approved", true).
		Update("is_manager", isManager)
	return err
}

// CountManagers returns the number of approved managers on a team.
func (r *TeamRepository) CountManagers(teamID uint) int64 {
	count, _ := facades.Orm().Query().Model(&models.TeamMember{}).
		Where("team_id", teamID).Where("is_manager", true).Where("is_approved", true).Count()
	return count
}

// CancelJoinRequest lets a user withdraw their own pending join request.
func (r *TeamRepository) CancelJoinRequest(teamID uint, userID uint) error {
	_, err := facades.Orm().Query().
		Where("team_id", teamID).Where("user_id", userID).Where("is_approved", false).
		Delete(&models.TeamMember{})
	return err
}

// JoinableTeams lists teams the user is not yet an approved member of,
// flagging the ones with a pending request.
func (r *TeamRepository) JoinableTeams(userID uint) ([]teamresp.JoinRow, error) {
	approved := map[uint]bool{}
	mine, _ := r.MemberTeams(userID)
	for _, t := range mine {
		approved[t.TeamID] = true
	}

	pending := map[uint]bool{}
	var prs []models.TeamMember
	_ = facades.Orm().Query().Where("user_id", userID).Where("is_approved", false).Find(&prs)
	for _, p := range prs {
		pending[p.TeamID] = true
	}

	var all []teamresp.JoinRow
	err := facades.Orm().Query().Table("teams").
		Select(`
			teams.id AS team_id,
			teams.name AS name
		`).
		Where("teams.deleted_at IS NULL").
		OrderBy("teams.name").
		Find(&all)
	if err != nil {
		return nil, err
	}

	var rows []teamresp.JoinRow
	for _, t := range all {
		if approved[t.TeamID] {
			continue
		}
		t.Pending = pending[t.TeamID]
		rows = append(rows, t)
	}
	return rows, nil
}
