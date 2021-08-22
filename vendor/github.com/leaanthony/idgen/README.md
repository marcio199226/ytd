# idgen

`idgen` only does 1 thing: generates unique uints. IDs may be released so they can be reused. It can provide "MaxUint" unique IDs.

## Why did you make this?

I needed a way of generating a load of IDs for Windows Menus & MenuItems. They will sometimes need to be regenerated and thus the IDs 
would need recalculating. Whilst it was probably ok to just increment a counter, I needed a way to guarantee that some long running
process wasn't going to hit that upper limit and screw up the application.

## Installation

`go get github.com/leaanthony/idgen`

## Usage

```go
package main

import "github.com/leaanthony/idgen"

func main() {
	
    // Create a new generator
    generator := idgen.New()	
    
    // Get an ID
    id, err := generator.NewID()
    if err != nil {
    	// We have run out of available IDs
    }
    
    limited := generator.NewWithMaximum(1)

    // Get the only ID available
    id, err = limited.NewID()
    // err == nil
    
    // Get another when no more are available
    _, err = limited.NewID()
    // err != nil 
	
    // Release the id we had
    limited.ReleaseID(id)
    
    // Now we can get the ID again
    id, err = limited.NewID()
    // err == nil
}
```
