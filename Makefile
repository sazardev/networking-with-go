# Print and web builds for the book.
#
# `pretty-pdf` alone cannot render the book's Mermaid diagrams (v0.8.0 prints
# them as literal text and blocks both network and file:// subresources), so
# `pdf` and `epub` first materialize build/docs with the checked-in SVGs from
# assets/diagrams inlined. `make diagrams` is the authoring step that produces
# those SVGs; it needs Node.js and is never used by CI.
PRETTY_PDF ?= pretty-pdf

.PHONY: pdf epub diagrams clean

pdf:
	go run tools/diagrams/main.go
	$(PRETTY_PDF) build --source build/docs --out networking-with-go.pdf

epub:
	go run tools/diagrams/main.go
	$(PRETTY_PDF) epub --source build/docs --out networking-with-go.epub

diagrams:
	go run tools/diagrams/main.go -render

clean:
	rm -rf build networking-with-go.pdf networking-with-go.epub
