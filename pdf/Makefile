
HTML2MD?="html2md"
PDF2HTMLEX?="pdf2htmlEX"

PDF = $(wildcard *.pdf)
MD = $(PDF:.pdf=.md)

all: $(MD)

%.md: %.html
	$(HTML2MD) -i $< > $@
	sed -i '/data:image/d' $@
	sed -i '/^$$/d' $@

%.html: %.pdf
	$(PDF2HTMLEX) --zoom 0 --process-outline 0 $<

clean:
	rm -f *.md *.html *.txt
