## donger cli (WIP)

### Description

Simple CLI tool to print dongers from Donger List. Should 

### TODOs

- [x] Marshal map of donger categories to JSON
- [x] Default randomization option
- [x] Copy donger to clipboard
- [ ] CLI option to print all categories
- [x] CLI options for categories based on JSON and add command to print random donger from category or just random donger
- [ ] CLI help
- [ ] Fix an issue with sometimes returning an empty donger categories when scraping Dongerlist
- [ ] Manual and automatic periodic updates for dongers JSON, printing diff after update
- [ ] Option to turn off automatic copying to clipboard (to overcome OS limitations like Windows WSL)
- [ ] Fix Windows support (for some reason clipboard is broken)

### Limitations

- Clibboard functionality does not work in Windows WSL