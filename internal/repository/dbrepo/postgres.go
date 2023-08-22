package dbrepo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/xuoxod/crew-app/internal/helpers"
	"github.com/xuoxod/crew-app/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

func (m *postgresDBRepo) CreateUser(res models.Member) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	var results = make(map[string]string)

	defer cancel()

	var newID, memberID int
	var fname, lname, email, phone string

	// Create new record in members table
	stmt := `insert into krxbyhhs.public.members(first_name, last_name, email, phone, password, created_at, updated_at)
	values($1,$2,$3,$4,$5,$6,$7) returning id, first_name, last_name, email, phone`

	hashedPassword, hashPasswordErr := helpers.HashPassword(res.Password)

	if hashPasswordErr != nil {
		fmt.Println("Error hashing password: ", hashPasswordErr.Error())
		return nil, hashPasswordErr
	}

	row := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		hashedPassword,
		time.Now(),
		time.Now(),
	)

	memberErr := row.Scan(&newID, &fname, &lname, &email, &phone)

	if memberErr != nil {
		fmt.Println("memberErr: ", memberErr.Error())
		return nil, memberErr
	}

	// Create new record in profiles table
	stmt = `
	insert into krxbyhhs.public.profiles(member_id, user_name, image_url, craft, display_name, updated_at, years_of_service)
	values($1, $2, $3, $4, $5, $6, $7) returning member_id`

	row = m.DB.QueryRowContext(ctx, stmt,
		newID,
		email,
		"https://via.placeholder.com/150/659403",
		"crew",
		email,
		time.Now(),
		0,
	)

	profileErr := row.Scan(&memberID)

	if profileErr != nil {
		fmt.Println("profileErr: ", profileErr.Error())
		return nil, profileErr
	}

	// Create new record in user_settings table
	stmt = `insert into krxbyhhs.public.user_settings(member_id) values($1) returning member_id`

	row = m.DB.QueryRowContext(ctx, stmt, newID)

	userSettingsErr := row.Scan(&memberID)

	if userSettingsErr != nil {
		fmt.Println("userSettingsErr: ", userSettingsErr.Error())
		return nil, userSettingsErr
	}

	results["ID"] = fmt.Sprintf("%d", newID)
	results["firstName"] = fname
	results["lastName"] = lname
	results["email"] = email
	results["phone"] = phone

	return results, nil
}

// User stuff

func (m *postgresDBRepo) GetUserByID(id int) (models.Member, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `select id, first_name, last_name, email, phone, password, access_level, created_at, updated_at, craft_id from users where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var u models.Member
	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Phone,
		&u.Password,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return u, err
	}

	return u, nil
}

func (m *postgresDBRepo) GetUserByEmail(email string) (models.Member, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `select id, first_name, last_name, email, phone, password, access_level, created_at, updated_at, craft_id from users where email = $1`

	row := m.DB.QueryRowContext(ctx, query, email)

	var u models.Member
	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return u, err
	}

	return u, nil
}

func (m *postgresDBRepo) UpdateUserProfile(u models.Member, p models.Profile) map[string]string {
	var results = make(map[string]string)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var fname, lname, email, phone, uname, dname, iurl, craft, address, city, state, profileUpdatedAt, createdAt string
	var memberId, yearsService int

	// Update member table
	memberQuery := `
		update members m set first_name = $1, last_name = $2, email = $3, phone = $4, updated_at = $5 where m.email = $6 returning m.id, m.first_name, m.last_name, m.email, m.phone, created_at
	`

	memberRows, memberErr := m.DB.QueryContext(ctx, memberQuery,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Phone,
		time.Now(),
		u.Email,
	)

	if memberErr != nil {
		fmt.Printf("\tMember Query Error: %s\n", memberErr.Error())
		results["memberQueryErr"] = memberErr.Error()
		return results
	}

	for memberRows.Next() {
		if err := memberRows.Scan(&memberId, &fname, &lname, &email, &phone, &createdAt); err != nil {
			fmt.Printf("\tMember Row Scan Error: %s\n", err.Error())
			results["memberRowsScanErr"] = err.Error()
			return results
		}
	}

	memberRerr := memberRows.Close()

	if memberRerr != nil {
		fmt.Printf("member rerr error:\t%s\n", memberRerr.Error())
		results["memberRerr"] = memberRerr.Error()
		return results
	}

	if err := memberRows.Err(); err != nil {
		fmt.Printf("Member Rows Error:\t%s\n", err.Error())
		results["memberRowsErr"] = err.Error()
		return results
	}

	// Update profile table

	proifleQuery := `
	update profiles set user_name = $1, image_url = $2, craft = $3, address = $4, city = $5, state = $6, display_name = $7, years_of_service = $8, updated_at = $9 where member_id = $10 returning user_name, image_url, craft, address, city, state, display_name, years_of_service, updated_at`

	profileRows, profileErr := m.DB.QueryContext(ctx, proifleQuery,
		p.UserName,
		p.ImageURL,
		p.Craft,
		p.Address,
		p.City,
		p.State,
		p.DisplayName,
		p.YOS,
		time.Now(),
		memberId,
	)

	if profileErr != nil {
		fmt.Printf("\tProfile Query Error: %s\n", profileErr.Error())
		results["profileQueryErr"] = profileErr.Error()
		return results
	}

	for profileRows.Next() {
		if err := profileRows.Scan(&uname, &iurl, &craft, &address, &city, &state, &dname, &yearsService, &profileUpdatedAt); err != nil {
			fmt.Printf("\tProfile Row Scan Error: %s\n", err.Error())
			results["profileRowsScanError"] = err.Error()
			return results
		}
	}

	profileRerr := profileRows.Close()

	if profileRerr != nil {
		fmt.Printf("profile rerr error:\t%s\n", profileRerr.Error())
		results["profileRerr"] = profileRerr.Error()
		return results
	}

	if err := profileRows.Err(); err != nil {
		fmt.Printf("Profile Rows Error:\t%s\n", err.Error())
		results["profileRowsErr"] = err.Error()
		return results
	}

	results["userID"] = fmt.Sprintf("%d", memberId)
	results["firstName"] = fname
	results["lastName"] = lname
	results["email"] = email
	results["phone"] = phone
	results["userName"] = uname
	results["displayName"] = dname
	results["imageUrl"] = iurl
	results["craft"] = craft
	results["address"] = address
	results["city"] = city
	results["state"] = state
	results["createdAt"] = fmt.Sprintf("%v", createdAt)
	results["profileUpdatedAt"] = fmt.Sprintf("%v", profileUpdatedAt)
	results["yearsService"] = fmt.Sprintf("%d", yearsService)

	return results
}

func (m *postgresDBRepo) AuthenticateUser(email, testPassword string) map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var results = make(map[string]string)
	var accessLevel, id, yos, memberId int
	var firstName, lastName, emailAddress, phone, hashedPassword, userName, imgUrl, craft, address, city, state, displayName string
	var createdAt, updatedAt time.Time
	var showProfile, showOnlineStatus, showAddress, showCity, showState, showDisplayName, showContactInfo, showPhone, showEmail, showCraft, showRun, showNotifications bool

	query := `select m.first_name, m.last_name, m.email, m.phone, m.access_level, m.created_at, m.updated_at, m.password, 
	p.user_name, p.image_url, p.craft, p.years_of_service, p.address, p.city, p.state, p.display_name, p.member_id, 
	us.show_profile, us.show_online_status, us.show_address, us.show_city, us.show_state, us.show_display_name, us.show_contact_info, us.show_phone, us.show_email, us.show_craft, us.show_run, us.show_notifications, us.member_id from members m 
	inner join profiles p on p.member_id = m.id inner join user_settings us on us.member_id = m.id where email = $1`

	row := m.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(&firstName, &lastName, &emailAddress, &phone, &accessLevel, &createdAt, &updatedAt, &hashedPassword, &userName, &imgUrl, &craft, &yos, &address, &city, &state, &displayName, &memberId, &showProfile, &showOnlineStatus, &showAddress, &showCity, &showState, &showDisplayName, &showContactInfo, &showPhone, &showEmail, &showCraft, &showRun, &showNotifications, &id)

	if err != nil {
		log.Printf("\n\tScan error:\n\t%s\n", err.Error())
		results["scanerr"] = err.Error()
		// return results
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))

	if err == bcrypt.ErrMismatchedHashAndPassword {
		log.Println("bcrypt error:\t", err.Error())
		results["err"] = err.Error()
		return results
	} else if err != nil {
		log.Println("bcrypt error:\t", err.Error())
		results["err"] = err.Error()
		return results
	}

	if userName == "" {
		userName = "Create username"
	}

	if imgUrl == "" {
		imgUrl = "Upload photo"
	}

	if address == "" {
		address = "Enter street address"
	}

	if city == "" {
		city = "Enter home town"
	}

	if state == "" {
		state = "Enter state"
	}

	if displayName == "" {
		displayName = "Create display name"
	}

	results["userID"] = fmt.Sprintf("%d", memberId)
	results["firstName"] = firstName
	results["lastName"] = lastName
	results["email"] = email
	results["phone"] = phone
	results["userName"] = userName
	results["imgUrl"] = imgUrl
	results["craft"] = craft
	results["address"] = address
	results["city"] = city
	results["state"] = state
	results["displayName"] = displayName
	results["accessLevel"] = fmt.Sprintf("%d", accessLevel)
	results["createdAt"] = createdAt.String()
	results["updatedAt"] = updatedAt.String()
	results["yos"] = fmt.Sprintf("%d", yos)
	results["showProfile"] = fmt.Sprintf("%t", showProfile)
	results["showOnlineStatus"] = fmt.Sprintf("%t", showOnlineStatus)
	results["showAddress"] = fmt.Sprintf("%t", showAddress)
	results["showCity"] = fmt.Sprintf("%t", showCity)
	results["showState"] = fmt.Sprintf("%t", showState)
	results["showDisplayName"] = fmt.Sprintf("%t", showDisplayName)
	results["showContactInfo"] = fmt.Sprintf("%t", showContactInfo)
	results["showPhone"] = fmt.Sprintf("%t", showPhone)
	results["showEmail"] = fmt.Sprintf("%t", showEmail)
	results["showCraft"] = fmt.Sprintf("%t", showCraft)
	results["showRun"] = fmt.Sprintf("%t", showRun)
	results["showNotifications"] = fmt.Sprintf("%t", showNotifications)
	results["err"] = ""
	return results
}

func (m *postgresDBRepo) UpdateUserSettings(p models.UserSettings) map[string]string {

	results := make(map[string]string)

	return results
}
