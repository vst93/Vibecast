package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"vibecast/internal/db"
)

// generateInviteCode generates a random 12-character alphanumeric invite code.
func generateInviteCode() string {
	return randomSuffix(12)
}

// handleOrg handles /api/org — GET: get current user's org info, POST: create org, DELETE: delete org
func (s *Server) handleOrg(w http.ResponseWriter, r *http.Request, user *db.User) {
	switch r.Method {
	case http.MethodGet:
		s.getMyOrg(w, r, user)
	case http.MethodPost:
		s.createOrg(w, r, user)
	case http.MethodDelete:
		s.deleteOrg(w, r, user)
	default:
		writeJSON(w, 405, jsonResp{Error: tMsg(r, "method_not_allowed")})
	}
}

// getMyOrg returns the current user's organization info (or null if none).
func (s *Server) getMyOrg(w http.ResponseWriter, r *http.Request, user *db.User) {
	org, err := db.GetUserOrganization(s.database, user.ID)
	if err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "internal_error")})
		return
	}
	if org == nil {
		writeJSON(w, 200, jsonResp{Data: map[string]interface{}{"hasOrg": false}})
		return
	}
	memberCount, _ := db.CountOrgMembers(s.database, org.ID)
	writeJSON(w, 200, jsonResp{Data: map[string]interface{}{
		"hasOrg":      true,
		"id":          org.ID,
		"name":        org.Name,
		"inviteCode":  org.InviteCode,
		"isOwner":     org.OwnerID == user.ID,
		"memberCount": memberCount,
	}})
}

// createOrg creates a new organization. The user must not already be in an org.
func (s *Server) createOrg(w http.ResponseWriter, r *http.Request, user *db.User) {
	// Check if user already has an org
	existing, err := db.GetUserOrganization(s.database, user.ID)
	if err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "internal_error")})
		return
	}
	if existing != nil {
		writeJSON(w, 403, jsonResp{Error: tMsg(r, "already_in_org")})
		return
	}

	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, 400, jsonResp{Error: tMsg(r, "invalid_json")})
		return
	}
	body.Name = strings.TrimSpace(body.Name)

	// Generate unique invite code
	var inviteCode string
	for i := 0; i < 30; i++ {
		inviteCode = generateInviteCode()
		ex, _ := db.GetOrganizationByInviteCode(s.database, inviteCode)
		if ex == nil {
			break
		}
		inviteCode = ""
	}
	if inviteCode == "" {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "internal_error")})
		return
	}

	org, err := db.CreateOrganization(s.database, user.ID, body.Name, inviteCode)
	if err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "create_org_failed")})
		return
	}

	writeJSON(w, 201, jsonResp{
		Message: "organization created",
		Data: map[string]interface{}{
			"id":         org.ID,
			"name":       org.Name,
			"inviteCode": org.InviteCode,
			"isOwner":    true,
		},
	})
}

// handleOrgAction handles /api/org/{action} sub-routes.
func (s *Server) handleOrgAction(w http.ResponseWriter, r *http.Request, user *db.User) {
	pathParts := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/api/org/"), "/", 2)
	action := pathParts[0]

	switch action {
	case "join":
		if r.Method != http.MethodPost {
			writeJSON(w, 405, jsonResp{Error: tMsg(r, "method_not_allowed")})
			return
		}
		s.joinOrg(w, r, user)
	case "leave":
		if r.Method != http.MethodPost {
			writeJSON(w, 405, jsonResp{Error: tMsg(r, "method_not_allowed")})
			return
		}
		s.leaveOrg(w, r, user)
	case "members":
		// GET /api/org/members — list members
		// DELETE /api/org/members/{userId} — remove member
		if len(pathParts) > 1 && pathParts[1] != "" {
			if r.Method != http.MethodDelete {
				writeJSON(w, 405, jsonResp{Error: tMsg(r, "method_not_allowed")})
				return
			}
			s.removeOrgMember(w, r, user, pathParts[1])
			return
		}
		if r.Method != http.MethodGet {
			writeJSON(w, 405, jsonResp{Error: tMsg(r, "method_not_allowed")})
			return
		}
		s.listOrgMembers(w, r, user)
	default:
		writeJSON(w, 404, jsonResp{Error: "not found"})
	}
}

// joinOrg lets a user join an org by invite code.
func (s *Server) joinOrg(w http.ResponseWriter, r *http.Request, user *db.User) {
	// Check if user already has an org
	existing, err := db.GetUserOrganization(s.database, user.ID)
	if err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "internal_error")})
		return
	}
	if existing != nil {
		writeJSON(w, 403, jsonResp{Error: tMsg(r, "already_in_org")})
		return
	}

	var body struct {
		InviteCode string `json:"inviteCode"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, 400, jsonResp{Error: tMsg(r, "invalid_json")})
		return
	}
	body.InviteCode = strings.TrimSpace(body.InviteCode)
	if body.InviteCode == "" {
		writeJSON(w, 400, jsonResp{Error: tMsg(r, "invite_code_required")})
		return
	}

	org, err := db.GetOrganizationByInviteCode(s.database, body.InviteCode)
	if err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "internal_error")})
		return
	}
	if org == nil {
		writeJSON(w, 404, jsonResp{Error: tMsg(r, "org_not_found")})
		return
	}

	if err := db.JoinOrganization(s.database, org.ID, user.ID); err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "join_org_failed")})
		return
	}

	writeJSON(w, 200, jsonResp{
		Message: "joined",
		Data: map[string]interface{}{
			"id":         org.ID,
			"name":       org.Name,
			"inviteCode": org.InviteCode,
			"isOwner":    false,
		},
	})
}

// leaveOrg lets a non-owner user leave their org.
func (s *Server) leaveOrg(w http.ResponseWriter, r *http.Request, user *db.User) {
	org, err := db.GetUserOrganization(s.database, user.ID)
	if err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "internal_error")})
		return
	}
	if org == nil {
		writeJSON(w, 400, jsonResp{Error: tMsg(r, "not_in_org")})
		return
	}
	if org.OwnerID == user.ID {
		writeJSON(w, 403, jsonResp{Error: tMsg(r, "owner_cannot_leave")})
		return
	}

	if err := db.LeaveOrganization(s.database, org.ID, user.ID); err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "leave_org_failed")})
		return
	}

	writeJSON(w, 200, jsonResp{Message: "left organization"})
}

// deleteOrg deletes an organization. Only the owner can delete, and only if no other members.
func (s *Server) deleteOrg(w http.ResponseWriter, r *http.Request, user *db.User) {
	org, err := db.GetUserOrganization(s.database, user.ID)
	if err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "internal_error")})
		return
	}
	if org == nil {
		writeJSON(w, 400, jsonResp{Error: tMsg(r, "not_in_org")})
		return
	}
	if org.OwnerID != user.ID {
		writeJSON(w, 403, jsonResp{Error: tMsg(r, "not_org_owner")})
		return
	}

	// Check no other members (owner is in org_members too, so count > 1 means others exist)
	memberCount, _ := db.CountOrgMembers(s.database, org.ID)
	if memberCount > 1 {
		writeJSON(w, 403, jsonResp{Error: tMsg(r, "org_has_members")})
		return
	}

	if err := db.DeleteOrganization(s.database, org.ID); err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "delete_org_failed")})
		return
	}

	writeJSON(w, 200, jsonResp{Message: "organization deleted"})
}

// listOrgMembers returns a paginated, searchable list of org members.
func (s *Server) listOrgMembers(w http.ResponseWriter, r *http.Request, user *db.User) {
	org, err := db.GetUserOrganization(s.database, user.ID)
	if err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "internal_error")})
		return
	}
	if org == nil {
		writeJSON(w, 400, jsonResp{Error: tMsg(r, "not_in_org")})
		return
	}

	page, perPage, offset, search := paginationParams(r)

	members, err := db.ListOrgMembersPaged(s.database, org.ID, search, perPage, offset)
	if err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "internal_error")})
		return
	}
	total, err := db.CountOrgMembersWithSearch(s.database, org.ID, search)
	if err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "internal_error")})
		return
	}

	var list []map[string]interface{}
	for _, m := range members {
		list = append(list, map[string]interface{}{
			"id":       m.ID,
			"userId":   m.UserID,
			"email":    m.Email,
			"isOwner":  m.IsOwner,
			"joinedAt": m.JoinedAt,
		})
	}
	if list == nil {
		list = []map[string]interface{}{}
	}

	writeJSON(w, 200, jsonResp{Data: map[string]interface{}{
		"items":   list,
		"total":   total,
		"page":    page,
		"perPage": perPage,
	}})
}

// removeOrgMember removes a member from the org. Only the owner can do this.
func (s *Server) removeOrgMember(w http.ResponseWriter, r *http.Request, user *db.User, targetUserIDStr string) {
	targetUserID, err := strconv.ParseInt(targetUserIDStr, 10, 64)
	if err != nil {
		writeJSON(w, 400, jsonResp{Error: tMsg(r, "invalid_user_id")})
		return
	}

	org, err := db.GetUserOrganization(s.database, user.ID)
	if err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "internal_error")})
		return
	}
	if org == nil {
		writeJSON(w, 400, jsonResp{Error: tMsg(r, "not_in_org")})
		return
	}
	if org.OwnerID != user.ID {
		writeJSON(w, 403, jsonResp{Error: tMsg(r, "not_org_owner")})
		return
	}
	if targetUserID == user.ID {
		writeJSON(w, 403, jsonResp{Error: tMsg(r, "cannot_remove_self")})
		return
	}

	// Verify target is in the same org
	isMember, err := db.IsOrgMember(s.database, org.ID, targetUserID)
	if err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "internal_error")})
		return
	}
	if !isMember {
		writeJSON(w, 404, jsonResp{Error: tMsg(r, "user_not_found")})
		return
	}

	if err := db.RemoveOrgMember(s.database, org.ID, targetUserID); err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "internal_error")})
		return
	}

	writeJSON(w, 200, jsonResp{Message: "member removed"})
}

// sameOrgAsSiteOwner checks if the current user is in the same org as the site owner.
// Returns true only if both are in the same org and the site has org_open enabled.
func (s *Server) sameOrgAsSiteOwner(site *db.Site, userID int64) bool {
	if !site.OrgOpen {
		return false
	}
	// Get site owner's org
	ownerOrg, err := db.GetUserOrganization(s.database, site.UserID)
	if err != nil || ownerOrg == nil {
		return false
	}
	// Check if current user is in the same org
	isMember, err := db.IsOrgMember(s.database, ownerOrg.ID, userID)
	if err != nil {
		return false
	}
	return isMember
}

var _ = fmt.Sprintf // keep fmt import
