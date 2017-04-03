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

func GETPart(partID string) (models.Part, error) {
	tempSession := GetDBSession()
	defer CloseDBSession(tempSession)
	coll := tempSession.DB(MongoDatabase).C(PartCollection)
	var part models.Part
	err := coll.FindId(IDToObjectID(partID)).One(&part)
	return part, err
}

func AddPart(part models.Part) error {
	tempSession := GetDBSession()
	defer CloseDBSession(tempSession)
	coll := tempSession.DB(MongoDatabase).C(PartCollection)
	err := coll.Insert(part)
	return err
}

func UpdatePart(partID string, part *models.Part) error {
	tempSession := GetDBSession()
	defer CloseDBSession(tempSession)
	coll := tempSession.DB(MongoDatabase).C(PartCollection)
	err := coll.UpdateId(IDToObjectID(partID), &part)
	return err
}

// GetPartFile - return a file by fileId
func GetPartFile(fileID string) ([]byte, error) {
	tempSession := GetDBSession()
	db := tempSession.DB(MongoDatabase)
	file, _ := db.GridFS("fs").OpenId(IDToObjectID(fileID))
	var fileSize int
	fileSize = int(file.Size())
	returnFile := make([]byte, fileSize)
	n, err := file.Read(returnFile)

	if err != nil {
		return returnFile, err
	}

	if n != fileSize {

	}

	file.Close()
	defer CloseDBSession(tempSession)
	return returnFile, err
}

// SavePartFile - Save a PartFile to the db
func SavePartFile(companyID string, partID string, fileAttributes models.PartFile, file multipart.File) error {

	tempSession := GetNewSession()
	db := tempSession.DB(MongoDatabase)
	defer CloseDBSession(tempSession)
	fileAttributes.FileID = bson.NewObjectId()
	fileAttributes.Created = time.Now()
	gridFile, err := db.GridFS("fs").Create(fileAttributes.FileName)

	if err != nil {
		return err
	}

	gridFile.SetId(fileAttributes.FileID)
	gridFile.SetContentType(fileAttributes.FileExtension)
	gridFile.SetName(fileAttributes.FileName)
	gridFile.SetMeta(bson.M{"machineID": fileAttributes.MachineID})
	gridFile.SetMeta(bson.M{"companyID": companyID})
	gridFile.SetMeta(bson.M{"partID": partID})

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
