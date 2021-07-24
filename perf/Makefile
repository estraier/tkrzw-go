# Makefile for the performance checker for Tkrzw

GOCMD := go
RUNENV := LD_LIBRARY_PATH=.:/lib:/usr/lib:/usr/local/lib:$(HOME)/lib

build :
	$(RUNENV) $(GOCMD) build perf.go

run :
	$(RUNENV) $(GOCMD) run perf.go

vet :
	$(RUNENV) $(GOCMD) vet

fmt :
	$(RUNENV) $(GOCMD) fmt

clean :
	rm -rf casket* *.tkh *.tkt *.tks *~ hoge moge tako ika uni \
	  perf
