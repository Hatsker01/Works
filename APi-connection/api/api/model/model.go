package model

type User struct {
	Id                   string     `protobuf:"bytes,1,opt,name=id,proto3" json:"id"`
	FirstName            string     `protobuf:"bytes,2,opt,name=first_name,json=firstName,proto3" json:"first_name"`
	LastName             string     `protobuf:"bytes,3,opt,name=last_name,json=lastName,proto3" json:"last_name"`
	Email                string     `protobuf:"bytes,4,opt,name=email,proto3" json:"email"`
	Bio                  string     `protobuf:"bytes,5,opt,name=bio,proto3" json:"bio"`
	PhoneNumbers         []string   `protobuf:"bytes,6,rep,name=phone_numbers,json=phoneNumbers,proto3" json:"phone_numbers"`
	Address              []*Address `protobuf:"bytes,7,rep,name=address,proto3" json:"address"`
	Status               string     `protobuf:"bytes,8,opt,name=status,proto3" json:"status"`
	CreatedAt            string     `protobuf:"bytes,9,opt,name=created_at,json=createdAt,proto3" json:"created_at"`
	UpdatedAt            string     `protobuf:"bytes,10,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at"`
	DeletedAt            string     `protobuf:"bytes,11,opt,name=deleted_at,json=deletedAt,proto3" json:"deleted_at"`
	Posts                []*Post    `protobuf:"bytes,12,rep,name=posts,proto3" json:"posts"`
}

type Address struct {
	City                 string   `protobuf:"bytes,1,opt,name=city,proto3" json:"city"`
	Country              string   `protobuf:"bytes,2,opt,name=country,proto3" json:"country"`
	District             string   `protobuf:"bytes,3,opt,name=district,proto3" json:"district"`
	PostalCode           int64    `protobuf:"varint,4,opt,name=postal_code,json=postalCode,proto3" json:"postal_code"`

}

type Post struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name"`
	Description          string   `protobuf:"bytes,3,opt,name=description,proto3" json:"description"`
	UserId               string   `protobuf:"bytes,4,opt,name=user_id,json=userId,proto3" json:"user_id"`
	Medias               []*Media `protobuf:"bytes,5,rep,name=medias,proto3" json:"medias"`
}
type Media struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id"`
	Type                 string   `protobuf:"bytes,2,opt,name=type,proto3" json:"type"`
	Link                 string   `protobuf:"bytes,3,opt,name=link,proto3" json:"link"`
}




