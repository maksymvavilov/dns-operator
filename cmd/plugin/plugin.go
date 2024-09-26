package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kuadrant/dns-operator/api/v1alpha1"
)

const (
	helpCMD = "help"
)

var (
	scheme    = runtime.NewScheme()
	k8sClient client.Client
	ctx       context.Context
)

func init() {
	utilruntime.Must(v1alpha1.AddToScheme(scheme))
}

func main() {
	// override error message
	flag.CommandLine.Usage = func() {
		help()
		os.Exit(1)
	}
	// parse the current set of args
	flag.Parse()

	// check for subcommand
	if flag.NArg() == 0 { // it could be > 1 if it carries flags for the subcommand that weer not parsed yet
		help()
		os.Exit(1)
	}

	// set up a client
	ctx = context.Background()
	var err error
	k8sClient, err = client.New(controllerruntime.GetConfigOrDie(), client.Options{
		Scheme: scheme,
	})
	if err != nil {
		fmt.Printf("Failed to create k8s client: %v\n", err)
		os.Exit(1)
	}

	switch flag.Arg(0) {
	case helpCMD:
		help()
	case deleteOwnerCMD:
		deleteOwner()
	case getZoneRecordsCMD:
		getZoneRecordsCmd()
	default:
		unknownCMD()
	}
}

// Default error and a help message. Displays usage of this cli. Must define all subcommands
func help() {
	fmt.Println("Usage of dns plugin: \n\t" +
		deleteOwnerCMD + "\t Delete specified owner ID from the provider (requires running controller)" +
		"\nFor more details use -h or -help flags on the command")
}

func unknownCMD() {
	fmt.Println("Unknown command \nPlease refer to \"help\" for usage")
}
