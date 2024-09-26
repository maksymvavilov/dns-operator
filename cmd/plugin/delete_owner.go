package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/rand"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/kuadrant/dns-operator/api/v1alpha1"
	"github.com/kuadrant/dns-operator/cmd/plugin/fixtures"
	"github.com/kuadrant/dns-operator/internal/controller"
)

const (
	deleteOwnerCMD = "delete-owner"
	unsetOwnerID   = "unset"
	defaultNS      = "kuadrant-system" // if this ns is not present, most likely controller is not running
	defaultHost    = "root.host"
	fixtureNS      = "targetNS"
	fixtureOwner   = "ownerID"
	fixtureHost    = "rootHost"
)

var (
	ownerID, rootHost, targetNS string
)

// Defines logic for the delete-owner subcommand
func deleteOwner() {
	// redefine flag set to correctly parse flags
	deleteOwnerSubcommand := flag.NewFlagSet(deleteOwnerCMD, flag.ExitOnError)

	// define flags for a new flag set
	deleteOwnerSubcommand.StringVar(&ownerID, "ownerID", unsetOwnerID, "The owner ID to be removed [required]")
	deleteOwnerSubcommand.StringVar(&rootHost, "rootHost", defaultHost, "Arbitrary root host to publish the record [optional]")
	deleteOwnerSubcommand.StringVar(&targetNS, "targetNS", defaultNS, "Arbitrary namespace to create the record in [optional]")

	// parse flags while consuming subcommand token
	_ = deleteOwnerSubcommand.Parse(flag.Args()[1:])

	// if the owner ID is not set - bail while outputting usage of this subcommand.
	if ownerID == unsetOwnerID {
		deleteOwnerSubcommand.PrintDefaults()
		os.Exit(1)
	}

	// make root host random to prevent conflicts in the provider
	rootHost = rand.String(10) + "." + rootHost

	// read record from the const
	record := &v1alpha1.DNSRecord{}
	err := resourceFromConst(fixtures.DNSRecordOwnerDeletion, record, func(s string) string {
		// map flags to placeholders in the string
		switch s {
		case fixtureNS:
			return targetNS
		case fixtureOwner:
			return ownerID
		case fixtureHost:
			return rootHost
		default:
			return ""
		}
	})
	if err != nil {
		fmt.Printf("Failed parsing fixture: %v\n", err)
		os.Exit(1)
	}

	// create record
	err = k8sClient.Create(ctx, record)
	if err != nil {
		fmt.Printf("Failed to create DNS Record: %v\n", err)
	}

	// wait for the finalizer
	end := time.Now().Add(time.Second * 5)
	var deleted bool
	for time.Now().Before(end) {
		err = k8sClient.Get(ctx, client.ObjectKeyFromObject(record), record)
		// not ignoring not found errors as Record MUST be present
		if err != nil {
			fmt.Printf("Failed to get DNS Record: %v\n", err)
			os.Exit(1)
		}
		if controllerutil.ContainsFinalizer(record, controller.DNSRecordFinalizer) {
			// mark for deletion once we have finalizer
			err = k8sClient.Delete(ctx, record)
			deleted = true
			if err != nil {
				fmt.Printf("Failed to delete DNS Record: %v\n", err)
				os.Exit(1)
			}
			break
		}
	}

	// timeout waiting for the finalizer
	if !deleted {
		fmt.Println("Timeout waiting fot the finalizer on  the DNS Record.\n" +
			"Ensure the dns-operator is running and try again.")
		// we still want to clean up DNSRecord CR
		err = k8sClient.Delete(ctx, record)
		if client.IgnoreNotFound(err) != nil {
			fmt.Printf("Failed to delete DNS Record: %v\n", err)
		}
		os.Exit(1)
	}

	fmt.Println("OwnerID " + ownerID + " was successfully cleaned up")
}

// Dnstead of using .yaml it is better to hold string in memory.
// Executable must be moved to PATH for proper integration with kubectl, so reading from the file is not reliable
func resourceFromConst(constant string, destObject runtime.Object, expandFunc func(string) string) error {
	decode := serializer.NewCodecFactory(scheme).UniversalDeserializer().Decode
	stream := []byte(os.Expand(constant, expandFunc))
	_, _, err := decode(stream, nil, destObject)
	return err
}
