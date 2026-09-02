import re, json, pathlib
gift_path = pathlib.Path("scripts/hindi_bank.gift")
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
    m_ch = re.search(r"\[chapter:(.*?)\]", qtext_raw)
    if m_ch: chapter = m_ch.group(1).strip()
    m_diff = re.search(r"\[difficulty:(.*?)\]", qtext_raw)
    if m_diff:
        difficulty = m_diff.group(1).strip().lower()
        if difficulty not in ("easy","medium","hard"): difficulty="medium"
    m_tags = re.search(r"\[tags:(.*?)\]", qtext_raw)
    if m_tags: tags = [t.strip() for t in m_tags.group(1).split(",") if t.strip()]
    clean_q = re.sub(r"\[chapter:.*?\]", "", qtext_raw)
    clean_q = re.sub(r"\[difficulty:.*?\]", "", clean_q)
    clean_q = re.sub(r"\[tags:.*?\]", "", clean_q)
    clean_q = clean_q.strip()
    # handle {=...} underline marker in question text: remove braces
    clean_q = re.sub(r"\{=([^}]*)\}", r"\1", clean_q)
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
# avoid cp1252 print issues
print("chapters", len(set(q["chapter"] for q in qs)))
print(Counter(q["difficulty"] for q in qs))
def esc(s): return s.replace("'", "''")
sid="00000000-0000-0000-0000-000000000001"
hin_id="4f167b16-cc29-415c-9266-1436a271f70e"
c10_id="aca6ed6e-2263-4841-be61-10db6a280f21"
out = pathlib.Path("scripts/hindi_391_seed.sql")
with open(out,"w",encoding="utf-8") as f:
    f.write("DO $$ DECLARE\n  sid UUID := '00000000-0000-0000-0000-000000000001';\n  hin_id UUID := '4f167b16-cc29-415c-9266-1436a271f70e';\n  c10_id UUID := 'aca6ed6e-2263-4841-be61-10db6a280f21';\nBEGIN\n")
    f.write("  INSERT INTO class_subjects (class_id, subject_id) VALUES (c10_id, hin_id) ON CONFLICT DO NOTHING;\n")
    f.write("  DELETE FROM questions WHERE subject_id=hin_id AND deleted_at IS NULL;\n")
    f.write("  INSERT INTO questions (school_id, subject_id, question_type, question_text, options, answer, marks, difficulty, chapters, tags, is_active) VALUES\n")
    rows=[]
    for q in qs:
        chapter_json=json.dumps([q["chapter"]] if q["chapter"] else [], ensure_ascii=False)
        tags_json=json.dumps(q["tags"], ensure_ascii=False)
        opts_json=json.dumps(q["options"], ensure_ascii=False)
        rows.append(f"  (sid, hin_id, 'mcq', '{esc(q['qtext'])}', '{esc(opts_json)}', '{q['ans_key']}', {q['marks']}, '{q['difficulty']}', '{esc(chapter_json)}', '{esc(tags_json)}', true)")
    f.write(",\n".join(rows)+";\nEND $$;\n")
print(f"Wrote {out} {len(rows)}")
