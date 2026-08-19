package db

import (
	"path/filepath"
	"testing"
)

func TestCreateOrganizationRollsBackWhenOwnerAlreadyBelongsToOne(t *testing.T) {
	database, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	user, err := CreateUser(database, "owner@example.com", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := CreateOrganization(database, user.ID, "First", "invite-first"); err != nil {
		t.Fatal(err)
	}
	if _, err := CreateOrganization(database, user.ID, "Second", "invite-second"); err == nil {
		t.Fatal("expected second organization creation to fail")
	}

	var count int
	if err := database.QueryRow(`SELECT COUNT(*) FROM organizations WHERE owner_id = ?`, user.ID).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("failed organization was not rolled back: got %d rows", count)
	}
}

func TestJoinOrganizationReportsMembershipConflict(t *testing.T) {
	database, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	ownerA, _ := CreateUser(database, "a@example.com", "hash", false)
	ownerB, _ := CreateUser(database, "b@example.com", "hash", false)
	member, _ := CreateUser(database, "member@example.com", "hash", false)
	orgA, err := CreateOrganization(database, ownerA.ID, "A", "invite-a")
	if err != nil {
		t.Fatal(err)
	}
	orgB, err := CreateOrganization(database, ownerB.ID, "B", "invite-b")
	if err != nil {
		t.Fatal(err)
	}
	if err := JoinOrganization(database, orgA.ID, member.ID); err != nil {
		t.Fatal(err)
	}
	if err := JoinOrganization(database, orgB.ID, member.ID); err == nil {
		t.Fatal("expected joining a second organization to fail")
	}
}
