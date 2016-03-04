# Writes 80 column line
function write_line(){
	echo "--------------------------------------------------------------------------------"
}

# Check if installer is run as root, exiting otherwise
function check_root(){
	if [[ $EUID != 0 ]]; then
		echo "This installer must be run with root rights." 1>&2
		exit 1
	fi
}

# Check if aptitude is installed; if not, install it
function check_or_install_aptitude(){
	aptpath=$(which aptitude)

	if [ ! $aptpath ]; then
		echo "aptitude not found, installing..."

		write_line
		apt-get install aptitude
		write_line

		if [ $? != 0 ]; then
			echo "Failed to install aptitude. Can not download packages."
			exit 1
		fi
	else
		echo "aptitude found."
	fi
}

# Install a .deb package (residing in working directory) given package name
# .deb files are found as <given_name>*.deb
function install_pkg(){
	pkg_name=$1
	echo "Attempting to install $pkg_name..."

	pkg_ok=$(dpkg-query -W --showformat='${Status}\n' $pkg_name | grep "install ok installed")

	if [[ ! $pkg_ok ]]; then
		echo "$pkg_name is not alreadt installed."
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
