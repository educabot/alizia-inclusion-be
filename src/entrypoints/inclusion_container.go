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

	ListStudents            inclusionuc.ListStudents
	CreateStudent           inclusionuc.CreateStudent
	UpdateStudent           inclusionuc.UpdateStudent
	DeleteStudent           inclusionuc.DeleteStudent
	ListAdaptations         inclusionuc.ListAdaptations
	GetAdaptation           inclusionuc.GetAdaptation
	CreateAdaptation        inclusionuc.CreateAdaptation
	UpdateAdaptation        inclusionuc.UpdateAdaptation
	DeleteAdaptation        inclusionuc.DeleteAdaptation
	ListAdaptationResources inclusionuc.ListAdaptationResources
	GetChatHistory          inclusionuc.GetChatHistory
}
