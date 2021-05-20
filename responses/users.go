package responses

import "go.mongodb.org/mongo-driver/bson/primitive"

type UsersResponse struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Nickname    string             `bson:"nickname" json:"nickname"`       //老师昵称
	Photo       string             `bson:"photo" json:"photo"`             //老师头像
	Phonenumber string             `bson:"phonenumber" json:"phonenumber"` //老师手机号码
	TalkmateId  string             `bson:"talkmate_id" json:"talkmate_id"` //talkmate_id
}
