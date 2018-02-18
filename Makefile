PROGRAM=	alexa_go_theoffice


all: $(PROGRAM) package

$(PROGRAM):
	@echo "                -=[ Building $(PROGRAM) ]=-"
	GOOS=linux GOARCH=amd64 go build $(PROGRAM).go
	@echo "                       -=[ DONE ]=-"
	
package:
	@echo "            -=[ Compressing $(PROGRAM).zip ]=-"
	zip $(PROGRAM).zip $(PROGRAM)
	@echo "                       -=[ DONE ]=-"
	
clean:
	@echo "                    -=[ Cleaning Up ]=-"
	rm -f $(PROGRAM).zip $(PROGRAM)
	@echo "                       -=[ DONE ]=-"
