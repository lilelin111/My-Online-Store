content = open("store/store.go","r",encoding="utf-8").read()
content = content.replace('\", \"date\" \"','\", \\\"date\\\" \"')
open("store/store.go","w",encoding="utf-8").write(content)
print("fixed")
