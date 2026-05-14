package entrypoints

import (
	inclusionuc "github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
)

type InclusionContainer struct {
	GetStudentProfile     inclusionuc.GetStudentProfile
	UpsertStudentProfile  inclusionuc.UpsertStudentProfile
	ListClassroomStudents inclusionuc.ListClassroomStudents
	RecommendDevice       inclusionuc.RecommendDevice
	AssistClassroom       inclusionuc.AssistClassroom
}
