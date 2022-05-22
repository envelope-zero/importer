package api

import (
	"fmt"

	"github.com/envelope-zero/importer/pkg/types"
	"github.com/google/uuid"
)

// TODO: Make the creator analyze the creation structs (with names instead of IDs) and query the API
// Creator should have a running list of resources already stored as map of name: id
//
// Resource creation should be accompanied with a permanentFail bool and a numFails counter to make 3 retries
// and for 404’s like a parent resource not being found

func CreateResources(id uuid.UUID, ch <-chan types.ResourceCerate, done chan<- bool) {
	fmt.Println("Receiving resources…")
	for {
		r, more := <-ch
		if more {
			fmt.Printf("%#v\n", r)
		} else {
			fmt.Println("Completed!")
			done <- true
			break
		}
	}
}
