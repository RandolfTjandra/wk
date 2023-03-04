# Installation package params
PROJECT_NAME := wk 
DIST_DIR  := dist

build-wk:
	rm -rf ./$(DIST_DIR)/wk
	mkdir -p $(DIST_DIR)/wk
	go build -o $(DIST_DIR)/wk/wk \
			 ./wk/main.go \
			 ./wk/model.go \
			 ./wk/command.go \
			 ./wk/mock_command.go \
			 ./wk/views.go

package-wk: build-wk
	cd $(DIST_DIR)/wk && zip wk.zip ./wk 

package-all: \
package-wk


