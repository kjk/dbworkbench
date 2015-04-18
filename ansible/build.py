#!/usr/bin/env python

import os, sys, shutil,zipfile, subprocess

pj = os.path.join

script_dir = os.path.realpath(os.path.dirname(__file__))

gopath = os.environ["GOPATH"]
src_dir = os.path.dirname(script_dir)

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

def git_trunk_sha1():
    return subprocess.check_output(["git", "log", "-1", "--pretty=format:%H"])

if __name__ == "__main__":
    os.chdir(src_dir)
    git_ensure_clean()
    subprocess.check_output(["./scripts/build_linux.sh"])
