package datasource

import (
	models "bunker/models"

	"errors"

	"gopkg.in/mgo.v2/bson"
)

// AddCompany - Add a new company
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

// UpdateCompany - Update a company profile
func UpdateCompany(companyID string, company *models.Company) error {
	tempSession := GetDBSession()
	defer CloseDBSession(tempSession)
	coll := tempSession.DB(MongoDatabase).C(CompanyCollection)
	err := coll.Update(bson.M{"_id": companyID}, bson.M{"$set": bson.M{"name": company.Name, "address1": company.Address1, "address2": company.Address2, "city": company.City, "state": company.State, "postal": company.Postal, "phoneNumber": company.PhoneNumber, "contactEmail": company.ContactEmail, "isSuspended": company.IsSuspended, "isActive": company.IsActive}})

	return err
}

// InActivateCompany - Change company activation
func InActivateCompany(companyID string) error {
	tempSession := GetDBSession()
	defer CloseDBSession(tempSession)
	coll := tempSession.DB(MongoDatabase).C(CompanyCollection)
	err := coll.Update(bson.M{"_id": companyID}, bson.M{"$set": bson.M{"isActive": false}})

	if err != nil {
		return err
	}

	return nil
}

// GetCompanyByName - Get a company by company name
func GetCompanyByName(name string) []models.Company {
	tempSession := GetDBSession()
	defer CloseDBSession(tempSession)
	coll := tempSession.DB(MongoDatabase).C(CompanyCollection)
	var companies []models.Company
	coll.Find(bson.M{"name": name}).All(&companies)

	return companies
}

// GetCompany - return a company by Id
func GetCompany(companyID string) (models.Company, error) {
	tempSession := GetDBSession()
	defer CloseDBSession(tempSession)
	coll := tempSession.DB(MongoDatabase).C(CompanyCollection)
	var company models.Company
	err := coll.Find(bson.M{"_id": IDToObjectID(companyID)}).One(&company)

	return company, err
}

// nameAvailable - Check that a company name is not already saved
func nameAvailable(companyName string) bool {
	tempSession := GetDBSession()
	defer CloseDBSession(tempSession)

	coll := tempSession.DB(MongoDatabase).C(CompanyCollection)
	totalUsed, _ := coll.Find(bson.M{"name": companyName}).Count()

	return totalUsed == 0
}
