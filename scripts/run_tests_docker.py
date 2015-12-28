#!/usr/local/bin/python3
# on mac: brew install python3

import sys, os, os.path, time
import urllib.request, subprocess

"""
High-level overview:
- start docker container with our image (built with ./scripts/build_docker* scripts)
- run tests pointing to a database running in that container
- stop container

Repeat for all docker images. 
"""

g_imageName = "dbhero/mysql-55"
g_containerName = "mysql-55-for-tests"

kStatusRunning = "running"
kStatusExited = "exited"

def print_cmd(cmd):
  print("cmd:" + " ".join(cmd))

def run_cmd(cmd):
  print_cmd(cmd)
  subprocess.run(cmd, check=True)

def run_cmd_out(cmd):
  print_cmd(cmd)
  s = subprocess.check_output(cmd)
  return s.decode("utf-8")

def verify_docker_running():  
  try:
    run_cmd(["docker", "ps"])
  except:
    print("docker is not running! must run docker")
    sys.exit(10)

def get_docker_machine_ip():
  ip = run_cmd_out(["docker-machine", "ip", "default"])
  return ip.strip()

# returns container id and status (running, exited) for a container
# started with a given name
# returns None if no container of that name
def docker_ps(containerName):  
  s = run_cmd_out(["docker", "ps", "-a"])
  lines = s.split("\n")
  #print(lines)  
  if len(lines) < 2:
    return None
  lines = lines[1:]
  for l in lines:
    # imperfect heuristic 
    if containerName in l:
      status = kStatusRunning
      # probably imperfect heuristic
      if "Exited" in l:
        status = kStatusExited
      parts = l.split()
      return (parts[0], status)
  return None

def remove_container(containerName):
  res = docker_ps(containerName)
  if res is None:
    return
  (containerId, status) = res
  print("id: %s, status: %s" % (containerId, status))
  if status == kStatusRunning:
    run_cmd(["docker", "stop", containerId])
  run_cmd(["docker", "rm", containerId])

def start_fresh_container(imageName, containerName, portMapping):
  remove_container(containerName)
  cmd = ["docker", "run", "-d", "--name=" + containerName, "-p", portMapping, imageName]
  run_cmd(cmd)
  wait_for_container(containerName)

def wait_for_container(containerName):
  # 8 secs is a heuristic
  timeOut = 8
  print("waiting %s secs for container to start" % timeOut, end="", flush=True)
  while timeOut > 0:
    print(".", end="", flush=True)
    time.sleep(1)
    timeOut -= 1
  print("")

#DBHERO_TEST_CONN="root@tcp(192.168.99.100:7100)/world?parseTime=true" godep go test .
def run_tests(dbConnURL):
  timeoutInSecs = 25
  print("conn: '%s'" % dbConnURL)
  env = os.environ.copy()
  env["DBHERO_TEST_CONN"] = dbConnURL
  if True:
    cmd = ["godep", "go", "test", "."]
    print_cmd(cmd)
  else:
    cmd = "godep go test ."
  p = subprocess.Popen(cmd, env=env, stderr=subprocess.STDOUT,stdout=subprocess.PIPE)
  p.wait(timeoutInSecs)   
  s = p.stdout.read().decode("utf-8")
  print(s)
  # TODO: remove dbworkbench.test 
  
def mysql_conn(ip, port, dbName):
  return "root@tcp(%s:%s)/%s?parseTime=true" % (ip, port, dbName)

def pg_conn(ip, port, dbName):
  return "postgres://postgres@%s:%s/world?sslmode=disable" % (ip, port, dbName)

def main():
  verify_docker_running()
  ip = get_docker_machine_ip()
  start_fresh_container(g_imageName, g_containerName, "7100:3306")
  conn = mysql_conn(ip, "7100", "world")
  run_tests(conn)
  remove_container(g_containerName)  

if __name__ == "__main__":
  main()
