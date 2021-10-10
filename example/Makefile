# Makefile for sample programs for Tkrzw

GOCMD := go
RUNENV := LD_LIBRARY_PATH=.:/lib:/usr/lib:/usr/local/lib:$(HOME)/lib

build :
	$(RUNENV) $(GOCMD) build example1.go
	$(RUNENV) $(GOCMD) build example2.go
	$(RUNENV) $(GOCMD) build example3.go
	$(RUNENV) $(GOCMD) build example4.go

run :
	$(RUNENV) $(GOCMD) run example1.go
	$(RUNENV) $(GOCMD) run example2.go
	$(RUNENV) $(GOCMD) run example3.go
	$(RUNENV) $(GOCMD) run example4.go

vet :
	$(RUNENV) $(GOCMD) vet

fmt :
	$(RUNENV) $(GOCMD) fmt

clean :
	rm -rf casket* *.tkh *.tkt *.tks *~ hoge moge tako ika uni \
	  example1 example2 example3 example4
