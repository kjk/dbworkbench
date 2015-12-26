#!/usr/local/bin/python3
# on mac: brew install python3
#
# helper script to update all version numbers (mac, windows, frontend)

import sys, os, os.path, re

pj = os.path.join

def match_replace_group(s, m, group, replacement):
	b = m.start(1)
	e = m.end(1)
	return s[:b] + replacement + s[e:]
	
def replace_string(s, ver):
	r = re.compile('''<string>([\d\.^<]+)</string>''', re.MULTILINE | re.DOTALL | re.IGNORECASE)
	m = r.search(s)
	return match_replace_group(s, m, 1, ver)

def is_plist_ver_key(s):
	if "CFBundleShortVersionString" in s:
		return True
	return "CFBundleVersion" in s

def update_mac_ver(ver):
	path = pj("mac", "dbHero", "Info.plist")
	lines = open(path, "r").readlines()
	res = []
	ver_is_next = False
	for l in lines:
		if ver_is_next:
			l = replace_string(l, ver)
			ver_is_next = False
		else:
			ver_is_next = is_plist_ver_key(l)
		res.append(l)
	s = "".join(res)
	open(path, "w").write(s)

def update_win_ver_in_file(ver, path):
	while len(ver.split(".")) < 4:
		ver += ".0"
	s = open(path).read()

	r = re.compile('''AssemblyVersion\(\"([\d\.^"]+)"''', re.MULTILINE | re.DOTALL | re.IGNORECASE)
	m = r.search(s)
	s = match_replace_group(s, m, 1, ver)

	r = re.compile('''AssemblyFileVersion\(\"([\d\.^"]+)"''', re.MULTILINE | re.DOTALL | re.IGNORECASE)
	m = r.search(s)
	s = match_replace_group(s, m, 1, ver)

	s = s.replace("\n", "\r\n")
	open(path, "w").write(s)

def update_win_ver(ver):
    path = pj("win", "dbhero", "Properties", "AssemblyInfo.cs")
    update_win_ver_in_file(ver, path)
    path = pj("win-cef", "dbhero", "Properties", "AssemblyInfo.cs")
    update_win_ver_in_file(ver, path)

def update_frontend_ver(ver):
	path = pj("s", "index.html")
	s = open(path).read()

	r = re.compile('''gVersionNumber = \"([\d\.^"]+)";''', re.MULTILINE | re.DOTALL | re.IGNORECASE)
	m = r.search(s)
	s = match_replace_group(s, m, 1, ver)

	open(path, "w").write(s)

def usage_and_exit():
	print("usage: python ./scripts/update_ver.py ${ver}")
	sys.exit(1)

def fatal(msg):
	print(msg)
	sys.exit(1)

def is_num(s):
	try:
		n = int(s)
		return True
	except:
		return False

def verify_ver(ver):
	parts = ver.split(".")
	if len(parts) > 4:
		fatal("'%s' is not a valid version number (too manu parts)" % ver)
	for part in parts:
		if not is_num(part):
			fatal("'%s' is not a valid version number" % ver)
	
def main():
	if len(sys.argv) != 2:
		usage_and_exit()
	ver = sys.argv[1].rstrip(".0")
	verify_ver(ver)
	update_mac_ver(ver)
	update_win_ver(ver)
	update_frontend_ver(ver)
	print("updated versions to '%s'!" % ver)

if __name__ == "__main__":
	main()
