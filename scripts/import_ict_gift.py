import re, json, pathlib
levels = [
    ("scripts/ict_level1_bank.gift", "ICT L1", "ict_level1"),
    ("scripts/ict_level2_bank.gift", "ICT L2", "ict_level2"),
    ("scripts/ict_level3_bank.gift", "ICT L3", "ict_level3"),
]
pattern = re.compile(r"::(?P<title>[^:]+)::\s*(?P<header>.*?)\{\s*(?P<opts>.*?)\s*\n\}", re.DOTALL)
def parse(m, prefix, level_tag):
    qtext_raw = m.group("header")
    title = m.group("title").strip()
    marks = "2"
    mm = re.search(r"\{#(\d+)\}", qtext_raw)
    if mm:
        marks = mm.group(1)
    opts_raw = m.group("opts")
    chapter = None
    difficulty = "medium"
    tags=[]
    m_ch=re.search(r"\[chapter:(.*?)\]", qtext_raw)
    if m_ch: chapter = f"{prefix}: " + m_ch.group(1).strip()
    m_diff=re.search(r"\[difficulty:(.*?)\]", qtext_raw)
    if m_diff:
        d=m_diff.group(1).strip().lower()
        if d in ("easy","medium","hard"): difficulty=d
        elif d=="average": difficulty="medium"
        elif d=="difficult": difficulty="hard"
        else: difficulty="medium"
    m_tags=re.search(r"\[tags:(.*?)\]", qtext_raw)
    if m_tags: tags=[t.strip() for t in m_tags.group(1).split(",") if t.strip()]
    tags.append(level_tag)
    clean_q=re.sub(r"\[chapter:.*?\]", "", qtext_raw)
    clean_q=re.sub(r"\[difficulty:.*?\]", "", clean_q)
    clean_q=re.sub(r"\[tags:.*?\]", "", clean_q)
    clean_q=clean_q.strip()
    clean_q=re.sub(r"\{=([^}]*)\}", r"\1", clean_q)
    options=[]
    for line in opts_raw.splitlines():
        line=line.strip()
        if not line: continue
        if line.startswith("="):
            options.append((line[1:].strip(), True))
        elif line.startswith("~"):
            options.append((line[1:].strip(), False))
    if not any(c for _,c in options):
        return None
    keys="ABCDE"
    opts_json=[]
    ans_key=None
    for i,(txt, correct) in enumerate(options[:5]):
        key=keys[i]
        opts_json.append({"key":key,"value":txt,"correct":correct})
        if correct: ans_key=key
    return {"title":title,"qtext":clean_q,"marks":marks,"chapter":chapter,"difficulty":difficulty,"tags":tags,"options":opts_json,"ans_key":ans_key}

all_qs=[]
seen=set()
for path, prefix, level_tag in levels:
    p = pathlib.Path(path)
    if not p.exists():
        alt = pathlib.Path(r"C:\Users\MDRS Bahaddurghatta\Downloads") / pathlib.Path(path).name.replace("ict_level1_bank", "ict_level1_questions_v2").replace("ict_level2_bank", "ict_level2_questions_v2").replace("ict_level3_bank", "ict_level3_questions_v2")
        # fallback names
        cands = [p,
                 pathlib.Path(r"C:\Users\MDRS Bahaddurghatta\Downloads\ict_level1_questions_v2.gift") if "level1" in path else pathlib.Path(r"C:\Users\MDRS Bahaddurghatta\Downloads\ict_level2_questions_v2.gift") if "level2" in path else pathlib.Path(r"C:\Users\MDRS Bahaddurghatta\Downloads\ict_level3_questions_v2.gift")]
        for c in cands:
            if c.exists():
                p = c
                break
    text = p.read_text(encoding="utf-8")
    n=0
    for m in pattern.finditer(text):
        q=parse(m, prefix, level_tag)
        if q and q["qtext"] not in seen:
            seen.add(q["qtext"])
            q["level"]=prefix
            all_qs.append(q)
            n+=1
    print(f"{prefix}: {n} unique from {p}")

from collections import Counter
print(f"Total {len(all_qs)}")
print(Counter(q["difficulty"] for q in all_qs))
print(Counter(q["level"] for q in all_qs))
def esc(s): return s.replace("'", "''")
out = pathlib.Path("scripts/ict_levels_seed.sql")
with open(out,"w",encoding="utf-8") as f:
    f.write("DO $$ DECLARE\n  sid UUID := '00000000-0000-0000-0000-000000000001';\n  cs_id UUID := 'a8004612-55ef-4304-bbc8-74295be72476';\nBEGIN\n")
    f.write("  INSERT INTO class_subjects (class_id, subject_id) SELECT c.id, cs_id FROM classes c WHERE c.deleted_at IS NULL AND c.name IN ('Class 6','Class 7','Class 8','Class 9','Class 10') ON CONFLICT DO NOTHING;\n")
    f.write("  DELETE FROM quiz_responses WHERE question_id IN (SELECT id FROM questions WHERE subject_id=cs_id AND chapters::text LIKE '%ICT L%');\n")
    f.write("  DELETE FROM quiz_questions WHERE question_id IN (SELECT id FROM questions WHERE subject_id=cs_id AND chapters::text LIKE '%ICT L%');\n")
    f.write("  DELETE FROM questions WHERE subject_id=cs_id AND chapters::text LIKE '%ICT L%' AND deleted_at IS NULL;\n")
    f.write("  INSERT INTO questions (school_id, subject_id, question_type, question_text, options, answer, marks, difficulty, chapters, tags, is_active) VALUES\n")
    rows=[]
    for q in all_qs:
        chapter_json=json.dumps([q["chapter"]] if q["chapter"] else [], ensure_ascii=False)
        tags_json=json.dumps(q["tags"], ensure_ascii=False)
        opts_json=json.dumps(q["options"], ensure_ascii=False)
        rows.append(f"  (sid, cs_id, 'mcq', '{esc(q['qtext'])}', '{esc(opts_json)}', '{q['ans_key']}', {q['marks']}, '{q['difficulty']}', '{esc(chapter_json)}', '{esc(tags_json)}', true)")
    f.write(",\n".join(rows)+";\nEND $$;\n")
print(f"Wrote {out} {len(rows)}")
