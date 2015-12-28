#!/usr/local/bin/python3
# on mac: brew install python3

import sys, os, os.path
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

def run_cmd(cmd):
  print("cmd: " + " ".join(cmd))
  subprocess.run(cmd, check=True)

def run_cmd_out(cmd):
  print("cmd: " + " ".join(cmd))
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
  
def main():
  verify_docker_running()
  ip = get_docker_machine_ip()
  start_fresh_container(g_imageName, g_containerName, "7100:3306")

if __name__ == "__main__":
  main()
