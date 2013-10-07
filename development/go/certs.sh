#!/bin/bash
# ./certs name country state location organization department common_name email_address
rm -rf certs
mkdir -p certs

make_cert() {
	echo "Creating $1 cert"
	openssl req -new -nodes -x509 -out certs/$1.pem -keyout certs/$1.key -days 365 -subj "/C=$2/ST=$3/L=$4/O=$5/OU=$6/CN=$7/emailAddress=$8"
}

if [ $# -lt 1 ] ; then
	ecgo "Creating generic certs"
	make_cert "server" "US" "FL" "FTL" "ACME" "SYSADM" "darq.lab" "root@localhost"
	make_cert "client" "US" "FL" "FTL" "ACME" "SYSADM" "darq.lab" "root@localhost"
else
	make_cert $1 $2 $3 $4 $5 $6 $7 $8
fi
