package datasource

import (
	models "bunker/models"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func AddMachine(machine *models.Machine) error {
	tempSession := GetDBSession()
	defer CloseDBSession(tempSession)

	machine.Created = time.Now()
	machine.ID = bson.NewObjectId()

	coll := tempSession.DB(MongoDatabase).C(MachineCollection)
	err := coll.Insert(machine)

	return err
}

func GetMachineById(companyID string, machineID string) (models.Machine, error) {
	tempSession := GetDBSession()
	defer CloseDBSession(tempSession)
	coll := tempSession.DB(MongoDatabase).C(MachineCollection)

	var machine models.Machine
	err := coll.FindId(IDToObjectID(machineID)).One(&machine)

	if machine.CompanyID != IDToObjectID(companyID) {
		invalid := models.Machine{}
		return invalid, err
	}

	return machine, err
}

func GetMachinesByCompany(companyID string) ([]models.Machine, error) {
	tempSession := GetDBSession()
	defer CloseDBSession(tempSession)
	coll := tempSession.DB(MongoDatabase).C(MachineCollection)

	var machines []models.Machine
	err := coll.Find(bson.M{"companyID": companyID, "isActive": true}).All(&machines)

	return machines, err
}

func InactivateMachine(companyID string, machineID string) error {
	tempSession := GetDBSession()
	defer CloseDBSession(tempSession)
	coll := tempSession.DB(MongoDatabase).C(MachineCollection)

	err := coll.Update(bson.M{"_id": machineID, "companyID": IDToObjectID(companyID), "isActive": true}, bson.M{"$set": bson.M{"isActive": false}})

	return err
}

func UpdateMachine(machineID string, machine models.Machine) error {
	tempSession := GetDBSession()
	defer CloseDBSession(tempSession)
	coll := tempSession.DB(MongoDatabase).C(MachineCollection)

	err := coll.UpdateId(IDToObjectID(machineID), bson.M{"$set": bson.M{"name": machine.Name, "description": machine.Description, "location": machine.Location, "isActive": machine.IsActive}})

	return err
}
