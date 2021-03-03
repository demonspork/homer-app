package service

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/Jeffail/gabs/v2"
	uuid "github.com/satori/go.uuid"
	"github.com/sipcapture/homer-app/model"
	"github.com/sipcapture/homer-app/utils/heputils"
)

type DashBoardService struct {
	ServiceConfig
}

// this method gets all users from database
func (us *DashBoardService) GetDashBoardsLists(username string) (string, error) {
	var userSettings []*model.TableUserSettings
	if err := us.Session.Table("user_settings").Where("category = 'dashboard' AND (username = ? OR (data ->> 'shared' = 'true' OR data ->> 'shared' = '1'))", username).
		Find(&userSettings).Error; err != nil {
		return "", err
	}
	data, _ := json.Marshal(userSettings)
	dashboardElement := model.DashBoardElement{
		CssClass: "fa",
		Href:     "",
		Id:       "",
		Name:     "undefined",
		Param:    "",
		Shared:   0,
		Type:     0,
		Weight:   10,
	}
	rows, _ := gabs.ParseJSON(data)
	dashboardList := []model.DashBoardElement{}
	count := 0
	for _, row := range rows.Children() {
		count++
		if row.Exists("param") {
			dashboardElement.Href = row.S("param").Data().(string)
			dashboardElement.Id = row.S("param").Data().(string)
		}
		if row.Exists("data") {
			dashboardElement.Name = row.S("data", "name").Data().(string)
			dashboardElement.Weight = row.S("data", "weight").Data().(float64)

			if row.Exists("data", "shared") {
				if heputils.CheckBoolValue(row.S("data", "shared").Data()) {
					dashboardElement.Shared = 1
				} else {
					dashboardElement.Shared = 0
				}
			}
		}
		dashboardList = append(dashboardList, dashboardElement)
	}

	sort.Slice(dashboardList[:], func(i, j int) bool {
		return dashboardList[i].Name < dashboardList[j].Name
	})

	reply := gabs.New()
	reply.Set(count, "total")
	reply.Set(dashboardList, "data")
	reply.Set("ok", "status")
	reply.Set("ok", "auth")

	return reply.String(), nil
}

// this method gets all users from database
func (us *DashBoardService) GetDashBoard(username, param string) (string, error) {
	var userSettings model.TableUserSettings
	if err := us.Session.Table("user_settings").Where("(username = ? OR (data ->> 'shared' = 'true' OR data ->> 'shared' = '1')) AND category = 'dashboard' and param = ? ", username, param).
		Find(&userSettings).Error; err != nil {
		return "", err
	}

	//data, _ := json.Marshal(userSettings)
	data := userSettings.Data
	reply := gabs.New()
	reply.Set(1, "total")
	reply.Set(data, "data")
	reply.Set("ok", "status")
	reply.Set("ok", "auth")

	return reply.String(), nil
}

// this method gets all users from database
func (us *DashBoardService) InsertDashboard(username, dashboardId string, data json.RawMessage) (string, error) {

	newDashboard := model.TableUserSettings{}
	u2 := uuid.NewV4()
	newDashboard.GUID = u2.String()
	newDashboard.Param = dashboardId
	newDashboard.PartId = 10
	newDashboard.UserName = username
	newDashboard.Category = "dashboard"
	newDashboard.CreateDate = time.Now()
	newDashboard.Data = data

	var userSettings model.TableUserSettings
	result := us.Session.Table("user_settings").Where("(username != ? AND category = 'dashboard' and param = ? and partid = ?)", newDashboard.UserName,
		newDashboard.Param, newDashboard.PartId).Find(&userSettings)

	if result.RowsAffected > 0 {
		return "", fmt.Errorf("this is not my dashboard")
	}

	if err := us.Session.Debug().Table("user_settings").Where("username = ? AND category = 'dashboard' and param = ? and partid = ? ", newDashboard.UserName, newDashboard.Param, newDashboard.PartId).
		Delete(&model.TableUserSettings{}).Error; err != nil {
		return "", err
	}

	if err := us.Session.Debug().Table("user_settings").Save(&newDashboard).Error; err != nil {
		return "", err
	}

	reply := gabs.New()
	reply.Set(3, "total")
	reply.Set("ok", "status")
	reply.Set("ok", "auth")

	return reply.String(), nil
}

// this method gets all users from database
func (us *DashBoardService) UpdateDashboard(username, dashboardId string, data json.RawMessage) (string, error) {

	newDashboard := model.TableUserSettings{}
	u2 := uuid.NewV4()
	newDashboard.GUID = u2.String()
	newDashboard.Param = dashboardId
	newDashboard.PartId = 10
	newDashboard.UserName = username
	newDashboard.Category = "dashboard"
	newDashboard.CreateDate = time.Now()
	newDashboard.Data = data

	var userSettings model.TableUserSettings
	result := us.Session.Table("user_settings").Where("(username = ? AND category = 'dashboard' and param = ? and partid = ?)", newDashboard.UserName,
		newDashboard.Param, newDashboard.PartId).Find(&userSettings)

	if result.RowsAffected == 0 {
		return "", fmt.Errorf("dashboard doesn't exist")
	}

	if err := us.Session.Debug().Table("user_settings").Update(&newDashboard).Error; err != nil {
		return "", err
	}

	reply := gabs.New()
	reply.Set(3, "total")
	reply.Set("ok", "status")
	reply.Set("ok", "auth")

	return reply.String(), nil
}

// this method gets all users from database
func (us *DashBoardService) DeleteDashboard(username, dashboardId string) (string, error) {
	if err := us.Session.Debug().Table("user_settings").Where("username = ? AND category = 'dashboard' and param = ? and partid = ? ", username, dashboardId, 10).
		Delete(&model.TableUserSettings{}).Error; err != nil {
		return "", err
	}
	reply := gabs.New()
	reply.Set(3, "total")
	reply.Set("ok", "status")
	reply.Set("ok", "auth")
	return reply.String(), nil
}
