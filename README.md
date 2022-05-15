## donger cli

### Description

Simple CLI tool to copy dongers from Dongerlist to clipboard. Should work on all mainstream OS.

### Usage

```bash
$ donger                        will print random donger from random category
$ donger -category=CATEGORY     will print random donger from chosen category
$ donger -list                  will print all available donger categories
```

### TODOs

- [x] Marshal map of donger categories to JSON
- [x] Default randomization option
- [x] Copy donger to clipboard
- [x] CLI option to print all categories
- [x] CLI options for categories based on JSON and add command to print random donger from category or just random donger
- [ ] CLI help
- [ ] Fix an issue with sometimes returning an empty donger categories when scraping Dongerlist
- [ ] Manual and automatic periodic updates for dongers JSON, printing diff after update
- [ ] Option to turn off automatic copying to clipboard (to overcome limited OS environments)

### Limitations

- Clibboard functionality does not work out of the box in Windows WSL due to absence of X11 and maybe some other exotic environments