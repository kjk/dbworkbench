#!/usr/bin/env python

import os, sys, shutil,zipfile, subprocess

pj = os.path.join

script_dir = os.path.realpath(os.path.dirname(__file__))

gopath = os.environ["GOPATH"]
top_dir = os.path.join(os.path.dirname(script_dir))
src_dir = pj(top_dir, "cmd", "import_stack_overflow")

assert os.path.exists(src_dir), "%s doesn't exist" % src_dir
assert os.path.exists(pj(src_dir, "main.go")), "%s doesn't exist" % pj(src_dir, "main.go")


def abort(s):
    print(s)
    sys.exit(1)


def git_ensure_clean():
    out = subprocess.check_output(["git", "status", "--porcelain"])
    if len(out) != 0:
        print("won't deploy because repo has uncommitted changes:")
        print(out)
        sys.exit(1)


if __name__ == "__main__":
    os.chdir(top_dir)
    #git_ensure_clean()
    subprocess.check_output(["./scripts/build_import.sh"])
