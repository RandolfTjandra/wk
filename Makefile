# Installation package params
DIST_DIR  := dist

build-wk:
	rm -rf ./$(DIST_DIR)/wk
	mkdir -p $(DIST_DIR)/wk
	go build -o $(DIST_DIR)/wk/wk \
			 ./wk/main.go \
			 ./wk/model.go \
			 ./wk/command.go \
			 ./wk/mock_command.go \
			 ./wk/navigation.go \
			 ./wk/views.go \
			 ./wk/index_view.go \
			 ./wk/summary_view.go \
			 ./wk/assignments_view.go \
			 ./wk/reviews_view.go \
			 ./wk/account_view.go

package-wk: build-wk
	cd $(DIST_DIR)/wk && zip wk.zip ./wk 

package-all: \
package-wk

local: build-wk
	cp -f $(DIST_DIR)/wk/wk /usr/local/bin/

