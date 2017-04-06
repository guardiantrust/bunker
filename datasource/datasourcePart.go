package datasource

import (
	"bufio"
	models "bunker/models"
	"errors"
	"io"
	multipart "mime/multipart"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// GETPart - get a part  by id
func GETPart(partID string) (models.Part, error) {
	tempSession := GetDBSession()
	defer CloseDBSession(tempSession)
	coll := tempSession.DB(MongoDatabase).C(PartCollection)
	var part models.Part
	err := coll.FindId(IDToObjectID(partID)).One(&part)
	return part, err
}

func GETPartsByMachine(machineID string) ([]models.Part, error) {
	tempSession := GetDBSession()
	defer CloseDBSession(tempSession)
	coll := tempSession.DB(MongoDatabase).C(PartCollection)
	var parts []models.Part
	err := coll.Find(

		bson.M{"$elemMatch": bson.M{"machineID": IDToObjectID(machineID),
			"processed": bson.M{"$lt": time.Now()},
		}}).All(&parts)

	if err != nil {
		return parts, err
	}

	return parts, err
}

func GetPartsByCompany(companyID string) ([]models.Part, error) {
	tempSession := GetDBSession()
	defer CloseDBSession(tempSession)
	coll := tempSession.DB(MongoDatabase).C(PartCollection)
	var parts []models.Part

	err := coll.Find(bson.M{"companyID": companyID, "files": bson.M{"$elemMatch": bson.M{"processed": bson.M{"$lt": time.Now()}}}}).All(&parts)

	return parts, err
}

// GETPartByBarcode - retrieve a part by barcode
func GETPartByBarcode(machineID string, barcode string) (models.Part, error) {
	tempSession := GetDBSession()
	defer CloseDBSession(tempSession)
	coll := tempSession.DB(MongoDatabase).C(PartCollection)
	var part models.Part
	err := coll.Find(bson.M{"barcode": barcode,
		"$elemMatch": bson.M{
			"machineID": bson.ObjectIdHex(machineID),
		},
	}).One(&part)

	return part, err
}

// AddPart - Add a part
func AddPart(part *models.Part) error {
	tempSession := GetDBSession()
	defer CloseDBSession(tempSession)

	part.ID = bson.NewObjectId()
	part.Created = time.Now()
	if len(part.BatchID) == 0 {
		part.BatchID = bson.NewObjectId()
	}

	for _, f := range part.Files {
		f.FileID = bson.NewObjectId()
		f.Created = time.Now()
	}

	coll := tempSession.DB(MongoDatabase).C(PartCollection)
	err := coll.Insert(part)
	return err
}

// ProcessPart - Mark a Part as processed
func ProcessPart(partID string, partFileID string) error {
	tempSession := GetDBSession()
	defer CloseDBSession(tempSession)
	coll := tempSession.DB(MongoDatabase).C(PartCollection)

	changes := bson.M{
		"files": bson.M{
			"$elemMatch": bson.M{
				"id": bson.ObjectIdHex(partFileID),
			},
		},
		"$set": bson.M{
			"processed": time.Now(),
		},
	}

	err := coll.UpdateId(IDToObjectID(partID), changes)
	return err
}

// UpdatePart - Update a part
func UpdatePart(partID string, part *models.Part) error {
	tempSession := GetDBSession()
	defer CloseDBSession(tempSession)
	coll := tempSession.DB(MongoDatabase).C(PartCollection)
	err := coll.UpdateId(IDToObjectID(partID), &part)
	return err
}

// DeletePart - Delete a part - all files are deleted as well
func DeletePart(partID string) error {
	tempSession := GetDBSession()
	defer CloseDBSession(tempSession)
	coll := tempSession.DB(MongoDatabase).C(PartCollection)
	err := deletePartFile(partID)

	if err != nil {
		return err
	}

	err = coll.RemoveId(IDToObjectID(partID))

	return err
}

func deletePartFile(partID string) error {
	tempSession := GetDBSession()
	defer CloseDBSession(tempSession)
	gridFS := tempSession.DB(MongoDatabase).GridFS("fs")

	part, err := GETPart(partID)

	for _, f := range part.Files {
		gridFS.RemoveId(f.FileID)
	}

	return err
}

// GetPartFile - return a file by fileId
func GetPartFile(fileID string) ([]byte, models.PartFile, error) {
	tempSession := GetDBSession()
	defer CloseDBSession(tempSession)
	db := tempSession.DB(MongoDatabase)
	file, _ := db.GridFS("fs").OpenId(IDToObjectID(fileID))
	var fileSize int
	fileSize = int(file.Size())
	returnFile := make([]byte, fileSize)
	_, err := file.Read(returnFile)

	file.Close()
	coll := tempSession.DB(MongoDatabase).C(PartCollection)

	var partFile models.PartFile
	err = coll.Find(bson.M{"fileID": fileID}).One(&partFile)

	return returnFile, partFile, err
}

// SavePartFile - Save a PartFile to the db
func SavePartFile(companyID string, fileID string, fileName string, file multipart.File) error {

	tempSession := GetNewSession()
	db := tempSession.DB(MongoDatabase)
	defer CloseDBSession(tempSession)

	gridFile, err := db.GridFS("fs").Create(fileName)

	if err != nil {
		return err
	}

	gridFile.SetId(fileID)
	gridFile.SetName(fileName)
	gridFile.SetMeta(bson.M{"companyID": companyID})

	reader := bufio.NewReader(file)
	defer func() { file.Close() }()

	// make a buffer to keep chunks that are read
	buf := make([]byte, 1024)
	for {
		// read a chunk
		n, err := reader.Read(buf)

		if err != nil && err != io.EOF {
			return errors.New("Could not read the input file")
		}

		if n == 0 {
			break
		}

		// write a chunk
		if _, err := gridFile.Write(buf[:n]); err != nil {
			return errors.New("Could not write to GridFs for " + gridFile.Name())
		}
	}

	gridFile.Close()

	return err
}
