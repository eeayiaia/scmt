# SCMT shell script utils
# Used by scripts and plugins for SuperK Cluster Management Toolkit

# Writes 80 column line
function write_line(){
	echo "--------------------------------------------------------------------------------"
}

# Check if script is run as root, exiting otherwise
function check_root(){
	if [[ $EUID != 0 ]]; then
		echo "This script must be run as root" 1>&2
		exit 1
	fi
}

# Check if variable is set, exiting otherwise
function assertIsSet(){
    if [[ ! ${!1} && ${!1-_} ]]; then
        echo "$1 is unset, exiting" 1>&2
        exit 1
    fi
}

# Check if script was invoked by SCMT, exiting otherwise
function check_invoked_by_scmt(){
	if [[ ! ${INVOKED_BY_SCMT} ]]; then
		echo "Error: this script is intended to be invoked by the SuperK "\
			"Cluster Management Toolkit, not used manually." 1>&2
		exit 1
	fi

	check_root
}

# Delete file safely
function delete_file(){
	local FILE=$1

	if [[ ! -z $2 ]]; then
		echo "Too many parameters passed to delete_file" 1>&2
		return
	fi

	if [[ ! -f "$FILE" ]]; then
		echo "delete_file: No such file:\n $FILE" 1>&2
	else
		echo "deleting file:\n$FILE"
		rm -- "${FILE:?}"
	fi
}

# Delete directory safely
function delete_directory(){
	local DIRECTORY=$1

	if [[ ! -z $2 ]]; then
		echo "Too many parameters passed to delete_directory." 1>&2
		return
	fi

	if [[ ! -d "$DIRECTORY" ]]; then
		echo "delete_directory: No such directory:\n $DIRECTORY" 1>&2
	else
		echo "deleting directory:\n$DIRECTORY"
		rm -rf -- "${DIRECTORY:?}"
	fi
}

# Backup file/directory with timestamp
# Parameter 1: File/directory to backup
# Output: BACKUP_OUTPUT will contain path to backup file/directory
function backup_file(){
	local BACKUP_FOLDER=~/.scmt-backup
	local DATE_STAMP=$(date "+%b_%d_%Y_%H:%M:%S")
	local BACKUP_FILE=$1
	local BACKUP_FILE_NAME=$(basename $BACKUP_FILE)

	BACKUP_OUTPUT=$BACKUP_FOLDER/$BACKUP_FILE_NAME-$DATE_STAMP

	if [[ ! -d $BACKUP_FOLDER ]]; then
		mkdir $BACKUP_FOLDER
	fi

	if [[ -d $BACKUP_FILE ]]; then
		echo "Backing up directory $BACKUP_FILE to $BACKUP_OUTPUT..."
		cp -r $BACKUP_FILE $BACKUP_OUTPUT
	elif [[ -f $BACKUP_FILE ]]; then
		echo "Backing up file $BACKUP_FILE to $BACKUP_OUTPUT..."
		cp $BACKUP_FILE $BACKUP_OUTPUT
	else
		echo "Cannot backup $BACKUP_FILE: path is not file or directory" 1>&2
	fi
}

# Creates a new user
# Parameter 1: username
# Parameter 2: password
# Parameter 3: uid
function create_user(){
	if [[ ! $1 ]]; then
		echo "Failed to create user: no username was provided."
		exit 1;
	fi

	if [[ ! $2 ]]; then
		echo "Failed to create user: no password was provided."
		exit 1;
	fi

	if [[ ! $3 ]]; then
		echo "Failed to create user: no UID was provided."
		exit 1;
	fi

	echo "Creating user '$1' with uid '$3'"
	adduser $1 --gecos "First Last,RoomNumber,WorkPhone,HomePhone" \
		--disabled-password --uid $3
	echo "$1:$2" | sudo chpasswd
}

