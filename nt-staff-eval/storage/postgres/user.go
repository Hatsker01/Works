package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"

	"github.com/Hatsker01/nt-staff-eval/pkg/structs"
	"github.com/Hatsker01/nt-staff-eval/pkg/utils"
	"github.com/Hatsker01/nt-staff-eval/storage/repo"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) repo.UserRepoInterface {
	return &userRepository{
		db: db,
	}
}

func (r userRepository) CreateUser(user structs.CreateUser) (structs.UserStruct, error) {
	branchId := utils.IntToNullInt64(user.BranchId)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return structs.UserStruct{}, err
	}

	err = r.db.QueryRow(`INSERT INTO users_client(
	 id, spec_id, first_name, last_name, email, password, branch_id, gender, role_id)
	 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`,
		user.Id,
		user.SpecId,
		user.FirstName,
		user.LastName,
		user.Email,
		string(hashedPassword),
		branchId,
		user.Gender,
		user.RoleId,
	).Scan(&user.Id)
	if err != nil {
		return structs.UserStruct{}, err
	}
	userClient, err := r.GetUser(user.Id)
	if err != nil {
		return structs.UserStruct{}, err
	}
	return userClient, nil
}

func (r userRepository) GetUser(id string) (structs.UserStruct, error) {
	var (
		user             structs.UserStruct
		scoreInfo        structs.UserScoreInfo
		addedAt          sql.NullTime
		branchId         sql.NullInt64
		cover            sql.NullString
		workTime         sql.NullString
		phone            sql.NullString
		address          sql.NullString
		email            sql.NullString
		socialMedias     sql.NullString
		additionalInform sql.NullString
		birthday         sql.NullString
	)

	err := r.db.QueryRow(`select  id, spec_id, first_name, last_name, email, password, cover, branch_id, birthday, gender, added_at, 
       role_id, phone, address, work_time, social_medias, additional_informs, created_at, updated_at from users_client where deleted_at is null and id=$1`, id).
		Scan(
			&user.Id,
			&user.SpecId,
			&user.FirstName,
			&user.LastName,
			&email,
			&user.Password,
			&cover,
			&branchId,
			&birthday,
			&user.Gender,
			&addedAt,
			&user.Role.Id,
			&phone,
			&address,
			&workTime,
			&socialMedias,
			&additionalInform,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
	if err != nil {
		return structs.UserStruct{}, err
	}
	if addedAt.Valid {
		user.AddedAt = addedAt.Time
	}
	if branchId.Valid {
		user.Branch.Id = branchId.Int64
	}
	if cover.Valid {
		user.Cover = cover.String
	}
	if workTime.Valid {
		user.WorkTime = workTime.String
	}
	if phone.Valid {
		user.Phone = phone.String
	}
	if address.Valid {
		user.Address = address.String
	}
	if email.Valid {
		user.Email = email.String
	}
	if birthday.Valid {
		user.Birthday = birthday.String
	}
	if socialMedias.Valid {
		err = json.Unmarshal([]byte(socialMedias.String), &user.SocialMedias)
		if err != nil {
			return structs.UserStruct{}, err
		}
	}
	if additionalInform.Valid {
		err = json.Unmarshal([]byte(additionalInform.String), &user.AdditionalInforms)
		if err != nil {
			return structs.UserStruct{}, err
		}
	}

	scoreInfo.StaffAverage, err = r.StaffAverage(user.Id)
	if err != nil {
		return structs.UserStruct{}, err
	}

	scoreInfo.ClientAverage, err = r.ClientAverage(user.Id)
	if err != nil {
		return structs.UserStruct{}, err
	}
	user.ScoreInfo = scoreInfo

	return user, nil
}

func (r userRepository) GetListUsers(filters map[string]string, page, limit int) ([]structs.UserListResp, int, error) {
	offset := (page - 1) * limit
	fmt.Println(filters)

	sb := sqlbuilder.NewSelectBuilder()

	sb.Select("u.id", "u.spec_id", "u.first_name", "u.last_name", "u.cover",
		"u.role_id", "social_medias")
	sb.From("users_client u")

	if value, ok := filters["section"]; ok {
		sb.Join("roles r", "r.id = u.role_id")
		sb.Where(sb.Equal("r.section_id", value))
		sb.Where(sb.IsNull("u.deleted_at"))
	}

	if value, ok := filters["searchId"]; ok {
		//value = strings.Trim(value, " ")
		value = strings.Replace(value, " ", "", len(value))
		sb.Where("(u.spec_id::varchar(64) ilike " + fmt.Sprintf("'%%%s%%'", value) +
			" or " + "u.first_name::varchar(64) || u.last_name::varchar(64) ilike " + fmt.Sprintf("'%%%s%%') ", value))
	}

	if value, ok := filters["branchId"]; ok {
		sb.Join("branches b", "b.id = u.branch_id")
		sb.Where(sb.Equal("u.branch_id", value))
	}
	sb.Offset(offset)
	sb.Limit(limit)
	query, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)
	fmt.Println(query)
	rows, err := r.db.Queryx(query, args...)
	if err != nil {
		return nil, 0, err
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var (
		users []structs.UserListResp
		count int
	)

	for rows.Next() {
		var (
			user         structs.UserListResp
			covNull      sql.NullString
			socialMedias sql.NullString
		)

		err = rows.Scan(
			&user.Id,
			&user.SpecId,
			&user.FirstName,
			&user.LastName,
			&covNull,
			&user.Role.Id,
			&socialMedias,
		)

		if err != nil {
			return nil, 0, err
		}
		if covNull.Valid {
			user.Cover = covNull.String
		}
		if socialMedias.Valid {
			err = json.Unmarshal([]byte(socialMedias.String), &user.SocialMedias)
			if err != nil {
				return nil, 0, err
			}
		}
		user.ScoreInfo.StaffAverage, err = r.StaffAverage(user.Id)
		if err != nil {
			return nil, 0, err
		}
		user.ScoreInfo.ClientAverage, err = r.ClientAverage(user.Id)
		if err != nil {
			return nil, 0, err
		}
		var score int
		err = r.db.QueryRow(`select count(*) from rated 
                where user_id=$1 and is_staff=false`, user.Id,
		).Scan(&score)
		user.ScoreInfo.NumberOfClientRated = score
		if err != nil {
			return nil, 0, err
		}
		err = r.db.QueryRow(`select count(*) from rated 
                where user_id=$1 and is_staff=true`, user.Id,
		).Scan(&score)
		user.ScoreInfo.NumberOfStaffRated = score
		if err != nil {
			return nil, 0, err
		}


		if filters["placeTop"] != "" {
			top, err := strconv.Atoi(filters["placeTop"])
			fmt.Println(top)
			if err != nil {
				return nil, 0, err
			}

			var place sql.NullInt64
			q := `with ca as (
					with stars_client as
						(select e.star::float as stars_client, r.user_id as user_id_c
						from rated as r
							join rated_evaluations re on r.id = re.rated_id
							join evaluations e on re.evaluation_id = e.id
						where r.is_staff = false
						group by re.rated_id, e.star, r.user_id),
					
						 stars_staff as
							(select e.star::float as stars_staff, r.user_id as user_id_s
						from rated as r
							join rated_evaluations re on r.id = re.rated_id
							join evaluations e on re.evaluation_id = e.id
						where r.is_staff = true
						group by re.rated_id, e.star, r.user_id)
					
						select (sum(stars_client) / count(stars_client))::decimal(10, 2) as average_client, u.id
						from stars_client
							right join users_client u on u.id = stars_client.user_id_c
							left join stars_staff on stars_staff.user_id_s = u.id
						where u.deleted_at is null and stars_client is not null
						group by user_id_c, u.id order by average_client desc limit $1) select count(*) from ca where ca.id = $2`
			err = r.db.QueryRow(q, top, user.Id).Scan(&place)
			fmt.Println(place)
			if err != nil {
				return nil, 0, err
			}
			fmt.Println(q)
			fmt.Println(top,user.Id)
			if place.Valid && place.Int64 > 0 {
				user.UserPlaceInTop = top
			}
		}
		

		users = append(users, user)
	}
	

	countAll := sqlbuilder.NewSelectBuilder()

	countAll.Select("count(*)")
	countAll.From("users_client u")
	fmt.Println("\n\n\n",users)
	if value, ok := filters["section"]; ok {
		countAll.Join("roles r", "r.id = u.role_id")
		countAll.Where(countAll.Equal("r.section_id", value))
	}
	fmt.Println(users[0])
	

	if value, ok := filters["searchId"]; ok {
		//value = strings.Trim(value, " ")
		value = strings.Replace(value, " ", "", len(value))
		countAll.Where("(u.spec_id::varchar(64) ilike " + fmt.Sprintf("'%%%s%%'", value) +
			" or " + "u.first_name::varchar(64) || u.last_name::varchar(64) ilike " + fmt.Sprintf("'%%%s%%')", value))
	}

	fmt.Println("\n\n\n",users[0])

	if value, ok := filters["branchId"]; ok {
		countAll.Join("branches b", "b.id = u.branch_id")
		countAll.Where(countAll.Equal("u.branch_id", value))
	}
	fmt.Println("\n\n\n",users[0])
	query, args = countAll.BuildWithFlavor(sqlbuilder.PostgreSQL)
	rows, err = r.db.Queryx(query, args...)
	if err != nil {
		return nil, 0, err
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	fmt.Println("\n\n\nLalaldsflkfasdlkjfhaskjd\n\n")
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return nil, 0, err
		}
	}
	fmt.Println("\n\nAssalomu alaykum\n\n")

	return users, count, nil
}

func (r userRepository) GetTopUsers(page, limit int) ([]structs.UserStruct, int, error) {
	offset := (page - 1) * limit

	query := `with stars_client as
    (select e.star::float as stars_client, r.user_id as user_id_c
	from rated as r
		join rated_evaluations re on r.id = re.rated_id
		join evaluations e on re.evaluation_id = e.id
	where r.is_staff = false
	group by re.rated_id, e.star, r.user_id),
	
	 stars_staff as
		(select e.star::float as stars_staff, r.user_id as user_id_s
	from rated as r
		join rated_evaluations re on r.id = re.rated_id
		join evaluations e on re.evaluation_id = e.id
	where r.is_staff = true
	group by re.rated_id, e.star, r.user_id)

	select (sum(stars_client) / count(stars_client))::decimal(10, 2) as average_client, (sum(stars_staff) / count(stars_staff))::decimal(10, 2) as average_staff,
		   u.id, u.spec_id, u.first_name, u.last_name, u.email, u.password, u.cover, u.branch_id, u.birthday, u.gender, u.added_at,
		   u.role_id, u.phone, u.address, u.work_time, u.social_medias, u.additional_informs, u.created_at, u.updated_at
	from stars_client
		right join users_client u on u.id = stars_client.user_id_c
		left join stars_staff on stars_staff.user_id_s = u.id
	where u.deleted_at is null and stars_client is not null
	group by user_id_c, u.id 
	order by average_client desc`

	query += ` limit $1 offset $2`
	// fmt.Println(query, limit, offset)
	rows, err := r.db.Queryx(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	fmt.Println("1")
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	fmt.Println("2")
	defer rows.Close()

	var (
		users []structs.UserStruct
		count int
	)

	for rows.Next() {
		var (
			user             structs.UserStruct
			averageStaff     sql.NullFloat64
			averageClient    sql.NullFloat64
			covNull          sql.NullString
			eNull            sql.NullString
			workTime         sql.NullString
			branchId         sql.NullInt64
			phone            sql.NullString
			address          sql.NullString
			socialMedias     sql.NullString
			additionalInform sql.NullString
		)

		err = rows.Scan(
			&averageClient,
			&averageStaff,
			&user.Id,
			&user.SpecId,
			&user.FirstName,
			&user.LastName,
			&eNull,
			&user.Password,
			&covNull,
			&branchId,
			&user.Birthday,
			&user.Gender,
			&user.AddedAt,
			&user.Role.Id,
			&phone,
			&address,
			&workTime,
			&socialMedias,
			&additionalInform,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			return nil, 0, err
		}
		if averageClient.Valid {
			user.ScoreInfo.ClientAverage = averageClient.Float64
		}
		if averageStaff.Valid {
			user.ScoreInfo.StaffAverage = averageStaff.Float64
		}
		if covNull.Valid {
			user.Cover = covNull.String
		}
		if eNull.Valid {
			user.Email = eNull.String
		}
		if workTime.Valid {
			user.WorkTime = workTime.String
		}
		if branchId.Valid {
			user.Branch.Id = branchId.Int64
		}
		if phone.Valid {
			user.Phone = phone.String
		}
		if address.Valid {
			user.Address = address.String
		}
		fmt.Println("3")
		if socialMedias.Valid {
			err = json.Unmarshal([]byte(socialMedias.String), &user.SocialMedias)
			if err != nil {
				return nil, 0, err
			}
		}
		if additionalInform.Valid {
			err = json.Unmarshal([]byte(additionalInform.String), &user.AdditionalInforms)
			if err != nil {
				return nil, 0, err
			}
		}

		users = append(users, user)
	}
	fmt.Println(users)


	query = `select count(*) from users_client where deleted_at is null`
	fmt.Println(query)
	rows, err = r.db.Queryx(query)
	if err != nil {
		return nil, 0, err
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return nil, 0, err
		}
	}
	fmt.Println(users)
	return users, count, nil
}

func (r userRepository) UpdateUser(user structs.UpdateUserFromUser) (structs.UserStruct, error) {
	result, err := r.db.Exec("update users_client set  gender=$1, phone=$2, address=$3, birthday=$4, social_medias=$5, first_name=$6, last_name=$7, updated_at=$8 where id=$9 and deleted_at is null",
		user.Gender, user.Phone, user.Address, user.Birthday, user.SocialMedias, user.FirstName, user.LastName, time.Now().UTC(), user.ID,
	)
	if err != nil {
		return structs.UserStruct{}, err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return structs.UserStruct{}, sql.ErrNoRows
	}

	var newUser structs.UserStruct

	newUser, err = r.GetUser(user.ID)

	if err != nil {
		return structs.UserStruct{}, err
	}

	return newUser, nil
}

func (r userRepository) StaffAverage(id string) (float64, error) {
	var StaffAverage sql.NullFloat64

	err := r.db.QueryRow(`
	with stars as (select e.star::float as stars
		from rated as r
			join rated_evaluations re on r.id = re.rated_id
			join evaluations e on re.evaluation_id = e.id
		where r.is_staff = true and r.user_id = $1
		group by re.rated_id, e.star)
	
	select (sum(stars) / count(stars))::decimal(10, 2) as average from stars
	`, id).Scan(&StaffAverage)
	if err != nil {
		return 0, err
	}

	if !StaffAverage.Valid {
		StaffAverage.Float64 = 0
	}

	return StaffAverage.Float64, nil
}

func (r userRepository) ClientAverage(id string) (float64, error) {
	var ClientAverage sql.NullFloat64

	err := r.db.QueryRow(`
	with stars as (select e.star::float as stars
		from rated as r
			join rated_evaluations re on r.id = re.rated_id
			join evaluations e on re.evaluation_id = e.id
		where r.is_staff = false and r.user_id = $1
		group by re.rated_id, e.star)
	
	select (sum(stars) / count(stars))::decimal(10, 2) as average from stars
	`, id).Scan(&ClientAverage)
	if err != nil {
		return 0, err
	}

	if !ClientAverage.Valid {
		ClientAverage.Float64 = 0
	}

	return ClientAverage.Float64, nil
}

func (r userRepository) LoginUser(user structs.LoginUser) (structs.UserLoginResp, error) {
	var resp structs.UserLoginResp
	err := r.db.QueryRow(`select  id, spec_id, first_name, last_name, email, password 
       from users_client where deleted_at is null and email=$1`, user.Email).
		Scan(
			&resp.Id,
			&resp.SpecId,
			&resp.FirstName,
			&resp.LastName,
			&resp.Email,
			&resp.Password,
		)
	if err != nil {
		return structs.UserLoginResp{}, errors.New("wrong username or password")
	}

	return resp, nil
}

func (r userRepository) UpdateToken(user structs.UserAuth) error {
	err := r.db.QueryRow(`select user_id from users_auth
	where user_id=$1`, user.Id).Scan(&user.Id)
	if err == sql.ErrNoRows {
		id, err := uuid.NewV4()
		if err != nil {
			return err
		}
		r.db.QueryRow(`INSERT INTO users_auth(id, access_token, refresh_token, user_id) 
			VALUES($1, $2, $3, $4)`, id.String(), user.AccessToken, user.RefreshToken, user.Id,
		)

		return nil
	}

	result, err := r.db.Exec(`update users_auth set access_token=$1, refresh_token=$2 where user_id=$3`,
		user.AccessToken,
		user.RefreshToken,
		user.Id)
	if err != nil {
		return err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return sql.ErrNoRows
	}

	return nil
}
func (a userRepository) HandleToken(token string) (user structs.UserStruct, err error) {
	err = a.db.QueryRow(`select user_id from users_auth where access_token = $1`, token).Scan(&user.Id)
	if err != nil {
		return structs.UserStruct{}, err
	}
	user, err = a.GetUser(user.Id)
	if err != nil {
		return structs.UserStruct{}, err
	}
	return
}

func (a userRepository) ChangePassword(user structs.UserStruct, password structs.UpdatePassword) (changed bool, err error) {
	if user.Password == password.OldPassword {
		_, err := a.db.Exec("update users_client set password = $1 where id = $2", password.NewPassword, user.Id)
		if err != nil {
			return false, err
		}
		return true, nil
	} else {
		return false, err
	}
	return
}
