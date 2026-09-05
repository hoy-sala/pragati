import json, sys
p = sys.argv[1] if len(sys.argv) > 1 else "/tmp/indcomp.geojson"
fc = json.load(open(p, encoding="utf-8"))
print("features:", len(fc["features"]))
for f in fc["features"]:
    props = f.get("properties", {})
    g = f["geometry"]
    s = {"n": 0, "minlat": 90, "maxlat": -90, "minlng": 180, "maxlng": -180}
    def walk(c):
        if isinstance(c[0], (int, float)):
            lng, lat = c[0], c[1]
            s["n"] += 1
            s["minlat"] = min(s["minlat"], lat); s["maxlat"] = max(s["maxlat"], lat)
            s["minlng"] = min(s["minlng"], lng); s["maxlng"] = max(s["maxlng"], lng)
        else:
            for x in c:
                walk(x)
    walk(g["coordinates"])
    print(props.get("ST_NM") or props.get("NAME_1") or props.get("name") or props.get("State"), "|", g["type"], "| pts:", s["n"], "| bbox:", round(s["minlng"],2), round(s["minlat"],2), round(s["maxlng"],2), round(s["maxlat"],2))
