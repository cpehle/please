"""A Python script to rewrite hashes in BUILD files."""

import ast


# These are templated in by Go. It's a bit hacky but is a way of avoiding
# passing arbitrary arguments through Go / C calls.
FILENAME = '__FILENAME__'
TARGETS = {__TARGETS__}
PLATFORM = '__PLATFORM__'


def is_a_target(node):
    """Returns the name of a node if it's a target that we're interested in."""
    if isinstance(node, ast.Expr) and isinstance(node.value, ast.Call):
        for keyword in node.value.keywords:
            if keyword.arg == 'name':
                if isinstance(keyword.value, ast.Str) and keyword.value.s in TARGETS:
                    return keyword.value.s


with _open(FILENAME) as f:
    lines = f.readlines()
    tree = ast.parse(''.join(lines), filename=FILENAME)

for node in ast.iter_child_nodes(tree):
    name = is_a_target(node)
    if name:
        for keyword in node.value.keywords:
            if keyword.arg == 'hashes' and isinstance(keyword.value, ast.List):
                candidates = {dep.s: dep.lineno - 1 for dep in keyword.value.elts
                              if isinstance(dep, ast.Str)}
                # Filter by any leading platform (i.e. linux_amd64: abcdef12345).
                platform_candidates = {k: v for k, v in candidates.items() if PLATFORM in k}
                if len(platform_candidates) == 1:
                    candidates = platform_candidates
                # Should really do something here about multiple hashes and working out which
                # is which...
                current, lineno = candidates.popitem()
                lines[lineno] = lines[lineno].replace(current, TARGETS[name])

with _open(FILENAME, 'w') as f:
    for line in lines:
        f.write(line)
