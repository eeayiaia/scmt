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
function backup_file(){
	BACKUP_FOLDER=~/.scmt-backup
	DATE_STAMP=$(date "+%b_%d_%Y_%H:%M:%S")
	BACKUP_FILE_NAME=$1
	if [[ ! -d $BACKUP_FOLDER ]]; then
		mkdir $BACKUP_FOLDER
	fi

	cp $BACKUP_FILE_NAME $BACKUP_FOLDER/$BACKUP_FILE_NAME-$DATE_STAMP
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
