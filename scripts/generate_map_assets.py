#!/usr/bin/env python3
"""Generate SVG map assets + map_places seed for the Map Quiz.

- frontend/static/maps/india.json      (country outline, projected)
- frontend/static/maps/karnataka.json  (state outline, projected)
- scripts/map_places_seed.sql          (places with lat/lng + svg_x/svg_y)

Projection: linear fit of lng/lat bbox into viewBox (aspect preserved).
Pins use the SAME projection, so they always align with outlines.
"""
import json, math, pathlib

ROOT = pathlib.Path(__file__).parent.parent
MAP_DIR = ROOT / "frontend" / "static" / "maps"
MAP_DIR.mkdir(parents=True, exist_ok=True)

# ——— India country outline (lng,lat), from world.geo.json IND ———
INDIA_OUTLINE = [
[77.837451,35.49401],[78.912269,34.321936],[78.811086,33.506198],[79.208892,32.994395],[79.176129,32.48378],[78.458446,32.618164],[78.738894,31.515906],[79.721367,30.882715],[81.111256,30.183481],[80.476721,29.729865],[80.088425,28.79447],[81.057203,28.416095],[81.999987,27.925479],[83.304249,27.364506],[84.675018,27.234901],[85.251779,26.726198],[86.024393,26.630985],[87.227472,26.397898],[88.060238,26.414615],[88.174804,26.810405],[88.043133,27.445819],[88.120441,27.876542],[88.730326,28.086865],[88.814248,27.299316],[88.835643,27.098966],[89.744528,26.719403],[90.373275,26.875724],[91.217513,26.808648],[92.033484,26.83831],[92.103712,27.452614],[91.696657,27.771742],[92.503119,27.896876],[93.413348,28.640629],[94.56599,29.277438],[95.404802,29.031717],[96.117679,29.452802],[96.586591,28.83098],[96.248833,28.411031],[97.327114,28.261583],[97.402561,27.882536],[97.051989,27.699059],[97.133999,27.083774],[96.419366,27.264589],[95.124768,26.573572],[95.155153,26.001307],[94.603249,25.162495],[94.552658,24.675238],[94.106742,23.850741],[93.325188,24.078556],[93.286327,23.043658],[93.060294,22.703111],[93.166128,22.27846],[92.672721,22.041239],[92.146035,23.627499],[91.869928,23.624346],[91.706475,22.985264],[91.158963,23.503527],[91.46773,24.072639],[91.915093,24.130414],[92.376202,24.976693],[91.799596,25.147432],[90.872211,25.132601],[89.920693,25.26975],[89.832481,25.965082],[89.355094,26.014407],[88.563049,26.446526],[88.209789,25.768066],[88.931554,25.238692],[88.306373,24.866079],[88.084422,24.501657],[88.69994,24.233715],[88.52977,23.631142],[88.876312,22.879146],[89.031961,22.055708],[88.888766,21.690588],[88.208497,21.703172],[86.975704,21.495562],[87.033169,20.743308],[86.499351,20.151638],[85.060266,19.478579],[83.941006,18.30201],[83.189217,17.671221],[82.192792,17.016636],[82.191242,16.556664],[81.692719,16.310219],[80.791999,15.951972],[80.324896,15.899185],[80.025069,15.136415],[80.233274,13.835771],[80.286294,13.006261],[79.862547,12.056215],[79.857999,10.357275],[79.340512,10.308854],[78.885345,9.546136],[79.18972,9.216544],[78.277941,8.933047],[77.941165,8.252959],[77.539898,7.965535],[76.592979,8.899276],[76.130061,10.29963],[75.746467,11.308251],[75.396101,11.781245],[74.864816,12.741936],[74.616717,13.992583],[74.443859,14.617222],[73.534199,15.990652],[73.119909,17.92857],[72.820909,19.208234],[72.824475,20.419503],[72.630533,21.356009],[71.175273,20.757441],[70.470459,20.877331],[69.16413,22.089298],[69.644928,22.450775],[69.349597,22.84318],[68.176645,23.691965],[68.842599,24.359134],[71.04324,24.356524],[70.844699,25.215102],[70.282873,25.722229],[70.168927,26.491872],[69.514393,26.940966],[70.616496,27.989196],[71.777666,27.91318],[72.823752,28.961592],[73.450638,29.976413],[74.42138,30.979815],[74.405929,31.692639],[75.258642,32.271105],[74.451559,32.7649],[74.104294,33.441473],[73.749948,34.317699],[74.240203,34.748887],[75.757061,34.504923],[76.871722,34.653544],
]

# ——— Karnataka schematic outline (lng,lat), clockwise from Karwar coast ———
KARNATAKA_OUTLINE = [
[74.13,14.81],[74.28,14.66],[74.40,14.42],[74.53,14.00],[74.61,13.63],[74.74,13.34],[74.86,12.87],[74.95,12.55],[75.10,12.30],[75.35,12.10],[75.65,11.95],[75.95,11.75],[76.30,11.80],[76.65,11.85],[77.00,11.95],[77.25,12.10],[77.50,12.30],[77.75,12.60],[78.00,12.85],[78.30,13.00],[78.55,13.15],[78.60,13.40],[78.45,13.70],[78.30,14.00],[78.10,14.30],[77.90,14.60],[77.70,15.00],[77.55,15.40],[77.50,15.80],[77.55,16.20],[77.60,16.60],[77.55,17.00],[77.52,17.40],[77.50,17.80],[77.30,17.95],[77.00,17.90],[76.70,17.70],[76.40,17.40],[76.10,17.10],[75.80,16.85],[75.50,16.60],[75.20,16.40],[74.95,16.20],[74.70,16.00],[74.50,15.85],[74.30,15.60],[74.15,15.30],[74.05,15.00],
]

def make_projection(outline, W, H, pad=24):
    lngs = [p[0] for p in outline]
    lats = [p[1] for p in outline]
    min_lng, max_lng = min(lngs), max(lngs)
    min_lat, max_lat = min(lats), max(lats)
    mean_lat = (min_lat + max_lat) / 2
    kx = math.cos(math.radians(mean_lat))
    w_deg = (max_lng - min_lng) * kx
    h_deg = max_lat - min_lat
    scale = min((W - 2 * pad) / w_deg, (H - 2 * pad) / h_deg)
    ox = pad + ((W - 2 * pad) - w_deg * scale) / 2
    oy = pad + ((H - 2 * pad) - h_deg * scale) / 2
    def proj(lng, lat):
        x = ox + (lng - min_lng) * kx * scale
        y = oy + (max_lat - lat) * scale
        return (round(x, 1), round(y, 1))
    return proj

INDIA_W, INDIA_H = 760, 860
KA_W, KA_H = 560, 760

india_proj = make_projection(INDIA_OUTLINE, INDIA_W, INDIA_H)
ka_proj = make_projection(KARNATAKA_OUTLINE, KA_W, KA_H)

def path_d(outline, proj):
    pts = [proj(lng, lat) for lng, lat in outline]
    d = f"M {pts[0][0]} {pts[0][1]} " + " ".join(f"L {x} {y}" for x, y in pts[1:]) + " Z"
    return d

(MAP_DIR / "india.json").write_text(json.dumps({
    "viewBox": [0, 0, INDIA_W, INDIA_H],
    "outlines": [{"id": "india", "d": path_d(INDIA_OUTLINE, india_proj)}],
    "schematic": False,
}, ensure_ascii=False), encoding="utf-8")

(MAP_DIR / "karnataka.json").write_text(json.dumps({
    "viewBox": [0, 0, KA_W, KA_H],
    "outlines": [{"id": "karnataka", "d": path_d(KARNATAKA_OUTLINE, ka_proj)}],
    "schematic": True,
    "note": "Schematic outline for practice, not to survey scale.",
}, ensure_ascii=False), encoding="utf-8")

print("wrote map JSON assets")

# ——— Places: (name, kind, category, map, lat, lng, state, district, why_in_news, news_date, exam_tags) ———
# Categories double as chapters: "Maps:<Category>"
P = []
def add(name, kind, cat, mp, lat, lng, state="", district="", news="", date=None, tags=None):
    P.append((name, kind, cat, mp, lat, lng, state, district, news, date, tags or []))

U = ["upsc"]; K = ["kpsc"]; UK = ["upsc", "kpsc"]; S = ["ssc", "banking"]; ALL = ["upsc", "kpsc", "ssc"]

# Karnataka Special (karnataka map)
add("Bengaluru", "city", "Karnataka Special", "karnataka", 12.97, 77.59, "Karnataka", "Bengaluru Urban", "State capital; IT capital", None, UK)
add("Mysuru", "city", "Karnataka Special", "karnataka", 12.30, 76.65, "Karnataka", "Mysuru", "Dasara; Palace", None, K)
add("Belagavi", "city", "Karnataka Special", "karnataka", 15.85, 74.50, "Karnataka", "Belagavi", "Winter session; Suvarna Soudha", None, K)
add("Kalaburagi", "city", "Karnataka Special", "karnataka", 17.32, 76.83, "Karnataka", "Kalaburagi", "Sannati stupa region", None, K)
add("Mangaluru", "port", "Karnataka Special", "karnataka", 12.87, 74.84, "Karnataka", "Dakshina Kannada", "Major port", None, UK)
add("Karwar", "port", "Karnataka Special", "karnataka", 14.81, 74.13, "Karnataka", "Uttara Kannada", "INS Kadamba naval base", None, UK)
add("Hampi", "heritage", "Karnataka Special", "karnataka", 15.33, 76.46, "Karnataka", "Vijayanagara", "UNESCO site", None, UK)
add("Pattadakal", "heritage", "Karnataka Special", "karnataka", 15.95, 75.82, "Karnataka", "Bagalkote", "UNESCO site", None, UK)
add("Belur", "heritage", "Karnataka Special", "karnataka", 13.16, 75.86, "Karnataka", "Hassan", "Hoysala UNESCO 2023", "2023-09-18", UK)
add("Halebidu", "heritage", "Karnataka Special", "karnataka", 13.21, 75.99, "Karnataka", "Hassan", "Hoysala UNESCO 2023", "2023-09-18", UK)
add("Jog Falls", "waterfall", "Karnataka Special", "karnataka", 14.22, 74.81, "Karnataka", "Shivamogga", "Sharavathi; 253 m", None, K)
add("Krishna Raja Sagara", "dam", "Karnataka Special", "karnataka", 12.42, 76.57, "Karnataka", "Mandya", "Kaveri reservoir", None, K)
add("Almatti Dam", "dam", "Karnataka Special", "karnataka", 16.33, 75.88, "Karnataka", "Vijayapura", "Krishna; dispute in news", "2024-06-01", K)
add("Kolar Gold Fields", "industrial", "Karnataka Special", "karnataka", 12.96, 78.27, "Karnataka", "Kolar", "Historic gold fields", None, K)
add("Shravanabelagola", "heritage", "Karnataka Special", "karnataka", 12.85, 76.48, "Karnataka", "Hassan", "Gomateshwara", None, K)
add("Ballari", "city", "Karnataka Special", "karnataka", 15.14, 76.92, "Karnataka", "Ballari", "Steel; Vijayanagara district split", None, K)
add("Raichur", "city", "Karnataka Special", "karnataka", 16.20, 77.35, "Karnataka", "Raichur", "Thermal power; Maski nearby", None, K)
add("Bidar", "city", "Karnataka Special", "karnataka", 15.91, 77.51, "Karnataka", "Bidar", "Bidriware; fort", None, K)
add("Shivamogga", "city", "Karnataka Special", "karnataka", 14.16, 75.63, "Karnataka", "Shivamogga", "Malnad gateway; Agumbe", None, K)
add("Tumakuru", "city", "Karnataka Special", "karnataka", 13.33, 77.10, "Karnataka", "Tumakuru", "Pavagada solar nearby", None, K)

# National Parks & Reserves (india map)
add("Bandipur NP", "park", "National Parks", "india", 11.66, 76.63, "Karnataka", "Chamarajanagara", "Tiger reserve", None, UK)
add("Nagarahole NP", "park", "National Parks", "india", 12.05, 76.15, "Karnataka", "Mysuru", "Tiger reserve", None, UK)
add("BRT Tiger Reserve", "park", "National Parks", "india", 11.95, 77.15, "Karnataka", "Chamarajanagara", "Tiger reserve", None, UK)
add("Bhadra WLS", "park", "National Parks", "india", 13.45, 75.65, "Karnataka", "Chikkamagaluru", "Tiger reserve", None, K)
add("Kali Tiger Reserve", "park", "National Parks", "india", 15.05, 74.45, "Karnataka", "Uttara Kannada", "Dandeli-Anshi", None, K)
add("Jim Corbett NP", "park", "National Parks", "india", 29.53, 78.77, "Uttarakhand", "", "Oldest NP (1936)", None, U)
add("Kaziranga NP", "park", "National Parks", "india", 26.57, 93.17, "Assam", "", "One-horned rhino; UNESCO", None, U)
add("Gir NP", "park", "National Parks", "india", 21.13, 70.82, "Gujarat", "", "Asiatic lion", None, U)
add("Sundarbans NP", "park", "National Parks", "india", 21.94, 88.90, "West Bengal", "", "Mangrove; tiger; Ramsar", None, U)
add("Keoladeo Ghana NP", "park", "National Parks", "india", 27.15, 77.50, "Rajasthan", "", "Ramsar 1981; Montreux", None, U)
add("Kanha NP", "park", "National Parks", "india", 22.33, 80.63, "Madhya Pradesh", "", "Tiger; barasingha", None, U)
add("Ranthambore NP", "park", "National Parks", "india", 26.02, 76.38, "Rajasthan", "", "Tiger reserve", None, U)
add("Periyar NP", "park", "National Parks", "india", 9.46, 77.24, "Kerala", "", "Tiger; elephant", None, U)
add("Manas NP", "park", "National Parks", "india", 26.66, 91.00, "Assam", "", "UNESCO; BTC", None, U)
add("Nanda Devi NP", "park", "National Parks", "india", 30.37, 79.97, "Uttarakhand", "", "UNESCO; biosphere", None, U)
add("Silent Valley NP", "park", "National Parks", "india", 11.13, 76.43, "Kerala", "", "Evergreen; lion-tailed macaque", None, U)
add("Great Himalayan NP", "park", "National Parks", "india", 31.75, 77.55, "Himachal Pradesh", "", "UNESCO 2014", None, U)
add("Desert NP", "park", "National Parks", "india", 25.70, 70.55, "Rajasthan", "", "Great Indian Bustard", None, U)
add("Simlipal NP", "park", "National Parks", "india", 21.93, 86.38, "Odisha", "", "Tiger; Similipal fires in news", "2024-03-01", U)
add("Madhav NP (TR)", "park", "National Parks", "india", 25.50, 77.75, "Madhya Pradesh", "", "58th tiger reserve (2025)", "2025-03-09", U)

# Nuclear & Power (india)
add("Kaiga Nuclear Plant", "power", "Nuclear & Power", "india", 14.86, 74.44, "Karnataka", "Uttara Kannada", "800 MW unit in news", "2024-01-01", UK)
add("Kudankulam Nuclear Plant", "power", "Nuclear & Power", "india", 8.17, 77.71, "Tamil Nadu", "", "Largest nuclear plant", None, U)
add("Tarapur Nuclear Plant", "power", "Nuclear & Power", "india", 19.83, 72.66, "Maharashtra", "", "First nuclear plant (1969)", None, U)
add("Kalpakkam (MAPS)", "power", "Nuclear & Power", "india", 12.55, 80.17, "Tamil Nadu", "", "PFBR; fast breeder", "2024-03-04", U)
add("Rawatbhata Nuclear Plant", "power", "Nuclear & Power", "india", 24.87, 76.60, "Rajasthan", "", "RAPS", None, U)
add("Pavagada Solar Park", "power", "Nuclear & Power", "india", 14.28, 77.27, "Karnataka", "Tumakuru", "2 GW solar park", None, K)
add("Khavda Solar Park", "power", "Nuclear & Power", "india", 23.55, 68.85, "Gujarat", "", "World's largest RE park", "2024-02-01", U)
add("Tehri Dam", "power", "Nuclear & Power", "india", 30.37, 78.47, "Uttarakhand", "", "Tallest dam (260 m)", None, S)
add("Bhakra Nangal Dam", "power", "Nuclear & Power", "india", 31.41, 76.43, "Punjab", "", "Sutlej; multipurpose", None, S)
add("Sardar Sarovar Dam", "power", "Nuclear & Power", "india", 21.83, 73.74, "Gujarat", "", "Narmada", None, S)
add("Hirakud Dam", "power", "Nuclear & Power", "india", 21.57, 83.87, "Odisha", "", "Mahanadi; longest dam", None, S)
add("Nagarjuna Sagar Dam", "power", "Nuclear & Power", "india", 16.57, 79.31, "Telangana", "", "Krishna", None, S)

# IITs & Institutions (india) — 23 IITs + IISc
add("IIT Bombay", "institute", "IITs & Institutions", "india", 19.13, 72.91, "Maharashtra", "Mumbai", "197 institutions: 1958", None, U)
add("IIT Delhi", "institute", "IITs & Institutions", "india", 28.54, 77.19, "Delhi", "", "1961", None, U)
add("IIT Madras", "institute", "IITs & Institutions", "india", 12.99, 80.23, "Tamil Nadu", "Chennai", "1959; top NIRF", None, U)
add("IIT Kanpur", "institute", "IITs & Institutions", "india", 26.51, 80.23, "Uttar Pradesh", "", "1959", None, U)
add("IIT Kharagpur", "institute", "IITs & Institutions", "india", 22.31, 87.30, "West Bengal", "", "First IIT (1951)", None, U)
add("IIT Roorkee", "institute", "IITs & Institutions", "india", 29.86, 77.89, "Uttarakhand", "", "1847 Thomason; IIT 2001", None, U)
add("IIT Guwahati", "institute", "IITs & Institutions", "india", 26.19, 91.69, "Assam", "", "1994", None, U)
add("IIT Hyderabad", "institute", "IITs & Institutions", "india", 17.49, 78.21, "Telangana", "", "2008", None, U)
add("IIT Dharwad", "institute", "IITs & Institutions", "india", 15.43, 74.98, "Karnataka", "Dharwad", "2016; Karnataka's IIT", None, K)
add("IISc Bengaluru", "institute", "IITs & Institutions", "india", 13.02, 77.56, "Karnataka", "Bengaluru Urban", "1909; top research", None, UK)
add("IIT Indore", "institute", "IITs & Institutions", "india", 22.52, 75.92, "Madhya Pradesh", "", "2009", None, S)
add("IIT Patna", "institute", "IITs & Institutions", "india", 25.53, 84.85, "Bihar", "", "2008", None, S)
add("IIT Bhubaneswar", "institute", "IITs & Institutions", "india", 20.24, 85.67, "Odisha", "", "2008", None, S)
add("IIT Gandhinagar", "institute", "IITs & Institutions", "india", 23.22, 72.63, "Gujarat", "", "2008", None, S)
add("IIT Ropar", "institute", "IITs & Institutions", "india", 30.95, 76.47, "Punjab", "", "2008", None, S)
add("IIT Mandi", "institute", "IITs & Institutions", "india", 31.70, 76.98, "Himachal Pradesh", "", "2009", None, S)
add("IIT Jodhpur", "institute", "IITs & Institutions", "india", 26.47, 73.11, "Rajasthan", "", "2008", None, S)
add("IIT Varanasi (BHU)", "institute", "IITs & Institutions", "india", 25.26, 82.98, "Uttar Pradesh", "", "1919 BENCO; IIT 2012", None, S)
add("IIT Tirupati", "institute", "IITs & Institutions", "india", 13.72, 79.33, "Andhra Pradesh", "", "2015", None, S)
add("IIT Palakkad", "institute", "IITs & Institutions", "india", 10.80, 76.65, "Kerala", "", "2015", None, S)
add("IIT Dhanbad (ISM)", "institute", "IITs & Institutions", "india", 23.81, 86.44, "Jharkhand", "", "1926 ISM; IIT 2016", None, S)
add("IIT Bhilai", "institute", "IITs & Institutions", "india", 21.19, 81.65, "Chhattisgarh", "", "2016", None, S)
add("IIT Jammu", "institute", "IITs & Institutions", "india", 32.63, 74.89, "J&K", "", "2016", None, S)
add("IIT Goa", "institute", "IITs & Institutions", "india", 15.42, 73.99, "Goa", "", "2016", None, S)

# Heritage UNESCO (india)
add("Taj Mahal", "heritage", "Heritage Sites", "india", 27.17, 78.01, "Uttar Pradesh", "Agra", "UNESCO 1983", None, U)
add("Qutub Minar", "heritage", "Heritage Sites", "india", 28.52, 77.18, "Delhi", "", "UNESCO 1993", None, U)
add("Red Fort", "heritage", "Heritage Sites", "india", 28.65, 77.24, "Delhi", "", "UNESCO 2007", None, U)
add("Khajuraho", "heritage", "Heritage Sites", "india", 24.85, 79.92, "Madhya Pradesh", "", "UNESCO 1986", None, U)
add("Konark Sun Temple", "heritage", "Heritage Sites", "india", 19.89, 86.09, "Odisha", "", "UNESCO 1984", None, U)
add("Ajanta Caves", "heritage", "Heritage Sites", "india", 20.55, 75.70, "Maharashtra", "", "UNESCO 1983", None, U)
add("Ellora Caves", "heritage", "Heritage Sites", "india", 20.02, 75.17, "Maharashtra", "", "UNESCO 1983", None, U)
add("Sanchi Stupa", "heritage", "Heritage Sites", "india", 23.47, 77.73, "Madhya Pradesh", "", "UNESCO 1989", None, U)
add("Mahabalipuram", "heritage", "Heritage Sites", "india", 12.61, 80.19, "Tamil Nadu", "", "UNESCO 1984", None, U)
add("Thanjavur (Brihadisvara)", "heritage", "Heritage Sites", "india", 10.78, 79.13, "Tamil Nadu", "", "Great Living Chola; UNESCO", None, U)
add("Raigad Fort", "heritage", "Heritage Sites", "india", 18.23, 73.44, "Maharashtra", "", "Maratha Landscapes UNESCO 2025", "2025-07-11", U)
add("Shivneri Fort", "heritage", "Heritage Sites", "india", 19.19, 73.85, "Maharashtra", "", "Maratha Landscapes UNESCO 2025", "2025-07-11", U)
add("Sindhudurg Fort", "heritage", "Heritage Sites", "india", 16.04, 73.46, "Maharashtra", "", "Maratha Landscapes UNESCO 2025", "2025-07-11", U)
add("Gingee Fort", "heritage", "Heritage Sites", "india", 12.25, 79.41, "Tamil Nadu", "", "Only non-Maharashtra Maratha site", "2025-07-11", U)
add("Moidams of Charaideo", "heritage", "Heritage Sites", "india", 27.31, 94.87, "Assam", "", "UNESCO 2024 (43rd)", "2024-07-26", U)
add("Jantar Mantar Jaipur", "heritage", "Heritage Sites", "india", 26.92, 75.82, "Rajasthan", "", "UNESCO 2010", None, S)
add("Fatehpur Sikri", "heritage", "Heritage Sites", "india", 27.09, 77.66, "Uttar Pradesh", "", "UNESCO 1986", None, S)
add("Hampi (Vijaya Vittala)", "heritage", "Heritage Sites", "india", 15.33, 76.46, "Karnataka", "Vijayanagara", "UNESCO 1986", None, UK)
add("Sarnath", "heritage", "Heritage Sites", "india", 25.37, 83.02, "Uttar Pradesh", "", "Tentative; Buddha first sermon", None, U)
add("Dholavira", "heritage", "Heritage Sites", "india", 23.88, 70.21, "Gujarat", "", "UNESCO 2021; Harappan", None, U)

# Ashokan Edicts (india)
add("Maski Edict", "edict", "Ashokan Edicts", "india", 15.96, 76.66, "Karnataka", "Raichur", "Only edict naming Ashoka", None, UK)
add("Brahmagiri Edict", "edict", "Ashokan Edicts", "india", 14.80, 76.80, "Karnataka", "Chitradurga", "Minor rock edicts", None, K)
add("Sannati Stupa & Edict", "edict", "Ashokan Edicts", "india", 16.84, 76.95, "Karnataka", "Kalaburagi", "Buddhist site", None, K)
add("Girnar Edict", "edict", "Ashokan Edicts", "india", 21.52, 70.45, "Gujarat", "", "Major rock edict", None, U)
add("Dhauli Edict", "edict", "Ashokan Edicts", "india", 20.19, 85.84, "Odisha", "", "Kalinga; elephant carving", None, U)
add("Jaugada Edict", "edict", "Ashokan Edicts", "india", 19.55, 84.95, "Odisha", "", "Major rock edict", None, U)
add("Kalsi Edict", "edict", "Ashokan Edicts", "india", 30.52, 77.85, "Uttarakhand", "", "Only complete set of 14", None, U)
add("Shahbazgarhi Edict", "edict", "Ashokan Edicts", "india", 34.00, 72.60, "Pakistan (historic India)", "", "Kharosthi script", None, U)
add("Sarnath Pillar", "edict", "Ashokan Edicts", "india", 25.37, 83.02, "Uttar Pradesh", "", "Lion capital = emblem", None, U)
add("Lauriya Nandangarh", "edict", "Ashokan Edicts", "india", 26.55, 84.40, "Bihar", "", "Pillar edicts", None, U)

# Neolithic & Megalithic (india)
add("Sanganakallu", "neolithic", "Neolithic Sites", "india", 15.20, 76.88, "Karnataka", "Ballari", "Neolithic; dolerite dykes", None, K)
add("Tekkalakota", "neolithic", "Neolithic Sites", "india", 15.32, 76.88, "Karnataka", "Ballari", "Gold + Neolithic", None, K)
add("Hallur", "neolithic", "Neolithic Sites", "india", 15.02, 75.63, "Karnataka", "Haveri", "Earliest horse gram", None, K)
add("Hire Benakal", "neolithic", "Neolithic Sites", "india", 15.45, 76.47, "Karnataka", "Raichur", "Megalithic dolmens", None, K)
add("Burzahom", "neolithic", "Neolithic Sites", "india", 34.16, 74.90, "J&K", "", "Pit dwellings", None, U)
add("Chirand", "neolithic", "Neolithic Sites", "india", 25.73, 85.33, "Bihar", "", "Neolithic-Chalcolithic", None, U)
add("Mehrgarh", "neolithic", "Neolithic Sites", "india", 29.38, 67.61, "Baluchistan (historic)", "", "9000-yr farming (context)", None, U)
add("Utnur", "neolithic", "Neolithic Sites", "india", 16.53, 77.33, "Karnataka", "Raichur", "Cattle pastoralism", None, K)
add("Budhihal", "neolithic", "Neolithic Sites", "india", 15.05, 76.55, "Karnataka", "Ballari", "Ashmound site", None, K)
add("Kupgal", "neolithic", "Neolithic Sites", "india", 15.15, 76.90, "Karnataka", "Ballari", "Rock art + Neolithic", None, K)

# Ports & Capitals (india)
add("Mumbai Port", "port", "Ports & Capitals", "india", 18.93, 72.93, "Maharashtra", "", "Gateway of India; largest", None, S)
add("Mundra Port", "port", "Ports & Capitals", "india", 22.73, 69.70, "Gujarat", "", "Largest private port (Adani)", None, U)
add("Vizhinjam Port", "port", "Ports & Capitals", "india", 8.38, 77.00, "Kerala", "", "Transshipment; 2025 ops", "2025-05-02", U)
add("Kolkata Port", "port", "Ports & Capitals", "india", 22.68, 88.06, "West Bengal", "", "Tea Port; riverine", None, S)
add("Visakhapatnam Port", "port", "Ports & Capitals", "india", 17.68, 83.28, "Andhra Pradesh", "", "Natural harbour; landlocked", None, S)
add("Chennai Port", "port", "Ports & Capitals", "india", 13.10, 80.29, "Tamil Nadu", "", "Artificial harbour", None, S)
add("Paradip Port", "port", "Ports & Capitals", "india", 20.26, 86.68, "Odisha", "", "Iron ore", None, S)
add("New Delhi", "capital", "Ports & Capitals", "india", 28.61, 77.20, "Delhi", "", "National capital", None, S)
add("Jaipur", "capital", "Ports & Capitals", "india", 26.91, 75.78, "Rajasthan", "", "Pink City", None, S)
add("Lucknow", "capital", "Ports & Capitals", "india", 26.84, 80.94, "Uttar Pradesh", "", "Awadh capital", None, S)
add("Bhopal", "capital", "Ports & Capitals", "india", 23.25, 77.41, "Madhya Pradesh", "", "Lake city", None, S)
add("Gandhinagar", "capital", "Ports & Capitals", "india", 23.22, 72.63, "Gujarat", "", "Planned capital", None, S)

# Ramsar Wetlands (india) — latest 2025-26 first
add("Glaw Lake", "ramsar", "Ramsar Wetlands", "india", 28.30, 94.50, "Arunachal Pradesh", "", "101st Ramsar (Aug 2026)", "2026-08-03", U)
add("JPN Bird Sanctuary (Surha Tal)", "ramsar", "Ramsar Wetlands", "india", 25.80, 84.20, "Uttar Pradesh", "Ballia", "100th Ramsar (2026)", "2026-01-19", U)
add("Patna Bird Sanctuary", "ramsar", "Ramsar Wetlands", "india", 27.60, 78.70, "Uttar Pradesh", "Etah", "Ramsar 2026", "2026-01-01", U)
add("Chhari-Dhand", "ramsar", "Ramsar Wetlands", "india", 23.60, 69.20, "Gujarat", "Kutch", "Ramsar 2026", "2026-01-01", U)
add("Khichan Wetland", "ramsar", "Ramsar Wetlands", "india", 27.13, 72.42, "Rajasthan", "Phalodi", "Ramsar 2025; demoiselle cranes", "2025-06-04", U)
add("Menar Wetland", "ramsar", "Ramsar Wetlands", "india", 24.66, 73.93, "Rajasthan", "Udaipur", "Ramsar 2025; bird village", "2025-06-04", U)
add("Gokul Jalashay", "ramsar", "Ramsar Wetlands", "india", 25.90, 85.90, "Bihar", "", "Ramsar 2025", "2025-05-13", U)
add("Kopra Reservoir", "ramsar", "Ramsar Wetlands", "india", 22.00, 82.10, "Chhattisgarh", "Bilaspur", "Ramsar 2025", "2025-08-08", U)
add("Siliserh Lake", "ramsar", "Ramsar Wetlands", "india", 27.50, 76.60, "Rajasthan", "Alwar", "Ramsar 2025; Sariska buffer", "2025-08-08", U)
add("Chilika Lake", "ramsar", "Ramsar Wetlands", "india", 19.70, 85.30, "Odisha", "", "First Ramsar 1981; lagoon", None, U)
add("Keoladeo Ghana Wetland", "ramsar", "Ramsar Wetlands", "india", 27.15, 77.50, "Rajasthan", "", "First Ramsar 1981; Montreux", None, U)
add("Loktak Lake", "ramsar", "Ramsar Wetlands", "india", 24.50, 93.77, "Manipur", "", "Phumdis; Montreux", None, U)
add("Vembanad-Kol", "ramsar", "Ramsar Wetlands", "india", 9.60, 76.40, "Kerala", "", "Largest Kerala wetland", None, U)
add("Kolleru Lake", "ramsar", "Ramsar Wetlands", "india", 16.60, 81.20, "Andhra Pradesh", "", "Godavari-Krishna basin", None, U)
add("Deepor Beel", "ramsar", "Ramsar Wetlands", "india", 26.11, 91.65, "Assam", "", "Brahmaputra floodplain", None, U)
add("Ranganathittu BS", "ramsar", "Ramsar Wetlands", "india", 12.42, 76.65, "Karnataka", "Mandya", "Karnataka's Ramsar 2022", "2022-02-15", K)
add("East Kolkata Wetlands", "ramsar", "Ramsar Wetlands", "india", 22.55, 88.45, "West Bengal", "", "Wastewater aquaculture", None, U)
add("Sambhar Lake", "ramsar", "Ramsar Wetlands", "india", 26.90, 75.08, "Rajasthan", "", "Salt lake; flamingos", None, U)
add("Wular Lake", "ramsar", "Ramsar Wetlands", "india", 34.30, 74.55, "J&K", "", "Largest freshwater lake", None, U)
add("Tso Moriri", "ramsar", "Ramsar Wetlands", "india", 32.90, 78.30, "Ladakh", "", "High-altitude lake", None, U)

# Dams & Rivers (india)
add("KRS Dam", "dam", "Dams & Rivers", "india", 12.42, 76.57, "Karnataka", "Mandya", "Kaveri", None, K)
add("Tungabhadra Dam", "dam", "Dams & Rivers", "india", 15.30, 76.35, "Karnataka", "Vijayanagara", "Pampa Sagar", None, K)
add("Linganamakki Dam", "dam", "Dams & Rivers", "india", 14.15, 74.85, "Karnataka", "Shivamogga", "Sharavathi", None, K)
add("Supa Dam", "dam", "Dams & Rivers", "india", 15.28, 74.52, "Karnataka", "Uttara Kannada", "Kali river", None, K)
add("Bhakra Dam", "dam", "Dams & Rivers", "india", 31.41, 76.43, "Himachal Pradesh", "", "Sutlej; highest gravity", None, S)
add("Tehri Dam", "dam", "Dams & Rivers", "india", 30.37, 78.47, "Uttarakhand", "", "Tallest 260 m", None, S)
add("Sardar Sarovar Dam", "dam", "Dams & Rivers", "india", 21.83, 73.74, "Gujarat", "", "Narmada", None, S)
add("Hirakud Dam", "dam", "Dams & Rivers", "india", 21.57, 83.87, "Odisha", "", "Mahanadi; longest", None, S)
add("Nagarjuna Sagar Dam", "dam", "Dams & Rivers", "india", 16.57, 79.31, "Telangana", "", "Krishna", None, S)
add("Kosi Barrage", "dam", "Dams & Rivers", "india", 26.52, 86.93, "Bihar", "", "Sorrow of Bihar", None, S)

# Freedom & Battles (india)
add("Plassey", "battle", "Freedom & Battles", "india", 23.78, 88.25, "West Bengal", "", "1757; Clive vs Siraj", None, U)
add("Buxar", "battle", "Freedom & Battles", "india", 25.57, 83.98, "Bihar", "", "1764", None, U)
add("Talikota", "battle", "Freedom & Battles", "india", 16.48, 76.32, "Karnataka", "Vijayapura", "1565; Vijayanagara fall", None, UK)
add("Srirangapatna", "battle", "Freedom & Battles", "india", 12.42, 76.68, "Karnataka", "Mandya", "Tipu 1799", None, UK)
add("Jallianwala Bagh", "battle", "Freedom & Battles", "india", 31.62, 74.87, "Punjab", "Amritsar", "1919 massacre", None, U)
add("Dandi", "battle", "Freedom & Battles", "india", 20.89, 72.80, "Gujarat", "", "Salt March 1930", None, U)
add("Sabarmati Ashram", "battle", "Freedom & Battles", "india", 23.02, 72.57, "Gujarat", "", "Gandhi 1917-30", None, U)
add("Cellular Jail", "battle", "Freedom & Battles", "india", 11.67, 92.74, "A&N Islands", "", "Kala Pani", None, U)
add("Meerut 1857", "battle", "Freedom & Battles", "india", 28.98, 77.70, "Uttar Pradesh", "", "Revolt began May 1857", None, U)
add("Gwalior Fort 1858", "battle", "Freedom & Battles", "india", 26.23, 78.16, "Madhya Pradesh", "", "Rani Lakshmibai fell", None, U)

# Space & Science (india)
add("Sriharikota (SHAR)", "space", "Space & Science", "india", 13.71, 80.23, "Andhra Pradesh", "", "Satish Dhawan Centre", None, U)
add("Thumba (VSSC)", "space", "Space & Science", "india", 8.52, 76.86, "Kerala", "", "First rocket 1963", None, U)
add("Byalalu DSN", "space", "Space & Science", "india", 12.90, 77.36, "Karnataka", "Bengaluru Rural", "Deep-space antennas", None, K)
add("Pokhran", "space", "Space & Science", "india", 27.09, 71.75, "Rajasthan", "", "Nuclear tests 1974/1998", None, U)
add("Chandipur (ITR)", "space", "Space & Science", "india", 21.44, 87.01, "Odisha", "", "Missile test range", None, U)
add("Kodaikanal Observatory", "space", "Space & Science", "india", 10.23, 77.46, "Tamil Nadu", "", "Solar observatory", None, S)
add("Hanle Observatory", "space", "Space & Science", "india", 32.77, 78.96, "Ladakh", "", "Highest observatory; dark-sky reserve", "2022-01-01", U)
add("Jantar Mantar Delhi", "space", "Space & Science", "india", 28.62, 77.21, "Delhi", "", "1724 observatory", None, S)

# Tiger Reserves top-up (india) — a few beyond parks list
add("Sariska TR", "park", "National Parks", "india", 27.30, 76.43, "Rajasthan", "", "Tiger reserve", None, S)
add("Dudhwa TR", "park", "National Parks", "india", 28.47, 80.68, "Uttar Pradesh", "", "Tiger; rhino reintro", None, S)

print(f"places: {len(P)}")
for cat in sorted(set(p[2] for p in P)):
    print(" ", cat, sum(1 for p in P if p[2] == cat))

# ——— emit SQL ———
def esc(s): return (s or "").replace("'", "''")

def proj_for(mp):
    return india_proj if mp == "india" else ka_proj

rows = []
for (name, kind, cat, mp, lat, lng, state, district, news, date, tags) in P:
    x, y = proj_for(mp)(lng, lat)
    tags_json = json.dumps(tags, ensure_ascii=False)
    date_sql = f"'{date}'" if date else "NULL"
    news_sql = f"'{esc(news)}'" if news else "NULL"
    rows.append(
        "  ('00000000-0000-0000-0000-000000000001', "
        f"'{esc(name)}', '{esc(kind)}', 'Maps:{esc(cat)}', '{mp}', "
        f"{lat}, {lng}, {x}, {y}, '{esc(state)}', '{esc(district)}', "
        f"{news_sql}, {date_sql}, '{esc(tags_json)}', true)"
    )

cols = "(school_id, name, kind, category, map, lat, lng, svg_x, svg_y, state, district, why_in_news, news_date, exam_tags, is_active)"
sql = (
    "DO $$ BEGIN\n"
    "  DELETE FROM map_places WHERE category LIKE 'Maps:%';\n"
    "  INSERT INTO map_places " + cols + " VALUES\n"
    + ",\n".join(rows) + ";\nEND $$;\n"
)
out = pathlib.Path(__file__).parent / "map_places_seed.sql"
out.write_text(sql, encoding="utf-8")
print(f"wrote {out} rows={len(rows)}")
