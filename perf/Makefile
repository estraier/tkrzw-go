# Makefile for the performance checker for Tkrzw

GOCMD := go
RUNENV := LD_LIBRARY_PATH=.:/lib:/usr/lib:/usr/local/lib:$(HOME)/lib:$(HOME)/local/lib:$(LD_LIBRARY_PATH)

build :
	$(RUNENV) $(GOCMD) get
	$(RUNENV) $(GOCMD) build perf.go

run :
	$(RUNENV) $(GOCMD) get
	$(RUNENV) $(GOCMD) run perf.go --path casket.tkh --params "num_buckets=100000" \
	  --iter 20000 --threads 5
	$(RUNENV) $(GOCMD) run perf.go --path casket.tkh --params "num_buckets=100000" \
	  --iter 20000 --threads 5 --random
	$(RUNENV) $(GOCMD) run perf.go --path casket.tkt --params "key_comparator=decimal" \
	  --iter 20000 --threads 5
	$(RUNENV) $(GOCMD) run perf.go --path casket.tkt --params "key_comparator=decimal" \
	  --iter 20000 --threads 5 --random
	$(RUNENV) $(GOCMD) run perf.go --path casket.tks --params "step_unit=3" \
	  --iter 20000 --threads 5
	$(RUNENV) $(GOCMD) run perf.go --path casket.tks --params "step_unit=3" \
	  --iter 20000 --threads 5 --random
	$(RUNENV) $(GOCMD) run perf.go --params "dbm=tiny,num_buckets=100000" \
	  --iter 20000 --threads 5 --random
	$(RUNENV) $(GOCMD) run perf.go --params "dbm=baby,key_comparator=decimal" \
	  --iter 20000 --threads 5 --random
	$(RUNENV) $(GOCMD) run perf.go --params "dbm=stdhash,num_buckets=100000" \
	  --iter 20000 --threads 5 --random
	$(RUNENV) $(GOCMD) run perf.go --params "dbm=stdtree" \
	  --iter 20000 --threads 5 --random

vet :
	$(RUNENV) $(GOCMD) get
	$(RUNENV) $(GOCMD) vet

fmt :
	$(RUNENV) $(GOCMD) get
	$(RUNENV) $(GOCMD) fmt

clean :
	rm -rf perf casket* *.tkh *.tkt *.tks *~ hoge moge tako ika uni
