//go:build integration

package inclusion_test

import (
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
)

// seedClassroom inserts a classroom in TestOrgID and returns its id.
func seedClassroom(t *testing.T, tx *gorm.DB, name string) int64 {
	t.Helper()
	c := entities.Classroom{OrganizationID: testutil.TestOrgID, Name: name}
	require.NoError(t, tx.Create(&c).Error)
	return c.ID
}

// seedStudent inserts a student (creating a classroom on the fly) and returns its id.
func seedStudent(t *testing.T, tx *gorm.DB, name string) int64 {
	t.Helper()
	classroomID := seedClassroom(t, tx, name+"-aula")
	s := entities.Student{OrganizationID: testutil.TestOrgID, ClassroomID: classroomID, Name: name}
	require.NoError(t, tx.Create(&s).Error)
	return s.ID
}

// seedStudentProfile inserts a profile for a student and returns the profile id.
// difficulties is normalised to a non-nil slice: the column is NOT NULL DEFAULT '{}'
// and pq.StringArray(nil) would serialise to NULL.
func seedStudentProfile(t *testing.T, tx *gorm.DB, studentID int64, difficulties ...string) int64 {
	t.Helper()
	p := entities.StudentProfile{StudentID: studentID, Difficulties: append(pq.StringArray{}, difficulties...)}
	require.NoError(t, tx.Create(&p).Error)
	return p.ID
}

// seedUser inserts a user in TestOrgID and returns its id.
func seedUser(t *testing.T, tx *gorm.DB, name, role string) int64 {
	t.Helper()
	u := entities.User{OrganizationID: testutil.TestOrgID, Email: name + "@test.com", Name: name, Role: role}
	require.NoError(t, tx.Create(&u).Error)
	return u.ID
}
