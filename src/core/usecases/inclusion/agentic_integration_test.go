//go:build integration

// Test e2e del Caso 2 + privacidad por docente contra una Postgres real.
//
// Ejercita la lógica nueva a través de los repositorios reales (SQL real) y del
// inclusionDispatcher, sin LLM ni JWT, de forma determinística:
//   - adaptaciones privadas del docente (filtro teacher_id)
//   - notas privadas del docente (filtro user_id; legacy NULL invisible)
//   - create_classroom (normaliza "tercero B" -> 3ro B; idempotente)
//   - create_student (classroom_id; idempotente dentro del aula)
//   - find_student_by_name (busca en toda la org)
//   - get_past_adaptations (solo las del docente del turno)
//
// Uso:
//   INTEGRATION_DATABASE_URL="postgres://postgres:postgres@localhost:5481/alizia_inclusion?sslmode=disable" \
//     go test -tags integration ./src/core/usecases/inclusion/ -run Integration -v
package inclusion

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	repoinc "github.com/educabot/alizia-inclusion-be/src/repositories/inclusion"
	repomgmt "github.com/educabot/alizia-inclusion-be/src/repositories/management"
)

func openTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("INTEGRATION_DATABASE_URL")
	if dsn == "" {
		t.Skip("INTEGRATION_DATABASE_URL no seteada; salteando test de integración")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	require.NoError(t, err)
	return db
}

// seedTeacher inserta un docente (users) y devuelve su id.
func seedTeacher(t *testing.T, db *gorm.DB, orgID uuid.UUID, email string) int64 {
	t.Helper()
	var id int64
	err := db.Raw(
		`INSERT INTO users (organization_id, email, name, password_hash, role)
		 VALUES (?, ?, ?, '', 'teacher') RETURNING id`,
		orgID, email, email,
	).Scan(&id).Error
	require.NoError(t, err)
	return id
}

func dispatch(t *testing.T, d inclusionDispatcher, orgID uuid.UUID, name, args string) string {
	t.Helper()
	out, err := d.Dispatch(context.Background(), orgID, providers.ToolCall{Name: name, Arguments: args})
	require.NoError(t, err)
	return out
}

func TestIntegration_Caso2AndTeacherPrivacy(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()
	orgID := uuid.New() // org aislada por test
	require.NoError(t, db.Exec(`INSERT INTO organizations (id, name) VALUES (?, ?)`, orgID, "e2e org").Error)

	students := repoinc.NewStudentRepo(db)
	profiles := repoinc.NewStudentProfileRepo(db)
	classrooms := repomgmt.NewClassroomRepo(db)
	adaptations := repoinc.NewAdaptationRepo(db)
	notes := repoinc.NewStudentNoteRepo(db)

	teacherA := seedTeacher(t, db, orgID, "a+"+orgID.String()+"@e2e.test")
	teacherB := seedTeacher(t, db, orgID, "b+"+orgID.String()+"@e2e.test")

	dispA := inclusionDispatcher{students: students, profiles: profiles, classrooms: classrooms, adaptations: adaptations, userID: teacherA}

	// ---- create_classroom: normaliza "tercero B" -> "3ro B" + idempotente ----
	out := dispatch(t, dispA, orgID, "create_classroom", `{"grade":"tercero B"}`)
	var room entities.Classroom
	require.NoError(t, json.Unmarshal([]byte(out), &room))
	assert.Equal(t, "3ro B", room.Name)
	require.NotNil(t, room.Grade)
	assert.Equal(t, "3ro", *room.Grade)
	require.NotNil(t, room.Section)
	assert.Equal(t, "B", *room.Section)
	assert.Positive(t, room.ID)

	out2 := dispatch(t, dispA, orgID, "create_classroom", `{"grade":"3ro B"}`)
	var roomAgain entities.Classroom
	require.NoError(t, json.Unmarshal([]byte(out2), &roomAgain))
	assert.Equal(t, room.ID, roomAgain.ID, "create_classroom debe ser idempotente por nombre")

	// ---- create_student: en el aula + perfil; idempotente dentro del aula ----
	createArgs := `{"name":"Lucas Pérez","classroom_id":` + itoa(room.ID) + `,"difficulties":["le cuesta sostener la atención"]}`
	out = dispatch(t, dispA, orgID, "create_student", createArgs)
	var lucas entities.Student
	require.NoError(t, json.Unmarshal([]byte(out), &lucas))
	assert.Positive(t, lucas.ID)
	assert.Equal(t, room.ID, lucas.ClassroomID)

	// mismo nombre (normalizado, con espacios/acentos) -> devuelve el existente
	out = dispatch(t, dispA, orgID, "create_student", `{"name":"  lucas perez ","classroom_id":`+itoa(room.ID)+`}`)
	var lucasAgain entities.Student
	require.NoError(t, json.Unmarshal([]byte(out), &lucasAgain))
	assert.Equal(t, lucas.ID, lucasAgain.ID, "create_student debe ser idempotente dentro del aula")

	// verificación directa en DB: un solo Lucas en el aula
	var lucasCount int64
	require.NoError(t, db.Model(&entities.Student{}).
		Where("organization_id = ? AND classroom_id = ?", orgID, room.ID).Count(&lucasCount).Error)
	assert.Equal(t, int64(1), lucasCount)

	// ---- find_student_by_name: encuentra en toda la org ----
	out = dispatch(t, dispA, orgID, "find_student_by_name", `{"name":"lucas"}`)
	var found struct {
		Students []struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"students"`
	}
	require.NoError(t, json.Unmarshal([]byte(out), &found))
	require.Len(t, found.Students, 1)
	assert.Equal(t, lucas.ID, found.Students[0].ID)

	// ---- adaptaciones privadas: A ve solo las suyas ----
	require.NoError(t, adaptations.Create(ctx, &entities.Adaptation{
		OrganizationID: orgID, StudentID: &lucas.ID, TeacherID: teacherA, Subject: "Lengua", Status: "en_curso",
	}))
	require.NoError(t, adaptations.Create(ctx, &entities.Adaptation{
		OrganizationID: orgID, StudentID: &lucas.ID, TeacherID: teacherB, Subject: "Matemática", Status: "en_curso",
	}))

	listA, err := adaptations.List(ctx, providers.AdaptationFilter{OrgID: orgID, TeacherID: &teacherA})
	require.NoError(t, err)
	require.Len(t, listA, 1)
	assert.Equal(t, "Lengua", listA[0].Subject)

	listAll, err := adaptations.List(ctx, providers.AdaptationFilter{OrgID: orgID})
	require.NoError(t, err)
	assert.Len(t, listAll, 2, "sin filtro de docente se ven las dos (caso dashboard)")

	// get_past_adaptations (turno de A) -> solo Lengua
	out = dispatch(t, dispA, orgID, "get_past_adaptations", `{"student_id":`+itoa(lucas.ID)+`}`)
	var past struct {
		Adaptations []struct {
			Subject string `json:"subject"`
		} `json:"adaptations"`
	}
	require.NoError(t, json.Unmarshal([]byte(out), &past))
	require.Len(t, past.Adaptations, 1)
	assert.Equal(t, "Lengua", past.Adaptations[0].Subject)

	// ---- notas privadas: A ve solo las suyas; legacy (user_id NULL) invisible ----
	require.NoError(t, notes.Create(ctx, &entities.StudentNote{OrganizationID: orgID, StudentID: lucas.ID, UserID: teacherA, Content: "nota de A", Type: "seguimiento"}))
	require.NoError(t, notes.Create(ctx, &entities.StudentNote{OrganizationID: orgID, StudentID: lucas.ID, UserID: teacherB, Content: "nota de B", Type: "seguimiento"}))
	// nota legacy sin dueño (user_id NULL) insertada a mano
	require.NoError(t, db.Exec(
		`INSERT INTO student_notes (student_id, organization_id, content, type, internal) VALUES (?, ?, 'legacy', 'seguimiento', true)`,
		lucas.ID, orgID,
	).Error)

	notesA, err := notes.ListByStudent(ctx, orgID, lucas.ID, teacherA)
	require.NoError(t, err)
	require.Len(t, notesA, 1, "A solo ve su nota (ni la de B ni la legacy NULL)")
	assert.Equal(t, "nota de A", notesA[0].Content)

	t.Cleanup(func() { cleanupOrg(db, orgID) })
}

func itoa(n int64) string { return strconv.FormatInt(n, 10) }

func cleanupOrg(db *gorm.DB, orgID uuid.UUID) {
	db.Exec(`DELETE FROM student_notes WHERE organization_id = ?`, orgID)
	db.Exec(`DELETE FROM adaptations WHERE organization_id = ?`, orgID)
	db.Exec(`DELETE FROM student_profiles WHERE student_id IN (SELECT id FROM students WHERE organization_id = ?)`, orgID)
	db.Exec(`DELETE FROM students WHERE organization_id = ?`, orgID)
	db.Exec(`DELETE FROM classrooms WHERE organization_id = ?`, orgID)
	db.Exec(`DELETE FROM users WHERE organization_id = ?`, orgID)
	db.Exec(`DELETE FROM organizations WHERE id = ?`, orgID)
}
