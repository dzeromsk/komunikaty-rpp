DOWNLOAD?=go run download.go
DIFF?=go run diff.go
ADDINDEX?=go run addindex.go
PUBLISH?=go run publish.go
HUGO?=$(shell which hugo)

INPUT:=pdf/
OUTPUT:=content/docs/nbp/

update: download convert unfix diff addindex fix

download:
	$(DOWNLOAD) --dir $(INPUT)

convert:
	$(MAKE) -C $(INPUT) all

diff:
	$(DIFF) --input $(INPUT) --output $(OUTPUT)/

addindex:
	$(ADDINDEX) --dir $(OUTPUT)/

fix:
	mkdir -p $(OUTPUT)/200x/
	cp -r $(OUTPUT)/200[0-9] $(OUTPUT)/200x/
	rm -r $(OUTPUT)/200[0-9]
	mkdir -p $(OUTPUT)/201x/
	cp -r $(OUTPUT)/201[0-9] $(OUTPUT)/201x/
	rm -r $(OUTPUT)/201[0-9]

unfix:
	mv -f $(OUTPUT)/200x/* $(OUTPUT)/
	mv -f $(OUTPUT)/201x/* $(OUTPUT)/

build:
	$(HUGO) --minify

publish: public
	$(PUBLISH) --dir public/
