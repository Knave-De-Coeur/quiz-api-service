package services

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"quiz-api-service/internal/pkg"

	"quiz-api-service/internal/api"
)

type UserService struct {
	DBConn   *gorm.DB
	logger   *zap.Logger
	settings UserServiceSettings
}

// UserServiceSettings used to affect code flow
type UserServiceSettings struct {
	Port     int
	Hostname string
}

type UserServices interface {
	InsertUser(user *api.User) error
	GetUsers() ([]api.User, error)
	GetUserByUsername(username string) (*api.User, error)
	GetUserByID(uID uint) (*api.User, error)
	Login(request api.LoginRequest) (*api.User, error)
}

func NewUserService(dbConn *gorm.DB, logger *zap.Logger, settings UserServiceSettings) *UserService {
	return &UserService{
		DBConn:   dbConn,
		logger:   logger,
		settings: settings,
	}
}

// InsertUser inserts new user in users table from data passed in arg.
func (service *UserService) InsertUser(req *api.User) error {

	user := &pkg.User{
		Name:     req.Name,
		Username: req.Username,
		Age:      req.Age,
		Password: req.Password,
	}

	res := service.DBConn.Select("name", "age", "username", "password").Create(user)
	if res.Error != nil {
		service.logger.Error("something went wrong inserting user", zap.Any("user", user), zap.Error(res.Error))
		return res.Error
	}

	service.logger.Debug("rows inserted", zap.Int64("rowsAffected", res.RowsAffected))

	return nil
}

// GetUsers returns list of users in db.
func (service *UserService) GetUsers() ([]api.User, error) {

	var users []api.User

	// Get all records
	res := service.DBConn.Select("name", "age", "username", "created_at", "updated_at", "id").Find(&users)
	if res.Error != nil {
		service.logger.Error("something went wrong getting all players", zap.Error(res.Error))
		return nil, res.Error
	}

	service.logger.Debug("users grabbed", zap.Int64("number", res.RowsAffected))

	return users, nil
}

// GetUserByUsername attempts to retrieve a single row from the users table.
func (service *UserService) GetUserByUsername(username string) (*api.User, error) {

	var user api.User
	// Get all records
	res := service.DBConn.Select("name", "age", "username", "password", "last_login_time_stamp").Where("username = ?", username).First(&user)
	if res.Error != nil {
		service.logger.Error("something went wrong getting player by username", zap.Error(res.Error), zap.String("username", username))
		return nil, res.Error
	}

	service.logger.Debug("user grabbed", zap.Any("user", user))

	return &user, nil
}

// GetUserByID grabs from table by id
func (service *UserService) GetUserByID(uID uint) (*api.User, error) {

	var user api.User
	// Get all records
	res := service.DBConn.Select("name", "age", "username", "password", "last_login_time_stamp").Where("id = ?", uID).First(&user)
	if res.Error != nil {
		service.logger.Error("something went wrong getting player by ID", zap.Error(res.Error))
		return nil, res.Error
	}

	service.logger.Debug("user grabbed", zap.Any("user", user))

	return &user, nil
}

// Login is a wrapper for the GetUserByUsername that also validates the password
func (service *UserService) Login(request api.LoginRequest) (*api.User, error) {

	user, err := service.GetUserByUsername(request.Username)
	if err != nil {
		return nil, err
	}

	if user.Password != request.Password {
		return nil, fmt.Errorf("invalid passord for user")
	}

	return user, nil
}

// Compare stats endpoint func that returns the message with how the user did compare to others
// func compareUserScores(res http.ResponseWriter, req *http.Request) {
// 	reqBody, _ := ioutil.ReadAll(req.Body)
// 	var compareUsersRequest api.CompareUsersRequest
// 	_ = json.Unmarshal(reqBody, &compareUsersRequest)
//
// 	user := searchUsersByID(compareUsersRequest.UserID)
//
// 	message := ""
// 	errorFound := false
// 	if len(user.SubmittedAnswers) < 1 {
// 		message = "Start playing to compare results!"
// 		errorFound = true
// 	} else {
// 		x := getUserComparisonScore(user)
//
// 		negative := math.Signbit(x)
//
// 		userScoreComparison := strconv.FormatFloat(x, 'f', 0, 64)
//
// 		if negative {
// 			message = "You did " + userScoreComparison + "% worse than everyone!"
// 		} else {
// 			message = "You did " + userScoreComparison + "% better than everyone!"
// 		}
// 	}
//
// 	responseMessage := api.Response{
// 		Message: message,
// 		Error:   errorFound,
// 	}
//
// 	_ = json.NewEncoder(res).Encode(responseMessage)
// }

// This calculates the comparison percentage the user has from other users
// func getUserComparisonScore(currentUser pkg.User) float64 {
//
// 	var listOfScores []int
//
// 	var sumPercentages int
//
// 	for i := range ListOfUsers {
// 		if ListOfUsers[i].ID != currentUser.ID {
// 			scorePercentage := ListOfUsers[i].Score * 20
// 			sumPercentages += scorePercentage
// 			listOfScores = append(listOfScores, scorePercentage)
// 		}
// 	}
//
// 	averagePercentage := float64(sumPercentages) / (float64(len(listOfScores)))
//
// 	scorePercentage := float64(currentUser.Score * 20)
//
// 	x := scorePercentage - averagePercentage
//
// 	return x
// }