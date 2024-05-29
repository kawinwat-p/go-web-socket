package entities

type TeacherData struct{
	Point int32 `json:"give_point" bson:"give_point"`
}

type AlertMessage struct{
	Teacher TeacherData `json:"teacher" bson:"teacher"`
}