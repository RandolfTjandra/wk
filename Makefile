# Installation package params
PROJECT_NAME := wk 
DIST_DIR  := dist

build-wk:
	rm -rf ./$(DIST_DIR)/wk
	mkdir -p $(DIST_DIR)/wk
	go build -o $(DIST_DIR)/wk/bootstrap \
			 ./wk/main.go \
			 ./wk/model.go \
			 ./wk/command.go

package-wk: build-wk
	cd $(DIST_DIR)/wk && zip wk.zip ./bootstrap 


package-all: \
package-wk


