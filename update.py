import os
import sys

os.system('go clean -cache')
os.system('go clean -modcache')

hasTest = False
hasNoTidy = False 

if len(sys.argv) >= 2:
    hasTest = sys.argv[1] == 't'
    hasNoTidy = sys.argv[1] == 'n'

if len(sys.argv) >= 3:
    hasTest = hasTest or sys.argv[2] == 't'
    hasNoTidy = hasNoTidy or sys.argv[2] == 'n'

cmd = 'go get'
if hasTest:
    cmd = cmd + ' -t'
cmd = cmd + ' -u ./...'
os.system(cmd)

if not hasNoTidy:
    os.system('go mod tidy')