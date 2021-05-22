## General 
* Specify your credentals in [Settings file](cmd/settings.json)
* Run one of the following commands:
  * create - creates new container.
  * open - extracts container and prints content.
  * add-signature - adds new signature to existing container.
  * remove-signature - removes existing signature by id. if no signature by given id found, error is returned.

## Commands and parameters:

### create
> go run main.go create < comma separated list of files to add container > < container's path where to save >

### open
> go run main.go open < container's path to extract >

### add-signature 
> go run main.go add-signature < container's path to add new signature >

### remove-signature
> go run main.go remove-signature < container's path to remove signature from > < signature id >