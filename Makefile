# Makefile for Tkrzw for Go

PACKAGE := tkrzw-go
VERSION := 0.1.27
PACKAGEDIR := $(PACKAGE)-$(VERSION)
PACKAGETGZ := $(PACKAGE)-$(VERSION).tar.gz

GOCMD := go
RUNENV := LD_LIBRARY_PATH=.:/lib:/usr/lib:/usr/local/lib:$(HOME)/lib

build :
	$(RUNENV) $(GOCMD) get
	$(RUNENV) $(GOCMD) build
	[ ! -f perf/Makefile ] || cd perf && $(MAKE)
	[ ! -f wicked/Makefile ] || cd wicked && $(MAKE)
	@printf '\n'
	@printf '#================================================================\n'
	@printf '# Build is OK.\n'
	@printf '#================================================================\n'

check : test runperf runwicked
	@printf '\n'
	@printf '#================================================================\n'
	@printf '# Checking completed.\n'
	@printf '#================================================================\n'

test :
	$(RUNENV) $(GOCMD) get
	$(RUNENV) $(GOCMD) test -v

runperf :
	[ ! -f perf/Makefile ] || cd perf && $(MAKE) run

runwicked :
	[ ! -f wicked/Makefile ] || cd wicked && $(MAKE) run

vet :
	$(RUNENV) $(GOCMD) get
	$(RUNENV) $(GOCMD) vet

fmt :
	$(RUNENV) $(GOCMD) get
	$(RUNENV) $(GOCMD) fmt

clean :
	rm -rf casket* *.tkh *.tkt *.tks *~ hoge moge tako ika uni go.sum
	[ ! -f perf/Makefile ] || cd perf && $(MAKE) clean
	[ ! -f wicked/Makefile ] || cd wicked && $(MAKE) clean

install :
	@printf '\n'
	@printf '#================================================================\n'
	@printf '# Installation is not necessary.\n'
	@printf '# Just import "github.com/estraier/tkrzw-go" in your Go code.\n'
	@printf '#================================================================\n'

dist :
	$(MAKE) fmt
	[ ! -f perf/Makefile ] || cd perf && $(MAKE) fmt
	[ ! -f wicked/Makefile ] || cd wicked && $(MAKE) fmt
	[ ! -f example1/Makefile ] || cd example1 && $(MAKE) fmt
	[ ! -f example2/Makefile ] || cd example2 && $(MAKE) fmt
	[ ! -f example3/Makefile ] || cd example3 && $(MAKE) fmt
	[ ! -f example4/Makefile ] || cd example4 && $(MAKE) fmt
	$(MAKE) distclean
	rm -Rf "../$(PACKAGEDIR)" "../$(PACKAGETGZ)"
	cd .. && cp -R tkrzw-go $(PACKAGEDIR) && \
	  tar --exclude=".*" -cvf - $(PACKAGEDIR) | gzip -c > $(PACKAGETGZ)
	rm -Rf "../$(PACKAGEDIR)"
	sync ; sync

distclean : clean apidocclean
	[ ! -f perf/Makefile ] || cd perf && $(MAKE) clean
	[ ! -f wicked/Makefile ] || cd wicked && $(MAKE) clean
	[ ! -f example1/Makefile ] || cd example1 && $(MAKE) clean
	[ ! -f example2/Makefile ] || cd example2 && $(MAKE) clean
	[ ! -f example3/Makefile ] || cd example3 && $(MAKE) clean
	[ ! -f example4/Makefile ] || cd example4 && $(MAKE) clean

apidoc :
	rm -rf api-doc
	PATH=$$PATH:$$HOME/go/bin:$$HOME/.local/bin \
	  godoc -http "localhost:8080" -play -goroot /usr/share/go & sleep 2
	mkdir api-doc
	curl -s "http://localhost:8080/lib/godoc/style.css" > api-doc/style.css
	echo '#topbar { display: none; }' >> api-doc/style.css
	echo '#short-nav, #pkg-subdirectories, .pkg-dir { display: none; }' >> api-doc/style.css
	echo 'div.param { margin-left: 2.5ex; max-width: 48rem; }' >> api-doc/style.css
	echo 'div.param .tag { font-size: 80%; opacity: 0.8; }' >> api-doc/style.css
	echo 'div.param .name { font-family: monospace; }' >> api-doc/style.css
	echo 'div.list { display: list-item; list-style: circle outside; }' >> api-doc/style.css
	echo 'div.list { margin-left: 4.5ex; max-width: 48rem; }' >> api-doc/style.css
	curl -s "http://localhost:8080/pkg/github.com/estraier/tkrzw-go/" |\
	  grep -v '^<script.*</script>$$' |\
	  sed -e 's/\/[a-z\/]*style.css/style.css/' \
	    -e 's/\/pkg\/builtin\/#/#/' \
	    -e 's/^\(@param\) \+\([a-zA-Z0-9_]\+\) \+\(.*\)/<div class="param"><span class="tag">\1<\/span> <span class="name">\2<\/span> \3<\/div>/' \
	    -e 's/^\(@return\) \+\(.*\)/<div class="param"><span class="tag">\1<\/span> \2<\/div>/' \
	    -e 's/^- \(.*\)/<div class="list">\1<\/div>/' > api-doc/index.html
	killall godoc

apidocclean :
	rm -rf api-doc
