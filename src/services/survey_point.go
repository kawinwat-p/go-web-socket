package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"websocketjingjing/domain/entities"
	"websocketjingjing/domain/repositories"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type sheetApiClient struct {
	sheetService *sheets.Service
}

type surveyPointService struct {
	client           sheetApiClient
	userRepo         repositories.IUsersRepository
	alertMessageRepo repositories.IAlertMessageRepositories
	redisrepo        repositories.IRedisConnectionRepository
}

type ISurveyPointService interface {
	AddCredits() entities.ResponseModel
	SetAlertMessageRedis() error
}

func NewSurveyPointService(userRepo repositories.IUsersRepository, alertMessageRepo repositories.IAlertMessageRepositories, redisRepo repositories.IRedisConnectionRepository) ISurveyPointService {
	return &surveyPointService{
		client:           *NewSheetApiClient("credentials.json"),
		userRepo:         userRepo,
		alertMessageRepo: alertMessageRepo,
		redisrepo:        redisRepo,
	}
}

func NewSheetApiClient(filePath string) *sheetApiClient {
	ctx := context.Background()

	serviceAccountKey := "credentials.json"

	sheetServices, err := sheets.NewService(ctx, option.WithCredentialsFile(serviceAccountKey))

	if err != nil {
		return nil
	}

	return &sheetApiClient{sheetService: sheetServices}
}

func (client sheetApiClient) GetSheet(sheetId string, sheetName string) ([][]interface{}, error) {
	data, err := client.sheetService.Spreadsheets.Values.Get(sheetId, sheetName).Do()

	if err != nil {
		return nil, err
	}

	return data.Values, nil
}

func (sv surveyPointService) AddCredits() entities.ResponseModel {
	sheetId := os.Getenv("SHEET_ID")
	sheetName := os.Getenv("SHEET_NAME")

	data, err := sv.client.GetSheet(sheetId, sheetName)

	if err != nil {
		return entities.ResponseModel{Message: err.Error(), Status: 401}
	}

	users := []string{}

	if len(data) <= 0 {
		return entities.ResponseModel{Message: "no data found in sheet", Status: 401}
	} else {
		var credit int32
		redisData := sv.redisrepo.GetRedisAlertMessageData()
		if redisData == nil {
			credit, err = sv.alertMessageRepo.GetSurveyCredits()
			if err != nil {
				credit = 10000
			}
		} else {
			credit = redisData.Teacher.Point
		}
		for i, row := range data {
			if i == 0 {
				continue
			}

			uid := row[5].(string)

			var alreadyPaid string

			if len(row) >= 14 {
				alreadyPaid = row[13].(string)
			} else {
				alreadyPaid = ""
			}

			if alreadyPaid == "จ่ายแล้ว" {
				fmt.Println(uid + alreadyPaid)
				continue
			} else if sv.userRepo.UserExist(uid) {
				sv.userRepo.AddCredits(uid, credit)
				fmt.Println("Adding", credit, "point to", uid)
				users = append(users, uid)

				if len(row) >= 14 {
					row[13] = "จ่ายแล้ว"
				} else {
					row = append(row, "จ่ายแล้ว")
				}

				data[i] = row
			} else {
				fmt.Println("User with", uid, "uid not found")
				if len(row) >= 14 {
					row[13] = "ไม่พบผู้ใช้งาน"
				} else {
					row = append(row, "ไม่พบผู้ใช้งาน")
				}

				data[i] = row
			}
		}
	}
	_, err = sv.client.sheetService.Spreadsheets.Values.Update(sheetId, sheetName, &sheets.ValueRange{Values: data}).ValueInputOption("RAW").Do()

	if err != nil {
		return entities.ResponseModel{
			Message: err.Error(),
			Status:  401,
		}
	}

	if len(users) > 0 {
		return entities.ResponseModel{
			Message: "Successfully add point to the user",
			Status:  200,
		}
	} else {
		return entities.ResponseModel{
			Message: "No user need to be paid",
			Status:  200,
		}
	}

}

func (sv surveyPointService) SetAlertMessageRedis() error {
	data, err := sv.alertMessageRepo.GetAll()

	if err != nil {
		return err
	}

	dataJson, err := json.Marshal(data)

	if err != nil {
		return err
	}

	cache := sv.redisrepo.SetRedisData(dataJson)

	if !cache {
		return fmt.Errorf("cannot set redis data")
	}

	return nil
}
