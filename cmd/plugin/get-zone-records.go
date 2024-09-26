package main

import "flag"

const (
	getZoneRecordsCMD = "get-zone-records"
)

var (
	policyName, host string
)

func getZoneRecordsCmd() {
	// redefine flag set to correctly parse flags
	getRecordsSubcommand := flag.NewFlagSet(getZoneRecordsCMD, flag.ExitOnError)

	// define flags for a new flag set
	getRecordsSubcommand.StringVar(&ownerID, "ownerID", unsetOwnerID, "The owner ID to be removed [required]")
	getRecordsSubcommand.StringVar(&rootHost, "rootHost", defaultHost, "Arbitrary root host to publish the record [optional]")
	getRecordsSubcommand.StringVar(&targetNS, "targetNS", defaultNS, "Arbitrary namespace to create the record in [optional]")

	// parse flags while consuming subcommand token
	_ = getRecordsSubcommand.Parse(flag.Args()[1:])

}
