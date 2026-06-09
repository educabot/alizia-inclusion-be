package main

import (
	"github.com/educabot/alizia-inclusion-be/config"
	authuc "github.com/educabot/alizia-inclusion-be/src/core/usecases/auth"
	cataloguc "github.com/educabot/alizia-inclusion-be/src/core/usecases/catalog"
	dashuc "github.com/educabot/alizia-inclusion-be/src/core/usecases/dashboard"
	inclusionuc "github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	mgmtuc "github.com/educabot/alizia-inclusion-be/src/core/usecases/management"
)

type UseCases struct {
	ListRamps   cataloguc.ListRamps
	GetRamp     cataloguc.GetRamp
	ListDevices cataloguc.ListDevices
	GetDevice   cataloguc.GetDevice

	GetStudentProfile     inclusionuc.GetStudentProfile
	UpsertStudentProfile  inclusionuc.UpsertStudentProfile
	ListClassroomStudents inclusionuc.ListClassroomStudents
	RecommendDevice       inclusionuc.RecommendDevice
	AssistClassroom       inclusionuc.AssistClassroom
	OpenSession           inclusionuc.OpenSession
	BuildPromptContext    inclusionuc.BuildPromptContext

	GetMe authuc.GetMe

	ListClassrooms  mgmtuc.ListClassrooms
	GetClassroom    mgmtuc.GetClassroom
	CreateClassroom mgmtuc.CreateClassroom
	UpdateClassroom mgmtuc.UpdateClassroom
	DeleteClassroom mgmtuc.DeleteClassroom
	ListTeachers    mgmtuc.ListTeachers

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

	GetMetrics dashuc.GetMetrics
	GetAIUsage dashuc.GetAIUsage
}

func NewUseCases(repos *Repositories, cfg *config.Config) *UseCases {
	return &UseCases{
		ListRamps:   cataloguc.NewListRamps(repos.Ramps),
		GetRamp:     cataloguc.NewGetRamp(repos.Ramps),
		ListDevices: cataloguc.NewListDevices(repos.Devices),
		GetDevice:   cataloguc.NewGetDevice(repos.Devices),

		GetStudentProfile:     inclusionuc.NewGetStudentProfile(repos.Students),
		UpsertStudentProfile:  inclusionuc.NewUpsertStudentProfile(repos.Students, repos.StudentProfiles),
		ListClassroomStudents: inclusionuc.NewListClassroomStudents(repos.Students),
		RecommendDevice:       inclusionuc.NewRecommendDevice(repos.AI, repos.Students, repos.Devices, repos.Ramps, repos.Conversations, repos.AIUsage),
		AssistClassroom:       inclusionuc.NewAssistClassroom(repos.AI, repos.Students, repos.Devices, repos.Conversations, repos.ConversationSummaries, repos.Adaptations, repos.AIUsage, cfg.AIAgenticEnabled),
		OpenSession:           inclusionuc.NewOpenSession(repos.Students, repos.ConversationSummaries),
		BuildPromptContext: inclusionuc.NewBuildPromptContext(
			repos.Students, repos.TeacherProfiles, repos.Situations, repos.Diagnoses, repos.PPIs,
			repos.Adaptations, repos.Classrooms, repos.Devices, repos.ConversationSummaries,
		),

		GetMe: authuc.NewGetMe(repos.Users),

		ListClassrooms:  mgmtuc.NewListClassrooms(repos.Classrooms),
		GetClassroom:    mgmtuc.NewGetClassroom(repos.Classrooms),
		CreateClassroom: mgmtuc.NewCreateClassroom(repos.Classrooms),
		UpdateClassroom: mgmtuc.NewUpdateClassroom(repos.Classrooms),
		DeleteClassroom: mgmtuc.NewDeleteClassroom(repos.Classrooms),
		ListTeachers:    mgmtuc.NewListTeachers(repos.Users),

		ListStudents:            inclusionuc.NewListStudents(repos.Students),
		CreateStudent:           inclusionuc.NewCreateStudent(repos.Students),
		UpdateStudent:           inclusionuc.NewUpdateStudent(repos.Students),
		DeleteStudent:           inclusionuc.NewDeleteStudent(repos.Students),
		ListAdaptations:         inclusionuc.NewListAdaptations(repos.Adaptations),
		GetAdaptation:           inclusionuc.NewGetAdaptation(repos.Adaptations),
		CreateAdaptation:        inclusionuc.NewCreateAdaptation(repos.Adaptations),
		UpdateAdaptation:        inclusionuc.NewUpdateAdaptation(repos.Adaptations),
		DeleteAdaptation:        inclusionuc.NewDeleteAdaptation(repos.Adaptations),
		ListAdaptationResources: inclusionuc.NewListAdaptationResources(repos.AdaptationResources),
		ExportAdaptation:        inclusionuc.NewExportAdaptation(repos.Adaptations),
		GetChatHistory:          inclusionuc.NewGetChatHistory(repos.Conversations),

		GetMetrics: dashuc.NewGetMetrics(repos.Students, repos.Adaptations, repos.Classrooms),
		GetAIUsage: dashuc.NewGetAIUsage(repos.AIUsage),
	}
}
