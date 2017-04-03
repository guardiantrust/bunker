package datasource

import (
	models "bunker/models"

	"gopkg.in/mgo.v2/bson"
)

// AddBatch - Add a new Batch
func AddBatch(batch *models.Batch) error {
	tempSession := GetDBSession()
	defer CloseDBSession(tempSession)
	coll := tempSession.DB(MongoDatabase).C(BatchCollection)
	batch.BatchID = bson.NewObjectId()
	coll.Insert(&batch)
	return nil
}

// GetBatchByID - Get a batch
func GetBatchByID(batchID string) (models.Batch, error) {
	tempSession := GetDBSession()
	defer CloseDBSession(tempSession)
	coll := tempSession.DB(MongoDatabase).C(BatchCollection)
	var batch models.Batch
	err := coll.FindId(IDToObjectID(batchID)).One(&batch)
	return batch, err
}

// GetBatchesByCompanyID - Get all the batches by a companies ID
func GetBatchesByCompanyID(companyID string) ([]models.Batch, error) {
	tempSession := GetDBSession()
	defer CloseDBSession(tempSession)
	coll := tempSession.DB(MongoDatabase).C(BatchCollection)
	var batches []models.Batch
	err := coll.Find(bson.M{"companyID": IDToObjectID(companyID)}).All(&batches)
	return batches, err
}

// UpdateBatch - Update a batch
func UpdateBatch(batchID string, batch *models.Batch) error {
	tempSession := GetDBSession()
	defer CloseDBSession(tempSession)
	coll := tempSession.DB(MongoDatabase).C(BatchCollection)
	err := coll.Update(bson.M{"_id": IDToObjectID(batchID)}, bson.M{"$set": bson.M{"name": batch.CompanyID, "description": batch.Description, "attributes": batch.Attributes}})
	return err
}

// DeleteBatchByID - Deletes a batch by its ID
func DeleteBatchByID(batchID string) error {
	tempSession := GetDBSession()
	defer CloseDBSession(tempSession)
	coll := tempSession.DB(MongoDatabase).C(BatchCollection)
	err := coll.RemoveId(IDToObjectID(batchID))
	return err
}
