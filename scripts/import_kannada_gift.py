import re, json, pathlib
candidates = [pathlib.Path("scripts/kannada_bank.gift"), pathlib.Path(r"C:\Users\MDRS Bahaddurghatta\Downloads\10th_std_kannada_lba_all_mcqs-v2.gift"), pathlib.Path(r"C:\Users\MDRS Bahaddurghatta\Downloads\10th_std_kannada_lba_all_mcqs.gift")]
gift_path = None
for p in candidates:
    if p.exists():
        gift_path = p
        break
if gift_path is None:
    print("No gift file found")
    exit(1)
text = gift_path.read_text(encoding="utf-8")
pattern = re.compile(r"::(?P<title>[^:]+)::\s*(?P<qtext>.*?)\s*\{#(?P<marks>\d+)\}\s*\{\s*(?P<opts>.*?)\s*\n\}", re.DOTALL)
def parse(m):
    qtext_raw = m.group("qtext")
    title = m.group("title").strip()
    marks = m.group("marks")
    opts_raw = m.group("opts")
    chapter = None
    difficulty = "medium"
    tags=[]
    m_ch=re.search(r"\[chapter:(.*?)\]", qtext_raw)
    if m_ch: chapter=m_ch.group(1).strip()
    m_diff=re.search(r"\[difficulty:(.*?)\]", qtext_raw)
    if m_diff:
        d=m_diff.group(1).strip().lower()
        if d in ("easy","medium","hard"): difficulty=d
        elif d=="average": difficulty="medium"
        elif d=="difficult": difficulty="hard"
        else: difficulty="medium"
    m_tags=re.search(r"\[tags:(.*?)\]", qtext_raw)
    if m_tags: tags=[t.strip() for t in m_tags.group(1).split(",") if t.strip()]
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
    keys="ABCD"
    opts_json=[]
    ans_key=None
    for i,(txt, correct) in enumerate(options[:4]):
        key=keys[i]
        opts_json.append({"key":key,"value":txt,"correct":correct})
        if correct: ans_key=key
    return {"title":title,"qtext":clean_q,"marks":marks,"chapter":chapter,"difficulty":difficulty,"tags":tags,"options":opts_json,"ans_key":ans_key}

qs=[]
for m in pattern.finditer(text):
    q=parse(m)
    if q: qs.append(q)
print(f"Parsed {len(qs)}")
from collections import Counter
print(Counter(q["difficulty"] for q in qs))
print("chapters", len(set(q["chapter"] for q in qs)))
def esc(s): return s.replace("'", "''")
sid="00000000-0000-0000-0000-000000000001"
kan_id="4c526c8b-60b4-464f-af4f-85611b50aaf5"
c10_id="aca6ed6e-2263-4841-be61-10db6a280f21"
out = pathlib.Path("scripts/kannada_150_seed.sql")
with open(out,"w",encoding="utf-8") as f:
    f.write("DO $$ DECLARE\n  sid UUID := '00000000-0000-0000-0000-000000000001';\n  kan_id UUID := '4c526c8b-60b4-464f-af4f-85611b50aaf5';\n  c10_id UUID := 'aca6ed6e-2263-4841-be61-10db6a280f21';\nBEGIN\n")
    f.write("  INSERT INTO class_subjects (class_id, subject_id) VALUES (c10_id, kan_id) ON CONFLICT DO NOTHING;\n")
    f.write("  DELETE FROM quiz_responses WHERE question_id IN (SELECT id FROM questions WHERE subject_id=kan_id);\n")
    f.write("  DELETE FROM quiz_questions WHERE question_id IN (SELECT id FROM questions WHERE subject_id=kan_id);\n")
    f.write("  DELETE FROM questions WHERE subject_id=kan_id AND deleted_at IS NULL;\n")
    f.write("  INSERT INTO questions (school_id, subject_id, question_type, question_text, options, answer, marks, difficulty, chapters, tags, is_active) VALUES\n")
    rows=[]
    for q in qs:
        chapter_json=json.dumps([q["chapter"]] if q["chapter"] else [], ensure_ascii=False)
        tags_json=json.dumps(q["tags"], ensure_ascii=False)
        opts_json=json.dumps(q["options"], ensure_ascii=False)
        rows.append(f"  (sid, kan_id, 'mcq', '{esc(q['qtext'])}', '{esc(opts_json)}', '{q['ans_key']}', {q['marks']}, '{q['difficulty']}', '{esc(chapter_json)}', '{esc(tags_json)}', true)")
    f.write(",\n".join(rows)+";\nEND $$;\n")
print(f"Wrote {out} {len(rows)}")
