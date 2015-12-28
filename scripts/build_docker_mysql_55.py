#!/usr/local/bin/python3
# on mac: brew install python3

import sys, os, os.path
import urllib.request, zipfile, subprocess

g_world_db_url = "http://downloads.mysql.com/docs/world.sql.zip"

g_tmp_data_dir = "tmp_data"

g_world_db_zip_path = os.path.join(g_tmp_data_dir, "world-mysql.sql.zip")
g_world_db_mysql_zip_sha1 = "c66bce8b37b253f0437755c078290e4c9c1e2104"
g_world_db_mysql_path = os.path.join(g_tmp_data_dir, "world-mysql.sql")

def fatal(s):
  print(s)
  sys.exit(1)

def sha1OfFile(filepath):
    import hashlib
    with open(filepath, 'rb') as f:
        return hashlib.sha1(f.read()).hexdigest()

def verifyFileSha1(path, expectedSha1):
  sha1 = sha1OfFile(path)
  if sha1 != expectedSha1:
    fatal("unexpected sha1 of '%s'. Wanted: '%s', is: '%s'" % (path, expectedSha1, sha1))

def dl_world_mysql():
  if os.path.exists(g_world_db_zip_path):
    print("dl_world_mysql: '%s' already exists\n" % g_world_db_zip_path)
    verifyFileSha1(g_world_db_zip_path, g_world_db_mysql_zip_sha1)
    return
  print("downloading %s\n" % g_world_db_url)
  urllib.request.urlretrieve(g_world_db_url, g_world_db_zip_path)
  verifyFileSha1(g_world_db_zip_path, g_world_db_mysql_zip_sha1)

def dl_and_unzip_world_mysql():
  if os.path.exists(g_world_db_mysql_path):
    print("dl_and_unzip_world_mysql: '%s' already exists\n" % g_world_db_mysql_path)
    return
  dl_world_mysql()
  print("dl_and_unzip_world_mysql: extracting world.sql")
  zf = zipfile.ZipFile(g_world_db_zip_path, "r")
  d = zf.read("world.sql")
  open(g_world_db_mysql_path, "wb").write(d)

def create_docker_image():
  print("building docker image")
  #docker build -t dbhero/mysql-55 -f scripts/mysql-55.dockerfile .
  subprocess.run(["docker", "build", "-t", "dbhero/mysql-55", "-f", "scripts/mysql-55.dockerfile", "."])

def verify_docker_running():
  try:
    subprocess.check_output("docker ps", shell=True)
  except:
    print("docker is not running! must run docker")
    sys.exit(10)

def main():
  verify_docker_running()
  os.makedirs(g_tmp_data_dir, exist_ok=True)
  dl_and_unzip_world_mysql()
  create_docker_image()
  
if __name__ == "__main__":
  main()
