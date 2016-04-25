package database
//Accessing information about devices
import (
	"errors"
	log "github.com/Sirupsen/logrus"
)

func AllDeviceStatus() {
	//This is the function to be run when a user wants to list node status..
	/*db, _ := NewConnection()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM devices")*/
}

//Returns all installed plugins on master
func GetAllInstalledPlugins() ([]string, error) {
	db, _ := NewConnection()
	defer db.Close()
	//get names of all installed plugins on master
	rows, err := db.Query("SELECT name FROM plugins WHERE installedOnMaster = 1") 
	defer rows.Close()
	if err != nil {
		return nil, errors.New("Could not get all installed plugins from database")
	}
	result := make([]string, 0)
	//iterate rows and append them to result
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			Log.WithFields(log.Fields{
				"row":	name,
			}).Warn("Could not parse plugin name")
			break
		}
		result = append(result, name)
	}
	err = rows.Err()
	if err != nil {
		return result, errors.New("Failed to get plugin names from database")
	}
	return result, nil

}