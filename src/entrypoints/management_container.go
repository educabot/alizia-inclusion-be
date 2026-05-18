package entrypoints

import mgmtuc "github.com/educabot/alizia-inclusion-be/src/core/usecases/management"

type ManagementContainer struct {
	ListClassrooms  mgmtuc.ListClassrooms
	GetClassroom    mgmtuc.GetClassroom
	CreateClassroom mgmtuc.CreateClassroom
	UpdateClassroom mgmtuc.UpdateClassroom
	DeleteClassroom mgmtuc.DeleteClassroom
	ListTeachers    mgmtuc.ListTeachers
}
