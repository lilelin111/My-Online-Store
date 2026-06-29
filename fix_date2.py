import pathlib
t = pathlib.Path("store/store.go").read_text("utf-8")
# Fix the unescaped quotes around date in the SQL string
# Replace: "date" TIMESTAMP (in SQL context only)
old = chr(34) + "date" + chr(34) + " TIMESTAMP"
new = '\\' + chr(34) + "date" + '\\' + chr(34) + " TIMESTAMP"
t = t.replace(old, new)
pathlib.Path("store/store.go").write_text(t, "utf-8")
print("done")
