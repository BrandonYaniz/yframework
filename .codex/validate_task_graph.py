#!/usr/bin/env python3
from pathlib import Path
import re, sys
root=Path(__file__).resolve().parent
files=list((root/'TASKS').glob('[0-9][0-9][0-9]_*.md'))
ids=sorted(int(p.name[:3]) for p in files)
want=list(range(1,882))
if ids != want:
    missing=sorted(set(want)-set(ids)); dup=[]
    print('task graph invalid', 'missing=', missing, file=sys.stderr); sys.exit(1)
required=['README.md','INSTRUCTIONS.md','DECISION_AUTHORITY.md','ARCHITECTURE_INDEX.md','CURRENT_MILESTONE.md','ACCEPTANCE.md']
for x in required:
    if not (root/x).exists(): print('missing',x,file=sys.stderr); sys.exit(1)
print('ok: 881 task files and required .codex control files present')
