package main

import (
	cataloguc "github.com/educabot/alizia-inclusion-be/src/core/usecases/catalog"
	inclusionuc "github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
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
}

func NewUseCases(repos *Repositories) *UseCases {
	return &UseCases{
		ListRamps:   cataloguc.NewListRamps(repos.Ramps),
		GetRamp:     cataloguc.NewGetRamp(repos.Ramps),
		ListDevices: cataloguc.NewListDevices(repos.Devices),
		GetDevice:   cataloguc.NewGetDevice(repos.Devices),

		GetStudentProfile:     inclusionuc.NewGetStudentProfile(repos.Students),
		UpsertStudentProfile:  inclusionuc.NewUpsertStudentProfile(repos.Students, repos.StudentProfiles),
		ListClassroomStudents: inclusionuc.NewListClassroomStudents(repos.Students),
		RecommendDevice:       inclusionuc.NewRecommendDevice(repos.AI, repos.Students, repos.Devices, repos.Ramps),
		AssistClassroom:       inclusionuc.NewAssistClassroom(repos.AI, repos.Students, repos.Devices),
	}
}
