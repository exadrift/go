# v0.0.26
- fixed line counting logic on textblock
# v0.0.25
- added ability to send styled text to text block
# v0.0.24
- fix rgb color bg vs fg specifier
# v0.0.23
- allow input of rune arrays
# v0.0.22
- cache counts on TextBlock, refactor TextBlock into its own object
# v0.0.21
- count number of lines for a given render width on a TextBlock
# v0.0.20
- added ability to add style to a style list and remove conflicting styles
# v0.0.19
- filter styletypes from style array
# v0.0.18
- completely redesign style system
# v0.0.17
- style overrides need to be reapplied on any reset condition
- make sure even zero length rows are padded out appropriately
# v0.0.16
- add default fg and bg colors
# v0.0.15
- color palette support
# v0.0.14
- default minimum rows to 1 when rendering
# v0.0.13
- add rune array support to T()
# v0.0.12
- render at least one row always
# v0.0.11
- added text wrap feature which wraps text into multiple rows, handling (and removing) newlines
# v0.0.10
- detect whether or not text will require scrolling
# v0.0.9
- compose styles directly into Text
# v0.0.8
- ability to compose text from other text
# v0.0.7
- moved text into an object
# v0.0.6
- add style array concept
# v0.0.5
- change naming to default styles
# v0.0.4
- apply style overrides in render instead of generating new object
# v0.0.3
- added generic extend function which does not affect immutability
# v0.0.2
- reworked into multiple sub-packages
- added text styling system
# v0.0.1
- initial release of ANSI encoder
