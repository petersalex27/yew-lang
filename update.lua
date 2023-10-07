require "os"

os.execute('go clean -cache')
os.execute('go clean -modcache')

local hasTest, doTidy = false, true

if #arg >= 2 then
  hasTest, doTidy = arg[1] == 't', arg[1] ~= 'n'
end

if #arg >= 3 then
  hasTest = hasTest or arg[2] == 't'
  doTidy = doTidy and arg[2] ~= 'n'
end

local cmd = 'go get'
if hasTest then
  cmd = cmd .. '-t'
end
cmd = cmd .. ' -u ./...'
os.execute(cmd)

if doTidy then
    os.execute('go mod tidy')
end