package postgres

import (
	"fmt"
	"time"

	pb "github.com/Hatsker01/Works/APi-connection/user-service/genproto"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type userRepo struct {
	db *sqlx.DB
}

//NewUserRepo ...
func NewUserRepo(db *sqlx.DB) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) CreateUser(user *pb.User) (*pb.User, error) {
	var (
		time_at = time.Now()
		ruser   = pb.User{}
	)
	insertUserQuery := `INSERT INTO users (id, first_name, last_name, email, bio, status, created_at, phone_number) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id, first_name, last_name, email, bio, phone_number, status, created_at`
	err := r.db.QueryRow(insertUserQuery, user.Id, user.FirstName, user.LastName, user.Email, user.Bio, user.Status, time_at, pq.Array(user.PhoneNumbers)).Scan(
		&ruser.Id,
		&ruser.FirstName,
		&ruser.LastName,
		&ruser.Email,
		&ruser.Bio,
		pq.Array(&ruser.PhoneNumbers),
		&ruser.Status,
		&ruser.CreatedAt,
	)
	if err != nil {
		return &pb.User{}, err
	}
	for _, value := range user.Address {
		var adrsId string
		insertAdressQuery := `INSERT INTO adress (user_id, city, country, district, postal_code)
			VALUES ($1, $2, $3, $4, $5) RETURNING user_id`
		err = r.db.QueryRow(insertAdressQuery, user.Id, value.City, value.Country, value.District, value.PostalCode).Scan(&adrsId)
		if err != nil {
			return nil, err
		}
	}
	return &ruser, nil
}

func (r *userRepo) UpdateUser(user *pb.User) (string, error) {
	var (
		timeAt = time.Now()
		// ruser = pb.User{}
	)

	insertUserQuery := `UPDATE users SET first_name = $1, last_name = $2, email = $3, bio = $4, status =$5,
	updated_at = $6, phone_number = $7 where id = $8`
	_, err := r.db.Query(insertUserQuery, user.FirstName, user.LastName, user.Email, user.Bio, user.Status, timeAt, pq.Array(user.PhoneNumbers), user.Id)
	if err != nil {
		return "", err
	}

	for _, value := range user.Address {
		insertAdressQuery := `UPDATE adress SET city = $1, country = $2, district = $3, postal_code = $4 
			WHERE user_id = $5`
		_, err = r.db.Exec(insertAdressQuery, value.City, value.Country, value.District, value.PostalCode, user.Id)
		if err != nil {
			return "", err
		}
	}
	return "", err
}

func (r *userRepo) GetUserById(ID string) (*pb.User, error) {
	var ruser pb.User
	getByIdQuery := `SELECT id, first_name, last_name, email, bio, status, created_at, phone_number FROM users WHERE id = $1`
	err := r.db.QueryRow(getByIdQuery, ID).Scan(
		&ruser.Id,
		&ruser.FirstName,
		&ruser.LastName,
		&ruser.Email,
		&ruser.Bio,
		&ruser.Status,
		&ruser.CreatedAt,
		pq.Array(&ruser.PhoneNumbers),
	)
	if err != nil {
		return &pb.User{}, err
	}
	getByIdAdressQuery := `SELECT city, country, district, postal_code FROM adress WHERE user_id = $1`
	rows, err := r.db.Query(getByIdAdressQuery, ID)
	if err != nil {
		return nil, err
	}
	var tempUser pb.User
	for rows.Next() {
		var adressById pb.Address
		err = rows.Scan(
			&adressById.City,
			&adressById.Country,
			&adressById.District,
			&adressById.PostalCode,
		)
		if err != nil {
			return nil, err
		}
		tempUser.Address = append(tempUser.Address, &adressById)
	}
	ruser.Address = tempUser.Address
	return &ruser, nil
}

func (r *userRepo) GetAllUser() ([]*pb.User, error) {
	var ruser1 []*pb.User

	getByIdQuery := `SELECT id, first_name, last_name, email, bio, status, created_at, phone_number FROM users`
	rowss, err := r.db.Query(getByIdQuery)

	if err != nil {
		return nil, err
	}

	for rowss.Next() {
		var ruser pb.User
		err = rowss.Scan(
			&ruser.Id,
			&ruser.FirstName,
			&ruser.LastName,
			&ruser.Email,
			&ruser.Bio,
			&ruser.Status,
			&ruser.CreatedAt,
			pq.Array(&ruser.PhoneNumbers),
		)
		if err != nil {
			return nil, err
		}

		getByIdAdressQuery := `SELECT city, country, district, postal_code FROM adress where user_id=$1`
		rows, err := r.db.Query(getByIdAdressQuery, ruser.Id)

		if err != nil {
			return nil, err
		}

		//var tempUser pb.User
		for rows.Next() {
			var adressById pb.Address
			err = rows.Scan(
				&adressById.City,
				&adressById.Country,
				&adressById.District,
				&adressById.PostalCode,
			)

			if err != nil {
				return nil, err
			}

			ruser.Address = append(ruser.Address, &adressById)
		}

		ruser1 = append(ruser1, &ruser)
	}

	return ruser1, nil
}

func (r *userRepo) GetUserFromPost(userID string) (*pb.GetUserFromPostResponse, error) {
	var user *pb.GetUserFromPostResponse

	query := `SELECT first_name, last_name FROM users WHERE id = $1`

	err := r.db.QueryRow(query, userID).Scan(
		&user.FirstName,
		&user.LastName,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
func (r *userRepo) DeleteUser(userId string) (pb.Empty, error) {
	//var user *pb.GetUserFromPostResponse
	queryL := `DELETE FROM users where id = $1`
	_, err := r.db.Exec(queryL, userId)
	if err != nil {
		fmt.Println(err)
		return pb.Empty{}, err
	}

	return pb.Empty{}, err
}

func (r *userRepo) GetListUsers(limit, page int64) (*pb.GetUserResponse, error) {

	var (
		users []*pb.User
		count int64
	)
	offset := (page - 1) * limit

	query := `SELECT id,first_name,last_name from users order by first_name OFFSET $1 limit $2`
	rows, err := r.db.Query(query, offset, limit)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user pb.User
		err := rows.Scan(
			&user.Id,
			&user.FirstName,
			&user.LastName,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	countQuery := `SELECT count(*) FROM users`
	err = r.db.QueryRow(countQuery).Scan(&count)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserResponse{
		Users: users,
		Count: count,
	}, nil

}
