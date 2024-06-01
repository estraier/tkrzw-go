# Makefile for sample programs for Tkrzw

GOCMD := go
RUNENV := LD_LIBRARY_PATH=.:/lib:/usr/lib:/usr/local/lib:$(HOME)/lib:$(HOME)/local/lib:$(LD_LIBRARY_PATH)

build :
	$(RUNENV) $(GOCMD) get
	$(RUNENV) $(GOCMD) build main.go

run :
	$(RUNENV) $(GOCMD) get
	$(RUNENV) $(GOCMD) run main.go

vet :
	$(RUNENV) $(GOCMD) get
	$(RUNENV) $(GOCMD) vet

fmt :
	$(RUNENV) $(GOCMD) get
	$(RUNENV) $(GOCMD) fmt

clean :
	rm -rf casket* *.tkh *.tkt *.tks *~ hoge moge tako ika uni \
	  main
