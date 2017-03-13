package datasource

import (
	models "bunker/models"

	"errors"

	"gopkg.in/mgo.v2/bson"
)

func AddCompany(company *models.Company) error {
	
	valid := nameAvailable(company.Name)

	if !valid {
		return errors.New("Company name already taken")
	}

	tempSession := GetDBSession()
	defer CloseDBSession(tempSession)
	coll := tempSession.DB(MongoDatabase).C(CompanyCollection)
	coll.Insert(&company)
	return nil
}

func GetCompanyByName(name string) []models.Company {
	tempSession := GetDBSession()
	defer CloseDBSession(tempSession)
	coll := tempSession.DB(MongoDatabase).C(CompanyCollection)
	var companies []models.Company
	coll.Find(bson.M{"name": name}).All(&companies)

	return companies
}

func nameAvailable(companyName string) bool {
	tempSession := GetDBSession()
	defer CloseDBSession(tempSession)

	coll := tempSession.DB(MongoDatabase).C(CompanyCollection)
	totalUsed, _ := coll.Find(bson.M{"name": companyName}).Count()

	return totalUsed == 0
}
