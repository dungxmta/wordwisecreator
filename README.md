### Generate epub wordwise with ruby tag

- Based on [xnohat's repo](https://github.com/xnohat/wordwisecreator) using PHP
- It worked but quite slowww... .e.g book with 2.837.437 words took more than 2 hours

### TL;DR

- Using Go
- Pros: same file ~ 5s
- Cons:
  + only support epub
  + need prepare file manually before run app

#### Usage

- Prepare:
  + Edit .epub with `Calibre` -> fix html + beautiful all files
  + Change file type .epub to .epub.zip and extract it to folder `./data/extract/`
  + TODO: edit config in code...

- Run:
  + `go run .`
  + Replace old source (html, xhtml) with new file in `./data/extract/*`
  + zip folder extract (only file in it)
  + .zip -> .epub

#### WIP

- Note4dev:
  + Load stopwords from txt
  + Load wordwise dict from csv
  + Load source (html, xhtml extract from epub)
  + Run 10 worker (1 worker / 1 file)

- TODO:
  + extract epub to get source
  + path config
  + hint level config
  + ...
