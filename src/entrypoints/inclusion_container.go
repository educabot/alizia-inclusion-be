package entrypoints

import (
	inclusionuc "github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
)

type InclusionContainer struct {
	GetStudentProfile        inclusionuc.GetStudentProfile
	UpsertStudentProfile     inclusionuc.UpsertStudentProfile
	ListClassroomStudents    inclusionuc.ListClassroomStudents
	RecommendDevice          inclusionuc.RecommendDevice
	AssistClassroom          inclusionuc.AssistClassroom
	OpenSession              inclusionuc.OpenSession
	BuildPromptContext       inclusionuc.BuildPromptContext
	SearchPedagogicalContent inclusionuc.SearchPedagogicalContent

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
	ExportAdaptation        inclusionuc.ExportAdaptation
	GetChatHistory          inclusionuc.GetChatHistory
	DeleteConversation      inclusionuc.DeleteConversation
	RenameConversation      inclusionuc.RenameConversation
}
