# Writes 80 column line
function write_line(){
	echo "--------------------------------------------------------------------------------"
}

# Check if installer is run as root, exiting otherwise
function check_root(){
	if [[ $EUID != 0 ]]; then
		echo "This installer must be run with root privileges." 1>&2
		exit 1
	fi
}

# Backup file/directory with timestamp
# Parameter 1: File/directory to backup
# Output: BACKUP_OUTPUT will contain path to backup file/directory
function backup_file(){
	BACKUP_FOLDER=~/.scmt-backup
	DATE_STAMP=$(date "+%b_%d_%Y_%H:%M:%S")
	BACKUP_FILE=$1
	BACKUP_FILE_NAME=$(basename $BACKUP_FILE)

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
	adduser $1 --gecos "First Last,RoomNumber,WorkPhone,HomePhone" --disabled-password --uid $3
	echo "$1:$2" | sudo chpasswd
}

# Check if aptitude is installed; if not, install it
# (aptitude is used to unconditionally download .deb packages to working
# directory without installing)
function check_or_install_aptitude(){
	aptpath=$(which aptitude)

	if [ ! $aptpath ]; then
		echo "aptitude not found, installing..."

		write_line
		apt-get install aptitude
		install_success=$?
		write_line

		if [ $install_success != 0 ]; then
			echo "Failed to install aptitude. Can not download packages."
			exit 1
		fi
	else
		echo "aptitude found."
	fi
}

# Install a .deb package (residing in working directory) given package name
# Usage: install_pkg package_name
# .deb files are found as <given_name>*.deb
function install_pkg(){
	pkg_name=$1
	echo "Attempting to install $pkg_name..."

	pkg_ok=$(dpkg-query -W --showformat='${Status}\n' $pkg_name | grep "install ok installed")

	if [[ ! $pkg_ok ]]; then
		echo "$pkg_name is not already installed."
		pkg_filename=$(find . -name $pkg_name*.deb)

		if [[ $pkg_filename ]]; then
			echo "Found $pkg_name: $pkg_filename"
			echo "INSTALLING $pkg_name"
			write_line
			dpkg -i $pkg_filename
			install_success=$?
			write_line

			if [ $install_success == 0 ]; then
				echo "$pkg_name installed successfully"
			else
				echo "Error: failed to install $pkg_name"
				exit 1
			fi

			write_line
		else
			echo "Error: could not find package '$pkg_name*.deb'"
			exit 1
		fi
	else
		echo "$pkg_name is already installed."
	fi
}
