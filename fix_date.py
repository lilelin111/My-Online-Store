import os
path = r"store/store.go"
with open(path, "r", encoding="utf-8") as f:
    content = f.read()

# Fix the date issue - the pattern is:
# "note TEXT DEFAULT '', "date" TIMESTAMP NOT NULL
# should be:
# "note TEXT DEFAULT '', " + `"date"` + " TIMESTAMP NOT NULL
# This uses Go raw string for the date part, avoiding escaping issues

old_line = '"note TEXT DEFAULT '', "date" TIMESTAMP NOT NULL, total DOUBLE PRECISION NOT NULL DEFAULT 0" +'
new_line = '"note TEXT DEFAULT '', " + `"date"` + " TIMESTAMP NOT NULL, total DOUBLE PRECISION NOT NULL DEFAULT 0" +'
content = content.replace(old_line, new_line)

with open(path, "w", encoding="utf-8") as f:
    f.write(content)
print("fixed")
